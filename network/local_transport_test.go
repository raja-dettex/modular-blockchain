package network

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestConnect(t *testing.T) {
	tra := NewLocalTransport("A")
	trb := NewLocalTransport("B")
	tra.Connect(trb)
	trb.Connect(tra)
	assert.Equal(t, tra.peers[trb.Addr()], trb)
	assert.Equal(t, trb.peers[tra.Addr()], tra)
}

func TestSendMessage(t *testing.T) {
	tra := NewLocalTransport("A")
	trb := NewLocalTransport("B")
	tra.Connect(trb)
	trb.Connect(tra)
	Payload := []byte("hello here")
	assert.Nil(t, tra.SendMessage(trb.Addr(), Payload))
	message := <-trb.Consume()
	msg := make([]byte, len(Payload))
	n, err := message.Payload.Read(msg)
	assert.Nil(t, err)
	assert.Equal(t, n, len(Payload))
	assert.Equal(t, msg, Payload)
	assert.Equal(t, message.From, tra.Addr())
}

func TestBroadcast(t *testing.T) {
	tra := NewLocalTransport("A")
	trb := NewLocalTransport("B")
	trc := NewLocalTransport("C")
	tra.Connect(trb)
	tra.Connect(trc)
	payload := []byte("foo")
	err := tra.Broadcast(payload)
	fmt.Println(err)
	assert.Nil(t, err)
	rpcB := <-trb.Consume()
	msg := make([]byte, len(payload))
	n, err := rpcB.Payload.Read(msg)
	assert.Nil(t, err)
	assert.Equal(t, n, len(payload))
	assert.Equal(t, msg, payload)
	rpcC := <-trc.Consume()
	n, err = rpcC.Payload.Read(msg)
	assert.Nil(t, err)
	assert.Equal(t, n, len(payload))
	assert.Equal(t, msg, payload)
}
