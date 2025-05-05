package data

type MessageType int

const (
	JOIN MessageType = iota
	FORWARD_JOIN
	DISCONNECT
	NEIGHTBOR
	NEIGHTBOR_REPLY
	SHUFFLE
	SHUFFLE_REPLY
)

type Message struct {
	Type    MessageType
	Payload any
}

type Join struct {
	NodeID,
	NodeAddress string
}

type ForwardJoin struct {
	NodeID,
	NodeAddress string
	TTL int
}

type Disconnect struct {
	NodeID,
	NodeAddress string
}

type Neighbor struct {
}

type NeighborReply struct {
}

type Shuffle struct {
}

type ShuffleReply struct {
}
