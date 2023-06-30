package fsm

import (
	"bytes"
	"fmt"
	"strings"

	"golang.org/x/exp/constraints"
)

// Visualize outputs a visualization of a FSM in Graphviz format.
func Visualize[E constraints.Ordered, S constraints.Ordered](fsm Visualizer[E, S]) string {
	var buf strings.Builder

	// we sort the key alphabetically to have a reproducible graph output
	sortedEKeys := getSortedTransitionKeys(fsm.Trans())
	sortedStateKeys, _ := getSortedStates(fsm.Trans())

	writeHeaderLine(&buf)
	writeTransitions(&buf, fsm.Current(), sortedEKeys, fsm.Trans())
	writeStates(&buf, sortedStateKeys)
	writeFooter(&buf)

	return buf.String()
}

func writeHeaderLine(buf *strings.Builder) {
	buf.WriteString("digraph fsm {")
	buf.WriteString("\n")
}

func writeTransitions[E constraints.Ordered, S constraints.Ordered](buf *strings.Builder, current S, sortedEKeys []eKey[E, S], transitions map[eKey[E, S]]S) {
	b := bytes.Buffer{}
	// make sure the current state is at top
	for _, k := range sortedEKeys {
		v := transitions[k]
		line := fmt.Sprintf(`    "%v" -> "%v" [ label = "%v" ];`, k.src, v, k.event)
		if k.src == current {
			buf.WriteString(line)
		} else {
			b.WriteString(line)
		}
		b.WriteString("\n")
	}
	if b.Len() > 0 {
		buf.Write(b.Bytes())
	}
	buf.WriteString("\n")
}

func writeStates[S constraints.Ordered](buf *strings.Builder, sortedStateKeys []S) {
	for _, k := range sortedStateKeys {
		buf.WriteString(fmt.Sprintf(`    "%v";`, k))
		buf.WriteString("\n")
	}
}

func writeFooter(buf *strings.Builder) {
	buf.WriteString(fmt.Sprintln("}"))
}
