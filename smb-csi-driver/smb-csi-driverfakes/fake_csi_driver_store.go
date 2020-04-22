// Code generated by counterfeiter. DO NOT EDIT.
package smbcsidriverfakes

import (
	"sync"

	"code.cloudfoundry.org/smb-csi-driver/nodeserver"
	"github.com/container-storage-interface/spec/lib/go/csi"
)

type FakeCSIDriverStore struct {
	CreateStub        func(string, *csi.NodePublishVolumeRequest) error
	createMutex       sync.RWMutex
	createArgsForCall []struct {
		arg1 string
		arg2 *csi.NodePublishVolumeRequest
	}
	createReturns struct {
		result1 error
	}
	createReturnsOnCall map[int]struct {
		result1 error
	}
	DeleteStub        func(string)
	deleteMutex       sync.RWMutex
	deleteArgsForCall []struct {
		arg1 string
	}
	GetStub        func(string, *csi.NodePublishVolumeRequest) (bool, bool, error)
	getMutex       sync.RWMutex
	getArgsForCall []struct {
		arg1 string
		arg2 *csi.NodePublishVolumeRequest
	}
	getReturns struct {
		result1 bool
		result2 bool
		result3 error
	}
	getReturnsOnCall map[int]struct {
		result1 bool
		result2 bool
		result3 error
	}
	invocations      map[string][][]interface{}
	invocationsMutex sync.RWMutex
}

func (fake *FakeCSIDriverStore) Create(arg1 string, arg2 *csi.NodePublishVolumeRequest) error {
	fake.createMutex.Lock()
	ret, specificReturn := fake.createReturnsOnCall[len(fake.createArgsForCall)]
	fake.createArgsForCall = append(fake.createArgsForCall, struct {
		arg1 string
		arg2 *csi.NodePublishVolumeRequest
	}{arg1, arg2})
	fake.recordInvocation("Create", []interface{}{arg1, arg2})
	fake.createMutex.Unlock()
	if fake.CreateStub != nil {
		return fake.CreateStub(arg1, arg2)
	}
	if specificReturn {
		return ret.result1
	}
	fakeReturns := fake.createReturns
	return fakeReturns.result1
}

func (fake *FakeCSIDriverStore) CreateCallCount() int {
	fake.createMutex.RLock()
	defer fake.createMutex.RUnlock()
	return len(fake.createArgsForCall)
}

func (fake *FakeCSIDriverStore) CreateCalls(stub func(string, *csi.NodePublishVolumeRequest) error) {
	fake.createMutex.Lock()
	defer fake.createMutex.Unlock()
	fake.CreateStub = stub
}

func (fake *FakeCSIDriverStore) CreateArgsForCall(i int) (string, *csi.NodePublishVolumeRequest) {
	fake.createMutex.RLock()
	defer fake.createMutex.RUnlock()
	argsForCall := fake.createArgsForCall[i]
	return argsForCall.arg1, argsForCall.arg2
}

func (fake *FakeCSIDriverStore) CreateReturns(result1 error) {
	fake.createMutex.Lock()
	defer fake.createMutex.Unlock()
	fake.CreateStub = nil
	fake.createReturns = struct {
		result1 error
	}{result1}
}

func (fake *FakeCSIDriverStore) CreateReturnsOnCall(i int, result1 error) {
	fake.createMutex.Lock()
	defer fake.createMutex.Unlock()
	fake.CreateStub = nil
	if fake.createReturnsOnCall == nil {
		fake.createReturnsOnCall = make(map[int]struct {
			result1 error
		})
	}
	fake.createReturnsOnCall[i] = struct {
		result1 error
	}{result1}
}

func (fake *FakeCSIDriverStore) Delete(arg1 string) {
	fake.deleteMutex.Lock()
	fake.deleteArgsForCall = append(fake.deleteArgsForCall, struct {
		arg1 string
	}{arg1})
	fake.recordInvocation("Delete", []interface{}{arg1})
	fake.deleteMutex.Unlock()
	if fake.DeleteStub != nil {
		fake.DeleteStub(arg1)
	}
}

func (fake *FakeCSIDriverStore) DeleteCallCount() int {
	fake.deleteMutex.RLock()
	defer fake.deleteMutex.RUnlock()
	return len(fake.deleteArgsForCall)
}

func (fake *FakeCSIDriverStore) DeleteCalls(stub func(string)) {
	fake.deleteMutex.Lock()
	defer fake.deleteMutex.Unlock()
	fake.DeleteStub = stub
}

func (fake *FakeCSIDriverStore) DeleteArgsForCall(i int) string {
	fake.deleteMutex.RLock()
	defer fake.deleteMutex.RUnlock()
	argsForCall := fake.deleteArgsForCall[i]
	return argsForCall.arg1
}

func (fake *FakeCSIDriverStore) Get(arg1 string, arg2 *csi.NodePublishVolumeRequest) (bool, bool, error) {
	fake.getMutex.Lock()
	ret, specificReturn := fake.getReturnsOnCall[len(fake.getArgsForCall)]
	fake.getArgsForCall = append(fake.getArgsForCall, struct {
		arg1 string
		arg2 *csi.NodePublishVolumeRequest
	}{arg1, arg2})
	fake.recordInvocation("Get", []interface{}{arg1, arg2})
	fake.getMutex.Unlock()
	if fake.GetStub != nil {
		return fake.GetStub(arg1, arg2)
	}
	if specificReturn {
		return ret.result1, ret.result2, ret.result3
	}
	fakeReturns := fake.getReturns
	return fakeReturns.result1, fakeReturns.result2, fakeReturns.result3
}

func (fake *FakeCSIDriverStore) GetCallCount() int {
	fake.getMutex.RLock()
	defer fake.getMutex.RUnlock()
	return len(fake.getArgsForCall)
}

func (fake *FakeCSIDriverStore) GetCalls(stub func(string, *csi.NodePublishVolumeRequest) (bool, bool, error)) {
	fake.getMutex.Lock()
	defer fake.getMutex.Unlock()
	fake.GetStub = stub
}

func (fake *FakeCSIDriverStore) GetArgsForCall(i int) (string, *csi.NodePublishVolumeRequest) {
	fake.getMutex.RLock()
	defer fake.getMutex.RUnlock()
	argsForCall := fake.getArgsForCall[i]
	return argsForCall.arg1, argsForCall.arg2
}

func (fake *FakeCSIDriverStore) GetReturns(result1 bool, result2 bool, result3 error) {
	fake.getMutex.Lock()
	defer fake.getMutex.Unlock()
	fake.GetStub = nil
	fake.getReturns = struct {
		result1 bool
		result2 bool
		result3 error
	}{result1, result2, result3}
}

func (fake *FakeCSIDriverStore) GetReturnsOnCall(i int, result1 bool, result2 bool, result3 error) {
	fake.getMutex.Lock()
	defer fake.getMutex.Unlock()
	fake.GetStub = nil
	if fake.getReturnsOnCall == nil {
		fake.getReturnsOnCall = make(map[int]struct {
			result1 bool
			result2 bool
			result3 error
		})
	}
	fake.getReturnsOnCall[i] = struct {
		result1 bool
		result2 bool
		result3 error
	}{result1, result2, result3}
}

func (fake *FakeCSIDriverStore) Invocations() map[string][][]interface{} {
	fake.invocationsMutex.RLock()
	defer fake.invocationsMutex.RUnlock()
	fake.createMutex.RLock()
	defer fake.createMutex.RUnlock()
	fake.deleteMutex.RLock()
	defer fake.deleteMutex.RUnlock()
	fake.getMutex.RLock()
	defer fake.getMutex.RUnlock()
	copiedInvocations := map[string][][]interface{}{}
	for key, value := range fake.invocations {
		copiedInvocations[key] = value
	}
	return copiedInvocations
}

func (fake *FakeCSIDriverStore) recordInvocation(key string, args []interface{}) {
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

var _ nodeserver.CSIDriverStore = new(FakeCSIDriverStore)
