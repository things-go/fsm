# fsm

Finite State Machine 

[![Go.Dev reference](https://img.shields.io/badge/go.dev-reference-blue?logo=go&logoColor=white)](https://pkg.go.dev/github.com/things-go/fsm?tab=doc)
[![codecov](https://codecov.io/gh/things-go/fsm/branch/main/graph/badge.svg)](https://codecov.io/gh/things-go/fsm)
[![Tests](https://github.com/things-go/fsm/actions/workflows/ci.yml/badge.svg)](https://github.com/things-go/fsm/actions/workflows/ci.yml)
[![Go Report Card](https://goreportcard.com/badge/github.com/things-go/fsm)](https://goreportcard.com/report/github.com/things-go/fsm)
[![Licence](https://img.shields.io/github/license/things-go/fsm)](https://raw.githubusercontent.com/things-go/fsm/main/LICENSE)
[![Tag](https://img.shields.io/github/v/tag/things-go/fsm)](https://github.com/things-go/fsm/tags)


## Features

## Usage

### Installation

Use go get.
```bash
    go get github.com/things-go/fsm
```

Then import the package into your own code.
```bash
    import "github.com/things-go/fsm"
```

### Example

[embedmd]:# (examples/generic.go go)
```go
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
	f := fsm.NewSafeFsm[MyEvent, MyState](
		IsClosed,
		fsm.NewTransition([]fsm.Transform[MyEvent, MyState]{
			{Event: Open, Src: []MyState{IsClosed}, Dst: IsOpen},
			{Event: Close, Src: []MyState{IsOpen}, Dst: IsClosed},
		}),
	)
	fmt.Println(f.Current())
	err := f.Trigger(Open)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(f.Current())
	err = f.Trigger(Close)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(f.Current())
	// Output:
	// closed
	// opened
	// closed
	fmt.Println(fsm.VisualizeGraphviz[MyEvent, MyState](f))
	// digraph fsm {
	//    "closed" -> "opened" [ label = "open" ];
	//    "opened" -> "closed" [ label = "close" ];
	//
	//    "closed";
	//    "opened";
	// }
}
```

## License

This project is under MIT License. See the [LICENSE](LICENSE) file for the full license text.
