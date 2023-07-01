package fsm

import (
	"bytes"
	"fmt"
	"strings"

	"golang.org/x/exp/constraints"
)

// VisualizeGraphviz outputs a visualization of a Fsm in Graphviz format.
func VisualizeGraphviz[E constraints.Ordered, S constraints.Ordered](fsm Visualizer[E, S]) (string, error) {
	v := newVisualizeGraphvizBuilder(fsm).
		writeHeaderLine().
		writeTransitions().
		writeStates().
		writeFooter()
	if v.Err() != nil {
		return "", v.Err()
	}
	return v.String(), nil
}

type visualizeGraphvizBuilder[E constraints.Ordered, S constraints.Ordered] struct {
	fsm                  Visualizer[E, S]
	sortedTriggerSources []TriggerSource[E, S] // we sort the key alphabetically to have a reproducible graph output
	sortedStates         []S
	buf                  strings.Builder
	err                  error
}

func newVisualizeGraphvizBuilder[E constraints.Ordered, S constraints.Ordered](fsm Visualizer[E, S]) *visualizeGraphvizBuilder[E, S] {
	return &visualizeGraphvizBuilder[E, S]{
		fsm:                  fsm,
		sortedTriggerSources: fsm.SortedTriggerSource(),
		sortedStates:         fsm.SortedStates(),
	}
}

func (v *visualizeGraphvizBuilder[E, S]) writeHeaderLine() *visualizeGraphvizBuilder[E, S] {
	if v.err != nil {
		return v
	}
	v.buf.WriteString("digraph fsm {")
	v.buf.WriteString("\n")
	if v.fsm.Name() != "" {
		v.buf.WriteString(fmt.Sprintf(`    label="%s"`, v.fsm.Name()))
		v.buf.WriteString("\n")
	}
	return v
}

func (v *visualizeGraphvizBuilder[E, S]) writeTransitions() *visualizeGraphvizBuilder[E, S] {
	if v.err != nil {
		return v
	}
	b := bytes.Buffer{}
	// make sure the current state is at top
	for _, ts := range v.sortedTriggerSources {
		dst, err := v.fsm.Transform(ts.State(), ts.Event())
		if err != nil {
			return v.setErr(err)
		}
		line := fmt.Sprintf(`    "%v" -> "%v" [ label = "%v" ];`, v.fsm.StateName(ts.State()), v.fsm.StateName(dst), v.fsm.EventName(ts.Event()))
		if ts.State() == v.fsm.Current() {
			v.buf.WriteString(line)
		} else {
			b.WriteString(line)
		}
		b.WriteString("\n")
	}
	if b.Len() > 0 {
		v.buf.Write(b.Bytes())
	}
	v.buf.WriteString("\n")
	return v
}

func (v *visualizeGraphvizBuilder[E, S]) writeStates() *visualizeGraphvizBuilder[E, S] {
	if v.err != nil {
		return v
	}
	for _, state := range v.sortedStates {
		v.buf.WriteString(fmt.Sprintf(`    "%v";`, v.fsm.StateName(state)))
		v.buf.WriteString("\n")
	}
	return v
}

func (v *visualizeGraphvizBuilder[E, S]) writeFooter() *visualizeGraphvizBuilder[E, S] {
	if v.err != nil {
		return v
	}
	v.buf.WriteString(fmt.Sprintln("}"))
	return v
}

func (v *visualizeGraphvizBuilder[E, S]) setErr(err error) *visualizeGraphvizBuilder[E, S] {
	v.err = err
	return v
}

func (v *visualizeGraphvizBuilder[E, S]) Err() error {
	return v.err
}

func (v *visualizeGraphvizBuilder[E, S]) String() string {
	return v.buf.String()
}
