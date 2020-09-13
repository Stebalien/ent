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

var dsPath = "~/.lotus/chain/datastore"

// Lifted from lotus/node/repo/fsrepo_ds.go
func chainBadgerDs(path string) (datastore.Batching, error) {
	opts := badger.DefaultOptions
	opts.GcInterval = 0 // disable GC for chain datastore

	opts.Options = dgbadger.DefaultOptions("").WithTruncate(true).
		WithValueThreshold(1 << 10)

	return badger.NewDatastore(path, &opts)
}

func loadLotusChainBstore(ctx context.Context) (blockstore.Blockstore, error) {
	// load chain datastore
	path, err := homedir.Expand(dsPath)
	if err != nil {
		return nil, err
	}
	ds, err := chainBadgerDs(path)
	if err != nil {
		return nil, err
	}
	return blockstore.NewBlockstore(ds), nil

}

// LoadCborStore loads the ~/.lotus chain datastore for chain traversal and state loading
func LoadCborStore(ctx context.Context) (cbornode.IpldStore, error) {
	bs, err := loadLotusChainBstore(ctx)
	if err != nil {
		return nil, err
	}
	return cbornode.NewCborStore(bs), nil
}

// ChainStateIterator moves from tip to genesis emiting parent state roots of blocks
type ChainStateIterator struct {
	currBlock *types.BlockHeader
}

func NewChainStateIterator(ctx context.Context, tipCid cid.Cid) (*ChainStateIterator, error) {
	bs, err := loadLotusChainBstore(ctx)
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
	}, nil
}

func (it *ChainStateIterator) Done() bool {
	if it.currBlock.Height == abi.ChainEpoch(0) {
		return true
	}
	return false
}

// Return the parent state root of the current block
func (it *ChainStateIterator) Val() cid.Cid {
	return it.currBlock.ParentStateRoot
}

// Moves iterator backwards towards genesis.  Noop at genesis
func (it *ChainStateIterator) Step(ctx context.Context) error {
	if it.Done() { // noop
		return nil
	}
	parents := it.currBlock.Parents
	// we don't care which, take the first one
	bs, err := loadLotusChainBstore(ctx)
	if err != nil {
		return err
	}
	raw, err := bs.Get(parents[0])
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