// Code generated by MockGen. DO NOT EDIT.
// Source: sender.go

// Package mock is a generated GoMock package.
package mock

import (
	context "context"
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
	models "github.com/kosalnik/metrics/internal/models"
)

// MockSender is a mock of Sender interface.
type MockSender struct {
	ctrl     *gomock.Controller
	recorder *MockSenderMockRecorder
}

// MockSenderMockRecorder is the mock recorder for MockSender.
type MockSenderMockRecorder struct {
	mock *MockSender
}

// NewMockSender creates a new mock instance.
func NewMockSender(ctrl *gomock.Controller) *MockSender {
	mock := &MockSender{ctrl: ctrl}
	mock.recorder = &MockSenderMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockSender) EXPECT() *MockSenderMockRecorder {
	return m.recorder
}

// SendBatch mocks base method.
func (m *MockSender) SendBatch(ctx context.Context, list []models.Metrics) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SendBatch", ctx, list)
	ret0, _ := ret[0].(error)
	return ret0
}

// SendBatch indicates an expected call of SendBatch.
func (mr *MockSenderMockRecorder) SendBatch(ctx, list interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SendBatch", reflect.TypeOf((*MockSender)(nil).SendBatch), ctx, list)
}

// SendCounter mocks base method.
func (m *MockSender) SendCounter(k string, v int64) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "SendCounter", k, v)
}

// SendCounter indicates an expected call of SendCounter.
func (mr *MockSenderMockRecorder) SendCounter(k, v interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SendCounter", reflect.TypeOf((*MockSender)(nil).SendCounter), k, v)
}

// SendGauge mocks base method.
func (m *MockSender) SendGauge(k string, v float64) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "SendGauge", k, v)
}

// SendGauge indicates an expected call of SendGauge.
func (mr *MockSenderMockRecorder) SendGauge(k, v interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SendGauge", reflect.TypeOf((*MockSender)(nil).SendGauge), k, v)
}
