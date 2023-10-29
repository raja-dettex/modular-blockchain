package network

import (
	"bytes"
	"fmt"
	"io"
	"net"
)

type TCPPeer struct {
	conn net.Conn
}

func (peer *TCPPeer) sendMsg(b []byte) error {
	_, err := peer.conn.Write(b)
	if err == io.EOF {
		return nil
	}
	if err != nil {
		return err
	}
	return nil
}

func (peer *TCPPeer) HandleConn(rpcCh chan RPC) {
	buff := make([]byte, 2048)
	for {
		n, err := peer.conn.Read(buff)
		if err != nil {
			fmt.Printf("read from byte stream error : %+v\n", err)
		}
		rpcCh <- RPC{
			From:    peer.conn.RemoteAddr(),
			Payload: bytes.NewReader(buff[:n]),
		}
		// go func() {
		// 	rpcCh <- RPC{
		// 		From:    peer.conn.RemoteAddr(),
		// 		Payload: bytes.NewReader(buff[:n]),
		// 	}
		// }()
	}
}

type TCPTransport struct {
	peerCh     chan *TCPPeer
	listenaddr string
	ln         net.Listener
}

func NewTCPTransport(addr string, peerCh chan *TCPPeer) *TCPTransport {
	return &TCPTransport{
		peerCh:     peerCh,
		listenaddr: addr,
	}
}

func (t *TCPTransport) Start() error {
	ln, err := net.Listen("tcp", t.listenaddr)
	if err != nil {
		return err
	}
	t.ln = ln
	fmt.Printf("listener starting on addr %v\n", t.listenaddr)
	go t.AcceptLoop()
	return nil
}

// func (t *TCPTransport) handleConn(peer *TCPPeer) {
// 	buff := make([]byte, 2048)
// 	for {
// 		n, err := peer.conn.Read(buff)
// 		if err != nil {
// 			fmt.Printf("Read from byte stream error %v\n ", err)
// 		}
// 		msg := buff[:n]
// 		fmt.Println(string(msg))
// 	}
// }

func (t *TCPTransport) AcceptLoop() {
	for {
		conn, err := t.ln.Accept()
		if err != nil {
			fmt.Printf("listen to conn %v error : %v\n", conn, err)
		}
		fmt.Printf("new connction -> %+v\n", conn)
		peer := &TCPPeer{
			conn: conn,
		}

		t.peerCh <- peer
		// time.Sleep(time.Second * 2)
		// getPeer := <-t.peerCh
		// fmt.Println(getPeer)

	}
}
