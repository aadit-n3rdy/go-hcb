package replica

import (
	"github.com/aadit-n3rdy/hotCrossBuns/ds"
)

type Replica struct {
	id         int
	viewNumber int
	lockedQC   ds.QuroumCertificate
	prepareQC  ds.QuroumCertificate
	branch     *ds.Branch
}

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
