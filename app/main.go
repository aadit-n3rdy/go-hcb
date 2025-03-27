package main

import (
	"flag"

	"github.com/aadit-n3rdy/hotCrossBuns/replica"
)

func main() {
	// set up the command line flag
	numbPtr := flag.Int("id", 0, "Stores id of the replica")
	flag.Parse()

	rep := new(replica.Replica)

	// start the replica and wait till all replicas have connected
	rep.New(*numbPtr)
	rep.Start()

	go rep.MenuDriver()

	// infinite loop
	for {
	}

}
