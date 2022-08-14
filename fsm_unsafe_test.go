package fsm

import (
	"testing"
)

func TestUnsafeState(t *testing.T) {
	fsm := NewUnsafeFSM(
		"walking",
		Transforms[string, string]{
			{Event: "walk", Src: []string{"start"}, Dst: "walking"},
		},
	)

	if !fsm.HasEvent("walk") {
		t.Error("expected support event 'walk'")
	}
	if fsm.HasEvent("nosupport") {
		t.Error("expected not support event 'nosupport'")
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
		t.Error("expected 'ErrNonExistEvent' with correct event")
	}
}
