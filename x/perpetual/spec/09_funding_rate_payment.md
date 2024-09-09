<!--
order: 9
-->

# Funding Rate Calculation and Payment

The funding rate is a critical mechanism in the perpetual contracts system designed to incentivize traders to take positions that counterbalance the market trend. This helps to maintain a balanced market, encouraging traders to engage in positions where the liabilities of long and short traders are not too heavily skewed in one direction.

## 1. **Overview of Funding Rate Calculation**

The funding rate is computed at each block using the total long and short liabilities, which form what is known as the **Open Interest (OI)**.

- **Long traders pay short traders** if the long liabilities exceed short liabilities.
- **Short traders pay long traders** if the short liabilities exceed long liabilities.

This system is designed to **incentivize traders to take on risk against the prevailing market trend**. For example, if the majority of traders are long, the funding rate will encourage some traders to open short positions by requiring long traders to pay a funding fee to short traders.

### 2. **Key Parameters Governing the Funding Rate**

The funding rate is influenced by several module governance parameters:

- **Base Funding Fee Rate (`funding_fee_base_rate`)**:  
  The base rate used to calculate the funding rate.  
  Example value: `"0.000300000000000000"`

- **Funding Fee Collection Address (`funding_fee_collection_address`)**:  
  The address where collected funding fees are deposited.  
  Example address: `elys1qqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqnrec2l`

- **Maximum Funding Fee Rate (`funding_fee_max_rate`)**:  
  The upper cap on the funding fee rate to avoid excessive payments.  
  Example value: `"0.001000000000000000"`

- **Minimum Funding Fee Rate (`funding_fee_min_rate`)**:  
  The lower limit on the funding fee rate to prevent negative rate extremes.  
  Example value: `"-0.001000000000000000"`

### 3. **Funding Rate Formula**

The funding rate is calculated as follows:

- **When Long Liabilities Exceed Short Liabilities:**

  ```
  funding_rate = min(max(base_rate * (long_liabilities / short_liabilities), min_rate), max_rate)
  ```

- **When Short Liabilities Exceed Long Liabilities:**
  ```
  funding_rate = min(max(base_rate * (short_liabilities / long_liabilities) * -1, min_rate), max_rate)
  ```

This formula ensures that the funding rate stays within defined boundaries, incentivizing traders to take positions in opposition to the prevailing market trend, thereby promoting liquidity and market balance.

### 4. **Funding Fee Collection**

At every block, the calculated funding fees are collected and sent to the **Funding Fee Collection Address**. The details of this process can be found in the source code:  
[Funding Fee Collection Code](https://github.com/elys-network/elys/blob/main/x/perpetual/keeper/handle_funding_fee_collection.go).

The collected funding fees act as the source for redistribution to the traders, depending on the market's funding rate outcome.

### 5. **Funding Fee Distribution**

Once collected, the funding fees are redistributed as follows:

- If the funding rate is positive, **long traders pay short traders**.
- If the funding rate is negative, **short traders pay long traders**.

The goal of this redistribution is to create financial incentives that encourage traders to take positions that go against the dominant market direction, thus helping to correct imbalances in the open interest.

For more technical details on how the redistribution process works, refer to the following code reference:  
[Funding Fee Distribution Code](https://github.com/elys-network/elys/blob/main/x/perpetual/keeper/handle_funding_fee_distribution.go).
