// Code generated by counterfeiter. DO NOT EDIT.
package dbfakes

import (
	"sync"
	"time"

	"github.com/concourse/concourse/atc/db"
)

type FakeContainerRepository struct {
	DestroyDirtyInMemoryBuildContainersStub        func() (int, error)
	destroyDirtyInMemoryBuildContainersMutex       sync.RWMutex
	destroyDirtyInMemoryBuildContainersArgsForCall []struct {
	}
	destroyDirtyInMemoryBuildContainersReturns struct {
		result1 int
		result2 error
	}
	destroyDirtyInMemoryBuildContainersReturnsOnCall map[int]struct {
		result1 int
		result2 error
	}
	DestroyFailedContainersStub        func() (int, error)
	destroyFailedContainersMutex       sync.RWMutex
	destroyFailedContainersArgsForCall []struct {
	}
	destroyFailedContainersReturns struct {
		result1 int
		result2 error
	}
	destroyFailedContainersReturnsOnCall map[int]struct {
		result1 int
		result2 error
	}
	DestroyUnknownContainersStub        func(string, []string) (int, error)
	destroyUnknownContainersMutex       sync.RWMutex
	destroyUnknownContainersArgsForCall []struct {
		arg1 string
		arg2 []string
	}
	destroyUnknownContainersReturns struct {
		result1 int
		result2 error
	}
	destroyUnknownContainersReturnsOnCall map[int]struct {
		result1 int
		result2 error
	}
	FindDestroyingContainersStub        func(string) ([]string, error)
	findDestroyingContainersMutex       sync.RWMutex
	findDestroyingContainersArgsForCall []struct {
		arg1 string
	}
	findDestroyingContainersReturns struct {
		result1 []string
		result2 error
	}
	findDestroyingContainersReturnsOnCall map[int]struct {
		result1 []string
		result2 error
	}
	FindOrphanedContainersStub        func() ([]db.CreatingContainer, []db.CreatedContainer, []db.DestroyingContainer, error)
	findOrphanedContainersMutex       sync.RWMutex
	findOrphanedContainersArgsForCall []struct {
	}
	findOrphanedContainersReturns struct {
		result1 []db.CreatingContainer
		result2 []db.CreatedContainer
		result3 []db.DestroyingContainer
		result4 error
	}
	findOrphanedContainersReturnsOnCall map[int]struct {
		result1 []db.CreatingContainer
		result2 []db.CreatedContainer
		result3 []db.DestroyingContainer
		result4 error
	}
	RemoveDestroyingContainersStub        func(string, []string) (int, error)
	removeDestroyingContainersMutex       sync.RWMutex
	removeDestroyingContainersArgsForCall []struct {
		arg1 string
		arg2 []string
	}
	removeDestroyingContainersReturns struct {
		result1 int
		result2 error
	}
	removeDestroyingContainersReturnsOnCall map[int]struct {
		result1 int
		result2 error
	}
	RemoveMissingContainersStub        func(time.Duration) (int, error)
	removeMissingContainersMutex       sync.RWMutex
	removeMissingContainersArgsForCall []struct {
		arg1 time.Duration
	}
	removeMissingContainersReturns struct {
		result1 int
		result2 error
	}
	removeMissingContainersReturnsOnCall map[int]struct {
		result1 int
		result2 error
	}
	UpdateContainersMissingSinceStub        func(string, []string) error
	updateContainersMissingSinceMutex       sync.RWMutex
	updateContainersMissingSinceArgsForCall []struct {
		arg1 string
		arg2 []string
	}
	updateContainersMissingSinceReturns struct {
		result1 error
	}
	updateContainersMissingSinceReturnsOnCall map[int]struct {
		result1 error
	}
	invocations      map[string][][]interface{}
	invocationsMutex sync.RWMutex
}

func (fake *FakeContainerRepository) DestroyDirtyInMemoryBuildContainers() (int, error) {
	fake.destroyDirtyInMemoryBuildContainersMutex.Lock()
	ret, specificReturn := fake.destroyDirtyInMemoryBuildContainersReturnsOnCall[len(fake.destroyDirtyInMemoryBuildContainersArgsForCall)]
	fake.destroyDirtyInMemoryBuildContainersArgsForCall = append(fake.destroyDirtyInMemoryBuildContainersArgsForCall, struct {
	}{})
	stub := fake.DestroyDirtyInMemoryBuildContainersStub
	fakeReturns := fake.destroyDirtyInMemoryBuildContainersReturns
	fake.recordInvocation("DestroyDirtyInMemoryBuildContainers", []interface{}{})
	fake.destroyDirtyInMemoryBuildContainersMutex.Unlock()
	if stub != nil {
		return stub()
	}
	if specificReturn {
		return ret.result1, ret.result2
	}
	return fakeReturns.result1, fakeReturns.result2
}

func (fake *FakeContainerRepository) DestroyDirtyInMemoryBuildContainersCallCount() int {
	fake.destroyDirtyInMemoryBuildContainersMutex.RLock()
	defer fake.destroyDirtyInMemoryBuildContainersMutex.RUnlock()
	return len(fake.destroyDirtyInMemoryBuildContainersArgsForCall)
}

func (fake *FakeContainerRepository) DestroyDirtyInMemoryBuildContainersCalls(stub func() (int, error)) {
	fake.destroyDirtyInMemoryBuildContainersMutex.Lock()
	defer fake.destroyDirtyInMemoryBuildContainersMutex.Unlock()
	fake.DestroyDirtyInMemoryBuildContainersStub = stub
}

func (fake *FakeContainerRepository) DestroyDirtyInMemoryBuildContainersReturns(result1 int, result2 error) {
	fake.destroyDirtyInMemoryBuildContainersMutex.Lock()
	defer fake.destroyDirtyInMemoryBuildContainersMutex.Unlock()
	fake.DestroyDirtyInMemoryBuildContainersStub = nil
	fake.destroyDirtyInMemoryBuildContainersReturns = struct {
		result1 int
		result2 error
	}{result1, result2}
}

func (fake *FakeContainerRepository) DestroyDirtyInMemoryBuildContainersReturnsOnCall(i int, result1 int, result2 error) {
	fake.destroyDirtyInMemoryBuildContainersMutex.Lock()
	defer fake.destroyDirtyInMemoryBuildContainersMutex.Unlock()
	fake.DestroyDirtyInMemoryBuildContainersStub = nil
	if fake.destroyDirtyInMemoryBuildContainersReturnsOnCall == nil {
		fake.destroyDirtyInMemoryBuildContainersReturnsOnCall = make(map[int]struct {
			result1 int
			result2 error
		})
	}
	fake.destroyDirtyInMemoryBuildContainersReturnsOnCall[i] = struct {
		result1 int
		result2 error
	}{result1, result2}
}

func (fake *FakeContainerRepository) DestroyFailedContainers() (int, error) {
	fake.destroyFailedContainersMutex.Lock()
	ret, specificReturn := fake.destroyFailedContainersReturnsOnCall[len(fake.destroyFailedContainersArgsForCall)]
	fake.destroyFailedContainersArgsForCall = append(fake.destroyFailedContainersArgsForCall, struct {
	}{})
	stub := fake.DestroyFailedContainersStub
	fakeReturns := fake.destroyFailedContainersReturns
	fake.recordInvocation("DestroyFailedContainers", []interface{}{})
	fake.destroyFailedContainersMutex.Unlock()
	if stub != nil {
		return stub()
	}
	if specificReturn {
		return ret.result1, ret.result2
	}
	return fakeReturns.result1, fakeReturns.result2
}

func (fake *FakeContainerRepository) DestroyFailedContainersCallCount() int {
	fake.destroyFailedContainersMutex.RLock()
	defer fake.destroyFailedContainersMutex.RUnlock()
	return len(fake.destroyFailedContainersArgsForCall)
}

func (fake *FakeContainerRepository) DestroyFailedContainersCalls(stub func() (int, error)) {
	fake.destroyFailedContainersMutex.Lock()
	defer fake.destroyFailedContainersMutex.Unlock()
	fake.DestroyFailedContainersStub = stub
}

func (fake *FakeContainerRepository) DestroyFailedContainersReturns(result1 int, result2 error) {
	fake.destroyFailedContainersMutex.Lock()
	defer fake.destroyFailedContainersMutex.Unlock()
	fake.DestroyFailedContainersStub = nil
	fake.destroyFailedContainersReturns = struct {
		result1 int
		result2 error
	}{result1, result2}
}

func (fake *FakeContainerRepository) DestroyFailedContainersReturnsOnCall(i int, result1 int, result2 error) {
	fake.destroyFailedContainersMutex.Lock()
	defer fake.destroyFailedContainersMutex.Unlock()
	fake.DestroyFailedContainersStub = nil
	if fake.destroyFailedContainersReturnsOnCall == nil {
		fake.destroyFailedContainersReturnsOnCall = make(map[int]struct {
			result1 int
			result2 error
		})
	}
	fake.destroyFailedContainersReturnsOnCall[i] = struct {
		result1 int
		result2 error
	}{result1, result2}
}

func (fake *FakeContainerRepository) DestroyUnknownContainers(arg1 string, arg2 []string) (int, error) {
	var arg2Copy []string
	if arg2 != nil {
		arg2Copy = make([]string, len(arg2))
		copy(arg2Copy, arg2)
	}
	fake.destroyUnknownContainersMutex.Lock()
	ret, specificReturn := fake.destroyUnknownContainersReturnsOnCall[len(fake.destroyUnknownContainersArgsForCall)]
	fake.destroyUnknownContainersArgsForCall = append(fake.destroyUnknownContainersArgsForCall, struct {
		arg1 string
		arg2 []string
	}{arg1, arg2Copy})
	stub := fake.DestroyUnknownContainersStub
	fakeReturns := fake.destroyUnknownContainersReturns
	fake.recordInvocation("DestroyUnknownContainers", []interface{}{arg1, arg2Copy})
	fake.destroyUnknownContainersMutex.Unlock()
	if stub != nil {
		return stub(arg1, arg2)
	}
	if specificReturn {
		return ret.result1, ret.result2
	}
	return fakeReturns.result1, fakeReturns.result2
}

func (fake *FakeContainerRepository) DestroyUnknownContainersCallCount() int {
	fake.destroyUnknownContainersMutex.RLock()
	defer fake.destroyUnknownContainersMutex.RUnlock()
	return len(fake.destroyUnknownContainersArgsForCall)
}

func (fake *FakeContainerRepository) DestroyUnknownContainersCalls(stub func(string, []string) (int, error)) {
	fake.destroyUnknownContainersMutex.Lock()
	defer fake.destroyUnknownContainersMutex.Unlock()
	fake.DestroyUnknownContainersStub = stub
}

func (fake *FakeContainerRepository) DestroyUnknownContainersArgsForCall(i int) (string, []string) {
	fake.destroyUnknownContainersMutex.RLock()
	defer fake.destroyUnknownContainersMutex.RUnlock()
	argsForCall := fake.destroyUnknownContainersArgsForCall[i]
	return argsForCall.arg1, argsForCall.arg2
}

func (fake *FakeContainerRepository) DestroyUnknownContainersReturns(result1 int, result2 error) {
	fake.destroyUnknownContainersMutex.Lock()
	defer fake.destroyUnknownContainersMutex.Unlock()
	fake.DestroyUnknownContainersStub = nil
	fake.destroyUnknownContainersReturns = struct {
		result1 int
		result2 error
	}{result1, result2}
}

func (fake *FakeContainerRepository) DestroyUnknownContainersReturnsOnCall(i int, result1 int, result2 error) {
	fake.destroyUnknownContainersMutex.Lock()
	defer fake.destroyUnknownContainersMutex.Unlock()
	fake.DestroyUnknownContainersStub = nil
	if fake.destroyUnknownContainersReturnsOnCall == nil {
		fake.destroyUnknownContainersReturnsOnCall = make(map[int]struct {
			result1 int
			result2 error
		})
	}
	fake.destroyUnknownContainersReturnsOnCall[i] = struct {
		result1 int
		result2 error
	}{result1, result2}
}

func (fake *FakeContainerRepository) FindDestroyingContainers(arg1 string) ([]string, error) {
	fake.findDestroyingContainersMutex.Lock()
	ret, specificReturn := fake.findDestroyingContainersReturnsOnCall[len(fake.findDestroyingContainersArgsForCall)]
	fake.findDestroyingContainersArgsForCall = append(fake.findDestroyingContainersArgsForCall, struct {
		arg1 string
	}{arg1})
	stub := fake.FindDestroyingContainersStub
	fakeReturns := fake.findDestroyingContainersReturns
	fake.recordInvocation("FindDestroyingContainers", []interface{}{arg1})
	fake.findDestroyingContainersMutex.Unlock()
	if stub != nil {
		return stub(arg1)
	}
	if specificReturn {
		return ret.result1, ret.result2
	}
	return fakeReturns.result1, fakeReturns.result2
}

func (fake *FakeContainerRepository) FindDestroyingContainersCallCount() int {
	fake.findDestroyingContainersMutex.RLock()
	defer fake.findDestroyingContainersMutex.RUnlock()
	return len(fake.findDestroyingContainersArgsForCall)
}

func (fake *FakeContainerRepository) FindDestroyingContainersCalls(stub func(string) ([]string, error)) {
	fake.findDestroyingContainersMutex.Lock()
	defer fake.findDestroyingContainersMutex.Unlock()
	fake.FindDestroyingContainersStub = stub
}

func (fake *FakeContainerRepository) FindDestroyingContainersArgsForCall(i int) string {
	fake.findDestroyingContainersMutex.RLock()
	defer fake.findDestroyingContainersMutex.RUnlock()
	argsForCall := fake.findDestroyingContainersArgsForCall[i]
	return argsForCall.arg1
}

func (fake *FakeContainerRepository) FindDestroyingContainersReturns(result1 []string, result2 error) {
	fake.findDestroyingContainersMutex.Lock()
	defer fake.findDestroyingContainersMutex.Unlock()
	fake.FindDestroyingContainersStub = nil
	fake.findDestroyingContainersReturns = struct {
		result1 []string
		result2 error
	}{result1, result2}
}

func (fake *FakeContainerRepository) FindDestroyingContainersReturnsOnCall(i int, result1 []string, result2 error) {
	fake.findDestroyingContainersMutex.Lock()
	defer fake.findDestroyingContainersMutex.Unlock()
	fake.FindDestroyingContainersStub = nil
	if fake.findDestroyingContainersReturnsOnCall == nil {
		fake.findDestroyingContainersReturnsOnCall = make(map[int]struct {
			result1 []string
			result2 error
		})
	}
	fake.findDestroyingContainersReturnsOnCall[i] = struct {
		result1 []string
		result2 error
	}{result1, result2}
}

func (fake *FakeContainerRepository) FindOrphanedContainers() ([]db.CreatingContainer, []db.CreatedContainer, []db.DestroyingContainer, error) {
	fake.findOrphanedContainersMutex.Lock()
	ret, specificReturn := fake.findOrphanedContainersReturnsOnCall[len(fake.findOrphanedContainersArgsForCall)]
	fake.findOrphanedContainersArgsForCall = append(fake.findOrphanedContainersArgsForCall, struct {
	}{})
	stub := fake.FindOrphanedContainersStub
	fakeReturns := fake.findOrphanedContainersReturns
	fake.recordInvocation("FindOrphanedContainers", []interface{}{})
	fake.findOrphanedContainersMutex.Unlock()
	if stub != nil {
		return stub()
	}
	if specificReturn {
		return ret.result1, ret.result2, ret.result3, ret.result4
	}
	return fakeReturns.result1, fakeReturns.result2, fakeReturns.result3, fakeReturns.result4
}

func (fake *FakeContainerRepository) FindOrphanedContainersCallCount() int {
	fake.findOrphanedContainersMutex.RLock()
	defer fake.findOrphanedContainersMutex.RUnlock()
	return len(fake.findOrphanedContainersArgsForCall)
}

func (fake *FakeContainerRepository) FindOrphanedContainersCalls(stub func() ([]db.CreatingContainer, []db.CreatedContainer, []db.DestroyingContainer, error)) {
	fake.findOrphanedContainersMutex.Lock()
	defer fake.findOrphanedContainersMutex.Unlock()
	fake.FindOrphanedContainersStub = stub
}

func (fake *FakeContainerRepository) FindOrphanedContainersReturns(result1 []db.CreatingContainer, result2 []db.CreatedContainer, result3 []db.DestroyingContainer, result4 error) {
	fake.findOrphanedContainersMutex.Lock()
	defer fake.findOrphanedContainersMutex.Unlock()
	fake.FindOrphanedContainersStub = nil
	fake.findOrphanedContainersReturns = struct {
		result1 []db.CreatingContainer
		result2 []db.CreatedContainer
		result3 []db.DestroyingContainer
		result4 error
	}{result1, result2, result3, result4}
}

func (fake *FakeContainerRepository) FindOrphanedContainersReturnsOnCall(i int, result1 []db.CreatingContainer, result2 []db.CreatedContainer, result3 []db.DestroyingContainer, result4 error) {
	fake.findOrphanedContainersMutex.Lock()
	defer fake.findOrphanedContainersMutex.Unlock()
	fake.FindOrphanedContainersStub = nil
	if fake.findOrphanedContainersReturnsOnCall == nil {
		fake.findOrphanedContainersReturnsOnCall = make(map[int]struct {
			result1 []db.CreatingContainer
			result2 []db.CreatedContainer
			result3 []db.DestroyingContainer
			result4 error
		})
	}
	fake.findOrphanedContainersReturnsOnCall[i] = struct {
		result1 []db.CreatingContainer
		result2 []db.CreatedContainer
		result3 []db.DestroyingContainer
		result4 error
	}{result1, result2, result3, result4}
}

func (fake *FakeContainerRepository) RemoveDestroyingContainers(arg1 string, arg2 []string) (int, error) {
	var arg2Copy []string
	if arg2 != nil {
		arg2Copy = make([]string, len(arg2))
		copy(arg2Copy, arg2)
	}
	fake.removeDestroyingContainersMutex.Lock()
	ret, specificReturn := fake.removeDestroyingContainersReturnsOnCall[len(fake.removeDestroyingContainersArgsForCall)]
	fake.removeDestroyingContainersArgsForCall = append(fake.removeDestroyingContainersArgsForCall, struct {
		arg1 string
		arg2 []string
	}{arg1, arg2Copy})
	stub := fake.RemoveDestroyingContainersStub
	fakeReturns := fake.removeDestroyingContainersReturns
	fake.recordInvocation("RemoveDestroyingContainers", []interface{}{arg1, arg2Copy})
	fake.removeDestroyingContainersMutex.Unlock()
	if stub != nil {
		return stub(arg1, arg2)
	}
	if specificReturn {
		return ret.result1, ret.result2
	}
	return fakeReturns.result1, fakeReturns.result2
}

func (fake *FakeContainerRepository) RemoveDestroyingContainersCallCount() int {
	fake.removeDestroyingContainersMutex.RLock()
	defer fake.removeDestroyingContainersMutex.RUnlock()
	return len(fake.removeDestroyingContainersArgsForCall)
}

func (fake *FakeContainerRepository) RemoveDestroyingContainersCalls(stub func(string, []string) (int, error)) {
	fake.removeDestroyingContainersMutex.Lock()
	defer fake.removeDestroyingContainersMutex.Unlock()
	fake.RemoveDestroyingContainersStub = stub
}

func (fake *FakeContainerRepository) RemoveDestroyingContainersArgsForCall(i int) (string, []string) {
	fake.removeDestroyingContainersMutex.RLock()
	defer fake.removeDestroyingContainersMutex.RUnlock()
	argsForCall := fake.removeDestroyingContainersArgsForCall[i]
	return argsForCall.arg1, argsForCall.arg2
}

func (fake *FakeContainerRepository) RemoveDestroyingContainersReturns(result1 int, result2 error) {
	fake.removeDestroyingContainersMutex.Lock()
	defer fake.removeDestroyingContainersMutex.Unlock()
	fake.RemoveDestroyingContainersStub = nil
	fake.removeDestroyingContainersReturns = struct {
		result1 int
		result2 error
	}{result1, result2}
}

func (fake *FakeContainerRepository) RemoveDestroyingContainersReturnsOnCall(i int, result1 int, result2 error) {
	fake.removeDestroyingContainersMutex.Lock()
	defer fake.removeDestroyingContainersMutex.Unlock()
	fake.RemoveDestroyingContainersStub = nil
	if fake.removeDestroyingContainersReturnsOnCall == nil {
		fake.removeDestroyingContainersReturnsOnCall = make(map[int]struct {
			result1 int
			result2 error
		})
	}
	fake.removeDestroyingContainersReturnsOnCall[i] = struct {
		result1 int
		result2 error
	}{result1, result2}
}

func (fake *FakeContainerRepository) RemoveMissingContainers(arg1 time.Duration) (int, error) {
	fake.removeMissingContainersMutex.Lock()
	ret, specificReturn := fake.removeMissingContainersReturnsOnCall[len(fake.removeMissingContainersArgsForCall)]
	fake.removeMissingContainersArgsForCall = append(fake.removeMissingContainersArgsForCall, struct {
		arg1 time.Duration
	}{arg1})
	stub := fake.RemoveMissingContainersStub
	fakeReturns := fake.removeMissingContainersReturns
	fake.recordInvocation("RemoveMissingContainers", []interface{}{arg1})
	fake.removeMissingContainersMutex.Unlock()
	if stub != nil {
		return stub(arg1)
	}
	if specificReturn {
		return ret.result1, ret.result2
	}
	return fakeReturns.result1, fakeReturns.result2
}

func (fake *FakeContainerRepository) RemoveMissingContainersCallCount() int {
	fake.removeMissingContainersMutex.RLock()
	defer fake.removeMissingContainersMutex.RUnlock()
	return len(fake.removeMissingContainersArgsForCall)
}

func (fake *FakeContainerRepository) RemoveMissingContainersCalls(stub func(time.Duration) (int, error)) {
	fake.removeMissingContainersMutex.Lock()
	defer fake.removeMissingContainersMutex.Unlock()
	fake.RemoveMissingContainersStub = stub
}

func (fake *FakeContainerRepository) RemoveMissingContainersArgsForCall(i int) time.Duration {
	fake.removeMissingContainersMutex.RLock()
	defer fake.removeMissingContainersMutex.RUnlock()
	argsForCall := fake.removeMissingContainersArgsForCall[i]
	return argsForCall.arg1
}

func (fake *FakeContainerRepository) RemoveMissingContainersReturns(result1 int, result2 error) {
	fake.removeMissingContainersMutex.Lock()
	defer fake.removeMissingContainersMutex.Unlock()
	fake.RemoveMissingContainersStub = nil
	fake.removeMissingContainersReturns = struct {
		result1 int
		result2 error
	}{result1, result2}
}

func (fake *FakeContainerRepository) RemoveMissingContainersReturnsOnCall(i int, result1 int, result2 error) {
	fake.removeMissingContainersMutex.Lock()
	defer fake.removeMissingContainersMutex.Unlock()
	fake.RemoveMissingContainersStub = nil
	if fake.removeMissingContainersReturnsOnCall == nil {
		fake.removeMissingContainersReturnsOnCall = make(map[int]struct {
			result1 int
			result2 error
		})
	}
	fake.removeMissingContainersReturnsOnCall[i] = struct {
		result1 int
		result2 error
	}{result1, result2}
}

func (fake *FakeContainerRepository) UpdateContainersMissingSince(arg1 string, arg2 []string) error {
	var arg2Copy []string
	if arg2 != nil {
		arg2Copy = make([]string, len(arg2))
		copy(arg2Copy, arg2)
	}
	fake.updateContainersMissingSinceMutex.Lock()
	ret, specificReturn := fake.updateContainersMissingSinceReturnsOnCall[len(fake.updateContainersMissingSinceArgsForCall)]
	fake.updateContainersMissingSinceArgsForCall = append(fake.updateContainersMissingSinceArgsForCall, struct {
		arg1 string
		arg2 []string
	}{arg1, arg2Copy})
	stub := fake.UpdateContainersMissingSinceStub
	fakeReturns := fake.updateContainersMissingSinceReturns
	fake.recordInvocation("UpdateContainersMissingSince", []interface{}{arg1, arg2Copy})
	fake.updateContainersMissingSinceMutex.Unlock()
	if stub != nil {
		return stub(arg1, arg2)
	}
	if specificReturn {
		return ret.result1
	}
	return fakeReturns.result1
}

func (fake *FakeContainerRepository) UpdateContainersMissingSinceCallCount() int {
	fake.updateContainersMissingSinceMutex.RLock()
	defer fake.updateContainersMissingSinceMutex.RUnlock()
	return len(fake.updateContainersMissingSinceArgsForCall)
}

func (fake *FakeContainerRepository) UpdateContainersMissingSinceCalls(stub func(string, []string) error) {
	fake.updateContainersMissingSinceMutex.Lock()
	defer fake.updateContainersMissingSinceMutex.Unlock()
	fake.UpdateContainersMissingSinceStub = stub
}

func (fake *FakeContainerRepository) UpdateContainersMissingSinceArgsForCall(i int) (string, []string) {
	fake.updateContainersMissingSinceMutex.RLock()
	defer fake.updateContainersMissingSinceMutex.RUnlock()
	argsForCall := fake.updateContainersMissingSinceArgsForCall[i]
	return argsForCall.arg1, argsForCall.arg2
}

func (fake *FakeContainerRepository) UpdateContainersMissingSinceReturns(result1 error) {
	fake.updateContainersMissingSinceMutex.Lock()
	defer fake.updateContainersMissingSinceMutex.Unlock()
	fake.UpdateContainersMissingSinceStub = nil
	fake.updateContainersMissingSinceReturns = struct {
		result1 error
	}{result1}
}

func (fake *FakeContainerRepository) UpdateContainersMissingSinceReturnsOnCall(i int, result1 error) {
	fake.updateContainersMissingSinceMutex.Lock()
	defer fake.updateContainersMissingSinceMutex.Unlock()
	fake.UpdateContainersMissingSinceStub = nil
	if fake.updateContainersMissingSinceReturnsOnCall == nil {
		fake.updateContainersMissingSinceReturnsOnCall = make(map[int]struct {
			result1 error
		})
	}
	fake.updateContainersMissingSinceReturnsOnCall[i] = struct {
		result1 error
	}{result1}
}

func (fake *FakeContainerRepository) Invocations() map[string][][]interface{} {
	fake.invocationsMutex.RLock()
	defer fake.invocationsMutex.RUnlock()
	fake.destroyDirtyInMemoryBuildContainersMutex.RLock()
	defer fake.destroyDirtyInMemoryBuildContainersMutex.RUnlock()
	fake.destroyFailedContainersMutex.RLock()
	defer fake.destroyFailedContainersMutex.RUnlock()
	fake.destroyUnknownContainersMutex.RLock()
	defer fake.destroyUnknownContainersMutex.RUnlock()
	fake.findDestroyingContainersMutex.RLock()
	defer fake.findDestroyingContainersMutex.RUnlock()
	fake.findOrphanedContainersMutex.RLock()
	defer fake.findOrphanedContainersMutex.RUnlock()
	fake.removeDestroyingContainersMutex.RLock()
	defer fake.removeDestroyingContainersMutex.RUnlock()
	fake.removeMissingContainersMutex.RLock()
	defer fake.removeMissingContainersMutex.RUnlock()
	fake.updateContainersMissingSinceMutex.RLock()
	defer fake.updateContainersMissingSinceMutex.RUnlock()
	copiedInvocations := map[string][][]interface{}{}
	for key, value := range fake.invocations {
		copiedInvocations[key] = value
	}
	return copiedInvocations
}

func (fake *FakeContainerRepository) recordInvocation(key string, args []interface{}) {
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

var _ db.ContainerRepository = new(FakeContainerRepository)
