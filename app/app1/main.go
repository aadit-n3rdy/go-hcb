package main

import (
	"fmt"

	"github.com/aadit-n3rdy/hotCrossBuns/replica"
)

func main() {
	rep := new(replica.Replica)

	rep.New(1)
	rep.Start()

	for _, rep := range rep.ReturnOtherReplicas() {
		fmt.Println(rep)
	}
}
