package lib

import (
	"context"

	dgbadger "github.com/dgraph-io/badger/v2"
	"github.com/filecoin-project/go-state-types/abi"
	"github.com/filecoin-project/lotus/chain/types"
	cid "github.com/ipfs/go-cid"
	datastore "github.com/ipfs/go-datastore"
	badger "github.com/ipfs/go-ds-badger2"
	"github.com/ipfs/go-ipfs-blockstore"
	"github.com/ipfs/go-ipld-cbor"
	"github.com/mitchellh/go-homedir"
)

var lotusPath = "~/.lotus/datastore/chain"
var entPath = "~/.ent/datastore/chain"

type Chain struct {
	cachedBs blockstore.Blockstore
}

// Lifted from lotus/node/repo/fsrepo_ds.go
func chainBadgerDs(path string) (datastore.Batching, error) {
	opts := badger.DefaultOptions
	opts.GcInterval = 0 // disable GC for chain datastore

	opts.Options = dgbadger.DefaultOptions("").WithTruncate(true).
		WithValueThreshold(1 << 10)

	return badger.NewDatastore(path, &opts)
}

func (c *Chain) loadRedirectBstore(ctx context.Context) (blockstore.Blockstore, error) {
	if c.cachedBs != nil {
		return c.cachedBs, nil
	}
	// load lotus chain datastore
	lotusExpPath, err := homedir.Expand(lotusPath)
	if err != nil {
		return nil, err
	}
	lotusDS, err := chainBadgerDs(lotusExpPath)
	if err != nil {
		return nil, err
	}

	// load ent chain datastore
	entExpPath, err := homedir.Expand(entPath)
	if err != nil {
		return nil, err
	}
	entDS, err := chainBadgerDs(entExpPath)
	if err != nil {
		return nil, err 
	}
	return NewRedirectBlockstore(blockstore.NewBlockstore(lotusDS), blockstore.NewBlockstore(entDS)), nil
}

// LoadCborStore loads the ~/.lotus chain datastore for chain traversal and state loading
func (c *Chain) LoadCborStore(ctx context.Context) (cbornode.IpldStore, error) {
	bs, err := c.loadRedirectBstore(ctx)
	if err != nil {
		return nil, err
	}
	return cbornode.NewCborStore(bs), nil
}

// ChainStateIterator moves from tip to genesis emiting parent state roots of blocks
type ChainStateIterator struct {
	bs        blockstore.Blockstore
	currBlock *types.BlockHeader
}

type IterVal struct {
	Height int64
	State cid.Cid
}

func (c *Chain) NewChainStateIterator(ctx context.Context, tipCid cid.Cid) (*ChainStateIterator, error) {
	bs, err := c.loadRedirectBstore(ctx)
	if err != nil {
		return nil, err
	}
	// get starting block
	raw, err := bs.Get(tipCid)
	if err != nil {
		return nil, err
	}

	blk, err := types.DecodeBlock(raw.RawData())
	if err != nil {
		return nil, err
	}

	return &ChainStateIterator{
		currBlock: blk,
		bs:        bs,
	}, nil
}

func (it *ChainStateIterator) Done() bool {
	if it.currBlock.Height == abi.ChainEpoch(0) {
		return true
	}
	return false
}

// Return the parent state root of the current block
func (it *ChainStateIterator) Val() IterVal {
	return IterVal{
		State: it.currBlock.ParentStateRoot,
		Height: int64(it.currBlock.Height),
	}
}

// Moves iterator backwards towards genesis.  Noop at genesis
func (it *ChainStateIterator) Step(ctx context.Context) error {
	if it.Done() { // noop
		return nil
	}
	parents := it.currBlock.Parents
	// we don't care which, take the first one
	raw, err := it.bs.Get(parents[0])
	if err != nil {
		return err
	}
	nextBlock, err := types.DecodeBlock(raw.RawData())
	if err != nil {
		return err
	}
	it.currBlock = nextBlock
	return nil
}
