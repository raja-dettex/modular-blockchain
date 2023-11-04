package main

import (
	"bytes"
	"encoding/gob"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/raja-dettex/modular-blockchain/core"
	"github.com/raja-dettex/modular-blockchain/crypto"
	"github.com/raja-dettex/modular-blockchain/network"
	"github.com/raja-dettex/modular-blockchain/types"
	"github.com/raja-dettex/modular-blockchain/utils"
)

func main() {
	peers := []string{":3000"}
	pk := crypto.GeneratePrivateKey()
	localNode := makeServer("LOCAL_NODE", &pk, ":3000", peers, ":9000")
	go localNode.Start()
	// remoteNode := makeServer("REMOTE_NODE", nil, ":4000", []string{":5000"})
	// go remoteNode.Start()
	// remoteNodeB := makeServer("REMOTE_NODE_B", nil, ":5000", []string{})
	// go remoteNodeB.Start()
	go func() {
		time.Sleep(time.Second * 13)
		remoteNodeLate := makeServer("REMOTE_NODE_LATE", nil, ":6000", []string{":3000"}, "")
		go remoteNodeLate.Start()
	}()

	time.Sleep(time.Second * 1)
	if err := sendTransaction(pk); err != nil {
		fmt.Println(err)
	}
	// panic("here")

	//var blockingCh chan interface{}
	// tcpTransport := network.NewTCPTransport(":3000")
	// if err := tcpTransport.Start(); err != nil {
	// 	log.Fatal(err)
	// }
	// time.Sleep(time.Second * 2)
	// for {
	// 	go testConn()
	// 	time.Sleep(time.Second * 2)
	// }

	// tickerInterval := time.NewTicker(time.Second * 1)
	// collectionOwnerPrivKey := crypto.GeneratePrivateKey()
	// colletionHash := createCollectionTx(collectionOwnerPrivKey)
	// go func() {
	// 	for i := 0; i < 20; i++ {
	// 		<-tickerInterval.C
	// 		go nftMinter(collectionOwnerPrivKey, colletionHash)
	// 	}
	// }()

	select {}
	//<-blockingCh

}

func sendTransaction(privKey crypto.PrivateKey) error {
	tx := core.NewTransaction(nil)
	// privKey := crypto.GeneratePrivateKey()
	toPrivKey := crypto.GeneratePrivateKey()
	tx.To = toPrivKey.GeneratePublicKey()
	tx.Value = 666
	if err := tx.Sign(privKey); err != nil {
		return err
	}
	buff := new(bytes.Buffer)
	if err := gob.NewEncoder(buff).Encode(tx); err != nil {
		return err
	}
	req, err := http.NewRequest("POST", "http://localhost:9000/tx", buff)
	if err != nil {
		panic(err)
	}
	client := http.DefaultClient
	_, err = client.Do(req)
	return err

}

func createCollectionTx(privKey crypto.PrivateKey) types.Hash {
	//data := Contract()
	tx := core.NewTransaction(nil)
	tx.TxInnner = core.CollectionTx{
		Fee:      200,
		MetaData: []byte("collection nft"),
	}
	if err := tx.Sign(privKey); err != nil {
		panic(err)
	}
	buff := &bytes.Buffer{}
	if err := tx.Encode(core.NewGobTxEncoder(buff)); err != nil {
		panic(err)
	}
	req, err := http.NewRequest("POST", "http://localhost:9000/tx", buff)
	if err != nil {
		panic(err)
	}
	client := http.DefaultClient
	_, err = client.Do(req)
	if err != nil {
		panic(err)
	}
	return tx.Hash(core.TransactionHashesr{})
}

func nftMinter(privKey crypto.PrivateKey, collectionHash types.Hash) {
	//data := Contract()
	tx := core.NewTransaction(nil)

	metaData := map[string]any{
		"hello": "world",
		"age":   34,
		"color": "green",
	}
	metaBuff := new(bytes.Buffer)
	if err := json.NewEncoder(metaBuff).Encode(metaData); err != nil {
		panic(err)
	}
	tx.TxInnner = core.MintTx{
		Fee:             200,
		NFT:             utils.RandomHash(utils.RandomBytes(32)),
		MetaData:        metaBuff.Bytes(),
		Collection:      collectionHash,
		CollectionOwner: privKey.GeneratePublicKey(),
	}
	if err := tx.Sign(privKey); err != nil {
		panic(err)
	}
	buff := &bytes.Buffer{}
	if err := tx.Encode(core.NewGobTxEncoder(buff)); err != nil {
		panic(err)
	}
	req, err := http.NewRequest("POST", "http://localhost:9000/tx", buff)
	if err != nil {
		panic(err)
	}
	client := http.DefaultClient
	_, err = client.Do(req)
	if err != nil {
		panic(err)
	}
}
func makeServer(ID string, privKey *crypto.PrivateKey, addr string, peers []string, apiListenAddr string) *network.Server {
	opts := network.ServerOpts{
		ApiListenAddr: apiListenAddr,
		ListenAddr:    addr,
		ID:            ID,
		PrivateKey:    privKey,
		SeedNodes:     peers,
	}
	server, err := network.NewServer(opts)
	if err != nil {
		log.Fatal(err)
	}
	return server
}

func Contract() []byte {
	keyBytes := []byte{0x061, 0x0c, 0x064, 0x0c, 0x02, 0x0a, 0x0d, 0x0ae}
	data := []byte{0x02, 0x0a, 0x03, 0x0a, 0x0e, 0x061, 0x0c, 0x064, 0x0c, 0x02, 0x0a, 0x0d, 0x0f}
	data = append(data, keyBytes...)
	return data
}

// var (
// 	transports = []network.Transport{
// 		network.NewLocalTransport("LOCAL"),
// 		//network.NewLocalTransport("REMOTE_A"),
// 		// network.NewLocalTransport("REMOTE_B"),
// 		// network.NewLocalTransport("REMOTE_C"),
// 		network.NewLocalTransport("REMOTE_LATE"),
// 	}
// )

// func main() {

// 	//initTransports(transports)
// 	localNode := transports[0]
// 	lateTr := transports[1]

// 	// go func() {
// 	// 	i := 0
// 	// 	for {
// 	// 		sendTransaction(traRemoteA, traLocal.Addr(), i)
// 	// 		//traRemote.SendMessage(traLocal.Addr(), []byte("hello world"))
// 	// 		i++
// 	// 		time.Sleep(time.Second * 2)
// 	// 	}
// 	// }()

// 	go func() {
// 		//panic("here!!!!")
// 		time.Sleep(time.Second * 7)

// 		lateServer := makeServer(string(lateTr.Addr()), nil, lateTr)
// 		// if err := sendGetStatusMessage(trLate, "REMOTE_C"); err != nil {
// 		// 	log.Fatal(err)
// 		// }
// 		//panic("here!!!")
// 		go lateServer.Start()
// 	}()
// 	privKey := crypto.GeneratePrivateKey()

// 	localServer := makeServer("LOCAL", &privKey, localNode)
// 	localServer.Start()
// }

// func initTransports(trs []network.Transport) {
// 	//i := 0
// 	for i := 1; i < len(trs)-1; i++ {
// 		s := makeServer(fmt.Sprintf("REMOTE_%v", i), nil, trs[i])
// 		go s.Start()
// 		//i++
// 	}
// }

// func sendTransaction(from network.Transport, to network.NetAdddr, i int) error {
// 	privKey := crypto.GeneratePrivateKey()
// 	// data := []byte(fmt.Sprintf("transaction : [%v]", i*100))
// 	data := Contract()
// 	tx := core.NewTransaction(data)
// 	if err := tx.Sign(privKey); err != nil {
// 		return err
// 	}
// 	buff := &bytes.Buffer{}
// 	if err := tx.Encode(core.NewGobTxEncoder(buff)); err != nil {
// 		return err
// 	}
// 	msg := network.NewMessage(network.MessageTypeTX, buff.Bytes())
// 	msgByte, err := msg.Bytes()
// 	if err != nil {
// 		return err
// 	}
// 	from.SendMessage(to, msgByte)
// 	return nil
// }

// func sendGetStatusMessage(tr network.Transport, to network.NetAdddr) error {
// 	getStatusMessage := &network.GetStatusMessage{}
// 	buff := &bytes.Buffer{}
// 	if err := gob.NewEncoder(buff).Encode(getStatusMessage); err != nil {
// 		return err
// 	}
// 	msg := network.NewMessage(network.MessageTypeGetStatusMessage, buff.Bytes())
// 	msgBytes, err := msg.Bytes()
// 	if err != nil {
// 		return err
// 	}
// 	return tr.SendMessage(to, msgBytes)
// }
