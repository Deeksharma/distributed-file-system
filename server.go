package main

import (
	"bytes"
	"distributed-file-system/p2p"
	"encoding/gob"
	"fmt"
	"io"
	"log"
	"sync"
)

type FileServerOpts struct {
	StorageRoot       string
	PathTransformFunc ParthTransformFunc
	Transport         p2p.Transport
	BootstrapNodes    []string
}

// FileServer will be handling all the command line or internal requests
type FileServer struct {
	FileServerOpts

	peerLock sync.Mutex
	peers    map[string]p2p.Peer

	store  *Store
	quitch chan struct{}
}

func NewFileServer(options FileServerOpts) *FileServer {
	return &FileServer{
		FileServerOpts: options,
		store: NewStore(StoreOpts{
			Root:              options.StorageRoot,
			PathTransformFunc: options.PathTransformFunc,
		}),
		quitch: make(chan struct{}),
		peers:  make(map[string]p2p.Peer),
	}
}

// Message will be sent over the wire
type Message struct {
	From    string
	Payload any
}
type DataMessage struct {
	Key  string
	Data []byte
}

func (s *FileServer) broadcast(msg *Message) error {
	//return gob.NewEncoder(p.conn).Encode(p)

	peers := []io.Writer{} // Peer interface embeds net.Con interface which embeds io.write as well

	for _, peer := range s.peers {
		peers = append(peers, peer)
	}

	mw := io.MultiWriter(peers...)
	return gob.NewEncoder(mw).Encode(msg)

	////buf := new(bytes.Buffer)
	//for _, peer := range s.peers {
	//	if err := gob.NewEncoder(peer).Encode(p); err != nil {
	//		return err
	//	}
	//	//if err := gob.NewEncoder(buf).Encode(p); err != nil {
	//	//	return err
	//	//}
	//	//
	//	//fmt.Println("Sending the bytes data to the peer", buf.Bytes())
	//	//peer.Send(buf.Bytes())
	//}
	//return nil
}

func (s *FileServer) StoreData(key string, r io.Reader) error {
	// 1. Store this file to the disk
	// 2. broadcast this file to all known peers in the network
	// 3. we need to encrypt both the keys and the data

	//p := &DataMessage{Key: key}

	buf := new(bytes.Buffer)

	tee := io.TeeReader(r, buf)

	if err := s.store.Write(key, tee); err != nil {
		return err
	}

	// after writing to the disk the reader is empty so we are using tee reader

	_, err := io.Copy(buf, r)
	if err != nil {
		return err
	}

	p := &DataMessage{Key: key, Data: buf.Bytes()}
	fmt.Println(buf.Bytes()) // this is empty so only key is being sent over the network
	return s.broadcast(&Message{
		From:    s.Transport.ListenAddress(),
		Payload: p,
	})
}

// this peer can work for other types of transport too, like UDP, HTTP, GRPC etc.
func (s *FileServer) OnPeer(p p2p.Peer) error {
	s.peerLock.Lock()
	defer s.peerLock.Unlock()

	s.peers[p.RemoteAddr().String()] = p
	log.Printf("connected with remote peer %s", p.RemoteAddr().String())
	return nil
}

func (s *FileServer) Stop() {
	close(s.quitch)
}

func (s *FileServer) handleMessage(msg *Message) error {
	switch v := msg.Payload.(type) {
	case *DataMessage:
		fmt.Printf("Received DataMessage from %+v\n", v)
	}
	return nil
}

func (s *FileServer) loop() {
	defer func() {
		log.Println("file server stopped due to user shut down")
		s.Transport.Close()
	}()

	for {
		select {
		case msg := <-s.Transport.Consume():
			var message Message
			fmt.Printf("Received message: %s\n", msg)

			if err := gob.NewDecoder(bytes.NewReader(msg.Payload)).Decode(&message); err != nil {
				log.Println("failed to decode the message" + err.Error())
			}

			if err := s.handleMessage(&message); err != nil {
				log.Println(err)
			}
			//fmt.Printf("Received message key: %s\n", message.From)
			//fmt.Printf("Received message data: %s\n", message.Payload)
		case <-s.quitch:
			return

		}
	}
}

func (s *FileServer) bootstrapNetwork() error {
	for _, addr := range s.BootstrapNodes {
		if len(addr) == 0 {
			continue
		}
		fmt.Println("Attempting to connect to ", addr)
		go func(addr string) {
			err := s.Transport.Dial(addr)
			if err != nil {
				log.Println("Dial error:", err)
			}
		}(addr)
	}
	return nil
}

func (s *FileServer) Start() error {

	if err := s.Transport.ListenAndAccept(); err != nil {
		return err
	}

	if len(s.BootstrapNodes) != 0 {
		s.bootstrapNetwork()
	}
	s.loop()

	return nil
}
