<!--
order: 3
-->

# Oracle Pools

Oracle Pools are liquidity pools designed to provide low slippage for two assets that are expected to have a tightly correlated price ratio. The AMM (Automated Market Maker) achieves this low slippage around the expected price ratio, though price impact and slippage increase significantly as liquidity becomes more unbalanced.

This documentation details the implementation of the Solidly stableswap curve, a type of Constant Function Market Maker (CFMM) with the invariant equation:

$f(x, y) = xy(x^2 + y^2) = k$

For multiple assets, this is generalized to:

$f(a_1, ..., a_n) = a_1 \cdot ... \cdot a_n (a_1^2 + ... + a_n^2)$

## Pool Configuration

### Scaling Factor Handling

Scaling factors in Oracle pools are analogous to asset weights, affecting the valuation of assets within the pool. A governor can adjust these scaling factors, particularly when assets do not maintain a fixed ratio but instead vary predictably (e.g., non-rebasing LSTs).

By default, pools do not have a governor, as fixed ratios are sufficient for most stablecoin pools (e.g., 1:1 ratios). Scaling factors should typically be adjusted infrequently and with LPs (Liquidity Providers) being informed of the associated risks.

### Practical Considerations

1. **Precision Differences**: Different pegged coins may have varying precision levels.
2. **Token Variants**: Tokens like `TwoFoo` might need to trade around a specific multiple (e.g., `1 TwoFoo = 2 Foo`).
3. **Staking Derivatives**: Value accrual within the token can dynamically change the expected price concentration.

Scaling factors map "raw coin units" to "AMM math units" by dividing the raw coin units by the scaling factor. For example, if `Foo` has a scaling factor of $10^6$ and `WrappedFoo` has a factor of 1, then the mapping is `raw coin units / scaling factor`.

Rounding is handled by an enum `rounding mode` with three options: `RoundUp`, `RoundDown`, `RoundBankers`. These ensure precise rounding when dealing with scaled reserves.

### Implementation

#### Python Example for Scaling

```python
scaled_Foo_reserves = decimal_round(pool.Foo_liquidity / 10^6, RoundingMode)
descaled_Foo_reserves = scaled_Foo_reserves * 10^6
```

## Algorithm Details

The AMM pool interface requires the implementation of several stateful methods:

```golang
SwapOutAmtGivenIn(tokenIn sdk.Coins, tokenOutDenom string, spreadFactor osmomath.Dec) (tokenOut sdk.Coin, err error)
SwapInAmtGivenOut(tokenOut sdk.Coins, tokenInDenom string, spreadFactor osmomath.Dec) (tokenIn sdk.Coin, err error)
SpotPrice(baseAssetDenom string, quoteAssetDenom string) (osmomath.Dec, error)
JoinPool(tokensIn sdk.Coins, spreadFactor osmomath.Dec) (numShares osmomath.Int, err error)
JoinPoolNoSwap(tokensIn sdk.Coins, spreadFactor osmomath.Dec) (numShares osmomath.Int, err error)
ExitPool(numShares osmomath.Int, exitFee osmomath.Dec) (exitedCoins sdk.Coins, err error)
```

### CFMM Function

The CFMM equation is symmetric, allowing us to reorder the arguments without loss of generality. We simplify the CFMM function for specific assets by using variable substitution:

$$
\begin{equation}
    v =
    \begin{cases}
      1, & \text{if } n=2 \\
      \prod\negthinspace \negthinspace \thinspace^{n}_{i=3} \space a_i, & \text{otherwise}
    \end{cases}
  \end{equation}
$$

$$
\begin{equation}
    w =
    \begin{cases}
      0, & \text{if}\ n=2 \\
      \sum\negthinspace \negthinspace \thinspace^{n}_{i=3} \space {a_i^2}, & \text{otherwise}
    \end{cases}
  \end{equation}
$$

Thus, we define $g(x,y,v,w) = xyv(x^2 + y^2 + w) = f(x,y, a_3, ... a_n)$.

### Swaps

#### Direct Swap Solution

For a swap involving $a$ units of $x$, yielding $b$ units of $y$:

$g(x_0, y_0, v, w) = k = g(x_0 + a, y_0 - b, v, w)$

We solve for $y$ using the CFMM equation, which can be complex but is provable and can be approximated through iterative methods like binary search.

#### Iterative Search Solution

We perform a binary search to find $y$ such that $ h(x, y, w) = k' $, adjusting bounds iteratively for efficient computation.

#### Pseudocode for Binary Search

```python
def solve_y(x_0, y_0, w, x_in):
    x_f = x_0 + x_in
    err_tolerance = {"within factor of 10^-12", RoundUp}
    y_f = iterative_search(x_0, x_f, y_0, w, err_tolerance)
    y_out = y_0 - y_f
    return y_out

def iter_k_fn(x_f, y_0, w):
    def f(y_f):
        y_out = y_0 - y_f
        return -(y_out)**3 + 3 y_0 * y_out^2 - (x_f**2 + w + 3y_0**2) * y_out

def iterative_search(x_0, x_f, y_0, w, err_tolerance):
    target_k = target_k_fn(x_0, y_0, w, x_f)
    iter_k_calculator = iter_k_fn(x_f, y_0, w)
    lowerbound, upperbound = y_0, y_0
    k_ratio = bound_estimation_k0 / bound_estimation_target_k
    if k_ratio < 1:
        upperbound = ceil(y_0 / k_ratio)
    elif k_ratio > 1:
        lowerbound = 0
    else:
        return y_0
    max_iteration_count = 100
    return binary_search(lowerbound, upperbound, iter_k_calculator, target_k, err_tolerance)

def binary_search(lowerbound, upperbound, approximation_fn, target, max_iteration_count, err_tolerance):
    iter_count = 0
    cur_k_guess = 0
    while (not satisfies_bounds(cur_k_guess, target, err_tolerance)) and iter_count < max_iteration_count:
        iter_count += 1
        cur_y_guess = (lowerbound + upperbound) / 2
        cur_k_guess = approximation_fn(cur_y_guess)
        if cur_k_guess > target:
            upperbound = cur_y_guess
        else if cur_k_guess < target:
            lowerbound = cur_y_guess
    if iter_count == max_iteration_count:
        return Error("max iteration count reached")
    return cur_y_guess
```

### Spot Price

The spot price is approximated by a small swap amount $\epsilon$:

$\text{spot price} = \frac{\text{CalculateOutAmountGivenIn}(\epsilon)}{\epsilon}$

### LP Equations

#### JoinPoolNoSwap and ExitPool

These methods follow generic AMM techniques, maintaining consistency with the CFMM properties.

#### JoinPoolSingleAssetIn

For single asset joins:

```python
def JoinPoolSingleAssetIn(pool, tokenIn):
    spreadFactorApplicableFraction = 1 - (pool.ScaledLiquidityOf(tokenIn.Denom) / pool.SumOfAllScaledLiquidity())
    effectiveSpreadFactor = pool.SwapFee * spreadFactorApplicableFraction
    effectiveTokenIn = RoundDown(tokenIn * (1 - effectiveSpreadFactor))
    return BinarySearchSingleJoinLpShares(pool, effectiveTokenIn)
```

## Code Structure

The code is structured to facilitate unit testing for each pool interface method, custom message testing, and simulator integrations.

## Testing Strategy

- **Unit Tests**: For each pool interface method.
- **Msg Tests**: For custom messages like `CreatePool` and `SetScalingFactors`.
- **Simulator Integrations**: Testing pool creation, join and exit pool operations, single token join/exit, and CFMM adjustments.
- **Fuzz Testing**: To ensure the binary search algorithm and iterative approximation swap algorithm work correctly across various scales.

By following this documentation, developers can understand and implement the Solidly stableswap curve efficiently, ensuring low slippage and robust functionality for tightly correlated asset pairs.
