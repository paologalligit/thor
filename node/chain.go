package node

import (
	"github.com/vechain/thor/v2/bft"
	"github.com/vechain/thor/v2/block"
	"github.com/vechain/thor/v2/chain"
	"github.com/vechain/thor/v2/genesis"
	"github.com/vechain/thor/v2/muxdb"
	"github.com/vechain/thor/v2/state"
)

type Chain struct {
	db           *muxdb.MuxDB
	genesis      *genesis.Genesis
	engine       bft.Committer
	repo         *chain.Repository
	stater       *state.Stater
	genesisBlock *block.Block
}

func (c *Chain) Repo() *chain.Repository {
	return c.repo
}

func (c *Chain) Stater() *state.Stater {
	return c.stater
}

func (c *Chain) GenesisBlock() *block.Block {
	return c.genesisBlock
}
