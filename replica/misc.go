package replica

import (
	"fmt"
	"net/rpc"
)

// function to return list of other replicas
func (r *Replica) ReturnOtherReplicas() []*OtherReplica {
	return r.replicaList
}

// function to return id of replica
func (r *Replica) ReturnID() int {
	return r.id
}

// function to return id of leader replica
func (r *Replica) ReturnLeaderID() int {
	return r.leaderID
}

// function to return connection of a replica
func (r *Replica) returnConnection(id int) *rpc.Client {
	for _, rep := range r.replicaList {
		if rep.Id == id {
			return rep.Client
		}
	}

	return nil
}

// function to print command history
func (r *Replica) commandHist() {
	itr := r.branch.Tail

	// print the commands onto the terminal
	for itr != nil {
		fmt.Println(itr.ID, itr.Cmd)
		itr = itr.Parent
	}
}

// function to print menu driver code
func (r *Replica) MenuDriver() {
	inp := new(string)

	for {
		fmt.Scanf("%s", inp)

		switch *inp {
		case "leader":
			fmt.Println("Leader ID: ", r.leaderID)

		case "id":
			fmt.Println("Replica ID: ", r.id)

		case "hist":
			r.commandHist()

		case "view":
			fmt.Println(r.viewNumber)

		case "h":
			fmt.Println("Help")
			fmt.Println("hist: Command history")
			fmt.Println("leader: Leader Id")
			fmt.Println("id: ID of replica")
			fmt.Println("view: Current View Number")
		}
	}
}
