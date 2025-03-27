package replica

import (
	"net/rpc"

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

	isLeader   bool
	leaderID   int
	successors []int

	viewNumber int
	lockedQC   ds.QuroumCertificate
	prepareQC  ds.QuroumCertificate
	branch     *ds.Branch
}
