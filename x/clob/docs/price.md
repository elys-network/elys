Function: SetTwapPrices

Purpose:

This function is responsible for recording and updating price and volume data based on individual trades within the context of a specific block height. Its primary goal is to maintain data points that can later be used to calculate a TWAP. It ensures that for each block where trading occurs, there's a record reflecting the volume-weighted average price within that block and also updates a running cumulative price integral used for efficient TWAP calculation across blocks. It also handles pruning old data.

Inputs:

ctx sdk.Context: Provides access to block height, block time, and KV stores.
trade types.Trade: Contains details of the trade that just occurred, including MarketId, Quantity (magnitude), and Price.
Outputs/Side Effects:

None (void function, returns no value).
Modifies State: Creates or updates records in the persistent KVStore (k.storeService) under a key derived from MarketId and BlockHeight. It stores a types.TwapPrice object.
Modifies Transient State: Mirrors the persistent store update in the transient store (k.transientStoreService) for potentially faster access within the same block.
Deletes State: Prunes old TwapPrice records from the persistent store based on a configured time window (market.QuoteDenomTwapPricesWindow).
Can Panic: Will panic if certain conditions are met (e.g., market not found, timestamp errors, unmarshalling fails).
Logic Breakdown:

Get Stores & Key: Accesses the persistent and transient KV stores and creates a unique key for the current trade.MarketId and ctx.BlockHeight().
Check for Existing Record (Current Block): It attempts to retrieve a TwapPrice record using the key for the current block.
If a record exists:
It means other trades have already happened in this block.
It updates the record by calculating a new AverageTradePrice for the block. This is done by taking the previous total value within the block (oldAvgPrice * oldVolume), adding the value of the current trade (tradePrice * tradeQty), and dividing by the new total volume for the block (oldVolume + tradeQty).
It updates the TotalVolume for the block by adding the absolute quantity of the current trade.
It saves the updated record back to both persistent and transient stores.
Note: The CumulativePrice is not updated when just adding trades within the same block, as it only changes at the boundary between blocks.
If no record exists:
This is the first trade recorded for this market in the current block.
It sets the record's AverageTradePrice to the current trade.Price and TotalVolume to the trade.Quantity.
Cumulative Price Calculation:
It finds the most recent previous TwapPrice record (from a prior block) using a reverse iterator.
If no previous record exists (e.g., first ever trade for the market), the CumulativePrice starts at zero.
If a previous record (lastTwapPrice) exists, it calculates the time elapsed since that record (timeDelta = currentTimestamp - lastTimestamp).
It adds the price contribution from the previous interval (lastTwapPrice.AverageTradePrice * timeDelta) to the previous cumulative total (lastTwapPrice.CumulativePrice) to get the new currentTwapPrice.CumulativePrice. This effectively integrates the price over time.
It performs a sanity check panic if timeDelta <= 0.
It saves the newly created record (with updated cumulative price) to persistent and transient stores.
Pruning: It then iterates forward from the beginning of the TWAP records for this market and deletes any record whose timestamp is older than the allowed window (currentTimestamp - market.QuoteDenomTwapPricesWindow), ensuring the stored data doesn't grow indefinitely.
Function: GetCurrentTwapPrice

Purpose:

This function calculates the Time-Weighted Average Price (TWAP) for a given market over the time window implicitly defined by the oldest available (non-pruned) TWAP record and the most recent TWAP record stored by SetTwapPrices. It leverages the pre-calculated cumulative price values for efficiency.

Inputs:

ctx sdk.Context: Provides access to the KV store.
marketId uint64: The ID of the market for which to calculate the TWAP.
Outputs:

math.LegacyDec: The calculated TWAP value. Returns zero if no data exists or the time delta is zero.
Can Panic: Will panic if the timestamp of the last record is somehow less than or equal to the timestamp of the first record (indicating corrupted data or clock issues).
Logic Breakdown:

Find Latest Record: Uses a reverse iterator on the market's TWAP records to find the most recent entry (lastTwapPrice). If no records exist, it initializes a default lastTwapPrice struct with current block/time and zero values.
Find Earliest Record: Uses a forward iterator to find the oldest available entry (firstTwapPrice) for the market. This relies on the pruning mechanism in SetTwapPrices having removed records older than the desired TWAP window. If no records exist, it initializes a default firstTwapPrice.
Calculate Deltas:
Calculates the difference in cumulative price: num = lastTwapPrice.CumulativePrice - firstTwapPrice.CumulativePrice. This represents the total "price * time" accumulated over the interval between the first and last record. (Mathematically: ∫ P(t) dt from t_first to t_last).
Calculates the time difference: timeDelta = lastTwapPrice.Timestamp - firstTwapPrice.Timestamp.
Handle Edge Cases:
If the cumulative price delta (num) is zero (e.g., price was always zero or no trades occurred), it returns math.LegacyZeroDec().
If the time delta is zero or negative (which indicates a data inconsistency or clock issue), it panics.
Calculate TWAP:
Returns the final TWAP by dividing the cumulative price delta by the time delta: num.Quo(timeDelta). (Mathematically: [∫ P(t) dt] / Δt).