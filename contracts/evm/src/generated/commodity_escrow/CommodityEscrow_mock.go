// Code generated — DO NOT EDIT.

//go:build !wasip1

package commodity_escrow

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

// CommodityEscrowMock is a mock implementation of CommodityEscrow for testing.
type CommodityEscrowMock struct {
}

// NewCommodityEscrowMock creates a new CommodityEscrowMock for testing.
func NewCommodityEscrowMock(address common.Address, clientMock *evmmock.ClientCapability) *CommodityEscrowMock {
	mock := &CommodityEscrowMock{}

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
