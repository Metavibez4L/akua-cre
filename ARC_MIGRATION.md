# Arc Testnet Migration

Adds Circle's Arc (USDC-native-gas L1) as a CRE target alongside Arbitrum. Arbitrum targets are unchanged.

## Why this is mostly config

Both workflows resolve their chain from `chainName` in the config file via `evm.ChainSelectorFromName`, and
`network` is only a label (plus a testnet check in `validateConfigForRuntime`). No workflow code changes are needed.

## Arc facts used here

| Item | Value |
| --- | --- |
| CRE chain name | `arc-testnet` |
| Chain ID (testnet / mainnet) | `5042002` / `5042` |
| Testnet RPC | `https://rpc.testnet.arc.io` |
| Min versions for Arc Testnet | CRE CLI v1.0.7+, cre-sdk-go v1.1.4+ |
| Forwarder (simulation, MockKeystoneForwarder) | `0x6E9EE680ef59ef64Aa8C7371279c27E496b5eDc1` |
| Forwarder (production, KeystoneForwarder) | `0x76c9cf548b4179F8901cda1f8623568b58215E62` |
| USDC (ERC-20 interface, same on testnet and mainnet) | `0x3600000000000000000000000000000000000000` |
| EURC testnet | `0x89B50855Aa3bE2F677cD6303Cec089B5F319D72a` |

Arc mainnet is not in CRE's supported-networks list yet, so there is no production Arc target.
Verify against https://docs.chain.link/cre/supported-networks-go and the forwarder directory before relying on these.

## What changed

- `project.yaml`: new `staging-arc-settings` target (RPC for `arc-testnet`).
- `escrow-flow/workflow.yaml`, new `dta-flow/workflow.yaml`: `staging-arc-settings` target pointing at `config.staging-arc.json`.
- `escrow-flow/config.staging-arc.json`, `dta-flow/config.staging-arc.json`: `chainName: arc-testnet`, `network: arcTestnet`.
  Contract addresses are placeholders until the contracts are deployed on Arc Testnet.
- `dta-flow/workflow.yaml` is new because dta-flow had no workflow settings file; the existing `staging-settings`
  and `production-settings` targets are reproduced from the existing config paths.

## Steps still to do (run locally)

1. Bump the SDK (the pinned `cre-sdk-go v1.0.0` predates Arc support) and update the CLI:

   ```
   cre update
   go get github.com/smartcontractkit/cre-sdk-go@latest
   go get github.com/smartcontractkit/cre-sdk-go/capabilities/blockchain/evm@latest
   go get github.com/smartcontractkit/cre-sdk-go/capabilities/networking/http@latest
   go mod tidy && go test ./...
   ```

2. Deploy `CommodityEscrow`, `WarehouseReceiptRegistry`, `CommodityDTARequestManager`, `PriceOracleRouter` and the
   token to Arc Testnet (Foundry/Hardhat work as on any EVM chain), then fill the placeholders in both
   `config.staging-arc.json` files. Startup validation rejects `UPDATE-WITH-...` values, which is intentional.

3. Simulate:

   ```
   cd escrow-flow && cre workflow simulate . --target staging-arc-settings
   cd ../dta-flow  && cre workflow simulate . --target staging-arc-settings
   ```

## Things to check in the contracts

- **USDC decimals.** On Arc, native gas USDC uses 18 decimals while the ERC-20 interface uses 6. Use only the ERC-20
  interface (and `decimals()`) for `paymentToken` / `paymentAmount` handling, never `msg.value` math.
- **Write path.** CRE on-chain writes are delivered through the chain's forwarder. Confirm the DTA manager's
  status-update path accepts the Arc forwarder address above, and add the role-based access control that is
  already a TODO on `updateRequestStatus`.
- **Oracle.** Replace the placeholder pricing with Chainlink Data Feeds on Arc where a feed exists for the asset.
