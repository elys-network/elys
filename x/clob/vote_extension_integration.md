# CLOB Vote Extension Integration

## Overview
This document describes how to integrate the in-memory order matching service with vote extensions in the application.

## Components Created

### 1. In-Memory Order Matching Service (`memory_orderbook.go`)
- `MemoryOrderBook`: Manages in-memory order book for fast matching
- `GetOperationsToPropose`: Main function that loads orders from state and matches them
- Returns `[]types.MatchedOrderExecution` for vote extensions

### 2. Vote Extension Handler (`abci/vote_extension.go`)
- `ExtendVoteHandler`: Creates vote extensions with matched orders
- `VerifyVoteExtensionHandler`: Validates vote extensions from other validators
- Calls `GetOperationsToPropose` to get matched orders

### 3. Proposal Handler (`abci/proposal.go`)
- `PrepareProposalHandler`: Aggregates vote extensions and creates tx
- `ProcessProposalHandler`: Validates proposals with matched orders
- Creates `MsgExecuteMatchedOrders` transaction

### 4. Message Server (`msg_server_execute_matched_orders.go`)
- Handles `MsgExecuteMatchedOrders` messages
- Uses existing `Exchange` function to execute trades
- Maintains consistency with existing order matching logic

## Integration Steps

To integrate this into your app, add the following to `app/app.go`:

```go
import (
    clobabci "github.com/elys-network/elys/v7/x/clob/abci"
)

// In your app initialization:

// Create vote extension handler
voteExtHandler := clobabci.NewVoteExtensionHandler(
    app.Logger(),
    app.ClobKeeper,
)
app.SetExtendVoteHandler(voteExtHandler.ExtendVoteHandler())
app.SetVerifyVoteExtensionHandler(voteExtHandler.VerifyVoteExtensionHandler())

// Create proposal handler
proposalHandler := clobabci.NewProposalHandler(
    app.Logger(),
    app.ClobKeeper,
    app.txConfig.TxEncoder(),
    app.txConfig.TxDecoder(),
)
app.SetPrepareProposal(proposalHandler.PrepareProposalHandler())
app.SetProcessProposal(proposalHandler.ProcessProposalHandler())
```

## How It Works

1. **During Block Extension (ExtendVote)**:
   - Each validator loads orders from state into memory
   - Matches buy and sell orders using price-time priority
   - Creates vote extension with matched orders

2. **During Block Preparation (PrepareProposal)**:
   - Proposer aggregates vote extensions from all validators
   - Deduplicates matched orders
   - Creates `MsgExecuteMatchedOrders` transaction
   - Adds transaction to the block

3. **During Block Execution**:
   - `ExecuteMatchedOrders` message server processes the transaction
   - Calls `Exchange` function for each matched order
   - Updates perpetual positions and subaccounts
   - Emits trade events

## Key Features

- **Deterministic Matching**: All validators compute the same matches
- **Reuses Existing Logic**: Uses the existing `Exchange` function
- **Order Management**: Automatically handles order updates/deletions
- **Event Emission**: Emits trade events for monitoring

## Testing

To test the integration:

1. Place buy and sell orders that cross
2. Wait for the next block
3. Check that orders are matched and executed
4. Verify trade events in the block

## Notes

- The `GetOperationsToPropose` function is called during vote extension
- Matched orders are executed using the existing `Exchange` function
- Order removal is handled automatically by the existing order matching flow
- The system maintains backward compatibility with existing order types