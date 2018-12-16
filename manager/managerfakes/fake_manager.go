// Code generated by counterfeiter. DO NOT EDIT.
package managerfakes

import (
	sync "sync"

	manager "github.com/ankeesler/anwork/manager"
	task2 "github.com/ankeesler/anwork/task2"
)

type FakeManager struct {
	CreateStub        func(string) error
	createMutex       sync.RWMutex
	createArgsForCall []struct {
		arg1 string
	}
	createReturns struct {
		result1 error
	}
	createReturnsOnCall map[int]struct {
		result1 error
	}
	DeleteStub        func(string) error
	deleteMutex       sync.RWMutex
	deleteArgsForCall []struct {
		arg1 string
	}
	deleteReturns struct {
		result1 error
	}
	deleteReturnsOnCall map[int]struct {
		result1 error
	}
	EventsStub        func() ([]*task2.Event, error)
	eventsMutex       sync.RWMutex
	eventsArgsForCall []struct {
	}
	eventsReturns struct {
		result1 []*task2.Event
		result2 error
	}
	eventsReturnsOnCall map[int]struct {
		result1 []*task2.Event
		result2 error
	}
	FindByIDStub        func(int) (*task2.Task, error)
	findByIDMutex       sync.RWMutex
	findByIDArgsForCall []struct {
		arg1 int
	}
	findByIDReturns struct {
		result1 *task2.Task
		result2 error
	}
	findByIDReturnsOnCall map[int]struct {
		result1 *task2.Task
		result2 error
	}
	FindByNameStub        func(string) (*task2.Task, error)
	findByNameMutex       sync.RWMutex
	findByNameArgsForCall []struct {
		arg1 string
	}
	findByNameReturns struct {
		result1 *task2.Task
		result2 error
	}
	findByNameReturnsOnCall map[int]struct {
		result1 *task2.Task
		result2 error
	}
	NoteStub        func(string, string) error
	noteMutex       sync.RWMutex
	noteArgsForCall []struct {
		arg1 string
		arg2 string
	}
	noteReturns struct {
		result1 error
	}
	noteReturnsOnCall map[int]struct {
		result1 error
	}
	RenameStub        func(string, string) error
	renameMutex       sync.RWMutex
	renameArgsForCall []struct {
		arg1 string
		arg2 string
	}
	renameReturns struct {
		result1 error
	}
	renameReturnsOnCall map[int]struct {
		result1 error
	}
	ResetStub        func() error
	resetMutex       sync.RWMutex
	resetArgsForCall []struct {
	}
	resetReturns struct {
		result1 error
	}
	resetReturnsOnCall map[int]struct {
		result1 error
	}
	SetPriorityStub        func(string, int) error
	setPriorityMutex       sync.RWMutex
	setPriorityArgsForCall []struct {
		arg1 string
		arg2 int
	}
	setPriorityReturns struct {
		result1 error
	}
	setPriorityReturnsOnCall map[int]struct {
		result1 error
	}
	SetStateStub        func(string, task2.State) error
	setStateMutex       sync.RWMutex
	setStateArgsForCall []struct {
		arg1 string
		arg2 task2.State
	}
	setStateReturns struct {
		result1 error
	}
	setStateReturnsOnCall map[int]struct {
		result1 error
	}
	TasksStub        func() ([]*task2.Task, error)
	tasksMutex       sync.RWMutex
	tasksArgsForCall []struct {
	}
	tasksReturns struct {
		result1 []*task2.Task
		result2 error
	}
	tasksReturnsOnCall map[int]struct {
		result1 []*task2.Task
		result2 error
	}
	invocations      map[string][][]interface{}
	invocationsMutex sync.RWMutex
}

func (fake *FakeManager) Create(arg1 string) error {
	fake.createMutex.Lock()
	ret, specificReturn := fake.createReturnsOnCall[len(fake.createArgsForCall)]
	fake.createArgsForCall = append(fake.createArgsForCall, struct {
		arg1 string
	}{arg1})
	fake.recordInvocation("Create", []interface{}{arg1})
	fake.createMutex.Unlock()
	if fake.CreateStub != nil {
		return fake.CreateStub(arg1)
	}
	if specificReturn {
		return ret.result1
	}
	fakeReturns := fake.createReturns
	return fakeReturns.result1
}

func (fake *FakeManager) CreateCallCount() int {
	fake.createMutex.RLock()
	defer fake.createMutex.RUnlock()
	return len(fake.createArgsForCall)
}

func (fake *FakeManager) CreateCalls(stub func(string) error) {
	fake.createMutex.Lock()
	defer fake.createMutex.Unlock()
	fake.CreateStub = stub
}

func (fake *FakeManager) CreateArgsForCall(i int) string {
	fake.createMutex.RLock()
	defer fake.createMutex.RUnlock()
	argsForCall := fake.createArgsForCall[i]
	return argsForCall.arg1
}

func (fake *FakeManager) CreateReturns(result1 error) {
	fake.createMutex.Lock()
	defer fake.createMutex.Unlock()
	fake.CreateStub = nil
	fake.createReturns = struct {
		result1 error
	}{result1}
}

func (fake *FakeManager) CreateReturnsOnCall(i int, result1 error) {
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

func (fake *FakeManager) Delete(arg1 string) error {
	fake.deleteMutex.Lock()
	ret, specificReturn := fake.deleteReturnsOnCall[len(fake.deleteArgsForCall)]
	fake.deleteArgsForCall = append(fake.deleteArgsForCall, struct {
		arg1 string
	}{arg1})
	fake.recordInvocation("Delete", []interface{}{arg1})
	fake.deleteMutex.Unlock()
	if fake.DeleteStub != nil {
		return fake.DeleteStub(arg1)
	}
	if specificReturn {
		return ret.result1
	}
	fakeReturns := fake.deleteReturns
	return fakeReturns.result1
}

func (fake *FakeManager) DeleteCallCount() int {
	fake.deleteMutex.RLock()
	defer fake.deleteMutex.RUnlock()
	return len(fake.deleteArgsForCall)
}

func (fake *FakeManager) DeleteCalls(stub func(string) error) {
	fake.deleteMutex.Lock()
	defer fake.deleteMutex.Unlock()
	fake.DeleteStub = stub
}

func (fake *FakeManager) DeleteArgsForCall(i int) string {
	fake.deleteMutex.RLock()
	defer fake.deleteMutex.RUnlock()
	argsForCall := fake.deleteArgsForCall[i]
	return argsForCall.arg1
}

func (fake *FakeManager) DeleteReturns(result1 error) {
	fake.deleteMutex.Lock()
	defer fake.deleteMutex.Unlock()
	fake.DeleteStub = nil
	fake.deleteReturns = struct {
		result1 error
	}{result1}
}

func (fake *FakeManager) DeleteReturnsOnCall(i int, result1 error) {
	fake.deleteMutex.Lock()
	defer fake.deleteMutex.Unlock()
	fake.DeleteStub = nil
	if fake.deleteReturnsOnCall == nil {
		fake.deleteReturnsOnCall = make(map[int]struct {
			result1 error
		})
	}
	fake.deleteReturnsOnCall[i] = struct {
		result1 error
	}{result1}
}

func (fake *FakeManager) Events() ([]*task2.Event, error) {
	fake.eventsMutex.Lock()
	ret, specificReturn := fake.eventsReturnsOnCall[len(fake.eventsArgsForCall)]
	fake.eventsArgsForCall = append(fake.eventsArgsForCall, struct {
	}{})
	fake.recordInvocation("Events", []interface{}{})
	fake.eventsMutex.Unlock()
	if fake.EventsStub != nil {
		return fake.EventsStub()
	}
	if specificReturn {
		return ret.result1, ret.result2
	}
	fakeReturns := fake.eventsReturns
	return fakeReturns.result1, fakeReturns.result2
}

func (fake *FakeManager) EventsCallCount() int {
	fake.eventsMutex.RLock()
	defer fake.eventsMutex.RUnlock()
	return len(fake.eventsArgsForCall)
}

func (fake *FakeManager) EventsCalls(stub func() ([]*task2.Event, error)) {
	fake.eventsMutex.Lock()
	defer fake.eventsMutex.Unlock()
	fake.EventsStub = stub
}

func (fake *FakeManager) EventsReturns(result1 []*task2.Event, result2 error) {
	fake.eventsMutex.Lock()
	defer fake.eventsMutex.Unlock()
	fake.EventsStub = nil
	fake.eventsReturns = struct {
		result1 []*task2.Event
		result2 error
	}{result1, result2}
}

func (fake *FakeManager) EventsReturnsOnCall(i int, result1 []*task2.Event, result2 error) {
	fake.eventsMutex.Lock()
	defer fake.eventsMutex.Unlock()
	fake.EventsStub = nil
	if fake.eventsReturnsOnCall == nil {
		fake.eventsReturnsOnCall = make(map[int]struct {
			result1 []*task2.Event
			result2 error
		})
	}
	fake.eventsReturnsOnCall[i] = struct {
		result1 []*task2.Event
		result2 error
	}{result1, result2}
}

func (fake *FakeManager) FindByID(arg1 int) (*task2.Task, error) {
	fake.findByIDMutex.Lock()
	ret, specificReturn := fake.findByIDReturnsOnCall[len(fake.findByIDArgsForCall)]
	fake.findByIDArgsForCall = append(fake.findByIDArgsForCall, struct {
		arg1 int
	}{arg1})
	fake.recordInvocation("FindByID", []interface{}{arg1})
	fake.findByIDMutex.Unlock()
	if fake.FindByIDStub != nil {
		return fake.FindByIDStub(arg1)
	}
	if specificReturn {
		return ret.result1, ret.result2
	}
	fakeReturns := fake.findByIDReturns
	return fakeReturns.result1, fakeReturns.result2
}

func (fake *FakeManager) FindByIDCallCount() int {
	fake.findByIDMutex.RLock()
	defer fake.findByIDMutex.RUnlock()
	return len(fake.findByIDArgsForCall)
}

func (fake *FakeManager) FindByIDCalls(stub func(int) (*task2.Task, error)) {
	fake.findByIDMutex.Lock()
	defer fake.findByIDMutex.Unlock()
	fake.FindByIDStub = stub
}

func (fake *FakeManager) FindByIDArgsForCall(i int) int {
	fake.findByIDMutex.RLock()
	defer fake.findByIDMutex.RUnlock()
	argsForCall := fake.findByIDArgsForCall[i]
	return argsForCall.arg1
}

func (fake *FakeManager) FindByIDReturns(result1 *task2.Task, result2 error) {
	fake.findByIDMutex.Lock()
	defer fake.findByIDMutex.Unlock()
	fake.FindByIDStub = nil
	fake.findByIDReturns = struct {
		result1 *task2.Task
		result2 error
	}{result1, result2}
}

func (fake *FakeManager) FindByIDReturnsOnCall(i int, result1 *task2.Task, result2 error) {
	fake.findByIDMutex.Lock()
	defer fake.findByIDMutex.Unlock()
	fake.FindByIDStub = nil
	if fake.findByIDReturnsOnCall == nil {
		fake.findByIDReturnsOnCall = make(map[int]struct {
			result1 *task2.Task
			result2 error
		})
	}
	fake.findByIDReturnsOnCall[i] = struct {
		result1 *task2.Task
		result2 error
	}{result1, result2}
}

func (fake *FakeManager) FindByName(arg1 string) (*task2.Task, error) {
	fake.findByNameMutex.Lock()
	ret, specificReturn := fake.findByNameReturnsOnCall[len(fake.findByNameArgsForCall)]
	fake.findByNameArgsForCall = append(fake.findByNameArgsForCall, struct {
		arg1 string
	}{arg1})
	fake.recordInvocation("FindByName", []interface{}{arg1})
	fake.findByNameMutex.Unlock()
	if fake.FindByNameStub != nil {
		return fake.FindByNameStub(arg1)
	}
	if specificReturn {
		return ret.result1, ret.result2
	}
	fakeReturns := fake.findByNameReturns
	return fakeReturns.result1, fakeReturns.result2
}

func (fake *FakeManager) FindByNameCallCount() int {
	fake.findByNameMutex.RLock()
	defer fake.findByNameMutex.RUnlock()
	return len(fake.findByNameArgsForCall)
}

func (fake *FakeManager) FindByNameCalls(stub func(string) (*task2.Task, error)) {
	fake.findByNameMutex.Lock()
	defer fake.findByNameMutex.Unlock()
	fake.FindByNameStub = stub
}

func (fake *FakeManager) FindByNameArgsForCall(i int) string {
	fake.findByNameMutex.RLock()
	defer fake.findByNameMutex.RUnlock()
	argsForCall := fake.findByNameArgsForCall[i]
	return argsForCall.arg1
}

func (fake *FakeManager) FindByNameReturns(result1 *task2.Task, result2 error) {
	fake.findByNameMutex.Lock()
	defer fake.findByNameMutex.Unlock()
	fake.FindByNameStub = nil
	fake.findByNameReturns = struct {
		result1 *task2.Task
		result2 error
	}{result1, result2}
}

func (fake *FakeManager) FindByNameReturnsOnCall(i int, result1 *task2.Task, result2 error) {
	fake.findByNameMutex.Lock()
	defer fake.findByNameMutex.Unlock()
	fake.FindByNameStub = nil
	if fake.findByNameReturnsOnCall == nil {
		fake.findByNameReturnsOnCall = make(map[int]struct {
			result1 *task2.Task
			result2 error
		})
	}
	fake.findByNameReturnsOnCall[i] = struct {
		result1 *task2.Task
		result2 error
	}{result1, result2}
}

func (fake *FakeManager) Note(arg1 string, arg2 string) error {
	fake.noteMutex.Lock()
	ret, specificReturn := fake.noteReturnsOnCall[len(fake.noteArgsForCall)]
	fake.noteArgsForCall = append(fake.noteArgsForCall, struct {
		arg1 string
		arg2 string
	}{arg1, arg2})
	fake.recordInvocation("Note", []interface{}{arg1, arg2})
	fake.noteMutex.Unlock()
	if fake.NoteStub != nil {
		return fake.NoteStub(arg1, arg2)
	}
	if specificReturn {
		return ret.result1
	}
	fakeReturns := fake.noteReturns
	return fakeReturns.result1
}

func (fake *FakeManager) NoteCallCount() int {
	fake.noteMutex.RLock()
	defer fake.noteMutex.RUnlock()
	return len(fake.noteArgsForCall)
}

func (fake *FakeManager) NoteCalls(stub func(string, string) error) {
	fake.noteMutex.Lock()
	defer fake.noteMutex.Unlock()
	fake.NoteStub = stub
}

func (fake *FakeManager) NoteArgsForCall(i int) (string, string) {
	fake.noteMutex.RLock()
	defer fake.noteMutex.RUnlock()
	argsForCall := fake.noteArgsForCall[i]
	return argsForCall.arg1, argsForCall.arg2
}

func (fake *FakeManager) NoteReturns(result1 error) {
	fake.noteMutex.Lock()
	defer fake.noteMutex.Unlock()
	fake.NoteStub = nil
	fake.noteReturns = struct {
		result1 error
	}{result1}
}

func (fake *FakeManager) NoteReturnsOnCall(i int, result1 error) {
	fake.noteMutex.Lock()
	defer fake.noteMutex.Unlock()
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

func (fake *FakeManager) Rename(arg1 string, arg2 string) error {
	fake.renameMutex.Lock()
	ret, specificReturn := fake.renameReturnsOnCall[len(fake.renameArgsForCall)]
	fake.renameArgsForCall = append(fake.renameArgsForCall, struct {
		arg1 string
		arg2 string
	}{arg1, arg2})
	fake.recordInvocation("Rename", []interface{}{arg1, arg2})
	fake.renameMutex.Unlock()
	if fake.RenameStub != nil {
		return fake.RenameStub(arg1, arg2)
	}
	if specificReturn {
		return ret.result1
	}
	fakeReturns := fake.renameReturns
	return fakeReturns.result1
}

func (fake *FakeManager) RenameCallCount() int {
	fake.renameMutex.RLock()
	defer fake.renameMutex.RUnlock()
	return len(fake.renameArgsForCall)
}

func (fake *FakeManager) RenameCalls(stub func(string, string) error) {
	fake.renameMutex.Lock()
	defer fake.renameMutex.Unlock()
	fake.RenameStub = stub
}

func (fake *FakeManager) RenameArgsForCall(i int) (string, string) {
	fake.renameMutex.RLock()
	defer fake.renameMutex.RUnlock()
	argsForCall := fake.renameArgsForCall[i]
	return argsForCall.arg1, argsForCall.arg2
}

func (fake *FakeManager) RenameReturns(result1 error) {
	fake.renameMutex.Lock()
	defer fake.renameMutex.Unlock()
	fake.RenameStub = nil
	fake.renameReturns = struct {
		result1 error
	}{result1}
}

func (fake *FakeManager) RenameReturnsOnCall(i int, result1 error) {
	fake.renameMutex.Lock()
	defer fake.renameMutex.Unlock()
	fake.RenameStub = nil
	if fake.renameReturnsOnCall == nil {
		fake.renameReturnsOnCall = make(map[int]struct {
			result1 error
		})
	}
	fake.renameReturnsOnCall[i] = struct {
		result1 error
	}{result1}
}

func (fake *FakeManager) Reset() error {
	fake.resetMutex.Lock()
	ret, specificReturn := fake.resetReturnsOnCall[len(fake.resetArgsForCall)]
	fake.resetArgsForCall = append(fake.resetArgsForCall, struct {
	}{})
	fake.recordInvocation("Reset", []interface{}{})
	fake.resetMutex.Unlock()
	if fake.ResetStub != nil {
		return fake.ResetStub()
	}
	if specificReturn {
		return ret.result1
	}
	fakeReturns := fake.resetReturns
	return fakeReturns.result1
}

func (fake *FakeManager) ResetCallCount() int {
	fake.resetMutex.RLock()
	defer fake.resetMutex.RUnlock()
	return len(fake.resetArgsForCall)
}

func (fake *FakeManager) ResetCalls(stub func() error) {
	fake.resetMutex.Lock()
	defer fake.resetMutex.Unlock()
	fake.ResetStub = stub
}

func (fake *FakeManager) ResetReturns(result1 error) {
	fake.resetMutex.Lock()
	defer fake.resetMutex.Unlock()
	fake.ResetStub = nil
	fake.resetReturns = struct {
		result1 error
	}{result1}
}

func (fake *FakeManager) ResetReturnsOnCall(i int, result1 error) {
	fake.resetMutex.Lock()
	defer fake.resetMutex.Unlock()
	fake.ResetStub = nil
	if fake.resetReturnsOnCall == nil {
		fake.resetReturnsOnCall = make(map[int]struct {
			result1 error
		})
	}
	fake.resetReturnsOnCall[i] = struct {
		result1 error
	}{result1}
}

func (fake *FakeManager) SetPriority(arg1 string, arg2 int) error {
	fake.setPriorityMutex.Lock()
	ret, specificReturn := fake.setPriorityReturnsOnCall[len(fake.setPriorityArgsForCall)]
	fake.setPriorityArgsForCall = append(fake.setPriorityArgsForCall, struct {
		arg1 string
		arg2 int
	}{arg1, arg2})
	fake.recordInvocation("SetPriority", []interface{}{arg1, arg2})
	fake.setPriorityMutex.Unlock()
	if fake.SetPriorityStub != nil {
		return fake.SetPriorityStub(arg1, arg2)
	}
	if specificReturn {
		return ret.result1
	}
	fakeReturns := fake.setPriorityReturns
	return fakeReturns.result1
}

func (fake *FakeManager) SetPriorityCallCount() int {
	fake.setPriorityMutex.RLock()
	defer fake.setPriorityMutex.RUnlock()
	return len(fake.setPriorityArgsForCall)
}

func (fake *FakeManager) SetPriorityCalls(stub func(string, int) error) {
	fake.setPriorityMutex.Lock()
	defer fake.setPriorityMutex.Unlock()
	fake.SetPriorityStub = stub
}

func (fake *FakeManager) SetPriorityArgsForCall(i int) (string, int) {
	fake.setPriorityMutex.RLock()
	defer fake.setPriorityMutex.RUnlock()
	argsForCall := fake.setPriorityArgsForCall[i]
	return argsForCall.arg1, argsForCall.arg2
}

func (fake *FakeManager) SetPriorityReturns(result1 error) {
	fake.setPriorityMutex.Lock()
	defer fake.setPriorityMutex.Unlock()
	fake.SetPriorityStub = nil
	fake.setPriorityReturns = struct {
		result1 error
	}{result1}
}

func (fake *FakeManager) SetPriorityReturnsOnCall(i int, result1 error) {
	fake.setPriorityMutex.Lock()
	defer fake.setPriorityMutex.Unlock()
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

func (fake *FakeManager) SetState(arg1 string, arg2 task2.State) error {
	fake.setStateMutex.Lock()
	ret, specificReturn := fake.setStateReturnsOnCall[len(fake.setStateArgsForCall)]
	fake.setStateArgsForCall = append(fake.setStateArgsForCall, struct {
		arg1 string
		arg2 task2.State
	}{arg1, arg2})
	fake.recordInvocation("SetState", []interface{}{arg1, arg2})
	fake.setStateMutex.Unlock()
	if fake.SetStateStub != nil {
		return fake.SetStateStub(arg1, arg2)
	}
	if specificReturn {
		return ret.result1
	}
	fakeReturns := fake.setStateReturns
	return fakeReturns.result1
}

func (fake *FakeManager) SetStateCallCount() int {
	fake.setStateMutex.RLock()
	defer fake.setStateMutex.RUnlock()
	return len(fake.setStateArgsForCall)
}

func (fake *FakeManager) SetStateCalls(stub func(string, task2.State) error) {
	fake.setStateMutex.Lock()
	defer fake.setStateMutex.Unlock()
	fake.SetStateStub = stub
}

func (fake *FakeManager) SetStateArgsForCall(i int) (string, task2.State) {
	fake.setStateMutex.RLock()
	defer fake.setStateMutex.RUnlock()
	argsForCall := fake.setStateArgsForCall[i]
	return argsForCall.arg1, argsForCall.arg2
}

func (fake *FakeManager) SetStateReturns(result1 error) {
	fake.setStateMutex.Lock()
	defer fake.setStateMutex.Unlock()
	fake.SetStateStub = nil
	fake.setStateReturns = struct {
		result1 error
	}{result1}
}

func (fake *FakeManager) SetStateReturnsOnCall(i int, result1 error) {
	fake.setStateMutex.Lock()
	defer fake.setStateMutex.Unlock()
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

func (fake *FakeManager) Tasks() ([]*task2.Task, error) {
	fake.tasksMutex.Lock()
	ret, specificReturn := fake.tasksReturnsOnCall[len(fake.tasksArgsForCall)]
	fake.tasksArgsForCall = append(fake.tasksArgsForCall, struct {
	}{})
	fake.recordInvocation("Tasks", []interface{}{})
	fake.tasksMutex.Unlock()
	if fake.TasksStub != nil {
		return fake.TasksStub()
	}
	if specificReturn {
		return ret.result1, ret.result2
	}
	fakeReturns := fake.tasksReturns
	return fakeReturns.result1, fakeReturns.result2
}

func (fake *FakeManager) TasksCallCount() int {
	fake.tasksMutex.RLock()
	defer fake.tasksMutex.RUnlock()
	return len(fake.tasksArgsForCall)
}

func (fake *FakeManager) TasksCalls(stub func() ([]*task2.Task, error)) {
	fake.tasksMutex.Lock()
	defer fake.tasksMutex.Unlock()
	fake.TasksStub = stub
}

func (fake *FakeManager) TasksReturns(result1 []*task2.Task, result2 error) {
	fake.tasksMutex.Lock()
	defer fake.tasksMutex.Unlock()
	fake.TasksStub = nil
	fake.tasksReturns = struct {
		result1 []*task2.Task
		result2 error
	}{result1, result2}
}

func (fake *FakeManager) TasksReturnsOnCall(i int, result1 []*task2.Task, result2 error) {
	fake.tasksMutex.Lock()
	defer fake.tasksMutex.Unlock()
	fake.TasksStub = nil
	if fake.tasksReturnsOnCall == nil {
		fake.tasksReturnsOnCall = make(map[int]struct {
			result1 []*task2.Task
			result2 error
		})
	}
	fake.tasksReturnsOnCall[i] = struct {
		result1 []*task2.Task
		result2 error
	}{result1, result2}
}

func (fake *FakeManager) Invocations() map[string][][]interface{} {
	fake.invocationsMutex.RLock()
	defer fake.invocationsMutex.RUnlock()
	fake.createMutex.RLock()
	defer fake.createMutex.RUnlock()
	fake.deleteMutex.RLock()
	defer fake.deleteMutex.RUnlock()
	fake.eventsMutex.RLock()
	defer fake.eventsMutex.RUnlock()
	fake.findByIDMutex.RLock()
	defer fake.findByIDMutex.RUnlock()
	fake.findByNameMutex.RLock()
	defer fake.findByNameMutex.RUnlock()
	fake.noteMutex.RLock()
	defer fake.noteMutex.RUnlock()
	fake.renameMutex.RLock()
	defer fake.renameMutex.RUnlock()
	fake.resetMutex.RLock()
	defer fake.resetMutex.RUnlock()
	fake.setPriorityMutex.RLock()
	defer fake.setPriorityMutex.RUnlock()
	fake.setStateMutex.RLock()
	defer fake.setStateMutex.RUnlock()
	fake.tasksMutex.RLock()
	defer fake.tasksMutex.RUnlock()
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

var _ manager.Manager = new(FakeManager)
