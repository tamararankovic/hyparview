package hyparview

import (
	"log"
	"math/rand"
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
	msgHandlers map[data.MessageType]func(received connections.MsgReceived) error
}

func NewHyParView(config HyParViewConfig, self data.Node, connManager connections.ConnManager) HyParView {
	hv := HyParView{
		self:        self,
		config:      config,
		activeView:  make([]Peer, 0),
		passiveView: make([]Peer, 0),
		peerUp:      make(chan Peer),
		peerDown:    make(chan Peer),
		connManager: connManager,
	}
	hv.msgHandlers = map[data.MessageType]func(received connections.MsgReceived) error{
		data.JOIN:         hv.onJoin,
		data.DISCONNECT:   hv.onDisconnect,
		data.FORWARD_JOIN: hv.onForwardJoin,
	}
	return hv
}

func (h *HyParView) Join(contactNodeAddress string) error {
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

func (h *HyParView) GetPeers() []Peer {
	return h.activeView
}

func (h *HyParView) OnPeerUp(handler func(peer Peer)) connections.Subscription {
	return h.connManager.OnConnUp(func(conn connections.Conn) {
		if peer := h.getPeer(conn); peer != nil {
			handler(*peer)
		}
	})
}

func (h *HyParView) OnPeerDown(handler func(peer Peer)) connections.Subscription {
	return h.connManager.OnConnDown(func(conn connections.Conn) {
		if peer := h.getPeer(conn); peer != nil {
			handler(*peer)
		}
	})
}

func (h *HyParView) onReeive(received connections.MsgReceived) {
	handler := h.msgHandlers[received.Msg.Type]
	if handler == nil {
		log.Printf("no handler found for message type %s", received.Msg.Type)
		return
	}
	err := handler(received)
	if err != nil {
		log.Println(err)
	}
}

func (h *HyParView) disconnectRandomPeer() error {
	disconnectPeer := h.selectRandomPeer()
	if disconnectPeer == nil {
		return nil
	}
	h.deletePeer(*disconnectPeer)
	disconnectMsg := data.Message{
		Type: data.DISCONNECT,
		Payload: data.Disconnect{
			NodeID:      h.self.ID,
			NodeAddress: h.self.Address,
		},
	}
	err := disconnectPeer.conn.Send(disconnectMsg, h.self)
	if err != nil {
		return err
	}
	err = h.connManager.Disconnect(disconnectPeer.conn)
	if err != nil {
		log.Println(err)
	}
	return nil
}

func (h *HyParView) getPeer(conn connections.Conn) *Peer {
	index := slices.IndexFunc(h.activeView, func(peer Peer) bool {
		return peer.conn == conn
	})
	if index < 0 {
		return nil
	}
	return &h.activeView[index]
}

func (h *HyParView) deletePeer(peer Peer) {
	index := slices.IndexFunc(h.activeView, func(p Peer) bool {
		return p.node.ID == peer.node.ID
	})
	if index < 0 {
		return
	}
	h.activeView = slices.Delete(h.activeView, index, index+1)
}

func (h *HyParView) activeViewSize() int {
	return h.config.Fanout + 1
}

func (h *HyParView) activeViewFull() bool {
	return len(h.activeView) == h.activeViewSize()
}

func (h *HyParView) selectRandomPeer() *Peer {
	if len(h.activeView) == 0 {
		return nil
	}
	index := rand.Intn(len(h.activeView))
	return &h.activeView[index]
}
