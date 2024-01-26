# Perpetual module

Perpetual module is to support perpetual trading on molecule token pools.
At initial, isolated perpetual is to be supported and cross perpetual will be supported later.

## Pool

Perpetual will use amm pool and will have position records. When the trader closes the position or liquidated, pool tokens will be transferred to the perpetual trader for P&L.

## Perpetual limit

Based on pool size perpetual should be limited.

## Race condition between amm & perpetual

Pool could have lack of balance. Therefore why we have to keep a healthy buffer when setting perpetual position. We should not allow more than 50% of the pool to be borrowed for perpetual.

The position should be auto-closed before the pool becomes insufficient to cover.

Any action which affects the health of the position shud trigger position closing.

## Risks

We are not going to offer perpetual on shallow pools.

## Oracle

Oracle should use average price to prevent too big liquidations sudden dump for a short time.

If we have 3 oracle price sources for example and one experiences a massive candle anomaly, We will need to add exceptions for this case.

## Reference codebases for perpetual

- TBD

## Notes

- Perpetual code to be the base code for LP leveraging - Only real difference being that for perpetual the borrow is from the pool liqudity itself while for LP leveraging the borrow is from base currency deposit

- Ultimately when we have cross perpetual, that’s when perpetual positions and LP positions can interact
