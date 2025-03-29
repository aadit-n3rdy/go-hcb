package replica

import (
	"github.com/aadit-n3rdy/hotCrossBuns/ds"
)

// function to check if passed node extends from the lockedQC
func (r *Replica) safetyRule(node *ds.Node, qc *ds.QuroumCertificate) bool {

	flag := false

	// check if node extends from locked qc
	itr := r.lockedQC.Node

	for itr.ID < node.ID {
		if itr.ID == node.ID {
			flag = flag || r.compareStringArrays(itr.Cmd, node.Cmd)
			break
		}
		itr = *itr.Parent
	}

	return flag
}

// helper function to compare array of strings
func (r *Replica) compareStringArrays(a []string, b []string) bool {
	if len(a) == len(b) {
		for i := range len(a) {
			if a[i] != b[i] {
				return false
			}
		}

		return true
	}

	return false
}

// function to check if liveness rule is satisfied
func (r *Replica) livenessRule(qview int) bool {
	return qview > r.lockedQC.ViewNumber
}

// function to create a message
func (r *Replica) Msg(mtype int, node ds.Node, qc *ds.QuroumCertificate) *ds.Message {
	msg := new(ds.Message)

	msg.Type = mtype
	msg.CurView = r.viewNumber
	msg.Node = node
	msg.Justify = qc

	return msg
}

// function to create a message and sign it
func (r *Replica) VoteMsg(mtype int, node ds.Node, qc *ds.QuroumCertificate) *ds.Message {
	msg := r.Msg(mtype, node, qc)

	msg.PartialSig[r.id] = 1

	return msg
}

// function to create a leaf for the current branch
func (r *Replica) CreateLeaf(parent *ds.Node, cmd []string) *ds.Node {
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
	return r.safetyRule(node, qc) && r.livenessRule(qc.ViewNumber)
}
