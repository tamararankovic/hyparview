package hyparview

import (
	"slices"

	"github.com/tamararankovic/hyparview/connections"
	"github.com/tamararankovic/hyparview/data"
)

type HyParView struct {
	self        data.Node
	config      HyParViewConfig
	activeView  []Peer
	passiveView []Peer
	connManager connections.ConnManager
	peerUp      chan Peer
	peerDown    chan Peer
}

func NewHyParView(config HyParViewConfig, connManager connections.ConnManager) HyParView {
	return HyParView{
		self: data.Node{
			ID:      config.NodeID,
			Address: config.NodeAddress,
		},
		config:      config,
		activeView:  make([]Peer, 0),
		passiveView: make([]Peer, 0),
		peerUp:      make(chan Peer),
		peerDown:    make(chan Peer),
		connManager: connManager,
	}
}

func (h HyParView) Join(contactNodeAddress string) error {
	conn, err := h.connManager.Connect(contactNodeAddress)
	if err != nil {
		return err
	}
	err = h.connManager.StartAcceptingConns()
	if err != nil {
		return err
	}
	sub := h.connManager.OnReceive(h.onReeive)
	msg := data.Message{
		Type: data.JOIN,
		Payload: data.Join{
			NodeID:      h.self.ID,
			NodeAddress: h.self.Address,
		},
	}
	err = conn.Send(msg, h.self)
	if err != nil {
		sub.Unsubscribe()
		h.connManager.StopAcceptingConns()
		return err
	}
	return nil
}

func (h HyParView) GetPeers() []Peer {
	return h.activeView
}

func (h HyParView) OnPeerUp(handler func(peer Peer)) connections.Subscription {
	return h.connManager.OnConnUp(func(conn connections.Conn) {
		if peer := h.getPeer(conn); peer != nil {
			handler(*peer)
		}
	})
}

func (h HyParView) OnPeerDown(handler func(peer Peer)) connections.Subscription {
	return h.connManager.OnConnDown(func(conn connections.Conn) {
		if peer := h.getPeer(conn); peer != nil {
			handler(*peer)
		}
	})
}

func (h HyParView) onReeive(msg connections.MsgReceived) {
	// todo: handle msg by type
}

func (h HyParView) getPeer(conn connections.Conn) *Peer {
	index := slices.IndexFunc(h.activeView, func(peer Peer) bool {
		return peer.conn == conn
	})
	if index < 0 {
		return nil
	}
	return &h.activeView[index]
}
