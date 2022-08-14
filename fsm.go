package fsm

import (
	"sync"

	"golang.org/x/exp/constraints"
)

// FSM is the state machine that holds the current state.
// E ist the event
// S is the state
type FSM[E constraints.Ordered, S constraints.Ordered] struct {
	*FSMUnsafe[E, S]

	// stateMu guards access to the current state.
	stateMu sync.RWMutex
}

// New constructs a generic FSM with an initial state S, for events E.
// E is the event type
// S is the state type.
func New[E constraints.Ordered, S constraints.Ordered](initState S, transitions Transforms[E, S]) *FSM[E, S] {
	return &FSM[E, S]{
		FSMUnsafe: NewUnsafeFSM[E, S](initState, transitions),
	}
}

// NewFromTransitions constructs a generic FSM with an initial state S, for events E.
// E is the event type
// S is the state type.
func NewFromTransitions[E constraints.Ordered, S constraints.Ordered](initState S, ts Translator[E, S]) *FSM[E, S] {
	return &FSM[E, S]{
		FSMUnsafe: NewUnsafeFSMFromTransitions(initState, ts),
	}
}

// Clone the FSM without Mutex.
func (f *FSM[E, S]) Clone() *FSM[E, S] {
	return &FSM[E, S]{
		FSMUnsafe: f.FSMUnsafe.Clone(),
	}
}

// CloneWithState clone the FSM with new state without Mutex.
func (f *FSM[E, S]) CloneWithState(newState S) *FSM[E, S] {
	return &FSM[E, S]{
		FSMUnsafe: f.FSMUnsafe.CloneWithState(newState),
	}
}

// Current returns the current state.
func (f *FSM[E, S]) Current() S {
	f.stateMu.RLock()
	defer f.stateMu.RUnlock()
	return f.FSMUnsafe.Current()
}

// Is returns true if state is the current state.
func (f *FSM[E, S]) Is(state S) bool {
	f.stateMu.RLock()
	defer f.stateMu.RUnlock()
	return f.FSMUnsafe.Is(state)
}

// SetState allows the user to move to the given state from current state.
func (f *FSM[E, S]) SetState(state S) {
	f.stateMu.Lock()
	defer f.stateMu.Unlock()
	f.FSMUnsafe.SetState(state)
}

// Trigger call a state transition with the named event.
// It will return nil if src state change to dst state success or one of these errors:
//
// - ErrInappropriateEvent: event X inappropriate in the state Y
// - ErrNonExistEvent: event X does not exist
func (f *FSM[E, S]) Trigger(event E) error {
	f.stateMu.Lock()
	defer f.stateMu.Unlock()
	return f.FSMUnsafe.Trigger(event)
}

// ShouldTrigger return dst state transition with the named event and src state.
// It will return if src state change to dst state success or holds the same as the current state:
func (f *FSM[E, S]) ShouldTrigger(event E) {
	f.stateMu.Lock()
	defer f.stateMu.Unlock()
	f.FSMUnsafe.ShouldTrigger(event)
}
