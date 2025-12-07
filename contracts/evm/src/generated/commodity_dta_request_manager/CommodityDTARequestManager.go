// Code generated — DO NOT EDIT.

package commodity_dta_request_manager

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"math/big"
	"reflect"
	"strings"

	ethereum "github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/event"
	"github.com/ethereum/go-ethereum/rpc"
	"google.golang.org/protobuf/types/known/emptypb"

	pb2 "github.com/smartcontractkit/chainlink-protos/cre/go/sdk"
	"github.com/smartcontractkit/chainlink-protos/cre/go/values/pb"
	"github.com/smartcontractkit/cre-sdk-go/capabilities/blockchain/evm"
	"github.com/smartcontractkit/cre-sdk-go/capabilities/blockchain/evm/bindings"
	"github.com/smartcontractkit/cre-sdk-go/cre"
)

var (
	_ = bytes.Equal
	_ = errors.New
	_ = fmt.Sprintf
	_ = big.NewInt
	_ = strings.NewReader
	_ = ethereum.NotFound
	_ = bind.Bind
	_ = common.Big1
	_ = types.BloomLookup
	_ = event.NewSubscription
	_ = abi.ConvertType
	_ = emptypb.Empty{}
	_ = pb.NewBigIntFromInt
	_ = pb2.AggregationType_AGGREGATION_TYPE_COMMON_PREFIX
	_ = bindings.FilterOptions{}
	_ = evm.FilterLogTriggerRequest{}
	_ = cre.ResponseBufferTooSmall
	_ = rpc.API{}
	_ = json.Unmarshal
	_ = reflect.Bool
)

var CommodityDTARequestManagerMetaData = &bind.MetaData{
	ABI: "[{\"type\":\"constructor\",\"inputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"createSubscriptionRequest\",\"inputs\":[{\"name\":\"quantityTons\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"paymentToken\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"createRedemptionRequest\",\"inputs\":[{\"name\":\"quantityTons\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"updateRequestStatus\",\"inputs\":[{\"name\":\"id\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"newStatus\",\"type\":\"uint8\",\"internalType\":\"uint8\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"getRequest\",\"inputs\":[{\"name\":\"id\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[{\"name\":\"\",\"type\":\"tuple\",\"internalType\":\"structCommodityDTARequestManager.Request\",\"components\":[{\"name\":\"id\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"requester\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"requestType\",\"type\":\"uint8\",\"internalType\":\"uint8\"},{\"name\":\"quantityTons\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"paymentToken\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"paymentAmount\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"status\",\"type\":\"uint8\",\"internalType\":\"uint8\"},{\"name\":\"createdAt\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"updatedAt\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getNextRequestId\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"event\",\"name\":\"SubscriptionRequested\",\"inputs\":[{\"name\":\"id\",\"type\":\"uint256\",\"indexed\":true,\"internalType\":\"uint256\"},{\"name\":\"requester\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"quantityTons\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"},{\"name\":\"paymentToken\",\"type\":\"address\",\"indexed\":false,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"RedemptionRequested\",\"inputs\":[{\"name\":\"id\",\"type\":\"uint256\",\"indexed\":true,\"internalType\":\"uint256\"},{\"name\":\"requester\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"quantityTons\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"RequestStatusUpdated\",\"inputs\":[{\"name\":\"id\",\"type\":\"uint256\",\"indexed\":true,\"internalType\":\"uint256\"},{\"name\":\"status\",\"type\":\"uint8\",\"indexed\":false,\"internalType\":\"uint8\"}],\"anonymous\":false}]",
}

// Structs
type Request struct {
	Id            *big.Int
	Requester     common.Address
	RequestType   uint8
	QuantityTons  *big.Int
	PaymentToken  common.Address
	PaymentAmount *big.Int
	Status        uint8
	CreatedAt     *big.Int
	UpdatedAt     *big.Int
}

// Contract Method Inputs
type CreateSubscriptionRequestInput struct {
	QuantityTons *big.Int
	PaymentToken common.Address
}

type CreateRedemptionRequestInput struct {
	QuantityTons *big.Int
}

type UpdateRequestStatusInput struct {
	Id        *big.Int
	NewStatus uint8
}

type GetRequestInput struct {
	Id *big.Int
}

// Contract Method Outputs

// Errors

// Events
type SubscriptionRequestedTopics struct {
	Id        *big.Int
	Requester common.Address
}

type SubscriptionRequestedDecoded struct {
	Id           *big.Int
	Requester    common.Address
	QuantityTons *big.Int
	PaymentToken common.Address
}

type RedemptionRequestedTopics struct {
	Id        *big.Int
	Requester common.Address
}

type RedemptionRequestedDecoded struct {
	Id           *big.Int
	Requester    common.Address
	QuantityTons *big.Int
}

type RequestStatusUpdatedTopics struct {
	Id *big.Int
}

type RequestStatusUpdatedDecoded struct {
	Id     *big.Int
	Status uint8
}

// Main Binding Type for CommodityDTARequestManager
type CommodityDTARequestManager struct {
	Address common.Address
	Options *bindings.ContractInitOptions
	ABI     *abi.ABI
	client  *evm.Client
	Codec   CommodityDTARequestManagerCodec
}

type CommodityDTARequestManagerCodec struct{}

func NewCommodityDTARequestManager(
	client *evm.Client,
	address common.Address,
	options *bindings.ContractInitOptions,
) (*CommodityDTARequestManager, error) {
	parsed, err := abi.JSON(strings.NewReader(CommodityDTARequestManagerMetaData.ABI))
	if err != nil {
		return nil, err
	}
	return &CommodityDTARequestManager{
		Address: address,
		Options: options,
		ABI:     &parsed,
		client:  client,
		Codec:   CommodityDTARequestManagerCodec{},
	}, nil
}

// Log Trigger Methods
func (c *CommodityDTARequestManager) LogTriggerSubscriptionRequestedLog(
	chainSelector uint64,
	confidenceLevel evm.ConfidenceLevel,
	topics []SubscriptionRequestedTopics,
) (*cre.Trigger[*bindings.DecodedLog[SubscriptionRequestedDecoded]], error) {
	filterOptions := []bindings.FilterOption{}

	for _, topic := range topics {
		var topicBytes [][]byte
		if topic.Id != nil {
			idBytes := common.LeftPadBytes(topic.Id.Bytes(), 32)
			topicBytes = append(topicBytes, idBytes)
		}
		if topic.Requester != (common.Address{}) {
			requesterBytes := topic.Requester.Bytes()
			topicBytes = append(topicBytes, common.LeftPadBytes(requesterBytes, 32))
		}
		if len(topicBytes) > 0 {
			filterOptions = append(filterOptions, bindings.WithTopicFilter(topicBytes...))
		}
	}

	filter, err := bindings.NewFilterLogTriggerRequest(
		c.client,
		chainSelector,
		confidenceLevel,
		c.Address,
		c.ABI.Events["SubscriptionRequested"],
		filterOptions...,
	)
	if err != nil {
		return nil, err
	}

	return cre.NewFilterLogTrigger(
		filter,
		bindings.NewLogAdapter[SubscriptionRequestedDecoded](c.ABI, "SubscriptionRequested"),
	), nil
}

func (c *CommodityDTARequestManager) LogTriggerRedemptionRequestedLog(
	chainSelector uint64,
	confidenceLevel evm.ConfidenceLevel,
	topics []RedemptionRequestedTopics,
) (*cre.Trigger[*bindings.DecodedLog[RedemptionRequestedDecoded]], error) {
	filterOptions := []bindings.FilterOption{}

	for _, topic := range topics {
		var topicBytes [][]byte
		if topic.Id != nil {
			idBytes := common.LeftPadBytes(topic.Id.Bytes(), 32)
			topicBytes = append(topicBytes, idBytes)
		}
		if topic.Requester != (common.Address{}) {
			requesterBytes := topic.Requester.Bytes()
			topicBytes = append(topicBytes, common.LeftPadBytes(requesterBytes, 32))
		}
		if len(topicBytes) > 0 {
			filterOptions = append(filterOptions, bindings.WithTopicFilter(topicBytes...))
		}
	}

	filter, err := bindings.NewFilterLogTriggerRequest(
		c.client,
		chainSelector,
		confidenceLevel,
		c.Address,
		c.ABI.Events["RedemptionRequested"],
		filterOptions...,
	)
	if err != nil {
		return nil, err
	}

	return cre.NewFilterLogTrigger(
		filter,
		bindings.NewLogAdapter[RedemptionRequestedDecoded](c.ABI, "RedemptionRequested"),
	), nil
}

func (c *CommodityDTARequestManager) LogTriggerRequestStatusUpdatedLog(
	chainSelector uint64,
	confidenceLevel evm.ConfidenceLevel,
	topics []RequestStatusUpdatedTopics,
) (*cre.Trigger[*bindings.DecodedLog[RequestStatusUpdatedDecoded]], error) {
	filterOptions := []bindings.FilterOption{}

	for _, topic := range topics {
		var topicBytes [][]byte
		if topic.Id != nil {
			idBytes := common.LeftPadBytes(topic.Id.Bytes(), 32)
			topicBytes = append(topicBytes, idBytes)
		}
		if len(topicBytes) > 0 {
			filterOptions = append(filterOptions, bindings.WithTopicFilter(topicBytes...))
		}
	}

	filter, err := bindings.NewFilterLogTriggerRequest(
		c.client,
		chainSelector,
		confidenceLevel,
		c.Address,
		c.ABI.Events["RequestStatusUpdated"],
		filterOptions...,
	)
	if err != nil {
		return nil, err
	}

	return cre.NewFilterLogTrigger(
		filter,
		bindings.NewLogAdapter[RequestStatusUpdatedDecoded](c.ABI, "RequestStatusUpdated"),
	), nil
}

// EVM Call Methods
func (c *CommodityDTARequestManager) GetRequest(
	chainSelector uint64,
	id *big.Int,
) (*evm.ContractRead[Request], error) {
	input := GetRequestInput{Id: id}
	return evm.NewContractRead[GetRequestInput, Request](
		c.client,
		chainSelector,
		c.Address,
		c.ABI,
		"getRequest",
		input,
	)
}

func (c *CommodityDTARequestManager) GetNextRequestId(
	chainSelector uint64,
) (*evm.ContractRead[struct{}, *big.Int], error) {
	return evm.NewContractRead[struct{}, *big.Int](
		c.client,
		chainSelector,
		c.Address,
		c.ABI,
		"getNextRequestId",
		struct{}{},
	)
}

// EVM Transaction Methods
func (c *CommodityDTARequestManager) CreateSubscriptionRequest(
	chainSelector uint64,
	quantityTons *big.Int,
	paymentToken common.Address,
	gasLimit uint64,
) (*evm.ContractWrite[CreateSubscriptionRequestInput], error) {
	input := CreateSubscriptionRequestInput{
		QuantityTons: quantityTons,
		PaymentToken: paymentToken,
	}
	return evm.NewContractWrite[CreateSubscriptionRequestInput](
		c.client,
		chainSelector,
		c.Address,
		c.ABI,
		"createSubscriptionRequest",
		input,
		gasLimit,
	)
}

func (c *CommodityDTARequestManager) CreateRedemptionRequest(
	chainSelector uint64,
	quantityTons *big.Int,
	gasLimit uint64,
) (*evm.ContractWrite[CreateRedemptionRequestInput], error) {
	input := CreateRedemptionRequestInput{
		QuantityTons: quantityTons,
	}
	return evm.NewContractWrite[CreateRedemptionRequestInput](
		c.client,
		chainSelector,
		c.Address,
		c.ABI,
		"createRedemptionRequest",
		input,
		gasLimit,
	)
}

func (c *CommodityDTARequestManager) UpdateRequestStatus(
	chainSelector uint64,
	id *big.Int,
	newStatus uint8,
	gasLimit uint64,
) (*evm.ContractWrite[UpdateRequestStatusInput], error) {
	input := UpdateRequestStatusInput{
		Id:        id,
		NewStatus: newStatus,
	}
	return evm.NewContractWrite[UpdateRequestStatusInput](
		c.client,
		chainSelector,
		c.Address,
		c.ABI,
		"updateRequestStatus",
		input,
		gasLimit,
	)
}

