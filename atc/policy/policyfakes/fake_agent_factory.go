// Code generated by counterfeiter. DO NOT EDIT.
package policyfakes

import (
	"sync"

	lager "code.cloudfoundry.org/lager/v3"
	"github.com/concourse/concourse/atc/policy"
)

type FakeAgentFactory struct {
	DescriptionStub        func() string
	descriptionMutex       sync.RWMutex
	descriptionArgsForCall []struct {
	}
	descriptionReturns struct {
		result1 string
	}
	descriptionReturnsOnCall map[int]struct {
		result1 string
	}
	IsConfiguredStub        func() bool
	isConfiguredMutex       sync.RWMutex
	isConfiguredArgsForCall []struct {
	}
	isConfiguredReturns struct {
		result1 bool
	}
	isConfiguredReturnsOnCall map[int]struct {
		result1 bool
	}
	NewAgentStub        func(lager.Logger) (policy.Agent, error)
	newAgentMutex       sync.RWMutex
	newAgentArgsForCall []struct {
		arg1 lager.Logger
	}
	newAgentReturns struct {
		result1 policy.Agent
		result2 error
	}
	newAgentReturnsOnCall map[int]struct {
		result1 policy.Agent
		result2 error
	}
	invocations      map[string][][]interface{}
	invocationsMutex sync.RWMutex
}

func (fake *FakeAgentFactory) Description() string {
	fake.descriptionMutex.Lock()
	ret, specificReturn := fake.descriptionReturnsOnCall[len(fake.descriptionArgsForCall)]
	fake.descriptionArgsForCall = append(fake.descriptionArgsForCall, struct {
	}{})
	stub := fake.DescriptionStub
	fakeReturns := fake.descriptionReturns
	fake.recordInvocation("Description", []interface{}{})
	fake.descriptionMutex.Unlock()
	if stub != nil {
		return stub()
	}
	if specificReturn {
		return ret.result1
	}
	return fakeReturns.result1
}

func (fake *FakeAgentFactory) DescriptionCallCount() int {
	fake.descriptionMutex.RLock()
	defer fake.descriptionMutex.RUnlock()
	return len(fake.descriptionArgsForCall)
}

func (fake *FakeAgentFactory) DescriptionCalls(stub func() string) {
	fake.descriptionMutex.Lock()
	defer fake.descriptionMutex.Unlock()
	fake.DescriptionStub = stub
}

func (fake *FakeAgentFactory) DescriptionReturns(result1 string) {
	fake.descriptionMutex.Lock()
	defer fake.descriptionMutex.Unlock()
	fake.DescriptionStub = nil
	fake.descriptionReturns = struct {
		result1 string
	}{result1}
}

func (fake *FakeAgentFactory) DescriptionReturnsOnCall(i int, result1 string) {
	fake.descriptionMutex.Lock()
	defer fake.descriptionMutex.Unlock()
	fake.DescriptionStub = nil
	if fake.descriptionReturnsOnCall == nil {
		fake.descriptionReturnsOnCall = make(map[int]struct {
			result1 string
		})
	}
	fake.descriptionReturnsOnCall[i] = struct {
		result1 string
	}{result1}
}

func (fake *FakeAgentFactory) IsConfigured() bool {
	fake.isConfiguredMutex.Lock()
	ret, specificReturn := fake.isConfiguredReturnsOnCall[len(fake.isConfiguredArgsForCall)]
	fake.isConfiguredArgsForCall = append(fake.isConfiguredArgsForCall, struct {
	}{})
	stub := fake.IsConfiguredStub
	fakeReturns := fake.isConfiguredReturns
	fake.recordInvocation("IsConfigured", []interface{}{})
	fake.isConfiguredMutex.Unlock()
	if stub != nil {
		return stub()
	}
	if specificReturn {
		return ret.result1
	}
	return fakeReturns.result1
}

func (fake *FakeAgentFactory) IsConfiguredCallCount() int {
	fake.isConfiguredMutex.RLock()
	defer fake.isConfiguredMutex.RUnlock()
	return len(fake.isConfiguredArgsForCall)
}

func (fake *FakeAgentFactory) IsConfiguredCalls(stub func() bool) {
	fake.isConfiguredMutex.Lock()
	defer fake.isConfiguredMutex.Unlock()
	fake.IsConfiguredStub = stub
}

func (fake *FakeAgentFactory) IsConfiguredReturns(result1 bool) {
	fake.isConfiguredMutex.Lock()
	defer fake.isConfiguredMutex.Unlock()
	fake.IsConfiguredStub = nil
	fake.isConfiguredReturns = struct {
		result1 bool
	}{result1}
}

func (fake *FakeAgentFactory) IsConfiguredReturnsOnCall(i int, result1 bool) {
	fake.isConfiguredMutex.Lock()
	defer fake.isConfiguredMutex.Unlock()
	fake.IsConfiguredStub = nil
	if fake.isConfiguredReturnsOnCall == nil {
		fake.isConfiguredReturnsOnCall = make(map[int]struct {
			result1 bool
		})
	}
	fake.isConfiguredReturnsOnCall[i] = struct {
		result1 bool
	}{result1}
}

func (fake *FakeAgentFactory) NewAgent(arg1 lager.Logger) (policy.Agent, error) {
	fake.newAgentMutex.Lock()
	ret, specificReturn := fake.newAgentReturnsOnCall[len(fake.newAgentArgsForCall)]
	fake.newAgentArgsForCall = append(fake.newAgentArgsForCall, struct {
		arg1 lager.Logger
	}{arg1})
	stub := fake.NewAgentStub
	fakeReturns := fake.newAgentReturns
	fake.recordInvocation("NewAgent", []interface{}{arg1})
	fake.newAgentMutex.Unlock()
	if stub != nil {
		return stub(arg1)
	}
	if specificReturn {
		return ret.result1, ret.result2
	}
	return fakeReturns.result1, fakeReturns.result2
}

func (fake *FakeAgentFactory) NewAgentCallCount() int {
	fake.newAgentMutex.RLock()
	defer fake.newAgentMutex.RUnlock()
	return len(fake.newAgentArgsForCall)
}

func (fake *FakeAgentFactory) NewAgentCalls(stub func(lager.Logger) (policy.Agent, error)) {
	fake.newAgentMutex.Lock()
	defer fake.newAgentMutex.Unlock()
	fake.NewAgentStub = stub
}

func (fake *FakeAgentFactory) NewAgentArgsForCall(i int) lager.Logger {
	fake.newAgentMutex.RLock()
	defer fake.newAgentMutex.RUnlock()
	argsForCall := fake.newAgentArgsForCall[i]
	return argsForCall.arg1
}

func (fake *FakeAgentFactory) NewAgentReturns(result1 policy.Agent, result2 error) {
	fake.newAgentMutex.Lock()
	defer fake.newAgentMutex.Unlock()
	fake.NewAgentStub = nil
	fake.newAgentReturns = struct {
		result1 policy.Agent
		result2 error
	}{result1, result2}
}

func (fake *FakeAgentFactory) NewAgentReturnsOnCall(i int, result1 policy.Agent, result2 error) {
	fake.newAgentMutex.Lock()
	defer fake.newAgentMutex.Unlock()
	fake.NewAgentStub = nil
	if fake.newAgentReturnsOnCall == nil {
		fake.newAgentReturnsOnCall = make(map[int]struct {
			result1 policy.Agent
			result2 error
		})
	}
	fake.newAgentReturnsOnCall[i] = struct {
		result1 policy.Agent
		result2 error
	}{result1, result2}
}

func (fake *FakeAgentFactory) Invocations() map[string][][]interface{} {
	fake.invocationsMutex.RLock()
	defer fake.invocationsMutex.RUnlock()
	fake.descriptionMutex.RLock()
	defer fake.descriptionMutex.RUnlock()
	fake.isConfiguredMutex.RLock()
	defer fake.isConfiguredMutex.RUnlock()
	fake.newAgentMutex.RLock()
	defer fake.newAgentMutex.RUnlock()
	copiedInvocations := map[string][][]interface{}{}
	for key, value := range fake.invocations {
		copiedInvocations[key] = value
	}
	return copiedInvocations
}

func (fake *FakeAgentFactory) recordInvocation(key string, args []interface{}) {
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

var _ policy.AgentFactory = new(FakeAgentFactory)
