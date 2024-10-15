<!--
order: 9
-->

# Funding Calculation and Payment

The funding rate is a critical mechanism in the perpetual contracts system designed to incentivize traders to take positions that counterbalance the market trend. This helps to maintain a balanced market, encouraging traders to engage in positions where the liabilities of long and short traders are not too heavily skewed in one direction.

## 1. **Overview of Funding Rate Calculation**

The funding rate is computed at each block using the total long and short liabilities, which form what is known as the **Open Interest (OI)**.

- **Long traders pay short traders** if the long liabilities exceed short liabilities.
- **Short traders pay long traders** if the short liabilities exceed long liabilities.

This system is designed to **incentivize traders to take on risk against the prevailing market trend**. For example, if the majority of traders are long, the funding rate will encourage some traders to open short positions by requiring long traders to pay a funding fee to short traders.

### 2. **Key Parameter Governing the Funding Rate**

The funding rate is influenced by several module governance parameters:

- **Fixed Funding Fee Rate (`fixed_funding_fee_rate`)**:  
  The base rate used to calculate the funding rate.  
  Example value: `"0.300000000000000"`

### 3. **Funding Rate Formula**

The funding rate is calculated as follows:

- **When Long Liabilities Exceed Short Liabilities:**

  ```
  funding_rate = fixed_rate * (long_custody - short_liability) / (long_custody + short_liability)
  ```

- **When Short Liabilities Exceed Long Liabilities:**
  ```
 funding_rate = fixed_rate * (short_liability - long_custody) / (long_custody + short_liability)
  ```

This formula ensures that the funding rate stays within defined boundaries, incentivizing traders to take positions in opposition to the prevailing market trend, thereby promoting liquidity and market balance.

### 4. **Funding Fee Collection**

At each block, the system tracks the calculated funding fees. If a user performs an interaction, such as closing a position or consolidating, the system triggers the collection of the accumulated fees and The details of this process can be found in the source code:  
[Funding Fee Collection Code](https://github.com/elys-network/elys/blob/main/x/perpetual/keeper/settle_funding_fee_collection.go).

The collected funding fees act as the source for redistribution to the traders, depending on the market's funding rate outcome.

### 5. **Funding Fee Distribution**

The distribution of funding fees is based on the calculated funding rate and serves to reward traders who are taking positions opposite the market trend.

- **Funding Rate Direction**:  
  The funding rate determines **which side of the market receives the collected funding payments**.
  - **Long custody is greater than short liabilities**: Short traders receive payments.
  - **Short liability is greater than long custody**: Long traders receive payments.

However, the **actual distribution amount** is not directly proportional to the funding rate alone. Instead, the amount that each position receives is calculated based on the size of that position relative to the total liabilities/custody on the side that is receiving the funding payments.

The formula for determining the amount distributed to each position is as follows:
  ```
  // For short positions

  payment_to_position = (position_liabilities / pool_liabilities) * total_funding_collected
  ```
  ```
  // For long positions

  payment_to_position = (position_custody / pool_custody) * total_funding_collected
  ```

Where:

- **`position_liabilities`**: The liabilities (position size) of the individual trader receiving the funding payment.
- **`pool_liabilities`**: The total liabilities on the side of the market receiving the funding payment (long or short, depending on the funding rate direction).
- **`position_custody`**: The custody (position size) of the individual trader receiving the funding payment.
- **`pool_custody`**: The total custody on the side of the market receiving the funding payment (long or short, depending on the funding rate direction).
- **`total_funding_collected`**: The total amount of funding fees collected during the block.

#### Example:

- If a trader holds 10% of the total liabilities on the receiving side, they will receive 10% of the total collected funding fees.

For more details on how the funding fees are distributed, refer to the following code reference:  
[Funding Fee Distribution Code](https://github.com/elys-network/elys/blob/main/x/perpetual/keeper/settle_funding_fee_distribution.go).
