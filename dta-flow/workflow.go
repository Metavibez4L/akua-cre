package main

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"math/big"
	"strings"
	"time"

	"github.com/ethereum/go-ethereum/common"

	"akua-cre/contracts/evm/src/generated/commodity_dta_request_manager"
	"akua-cre/contracts/evm/src/generated/commodity_escrow"
	"akua-cre/contracts/evm/src/generated/warehouse_receipt_registry"

	"github.com/smartcontractkit/cre-sdk-go/capabilities/blockchain/evm"
	"github.com/smartcontractkit/cre-sdk-go/capabilities/blockchain/evm/bindings"
	"github.com/smartcontractkit/cre-sdk-go/capabilities/networking/http"
	"github.com/smartcontractkit/cre-sdk-go/cre"
)

// Config defines the runtime configuration for the DTA workflow.
type Config struct {
	// HTTP webhook that will receive DTA lifecycle events
	WebhookURL string `json:"webhookUrl"`

	// Shared secret for webhook authentication
	WebhookSecret string `json:"webhookSecret"`

	// CRE chain selector name, e.g. "ethereum-testnet-sepolia-arbitrum-1"
	ChainName string `json:"chainName"`

	// Human-friendly network label, e.g. "arbitrumSepolia"
	Network string `json:"network"`

	// On-chain contract addresses
	DtaRequestManagerAddress string `json:"dtaRequestManagerAddress"`
	EscrowAddress            string `json:"escrowContractAddress"`
	ReceiptRegistryAddress   string `json:"receiptRegistryAddress"`
	PriceOracleAddress       string `json:"priceOracleAddress"`
	CocoaTonTokenAddress     string `json:"cocoaTonTokenAddress"`
	GasLimit                 uint64 `json:"gasLimit"`
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

// DtaEventPayload is the normalized payload that will be sent to the DTA webhook.
type DtaEventPayload struct {
	EventName string `json:"eventName"` // "DtaRequestCreated", "DtaRequestPriced", "DtaRequestApproved", "DtaRequestRejected", "DtaRequestSettled"
	Network   string `json:"network"`
	ChainName string `json:"chainName"`

	// DTA Request fields
	RequestID     string `json:"requestId"`
	RequestType   string `json:"requestType"` // "Subscription" or "Redemption"
	Status        string `json:"status"`      // "Pending", "Priced", "Approved", "Rejected", "Settled", "Cancelled"
	Requester     string `json:"requester"`
	QuantityTons  string `json:"quantityTons"`
	PaymentToken  string `json:"paymentToken,omitempty"`
	PaymentAmount string `json:"paymentAmount,omitempty"`
	Reason        string `json:"reason,omitempty"` // Rejection reason or other notes

	// Transaction metadata
	TxHash      string `json:"txHash,omitempty"`
	BlockNumber string `json:"blockNumber,omitempty"`
	Timestamp   uint64 `json:"timestamp,omitempty"`
}

// WebhookResult is a trivial type used for CRE HTTP consensus aggregation.
type WebhookResult struct {
	Success bool `json:"success"`
}

// Request status constants matching the Solidity contract
const (
	STATUS_PENDING   uint8 = 0
	STATUS_PRICED    uint8 = 1
	STATUS_APPROVED  uint8 = 2
	STATUS_REJECTED  uint8 = 3
	STATUS_SETTLED   uint8 = 4
	STATUS_CANCELLED uint8 = 5
)

// Request type constants
const (
	REQUEST_TYPE_SUBSCRIPTION uint8 = 0
	REQUEST_TYPE_REDEMPTION   uint8 = 1
)

// populateTxMetadata extracts transaction metadata from a decoded log and populates the DtaEventPayload.
func populateTxMetadata(p *DtaEventPayload, log *evm.Log) {
	if log == nil {
		return
	}

	if len(log.TxHash) > 0 {
		p.TxHash = common.BytesToHash(log.TxHash).Hex()
	}

	if log.BlockNumber != nil {
		p.BlockNumber = log.BlockNumber.String()
	}
}

// validateConfigForRuntime performs startup validation on configuration.
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
	if !common.IsHexAddress(cfg.DtaRequestManagerAddress) {
		return fmt.Errorf("invalid dtaRequestManagerAddress: %s", cfg.DtaRequestManagerAddress)
	}
	if !common.IsHexAddress(cfg.EscrowAddress) {
		return fmt.Errorf("invalid escrowContractAddress: %s", cfg.EscrowAddress)
	}
	if !common.IsHexAddress(cfg.ReceiptRegistryAddress) {
		return fmt.Errorf("invalid receiptRegistryAddress: %s", cfg.ReceiptRegistryAddress)
	}
	if !common.IsHexAddress(cfg.PriceOracleAddress) {
		return fmt.Errorf("invalid priceOracleAddress: %s", cfg.PriceOracleAddress)
	}
	if !common.IsHexAddress(cfg.CocoaTonTokenAddress) {
		return fmt.Errorf("invalid cocoaTonTokenAddress: %s", cfg.CocoaTonTokenAddress)
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

// InitWorkflow wires up log triggers for DTA request events and orchestrates the subscription/redemption lifecycle.
func InitWorkflow(config *Config, logger *slog.Logger, secretsProvider cre.SecretsProvider) (cre.Workflow[*Config], error) {
	// Validate configuration early so we fail-fast on misconfiguration.
	if err := validateConfigForRuntime(config); err != nil {
		return nil, fmt.Errorf("config validation failed: %w", err)
	}

	evmClient, err := config.NewEVMClient()
	if err != nil {
		return nil, fmt.Errorf("failed to create EVM client: %w", err)
	}

	dtaManagerAddress := common.HexToAddress(config.DtaRequestManagerAddress)
	escrowAddress := common.HexToAddress(config.EscrowAddress)
	receiptAddress := common.HexToAddress(config.ReceiptRegistryAddress)

	dtaManager, err := commodity_dta_request_manager.NewCommodityDTARequestManager(evmClient, dtaManagerAddress, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create CommodityDTARequestManager binding: %w", err)
	}

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

	// DTA event triggers
	subscriptionRequestedTrigger, err := dtaManager.LogTriggerSubscriptionRequestedLog(
		chainSelector,
		evm.ConfidenceLevel_CONFIDENCE_LEVEL_LATEST,
		[]commodity_dta_request_manager.SubscriptionRequestedTopics{},
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create SubscriptionRequested trigger: %w", err)
	}

	redemptionRequestedTrigger, err := dtaManager.LogTriggerRedemptionRequestedLog(
		chainSelector,
		evm.ConfidenceLevel_CONFIDENCE_LEVEL_LATEST,
		[]commodity_dta_request_manager.RedemptionRequestedTopics{},
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create RedemptionRequested trigger: %w", err)
	}

	workflow := cre.Workflow[*Config]{
		cre.Handler(subscriptionRequestedTrigger, onSubscriptionRequested),
		cre.Handler(redemptionRequestedTrigger, onRedemptionRequested),
	}

	return workflow, nil
}

// onSubscriptionRequested handles a new subscription request (buy commodity-backed tokens).
// It prices the request, runs compliance checks, and orchestrates settlement.
func onSubscriptionRequested(config *Config, runtime cre.Runtime, payload *bindings.DecodedLog[commodity_dta_request_manager.SubscriptionRequestedDecoded]) (string, error) {
	logger := runtime.Logger()
	start := time.Now()

	requestID := payload.Data.Id.String()

	logger.Info("processing SubscriptionRequested event",
		"eventName", "SubscriptionRequested",
		"network", config.Network,
		"requestId", requestID,
		"requester", payload.Data.Requester.Hex(),
		"quantityTons", payload.Data.QuantityTons.String(),
		"paymentToken", payload.Data.PaymentToken.Hex())

	// Send initial "request created" webhook
	p := &DtaEventPayload{
		EventName:    "DtaRequestCreated",
		Network:      config.Network,
		ChainName:    config.ChainName,
		RequestID:    requestID,
		RequestType:  "Subscription",
		Status:       "Pending",
		Requester:    payload.Data.Requester.Hex(),
		QuantityTons: payload.Data.QuantityTons.String(),
		PaymentToken: payload.Data.PaymentToken.Hex(),
	}
	populateTxMetadata(p, payload.Log)

	if err := sendWebhook(config, runtime, logger, p); err != nil {
		return "", err
	}

	// Fetch full request details from contract
	evmClient, err := config.NewEVMClient()
	if err != nil {
		return "", fmt.Errorf("failed to create EVM client: %w", err)
	}

	dtaManagerAddress := common.HexToAddress(config.DtaRequestManagerAddress)
	dtaManager, err := commodity_dta_request_manager.NewCommodityDTARequestManager(evmClient, dtaManagerAddress, nil)
	if err != nil {
		return "", fmt.Errorf("failed to create DTA manager binding: %w", err)
	}

	chainSelector, err := config.GetChainSelector()
	if err != nil {
		return "", fmt.Errorf("failed to get chain selector: %w", err)
	}

	// Placeholder: In a full implementation, we would:
	// 1. Call PriceOracleRouter to get current price
	// 2. Compute total payment amount = quantityTons * price
	// 3. Update request with computed payment amount

	// For now, we'll use a placeholder price calculation
	// TODO: Integrate with actual PriceOracleRouter contract
	logger.Info("pricing request (placeholder logic)", "requestId", requestID)

	// Run compliance checks
	approved, reason := runComplianceCheck(payload.Data.Requester, payload.Data.QuantityTons, logger)

	if !approved {
		// Reject the request
		logger.Warn("request rejected by compliance", "requestId", requestID, "reason", reason)

		// Update on-chain status to REJECTED
		updateStatusCall, err := dtaManager.UpdateRequestStatus(chainSelector, payload.Data.Id, STATUS_REJECTED, config.GasLimit)
		if err != nil {
			return "", fmt.Errorf("failed to create updateRequestStatus call: %w", err)
		}

		_, err = evm.ExecuteEVMCall(config, runtime, updateStatusCall, evm.NewNonConsensusExecutor()).Await()
		if err != nil {
			logger.Error("failed to execute updateRequestStatus on-chain", "err", err)
			return "", err
		}

		// Send rejection webhook
		rejectionPayload := &DtaEventPayload{
			EventName:    "DtaRequestRejected",
			Network:      config.Network,
			ChainName:    config.ChainName,
			RequestID:    requestID,
			RequestType:  "Subscription",
			Status:       "Rejected",
			Requester:    payload.Data.Requester.Hex(),
			QuantityTons: payload.Data.QuantityTons.String(),
			PaymentToken: payload.Data.PaymentToken.Hex(),
			Reason:       reason,
		}
		if err := sendWebhook(config, runtime, logger, rejectionPayload); err != nil {
			return "", err
		}

		logger.Info("SubscriptionRequested processed (rejected)",
			"requestId", requestID,
			"processingTimeMs", time.Since(start).Milliseconds())
		return "SubscriptionRequested rejected", nil
	}

	// Approve the request
	logger.Info("request approved by compliance", "requestId", requestID)

	// Update on-chain status to APPROVED
	updateStatusCall, err := dtaManager.UpdateRequestStatus(chainSelector, payload.Data.Id, STATUS_APPROVED, config.GasLimit)
	if err != nil {
		return "", fmt.Errorf("failed to create updateRequestStatus (approved) call: %w", err)
	}

	_, err = evm.ExecuteEVMCall(config, runtime, updateStatusCall, evm.NewNonConsensusExecutor()).Await()
	if err != nil {
		logger.Error("failed to execute updateRequestStatus (approved) on-chain", "err", err)
		return "", err
	}

	// Send approval webhook
	approvalPayload := &DtaEventPayload{
		EventName:    "DtaRequestApproved",
		Network:      config.Network,
		ChainName:    config.ChainName,
		RequestID:    requestID,
		RequestType:  "Subscription",
		Status:       "Approved",
		Requester:    payload.Data.Requester.Hex(),
		QuantityTons: payload.Data.QuantityTons.String(),
		PaymentToken: payload.Data.PaymentToken.Hex(),
	}
	if err := sendWebhook(config, runtime, logger, approvalPayload); err != nil {
		return "", err
	}

	// Orchestrate settlement
	// TODO: In a full implementation, this would:
	// 1. Create/fund an escrow via CommodityEscrow contract
	// 2. Mint tokens via CocoaTonToken or issue receipts via WarehouseReceiptRegistry
	// 3. Transfer tokens to requester
	// For now, we just log and mark as settled
	logger.Info("orchestrating settlement (placeholder)", "requestId", requestID)

	// Update status to SETTLED
	settleStatusCall, err := dtaManager.UpdateRequestStatus(chainSelector, payload.Data.Id, STATUS_SETTLED, config.GasLimit)
	if err != nil {
		return "", fmt.Errorf("failed to create updateRequestStatus (settled) call: %w", err)
	}

	_, err = evm.ExecuteEVMCall(config, runtime, settleStatusCall, evm.NewNonConsensusExecutor()).Await()
	if err != nil {
		logger.Error("failed to execute updateRequestStatus (settled) on-chain", "err", err)
		return "", err
	}

	// Send settlement webhook
	settlementPayload := &DtaEventPayload{
		EventName:    "DtaRequestSettled",
		Network:      config.Network,
		ChainName:    config.ChainName,
		RequestID:    requestID,
		RequestType:  "Subscription",
		Status:       "Settled",
		Requester:    payload.Data.Requester.Hex(),
		QuantityTons: payload.Data.QuantityTons.String(),
		PaymentToken: payload.Data.PaymentToken.Hex(),
	}
	if err := sendWebhook(config, runtime, logger, settlementPayload); err != nil {
		return "", err
	}

	logger.Info("SubscriptionRequested processed (settled)",
		"requestId", requestID,
		"processingTimeMs", time.Since(start).Milliseconds())
	return "SubscriptionRequested settled", nil
}

// onRedemptionRequested handles a new redemption request (redeem tokens for physical commodity or cash).
func onRedemptionRequested(config *Config, runtime cre.Runtime, payload *bindings.DecodedLog[commodity_dta_request_manager.RedemptionRequestedDecoded]) (string, error) {
	logger := runtime.Logger()
	start := time.Now()

	requestID := payload.Data.Id.String()

	logger.Info("processing RedemptionRequested event",
		"eventName", "RedemptionRequested",
		"network", config.Network,
		"requestId", requestID,
		"requester", payload.Data.Requester.Hex(),
		"quantityTons", payload.Data.QuantityTons.String())

	// Send initial "request created" webhook
	p := &DtaEventPayload{
		EventName:    "DtaRequestCreated",
		Network:      config.Network,
		ChainName:    config.ChainName,
		RequestID:    requestID,
		RequestType:  "Redemption",
		Status:       "Pending",
		Requester:    payload.Data.Requester.Hex(),
		QuantityTons: payload.Data.QuantityTons.String(),
	}
	populateTxMetadata(p, payload.Log)

	if err := sendWebhook(config, runtime, logger, p); err != nil {
		return "", err
	}

	// Run compliance checks
	approved, reason := runComplianceCheck(payload.Data.Requester, payload.Data.QuantityTons, logger)

	evmClient, err := config.NewEVMClient()
	if err != nil {
		return "", fmt.Errorf("failed to create EVM client: %w", err)
	}

	dtaManagerAddress := common.HexToAddress(config.DtaRequestManagerAddress)
	dtaManager, err := commodity_dta_request_manager.NewCommodityDTARequestManager(evmClient, dtaManagerAddress, nil)
	if err != nil {
		return "", fmt.Errorf("failed to create DTA manager binding: %w", err)
	}

	chainSelector, err := config.GetChainSelector()
	if err != nil {
		return "", fmt.Errorf("failed to get chain selector: %w", err)
	}

	if !approved {
		// Reject the request
		logger.Warn("redemption request rejected by compliance", "requestId", requestID, "reason", reason)

		updateStatusCall, err := dtaManager.UpdateRequestStatus(chainSelector, payload.Data.Id, STATUS_REJECTED, config.GasLimit)
		if err != nil {
			return "", fmt.Errorf("failed to create updateRequestStatus call: %w", err)
		}

		_, err = evm.ExecuteEVMCall(config, runtime, updateStatusCall, evm.NewNonConsensusExecutor()).Await()
		if err != nil {
			logger.Error("failed to execute updateRequestStatus on-chain", "err", err)
			return "", err
		}

		rejectionPayload := &DtaEventPayload{
			EventName:    "DtaRequestRejected",
			Network:      config.Network,
			ChainName:    config.ChainName,
			RequestID:    requestID,
			RequestType:  "Redemption",
			Status:       "Rejected",
			Requester:    payload.Data.Requester.Hex(),
			QuantityTons: payload.Data.QuantityTons.String(),
			Reason:       reason,
		}
		if err := sendWebhook(config, runtime, logger, rejectionPayload); err != nil {
			return "", err
		}

		logger.Info("RedemptionRequested processed (rejected)",
			"requestId", requestID,
			"processingTimeMs", time.Since(start).Milliseconds())
		return "RedemptionRequested rejected", nil
	}

	// Approve the request
	logger.Info("redemption request approved by compliance", "requestId", requestID)

	updateStatusCall, err := dtaManager.UpdateRequestStatus(chainSelector, payload.Data.Id, STATUS_APPROVED, config.GasLimit)
	if err != nil {
		return "", fmt.Errorf("failed to create updateRequestStatus (approved) call: %w", err)
	}

	_, err = evm.ExecuteEVMCall(config, runtime, updateStatusCall, evm.NewNonConsensusExecutor()).Await()
	if err != nil {
		logger.Error("failed to execute updateRequestStatus (approved) on-chain", "err", err)
		return "", err
	}

	approvalPayload := &DtaEventPayload{
		EventName:    "DtaRequestApproved",
		Network:      config.Network,
		ChainName:    config.ChainName,
		RequestID:    requestID,
		RequestType:  "Redemption",
		Status:       "Approved",
		Requester:    payload.Data.Requester.Hex(),
		QuantityTons: payload.Data.QuantityTons.String(),
	}
	if err := sendWebhook(config, runtime, logger, approvalPayload); err != nil {
		return "", err
	}

	// Orchestrate settlement (placeholder: burn tokens, release collateral)
	// TODO: Implement actual token burning and collateral release
	logger.Info("orchestrating redemption settlement (placeholder)", "requestId", requestID)

	settleStatusCall, err := dtaManager.UpdateRequestStatus(chainSelector, payload.Data.Id, STATUS_SETTLED, config.GasLimit)
	if err != nil {
		return "", fmt.Errorf("failed to create updateRequestStatus (settled) call: %w", err)
	}

	_, err = evm.ExecuteEVMCall(config, runtime, settleStatusCall, evm.NewNonConsensusExecutor()).Await()
	if err != nil {
		logger.Error("failed to execute updateRequestStatus (settled) on-chain", "err", err)
		return "", err
	}

	settlementPayload := &DtaEventPayload{
		EventName:    "DtaRequestSettled",
		Network:      config.Network,
		ChainName:    config.ChainName,
		RequestID:    requestID,
		RequestType:  "Redemption",
		Status:       "Settled",
		Requester:    payload.Data.Requester.Hex(),
		QuantityTons: payload.Data.QuantityTons.String(),
	}
	if err := sendWebhook(config, runtime, logger, settlementPayload); err != nil {
		return "", err
	}

	logger.Info("RedemptionRequested processed (settled)",
		"requestId", requestID,
		"processingTimeMs", time.Since(start).Milliseconds())
	return "RedemptionRequested settled", nil
}

// runComplianceCheck performs placeholder compliance validation.
// TODO: In production, this should call an external compliance API or implement real KYC/AML logic.
func runComplianceCheck(requester common.Address, quantityTons *big.Int, logger *slog.Logger) (approved bool, reason string) {
	// Placeholder compliance logic
	if requester == (common.Address{}) {
		return false, "Invalid requester address (zero address)"
	}

	if quantityTons == nil || quantityTons.Cmp(big.NewInt(0)) == 0 {
		return false, "Invalid quantity (zero or nil)"
	}

	// Reject requests for more than 10,000 tons as a placeholder limit
	maxTons := big.NewInt(10000)
	if quantityTons.Cmp(maxTons) > 0 {
		return false, fmt.Sprintf("Quantity exceeds maximum allowed (%s tons)", maxTons.String())
	}

	// TODO: Add real compliance checks here:
	// - KYC/AML verification via external API
	// - Sanctions screening
	// - Accredited investor verification
	// - Jurisdiction checks

	logger.Info("compliance check passed (placeholder)", "requester", requester.Hex(), "quantityTons", quantityTons.String())
	return true, ""
}

// sendWebhook uses CRE's HTTP capability to POST the event payload to the configured webhook URL.
func sendWebhook(config *Config, runtime cre.Runtime, logger *slog.Logger, payload *DtaEventPayload) error {
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

