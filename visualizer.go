package fsm

import (
	"golang.org/x/exp/constraints"
)

// Visualize outputs a visualization of a Fsm in the desired format.
type Visualizer[E constraints.Ordered, S constraints.Ordered] interface {
	Current() S
	Name() string
	Transform(srcState S, event E) (dstState S, err error)
	SortedTriggerSource() []TriggerSource[E, S]
	SortedStates() []S
	SortedEvents() []E
	EventName(event E) string
	StateName(state S) string
}

// VisualizeType the type of the visualization
type VisualizeType string

const (
	// Graphviz the type for graphviz output (http://www.webgraphviz.com/)
	Graphviz VisualizeType = "graphviz"
	// Mermaid the type for mermaid output (https://mermaid-js.github.io/mermaid/#/stateDiagram) in the stateDiagram form
	Mermaid VisualizeType = "mermaid"
	// MermaidStateDiagram the type for mermaid output (https://mermaid-js.github.io/mermaid/#/stateDiagram) in the stateDiagram form
	MermaidStateDiagram VisualizeType = "mermaid-state-diagram"
	// MermaidFlowChart the type for mermaid output (https://mermaid-js.github.io/mermaid/#/flowchart) in the flow chart form
	MermaidFlowChart VisualizeType = "mermaid-flow-chart"
)

// Visualize outputs a visualization of a Fsm in the desired format.
// If the type is not given it defaults to Graphviz
func Visualize[E constraints.Ordered, S constraints.Ordered](t VisualizeType, fsm Visualizer[E, S]) (string, error) {
	switch t {
	case Mermaid, MermaidStateDiagram:
		return VisualizeMermaid(StateDiagram, fsm)
	case MermaidFlowChart:
		return VisualizeMermaid(FlowChart, fsm)
	case Graphviz:
		fallthrough
	default:
		return VisualizeGraphviz(fsm)
	}
}
