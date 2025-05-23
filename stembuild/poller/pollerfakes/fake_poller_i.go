// Code generated by counterfeiter. DO NOT EDIT.
package pollerfakes

import (
	"sync"
	"time"

	"github.com/cloudfoundry/bosh-windows-stemcell-builder/stembuild/poller"
)

type FakePollerI struct {
	PollStub        func(time.Duration, func() (bool, error)) error
	pollMutex       sync.RWMutex
	pollArgsForCall []struct {
		arg1 time.Duration
		arg2 func() (bool, error)
	}
	pollReturns struct {
		result1 error
	}
	pollReturnsOnCall map[int]struct {
		result1 error
	}
	invocations      map[string][][]interface{}
	invocationsMutex sync.RWMutex
}

func (fake *FakePollerI) Poll(arg1 time.Duration, arg2 func() (bool, error)) error {
	fake.pollMutex.Lock()
	ret, specificReturn := fake.pollReturnsOnCall[len(fake.pollArgsForCall)]
	fake.pollArgsForCall = append(fake.pollArgsForCall, struct {
		arg1 time.Duration
		arg2 func() (bool, error)
	}{arg1, arg2})
	stub := fake.PollStub
	fakeReturns := fake.pollReturns
	fake.recordInvocation("Poll", []interface{}{arg1, arg2})
	fake.pollMutex.Unlock()
	if stub != nil {
		return stub(arg1, arg2)
	}
	if specificReturn {
		return ret.result1
	}
	return fakeReturns.result1
}

func (fake *FakePollerI) PollCallCount() int {
	fake.pollMutex.RLock()
	defer fake.pollMutex.RUnlock()
	return len(fake.pollArgsForCall)
}

func (fake *FakePollerI) PollCalls(stub func(time.Duration, func() (bool, error)) error) {
	fake.pollMutex.Lock()
	defer fake.pollMutex.Unlock()
	fake.PollStub = stub
}

func (fake *FakePollerI) PollArgsForCall(i int) (time.Duration, func() (bool, error)) {
	fake.pollMutex.RLock()
	defer fake.pollMutex.RUnlock()
	argsForCall := fake.pollArgsForCall[i]
	return argsForCall.arg1, argsForCall.arg2
}

func (fake *FakePollerI) PollReturns(result1 error) {
	fake.pollMutex.Lock()
	defer fake.pollMutex.Unlock()
	fake.PollStub = nil
	fake.pollReturns = struct {
		result1 error
	}{result1}
}

func (fake *FakePollerI) PollReturnsOnCall(i int, result1 error) {
	fake.pollMutex.Lock()
	defer fake.pollMutex.Unlock()
	fake.PollStub = nil
	if fake.pollReturnsOnCall == nil {
		fake.pollReturnsOnCall = make(map[int]struct {
			result1 error
		})
	}
	fake.pollReturnsOnCall[i] = struct {
		result1 error
	}{result1}
}

func (fake *FakePollerI) Invocations() map[string][][]interface{} {
	fake.invocationsMutex.RLock()
	defer fake.invocationsMutex.RUnlock()
	fake.pollMutex.RLock()
	defer fake.pollMutex.RUnlock()
	copiedInvocations := map[string][][]interface{}{}
	for key, value := range fake.invocations {
		copiedInvocations[key] = value
	}
	return copiedInvocations
}

func (fake *FakePollerI) recordInvocation(key string, args []interface{}) {
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

var _ poller.PollerI = new(FakePollerI)
