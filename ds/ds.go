package ds

const (
	NEW_VIEW int = iota
	PREPARE
	PRE_COMMIT
	COMMIT
	DECIDE
)

type Node struct {
	ID     int
	Cmd    string
	Parent *Node
}

type Branch struct {
	Root *Node
	Tail *Node
}

type QuroumCertificate struct {
	Type       int
	ViewNumber int
	Node       Node
	Sig        []int
}

type Message struct {
	Type       int
	CurView    int
	Node       Node
	Justify    QuroumCertificate
	PartialSig []int
}
