package fsm

import (
	"fmt"
	"strings"
	"testing"
)

func Test_Graphviz(t *testing.T) {
	fsmUnderTest := NewFsm[LampEvent, LampStatus](
		LampStatus_Closed,
		NewTransition([]Transform[LampEvent, LampStatus]{
			{Event: LampEvent_Open, Src: []LampStatus{LampStatus_Closed}, Dst: LampStatus_Opened},
			{Event: LampEvent_Close, Src: []LampStatus{LampStatus_Opened}, Dst: LampStatus_Closed},
			{Event: LampEvent_PartialClose, Src: []LampStatus{LampStatus_Intermediate}, Dst: LampStatus_Closed},
		}),
	)
	got, err := fsmUnderTest.Visualize(Graphviz)
	if err != nil {
		panic(err)
	}
	wanted := `
digraph fsm {
    "closed" -> "opened" [ label = "open" ];
    "intermediate" -> "closed" [ label = "partial-close" ];
    "opened" -> "closed" [ label = "close" ];

    "closed";
    "intermediate";
    "opened";
}`
	normalizedGot := strings.ReplaceAll(got, "\n", "")
	normalizedWanted := strings.ReplaceAll(wanted, "\n", "")
	if normalizedGot != normalizedWanted {
		t.Errorf("build graphivz graph failed. \nwanted \n%s\nand got \n%s\n", wanted, got)
		fmt.Println(normalizedGot)
		fmt.Println(normalizedWanted)
	}
}

func Test_Graphviz_CustomName(t *testing.T) {
	fsmUnderTest := NewSafeFsm[LampEvent, LampStatus](
		LampStatus_Closed,
		NewTransitionBuilder([]Transform[LampEvent, LampStatus]{
			{Name: formatEvent(LampEvent_Open), Event: LampEvent_Open, Src: []LampStatus{LampStatus_Closed}, Dst: LampStatus_Opened},
			{Name: formatEvent(LampEvent_Close), Event: LampEvent_Close, Src: []LampStatus{LampStatus_Opened}, Dst: LampStatus_Closed},
			{Name: formatEvent(LampEvent_PartialClose), Event: LampEvent_PartialClose, Src: []LampStatus{LampStatus_Intermediate}, Dst: LampStatus_Closed},
		}).
			Name("Lamp FSM").
			StateNames(map[LampStatus]string{
				LampStatus_Intermediate: formatState(LampStatus_Intermediate),
				LampStatus_Opened:       formatState(LampStatus_Opened),
				LampStatus_Closed:       formatState(LampStatus_Closed),
			}).
			Build(),
	)

	got, err := fsmUnderTest.Visualize(Graphviz)
	if err != nil {
		panic(err)
	}
	wanted := `
digraph fsm {
    label="Lamp FSM"
    ">closed" -> ">opened" [ label = "<open>" ];
    ">intermediate" -> ">closed" [ label = "<partial-close>" ];
    ">opened" -> ">closed" [ label = "<close>" ];

    ">closed";
    ">intermediate";
    ">opened";
}`
	normalizedGot := strings.ReplaceAll(got, "\n", "")
	normalizedWanted := strings.ReplaceAll(wanted, "\n", "")
	if normalizedGot != normalizedWanted {
		t.Errorf("build graphivz graph failed. \nwanted \n%s\nand got \n%s\n", wanted, got)
		fmt.Println(normalizedGot)
		fmt.Println(normalizedWanted)
	}
}
