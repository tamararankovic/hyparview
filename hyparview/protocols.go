package hyparview

import (
	"fmt"
	"log"

	"github.com/tamararankovic/hyparview/connections"
	"github.com/tamararankovic/hyparview/data"
)

func (h *HyParView) onJoin(received connections.MsgReceived) error {
	msg, ok := received.Msg.Payload.(data.Join)
	if !ok {
		return fmt.Errorf("msg %v not a join msg", received.Msg.Payload)
	}
	if h.activeViewFull() {
		err := h.disconnectRandomPeer()
		if err != nil {
			return err
		}
	}
	newPeer := Peer{
		node: data.Node{
			ID:      msg.NodeID,
			Address: msg.NodeAddress,
		},
		conn: received.Sender,
	}
	h.activeView = append(h.activeView, newPeer)
	forwardJoinMsg := data.Message{
		Type: data.FORWARD_JOIN,
		Payload: data.ForwardJoin{
			NodeID:      msg.NodeID,
			NodeAddress: msg.NodeAddress,
			TTL:         h.config.ARWL,
		},
	}
	for _, peer := range h.activeView {
		if peer.node.ID == newPeer.node.ID {
			continue
		}
		err := peer.conn.Send(forwardJoinMsg, h.self)
		if err != nil {
			log.Println(err)
		}
	}
	return nil
}

func (h *HyParView) onDisconnect(received connections.MsgReceived) error {
	msg, ok := received.Msg.Payload.(data.Disconnect)
	if !ok {
		return fmt.Errorf("msg %v not a disconnect msg", received.Msg.Payload)
	}
	peer := h.getPeer(received.Sender)
	if peer == nil {
		log.Printf("peer %s not in active view\n", received.Sender.GetAddress())
	}
	if peer.node.ID != msg.NodeID {
		log.Printf("node trying to disconnect not matching the message info: %v - %s\n", msg, received.Sender.GetAddress())
		return nil
	}
	h.deletePeer(*peer)
	return h.connManager.Disconnect(peer.conn)
}

func (h *HyParView) onForwardJoin(received connections.MsgReceived) error {
	msg, ok := received.Msg.Payload.(data.ForwardJoin)
	if !ok {
		return fmt.Errorf("msg %v not a forward join msg", received.Msg.Payload)
	}
	newPeer := Peer{
		node: data.Node{
			ID:      msg.NodeID,
			Address: msg.NodeAddress,
		},
		conn: nil,
	}
	addedToActiveView := false
	if msg.TTL == 0 || len(h.activeView) == 1 {
		conn, err := h.connManager.Connect(msg.NodeAddress)
		if err != nil {
			return err
		}
		newPeer.conn = conn
		h.activeView = append(h.activeView, newPeer)
		addedToActiveView = true
	} else if msg.TTL == h.config.PRWL {
		h.passiveView = append(h.passiveView, newPeer)
	}
	msg.TTL--
	if !addedToActiveView {
		randomPeer := h.selectRandomPeer()
		for randomPeer.node.Address == received.Sender.GetAddress() {
			randomPeer = h.selectRandomPeer()
		}
		return randomPeer.conn.Send(data.Message{
			Type:    data.FORWARD_JOIN,
			Payload: msg,
		}, h.self)
	}
	return nil
}
