package fsm

import (
	"sync"

	"golang.org/x/exp/constraints"
)

var _ IFsm[string, string] = (*SafeFsm[string, string])(nil)
var _ IFsm[int, string] = (*SafeFsm[int, string])(nil)

// SafeFsm is the state machine that holds the current state and mutex.
// E is the event
// S is the state
type SafeFsm[E constraints.Ordered, S constraints.Ordered] struct {
	// Transition contain events and source states to destination states.
	// This is immutable
	*Transition[E, S]
	// mu guards access to the current state.
	mu sync.RWMutex
	// current is the state that the Fsm is currently in.
	current S
}

// NewSafeFsm constructs a generic Fsm with an initial state S and a transition.
// E is the event type
// S is the state type.
func NewSafeFsm[E constraints.Ordered, S constraints.Ordered](initState S, ts *Transition[E, S]) IFsm[E, S] {
	return &SafeFsm[E, S]{
		current:    initState,
		Transition: ts,
	}
}
func (f *SafeFsm[E, S]) Clone() IFsm[E, S] {
	return f.Transition.CloneSafeFsm(f.current)
}
func (f *SafeFsm[E, S]) CloneNewState(newState S) IFsm[E, S] {
	return f.Transition.CloneSafeFsm(newState)
}
func (f *SafeFsm[E, S]) Current() S {
	f.mu.RLock()
	defer f.mu.RUnlock()
	return f.current
}
func (f *SafeFsm[E, S]) SetCurrent(newState S) {
	f.mu.Lock()
	defer f.mu.Unlock()
	f.current = newState
}
func (f *SafeFsm[E, S]) Is(state S) bool {
	f.mu.RLock()
	defer f.mu.RUnlock()
	return state == f.current
}
func (f *SafeFsm[E, S]) Trigger(event E) error {
	f.mu.Lock()
	defer f.mu.Unlock()
	dst, err := f.Transition.Transform(f.current, event)
	if err != nil {
		return err
	}
	f.current = dst
	return nil
}
func (f *SafeFsm[E, S]) MatchCurrentOccur(event E) bool {
	return f.Transition.MatchOccur(f.current, event)
}
func (f *SafeFsm[E, S]) MatchCurrentAllOccur(event ...E) bool {
	return f.Transition.MatchAllOccur(f.current, event...)
}
func (f *SafeFsm[E, S]) CurrentAvailEvents() []E {
	return f.Transition.AvailEvents(f.current)
}
func (f *SafeFsm[E, S]) Visualize(t VisualizeType) (string, error) {
	return Visualize[E, S](t, f)
}
