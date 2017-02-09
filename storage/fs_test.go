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
	"errors"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
	"github.com/uber-go/tally"
)

func TestFileSystemManagerShouldRunDuringBootstrap(t *testing.T) {
	database := newMockDatabase()
	mgr, err := newFileSystemManager(database, tally.NoopScope)
	require.NoError(t, err)
	now := database.Options().ClockOptions().NowFn()()
	require.False(t, mgr.ShouldRun(now))
	database.bs = bootstrapped
	require.True(t, mgr.ShouldRun(now))
}

func TestFileSystemManagerShouldRunWhileRunning(t *testing.T) {
	database := newMockDatabase()
	database.bs = bootstrapped
	fsm, err := newFileSystemManager(database, tally.NoopScope)
	require.NoError(t, err)
	mgr := fsm.(*fileSystemManager)
	now := database.Options().ClockOptions().NowFn()()
	require.True(t, mgr.ShouldRun(now))
	mgr.status = fileOpInProgress
	require.False(t, mgr.ShouldRun(now))
}

func TestFileSystemManagerRun(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	database := newMockDatabase()
	database.bs = bootstrapped
	fm := NewMockdatabaseFlushManager(ctrl)
	cm := NewMockdatabaseCleanupManager(ctrl)
	fsm, err := newFileSystemManager(database, tally.NoopScope)
	require.NoError(t, err)
	mgr := fsm.(*fileSystemManager)
	mgr.databaseFlushManager = fm
	mgr.databaseCleanupManager = cm

	ts := time.Now()
	gomock.InOrder(
		cm.EXPECT().Cleanup(ts).Return(errors.New("foo")),
		fm.EXPECT().Flush(ts).Return(errors.New("bar")),
	)

	mgr.Run(ts, false)
	require.Equal(t, fileOpNotStarted, mgr.status)
}