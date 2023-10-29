package network

import (
	"bytes"
	"fmt"
	"io"
	"net"
	"testing"
	"time"
)

type Stream struct {
	id   int64
	data io.Reader
}

func createStream() []byte {
	buff := make([]byte, 10)
	for i := 0; i < len(buff); i++ {
		buff = append(buff, 0x02)
	}
	return buff
}

func decode(msg *Stream) {
	buff := make([]byte, 10)
	for {
		n, err := msg.data.Read(buff)
		if err != nil {
			fmt.Printf("error %+v\n", err)
		}
		fmt.Printf("messge received with id %v and the bytes are %v\n", msg.id, buff[:n])
	}
}

func TestIOReadStream(t *testing.T) {
	buff := createStream()
	panic("before")
	fmt.Println(buff)
	panic("after")
	msg := &Stream{
		id:   int64(1),
		data: bytes.NewReader(buff),
	}
	decode(msg)
}

func TestMsgDeliveredToPeer(t *testing.T) {
	//var blockingCh chan interface{}
	peerCh := make(chan *TCPPeer, 1)
	rpcCh := make(chan RPC)
	tr := NewTCPTransport(":3000", peerCh)
	go tr.Start()
	// assert.Nil(t, err)
	time.Sleep(time.Second * 2)
	fmt.Println("sending connection")
	go testConn()
	peer := <-peerCh
	time.Sleep(time.Second * 2)

	fmt.Printf("new peer %v\n", peer)
	go peer.HandleConn(rpcCh)
	time.Sleep(time.Second * 2)
	go handleMsgChannel(rpcCh)
	time.Sleep(time.Second * 2)
}

func sendMultipleConn() {
	for {
		go testConn()
		time.Sleep(time.Second * 2)
	}
}

func handleMsgChannel(rpcCh chan RPC) {
	for {
		select {
		case msg := <-rpcCh:
			for {
				buff := bytes.Buffer{}
				_, err := msg.Payload.Read(buff.Bytes())
				if err == io.EOF {
					fmt.Printf("message receieved from %v and the message is %v\n", msg.From, buff.Bytes())
				}
				if err != nil {
					fmt.Printf("error reading from rpc stream %v \n", err)
				}
			}
		}
	}
}

func TestTcpPeerChannel(t *testing.T) {
	peerCh := make(chan *TCPPeer, 1)
	tr := NewTCPTransport(":3000", peerCh)
	go tr.Start()
	// assert.Nil(t, err)
	time.Sleep(time.Second * 2)
	fmt.Println("sending connection")
	go testConn()
	peer := <-peerCh
	fmt.Printf("new peer %v", peer)
}
func testConn() {
	conn, err := net.Dial("tcp", ":3000")
	if err != nil {
		panic(err)
	}
	_, err = conn.Write([]byte("hey you!!!"))
	if err != nil {
		panic(err)
	}
}
