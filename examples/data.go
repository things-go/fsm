//go:build ignore
// +build ignore

package main

import (
	"fmt"

	"github.com/things-go/fsm"
)

func main() {
	fsm1 := fsm.NewSafeFsm[string, string](
		"idle",
		fsm.NewTransition([]fsm.Transform[string, string]{
			{Event: "produce", Src: []string{"idle"}, Dst: "idle"},
			{Event: "consume", Src: []string{"idle"}, Dst: "idle"},
		}),
	)
	fmt.Println(fsm1.Current())

	err := fsm1.Trigger("produce")
	if err != nil {
		fmt.Println(err)
	}

	fmt.Println(fsm1.Current())

	err = fsm1.Trigger("consume")
	if err != nil {
		fmt.Println(err)
	}

	fmt.Println(fsm1.Current())
}
