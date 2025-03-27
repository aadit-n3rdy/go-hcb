package replica

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
