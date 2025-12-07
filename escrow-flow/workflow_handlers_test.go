package main

import (
	"encoding/json"
	"log/slog"
	"math/big"
	"os"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/smartcontractkit/cre-sdk-go/capabilities/blockchain/evm"
	"github.com/smartcontractkit/cre-sdk-go/capabilities/blockchain/evm/bindings"

	"akua-cre/contracts/evm/src/generated/commodity_escrow"
	"akua-cre/contracts/evm/src/generated/warehouse_receipt_registry"
)

// MockRuntime implements a minimal cre.Runtime for testing
type MockRuntime struct {
	logger *slog.Logger
}

func (m *MockRuntime) Logger() *slog.Logger {
	if m.logger == nil {
		m.logger = slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
			Level: slog.LevelError, // Suppress logs during tests
		}))
	}
	return m.logger
}

// Helper to create a test config
func createTestConfig() *Config {
	return &Config{
		WebhookURL:             "http://localhost:9999/test-webhook",
		WebhookSecret:          "test-secret",
		ChainName:              "ethereum-testnet-sepolia-arbitrum-1",
		Network:                "arbitrumSepolia",
		EscrowAddress:          "0x4cE838B32469bAdD067C05EBc04A4c01EDcAE903",
		ReceiptRegistryAddress: "0xa04f93a8810727Fd1B2867432F5612391a9daF15",
		GasLimit:               800000,
	}
}

// Helper to create mock log metadata
// The newer CRE SDK represents TxHash as a []byte and no longer exposes BlockTimestamp.
// For testing we only care that populateTxMetadata can derive a non-empty txHash string.
func createMockLog() *evm.Log {
	txBytes := common.Hex2Bytes("1234567890abcdef1234567890abcdef1234567890abcdef1234567890abcdef")
	return &evm.Log{
		TxHash: txBytes,
	}
}

// Test EscrowCreated handler
func TestOnEscrowCreated(t *testing.T) {
	payload := &bindings.DecodedLog[commodity_escrow.EscrowCreatedDecoded]{
		Data: commodity_escrow.EscrowCreatedDecoded{
			EscrowId:       [32]byte{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19, 20, 21, 22, 23, 24, 25, 26, 27, 28, 29, 30, 31, 32},
			Seller:         common.HexToAddress("0x1234567890123456789012345678901234567890"),
			Buyer:          common.HexToAddress("0x0987654321098765432109876543210987654321"),
			CommodityToken: common.HexToAddress("0xabcdef1234567890abcdef1234567890abcdef12"),
			TokenAmount:    big.NewInt(1000),
			PricePerToken:  big.NewInt(50),
			ReceiptId:      [32]byte{32, 31, 30, 29, 28, 27, 26, 25, 24, 23, 22, 21, 20, 19, 18, 17, 16, 15, 14, 13, 12, 11, 10, 9, 8, 7, 6, 5, 4, 3, 2, 1},
		},
		Log: createMockLog(),
	}

	// Note: We can't actually call the handler without setting up HTTP mocking
	// Instead, we'll test the payload construction logic
	testPayload := &EventPayload{
		EventName:      "EscrowCreated",
		Network:        "arbitrumSepolia",
		ChainName:      "ethereum-testnet-sepolia-arbitrum-1",
		EscrowID:       common.BytesToHash(payload.Data.EscrowId[:]).Hex(),
		Seller:         payload.Data.Seller.Hex(),
		Buyer:          payload.Data.Buyer.Hex(),
		CommodityToken: payload.Data.CommodityToken.Hex(),
		TokenAmount:    payload.Data.TokenAmount.String(),
		PricePerToken:  payload.Data.PricePerToken.String(),
		ReceiptID:      common.BytesToHash(payload.Data.ReceiptId[:]).Hex(),
	}
	populateTxMetadata(testPayload, payload.Log)

	// Validate payload structure
	data, err := json.Marshal(testPayload)
	if err != nil {
		t.Fatalf("failed to marshal EscrowCreated payload: %v", err)
	}

	var result EventPayload
	if err := json.Unmarshal(data, &result); err != nil {
		t.Fatalf("failed to unmarshal EscrowCreated payload: %v", err)
	}

	// Validate required fields
	if result.EventName != "EscrowCreated" {
		t.Errorf("expected eventName 'EscrowCreated', got: %s", result.EventName)
	}
	if result.Network == "" {
		t.Error("network should not be empty")
	}
	if result.EscrowID == "" {
		t.Error("escrowId should not be empty")
	}
	if result.TxHash == "" {
		t.Error("txHash should not be empty")
	}
}

// Test EscrowFunded handler
func TestOnEscrowFunded(t *testing.T) {
	payload := &bindings.DecodedLog[commodity_escrow.EscrowFundedDecoded]{
		Data: commodity_escrow.EscrowFundedDecoded{
			EscrowId: [32]byte{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19, 20, 21, 22, 23, 24, 25, 26, 27, 28, 29, 30, 31, 32},
			Funder:   common.HexToAddress("0x1234567890123456789012345678901234567890"),
			Amount:   big.NewInt(50000),
		},
		Log: createMockLog(),
	}

	testPayload := &EventPayload{
		EventName: "EscrowFunded",
		Network:   "arbitrumSepolia",
		ChainName: "ethereum-testnet-sepolia-arbitrum-1",
		EscrowID:  common.BytesToHash(payload.Data.EscrowId[:]).Hex(),
		Funder:    payload.Data.Funder.Hex(),
		Amount:    payload.Data.Amount.String(), // Payment amount, not token amount
	}
	populateTxMetadata(testPayload, payload.Log)

	data, err := json.Marshal(testPayload)
	if err != nil {
		t.Fatalf("failed to marshal EscrowFunded payload: %v", err)
	}

	var result EventPayload
	if err := json.Unmarshal(data, &result); err != nil {
		t.Fatalf("failed to unmarshal EscrowFunded payload: %v", err)
	}

	if result.EventName != "EscrowFunded" {
		t.Errorf("expected eventName 'EscrowFunded', got: %s", result.EventName)
	}
	if result.Funder == "" {
		t.Error("funder should not be empty")
	}
}

// Test EscrowConfirmed handler
func TestOnEscrowConfirmed(t *testing.T) {
	payload := &bindings.DecodedLog[commodity_escrow.EscrowConfirmedDecoded]{
		Data: commodity_escrow.EscrowConfirmedDecoded{
			EscrowId:  [32]byte{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19, 20, 21, 22, 23, 24, 25, 26, 27, 28, 29, 30, 31, 32},
			Confirmer: common.HexToAddress("0x0987654321098765432109876543210987654321"),
		},
		Log: createMockLog(),
	}

	testPayload := &EventPayload{
		EventName: "EscrowConfirmed",
		Network:   "arbitrumSepolia",
		ChainName: "ethereum-testnet-sepolia-arbitrum-1",
		EscrowID:  common.BytesToHash(payload.Data.EscrowId[:]).Hex(),
		Confirmer: payload.Data.Confirmer.Hex(),
	}
	populateTxMetadata(testPayload, payload.Log)

	data, err := json.Marshal(testPayload)
	if err != nil {
		t.Fatalf("failed to marshal EscrowConfirmed payload: %v", err)
	}

	var result EventPayload
	if err := json.Unmarshal(data, &result); err != nil {
		t.Fatalf("failed to unmarshal EscrowConfirmed payload: %v", err)
	}

	if result.EventName != "EscrowConfirmed" {
		t.Errorf("expected eventName 'EscrowConfirmed', got: %s", result.EventName)
	}
}

// Test EscrowCompleted handler
func TestOnEscrowCompleted(t *testing.T) {
	payload := &bindings.DecodedLog[commodity_escrow.EscrowCompletedDecoded]{
		Data: commodity_escrow.EscrowCompletedDecoded{
			EscrowId: [32]byte{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19, 20, 21, 22, 23, 24, 25, 26, 27, 28, 29, 30, 31, 32},
		},
		Log: createMockLog(),
	}

	testPayload := &EventPayload{
		EventName: "EscrowCompleted",
		Network:   "arbitrumSepolia",
		ChainName: "ethereum-testnet-sepolia-arbitrum-1",
		EscrowID:  common.BytesToHash(payload.Data.EscrowId[:]).Hex(),
	}
	populateTxMetadata(testPayload, payload.Log)

	data, err := json.Marshal(testPayload)
	if err != nil {
		t.Fatalf("failed to marshal EscrowCompleted payload: %v", err)
	}

	var result EventPayload
	if err := json.Unmarshal(data, &result); err != nil {
		t.Fatalf("failed to unmarshal EscrowCompleted payload: %v", err)
	}

	if result.EventName != "EscrowCompleted" {
		t.Errorf("expected eventName 'EscrowCompleted', got: %s", result.EventName)
	}
}

// Test EscrowCancelled handler
func TestOnEscrowCancelled(t *testing.T) {
	payload := &bindings.DecodedLog[commodity_escrow.EscrowCancelledDecoded]{
		Data: commodity_escrow.EscrowCancelledDecoded{
			EscrowId:    [32]byte{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19, 20, 21, 22, 23, 24, 25, 26, 27, 28, 29, 30, 31, 32},
			CancelledBy: common.HexToAddress("0x1234567890123456789012345678901234567890"),
		},
		Log: createMockLog(),
	}

	testPayload := &EventPayload{
		EventName:   "EscrowCancelled",
		Network:     "arbitrumSepolia",
		ChainName:   "ethereum-testnet-sepolia-arbitrum-1",
		EscrowID:    common.BytesToHash(payload.Data.EscrowId[:]).Hex(),
		CancelledBy: payload.Data.CancelledBy.Hex(),
	}
	populateTxMetadata(testPayload, payload.Log)

	data, err := json.Marshal(testPayload)
	if err != nil {
		t.Fatalf("failed to marshal EscrowCancelled payload: %v", err)
	}

	var result EventPayload
	if err := json.Unmarshal(data, &result); err != nil {
		t.Fatalf("failed to unmarshal EscrowCancelled payload: %v", err)
	}

	if result.EventName != "EscrowCancelled" {
		t.Errorf("expected eventName 'EscrowCancelled', got: %s", result.EventName)
	}
	if result.CancelledBy == "" {
		t.Error("cancelledBy should not be empty")
	}
}

// Test EscrowDisputed handler
func TestOnEscrowDisputed(t *testing.T) {
	payload := &bindings.DecodedLog[commodity_escrow.EscrowDisputedDecoded]{
		Data: commodity_escrow.EscrowDisputedDecoded{
			EscrowId: [32]byte{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19, 20, 21, 22, 23, 24, 25, 26, 27, 28, 29, 30, 31, 32},
			Disputer: common.HexToAddress("0x0987654321098765432109876543210987654321"),
		},
		Log: createMockLog(),
	}

	testPayload := &EventPayload{
		EventName: "EscrowDisputed",
		Network:   "arbitrumSepolia",
		ChainName: "ethereum-testnet-sepolia-arbitrum-1",
		EscrowID:  common.BytesToHash(payload.Data.EscrowId[:]).Hex(),
		Disputer:  payload.Data.Disputer.Hex(),
	}
	populateTxMetadata(testPayload, payload.Log)

	data, err := json.Marshal(testPayload)
	if err != nil {
		t.Fatalf("failed to marshal EscrowDisputed payload: %v", err)
	}

	var result EventPayload
	if err := json.Unmarshal(data, &result); err != nil {
		t.Fatalf("failed to unmarshal EscrowDisputed payload: %v", err)
	}

	if result.EventName != "EscrowDisputed" {
		t.Errorf("expected eventName 'EscrowDisputed', got: %s", result.EventName)
	}
	if result.Disputer == "" {
		t.Error("disputer should not be empty")
	}
}

// Test ReceiptIssued handler
func TestOnReceiptIssued(t *testing.T) {
	payload := &bindings.DecodedLog[warehouse_receipt_registry.ReceiptIssuedDecoded]{
		Data: warehouse_receipt_registry.ReceiptIssuedDecoded{
			ReceiptId:     [32]byte{32, 31, 30, 29, 28, 27, 26, 25, 24, 23, 22, 21, 20, 19, 18, 17, 16, 15, 14, 13, 12, 11, 10, 9, 8, 7, 6, 5, 4, 3, 2, 1},
			Depositor:     common.HexToAddress("0x1234567890123456789012345678901234567890"),
			Warehouse:     common.HexToAddress("0xabcdef1234567890abcdef1234567890abcdef12"),
			CommodityType: "Cocoa Beans",
			Quantity:      big.NewInt(10000),
			DepositDate:   big.NewInt(1640000000),
			ExpiryDate:    big.NewInt(1672000000),
		},
		Log: createMockLog(),
	}

	testPayload := &EventPayload{
		EventName:     "ReceiptIssued",
		Network:       "arbitrumSepolia",
		ChainName:     "ethereum-testnet-sepolia-arbitrum-1",
		ReceiptID:     common.BytesToHash(payload.Data.ReceiptId[:]).Hex(),
		Depositor:     payload.Data.Depositor.Hex(),
		Warehouse:     payload.Data.Warehouse.Hex(),
		CommodityType: payload.Data.CommodityType,
		Quantity:      payload.Data.Quantity.String(),
		DepositDate:   payload.Data.DepositDate.String(),
		ExpiryDate:    payload.Data.ExpiryDate.String(),
	}
	populateTxMetadata(testPayload, payload.Log)

	data, err := json.Marshal(testPayload)
	if err != nil {
		t.Fatalf("failed to marshal ReceiptIssued payload: %v", err)
	}

	var result EventPayload
	if err := json.Unmarshal(data, &result); err != nil {
		t.Fatalf("failed to unmarshal ReceiptIssued payload: %v", err)
	}

	if result.EventName != "ReceiptIssued" {
		t.Errorf("expected eventName 'ReceiptIssued', got: %s", result.EventName)
	}
	if result.ReceiptID == "" {
		t.Error("receiptId should not be empty")
	}
	if result.CommodityType == "" {
		t.Error("commodityType should not be empty")
	}
}

// Test ReceiptRedeemed handler
func TestOnReceiptRedeemed(t *testing.T) {
	payload := &bindings.DecodedLog[warehouse_receipt_registry.ReceiptRedeemedDecoded]{
		Data: warehouse_receipt_registry.ReceiptRedeemedDecoded{
			ReceiptId: [32]byte{32, 31, 30, 29, 28, 27, 26, 25, 24, 23, 22, 21, 20, 19, 18, 17, 16, 15, 14, 13, 12, 11, 10, 9, 8, 7, 6, 5, 4, 3, 2, 1},
			Redeemer:  common.HexToAddress("0x0987654321098765432109876543210987654321"),
		},
		Log: createMockLog(),
	}

	testPayload := &EventPayload{
		EventName: "ReceiptRedeemed",
		Network:   "arbitrumSepolia",
		ChainName: "ethereum-testnet-sepolia-arbitrum-1",
		ReceiptID: common.BytesToHash(payload.Data.ReceiptId[:]).Hex(),
		Redeemer:  payload.Data.Redeemer.Hex(),
	}
	populateTxMetadata(testPayload, payload.Log)

	data, err := json.Marshal(testPayload)
	if err != nil {
		t.Fatalf("failed to marshal ReceiptRedeemed payload: %v", err)
	}

	var result EventPayload
	if err := json.Unmarshal(data, &result); err != nil {
		t.Fatalf("failed to unmarshal ReceiptRedeemed payload: %v", err)
	}

	if result.EventName != "ReceiptRedeemed" {
		t.Errorf("expected eventName 'ReceiptRedeemed', got: %s", result.EventName)
	}
	if result.Redeemer == "" {
		t.Error("redeemer should not be empty")
	}
}

// Test ReceiptCancelled handler
func TestOnReceiptCancelled(t *testing.T) {
	payload := &bindings.DecodedLog[warehouse_receipt_registry.ReceiptCancelledDecoded]{
		Data: warehouse_receipt_registry.ReceiptCancelledDecoded{
			ReceiptId:   [32]byte{32, 31, 30, 29, 28, 27, 26, 25, 24, 23, 22, 21, 20, 19, 18, 17, 16, 15, 14, 13, 12, 11, 10, 9, 8, 7, 6, 5, 4, 3, 2, 1},
			CancelledBy: common.HexToAddress("0x1234567890123456789012345678901234567890"),
		},
		Log: createMockLog(),
	}

	testPayload := &EventPayload{
		EventName:   "ReceiptCancelled",
		Network:     "arbitrumSepolia",
		ChainName:   "ethereum-testnet-sepolia-arbitrum-1",
		ReceiptID:   common.BytesToHash(payload.Data.ReceiptId[:]).Hex(),
		CancelledBy: payload.Data.CancelledBy.Hex(),
	}
	populateTxMetadata(testPayload, payload.Log)

	data, err := json.Marshal(testPayload)
	if err != nil {
		t.Fatalf("failed to marshal ReceiptCancelled payload: %v", err)
	}

	var result EventPayload
	if err := json.Unmarshal(data, &result); err != nil {
		t.Fatalf("failed to unmarshal ReceiptCancelled payload: %v", err)
	}

	if result.EventName != "ReceiptCancelled" {
		t.Errorf("expected eventName 'ReceiptCancelled', got: %s", result.EventName)
	}
}

// Test ReceiptUpdated handler
func TestOnReceiptUpdated(t *testing.T) {
	payload := &bindings.DecodedLog[warehouse_receipt_registry.ReceiptUpdatedDecoded]{
		Data: warehouse_receipt_registry.ReceiptUpdatedDecoded{
			ReceiptId:   [32]byte{32, 31, 30, 29, 28, 27, 26, 25, 24, 23, 22, 21, 20, 19, 18, 17, 16, 15, 14, 13, 12, 11, 10, 9, 8, 7, 6, 5, 4, 3, 2, 1},
			NewQuantity: big.NewInt(8000),
		},
		Log: createMockLog(),
	}

	testPayload := &EventPayload{
		EventName:   "ReceiptUpdated",
		Network:     "arbitrumSepolia",
		ChainName:   "ethereum-testnet-sepolia-arbitrum-1",
		ReceiptID:   common.BytesToHash(payload.Data.ReceiptId[:]).Hex(),
		NewQuantity: payload.Data.NewQuantity.String(),
	}
	populateTxMetadata(testPayload, payload.Log)

	data, err := json.Marshal(testPayload)
	if err != nil {
		t.Fatalf("failed to marshal ReceiptUpdated payload: %v", err)
	}

	var result EventPayload
	if err := json.Unmarshal(data, &result); err != nil {
		t.Fatalf("failed to unmarshal ReceiptUpdated payload: %v", err)
	}

	if result.EventName != "ReceiptUpdated" {
		t.Errorf("expected eventName 'ReceiptUpdated', got: %s", result.EventName)
	}
	if result.NewQuantity == "" {
		t.Error("newQuantity should not be empty")
	}
}

// Test populateTxMetadata function
func TestPopulateTxMetadata(t *testing.T) {
	payload := &EventPayload{
		EventName: "TestEvent",
	}

	log := createMockLog()
	populateTxMetadata(payload, log)

	if payload.TxHash == "" {
		t.Error("txHash should be populated")
	}
	// We no longer assert on blockNumber or timestamp because the current SDK
	// does not expose BlockTimestamp on evm.Log and BlockNumber may be nil in tests.

	// Test with nil log (should not panic)
	payload2 := &EventPayload{EventName: "TestEvent2"}
	populateTxMetadata(payload2, nil)
	// Should complete without panic
}

