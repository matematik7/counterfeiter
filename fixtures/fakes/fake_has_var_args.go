// This file was generated by counterfeiter
package fakes

import (
	"sync"

	"github.com/matematik7/counterfeiter/fixtures"
)

type FakeHasVarArgs struct {
	DoThingsStub        func(int, ...string) int
	doThingsMutex       sync.RWMutex
	doThingsArgsForCall []struct {
		arg1 int
		arg2 []string
	}
	doThingsReturns struct {
		result1 int
	}
	DoMoreThingsStub        func(int, int, ...string) int
	doMoreThingsMutex       sync.RWMutex
	doMoreThingsArgsForCall []struct {
		arg1 int
		arg2 int
		arg3 []string
	}
	doMoreThingsReturns struct {
		result1 int
	}
}

func (fake *FakeHasVarArgs) DoThings(arg1 int, arg2 ...string) int {
	fake.doThingsMutex.Lock()
	fake.doThingsArgsForCall = append(fake.doThingsArgsForCall, struct {
		arg1 int
		arg2 []string
	}{arg1, arg2})
	fake.doThingsMutex.Unlock()
	if fake.DoThingsStub != nil {
		return fake.DoThingsStub(arg1, arg2...)
	} else {
		return fake.doThingsReturns.result1
	}
}

func (fake *FakeHasVarArgs) DoThingsCallCount() int {
	fake.doThingsMutex.RLock()
	defer fake.doThingsMutex.RUnlock()
	return len(fake.doThingsArgsForCall)
}

func (fake *FakeHasVarArgs) DoThingsArgsForCall(i int) (int, []string) {
	fake.doThingsMutex.RLock()
	defer fake.doThingsMutex.RUnlock()
	return fake.doThingsArgsForCall[i].arg1, fake.doThingsArgsForCall[i].arg2
}

func (fake *FakeHasVarArgs) DoThingsReturns(result1 int) {
	fake.DoThingsStub = nil
	fake.doThingsReturns = struct {
		result1 int
	}{result1}
}

func (fake *FakeHasVarArgs) DoMoreThings(arg1 int, arg2 int, arg3 ...string) int {
	fake.doMoreThingsMutex.Lock()
	fake.doMoreThingsArgsForCall = append(fake.doMoreThingsArgsForCall, struct {
		arg1 int
		arg2 int
		arg3 []string
	}{arg1, arg2, arg3})
	fake.doMoreThingsMutex.Unlock()
	if fake.DoMoreThingsStub != nil {
		return fake.DoMoreThingsStub(arg1, arg2, arg3...)
	} else {
		return fake.doMoreThingsReturns.result1
	}
}

func (fake *FakeHasVarArgs) DoMoreThingsCallCount() int {
	fake.doMoreThingsMutex.RLock()
	defer fake.doMoreThingsMutex.RUnlock()
	return len(fake.doMoreThingsArgsForCall)
}

func (fake *FakeHasVarArgs) DoMoreThingsArgsForCall(i int) (int, int, []string) {
	fake.doMoreThingsMutex.RLock()
	defer fake.doMoreThingsMutex.RUnlock()
	return fake.doMoreThingsArgsForCall[i].arg1, fake.doMoreThingsArgsForCall[i].arg2, fake.doMoreThingsArgsForCall[i].arg3
}

func (fake *FakeHasVarArgs) DoMoreThingsReturns(result1 int) {
	fake.DoMoreThingsStub = nil
	fake.doMoreThingsReturns = struct {
		result1 int
	}{result1}
}

var _ fixtures.HasVarArgs = new(FakeHasVarArgs)
