package connections

import (
	"github.com/tamararankovic/hyparview/data"
)

type Conn interface {
	GetAddress() string
	Send(msg data.Message, sender data.Node) error
	onReceive(handler func(msg data.Message))
	disconnect() error
	onDisconnect(handler func())
}
