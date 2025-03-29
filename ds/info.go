package ds

// struct that informs participant replica of leader replica id
type LeaderBroadcast struct {
	LId int
}

// struct that informs current leader of a potential successor
type LeaderAttempt struct {
	Id int
}

// struct to inform next leader of its successors
type FutureGenerations struct {
	Successors []int
}
