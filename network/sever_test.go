package network

// var localNode = NewLocalTransport("LOCAL")
// var remoteLate = NewLocalTransport("REMOTE_LATE")
// var remoteB = NewLocalTransport("REMOTE_B")
// var remoteC = NewLocalTransport("REMOTE_C")
// var remoteD = NewLocalTransport("REMOTE_D")
// var transports = []Transport{
// 	localNode, remoteLate, remoteB, remoteC, remoteD,
// }

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

// func TestStartServer(t *testing.T) {

// 	initTranports(t, transports)

// }

// func makeServer(t *testing.T, tr Transport, ID string, privKey *crypto.PrivateKey) *Server {
// 	opts := ServerOpts{
// 		ID:         ID,
// 		Transport:  tr,
// 		Transports: transports,
// 		PrivateKey: privKey,
// 	}
// 	server, err := NewServer(opts)
// 	assert.Nil(t, err)
// 	// fmt.Println(server)
// 	return server
// }

// func initTranports(t *testing.T, rtransports []Transport) {
// 	var s *Server
// 	for i := 0; i < len(transports); i++ {
// 		if i ==  {
// 			privKey := crypto.GeneratePrivateKey()
// 			s = makeServer(t, transports[i], string(transports[i].Addr()), &privKey)
// 		} else {
// 			s = makeServer(t, transports[i], string(transports[i].Addr()), nil)
// 		}
// 		go s.Start()
// 	}
// }
