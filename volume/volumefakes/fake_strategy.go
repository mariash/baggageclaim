// Code generated by counterfeiter. DO NOT EDIT.
package volumefakes

import (
	sync "sync"

	lager "code.cloudfoundry.org/lager"
	volume "github.com/concourse/baggageclaim/volume"
)

type FakeStrategy struct {
	MaterializeStub        func(lager.Logger, string, volume.Filesystem, volume.Streamer) (volume.FilesystemInitVolume, error)
	materializeMutex       sync.RWMutex
	materializeArgsForCall []struct {
		arg1 lager.Logger
		arg2 string
		arg3 volume.Filesystem
		arg4 volume.Streamer
	}
	materializeReturns struct {
		result1 volume.FilesystemInitVolume
		result2 error
	}
	materializeReturnsOnCall map[int]struct {
		result1 volume.FilesystemInitVolume
		result2 error
	}
	invocations      map[string][][]interface{}
	invocationsMutex sync.RWMutex
}

func (fake *FakeStrategy) Materialize(arg1 lager.Logger, arg2 string, arg3 volume.Filesystem, arg4 volume.Streamer) (volume.FilesystemInitVolume, error) {
	fake.materializeMutex.Lock()
	ret, specificReturn := fake.materializeReturnsOnCall[len(fake.materializeArgsForCall)]
	fake.materializeArgsForCall = append(fake.materializeArgsForCall, struct {
		arg1 lager.Logger
		arg2 string
		arg3 volume.Filesystem
		arg4 volume.Streamer
	}{arg1, arg2, arg3, arg4})
	fake.recordInvocation("Materialize", []interface{}{arg1, arg2, arg3, arg4})
	fake.materializeMutex.Unlock()
	if fake.MaterializeStub != nil {
		return fake.MaterializeStub(arg1, arg2, arg3, arg4)
	}
	if specificReturn {
		return ret.result1, ret.result2
	}
	fakeReturns := fake.materializeReturns
	return fakeReturns.result1, fakeReturns.result2
}

func (fake *FakeStrategy) MaterializeCallCount() int {
	fake.materializeMutex.RLock()
	defer fake.materializeMutex.RUnlock()
	return len(fake.materializeArgsForCall)
}

func (fake *FakeStrategy) MaterializeCalls(stub func(lager.Logger, string, volume.Filesystem, volume.Streamer) (volume.FilesystemInitVolume, error)) {
	fake.materializeMutex.Lock()
	defer fake.materializeMutex.Unlock()
	fake.MaterializeStub = stub
}

func (fake *FakeStrategy) MaterializeArgsForCall(i int) (lager.Logger, string, volume.Filesystem, volume.Streamer) {
	fake.materializeMutex.RLock()
	defer fake.materializeMutex.RUnlock()
	argsForCall := fake.materializeArgsForCall[i]
	return argsForCall.arg1, argsForCall.arg2, argsForCall.arg3, argsForCall.arg4
}

func (fake *FakeStrategy) MaterializeReturns(result1 volume.FilesystemInitVolume, result2 error) {
	fake.materializeMutex.Lock()
	defer fake.materializeMutex.Unlock()
	fake.MaterializeStub = nil
	fake.materializeReturns = struct {
		result1 volume.FilesystemInitVolume
		result2 error
	}{result1, result2}
}

func (fake *FakeStrategy) MaterializeReturnsOnCall(i int, result1 volume.FilesystemInitVolume, result2 error) {
	fake.materializeMutex.Lock()
	defer fake.materializeMutex.Unlock()
	fake.MaterializeStub = nil
	if fake.materializeReturnsOnCall == nil {
		fake.materializeReturnsOnCall = make(map[int]struct {
			result1 volume.FilesystemInitVolume
			result2 error
		})
	}
	fake.materializeReturnsOnCall[i] = struct {
		result1 volume.FilesystemInitVolume
		result2 error
	}{result1, result2}
}

func (fake *FakeStrategy) Invocations() map[string][][]interface{} {
	fake.invocationsMutex.RLock()
	defer fake.invocationsMutex.RUnlock()
	fake.materializeMutex.RLock()
	defer fake.materializeMutex.RUnlock()
	copiedInvocations := map[string][][]interface{}{}
	for key, value := range fake.invocations {
		copiedInvocations[key] = value
	}
	return copiedInvocations
}

func (fake *FakeStrategy) recordInvocation(key string, args []interface{}) {
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

var _ volume.Strategy = new(FakeStrategy)
