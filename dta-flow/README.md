# DTA Flow – Digital Transfer Agent Workflow

## Overview

The **DTA (Digital Transfer Agent) Flow** is a Chainlink CRE workflow that orchestrates the subscription and redemption lifecycle for commodity-backed tokens in the Akua platform. Unlike the existing `escrow-flow` (which is analytics-only), `dta-flow` is an **active orchestrator** that:

- Listens for on-chain DTA request events (subscription/redemption)
- Prices requests using the on-chain `PriceOracleRouter`
- Runs placeholder compliance checks (with structure for external API integration)
- Executes settlement via EVM contract calls (minting/burning tokens, managing collateral)
- Posts lifecycle updates to the frontend via webhook

## Architecture

### Workflow Components

1. **Event Listeners**
   - `SubscriptionRequested` – User wants to buy commodity-backed tokens
   - `RedemptionRequested` – User wants to redeem tokens for physical commodity or cash

2. **Orchestration Steps**
   - **Pricing**: Fetch NAV-like price from `PriceOracleRouter` (placeholder for now)
   - **Compliance**: Validate requester, quantity, and eligibility (placeholder with TODOs for KYC/AML integration)
   - **Approval/Rejection**: Update request status on-chain via `updateRequestStatus`
   - **Settlement**:
     - Subscription: Interact with `CommodityEscrow`, `WarehouseReceiptRegistry`, `CocoaTonToken` to mint tokens
     - Redemption: Burn tokens and release collateral (placeholder implementation)

3. **Webhook Integration**
   - All lifecycle transitions (`Pending → Priced → Approved → Settled`) are posted to `/api/dta-webhook`
   - Frontend displays DTA requests in the **DTA Request Lifecycle** panel on the dashboard

## Contract

### CommodityDTARequestManager

**Deployed Address**: `0x4A52E850F9B76DA35D6B1E9055D60c007d396feF` (Arbitrum Sepolia)  
**Verification**: https://sepolia.arbiscan.io/address/0x4A52E850F9B76DA35D6B1E9055D60c007d396feF#code

The on-chain smart contract (`akua-cre/contracts/evm/src/CommodityDTARequestManager.sol`) manages:

- **Request Struct**: `id`, `requester`, `requestType`, `quantityTons`, `paymentToken`, `paymentAmount`, `status`, `createdAt`, `updatedAt`
- **Request Types**: `0 = Subscription`, `1 = Redemption`
- **Status Lifecycle**:
  - `0 = Pending` – Request created
  - `1 = Priced` – Price calculated
  - `2 = Approved` – Passed compliance
  - `3 = Rejected` – Failed compliance
  - `4 = Settled` – Settlement complete
  - `5 = Cancelled` – Request cancelled

- **Events**:
  - `SubscriptionRequested(id, requester, quantityTons, paymentToken)`
  - `RedemptionRequested(id, requester, quantityTons)`
  - `RequestStatusUpdated(id, status)`

- **Functions**:
  - `createSubscriptionRequest(quantityTons, paymentToken) → requestId`
  - `createRedemptionRequest(quantityTons) → requestId`
  - `updateRequestStatus(id, status)` (TODO: add role-based access control)
  - `getRequest(id) → Request`

## Configuration

### Staging (`config.staging.json`)

```json
{
  "webhookUrl": "http://localhost:3002/api/dta-webhook",
  "webhookSecret": "akua-dta-staging-secret",
  "chainName": "ethereum-testnet-sepolia-arbitrum-1",
  "network": "arbitrumSepolia",
  "dtaRequestManagerAddress": "UPDATE-WITH-DTA-REQUEST-MANAGER-ADDRESS",
  "escrowContractAddress": "0x4cE838B32469bAdD067C05EBc04A4c01EDcAE903",
  "receiptRegistryAddress": "0xa04f93a8810727Fd1B2867432F5612391a9daF15",
  "priceOracleAddress": "0x6199C8bf5be5513877899d9EB2cDd892F87D6066",
  "cocoaTonTokenAddress": "UPDATE-WITH-COCOA-TON-TOKEN-ADDRESS",
  "gasLimit": 800000
}
```

### Production (`config.production.json`)

```json
{
  "webhookUrl": "https://your-production-domain.com/api/dta-webhook",
  "webhookSecret": "CHANGE-THIS-TO-SECURE-RANDOM-SECRET-IN-PRODUCTION",
  "chainName": "ethereum-mainnet-arbitrum-1",
  "network": "arbitrumOne",
  "dtaRequestManagerAddress": "UPDATE-WITH-MAINNET-DTA-REQUEST-MANAGER-ADDRESS",
  "escrowContractAddress": "UPDATE-WITH-MAINNET-ESCROW-ADDRESS",
  "receiptRegistryAddress": "UPDATE-WITH-MAINNET-RECEIPT-REGISTRY-ADDRESS",
  "priceOracleAddress": "UPDATE-WITH-MAINNET-PRICE-ORACLE-ADDRESS",
  "cocoaTonTokenAddress": "UPDATE-WITH-MAINNET-COCOA-TON-TOKEN-ADDRESS",
  "gasLimit": 800000
}
```

## Webhook Payload

Example JSON posted to `/api/dta-webhook`:

```json
{
  "eventName": "DtaRequestSettled",
  "network": "arbitrumSepolia",
  "chainName": "ethereum-testnet-sepolia-arbitrum-1",
  "requestId": "1",
  "requestType": "Subscription",
  "status": "Settled",
  "requester": "0x1234567890123456789012345678901234567890",
  "quantityTons": "100",
  "paymentToken": "0xabcdabcdabcdabcdabcdabcdabcdabcdabcdabcd",
  "paymentAmount": "0",
  "reason": "",
  "txHash": "0x...",
  "blockNumber": "12345678",
  "timestamp": 1733503200
}
```

## Deployment

### ✅ Contract Deployed

**CommodityDTARequestManager**: `0x4A52E850F9B76DA35D6B1E9055D60c007d396feF` (Arbitrum Sepolia)

### Prerequisites for CRE Workflow Deployment

1. ✅ `CommodityDTARequestManager` deployed to Arbitrum Sepolia
2. ✅ Config files updated with deployed contract address
3. Set `CRE_DTA_WEBHOOK_SECRET` in frontend `.env.local`:
   ```bash
   CRE_DTA_WEBHOOK_SECRET=akua-dta-staging-secret
   ```

### Deploy to CRE

**⏸️ DEPLOYMENT CURRENTLY BLOCKED**

CRE workflow deployment is in **early access** and not yet available for testnet/mainnet.

```bash
# Will run this when CRE access is granted:
cd akua-cre/dta-flow
cre workflow deploy . --target staging-settings

# Production deployment (requires HTTPS webhook URL)
cre workflow deploy . --target production-settings
```

**Request Access**: https://cre.chain.link/request-access

### Manual Testing (Available Now)

Since CRE deployment is blocked, you can test manually:

1. **Test Contract Directly**:
   ```bash
   cast send 0x4A52E850F9B76DA35D6B1E9055D60c007d396feF \
     "createSubscriptionRequest(uint256,address)" \
     100 0xYourUSDCAddress \
     --rpc-url https://sepolia-rollup.arbitrum.io/rpc \
     --private-key $PRIVATE_KEY
   ```

2. **Simulate Webhook Events**:
   ```bash
   curl -X POST http://localhost:3002/api/dta-webhook \
     -H "Content-Type: application/json" \
     -H "X-CRE-Webhook-Secret: akua-dta-staging-secret" \
     -d '{"eventName":"DtaRequestCreated",...}'
   ```

3. **Verify in Frontend**: Check `/dashboard` → **DTA Request Lifecycle** panel

**See**: `DTA_MANUAL_TESTING_GUIDE.md` in repo root for complete testing instructions.

### Future: Automated Deployment (When Access Granted)

Once CRE access is available:

1. Check CRE dashboard for workflow status
2. Monitor frontend `/dashboard` → **DTA Request Lifecycle** panel automatically
3. Create test requests - CRE will handle orchestration
4. Verify lifecycle events appear automatically in the DTA panel

## Testing

Run Go tests:

```bash
cd akua-cre/dta-flow
go test -v
```

### Manual Testing

1. Deploy `CommodityDTARequestManager` to Arbitrum Sepolia
2. Start frontend dev server (`npm run dev`)
3. Create a subscription request via contract interaction
4. Watch the CRE workflow logs and frontend DTA panel for lifecycle updates

## TODOs

### High Priority

- **Access Control**: Add role-based access control to `updateRequestStatus` in `CommodityDTARequestManager.sol`
- **PriceOracleRouter Integration**: Replace placeholder pricing logic with actual on-chain price fetching
- **Compliance API**: Integrate external KYC/AML/compliance API in `runComplianceCheck`
- **Settlement Logic**: Implement real token minting/burning and collateral management for subscriptions and redemptions

### Medium Priority

- **Gas Optimization**: Review gas costs for `updateRequestStatus` and batch updates if needed
- **Error Handling**: Add retries and dead-letter queue for failed webhook deliveries
- **Monitoring**: Add metrics and alerts for failed compliance checks or settlement errors

### Low Priority

- **Request Cancellation**: Add user-initiated cancellation flow
- **Fee Calculation**: Implement fee structure for DTA operations
- **Multi-token Support**: Extend beyond USDC for payment tokens

## Architecture Diagram

```
┌─────────────────────────────────────────────────────────────────┐
│                      On-Chain (Arbitrum)                        │
│                                                                 │
│  ┌────────────────────────────────────────────────────────┐    │
│  │ CommodityDTARequestManager                             │    │
│  │                                                        │    │
│  │ • createSubscriptionRequest()                          │    │
│  │ • createRedemptionRequest()                            │    │
│  │ • updateRequestStatus()                                │    │
│  │ • getRequest()                                         │    │
│  └─────────────┬──────────────────────────────────────────┘    │
│                │ Events:                                        │
│                │ • SubscriptionRequested                        │
│                │ • RedemptionRequested                          │
│                │ • RequestStatusUpdated                         │
└────────────────┼──────────────────────────────────────────────┘
                 │
                 ▼
┌─────────────────────────────────────────────────────────────────┐
│               Chainlink CRE (dta-flow)                          │
│                                                                 │
│  ┌──────────────────────────────────────────────────────────┐  │
│  │ Event Handlers                                           │  │
│  │                                                          │  │
│  │ onSubscriptionRequested:                                 │  │
│  │   1. Fetch request details                               │  │
│  │   2. Calculate price (PriceOracleRouter)                 │  │
│  │   3. Run compliance checks                               │  │
│  │   4. Update status (Approved/Rejected)                   │  │
│  │   5. Execute settlement (mint tokens)                    │  │
│  │   6. Update status (Settled)                             │  │
│  │   7. POST to /api/dta-webhook                            │  │
│  │                                                          │  │
│  │ onRedemptionRequested:                                   │  │
│  │   1. Fetch request details                               │  │
│  │   2. Run compliance checks                               │  │
│  │   3. Update status (Approved/Rejected)                   │  │
│  │   4. Execute settlement (burn tokens)                    │  │
│  │   5. Update status (Settled)                             │  │
│  │   6. POST to /api/dta-webhook                            │  │
│  └──────────────┬───────────────────────────────────────────┘  │
│                 │                                               │
└─────────────────┼───────────────────────────────────────────────┘
                  │ HTTP POST
                  │ (webhook)
                  ▼
┌─────────────────────────────────────────────────────────────────┐
│               Frontend (Next.js)                                │
│                                                                 │
│  /api/dta-webhook (POST)                                        │
│    ↓                                                            │
│  Global DTA Event Store (in-memory)                             │
│    ↓                                                            │
│  /api/dta-webhook (GET)                                         │
│    ↓                                                            │
│  useDtaEvents hook                                              │
│    ↓                                                            │
│  DtaRequestsPanel component                                     │
│    • Displays requests in table                                 │
│    • Filters by status, type, network                           │
│    • Color-coded status badges                                  │
└─────────────────────────────────────────────────────────────────┘
```

## Comparison: DTA Flow vs. Escrow Flow

| Feature                     | `escrow-flow`              | `dta-flow`                      |
|-----------------------------|----------------------------|---------------------------------|
| **Purpose**                 | Analytics-only monitoring  | Active orchestration            |
| **On-chain Writes**         | No                         | Yes (status updates, settlement)|
| **Webhook Endpoint**        | `/api/cre-webhook`         | `/api/dta-webhook`              |
| **Events Monitored**        | Escrow & Receipt events    | DTA request events              |
| **Business Logic**          | None (pure event forwarding)| Pricing, compliance, settlement|
| **Configuration Secret**    | `CRE_WEBHOOK_SECRET`       | `CRE_DTA_WEBHOOK_SECRET`        |
| **Frontend Component**      | `CreEventsFeed`            | `DtaRequestsPanel`              |

## Support

For questions or issues:
- Review the main project README at repo root
- Check `CHAINLINK_CRE_INTEGRATION.md` for CRE setup
- See `PROJECT_OVERVIEW.md` for high-level architecture

