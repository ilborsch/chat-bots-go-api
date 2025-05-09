// Code generated by MockGen. DO NOT EDIT.
// Source: chat-bots-api/internal/repository (interfaces: ChatBotRepository)

// Package mocks is a generated GoMock package.
package mocks

import (
	domain "chat-bots-api/domain"
	context "context"
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
)

// MockChatBotRepository is a mock of ChatBotRepository interface.
type MockChatBotRepository struct {
	ctrl     *gomock.Controller
	recorder *MockChatBotRepositoryMockRecorder
}

// MockChatBotRepositoryMockRecorder is the mock recorder for MockChatBotRepository.
type MockChatBotRepositoryMockRecorder struct {
	mock *MockChatBotRepository
}

// NewMockChatBotRepository creates a new mock instance.
func NewMockChatBotRepository(ctrl *gomock.Controller) *MockChatBotRepository {
	mock := &MockChatBotRepository{ctrl: ctrl}
	mock.recorder = &MockChatBotRepositoryMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockChatBotRepository) EXPECT() *MockChatBotRepositoryMockRecorder {
	return m.recorder
}

// ChatBot mocks base method.
func (m *MockChatBotRepository) ChatBot(arg0 context.Context, arg1, arg2 int64) (domain.ChatBot, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ChatBot", arg0, arg1, arg2)
	ret0, _ := ret[0].(domain.ChatBot)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ChatBot indicates an expected call of ChatBot.
func (mr *MockChatBotRepositoryMockRecorder) ChatBot(arg0, arg1, arg2 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ChatBot", reflect.TypeOf((*MockChatBotRepository)(nil).ChatBot), arg0, arg1, arg2)
}

// RemoveChatBot mocks base method.
func (m *MockChatBotRepository) RemoveChatBot(arg0 context.Context, arg1, arg2 int64) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "RemoveChatBot", arg0, arg1, arg2)
	ret0, _ := ret[0].(error)
	return ret0
}

// RemoveChatBot indicates an expected call of RemoveChatBot.
func (mr *MockChatBotRepositoryMockRecorder) RemoveChatBot(arg0, arg1, arg2 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "RemoveChatBot", reflect.TypeOf((*MockChatBotRepository)(nil).RemoveChatBot), arg0, arg1, arg2)
}

// SaveChatBot mocks base method.
func (m *MockChatBotRepository) SaveChatBot(arg0 context.Context, arg1 domain.ChatBot) (int64, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SaveChatBot", arg0, arg1)
	ret0, _ := ret[0].(int64)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// SaveChatBot indicates an expected call of SaveChatBot.
func (mr *MockChatBotRepositoryMockRecorder) SaveChatBot(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SaveChatBot", reflect.TypeOf((*MockChatBotRepository)(nil).SaveChatBot), arg0, arg1)
}

// UpdateChatBot mocks base method.
func (m *MockChatBotRepository) UpdateChatBot(arg0 context.Context, arg1, arg2 int64, arg3 domain.ChatBot) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UpdateChatBot", arg0, arg1, arg2, arg3)
	ret0, _ := ret[0].(error)
	return ret0
}

// UpdateChatBot indicates an expected call of UpdateChatBot.
func (mr *MockChatBotRepositoryMockRecorder) UpdateChatBot(arg0, arg1, arg2, arg3 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdateChatBot", reflect.TypeOf((*MockChatBotRepository)(nil).UpdateChatBot), arg0, arg1, arg2, arg3)
}

// UserChatBots mocks base method.
func (m *MockChatBotRepository) UserChatBots(arg0 context.Context, arg1 int64) ([]domain.ChatBot, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UserChatBots", arg0, arg1)
	ret0, _ := ret[0].([]domain.ChatBot)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// UserChatBots indicates an expected call of UserChatBots.
func (mr *MockChatBotRepositoryMockRecorder) UserChatBots(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UserChatBots", reflect.TypeOf((*MockChatBotRepository)(nil).UserChatBots), arg0, arg1)
}
