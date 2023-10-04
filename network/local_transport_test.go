package network

import (
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
	assert.Equal(t, message.Payload, Payload)
	assert.Equal(t, message.From, trb.Addr())
}
