// Code generated by counterfeiter. DO NOT EDIT.
package packagerfakes

import (
	"sync"

	"github.com/cloudfoundry/bosh-windows-stemcell-builder/stembuild/package_stemcell/packager"
)

type FakeIaasClient struct {
	EjectCDRomStub        func(string, string) error
	ejectCDRomMutex       sync.RWMutex
	ejectCDRomArgsForCall []struct {
		arg1 string
		arg2 string
	}
	ejectCDRomReturns struct {
		result1 error
	}
	ejectCDRomReturnsOnCall map[int]struct {
		result1 error
	}
	ExportVMStub        func(string, string) error
	exportVMMutex       sync.RWMutex
	exportVMArgsForCall []struct {
		arg1 string
		arg2 string
	}
	exportVMReturns struct {
		result1 error
	}
	exportVMReturnsOnCall map[int]struct {
		result1 error
	}
	FindVMStub        func(string) error
	findVMMutex       sync.RWMutex
	findVMArgsForCall []struct {
		arg1 string
	}
	findVMReturns struct {
		result1 error
	}
	findVMReturnsOnCall map[int]struct {
		result1 error
	}
	ListDevicesStub        func(string) ([]string, error)
	listDevicesMutex       sync.RWMutex
	listDevicesArgsForCall []struct {
		arg1 string
	}
	listDevicesReturns struct {
		result1 []string
		result2 error
	}
	listDevicesReturnsOnCall map[int]struct {
		result1 []string
		result2 error
	}
	RemoveDeviceStub        func(string, string) error
	removeDeviceMutex       sync.RWMutex
	removeDeviceArgsForCall []struct {
		arg1 string
		arg2 string
	}
	removeDeviceReturns struct {
		result1 error
	}
	removeDeviceReturnsOnCall map[int]struct {
		result1 error
	}
	ValidateCredentialsStub        func() error
	validateCredentialsMutex       sync.RWMutex
	validateCredentialsArgsForCall []struct {
	}
	validateCredentialsReturns struct {
		result1 error
	}
	validateCredentialsReturnsOnCall map[int]struct {
		result1 error
	}
	ValidateUrlStub        func() error
	validateUrlMutex       sync.RWMutex
	validateUrlArgsForCall []struct {
	}
	validateUrlReturns struct {
		result1 error
	}
	validateUrlReturnsOnCall map[int]struct {
		result1 error
	}
	invocations      map[string][][]interface{}
	invocationsMutex sync.RWMutex
}

func (fake *FakeIaasClient) EjectCDRom(arg1 string, arg2 string) error {
	fake.ejectCDRomMutex.Lock()
	ret, specificReturn := fake.ejectCDRomReturnsOnCall[len(fake.ejectCDRomArgsForCall)]
	fake.ejectCDRomArgsForCall = append(fake.ejectCDRomArgsForCall, struct {
		arg1 string
		arg2 string
	}{arg1, arg2})
	stub := fake.EjectCDRomStub
	fakeReturns := fake.ejectCDRomReturns
	fake.recordInvocation("EjectCDRom", []interface{}{arg1, arg2})
	fake.ejectCDRomMutex.Unlock()
	if stub != nil {
		return stub(arg1, arg2)
	}
	if specificReturn {
		return ret.result1
	}
	return fakeReturns.result1
}

func (fake *FakeIaasClient) EjectCDRomCallCount() int {
	fake.ejectCDRomMutex.RLock()
	defer fake.ejectCDRomMutex.RUnlock()
	return len(fake.ejectCDRomArgsForCall)
}

func (fake *FakeIaasClient) EjectCDRomCalls(stub func(string, string) error) {
	fake.ejectCDRomMutex.Lock()
	defer fake.ejectCDRomMutex.Unlock()
	fake.EjectCDRomStub = stub
}

func (fake *FakeIaasClient) EjectCDRomArgsForCall(i int) (string, string) {
	fake.ejectCDRomMutex.RLock()
	defer fake.ejectCDRomMutex.RUnlock()
	argsForCall := fake.ejectCDRomArgsForCall[i]
	return argsForCall.arg1, argsForCall.arg2
}

func (fake *FakeIaasClient) EjectCDRomReturns(result1 error) {
	fake.ejectCDRomMutex.Lock()
	defer fake.ejectCDRomMutex.Unlock()
	fake.EjectCDRomStub = nil
	fake.ejectCDRomReturns = struct {
		result1 error
	}{result1}
}

func (fake *FakeIaasClient) EjectCDRomReturnsOnCall(i int, result1 error) {
	fake.ejectCDRomMutex.Lock()
	defer fake.ejectCDRomMutex.Unlock()
	fake.EjectCDRomStub = nil
	if fake.ejectCDRomReturnsOnCall == nil {
		fake.ejectCDRomReturnsOnCall = make(map[int]struct {
			result1 error
		})
	}
	fake.ejectCDRomReturnsOnCall[i] = struct {
		result1 error
	}{result1}
}

func (fake *FakeIaasClient) ExportVM(arg1 string, arg2 string) error {
	fake.exportVMMutex.Lock()
	ret, specificReturn := fake.exportVMReturnsOnCall[len(fake.exportVMArgsForCall)]
	fake.exportVMArgsForCall = append(fake.exportVMArgsForCall, struct {
		arg1 string
		arg2 string
	}{arg1, arg2})
	stub := fake.ExportVMStub
	fakeReturns := fake.exportVMReturns
	fake.recordInvocation("ExportVM", []interface{}{arg1, arg2})
	fake.exportVMMutex.Unlock()
	if stub != nil {
		return stub(arg1, arg2)
	}
	if specificReturn {
		return ret.result1
	}
	return fakeReturns.result1
}

func (fake *FakeIaasClient) ExportVMCallCount() int {
	fake.exportVMMutex.RLock()
	defer fake.exportVMMutex.RUnlock()
	return len(fake.exportVMArgsForCall)
}

func (fake *FakeIaasClient) ExportVMCalls(stub func(string, string) error) {
	fake.exportVMMutex.Lock()
	defer fake.exportVMMutex.Unlock()
	fake.ExportVMStub = stub
}

func (fake *FakeIaasClient) ExportVMArgsForCall(i int) (string, string) {
	fake.exportVMMutex.RLock()
	defer fake.exportVMMutex.RUnlock()
	argsForCall := fake.exportVMArgsForCall[i]
	return argsForCall.arg1, argsForCall.arg2
}

func (fake *FakeIaasClient) ExportVMReturns(result1 error) {
	fake.exportVMMutex.Lock()
	defer fake.exportVMMutex.Unlock()
	fake.ExportVMStub = nil
	fake.exportVMReturns = struct {
		result1 error
	}{result1}
}

func (fake *FakeIaasClient) ExportVMReturnsOnCall(i int, result1 error) {
	fake.exportVMMutex.Lock()
	defer fake.exportVMMutex.Unlock()
	fake.ExportVMStub = nil
	if fake.exportVMReturnsOnCall == nil {
		fake.exportVMReturnsOnCall = make(map[int]struct {
			result1 error
		})
	}
	fake.exportVMReturnsOnCall[i] = struct {
		result1 error
	}{result1}
}

func (fake *FakeIaasClient) FindVM(arg1 string) error {
	fake.findVMMutex.Lock()
	ret, specificReturn := fake.findVMReturnsOnCall[len(fake.findVMArgsForCall)]
	fake.findVMArgsForCall = append(fake.findVMArgsForCall, struct {
		arg1 string
	}{arg1})
	stub := fake.FindVMStub
	fakeReturns := fake.findVMReturns
	fake.recordInvocation("FindVM", []interface{}{arg1})
	fake.findVMMutex.Unlock()
	if stub != nil {
		return stub(arg1)
	}
	if specificReturn {
		return ret.result1
	}
	return fakeReturns.result1
}

func (fake *FakeIaasClient) FindVMCallCount() int {
	fake.findVMMutex.RLock()
	defer fake.findVMMutex.RUnlock()
	return len(fake.findVMArgsForCall)
}

func (fake *FakeIaasClient) FindVMCalls(stub func(string) error) {
	fake.findVMMutex.Lock()
	defer fake.findVMMutex.Unlock()
	fake.FindVMStub = stub
}

func (fake *FakeIaasClient) FindVMArgsForCall(i int) string {
	fake.findVMMutex.RLock()
	defer fake.findVMMutex.RUnlock()
	argsForCall := fake.findVMArgsForCall[i]
	return argsForCall.arg1
}

func (fake *FakeIaasClient) FindVMReturns(result1 error) {
	fake.findVMMutex.Lock()
	defer fake.findVMMutex.Unlock()
	fake.FindVMStub = nil
	fake.findVMReturns = struct {
		result1 error
	}{result1}
}

func (fake *FakeIaasClient) FindVMReturnsOnCall(i int, result1 error) {
	fake.findVMMutex.Lock()
	defer fake.findVMMutex.Unlock()
	fake.FindVMStub = nil
	if fake.findVMReturnsOnCall == nil {
		fake.findVMReturnsOnCall = make(map[int]struct {
			result1 error
		})
	}
	fake.findVMReturnsOnCall[i] = struct {
		result1 error
	}{result1}
}

func (fake *FakeIaasClient) ListDevices(arg1 string) ([]string, error) {
	fake.listDevicesMutex.Lock()
	ret, specificReturn := fake.listDevicesReturnsOnCall[len(fake.listDevicesArgsForCall)]
	fake.listDevicesArgsForCall = append(fake.listDevicesArgsForCall, struct {
		arg1 string
	}{arg1})
	stub := fake.ListDevicesStub
	fakeReturns := fake.listDevicesReturns
	fake.recordInvocation("ListDevices", []interface{}{arg1})
	fake.listDevicesMutex.Unlock()
	if stub != nil {
		return stub(arg1)
	}
	if specificReturn {
		return ret.result1, ret.result2
	}
	return fakeReturns.result1, fakeReturns.result2
}

func (fake *FakeIaasClient) ListDevicesCallCount() int {
	fake.listDevicesMutex.RLock()
	defer fake.listDevicesMutex.RUnlock()
	return len(fake.listDevicesArgsForCall)
}

func (fake *FakeIaasClient) ListDevicesCalls(stub func(string) ([]string, error)) {
	fake.listDevicesMutex.Lock()
	defer fake.listDevicesMutex.Unlock()
	fake.ListDevicesStub = stub
}

func (fake *FakeIaasClient) ListDevicesArgsForCall(i int) string {
	fake.listDevicesMutex.RLock()
	defer fake.listDevicesMutex.RUnlock()
	argsForCall := fake.listDevicesArgsForCall[i]
	return argsForCall.arg1
}

func (fake *FakeIaasClient) ListDevicesReturns(result1 []string, result2 error) {
	fake.listDevicesMutex.Lock()
	defer fake.listDevicesMutex.Unlock()
	fake.ListDevicesStub = nil
	fake.listDevicesReturns = struct {
		result1 []string
		result2 error
	}{result1, result2}
}

func (fake *FakeIaasClient) ListDevicesReturnsOnCall(i int, result1 []string, result2 error) {
	fake.listDevicesMutex.Lock()
	defer fake.listDevicesMutex.Unlock()
	fake.ListDevicesStub = nil
	if fake.listDevicesReturnsOnCall == nil {
		fake.listDevicesReturnsOnCall = make(map[int]struct {
			result1 []string
			result2 error
		})
	}
	fake.listDevicesReturnsOnCall[i] = struct {
		result1 []string
		result2 error
	}{result1, result2}
}

func (fake *FakeIaasClient) RemoveDevice(arg1 string, arg2 string) error {
	fake.removeDeviceMutex.Lock()
	ret, specificReturn := fake.removeDeviceReturnsOnCall[len(fake.removeDeviceArgsForCall)]
	fake.removeDeviceArgsForCall = append(fake.removeDeviceArgsForCall, struct {
		arg1 string
		arg2 string
	}{arg1, arg2})
	stub := fake.RemoveDeviceStub
	fakeReturns := fake.removeDeviceReturns
	fake.recordInvocation("RemoveDevice", []interface{}{arg1, arg2})
	fake.removeDeviceMutex.Unlock()
	if stub != nil {
		return stub(arg1, arg2)
	}
	if specificReturn {
		return ret.result1
	}
	return fakeReturns.result1
}

func (fake *FakeIaasClient) RemoveDeviceCallCount() int {
	fake.removeDeviceMutex.RLock()
	defer fake.removeDeviceMutex.RUnlock()
	return len(fake.removeDeviceArgsForCall)
}

func (fake *FakeIaasClient) RemoveDeviceCalls(stub func(string, string) error) {
	fake.removeDeviceMutex.Lock()
	defer fake.removeDeviceMutex.Unlock()
	fake.RemoveDeviceStub = stub
}

func (fake *FakeIaasClient) RemoveDeviceArgsForCall(i int) (string, string) {
	fake.removeDeviceMutex.RLock()
	defer fake.removeDeviceMutex.RUnlock()
	argsForCall := fake.removeDeviceArgsForCall[i]
	return argsForCall.arg1, argsForCall.arg2
}

func (fake *FakeIaasClient) RemoveDeviceReturns(result1 error) {
	fake.removeDeviceMutex.Lock()
	defer fake.removeDeviceMutex.Unlock()
	fake.RemoveDeviceStub = nil
	fake.removeDeviceReturns = struct {
		result1 error
	}{result1}
}

func (fake *FakeIaasClient) RemoveDeviceReturnsOnCall(i int, result1 error) {
	fake.removeDeviceMutex.Lock()
	defer fake.removeDeviceMutex.Unlock()
	fake.RemoveDeviceStub = nil
	if fake.removeDeviceReturnsOnCall == nil {
		fake.removeDeviceReturnsOnCall = make(map[int]struct {
			result1 error
		})
	}
	fake.removeDeviceReturnsOnCall[i] = struct {
		result1 error
	}{result1}
}

func (fake *FakeIaasClient) ValidateCredentials() error {
	fake.validateCredentialsMutex.Lock()
	ret, specificReturn := fake.validateCredentialsReturnsOnCall[len(fake.validateCredentialsArgsForCall)]
	fake.validateCredentialsArgsForCall = append(fake.validateCredentialsArgsForCall, struct {
	}{})
	stub := fake.ValidateCredentialsStub
	fakeReturns := fake.validateCredentialsReturns
	fake.recordInvocation("ValidateCredentials", []interface{}{})
	fake.validateCredentialsMutex.Unlock()
	if stub != nil {
		return stub()
	}
	if specificReturn {
		return ret.result1
	}
	return fakeReturns.result1
}

func (fake *FakeIaasClient) ValidateCredentialsCallCount() int {
	fake.validateCredentialsMutex.RLock()
	defer fake.validateCredentialsMutex.RUnlock()
	return len(fake.validateCredentialsArgsForCall)
}

func (fake *FakeIaasClient) ValidateCredentialsCalls(stub func() error) {
	fake.validateCredentialsMutex.Lock()
	defer fake.validateCredentialsMutex.Unlock()
	fake.ValidateCredentialsStub = stub
}

func (fake *FakeIaasClient) ValidateCredentialsReturns(result1 error) {
	fake.validateCredentialsMutex.Lock()
	defer fake.validateCredentialsMutex.Unlock()
	fake.ValidateCredentialsStub = nil
	fake.validateCredentialsReturns = struct {
		result1 error
	}{result1}
}

func (fake *FakeIaasClient) ValidateCredentialsReturnsOnCall(i int, result1 error) {
	fake.validateCredentialsMutex.Lock()
	defer fake.validateCredentialsMutex.Unlock()
	fake.ValidateCredentialsStub = nil
	if fake.validateCredentialsReturnsOnCall == nil {
		fake.validateCredentialsReturnsOnCall = make(map[int]struct {
			result1 error
		})
	}
	fake.validateCredentialsReturnsOnCall[i] = struct {
		result1 error
	}{result1}
}

func (fake *FakeIaasClient) ValidateUrl() error {
	fake.validateUrlMutex.Lock()
	ret, specificReturn := fake.validateUrlReturnsOnCall[len(fake.validateUrlArgsForCall)]
	fake.validateUrlArgsForCall = append(fake.validateUrlArgsForCall, struct {
	}{})
	stub := fake.ValidateUrlStub
	fakeReturns := fake.validateUrlReturns
	fake.recordInvocation("ValidateUrl", []interface{}{})
	fake.validateUrlMutex.Unlock()
	if stub != nil {
		return stub()
	}
	if specificReturn {
		return ret.result1
	}
	return fakeReturns.result1
}

func (fake *FakeIaasClient) ValidateUrlCallCount() int {
	fake.validateUrlMutex.RLock()
	defer fake.validateUrlMutex.RUnlock()
	return len(fake.validateUrlArgsForCall)
}

func (fake *FakeIaasClient) ValidateUrlCalls(stub func() error) {
	fake.validateUrlMutex.Lock()
	defer fake.validateUrlMutex.Unlock()
	fake.ValidateUrlStub = stub
}

func (fake *FakeIaasClient) ValidateUrlReturns(result1 error) {
	fake.validateUrlMutex.Lock()
	defer fake.validateUrlMutex.Unlock()
	fake.ValidateUrlStub = nil
	fake.validateUrlReturns = struct {
		result1 error
	}{result1}
}

func (fake *FakeIaasClient) ValidateUrlReturnsOnCall(i int, result1 error) {
	fake.validateUrlMutex.Lock()
	defer fake.validateUrlMutex.Unlock()
	fake.ValidateUrlStub = nil
	if fake.validateUrlReturnsOnCall == nil {
		fake.validateUrlReturnsOnCall = make(map[int]struct {
			result1 error
		})
	}
	fake.validateUrlReturnsOnCall[i] = struct {
		result1 error
	}{result1}
}

func (fake *FakeIaasClient) Invocations() map[string][][]interface{} {
	fake.invocationsMutex.RLock()
	defer fake.invocationsMutex.RUnlock()
	fake.ejectCDRomMutex.RLock()
	defer fake.ejectCDRomMutex.RUnlock()
	fake.exportVMMutex.RLock()
	defer fake.exportVMMutex.RUnlock()
	fake.findVMMutex.RLock()
	defer fake.findVMMutex.RUnlock()
	fake.listDevicesMutex.RLock()
	defer fake.listDevicesMutex.RUnlock()
	fake.removeDeviceMutex.RLock()
	defer fake.removeDeviceMutex.RUnlock()
	fake.validateCredentialsMutex.RLock()
	defer fake.validateCredentialsMutex.RUnlock()
	fake.validateUrlMutex.RLock()
	defer fake.validateUrlMutex.RUnlock()
	copiedInvocations := map[string][][]interface{}{}
	for key, value := range fake.invocations {
		copiedInvocations[key] = value
	}
	return copiedInvocations
}

func (fake *FakeIaasClient) recordInvocation(key string, args []interface{}) {
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

var _ packager.IaasClient = new(FakeIaasClient)
