package replica

import (
	"fmt"
	"net"
	"net/rpc"
	"os"
	"time"

	"github.com/aadit-n3rdy/hotCrossBuns/ds"
)

// initialize the replica
func (r *Replica) New(rid int) {

	// open file for reading
	file, err := os.Open("replicas.txt")
	if err != nil {
		panic("Failed to read file")
	}
	defer file.Close()

	// read from file
	for {
		var id int
		var ipAddr string
		var port string

		// read a line from the file if it exists, otherwise exit the loop
		items, err := fmt.Fscanf(file, "%d %s %s\n", &id, &ipAddr, &port)
		if err != nil {
			fmt.Println("Finished reading file")
			break
		}

		if items != 3 {
			panic("Failed to parse line correctly")
		}

		// if id matches replica store the details into the replica object
		if id == rid {
			r.id = rid
			r.ipAddr = ipAddr
			r.port = port
		} else {
			// create an other object and append it to the array
			other := new(OtherReplica)
			other.Id = id
			other.IpAddr = ipAddr
			other.Port = port

			r.replicaList = append(r.replicaList, other)
		}
	}

	// if rid is 0, declare self as leader
	if r.id == 0 {
		r.isLeader = true
	} else {
		r.isLeader = false
	}

	// set view number to 0
	r.viewNumber = 0

	// create a branch and assign the root and tail pointers
	r.branch = new(ds.Branch)
	r.branch.Root = &ds.Node{ID: r.viewNumber, Cmd: "ROOT", Parent: nil}
	r.branch.Tail = r.branch.Root

}

// function to start the replica
func (r *Replica) Start() {

	// start listening
	go r.startServer()

	// sleep the thread to allow other goroutines to start
	time.Sleep(time.Second * 5)

	// buffer period to start rpc connections to the other replicas
	for !r.startMesh() {
		time.Sleep(time.Second * 2)
	}
}

// function to start server and register the rpc
func (r *Replica) startServer() {
	// listen on ip address and port specified
	listener, err := net.Listen("tcp", ":"+r.port)
	if err != nil {
		panic("Failed to listen on port")
	}

	// register the replica as an rpc server
	rpc.Register(r)

	// logging
	fmt.Printf("Replica %d is listening on %s:%s\n", r.id, r.ipAddr, r.port)

	// accept rpc connections
	for {
		rpc.Accept(listener)

	}
}

// function to connect to all other replicas in the list and establish the mesh
func (r *Replica) startMesh() bool {

	allConnected := true

	// establish a connection with all replicas
	for _, rep := range r.replicaList {

		// continue to the next replica if connection has already been established
		// and avoid duplicate connections by connecting to higher id replicas
		if rep.Client != nil {
			continue
		}

		address := rep.IpAddr + ":" + rep.Port
		client, err := rpc.Dial("tcp", address)
		if err != nil {
			fmt.Printf("Failed to connect to replica %d\n", rep.Id)
			allConnected = false
			continue
		}

		rep.Client = client
	}

	return allConnected
}
