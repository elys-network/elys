# Off-Chain Logic for Order Executions

## Check Spot Order Execution Criteria Logic

### Requirements

  - `SpotOrder` object contains both `order_price` and `order_type`.
  - `market_price` can be retrieved from the oracle module.
  - If `order_type` is `STOPLOSS` and `market_price` is less than or equal to `order_price` then trigger `executeOrders` function.
  - If `order_type` is `LIMITSELL` and `market_price` is greater than or equal to `order_price` then trigger `executeOrders` function.
  - If `order_type` is `LIMITBUY` and `market_price` is less than or equal to `order_price` then trigger `executeOrders` function.

## Check Perpetual Order Execution Criteria Logic

### Requirements

  - `PerpetualOrder` object contains both `trigger_price` and `perpetual_order_type`.
  - `market_price` can be retrieved from the oracle module.
  - For `LONG` positions: If `perpetual_order_type` is `LIMITOPEN` and `market_price` is less than or equal to `trigger_price` then trigger `executeOrders` function.
  - For `SHORT` positions: If `perpetual_order_type` is `LIMITOPEN` and `market_price` is greater than or equal to `trigger_price` then trigger `executeOrders` function.

### executeOrders function
- `executeOrders(spot_order_ids: []u64, perpetual_order_ids: []u64)`
  - spot_order_ids: list of spot order ids that needs to be executed, function will verify and execute the spot orders.
  - perpetual_order_ids: list of perpetual order ids that needs to be executed, function will verify and execute the perpetual orders.
