package replica

import (
	"net/rpc"
	"sync"

	"github.com/aadit-n3rdy/hotCrossBuns/ds"
)

type OtherReplica struct {
	Id     int
	Client *rpc.Client
	IpAddr string
	Port   string
}

type Replica struct {
	id     int
	ipAddr string
	port   string

	replicaList []*OtherReplica

	busy bool
	lock sync.Mutex

	isLeader   bool
	leaderID   int
	successors []int

	cmds       []string
	viewNumber int
	lockedQC   ds.QuroumCertificate
	prepareQC  ds.QuroumCertificate
	branch     *ds.Branch
}
