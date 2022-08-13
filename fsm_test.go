package fsm

import (
	"testing"

	"golang.org/x/exp/slices"
)

func TestClone(t *testing.T) {
	fsm := New(
		"close",
		Transforms[string, string]{
			{Event: "open", Src: []string{"closed"}, Dst: "open"},
			{Event: "close", Src: []string{"open"}, Dst: "closed"},
		},
	)
	fsm1 := fsm.Clone()
	if fsm1.Current() != fsm.Current() {
		t.Errorf("expected same current state")
	}

	fsm2 := fsm.CloneWithState("open")
	if fsm2.Current() != "open" {
		t.Error("expected state to be 'open'")
	}
}

func TestSameState(t *testing.T) {
	fsm := New(
		"start",
		Transforms[string, string]{
			{Event: "run", Src: []string{"start"}, Dst: "start"},
		},
	)
	err := fsm.Trigger("run")
	if err != nil {
		t.Errorf("expected trigger no error")
	}
	if fsm.Current() != "start" {
		t.Error("expected state to be 'start'")
	}
}

func TestState(t *testing.T) {
	fsm := New(
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

func TestAvailTransitionEvent(t *testing.T) {
	fsm := New(
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
}

func TestInappropriateEvent(t *testing.T) {
	fsm := New(
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

func TestNonExistEvent(t *testing.T) {
	fsm := New(
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

func TestMultipleSources(t *testing.T) {
	fsm := New(
		"one",
		Transforms[string, string]{
			{Event: "first", Src: []string{"one"}, Dst: "two"},
			{Event: "second", Src: []string{"two"}, Dst: "three"},
			{Event: "reset", Src: []string{"one", "two", "three"}, Dst: "one"},
		},
	)

	err := fsm.Trigger("first")
	if err != nil {
		t.Errorf("trigger failed %v", err)
	}
	if fsm.Current() != "two" {
		t.Error("expected state to be 'two'")
	}
	err = fsm.Trigger("reset")
	if err != nil {
		t.Errorf("transition failed %v", err)
	}
	if fsm.Current() != "one" {
		t.Error("expected state to be 'one'")
	}
	err = fsm.Trigger("first")
	if err != nil {
		t.Errorf("trigger failed %v", err)
	}
	err = fsm.Trigger("second")
	if err != nil {
		t.Errorf("trigger failed %v", err)
	}
	if fsm.Current() != "three" {
		t.Error("expected state to be 'three'")
	}
	err = fsm.Trigger("reset")
	if err != nil {
		t.Errorf("trigger failed %v", err)
	}
	if fsm.Current() != "one" {
		t.Error("expected state to be 'one'")
	}
}

func TestMultipleEvents(t *testing.T) {
	fsm := New(
		"start",
		Transforms[string, string]{
			{Event: "first", Src: []string{"start"}, Dst: "one"},
			{Event: "second", Src: []string{"start"}, Dst: "two"},
			{Event: "reset", Src: []string{"one"}, Dst: "reset_one"},
			{Event: "reset", Src: []string{"two"}, Dst: "reset_two"},
			{Event: "reset", Src: []string{"reset_one", "reset_two"}, Dst: "start"},
		},
	)

	err := fsm.Trigger("first")
	if err != nil {
		t.Errorf("trigger failed %v", err)
	}
	err = fsm.Trigger("reset")
	if err != nil {
		t.Errorf("trigger failed %v", err)
	}
	if fsm.Current() != "reset_one" {
		t.Error("expected state to be 'reset_one'")
	}
	err = fsm.Trigger("reset")
	if err != nil {
		t.Errorf("trigger failed %v", err)
	}
	if fsm.Current() != "start" {
		t.Error("expected state to be 'start'")
	}

	err = fsm.Trigger("second")
	if err != nil {
		t.Errorf("trigger failed %v", err)
	}
	err = fsm.Trigger("reset")
	if err != nil {
		t.Errorf("trigger failed %v", err)
	}
	if fsm.Current() != "reset_two" {
		t.Error("expected state to be 'reset_two'")
	}
	err = fsm.Trigger("reset")
	if err != nil {
		t.Errorf("trigger failed %v", err)
	}
	if fsm.Current() != "start" {
		t.Error("expected state to be 'start'")
	}
}
