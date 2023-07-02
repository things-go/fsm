package fsm

import (
	"fmt"
	"strings"

	"golang.org/x/exp/constraints"
)

const highlightingColor = "#00AA00"

// MermaidType the type of the mermaid diagram type
type MermaidType string

const (
	// FlowChart the diagram type for output in flowchart style (https://mermaid-js.github.io/mermaid/#/flowchart) (including current state)
	FlowChart MermaidType = "flowChart"
	// StateDiagram the diagram type for output in stateDiagram style (https://mermaid-js.github.io/mermaid/#/stateDiagram)
	StateDiagram MermaidType = "stateDiagram"
)

// VisualizeMermaid outputs a visualization of a Fsm in Mermaid format as specified by the graphType.
func VisualizeMermaid[E constraints.Ordered, S constraints.Ordered](t MermaidType, fsm Visualizer[E, S]) (string, error) {
	switch t {
	case FlowChart:
		return visualizeMermaidFlowChart(fsm)
	case StateDiagram:
		return visualizeMermaidStateDiagram(fsm)
	default:
		return "", fmt.Errorf("unknown MermaidDiagramType: %s", t)
	}
}

func visualizeMermaidStateDiagram[E constraints.Ordered, S constraints.Ordered](fsm Visualizer[E, S]) (string, error) {
	sortedTriggerSources := fsm.SortedTriggerSource()
	buf := strings.Builder{}
	if fsm.Name() != "" {
		buf.WriteString("---\n")
		buf.WriteString(fmt.Sprintf("title: %s\n", fsm.Name()))
		buf.WriteString("---\n")
	}
	buf.WriteString("stateDiagram-v2\n")
	buf.WriteString(fmt.Sprintln(`    [*] -->`, fsm.StateName(fsm.Current())))
	for _, ts := range sortedTriggerSources {
		dst, err := fsm.Transform(ts.State(), ts.Event())
		if err != nil {
			return "", err
		}
		buf.WriteString(fmt.Sprintf(`    %s --> %s: %s`, fsm.StateName(ts.State()), fsm.StateName(dst), fsm.EventName(ts.Event())))
		buf.WriteString("\n")
	}
	return buf.String(), nil
}

// visualizeMermaidFlowChart outputs a visualization of a Fsm in Mermaid format (including highlighting of current state).
func visualizeMermaidFlowChart[E constraints.Ordered, S constraints.Ordered](fsm Visualizer[E, S]) (string, error) {
	v := newVisualizeMermaidFlowChartBuilder(fsm).
		writeFlowChartGraphType().
		writeFlowChartStates().
		writeFlowChartTransitions().
		writeFlowChartHighlightCurrent()
	if v.Err() != nil {
		return "", v.Err()
	}
	return v.String(), nil
}

type visualizeMermaidFlowChartBuilder[E constraints.Ordered, S constraints.Ordered] struct {
	fsm                  Visualizer[E, S]
	sortedTriggerSources []TriggerSource[E, S] // we sort the key alphabetically to have a reproducible graph output
	sortedStates         []S
	statesId             map[S]string
	buf                  strings.Builder
	err                  error
}

func newVisualizeMermaidFlowChartBuilder[E constraints.Ordered, S constraints.Ordered](fsm Visualizer[E, S]) *visualizeMermaidFlowChartBuilder[E, S] {
	sortedTriggerSources := fsm.SortedTriggerSource()
	sortedStates := fsm.SortedStates()
	statesId := intoSortedStateId(sortedStates)
	return &visualizeMermaidFlowChartBuilder[E, S]{
		fsm:                  fsm,
		sortedTriggerSources: sortedTriggerSources,
		sortedStates:         sortedStates,
		statesId:             statesId,
	}
}

func (v *visualizeMermaidFlowChartBuilder[E, S]) writeFlowChartGraphType() *visualizeMermaidFlowChartBuilder[E, S] {
	if v.err != nil {
		return v
	}
	if v.fsm.Name() != "" {
		v.buf.WriteString("---\n")
		v.buf.WriteString(fmt.Sprintf("title: %s\n", v.fsm.Name()))
		v.buf.WriteString("---\n")
	}
	v.buf.WriteString("graph LR\n")
	return v
}

func (v *visualizeMermaidFlowChartBuilder[E, S]) writeFlowChartStates() *visualizeMermaidFlowChartBuilder[E, S] {
	if v.err != nil {
		return v
	}
	for _, state := range v.sortedStates {
		v.buf.WriteString(fmt.Sprintf(`    %s[%s]`, v.statesId[state], v.fsm.StateName(state)))
		v.buf.WriteString("\n")
	}
	v.buf.WriteString("\n")
	return v
}

func (v *visualizeMermaidFlowChartBuilder[E, S]) writeFlowChartTransitions() *visualizeMermaidFlowChartBuilder[E, S] {
	if v.err != nil {
		return v
	}
	for _, ts := range v.sortedTriggerSources {
		dst, err := v.fsm.Transform(ts.State(), ts.Event())
		if err != nil {
			return v.setErr(err)
		}
		v.buf.WriteString(fmt.Sprintf(`    %s --> |%s| %s`, v.statesId[ts.State()], v.fsm.EventName(ts.Event()), v.statesId[dst]))
		v.buf.WriteString("\n")
	}
	v.buf.WriteString("\n")
	return v
}

func (v *visualizeMermaidFlowChartBuilder[E, S]) writeFlowChartHighlightCurrent() *visualizeMermaidFlowChartBuilder[E, S] {
	if v.err != nil {
		return v
	}
	v.buf.WriteString(fmt.Sprintf(`    style %s fill:%s`, v.statesId[v.fsm.Current()], highlightingColor))
	v.buf.WriteString("\n")
	return v
}
func (v *visualizeMermaidFlowChartBuilder[E, S]) setErr(err error) *visualizeMermaidFlowChartBuilder[E, S] {
	v.err = err
	return v
}

func (v *visualizeMermaidFlowChartBuilder[E, S]) Err() error {
	return v.err
}

func (v *visualizeMermaidFlowChartBuilder[E, S]) String() string {
	return v.buf.String()
}

func intoSortedStateId[S constraints.Ordered](sortedStates []S) map[S]string {
	statesId := make(map[S]string)
	for i, state := range sortedStates {
		statesId[state] = fmt.Sprintf("id%d", i)
	}
	return statesId
}
