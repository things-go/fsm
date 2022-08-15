package fsm

import (
	"testing"

	"golang.org/x/exp/slices"
)

func TestUnsafeClone(t *testing.T) {
	ts := NewTranslator(Transforms[string, string]{
		{Event: "open", Src: []string{"closed"}, Dst: "open"},
		{Event: "close", Src: []string{"open"}, Dst: "closed"},
	})

	fsm := NewUnsafeFSMFromTransitions("close", ts)
	fsm1 := fsm.Clone()
	if fsm1.Current() != fsm.Current() {
		t.Errorf("expected same current state")
	}

	fsm2 := fsm.CloneWithState("open")
	if fsm2.Current() != "open" {
		t.Error("expected state to be 'open'")
	}
}

func TestUnsafeState(t *testing.T) {
	fsm := NewUnsafeFSM(
		"walking",
		Transforms[string, string]{
			{Event: "walk", Src: []string{"start"}, Dst: "walking"},
			{Event: "look", Src: []string{"walking"}, Dst: "walking"},
		},
	)

	if !fsm.ContainEvent("walk") {
		t.Error("expected support event 'walk'")
	}
	if fsm.ContainEvent("nosupport") {
		t.Error("expected not support event 'nosupport'")
	}
	if !fsm.ContainAllEvent("walk", "look") {
		t.Error("expected support all event 'walk' and 'look'")
	}
	if fsm.ContainAllEvent("walk", "nosupport") {
		t.Error("expected support all event 'walk' and 'nosupport'")
	}

	fsm.SetState("start")
	if fsm.Current() != "start" {
		t.Error("expected state to be 'start'")
	}
	if !fsm.IsCan("walk") {
		t.Error("expected event can occur in the current state.")
	}
	err := fsm.Trigger("walk")
	if err != nil {
		t.Error("trigger is expected no error")
	}
	if !fsm.Is("walking") {
		t.Error("expected state to be 'walking'")
	}
}

func TestUnsafeAvailTransitionEvent(t *testing.T) {
	fsm := NewUnsafeFSM(
		"closed",
		Transforms[string, string]{
			{Event: "open", Src: []string{"closed"}, Dst: "open"},
			{Event: "close", Src: []string{"open"}, Dst: "closed"},
			{Event: "middle", Src: []string{"closed"}, Dst: "middle"},
		},
	)
	events := fsm.AvailTransitionEvent()
	if !(slices.Contains(events, "middle") && slices.Contains(events, "open")) {
		t.Error("expected contain [middle, open] event with current state")
	}
	if !fsm.IsAllCan("middle", "open") {
		t.Error("expected contain all [middle, open] event with current state")
	}
	if fsm.IsAllCan("open", "close") {
		t.Error("expected not contain all [middle, open] event with current state")
	}
}
func TestUnsafeInappropriateEvent(t *testing.T) {
	fsm := NewUnsafeFSM(
		"closed",
		Transforms[string, string]{
			{Event: "open", Src: []string{"closed"}, Dst: "open"},
			{Event: "close", Src: []string{"open"}, Dst: "closed"},
		},
	)

	err := fsm.Trigger("close")
	if err != ErrInappropriateEvent {
		t.Error("expected 'ErrInappropriateEvent' with correct state and event")
	}
	historyState := fsm.Current()
	fsm.ShouldTrigger("close")
	if historyState != fsm.Current() {
		t.Error("ShouldTrigger expected hold original state with correct event")
	}
}

func TestUnsafeNonExistEvent(t *testing.T) {
	fsm := NewUnsafeFSM(
		"closed",
		Transforms[string, string]{
			{Event: "open", Src: []string{"closed"}, Dst: "open"},
			{Event: "close", Src: []string{"open"}, Dst: "closed"},
		},
	)

	err := fsm.Trigger("lock")
	if err != ErrNonExistEvent {
		t.Error("expected 'ErrNonExistEvent' with incorrect event")
	}
	historyState := fsm.Current()
	fsm.ShouldTrigger("lock")
	if historyState != fsm.Current() {
		t.Error("ShouldTrigger expected hold original state with incorrect event")
	}
}
