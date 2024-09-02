<!--
order: 6
-->

# Borrow Interest Calculation

### 1. **Overview**

This document specifies the formula and components used to calculate the `Borrow Interest` in a decentralized financial (DeFi) platform. The `Borrow Interest` is the amount that a user (or Margin Trading Position, MTP) owes on borrowed funds, determined by the borrowed amount, applicable interest rate, and other relevant factors. The **Take Profit Borrow Rate** is a number between 0 and 1, and its effect on the final borrow interest amount is to reduce the calculated interest.

### 2. **Formula for Borrow Interest**

The `Borrow Interest` is calculated using the following formula:

$$
\text{Borrow Interest} = \left\lceil \left( \text{Liabilities} + \text{Unpaid Collateral} \right) \times \text{Borrow Interest Rate} \times \frac{\text{Epoch Position}}{\text{Epoch Length}} \right\rceil \times \text{Take Profit Borrow Rate}
$$

### 3. **Component Descriptions**

- **\*Liabilities**: The `Liabilities` represent the total amount of funds that the user has borrowed in the base currency. This value is stored in the user's Margin Trading Position (MTP) and is updated as the user borrows or repays funds.
- **Unpaid Collateral**: The `Unpaid Collateral` refers to any interest amount that has accrued but not yet been paid. If the collateral is not in the base currency, it is converted to the base currency using an estimated swap value.
- **Borrow Interest Rate**: The `Borrow Interest Rate` is the rate at which interest is accrued on the borrowed amount. It is expressed as a decimal. This rate is dynamically calculated and updated based on pool health.
- **Epoch Position**: The `Epoch Position` refers to the current position within the epoch, which is a predefined time period over which interest accrual is calculated. The position helps prorate the interest if it’s being calculated mid-epoch. It is calculated based on the current time or block height relative to the start of the epoch.
- **Epoch Length**: The `Epoch Length` is the total duration of an epoch. This is used in conjunction with the `Epoch Position` to determine the proportion of the epoch that has elapsed. This is a fixed value set by the platform, representing the length of time or number of blocks in one epoch.
- **Take Profit Borrow Rate**: The `Take Profit Borrow Rate` is a multiplier between 0 and 1 that adjusts the final borrow interest amount downward. It is used to apply user reductions to the calculated interest. This rate is calculated based on the position take profit price. Higher is the take profit price and higher is the borrow interest rate, conversely lower is the take profit price and lower is the borrow interest rate. For example, a rate of 0.9 would reduce the interest to 90% of the initially calculated value.

### 4. **Detailed Calculation Steps**

1. **Initialize Variables:** Convert the `Borrow Interest Rate` to a rational number for precision, and prepare the `Liabilities` and `Unpaid Collateral` as rational numbers.

2. **Calculate Liabilities with Unpaid Collateral:**

   $$
   \text{Total Liabilities} = \text{Liabilities} + \text{Unpaid Collateral}
   $$

3. **Compute Initial Borrow Interest:**

   $$
   \text{Borrow Interest Initial} = \text{Total Liabilities} \times \text{Borrow Interest Rate}
   $$

4. **Prorate Borrow Interest Based on Epoch Position (if applicable):**

   $$
   \text{Prorated Borrow Interest} = \text{Borrow Interest Initial} \times \frac{\text{Epoch Position}}{\text{Epoch Length}}
   $$

5. **Final Borrow Interest Calculation:**

   $$
   \text{Borrow Interest} = \left\lceil \text{Prorated Borrow Interest} \right\rceil \times \text{Take Profit Borrow Rate}
   $$

   The **Take Profit Borrow Rate** is used as a reduction factor, where a value between 0 and 1 is applied to decrease the calculated borrow interest. For example, a rate of 0.9 would reduce the interest to 90% of the initially calculated value.

6. **Apply Minimum Interest Rule:** If the calculated interest is too low, ensure that a minimum of 1 unit is charged if the `Borrow Interest Rate` is not zero.

### 5. **Example Calculation**

Suppose:

- `Liabilities = 1000`
- `Unpaid Collateral = 50`
- `Borrow Interest Rate = 0.05`
- `Epoch Position = 5`
- `Epoch Length = 10`
- `Take Profit Borrow Rate = 0.9`

**Steps:**

1. Total Liabilities = 1000 + 50 = 1050
2. Borrow Interest Initial = 1050 \* 0.05 = 52.5
3. Prorated Borrow Interest = 52.5 \* (5/10) = 26.25
4. Final Borrow Interest = ⌈26.25⌉ \* 0.9 = 24

The final `Borrow Interest` would be **24 units**, reduced by the **Take Profit Borrow Rate** of 0.9.

### 6. **Error Handling**

- **Zero Interest Rate:** If the `Borrow Interest Rate` or `Liabilities` is zero, the function should return zero interest.
- **Invalid Inputs:** Ensure that all inputs (e.g., `Liabilities`, `Epoch Length`) are positive and non-zero, except where explicitly allowed (e.g., `Epoch Position` could be zero).
