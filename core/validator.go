package core

import (
	"errors"
	"fmt"
)

var (
	ErrBlockKnown = errors.New("this block is known")
)

type Validator interface {
	Validate(*Block) error
}

type BlockValidator struct {
	bc *Blockchain
}

func NewBlockValidator(bc *Blockchain) *BlockValidator {
	return &BlockValidator{
		bc: bc,
	}

}

func (v *BlockValidator) Validate(b *Block) error {
	if v.bc.HasBlock(uint32(b.Header.Height)) {
		// err := fmt.Errorf("blockchain already contains the block with height {%v} hash {%d}", b.Header.Height, b.Hash(BlockHasher{}))
		return ErrBlockKnown
	}
	if err := b.Verify(); err != nil {
		return err
	}
	if b.Header.Height > int32(v.bc.Height()+1) {
		err := fmt.Errorf("block height can not be greater than the height of the chain %s", b.Hash(BlockHasher{}))

		return err
	}
	prevHeader, err := v.bc.GetHeader(b.Header.Height - 1)
	if err != nil {
		return err
	}
	hash := BlockHasher{}.Hash(prevHeader)
	if hash != b.Header.PrevBlock {
		err := fmt.Errorf("invalid prev hash")

		return err
	}
	return nil
}
