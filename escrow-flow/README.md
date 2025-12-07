## Escrow Flow – CRE Analytics Workflow Overview

### Purpose

The **escrow-flow** workflow is a **read‑only, analytics pipeline** that listens to on‑chain events from Akua’s commodity contracts on Arbitrum, normalizes them into a stable JSON schema, and forwards them to a backend webhook for storage, analytics, and UI display.

It does **not** perform any on‑chain writes or business‑logic decisions (those live in `dta-flow`). Instead, it acts as a reliable **event firehose** for:

- **Commodity escrow lifecycle events**
- **Warehouse receipt lifecycle events**

### High‑Level Architecture

- **On‑chain (Arbitrum)**
  - `CommodityEscrow` contract
  - `WarehouseReceiptRegistry` contract
  - Emits events such as `EscrowCreated`, `EscrowFunded`, `ReceiptIssued`, etc.

- **Chainlink CRE (escrow-flow)**
  - Compiled as a WASM workflow (`main.go`)
  - Uses `cre-sdk-go` EVM log triggers to subscribe to contract events
  - For each event:
    - Decodes the log into a typed struct
    - Maps fields into a normalized `EventPayload`
    - Enriches with transaction metadata (`txHash`, `blockNumber`)
    - Sends a JSON POST to the configured webhook URL

- **Webhook / Backend**
  - Receives JSON payloads
  - Authenticates via `X-CRE-Webhook-Secret` header
  - Persists, indexes, and exposes events to the frontend (e.g. dashboard/event feed)

### Runtime Configuration (`Config`)

The workflow is configured at runtime via a JSON `Config` object (see `config.staging.json` / `config.production.json`):

- **`webhookUrl`**: HTTP(S) endpoint that will receive normalized Akua events.
- **`webhookSecret`**: Shared secret; sent as `X-CRE-Webhook-Secret` for simple authentication.
- **`chainName`**: CRE chain selector name, e.g. `ethereum-testnet-sepolia-arbitrum-1`.
- **`network`**: Human‑readable label used in logs/payloads, e.g. `arbitrumSepolia`.
- **`escrowContractAddress`**: Address of the `CommodityEscrow` contract.
- **`receiptRegistryAddress`**: Address of the `WarehouseReceiptRegistry` contract.
- **`gasLimit`**: Reserved for future on‑chain writes (not used in this read‑only workflow).

At startup, `InitWorkflow` calls `validateConfigForRuntime` to ensure:

- **Required fields are present** (`webhookUrl`, `chainName`, `network`).
- **Addresses are valid hex addresses** for the escrow and receipt contracts.
- **`gasLimit > 0`**.
- **Production safety**:
  - Non‑test networks must use **HTTPS** webhooks.
  - Non‑test networks **cannot** point to `localhost` / `127.0.0.1`.

See `ENVIRONMENTS.md` for the canonical staging and production parameters (chain names, network labels, and deployed contract addresses).

### Workflow Wiring (`InitWorkflow`)

The core wiring lives in `workflow.go`:

1. **EVM client setup**
   - Uses `Config.GetChainSelector()` and `Config.NewEVMClient()` to build an `evm.Client` for the configured chain.

2. **Contract bindings**
   - `commodity_escrow.NewCommodityEscrow(evmClient, escrowAddress, nil)`
   - `warehouse_receipt_registry.NewWarehouseReceiptRegistry(evmClient, receiptAddress, nil)`

3. **Log triggers**
   - For `CommodityEscrow`:
     - `EscrowCreated`
     - `EscrowFunded`
     - `EscrowConfirmed`
     - `EscrowCompleted`
     - `EscrowCancelled`
     - `EscrowDisputed`
   - For `WarehouseReceiptRegistry`:
     - `ReceiptIssued`
     - `ReceiptRedeemed`
     - `ReceiptCancelled`
     - `ReceiptUpdated`

   Each trigger:

   - Uses the configured **chain selector**.
   - Uses `CONFIDENCE_LEVEL_LATEST` (latest block, no deep confirmations).
   - Subscribes with default topic filters (empty slices in code; can be extended later).

4. **Handlers**

   All triggers are composed into a single `cre.Workflow[*Config]`:

   - `cre.Handler(escrowCreatedTrigger, onEscrowCreated)`
   - `cre.Handler(escrowFundedTrigger, onEscrowFunded)`
   - `cre.Handler(escrowConfirmedTrigger, onEscrowConfirmed)`
   - `cre.Handler(escrowCompletedTrigger, onEscrowCompleted)`
   - `cre.Handler(escrowCancelledTrigger, onEscrowCancelled)`
   - `cre.Handler(escrowDisputedTrigger, onEscrowDisputed)`
   - `cre.Handler(receiptIssuedTrigger, onReceiptIssued)`
   - `cre.Handler(receiptRedeemedTrigger, onReceiptRedeemed)`
   - `cre.Handler(receiptCancelledTrigger, onReceiptCancelled)`
   - `cre.Handler(receiptUpdatedTrigger, onReceiptUpdated)`

The WASM entrypoint in `main.go` simply runs this workflow with JSON config parsing:

- `wasm.NewRunner(cre.ParseJSON[Config]).Run(InitWorkflow)`

### Normalized Event Payload (`EventPayload`)

All handlers build a single normalized struct before sending to the webhook:

- **Common metadata**
  - `eventName`: Name of the event, e.g. `"EscrowCreated"`, `"ReceiptIssued"`.
  - `network`: Human label from config, e.g. `"arbitrumSepolia"`.
  - `chainName`: CRE chain name, e.g. `"ethereum-testnet-sepolia-arbitrum-1"`.

- **Escrow‑related fields**
  - `escrowId`
  - `seller`
  - `buyer`
  - `commodityToken`
  - `tokenAmount`
  - `pricePerToken`

- **Receipt‑related fields**
  - `receiptId`
  - `depositor`
  - `warehouse`
  - `commodityType`
  - `quantity`
  - `depositDate`
  - `expiryDate`

- **Event‑specific fields**
  - `funder`
  - `amount` (payment amount for `EscrowFunded`)
  - `confirmer`
  - `cancelledBy`
  - `disputer`
  - `redeemer`
  - `newQuantity`

- **Transaction metadata**
  - `txHash`
  - `blockNumber`
  - `timestamp` (currently left at default; CRE’s EVM log type does not expose block timestamp in this version).

The helper `populateTxMetadata` maps data from `evm.Log` into `txHash` and `blockNumber`:

- Converts `log.TxHash` (raw bytes) into a 0x‑prefixed hash string.
- Uses `log.BlockNumber.String()` where available.

### Webhook Delivery (`sendWebhook`)

All handlers eventually call `sendWebhook`, which:

- Uses CRE’s HTTP capability (`github.com/smartcontractkit/cre-sdk-go/capabilities/networking/http`).
- Marshals `EventPayload` to JSON.
- Sets headers:
  - `"Content-Type": "application/json"`
  - `"X-CRE-Webhook-Secret": cfg.WebhookSecret` (if configured)
- Sends a `POST` to `cfg.WebhookURL`.
- Validates response:
  - Logs and errors on any `HTTP >= 400`, truncating the body preview for safety.
  - On success, logs the status code and URL.
- Wraps the call with `cre.ConsensusAggregationFromTags[*WebhookResult]()` to be compatible with CRE’s consensus abstraction (even though this workflow is read‑only).

If sending fails (transport error or HTTP ≥ 400), the handler logs the error and returns a failure up the CRE runtime; the webhook consumer should treat these events as **not delivered**.

### Event Flows in Detail

#### Escrow Lifecycle (Happy Path)

Maps to **Scenario 1** in `SIMULATION_SCENARIOS.md`.

1. **`EscrowCreated`**
   - Emitted when a new escrow is created with a linked warehouse receipt.
   - Handler: `onEscrowCreated`
   - Builds payload with:
     - `escrowId` (hash of `EscrowId` bytes)
     - `receiptId` (hash of `ReceiptId` bytes)
     - `seller`, `buyer`
     - `commodityToken`, `tokenAmount`, `pricePerToken`
   - Sends `"EscrowCreated"` webhook.

2. **`EscrowFunded`**
   - Emitted when funds are deposited into the escrow.
   - Handler: `onEscrowFunded`
   - Adds:
     - `funder`
     - `amount` (ETH/USDC/USDT, as string)
   - Sends `"EscrowFunded"` webhook.

3. **`EscrowConfirmed`**
   - Emitted when a party confirms the escrow terms (e.g. buyer acknowledgement).
   - Handler: `onEscrowConfirmed`
   - Adds:
     - `confirmer`
   - Sends `"EscrowConfirmed"` webhook.

4. **`EscrowCompleted`**
   - Emitted when the escrow successfully closes and settles.
   - Handler: `onEscrowCompleted`
   - Sends `"EscrowCompleted"` webhook.

In the UI, this should appear as a linear progression of events for the same `escrowId`, allowing time‑series visualization of the deal lifecycle.

#### Escrow Cancellation Flow

Maps to **Scenario 2** in `SIMULATION_SCENARIOS.md`.

1. `EscrowCreated` (as above).
2. `EscrowCancelled`
   - Handler: `onEscrowCancelled`
   - Adds:
     - `cancelledBy` (address of actor who cancelled).
   - Sends `"EscrowCancelled"` webhook.

This allows the backend to distinguish **cancelled** escrows from completed ones and record who initiated the cancellation.

#### Escrow Dispute Flow

Maps to **Scenario 3** in `SIMULATION_SCENARIOS.md`.

1. `EscrowCreated`
2. `EscrowFunded`
3. `EscrowDisputed`
   - Handler: `onEscrowDisputed`
   - Adds:
     - `disputer`
   - Sends `"EscrowDisputed"` webhook.

UI can surface these as **disputed deals** requiring off‑chain or on‑chain resolution.

#### Warehouse Receipt Lifecycle

Maps to **Scenario 4** in `SIMULATION_SCENARIOS.md`.

1. **`ReceiptIssued`**
   - Emitted when a warehouse issues a new commodity receipt.
   - Handler: `onReceiptIssued`
   - Adds:
     - `receiptId`
     - `depositor`
     - `warehouse`
     - `commodityType`
     - `quantity`
     - `depositDate`, `expiryDate`
   - Sends `"ReceiptIssued"` webhook.

2. **`ReceiptUpdated`**
   - Emitted when quantity or other mutable fields change.
   - Handler: `onReceiptUpdated`
   - Adds:
     - `newQuantity`
   - Sends `"ReceiptUpdated"` webhook.

3. **`ReceiptRedeemed`**
   - Emitted when the receipt is redeemed for physical commodity or equivalent.
   - Handler: `onReceiptRedeemed`
   - Adds:
     - `redeemer`
   - Sends `"ReceiptRedeemed"` webhook.

This gives a full audit trail of warehouse inventory movements across issuance, modification, and redemption.

#### Receipt Cancellation Flow

Maps to **Scenario 5** in `SIMULATION_SCENARIOS.md`.

1. `ReceiptIssued`
2. `ReceiptCancelled`
   - Handler: `onReceiptCancelled`
   - Adds:
     - `cancelledBy`
   - Sends `"ReceiptCancelled"` webhook.

Cancelled receipts can be excluded from active inventory and flagged in the UI.

### Example Webhook Payload

An example (simplified) payload for `ReceiptIssued` might look like:

```json
{
  "eventName": "ReceiptIssued",
  "network": "arbitrumSepolia",
  "chainName": "ethereum-testnet-sepolia-arbitrum-1",
  "receiptId": "0xabc123...",
  "depositor": "0x1111...1111",
  "warehouse": "0x2222...2222",
  "commodityType": "Cocoa Beans",
  "quantity": "1000",
  "depositDate": "1733539200",
  "expiryDate": "1736144800",
  "txHash": "0x...",
  "blockNumber": "1234567"
}
```

Exact field values and units should be interpreted using the on‑chain contract ABI and the backend’s business rules (e.g. whether `quantity` is in tons, kilos, or token units).

### Simulation & Testing

For curated end‑to‑end tests, see `SIMULATION_SCENARIOS.md`. The recommended flow:

1. **Run CRE simulations**
   - From `akua-cre/escrow-flow`:
   - `cre workflow simulate . --target staging-settings`
   - Use real Arbitrum Sepolia transaction hashes and event indices from your test contracts.

2. **Verify webhook behavior**
   - Confirm the webhook endpoint:
     - Receives events in order.
     - Authenticates `X-CRE-Webhook-Secret`.
     - Correctly persists and surfaces events in the UI.

3. **Regression testing**
   - Re‑run the canonical scenarios:
     - After any handler changes in `workflow.go`.
     - After ABI/binding updates.
     - Before any production deployment.

### Relationship to `dta-flow`

- **`escrow-flow`**: Read‑only **analytics** workflow for escrow and receipt events. No on‑chain writes; focuses on monitoring and observability.
- **`dta-flow`**: Active **orchestration** workflow for DTA requests (subscription/redemption), including on‑chain status updates and settlement (see `dta-flow/README.md`).

Together, they provide:

- A **control plane** for DTA lifecycle (`dta-flow`).
- A comprehensive **data plane** for escrow and inventory analytics (`escrow-flow`).


