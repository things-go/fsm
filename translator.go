package fsm

import (
	"errors"

	"golang.org/x/exp/constraints"
	"golang.org/x/exp/maps"
	"golang.org/x/exp/slices"
)

var ErrInappropriateEvent = errors.New("fsm: event inappropriate in the state")
var ErrNonExistEvent = errors.New("fsm: event does not exist")

// Transform represents an event when initializing the FSM.
//
// The event can have one or more source states that is valid for performing
// the transition. If the FSM is in one of the source states it will end up in
// the specified destination state, calling all defined callbacks as it goes.
type Transform[E constraints.Ordered, S constraints.Ordered] struct {
	// Event is the event used when calling for a transition.
	Event E

	// Src is a slice of source states that the FSM must be in to perform a
	// state transition.
	Src []S

	// Dst is the destination state that the FSM will be in if the transition
	// succeeds.
	Dst S
}

// Transforms is a shorthand for defining the transition map in NewFSM.
type Transforms[E constraints.Ordered, S constraints.Ordered] []Transform[E, S]

// eKey is a struct key used for storing the transition map.
type eKey[E constraints.Ordered, S constraints.Ordered] struct {
	// event is the name of the event that the keys refers to.
	event E

	// src is the source from where the event can transition.
	src S
}

type Translator[E constraints.Ordered, S constraints.Ordered] map[eKey[E, S]]S

func NewTranslator[E constraints.Ordered, S constraints.Ordered](transitions Transforms[E, S]) Translator[E, S] {
	ts := Translator[E, S]{}
	for _, e := range transitions {
		for _, src := range e.Src {
			ts[eKey[E, S]{e.Event, src}] = e.Dst
		}
	}
	return ts
}

// Trigger return dst state transition with the named event and src state.
// It will return nil if src state change to dst state success or one of these errors:
//
// - ErrInappropriateEvent: event X inappropriate in the state Y
// - ErrNonExistEvent: event X does not exist
func (ts Translator[E, S]) Trigger(srcState S, event E) (dstState S, err error) {
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

// ShouldTrigger return dst state transition with the named event and src state.
// It will return if src state change to dst state success or holds the same as the src state:
func (ts Translator[E, S]) ShouldTrigger(srcState S, event E) (dstState S) {
	var err error

	dstState, err = ts.Trigger(srcState, event)
	if err != nil {
		dstState = srcState
	}
	return dstState
}

// IsCan returns true if event can occur in src state.
func (ts Translator[E, S]) IsCan(srcState S, event E) bool {
	_, ok := ts[eKey[E, S]{event, srcState}]
	return ok
}

// IsAllCan returns true if all the events can occur in src state.
func (ts Translator[E, S]) IsAllCan(srcState S, events ...E) bool {
	es := ts.AvailTransitionEvent(srcState)

	for _, event := range events {
		if !slices.Contains(es, event) {
			return false
		}
	}
	return true
}

// AvailTransitionEvent returns a list of available transition event in src state.
func (ts Translator[E, S]) AvailTransitionEvent(srcState S) []E {
	es := make(map[E]struct{})
	for key := range ts {
		if key.src == srcState {
			es[key.event] = struct{}{}
		}
	}
	return maps.Keys(es)
}

// ContainEvent returns true if support the event.
func (ts Translator[E, S]) ContainEvent(event E) bool {
	for key := range ts {
		if key.event == event {
			return true
		}
	}
	return false
}

// ContainAllEvent returns true if support all event.
func (ts Translator[E, S]) ContainAllEvent(events ...E) bool {
	es := make(map[E]struct{})
	for key := range ts {
		es[key.event] = struct{}{}
	}
	for _, event := range events {
		_, ok := es[event]
		if !ok {
			return false
		}
	}
	return true
}
