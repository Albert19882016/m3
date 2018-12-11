// Copyright (c) 2016 Uber Technologies, Inc.
//
// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in
// all copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
// THE SOFTWARE.

package storage

import (
	"fmt"
	"sync"
	"time"

	"github.com/m3db/m3/src/dbnode/clock"
	"github.com/m3db/m3/src/dbnode/persist"
	"github.com/m3db/m3/src/dbnode/persist/fs"
	"github.com/m3db/m3/src/dbnode/persist/fs/commitlog"
	"github.com/m3db/m3/src/dbnode/retention"
	xerrors "github.com/m3db/m3x/errors"
	"github.com/m3db/m3x/ident"
	xlog "github.com/m3db/m3x/log"
	"github.com/pborman/uuid"

	"github.com/uber-go/tally"
)

type commitLogFilesFn func(commitlog.Options) ([]persist.CommitlogFile, []commitlog.ErrorWithPath, error)
type sortedSnapshotMetadataFilesFn func(fs.Options) ([]fs.SnapshotMetadata, []fs.SnapshotMetadataErrorWithPaths, error)

type snapshotFilesFn func(filePathPrefix string, namespace ident.ID, shard uint32) (fs.FileSetFilesSlice, error)

type deleteFilesFn func(files []string) error

type deleteInactiveDirectoriesFn func(parentDirPath string, activeDirNames []string) error

// Narrow interface so as not to expose all the functionality of the commitlog
// to the cleanup manager.
type activeCommitlogs interface {
	ActiveLogs() ([]persist.CommitlogFile, error)
}

type cleanupManager struct {
	sync.RWMutex

	database         database
	activeCommitlogs activeCommitlogs

	opts                          Options
	nowFn                         clock.NowFn
	filePathPrefix                string
	commitLogsDir                 string
	commitLogFilesFn              commitLogFilesFn
	sortedSnapshotMetadataFilesFn sortedSnapshotMetadataFilesFn
	snapshotFilesFn               snapshotFilesFn

	deleteFilesFn               deleteFilesFn
	deleteInactiveDirectoriesFn deleteInactiveDirectoriesFn
	cleanupInProgress           bool
	metrics                     cleanupManagerMetrics
}

type cleanupManagerMetrics struct {
	status                      tally.Gauge
	corruptCommitlogFile        tally.Counter
	corruptSnapshotMetadataFile tally.Counter
	deletedCommitlogFile        tally.Counter
	deletedSnapshotMetadataFile tally.Counter
}

func newCleanupManagerMetrics(scope tally.Scope) cleanupManagerMetrics {
	clScope := scope.SubScope("commitlog")
	smScope := scope.SubScope("snapshot-metadata")
	return cleanupManagerMetrics{
		status:                      scope.Gauge("cleanup"),
		corruptCommitlogFile:        clScope.Counter("corrupt"),
		corruptSnapshotMetadataFile: smScope.Counter("corrupt"),
		deletedCommitlogFile:        clScope.Counter("deleted"),
		deletedSnapshotMetadataFile: smScope.Counter("deleted"),
	}
}

func newCleanupManager(
	database database, activeLogs activeCommitlogs, scope tally.Scope) databaseCleanupManager {
	opts := database.Options()
	filePathPrefix := opts.CommitLogOptions().FilesystemOptions().FilePathPrefix()
	commitLogsDir := fs.CommitLogsDirPath(filePathPrefix)

	return &cleanupManager{
		database:         database,
		activeCommitlogs: activeLogs,

		opts:                          opts,
		nowFn:                         opts.ClockOptions().NowFn(),
		filePathPrefix:                filePathPrefix,
		commitLogsDir:                 commitLogsDir,
		commitLogFilesFn:              commitlog.Files,
		sortedSnapshotMetadataFilesFn: fs.SortedSnapshotMetadataFiles,
		snapshotFilesFn:               fs.SnapshotFiles,
		deleteFilesFn:                 fs.DeleteFiles,
		deleteInactiveDirectoriesFn:   fs.DeleteInactiveDirectories,
		metrics:                       newCleanupManagerMetrics(scope),
	}
}

func (m *cleanupManager) Cleanup(t time.Time) error {
	m.Lock()
	m.cleanupInProgress = true
	m.Unlock()

	defer func() {
		m.Lock()
		m.cleanupInProgress = false
		m.Unlock()
	}()

	multiErr := xerrors.NewMultiError()
	if err := m.cleanupExpiredDataFiles(t); err != nil {
		multiErr = multiErr.Add(fmt.Errorf(
			"encountered errors when cleaning up data files for %v: %v", t, err))
	}

	if err := m.cleanupExpiredIndexFiles(t); err != nil {
		multiErr = multiErr.Add(fmt.Errorf(
			"encountered errors when cleaning up index files for %v: %v", t, err))
	}

	if err := m.deleteInactiveDataFiles(); err != nil {
		multiErr = multiErr.Add(fmt.Errorf(
			"encountered errors when deleting inactive data files for %v: %v", t, err))
	}

	if err := m.deleteInactiveDataSnapshotFiles(); err != nil {
		multiErr = multiErr.Add(fmt.Errorf(
			"encountered errors when deleting inactive snapshot files for %v: %v", t, err))
	}

	if err := m.deleteInactiveNamespaceFiles(); err != nil {
		multiErr = multiErr.Add(fmt.Errorf(
			"encountered errors when deleting inactive namespace files for %v: %v", t, err))
	}

	if err := m.cleanupSnapshotsAndCommitlogs(); err != nil {
		multiErr = multiErr.Add(fmt.Errorf(
			"encountered errors when cleaning up snapshot and commitlog files: %v", err))
	}

	return multiErr.FinalError()
}

func (m *cleanupManager) Report() {
	m.RLock()
	cleanupInProgress := m.cleanupInProgress
	m.RUnlock()

	if cleanupInProgress {
		m.metrics.status.Update(1)
	} else {
		m.metrics.status.Update(0)
	}
}

func (m *cleanupManager) deleteInactiveNamespaceFiles() error {
	var namespaceDirNames []string
	filePathPrefix := m.database.Options().CommitLogOptions().FilesystemOptions().FilePathPrefix()
	dataDirPath := fs.DataDirPath(filePathPrefix)
	namespaces, err := m.database.GetOwnedNamespaces()
	if err != nil {
		return err
	}

	for _, n := range namespaces {
		namespaceDirNames = append(namespaceDirNames, n.ID().String())
	}

	return m.deleteInactiveDirectoriesFn(dataDirPath, namespaceDirNames)
}

// deleteInactiveDataFiles will delete data files for shards that the node no longer owns
// which can occur in the case of topology changes
func (m *cleanupManager) deleteInactiveDataFiles() error {
	return m.deleteInactiveDataFileSetFiles(fs.NamespaceDataDirPath)
}

// deleteInactiveDataSnapshotFiles will delete snapshot files for shards that the node no longer owns
// which can occur in the case of topology changes
func (m *cleanupManager) deleteInactiveDataSnapshotFiles() error {
	return m.deleteInactiveDataFileSetFiles(fs.NamespaceSnapshotsDirPath)
}

func (m *cleanupManager) deleteInactiveDataFileSetFiles(filesetFilesDirPathFn func(string, ident.ID) string) error {
	multiErr := xerrors.NewMultiError()
	filePathPrefix := m.database.Options().CommitLogOptions().FilesystemOptions().FilePathPrefix()
	namespaces, err := m.database.GetOwnedNamespaces()
	if err != nil {
		return err
	}
	for _, n := range namespaces {
		var activeShards []string
		namespaceDirPath := filesetFilesDirPathFn(filePathPrefix, n.ID())
		for _, s := range n.GetOwnedShards() {
			shard := fmt.Sprintf("%d", s.ID())
			activeShards = append(activeShards, shard)
		}
		multiErr = multiErr.Add(m.deleteInactiveDirectoriesFn(namespaceDirPath, activeShards))
	}

	return multiErr.FinalError()
}

func (m *cleanupManager) cleanupExpiredDataFiles(t time.Time) error {
	multiErr := xerrors.NewMultiError()
	namespaces, err := m.database.GetOwnedNamespaces()
	if err != nil {
		return err
	}
	for _, n := range namespaces {
		if !n.Options().CleanupEnabled() {
			continue
		}
		earliestToRetain := retention.FlushTimeStart(n.Options().RetentionOptions(), t)
		shards := n.GetOwnedShards()
		multiErr = multiErr.Add(m.cleanupExpiredNamespaceDataFiles(earliestToRetain, shards))
	}
	return multiErr.FinalError()
}

func (m *cleanupManager) cleanupExpiredIndexFiles(t time.Time) error {
	namespaces, err := m.database.GetOwnedNamespaces()
	if err != nil {
		return err
	}
	multiErr := xerrors.NewMultiError()
	for _, n := range namespaces {
		if !n.Options().CleanupEnabled() || !n.Options().IndexOptions().Enabled() {
			continue
		}
		idx, err := n.GetIndex()
		if err != nil {
			multiErr = multiErr.Add(err)
			continue
		}
		multiErr = multiErr.Add(idx.CleanupExpiredFileSets(t))
	}
	return multiErr.FinalError()
}

func (m *cleanupManager) cleanupExpiredNamespaceDataFiles(earliestToRetain time.Time, shards []databaseShard) error {
	multiErr := xerrors.NewMultiError()
	for _, shard := range shards {
		if err := shard.CleanupExpiredFileSets(earliestToRetain); err != nil {
			multiErr = multiErr.Add(err)
		}
	}

	return multiErr.FinalError()
}

// List all the metadata files on disk
// Identify the most recent one
// Delete all snapshot files whose snapshot ID does not match the most recent id from the most recent metadata file
// Delete all the metadata files before the most recent one
// Delete all commitlogs before the one whose commilog we identified
func (m *cleanupManager) cleanupSnapshotsAndCommitlogs() error {
	namespaces, err := m.database.GetOwnedNamespaces()
	if err != nil {
		return err
	}

	fsOpts := m.opts.CommitLogOptions().FilesystemOptions()
	sortedSnapshotMetadatas, snapshotMetadataErrorsWithPaths, err := fs.SortedSnapshotMetadataFiles(fsOpts)
	if err != nil {
		return err
	}

	if len(sortedSnapshotMetadatas) == 0 {
		// No cleanup can be performed until we have at least one complete snapshot.
		return nil
	}

	var (
		filesToDelete      = []string{}
		mostRecentSnapshot = sortedSnapshotMetadatas[len(sortedSnapshotMetadatas)-1]
	)

	for _, ns := range namespaces {
		for _, s := range ns.GetOwnedShards() {
			// TODO: Need to make sure we can handle corrupt files.
			shardSnapshots, err := m.snapshotFilesFn(fsOpts.FilePathPrefix(), ns.ID(), s.ID())
			if err != nil {
				// TODO: Multierr?
				return err
			}

			for _, snapshot := range shardSnapshots {
				_, snapshotID, err := snapshot.SnapshotTimeAndID()
				if err != nil {
					// TODO: Multierr?
					return err
				}

				if !uuid.Equal(snapshotID, mostRecentSnapshot.ID.UUID) {
					// If the UUID of the snapshot files doesn't match the most recent snapshot
					// then its safe to delete because it means we have a more recently complete set.
					filesToDelete = append(filesToDelete, snapshot.AbsoluteFilepaths...)
				}
			}
		}
	}

	// TODO: Handle corrupt files
	// TODO: Handle active logs files (actually might not need this)
	files, commitlogErrorsWithPaths, err := m.commitLogFilesFn(m.opts.CommitLogOptions())
	if err != nil {
		return err
	}

	// Delete all commitlog files prior to the one captured by the most recent snapshot.
	for _, file := range files {
		if file.Index < mostRecentSnapshot.CommitlogIdentifier.Index {
			m.metrics.deletedCommitlogFile.Inc(1)
			filesToDelete = append(filesToDelete, file.FilePath)
		}
	}

	// Delete all snapshot metadatas prior to the most recent one.
	for _, snapshot := range sortedSnapshotMetadatas[:len(sortedSnapshotMetadatas)-1] {
		m.metrics.deletedSnapshotMetadataFile.Inc(1)
		filesToDelete = append(filesToDelete, snapshot.AbsoluteFilepaths()...)
	}

	// Delete corrupt snapshot metadata files.
	for _, errorWithPath := range snapshotMetadataErrorsWithPaths {
		m.metrics.corruptSnapshotMetadataFile.Inc(1)
		// TODO: Comment.
		m.opts.InstrumentOptions().Logger().WithFields(
			xlog.NewField("err", errorWithPath.Error),
			xlog.NewField("metadataFilePath", errorWithPath.MetadataFilePath),
			xlog.NewField("checkpointFilePath", errorWithPath.CheckpointFilePath),
		).Errorf(
			"encountered corrupt snapshot metadata file during cleanup, marking files for deletion")
		filesToDelete = append(filesToDelete, errorWithPath.MetadataFilePath)
		filesToDelete = append(filesToDelete, errorWithPath.CheckpointFilePath)
	}

	// Delete corrupt commitlog files.
	for _, errorWithPath := range commitlogErrorsWithPaths {
		m.metrics.corruptCommitlogFile.Inc(1)
		// If we were unable to read the commit log files info header, then we're forced to assume
		// that the file is corrupt and remove it. This can happen in situations where M3DB experiences
		// sudden shutdown.
		m.opts.InstrumentOptions().Logger().WithFields(
			xlog.NewField("err", errorWithPath.Error()),
			xlog.NewField("path", errorWithPath.Path()),
		).Errorf(
			"encountered corrupt commitlog file during cleanup, marking file for deletion: %s",
			errorWithPath.Error())
		filesToDelete = append(filesToDelete, errorWithPath.Path())
	}

	return nil
}

// commitLogTimes returns the earliest time before which the commit logs are expired,
// as well as a list of times we need to clean up commit log files for.
func (m *cleanupManager) commitLogTimes(t time.Time) ([]commitLogFileWithErrorAndPath, error) {
	namespaces, err := m.database.GetOwnedNamespaces()
	if err != nil {
		return nil, err
	}

	// We list the commit log files on disk before we determine what the currently active commitlog
	// is to ensure that the logic remains correct even if the commitlog is rotated while this
	// function is executing. For example, imagine the following commitlogs are on disk:
	//
	// [time1, time2, time3]
	//
	// If we call ActiveLogs first then it will return time3. Next, the commit log file rotates, and
	// after that we call commitLogFilesFn which returns: [time1, time2, time3, time4]. In this scenario
	// we would be allowed to delete commitlog files 1,2, and 4 which is not the desired behavior. Instead,
	// we list the commitlogs on disk first (which returns time1, time2, and time3) and *then* check what
	// the active file is. If the commitlog has not rotated, then ActiveLogs() will return time3 which
	// we will correctly avoid deleting, and if the commitlog has rotated, then ActiveLogs() will return
	// time4 which we wouldn't consider deleting anyways because it wasn't returned from the first call
	// to commitLogFilesFn.
	files, corruptFiles, err := m.commitLogFilesFn(m.opts.CommitLogOptions())
	if err != nil {
		return nil, err
	}

	activeCommitlogs, err := m.activeCommitlogs.ActiveLogs()
	if err != nil {
		return nil, err
	}

	shouldCleanupFile := func(f persist.CommitlogFile) (bool, error) {
		if commitlogsContainPath(activeCommitlogs, f.FilePath) {
			// An active commitlog should never satisfy all of the constraints
			// for deleting a commitlog, but skip them for posterity.
			return false, nil
		}

		for _, ns := range namespaces {
			var (
				start                      = f.Start
				duration                   = f.Duration
				ropts                      = ns.Options().RetentionOptions()
				nsBlocksStart, nsBlocksEnd = commitLogNamespaceBlockTimes(start, duration, ropts)
				needsFlush                 = ns.NeedsFlush(nsBlocksStart, nsBlocksEnd)
			)

			outOfRetention := nsBlocksEnd.Before(retention.FlushTimeStart(ropts, t))
			if outOfRetention {
				continue
			}

			if !needsFlush {
				// Data has been flushed to disk so the commit log file is
				// safe to clean up.
				continue
			}

			// Add commit log blockSize to the startTime because that is the latest
			// system time that the commit log file could contain data for. Note that
			// this is different than the latest datapoint timestamp that the commit
			// log file could contain data for (because of bufferPast/bufferFuture),
			// but the commit log files and snapshot files both deal with system time.
			isCapturedBySnapshot, err := ns.IsCapturedBySnapshot(
				nsBlocksStart, nsBlocksEnd, start.Add(duration))
			if err != nil {
				// Return error because we don't want to proceed since this is not a commitlog
				// file specific issue.
				return false, err
			}

			if !isCapturedBySnapshot {
				// The data has not been flushed and has also not been captured by
				// a snapshot, so it is not safe to clean up the commit log file.
				return false, nil
			}

			// All the data in the commit log file is captured by the snapshot files
			// so its safe to clean up.
		}

		return true, nil
	}

	filesToCleanup := make([]commitLogFileWithErrorAndPath, 0, len(files))
	for _, f := range files {
		shouldDelete, err := shouldCleanupFile(f)
		if err != nil {
			return nil, err
		}

		if shouldDelete {
			filesToCleanup = append(filesToCleanup, newCommitLogFileWithErrorAndPath(
				f, f.FilePath, nil))
		}
	}

	for _, errorWithPath := range corruptFiles {
		if commitlogsContainPath(activeCommitlogs, errorWithPath.Path()) {
			// Skip active commit log files as they may appear corrupt due to the
			// header info not being written out yet.
			continue
		}

		m.metrics.corruptCommitlogFile.Inc(1)
		// If we were unable to read the commit log files info header, then we're forced to assume
		// that the file is corrupt and remove it. This can happen in situations where M3DB experiences
		// sudden shutdown.
		m.opts.InstrumentOptions().Logger().Errorf(
			"encountered corrupt commitlog file during cleanup, marking file for deletion: %s",
			errorWithPath.Error())
		filesToCleanup = append(filesToCleanup, newCommitLogFileWithErrorAndPath(
			persist.CommitlogFile{}, errorWithPath.Path(), err))
	}

	return filesToCleanup, nil
}

// commitLogNamespaceBlockTimes returns the range of namespace block starts for which the
// given commit log block may contain data for.
//
// consider the situation where we have a single namespace, and a commit log with the following
// retention options:
//             buffer past | buffer future | block size
// namespace:    ns_bp     |     ns_bf     |    ns_bs
// commit log:     _       |      _        |    cl_bs
//
// for the commit log block with start time `t`, we can receive data for a range of namespace
// blocks depending on the namespace retention options. The range is given by these relationships:
//	- earliest ns block start = t.Add(-ns_bp).Truncate(ns_bs)
//  - latest ns block start   = t.Add(cl_bs).Add(ns_bf).Truncate(ns_bs)
// NB:
// - blockStart assumed to be aligned to commit log block size
func commitLogNamespaceBlockTimes(
	blockStart time.Time,
	commitlogBlockSize time.Duration,
	nsRetention retention.Options,
) (time.Time, time.Time) {
	earliest := blockStart.
		Add(-nsRetention.BufferPast()).
		Truncate(nsRetention.BlockSize())
	latest := blockStart.
		Add(commitlogBlockSize).
		Add(nsRetention.BufferFuture()).
		Truncate(nsRetention.BlockSize())
	return earliest, latest
}

type commitLogFileWithErrorAndPath struct {
	f    persist.CommitlogFile
	path string
	err  error
}

func newCommitLogFileWithErrorAndPath(
	f persist.CommitlogFile, path string, err error) commitLogFileWithErrorAndPath {
	return commitLogFileWithErrorAndPath{
		f:    f,
		path: path,
		err:  err,
	}
}

func commitlogsContainPath(commitlogs []persist.CommitlogFile, path string) bool {
	for _, f := range commitlogs {
		if path == f.FilePath {
			return true
		}
	}

	return false
}
