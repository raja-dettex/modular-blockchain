package api

import (
	"encoding/gob"
	"encoding/hex"
	"net/http"
	"strconv"

	"github.com/go-kit/log"
	"github.com/labstack/echo/v4"
	"github.com/raja-dettex/modular-blockchain/core"
	"github.com/raja-dettex/modular-blockchain/types"
)

type TxResponse struct {
	TxxCount uint
	TxxHash  []string
}

type Block struct {
	Hash          string
	Version       uint32
	DataHash      string
	PrevBlockHash string
	Height        int32
	Validator     string
	Signature     string
	TimeStamp     int64
	TxResponse    TxResponse
}

type Transaction struct {
	Data      string
	From      string
	Signature string
	Hash      string
	FirstSeen int64
}

type APIError struct {
	Error string
}

type ServerConfig struct {
	Logger     log.Logger
	ListenAddr string
}

type Server struct {
	Config ServerConfig
	bc     *core.Blockchain
	txChan chan *core.Transaction
}

func NewServer(cfg ServerConfig, bc *core.Blockchain, txChan chan *core.Transaction) *Server {
	return &Server{
		Config: cfg,
		bc:     bc,
		txChan: txChan,
	}
}

func (s *Server) Start() error {
	e := echo.New()
	e.GET("/block/:hashorid", s.handleGetBlock)
	e.GET("/tx/:hash", s.handleGetTx)
	e.POST("/tx", s.handlePostTx)
	return e.Start(s.Config.ListenAddr)
}

func (s *Server) handlePostTx(c echo.Context) error {
	tx := &core.Transaction{}
	if err := gob.NewDecoder(c.Request().Body).Decode(tx); err != nil {
		return c.JSON(http.StatusBadRequest, APIError{Error: err.Error()})
	}
	s.txChan <- tx
	return nil
}

func (s *Server) handleGetTx(c echo.Context) error {
	hash := c.Param("hash")
	txHash, err := hex.DecodeString(hash)
	if err != nil {
		//panic("here")
		return c.JSON(http.StatusBadRequest, APIError{Error: err.Error()})
	}
	tx, err := s.bc.GetTxByHash(types.HashFromBytes(txHash))
	if err != nil {
		//panic("here")
		return c.JSON(http.StatusBadRequest, APIError{Error: err.Error()})
	}
	jsonTx := txToJsonTx(tx)
	return c.JSON(http.StatusOK, jsonTx)
}

func (s *Server) handleGetBlock(c echo.Context) error {
	hashorid := c.Param("hashorid")

	// if id fetch the block from the height given
	id, err := strconv.Atoi(hashorid)
	if err != nil {
		hash, err := hex.DecodeString(hashorid)
		if err != nil {
			return err
		}
		block, err := s.bc.GetBlockByHash(types.HashFromBytes(hash))
		if err != nil {
			//panic("here")
			return c.JSON(http.StatusBadRequest, APIError{Error: err.Error()})
		}
		jsonBlock := blockToJsonBlock(block)
		return c.JSON(http.StatusOK, jsonBlock)

	}

	block, err := s.bc.GetBlock(uint32(id))
	if err != nil {
		//panic("here")
		return c.JSON(http.StatusBadRequest, APIError{Error: err.Error()})
	}
	jsonBlock := blockToJsonBlock(block)
	return c.JSON(http.StatusOK, jsonBlock)

}

func blockToJsonBlock(block *core.Block) Block {
	txResponse := TxResponse{
		TxxCount: uint(len(block.Transactions)),
		TxxHash:  make([]string, 0),
	}

	for _, tx := range block.Transactions {
		txResponse.TxxHash = append(txResponse.TxxHash, tx.Hash(core.TransactionHashesr{}).String())
	}
	jsonBlock := Block{
		Hash:          block.Hash(core.BlockHasher{}).String(),
		Version:       block.Header.Version,
		DataHash:      block.Header.DataHash.String(),
		PrevBlockHash: block.Header.PrevBlock.String(),
		Height:        block.Header.Height,
		Validator:     block.Validator.Address().String(),
		Signature:     block.Signature.String(),
		TxResponse:    txResponse,
		TimeStamp:     block.Header.Timestamp,
	}
	return jsonBlock
}

func txToJsonTx(tx *core.Transaction) Transaction {

	transaction := Transaction{
		Data:      string(tx.Data),
		From:      string(tx.From),
		Signature: tx.Signaure.String(),
		Hash:      tx.Hash(core.TransactionHashesr{}).String(),
		FirstSeen: tx.FirstSeen(),
	}
	return transaction
}
