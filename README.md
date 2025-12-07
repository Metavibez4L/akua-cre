## Akua CRE Workflows

### Overview

This repository contains **Chainlink CRE workflows** for the Akua commodity platform.  
The workflows listen to on‑chain events on Arbitrum, normalize them, and send JSON webhooks to the Akua backend and frontend for analytics and orchestration.

There are two primary workflows:

- **`escrow-flow`** – Read‑only **analytics** pipeline for escrow and warehouse receipt events.
- **`dta-flow`** – **Digital Transfer Agent (DTA)** workflow that orchestrates subscription and redemption requests for commodity‑backed tokens.

Both workflows are compiled to WASM using `cre-sdk-go` and deployed via the CRE CLI.

### Repository Structure

- **`escrow-flow/`**
  - `workflow.go` – CRE workflow wiring and event handlers for escrows and warehouse receipts.
  - `main.go` – WASM entrypoint (`wasm.NewRunner(cre.ParseJSON[Config]).Run(InitWorkflow)`).
  - `config.staging.json` / `config.production.json` – Workflow runtime configuration.
  - `SIMULATION_SCENARIOS.md` – Canonical CRE simulation scenarios for escrow & receipt lifecycle.
  - `README.md` – Detailed documentation of the escrow analytics flow.

- **`dta-flow/`**
  - `workflow.go` – DTA orchestration logic (subscription/redemption lifecycle, status updates).
  - `main.go` – WASM entrypoint for the DTA workflow.
  - `config.staging.json` / `config.production.json` – DTA runtime configuration.
  - `README.md` – Detailed DTA flow documentation (pricing, compliance, settlement, webhooks).

- **Root files**
  - `ENVIRONMENTS.md` – Describes staging vs production environments, chain names, and canonical contract addresses.
  - `project.yaml` – CRE project‑level settings (RPCs and targets).
  - `secrets.yaml` – CRE secrets mapping (names only, no values).
  - `go.mod`, `go.sum` – Go module definition and dependencies.

### High‑Level Architecture

At a high level the system works like this:

1. **On‑chain contracts (Arbitrum)**
   - `CommodityEscrow`, `WarehouseReceiptRegistry`, `CommodityDTARequestManager`, and related contracts emit events.
2. **Chainlink CRE workflows (this repo)**
   - `escrow-flow` subscribes to escrow & receipt events and forwards **normalized analytics events** to a webhook.
   - `dta-flow` subscribes to DTA request events and **orchestrates** status updates and settlement, emitting lifecycle webhooks.
3. **Backend & Frontend**
   - Webhooks terminate in Next.js API routes (e.g. `/api/cre-webhook`, `/api/dta-webhook`).
   - Data is stored and surfaced in dashboard components: CRE Event Feed, DTA Request Lifecycle panel, etc.

Simplified diagram:

```
On‑chain (Arbitrum)
  ├─ CommodityEscrow
  ├─ WarehouseReceiptRegistry
  └─ CommodityDTARequestManager
          │ events
          ▼
Chainlink CRE (this repo)
  ├─ escrow-flow  ──► POST /api/cre-webhook
  └─ dta-flow     ──► POST /api/dta-webhook
          │
          ▼
Akua Backend / Frontend
  ├─ Store & index events
  └─ Dashboard views (escrows, receipts, DTA requests)
```

### Workflows in Brief

- **Escrow Flow (`escrow-flow`)**
  - **Purpose**: Analytics‑only; no on‑chain writes.
  - **Events**:
    - Escrow: `EscrowCreated`, `EscrowFunded`, `EscrowConfirmed`, `EscrowCompleted`, `EscrowCancelled`, `EscrowDisputed`.
    - Receipts: `ReceiptIssued`, `ReceiptUpdated`, `ReceiptRedeemed`, `ReceiptCancelled`.
  - **Behavior**:
    - Decodes logs via generated Go bindings.
    - Maps to a unified `EventPayload` (escrow/receipt fields + tx metadata).
    - Sends JSON to the configured `webhookUrl` with `X-CRE-Webhook-Secret`.
  - **Docs**: See `escrow-flow/README.md` for full field‑level and scenario‑level documentation.

- **DTA Flow (`dta-flow`)**
  - **Purpose**: Orchestrates **subscription** and **redemption** requests for commodity‑backed tokens.
  - **Events**:
    - `SubscriptionRequested`, `RedemptionRequested` (from `CommodityDTARequestManager`).
  - **Behavior** (current implementation):
    - Sends `"DtaRequestCreated"` webhook when a request is created.
    - Runs placeholder compliance checks (`runComplianceCheck`).
    - Updates on‑chain status (`Pending → Approved/Rejected → Settled`) using `updateRequestStatus`.
    - Emits `"DtaRequestApproved"`, `"DtaRequestRejected"`, `"DtaRequestSettled"` webhooks to `/api/dta-webhook`.
    - Contains TODOs for real pricing (via `PriceOracleRouter`), KYC/AML integration, mint/burn & collateral settlement.
  - **Docs**: See `dta-flow/README.md` for full lifecycle, payload examples, and architecture diagram.

### Environments & Configuration

Environment‑specific details live in `ENVIRONMENTS.md` and the per‑workflow config files.

- **Staging**
  - Chain: `ethereum-testnet-sepolia-arbitrum-1` (Arbitrum Sepolia).
  - Network label: `arbitrumSepolia`.
  - Webhooks typically point to `http://localhost:3002/...`.
  - Canonical contract addresses (e.g. escrow v10, receipt registry) are documented in `ENVIRONMENTS.md`.

- **Production**
  - Chain: `ethereum-mainnet-arbitrum-1` (Arbitrum One).
  - Network label: `arbitrumOne`.
  - Webhooks **must** be HTTPS and **must not** point to `localhost`.

CRE CLI targets are defined as:

- **Project‑level** (`project.yaml`)
  - `staging-settings` / `production-settings` define RPC URLs per chain.
- **Workflow‑level** (`escrow-flow/workflow.yaml`, `dta-flow` configs)
  - Each workflow has its own `staging-settings` and `production-settings` target pointing at its config files.

### Getting Started (Local)

- **Prerequisites**
  - Go toolchain (version as per `go.mod`).
  - Chainlink CRE CLI (for simulation and deployment).
  - Access to an Arbitrum Sepolia RPC endpoint for staging.

- **Clone & install**

  ```bash
  git clone https://github.com/Metavibez4L/akua-cre.git
  cd akua-cre
  go test ./...
  ```

- **Simulate escrow workflow**

  ```bash
  cd escrow-flow
  cre workflow simulate . --target staging-settings
  ```

- **Simulate DTA workflow**

  ```bash
  cd ../dta-flow
  cre workflow simulate . --target staging-settings
  ```

Use the scenarios in `escrow-flow/SIMULATION_SCENARIOS.md` and the examples in `dta-flow/README.md` to drive end‑to‑end tests.

### Deployment Notes

- **Escrow Flow**
  - Ready for CRE deployment once you have CRE access.
  - Uses the webhook URL and secret from `escrow-flow/config.*.json`.

- **DTA Flow**
  - CRE deployment is currently gated by Chainlink CRE early‑access (see notes in `dta-flow/README.md`).
  - In the meantime, you can:
    - Interact with `CommodityDTARequestManager` directly on Arbitrum Sepolia.
    - Manually simulate webhook events against your Next.js API.

Before any mainnet deployment:

- Double‑check contract addresses.
- Ensure webhook URLs are HTTPS and backed by production infrastructure.
- Verify secrets are set in the target environment (e.g. `CRE_WEBHOOK_SECRET`, `CRE_DTA_WEBHOOK_SECRET`).

### Where to Go Next

- **Escrow analytics flow details**: `escrow-flow/README.md`
- **DTA orchestration flow details**: `dta-flow/README.md`
- **Environments & test strategy**: `ENVIRONMENTS.md`

This top‑level README is the landing page on GitHub; use the linked documents for deeper, flow‑specific detail.


