// Code generated by MockGen. DO NOT EDIT.
// Source: github.com/ilborsch/openai-go/openai (interfaces: OpenAIClient)

// Package mocks is a generated GoMock package.
package mocks

import (
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
	assistants "github.com/ilborsch/openai-go/openai/assistants"
	messages "github.com/ilborsch/openai-go/openai/assistants/messages"
	runs "github.com/ilborsch/openai-go/openai/assistants/runs"
	vecstores "github.com/ilborsch/openai-go/openai/assistants/vector-stores"
	message "github.com/ilborsch/openai-go/openai/chatgpt/message"
)

// MockOpenAIClient is a mock of OpenAIClient interface.
type MockOpenAIClient struct {
	ctrl     *gomock.Controller
	recorder *MockOpenAIClientMockRecorder
}

// MockOpenAIClientMockRecorder is the mock recorder for MockOpenAIClient.
type MockOpenAIClientMockRecorder struct {
	mock *MockOpenAIClient
}

// NewMockOpenAIClient creates a new mock instance.
func NewMockOpenAIClient(ctrl *gomock.Controller) *MockOpenAIClient {
	mock := &MockOpenAIClient{ctrl: ctrl}
	mock.recorder = &MockOpenAIClientMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockOpenAIClient) EXPECT() *MockOpenAIClientMockRecorder {
	return m.recorder
}

// AddMessageToThread mocks base method.
func (m *MockOpenAIClient) AddMessageToThread(arg0, arg1 string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "AddMessageToThread", arg0, arg1)
	ret0, _ := ret[0].(error)
	return ret0
}

// AddMessageToThread indicates an expected call of AddMessageToThread.
func (mr *MockOpenAIClientMockRecorder) AddMessageToThread(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "AddMessageToThread", reflect.TypeOf((*MockOpenAIClient)(nil).AddMessageToThread), arg0, arg1)
}

// AddVectorStoreFile mocks base method.
func (m *MockOpenAIClient) AddVectorStoreFile(arg0, arg1 string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "AddVectorStoreFile", arg0, arg1)
	ret0, _ := ret[0].(error)
	return ret0
}

// AddVectorStoreFile indicates an expected call of AddVectorStoreFile.
func (mr *MockOpenAIClientMockRecorder) AddVectorStoreFile(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "AddVectorStoreFile", reflect.TypeOf((*MockOpenAIClient)(nil).AddVectorStoreFile), arg0, arg1)
}

// CreateAssistant mocks base method.
func (m *MockOpenAIClient) CreateAssistant(arg0, arg1, arg2 string, arg3 []assistants.Tool) (string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateAssistant", arg0, arg1, arg2, arg3)
	ret0, _ := ret[0].(string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CreateAssistant indicates an expected call of CreateAssistant.
func (mr *MockOpenAIClientMockRecorder) CreateAssistant(arg0, arg1, arg2, arg3 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateAssistant", reflect.TypeOf((*MockOpenAIClient)(nil).CreateAssistant), arg0, arg1, arg2, arg3)
}

// CreateCompletion mocks base method.
func (m *MockOpenAIClient) CreateCompletion(arg0 []message.Message) (string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateCompletion", arg0)
	ret0, _ := ret[0].(string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CreateCompletion indicates an expected call of CreateCompletion.
func (mr *MockOpenAIClientMockRecorder) CreateCompletion(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateCompletion", reflect.TypeOf((*MockOpenAIClient)(nil).CreateCompletion), arg0)
}

// CreateRun mocks base method.
func (m *MockOpenAIClient) CreateRun(arg0, arg1 string) (string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateRun", arg0, arg1)
	ret0, _ := ret[0].(string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CreateRun indicates an expected call of CreateRun.
func (mr *MockOpenAIClientMockRecorder) CreateRun(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateRun", reflect.TypeOf((*MockOpenAIClient)(nil).CreateRun), arg0, arg1)
}

// CreateThread mocks base method.
func (m *MockOpenAIClient) CreateThread() (string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateThread")
	ret0, _ := ret[0].(string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CreateThread indicates an expected call of CreateThread.
func (mr *MockOpenAIClientMockRecorder) CreateThread() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateThread", reflect.TypeOf((*MockOpenAIClient)(nil).CreateThread))
}

// CreateVectorStore mocks base method.
func (m *MockOpenAIClient) CreateVectorStore(arg0 string) (string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateVectorStore", arg0)
	ret0, _ := ret[0].(string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CreateVectorStore indicates an expected call of CreateVectorStore.
func (mr *MockOpenAIClientMockRecorder) CreateVectorStore(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateVectorStore", reflect.TypeOf((*MockOpenAIClient)(nil).CreateVectorStore), arg0)
}

// DeleteAssistant mocks base method.
func (m *MockOpenAIClient) DeleteAssistant(arg0 string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DeleteAssistant", arg0)
	ret0, _ := ret[0].(error)
	return ret0
}

// DeleteAssistant indicates an expected call of DeleteAssistant.
func (mr *MockOpenAIClientMockRecorder) DeleteAssistant(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeleteAssistant", reflect.TypeOf((*MockOpenAIClient)(nil).DeleteAssistant), arg0)
}

// DeleteFile mocks base method.
func (m *MockOpenAIClient) DeleteFile(arg0 string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DeleteFile", arg0)
	ret0, _ := ret[0].(error)
	return ret0
}

// DeleteFile indicates an expected call of DeleteFile.
func (mr *MockOpenAIClientMockRecorder) DeleteFile(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeleteFile", reflect.TypeOf((*MockOpenAIClient)(nil).DeleteFile), arg0)
}

// DeleteVectorStore mocks base method.
func (m *MockOpenAIClient) DeleteVectorStore(arg0 string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DeleteVectorStore", arg0)
	ret0, _ := ret[0].(error)
	return ret0
}

// DeleteVectorStore indicates an expected call of DeleteVectorStore.
func (mr *MockOpenAIClientMockRecorder) DeleteVectorStore(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeleteVectorStore", reflect.TypeOf((*MockOpenAIClient)(nil).DeleteVectorStore), arg0)
}

// DeleteVectorStoreFile mocks base method.
func (m *MockOpenAIClient) DeleteVectorStoreFile(arg0, arg1 string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DeleteVectorStoreFile", arg0, arg1)
	ret0, _ := ret[0].(error)
	return ret0
}

// DeleteVectorStoreFile indicates an expected call of DeleteVectorStoreFile.
func (mr *MockOpenAIClientMockRecorder) DeleteVectorStoreFile(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeleteVectorStoreFile", reflect.TypeOf((*MockOpenAIClient)(nil).DeleteVectorStoreFile), arg0, arg1)
}

// GetAssistant mocks base method.
func (m *MockOpenAIClient) GetAssistant(arg0 string) (assistants.GetAssistantResponse, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetAssistant", arg0)
	ret0, _ := ret[0].(assistants.GetAssistantResponse)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetAssistant indicates an expected call of GetAssistant.
func (mr *MockOpenAIClientMockRecorder) GetAssistant(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetAssistant", reflect.TypeOf((*MockOpenAIClient)(nil).GetAssistant), arg0)
}

// GetRun mocks base method.
func (m *MockOpenAIClient) GetRun(arg0, arg1 string) (runs.GetRunResponse, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetRun", arg0, arg1)
	ret0, _ := ret[0].(runs.GetRunResponse)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetRun indicates an expected call of GetRun.
func (mr *MockOpenAIClientMockRecorder) GetRun(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetRun", reflect.TypeOf((*MockOpenAIClient)(nil).GetRun), arg0, arg1)
}

// GetThreadMessages mocks base method.
func (m *MockOpenAIClient) GetThreadMessages(arg0 string) (messages.ThreadMessages, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetThreadMessages", arg0)
	ret0, _ := ret[0].(messages.ThreadMessages)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetThreadMessages indicates an expected call of GetThreadMessages.
func (mr *MockOpenAIClientMockRecorder) GetThreadMessages(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetThreadMessages", reflect.TypeOf((*MockOpenAIClient)(nil).GetThreadMessages), arg0)
}

// GetVectorStoreFiles mocks base method.
func (m *MockOpenAIClient) GetVectorStoreFiles(arg0 string) (vecstores.GetVectorStoreFilesResponse, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetVectorStoreFiles", arg0)
	ret0, _ := ret[0].(vecstores.GetVectorStoreFilesResponse)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetVectorStoreFiles indicates an expected call of GetVectorStoreFiles.
func (mr *MockOpenAIClientMockRecorder) GetVectorStoreFiles(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetVectorStoreFiles", reflect.TypeOf((*MockOpenAIClient)(nil).GetVectorStoreFiles), arg0)
}

// LatestAssistantResponse mocks base method.
func (m *MockOpenAIClient) LatestAssistantResponse(arg0 string) (string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "LatestAssistantResponse", arg0)
	ret0, _ := ret[0].(string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// LatestAssistantResponse indicates an expected call of LatestAssistantResponse.
func (mr *MockOpenAIClientMockRecorder) LatestAssistantResponse(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "LatestAssistantResponse", reflect.TypeOf((*MockOpenAIClient)(nil).LatestAssistantResponse), arg0)
}

// Modify mocks base method.
func (m *MockOpenAIClient) Modify(arg0, arg1, arg2 string, arg3 float32) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Modify", arg0, arg1, arg2, arg3)
	ret0, _ := ret[0].(error)
	return ret0
}

// Modify indicates an expected call of Modify.
func (mr *MockOpenAIClientMockRecorder) Modify(arg0, arg1, arg2, arg3 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Modify", reflect.TypeOf((*MockOpenAIClient)(nil).Modify), arg0, arg1, arg2, arg3)
}

// UploadFile mocks base method.
func (m *MockOpenAIClient) UploadFile(arg0 string, arg1 []byte) (string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UploadFile", arg0, arg1)
	ret0, _ := ret[0].(string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// UploadFile indicates an expected call of UploadFile.
func (mr *MockOpenAIClientMockRecorder) UploadFile(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UploadFile", reflect.TypeOf((*MockOpenAIClient)(nil).UploadFile), arg0, arg1)
}
