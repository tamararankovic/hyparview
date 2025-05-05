package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/tamararankovic/hyparview/connections"
	"github.com/tamararankovic/hyparview/data"
	"github.com/tamararankovic/hyparview/hyparview"
)

func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	id := os.Args[1]
	address := os.Args[2]
	contactNodeAddress := os.Args[3]
	config := hyparview.Config{
		NodeID:             id,
		NodeAddress:        address,
		ContactNodeAddress: contactNodeAddress,
		HyParViewConfig: hyparview.HyParViewConfig{
			Fanout:          2,
			PartialViewSize: 10,
			ARWL:            3,
			PRWL:            2,
			ShuffleInterval: 10,
			Ka:              1,
			Kp:              1,
		},
	}
	self := data.Node{
		ID:      config.NodeID,
		Address: config.NodeAddress,
	}
	connManager := connections.NewConnManager(connections.NewTCPConn, connections.AcceptTcpConnsFn(self.Address))
	hv, err := hyparview.NewHyParView(config.HyParViewConfig, self, connManager)
	if err != nil {
		log.Fatal(err)
	}
	// time.Sleep(10 * time.Second)
	err = hv.Join(contactNodeAddress)
	if err != nil {
		log.Println(err)
	}

	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
	log.Println("Waiting for exit signal...")
	sig := <-sigs
	log.Println("Received signal:", sig)
}
