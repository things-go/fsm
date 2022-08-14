package fsm

import (
	"golang.org/x/exp/constraints"
)

// FSMUnsafe is the state machine that holds the current state.
// E ist the event
// S is the state
type FSMUnsafe[E constraints.Ordered, S constraints.Ordered] struct {
	// current is the state that the FSMUnsafe is currently in.
	current S

	// translator maps events and source states to destination states.
	// This is immutable
	translator Translator[E, S]
}

// NewUnsafeFSM constructs a generic FSMUnsafe with an initial state S, for events E.
// E is the event type
// S is the state type.
func NewUnsafeFSM[E constraints.Ordered, S constraints.Ordered](initState S, transitions Transforms[E, S]) *FSMUnsafe[E, S] {
	return NewUnsafeFSMFromTransitions(initState, NewTranslator[E, S](transitions))
}

// NewUnsafeFSMFromTransitions constructs a generic FSMUnsafe with an initial state S, for events E.
// E is the event type
// S is the state type.
func NewUnsafeFSMFromTransitions[E constraints.Ordered, S constraints.Ordered](initState S, ts Translator[E, S]) *FSMUnsafe[E, S] {
	return &FSMUnsafe[E, S]{
		current:    initState,
		translator: ts,
	}
}

// Clone the FSMUnsafe.
func (f *FSMUnsafe[E, S]) Clone() *FSMUnsafe[E, S] {
	return &FSMUnsafe[E, S]{
		current:    f.current,
		translator: f.translator,
	}
}

// CloneWithState clone the FSMUnsafe with new state.
func (f *FSMUnsafe[E, S]) CloneWithState(newState S) *FSMUnsafe[E, S] {
	return &FSMUnsafe[E, S]{
		current:    newState,
		translator: f.translator,
	}
}

// Current returns the current state.
func (f *FSMUnsafe[E, S]) Current() S {
	return f.current
}

// Is returns true if state is the current state.
func (f *FSMUnsafe[E, S]) Is(state S) bool {
	return state == f.current
}

// SetState allows the user to move to the given state from current state.
func (f *FSMUnsafe[E, S]) SetState(state S) {
	f.current = state
}

// Trigger call a state transition with the named event.
//
// It will return nil if src state change to dst state success or one of these errors:
//
// - ErrInappropriateEvent: event X inappropriate in the state Y
//
// - ErrNonExistEvent: event X does not exist
func (f *FSMUnsafe[E, S]) Trigger(event E) error {
	dst, err := f.translator.Trigger(event, f.current)
	if err != nil {
		return err
	}
	f.current = dst
	return nil
}

// IsCan returns true if event can occur in the current state.
func (f *FSMUnsafe[E, S]) IsCan(event E) bool {
	return f.translator.IsCan(event, f.current)
}

// AvailTransitionEvent returns a list of available transition event in the
// current state.
func (f *FSMUnsafe[E, S]) AvailTransitionEvent() []E {
	return f.translator.AvailTransitionEvent(f.current)
}

// HasEvent returns true if event has supported.
func (f *FSMUnsafe[E, S]) HasEvent(event E) bool {
	return f.translator.HasEvent(event)
}
