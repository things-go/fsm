package fsm

import (
	"fmt"
	"strings"
	"testing"
)

func Test_Graphviz(t *testing.T) {
	fsmUnderTest := NewFsm[string, string](
		"closed",
		NewTransition([]Transform[string, string]{
			{Event: "open", Src: []string{"closed"}, Dst: "open"},
			{Event: "close", Src: []string{"open"}, Dst: "closed"},
			{Event: "part-close", Src: []string{"intermediate"}, Dst: "closed"},
		}),
	)
	got, err := fsmUnderTest.Visualize(Graphviz)
	if err != nil {
		panic(err)
	}
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
		fmt.Println(normalizedGot)
		fmt.Println(normalizedWanted)
	}
}

func Test_Graphviz_CustomName(t *testing.T) {
	fsmUnderTest := NewSafeFsm[string, string](
		"closed",
		NewTransitionBuilder([]Transform[string, string]{
			{Name: "打开", Event: "open", Src: []string{"closed"}, Dst: "open"},
			{Name: "关闭", Event: "close", Src: []string{"open"}, Dst: "closed"},
			{Name: "部份关闭", Event: "part-close", Src: []string{"intermediate"}, Dst: "closed"},
		}).
			Name("开关灯状态转移").
			StateNames(map[string]string{
				"intermediate": "intermediate(初始态)",
				"closed":       "closed(关闭的)",
				"open":         "open(打开的)",
			}).
			Build(),
	)

	got, err := fsmUnderTest.Visualize(Graphviz)
	if err != nil {
		panic(err)
	}
	wanted := `
digraph fsm {
    label="开关灯状态转移"
    "closed(关闭的)" -> "open(打开的)" [ label = "打开" ];
    "intermediate(初始态)" -> "closed(关闭的)" [ label = "部份关闭" ];
    "open(打开的)" -> "closed(关闭的)" [ label = "关闭" ];

    "closed(关闭的)";
    "intermediate(初始态)";
    "open(打开的)";
}`
	normalizedGot := strings.ReplaceAll(got, "\n", "")
	normalizedWanted := strings.ReplaceAll(wanted, "\n", "")
	if normalizedGot != normalizedWanted {
		t.Errorf("build graphivz graph failed. \nwanted \n%s\nand got \n%s\n", wanted, got)
		fmt.Println(normalizedGot)
		fmt.Println(normalizedWanted)
	}
}
