# Akua CRE Workflow Environments

## Overview

This document defines the different environments used for Akua CRE workflow testing and deployment.

> **Arc Testnet** (`arc-testnet`, target `staging-arc-settings`) is available as an additional staging environment. See [ARC_MIGRATION.md](ARC_MIGRATION.md).

## Environment Configuration

### Staging (Testnet)
- **Purpose**: Development, testing, and simulation
- **Blockchain**: Arbitrum Sepolia Testnet
- **Chain Name**: `ethereum-testnet-sepolia-arbitrum-1`
- **Network Label**: `arbitrumSepolia`
- **RPC Endpoint**: `https://sepolia-rollup.arbitrum.io/rpc`
- **Webhook URL**: `http://localhost:3002/api/cre-webhook` (local dev)
- **Config File**: `config.staging.json`
- **CRE Workflow Name**: `escrow-flow-staging`

**Deployed Contracts (Arbitrum Sepolia)**:
- Escrow Contract (v10): `0x4cE838B32469bAdD067C05EBc04A4c01EDcAE903` ← **CURRENT - Use This**
- Receipt Registry: `0xa04f93a8810727Fd1B2867432F5612391a9daF15`

> **Migration Notice**: CommodityEscrow v10 deployed Dec 1, 2025 with critical security fixes. Previous versions are deprecated. See [CONTRACT_MIGRATION_GUIDE.md](../CONTRACT_MIGRATION_GUIDE.md) for details.

### Production (Mainnet)
- **Purpose**: Live production deployment
- **Blockchain**: Arbitrum One Mainnet
- **Chain Name**: `ethereum-mainnet-arbitrum-1`
- **Network Label**: `arbitrumOne`
- **RPC Endpoint**: TBD (requires secure RPC provider)
- **Webhook URL**: `https://your-production-domain.com/api/cre-webhook`
- **Config File**: `config.production.json`
- **CRE Workflow Name**: `escrow-flow-production`

**Deployed Contracts (Arbitrum One)**: TBD

## Canonical Test Scenarios

### 1. Full Escrow Lifecycle (Happy Path)
- EscrowCreated → EscrowFunded → EscrowConfirmed → EscrowCompleted
- Tests normal commodity purchase flow from creation to completion

### 2. Escrow Cancellation Flow
- EscrowCreated → EscrowCancelled
- Tests early termination by seller or buyer before funding

### 3. Escrow Dispute Flow
- EscrowCreated → EscrowFunded → EscrowDisputed
- Tests dispute resolution mechanism

### 4. Warehouse Receipt Lifecycle
- ReceiptIssued → ReceiptUpdated (quantity change) → ReceiptRedeemed
- Tests commodity storage and withdrawal flow

### 5. Receipt Cancellation
- ReceiptIssued → ReceiptCancelled
- Tests invalidation of warehouse receipts

## Testing Strategy

1. **Unit Tests**: Test individual event handlers with mock payloads
2. **Integration Tests**: Test webhook delivery with mock HTTP server
3. **CRE Simulations**: Use `cre workflow simulate` with real testnet transactions
4. **Frontend Validation**: Verify events display correctly in UI
5. **End-to-End Tests**: Full workflow from on-chain event to UI display

## Deployment Flow

```
Local Development (Staging Config)
    ↓
Unit Tests Pass
    ↓
Integration Tests Pass
    ↓
CRE Workflow Simulations (Staging)
    ↓
Frontend UI Validation
    ↓
Pre-Mainnet Checklist Complete
    ↓
Deploy to CRE Mainnet (Production Config)
```

## Important Notes

- **Never deploy test workflows to mainnet**: All simulation and testing must use staging configuration
- **Webhook secrets**: Use different secrets for staging vs production
- **Contract addresses**: Verify all addresses before mainnet deployment
- **RPC endpoints**: Use reliable, rate-limited RPC providers for production
- **Event visibility**: UI must clearly distinguish between staging and production events

