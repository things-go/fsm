package fsm

import (
	"errors"
	"strings"
	"testing"

	"golang.org/x/exp/slices"
)

func Test_Fsm_Clone(t *testing.T) {
	test_Fsm_Clone(t, NewSafeFsm[string, string])
	test_Fsm_Clone(t, NewFsm[string, string])
}

func test_Fsm_Clone(t *testing.T, newFsm func(initState string, ts *Transition[string, string]) IFsm[string, string]) {
	ts := NewTransition([]Transform[string, string]{
		{Event: "open", Src: []string{"closed"}, Dst: "open"},
		{Event: "close", Src: []string{"open"}, Dst: "closed"},
	})

	fsm := newFsm("close", ts)
	fsm1 := fsm.Clone()
	if fsm1.Current() != fsm.Current() {
		t.Errorf("expected same current state")
	}
	fsm2 := fsm.CloneNewState("open")
	if fsm2.Current() != "open" {
		t.Error("expected state to be 'open'")
	}
}

func Test_Fsm_SameState(t *testing.T) {
	test_Fsm_SameState(t, NewSafeFsm[string, string])
	test_Fsm_SameState(t, NewFsm[string, string])
}

func test_Fsm_SameState(t *testing.T, newFsm func(initState string, ts *Transition[string, string]) IFsm[string, string]) {
	fsm := newFsm(
		"start",
		NewTransition([]Transform[string, string]{
			{Event: "run", Src: []string{"start"}, Dst: "start"},
		}),
	)
	err := fsm.Trigger("run")
	if err != nil {
		t.Errorf("expected trigger no error")
	}
	if fsm.Current() != "start" {
		t.Error("expected state to be 'start'")
	}
}

func Test_Fsm_State(t *testing.T) {
	test_Fsm_State(t, NewSafeFsm[string, string])
	test_Fsm_State(t, NewFsm[string, string])
}

func test_Fsm_State(t *testing.T, newFsm func(initState string, ts *Transition[string, string]) IFsm[string, string]) {
	fsm := newFsm(
		"walking",
		NewTransition([]Transform[string, string]{
			{Event: "walk", Src: []string{"start"}, Dst: "walking"},
			{Event: "look", Src: []string{"walking"}, Dst: "walking"},
		}),
	)

	if !fsm.ContainsEvent("walk") {
		t.Error("expected support event 'walk'")
	}
	if fsm.ContainsEvent("nosupport") {
		t.Error("expected not support event 'nosupport'")
	}
	if !fsm.ContainsAllEvent("walk", "look") {
		t.Error("expected support all event 'walk' and 'look'")
	}
	if fsm.ContainsAllEvent("walk", "nosupport") {
		t.Error("expected support all event 'walk' and 'nosupport'")
	}

	if b, err := fsm.Match("start", "walking", "walk"); err != nil {
		t.Error("expected event src dst state match no error")
	} else {
		if !b {
			t.Error("expected event src dst state match")
		}
	}

	if b, err := fsm.Match("start", "walking", "look"); err == nil {
		t.Error("expected event src dst state match has error")
	} else {
		if b {
			t.Error("expected event src dst state match false")
		}
	}

	fsm.SetCurrent("start")
	if fsm.Current() != "start" {
		t.Error("expected state to be 'start'")
	}
	if !fsm.MatchCurrentOccur("walk") {
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

func Test_Fsm_AvailEvents(t *testing.T) {
	test_Fsm_AvailEvents(t, NewSafeFsm[string, string])
	test_Fsm_AvailEvents(t, NewFsm[string, string])
}

func test_Fsm_AvailEvents(t *testing.T, newFsm func(initState string, ts *Transition[string, string]) IFsm[string, string]) {
	fsm := newFsm(
		"closed",
		NewTransition([]Transform[string, string]{
			{Event: "open", Src: []string{"closed"}, Dst: "open"},
			{Event: "close", Src: []string{"open"}, Dst: "closed"},
			{Event: "middle", Src: []string{"closed"}, Dst: "middle"},
		}),
	)
	events := fsm.CurrentAvailEvents()
	if !(slices.Contains(events, "middle") && slices.Contains(events, "open")) {
		t.Error("expected contain [middle, open] event with current state")
	}
	sortedEvents := fsm.SortedEvents()
	if !slices.Equal(sortedEvents, []string{"close", "middle", "open"}) {
		t.Error("expected sort event [close, middle, open] event with current state")
	}

	if !fsm.MatchCurrentAllOccur("middle", "open") {
		t.Error("expected contain all [middle, open] event with current state")
	}
	if fsm.MatchCurrentAllOccur("open", "close") {
		t.Error("expected not contain all [middle, open] event with current state")
	}
}
func Test_Fsm_InappropriateEvent(t *testing.T) {
	test_Fsm_InappropriateEvent(t, NewSafeFsm[string, string])
	test_Fsm_InappropriateEvent(t, NewFsm[string, string])
}

func test_Fsm_InappropriateEvent(t *testing.T, newFsm func(initState string, ts *Transition[string, string]) IFsm[string, string]) {
	fsm := newFsm(
		"closed",
		NewTransition([]Transform[string, string]{
			{Event: "open", Src: []string{"closed"}, Dst: "open"},
			{Event: "close", Src: []string{"open"}, Dst: "closed"},
		}),
	)

	err := fsm.Trigger("close")
	if err != ErrInappropriateEvent {
		t.Error("expected 'ErrInappropriateEvent' with correct state and event")
	}
}

func Test_Fsm_NonExistEvent(t *testing.T) {
	test_Fsm_NonExistEvent(t, NewSafeFsm[string, string])
	test_Fsm_NonExistEvent(t, NewFsm[string, string])
}

func test_Fsm_NonExistEvent(t *testing.T, newFsm func(initState string, ts *Transition[string, string]) IFsm[string, string]) {
	fsm := newFsm(
		"closed",
		NewTransition([]Transform[string, string]{
			{Event: "open", Src: []string{"closed"}, Dst: "open"},
			{Event: "close", Src: []string{"open"}, Dst: "closed"},
		}),
	)

	err := fsm.Trigger("lock")
	if err != ErrNonExistEvent {
		t.Error("expected 'ErrNonExistEvent' with incorrect event")
	}
}
func Test_Fsm_MultipleSources(t *testing.T) {
	testFsm_MultipleSources(t, NewSafeFsm[string, string])
	testFsm_MultipleSources(t, NewFsm[string, string])
}
func testFsm_MultipleSources(t *testing.T, newFsm func(initState string, ts *Transition[string, string]) IFsm[string, string]) {
	fsm := newFsm(
		"one",
		NewTransition([]Transform[string, string]{
			{Event: "first", Src: []string{"one"}, Dst: "two"},
			{Event: "second", Src: []string{"two"}, Dst: "three"},
			{Event: "reset", Src: []string{"one", "two", "three"}, Dst: "one"},
		}),
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
func Test_Fsm_MultipleEvents(t *testing.T) {
	test_Fsm_MultipleEvents(t, NewSafeFsm[string, string])
	test_Fsm_MultipleEvents(t, NewFsm[string, string])
}
func test_Fsm_MultipleEvents(t *testing.T, newFsm func(initState string, ts *Transition[string, string]) IFsm[string, string]) {
	fsm := newFsm(
		"start",
		NewTransition([]Transform[string, string]{
			{Event: "first", Src: []string{"start"}, Dst: "one"},
			{Event: "second", Src: []string{"start"}, Dst: "two"},
			{Event: "reset", Src: []string{"one"}, Dst: "reset_one"},
			{Event: "reset", Src: []string{"two"}, Dst: "reset_two"},
			{Event: "reset", Src: []string{"reset_one", "reset_two"}, Dst: "start"},
		}),
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

func Test_Fsm_TranslateError(t *testing.T) {
	test_Fsm_TranslateError(t, NewSafeFsm[string, string])
	test_Fsm_TranslateError(t, NewFsm[string, string])
}

type testTranslatorError struct{}

func (testTranslatorError) Translate(err error) error {
	switch err {
	case ErrInappropriateEvent:
		return errors.New("one")
	case ErrNonExistEvent:
		return errors.New("two")
	}
	return err
}

func test_Fsm_TranslateError(t *testing.T, newFsm func(initState string, ts *Transition[string, string]) IFsm[string, string]) {
	fsm := newFsm(
		"closed",
		NewTransitionBuilder(
			[]Transform[string, string]{
				{Event: "open", Src: []string{"closed"}, Dst: "open"},
				{Event: "close", Src: []string{"open"}, Dst: "closed"},
			}).
			TranslatorError(&testTranslatorError{}).
			Build(),
	)
	err := fsm.Trigger("close")
	if err == nil {
		t.Error("expected a error")
	} else {
		if !strings.Contains(err.Error(), "one") {
			t.Error("expected a error <one>")
		}
	}
	err = fsm.Trigger("nosupport")
	if err == nil {
		t.Error("expected a error")
	} else {
		if !strings.Contains(err.Error(), "two") {
			t.Error("expected a error <two>")
		}
	}
}
