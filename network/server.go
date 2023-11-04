package network

import (
	"bytes"
	"encoding/gob"
	"fmt"
	"net"
	"os"
	"sync"
	"time"

	"github.com/go-kit/log"
	"github.com/raja-dettex/modular-blockchain/api"
	"github.com/raja-dettex/modular-blockchain/core"
	"github.com/raja-dettex/modular-blockchain/crypto"
	"github.com/raja-dettex/modular-blockchain/types"
)

var (
	defaultBlockTime = time.Second * 5
)

type ServerOpts struct {
	ApiListenAddr string
	SeedNodes     []string
	ListenAddr    string
	TcpTransport  *TCPTransport
	ID            string
	Logger        log.Logger
	RPCDecodeFunc RPCDecodeFunc
	RPCProcessor  RPCProcessor
	PrivateKey    *crypto.PrivateKey
	BlockTime     time.Duration
}
type Server struct {
	peerCh chan *TCPPeer

	mu           sync.RWMutex
	peerMap      map[string]*TCPPeer
	AccountState *core.AccountState

	Opts        ServerOpts
	isValidator bool
	chain       *core.Blockchain
	memPool     *TxPool
	rpcCh       chan RPC
	quitCh      chan struct{}
	txChan      chan *core.Transaction
}

func (s *Server) bootstrapNetwork() {
	for _, addr := range s.Opts.SeedNodes {
		if addr == s.Opts.ListenAddr {
			continue
		}
		go func(addr string) {
			conn, err := net.Dial("tcp", addr)
			if err != nil {
				s.Opts.Logger.Log("msg", "Can not connect to peer with address", " address", addr, "error", err)
			}
			s.peerCh <- &TCPPeer{
				conn: conn,
			}
		}(addr)

	}
}

func NewServer(opts ServerOpts) (*Server, error) {
	if opts.BlockTime == time.Duration(0) {
		opts.BlockTime = defaultBlockTime
	}
	if opts.RPCDecodeFunc == nil {
		opts.RPCDecodeFunc = DefaultRPCDecodeFunc
	}

	peerCh := make(chan *TCPPeer)
	tcpTransport := NewTCPTransport(opts.ListenAddr, peerCh)
	opts.TcpTransport = tcpTransport
	s := &Server{
		peerCh:      peerCh,
		peerMap:     make(map[string]*TCPPeer),
		isValidator: opts.PrivateKey != nil,
		memPool:     NewTxPool(100),
		rpcCh:       make(chan RPC),
		quitCh:      make(chan struct{}),
	}

	if opts.RPCProcessor == nil {
		opts.RPCProcessor = s
	}
	if opts.Logger == nil {
		opts.Logger = log.NewLogfmtLogger(os.Stderr)
		opts.Logger = log.With(opts.Logger, "ID", opts.ID)

	}
	s.Opts = opts
	ac := core.NewAccountState()
	if s.Opts.PrivateKey != nil {
		if err := ac.AddBalance(s.Opts.PrivateKey.GeneratePublicKey().Address(), 10000); err != nil {
			s.Opts.Logger.Log("err", err)
		}
	}
	s.AccountState = ac

	bc, err := core.NewBlockChain(s.Opts.Logger, genesisBlock(), ac)
	if err != nil {
		return nil, err
	}
	s.chain = bc

	// put the account state into the server

	if s.isValidator {
		s.Opts.Logger.Log("is validator", s.isValidator, "key", s.Opts.PrivateKey)
		go s.validatorLoop()
	}
	// for _, tr := range s.Opts.Transports {
	// 	if err := s.sendStatusMessage(tr); err != nil {
	// 		s.Opts.Logger.Log("send status to peer error", err)
	// 	}
	// }

	// send a get status message to its neighbour
	// s.sendStatusMessage()

	// start the api server if the config has a valid port

	// this channel is used to get the transactions from rpc to network and propagate it among nodes
	txChan := make(chan *core.Transaction)
	s.txChan = txChan
	if opts.ApiListenAddr != "" {
		apiServerConfig := api.ServerConfig{ListenAddr: opts.ApiListenAddr}
		apiServer := api.NewServer(apiServerConfig, s.chain, txChan)
		go apiServer.Start()
	}

	return s, nil
}

func (s *Server) Start() {
	s.Opts.Logger.Log("msg", "node started")
	s.Opts.TcpTransport.Start()
	s.bootstrapNetwork()
	// s.InitTransport()
free:
	for {
		select {
		case peer := <-s.peerCh:
			if _, ok := s.peerMap[peer.conn.RemoteAddr().String()]; !ok {
				s.peerMap[peer.conn.RemoteAddr().String()] = peer
				s.sendGetStatusMessage(peer)
			}

			go peer.HandleConn(s.rpcCh)
		case tx := <-s.txChan:
			if err := s.processTransaction(tx); err != nil {
				panic(err)
			}

		case rpc := <-s.rpcCh:
			msg, err := s.Opts.RPCDecodeFunc(rpc)
			if err != nil {
				s.Opts.Logger.Log("error", err)
				continue

				// if err != core.ErrBlockKnown {
				// 	s.Opts.Logger.Log("error", err)
				// }
			}
			if err = s.Opts.RPCProcessor.ProcessMessage(msg); err != nil {
				fmt.Printf("====> error %v\n", err)
				if err != core.ErrBlockKnown {
					s.Opts.Logger.Log("error", err)
				}
			}
		case <-s.quitCh:
			break free

		}
	}
	s.Opts.Logger.Log("msg", "server is shutting down")
}

func (s *Server) ProcessMessage(msg *DecodedMessage) error {

	switch t := msg.Data.(type) {
	case *core.Transaction:
		s.processTransaction(t)
	case *core.Block:
		if err := s.processBlock(t); err != nil {
			return err
		}
	case *GetStatusMessage:
		s.processGetStatusMessage(msg.From, t)
	case *StatusMessage:
		//fmt.Printf("received statusMessage %v", t)
		s.processStatusMessage(msg.From, t)
	case *GetBlocksMessage:
		s.Opts.Logger.Log("msg", "received blocks message")
		s.processGetBlocksMessage(msg.From, t)
	case *BlocksMessage:
		s.processBlocksMessage(msg.From, t)
	default:
		panic("here")
	}
	return nil
}

func (s *Server) processBlocksMessage(from net.Addr, data *BlocksMessage) error {
	s.Opts.Logger.Log("msg", "received blocks", "from", from, "blocks", data.Blocks)
	for _, block := range data.Blocks {
		if err := s.chain.AddBlock(block); err != nil {
			return err
		}
	}
	return nil
}

func (s *Server) processGetBlocksMessage(addr net.Addr, data *GetBlocksMessage) error {
	fmt.Printf("id : %v, received get blocks message +%v\n", s.Opts.ID, data)
	blocks := []*core.Block{}
	if data.From == 0 {
		for i := uint32(1); i <= data.To; i++ {
			b, err := s.chain.GetBlock(i)
			if err != nil {
				s.Opts.Logger.Log("msg", "can not get block", "error", err)
			}
			blocks = append(blocks, b)
		}
		return s.sendBlocksMessage(addr, blocks)
	}
	for i := data.From; i <= data.To; i++ {
		b, err := s.chain.GetBlock(i)
		if err != nil {
			s.Opts.Logger.Log("msg", "can not get block", "error", err)
		}
		blocks = append(blocks, b)
	}
	return s.sendBlocksMessage(addr, blocks)
}

func (s *Server) processGetStatusMessage(addr net.Addr, data *GetStatusMessage) error {

	fmt.Printf("id : %v, received get status message from %v, data+%v\n", s.Opts.ID, addr, data)
	//return s.sendStatusMessage(addr)

	if err := s.sendStatusMessage(addr); err != nil {
		s.Opts.Logger.Log("error", err)
	}
	return nil
}

func (s *Server) processStatusMessage(addr net.Addr, data *StatusMessage) error {
	fmt.Printf("id : %v, received  status message from %v, +%v\n", s.Opts.ID, addr, data)
	if s.chain.Height() >= data.CurrentHeight {
		err := fmt.Errorf("can not sync to a lower height local height -> %v, remote height -> %v", s.chain.Height(), data.CurrentHeight)
		s.Opts.Logger.Log("error", err)
		return err
	}

	return s.sendGetBlocksMessage(addr, data)
	// statusMessage := &StatusMessage{
	// 	ID : s.Opts.ID,
	// 	CurrentHeight: s.chain.Height(),
	// }
	// buff := &bytes.Buffer{}
	// if err := gob.NewEncoder(buff).Encode(statusMessage); err != nil {
	// 	return err
	// }
	// msg := NewMessage(network.MessageTypeStatusMessage, buff.Bytes())
	// msgBytes, err := msg.Bytes()
	// if err != nil {
	// 	return err
	// }
	// return s.Opts.Transport.SendMessage(addr, msgBytes)
	// return s.sendGetBlocksMessage(addr, data)
}

func (s *Server) processBlock(block *core.Block) error {
	if err := s.chain.AddBlock(block); err != nil {
		return err
	}
	s.broadcastBlock(block)
	return nil
}

func (s *Server) processTransaction(tx *core.Transaction) error {
	// check if the transaction is already in the mem pool
	txHash := tx.Hash(core.TransactionHashesr{})
	//s.Opts.Logger.Log("msg", "processing transaction", "hash", txHash)
	if s.memPool.Contains(txHash) {
		//s.Opts.Logger.Log("msg", "transaction is in the mempool", "hash", txHash)
		return nil
	}
	// check if the transaction is verified
	if err := tx.Verify(); err != nil {
		s.Opts.Logger.Log("error", err)
		return nil
	}

	tx.SetFirstSeen(time.Now().UnixNano())

	s.memPool.Add(tx)
	// s.Opts.Logger.Log("msg", "adding new tx to the mempool",
	// 	"hash", txHash,
	// 	"mempool length", s.memPool.AllCount())
	go s.broadcastTX(tx)
	return nil
}

func (s *Server) sendBlocksMessage(to net.Addr, blocks []*core.Block) error {
	blocksMessage := &BlocksMessage{
		Blocks: blocks,
	}
	fmt.Println(blocksMessage)
	buff := &bytes.Buffer{}
	if err := gob.NewEncoder(buff).Encode(blocksMessage); err != nil {
		s.Opts.Logger.Log("error", err)
		return err
	}
	msg := NewMessage(MessageTypeSyncBlocks, buff.Bytes())
	msgBytes, err := msg.Bytes()
	if err != nil {
		s.Opts.Logger.Log("error", err)
		return err
	}
	peer, ok := s.peerMap[to.String()]
	if !ok {
		err = fmt.Errorf("Can not send to unknown peer %v", to)
		s.Opts.Logger.Log("error", err)
		return err
	}
	if err := peer.sendMsg(msgBytes); err != nil {
		s.Opts.Logger.Log("error sending blocks to peers ", err)
		return err
	}
	return nil
}

func (s *Server) sendGetStatusMessage(peer *TCPPeer) error {
	getStatusMessage := new(GetStatusMessage)
	buff := &bytes.Buffer{}
	if err := gob.NewEncoder(buff).Encode(getStatusMessage); err != nil {
		return err
	}
	msg := NewMessage(MessageTypeGetStatusMessage, buff.Bytes())
	msgBytes, err := msg.Bytes()
	if err != nil {
		return err
	}
	return peer.sendMsg(msgBytes)

}

func (s *Server) sendStatusMessage(to net.Addr) error {
	stautusMessage := &StatusMessage{
		ID:            s.Opts.ID,
		CurrentHeight: s.chain.Height(),
	}
	buff := &bytes.Buffer{}
	if err := gob.NewEncoder(buff).Encode(stautusMessage); err != nil {
		s.Opts.Logger.Log("send status message error", err)
		return err
	}
	msg := NewMessage(MessageTypeStatusMessage, buff.Bytes())
	msgBytes, err := msg.Bytes()
	if err != nil {
		s.Opts.Logger.Log("send status message error", err)
		return err
	}
	peer, ok := s.peerMap[to.String()]
	if !ok {
		//s.Opts.Logger.Log("error", fmt.Errorf("can not status message to unknown peer"))
		return fmt.Errorf("can not send to unknown peer")
	}
	return peer.sendMsg(msgBytes)

}

func (s *Server) sendGetBlocksMessage(to net.Addr, data *StatusMessage) error {
	messageGetBlocks := &GetBlocksMessage{
		From: s.chain.Height(),
		To:   data.CurrentHeight,
	}
	buff := &bytes.Buffer{}
	if err := gob.NewEncoder(buff).Encode(messageGetBlocks); err != nil {
		s.Opts.Logger.Log("error here", err)
		return err
	}

	msg := NewMessage(MessageTypeGetBlocks, buff.Bytes())
	msgBytes, err := msg.Bytes()
	if err != nil {
		s.Opts.Logger.Log("error over here", err)
		return err
	}

	peer, ok := s.peerMap[to.String()]
	if !ok {
		err = fmt.Errorf("error sending message to unknown peer of address %v", to)
		s.Opts.Logger.Log("error there", err)
		return err
	}

	if err := peer.sendMsg(msgBytes); err != nil {
		s.Opts.Logger.Log("error in sending get blocks messge", err)
		return err
	}
	return nil
}

func (s *Server) validatorLoop() {
	ticker := time.NewTicker(s.Opts.BlockTime)
	s.Opts.Logger.Log("msg", "starting the validator", "Block Time", s.Opts.BlockTime)
	for {
		<-ticker.C
		s.CreateBlock()
	}
}

// func (s *Server) bootstrapNodes() error {
// 	for _, tr := range s.Opts.Transports {
// 		if s.Opts.Transport.Addr() != tr.Addr() {
// 			if err := s.Opts.Transport.Connect(tr); err != nil {
// 				return err
// 			}
// 			if err := s.sendGetStatusMessage(tr); err != nil {
// 				return err
// 			}
// 		}
// 	}
// 	return nil
// }

func (s *Server) CreateBlock() error {
	currHeight := s.chain.Height()
	currHeader, err := s.chain.GetHeader(int32(currHeight))
	if err != nil {
		s.Opts.Logger.Log("error", err)
		return err
	}
	// todo determine how to add transactions in a block, some complex scripts may be.
	txx := s.memPool.Pending()
	block, err := core.NewBlockFromPrevHeader(currHeader, txx)
	if err != nil {
		s.Opts.Logger.Log("error", err)
		return err
	}
	if err := block.Sign(*s.Opts.PrivateKey); err != nil {
		return err
	}
	err = s.chain.AddBlock(block)
	if err != nil {
		panic(err)
	}
	go s.broadcastBlock(block)
	s.memPool.ClearPending()
	return nil
}

func (s *Server) broadcastTX(tx *core.Transaction) error {
	buff := &bytes.Buffer{}
	if err := tx.Encode(core.NewGobTxEncoder(buff)); err != nil {
		return err
	}
	msg := Message{
		Header: MessageTypeTX,
		Data:   buff.Bytes(),
	}
	msgByte, err := msg.Bytes()
	if err != nil {
		return err
	}
	return s.broadcast(msgByte)
}

func (s *Server) broadcastBlock(block *core.Block) error {
	buff := &bytes.Buffer{}
	if err := block.Encoder(core.NewGobBlockEncoder(buff)); err != nil {
		return err
	}
	msg := NewMessage(MessageTypeBlock, buff.Bytes())
	msgByte, err := msg.Bytes()
	if err != nil {
		return err
	}
	return s.broadcast(msgByte)
}

func (s *Server) broadcast(data []byte) error {
	// for _, tr := range s.Opts.Transports {
	// 	if err := tr.Broadcast(data); err != nil {
	// 		return err
	// 	}
	// }
	s.mu.RLock()
	defer s.mu.RUnlock()
	for _, peer := range s.peerMap {
		if err := peer.sendMsg(data); err != nil {
			fmt.Printf("send messge to other peer error %v", err)
		}
	}
	return nil
}

func genesisBlock() *core.Block {
	header := &core.Header{
		Version:   1,
		DataHash:  types.Hash{},
		Timestamp: 000000,
		Height:    0,
	}
	block, _ := core.NewBlock(header, nil)
	return block
}
