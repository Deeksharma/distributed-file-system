package main

import (
	"bytes"
	"distributed-file-system/p2p"
	"log"
	"time"
)

func OnPeer(peer p2p.Peer) error {
	peer.Close()
	log.Println("doing some logic with the peer outside of TCP Transport")
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

	transport.OnPeer = s.OnPeer // you can assign functions directly here - i did not know this
	return s
}

func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	server1 := makeServer("3000")

	server2 := makeServer("4000", ":3000")

	go func() {
		log.Fatal(server1.Start())
	}()

	go server2.Start()

	time.Sleep(5 * time.Second)

	data := bytes.NewReader([]byte("my very big data"))
	server2.StoreData("myPrivateData", data)
	select {}
}
