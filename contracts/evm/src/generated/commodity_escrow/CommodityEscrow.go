// Code generated — DO NOT EDIT.

package commodity_escrow

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

var CommodityEscrowMetaData = &bind.MetaData{
	ABI: "[{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"escrowId\",\"type\":\"bytes32\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"cancelledBy\",\"type\":\"address\"}],\"name\":\"EscrowCancelled\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"escrowId\",\"type\":\"bytes32\"}],\"name\":\"EscrowCompleted\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"escrowId\",\"type\":\"bytes32\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"confirmer\",\"type\":\"address\"}],\"name\":\"EscrowConfirmed\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"escrowId\",\"type\":\"bytes32\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"seller\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"buyer\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"commodityToken\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"tokenAmount\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"pricePerToken\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"bytes32\",\"name\":\"receiptId\",\"type\":\"bytes32\"}],\"name\":\"EscrowCreated\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"escrowId\",\"type\":\"bytes32\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"disputer\",\"type\":\"address\"}],\"name\":\"EscrowDisputed\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"escrowId\",\"type\":\"bytes32\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"funder\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"EscrowFunded\",\"type\":\"event\"}]",
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

type EscrowCancelledTopics struct {
	EscrowId    [32]byte
	CancelledBy common.Address
}

type EscrowCancelledDecoded struct {
	EscrowId    [32]byte
	CancelledBy common.Address
}

type EscrowCompletedTopics struct {
	EscrowId [32]byte
}

type EscrowCompletedDecoded struct {
	EscrowId [32]byte
}

type EscrowConfirmedTopics struct {
	EscrowId  [32]byte
	Confirmer common.Address
}

type EscrowConfirmedDecoded struct {
	EscrowId  [32]byte
	Confirmer common.Address
}

type EscrowCreatedTopics struct {
	EscrowId [32]byte
	Seller   common.Address
	Buyer    common.Address
}

type EscrowCreatedDecoded struct {
	EscrowId       [32]byte
	Seller         common.Address
	Buyer          common.Address
	CommodityToken common.Address
	TokenAmount    *big.Int
	PricePerToken  *big.Int
	ReceiptId      [32]byte
}

type EscrowDisputedTopics struct {
	EscrowId [32]byte
	Disputer common.Address
}

type EscrowDisputedDecoded struct {
	EscrowId [32]byte
	Disputer common.Address
}

type EscrowFundedTopics struct {
	EscrowId [32]byte
	Funder   common.Address
}

type EscrowFundedDecoded struct {
	EscrowId [32]byte
	Funder   common.Address
	Amount   *big.Int
}

// Main Binding Type for CommodityEscrow
type CommodityEscrow struct {
	Address common.Address
	Options *bindings.ContractInitOptions
	ABI     *abi.ABI
	client  *evm.Client
	Codec   CommodityEscrowCodec
}

type CommodityEscrowCodec interface {
	EscrowCancelledLogHash() []byte
	EncodeEscrowCancelledTopics(evt abi.Event, values []EscrowCancelledTopics) ([]*evm.TopicValues, error)
	DecodeEscrowCancelled(log *evm.Log) (*EscrowCancelledDecoded, error)
	EscrowCompletedLogHash() []byte
	EncodeEscrowCompletedTopics(evt abi.Event, values []EscrowCompletedTopics) ([]*evm.TopicValues, error)
	DecodeEscrowCompleted(log *evm.Log) (*EscrowCompletedDecoded, error)
	EscrowConfirmedLogHash() []byte
	EncodeEscrowConfirmedTopics(evt abi.Event, values []EscrowConfirmedTopics) ([]*evm.TopicValues, error)
	DecodeEscrowConfirmed(log *evm.Log) (*EscrowConfirmedDecoded, error)
	EscrowCreatedLogHash() []byte
	EncodeEscrowCreatedTopics(evt abi.Event, values []EscrowCreatedTopics) ([]*evm.TopicValues, error)
	DecodeEscrowCreated(log *evm.Log) (*EscrowCreatedDecoded, error)
	EscrowDisputedLogHash() []byte
	EncodeEscrowDisputedTopics(evt abi.Event, values []EscrowDisputedTopics) ([]*evm.TopicValues, error)
	DecodeEscrowDisputed(log *evm.Log) (*EscrowDisputedDecoded, error)
	EscrowFundedLogHash() []byte
	EncodeEscrowFundedTopics(evt abi.Event, values []EscrowFundedTopics) ([]*evm.TopicValues, error)
	DecodeEscrowFunded(log *evm.Log) (*EscrowFundedDecoded, error)
}

func NewCommodityEscrow(
	client *evm.Client,
	address common.Address,
	options *bindings.ContractInitOptions,
) (*CommodityEscrow, error) {
	parsed, err := abi.JSON(strings.NewReader(CommodityEscrowMetaData.ABI))
	if err != nil {
		return nil, err
	}
	codec, err := NewCodec()
	if err != nil {
		return nil, err
	}
	return &CommodityEscrow{
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

func NewCodec() (CommodityEscrowCodec, error) {
	parsed, err := abi.JSON(strings.NewReader(CommodityEscrowMetaData.ABI))
	if err != nil {
		return nil, err
	}
	return &Codec{abi: &parsed}, nil
}

func (c *Codec) EscrowCancelledLogHash() []byte {
	return c.abi.Events["EscrowCancelled"].ID.Bytes()
}

func (c *Codec) EncodeEscrowCancelledTopics(
	evt abi.Event,
	values []EscrowCancelledTopics,
) ([]*evm.TopicValues, error) {
	var escrowIdRule []interface{}
	for _, v := range values {
		if reflect.ValueOf(v.EscrowId).IsZero() {
			escrowIdRule = append(escrowIdRule, common.Hash{})
			continue
		}
		fieldVal, err := bindings.PrepareTopicArg(evt.Inputs[0], v.EscrowId)
		if err != nil {
			return nil, err
		}
		escrowIdRule = append(escrowIdRule, fieldVal)
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
		escrowIdRule,
		cancelledByRule,
	)
	if err != nil {
		return nil, err
	}

	return bindings.PrepareTopics(rawTopics, evt.ID.Bytes()), nil
}

// DecodeEscrowCancelled decodes a log into a EscrowCancelled struct.
func (c *Codec) DecodeEscrowCancelled(log *evm.Log) (*EscrowCancelledDecoded, error) {
	event := new(EscrowCancelledDecoded)
	if err := c.abi.UnpackIntoInterface(event, "EscrowCancelled", log.Data); err != nil {
		return nil, err
	}
	var indexed abi.Arguments
	for _, arg := range c.abi.Events["EscrowCancelled"].Inputs {
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

func (c *Codec) EscrowCompletedLogHash() []byte {
	return c.abi.Events["EscrowCompleted"].ID.Bytes()
}

func (c *Codec) EncodeEscrowCompletedTopics(
	evt abi.Event,
	values []EscrowCompletedTopics,
) ([]*evm.TopicValues, error) {
	var escrowIdRule []interface{}
	for _, v := range values {
		if reflect.ValueOf(v.EscrowId).IsZero() {
			escrowIdRule = append(escrowIdRule, common.Hash{})
			continue
		}
		fieldVal, err := bindings.PrepareTopicArg(evt.Inputs[0], v.EscrowId)
		if err != nil {
			return nil, err
		}
		escrowIdRule = append(escrowIdRule, fieldVal)
	}

	rawTopics, err := abi.MakeTopics(
		escrowIdRule,
	)
	if err != nil {
		return nil, err
	}

	return bindings.PrepareTopics(rawTopics, evt.ID.Bytes()), nil
}

// DecodeEscrowCompleted decodes a log into a EscrowCompleted struct.
func (c *Codec) DecodeEscrowCompleted(log *evm.Log) (*EscrowCompletedDecoded, error) {
	event := new(EscrowCompletedDecoded)
	if err := c.abi.UnpackIntoInterface(event, "EscrowCompleted", log.Data); err != nil {
		return nil, err
	}
	var indexed abi.Arguments
	for _, arg := range c.abi.Events["EscrowCompleted"].Inputs {
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

func (c *Codec) EscrowConfirmedLogHash() []byte {
	return c.abi.Events["EscrowConfirmed"].ID.Bytes()
}

func (c *Codec) EncodeEscrowConfirmedTopics(
	evt abi.Event,
	values []EscrowConfirmedTopics,
) ([]*evm.TopicValues, error) {
	var escrowIdRule []interface{}
	for _, v := range values {
		if reflect.ValueOf(v.EscrowId).IsZero() {
			escrowIdRule = append(escrowIdRule, common.Hash{})
			continue
		}
		fieldVal, err := bindings.PrepareTopicArg(evt.Inputs[0], v.EscrowId)
		if err != nil {
			return nil, err
		}
		escrowIdRule = append(escrowIdRule, fieldVal)
	}
	var confirmerRule []interface{}
	for _, v := range values {
		if reflect.ValueOf(v.Confirmer).IsZero() {
			confirmerRule = append(confirmerRule, common.Hash{})
			continue
		}
		fieldVal, err := bindings.PrepareTopicArg(evt.Inputs[1], v.Confirmer)
		if err != nil {
			return nil, err
		}
		confirmerRule = append(confirmerRule, fieldVal)
	}

	rawTopics, err := abi.MakeTopics(
		escrowIdRule,
		confirmerRule,
	)
	if err != nil {
		return nil, err
	}

	return bindings.PrepareTopics(rawTopics, evt.ID.Bytes()), nil
}

// DecodeEscrowConfirmed decodes a log into a EscrowConfirmed struct.
func (c *Codec) DecodeEscrowConfirmed(log *evm.Log) (*EscrowConfirmedDecoded, error) {
	event := new(EscrowConfirmedDecoded)
	if err := c.abi.UnpackIntoInterface(event, "EscrowConfirmed", log.Data); err != nil {
		return nil, err
	}
	var indexed abi.Arguments
	for _, arg := range c.abi.Events["EscrowConfirmed"].Inputs {
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

func (c *Codec) EscrowCreatedLogHash() []byte {
	return c.abi.Events["EscrowCreated"].ID.Bytes()
}

func (c *Codec) EncodeEscrowCreatedTopics(
	evt abi.Event,
	values []EscrowCreatedTopics,
) ([]*evm.TopicValues, error) {
	var escrowIdRule []interface{}
	for _, v := range values {
		if reflect.ValueOf(v.EscrowId).IsZero() {
			escrowIdRule = append(escrowIdRule, common.Hash{})
			continue
		}
		fieldVal, err := bindings.PrepareTopicArg(evt.Inputs[0], v.EscrowId)
		if err != nil {
			return nil, err
		}
		escrowIdRule = append(escrowIdRule, fieldVal)
	}
	var sellerRule []interface{}
	for _, v := range values {
		if reflect.ValueOf(v.Seller).IsZero() {
			sellerRule = append(sellerRule, common.Hash{})
			continue
		}
		fieldVal, err := bindings.PrepareTopicArg(evt.Inputs[1], v.Seller)
		if err != nil {
			return nil, err
		}
		sellerRule = append(sellerRule, fieldVal)
	}
	var buyerRule []interface{}
	for _, v := range values {
		if reflect.ValueOf(v.Buyer).IsZero() {
			buyerRule = append(buyerRule, common.Hash{})
			continue
		}
		fieldVal, err := bindings.PrepareTopicArg(evt.Inputs[2], v.Buyer)
		if err != nil {
			return nil, err
		}
		buyerRule = append(buyerRule, fieldVal)
	}

	rawTopics, err := abi.MakeTopics(
		escrowIdRule,
		sellerRule,
		buyerRule,
	)
	if err != nil {
		return nil, err
	}

	return bindings.PrepareTopics(rawTopics, evt.ID.Bytes()), nil
}

// DecodeEscrowCreated decodes a log into a EscrowCreated struct.
func (c *Codec) DecodeEscrowCreated(log *evm.Log) (*EscrowCreatedDecoded, error) {
	event := new(EscrowCreatedDecoded)
	if err := c.abi.UnpackIntoInterface(event, "EscrowCreated", log.Data); err != nil {
		return nil, err
	}
	var indexed abi.Arguments
	for _, arg := range c.abi.Events["EscrowCreated"].Inputs {
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

func (c *Codec) EscrowDisputedLogHash() []byte {
	return c.abi.Events["EscrowDisputed"].ID.Bytes()
}

func (c *Codec) EncodeEscrowDisputedTopics(
	evt abi.Event,
	values []EscrowDisputedTopics,
) ([]*evm.TopicValues, error) {
	var escrowIdRule []interface{}
	for _, v := range values {
		if reflect.ValueOf(v.EscrowId).IsZero() {
			escrowIdRule = append(escrowIdRule, common.Hash{})
			continue
		}
		fieldVal, err := bindings.PrepareTopicArg(evt.Inputs[0], v.EscrowId)
		if err != nil {
			return nil, err
		}
		escrowIdRule = append(escrowIdRule, fieldVal)
	}
	var disputerRule []interface{}
	for _, v := range values {
		if reflect.ValueOf(v.Disputer).IsZero() {
			disputerRule = append(disputerRule, common.Hash{})
			continue
		}
		fieldVal, err := bindings.PrepareTopicArg(evt.Inputs[1], v.Disputer)
		if err != nil {
			return nil, err
		}
		disputerRule = append(disputerRule, fieldVal)
	}

	rawTopics, err := abi.MakeTopics(
		escrowIdRule,
		disputerRule,
	)
	if err != nil {
		return nil, err
	}

	return bindings.PrepareTopics(rawTopics, evt.ID.Bytes()), nil
}

// DecodeEscrowDisputed decodes a log into a EscrowDisputed struct.
func (c *Codec) DecodeEscrowDisputed(log *evm.Log) (*EscrowDisputedDecoded, error) {
	event := new(EscrowDisputedDecoded)
	if err := c.abi.UnpackIntoInterface(event, "EscrowDisputed", log.Data); err != nil {
		return nil, err
	}
	var indexed abi.Arguments
	for _, arg := range c.abi.Events["EscrowDisputed"].Inputs {
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

func (c *Codec) EscrowFundedLogHash() []byte {
	return c.abi.Events["EscrowFunded"].ID.Bytes()
}

func (c *Codec) EncodeEscrowFundedTopics(
	evt abi.Event,
	values []EscrowFundedTopics,
) ([]*evm.TopicValues, error) {
	var escrowIdRule []interface{}
	for _, v := range values {
		if reflect.ValueOf(v.EscrowId).IsZero() {
			escrowIdRule = append(escrowIdRule, common.Hash{})
			continue
		}
		fieldVal, err := bindings.PrepareTopicArg(evt.Inputs[0], v.EscrowId)
		if err != nil {
			return nil, err
		}
		escrowIdRule = append(escrowIdRule, fieldVal)
	}
	var funderRule []interface{}
	for _, v := range values {
		if reflect.ValueOf(v.Funder).IsZero() {
			funderRule = append(funderRule, common.Hash{})
			continue
		}
		fieldVal, err := bindings.PrepareTopicArg(evt.Inputs[1], v.Funder)
		if err != nil {
			return nil, err
		}
		funderRule = append(funderRule, fieldVal)
	}

	rawTopics, err := abi.MakeTopics(
		escrowIdRule,
		funderRule,
	)
	if err != nil {
		return nil, err
	}

	return bindings.PrepareTopics(rawTopics, evt.ID.Bytes()), nil
}

// DecodeEscrowFunded decodes a log into a EscrowFunded struct.
func (c *Codec) DecodeEscrowFunded(log *evm.Log) (*EscrowFundedDecoded, error) {
	event := new(EscrowFundedDecoded)
	if err := c.abi.UnpackIntoInterface(event, "EscrowFunded", log.Data); err != nil {
		return nil, err
	}
	var indexed abi.Arguments
	for _, arg := range c.abi.Events["EscrowFunded"].Inputs {
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

func (c CommodityEscrow) WriteReport(
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

func (c *CommodityEscrow) UnpackError(data []byte) (any, error) {
	switch common.Bytes2Hex(data[:4]) {
	default:
		return nil, errors.New("unknown error selector")
	}
}

// EscrowCancelledTrigger wraps the raw log trigger and provides decoded EscrowCancelledDecoded data
type EscrowCancelledTrigger struct {
	cre.Trigger[*evm.Log, *evm.Log]                  // Embed the raw trigger
	contract                        *CommodityEscrow // Keep reference for decoding
}

// Adapt method that decodes the log into EscrowCancelled data
func (t *EscrowCancelledTrigger) Adapt(l *evm.Log) (*bindings.DecodedLog[EscrowCancelledDecoded], error) {
	// Decode the log using the contract's codec
	decoded, err := t.contract.Codec.DecodeEscrowCancelled(l)
	if err != nil {
		return nil, fmt.Errorf("failed to decode EscrowCancelled log: %w", err)
	}

	return &bindings.DecodedLog[EscrowCancelledDecoded]{
		Log:  l,        // Original log
		Data: *decoded, // Decoded data
	}, nil
}

func (c *CommodityEscrow) LogTriggerEscrowCancelledLog(chainSelector uint64, confidence evm.ConfidenceLevel, filters []EscrowCancelledTopics) (cre.Trigger[*evm.Log, *bindings.DecodedLog[EscrowCancelledDecoded]], error) {
	event := c.ABI.Events["EscrowCancelled"]
	topics, err := c.Codec.EncodeEscrowCancelledTopics(event, filters)
	if err != nil {
		return nil, fmt.Errorf("failed to encode topics for EscrowCancelled: %w", err)
	}

	rawTrigger := evm.LogTrigger(chainSelector, &evm.FilterLogTriggerRequest{
		Addresses:  [][]byte{c.Address.Bytes()},
		Topics:     topics,
		Confidence: confidence,
	})

	return &EscrowCancelledTrigger{
		Trigger:  rawTrigger,
		contract: c,
	}, nil
}

func (c *CommodityEscrow) FilterLogsEscrowCancelled(runtime cre.Runtime, options *bindings.FilterOptions) (cre.Promise[*evm.FilterLogsReply], error) {
	if options == nil {
		return nil, errors.New("FilterLogs options are required.")
	}
	return c.client.FilterLogs(runtime, &evm.FilterLogsRequest{
		FilterQuery: &evm.FilterQuery{
			Addresses: [][]byte{c.Address.Bytes()},
			Topics: []*evm.Topics{
				{Topic: [][]byte{c.Codec.EscrowCancelledLogHash()}},
			},
			BlockHash: options.BlockHash,
			FromBlock: pb.NewBigIntFromInt(options.FromBlock),
			ToBlock:   pb.NewBigIntFromInt(options.ToBlock),
		},
	}), nil
}

// EscrowCompletedTrigger wraps the raw log trigger and provides decoded EscrowCompletedDecoded data
type EscrowCompletedTrigger struct {
	cre.Trigger[*evm.Log, *evm.Log]                  // Embed the raw trigger
	contract                        *CommodityEscrow // Keep reference for decoding
}

// Adapt method that decodes the log into EscrowCompleted data
func (t *EscrowCompletedTrigger) Adapt(l *evm.Log) (*bindings.DecodedLog[EscrowCompletedDecoded], error) {
	// Decode the log using the contract's codec
	decoded, err := t.contract.Codec.DecodeEscrowCompleted(l)
	if err != nil {
		return nil, fmt.Errorf("failed to decode EscrowCompleted log: %w", err)
	}

	return &bindings.DecodedLog[EscrowCompletedDecoded]{
		Log:  l,        // Original log
		Data: *decoded, // Decoded data
	}, nil
}

func (c *CommodityEscrow) LogTriggerEscrowCompletedLog(chainSelector uint64, confidence evm.ConfidenceLevel, filters []EscrowCompletedTopics) (cre.Trigger[*evm.Log, *bindings.DecodedLog[EscrowCompletedDecoded]], error) {
	event := c.ABI.Events["EscrowCompleted"]
	topics, err := c.Codec.EncodeEscrowCompletedTopics(event, filters)
	if err != nil {
		return nil, fmt.Errorf("failed to encode topics for EscrowCompleted: %w", err)
	}

	rawTrigger := evm.LogTrigger(chainSelector, &evm.FilterLogTriggerRequest{
		Addresses:  [][]byte{c.Address.Bytes()},
		Topics:     topics,
		Confidence: confidence,
	})

	return &EscrowCompletedTrigger{
		Trigger:  rawTrigger,
		contract: c,
	}, nil
}

func (c *CommodityEscrow) FilterLogsEscrowCompleted(runtime cre.Runtime, options *bindings.FilterOptions) (cre.Promise[*evm.FilterLogsReply], error) {
	if options == nil {
		return nil, errors.New("FilterLogs options are required.")
	}
	return c.client.FilterLogs(runtime, &evm.FilterLogsRequest{
		FilterQuery: &evm.FilterQuery{
			Addresses: [][]byte{c.Address.Bytes()},
			Topics: []*evm.Topics{
				{Topic: [][]byte{c.Codec.EscrowCompletedLogHash()}},
			},
			BlockHash: options.BlockHash,
			FromBlock: pb.NewBigIntFromInt(options.FromBlock),
			ToBlock:   pb.NewBigIntFromInt(options.ToBlock),
		},
	}), nil
}

// EscrowConfirmedTrigger wraps the raw log trigger and provides decoded EscrowConfirmedDecoded data
type EscrowConfirmedTrigger struct {
	cre.Trigger[*evm.Log, *evm.Log]                  // Embed the raw trigger
	contract                        *CommodityEscrow // Keep reference for decoding
}

// Adapt method that decodes the log into EscrowConfirmed data
func (t *EscrowConfirmedTrigger) Adapt(l *evm.Log) (*bindings.DecodedLog[EscrowConfirmedDecoded], error) {
	// Decode the log using the contract's codec
	decoded, err := t.contract.Codec.DecodeEscrowConfirmed(l)
	if err != nil {
		return nil, fmt.Errorf("failed to decode EscrowConfirmed log: %w", err)
	}

	return &bindings.DecodedLog[EscrowConfirmedDecoded]{
		Log:  l,        // Original log
		Data: *decoded, // Decoded data
	}, nil
}

func (c *CommodityEscrow) LogTriggerEscrowConfirmedLog(chainSelector uint64, confidence evm.ConfidenceLevel, filters []EscrowConfirmedTopics) (cre.Trigger[*evm.Log, *bindings.DecodedLog[EscrowConfirmedDecoded]], error) {
	event := c.ABI.Events["EscrowConfirmed"]
	topics, err := c.Codec.EncodeEscrowConfirmedTopics(event, filters)
	if err != nil {
		return nil, fmt.Errorf("failed to encode topics for EscrowConfirmed: %w", err)
	}

	rawTrigger := evm.LogTrigger(chainSelector, &evm.FilterLogTriggerRequest{
		Addresses:  [][]byte{c.Address.Bytes()},
		Topics:     topics,
		Confidence: confidence,
	})

	return &EscrowConfirmedTrigger{
		Trigger:  rawTrigger,
		contract: c,
	}, nil
}

func (c *CommodityEscrow) FilterLogsEscrowConfirmed(runtime cre.Runtime, options *bindings.FilterOptions) (cre.Promise[*evm.FilterLogsReply], error) {
	if options == nil {
		return nil, errors.New("FilterLogs options are required.")
	}
	return c.client.FilterLogs(runtime, &evm.FilterLogsRequest{
		FilterQuery: &evm.FilterQuery{
			Addresses: [][]byte{c.Address.Bytes()},
			Topics: []*evm.Topics{
				{Topic: [][]byte{c.Codec.EscrowConfirmedLogHash()}},
			},
			BlockHash: options.BlockHash,
			FromBlock: pb.NewBigIntFromInt(options.FromBlock),
			ToBlock:   pb.NewBigIntFromInt(options.ToBlock),
		},
	}), nil
}

// EscrowCreatedTrigger wraps the raw log trigger and provides decoded EscrowCreatedDecoded data
type EscrowCreatedTrigger struct {
	cre.Trigger[*evm.Log, *evm.Log]                  // Embed the raw trigger
	contract                        *CommodityEscrow // Keep reference for decoding
}

// Adapt method that decodes the log into EscrowCreated data
func (t *EscrowCreatedTrigger) Adapt(l *evm.Log) (*bindings.DecodedLog[EscrowCreatedDecoded], error) {
	// Decode the log using the contract's codec
	decoded, err := t.contract.Codec.DecodeEscrowCreated(l)
	if err != nil {
		return nil, fmt.Errorf("failed to decode EscrowCreated log: %w", err)
	}

	return &bindings.DecodedLog[EscrowCreatedDecoded]{
		Log:  l,        // Original log
		Data: *decoded, // Decoded data
	}, nil
}

func (c *CommodityEscrow) LogTriggerEscrowCreatedLog(chainSelector uint64, confidence evm.ConfidenceLevel, filters []EscrowCreatedTopics) (cre.Trigger[*evm.Log, *bindings.DecodedLog[EscrowCreatedDecoded]], error) {
	event := c.ABI.Events["EscrowCreated"]
	topics, err := c.Codec.EncodeEscrowCreatedTopics(event, filters)
	if err != nil {
		return nil, fmt.Errorf("failed to encode topics for EscrowCreated: %w", err)
	}

	rawTrigger := evm.LogTrigger(chainSelector, &evm.FilterLogTriggerRequest{
		Addresses:  [][]byte{c.Address.Bytes()},
		Topics:     topics,
		Confidence: confidence,
	})

	return &EscrowCreatedTrigger{
		Trigger:  rawTrigger,
		contract: c,
	}, nil
}

func (c *CommodityEscrow) FilterLogsEscrowCreated(runtime cre.Runtime, options *bindings.FilterOptions) (cre.Promise[*evm.FilterLogsReply], error) {
	if options == nil {
		return nil, errors.New("FilterLogs options are required.")
	}
	return c.client.FilterLogs(runtime, &evm.FilterLogsRequest{
		FilterQuery: &evm.FilterQuery{
			Addresses: [][]byte{c.Address.Bytes()},
			Topics: []*evm.Topics{
				{Topic: [][]byte{c.Codec.EscrowCreatedLogHash()}},
			},
			BlockHash: options.BlockHash,
			FromBlock: pb.NewBigIntFromInt(options.FromBlock),
			ToBlock:   pb.NewBigIntFromInt(options.ToBlock),
		},
	}), nil
}

// EscrowDisputedTrigger wraps the raw log trigger and provides decoded EscrowDisputedDecoded data
type EscrowDisputedTrigger struct {
	cre.Trigger[*evm.Log, *evm.Log]                  // Embed the raw trigger
	contract                        *CommodityEscrow // Keep reference for decoding
}

// Adapt method that decodes the log into EscrowDisputed data
func (t *EscrowDisputedTrigger) Adapt(l *evm.Log) (*bindings.DecodedLog[EscrowDisputedDecoded], error) {
	// Decode the log using the contract's codec
	decoded, err := t.contract.Codec.DecodeEscrowDisputed(l)
	if err != nil {
		return nil, fmt.Errorf("failed to decode EscrowDisputed log: %w", err)
	}

	return &bindings.DecodedLog[EscrowDisputedDecoded]{
		Log:  l,        // Original log
		Data: *decoded, // Decoded data
	}, nil
}

func (c *CommodityEscrow) LogTriggerEscrowDisputedLog(chainSelector uint64, confidence evm.ConfidenceLevel, filters []EscrowDisputedTopics) (cre.Trigger[*evm.Log, *bindings.DecodedLog[EscrowDisputedDecoded]], error) {
	event := c.ABI.Events["EscrowDisputed"]
	topics, err := c.Codec.EncodeEscrowDisputedTopics(event, filters)
	if err != nil {
		return nil, fmt.Errorf("failed to encode topics for EscrowDisputed: %w", err)
	}

	rawTrigger := evm.LogTrigger(chainSelector, &evm.FilterLogTriggerRequest{
		Addresses:  [][]byte{c.Address.Bytes()},
		Topics:     topics,
		Confidence: confidence,
	})

	return &EscrowDisputedTrigger{
		Trigger:  rawTrigger,
		contract: c,
	}, nil
}

func (c *CommodityEscrow) FilterLogsEscrowDisputed(runtime cre.Runtime, options *bindings.FilterOptions) (cre.Promise[*evm.FilterLogsReply], error) {
	if options == nil {
		return nil, errors.New("FilterLogs options are required.")
	}
	return c.client.FilterLogs(runtime, &evm.FilterLogsRequest{
		FilterQuery: &evm.FilterQuery{
			Addresses: [][]byte{c.Address.Bytes()},
			Topics: []*evm.Topics{
				{Topic: [][]byte{c.Codec.EscrowDisputedLogHash()}},
			},
			BlockHash: options.BlockHash,
			FromBlock: pb.NewBigIntFromInt(options.FromBlock),
			ToBlock:   pb.NewBigIntFromInt(options.ToBlock),
		},
	}), nil
}

// EscrowFundedTrigger wraps the raw log trigger and provides decoded EscrowFundedDecoded data
type EscrowFundedTrigger struct {
	cre.Trigger[*evm.Log, *evm.Log]                  // Embed the raw trigger
	contract                        *CommodityEscrow // Keep reference for decoding
}

// Adapt method that decodes the log into EscrowFunded data
func (t *EscrowFundedTrigger) Adapt(l *evm.Log) (*bindings.DecodedLog[EscrowFundedDecoded], error) {
	// Decode the log using the contract's codec
	decoded, err := t.contract.Codec.DecodeEscrowFunded(l)
	if err != nil {
		return nil, fmt.Errorf("failed to decode EscrowFunded log: %w", err)
	}

	return &bindings.DecodedLog[EscrowFundedDecoded]{
		Log:  l,        // Original log
		Data: *decoded, // Decoded data
	}, nil
}

func (c *CommodityEscrow) LogTriggerEscrowFundedLog(chainSelector uint64, confidence evm.ConfidenceLevel, filters []EscrowFundedTopics) (cre.Trigger[*evm.Log, *bindings.DecodedLog[EscrowFundedDecoded]], error) {
	event := c.ABI.Events["EscrowFunded"]
	topics, err := c.Codec.EncodeEscrowFundedTopics(event, filters)
	if err != nil {
		return nil, fmt.Errorf("failed to encode topics for EscrowFunded: %w", err)
	}

	rawTrigger := evm.LogTrigger(chainSelector, &evm.FilterLogTriggerRequest{
		Addresses:  [][]byte{c.Address.Bytes()},
		Topics:     topics,
		Confidence: confidence,
	})

	return &EscrowFundedTrigger{
		Trigger:  rawTrigger,
		contract: c,
	}, nil
}

func (c *CommodityEscrow) FilterLogsEscrowFunded(runtime cre.Runtime, options *bindings.FilterOptions) (cre.Promise[*evm.FilterLogsReply], error) {
	if options == nil {
		return nil, errors.New("FilterLogs options are required.")
	}
	return c.client.FilterLogs(runtime, &evm.FilterLogsRequest{
		FilterQuery: &evm.FilterQuery{
			Addresses: [][]byte{c.Address.Bytes()},
			Topics: []*evm.Topics{
				{Topic: [][]byte{c.Codec.EscrowFundedLogHash()}},
			},
			BlockHash: options.BlockHash,
			FromBlock: pb.NewBigIntFromInt(options.FromBlock),
			ToBlock:   pb.NewBigIntFromInt(options.ToBlock),
		},
	}), nil
}
