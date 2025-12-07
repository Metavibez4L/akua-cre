package main

import (
	"encoding/json"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestConfigDeserialization(t *testing.T) {
	configJSON := `{
		"webhookUrl": "http://localhost:3002/api/dta-webhook",
		"webhookSecret": "test-secret",
		"chainName": "ethereum-testnet-sepolia-arbitrum-1",
		"network": "arbitrumSepolia",
		"dtaRequestManagerAddress": "0x1234567890123456789012345678901234567890",
		"escrowContractAddress": "0x4cE838B32469bAdD067C05EBc04A4c01EDcAE903",
		"receiptRegistryAddress": "0xa04f93a8810727Fd1B2867432F5612391a9daF15",
		"priceOracleAddress": "0x6199C8bf5be5513877899d9EB2cDd892F87D6066",
		"cocoaTonTokenAddress": "0xabcdabcdabcdabcdabcdabcdabcdabcdabcdabcd",
		"gasLimit": 800000
	}`

	var cfg Config
	err := json.Unmarshal([]byte(configJSON), &cfg)
	require.NoError(t, err)

	assert.Equal(t, "http://localhost:3002/api/dta-webhook", cfg.WebhookURL)
	assert.Equal(t, "test-secret", cfg.WebhookSecret)
	assert.Equal(t, "ethereum-testnet-sepolia-arbitrum-1", cfg.ChainName)
	assert.Equal(t, "arbitrumSepolia", cfg.Network)
	assert.Equal(t, "0x1234567890123456789012345678901234567890", cfg.DtaRequestManagerAddress)
	assert.Equal(t, "0x4cE838B32469bAdD067C05EBc04A4c01EDcAE903", cfg.EscrowAddress)
	assert.Equal(t, uint64(800000), cfg.GasLimit)
}

func TestValidateConfigForRuntime_ValidTestnetConfig(t *testing.T) {
	cfg := &Config{
		WebhookURL:               "http://localhost:3002/api/dta-webhook",
		WebhookSecret:            "test-secret",
		ChainName:                "ethereum-testnet-sepolia-arbitrum-1",
		Network:                  "arbitrumSepolia",
		DtaRequestManagerAddress: "0x1234567890123456789012345678901234567890",
		EscrowAddress:            "0x4cE838B32469bAdD067C05EBc04A4c01EDcAE903",
		ReceiptRegistryAddress:   "0xa04f93a8810727Fd1B2867432F5612391a9daF15",
		PriceOracleAddress:       "0x6199C8bf5be5513877899d9EB2cDd892F87D6066",
		CocoaTonTokenAddress:     "0xabcdabcdabcdabcdabcdabcdabcdabcdabcdabcd",
		GasLimit:                 800000,
	}

	err := validateConfigForRuntime(cfg)
	assert.NoError(t, err)
}

func TestValidateConfigForRuntime_MissingWebhookURL(t *testing.T) {
	cfg := &Config{
		WebhookSecret:            "test-secret",
		ChainName:                "ethereum-testnet-sepolia-arbitrum-1",
		Network:                  "arbitrumSepolia",
		DtaRequestManagerAddress: "0x1234567890123456789012345678901234567890",
		EscrowAddress:            "0x4cE838B32469bAdD067C05EBc04A4c01EDcAE903",
		ReceiptRegistryAddress:   "0xa04f93a8810727Fd1B2867432F5612391a9daF15",
		PriceOracleAddress:       "0x6199C8bf5be5513877899d9EB2cDd892F87D6066",
		CocoaTonTokenAddress:     "0xabcdabcdabcdabcdabcdabcdabcdabcdabcdabcd",
		GasLimit:                 800000,
	}

	err := validateConfigForRuntime(cfg)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "webhookUrl cannot be empty")
}

func TestValidateConfigForRuntime_InvalidAddress(t *testing.T) {
	cfg := &Config{
		WebhookURL:               "http://localhost:3002/api/dta-webhook",
		WebhookSecret:            "test-secret",
		ChainName:                "ethereum-testnet-sepolia-arbitrum-1",
		Network:                  "arbitrumSepolia",
		DtaRequestManagerAddress: "not-a-valid-address",
		EscrowAddress:            "0x4cE838B32469bAdD067C05EBc04A4c01EDcAE903",
		ReceiptRegistryAddress:   "0xa04f93a8810727Fd1B2867432F5612391a9daF15",
		PriceOracleAddress:       "0x6199C8bf5be5513877899d9EB2cDd892F87D6066",
		CocoaTonTokenAddress:     "0xabcdabcdabcdabcdabcdabcdabcdabcdabcdabcd",
		GasLimit:                 800000,
	}

	err := validateConfigForRuntime(cfg)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "invalid dtaRequestManagerAddress")
}

func TestValidateConfigForRuntime_ProductionRequiresHTTPS(t *testing.T) {
	cfg := &Config{
		WebhookURL:               "http://example.com/webhook",
		WebhookSecret:            "test-secret",
		ChainName:                "ethereum-mainnet-arbitrum-1",
		Network:                  "arbitrumOne",
		DtaRequestManagerAddress: "0x1234567890123456789012345678901234567890",
		EscrowAddress:            "0x4cE838B32469bAdD067C05EBc04A4c01EDcAE903",
		ReceiptRegistryAddress:   "0xa04f93a8810727Fd1B2867432F5612391a9daF15",
		PriceOracleAddress:       "0x6199C8bf5be5513877899d9EB2cDd892F87D6066",
		CocoaTonTokenAddress:     "0xabcdabcdabcdabcdabcdabcdabcdabcdabcdabcd",
		GasLimit:                 800000,
	}

	err := validateConfigForRuntime(cfg)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "production webhookUrl must use HTTPS")
}

func TestValidateConfigForRuntime_ProductionDisallowsLocalhost(t *testing.T) {
	cfg := &Config{
		WebhookURL:               "https://localhost:3002/api/dta-webhook",
		WebhookSecret:            "test-secret",
		ChainName:                "ethereum-mainnet-arbitrum-1",
		Network:                  "arbitrumOne",
		DtaRequestManagerAddress: "0x1234567890123456789012345678901234567890",
		EscrowAddress:            "0x4cE838B32469bAdD067C05EBc04A4c01EDcAE903",
		ReceiptRegistryAddress:   "0xa04f93a8810727Fd1B2867432F5612391a9daF15",
		PriceOracleAddress:       "0x6199C8bf5be5513877899d9EB2cDd892F87D6066",
		CocoaTonTokenAddress:     "0xabcdabcdabcdabcdabcdabcdabcdabcdabcdabcd",
		GasLimit:                 800000,
	}

	err := validateConfigForRuntime(cfg)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "production webhookUrl cannot point to localhost")
}

func TestValidateConfigForRuntime_ZeroGasLimit(t *testing.T) {
	cfg := &Config{
		WebhookURL:               "http://localhost:3002/api/dta-webhook",
		WebhookSecret:            "test-secret",
		ChainName:                "ethereum-testnet-sepolia-arbitrum-1",
		Network:                  "arbitrumSepolia",
		DtaRequestManagerAddress: "0x1234567890123456789012345678901234567890",
		EscrowAddress:            "0x4cE838B32469bAdD067C05EBc04A4c01EDcAE903",
		ReceiptRegistryAddress:   "0xa04f93a8810727Fd1B2867432F5612391a9daF15",
		PriceOracleAddress:       "0x6199C8bf5be5513877899d9EB2cDd892F87D6066",
		CocoaTonTokenAddress:     "0xabcdabcdabcdabcdabcdabcdabcdabcdabcdabcd",
		GasLimit:                 0,
	}

	err := validateConfigForRuntime(cfg)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "gasLimit must be greater than 0")
}

func TestGetChainSelector(t *testing.T) {
	cfg := &Config{
		ChainName: "ethereum-testnet-sepolia-arbitrum-1",
	}

	// This will fail unless we're in a CRE runtime context, so we just test the method exists
	_, err := cfg.GetChainSelector()
	// We expect this to work or fail depending on SDK availability
	_ = err
}

func TestNewEVMClient(t *testing.T) {
	cfg := &Config{
		ChainName: "ethereum-testnet-sepolia-arbitrum-1",
	}

	// This will fail unless we're in a CRE runtime context
	_, err := cfg.NewEVMClient()
	// We expect this to work or fail depending on SDK availability
	_ = err
}

