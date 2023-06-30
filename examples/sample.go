//go:build ignore
// +build ignore

package main

import (
	"fmt"

	"github.com/things-go/fsm"
)

func main() {
	fsm1 := fsm.NewSafeFsm[string, string](
		"closed",
		fsm.NewTransition([]fsm.Transform[string, string]{
			{Event: "open", Src: []string{"closed"}, Dst: "open"},
			{Event: "close", Src: []string{"open"}, Dst: "closed"},
		}),
	)
	fmt.Println(fsm1.Current())

	err := fsm1.Trigger("open")
	if err != nil {
		fmt.Println(err)
	}

	fmt.Println(fsm1.Current())

	err = fsm1.Trigger("close")
	if err != nil {
		fmt.Println(err)
	}

	fmt.Println(fsm1.Current())
}
