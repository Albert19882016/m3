// Code generated by MockGen. DO NOT EDIT.
// Source: github.com/m3db/m3/src/dbnode/persist/fs/commitlog/types.go

// Copyright (c) 2018 Uber Technologies, Inc.
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

// Package commitlog is a generated GoMock package.
package commitlog

import (
	"reflect"
	"time"

	"github.com/m3db/m3/src/dbnode/clock"
	"github.com/m3db/m3/src/dbnode/persist/fs"
	"github.com/m3db/m3/src/dbnode/ts"
	"github.com/m3db/m3x/context"
	"github.com/m3db/m3x/ident"
	"github.com/m3db/m3x/instrument"
	"github.com/m3db/m3x/pool"
	time0 "github.com/m3db/m3x/time"

	"github.com/golang/mock/gomock"
)

// MockCommitLog is a mock of CommitLog interface
type MockCommitLog struct {
	ctrl     *gomock.Controller
	recorder *MockCommitLogMockRecorder
}

// MockCommitLogMockRecorder is the mock recorder for MockCommitLog
type MockCommitLogMockRecorder struct {
	mock *MockCommitLog
}

// NewMockCommitLog creates a new mock instance
func NewMockCommitLog(ctrl *gomock.Controller) *MockCommitLog {
	mock := &MockCommitLog{ctrl: ctrl}
	mock.recorder = &MockCommitLogMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockCommitLog) EXPECT() *MockCommitLogMockRecorder {
	return m.recorder
}

// Open mocks base method
func (m *MockCommitLog) Open() error {
	ret := m.ctrl.Call(m, "Open")
	ret0, _ := ret[0].(error)
	return ret0
}

// Open indicates an expected call of Open
func (mr *MockCommitLogMockRecorder) Open() *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Open", reflect.TypeOf((*MockCommitLog)(nil).Open))
}

// Write mocks base method
func (m *MockCommitLog) Write(ctx context.Context, series ts.Series, datapoint ts.Datapoint, unit time0.Unit, annotation ts.Annotation) error {
	ret := m.ctrl.Call(m, "Write", ctx, series, datapoint, unit, annotation)
	ret0, _ := ret[0].(error)
	return ret0
}

// Write indicates an expected call of Write
func (mr *MockCommitLogMockRecorder) Write(ctx, series, datapoint, unit, annotation interface{}) *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Write", reflect.TypeOf((*MockCommitLog)(nil).Write), ctx, series, datapoint, unit, annotation)
}

// WriteBatch mocks base method
func (m *MockCommitLog) WriteBatch(ctx context.Context, writes ts.WriteBatch) error {
	ret := m.ctrl.Call(m, "WriteBatch", ctx, writes)
	ret0, _ := ret[0].(error)
	return ret0
}

// WriteBatch indicates an expected call of WriteBatch
func (mr *MockCommitLogMockRecorder) WriteBatch(ctx, writes interface{}) *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "WriteBatch", reflect.TypeOf((*MockCommitLog)(nil).WriteBatch), ctx, writes)
}

// Close mocks base method
func (m *MockCommitLog) Close() error {
	ret := m.ctrl.Call(m, "Close")
	ret0, _ := ret[0].(error)
	return ret0
}

// Close indicates an expected call of Close
func (mr *MockCommitLogMockRecorder) Close() *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Close", reflect.TypeOf((*MockCommitLog)(nil).Close))
}

// ActiveLogs mocks base method
func (m *MockCommitLog) ActiveLogs() ([]fs.CommitlogFile, error) {
	ret := m.ctrl.Call(m, "ActiveLogs")
	ret0, _ := ret[0].([]fs.CommitlogFile)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ActiveLogs indicates an expected call of ActiveLogs
func (mr *MockCommitLogMockRecorder) ActiveLogs() *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ActiveLogs", reflect.TypeOf((*MockCommitLog)(nil).ActiveLogs))
}

// RotateLogs mocks base method
func (m *MockCommitLog) RotateLogs() (fs.CommitlogFile, error) {
	ret := m.ctrl.Call(m, "RotateLogs")
	ret0, _ := ret[0].(fs.CommitlogFile)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// RotateLogs indicates an expected call of RotateLogs
func (mr *MockCommitLogMockRecorder) RotateLogs() *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "RotateLogs", reflect.TypeOf((*MockCommitLog)(nil).RotateLogs))
}

// MockIterator is a mock of Iterator interface
type MockIterator struct {
	ctrl     *gomock.Controller
	recorder *MockIteratorMockRecorder
}

// MockIteratorMockRecorder is the mock recorder for MockIterator
type MockIteratorMockRecorder struct {
	mock *MockIterator
}

// NewMockIterator creates a new mock instance
func NewMockIterator(ctrl *gomock.Controller) *MockIterator {
	mock := &MockIterator{ctrl: ctrl}
	mock.recorder = &MockIteratorMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockIterator) EXPECT() *MockIteratorMockRecorder {
	return m.recorder
}

// Next mocks base method
func (m *MockIterator) Next() bool {
	ret := m.ctrl.Call(m, "Next")
	ret0, _ := ret[0].(bool)
	return ret0
}

// Next indicates an expected call of Next
func (mr *MockIteratorMockRecorder) Next() *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Next", reflect.TypeOf((*MockIterator)(nil).Next))
}

// Current mocks base method
func (m *MockIterator) Current() (ts.Series, ts.Datapoint, time0.Unit, ts.Annotation) {
	ret := m.ctrl.Call(m, "Current")
	ret0, _ := ret[0].(ts.Series)
	ret1, _ := ret[1].(ts.Datapoint)
	ret2, _ := ret[2].(time0.Unit)
	ret3, _ := ret[3].(ts.Annotation)
	return ret0, ret1, ret2, ret3
}

// Current indicates an expected call of Current
func (mr *MockIteratorMockRecorder) Current() *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Current", reflect.TypeOf((*MockIterator)(nil).Current))
}

// Err mocks base method
func (m *MockIterator) Err() error {
	ret := m.ctrl.Call(m, "Err")
	ret0, _ := ret[0].(error)
	return ret0
}

// Err indicates an expected call of Err
func (mr *MockIteratorMockRecorder) Err() *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Err", reflect.TypeOf((*MockIterator)(nil).Err))
}

// Close mocks base method
func (m *MockIterator) Close() {
	m.ctrl.Call(m, "Close")
}

// Close indicates an expected call of Close
func (mr *MockIteratorMockRecorder) Close() *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Close", reflect.TypeOf((*MockIterator)(nil).Close))
}

// MockOptions is a mock of Options interface
type MockOptions struct {
	ctrl     *gomock.Controller
	recorder *MockOptionsMockRecorder
}

// MockOptionsMockRecorder is the mock recorder for MockOptions
type MockOptionsMockRecorder struct {
	mock *MockOptions
}

// NewMockOptions creates a new mock instance
func NewMockOptions(ctrl *gomock.Controller) *MockOptions {
	mock := &MockOptions{ctrl: ctrl}
	mock.recorder = &MockOptionsMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockOptions) EXPECT() *MockOptionsMockRecorder {
	return m.recorder
}

// Validate mocks base method
func (m *MockOptions) Validate() error {
	ret := m.ctrl.Call(m, "Validate")
	ret0, _ := ret[0].(error)
	return ret0
}

// Validate indicates an expected call of Validate
func (mr *MockOptionsMockRecorder) Validate() *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Validate", reflect.TypeOf((*MockOptions)(nil).Validate))
}

// SetClockOptions mocks base method
func (m *MockOptions) SetClockOptions(value clock.Options) Options {
	ret := m.ctrl.Call(m, "SetClockOptions", value)
	ret0, _ := ret[0].(Options)
	return ret0
}

// SetClockOptions indicates an expected call of SetClockOptions
func (mr *MockOptionsMockRecorder) SetClockOptions(value interface{}) *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SetClockOptions", reflect.TypeOf((*MockOptions)(nil).SetClockOptions), value)
}

// ClockOptions mocks base method
func (m *MockOptions) ClockOptions() clock.Options {
	ret := m.ctrl.Call(m, "ClockOptions")
	ret0, _ := ret[0].(clock.Options)
	return ret0
}

// ClockOptions indicates an expected call of ClockOptions
func (mr *MockOptionsMockRecorder) ClockOptions() *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ClockOptions", reflect.TypeOf((*MockOptions)(nil).ClockOptions))
}

// SetInstrumentOptions mocks base method
func (m *MockOptions) SetInstrumentOptions(value instrument.Options) Options {
	ret := m.ctrl.Call(m, "SetInstrumentOptions", value)
	ret0, _ := ret[0].(Options)
	return ret0
}

// SetInstrumentOptions indicates an expected call of SetInstrumentOptions
func (mr *MockOptionsMockRecorder) SetInstrumentOptions(value interface{}) *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SetInstrumentOptions", reflect.TypeOf((*MockOptions)(nil).SetInstrumentOptions), value)
}

// InstrumentOptions mocks base method
func (m *MockOptions) InstrumentOptions() instrument.Options {
	ret := m.ctrl.Call(m, "InstrumentOptions")
	ret0, _ := ret[0].(instrument.Options)
	return ret0
}

// InstrumentOptions indicates an expected call of InstrumentOptions
func (mr *MockOptionsMockRecorder) InstrumentOptions() *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "InstrumentOptions", reflect.TypeOf((*MockOptions)(nil).InstrumentOptions))
}

// SetBlockSize mocks base method
func (m *MockOptions) SetBlockSize(value time.Duration) Options {
	ret := m.ctrl.Call(m, "SetBlockSize", value)
	ret0, _ := ret[0].(Options)
	return ret0
}

// SetBlockSize indicates an expected call of SetBlockSize
func (mr *MockOptionsMockRecorder) SetBlockSize(value interface{}) *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SetBlockSize", reflect.TypeOf((*MockOptions)(nil).SetBlockSize), value)
}

// BlockSize mocks base method
func (m *MockOptions) BlockSize() time.Duration {
	ret := m.ctrl.Call(m, "BlockSize")
	ret0, _ := ret[0].(time.Duration)
	return ret0
}

// BlockSize indicates an expected call of BlockSize
func (mr *MockOptionsMockRecorder) BlockSize() *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "BlockSize", reflect.TypeOf((*MockOptions)(nil).BlockSize))
}

// SetFilesystemOptions mocks base method
func (m *MockOptions) SetFilesystemOptions(value fs.Options) Options {
	ret := m.ctrl.Call(m, "SetFilesystemOptions", value)
	ret0, _ := ret[0].(Options)
	return ret0
}

// SetFilesystemOptions indicates an expected call of SetFilesystemOptions
func (mr *MockOptionsMockRecorder) SetFilesystemOptions(value interface{}) *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SetFilesystemOptions", reflect.TypeOf((*MockOptions)(nil).SetFilesystemOptions), value)
}

// FilesystemOptions mocks base method
func (m *MockOptions) FilesystemOptions() fs.Options {
	ret := m.ctrl.Call(m, "FilesystemOptions")
	ret0, _ := ret[0].(fs.Options)
	return ret0
}

// FilesystemOptions indicates an expected call of FilesystemOptions
func (mr *MockOptionsMockRecorder) FilesystemOptions() *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "FilesystemOptions", reflect.TypeOf((*MockOptions)(nil).FilesystemOptions))
}

// SetFlushSize mocks base method
func (m *MockOptions) SetFlushSize(value int) Options {
	ret := m.ctrl.Call(m, "SetFlushSize", value)
	ret0, _ := ret[0].(Options)
	return ret0
}

// SetFlushSize indicates an expected call of SetFlushSize
func (mr *MockOptionsMockRecorder) SetFlushSize(value interface{}) *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SetFlushSize", reflect.TypeOf((*MockOptions)(nil).SetFlushSize), value)
}

// FlushSize mocks base method
func (m *MockOptions) FlushSize() int {
	ret := m.ctrl.Call(m, "FlushSize")
	ret0, _ := ret[0].(int)
	return ret0
}

// FlushSize indicates an expected call of FlushSize
func (mr *MockOptionsMockRecorder) FlushSize() *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "FlushSize", reflect.TypeOf((*MockOptions)(nil).FlushSize))
}

// SetStrategy mocks base method
func (m *MockOptions) SetStrategy(value Strategy) Options {
	ret := m.ctrl.Call(m, "SetStrategy", value)
	ret0, _ := ret[0].(Options)
	return ret0
}

// SetStrategy indicates an expected call of SetStrategy
func (mr *MockOptionsMockRecorder) SetStrategy(value interface{}) *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SetStrategy", reflect.TypeOf((*MockOptions)(nil).SetStrategy), value)
}

// Strategy mocks base method
func (m *MockOptions) Strategy() Strategy {
	ret := m.ctrl.Call(m, "Strategy")
	ret0, _ := ret[0].(Strategy)
	return ret0
}

// Strategy indicates an expected call of Strategy
func (mr *MockOptionsMockRecorder) Strategy() *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Strategy", reflect.TypeOf((*MockOptions)(nil).Strategy))
}

// SetFlushInterval mocks base method
func (m *MockOptions) SetFlushInterval(value time.Duration) Options {
	ret := m.ctrl.Call(m, "SetFlushInterval", value)
	ret0, _ := ret[0].(Options)
	return ret0
}

// SetFlushInterval indicates an expected call of SetFlushInterval
func (mr *MockOptionsMockRecorder) SetFlushInterval(value interface{}) *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SetFlushInterval", reflect.TypeOf((*MockOptions)(nil).SetFlushInterval), value)
}

// FlushInterval mocks base method
func (m *MockOptions) FlushInterval() time.Duration {
	ret := m.ctrl.Call(m, "FlushInterval")
	ret0, _ := ret[0].(time.Duration)
	return ret0
}

// FlushInterval indicates an expected call of FlushInterval
func (mr *MockOptionsMockRecorder) FlushInterval() *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "FlushInterval", reflect.TypeOf((*MockOptions)(nil).FlushInterval))
}

// SetBacklogQueueSize mocks base method
func (m *MockOptions) SetBacklogQueueSize(value int) Options {
	ret := m.ctrl.Call(m, "SetBacklogQueueSize", value)
	ret0, _ := ret[0].(Options)
	return ret0
}

// SetBacklogQueueSize indicates an expected call of SetBacklogQueueSize
func (mr *MockOptionsMockRecorder) SetBacklogQueueSize(value interface{}) *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SetBacklogQueueSize", reflect.TypeOf((*MockOptions)(nil).SetBacklogQueueSize), value)
}

// BacklogQueueSize mocks base method
func (m *MockOptions) BacklogQueueSize() int {
	ret := m.ctrl.Call(m, "BacklogQueueSize")
	ret0, _ := ret[0].(int)
	return ret0
}

// BacklogQueueSize indicates an expected call of BacklogQueueSize
func (mr *MockOptionsMockRecorder) BacklogQueueSize() *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "BacklogQueueSize", reflect.TypeOf((*MockOptions)(nil).BacklogQueueSize))
}

// SetBacklogQueueChannelSize mocks base method
func (m *MockOptions) SetBacklogQueueChannelSize(value int) Options {
	ret := m.ctrl.Call(m, "SetBacklogQueueChannelSize", value)
	ret0, _ := ret[0].(Options)
	return ret0
}

// SetBacklogQueueChannelSize indicates an expected call of SetBacklogQueueChannelSize
func (mr *MockOptionsMockRecorder) SetBacklogQueueChannelSize(value interface{}) *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SetBacklogQueueChannelSize", reflect.TypeOf((*MockOptions)(nil).SetBacklogQueueChannelSize), value)
}

// BacklogQueueChannelSize mocks base method
func (m *MockOptions) BacklogQueueChannelSize() int {
	ret := m.ctrl.Call(m, "BacklogQueueChannelSize")
	ret0, _ := ret[0].(int)
	return ret0
}

// BacklogQueueChannelSize indicates an expected call of BacklogQueueChannelSize
func (mr *MockOptionsMockRecorder) BacklogQueueChannelSize() *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "BacklogQueueChannelSize", reflect.TypeOf((*MockOptions)(nil).BacklogQueueChannelSize))
}

// SetBytesPool mocks base method
func (m *MockOptions) SetBytesPool(value pool.CheckedBytesPool) Options {
	ret := m.ctrl.Call(m, "SetBytesPool", value)
	ret0, _ := ret[0].(Options)
	return ret0
}

// SetBytesPool indicates an expected call of SetBytesPool
func (mr *MockOptionsMockRecorder) SetBytesPool(value interface{}) *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SetBytesPool", reflect.TypeOf((*MockOptions)(nil).SetBytesPool), value)
}

// BytesPool mocks base method
func (m *MockOptions) BytesPool() pool.CheckedBytesPool {
	ret := m.ctrl.Call(m, "BytesPool")
	ret0, _ := ret[0].(pool.CheckedBytesPool)
	return ret0
}

// BytesPool indicates an expected call of BytesPool
func (mr *MockOptionsMockRecorder) BytesPool() *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "BytesPool", reflect.TypeOf((*MockOptions)(nil).BytesPool))
}

// SetReadConcurrency mocks base method
func (m *MockOptions) SetReadConcurrency(concurrency int) Options {
	ret := m.ctrl.Call(m, "SetReadConcurrency", concurrency)
	ret0, _ := ret[0].(Options)
	return ret0
}

// SetReadConcurrency indicates an expected call of SetReadConcurrency
func (mr *MockOptionsMockRecorder) SetReadConcurrency(concurrency interface{}) *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SetReadConcurrency", reflect.TypeOf((*MockOptions)(nil).SetReadConcurrency), concurrency)
}

// ReadConcurrency mocks base method
func (m *MockOptions) ReadConcurrency() int {
	ret := m.ctrl.Call(m, "ReadConcurrency")
	ret0, _ := ret[0].(int)
	return ret0
}

// ReadConcurrency indicates an expected call of ReadConcurrency
func (mr *MockOptionsMockRecorder) ReadConcurrency() *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ReadConcurrency", reflect.TypeOf((*MockOptions)(nil).ReadConcurrency))
}

// SetIdentifierPool mocks base method
func (m *MockOptions) SetIdentifierPool(value ident.Pool) Options {
	ret := m.ctrl.Call(m, "SetIdentifierPool", value)
	ret0, _ := ret[0].(Options)
	return ret0
}

// SetIdentifierPool indicates an expected call of SetIdentifierPool
func (mr *MockOptionsMockRecorder) SetIdentifierPool(value interface{}) *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SetIdentifierPool", reflect.TypeOf((*MockOptions)(nil).SetIdentifierPool), value)
}

// IdentifierPool mocks base method
func (m *MockOptions) IdentifierPool() ident.Pool {
	ret := m.ctrl.Call(m, "IdentifierPool")
	ret0, _ := ret[0].(ident.Pool)
	return ret0
}

// IdentifierPool indicates an expected call of IdentifierPool
func (mr *MockOptionsMockRecorder) IdentifierPool() *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "IdentifierPool", reflect.TypeOf((*MockOptions)(nil).IdentifierPool))
}
