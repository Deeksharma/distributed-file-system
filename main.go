package main

import (
	"distributed-file-system/p2p"
	"fmt"
	"log"
)

// TIP <p>To run your code, right-click the code and select <b>Run</b>.</p> <p>Alternatively, click
// the <icon src="AllIcons.Actions.Execute"/> icon in the gutter and select the <b>Run</b> menu item from here.</p>
func OnPeer(peer p2p.Peer) error {
	peer.Close()
	fmt.Println("doing some logic with the peer outside of TCP Transport")
	return nil
	//return fmt.Errorf("error connecting to the peer %+v\n", peer)
}

func main() {
	tcpOpts := p2p.TCPTransportOps{
		ListenAddr:    ":3000",
		HandshakeFunc: p2p.NOPHandshakeFunc,
		Decoder:       &p2p.DefaultDecoder{},
		OnPeer:        OnPeer,
	}
	transport := p2p.NewTCPTransport(tcpOpts)
	//log.Fatal(transport.ListenAndAccept())

	go func() {
		for {
			msg := <-transport.Consume()
			fmt.Println(string(msg.Payload))
		}
	}()

	if err := transport.ListenAndAccept(); err != nil {
		log.Fatal(err)
	}
	select {}
}
