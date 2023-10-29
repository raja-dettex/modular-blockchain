package network

// func TestBlocksMessage(t *testing.T) {
// 	txx := []*core.Transaction{}
// 	dataHash, err := core.CalculateDataHash(txx)
// 	assert.Nil(t, err)
// 	header := &core.Header{
// 		Version: uint32(1),
// 		DataHash: dataHash,
// 		PrevBlock: types.Hash{},
// 		Timestamp: time.Now().UnixNano(),
// 		Height: 0,
// 	}
// 	privKey := crypto.GeneratePrivateKey()
// 	block1, err  := core.NewBlock(header, txx)
// 	assert.Nil(t, err)
// 	block1.Sign(privKey)

// 	block2, err  := core.NewBlockFromPrevHeader(header, []*core.Transaction{})
// 	assert.Nil(t, err)
// 	block2.Sign(privKey)
// 	block3, err  := core.NewBlockFromPrevHeader(block2.Header, []*core.Transaction{})
// 	assert.Nil(t, err)
// 	block3.Sign(privKey)
// 	core.NewBlockchainWith

// }
