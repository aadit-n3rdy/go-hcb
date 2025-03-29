package replica

import (
	"errors"
	"fmt"
	"time"

	"github.com/aadit-n3rdy/hotCrossBuns/ds"
)

// function to inform that replica is leader
func (r *Replica) InformofLeaderPosition(req *ds.LeaderBroadcast, rep *bool) error {
	if req == nil {
		return errors.New("Informing replica of leader position failed as request should not be nil")
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

	return nil
}

// function to hand over leader position along with other data
func (r *Replica) handOverLeader() {

	for {
		// if r is a leader,not in critical section and successors exist
		if r.isLeader {
			//fmt.Println("Checking if handover should be done")
			if !r.busy && len(r.successors) != 0 {

				// logging handing over
				fmt.Println("Handing over started")

				// get connection of next successor
				client := r.returnConnection(r.successors[0])

				// renounce leader position
				r.isLeader = false
				r.leaderID = r.successors[0]

				// make next replica a leader and send it successor information, and commit history
				reply := false
				client.Call("Replica.BecomeLeader", &ds.FutureGenerations{Successors: r.successors[1:]}, &reply)
				client.Call("Replica.SendCommitHistory", r.branch, &reply)

			}

			time.Sleep(time.Second * 2)

		}
	}

}

// function to send replica commit history
func (r *Replica) SendCommitHistory(req *ds.Branch, rep *bool) error {
	if req == nil {
		return errors.New("SendCommitHistory failed because req is nil")
	}

	// replace the current commit history if there are more commits on the sent branch
	if req.Tail.ID > r.branch.Tail.ID {

		// replace the branch and view number
		r.branch = req
		r.viewNumber = r.branch.Tail.ID
		fmt.Println("Commit history overwritten")
	} else {
		fmt.Println("Commit history remains same")
	}

	return nil
}

// function to inform next leader of it succeeding leaders
func (r *Replica) BecomeLeader(req *ds.FutureGenerations, rep *bool) error {
	if req == nil {
		return errors.New("GiveSuccessors failed as req is nil")
	}

	// declare self as leader
	r.leaderID = r.id
	r.isLeader = true
	r.successors = req.Successors

	// inform other replicas of leader position
	for _, rep := range r.replicaList {
		reply := false
		err := rep.Client.Call("Replica.InformofLeaderPosition", &ds.LeaderBroadcast{LId: r.id}, &reply)
		if err != nil {
			return errors.New("Failed to inform replica of leadership status")
		}
	}

	// store successor list in replica and set reply to true
	r.successors = req.Successors
	*rep = true

	return nil
}

// function to add execute a command
func (r *Replica) ExecuteCommand(req *string, rep *bool) error {
	if req == nil {
		return errors.New("ExecuteCommand failed as req is nil")
	}

	// print command execution statement
	fmt.Println("Executing Command: ", *req)

	// store the command into the node array and set bool to true
	r.cmds = append(r.cmds, *req)
	*rep = true

	// if length of node array is 3 then commit
	if len(r.cmds) == 3 {
		r.Commit(rep, rep)
	}

	return nil
}

// function to pull changes
func (r *Replica) Pull(req *bool, rep *bool) error {
	fmt.Println("Pulling Changes")

	// attempt leadership if not leader
	if !r.isLeader {

		// attempt to become leader
		fmt.Println("Attempting to become leader")

		leader := r.returnConnection(r.leaderID)
		reply := false

		// attempt leadership
		err := leader.Call("Replica.AttemptLeadership", &ds.LeaderAttempt{Id: r.id}, &reply)
		if err != nil {
			return errors.New("Failed to attempt leadership")
		}

		// wait till elected as leader
		for !r.isLeader {
			fmt.Println("Waiting to become leader")
			time.Sleep(time.Second * 2)
		}

		// inform other replicas of leader change and log it
		fmt.Println("Replica has become leader")
		for _, orep := range r.replicaList {
			orep.Client.Call("Replica.InformofLeaderPosition", &ds.LeaderBroadcast{LId: r.id}, &reply)
		}

	}

	fmt.Println("Changes pulled")

	return nil
}

// function to commit a node
func (r *Replica) Commit(req *bool, rep *bool) error {
	if !r.isLeader {
		r.Pull(req, req)
	}

	// lock the thread
	r.lock.Lock()
	defer r.lock.Unlock()

	// set replica to busy state
	r.busy = true

	// commiting a node to branch
	fmt.Println("Commiting changes")

	// increment the view number
	r.viewNumber += 1

	// create a node and store the current cmds array into it
	node := &ds.Node{ID: r.viewNumber, Cmd: r.cmds, Parent: r.branch.Tail}

	// add the node to the branch and clear the node array
	r.branch.Tail = node
	r.cmds = []string{}

	// set bool to true
	*rep = true

	// exit critical section
	r.busy = false

	return nil
}
