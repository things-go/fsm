package fsm

import (
	"errors"

	"golang.org/x/exp/constraints"
)

var ErrInappropriateEvent = errors.New("fsm: event inappropriate in the state")
var ErrNonExistEvent = errors.New("fsm: event does not exist")

// Transition represents an event when initializing the FSM.
//
// The event can have one or more source states that is valid for performing
// the transition. If the FSM is in one of the source states it will end up in
// the specified destination state, calling all defined callbacks as it goes.
type Transition[E constraints.Ordered, S constraints.Ordered] struct {
	// Event is the event used when calling for a transition.
	Event E

	// Src is a slice of source states that the FSM must be in to perform a
	// state transition.
	Src []S

	// Dst is the destination state that the FSM will be in if the transition
	// succeeds.
	Dst S
}

// Transitions is a shorthand for defining the transition map in NewFSM.
type Transitions[E constraints.Ordered, S constraints.Ordered] []Transition[E, S]

// eKey is a struct key used for storing the transition map.
type eKey[E constraints.Ordered, S constraints.Ordered] struct {
	// event is the name of the event that the keys refers to.
	event E

	// src is the source from where the event can transition.
	src S
}

type TransitionsMap[E constraints.Ordered, S constraints.Ordered] map[eKey[E, S]]S

func NewTransitions[E constraints.Ordered, S constraints.Ordered](transitions Transitions[E, S]) TransitionsMap[E, S] {
	ts := TransitionsMap[E, S]{}
	for _, e := range transitions {
		for _, src := range e.Src {
			ts[eKey[E, S]{e.Event, src}] = e.Dst
		}
	}
	return ts
}

// IsCan returns true if event can occur in src state.
func (ts TransitionsMap[E, S]) IsCan(event E, srcState S) bool {
	_, ok := ts[eKey[E, S]{event, srcState}]
	return ok
}

// AvailTransitionEvent returns a list of available transition event in src state.
func (ts TransitionsMap[E, S]) AvailTransitionEvent(srcState S) []E {
	var events []E
	for key := range ts {
		if key.src == srcState {
			events = append(events, key.event)
		}
	}
	return events
}

// Trigger return dst state transition with the named event and src state.
//
// It will return nil if src state change to dst state success or one of these errors:
//
// - ErrInappropriateEvent: event X inappropriate in the state Y
//
// - ErrNonExistEvent: event X does not exist
func (ts TransitionsMap[E, S]) Trigger(event E, srcState S) (dstState S, err error) {
	var ok bool

	dstState, ok = ts[eKey[E, S]{event, srcState}]
	if !ok {
		for ek := range ts {
			if ek.event == event {
				return dstState, ErrInappropriateEvent
			}
		}
		return dstState, ErrNonExistEvent
	}
	return dstState, nil
}
