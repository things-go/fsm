package fsm

import (
	"fmt"
	"strings"
	"testing"
)

func TestGraphvizOutput(t *testing.T) {
	fsmUnderTest := New(
		"closed",
		Transforms[string, string]{
			{Event: "open", Src: []string{"closed"}, Dst: "open"},
			{Event: "close", Src: []string{"open"}, Dst: "closed"},
			{Event: "part-close", Src: []string{"intermediate"}, Dst: "closed"},
		},
	)
	got := Visualize[string, string](fsmUnderTest)
	wanted := `
digraph fsm {
    "closed" -> "open" [ label = "open" ];
    "intermediate" -> "closed" [ label = "part-close" ];
    "open" -> "closed" [ label = "close" ];

    "closed";
    "intermediate";
    "open";
}`
	normalizedGot := strings.ReplaceAll(got, "\n", "")
	normalizedWanted := strings.ReplaceAll(wanted, "\n", "")
	if normalizedGot != normalizedWanted {
		t.Errorf("build graphivz graph failed. \nwanted \n%s\nand got \n%s\n", wanted, got)
		fmt.Println([]byte(normalizedGot))
		fmt.Println([]byte(normalizedWanted))
	}
}
