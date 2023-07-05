package fsm

import (
	"errors"
	"fmt"
	"strings"
	"testing"

	"golang.org/x/exp/slices"
)

type LampEvent string

const (
	LampEvent_Open         LampEvent = "open"
	LampEvent_Close        LampEvent = "close"
	LampEvent_PartialClose LampEvent = "partial-close"
	LampEvent_PartialOpen  LampEvent = "partial-open"
	LampEvent_Look         LampEvent = "look"
)

func (l LampEvent) String() string { return string(l) }

func formatEvent[E fmt.Stringer](e E) string {
	return fmt.Sprintf("<%s>", e.String())
}

type LampStatus string

const (
	LampStatus_Intermediate LampStatus = "intermediate"
	LampStatus_Opened       LampStatus = "opened"
	LampStatus_Closed       LampStatus = "closed"
)

func (l LampStatus) String() string { return string(l) }

func formatState[S fmt.Stringer](e S) string {
	return fmt.Sprintf(">%s", e.String())
}

type testTranslatorError struct{}

func (testTranslatorError) Translate(err error) error {
	switch err {
	case ErrInappropriateEvent:
		return errors.New("err1")
	case ErrNonExistEvent:
		return errors.New("err2")
	}
	return err
}

func Test_Fsm_Clone(t *testing.T) {
	test_Fsm_Clone(t, NewSafeFsm[LampEvent, LampStatus])
	test_Fsm_Clone(t, NewFsm[LampEvent, LampStatus])
}

func test_Fsm_Clone(t *testing.T, newFsm func(initState LampStatus, ts ITransition[LampEvent, LampStatus]) IFsm[LampEvent, LampStatus]) {
	fsm := newFsm(
		LampStatus_Closed,
		NewTransition([]Transform[LampEvent, LampStatus]{
			{Event: LampEvent_Open, Src: []LampStatus{LampStatus_Closed}, Dst: LampStatus_Opened},
			{Event: LampEvent_Close, Src: []LampStatus{LampStatus_Opened}, Dst: LampStatus_Closed},
		}),
	)
	fsm1 := fsm.Clone()
	if fsm1.Current() != fsm.Current() {
		t.Errorf("expected same current state")
	}
	fsm2 := fsm.CloneNewState(LampStatus_Opened)
	if fsm2.Current() != LampStatus_Opened {
		t.Error("expected state to be 'opened'")
	}
}

func Test_Fsm_SameState(t *testing.T) {
	test_Fsm_SameState(t, NewSafeFsm[LampEvent, LampStatus])
	test_Fsm_SameState(t, NewFsm[LampEvent, LampStatus])
}

func test_Fsm_SameState(t *testing.T, newFsm func(initState LampStatus, ts ITransition[LampEvent, LampStatus]) IFsm[LampEvent, LampStatus]) {
	fsm := newFsm(
		LampStatus_Closed,
		NewTransition([]Transform[LampEvent, LampStatus]{
			{Event: LampEvent_Close, Src: []LampStatus{LampStatus_Closed}, Dst: LampStatus_Closed},
		}),
	)
	err := fsm.Trigger(LampEvent_Close)
	if err != nil {
		t.Errorf("expected trigger no error")
	}
	if fsm.Current() != LampStatus_Closed {
		t.Error("expected state to be 'closed'")
	}
}

func Test_Fsm_State(t *testing.T) {
	test_Fsm_State(t, NewSafeFsm[LampEvent, LampStatus])
	test_Fsm_State(t, NewFsm[LampEvent, LampStatus])
}

func test_Fsm_State(t *testing.T, newFsm func(initState LampStatus, ts ITransition[LampEvent, LampStatus]) IFsm[LampEvent, LampStatus]) {
	fsm := newFsm(
		LampStatus_Closed,
		NewTransition([]Transform[LampEvent, LampStatus]{
			{Event: LampEvent_Close, Src: []LampStatus{LampStatus_Opened}, Dst: LampStatus_Closed},
			{Event: LampEvent_Look, Src: []LampStatus{LampStatus_Closed}, Dst: LampStatus_Closed},
		}),
	)

	if !fsm.ContainsEvent(LampEvent_Close) {
		t.Error("expected support event 'close'")
	}
	if fsm.ContainsEvent(LampEvent_PartialClose) {
		t.Error("expected not support event 'partial-close'")
	}
	if !fsm.ContainsAllEvent(LampEvent_Close, LampEvent_Look) {
		t.Error("expected support all event 'close' and 'look'")
	}
	if fsm.ContainsAllEvent(LampEvent_Close, LampEvent_PartialClose) {
		t.Error("expected not support all event 'close' and 'partial-close'")
	}

	if b, err := fsm.Match(LampStatus_Opened, LampStatus_Closed, LampEvent_Close); err != nil {
		t.Error("expected event src dst state match no error")
	} else {
		if !b {
			t.Error("expected event src dst state match")
		}
	}
	if b, err := fsm.Match(LampStatus_Opened, LampStatus_Closed, LampEvent_Look); err == nil {
		t.Error("expected event src dst state match has error")
	} else {
		if b {
			t.Error("expected event src dst state match false")
		}
	}
	fsm.SetCurrent(LampStatus_Opened)
	if fsm.Current() != LampStatus_Opened {
		t.Error("expected state to be 'opened'")
	}
	if !fsm.MatchCurrentOccur(LampEvent_Close) {
		t.Error("expected event can occur in the current state.")
	}

	err := fsm.Trigger(LampEvent_Close)
	if err != nil {
		t.Error("trigger is expected no error")
	}
	if !fsm.Is(LampStatus_Closed) {
		t.Error("expected state to be 'closed'")
	}
}

func Test_Fsm_Avail(t *testing.T) {
	test_Fsm_Avail(t, NewSafeFsm[LampEvent, LampStatus])
	test_Fsm_Avail(t, NewFsm[LampEvent, LampStatus])
}

func test_Fsm_Avail(t *testing.T, newFsm func(initState LampStatus, ts ITransition[LampEvent, LampStatus]) IFsm[LampEvent, LampStatus]) {
	fsm := newFsm(
		LampStatus_Closed,
		NewTransition([]Transform[LampEvent, LampStatus]{
			{Event: LampEvent_Open, Src: []LampStatus{LampStatus_Closed}, Dst: LampStatus_Opened},
			{Event: LampEvent_Close, Src: []LampStatus{LampStatus_Opened}, Dst: LampStatus_Closed},
			{Event: LampEvent_PartialOpen, Src: []LampStatus{LampStatus_Closed}, Dst: LampStatus_Intermediate},
		}),
	)
	events := fsm.CurrentAvailEvents()
	if !(slices.Contains(events, LampEvent_Open) && slices.Contains(events, LampEvent_PartialOpen)) {
		t.Error("expected contain [open, partial-open] event with current state")
	}
	sortedEvents := fsm.SortedEvents()
	if !slices.Equal(sortedEvents, []LampEvent{LampEvent_Close, LampEvent_Open, LampEvent_PartialOpen}) {
		t.Error("expected sort event [close, open, partial-open] event with current state")
	}
	availSourceStates := fsm.AvailSourceStates(LampEvent_PartialOpen)
	if !slices.Contains(availSourceStates, LampStatus_Closed) {
		t.Error("expected avail source state [closed] with the event")
	}
	if !fsm.MatchCurrentAllOccur(LampEvent_PartialOpen, LampEvent_Open) {
		t.Error("expected contain all [partial-open, open] event with current state")
	}
	if fsm.MatchCurrentAllOccur(LampEvent_PartialOpen, LampEvent_Close) {
		t.Error("expected not contain all [partial-open, open] event with current state")
	}
}
func Test_Fsm_NonExistEvent_InappropriateEvent(t *testing.T) {
	test_Fsm_NonExistEvent_InappropriateEvent(t, NewSafeFsm[LampEvent, LampStatus])
	test_Fsm_NonExistEvent_InappropriateEvent(t, NewFsm[LampEvent, LampStatus])
}

func test_Fsm_NonExistEvent_InappropriateEvent(t *testing.T, newFsm func(initState LampStatus, ts ITransition[LampEvent, LampStatus]) IFsm[LampEvent, LampStatus]) {
	fsm := newFsm(
		LampStatus_Closed,
		NewTransition([]Transform[LampEvent, LampStatus]{
			{Event: LampEvent_Open, Src: []LampStatus{LampStatus_Closed}, Dst: LampStatus_Opened},
			{Event: LampEvent_Close, Src: []LampStatus{LampStatus_Opened}, Dst: LampStatus_Closed},
		}),
	)

	err := fsm.Trigger(LampEvent_PartialClose)
	if err != ErrNonExistEvent {
		t.Error("expected 'ErrNonExistEvent' with incorrect event")
	}
	err = fsm.Trigger(LampEvent_Close)
	if err != ErrInappropriateEvent {
		t.Error("expected 'ErrInappropriateEvent' with correct state and event")
	}
}

func Test_Fsm_TranslateError(t *testing.T) {
	test_Fsm_TranslateError(t, NewSafeFsm[LampEvent, LampStatus])
	test_Fsm_TranslateError(t, NewFsm[LampEvent, LampStatus])
}

func test_Fsm_TranslateError(t *testing.T, newFsm func(initState LampStatus, ts ITransition[LampEvent, LampStatus]) IFsm[LampEvent, LampStatus]) {
	fsm := newFsm(
		LampStatus_Closed,
		NewTransitionBuilder(
			[]Transform[LampEvent, LampStatus]{
				{Event: LampEvent_Open, Src: []LampStatus{LampStatus_Closed}, Dst: LampStatus_Opened},
				{Event: LampEvent_Close, Src: []LampStatus{LampStatus_Opened}, Dst: LampStatus_Closed},
			}).
			TranslatorError(&testTranslatorError{}).
			Build(),
	)
	err := fsm.Trigger(LampEvent_Close)
	if err == nil {
		t.Error("expected a error")
	} else {
		if !strings.Contains(err.Error(), "err1") {
			t.Error("expected a error <err1>")
		}
	}
	err = fsm.Trigger(LampEvent_PartialClose)
	if err == nil {
		t.Error("expected a error")
	} else {
		if !strings.Contains(err.Error(), "err2") {
			t.Error("expected a error <err2>")
		}
	}

}

const (
	eventFirst  = "first"
	eventSecond = "second"
	eventReset  = "reset"
)

const (
	statusStart      = "start"
	statusOne        = "one"
	statusTwo        = "two"
	statusThree      = "three"
	statusResetOne   = "reset-one"
	statusResetTwo   = "reset-two"
	statusResetThree = "reset-three"
)

func Test_Fsm_MultipleSources(t *testing.T) {
	testFsm_MultipleSources(t, NewSafeFsm[string, string])
	testFsm_MultipleSources(t, NewFsm[string, string])
}
func testFsm_MultipleSources(t *testing.T, newFsm func(initState string, ts ITransition[string, string]) IFsm[string, string]) {
	fsm := newFsm(
		statusOne,
		NewTransition([]Transform[string, string]{
			{Event: eventFirst, Src: []string{statusOne}, Dst: statusTwo},
			{Event: eventSecond, Src: []string{statusTwo}, Dst: statusThree},
			{Event: eventReset, Src: []string{statusOne, statusTwo, statusThree}, Dst: statusOne},
		}),
	)

	err := fsm.Trigger(eventFirst)
	if err != nil {
		t.Errorf("trigger failed %v", err)
	}
	if fsm.Current() != statusTwo {
		t.Error("expected state to be 'two'")
	}
	err = fsm.Trigger(eventReset)
	if err != nil {
		t.Errorf("transition failed %v", err)
	}
	if fsm.Current() != statusOne {
		t.Error("expected state to be 'one'")
	}
	err = fsm.Trigger(eventFirst)
	if err != nil {
		t.Errorf("trigger failed %v", err)
	}
	err = fsm.Trigger(eventSecond)
	if err != nil {
		t.Errorf("trigger failed %v", err)
	}
	if fsm.Current() != statusThree {
		t.Errorf("expected state to be '%s'", statusThree)
	}
	err = fsm.Trigger(eventReset)
	if err != nil {
		t.Errorf("trigger failed %v", err)
	}
	if fsm.Current() != statusOne {
		t.Errorf("expected state to be '%s'", statusOne)
	}
}
func Test_Fsm_MultipleEvents(t *testing.T) {
	test_Fsm_MultipleEvents(t, NewSafeFsm[string, string])
	test_Fsm_MultipleEvents(t, NewFsm[string, string])
}
func test_Fsm_MultipleEvents(t *testing.T, newFsm func(initState string, ts ITransition[string, string]) IFsm[string, string]) {
	fsm := newFsm(
		statusStart,
		NewTransition([]Transform[string, string]{
			{Event: eventFirst, Src: []string{statusStart}, Dst: statusOne},
			{Event: eventSecond, Src: []string{statusStart}, Dst: statusTwo},
			{Event: eventReset, Src: []string{statusOne}, Dst: statusResetOne},
			{Event: eventReset, Src: []string{statusTwo}, Dst: statusResetTwo},
			{Event: eventReset, Src: []string{statusResetOne, statusResetTwo}, Dst: statusStart},
		}),
	)

	err := fsm.Trigger(eventFirst)
	if err != nil {
		t.Errorf("trigger failed %v", err)
	}
	err = fsm.Trigger(eventReset)
	if err != nil {
		t.Errorf("trigger failed %v", err)
	}
	if fsm.Current() != statusResetOne {
		t.Errorf("expected state to be '%s'", statusResetOne)
	}
	err = fsm.Trigger(eventReset)
	if err != nil {
		t.Errorf("trigger failed %v", err)
	}
	if fsm.Current() != statusStart {
		t.Errorf("expected state to be '%s'", statusStart)
	}

	err = fsm.Trigger(eventSecond)
	if err != nil {
		t.Errorf("trigger failed %v", err)
	}
	err = fsm.Trigger(eventReset)
	if err != nil {
		t.Errorf("trigger failed %v", err)
	}
	if fsm.Current() != statusResetTwo {
		t.Errorf("expected state to be '%s'", statusResetTwo)
	}
	err = fsm.Trigger(eventReset)
	if err != nil {
		t.Errorf("trigger failed %v", err)
	}
	if fsm.Current() != statusStart {
		t.Errorf("expected state to be '%s'", statusStart)
	}
}
