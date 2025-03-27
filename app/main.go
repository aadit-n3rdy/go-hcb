package main

import (
	"flag"

	"github.com/aadit-n3rdy/hotCrossBuns/ds"
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

	//NOTE: temp code to send leader id and hand over leadership
	if rep.ReturnID() == 0 {
		reply := new(bool)
		for _, orep := range rep.ReturnOtherReplicas() {
			orep.Client.Call("Replica.InformLeader", &ds.LeaderBroadcast{LId: rep.ReturnID()}, &reply)
		}
	} else {
		reply := new(bool)

		//NOTE: attempt succession if you are not leader
		for _, orep := range rep.ReturnOtherReplicas() {
			if orep.Id == rep.ReturnLeaderID() {
				orep.Client.Call("Replica.AttemptLeadership", &ds.LeaderAttempt{Id: rep.ReturnID()}, &reply)
			}
		}
	}

	// infinite loop
	for {
	}

}
