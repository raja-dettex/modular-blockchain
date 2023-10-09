package main

import (
	"bytes"
	"fmt"
	"log"
	"time"

	"github.com/raja-dettex/modular-blockchain/core"
	"github.com/raja-dettex/modular-blockchain/crypto"
	"github.com/raja-dettex/modular-blockchain/network"
)

func main() {
	traLocal := network.NewLocalTransport("LOCAL")
	traRemoteA := network.NewLocalTransport("REMOTE_A")
	traRemoteB := network.NewLocalTransport("REMOTE_B")
	traRemoteC := network.NewLocalTransport("REMOTE_C")
	traLocal.Connect(traRemoteA)
	traRemoteA.Connect(traLocal)
	traRemoteA.Connect(traRemoteB)
	traRemoteB.Connect(traRemoteC)
	trs := []network.Transport{traRemoteA, traRemoteB, traRemoteC}

	initTransports(trs)
	go func() {
		i := 0
		for {
			sendTransaction(traRemoteA, traLocal.Addr(), i)
			//traRemote.SendMessage(traLocal.Addr(), []byte("hello world"))
			i++
			time.Sleep(time.Second * 2)
		}
	}()

	go func() {
		time.Sleep(time.Second * 7)
		trLate := network.NewLocalTransport("REMOTE_LATE")
		traRemoteC.Connect(trLate)
		lateServer := makeServer(string(trLate.Addr()), nil, trLate)
		go lateServer.Start()
	}()
	privKey := crypto.GeneratePrivateKey()

	localServer := makeServer("LOCAL", &privKey, traLocal)
	localServer.Start()
}

func makeServer(ID string, privKey *crypto.PrivateKey, tr network.Transport) *network.Server {
	opts := network.ServerOpts{
		ID:         ID,
		PrivateKey: privKey,
		Transports: []network.Transport{tr},
	}
	server, err := network.NewServer(opts)
	if err != nil {
		log.Fatal(err)
	}
	return server
}

func initTransports(trs []network.Transport) {
	i := 0
	for _, tr := range trs {
		s := makeServer(fmt.Sprintf("REMOTE_%v", i), nil, tr)
		go s.Start()
		i++
	}
}

func sendTransaction(from network.Transport, to network.NetAdddr, i int) error {
	privKey := crypto.GeneratePrivateKey()
	data := []byte(fmt.Sprintf("transaction : [%v]", i*100))
	tx := core.NewTransaction(data)
	if err := tx.Sign(privKey); err != nil {
		return err
	}
	buff := &bytes.Buffer{}
	if err := tx.Encode(core.NewGobTxEncoder(buff)); err != nil {
		return err
	}
	msg := network.NewMessage(network.MessageTypeTX, buff.Bytes())
	msgByte, err := msg.Bytes()
	if err != nil {
		return err
	}
	from.SendMessage(to, msgByte)
	return nil
}
