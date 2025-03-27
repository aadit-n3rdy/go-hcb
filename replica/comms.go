package replica

import (
	"errors"
	"fmt"

	"github.com/aadit-n3rdy/hotCrossBuns/ds"
)

// function to inform that replica is leader
func (r *Replica) InformLeader(req *ds.LeaderBroadcast, rep *bool) error {
	if req == nil {
		return errors.New("InformLeader failed as request should not be nil")
	}

	// store leader id and set bool to true
	r.leaderID = req.LId
	*rep = true

	//NOTE: temp code to print id of leader
	fmt.Printf("Replica %d has accepted %d as leader\n", r.id, req.LId)

	return nil
}

// function to attempt to become leader
func (r *Replica) AttemptLeadership(req *ds.LeaderAttempt, rep *bool) error {
	if req == nil {
		return errors.New("AttemptLeadership failed as request field was nil")
	}

	// append id of current replica to the successor list and set bool to true
	r.successors = append(r.successors, req.Id)
	*rep = true

	//NOTE: temp code to print list of successors
	fmt.Println("Successor list ", r.successors)

	return nil
}

// function to inform next leader of it succeeding leaders
func (r *Replica) GiveSuccessors(req *ds.FutureGenerations, rep *bool) error {
	if req == nil {
		return errors.New("GiveSuccessors failed as req is nil")
	}

	// store successor list in replica and set reply to true
	r.successors = req.Successors
	*rep = true

	return nil
}
