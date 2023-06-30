package fsm

import (
	"sync"

	"golang.org/x/exp/constraints"
)

type IFsm[E constraints.Ordered, S constraints.Ordered] interface {
	// Clone the FSM without Mutex.
	Clone() IFsm[E, S]
	// CloneWithState clone the FSM with new state without Mutex.
	CloneWithState(newState S) IFsm[E, S]
	// Current returns the current state.
	Current() S
	// Is returns true if state is the current state.
	Is(state S) bool
	// SetState allows the user to move to the given state from current state.
	SetState(state S)
	// Trigger call a state transition with the named event.
	// It will return nil if src state change to dst state success or one of these errors:
	// - ErrInappropriateEvent: event X inappropriate in the state Y
	// - ErrNonExistEvent: event X does not exist
	Trigger(event E) error
	// ShouldTrigger return dst state transition with the named event and src state.
	// It will return if src state change to dst state success or holds the same as the current state:
	ShouldTrigger(event E)
	// IsCan returns true if event can occur in the current state.
	IsCan(event E) bool
	// IsAllCan returns true if all the events can occur in src state.
	IsAllCan(event ...E) bool
	// AvailTransitionEvent returns a list of available transition event in the
	// current state.
	AvailTransitionEvent() []E
	// ContainEvent returns true if support the event.
	ContainEvent(event E) bool
	// ContainAllEvent returns true if support all event.
	ContainAllEvent(events ...E) bool
	// Trans return Trans
	Trans() map[eKey[E, S]]S
}

var _ IFsm[string, string] = (*FSM[string, string])(nil)
var _ IFsm[int, string] = (*FSM[int, string])(nil)

// FSM is the state machine that holds the current state.
// E ist the event
// S is the state
type FSM[E constraints.Ordered, S constraints.Ordered] struct {
	// Translator contain events and source states to destination states.
	// This is immutable
	Translator[E, S]
	// mu guards access to the current state.
	mu sync.RWMutex
	// current is the state that the FSM is currently in.
	current S
}

// New constructs a generic FSM with an initial state S, for events E.
// E is the event type
// S is the state type.
func New[E constraints.Ordered, S constraints.Ordered](initState S, transitions []Transform[E, S]) IFsm[E, S] {
	return NewFromTransitions[E, S](initState, NewTranslator[E, S](transitions))
}

// NewFromTransitions constructs a generic FSM with an initial state S, for events E.
func NewFromTransitions[E constraints.Ordered, S constraints.Ordered](initState S, ts Translator[E, S]) IFsm[E, S] {
	return &FSM[E, S]{
		current:    initState,
		Translator: ts,
	}
}

// Clone the FSM without Mutex.
func (f *FSM[E, S]) Clone() IFsm[E, S] {
	return &FSM[E, S]{
		current:    f.current,
		Translator: f.Translator,
	}
}
func (f *FSM[E, S]) CloneWithState(newState S) IFsm[E, S] {
	return &FSM[E, S]{
		current:    newState,
		Translator: f.Translator,
	}
}
func (f *FSM[E, S]) Current() S {
	f.mu.RLock()
	defer f.mu.RUnlock()
	return f.current
}
func (f *FSM[E, S]) Is(state S) bool {
	f.mu.RLock()
	defer f.mu.RUnlock()
	return state == f.current
}
func (f *FSM[E, S]) SetState(state S) {
	f.mu.Lock()
	defer f.mu.Unlock()
	f.current = state
}
func (f *FSM[E, S]) Trigger(event E) error {
	f.mu.Lock()
	defer f.mu.Unlock()
	dst, err := f.Translator.Trigger(f.current, event)
	if err != nil {
		return err
	}
	f.current = dst
	return nil
}
func (f *FSM[E, S]) ShouldTrigger(event E) {
	f.mu.Lock()
	defer f.mu.Unlock()
	f.current = f.Translator.ShouldTrigger(f.current, event)
}
func (f *FSM[E, S]) IsCan(event E) bool {
	return f.Translator.IsCan(f.current, event)
}
func (f *FSM[E, S]) IsAllCan(event ...E) bool {
	return f.Translator.IsAllCan(f.current, event...)
}
func (f *FSM[E, S]) AvailTransitionEvent() []E {
	return f.Translator.AvailTransitionEvent(f.current)
}

var _ IFsm[string, string] = (*FSMUnsafe[string, string])(nil)
var _ IFsm[int, string] = (*FSMUnsafe[int, string])(nil)

// FSMUnsafe is the state machine that holds the current state.
// E ist the event
// S is the state
type FSMUnsafe[E constraints.Ordered, S constraints.Ordered] struct {
	// translator contain events and source states to destination states.
	// This is immutable
	Translator[E, S]
	// current is the state that the FSMUnsafe is currently in.
	current S
}

// NewUnsafeFSM constructs a generic FSMUnsafe with an initial state S, for events E.
// E is the event type
// S is the state type.
func NewUnsafeFSM[E constraints.Ordered, S constraints.Ordered](initState S, transitions []Transform[E, S]) IFsm[E, S] {
	return NewUnsafeFSMFromTransitions(initState, NewTranslator[E, S](transitions))
}

// NewUnsafeFSMFromTransitions constructs a generic FSMUnsafe with an initial state S, for events E.
func NewUnsafeFSMFromTransitions[E constraints.Ordered, S constraints.Ordered](initState S, ts Translator[E, S]) IFsm[E, S] {
	return &FSMUnsafe[E, S]{
		current:    initState,
		Translator: ts,
	}
}
func (f *FSMUnsafe[E, S]) Clone() IFsm[E, S] {
	return &FSMUnsafe[E, S]{
		current:    f.current,
		Translator: f.Translator,
	}
}
func (f *FSMUnsafe[E, S]) CloneWithState(newState S) IFsm[E, S] {
	return &FSMUnsafe[E, S]{
		current:    newState,
		Translator: f.Translator,
	}
}
func (f *FSMUnsafe[E, S]) Current() S {
	return f.current
}
func (f *FSMUnsafe[E, S]) Is(state S) bool {
	return state == f.current
}
func (f *FSMUnsafe[E, S]) SetState(state S) {
	f.current = state
}
func (f *FSMUnsafe[E, S]) Trigger(event E) error {
	dst, err := f.Translator.Trigger(f.current, event)
	if err != nil {
		return err
	}
	f.current = dst
	return nil
}
func (f *FSMUnsafe[E, S]) ShouldTrigger(event E) {
	f.current = f.Translator.ShouldTrigger(f.current, event)
}
func (f *FSMUnsafe[E, S]) IsCan(event E) bool {
	return f.Translator.IsCan(f.current, event)
}
func (f *FSMUnsafe[E, S]) IsAllCan(event ...E) bool {
	return f.Translator.IsAllCan(f.current, event...)
}
func (f *FSMUnsafe[E, S]) AvailTransitionEvent() []E {
	return f.Translator.AvailTransitionEvent(f.current)
}
