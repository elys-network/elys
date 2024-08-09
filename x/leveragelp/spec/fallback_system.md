# Fallback System Specification

## Overview
The fallback system is designed to handle position liquidations and closings in cases where the regular process(incentivized liquidation system) is either delayed or skipped. It systematically traverses a defined number of blocks to ensure that unhealthy positions are liquidated and positions are closed if they fall under a specified stop-loss price.

## Functionality

### 1. Epoch Check and Fallback Activation
- The fallback system activates when `epochPosition == 0` and `params.FallbackEnabled` is set to `true`.
- This condition checks if the current epoch has passed and whether fallback operations are permitted.

### 2. Pagination Setup
- A `PageRequest` is created with the following properties:
  - `Limit`: Set to `params.NumberPerBlock`, defining how many positions are processed per block.
- The current offset (the number of positions already processed) is retrieved from the state.

### 3. Position Retrieval
- The fallback system fetches positions using the `GetPositions` method, paginated by the `PageRequest` object.
- The offset is used to determine which positions have already been processed and where to continue.

### 4. Offset Management
- After processing positions, the system updates the offset:
  - If the current offset plus the `NumberPerBlock` is greater than or equal to the total open position count, the offset is deleted, indicating that all positions have been processed.
  - Otherwise, the offset is updated to continue processing the next batch of positions in subsequent fallback operations.

## Parameters
- **FallbackEnabled:** A boolean flag that enables or disables the fallback system.
- **NumberPerBlock:** Defines how many positions are processed per block during fallback operations.

