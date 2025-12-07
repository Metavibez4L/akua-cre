# CRE Workflow Simulation Scenarios

This document contains curated test scenarios for simulating the Akua escrow workflow using real Arbitrum Sepolia testnet transactions.

## How to Run Simulations

From the `akua-cre/escrow-flow` directory, run:

```bash
cre workflow simulate . --target staging-settings
```

When prompted, select the appropriate trigger type (Log Trigger) and provide the transaction hash and event index from the scenarios below.

## Canonical Test Scenarios

### Scenario 1: Full Escrow Lifecycle (Happy Path)

**Description**: Complete escrow flow from creation through funding, confirmation, and completion.

#### Events to Simulate

1. **EscrowCreated**
   - Transaction Hash: `TBD - Deploy and create test escrow on Arbitrum Sepolia`
   - Event Index: `0`
   - Expected Behavior: Webhook receives EscrowCreated event with seller, buyer, commodity token, amounts, and receipt ID

2. **EscrowFunded**
   - Transaction Hash: `TBD - Fund the test escrow`
   - Event Index: `0`
   - Expected Behavior: Webhook receives EscrowFunded event with funder address and amount

3. **EscrowConfirmed**
   - Transaction Hash: `TBD - Confirm the test escrow`
   - Event Index: `0`
   - Expected Behavior: Webhook receives EscrowConfirmed event with confirmer address

4. **EscrowCompleted**
   - Transaction Hash: `TBD - Complete the test escrow`
   - Event Index: `0`
   - Expected Behavior: Webhook receives EscrowCompleted event

### Scenario 2: Escrow Cancellation Flow

**Description**: Escrow is created but cancelled before completion.

#### Events to Simulate

1. **EscrowCreated**
   - Transaction Hash: `TBD`
   - Event Index: `0`
   - Expected Behavior: Webhook receives EscrowCreated event

2. **EscrowCancelled**
   - Transaction Hash: `TBD`
   - Event Index: `0`
   - Expected Behavior: Webhook receives EscrowCancelled event with cancelledBy address

### Scenario 3: Escrow Dispute Flow

**Description**: Escrow is funded but then disputed.

#### Events to Simulate

1. **EscrowCreated**
   - Transaction Hash: `TBD`
   - Event Index: `0`
   - Expected Behavior: Webhook receives EscrowCreated event

2. **EscrowFunded**
   - Transaction Hash: `TBD`
   - Event Index: `0`
   - Expected Behavior: Webhook receives EscrowFunded event

3. **EscrowDisputed**
   - Transaction Hash: `TBD`
   - Event Index: `0`
   - Expected Behavior: Webhook receives EscrowDisputed event with disputer address

### Scenario 4: Warehouse Receipt Lifecycle

**Description**: Complete warehouse receipt flow from issuance through updates to redemption.

#### Events to Simulate

1. **ReceiptIssued**
   - Transaction Hash: `TBD - Issue test warehouse receipt`
   - Event Index: `0`
   - Expected Behavior: Webhook receives ReceiptIssued event with depositor, warehouse, commodity type, quantity, and dates

2. **ReceiptUpdated**
   - Transaction Hash: `TBD - Update receipt quantity`
   - Event Index: `0`
   - Expected Behavior: Webhook receives ReceiptUpdated event with new quantity

3. **ReceiptRedeemed**
   - Transaction Hash: `TBD - Redeem the receipt`
   - Event Index: `0`
   - Expected Behavior: Webhook receives ReceiptRedeemed event with redeemer address

### Scenario 5: Receipt Cancellation

**Description**: Warehouse receipt is issued but then cancelled.

#### Events to Simulate

1. **ReceiptIssued**
   - Transaction Hash: `TBD`
   - Event Index: `0`
   - Expected Behavior: Webhook receives ReceiptIssued event

2. **ReceiptCancelled**
   - Transaction Hash: `TBD`
   - Event Index: `0`
   - Expected Behavior: Webhook receives ReceiptCancelled event with cancelledBy address

## Setting Up Test Data

To populate this document with actual transaction hashes:

1. **Deploy test contracts** (if not already deployed to Arbitrum Sepolia)
   - Contracts are deployed at:
     - Escrow: `0x4cE838B32469bAdD067C05EBc04A4c01EDcAE903`
     - Receipt Registry: `0xa04f93a8810727Fd1B2867432F5612391a9daF15`

2. **Execute test transactions** on Arbitrum Sepolia:
   ```bash
   # Using Foundry/Cast or Hardhat
   # Example: Create an escrow
   cast send <ESCROW_ADDRESS> "createEscrow(...)" --rpc-url $ARBITRUM_SEPOLIA_RPC --private-key $TEST_PRIVATE_KEY
   ```

3. **Find transaction hashes** on Arbiscan Sepolia:
   - Visit https://sepolia.arbiscan.io/
   - Search for your contract addresses
   - Copy transaction hashes that emitted the relevant events

4. **Update this document** with the actual transaction hashes and event indices

5. **Run simulations** to verify each scenario works end-to-end

## Validation Checklist

For each simulation, verify:

- [ ] CRE workflow executes without errors
- [ ] Webhook endpoint receives the event
- [ ] Event payload contains all expected fields
- [ ] Event appears in frontend dashboard/history view
- [ ] Event is correctly labeled as "Simulated" or shows staging environment
- [ ] Transaction hash links to correct Arbiscan page
- [ ] All addresses and IDs are formatted correctly
- [ ] Timestamps are present and reasonable

## Frontend Verification

After running simulations:

1. Open frontend at `http://localhost:3002/dashboard`
2. Verify events appear in the CRE Event Feed
3. Check that simulated/staging events are visually distinguished
4. Test filtering and search functionality
5. Verify transaction links work correctly

## Continuous Testing

These scenarios should be re-run:

- After any changes to `workflow.go` event handlers
- After updates to contract ABIs or bindings
- Before deploying any new workflow versions
- When testing webhook authentication changes
- Before mainnet deployment (as final validation)

## Notes

- Always use the `staging-settings` target for simulations
- Never use production config for testing
- Keep this document updated as new test transactions are created
- Archive old transaction hashes if contracts are redeployed
- Document any edge cases or unusual behaviors discovered during simulation

