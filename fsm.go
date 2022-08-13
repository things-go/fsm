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
	// stateMu guards access to the current state.
	stateMu sync.RWMutex

	// translator maps events and source states to destination states.
	// This is immutable
	translator Translator[E, S]
}

// New constructs a generic FSM with an initial state S, for events E.
// E is the event type
// S is the state type.
func New[E constraints.Ordered, S constraints.Ordered](initState S, transitions Transforms[E, S]) *FSM[E, S] {
	return NewFromTransitions(initState, NewTranslator[E, S](transitions))
}

// NewFromTransitions constructs a generic FSM with an initial state S, for events E.
// E is the event type
// S is the state type.
func NewFromTransitions[E constraints.Ordered, S constraints.Ordered](initState S, ts Translator[E, S]) *FSM[E, S] {
	return &FSM[E, S]{
		current:    initState,
		translator: ts,
	}
}

// Clone the FSM.
func (f *FSM[E, S]) Clone() *FSM[E, S] {
	return &FSM[E, S]{
		current:    f.current,
		translator: f.translator,
	}
}

// CloneWithNewState clone the FSM with new state.
func (f *FSM[E, S]) CloneWithState(newState S) *FSM[E, S] {
	return &FSM[E, S]{
		current:    newState,
		translator: f.translator,
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

// IsCan returns true if event can occur in the current state.
func (f *FSM[E, S]) IsCan(event E) bool {
	return f.translator.IsCan(event, f.current)
}

// AvailTransitionEvent returns a list of available transition event in the
// current state.
func (f *FSM[E, S]) AvailTransitionEvent() []E {
	return f.translator.AvailTransitionEvent(f.current)
}

// HasEvent returns true if event has supported.
func (f *FSM[E, S]) HasEvent(event E) bool {
	return f.translator.HasEvent(event)
}

// Trigger call a state transition with the named event.
//
// It will return nil if src state change to dst state success or one of these errors:
//
// - ErrInappropriateEvent: event X inappropriate in the state Y
//
// - ErrNonExistEvent: event X does not exist
func (f *FSM[E, S]) Trigger(event E) error {
	dst, err := f.translator.Trigger(event, f.current)
	if err != nil {
		return err
	}
	f.stateMu.Lock()
	defer f.stateMu.Unlock()
	f.current = dst
	return nil
}
