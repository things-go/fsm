//go:build ignore
// +build ignore

package main

import (
	"fmt"

	"github.com/things-go/fsm"
)

func main() {
	f := fsm.New(
		"idle",
		fsm.Transforms[string, string]{
			{Event: "scan", Src: []string{"idle"}, Dst: "scanning"},
			{Event: "working", Src: []string{"scanning"}, Dst: "scanning"},
			{Event: "situation", Src: []string{"scanning"}, Dst: "scanning"},
			{Event: "situation", Src: []string{"idle"}, Dst: "idle"},
			{Event: "finish", Src: []string{"scanning"}, Dst: "idle"},
		},
	)
	fmt.Println(f.Current())

	err := f.Trigger("scan")
	if err != nil {
		fmt.Println(err)
	}

	fmt.Println("1:" + f.Current())

	err = f.Trigger("working")
	if err != nil {
		fmt.Println(err)
	}

	fmt.Println("2:" + f.Current())

	err = f.Trigger("situation")
	if err != nil {
		fmt.Println(err)
	}

	fmt.Println("3:" + f.Current())

	err = f.Trigger("finish")
	if err != nil {
		fmt.Println(err)
	}

	fmt.Println("4:" + f.Current())
}
