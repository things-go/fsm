package fsm

import (
	"sync"

	"golang.org/x/exp/constraints"
)

// FSM is the state machine that holds the current state.
// E ist the event
// S is the state
type FSM[E constraints.Ordered, S constraints.Ordered] struct {
	// current is the state that the FSM is currently in.
	current S

	// Translator contain events and source states to destination states.
	// This is immutable
	Translator[E, S]

	// stateMu guards access to the current state.
	stateMu sync.RWMutex
}

// New constructs a generic FSM with an initial state S, for events E.
// E is the event type
// S is the state type.
func New[E constraints.Ordered, S constraints.Ordered](initState S, transitions Transforms[E, S]) *FSM[E, S] {
	return NewFromTransitions[E, S](initState, NewTranslator[E, S](transitions))
}

// NewFromTransitions constructs a generic FSM with an initial state S, for events E.
// E is the event type
// S is the state type.
func NewFromTransitions[E constraints.Ordered, S constraints.Ordered](initState S, ts Translator[E, S]) *FSM[E, S] {
	return &FSM[E, S]{
		current:    initState,
		Translator: ts,
	}
}

// Clone the FSM without Mutex.
func (f *FSM[E, S]) Clone() *FSM[E, S] {
	return &FSM[E, S]{
		current:    f.current,
		Translator: f.Translator,
	}
}

// CloneWithState clone the FSM with new state without Mutex.
func (f *FSM[E, S]) CloneWithState(newState S) *FSM[E, S] {
	return &FSM[E, S]{
		current:    newState,
		Translator: f.Translator,
	}
}

// Current returns the current state.
func (f *FSM[E, S]) Current() S {
	f.stateMu.RLock()
	defer f.stateMu.RUnlock()
	return f.current
}

// Is returns true if state is the current state.
func (f *FSM[E, S]) Is(state S) bool {
	f.stateMu.RLock()
	defer f.stateMu.RUnlock()
	return state == f.current
}

// SetState allows the user to move to the given state from current state.
func (f *FSM[E, S]) SetState(state S) {
	f.stateMu.Lock()
	defer f.stateMu.Unlock()
	f.current = state
}

// Trigger call a state transition with the named event.
// It will return nil if src state change to dst state success or one of these errors:
//
// - ErrInappropriateEvent: event X inappropriate in the state Y
// - ErrNonExistEvent: event X does not exist
func (f *FSM[E, S]) Trigger(event E) error {
	f.stateMu.Lock()
	defer f.stateMu.Unlock()
	dst, err := f.Translator.Trigger(f.current, event)
	if err != nil {
		return err
	}
	f.current = dst
	return nil
}

// ShouldTrigger return dst state transition with the named event and src state.
// It will return if src state change to dst state success or holds the same as the current state:
func (f *FSM[E, S]) ShouldTrigger(event E) {
	f.stateMu.Lock()
	defer f.stateMu.Unlock()
	f.current = f.Translator.ShouldTrigger(f.current, event)
}

// IsCan returns true if event can occur in the current state.
func (f *FSM[E, S]) IsCan(event E) bool {
	return f.Translator.IsCan(f.current, event)
}

// IsAllCan returns true if all the events can occur in src state.
func (f *FSM[E, S]) IsAllCan(event ...E) bool {
	return f.Translator.IsAllCan(f.current, event...)
}

// AvailTransitionEvent returns a list of available transition event in the
// current state.
func (f *FSM[E, S]) AvailTransitionEvent() []E {
	return f.Translator.AvailTransitionEvent(f.current)
}
