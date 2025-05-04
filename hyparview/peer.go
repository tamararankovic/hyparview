package hyparview

import (
	"github.com/tamararankovic/hyparview/connections"
	"github.com/tamararankovic/hyparview/data"
)

type Peer struct {
	node data.Node
	conn connections.Conn
}
