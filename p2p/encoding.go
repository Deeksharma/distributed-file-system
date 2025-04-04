package p2p

import (
	"encoding/gob"
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
	buf := make([]byte, 1024) // this is very small for large files so we will stream data
	n, err := reader.Read(buf)

	if err != nil {
		return err
	}

	msg.Payload = buf[:n]
	return nil
}
