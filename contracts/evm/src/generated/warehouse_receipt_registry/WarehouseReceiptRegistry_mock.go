// Code generated — DO NOT EDIT.

//go:build !wasip1

package warehouse_receipt_registry

import (
	"errors"
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	evmmock "github.com/smartcontractkit/cre-sdk-go/capabilities/blockchain/evm/mock"
)

var (
	_ = errors.New
	_ = fmt.Errorf
	_ = big.NewInt
	_ = common.Big1
)

// WarehouseReceiptRegistryMock is a mock implementation of WarehouseReceiptRegistry for testing.
type WarehouseReceiptRegistryMock struct {
}

// NewWarehouseReceiptRegistryMock creates a new WarehouseReceiptRegistryMock for testing.
func NewWarehouseReceiptRegistryMock(address common.Address, clientMock *evmmock.ClientCapability) *WarehouseReceiptRegistryMock {
	mock := &WarehouseReceiptRegistryMock{}

	codec, err := NewCodec()
	if err != nil {
		panic("failed to create codec for mock: " + err.Error())
	}

	abi := codec.(*Codec).abi
	_ = abi

	funcMap := map[string]func([]byte) ([]byte, error){}

	evmmock.AddContractMock(address, clientMock, funcMap, nil)
	return mock
}
