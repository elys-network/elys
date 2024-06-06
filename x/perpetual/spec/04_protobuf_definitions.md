<!--
order: 4
-->

# Protobuf Definitions

The following sections provide detailed protobuf definitions used in the Elys Network's Perpetual trading module.

## Genesis State Definition

The `GenesisState` message defines the initial state of the perpetual module at genesis.

### GenesisState

- **Params params**: Parameters for the module.
- **repeated Pool pool_list**: List of pools.
- **repeated MTP mtp_list**: List of MTPs.
- **repeated string address_whitelist**: List of whitelisted addresses.

```proto
message GenesisState {
  Params params = 1 [(gogoproto.nullable) = false];
  repeated Pool pool_list = 2 [(gogoproto.nullable) = false];
  repeated MTP mtp_list = 3 [(gogoproto.nullable) = false];
  repeated string address_whitelist = 4;
}
```

## Module Parameters

The `Params` message defines the configuration parameters for the perpetual module.

### Params

- **string leverage_max**: Maximum leverage allowed.
- **string borrow_interest_rate_max**: Maximum borrow interest rate.
- **string borrow_interest_rate_min**: Minimum borrow interest rate.
- **string borrow_interest_rate_increase**: Borrow interest rate increase factor.
- **string borrow_interest_rate_decrease**: Borrow interest rate decrease factor.
- **string health_gain_factor**: Health gain factor.
- **int64 epoch_length**: Length of an epoch.
- **int64 max_open_positions**: Maximum number of open positions.
- **string pool_open_threshold**: Pool open threshold.
- **string force_close_fund_percentage**: Percentage of force close fund.
- **string force_close_fund_address**: Address of the force close fund.
- **string incremental_borrow_interest_payment_fund_percentage**: Percentage for incremental borrow interest payment fund.
- **string incremental_borrow_interest_payment_fund_address**: Address for incremental borrow interest payment fund.
- **string safety_factor**: Safety factor.
- **bool incremental_borrow_interest_payment_enabled**: Flag to enable incremental borrow interest payment.
- **bool whitelisting_enabled**: Flag to enable whitelisting.
- **string invariant_check_epoch**: Epoch for invariant check.
- **string take_profit_borrow_interest_rate_min**: Minimum take profit borrow interest rate.
- **string funding_fee_base_rate**: Base rate for funding fee.
- **string funding_fee_max_rate**: Maximum rate for funding fee.
- **string funding_fee_min_rate**: Minimum rate for funding fee.
- **string funding_fee_collection_address**: Address for funding fee collection.
- **string swap_fee**: Swap fee.

```proto
message Params {
  option (gogoproto.goproto_stringer) = false;
  string leverage_max = 1 [
    (gogoproto.customtype) = "github.com/cosmos/cosmos-sdk/types.Dec",
    (gogoproto.nullable) = false
  ];
  string borrow_interest_rate_max = 2 [
    (gogoproto.customtype) = "github.com/cosmos/cosmos-sdk/types.Dec",
    (gogoproto.nullable) = false
  ];
  string borrow_interest_rate_min = 3 [
    (gogoproto.customtype) = "github.com/cosmos/cosmos-sdk/types.Dec",
    (gogoproto.nullable) = false
  ];
  string borrow_interest_rate_increase = 4 [
    (gogoproto.customtype) = "github.com/cosmos/cosmos-sdk/types.Dec",
    (gogoproto.nullable) = false
  ];
  string borrow_interest_rate_decrease = 5 [
    (gogoproto.customtype) = "github.com/cosmos/cosmos-sdk/types.Dec",
    (gogoproto.nullable) = false
  ];
  string health_gain_factor = 6 [
    (gogoproto.customtype) = "github.com/cosmos/cosmos-sdk/types.Dec",
    (gogoproto.nullable) = false
  ];
  int64 epoch_length = 7;
  int64 max_open_positions = 8;
  string pool_open_threshold = 9 [
    (gogoproto.customtype) = "github.com/cosmos/cosmos-sdk/types.Dec",
    (gogoproto.nullable) = false
  ];
  string force_close_fund_percentage = 10 [
    (gogoproto.customtype) = "github.com/cosmos/cosmos-sdk/types.Dec",
    (gogoproto.nullable) = false
  ];
  string force_close_fund_address = 11;
  string incremental_borrow_interest_payment_fund_percentage = 12 [
    (gogoproto.customtype) = "github.com/cosmos/cosmos-sdk/types.Dec",
    (gogoproto.nullable) = false
  ];
  string incremental_borrow_interest_payment_fund_address = 13;
  string safety_factor = 14 [
    (gogoproto.customtype) = "github.com/cosmos/cosmos-sdk/types.Dec",
    (gogoproto.nullable) = false
  ];
  bool incremental_borrow_interest_payment_enabled = 15;
  bool whitelisting_enabled = 16;
  string invariant_check_epoch = 17;
  string take_profit_borrow_interest_rate_min = 18 [
    (gogoproto.customtype) = "github.com/cosmos/cosmos-sdk/types.Dec",
    (gogoproto.nullable) = false
  ];
  string funding_fee_base_rate = 19 [
    (gogoproto.customtype) = "github.com/cosmos/cosmos-sdk/types.Dec",
    (gogoproto.nullable) = false
  ];
  string funding_fee_max_rate = 20 [
    (gogoproto.customtype) = "github.com/cosmos/cosmos-sdk/types.Dec",
    (gogoproto.nullable) = false
  ];
  string funding_fee_min_rate = 21 [
    (gogoproto.customtype) = "github.com/cosmos/cosmos-sdk/types.Dec",
    (gogoproto.nullable) = false
  ];
  string funding_fee_collection_address = 22;
  string swap_fee = 23 [
    (gogoproto.customtype) = "github.com/cosmos/cosmos-sdk/types.Dec",
    (gogoproto.nullable) = false
  ];
}
```

## Pool Definitions

The `Pool` and `PoolAsset` messages define the structure of pools and their assets in the perpetual module.

### PoolAsset

- **string liabilities**: Liabilities of the asset.
- **string custody**: Custody of the asset.
- **string take_profit_liabilities**: Take profit liabilities of the asset.
- **string take_profit_custody**: Take profit custody of the asset.
- **string asset_balance**: Balance of the asset.
- **string block_borrow_interest**: Block borrow interest.
- **string asset_denom**: Denomination of the asset.

```proto
message PoolAsset {
  string liabilities = 1 [
    (gogoproto.customtype) = "github.com/cosmos/cosmos-sdk/types.Int",
    (gogoproto.nullable) = false
  ];
  string custody = 2 [
    (gogoproto.customtype) = "github.com/cosmos/cosmos-sdk/types.Int",
    (gogoproto.nullable) = false
  ];
  string take_profit_liabilities = 3 [
    (gogoproto.customtype) = "github.com/cosmos/cosmos-sdk/types.Int",
    (gogoproto.nullable) = false
  ];
  string take_profit_custody = 4 [
    (gogoproto.customtype) = "github.com/cosmos/cosmos-sdk/types.Int",
    (gogoproto.nullable) = false
  ];
  string asset_balance = 5 [
    (gogoproto.customtype) = "github.com/cosmos/cosmos-sdk/types.Int",
    (gogoproto.nullable) = false
  ];
  string block_borrow_interest = 6 [
    (gogoproto.customtype) = "github.com/cosmos/cosmos-sdk/types.Int",
    (gogoproto.nullable) = false
  ];
  string asset_denom = 7;
}
```

### Pool

- **uint64 amm_pool_id**: AMM pool ID.
- **string health**: Health of the pool.
- **bool enabled**: Flag indicating if the pool is enabled.
- **bool closed**: Flag indicating if the pool is closed.
- **string borrow_interest_rate**: Borrow interest rate.
- **repeated PoolAsset pool_assets_long**: List of long position assets.
- **repeated PoolAsset pool_assets_short**: List of short position assets.
- **int64 last_height_borrow_interest_rate_computed**: Last height at which borrow interest rate was computed.
- **string funding_rate**: Funding rate.

```proto
message Pool {
  uint64 amm_pool_id = 1;
  string health = 2 [
    (gogoproto.customtype) = "github.com/cosmos/cosmos-sdk/types.Dec",
    (gogoproto.nullable) = false
  ];
  bool enabled = 3;
  bool closed = 4;
  string borrow_interest_rate = 5 [
    (gogoproto.customtype) = "github.com/cosmos/cosmos-sdk/types.Dec",
    (gogoproto.nullable) = false
  ];
  repeated PoolAsset pool_assets_long = 6 [(gogoproto.nullable) = false];
  repeated

 PoolAsset pool_assets_short = 7 [(gogoproto.nullable) = false];
  int64 last_height_borrow_interest_rate_computed = 8;
  string funding_rate = 9 [
    (gogoproto.customtype) = "github.com/cosmos/cosmos-sdk/types.Dec",
    (gogoproto.nullable) = false
  ];
}
```

## Query Service Definitions

The `Query` service defines the gRPC querier service for the perpetual module.

### Query Service

- **Params**: Retrieves the parameters of the module.
- **GetPositions**: Retrieves a list of positions.
- **GetPositionsByPool**: Retrieves a list of MTP positions by pool.
- **GetStatus**: Retrieves the total number of open and lifetime MTPs.
- **GetPositionsForAddress**: Retrieves a list of MTP positions for a given address.
- **GetWhitelist**: Retrieves a list of whitelisted addresses.
- **IsWhitelisted**: Checks if an address is whitelisted.
- **Pool**: Retrieves a single pool given its index.
- **Pools**: Retrieves a list of all pools.
- **MTP**: Retrieves a single MTP position given its address and ID.
- **OpenEstimation**: Provides an estimation of a new open position.

### ParamsRequest

**Request**: None

**Response**:

- **Params params**: The parameters of the module.

```proto
message ParamsRequest {}

message ParamsResponse {
  Params params = 1 [(gogoproto.nullable) = false];
}
```

### PositionsRequest

**Request**:

- **cosmos.base.query.v1beta1.PageRequest pagination**: Pagination for the request.

**Response**:

- **repeated MTP mtps**: List of MTPs.
- **cosmos.base.query.v1beta1.PageResponse pagination**: Pagination for the response.

```proto
message PositionsRequest {
  cosmos.base.query.v1beta1.PageRequest pagination = 1;
}

message PositionsResponse {
  repeated MTP mtps = 1;
  cosmos.base.query.v1beta1.PageResponse pagination = 2;
}
```

### PositionsByPoolRequest

**Request**:

- **uint64 amm_pool_id**: AMM pool ID.
- **cosmos.base.query.v1beta1.PageRequest pagination**: Pagination for the request.

**Response**:

- **repeated MTP mtps**: List of MTPs by pool.
- **cosmos.base.query.v1beta1.PageResponse pagination**: Pagination for the response.

```proto
message PositionsByPoolRequest {
  uint64 amm_pool_id = 1;
  cosmos.base.query.v1beta1.PageRequest pagination = 2;
}

message PositionsByPoolResponse {
  repeated MTP mtps = 1;
  cosmos.base.query.v1beta1.PageResponse pagination = 2;
}
```

### StatusRequest

**Request**: None

**Response**:

- **uint64 open_mtp_count**: Number of open MTPs.
- **uint64 lifetime_mtp_count**: Number of lifetime MTPs.

```proto
message StatusRequest {}

message StatusResponse {
  uint64 open_mtp_count = 1;
  uint64 lifetime_mtp_count = 2;
}
```

### PositionsForAddressRequest

**Request**:

- **string address**: Address to query positions for.
- **cosmos.base.query.v1beta1.PageRequest pagination**: Pagination for the request.

**Response**:

- **repeated MTP mtps**: List of MTPs for the address.
- **cosmos.base.query.v1beta1.PageResponse pagination**: Pagination for the response.

```proto
message PositionsForAddressRequest {
  string address = 1;
  cosmos.base.query.v1beta1.PageRequest pagination = 2;
}

message PositionsForAddressResponse {
  repeated MTP mtps = 1;
  cosmos.base.query.v1beta1.PageResponse pagination = 2;
}
```

### WhitelistRequest

**Request**:

- **cosmos.base.query.v1beta1.PageRequest pagination**: Pagination for the request.

**Response**:

- **repeated string whitelist**: List of whitelisted addresses.
- **cosmos.base.query.v1beta1.PageResponse pagination**: Pagination for the response.

```proto
message WhitelistRequest {
  cosmos.base.query.v1beta1.PageRequest pagination = 1;
}

message WhitelistResponse {
  repeated string whitelist = 1;
  cosmos.base.query.v1beta1.PageResponse pagination = 2;
}
```

### IsWhitelistedRequest

**Request**:

- **string address**: Address to check.

**Response**:

- **string address**: Address checked.
- **bool is_whitelisted**: Whether the address is whitelisted.

```proto
message IsWhitelistedRequest {
  string address = 1;
}

message IsWhitelistedResponse {
  string address = 1;
  bool is_whitelisted = 2;
}
```

### QueryGetPoolRequest

**Request**:

- **uint64 index**: Index of the pool to retrieve.

**Response**:

- **Pool pool**: The pool retrieved.

```proto
message QueryGetPoolRequest {
  uint64 index = 1;
}

message QueryGetPoolResponse {
  Pool pool = 1 [(gogoproto.nullable) = false];
}
```

### QueryAllPoolRequest

**Request**:

- **cosmos.base.query.v1beta1.PageRequest pagination**: Pagination for the request.

**Response**:

- **repeated Pool pool**: List of all pools.
- **cosmos.base.query.v1beta1.PageResponse pagination**: Pagination for the response.

```proto
message QueryAllPoolRequest {
  cosmos.base.query.v1beta1.PageRequest pagination = 1;
}

message QueryAllPoolResponse {
  repeated Pool pool = 1 [(gogoproto.nullable) = false];
  cosmos.base.query.v1beta1.PageResponse pagination = 2;
}
```

### MTPRequest

**Request**:

- **string address**: Address associated with the MTP.
- **uint64 id**: ID of the MTP.

**Response**:

- **MTP mtp**: The MTP retrieved.

```proto
message MTPRequest {
  string address = 1;
  uint64 id = 2;
}

message MTPResponse {
  MTP mtp = 1;
}
```

### QueryOpenEstimationRequest

**Request**:

- **Position position**: Position to estimate.
- **string leverage**: Leverage for the position.
- **string trading_asset**: Trading asset.
- **cosmos.base.v1beta1.Coin collateral**: Collateral.
- **string discount**: Discount for the position.
- **string take_profit_price**: Take profit price.

**Response**:

- **Position position**: Estimated position.
- **string leverage**: Leverage for the position.
- **string trading_asset**: Trading asset.
- **cosmos.base.v1beta1.Coin collateral**: Collateral.
- **cosmos.base.v1beta1.Coin min_collateral**: Minimum collateral.
- **bool valid_collateral**: Whether the collateral is valid.
- **cosmos.base.v1beta1.Coin position_size**: Position size.
- **string swap_fee**: Swap fee.
- **string discount**: Discount for the position.
- **string open_price**: Open price.
- **string take_profit_price**: Take profit price.
- **string liquidation_price**: Liquidation price.
- **string estimated_pnl**: Estimated PnL.
- **string estimated_pnl_denom**: Denomination of the estimated PnL.
- **cosmos.base.v1beta1.Coin available_liquidity**: Available liquidity.
- **string slippage**: Slippage.
- **string weight_balance_ratio**: Weight balance ratio.
- **string borrow_interest_rate**: Borrow interest rate.
- **string funding_rate**: Funding rate.
- **string price_impact**: Price impact.

```proto
message QueryOpenEstimationRequest {
  Position position = 1;
  string leverage = 2 [
    (gogoproto.customtype) = "github.com/cosmos/cosmos-sdk/types.Dec",
    (gogoproto.nullable) = false
  ];
  string trading_asset = 3;
  cosmos.base.v1beta1.Coin collateral = 4 [
    (gogoproto.nullable) = false
  ];
  string discount = 5 [
    (gogoproto.customtype) = "github.com/cosmos/cosmos-sdk/types.Dec",
    (gogoproto.nullable) = false
  ];
  string take_profit_price = 6 [
    (gogoproto.customtype) = "github.com/cosmos/cosmos-sdk/types.Dec",
    (gogoproto.nullable) = false
  ];
}

message QueryOpenEstimationResponse {
  Position position = 1;
  string leverage = 2 [
    (gogoproto.customtype) = "github.com/cosmos/cosmos-sdk/types.Dec",
    (gogoproto.nullable) = false
  ];
  string trading_asset = 3;
  cosmos.base.v1beta1.Coin collateral = 4 [
    (gogoproto.nullable) = false
  ];
  cosmos.base.v1beta1.Coin min_collateral = 5 [
    (gogoproto.nullable) = false
  ];
  bool valid_collateral = 6;
  cosmos.base.v1beta1.Coin position_size = 7 [
    (gogoproto.null

able) = false
  ];
  string swap_fee = 8 [
    (gogoproto.customtype) = "github.com/cosmos/cosmos-sdk/types.Dec",
    (gogoproto.nullable) = false
  ];
  string discount = 9 [
    (gogoproto.customtype) = "github.com/cosmos/cosmos-sdk/types.Dec",
    (gogoproto.nullable) = false
  ];
  string open_price = 10 [
    (gogoproto.customtype) = "github.com/cosmos/cosmos-sdk/types.Dec",
    (gogoproto.nullable) = false
  ];
  string take_profit_price = 11 [
    (gogoproto.customtype) = "github.com/cosmos/cosmos-sdk/types.Dec",
    (gogoproto.nullable) = false
  ];
  string liquidation_price = 12 [
    (gogoproto.customtype) = "github.com/cosmos/cosmos-sdk/types.Dec",
    (gogoproto.nullable) = false
  ];
  string estimated_pnl = 13 [
    (gogoproto.customtype) = "github.com/cosmos/cosmos-sdk/types.Int",
    (gogoproto.nullable) = false
  ];
  string estimated_pnl_denom = 14;
  cosmos.base.v1beta1.Coin available_liquidity = 15 [
    (gogoproto.nullable) = false
  ];
  string slippage = 16 [
    (gogoproto.customtype) = "github.com/cosmos/cosmos-sdk/types.Dec",
    (gogoproto.nullable) = false
  ];
  string weight_balance_ratio = 17 [
    (gogoproto.customtype) = "github.com/cosmos/cosmos-sdk/types.Dec",
    (gogoproto.nullable) = false
  ];
  string borrow_interest_rate = 18 [
    (gogoproto.customtype) = "github.com/cosmos/cosmos-sdk/types.Dec",
    (gogoproto.nullable) = false
  ];
  string funding_rate = 19 [
    (gogoproto.customtype) = "github.com/cosmos/cosmos-sdk/types.Dec",
    (gogoproto.nullable) = false
  ];
  string price_impact = 20 [
    (gogoproto.customtype) = "github.com/cosmos/cosmos-sdk/types.Dec",
    (gogoproto.nullable) = false
  ];
}
```

## Transaction Service Definitions

The `Msg` service defines the gRPC transaction service for the perpetual module.

### Msg Service

- **Open**: Opens a new position.
- **BrokerOpen**: Broker opens a new position.
- **Close**: Closes an existing position.
- **BrokerClose**: Broker closes an existing position.
- **UpdateParams**: Updates the module parameters.
- **Whitelist**: Adds an address to the whitelist.
- **Dewhitelist**: Removes an address from the whitelist.

### MsgOpen

**Request**:

- **string creator**: Creator of the message.
- **Position position**: Position to open.
- **string leverage**: Leverage for the position.
- **string trading_asset**: Trading asset.
- **cosmos.base.v1beta1.Coin collateral**: Collateral.
- **string take_profit_price**: Take profit price.

**Response**:

- **uint64 id**: ID of the opened position.

```proto
message MsgOpen {
  string creator = 1;
  Position position = 2;
  string leverage = 3 [
    (gogoproto.customtype) = "github.com/cosmos/cosmos-sdk/types.Dec",
    (gogoproto.nullable) = false
  ];
  string trading_asset = 4;
  cosmos.base.v1beta1.Coin collateral = 5 [
    (gogoproto.nullable) = false
  ];
  string take_profit_price = 6 [
    (gogoproto.customtype) = "github.com/cosmos/cosmos-sdk/types.Dec",
    (gogoproto.nullable) = false
  ];
}

message MsgOpenResponse {
  uint64 id = 1;
}
```

### MsgBrokerOpen

**Request**:

- **string creator**: Creator of the message.
- **Position position**: Position to open.
- **string leverage**: Leverage for the position.
- **string trading_asset**: Trading asset.
- **cosmos.base.v1beta1.Coin collateral**: Collateral.
- **string take_profit_price**: Take profit price.
- **string owner**: Owner of the position.

**Response**:

- **uint64 id**: ID of the opened position.

```proto
message MsgBrokerOpen {
  string creator = 1;
  Position position = 2;
  string leverage = 3 [
    (gogoproto.customtype) = "github.com/cosmos/cosmos-sdk/types.Dec",
    (gogoproto.nullable) = false
  ];
  string trading_asset = 4;
  cosmos.base.v1beta1.Coin collateral = 5 [
    (gogoproto.nullable) = false
  ];
  string take_profit_price = 6 [
    (gogoproto.customtype) = "github.com/cosmos/cosmos-sdk/types.Dec",
    (gogoproto.nullable) = false
  ];
  string owner = 7;
}

message MsgOpenResponse {
  uint64 id = 1;
}
```

### MsgClose

**Request**:

- **string creator**: Creator of the message.
- **uint64 id**: ID of the position to close.
- **string amount**: Amount to close.

**Response**:

- **uint64 id**: ID of the closed position.
- **string amount**: Closed amount.

```proto
message MsgClose {
  string creator = 1;
  uint64 id = 2;
  string amount = 3 [
    (gogoproto.customtype) = "github.com/cosmos/cosmos-sdk/types.Int",
    (gogoproto.nullable) = false
  ];
}

message MsgCloseResponse {
  uint64 id = 1;
  string amount = 2 [
    (gogoproto.customtype) = "github.com/cosmos/cosmos-sdk/types.Int",
    (gogoproto.nullable) = false
  ];
}
```

### MsgBrokerClose

**Request**:

- **string creator**: Creator of the message.
- **uint64 id**: ID of the position to close.
- **string amount**: Amount to close.
- **string owner**: Owner of the position.

**Response**:

- **uint64 id**: ID of the closed position.
- **string amount**: Closed amount.

```proto
message MsgBrokerClose {
  string creator = 1;
  uint64 id = 2;
  string amount = 3 [
    (gogoproto.customtype) = "github.com/cosmos/cosmos-sdk/types.Int",
    (gogoproto.nullable) = false
  ];
  string owner = 4;
}

message MsgCloseResponse {
  uint64 id = 1;
  string amount = 2 [
    (gogoproto.customtype) = "github.com/cosmos/cosmos-sdk/types.Int",
    (gogoproto.nullable) = false
  ];
}
```

### MsgUpdateParams

**Request**:

- **string authority**: Authority to update the parameters.
- **Params params**: New parameters.

**Response**: None

```proto
message MsgUpdateParams {
  string authority = 1;
  Params params = 2;
}

message MsgUpdateParamsResponse {}
```

### MsgWhitelist

**Request**:

- **string authority**: Authority to whitelist an address.
- **string whitelisted_address**: Address to be whitelisted.

**Response**: None

```proto
message MsgWhitelist {
  string authority = 1;
  string whitelisted_address = 2;
}

message MsgWhitelistResponse {}
```

### MsgDewhitelist

**Request**:

- **string authority**: Authority to dewhitelist an address.
- **string whitelisted_address**: Address to be dewhitelisted.

**Response**: None

```proto
message MsgDewhitelist {
  string authority = 1;
  string whitelisted_address = 2;
}

message MsgDewhitelistResponse {}
```

## Type Definitions

The `MTP` and `Position` messages define the structure of positions and MTPs in the perpetual module.

### MTP

- **string address**: Address associated with the MTP.
- **string collateral_asset**: Collateral asset.
- **string trading_asset**: Trading asset.
- **string liabilities_asset**: Liabilities asset.
- **string custody_asset**: Custody asset.
- **string collateral**: Collateral amount.
- **string liabilities**: Liabilities amount.
- **string borrow_interest_paid_collateral**: Borrow interest paid in collateral.
- **string borrow_interest_paid_custody**: Borrow interest paid in custody.
- **string borrow_interest_unpaid_collateral**: Unpaid borrow interest in collateral.
- **string custody**: Custody amount.
- **string take_profit_liabilities**: Take profit liabilities amount.
- **string take_profit_custody**: Take profit custody amount.
- **string leverage**: Leverage for the MTP.
- **string mtp_health**: Health of the MTP.
- **Position position**: Position type (LONG or SHORT).
- **uint64 id**: ID of the MTP.
- **uint64 amm_pool_id**: AMM pool ID.

- **string consolidate_leverage**: Consolidated leverage.
- **string sum_collateral**: Sum of collateral.
- **string take_profit_price**: Take profit price.
- **string take_profit_borrow_rate**: Take profit borrow rate.
- **string funding_fee_paid_collateral**: Funding fee paid in collateral.
- **string funding_fee_paid_custody**: Funding fee paid in custody.
- **string funding_fee_received_collateral**: Funding fee received in collateral.
- **string funding_fee_received_custody**: Funding fee received in custody.
- **string open_price**: Open price.

```proto
message MTP {
  string address = 1;
  string collateral_asset = 2;
  string trading_asset = 3;
  string liabilities_asset = 4;
  string custody_asset = 5;
  string collateral = 6 [
    (gogoproto.customtype) = "github.com/cosmos/cosmos-sdk/types.Int",
    (gogoproto.nullable) = false
  ];
  string liabilities = 7 [
    (gogoproto.customtype) = "github.com/cosmos/cosmos-sdk/types.Int",
    (gogoproto.nullable) = false
  ];
  string borrow_interest_paid_collateral = 8 [
    (gogoproto.customtype) = "github.com/cosmos/cosmos-sdk/types.Int",
    (gogoproto.nullable) = false
  ];
  string borrow_interest_paid_custody = 9 [
    (gogoproto.customtype) = "github.com/cosmos/cosmos-sdk/types.Int",
    (gogoproto.nullable) = false
  ];
  string borrow_interest_unpaid_collateral = 10 [
    (gogoproto.customtype) = "github.com/cosmos/cosmos-sdk/types.Int",
    (gogoproto.nullable) = false
  ];
  string custody = 11 [
    (gogoproto.customtype) = "github.com/cosmos/cosmos-sdk/types.Int",
    (gogoproto.nullable) = false
  ];
  string take_profit_liabilities = 12 [
    (gogoproto.customtype) = "github.com/cosmos/cosmos-sdk/types.Int",
    (gogoproto.nullable) = false
  ];
  string take_profit_custody = 13 [
    (gogoproto.customtype) = "github.com/cosmos/cosmos-sdk/types.Int",
    (gogoproto.nullable) = false
  ];
  string leverage = 14 [
    (gogoproto.customtype) = "github.com/cosmos/cosmos-sdk/types.Dec",
    (gogoproto.nullable) = false
  ];
  string mtp_health = 15 [
    (gogoproto.customtype) = "github.com/cosmos/cosmos-sdk/types.Dec",
    (gogoproto.nullable) = false
  ];
  Position position = 16;
  uint64 id = 17;
  uint64 amm_pool_id = 18;
  string consolidate_leverage = 19 [
    (gogoproto.customtype) = "github.com/cosmos/cosmos-sdk/types.Dec",
    (gogoproto.nullable) = false
  ];
  string sum_collateral = 20 [
    (gogoproto.customtype) = "github.com/cosmos/cosmos-sdk/types.Int",
    (gogoproto.nullable) = false
  ];
  string take_profit_price = 21 [
    (gogoproto.customtype) = "github.com/cosmos/cosmos-sdk/types.Dec",
    (gogoproto.nullable) = false
  ];
  string take_profit_borrow_rate = 22 [
    (gogoproto.customtype) = "github.com/cosmos/cosmos-sdk/types.Dec",
    (gogoproto.nullable) = false
  ];
  string funding_fee_paid_collateral = 23 [
    (gogoproto.customtype) = "github.com/cosmos/cosmos-sdk/types.Int",
    (gogoproto.nullable) = false
  ];
  string funding_fee_paid_custody = 24 [
    (gogoproto.customtype) = "github.com/cosmos/cosmos-sdk/types.Int",
    (gogoproto.nullable) = false
  ];
  string funding_fee_received_collateral = 25 [
    (gogoproto.customtype) = "github.com/cosmos/cosmos-sdk/types.Int",
    (gogoproto.nullable) = false
  ];
  string funding_fee_received_custody = 26 [
    (gogoproto.customtype) = "github.com/cosmos/cosmos-sdk/types.Int",
    (gogoproto.nullable) = false
  ];
  string open_price = 27 [
    (gogoproto.customtype) = "github.com/cosmos/cosmos-sdk/types.Dec",
    (gogoproto.nullable) = false
  ];
}
```

### Position

The `Position` enum defines the possible position types in the perpetual module.

```proto
enum Position {
  UNSPECIFIED = 0;
  LONG = 1;
  SHORT = 2;
}
```

### WhiteList

- **repeated string validator_list**: List of validators.

```proto
message WhiteList {
  repeated string validator_list = 1;
}
```
