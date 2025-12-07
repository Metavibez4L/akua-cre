package main

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"strings"
	"time"

	"github.com/ethereum/go-ethereum/common"

	"akua-cre/contracts/evm/src/generated/commodity_escrow"
	"akua-cre/contracts/evm/src/generated/warehouse_receipt_registry"

	"github.com/smartcontractkit/cre-sdk-go/capabilities/blockchain/evm"
	"github.com/smartcontractkit/cre-sdk-go/capabilities/blockchain/evm/bindings"
	"github.com/smartcontractkit/cre-sdk-go/capabilities/networking/http"
	"github.com/smartcontractkit/cre-sdk-go/cre"
)

// Config defines the runtime configuration for the Akua analytics workflow.
type Config struct {
	// HTTP webhook that will receive normalized Akua events
	WebhookURL string `json:"webhookUrl"`

	// Shared secret for webhook authentication
	WebhookSecret string `json:"webhookSecret"`

	// CRE chain selector name, e.g. "ethereum-testnet-sepolia-arbitrum-1"
	ChainName string `json:"chainName"`

	// Human-friendly network label, e.g. "arbitrumSepolia"
	Network string `json:"network"`

	// On-chain contract addresses
	EscrowAddress          string `json:"escrowContractAddress"`
	ReceiptRegistryAddress string `json:"receiptRegistryAddress"`
	GasLimit               uint64 `json:"gasLimit"` // reserved for future on-chain writes
}

// GetChainSelector resolves the CRE chain selector for the configured chain.
func (c *Config) GetChainSelector() (uint64, error) {
	return evm.ChainSelectorFromName(c.ChainName)
}

// NewEVMClient constructs a CRE EVM client for the configured chain.
func (c *Config) NewEVMClient() (*evm.Client, error) {
	selector, err := c.GetChainSelector()
	if err != nil {
		return nil, err
	}
	return &evm.Client{
		ChainSelector: selector,
	}, nil
}

// EventPayload is the normalized payload that will be sent to the webhook.
type EventPayload struct {
	EventName string `json:"eventName"`
	Network   string `json:"network"`
	ChainName string `json:"chainName"`

	// Escrow fields
	EscrowID       string `json:"escrowId,omitempty"`
	Seller         string `json:"seller,omitempty"`
	Buyer          string `json:"buyer,omitempty"`
	CommodityToken string `json:"commodityToken,omitempty"`
	TokenAmount    string `json:"tokenAmount,omitempty"`
	PricePerToken  string `json:"pricePerToken,omitempty"`

	// Receipt fields
	ReceiptID     string `json:"receiptId,omitempty"`
	Depositor     string `json:"depositor,omitempty"`
	Warehouse     string `json:"warehouse,omitempty"`
	CommodityType string `json:"commodityType,omitempty"`
	Quantity      string `json:"quantity,omitempty"`
	DepositDate   string `json:"depositDate,omitempty"`
	ExpiryDate    string `json:"expiryDate,omitempty"`

	// Event-specific fields
	Funder      string `json:"funder,omitempty"`
	Amount      string `json:"amount,omitempty"` // Payment amount (e.g., ETH/USDC paid in EscrowFunded)
	Confirmer   string `json:"confirmer,omitempty"`
	CancelledBy string `json:"cancelledBy,omitempty"`
	Disputer    string `json:"disputer,omitempty"`
	Redeemer    string `json:"redeemer,omitempty"`
	NewQuantity string `json:"newQuantity,omitempty"`

	// Transaction metadata
	TxHash      string `json:"txHash,omitempty"`
	BlockNumber string `json:"blockNumber,omitempty"`
	Timestamp   uint64 `json:"timestamp,omitempty"`
}

// WebhookResult is a trivial type used for CRE HTTP consensus aggregation.
type WebhookResult struct {
	Success bool `json:"success"`
}

// populateTxMetadata extracts transaction metadata from a decoded log and populates the EventPayload.
// NOTE: Newer versions of the CRE EVM log type represent hashes and block numbers differently:
//   - TxHash is now a []byte instead of a common.Hash
//   - BlockTimestamp has been removed
// We only populate fields that are available in the current SDK to keep the workflow forward compatible.
func populateTxMetadata(p *EventPayload, log *evm.Log) {
	if log == nil {
		return
	}

	// TxHash is exposed as raw bytes – convert to a 0x-prefixed hex string for the payload.
	if len(log.TxHash) > 0 {
		p.TxHash = common.BytesToHash(log.TxHash).Hex()
	}

	// BlockNumber may be represented as a big-int like type (e.g. *pb.BigInt) with String() support.
	if log.BlockNumber != nil {
		p.BlockNumber = log.BlockNumber.String()
	}

	// BlockTimestamp is no longer available on evm.Log in the current SDK, so we intentionally
	// leave EventPayload.Timestamp at its zero value. If we later need wall-clock time we can
	// derive it from an external block lookup inside the webhook consumer.
}

// validateConfigForRuntime performs startup validation on configuration.
// It mirrors the validation used in tests and adds a few runtime safety checks.
func validateConfigForRuntime(cfg *Config) error {
	if cfg.WebhookURL == "" {
		return fmt.Errorf("webhookUrl cannot be empty")
	}
	if cfg.ChainName == "" {
		return fmt.Errorf("chainName cannot be empty")
	}
	if cfg.Network == "" {
		return fmt.Errorf("network cannot be empty")
	}
	if !common.IsHexAddress(cfg.EscrowAddress) {
		return fmt.Errorf("invalid escrowContractAddress: %s", cfg.EscrowAddress)
	}
	if !common.IsHexAddress(cfg.ReceiptRegistryAddress) {
		return fmt.Errorf("invalid receiptRegistryAddress: %s", cfg.ReceiptRegistryAddress)
	}
	if cfg.GasLimit == 0 {
		return fmt.Errorf("gasLimit must be greater than 0")
	}

	// For non-testnet networks (i.e. presumed production), enforce HTTPS and disallow localhost.
	if !strings.Contains(strings.ToLower(cfg.Network), "sepolia") &&
		!strings.Contains(strings.ToLower(cfg.Network), "testnet") {
		if !strings.HasPrefix(strings.ToLower(cfg.WebhookURL), "https://") {
			return fmt.Errorf("production webhookUrl must use HTTPS, got: %s", cfg.WebhookURL)
		}
		if strings.Contains(cfg.WebhookURL, "localhost") || strings.Contains(cfg.WebhookURL, "127.0.0.1") {
			return fmt.Errorf("production webhookUrl cannot point to localhost: %s", cfg.WebhookURL)
		}
	}

	return nil
}

// InitWorkflow wires up log triggers for Akua contracts on Arbitrum Sepolia
// and forwards decoded events to the configured webhook as JSON.
func InitWorkflow(config *Config, logger *slog.Logger, secretsProvider cre.SecretsProvider) (cre.Workflow[*Config], error) {
	// Validate configuration early so we fail-fast on misconfiguration.
	if err := validateConfigForRuntime(config); err != nil {
		return nil, fmt.Errorf("config validation failed: %w", err)
	}

	evmClient, err := config.NewEVMClient()
	if err != nil {
		return nil, fmt.Errorf("failed to create EVM client: %w", err)
	}

	escrowAddress := common.HexToAddress(config.EscrowAddress)
	receiptAddress := common.HexToAddress(config.ReceiptRegistryAddress)

	escrowContract, err := commodity_escrow.NewCommodityEscrow(evmClient, escrowAddress, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create CommodityEscrow binding: %w", err)
	}

	receiptRegistry, err := warehouse_receipt_registry.NewWarehouseReceiptRegistry(evmClient, receiptAddress, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create WarehouseReceiptRegistry binding: %w", err)
	}

	chainSelector, err := config.GetChainSelector()
	if err != nil {
		return nil, fmt.Errorf("failed to get chain selector: %w", err)
	}

	// Escrow events
	escrowCreatedTrigger, err := escrowContract.LogTriggerEscrowCreatedLog(
		chainSelector,
		evm.ConfidenceLevel_CONFIDENCE_LEVEL_LATEST,
		[]commodity_escrow.EscrowCreatedTopics{},
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create EscrowCreated trigger: %w", err)
	}

	escrowFundedTrigger, err := escrowContract.LogTriggerEscrowFundedLog(
		chainSelector,
		evm.ConfidenceLevel_CONFIDENCE_LEVEL_LATEST,
		[]commodity_escrow.EscrowFundedTopics{},
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create EscrowFunded trigger: %w", err)
	}

	escrowConfirmedTrigger, err := escrowContract.LogTriggerEscrowConfirmedLog(
		chainSelector,
		evm.ConfidenceLevel_CONFIDENCE_LEVEL_LATEST,
		[]commodity_escrow.EscrowConfirmedTopics{},
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create EscrowConfirmed trigger: %w", err)
	}

	escrowCompletedTrigger, err := escrowContract.LogTriggerEscrowCompletedLog(
		chainSelector,
		evm.ConfidenceLevel_CONFIDENCE_LEVEL_LATEST,
		[]commodity_escrow.EscrowCompletedTopics{},
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create EscrowCompleted trigger: %w", err)
	}

	escrowCancelledTrigger, err := escrowContract.LogTriggerEscrowCancelledLog(
		chainSelector,
		evm.ConfidenceLevel_CONFIDENCE_LEVEL_LATEST,
		[]commodity_escrow.EscrowCancelledTopics{},
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create EscrowCancelled trigger: %w", err)
	}

	escrowDisputedTrigger, err := escrowContract.LogTriggerEscrowDisputedLog(
		chainSelector,
		evm.ConfidenceLevel_CONFIDENCE_LEVEL_LATEST,
		[]commodity_escrow.EscrowDisputedTopics{},
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create EscrowDisputed trigger: %w", err)
	}

	// Receipt events
	receiptIssuedTrigger, err := receiptRegistry.LogTriggerReceiptIssuedLog(
		chainSelector,
		evm.ConfidenceLevel_CONFIDENCE_LEVEL_LATEST,
		[]warehouse_receipt_registry.ReceiptIssuedTopics{},
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create ReceiptIssued trigger: %w", err)
	}

	receiptRedeemedTrigger, err := receiptRegistry.LogTriggerReceiptRedeemedLog(
		chainSelector,
		evm.ConfidenceLevel_CONFIDENCE_LEVEL_LATEST,
		[]warehouse_receipt_registry.ReceiptRedeemedTopics{},
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create ReceiptRedeemed trigger: %w", err)
	}

	receiptCancelledTrigger, err := receiptRegistry.LogTriggerReceiptCancelledLog(
		chainSelector,
		evm.ConfidenceLevel_CONFIDENCE_LEVEL_LATEST,
		[]warehouse_receipt_registry.ReceiptCancelledTopics{},
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create ReceiptCancelled trigger: %w", err)
	}

	receiptUpdatedTrigger, err := receiptRegistry.LogTriggerReceiptUpdatedLog(
		chainSelector,
		evm.ConfidenceLevel_CONFIDENCE_LEVEL_LATEST,
		[]warehouse_receipt_registry.ReceiptUpdatedTopics{},
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create ReceiptUpdated trigger: %w", err)
	}

	workflow := cre.Workflow[*Config]{
		cre.Handler(escrowCreatedTrigger, onEscrowCreated),
		cre.Handler(escrowFundedTrigger, onEscrowFunded),
		cre.Handler(escrowConfirmedTrigger, onEscrowConfirmed),
		cre.Handler(escrowCompletedTrigger, onEscrowCompleted),
		cre.Handler(escrowCancelledTrigger, onEscrowCancelled),
		cre.Handler(escrowDisputedTrigger, onEscrowDisputed),
		cre.Handler(receiptIssuedTrigger, onReceiptIssued),
		cre.Handler(receiptRedeemedTrigger, onReceiptRedeemed),
		cre.Handler(receiptCancelledTrigger, onReceiptCancelled),
		cre.Handler(receiptUpdatedTrigger, onReceiptUpdated),
	}

	return workflow, nil
}

func onEscrowCreated(config *Config, runtime cre.Runtime, payload *bindings.DecodedLog[commodity_escrow.EscrowCreatedDecoded]) (string, error) {
	logger := runtime.Logger()
	start := time.Now()

	escrowID := common.BytesToHash(payload.Data.EscrowId[:]).Hex()
	receiptID := common.BytesToHash(payload.Data.ReceiptId[:]).Hex()

	logger.Info("processing EscrowCreated event",
		"eventName", "EscrowCreated",
		"network", config.Network,
		"escrowId", escrowID,
		"receiptId", receiptID,
		"seller", payload.Data.Seller.Hex(),
		"buyer", payload.Data.Buyer.Hex())

	p := &EventPayload{
		EventName:      "EscrowCreated",
		Network:        config.Network,
		ChainName:      config.ChainName,
		EscrowID:       escrowID,
		Seller:         payload.Data.Seller.Hex(),
		Buyer:          payload.Data.Buyer.Hex(),
		CommodityToken: payload.Data.CommodityToken.Hex(),
		TokenAmount:    payload.Data.TokenAmount.String(),
		PricePerToken:  payload.Data.PricePerToken.String(),
		ReceiptID:      receiptID,
	}
	populateTxMetadata(p, payload.Log)

	if err := sendWebhook(config, runtime, logger, p); err != nil {
		return "", err
	}
	logger.Info("EscrowCreated event processed",
		"escrowId", escrowID,
		"processingTimeMs", time.Since(start).Milliseconds())
	return "EscrowCreated forwarded", nil
}

func onEscrowFunded(config *Config, runtime cre.Runtime, payload *bindings.DecodedLog[commodity_escrow.EscrowFundedDecoded]) (string, error) {
	logger := runtime.Logger()
	start := time.Now()

	escrowID := common.BytesToHash(payload.Data.EscrowId[:]).Hex()

	logger.Info("processing EscrowFunded event",
		"eventName", "EscrowFunded",
		"network", config.Network,
		"escrowId", escrowID,
		"funder", payload.Data.Funder.Hex(),
		"amount", payload.Data.Amount.String())

	p := &EventPayload{
		EventName: "EscrowFunded",
		Network:   config.Network,
		ChainName: config.ChainName,
		EscrowID:  escrowID,
		Funder:    payload.Data.Funder.Hex(),
		Amount:    payload.Data.Amount.String(), // Payment amount in ETH/USDC/USDT
	}
	populateTxMetadata(p, payload.Log)

	if err := sendWebhook(config, runtime, logger, p); err != nil {
		return "", err
	}
	logger.Info("EscrowFunded event processed",
		"escrowId", escrowID,
		"processingTimeMs", time.Since(start).Milliseconds())
	return "EscrowFunded forwarded", nil
}

func onEscrowConfirmed(config *Config, runtime cre.Runtime, payload *bindings.DecodedLog[commodity_escrow.EscrowConfirmedDecoded]) (string, error) {
	logger := runtime.Logger()
	start := time.Now()

	escrowID := common.BytesToHash(payload.Data.EscrowId[:]).Hex()

	logger.Info("processing EscrowConfirmed event",
		"eventName", "EscrowConfirmed",
		"network", config.Network,
		"escrowId", escrowID,
		"confirmer", payload.Data.Confirmer.Hex())

	p := &EventPayload{
		EventName: "EscrowConfirmed",
		Network:   config.Network,
		ChainName: config.ChainName,
		EscrowID:  escrowID,
		Confirmer: payload.Data.Confirmer.Hex(),
	}
	populateTxMetadata(p, payload.Log)

	if err := sendWebhook(config, runtime, logger, p); err != nil {
		return "", err
	}
	logger.Info("EscrowConfirmed event processed",
		"escrowId", escrowID,
		"processingTimeMs", time.Since(start).Milliseconds())
	return "EscrowConfirmed forwarded", nil
}

func onEscrowCompleted(config *Config, runtime cre.Runtime, payload *bindings.DecodedLog[commodity_escrow.EscrowCompletedDecoded]) (string, error) {
	logger := runtime.Logger()
	start := time.Now()

	escrowID := common.BytesToHash(payload.Data.EscrowId[:]).Hex()

	logger.Info("processing EscrowCompleted event",
		"eventName", "EscrowCompleted",
		"network", config.Network,
		"escrowId", escrowID)

	p := &EventPayload{
		EventName: "EscrowCompleted",
		Network:   config.Network,
		ChainName: config.ChainName,
		EscrowID:  escrowID,
	}
	populateTxMetadata(p, payload.Log)

	if err := sendWebhook(config, runtime, logger, p); err != nil {
		return "", err
	}
	logger.Info("EscrowCompleted event processed",
		"escrowId", escrowID,
		"processingTimeMs", time.Since(start).Milliseconds())
	return "EscrowCompleted forwarded", nil
}

func onEscrowCancelled(config *Config, runtime cre.Runtime, payload *bindings.DecodedLog[commodity_escrow.EscrowCancelledDecoded]) (string, error) {
	logger := runtime.Logger()
	start := time.Now()

	escrowID := common.BytesToHash(payload.Data.EscrowId[:]).Hex()

	logger.Info("processing EscrowCancelled event",
		"eventName", "EscrowCancelled",
		"network", config.Network,
		"escrowId", escrowID,
		"cancelledBy", payload.Data.CancelledBy.Hex())

	p := &EventPayload{
		EventName:   "EscrowCancelled",
		Network:     config.Network,
		ChainName:   config.ChainName,
		EscrowID:    escrowID,
		CancelledBy: payload.Data.CancelledBy.Hex(),
	}
	populateTxMetadata(p, payload.Log)

	if err := sendWebhook(config, runtime, logger, p); err != nil {
		return "", err
	}
	logger.Info("EscrowCancelled event processed",
		"escrowId", escrowID,
		"processingTimeMs", time.Since(start).Milliseconds())
	return "EscrowCancelled forwarded", nil
}

func onEscrowDisputed(config *Config, runtime cre.Runtime, payload *bindings.DecodedLog[commodity_escrow.EscrowDisputedDecoded]) (string, error) {
	logger := runtime.Logger()
	start := time.Now()

	escrowID := common.BytesToHash(payload.Data.EscrowId[:]).Hex()

	logger.Info("processing EscrowDisputed event",
		"eventName", "EscrowDisputed",
		"network", config.Network,
		"escrowId", escrowID,
		"disputer", payload.Data.Disputer.Hex())

	p := &EventPayload{
		EventName: "EscrowDisputed",
		Network:   config.Network,
		ChainName: config.ChainName,
		EscrowID:  escrowID,
		Disputer:  payload.Data.Disputer.Hex(),
	}
	populateTxMetadata(p, payload.Log)

	if err := sendWebhook(config, runtime, logger, p); err != nil {
		return "", err
	}
	logger.Info("EscrowDisputed event processed",
		"escrowId", escrowID,
		"processingTimeMs", time.Since(start).Milliseconds())
	return "EscrowDisputed forwarded", nil
}

func onReceiptIssued(config *Config, runtime cre.Runtime, payload *bindings.DecodedLog[warehouse_receipt_registry.ReceiptIssuedDecoded]) (string, error) {
	logger := runtime.Logger()
	start := time.Now()

	receiptID := common.BytesToHash(payload.Data.ReceiptId[:]).Hex()

	logger.Info("processing ReceiptIssued event",
		"eventName", "ReceiptIssued",
		"network", config.Network,
		"receiptId", receiptID,
		"depositor", payload.Data.Depositor.Hex(),
		"warehouse", payload.Data.Warehouse.Hex(),
		"commodityType", payload.Data.CommodityType,
		"quantity", payload.Data.Quantity.String())

	p := &EventPayload{
		EventName:     "ReceiptIssued",
		Network:       config.Network,
		ChainName:     config.ChainName,
		ReceiptID:     receiptID,
		Depositor:     payload.Data.Depositor.Hex(),
		Warehouse:     payload.Data.Warehouse.Hex(),
		CommodityType: payload.Data.CommodityType,
		Quantity:      payload.Data.Quantity.String(),
		DepositDate:   payload.Data.DepositDate.String(),
		ExpiryDate:    payload.Data.ExpiryDate.String(),
	}
	populateTxMetadata(p, payload.Log)

	if err := sendWebhook(config, runtime, logger, p); err != nil {
		return "", err
	}
	logger.Info("ReceiptIssued event processed",
		"receiptId", receiptID,
		"processingTimeMs", time.Since(start).Milliseconds())
	return "ReceiptIssued forwarded", nil
}

func onReceiptRedeemed(config *Config, runtime cre.Runtime, payload *bindings.DecodedLog[warehouse_receipt_registry.ReceiptRedeemedDecoded]) (string, error) {
	logger := runtime.Logger()
	start := time.Now()

	receiptID := common.BytesToHash(payload.Data.ReceiptId[:]).Hex()

	logger.Info("processing ReceiptRedeemed event",
		"eventName", "ReceiptRedeemed",
		"network", config.Network,
		"receiptId", receiptID,
		"redeemer", payload.Data.Redeemer.Hex())

	p := &EventPayload{
		EventName: "ReceiptRedeemed",
		Network:   config.Network,
		ChainName: config.ChainName,
		ReceiptID: receiptID,
		Redeemer:  payload.Data.Redeemer.Hex(),
	}
	populateTxMetadata(p, payload.Log)

	if err := sendWebhook(config, runtime, logger, p); err != nil {
		return "", err
	}
	logger.Info("ReceiptRedeemed event processed",
		"receiptId", receiptID,
		"processingTimeMs", time.Since(start).Milliseconds())
	return "ReceiptRedeemed forwarded", nil
}

func onReceiptCancelled(config *Config, runtime cre.Runtime, payload *bindings.DecodedLog[warehouse_receipt_registry.ReceiptCancelledDecoded]) (string, error) {
	logger := runtime.Logger()
	start := time.Now()

	receiptID := common.BytesToHash(payload.Data.ReceiptId[:]).Hex()

	logger.Info("processing ReceiptCancelled event",
		"eventName", "ReceiptCancelled",
		"network", config.Network,
		"receiptId", receiptID,
		"cancelledBy", payload.Data.CancelledBy.Hex())

	p := &EventPayload{
		EventName:   "ReceiptCancelled",
		Network:     config.Network,
		ChainName:   config.ChainName,
		ReceiptID:   receiptID,
		CancelledBy: payload.Data.CancelledBy.Hex(),
	}
	populateTxMetadata(p, payload.Log)

	if err := sendWebhook(config, runtime, logger, p); err != nil {
		return "", err
	}
	logger.Info("ReceiptCancelled event processed",
		"receiptId", receiptID,
		"processingTimeMs", time.Since(start).Milliseconds())
	return "ReceiptCancelled forwarded", nil
}

func onReceiptUpdated(config *Config, runtime cre.Runtime, payload *bindings.DecodedLog[warehouse_receipt_registry.ReceiptUpdatedDecoded]) (string, error) {
	logger := runtime.Logger()
	start := time.Now()

	receiptID := common.BytesToHash(payload.Data.ReceiptId[:]).Hex()

	logger.Info("processing ReceiptUpdated event",
		"eventName", "ReceiptUpdated",
		"network", config.Network,
		"receiptId", receiptID,
		"newQuantity", payload.Data.NewQuantity.String())

	p := &EventPayload{
		EventName:   "ReceiptUpdated",
		Network:     config.Network,
		ChainName:   config.ChainName,
		ReceiptID:   receiptID,
		NewQuantity: payload.Data.NewQuantity.String(),
	}
	populateTxMetadata(p, payload.Log)

	if err := sendWebhook(config, runtime, logger, p); err != nil {
		return "", err
	}
	logger.Info("ReceiptUpdated event processed",
		"receiptId", receiptID,
		"processingTimeMs", time.Since(start).Milliseconds())
	return "ReceiptUpdated forwarded", nil
}

// sendWebhook uses CRE's HTTP capability to POST the event payload to the configured webhook URL.
func sendWebhook(config *Config, runtime cre.Runtime, logger *slog.Logger, payload *EventPayload) error {
	client := &http.Client{}

	_, err := http.SendRequest(
		config,
		runtime,
		client,
		func(cfg *Config, lg *slog.Logger, sr *http.SendRequester) (*WebhookResult, error) {
			body, err := json.Marshal(payload)
			if err != nil {
				lg.Error("failed to marshal webhook payload", "err", err)
				return nil, err
			}

			headers := map[string]string{
				"Content-Type": "application/json",
			}
			if cfg.WebhookSecret != "" {
				headers["X-CRE-Webhook-Secret"] = cfg.WebhookSecret
			}

			resp, err := sr.SendRequest(&http.Request{
				Method:  "POST",
				Url:     cfg.WebhookURL,
				Body:    body,
				Headers: headers,
			}).Await()
			if err != nil {
				lg.Error("webhook transport failed", "err", err)
				return nil, err
			}

			if resp.StatusCode >= 400 {
				bodyPreview := string(resp.Body)
				if len(bodyPreview) > 256 {
					bodyPreview = bodyPreview[:256] + "..."
				}
				lg.Error("webhook HTTP error",
					"status", resp.StatusCode,
					"body", bodyPreview,
					"url", cfg.WebhookURL)
				return nil, fmt.Errorf("webhook returned HTTP %d", resp.StatusCode)
			}

			lg.Info("webhook delivered successfully",
				"status", resp.StatusCode,
				"url", cfg.WebhookURL)

			return &WebhookResult{Success: true}, nil
		},
		cre.ConsensusAggregationFromTags[*WebhookResult](),
	).Await()

	if err != nil {
		logger.Error("failed to send webhook", "err", err)
		return err
	}

	return nil
}
