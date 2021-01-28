package lib

import (
	address "github.com/filecoin-project/go-address"
	states3 "github.com/filecoin-project/specs-actors/v3/actors/states"
)

type StateTree interface {
	GetActor(addr address.Address) (*states3.Actor, bool, error)
	SetActor(addr address.Address, actor *states3.Actor) error
	ForEach(fn func(addr address.Address, actor *states3.Actor) error) error
	ForEachKey(fn func(addr address.Address) error) error
}
