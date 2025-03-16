package p2p

import (
	"encoding/gob"
	"fmt"
	"io"
)

type Decoder interface {
	Decode(io.Reader, *RPC) error
}

type GOBDecoder struct {
}

func (dec *GOBDecoder) Decode(reader io.Reader, msg *RPC) error {
	return gob.NewDecoder(reader).Decode(msg)
}

type DefaultDecoder struct {
}

func (dec *DefaultDecoder) Decode(reader io.Reader, msg *RPC) error {
	buf := make([]byte, 2000)
	n, err := reader.Read(buf)
	if err != nil {
		return err
	}
	//buf = buf[:n]
	fmt.Println(string(buf[:n]))

	msg.Payload = buf[:n]
	return nil
}
