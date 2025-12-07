// Code generated — DO NOT EDIT.

package warehouse_receipt_registry

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

var WarehouseReceiptRegistryMetaData = &bind.MetaData{
	ABI: "[{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"receiptId\",\"type\":\"bytes32\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"cancelledBy\",\"type\":\"address\"}],\"name\":\"ReceiptCancelled\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"receiptId\",\"type\":\"bytes32\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"depositor\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"warehouse\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"string\",\"name\":\"commodityType\",\"type\":\"string\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"quantity\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"depositDate\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"expiryDate\",\"type\":\"uint256\"}],\"name\":\"ReceiptIssued\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"receiptId\",\"type\":\"bytes32\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"redeemer\",\"type\":\"address\"}],\"name\":\"ReceiptRedeemed\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"receiptId\",\"type\":\"bytes32\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"newQuantity\",\"type\":\"uint256\"}],\"name\":\"ReceiptUpdated\",\"type\":\"event\"}]",
}

// Structs

// Contract Method Inputs

// Contract Method Outputs

// Errors

// Events
// The <Event>Topics struct should be used as a filter (for log triggers).
// Note: It is only possible to filter on indexed fields.
// Indexed (string and bytes) fields will be of type common.Hash.
// They need to he (crypto.Keccak256) hashed and passed in.
// Indexed (tuple/slice/array) fields can be passed in as is, the Encode<Event>Topics function will handle the hashing.
//
// The <Event>Decoded struct will be the result of calling decode (Adapt) on the log trigger result.
// Indexed dynamic type fields will be of type common.Hash.

type ReceiptCancelledTopics struct {
	ReceiptId   [32]byte
	CancelledBy common.Address
}

type ReceiptCancelledDecoded struct {
	ReceiptId   [32]byte
	CancelledBy common.Address
}

type ReceiptIssuedTopics struct {
	ReceiptId [32]byte
	Depositor common.Address
	Warehouse common.Address
}

type ReceiptIssuedDecoded struct {
	ReceiptId     [32]byte
	Depositor     common.Address
	Warehouse     common.Address
	CommodityType string
	Quantity      *big.Int
	DepositDate   *big.Int
	ExpiryDate    *big.Int
}

type ReceiptRedeemedTopics struct {
	ReceiptId [32]byte
	Redeemer  common.Address
}

type ReceiptRedeemedDecoded struct {
	ReceiptId [32]byte
	Redeemer  common.Address
}

type ReceiptUpdatedTopics struct {
	ReceiptId [32]byte
}

type ReceiptUpdatedDecoded struct {
	ReceiptId   [32]byte
	NewQuantity *big.Int
}

// Main Binding Type for WarehouseReceiptRegistry
type WarehouseReceiptRegistry struct {
	Address common.Address
	Options *bindings.ContractInitOptions
	ABI     *abi.ABI
	client  *evm.Client
	Codec   WarehouseReceiptRegistryCodec
}

type WarehouseReceiptRegistryCodec interface {
	ReceiptCancelledLogHash() []byte
	EncodeReceiptCancelledTopics(evt abi.Event, values []ReceiptCancelledTopics) ([]*evm.TopicValues, error)
	DecodeReceiptCancelled(log *evm.Log) (*ReceiptCancelledDecoded, error)
	ReceiptIssuedLogHash() []byte
	EncodeReceiptIssuedTopics(evt abi.Event, values []ReceiptIssuedTopics) ([]*evm.TopicValues, error)
	DecodeReceiptIssued(log *evm.Log) (*ReceiptIssuedDecoded, error)
	ReceiptRedeemedLogHash() []byte
	EncodeReceiptRedeemedTopics(evt abi.Event, values []ReceiptRedeemedTopics) ([]*evm.TopicValues, error)
	DecodeReceiptRedeemed(log *evm.Log) (*ReceiptRedeemedDecoded, error)
	ReceiptUpdatedLogHash() []byte
	EncodeReceiptUpdatedTopics(evt abi.Event, values []ReceiptUpdatedTopics) ([]*evm.TopicValues, error)
	DecodeReceiptUpdated(log *evm.Log) (*ReceiptUpdatedDecoded, error)
}

func NewWarehouseReceiptRegistry(
	client *evm.Client,
	address common.Address,
	options *bindings.ContractInitOptions,
) (*WarehouseReceiptRegistry, error) {
	parsed, err := abi.JSON(strings.NewReader(WarehouseReceiptRegistryMetaData.ABI))
	if err != nil {
		return nil, err
	}
	codec, err := NewCodec()
	if err != nil {
		return nil, err
	}
	return &WarehouseReceiptRegistry{
		Address: address,
		Options: options,
		ABI:     &parsed,
		client:  client,
		Codec:   codec,
	}, nil
}

type Codec struct {
	abi *abi.ABI
}

func NewCodec() (WarehouseReceiptRegistryCodec, error) {
	parsed, err := abi.JSON(strings.NewReader(WarehouseReceiptRegistryMetaData.ABI))
	if err != nil {
		return nil, err
	}
	return &Codec{abi: &parsed}, nil
}

func (c *Codec) ReceiptCancelledLogHash() []byte {
	return c.abi.Events["ReceiptCancelled"].ID.Bytes()
}

func (c *Codec) EncodeReceiptCancelledTopics(
	evt abi.Event,
	values []ReceiptCancelledTopics,
) ([]*evm.TopicValues, error) {
	var receiptIdRule []interface{}
	for _, v := range values {
		if reflect.ValueOf(v.ReceiptId).IsZero() {
			receiptIdRule = append(receiptIdRule, common.Hash{})
			continue
		}
		fieldVal, err := bindings.PrepareTopicArg(evt.Inputs[0], v.ReceiptId)
		if err != nil {
			return nil, err
		}
		receiptIdRule = append(receiptIdRule, fieldVal)
	}
	var cancelledByRule []interface{}
	for _, v := range values {
		if reflect.ValueOf(v.CancelledBy).IsZero() {
			cancelledByRule = append(cancelledByRule, common.Hash{})
			continue
		}
		fieldVal, err := bindings.PrepareTopicArg(evt.Inputs[1], v.CancelledBy)
		if err != nil {
			return nil, err
		}
		cancelledByRule = append(cancelledByRule, fieldVal)
	}

	rawTopics, err := abi.MakeTopics(
		receiptIdRule,
		cancelledByRule,
	)
	if err != nil {
		return nil, err
	}

	return bindings.PrepareTopics(rawTopics, evt.ID.Bytes()), nil
}

// DecodeReceiptCancelled decodes a log into a ReceiptCancelled struct.
func (c *Codec) DecodeReceiptCancelled(log *evm.Log) (*ReceiptCancelledDecoded, error) {
	event := new(ReceiptCancelledDecoded)
	if err := c.abi.UnpackIntoInterface(event, "ReceiptCancelled", log.Data); err != nil {
		return nil, err
	}
	var indexed abi.Arguments
	for _, arg := range c.abi.Events["ReceiptCancelled"].Inputs {
		if arg.Indexed {
			if arg.Type.T == abi.TupleTy {
				// abigen throws on tuple, so converting to bytes to
				// receive back the common.Hash as is instead of error
				arg.Type.T = abi.BytesTy
			}
			indexed = append(indexed, arg)
		}
	}
	// Convert [][]byte → []common.Hash
	topics := make([]common.Hash, len(log.Topics))
	for i, t := range log.Topics {
		topics[i] = common.BytesToHash(t)
	}

	if err := abi.ParseTopics(event, indexed, topics[1:]); err != nil {
		return nil, err
	}
	return event, nil
}

func (c *Codec) ReceiptIssuedLogHash() []byte {
	return c.abi.Events["ReceiptIssued"].ID.Bytes()
}

func (c *Codec) EncodeReceiptIssuedTopics(
	evt abi.Event,
	values []ReceiptIssuedTopics,
) ([]*evm.TopicValues, error) {
	var receiptIdRule []interface{}
	for _, v := range values {
		if reflect.ValueOf(v.ReceiptId).IsZero() {
			receiptIdRule = append(receiptIdRule, common.Hash{})
			continue
		}
		fieldVal, err := bindings.PrepareTopicArg(evt.Inputs[0], v.ReceiptId)
		if err != nil {
			return nil, err
		}
		receiptIdRule = append(receiptIdRule, fieldVal)
	}
	var depositorRule []interface{}
	for _, v := range values {
		if reflect.ValueOf(v.Depositor).IsZero() {
			depositorRule = append(depositorRule, common.Hash{})
			continue
		}
		fieldVal, err := bindings.PrepareTopicArg(evt.Inputs[1], v.Depositor)
		if err != nil {
			return nil, err
		}
		depositorRule = append(depositorRule, fieldVal)
	}
	var warehouseRule []interface{}
	for _, v := range values {
		if reflect.ValueOf(v.Warehouse).IsZero() {
			warehouseRule = append(warehouseRule, common.Hash{})
			continue
		}
		fieldVal, err := bindings.PrepareTopicArg(evt.Inputs[2], v.Warehouse)
		if err != nil {
			return nil, err
		}
		warehouseRule = append(warehouseRule, fieldVal)
	}

	rawTopics, err := abi.MakeTopics(
		receiptIdRule,
		depositorRule,
		warehouseRule,
	)
	if err != nil {
		return nil, err
	}

	return bindings.PrepareTopics(rawTopics, evt.ID.Bytes()), nil
}

// DecodeReceiptIssued decodes a log into a ReceiptIssued struct.
func (c *Codec) DecodeReceiptIssued(log *evm.Log) (*ReceiptIssuedDecoded, error) {
	event := new(ReceiptIssuedDecoded)
	if err := c.abi.UnpackIntoInterface(event, "ReceiptIssued", log.Data); err != nil {
		return nil, err
	}
	var indexed abi.Arguments
	for _, arg := range c.abi.Events["ReceiptIssued"].Inputs {
		if arg.Indexed {
			if arg.Type.T == abi.TupleTy {
				// abigen throws on tuple, so converting to bytes to
				// receive back the common.Hash as is instead of error
				arg.Type.T = abi.BytesTy
			}
			indexed = append(indexed, arg)
		}
	}
	// Convert [][]byte → []common.Hash
	topics := make([]common.Hash, len(log.Topics))
	for i, t := range log.Topics {
		topics[i] = common.BytesToHash(t)
	}

	if err := abi.ParseTopics(event, indexed, topics[1:]); err != nil {
		return nil, err
	}
	return event, nil
}

func (c *Codec) ReceiptRedeemedLogHash() []byte {
	return c.abi.Events["ReceiptRedeemed"].ID.Bytes()
}

func (c *Codec) EncodeReceiptRedeemedTopics(
	evt abi.Event,
	values []ReceiptRedeemedTopics,
) ([]*evm.TopicValues, error) {
	var receiptIdRule []interface{}
	for _, v := range values {
		if reflect.ValueOf(v.ReceiptId).IsZero() {
			receiptIdRule = append(receiptIdRule, common.Hash{})
			continue
		}
		fieldVal, err := bindings.PrepareTopicArg(evt.Inputs[0], v.ReceiptId)
		if err != nil {
			return nil, err
		}
		receiptIdRule = append(receiptIdRule, fieldVal)
	}
	var redeemerRule []interface{}
	for _, v := range values {
		if reflect.ValueOf(v.Redeemer).IsZero() {
			redeemerRule = append(redeemerRule, common.Hash{})
			continue
		}
		fieldVal, err := bindings.PrepareTopicArg(evt.Inputs[1], v.Redeemer)
		if err != nil {
			return nil, err
		}
		redeemerRule = append(redeemerRule, fieldVal)
	}

	rawTopics, err := abi.MakeTopics(
		receiptIdRule,
		redeemerRule,
	)
	if err != nil {
		return nil, err
	}

	return bindings.PrepareTopics(rawTopics, evt.ID.Bytes()), nil
}

// DecodeReceiptRedeemed decodes a log into a ReceiptRedeemed struct.
func (c *Codec) DecodeReceiptRedeemed(log *evm.Log) (*ReceiptRedeemedDecoded, error) {
	event := new(ReceiptRedeemedDecoded)
	if err := c.abi.UnpackIntoInterface(event, "ReceiptRedeemed", log.Data); err != nil {
		return nil, err
	}
	var indexed abi.Arguments
	for _, arg := range c.abi.Events["ReceiptRedeemed"].Inputs {
		if arg.Indexed {
			if arg.Type.T == abi.TupleTy {
				// abigen throws on tuple, so converting to bytes to
				// receive back the common.Hash as is instead of error
				arg.Type.T = abi.BytesTy
			}
			indexed = append(indexed, arg)
		}
	}
	// Convert [][]byte → []common.Hash
	topics := make([]common.Hash, len(log.Topics))
	for i, t := range log.Topics {
		topics[i] = common.BytesToHash(t)
	}

	if err := abi.ParseTopics(event, indexed, topics[1:]); err != nil {
		return nil, err
	}
	return event, nil
}

func (c *Codec) ReceiptUpdatedLogHash() []byte {
	return c.abi.Events["ReceiptUpdated"].ID.Bytes()
}

func (c *Codec) EncodeReceiptUpdatedTopics(
	evt abi.Event,
	values []ReceiptUpdatedTopics,
) ([]*evm.TopicValues, error) {
	var receiptIdRule []interface{}
	for _, v := range values {
		if reflect.ValueOf(v.ReceiptId).IsZero() {
			receiptIdRule = append(receiptIdRule, common.Hash{})
			continue
		}
		fieldVal, err := bindings.PrepareTopicArg(evt.Inputs[0], v.ReceiptId)
		if err != nil {
			return nil, err
		}
		receiptIdRule = append(receiptIdRule, fieldVal)
	}

	rawTopics, err := abi.MakeTopics(
		receiptIdRule,
	)
	if err != nil {
		return nil, err
	}

	return bindings.PrepareTopics(rawTopics, evt.ID.Bytes()), nil
}

// DecodeReceiptUpdated decodes a log into a ReceiptUpdated struct.
func (c *Codec) DecodeReceiptUpdated(log *evm.Log) (*ReceiptUpdatedDecoded, error) {
	event := new(ReceiptUpdatedDecoded)
	if err := c.abi.UnpackIntoInterface(event, "ReceiptUpdated", log.Data); err != nil {
		return nil, err
	}
	var indexed abi.Arguments
	for _, arg := range c.abi.Events["ReceiptUpdated"].Inputs {
		if arg.Indexed {
			if arg.Type.T == abi.TupleTy {
				// abigen throws on tuple, so converting to bytes to
				// receive back the common.Hash as is instead of error
				arg.Type.T = abi.BytesTy
			}
			indexed = append(indexed, arg)
		}
	}
	// Convert [][]byte → []common.Hash
	topics := make([]common.Hash, len(log.Topics))
	for i, t := range log.Topics {
		topics[i] = common.BytesToHash(t)
	}

	if err := abi.ParseTopics(event, indexed, topics[1:]); err != nil {
		return nil, err
	}
	return event, nil
}

func (c WarehouseReceiptRegistry) WriteReport(
	runtime cre.Runtime,
	report *cre.Report,
	gasConfig *evm.GasConfig,
) cre.Promise[*evm.WriteReportReply] {
	return c.client.WriteReport(runtime, &evm.WriteCreReportRequest{
		Receiver:  c.Address.Bytes(),
		Report:    report,
		GasConfig: gasConfig,
	})
}

func (c *WarehouseReceiptRegistry) UnpackError(data []byte) (any, error) {
	switch common.Bytes2Hex(data[:4]) {
	default:
		return nil, errors.New("unknown error selector")
	}
}

// ReceiptCancelledTrigger wraps the raw log trigger and provides decoded ReceiptCancelledDecoded data
type ReceiptCancelledTrigger struct {
	cre.Trigger[*evm.Log, *evm.Log]                           // Embed the raw trigger
	contract                        *WarehouseReceiptRegistry // Keep reference for decoding
}

// Adapt method that decodes the log into ReceiptCancelled data
func (t *ReceiptCancelledTrigger) Adapt(l *evm.Log) (*bindings.DecodedLog[ReceiptCancelledDecoded], error) {
	// Decode the log using the contract's codec
	decoded, err := t.contract.Codec.DecodeReceiptCancelled(l)
	if err != nil {
		return nil, fmt.Errorf("failed to decode ReceiptCancelled log: %w", err)
	}

	return &bindings.DecodedLog[ReceiptCancelledDecoded]{
		Log:  l,        // Original log
		Data: *decoded, // Decoded data
	}, nil
}

func (c *WarehouseReceiptRegistry) LogTriggerReceiptCancelledLog(chainSelector uint64, confidence evm.ConfidenceLevel, filters []ReceiptCancelledTopics) (cre.Trigger[*evm.Log, *bindings.DecodedLog[ReceiptCancelledDecoded]], error) {
	event := c.ABI.Events["ReceiptCancelled"]
	topics, err := c.Codec.EncodeReceiptCancelledTopics(event, filters)
	if err != nil {
		return nil, fmt.Errorf("failed to encode topics for ReceiptCancelled: %w", err)
	}

	rawTrigger := evm.LogTrigger(chainSelector, &evm.FilterLogTriggerRequest{
		Addresses:  [][]byte{c.Address.Bytes()},
		Topics:     topics,
		Confidence: confidence,
	})

	return &ReceiptCancelledTrigger{
		Trigger:  rawTrigger,
		contract: c,
	}, nil
}

func (c *WarehouseReceiptRegistry) FilterLogsReceiptCancelled(runtime cre.Runtime, options *bindings.FilterOptions) (cre.Promise[*evm.FilterLogsReply], error) {
	if options == nil {
		return nil, errors.New("FilterLogs options are required.")
	}
	return c.client.FilterLogs(runtime, &evm.FilterLogsRequest{
		FilterQuery: &evm.FilterQuery{
			Addresses: [][]byte{c.Address.Bytes()},
			Topics: []*evm.Topics{
				{Topic: [][]byte{c.Codec.ReceiptCancelledLogHash()}},
			},
			BlockHash: options.BlockHash,
			FromBlock: pb.NewBigIntFromInt(options.FromBlock),
			ToBlock:   pb.NewBigIntFromInt(options.ToBlock),
		},
	}), nil
}

// ReceiptIssuedTrigger wraps the raw log trigger and provides decoded ReceiptIssuedDecoded data
type ReceiptIssuedTrigger struct {
	cre.Trigger[*evm.Log, *evm.Log]                           // Embed the raw trigger
	contract                        *WarehouseReceiptRegistry // Keep reference for decoding
}

// Adapt method that decodes the log into ReceiptIssued data
func (t *ReceiptIssuedTrigger) Adapt(l *evm.Log) (*bindings.DecodedLog[ReceiptIssuedDecoded], error) {
	// Decode the log using the contract's codec
	decoded, err := t.contract.Codec.DecodeReceiptIssued(l)
	if err != nil {
		return nil, fmt.Errorf("failed to decode ReceiptIssued log: %w", err)
	}

	return &bindings.DecodedLog[ReceiptIssuedDecoded]{
		Log:  l,        // Original log
		Data: *decoded, // Decoded data
	}, nil
}

func (c *WarehouseReceiptRegistry) LogTriggerReceiptIssuedLog(chainSelector uint64, confidence evm.ConfidenceLevel, filters []ReceiptIssuedTopics) (cre.Trigger[*evm.Log, *bindings.DecodedLog[ReceiptIssuedDecoded]], error) {
	event := c.ABI.Events["ReceiptIssued"]
	topics, err := c.Codec.EncodeReceiptIssuedTopics(event, filters)
	if err != nil {
		return nil, fmt.Errorf("failed to encode topics for ReceiptIssued: %w", err)
	}

	rawTrigger := evm.LogTrigger(chainSelector, &evm.FilterLogTriggerRequest{
		Addresses:  [][]byte{c.Address.Bytes()},
		Topics:     topics,
		Confidence: confidence,
	})

	return &ReceiptIssuedTrigger{
		Trigger:  rawTrigger,
		contract: c,
	}, nil
}

func (c *WarehouseReceiptRegistry) FilterLogsReceiptIssued(runtime cre.Runtime, options *bindings.FilterOptions) (cre.Promise[*evm.FilterLogsReply], error) {
	if options == nil {
		return nil, errors.New("FilterLogs options are required.")
	}
	return c.client.FilterLogs(runtime, &evm.FilterLogsRequest{
		FilterQuery: &evm.FilterQuery{
			Addresses: [][]byte{c.Address.Bytes()},
			Topics: []*evm.Topics{
				{Topic: [][]byte{c.Codec.ReceiptIssuedLogHash()}},
			},
			BlockHash: options.BlockHash,
			FromBlock: pb.NewBigIntFromInt(options.FromBlock),
			ToBlock:   pb.NewBigIntFromInt(options.ToBlock),
		},
	}), nil
}

// ReceiptRedeemedTrigger wraps the raw log trigger and provides decoded ReceiptRedeemedDecoded data
type ReceiptRedeemedTrigger struct {
	cre.Trigger[*evm.Log, *evm.Log]                           // Embed the raw trigger
	contract                        *WarehouseReceiptRegistry // Keep reference for decoding
}

// Adapt method that decodes the log into ReceiptRedeemed data
func (t *ReceiptRedeemedTrigger) Adapt(l *evm.Log) (*bindings.DecodedLog[ReceiptRedeemedDecoded], error) {
	// Decode the log using the contract's codec
	decoded, err := t.contract.Codec.DecodeReceiptRedeemed(l)
	if err != nil {
		return nil, fmt.Errorf("failed to decode ReceiptRedeemed log: %w", err)
	}

	return &bindings.DecodedLog[ReceiptRedeemedDecoded]{
		Log:  l,        // Original log
		Data: *decoded, // Decoded data
	}, nil
}

func (c *WarehouseReceiptRegistry) LogTriggerReceiptRedeemedLog(chainSelector uint64, confidence evm.ConfidenceLevel, filters []ReceiptRedeemedTopics) (cre.Trigger[*evm.Log, *bindings.DecodedLog[ReceiptRedeemedDecoded]], error) {
	event := c.ABI.Events["ReceiptRedeemed"]
	topics, err := c.Codec.EncodeReceiptRedeemedTopics(event, filters)
	if err != nil {
		return nil, fmt.Errorf("failed to encode topics for ReceiptRedeemed: %w", err)
	}

	rawTrigger := evm.LogTrigger(chainSelector, &evm.FilterLogTriggerRequest{
		Addresses:  [][]byte{c.Address.Bytes()},
		Topics:     topics,
		Confidence: confidence,
	})

	return &ReceiptRedeemedTrigger{
		Trigger:  rawTrigger,
		contract: c,
	}, nil
}

func (c *WarehouseReceiptRegistry) FilterLogsReceiptRedeemed(runtime cre.Runtime, options *bindings.FilterOptions) (cre.Promise[*evm.FilterLogsReply], error) {
	if options == nil {
		return nil, errors.New("FilterLogs options are required.")
	}
	return c.client.FilterLogs(runtime, &evm.FilterLogsRequest{
		FilterQuery: &evm.FilterQuery{
			Addresses: [][]byte{c.Address.Bytes()},
			Topics: []*evm.Topics{
				{Topic: [][]byte{c.Codec.ReceiptRedeemedLogHash()}},
			},
			BlockHash: options.BlockHash,
			FromBlock: pb.NewBigIntFromInt(options.FromBlock),
			ToBlock:   pb.NewBigIntFromInt(options.ToBlock),
		},
	}), nil
}

// ReceiptUpdatedTrigger wraps the raw log trigger and provides decoded ReceiptUpdatedDecoded data
type ReceiptUpdatedTrigger struct {
	cre.Trigger[*evm.Log, *evm.Log]                           // Embed the raw trigger
	contract                        *WarehouseReceiptRegistry // Keep reference for decoding
}

// Adapt method that decodes the log into ReceiptUpdated data
func (t *ReceiptUpdatedTrigger) Adapt(l *evm.Log) (*bindings.DecodedLog[ReceiptUpdatedDecoded], error) {
	// Decode the log using the contract's codec
	decoded, err := t.contract.Codec.DecodeReceiptUpdated(l)
	if err != nil {
		return nil, fmt.Errorf("failed to decode ReceiptUpdated log: %w", err)
	}

	return &bindings.DecodedLog[ReceiptUpdatedDecoded]{
		Log:  l,        // Original log
		Data: *decoded, // Decoded data
	}, nil
}

func (c *WarehouseReceiptRegistry) LogTriggerReceiptUpdatedLog(chainSelector uint64, confidence evm.ConfidenceLevel, filters []ReceiptUpdatedTopics) (cre.Trigger[*evm.Log, *bindings.DecodedLog[ReceiptUpdatedDecoded]], error) {
	event := c.ABI.Events["ReceiptUpdated"]
	topics, err := c.Codec.EncodeReceiptUpdatedTopics(event, filters)
	if err != nil {
		return nil, fmt.Errorf("failed to encode topics for ReceiptUpdated: %w", err)
	}

	rawTrigger := evm.LogTrigger(chainSelector, &evm.FilterLogTriggerRequest{
		Addresses:  [][]byte{c.Address.Bytes()},
		Topics:     topics,
		Confidence: confidence,
	})

	return &ReceiptUpdatedTrigger{
		Trigger:  rawTrigger,
		contract: c,
	}, nil
}

func (c *WarehouseReceiptRegistry) FilterLogsReceiptUpdated(runtime cre.Runtime, options *bindings.FilterOptions) (cre.Promise[*evm.FilterLogsReply], error) {
	if options == nil {
		return nil, errors.New("FilterLogs options are required.")
	}
	return c.client.FilterLogs(runtime, &evm.FilterLogsRequest{
		FilterQuery: &evm.FilterQuery{
			Addresses: [][]byte{c.Address.Bytes()},
			Topics: []*evm.Topics{
				{Topic: [][]byte{c.Codec.ReceiptUpdatedLogHash()}},
			},
			BlockHash: options.BlockHash,
			FromBlock: pb.NewBigIntFromInt(options.FromBlock),
			ToBlock:   pb.NewBigIntFromInt(options.ToBlock),
		},
	}), nil
}
