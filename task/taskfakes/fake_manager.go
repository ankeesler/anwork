// Code generated by counterfeiter. DO NOT EDIT.
package taskfakes

import (
	"sync"

	"github.com/ankeesler/anwork/task"
)

type FakeManager struct {
	CreateStub        func(name string) error
	createMutex       sync.RWMutex
	createArgsForCall []struct {
		name string
	}
	createReturns struct {
		result1 error
	}
	createReturnsOnCall map[int]struct {
		result1 error
	}
	DeleteStub        func(name string) bool
	deleteMutex       sync.RWMutex
	deleteArgsForCall []struct {
		name string
	}
	deleteReturns struct {
		result1 bool
	}
	deleteReturnsOnCall map[int]struct {
		result1 bool
	}
	FindByIDStub        func(id int) *task.Task
	findByIDMutex       sync.RWMutex
	findByIDArgsForCall []struct {
		id int
	}
	findByIDReturns struct {
		result1 *task.Task
	}
	findByIDReturnsOnCall map[int]struct {
		result1 *task.Task
	}
	FindByNameStub        func(name string) *task.Task
	findByNameMutex       sync.RWMutex
	findByNameArgsForCall []struct {
		name string
	}
	findByNameReturns struct {
		result1 *task.Task
	}
	findByNameReturnsOnCall map[int]struct {
		result1 *task.Task
	}
	TasksStub        func() []*task.Task
	tasksMutex       sync.RWMutex
	tasksArgsForCall []struct{}
	tasksReturns     struct {
		result1 []*task.Task
	}
	tasksReturnsOnCall map[int]struct {
		result1 []*task.Task
	}
	NoteStub        func(name, note string) error
	noteMutex       sync.RWMutex
	noteArgsForCall []struct {
		name string
		note string
	}
	noteReturns struct {
		result1 error
	}
	noteReturnsOnCall map[int]struct {
		result1 error
	}
	SetPriorityStub        func(name string, priority int) error
	setPriorityMutex       sync.RWMutex
	setPriorityArgsForCall []struct {
		name     string
		priority int
	}
	setPriorityReturns struct {
		result1 error
	}
	setPriorityReturnsOnCall map[int]struct {
		result1 error
	}
	SetStateStub        func(name string, state task.State) error
	setStateMutex       sync.RWMutex
	setStateArgsForCall []struct {
		name  string
		state task.State
	}
	setStateReturns struct {
		result1 error
	}
	setStateReturnsOnCall map[int]struct {
		result1 error
	}
	EventsStub        func() []*task.Event
	eventsMutex       sync.RWMutex
	eventsArgsForCall []struct{}
	eventsReturns     struct {
		result1 []*task.Event
	}
	eventsReturnsOnCall map[int]struct {
		result1 []*task.Event
	}
	invocations      map[string][][]interface{}
	invocationsMutex sync.RWMutex
}

func (fake *FakeManager) Create(name string) error {
	fake.createMutex.Lock()
	ret, specificReturn := fake.createReturnsOnCall[len(fake.createArgsForCall)]
	fake.createArgsForCall = append(fake.createArgsForCall, struct {
		name string
	}{name})
	fake.recordInvocation("Create", []interface{}{name})
	fake.createMutex.Unlock()
	if fake.CreateStub != nil {
		return fake.CreateStub(name)
	}
	if specificReturn {
		return ret.result1
	}
	return fake.createReturns.result1
}

func (fake *FakeManager) CreateCallCount() int {
	fake.createMutex.RLock()
	defer fake.createMutex.RUnlock()
	return len(fake.createArgsForCall)
}

func (fake *FakeManager) CreateArgsForCall(i int) string {
	fake.createMutex.RLock()
	defer fake.createMutex.RUnlock()
	return fake.createArgsForCall[i].name
}

func (fake *FakeManager) CreateReturns(result1 error) {
	fake.CreateStub = nil
	fake.createReturns = struct {
		result1 error
	}{result1}
}

func (fake *FakeManager) CreateReturnsOnCall(i int, result1 error) {
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

func (fake *FakeManager) Delete(name string) bool {
	fake.deleteMutex.Lock()
	ret, specificReturn := fake.deleteReturnsOnCall[len(fake.deleteArgsForCall)]
	fake.deleteArgsForCall = append(fake.deleteArgsForCall, struct {
		name string
	}{name})
	fake.recordInvocation("Delete", []interface{}{name})
	fake.deleteMutex.Unlock()
	if fake.DeleteStub != nil {
		return fake.DeleteStub(name)
	}
	if specificReturn {
		return ret.result1
	}
	return fake.deleteReturns.result1
}

func (fake *FakeManager) DeleteCallCount() int {
	fake.deleteMutex.RLock()
	defer fake.deleteMutex.RUnlock()
	return len(fake.deleteArgsForCall)
}

func (fake *FakeManager) DeleteArgsForCall(i int) string {
	fake.deleteMutex.RLock()
	defer fake.deleteMutex.RUnlock()
	return fake.deleteArgsForCall[i].name
}

func (fake *FakeManager) DeleteReturns(result1 bool) {
	fake.DeleteStub = nil
	fake.deleteReturns = struct {
		result1 bool
	}{result1}
}

func (fake *FakeManager) DeleteReturnsOnCall(i int, result1 bool) {
	fake.DeleteStub = nil
	if fake.deleteReturnsOnCall == nil {
		fake.deleteReturnsOnCall = make(map[int]struct {
			result1 bool
		})
	}
	fake.deleteReturnsOnCall[i] = struct {
		result1 bool
	}{result1}
}

func (fake *FakeManager) FindByID(id int) *task.Task {
	fake.findByIDMutex.Lock()
	ret, specificReturn := fake.findByIDReturnsOnCall[len(fake.findByIDArgsForCall)]
	fake.findByIDArgsForCall = append(fake.findByIDArgsForCall, struct {
		id int
	}{id})
	fake.recordInvocation("FindByID", []interface{}{id})
	fake.findByIDMutex.Unlock()
	if fake.FindByIDStub != nil {
		return fake.FindByIDStub(id)
	}
	if specificReturn {
		return ret.result1
	}
	return fake.findByIDReturns.result1
}

func (fake *FakeManager) FindByIDCallCount() int {
	fake.findByIDMutex.RLock()
	defer fake.findByIDMutex.RUnlock()
	return len(fake.findByIDArgsForCall)
}

func (fake *FakeManager) FindByIDArgsForCall(i int) int {
	fake.findByIDMutex.RLock()
	defer fake.findByIDMutex.RUnlock()
	return fake.findByIDArgsForCall[i].id
}

func (fake *FakeManager) FindByIDReturns(result1 *task.Task) {
	fake.FindByIDStub = nil
	fake.findByIDReturns = struct {
		result1 *task.Task
	}{result1}
}

func (fake *FakeManager) FindByIDReturnsOnCall(i int, result1 *task.Task) {
	fake.FindByIDStub = nil
	if fake.findByIDReturnsOnCall == nil {
		fake.findByIDReturnsOnCall = make(map[int]struct {
			result1 *task.Task
		})
	}
	fake.findByIDReturnsOnCall[i] = struct {
		result1 *task.Task
	}{result1}
}

func (fake *FakeManager) FindByName(name string) *task.Task {
	fake.findByNameMutex.Lock()
	ret, specificReturn := fake.findByNameReturnsOnCall[len(fake.findByNameArgsForCall)]
	fake.findByNameArgsForCall = append(fake.findByNameArgsForCall, struct {
		name string
	}{name})
	fake.recordInvocation("FindByName", []interface{}{name})
	fake.findByNameMutex.Unlock()
	if fake.FindByNameStub != nil {
		return fake.FindByNameStub(name)
	}
	if specificReturn {
		return ret.result1
	}
	return fake.findByNameReturns.result1
}

func (fake *FakeManager) FindByNameCallCount() int {
	fake.findByNameMutex.RLock()
	defer fake.findByNameMutex.RUnlock()
	return len(fake.findByNameArgsForCall)
}

func (fake *FakeManager) FindByNameArgsForCall(i int) string {
	fake.findByNameMutex.RLock()
	defer fake.findByNameMutex.RUnlock()
	return fake.findByNameArgsForCall[i].name
}

func (fake *FakeManager) FindByNameReturns(result1 *task.Task) {
	fake.FindByNameStub = nil
	fake.findByNameReturns = struct {
		result1 *task.Task
	}{result1}
}

func (fake *FakeManager) FindByNameReturnsOnCall(i int, result1 *task.Task) {
	fake.FindByNameStub = nil
	if fake.findByNameReturnsOnCall == nil {
		fake.findByNameReturnsOnCall = make(map[int]struct {
			result1 *task.Task
		})
	}
	fake.findByNameReturnsOnCall[i] = struct {
		result1 *task.Task
	}{result1}
}

func (fake *FakeManager) Tasks() []*task.Task {
	fake.tasksMutex.Lock()
	ret, specificReturn := fake.tasksReturnsOnCall[len(fake.tasksArgsForCall)]
	fake.tasksArgsForCall = append(fake.tasksArgsForCall, struct{}{})
	fake.recordInvocation("Tasks", []interface{}{})
	fake.tasksMutex.Unlock()
	if fake.TasksStub != nil {
		return fake.TasksStub()
	}
	if specificReturn {
		return ret.result1
	}
	return fake.tasksReturns.result1
}

func (fake *FakeManager) TasksCallCount() int {
	fake.tasksMutex.RLock()
	defer fake.tasksMutex.RUnlock()
	return len(fake.tasksArgsForCall)
}

func (fake *FakeManager) TasksReturns(result1 []*task.Task) {
	fake.TasksStub = nil
	fake.tasksReturns = struct {
		result1 []*task.Task
	}{result1}
}

func (fake *FakeManager) TasksReturnsOnCall(i int, result1 []*task.Task) {
	fake.TasksStub = nil
	if fake.tasksReturnsOnCall == nil {
		fake.tasksReturnsOnCall = make(map[int]struct {
			result1 []*task.Task
		})
	}
	fake.tasksReturnsOnCall[i] = struct {
		result1 []*task.Task
	}{result1}
}

func (fake *FakeManager) Note(name string, note string) error {
	fake.noteMutex.Lock()
	ret, specificReturn := fake.noteReturnsOnCall[len(fake.noteArgsForCall)]
	fake.noteArgsForCall = append(fake.noteArgsForCall, struct {
		name string
		note string
	}{name, note})
	fake.recordInvocation("Note", []interface{}{name, note})
	fake.noteMutex.Unlock()
	if fake.NoteStub != nil {
		return fake.NoteStub(name, note)
	}
	if specificReturn {
		return ret.result1
	}
	return fake.noteReturns.result1
}

func (fake *FakeManager) NoteCallCount() int {
	fake.noteMutex.RLock()
	defer fake.noteMutex.RUnlock()
	return len(fake.noteArgsForCall)
}

func (fake *FakeManager) NoteArgsForCall(i int) (string, string) {
	fake.noteMutex.RLock()
	defer fake.noteMutex.RUnlock()
	return fake.noteArgsForCall[i].name, fake.noteArgsForCall[i].note
}

func (fake *FakeManager) NoteReturns(result1 error) {
	fake.NoteStub = nil
	fake.noteReturns = struct {
		result1 error
	}{result1}
}

func (fake *FakeManager) NoteReturnsOnCall(i int, result1 error) {
	fake.NoteStub = nil
	if fake.noteReturnsOnCall == nil {
		fake.noteReturnsOnCall = make(map[int]struct {
			result1 error
		})
	}
	fake.noteReturnsOnCall[i] = struct {
		result1 error
	}{result1}
}

func (fake *FakeManager) SetPriority(name string, priority int) error {
	fake.setPriorityMutex.Lock()
	ret, specificReturn := fake.setPriorityReturnsOnCall[len(fake.setPriorityArgsForCall)]
	fake.setPriorityArgsForCall = append(fake.setPriorityArgsForCall, struct {
		name     string
		priority int
	}{name, priority})
	fake.recordInvocation("SetPriority", []interface{}{name, priority})
	fake.setPriorityMutex.Unlock()
	if fake.SetPriorityStub != nil {
		return fake.SetPriorityStub(name, priority)
	}
	if specificReturn {
		return ret.result1
	}
	return fake.setPriorityReturns.result1
}

func (fake *FakeManager) SetPriorityCallCount() int {
	fake.setPriorityMutex.RLock()
	defer fake.setPriorityMutex.RUnlock()
	return len(fake.setPriorityArgsForCall)
}

func (fake *FakeManager) SetPriorityArgsForCall(i int) (string, int) {
	fake.setPriorityMutex.RLock()
	defer fake.setPriorityMutex.RUnlock()
	return fake.setPriorityArgsForCall[i].name, fake.setPriorityArgsForCall[i].priority
}

func (fake *FakeManager) SetPriorityReturns(result1 error) {
	fake.SetPriorityStub = nil
	fake.setPriorityReturns = struct {
		result1 error
	}{result1}
}

func (fake *FakeManager) SetPriorityReturnsOnCall(i int, result1 error) {
	fake.SetPriorityStub = nil
	if fake.setPriorityReturnsOnCall == nil {
		fake.setPriorityReturnsOnCall = make(map[int]struct {
			result1 error
		})
	}
	fake.setPriorityReturnsOnCall[i] = struct {
		result1 error
	}{result1}
}

func (fake *FakeManager) SetState(name string, state task.State) error {
	fake.setStateMutex.Lock()
	ret, specificReturn := fake.setStateReturnsOnCall[len(fake.setStateArgsForCall)]
	fake.setStateArgsForCall = append(fake.setStateArgsForCall, struct {
		name  string
		state task.State
	}{name, state})
	fake.recordInvocation("SetState", []interface{}{name, state})
	fake.setStateMutex.Unlock()
	if fake.SetStateStub != nil {
		return fake.SetStateStub(name, state)
	}
	if specificReturn {
		return ret.result1
	}
	return fake.setStateReturns.result1
}

func (fake *FakeManager) SetStateCallCount() int {
	fake.setStateMutex.RLock()
	defer fake.setStateMutex.RUnlock()
	return len(fake.setStateArgsForCall)
}

func (fake *FakeManager) SetStateArgsForCall(i int) (string, task.State) {
	fake.setStateMutex.RLock()
	defer fake.setStateMutex.RUnlock()
	return fake.setStateArgsForCall[i].name, fake.setStateArgsForCall[i].state
}

func (fake *FakeManager) SetStateReturns(result1 error) {
	fake.SetStateStub = nil
	fake.setStateReturns = struct {
		result1 error
	}{result1}
}

func (fake *FakeManager) SetStateReturnsOnCall(i int, result1 error) {
	fake.SetStateStub = nil
	if fake.setStateReturnsOnCall == nil {
		fake.setStateReturnsOnCall = make(map[int]struct {
			result1 error
		})
	}
	fake.setStateReturnsOnCall[i] = struct {
		result1 error
	}{result1}
}

func (fake *FakeManager) Events() []*task.Event {
	fake.eventsMutex.Lock()
	ret, specificReturn := fake.eventsReturnsOnCall[len(fake.eventsArgsForCall)]
	fake.eventsArgsForCall = append(fake.eventsArgsForCall, struct{}{})
	fake.recordInvocation("Events", []interface{}{})
	fake.eventsMutex.Unlock()
	if fake.EventsStub != nil {
		return fake.EventsStub()
	}
	if specificReturn {
		return ret.result1
	}
	return fake.eventsReturns.result1
}

func (fake *FakeManager) EventsCallCount() int {
	fake.eventsMutex.RLock()
	defer fake.eventsMutex.RUnlock()
	return len(fake.eventsArgsForCall)
}

func (fake *FakeManager) EventsReturns(result1 []*task.Event) {
	fake.EventsStub = nil
	fake.eventsReturns = struct {
		result1 []*task.Event
	}{result1}
}

func (fake *FakeManager) EventsReturnsOnCall(i int, result1 []*task.Event) {
	fake.EventsStub = nil
	if fake.eventsReturnsOnCall == nil {
		fake.eventsReturnsOnCall = make(map[int]struct {
			result1 []*task.Event
		})
	}
	fake.eventsReturnsOnCall[i] = struct {
		result1 []*task.Event
	}{result1}
}

func (fake *FakeManager) Invocations() map[string][][]interface{} {
	fake.invocationsMutex.RLock()
	defer fake.invocationsMutex.RUnlock()
	fake.createMutex.RLock()
	defer fake.createMutex.RUnlock()
	fake.deleteMutex.RLock()
	defer fake.deleteMutex.RUnlock()
	fake.findByIDMutex.RLock()
	defer fake.findByIDMutex.RUnlock()
	fake.findByNameMutex.RLock()
	defer fake.findByNameMutex.RUnlock()
	fake.tasksMutex.RLock()
	defer fake.tasksMutex.RUnlock()
	fake.noteMutex.RLock()
	defer fake.noteMutex.RUnlock()
	fake.setPriorityMutex.RLock()
	defer fake.setPriorityMutex.RUnlock()
	fake.setStateMutex.RLock()
	defer fake.setStateMutex.RUnlock()
	fake.eventsMutex.RLock()
	defer fake.eventsMutex.RUnlock()
	copiedInvocations := map[string][][]interface{}{}
	for key, value := range fake.invocations {
		copiedInvocations[key] = value
	}
	return copiedInvocations
}

func (fake *FakeManager) recordInvocation(key string, args []interface{}) {
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

var _ task.Manager = new(FakeManager)
