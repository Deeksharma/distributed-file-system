package main

import (
	"bytes"
	"distributed-file-system/p2p"
	"encoding/gob"
	"fmt"
	"io"
	"log"
	"sync"
	"time"
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

// Message will be sent over the wire
type Message struct {
	Payload any
}

// MessageStoreFile will be sent over the wire
type MessageStoreFile struct {
	Key  string
	Size int64
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

func (s *FileServer) Stop() {
	close(s.quitch)
}

func (s *FileServer) StoreData(key string, r io.Reader) error {
	// 1. Store this file to the disk
	// 2. broadcast this file to all known peers in the network
	// 3. we need to encrypt both the keys and the data

	buf := new(bytes.Buffer)
	tee := io.TeeReader(r, buf)
	size, err := s.store.Write(key, tee)
	if err != nil {
		return err
	}

	msgBuf := new(bytes.Buffer)
	msg := Message{
		Payload: MessageStoreFile{
			Key:  key,
			Size: size,
		},
	}
	if err := gob.NewEncoder(msgBuf).Encode(msg); err != nil {
		return err
	}

	// first we send the message and then stream the big file
	for _, p := range s.peers {
		log.Printf("StoreData function call for server %s, sending key - (localaddress: %s ---> remoteAddress %s) %v\n", s.Transport.ListenAddress(), p.LocalAddr().String(), p.RemoteAddr().String(), p)

		if err := p.Send(msgBuf.Bytes()); err != nil {
			return err
		}
	}

	time.Sleep(time.Second * 3)

	// let's say this is the file that needs to be streamed over the network
	//payload := []byte("this is a very large file")
	for _, p := range s.peers {
		log.Printf("StoreData function call for server %s, sending the file - (localaddress: %s ---> remoteAddress %s) %v\n", s.Transport.ListenAddress(), p.LocalAddr().String(), p.RemoteAddr().String(), p)

		//if err := p.Send(payload); err != nil {
		//	return err
		//}
		n, err := io.Copy(p, buf)
		if err != nil {
			return err
		}

		fmt.Printf("%d bytes written\n", n)
	}

	return nil

	//buf := new(bytes.Buffer)
	//
	//tee := io.TeeReader(r, buf)
	//
	//if err := s.store.Write(key, tee); err != nil {
	//	return err
	//}
	//
	//// after writing to the disk the reader is empty so we are using tee reader
	//
	//_, err := io.Copy(buf, r)
	//if err != nil {
	//	return err
	//}
	//
	//p := &DataMessage{Key: key, Data: buf.Bytes()}
	//log.Println(buf.Bytes()) // this is empty so only key is being sent over the network
	//return s.broadcast(&Message{
	//	From:    s.Transport.ListenAddress(),
	//	Payload: p,
	//})
}

// OnPeer this peer can work for other types of transport too, like UDP, HTTP, GRPC etc.
func (s *FileServer) OnPeer(p p2p.Peer) error {
	s.peerLock.Lock()
	defer s.peerLock.Unlock()

	s.peers[p.RemoteAddr().String()] = p
	log.Printf("On Peer function call for server %s - (localaddress: %s ---> remoteAddress %s) %v\n", s.Transport.ListenAddress(), p.LocalAddr().String(), p.RemoteAddr().String(), p)
	return nil
}

func (s *FileServer) bootstrapNetwork() error {
	for _, addr := range s.BootstrapNodes {
		if len(addr) == 0 {
			continue
		}
		log.Printf("Attempting to connect to %s\n", addr)
		go func(addr string) {
			err := s.Transport.Dial(addr)
			if err != nil {
				log.Println("Dial error:", err)
			}
		}(addr)
	}
	return nil
}

func (s *FileServer) handleMessage(from string, msg *Message) error {
	switch v := msg.Payload.(type) {
	case MessageStoreFile:
		return s.handleMessageStoreFile(from, v)
	}
	return nil
}

func (s *FileServer) handleMessageStoreFile(from string, msg MessageStoreFile) error {
	log.Printf("Inside handleMessageStoreFile Received message from %s: %+v\n", from, msg)
	peer, ok := s.peers[from]
	if !ok {
		return fmt.Errorf("peer (%s) could not be found int he peer list of server (%s)", from, s.Transport.ListenAddress())
	}

	if _, err := s.store.Write(msg.Key, io.LimitReader(peer, msg.Size)); err != nil {
		return err
	}

	peer.(*p2p.TCPPeer).Wg.Done()

	return nil
}

func (s *FileServer) loop() {
	defer func() {
		log.Println("file server stopped due to user shut down")
		s.Transport.Close()
	}()

	for {
		select {
		case rpc := <-s.Transport.Consume():
			var msg Message

			if err := gob.NewDecoder(bytes.NewReader(rpc.Payload)).Decode(&msg); err != nil {
				log.Println("failed to decode the message" + err.Error())
			}

			if err := s.handleMessage(rpc.From.String(), &msg); err != nil {
				log.Println("failed to handle the message" + err.Error())
				return
			}

			//log.Printf("Received key message %+v\n", msg.Payload)
			//
			//// we need to read from peer here
			//peer, ok := s.peers[rpc.From.String()]
			//log.Printf("Inside loop function call for server %s - (localaddress: %s ---> remoteAddress %s) %v\n", s.Transport.ListenAddress(), peer.LocalAddr().String(), peer.RemoteAddr().String(), peer)
			//
			//if !ok {
			//	panic("peer not found in peers map")
			//}
			//
			//b := make([]byte, 1000)
			//// read function of the underlying connection
			//// we are getting stuck while reading from the connection, why? why do we have to read the peer? this peer is the main gr port
			//// the problem is not that the thread is hanging, the problem is that 2 gr's are reading from the same conn, and there wuill be inconsistency which will read first
			//// for the first message it will be read by handleConn but the next message can be read by this loop
			////
			//log.Println("before the read ")
			//if _, err := peer.Read(b); err != nil {
			//	panic("failed to read from peer")
			//}
			////panic("this is a very large file")
			//log.Println("after the read ")
			//
			//log.Printf("Received stream data from peer: %s\n", string(b))
			//peer.(*p2p.TCPPeer).Wg.Done()
		case <-s.quitch:
			return

		}
	}
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
	//	//log.Println("Sending the bytes data to the peer", buf.Bytes())
	//	//peer.Send(buf.Bytes())
	//}
	//return nil
}

func init() {
	// any type that we are placing inside Message.Payload needs to be i
	gob.Register(MessageStoreFile{})
}
