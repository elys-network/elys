# Margin module

Margin module is to support margin trading on molecule token pools.
At initial, isolated margin is to be supported and cross margin will be supported later.

## Pool

Margin will use amm pool and will have position records. When the trader closes the position or liquidated, pool tokens will be transferred to the margin trader for P&L.

## Margin limit

Based on pool size margin should be limited.

## Race condition between amm & margin

Pool could have lack of balance. Therefore why we have to keep a healthy buffer when setting margin position. We should not allow more than 50% of the pool to be borrowed for margin.

The position should be auto-closed before the pool becomes insufficient to cover.

Any action which affects the health of the position shud trigger position closing.

## Risks

We are not going to offer margin on shallow pools.

## Oracle

Oracle should use average price to prevent too big liquidations sudden dump for a short time.

If we have 3 oracle price sources for example and one experiences a massive candle anomaly, We will need to add exceptions for this case.

## Reference codebases for margin

- Get reference code from Caner

## Notes

- Margin code to be the base code for LP leveraging - Only real difference being that for margin the borrow is from the pool liqudity itself while for LP leveraging the borrow is from USDC deposit

- Ultimately when we have cross margin, that’s when margin positions and LP positions can interact
