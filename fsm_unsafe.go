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

	// translator contain events and source states to destination states.
	// This is immutable
	Translator[E, S]
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
		Translator: ts,
	}
}

// Clone the FSMUnsafe.
func (f *FSMUnsafe[E, S]) Clone() *FSMUnsafe[E, S] {
	return &FSMUnsafe[E, S]{
		current:    f.current,
		Translator: f.Translator,
	}
}

// CloneWithState clone the FSMUnsafe with new state.
func (f *FSMUnsafe[E, S]) CloneWithState(newState S) *FSMUnsafe[E, S] {
	return &FSMUnsafe[E, S]{
		current:    newState,
		Translator: f.Translator,
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
	dst, err := f.Translator.Trigger(f.current, event)
	if err != nil {
		return err
	}
	f.current = dst
	return nil
}

// ShouldTrigger return dst state transition with the named event and src state.
// It will return if src state change to dst state success or holds the same as the current state:
func (f *FSMUnsafe[E, S]) ShouldTrigger(event E) {
	f.current = f.Translator.ShouldTrigger(f.current, event)
}

// IsCan returns true if event can occur in the current state.
func (f *FSMUnsafe[E, S]) IsCan(event E) bool {
	return f.Translator.IsCan(f.current, event)
}

// IsAllCan returns true if all the events can occur in src state.
func (f *FSMUnsafe[E, S]) IsAllCan(event ...E) bool {
	return f.Translator.IsAllCan(f.current, event...)
}

// AvailTransitionEvent returns a list of available transition event in the
// current state.
func (f *FSMUnsafe[E, S]) AvailTransitionEvent() []E {
	return f.Translator.AvailTransitionEvent(f.current)
}
