package network

import (
	"bytes"
	"crypto/elliptic"
	"encoding/binary"
	"encoding/gob"
	"fmt"
	"io"
	"net"

	"github.com/anthdm/projectx/core"
	"github.com/sirupsen/logrus"
)

type MessageType byte

const (
	MessageTypeTx        MessageType = 0x1
	MessageTypeBlock     MessageType = 0x2
	MessageTypeGetBlocks MessageType = 0x3
	MessageTypeStatus    MessageType = 0x4
	MessageTypeGetStatus MessageType = 0x5
	MessageTypeBlocks    MessageType = 0x6
)

type RPC struct {
	From    net.Addr //string
	Payload io.Reader
}

type Message struct {
	Header MessageType
	Data   []byte
}

func NewMessage(t MessageType, data []byte) *Message {
	return &Message{
		Header: t,
		Data:   data,
	}
}

func NewMessageGobEncoded(t MessageType, data any) (*Message, error) {
	buf := new(bytes.Buffer)
	if err := gob.NewEncoder(buf).Encode(data); err != nil {
		return nil, err
	}

	return &Message{
		Header: t,
		Data:   buf.Bytes(),
	}, nil
}

func NewMessageWithEncoder[T any](
	t MessageType,
	encFactory func(writer io.Writer) core.Encoder[T],
	data T,
) (*Message, error) {
	buf := new(bytes.Buffer)
	if err := encFactory(buf).Encode(data); err != nil {
		return nil, err
	}

	return &Message{
		Header: t,
		Data:   buf.Bytes(),
	}, nil
}

func (msg *Message) WriteTo(w io.Writer) (int64, error) {
	n, err := w.Write([]byte{byte(msg.Header)})
	if err != nil {
		return 0, err
	}
	if err = binary.Write(w, binary.LittleEndian, len(msg.Data)); err != nil {
		return 0, err
	}
	if n, err = w.Write(msg.Data); err != nil {
		return 0, err
	}

	return int64(1 + binary.Size(len(msg.Data)) + n), err
}

type DecodedMessage struct {
	From net.Addr
	Data any
}

type RPCDecodeFunc func(RPC) (*DecodedMessage, error)

func DefaultRPCDecodeFunc(rpc RPC) (*DecodedMessage, error) {
	buf := make([]byte, 1)
	if _, err := rpc.Payload.Read(buf); err != nil {
		return nil, err
	}
	var dataLen int
	if err := binary.Read(rpc.Payload, binary.LittleEndian, &dataLen); err != nil {
		return nil, err
	}
	msg := Message{
		Header: MessageType(buf[0]),
		Data:   make([]byte, dataLen),
	}
	if _, err := rpc.Payload.Read(msg.Data); err != nil {
		return nil, err
	}
	fmt.Printf("receiving message: %+v\n", msg)

	logrus.WithFields(logrus.Fields{
		"from": rpc.From,
		"type": msg.Header,
	}).Debug("new incoming message")

	switch msg.Header {
	case MessageTypeTx:
		tx := new(core.Transaction)
		if err := tx.Decode(core.NewGobTxDecoder(bytes.NewReader(msg.Data))); err != nil {
			return nil, err
		}

		return &DecodedMessage{
			From: rpc.From,
			Data: tx,
		}, nil
	case MessageTypeBlock:
		block := new(core.Block)
		if err := block.Decode(core.NewGobBlockDecoder(bytes.NewReader(msg.Data))); err != nil {
			return nil, err
		}

		return &DecodedMessage{
			From: rpc.From,
			Data: block,
		}, nil

	case MessageTypeGetStatus:
		return &DecodedMessage{
			From: rpc.From,
			Data: &GetStatusMessage{},
		}, nil

	case MessageTypeStatus:
		statusMessage := new(StatusMessage)
		if err := gob.NewDecoder(bytes.NewReader(msg.Data)).Decode(statusMessage); err != nil {
			return nil, err
		}

		return &DecodedMessage{
			From: rpc.From,
			Data: statusMessage,
		}, nil
	case MessageTypeGetBlocks:
		getBlocks := new(GetBlocksMessage)
		if err := gob.NewDecoder(bytes.NewReader(msg.Data)).Decode(getBlocks); err != nil {
			return nil, err
		}

		return &DecodedMessage{
			From: rpc.From,
			Data: getBlocks,
		}, nil

	case MessageTypeBlocks:
		blocks := new(BlocksMessage)
		if err := gob.NewDecoder(bytes.NewReader(msg.Data)).Decode(blocks); err != nil {
			return nil, err
		}

		return &DecodedMessage{
			From: rpc.From,
			Data: blocks,
		}, nil
	default:
		return nil, fmt.Errorf("invalid message header %x", msg.Header)
	}
}

type RPCProcessor interface {
	ProcessMessage(*DecodedMessage) error
}

func init() {
	gob.Register(elliptic.P256())
}
