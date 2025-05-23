// Code generated by counterfeiter. DO NOT EDIT.
package constructfakes

import (
	"sync"

	"github.com/cloudfoundry/bosh-windows-stemcell-builder/stembuild/construct"
)

type FakeConstructMessenger struct {
	CreateProvisionDirStartedStub        func()
	createProvisionDirStartedMutex       sync.RWMutex
	createProvisionDirStartedArgsForCall []struct {
	}
	CreateProvisionDirSucceededStub        func()
	createProvisionDirSucceededMutex       sync.RWMutex
	createProvisionDirSucceededArgsForCall []struct {
	}
	EnableWinRMStartedStub        func()
	enableWinRMStartedMutex       sync.RWMutex
	enableWinRMStartedArgsForCall []struct {
	}
	EnableWinRMSucceededStub        func()
	enableWinRMSucceededMutex       sync.RWMutex
	enableWinRMSucceededArgsForCall []struct {
	}
	ExecutePostRebootScriptStartedStub        func()
	executePostRebootScriptStartedMutex       sync.RWMutex
	executePostRebootScriptStartedArgsForCall []struct {
	}
	ExecutePostRebootScriptSucceededStub        func()
	executePostRebootScriptSucceededMutex       sync.RWMutex
	executePostRebootScriptSucceededArgsForCall []struct {
	}
	ExecutePostRebootWarningStub        func(string)
	executePostRebootWarningMutex       sync.RWMutex
	executePostRebootWarningArgsForCall []struct {
		arg1 string
	}
	ExecuteSetupScriptStartedStub        func()
	executeSetupScriptStartedMutex       sync.RWMutex
	executeSetupScriptStartedArgsForCall []struct {
	}
	ExecuteSetupScriptSucceededStub        func()
	executeSetupScriptSucceededMutex       sync.RWMutex
	executeSetupScriptSucceededArgsForCall []struct {
	}
	ExtractArtifactsStartedStub        func()
	extractArtifactsStartedMutex       sync.RWMutex
	extractArtifactsStartedArgsForCall []struct {
	}
	ExtractArtifactsSucceededStub        func()
	extractArtifactsSucceededMutex       sync.RWMutex
	extractArtifactsSucceededArgsForCall []struct {
	}
	LogOutUsersStartedStub        func()
	logOutUsersStartedMutex       sync.RWMutex
	logOutUsersStartedArgsForCall []struct {
	}
	LogOutUsersSucceededStub        func()
	logOutUsersSucceededMutex       sync.RWMutex
	logOutUsersSucceededArgsForCall []struct {
	}
	RebootHasFinishedStub        func()
	rebootHasFinishedMutex       sync.RWMutex
	rebootHasFinishedArgsForCall []struct {
	}
	RebootHasStartedStub        func()
	rebootHasStartedMutex       sync.RWMutex
	rebootHasStartedArgsForCall []struct {
	}
	ShutdownCompletedStub        func()
	shutdownCompletedMutex       sync.RWMutex
	shutdownCompletedArgsForCall []struct {
	}
	UploadArtifactsStartedStub        func()
	uploadArtifactsStartedMutex       sync.RWMutex
	uploadArtifactsStartedArgsForCall []struct {
	}
	UploadArtifactsSucceededStub        func()
	uploadArtifactsSucceededMutex       sync.RWMutex
	uploadArtifactsSucceededArgsForCall []struct {
	}
	UploadFileStartedStub        func(string)
	uploadFileStartedMutex       sync.RWMutex
	uploadFileStartedArgsForCall []struct {
		arg1 string
	}
	UploadFileSucceededStub        func()
	uploadFileSucceededMutex       sync.RWMutex
	uploadFileSucceededArgsForCall []struct {
	}
	ValidateVMConnectionStartedStub        func()
	validateVMConnectionStartedMutex       sync.RWMutex
	validateVMConnectionStartedArgsForCall []struct {
	}
	ValidateVMConnectionSucceededStub        func()
	validateVMConnectionSucceededMutex       sync.RWMutex
	validateVMConnectionSucceededArgsForCall []struct {
	}
	WaitingForShutdownStub        func()
	waitingForShutdownMutex       sync.RWMutex
	waitingForShutdownArgsForCall []struct {
	}
	WinRMDisconnectedForRebootStub        func()
	winRMDisconnectedForRebootMutex       sync.RWMutex
	winRMDisconnectedForRebootArgsForCall []struct {
	}
	invocations      map[string][][]interface{}
	invocationsMutex sync.RWMutex
}

func (fake *FakeConstructMessenger) CreateProvisionDirStarted() {
	fake.createProvisionDirStartedMutex.Lock()
	fake.createProvisionDirStartedArgsForCall = append(fake.createProvisionDirStartedArgsForCall, struct {
	}{})
	stub := fake.CreateProvisionDirStartedStub
	fake.recordInvocation("CreateProvisionDirStarted", []interface{}{})
	fake.createProvisionDirStartedMutex.Unlock()
	if stub != nil {
		fake.CreateProvisionDirStartedStub()
	}
}

func (fake *FakeConstructMessenger) CreateProvisionDirStartedCallCount() int {
	fake.createProvisionDirStartedMutex.RLock()
	defer fake.createProvisionDirStartedMutex.RUnlock()
	return len(fake.createProvisionDirStartedArgsForCall)
}

func (fake *FakeConstructMessenger) CreateProvisionDirStartedCalls(stub func()) {
	fake.createProvisionDirStartedMutex.Lock()
	defer fake.createProvisionDirStartedMutex.Unlock()
	fake.CreateProvisionDirStartedStub = stub
}

func (fake *FakeConstructMessenger) CreateProvisionDirSucceeded() {
	fake.createProvisionDirSucceededMutex.Lock()
	fake.createProvisionDirSucceededArgsForCall = append(fake.createProvisionDirSucceededArgsForCall, struct {
	}{})
	stub := fake.CreateProvisionDirSucceededStub
	fake.recordInvocation("CreateProvisionDirSucceeded", []interface{}{})
	fake.createProvisionDirSucceededMutex.Unlock()
	if stub != nil {
		fake.CreateProvisionDirSucceededStub()
	}
}

func (fake *FakeConstructMessenger) CreateProvisionDirSucceededCallCount() int {
	fake.createProvisionDirSucceededMutex.RLock()
	defer fake.createProvisionDirSucceededMutex.RUnlock()
	return len(fake.createProvisionDirSucceededArgsForCall)
}

func (fake *FakeConstructMessenger) CreateProvisionDirSucceededCalls(stub func()) {
	fake.createProvisionDirSucceededMutex.Lock()
	defer fake.createProvisionDirSucceededMutex.Unlock()
	fake.CreateProvisionDirSucceededStub = stub
}

func (fake *FakeConstructMessenger) EnableWinRMStarted() {
	fake.enableWinRMStartedMutex.Lock()
	fake.enableWinRMStartedArgsForCall = append(fake.enableWinRMStartedArgsForCall, struct {
	}{})
	stub := fake.EnableWinRMStartedStub
	fake.recordInvocation("EnableWinRMStarted", []interface{}{})
	fake.enableWinRMStartedMutex.Unlock()
	if stub != nil {
		fake.EnableWinRMStartedStub()
	}
}

func (fake *FakeConstructMessenger) EnableWinRMStartedCallCount() int {
	fake.enableWinRMStartedMutex.RLock()
	defer fake.enableWinRMStartedMutex.RUnlock()
	return len(fake.enableWinRMStartedArgsForCall)
}

func (fake *FakeConstructMessenger) EnableWinRMStartedCalls(stub func()) {
	fake.enableWinRMStartedMutex.Lock()
	defer fake.enableWinRMStartedMutex.Unlock()
	fake.EnableWinRMStartedStub = stub
}

func (fake *FakeConstructMessenger) EnableWinRMSucceeded() {
	fake.enableWinRMSucceededMutex.Lock()
	fake.enableWinRMSucceededArgsForCall = append(fake.enableWinRMSucceededArgsForCall, struct {
	}{})
	stub := fake.EnableWinRMSucceededStub
	fake.recordInvocation("EnableWinRMSucceeded", []interface{}{})
	fake.enableWinRMSucceededMutex.Unlock()
	if stub != nil {
		fake.EnableWinRMSucceededStub()
	}
}

func (fake *FakeConstructMessenger) EnableWinRMSucceededCallCount() int {
	fake.enableWinRMSucceededMutex.RLock()
	defer fake.enableWinRMSucceededMutex.RUnlock()
	return len(fake.enableWinRMSucceededArgsForCall)
}

func (fake *FakeConstructMessenger) EnableWinRMSucceededCalls(stub func()) {
	fake.enableWinRMSucceededMutex.Lock()
	defer fake.enableWinRMSucceededMutex.Unlock()
	fake.EnableWinRMSucceededStub = stub
}

func (fake *FakeConstructMessenger) ExecutePostRebootScriptStarted() {
	fake.executePostRebootScriptStartedMutex.Lock()
	fake.executePostRebootScriptStartedArgsForCall = append(fake.executePostRebootScriptStartedArgsForCall, struct {
	}{})
	stub := fake.ExecutePostRebootScriptStartedStub
	fake.recordInvocation("ExecutePostRebootScriptStarted", []interface{}{})
	fake.executePostRebootScriptStartedMutex.Unlock()
	if stub != nil {
		fake.ExecutePostRebootScriptStartedStub()
	}
}

func (fake *FakeConstructMessenger) ExecutePostRebootScriptStartedCallCount() int {
	fake.executePostRebootScriptStartedMutex.RLock()
	defer fake.executePostRebootScriptStartedMutex.RUnlock()
	return len(fake.executePostRebootScriptStartedArgsForCall)
}

func (fake *FakeConstructMessenger) ExecutePostRebootScriptStartedCalls(stub func()) {
	fake.executePostRebootScriptStartedMutex.Lock()
	defer fake.executePostRebootScriptStartedMutex.Unlock()
	fake.ExecutePostRebootScriptStartedStub = stub
}

func (fake *FakeConstructMessenger) ExecutePostRebootScriptSucceeded() {
	fake.executePostRebootScriptSucceededMutex.Lock()
	fake.executePostRebootScriptSucceededArgsForCall = append(fake.executePostRebootScriptSucceededArgsForCall, struct {
	}{})
	stub := fake.ExecutePostRebootScriptSucceededStub
	fake.recordInvocation("ExecutePostRebootScriptSucceeded", []interface{}{})
	fake.executePostRebootScriptSucceededMutex.Unlock()
	if stub != nil {
		fake.ExecutePostRebootScriptSucceededStub()
	}
}

func (fake *FakeConstructMessenger) ExecutePostRebootScriptSucceededCallCount() int {
	fake.executePostRebootScriptSucceededMutex.RLock()
	defer fake.executePostRebootScriptSucceededMutex.RUnlock()
	return len(fake.executePostRebootScriptSucceededArgsForCall)
}

func (fake *FakeConstructMessenger) ExecutePostRebootScriptSucceededCalls(stub func()) {
	fake.executePostRebootScriptSucceededMutex.Lock()
	defer fake.executePostRebootScriptSucceededMutex.Unlock()
	fake.ExecutePostRebootScriptSucceededStub = stub
}

func (fake *FakeConstructMessenger) ExecutePostRebootWarning(arg1 string) {
	fake.executePostRebootWarningMutex.Lock()
	fake.executePostRebootWarningArgsForCall = append(fake.executePostRebootWarningArgsForCall, struct {
		arg1 string
	}{arg1})
	stub := fake.ExecutePostRebootWarningStub
	fake.recordInvocation("ExecutePostRebootWarning", []interface{}{arg1})
	fake.executePostRebootWarningMutex.Unlock()
	if stub != nil {
		fake.ExecutePostRebootWarningStub(arg1)
	}
}

func (fake *FakeConstructMessenger) ExecutePostRebootWarningCallCount() int {
	fake.executePostRebootWarningMutex.RLock()
	defer fake.executePostRebootWarningMutex.RUnlock()
	return len(fake.executePostRebootWarningArgsForCall)
}

func (fake *FakeConstructMessenger) ExecutePostRebootWarningCalls(stub func(string)) {
	fake.executePostRebootWarningMutex.Lock()
	defer fake.executePostRebootWarningMutex.Unlock()
	fake.ExecutePostRebootWarningStub = stub
}

func (fake *FakeConstructMessenger) ExecutePostRebootWarningArgsForCall(i int) string {
	fake.executePostRebootWarningMutex.RLock()
	defer fake.executePostRebootWarningMutex.RUnlock()
	argsForCall := fake.executePostRebootWarningArgsForCall[i]
	return argsForCall.arg1
}

func (fake *FakeConstructMessenger) ExecuteSetupScriptStarted() {
	fake.executeSetupScriptStartedMutex.Lock()
	fake.executeSetupScriptStartedArgsForCall = append(fake.executeSetupScriptStartedArgsForCall, struct {
	}{})
	stub := fake.ExecuteSetupScriptStartedStub
	fake.recordInvocation("ExecuteSetupScriptStarted", []interface{}{})
	fake.executeSetupScriptStartedMutex.Unlock()
	if stub != nil {
		fake.ExecuteSetupScriptStartedStub()
	}
}

func (fake *FakeConstructMessenger) ExecuteSetupScriptStartedCallCount() int {
	fake.executeSetupScriptStartedMutex.RLock()
	defer fake.executeSetupScriptStartedMutex.RUnlock()
	return len(fake.executeSetupScriptStartedArgsForCall)
}

func (fake *FakeConstructMessenger) ExecuteSetupScriptStartedCalls(stub func()) {
	fake.executeSetupScriptStartedMutex.Lock()
	defer fake.executeSetupScriptStartedMutex.Unlock()
	fake.ExecuteSetupScriptStartedStub = stub
}

func (fake *FakeConstructMessenger) ExecuteSetupScriptSucceeded() {
	fake.executeSetupScriptSucceededMutex.Lock()
	fake.executeSetupScriptSucceededArgsForCall = append(fake.executeSetupScriptSucceededArgsForCall, struct {
	}{})
	stub := fake.ExecuteSetupScriptSucceededStub
	fake.recordInvocation("ExecuteSetupScriptSucceeded", []interface{}{})
	fake.executeSetupScriptSucceededMutex.Unlock()
	if stub != nil {
		fake.ExecuteSetupScriptSucceededStub()
	}
}

func (fake *FakeConstructMessenger) ExecuteSetupScriptSucceededCallCount() int {
	fake.executeSetupScriptSucceededMutex.RLock()
	defer fake.executeSetupScriptSucceededMutex.RUnlock()
	return len(fake.executeSetupScriptSucceededArgsForCall)
}

func (fake *FakeConstructMessenger) ExecuteSetupScriptSucceededCalls(stub func()) {
	fake.executeSetupScriptSucceededMutex.Lock()
	defer fake.executeSetupScriptSucceededMutex.Unlock()
	fake.ExecuteSetupScriptSucceededStub = stub
}

func (fake *FakeConstructMessenger) ExtractArtifactsStarted() {
	fake.extractArtifactsStartedMutex.Lock()
	fake.extractArtifactsStartedArgsForCall = append(fake.extractArtifactsStartedArgsForCall, struct {
	}{})
	stub := fake.ExtractArtifactsStartedStub
	fake.recordInvocation("ExtractArtifactsStarted", []interface{}{})
	fake.extractArtifactsStartedMutex.Unlock()
	if stub != nil {
		fake.ExtractArtifactsStartedStub()
	}
}

func (fake *FakeConstructMessenger) ExtractArtifactsStartedCallCount() int {
	fake.extractArtifactsStartedMutex.RLock()
	defer fake.extractArtifactsStartedMutex.RUnlock()
	return len(fake.extractArtifactsStartedArgsForCall)
}

func (fake *FakeConstructMessenger) ExtractArtifactsStartedCalls(stub func()) {
	fake.extractArtifactsStartedMutex.Lock()
	defer fake.extractArtifactsStartedMutex.Unlock()
	fake.ExtractArtifactsStartedStub = stub
}

func (fake *FakeConstructMessenger) ExtractArtifactsSucceeded() {
	fake.extractArtifactsSucceededMutex.Lock()
	fake.extractArtifactsSucceededArgsForCall = append(fake.extractArtifactsSucceededArgsForCall, struct {
	}{})
	stub := fake.ExtractArtifactsSucceededStub
	fake.recordInvocation("ExtractArtifactsSucceeded", []interface{}{})
	fake.extractArtifactsSucceededMutex.Unlock()
	if stub != nil {
		fake.ExtractArtifactsSucceededStub()
	}
}

func (fake *FakeConstructMessenger) ExtractArtifactsSucceededCallCount() int {
	fake.extractArtifactsSucceededMutex.RLock()
	defer fake.extractArtifactsSucceededMutex.RUnlock()
	return len(fake.extractArtifactsSucceededArgsForCall)
}

func (fake *FakeConstructMessenger) ExtractArtifactsSucceededCalls(stub func()) {
	fake.extractArtifactsSucceededMutex.Lock()
	defer fake.extractArtifactsSucceededMutex.Unlock()
	fake.ExtractArtifactsSucceededStub = stub
}

func (fake *FakeConstructMessenger) LogOutUsersStarted() {
	fake.logOutUsersStartedMutex.Lock()
	fake.logOutUsersStartedArgsForCall = append(fake.logOutUsersStartedArgsForCall, struct {
	}{})
	stub := fake.LogOutUsersStartedStub
	fake.recordInvocation("LogOutUsersStarted", []interface{}{})
	fake.logOutUsersStartedMutex.Unlock()
	if stub != nil {
		fake.LogOutUsersStartedStub()
	}
}

func (fake *FakeConstructMessenger) LogOutUsersStartedCallCount() int {
	fake.logOutUsersStartedMutex.RLock()
	defer fake.logOutUsersStartedMutex.RUnlock()
	return len(fake.logOutUsersStartedArgsForCall)
}

func (fake *FakeConstructMessenger) LogOutUsersStartedCalls(stub func()) {
	fake.logOutUsersStartedMutex.Lock()
	defer fake.logOutUsersStartedMutex.Unlock()
	fake.LogOutUsersStartedStub = stub
}

func (fake *FakeConstructMessenger) LogOutUsersSucceeded() {
	fake.logOutUsersSucceededMutex.Lock()
	fake.logOutUsersSucceededArgsForCall = append(fake.logOutUsersSucceededArgsForCall, struct {
	}{})
	stub := fake.LogOutUsersSucceededStub
	fake.recordInvocation("LogOutUsersSucceeded", []interface{}{})
	fake.logOutUsersSucceededMutex.Unlock()
	if stub != nil {
		fake.LogOutUsersSucceededStub()
	}
}

func (fake *FakeConstructMessenger) LogOutUsersSucceededCallCount() int {
	fake.logOutUsersSucceededMutex.RLock()
	defer fake.logOutUsersSucceededMutex.RUnlock()
	return len(fake.logOutUsersSucceededArgsForCall)
}

func (fake *FakeConstructMessenger) LogOutUsersSucceededCalls(stub func()) {
	fake.logOutUsersSucceededMutex.Lock()
	defer fake.logOutUsersSucceededMutex.Unlock()
	fake.LogOutUsersSucceededStub = stub
}

func (fake *FakeConstructMessenger) RebootHasFinished() {
	fake.rebootHasFinishedMutex.Lock()
	fake.rebootHasFinishedArgsForCall = append(fake.rebootHasFinishedArgsForCall, struct {
	}{})
	stub := fake.RebootHasFinishedStub
	fake.recordInvocation("RebootHasFinished", []interface{}{})
	fake.rebootHasFinishedMutex.Unlock()
	if stub != nil {
		fake.RebootHasFinishedStub()
	}
}

func (fake *FakeConstructMessenger) RebootHasFinishedCallCount() int {
	fake.rebootHasFinishedMutex.RLock()
	defer fake.rebootHasFinishedMutex.RUnlock()
	return len(fake.rebootHasFinishedArgsForCall)
}

func (fake *FakeConstructMessenger) RebootHasFinishedCalls(stub func()) {
	fake.rebootHasFinishedMutex.Lock()
	defer fake.rebootHasFinishedMutex.Unlock()
	fake.RebootHasFinishedStub = stub
}

func (fake *FakeConstructMessenger) RebootHasStarted() {
	fake.rebootHasStartedMutex.Lock()
	fake.rebootHasStartedArgsForCall = append(fake.rebootHasStartedArgsForCall, struct {
	}{})
	stub := fake.RebootHasStartedStub
	fake.recordInvocation("RebootHasStarted", []interface{}{})
	fake.rebootHasStartedMutex.Unlock()
	if stub != nil {
		fake.RebootHasStartedStub()
	}
}

func (fake *FakeConstructMessenger) RebootHasStartedCallCount() int {
	fake.rebootHasStartedMutex.RLock()
	defer fake.rebootHasStartedMutex.RUnlock()
	return len(fake.rebootHasStartedArgsForCall)
}

func (fake *FakeConstructMessenger) RebootHasStartedCalls(stub func()) {
	fake.rebootHasStartedMutex.Lock()
	defer fake.rebootHasStartedMutex.Unlock()
	fake.RebootHasStartedStub = stub
}

func (fake *FakeConstructMessenger) ShutdownCompleted() {
	fake.shutdownCompletedMutex.Lock()
	fake.shutdownCompletedArgsForCall = append(fake.shutdownCompletedArgsForCall, struct {
	}{})
	stub := fake.ShutdownCompletedStub
	fake.recordInvocation("ShutdownCompleted", []interface{}{})
	fake.shutdownCompletedMutex.Unlock()
	if stub != nil {
		fake.ShutdownCompletedStub()
	}
}

func (fake *FakeConstructMessenger) ShutdownCompletedCallCount() int {
	fake.shutdownCompletedMutex.RLock()
	defer fake.shutdownCompletedMutex.RUnlock()
	return len(fake.shutdownCompletedArgsForCall)
}

func (fake *FakeConstructMessenger) ShutdownCompletedCalls(stub func()) {
	fake.shutdownCompletedMutex.Lock()
	defer fake.shutdownCompletedMutex.Unlock()
	fake.ShutdownCompletedStub = stub
}

func (fake *FakeConstructMessenger) UploadArtifactsStarted() {
	fake.uploadArtifactsStartedMutex.Lock()
	fake.uploadArtifactsStartedArgsForCall = append(fake.uploadArtifactsStartedArgsForCall, struct {
	}{})
	stub := fake.UploadArtifactsStartedStub
	fake.recordInvocation("UploadArtifactsStarted", []interface{}{})
	fake.uploadArtifactsStartedMutex.Unlock()
	if stub != nil {
		fake.UploadArtifactsStartedStub()
	}
}

func (fake *FakeConstructMessenger) UploadArtifactsStartedCallCount() int {
	fake.uploadArtifactsStartedMutex.RLock()
	defer fake.uploadArtifactsStartedMutex.RUnlock()
	return len(fake.uploadArtifactsStartedArgsForCall)
}

func (fake *FakeConstructMessenger) UploadArtifactsStartedCalls(stub func()) {
	fake.uploadArtifactsStartedMutex.Lock()
	defer fake.uploadArtifactsStartedMutex.Unlock()
	fake.UploadArtifactsStartedStub = stub
}

func (fake *FakeConstructMessenger) UploadArtifactsSucceeded() {
	fake.uploadArtifactsSucceededMutex.Lock()
	fake.uploadArtifactsSucceededArgsForCall = append(fake.uploadArtifactsSucceededArgsForCall, struct {
	}{})
	stub := fake.UploadArtifactsSucceededStub
	fake.recordInvocation("UploadArtifactsSucceeded", []interface{}{})
	fake.uploadArtifactsSucceededMutex.Unlock()
	if stub != nil {
		fake.UploadArtifactsSucceededStub()
	}
}

func (fake *FakeConstructMessenger) UploadArtifactsSucceededCallCount() int {
	fake.uploadArtifactsSucceededMutex.RLock()
	defer fake.uploadArtifactsSucceededMutex.RUnlock()
	return len(fake.uploadArtifactsSucceededArgsForCall)
}

func (fake *FakeConstructMessenger) UploadArtifactsSucceededCalls(stub func()) {
	fake.uploadArtifactsSucceededMutex.Lock()
	defer fake.uploadArtifactsSucceededMutex.Unlock()
	fake.UploadArtifactsSucceededStub = stub
}

func (fake *FakeConstructMessenger) UploadFileStarted(arg1 string) {
	fake.uploadFileStartedMutex.Lock()
	fake.uploadFileStartedArgsForCall = append(fake.uploadFileStartedArgsForCall, struct {
		arg1 string
	}{arg1})
	stub := fake.UploadFileStartedStub
	fake.recordInvocation("UploadFileStarted", []interface{}{arg1})
	fake.uploadFileStartedMutex.Unlock()
	if stub != nil {
		fake.UploadFileStartedStub(arg1)
	}
}

func (fake *FakeConstructMessenger) UploadFileStartedCallCount() int {
	fake.uploadFileStartedMutex.RLock()
	defer fake.uploadFileStartedMutex.RUnlock()
	return len(fake.uploadFileStartedArgsForCall)
}

func (fake *FakeConstructMessenger) UploadFileStartedCalls(stub func(string)) {
	fake.uploadFileStartedMutex.Lock()
	defer fake.uploadFileStartedMutex.Unlock()
	fake.UploadFileStartedStub = stub
}

func (fake *FakeConstructMessenger) UploadFileStartedArgsForCall(i int) string {
	fake.uploadFileStartedMutex.RLock()
	defer fake.uploadFileStartedMutex.RUnlock()
	argsForCall := fake.uploadFileStartedArgsForCall[i]
	return argsForCall.arg1
}

func (fake *FakeConstructMessenger) UploadFileSucceeded() {
	fake.uploadFileSucceededMutex.Lock()
	fake.uploadFileSucceededArgsForCall = append(fake.uploadFileSucceededArgsForCall, struct {
	}{})
	stub := fake.UploadFileSucceededStub
	fake.recordInvocation("UploadFileSucceeded", []interface{}{})
	fake.uploadFileSucceededMutex.Unlock()
	if stub != nil {
		fake.UploadFileSucceededStub()
	}
}

func (fake *FakeConstructMessenger) UploadFileSucceededCallCount() int {
	fake.uploadFileSucceededMutex.RLock()
	defer fake.uploadFileSucceededMutex.RUnlock()
	return len(fake.uploadFileSucceededArgsForCall)
}

func (fake *FakeConstructMessenger) UploadFileSucceededCalls(stub func()) {
	fake.uploadFileSucceededMutex.Lock()
	defer fake.uploadFileSucceededMutex.Unlock()
	fake.UploadFileSucceededStub = stub
}

func (fake *FakeConstructMessenger) ValidateVMConnectionStarted() {
	fake.validateVMConnectionStartedMutex.Lock()
	fake.validateVMConnectionStartedArgsForCall = append(fake.validateVMConnectionStartedArgsForCall, struct {
	}{})
	stub := fake.ValidateVMConnectionStartedStub
	fake.recordInvocation("ValidateVMConnectionStarted", []interface{}{})
	fake.validateVMConnectionStartedMutex.Unlock()
	if stub != nil {
		fake.ValidateVMConnectionStartedStub()
	}
}

func (fake *FakeConstructMessenger) ValidateVMConnectionStartedCallCount() int {
	fake.validateVMConnectionStartedMutex.RLock()
	defer fake.validateVMConnectionStartedMutex.RUnlock()
	return len(fake.validateVMConnectionStartedArgsForCall)
}

func (fake *FakeConstructMessenger) ValidateVMConnectionStartedCalls(stub func()) {
	fake.validateVMConnectionStartedMutex.Lock()
	defer fake.validateVMConnectionStartedMutex.Unlock()
	fake.ValidateVMConnectionStartedStub = stub
}

func (fake *FakeConstructMessenger) ValidateVMConnectionSucceeded() {
	fake.validateVMConnectionSucceededMutex.Lock()
	fake.validateVMConnectionSucceededArgsForCall = append(fake.validateVMConnectionSucceededArgsForCall, struct {
	}{})
	stub := fake.ValidateVMConnectionSucceededStub
	fake.recordInvocation("ValidateVMConnectionSucceeded", []interface{}{})
	fake.validateVMConnectionSucceededMutex.Unlock()
	if stub != nil {
		fake.ValidateVMConnectionSucceededStub()
	}
}

func (fake *FakeConstructMessenger) ValidateVMConnectionSucceededCallCount() int {
	fake.validateVMConnectionSucceededMutex.RLock()
	defer fake.validateVMConnectionSucceededMutex.RUnlock()
	return len(fake.validateVMConnectionSucceededArgsForCall)
}

func (fake *FakeConstructMessenger) ValidateVMConnectionSucceededCalls(stub func()) {
	fake.validateVMConnectionSucceededMutex.Lock()
	defer fake.validateVMConnectionSucceededMutex.Unlock()
	fake.ValidateVMConnectionSucceededStub = stub
}

func (fake *FakeConstructMessenger) WaitingForShutdown() {
	fake.waitingForShutdownMutex.Lock()
	fake.waitingForShutdownArgsForCall = append(fake.waitingForShutdownArgsForCall, struct {
	}{})
	stub := fake.WaitingForShutdownStub
	fake.recordInvocation("WaitingForShutdown", []interface{}{})
	fake.waitingForShutdownMutex.Unlock()
	if stub != nil {
		fake.WaitingForShutdownStub()
	}
}

func (fake *FakeConstructMessenger) WaitingForShutdownCallCount() int {
	fake.waitingForShutdownMutex.RLock()
	defer fake.waitingForShutdownMutex.RUnlock()
	return len(fake.waitingForShutdownArgsForCall)
}

func (fake *FakeConstructMessenger) WaitingForShutdownCalls(stub func()) {
	fake.waitingForShutdownMutex.Lock()
	defer fake.waitingForShutdownMutex.Unlock()
	fake.WaitingForShutdownStub = stub
}

func (fake *FakeConstructMessenger) WinRMDisconnectedForReboot() {
	fake.winRMDisconnectedForRebootMutex.Lock()
	fake.winRMDisconnectedForRebootArgsForCall = append(fake.winRMDisconnectedForRebootArgsForCall, struct {
	}{})
	stub := fake.WinRMDisconnectedForRebootStub
	fake.recordInvocation("WinRMDisconnectedForReboot", []interface{}{})
	fake.winRMDisconnectedForRebootMutex.Unlock()
	if stub != nil {
		fake.WinRMDisconnectedForRebootStub()
	}
}

func (fake *FakeConstructMessenger) WinRMDisconnectedForRebootCallCount() int {
	fake.winRMDisconnectedForRebootMutex.RLock()
	defer fake.winRMDisconnectedForRebootMutex.RUnlock()
	return len(fake.winRMDisconnectedForRebootArgsForCall)
}

func (fake *FakeConstructMessenger) WinRMDisconnectedForRebootCalls(stub func()) {
	fake.winRMDisconnectedForRebootMutex.Lock()
	defer fake.winRMDisconnectedForRebootMutex.Unlock()
	fake.WinRMDisconnectedForRebootStub = stub
}

func (fake *FakeConstructMessenger) Invocations() map[string][][]interface{} {
	fake.invocationsMutex.RLock()
	defer fake.invocationsMutex.RUnlock()
	fake.createProvisionDirStartedMutex.RLock()
	defer fake.createProvisionDirStartedMutex.RUnlock()
	fake.createProvisionDirSucceededMutex.RLock()
	defer fake.createProvisionDirSucceededMutex.RUnlock()
	fake.enableWinRMStartedMutex.RLock()
	defer fake.enableWinRMStartedMutex.RUnlock()
	fake.enableWinRMSucceededMutex.RLock()
	defer fake.enableWinRMSucceededMutex.RUnlock()
	fake.executePostRebootScriptStartedMutex.RLock()
	defer fake.executePostRebootScriptStartedMutex.RUnlock()
	fake.executePostRebootScriptSucceededMutex.RLock()
	defer fake.executePostRebootScriptSucceededMutex.RUnlock()
	fake.executePostRebootWarningMutex.RLock()
	defer fake.executePostRebootWarningMutex.RUnlock()
	fake.executeSetupScriptStartedMutex.RLock()
	defer fake.executeSetupScriptStartedMutex.RUnlock()
	fake.executeSetupScriptSucceededMutex.RLock()
	defer fake.executeSetupScriptSucceededMutex.RUnlock()
	fake.extractArtifactsStartedMutex.RLock()
	defer fake.extractArtifactsStartedMutex.RUnlock()
	fake.extractArtifactsSucceededMutex.RLock()
	defer fake.extractArtifactsSucceededMutex.RUnlock()
	fake.logOutUsersStartedMutex.RLock()
	defer fake.logOutUsersStartedMutex.RUnlock()
	fake.logOutUsersSucceededMutex.RLock()
	defer fake.logOutUsersSucceededMutex.RUnlock()
	fake.rebootHasFinishedMutex.RLock()
	defer fake.rebootHasFinishedMutex.RUnlock()
	fake.rebootHasStartedMutex.RLock()
	defer fake.rebootHasStartedMutex.RUnlock()
	fake.shutdownCompletedMutex.RLock()
	defer fake.shutdownCompletedMutex.RUnlock()
	fake.uploadArtifactsStartedMutex.RLock()
	defer fake.uploadArtifactsStartedMutex.RUnlock()
	fake.uploadArtifactsSucceededMutex.RLock()
	defer fake.uploadArtifactsSucceededMutex.RUnlock()
	fake.uploadFileStartedMutex.RLock()
	defer fake.uploadFileStartedMutex.RUnlock()
	fake.uploadFileSucceededMutex.RLock()
	defer fake.uploadFileSucceededMutex.RUnlock()
	fake.validateVMConnectionStartedMutex.RLock()
	defer fake.validateVMConnectionStartedMutex.RUnlock()
	fake.validateVMConnectionSucceededMutex.RLock()
	defer fake.validateVMConnectionSucceededMutex.RUnlock()
	fake.waitingForShutdownMutex.RLock()
	defer fake.waitingForShutdownMutex.RUnlock()
	fake.winRMDisconnectedForRebootMutex.RLock()
	defer fake.winRMDisconnectedForRebootMutex.RUnlock()
	copiedInvocations := map[string][][]interface{}{}
	for key, value := range fake.invocations {
		copiedInvocations[key] = value
	}
	return copiedInvocations
}

func (fake *FakeConstructMessenger) recordInvocation(key string, args []interface{}) {
	fake.invocationsMutex.Lock()
	defer fake.invocationsMutex.Unlock()
	if fake.invocations == nil {
		fake.invocations = map[string][][]interface{}{}
	}
	if fake.invocations[key] == nil {
		fake.invocations[key] = [][]interface{}{}
	}
	fake.invocations[key] = append(fake.invocations[key], args)
}

var _ construct.ConstructMessenger = new(FakeConstructMessenger)
