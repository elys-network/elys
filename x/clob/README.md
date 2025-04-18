Function: SettleMarginAndRPnL

Purpose:

This function is the core engine for updating a single trader's perpetual position state immediately after a trade involving them has occurred. It takes the trader's state before the trade (oldPerpetual), the details of the trade itself (price, quantity magnitude, buyer/seller subaccounts), and a flag (isBuyer) indicating whether the function is currently processing the buyer's or seller's side of that trade.

Its responsibilities are to:

Calculate the trader's new position Quantity based on the trade.
Calculate the new average EntryPrice if the position size increased.
Calculate the required Margin for the new/updated position based on the current market Initial Margin Ratio (IMR).
Adjust the trader's locked margin by transferring funds (SendFromSubAccount for deductions, AddToSubAccount for refunds) between their subAccount and the central market.GetAccount().
Settle any Realized Profit or Loss (RPNL) incurred by this trader if the trade reduced, closed, or flipped their position, by calling SettleRealizedPnL (which also transfers funds between the subAccount and market.GetAccount()).
Return the updatedPerpetual state reflecting all these changes.
Breakdown of Cases Handled:

The function uses the signs of the oldPerpetual.Quantity and the calculated updatedPerpetual.Quantity, along with the isBuyer flag, to determine which scenario applies to the specific trader being processed.

Case 9: Open New Position (oldPerpetual.Quantity.IsZero()):

Scenario: The trader had no position in this market before this trade.
Logic:
Calculates the required initial margin based on the trade's value (trade.Quantity * trade.Price) and the current market.InitialMarginRatio.
Deducts this margin from the opening trader's subAccount to the market account (SendFromSubAccount).
Creates and returns a brand-new Perpetual position state with:
Quantity: +trade.Quantity if isBuyer, -trade.Quantity if !isBuyer.
EntryPrice: trade.Price.
Margin: The calculated required initial margin.
Owner: The correct subaccount owner.
RPNL: None is realized.
Example: Buyer B (starts 0) buys 5 units @ $100. This function runs for B (isBuyer=true). B pays initial margin (e.g., 5 * 100 * 10% = $50) from their subaccount to the market account. B's new position is Long +5, EP $100, Margin $50.
Case 7/3: Close Position (updatedPerpetual.Quantity.IsZero()):

Scenario: The trade perfectly offsets the trader's existing position, bringing it to zero.
Logic:
Refunds the entire margin held in oldPerpetual.Margin from the market account to the trader's subAccount (AddToSubAccount).
Calls SettleRealizedPnL using the entire (signed) oldPerpetual.Quantity and trade.Price to calculate the profit or loss from the closed position. This PNL amount is transferred between the market account and the trader's subAccount.
Returns a state with zero quantity and margin (the caller should handle deleting the actual position object).
Example: Seller S (starts Long +10 @ $100, Margin $100) sells 10 units @ $105. This function runs for S (isBuyer=false). S gets $100 margin refunded from market. PNL = +10 * (105-100) = +$50 is transferred from market to S. S's position is now closed.
Case 5: Increase Long (old > 0, new > 0, isBuyer = true):

Scenario: A buyer with an existing long position buys more.
Logic:
Calculates the new weighted average EntryPrice.
Recalculates the total required margin (newRequiredInitialMargin) for the new total quantity at the new average price, using the current IMR.
Calculates the additional margin needed (diff = newRequiredInitialMargin - oldMargin).
Deducts this diff from the buyer's subAccount to the market account (SendFromSubAccount).
Updates the perpetual state with the increased quantity, new average entry price, and the new total margin held.
RPNL: None is realized.
Example: Buyer B (starts +10 @ $100, Margin $100) buys 5 units @ $110. Avg EP becomes ~$103.33. New total margin needed is ~abs(15) * 103.33 * 0.1 = ~$155. Additional margin diff = $155 - $100 = $55 is deducted from B. B's position becomes +15, EP ~$103.33, Margin $155.
Case 1: Increase Short (old < 0, new < 0, isBuyer = false):

Scenario: A seller with an existing short position sells more.
Logic:
Calculates the new weighted average EntryPrice.
Recalculates the total required margin (newRequiredInitialMargin) for the new total absolute quantity at the new average price, using the current IMR.
Calculates the additional margin needed (diff = newRequiredInitialMargin - oldMargin).
Deducts this diff from the seller's subAccount to the market account (SendFromSubAccount).
Updates the perpetual state with the more negative quantity, new average entry price, and the new total margin held.
RPNL: None is realized.
Example: Seller S (starts -10 @ $100, Margin $100) sells 5 units @ $90. Avg EP becomes ~$96.67. New total margin needed is ~abs(-15) * 96.67 * 0.1 = ~$145. Additional margin diff = $145 - $100 = $45 is deducted from S. S's position becomes -15, EP ~$96.67, Margin $145.
Case 6: Decrease Long (old > 0, new > 0, isBuyer = false):

Scenario: A seller with an existing long position sells some of it.
Logic:
Calculates the positionClosed (signed quantity, positive here, e.g., +3).
Calls SettleRealizedPnL for positionClosed to settle the PNL on the part sold (transferred between market and seller).
Recalculates the required margin (newRequiredInitialMargin) for the remaining quantity at the original entry price, using the current IMR.
Calculates the margin to refund (diff = oldMargin - newRequiredInitialMargin).
Refunds this diff from the market account to the seller's subAccount (AddToSubAccount).
Updates the perpetual state with the reduced quantity, unchanged entry price, and the new (lower) margin held.
Example: Seller S (starts +8 @ $100, Margin $80) sells 3 units @ $105. PNL = +3 * (105-100) = +$15 transferred to S. Remaining Qty = +5. New margin needed = abs(5)1000.1 = $50. Margin refund diff = $80 - $50 = $30 transferred to S. S's position becomes +5, EP $100, Margin $50.
Case 2: Decrease Short (old < 0, new < 0, isBuyer = true):

Scenario: A buyer with an existing short position buys some back.
Logic:
Calculates the positionClosed (signed quantity, negative here, e.g., -4).
Calls SettleRealizedPnL for positionClosed to settle the PNL on the part bought back (transferred between market and buyer).
Recalculates the required margin (newRequiredInitialMargin) for the remaining quantity at the original entry price, using the current IMR.
Calculates the margin to refund (diff = oldMargin - newRequiredInitialMargin).
Refunds this diff from the market account to the buyer's subAccount (AddToSubAccount).
Updates the perpetual state with the less negative quantity, unchanged entry price, and the new (lower) margin held.
Example: Buyer B (starts -10 @ $100, Margin $100) buys 4 units @ $95. PNL = -4 * (95-100) = +$20 transferred to B. Remaining Qty = -6. New margin needed = abs(-6)1000.1 = $60. Margin refund diff = $100 - $60 = $40 transferred to B. B's position becomes -6, EP $100, Margin $60.
Case 8: Flip Long -> Short (old > 0, new < 0, isBuyer = false):

Scenario: A seller with a long position sells more than their quantity, ending up short.
Logic (Close Old Long + Open New Short):
Refund the entire oldPerpetual.Margin from market to seller.
Settle PNL for closing the entire old long position at trade.Price.
Calculate margin required for the new short position at trade.Price.
Deduct this new margin from the seller to the market.
Update perpetual state: Quantity is now negative, EntryPrice is trade.Price, Margin is the newly deducted amount.
Example: Seller S (starts +8 @ $100, Margin $80) sells 10 units @ $105. Refund $80 to S. PNL for +8*(105-100)=+$40 transferred to S. New Qty = -2. New EP = $105. New margin required = abs(-2)1050.1 = $21. Deduct $21 from S. S's position becomes -2, EP $105, Margin $21.
Case 4: Flip Short -> Long (old < 0, new > 0, isBuyer = true):

Scenario: A buyer with a short position buys more than their quantity, ending up long.
Logic (Close Old Short + Open New Long):
Refund the entire oldPerpetual.Margin from market to buyer.
Settle PNL for closing the entire old short position at trade.Price.
Calculate margin required for the new long position at trade.Price.
Deduct this new margin from the buyer to the market.
Update perpetual state: Quantity is now positive, EntryPrice is trade.Price, Margin is the newly deducted amount.
Example: Buyer B (starts -10 @ $100, Margin $100) buys 12 units @ $95. Refund $100 to B. PNL for -10*(95-100)=+$50 transferred to B. New Qty = +2. New EP = $95. New margin required = abs(2)950.1 = $19. Deduct $19 from B. B's position becomes +2, EP $95, Margin $19.

Example:

Seller had long position +8, and its reduced to 5 and a new buyer (no position held before) gets +3, so this new buyer will pay RPNL ?

The new buyer does not pay the seller's RPNL directly. Here's how SettleMarginAndRPnL handles it:

Run for Seller (isBuyer=false):

Input: oldPerpetual (Qty: +8), trade (Qty: 3, Price: e.g., $105), isBuyer=false.
Output: Falls into Case 6 (Decrease Long).
Effect: Seller's position becomes +5. Seller realizes PNL on the +3 sold portion (calculated as +3 * (105 - SellerEP)). This PNL amount is transferred from the market account to the seller's subaccount. Seller also gets a partial margin refund from the market account.
Run for Buyer (isBuyer=true):

Input: oldPerpetual (Qty: 0), trade (Qty: 3, Price: $105), isBuyer=true.
Output: Falls into Case 9 (Open New Position).
Effect: Buyer's position becomes +3. Buyer realizes no PNL. Buyer pays the initial margin required for the +3 position (calculated as abs(3) * 105 * IMR). This margin is transferred from the buyer's subaccount to the market account.
The market account acts as the central counterparty for these settlements. The buyer pays margin into it; it pays out the seller's PNL and margin refund. The system ensures these balance out (along with fees, funding, etc., handled elsewhere).


Function: SettleRealizedPnL

Purpose:

This function calculates and settles the realized profit or loss (RPNL) that occurs when a specific portion of a trader's perpetual position is closed due to a trade. It handles the actual transfer of funds corresponding to this PNL between the trader's subaccount (subAccount) and the central market account (market.GetAccount()), which acts as the settlement pool. This function focuses only on the PNL settlement aspect of closing a position portion.

Inputs:

ctx sdk.Context: The standard Cosmos SDK context, providing access to state and other modules.
market types.PerpetualMarket: The market object containing details like the quote currency denomination (QuoteDenom) and the address of the market's settlement account (GetAccount()).
positionClosed math.LegacyDec: Crucially, this represents the quantity of the position that was just closed by the trade. It must carry the sign of the original position being closed:
Positive (+) if closing part of a Long position.
Negative (-) if closing part of a Short position (i.e., buying back).
subAccount types.SubAccount: The specific subaccount belonging to the trader whose RPNL is being calculated and settled. This account will either receive funds (profit) or send funds (loss).
entryPrice math.LegacyDec: The average price at which the positionClosed quantity was originally entered (the entry price of the position before this closing trade).
tradePrice math.LegacyDec: The price at which the positionClosed quantity was executed in the current trade (i.e., the exit price for this portion).
Outputs:

error: Returns nil if the PNL calculation and fund transfer (if any) were successful. Returns an error if the underlying fund transfer (AddToSubAccount or SendFromSubAccount) fails (e.g., due to insufficient funds in the source account for the transfer).
Logic Breakdown:

Calculate RPNL:

It computes the realized profit or loss using the standard formula: realizedPnlDec = positionClosed * (tradePrice - entryPrice)
The use of the signed positionClosed quantity automatically yields the correct PNL sign:
Closing Long Profit: (+) * (Exit(+) - Entry(-)) = + PNL
Closing Long Loss: (+) * (Exit(-) - Entry(+)) = - PNL
Closing Short Profit: (-) * (Exit(-) - Entry(+)) = + PNL
Closing Short Loss: (-) * (Exit(+) - Entry(-)) = - PNL
It converts the decimal result (realizedPnlDec) to an integer (realizedPnl) using TruncateInt(), suitable for creating an sdk.Coin.
Handle Zero PNL:

It checks if !realizedPnl.IsZero(). If the calculated PNL is exactly zero (e.g., exit price equals entry price), no funds need to be transferred, and the function simply returns successfully.
Settle Non-Zero PNL (Fund Transfer):

If realizedPnl.IsPositive() (Profit for Trader):
This means the trader made money on the closed portion.
It calls k.AddToSubAccount(...) to transfer the realizedPnl amount (as an sdk.Coin) from the central market.GetAccount() to the trader's subAccount. The market pool pays the profit.
Else (realizedPnl.IsNegative()) (Loss for Trader):
This means the trader lost money on the closed portion.
It calls k.SendFromSubAccount(...) to transfer the magnitude of the loss (realizedPnl.Neg(), which makes it positive) from the trader's subAccount to the central market.GetAccount(). The trader pays the loss into the market pool.
Error Propagation: If either AddToSubAccount or SendFromSubAccount returns an error during the fund transfer, this function immediately propagates that error back to the caller.

How it's Used:

This function is designed to be a helper called by SettleMarginAndRPnL. When SettleMarginAndRPnL determines that a trade causes a position decrease, full closure, or a flip, it calculates the appropriate (signed) positionClosed quantity and then invokes SettleRealizedPnL to handle the corresponding profit or loss settlement for that closed portion. SettleMarginAndRPnL separately handles the margin adjustments (refunds/deductions).

Illustrative Examples:

Profit on Long Close: Closing +2 Long @ $110 (EP $100). positionClosed=+2, tradePrice=110, entryPrice=100. realizedPnl = +2 * (110-100) = +20. AddToSubAccount sends 20 from market to trader.
Loss on Short Close: Closing -5 Short @ $105 (EP $100). positionClosed=-5, tradePrice=105, entryPrice=100. realizedPnl = -5 * (105-100) = -50. SendFromSubAccount sends 50 (-50.Neg()) from trader to market.