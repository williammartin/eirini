// Code generated by counterfeiter. DO NOT EDIT.
package uaafakes

import (
	http "net/http"
	sync "sync"

	uaa "code.cloudfoundry.org/cli/api/uaa"
)

type FakeConnectionWrapper struct {
	MakeStub        func(*http.Request, *uaa.Response) error
	makeMutex       sync.RWMutex
	makeArgsForCall []struct {
		arg1 *http.Request
		arg2 *uaa.Response
	}
	makeReturns struct {
		result1 error
	}
	makeReturnsOnCall map[int]struct {
		result1 error
	}
	WrapStub        func(uaa.Connection) uaa.Connection
	wrapMutex       sync.RWMutex
	wrapArgsForCall []struct {
		arg1 uaa.Connection
	}
	wrapReturns struct {
		result1 uaa.Connection
	}
	wrapReturnsOnCall map[int]struct {
		result1 uaa.Connection
	}
	invocations      map[string][][]interface{}
	invocationsMutex sync.RWMutex
}

func (fake *FakeConnectionWrapper) Make(arg1 *http.Request, arg2 *uaa.Response) error {
	fake.makeMutex.Lock()
	ret, specificReturn := fake.makeReturnsOnCall[len(fake.makeArgsForCall)]
	fake.makeArgsForCall = append(fake.makeArgsForCall, struct {
		arg1 *http.Request
		arg2 *uaa.Response
	}{arg1, arg2})
	fake.recordInvocation("Make", []interface{}{arg1, arg2})
	fake.makeMutex.Unlock()
	if fake.MakeStub != nil {
		return fake.MakeStub(arg1, arg2)
	}
	if specificReturn {
		return ret.result1
	}
	fakeReturns := fake.makeReturns
	return fakeReturns.result1
}

func (fake *FakeConnectionWrapper) MakeCallCount() int {
	fake.makeMutex.RLock()
	defer fake.makeMutex.RUnlock()
	return len(fake.makeArgsForCall)
}

func (fake *FakeConnectionWrapper) MakeCalls(stub func(*http.Request, *uaa.Response) error) {
	fake.makeMutex.Lock()
	defer fake.makeMutex.Unlock()
	fake.MakeStub = stub
}

func (fake *FakeConnectionWrapper) MakeArgsForCall(i int) (*http.Request, *uaa.Response) {
	fake.makeMutex.RLock()
	defer fake.makeMutex.RUnlock()
	argsForCall := fake.makeArgsForCall[i]
	return argsForCall.arg1, argsForCall.arg2
}

func (fake *FakeConnectionWrapper) MakeReturns(result1 error) {
	fake.makeMutex.Lock()
	defer fake.makeMutex.Unlock()
	fake.MakeStub = nil
	fake.makeReturns = struct {
		result1 error
	}{result1}
}

func (fake *FakeConnectionWrapper) MakeReturnsOnCall(i int, result1 error) {
	fake.makeMutex.Lock()
	defer fake.makeMutex.Unlock()
	fake.MakeStub = nil
	if fake.makeReturnsOnCall == nil {
		fake.makeReturnsOnCall = make(map[int]struct {
			result1 error
		})
	}
	fake.makeReturnsOnCall[i] = struct {
		result1 error
	}{result1}
}

func (fake *FakeConnectionWrapper) Wrap(arg1 uaa.Connection) uaa.Connection {
	fake.wrapMutex.Lock()
	ret, specificReturn := fake.wrapReturnsOnCall[len(fake.wrapArgsForCall)]
	fake.wrapArgsForCall = append(fake.wrapArgsForCall, struct {
		arg1 uaa.Connection
	}{arg1})
	fake.recordInvocation("Wrap", []interface{}{arg1})
	fake.wrapMutex.Unlock()
	if fake.WrapStub != nil {
		return fake.WrapStub(arg1)
	}
	if specificReturn {
		return ret.result1
	}
	fakeReturns := fake.wrapReturns
	return fakeReturns.result1
}

func (fake *FakeConnectionWrapper) WrapCallCount() int {
	fake.wrapMutex.RLock()
	defer fake.wrapMutex.RUnlock()
	return len(fake.wrapArgsForCall)
}

func (fake *FakeConnectionWrapper) WrapCalls(stub func(uaa.Connection) uaa.Connection) {
	fake.wrapMutex.Lock()
	defer fake.wrapMutex.Unlock()
	fake.WrapStub = stub
}

func (fake *FakeConnectionWrapper) WrapArgsForCall(i int) uaa.Connection {
	fake.wrapMutex.RLock()
	defer fake.wrapMutex.RUnlock()
	argsForCall := fake.wrapArgsForCall[i]
	return argsForCall.arg1
}

func (fake *FakeConnectionWrapper) WrapReturns(result1 uaa.Connection) {
	fake.wrapMutex.Lock()
	defer fake.wrapMutex.Unlock()
	fake.WrapStub = nil
	fake.wrapReturns = struct {
		result1 uaa.Connection
	}{result1}
}

func (fake *FakeConnectionWrapper) WrapReturnsOnCall(i int, result1 uaa.Connection) {
	fake.wrapMutex.Lock()
	defer fake.wrapMutex.Unlock()
	fake.WrapStub = nil
	if fake.wrapReturnsOnCall == nil {
		fake.wrapReturnsOnCall = make(map[int]struct {
			result1 uaa.Connection
		})
	}
	fake.wrapReturnsOnCall[i] = struct {
		result1 uaa.Connection
	}{result1}
}

func (fake *FakeConnectionWrapper) Invocations() map[string][][]interface{} {
	fake.invocationsMutex.RLock()
	defer fake.invocationsMutex.RUnlock()
	fake.makeMutex.RLock()
	defer fake.makeMutex.RUnlock()
	fake.wrapMutex.RLock()
	defer fake.wrapMutex.RUnlock()
	copiedInvocations := map[string][][]interface{}{}
	for key, value := range fake.invocations {
		copiedInvocations[key] = value
	}
	return copiedInvocations
}

func (fake *FakeConnectionWrapper) recordInvocation(key string, args []interface{}) {
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

var _ uaa.ConnectionWrapper = new(FakeConnectionWrapper)
