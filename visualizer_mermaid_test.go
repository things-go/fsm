package fsm

import (
	"fmt"
	"strings"
	"testing"
)

func Test_MermaidStateDiagram(t *testing.T) {
	fsmUnderTest := NewSafeFsm[LampEvent, LampStatus](
		LampStatus_Closed,
		NewTransition([]Transform[LampEvent, LampStatus]{
			{Event: LampEvent_Open, Src: []LampStatus{LampStatus_Closed}, Dst: LampStatus_Opened},
			{Event: LampEvent_Close, Src: []LampStatus{LampStatus_Opened}, Dst: LampStatus_Closed},
			{Event: LampEvent_PartialClose, Src: []LampStatus{LampStatus_Intermediate}, Dst: LampStatus_Closed},
		}),
	)
	got, err := VisualizeMermaid[LampEvent, LampStatus](StateDiagram, fsmUnderTest)
	if err != nil {
		t.Errorf("got error for visualizing with type MERMAID: %s", err)
	}
	wanted := `
stateDiagram-v2
    [*] --> closed
    closed --> opened: open
    intermediate --> closed: partial-close
    opened --> closed: close
`
	normalizedGot := strings.ReplaceAll(got, "\n", "")
	normalizedWanted := strings.ReplaceAll(wanted, "\n", "")
	if normalizedGot != normalizedWanted {
		t.Errorf("build mermaid graph failed. \nwanted \n%s\nand got \n%s\n", wanted, got)
		fmt.Println(normalizedGot)
		fmt.Println(normalizedWanted)
	}
}

func Test_MermaidStateDiagram_CustomName(t *testing.T) {
	fsmUnderTest := NewSafeFsm[LampEvent, LampStatus](
		LampStatus_Closed,
		NewTransitionBuilder([]Transform[LampEvent, LampStatus]{
			{Event: LampEvent_Open, Src: []LampStatus{LampStatus_Closed}, Dst: LampStatus_Opened},
			{Event: LampEvent_Close, Src: []LampStatus{LampStatus_Opened}, Dst: LampStatus_Closed},
			{Event: LampEvent_PartialClose, Src: []LampStatus{LampStatus_Intermediate}, Dst: LampStatus_Closed},
		}).
			Name("Lamp FSM").
			Build(),
	)
	got, err := VisualizeMermaid[LampEvent, LampStatus](StateDiagram, fsmUnderTest)
	if err != nil {
		t.Errorf("got error for visualizing with type MERMAID: %s", err)
	}
	wanted := `
---
title: Lamp FSM
---
stateDiagram-v2
    [*] --> closed
    closed --> opened: open
    intermediate --> closed: partial-close
    opened --> closed: close
`
	normalizedGot := strings.ReplaceAll(got, "\n", "")
	normalizedWanted := strings.ReplaceAll(wanted, "\n", "")
	if normalizedGot != normalizedWanted {
		t.Errorf("build mermaid graph failed. \nwanted \n%s\nand got \n%s\n", wanted, got)
		fmt.Println(normalizedGot)
		fmt.Println(normalizedWanted)
	}
}

func Test_MermaidFlowChart(t *testing.T) {
	fsmUnderTest := NewSafeFsm[LampEvent, LampStatus](
		LampStatus_Closed,
		NewTransition([]Transform[LampEvent, LampStatus]{
			{Event: LampEvent_Open, Src: []LampStatus{LampStatus_Closed}, Dst: LampStatus_Opened},
			{Event: LampEvent_PartialOpen, Src: []LampStatus{LampStatus_Closed}, Dst: LampStatus_Intermediate},
			{Event: LampEvent_PartialOpen, Src: []LampStatus{LampStatus_Intermediate}, Dst: LampStatus_Opened},
			{Event: LampEvent_Close, Src: []LampStatus{LampStatus_Opened}, Dst: LampStatus_Closed},
			{Event: LampEvent_PartialClose, Src: []LampStatus{LampStatus_Intermediate}, Dst: LampStatus_Closed},
		}),
	)
	got, err := VisualizeMermaid[LampEvent, LampStatus](FlowChart, fsmUnderTest)
	if err != nil {
		t.Errorf("got error for visualizing with type MERMAID: %s", err)
	}
	wanted := `
graph LR
    id0[closed]
    id1[intermediate]
    id2[opened]

    id0 --> |open| id2
    id0 --> |partial-open| id1
    id1 --> |partial-close| id0
    id1 --> |partial-open| id2
    id2 --> |close| id0

    style id0 fill:#00AA00
`
	normalizedGot := strings.ReplaceAll(got, "\n", "")
	normalizedWanted := strings.ReplaceAll(wanted, "\n", "")
	if normalizedGot != normalizedWanted {
		t.Errorf("build mermaid graph failed. \nwanted \n%s\nand got \n%s\n", wanted, got)
		fmt.Println(normalizedGot)
		fmt.Println(normalizedWanted)
	}
}

func Test_MermaidFlowChart_CustomName(t *testing.T) {
	fsmUnderTest := NewSafeFsm[LampEvent, LampStatus](
		LampStatus_Closed,
		NewTransitionBuilder([]Transform[LampEvent, LampStatus]{
			{Event: LampEvent_Open, Src: []LampStatus{LampStatus_Closed}, Dst: LampStatus_Opened},
			{Event: LampEvent_PartialOpen, Src: []LampStatus{LampStatus_Closed}, Dst: LampStatus_Intermediate},
			{Event: LampEvent_PartialOpen, Src: []LampStatus{LampStatus_Intermediate}, Dst: LampStatus_Opened},
			{Event: LampEvent_Close, Src: []LampStatus{LampStatus_Opened}, Dst: LampStatus_Closed},
			{Event: LampEvent_PartialClose, Src: []LampStatus{LampStatus_Intermediate}, Dst: LampStatus_Closed},
		}).
			Name("Lamp FSM").
			Build(),
	)
	got, err := VisualizeMermaid[LampEvent, LampStatus](FlowChart, fsmUnderTest)
	if err != nil {
		t.Errorf("got error for visualizing with type MERMAID: %s", err)
	}
	wanted := `
---
title: Lamp FSM
---
graph LR
    id0[closed]
    id1[intermediate]
    id2[opened]

    id0 --> |open| id2
    id0 --> |partial-open| id1
    id1 --> |partial-close| id0
    id1 --> |partial-open| id2
    id2 --> |close| id0

    style id0 fill:#00AA00
`
	normalizedGot := strings.ReplaceAll(got, "\n", "")
	normalizedWanted := strings.ReplaceAll(wanted, "\n", "")
	if normalizedGot != normalizedWanted {
		t.Errorf("build mermaid graph failed. \nwanted \n%s\nand got \n%s\n", wanted, got)
		fmt.Println([]byte(normalizedGot))
		fmt.Println([]byte(normalizedWanted))
	}
}
