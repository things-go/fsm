package fsm

import "golang.org/x/exp/constraints"

var _ IFsm[string, string] = (*Fsm[string, string])(nil)
var _ IFsm[int, string] = (*Fsm[int, string])(nil)

// Fsm is the state machine that holds the current state.
// E is the event
// S is the state
type Fsm[E constraints.Ordered, S constraints.Ordered] struct {
	// Transition contain events and source states to destination states.
	// This is immutable
	*Transition[E, S]
	// current is the state that the Fsm is currently in.
	current S
}

// NewFsm constructs a generic Fsm with an initial state S and a transition.
// E is the event type
// S is the state type.
func NewFsm[E constraints.Ordered, S constraints.Ordered](initState S, ts *Transition[E, S]) IFsm[E, S] {
	return &Fsm[E, S]{
		current:    initState,
		Transition: ts,
	}
}
func (f *Fsm[E, S]) Clone() IFsm[E, S] {
	return f.Transition.CloneFsm(f.current)
}
func (f *Fsm[E, S]) CloneNewState(newState S) IFsm[E, S] {
	return f.Transition.CloneFsm(newState)
}
func (f *Fsm[E, S]) Current() S         { return f.current }
func (f *Fsm[E, S]) Is(state S) bool    { return state == f.current }
func (f *Fsm[E, S]) SetCurrent(state S) { f.current = state }
func (f *Fsm[E, S]) Trigger(event E) error {
	dst, err := f.Transition.Transform(f.current, event)
	if err != nil {
		return err
	}
	f.current = dst
	return nil
}
func (f *Fsm[E, S]) MatchCurrentOccur(event E) bool {
	return f.Transition.MatchOccur(f.current, event)
}
func (f *Fsm[E, S]) MatchCurrentAllOccur(event ...E) bool {
	return f.Transition.MatchAllOccur(f.current, event...)
}
func (f *Fsm[E, S]) CurrentAvailEvents() []E {
	return f.Transition.AvailEvents(f.current)
}
func (f *Fsm[E, S]) Visualize(t VisualizeType) (string, error) {
	return Visualize[E, S](t, f)
}
