//go:build ignore
// +build ignore

package main

import (
	"fmt"

	"github.com/things-go/fsm"
)

type MyEvent int

const (
	Close MyEvent = 1
	Open  MyEvent = 2
)

func (c MyEvent) String() string {
	switch c {
	case 1:
		return "close"
	case 2:
		return "open"
	default:
		return "none"
	}
}

type MyState int

func (c MyState) String() string {
	switch c {
	case 1:
		return "closed"
	case 2:
		return "opened"
	default:
		return "none"
	}
}

const (
	IsClosed MyState = 1
	IsOpen   MyState = 2
)

func main() {
	fsm1 := fsm.New(
		IsClosed,
		fsm.Transforms[MyEvent, MyState]{
			{Event: Open, Src: []MyState{IsClosed}, Dst: IsOpen},
			{Event: Close, Src: []MyState{IsOpen}, Dst: IsClosed},
		},
	)
	fmt.Println(fsm1.Current())
	err := fsm1.Trigger(Open)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(fsm1.Current())
	err = fsm1.Trigger(Close)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(fsm1.Current())
	// Output:
	// closed
	// open
	// closed
}
