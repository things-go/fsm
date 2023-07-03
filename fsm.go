package fsm

import (
	"golang.org/x/exp/constraints"
)

type Diagram interface {
	// Visualize outputs a visualization of a Fsm in the desired format.
	// If the type is not given it defaults to Graphviz
	Visualize(t VisualizeType) (string, error)
}

type IFsm[E constraints.Ordered, S constraints.Ordered] interface {
	// Clone the Fsm.
	Clone() IFsm[E, S]
	// CloneNewState clone the Fsm with new state.
	CloneNewState(newState S) IFsm[E, S]
	// Current returns the current state.
	Current() S
	// Is returns true if state match the current state.
	Is(state S) bool
	// SetCurrent allows the user to move to the given state from current state.
	SetCurrent(state S)
	// Trigger call a state transition with the named event and src state if success will change the current state.
	// It will return nil if src state change to dst state success or one of these errors:
	//
	// - ErrInappropriateEvent: event inappropriate in the src state.
	// - ErrNonExistEvent: event does not exist
	Trigger(event E) error
	// MatchOccur returns true if event can occur in the current state.
	MatchCurrentOccur(event E) bool
	// MatchAllOccur returns true if all the events can occur in current state.
	MatchCurrentAllOccur(event ...E) bool
	// AvailEvents returns a list of available transform event in current state.
	CurrentAvailEvents() []E

	ITransition[E, S]
	Diagram
}

type ErrorTranslator interface {
	Translate(err error) error
}
