# leveragelp module

leveragelp module is to support leveragelp trading on molecule token pools.
At initial, isolated leveragelp is to be supported and cross leveragelp will be supported later.

## Pool

leveragelp will use amm pool and will have position records. When the trader closes the position or liquidated, pool tokens will be transferred to the leveragelp trader for P&L.

## leveragelp limit

Based on pool size leveragelp should be limited.

## Race condition between amm & leveragelp

Pool could have lack of balance. Therefore why we have to keep a healthy buffer when setting leveragelp position. We should not allow more than 50% of the pool to be borrowed for leveragelp.

The position should be auto-closed before the pool becomes insufficient to cover.

Any action which affects the health of the position shud trigger position closing.

## Risks

We are not going to offer leveragelp on shallow pools.

## Oracle

Oracle should use average price to prevent too big liquidations sudden dump for a short time.

If we have 3 oracle price sources for example and one experiences a massive candle anomaly, We will need to add exceptions for this case.

## Reference codebases for leveragelp

- TBD

## Notes

- leveragelp code to be the base code for LP leveraging - Only real difference being that for leveragelp the borrow is from the pool liqudity itself while for LP leveraging the borrow is from base currency deposit

- Ultimately when we have cross leveragelp, that’s when leveragelp positions and LP positions can interact
