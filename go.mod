module github.com/zenground0/ent

go 1.13

require (
	github.com/dgraph-io/badger/v2 v2.2007.2
	github.com/filecoin-project/filecoin-ffi v0.30.4-0.20200910194244-f640612a1a1f // indirect
	github.com/filecoin-project/go-address v0.0.5
	github.com/filecoin-project/go-hamt-ipld/v2 v2.0.0
	github.com/filecoin-project/go-state-types v0.0.0-20210119062722-4adba5aaea71
	github.com/filecoin-project/lotus v1.5.0-pre2.0.20210128154146-e2eb3e6e3502
	github.com/filecoin-project/specs-actors v0.9.13
	github.com/filecoin-project/specs-actors/v2 v2.3.4
	github.com/filecoin-project/specs-actors/v3 v3.0.1-0.20210128055125-ab0632b1c8fa
	github.com/ipfs/go-block-format v0.0.2
	github.com/ipfs/go-cid v0.0.7
	github.com/ipfs/go-datastore v0.4.5
	github.com/ipfs/go-ds-badger2 v0.1.1-0.20200708190120-187fc06f714e
	github.com/ipfs/go-ipfs-blockstore v1.0.3
	github.com/ipfs/go-ipld-cbor v0.0.5
	github.com/mitchellh/go-homedir v1.1.0
	github.com/urfave/cli/v2 v2.2.0
	github.com/whyrusleeping/cbor-gen v0.0.0-20210118024343-169e9d70c0c2
	golang.org/x/xerrors v0.0.0-20200804184101-5ec99f83aff1
)

replace github.com/filecoin-project/filecoin-ffi => ./extern/filecoin-ffi
