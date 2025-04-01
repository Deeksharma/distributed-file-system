package main

import (
	"bytes"
	"distributed-file-system/p2p"
	"fmt"
	"log"
	"time"
)

// TIP <p>To run your code, right-click the code and select <b>Run</b>.</p> <p>Alternatively, click
// the <icon src="AllIcons.Actions.Execute"/> icon in the gutter and select the <b>Run</b> menu item from here.</p>
func OnPeer(peer p2p.Peer) error {
	peer.Close()
	fmt.Println("doing some logic with the peer outside of TCP Transport")
	return nil
	//return fmt.Errorf("error connecting to the peer %+v\n", peer)
}

func makeServer(listenAddr string, nodes ...string) *FileServer {
	transportOpts := p2p.TCPTransportOps{
		ListenAddr:    ":" + listenAddr,
		HandshakeFunc: p2p.NOPHandshakeFunc,
		Decoder:       &p2p.DefaultDecoder{},
	}

	transport := p2p.NewTCPTransport(transportOpts)

	fileServerOptions := FileServerOpts{
		StorageRoot:       listenAddr + "_network",
		PathTransformFunc: CASPathTransformFunc,
		Transport:         transport,
		BootstrapNodes:    nodes,
	}

	s := NewFileServer(fileServerOptions)

	transport.OnPeer = s.OnPeer
	return s
}

func main() {
	//tcpOpts := p2p.TCPTransportOps{
	//	ListenAddr:    ":3000",
	//	HandshakeFunc: p2p.NOPHandshakeFunc,
	//	Decoder:       &p2p.DefaultDecoder{},
	//	OnPeer:        OnPeer,
	//}
	//transport := p2p.NewTCPTransport(tcpOpts)
	////log.Fatal(transport.ListenAndAccept())
	//
	//go func() {
	//	for {
	//		msg := <-transport.Consume()
	//		fmt.Println(string(msg.DataMessage))
	//	}
	//}()
	//
	//if err := transport.ListenAndAccept(); err != nil {
	//	log.Fatal(err)
	//}
	//select {}

	server1 := makeServer("3000")

	server2 := makeServer("4000", ":3000")

	go func() {
		log.Fatal(server1.Start())
	}()

	go server2.Start()

	time.Sleep(5 * time.Second)
	//go func() {
	//	time.Sleep(5 * time.Second)
	//	s.Stop()
	//}()
	//
	//if err := s.Start(); err != nil {
	//	log.Fatal(err)
	//}

	data := bytes.NewReader([]byte("my very big data"))
	server2.StoreData("myPrivateData", data)
	select {}
}
