package network

// func TestHandleTransaction(t *testing.T) {
// 	tra := NewLocalTransport("A")
// 	opts := ServerOpts{
// 		Transports: []Transport{tra},
// 		PrivateKey: crypto.GeneratePrivateKey(),
// 	}
// 	s := NewServer(opts)
// 	tx := core.NewTransaction([]byte("foo"))
// 	tx.Sign(opts.PrivateKey)
// 	s.HandleTransaction(tx)
// 	assert.Equal(t, 1, s.memPool.Len())
// }
