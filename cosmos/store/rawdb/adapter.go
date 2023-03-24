package rawdb

import (
	"errors"

	storetypes "cosmossdk.io/store/types"
	"github.com/ethereum/go-ethereum/ethdb"
)

type Adapter interface {
	ethdb.KeyValueReader
	ethdb.KeyValueWriter
}

type adapter struct {
	storetypes.KVStore
}

// use prefix store as input with "evm" + "historical", kvstore prefix
func NewAdapter(kvstore storetypes.KVStore) Adapter {
	return &adapter{kvstore}
}

func (a *adapter) Put(key []byte, value []byte) error {
	a.KVStore.Set(key, value)
	return nil
}

func (a *adapter) Get(key []byte) ([]byte, error) {
	value := a.KVStore.Get(key)
	if value == nil {
		return nil, errors.New("no such key")
	}
	return value, nil
}

func (a *adapter) Has(key []byte) (bool, error) {
	return a.KVStore.Has(key), nil
}

func (a *adapter) Delete(key []byte) error {
	a.KVStore.Delete(key)
	return nil
}
