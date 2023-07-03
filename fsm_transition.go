package fsm

import (
	"errors"
	"fmt"
	"sort"

	"golang.org/x/exp/constraints"
	"golang.org/x/exp/maps"
	"golang.org/x/exp/slices"
)

var (
	ErrInappropriateEvent = errors.New("fsm: event inappropriate in the state")
	ErrNonExistEvent      = errors.New("fsm: event does not exist")
)

type ITransition[E constraints.Ordered, S constraints.Ordered] interface {
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
	// AvailSourceStates returns a list of available source state in the event.
	AvailSourceStates(event E) []S
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

// Transform represents an event when initializing the Fsm.
//
// The event can have one or more source states that is valid for performing
// the transition. If the Fsm is in one of the source states it will end up in
// the specified destination state.
type Transform[E constraints.Ordered, S constraints.Ordered] struct {
	// Name the event.
	Name string
	// Event is the event used when calling for the transform.
	Event E
	// Src is a slice of source states that the Fsm must be in to perform a
	// state transform.
	Src []S
	// Dst is the destination state that the Fsm will be in if the transform
	// succeeds.
	Dst S
}

// TriggerSource is storing the trigger source.
type TriggerSource[E constraints.Ordered, S constraints.Ordered] struct {
	// event is the name of the event that the keys refers to.
	event E
	// src is the source from where the event can transition.
	src S
}

// Event the trigger source event.
func (e *TriggerSource[E, S]) Event() E { return e.event }

// State the trigger source state.
func (e *TriggerSource[E, S]) State() S { return e.src }

// Transition contain events and source states to destination states.
// NOTE: This is immutable
type Transition[E constraints.Ordered, S constraints.Ordered] struct {
	// name is the name of the transition.
	name string
	// contain all support event and name.
	events map[E]string
	// contain all support state and name.
	states map[S]string
	// mapping map the trigger source to destination states.
	mapping map[TriggerSource[E, S]]S
	// translate error
	translate ErrorTranslator
}

type TransitionBuilder[E constraints.Ordered, S constraints.Ordered] struct {
	// name is the name of the transition.
	name string
	// transforms
	transforms []Transform[E, S]
	// contain all support state and name.
	states map[S]string
	// translate error
	translate ErrorTranslator
}

func NewTransitionBuilder[E constraints.Ordered, S constraints.Ordered](transforms []Transform[E, S]) *TransitionBuilder[E, S] {
	return &TransitionBuilder[E, S]{
		transforms: transforms,
	}
}

func (b *TransitionBuilder[E, S]) Name(name string) *TransitionBuilder[E, S] {
	b.name = name
	return b
}

func (b *TransitionBuilder[E, S]) StateNames(states map[S]string) *TransitionBuilder[E, S] {
	b.states = states
	return b
}

func (b *TransitionBuilder[E, S]) TranslatorError(translate ErrorTranslator) *TransitionBuilder[E, S] {
	b.translate = translate
	return b
}

func (b *TransitionBuilder[E, S]) Build() *Transition[E, S] {
	t := &Transition[E, S]{
		name:      b.name,
		events:    make(map[E]string),
		states:    make(map[S]string),
		mapping:   make(map[TriggerSource[E, S]]S),
		translate: b.translate,
	}
	for _, ts := range b.transforms {
		t.events[ts.Event] = ts.Name
		for _, src := range ts.Src {
			t.mapping[TriggerSource[E, S]{ts.Event, src}] = ts.Dst
			t.states[src] = ""
			t.states[ts.Dst] = ""
		}
	}
	for k, v := range b.states {
		t.states[k] = v
	}
	return t
}

// NewTransition new a transition instance.
func NewTransition[E constraints.Ordered, S constraints.Ordered](transforms []Transform[E, S]) *Transition[E, S] {
	return NewTransitionBuilder[E, S](transforms).
		Build()
}

// Name return the name of the transition.
func (t *Transition[E, S]) Name() string { return t.name }

// Transform return the dst state transition with the named event and src state.
// It will return nil if src state change to dst state success or one of these errors:
//
// - ErrInappropriateEvent: event inappropriate in the src state.
// - ErrNonExistEvent: event does not exist
func (t *Transition[E, S]) Transform(srcState S, event E) (dstState S, err error) {
	dstState, ok := t.mapping[TriggerSource[E, S]{event, srcState}]
	if !ok {
		for ts := range t.mapping {
			if ts.event == event {
				return dstState, t.translateError(ErrInappropriateEvent)
			}
		}
		return dstState, t.translateError(ErrNonExistEvent)
	}
	return dstState, nil
}

// Match reports whether it can be transform to dst state with the named event and src state.
func (t *Transition[E, S]) Match(srcState, dstState S, event E) (bool, error) {
	targetDstState, err := t.Transform(srcState, event)
	if err != nil {
		return false, err
	}
	return targetDstState == dstState, nil
}

// MatchOccur returns true if event can occur in src state.
func (t *Transition[E, S]) MatchOccur(srcState S, event E) bool {
	_, ok := t.mapping[TriggerSource[E, S]{event, srcState}]
	return ok
}

// MatchAllOccur returns true if all the events can occur in src state.
func (t *Transition[E, S]) MatchAllOccur(srcState S, events ...E) bool {
	occurEvents := t.availEvents(srcState)
	for _, e := range events {
		if _, ok := occurEvents[e]; !ok {
			return false
		}
	}
	return true
}

// ContainsEvent returns true if support the event.
func (t *Transition[E, S]) ContainsEvent(event E) bool {
	_, ok := t.events[event]
	return ok
}

// ContainsAllEvent returns true if support all event.
func (t *Transition[E, S]) ContainsAllEvent(events ...E) bool {
	for _, event := range events {
		_, ok := t.events[event]
		if !ok {
			return false
		}
	}
	return true
}

// AvailEvents returns a list of available transform event in src state.
func (t *Transition[E, S]) AvailEvents(srcState S) []E {
	events := t.availEvents(srcState)
	return maps.Keys(events)
}

// AvailSourceStates returns a list of available source state in the event.
func (t *Transition[E, S]) AvailSourceStates(event E) []S {
	srcs := make([]S, 0, 8)
	for ts := range t.mapping {
		if ts.event == event {
			srcs = append(srcs, ts.src)
		}
	}
	return srcs
}

// SortedTriggerSource return a list of sorted trigger source
func (t *Transition[E, S]) SortedTriggerSource() []TriggerSource[E, S] {
	triggerSources := maps.Keys(t.mapping)
	sort.Slice(triggerSources, func(i, j int) bool {
		if triggerSources[i].src == triggerSources[j].src {
			return triggerSources[i].event < triggerSources[j].event
		}
		return triggerSources[i].src < triggerSources[j].src
	})
	return triggerSources
}

// SortedStates return a list of sorted states.
func (t *Transition[E, S]) SortedStates() []S {
	states := maps.Keys(t.states)
	slices.Sort(states)
	return states
}

// SortedEvents return a list of sorted events.
func (t *Transition[E, S]) SortedEvents() []E {
	events := maps.Keys(t.events)
	slices.Sort(events)
	return events
}

// StateName returns a event name.
func (t *Transition[E, S]) EventName(event E) string {
	v, ok := t.events[event]
	if ok && v != "" {
		return v
	}
	return fmt.Sprintf("%v", event)
}

// StateName returns a state name.
func (t *Transition[E, S]) StateName(state S) string {
	v, ok := t.states[state]
	if ok && v != "" {
		return v
	}
	return fmt.Sprintf("%v", state)
}

// availEvents returns an available transform event in src state.
func (t *Transition[E, S]) availEvents(srcState S) map[E]struct{} {
	occurEvents := make(map[E]struct{})
	for ts := range t.mapping {
		if ts.src == srcState {
			occurEvents[ts.event] = struct{}{}
		}
	}
	return occurEvents
}

// availEvents returns an available transform event in src state.
func (t *Transition[E, S]) translateError(err error) error {
	if err == nil || t.translate == nil {
		return err
	}
	return t.translate.Translate(err)
}
