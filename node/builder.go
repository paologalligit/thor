package node

import (
	"fmt"
	"github.com/gorilla/mux"
	"github.com/vechain/thor/v2/api/accounts"
	"github.com/vechain/thor/v2/api/utils"
	"github.com/vechain/thor/v2/bft"
	"github.com/vechain/thor/v2/chain"
	"github.com/vechain/thor/v2/genesis"
	"github.com/vechain/thor/v2/muxdb"
	"github.com/vechain/thor/v2/packer"
	"github.com/vechain/thor/v2/state"
	"github.com/vechain/thor/v2/thor"
	"github.com/vechain/thor/v2/tx"
	"time"
)

func RegisterAccountsAPI(
	gasLimit uint64,
	forkConfig thor.ForkConfig,
	bftFunc func(repo *chain.Repository) bft.Committer,
) func(repo *chain.Repository, stater *state.Stater, router *mux.Router) {
	return func(repo *chain.Repository, stater *state.Stater, router *mux.Router) {
		accounts.New(repo, stater, gasLimit, forkConfig, bftFunc(repo)).Mount(router, "/accounts")
	}
}

type Builder struct {
	dbFunc       func() *muxdb.MuxDB
	genesisFunc  func() *genesis.Genesis
	engineFunc   func(repo *chain.Repository) bft.Committer
	router       *mux.Router
	chain        *Chain
	transactions []*tx.Transaction
}

func (b *Builder) WithDB(memFunc func() *muxdb.MuxDB) *Builder {
	b.dbFunc = memFunc
	return b
}

func (b *Builder) WithGenesis(genesisFunc func() *genesis.Genesis) *Builder {
	b.genesisFunc = genesisFunc
	return b
}

func (b *Builder) WithBFTEngine(engineFunc func(repo *chain.Repository) bft.Committer) *Builder {
	b.engineFunc = engineFunc
	return b
}

func (b *Builder) WithAPIs(apis ...utils.APIServer) *Builder {

	b.router = mux.NewRouter()

	for _, api := range apis {
		api.MountDefaultPath(b.router)
	}

	return b
}

func (b *Builder) Build() (*Node, error) {

	for idx, newTx := range b.transactions {
		err := b.storeTx(newTx)
		if err != nil {
			return nil, fmt.Errorf("unable to store tx no %d : %w", idx, err)
		}
	}

	return &Node{
		chain:  b.chain,
		router: b.router,
	}, nil
}

func (b *Builder) WithChain(chain *Chain) *Builder {
	b.chain = chain
	return b
}

func (b *Builder) WithTransactions(transactions ...*tx.Transaction) *Builder {
	b.transactions = append(b.transactions, transactions...)
	return b
}

func (b *Builder) storeTx(transaction *tx.Transaction) error {
	blkPacker := packer.New(b.chain.Repo(), b.chain.Stater(), genesis.DevAccounts()[0].Address, &genesis.DevAccounts()[0].Address, thor.NoFork)
	flow, err := blkPacker.Schedule(b.chain.Repo().BestBlockSummary(), uint64(time.Now().Unix()))
	if err != nil {
		return fmt.Errorf("unable to schedule tx: %w", err)
	}
	err = flow.Adopt(transaction)
	if err != nil {
		return fmt.Errorf("unable to adopt tx: %w", err)
	}
	newBlk, stage, receipts, err := flow.Pack(genesis.DevAccounts()[0].PrivateKey, 0, false)
	if err != nil {
		return fmt.Errorf("unable to pack tx: %w", err)
	}
	if _, err := stage.Commit(); err != nil {
		return fmt.Errorf("unable to commit tx: %w", err)
	}
	if err := b.chain.Repo().AddBlock(newBlk, receipts, 0); err != nil {
		return fmt.Errorf("unable to add tx to repo: %w", err)
	}
	if err := b.chain.Repo().SetBestBlockID(newBlk.Header().ID()); err != nil {
		return fmt.Errorf("unable to set best block: %w", err)
	}
	return nil
}
