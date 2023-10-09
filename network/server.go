package network

import (
	"bytes"
	"os"
	"time"

	"github.com/go-kit/log"
	"github.com/raja-dettex/modular-blockchain/core"
	"github.com/raja-dettex/modular-blockchain/crypto"
	"github.com/raja-dettex/modular-blockchain/types"
)

var (
	defaultBlockTime = time.Second * 5
)

type ServerOpts struct {
	ID            string
	Logger        log.Logger
	RPCDecodeFunc RPCDecodeFunc
	RPCProcessor  RPCProcessor
	Transports    []Transport
	PrivateKey    *crypto.PrivateKey
	BlockTime     time.Duration
}
type Server struct {
	Opts        ServerOpts
	isValidator bool
	chain       *core.Blockchain
	memPool     *TxPool
	rpcCh       chan RPC
	quitCh      chan struct{}
}

func NewServer(opts ServerOpts) (*Server, error) {
	if opts.BlockTime == time.Duration(0) {
		opts.BlockTime = defaultBlockTime
	}
	if opts.RPCDecodeFunc == nil {
		opts.RPCDecodeFunc = DefaultRPCDecodeFunc
	}

	s := &Server{
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
	bc, err := core.NewBlockChain(s.Opts.Logger, genesisBlock())
	if err != nil {
		return nil, err
	}
	s.chain = bc

	if s.isValidator {
		s.Opts.Logger.Log("is validator", s.isValidator, "key", s.Opts.PrivateKey)
		go s.validatorLoop()
	}

	return s, nil
}

func (s *Server) Start() {
	s.InitTransport()
free:
	for {
		select {
		case rpc := <-s.rpcCh:
			msg, err := s.Opts.RPCDecodeFunc(rpc)
			if err != nil {
				s.Opts.Logger.Log("error", err)
			}
			if err = s.Opts.RPCProcessor.ProcessMessage(msg); err != nil {
				s.Opts.Logger.Log("error", err)
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
		s.processBlock(t)
	}
	return nil
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

func (s *Server) validatorLoop() {
	ticker := time.NewTicker(s.Opts.BlockTime)
	s.Opts.Logger.Log("msg", "starting the validator", "Block Time", s.Opts.BlockTime)
	for {
		<-ticker.C
		s.CreateBlock()
	}
}

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
		return err
	}
	go s.broadcastBlock(block)
	s.memPool.ClearPending()
	return nil
}

func (s *Server) InitTransport() {
	for _, tra := range s.Opts.Transports {
		go func(tr Transport) {
			for rpc := range tr.Consume() {
				s.rpcCh <- rpc
			}
		}(tra)
	}
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
	for _, tr := range s.Opts.Transports {
		if err := tr.Broadcast(data); err != nil {
			return err
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
