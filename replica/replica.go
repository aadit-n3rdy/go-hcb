package replica

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/aadit-n3rdy/hotCrossBuns/ds"
)

type OtherReplica struct {
	Id     int
	Conn   net.Conn
	IpAddr string
	Port   string
}

type Replica struct {
	id     int
	ipAddr string
	port   string

	replicaList []*OtherReplica

	isLeader bool
	leaderID int

	viewNumber int
	lockedQC   ds.QuroumCertificate
	prepareQC  ds.QuroumCertificate
	branch     *ds.Branch
}

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
	}

}

// function to start the replica
func (r *Replica) Start() {

	// start listening
	go r.startServer()

	// sleep the thread to allow other goroutines to start
	time.Sleep(time.Second * 5)

	// buffer period to start connnections to the other replicas
	for !r.startMesh() {
		time.Sleep(time.Second * 2)
	}

}

// function to return list of other replicas
func (r *Replica) ReturnOtherReplicas() []*OtherReplica {
	return r.replicaList
}

// function to start server
func (r *Replica) startServer() {
	// listen on ip address and port specified
	listener, err := net.Listen("tcp", ":"+r.port)
	if err != nil {
		panic("Failed to listen on port")
	}

	// logging
	fmt.Printf("Replica %d is listening on %s:%s\n", r.id, r.ipAddr, r.port)

	for {

		// accept a connection to the server
		conn, err := listener.Accept()
		if err != nil {
			panic("Failed to accept connection")
		}

		go r.handleConnection(conn)

	}
}

// function to handle connection
func (r *Replica) handleConnection(conn net.Conn) {
	reader := bufio.NewReader(conn)

	// read the id from the connection
	idStr, err := reader.ReadString('\n')
	if err != nil {
		log.Printf("Error reading replica id: %v", err)
		conn.Close()
		return
	}

	// convert string id to interger id
	idStr = strings.TrimSpace(idStr)
	id, err := strconv.Atoi(idStr)
	if err != nil {
		log.Printf("Error converting replica id: %v", err)
		conn.Close()
		return
	}

	// store connection with corresponding replica
	for _, rep := range r.replicaList {
		if rep.Id == id {
			rep.Conn = conn
			fmt.Printf("Connection established with replica %d\n", id)
			break
		}
	}
}

// function to connect to all other replicas in the list and establish the mesh
func (r *Replica) startMesh() bool {

	allConnected := true

	// establish a connection with all replicas
	for _, rep := range r.replicaList {

		// continue to the next replica if connection has already been established
		if rep.Conn != nil {
			continue
		}

		address := rep.IpAddr + ":" + rep.Port
		conn, err := net.Dial("tcp", address)
		if err != nil {
			fmt.Printf("Failed to connect to replica %d\n", rep.Id)
			allConnected = false
			continue
		}

		// send replica id to connected replica
		_, err = conn.Write([]byte(strconv.Itoa(r.id) + "\n"))
		if err != nil {
			fmt.Printf("Failed to send id to replica %d\n", rep.Id)
			allConnected = false
			continue
		}

		rep.Conn = conn

		// handle send and receive
		go r.handleSend(conn)
		go r.handleRecv(conn)
	}

	return allConnected
}

// function to handle sending of data
func (r *Replica) handleSend(conn net.Conn) {}

// function to handle reading of data
func (r *Replica) handleRecv(conn net.Conn) {}

// function to check if passed node extends from the lockedQC
func (r *Replica) safetyRule(node *ds.Node) bool {

	// node is in the future
	if node.ID > r.branch.Tail.ID {
		var itr *ds.Node
		itr = node

		for itr.ID > r.branch.Tail.ID {
			itr = itr.Parent
		}

		return itr.Cmd == node.Cmd
	}

	// node is in the past
	return false
}

// function to check if liveness rule is satisfied
func (r *Replica) livenessRule(qview int) bool {
	return qview > r.lockedQC.ViewNumber
}

// function to create a message
func (r *Replica) Msg(mtype int, node ds.Node, qc ds.QuroumCertificate) *ds.Message {
	msg := new(ds.Message)

	msg.Type = mtype
	msg.CurView = r.viewNumber
	msg.Node = node
	msg.Justify = qc

	return msg
}

// function to create a message and sign it
func (r *Replica) VoteMsg(mtype int, node ds.Node, qc ds.QuroumCertificate) *ds.Message {
	msg := new(ds.Message)

	msg.PartialSig[r.id] = 1

	return msg
}

// function to create a leaf for the current branch
func (r *Replica) CreateLeaf(parent *ds.Node, cmd string) *ds.Node {
	b := new(ds.Node)

	b.ID = r.viewNumber
	b.Parent = parent
	b.Cmd = cmd

	return b
}

// function to create a quorum certificate
func (r *Replica) QC(msg *ds.Message) *ds.QuroumCertificate {
	qc := new(ds.QuroumCertificate)

	qc.Type = msg.Type
	qc.ViewNumber = msg.CurView
	qc.Node = msg.Node
	qc.Sig = msg.PartialSig

	return qc
}

func (r *Replica) MatchingMessage(msg *ds.Message, mtype int, mview int) bool {
	return (msg.CurView == mview) && (msg.Type == mtype)
}

func (r *Replica) MatchingQC(qc *ds.QuroumCertificate, qtype int, qview int) bool {
	return (qc.ViewNumber == qview) && (qc.Type == qtype)
}

func (r *Replica) SafeNode(node *ds.Node, qc *ds.QuroumCertificate) bool {
	return r.safetyRule(node) && r.livenessRule(qc.ViewNumber)
}
