package fsm

import (
	"golang.org/x/exp/constraints"
)

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

	// transition support.

	// Name return the name of the transition.
	Name() string
	// Transform return the dst state transition with the named event and src state.
	// It will return nil if src state change to dst state success or one of these errors:
	//
	// - ErrInappropriateEvent: event inappropriate in the src state.
	// - ErrNonExistEvent: event does not exist
	Transform(srcState S, event E) (dstState S, err error)
	// Match reports whether it can be transform to dst state with the named event and src state.
	Match(srcState, dstState S, event E) (bool, error)
	// MatchOccur returns true if event can occur in src state.
	MatchOccur(srcState S, event E) bool
	// MatchAllOccur returns true if all the events can occur in src state.
	MatchAllOccur(srcState S, events ...E) bool
	// ContainsEvent returns true if support the event.
	ContainsEvent(event E) bool
	// ContainsAllEvent returns true if support all event.
	ContainsAllEvent(events ...E) bool
	// AvailEvents returns a list of available transform event in src state.
	AvailEvents(srcState S) []E
	// SortedTriggerSource return a list of sorted trigger source
	SortedTriggerSource() []TriggerSource[E, S]
	// SortedStates return a list of sorted states.
	SortedStates() []S
	// SortedEvents return a list of sorted events.
	SortedEvents() []E
	// StateName returns a event name.
	EventName(event E) string
	// StateName returns a state name.
	StateName(state S) string
}

type ErrorTranslator interface {
	Translate(err error) error
}
