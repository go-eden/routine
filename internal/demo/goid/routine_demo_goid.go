package main

import (
	"fmt"
	"github.com/go-eden/routine"
)

func main() {
	goid := routine.Goid()
	fmt.Printf("curr goid: %d\n", goid)
}
