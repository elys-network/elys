package types

import (
	"errors"
	fmt "fmt"
	"sort"
	"strings"
	"time"

	errorsmod "cosmossdk.io/errors"
	"cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func (p *Pool) addToPoolAssetBalances(coins sdk.Coins) error {
	for _, coin := range coins {
		i, poolAsset, err := p.GetPoolAssetAndIndex(coin.Denom)
		if err != nil {
			return err
		}
		poolAsset.Token.Amount = poolAsset.Token.Amount.Add(coin.Amount)
		p.PoolAssets[i] = poolAsset
	}
	return nil
}

func (p Pool) parsePoolAssets(tokensA sdk.Coins, tokenBDenom string) (
	tokenA sdk.Coin, Aasset PoolAsset, Basset PoolAsset, err error,
) {
	if len(tokensA) != 1 {
		return tokenA, Aasset, Basset, errors.New("expected tokensB to be of length one")
	}
	Aasset, Basset, err = p.parsePoolAssetsByDenoms(tokensA[0].Denom, tokenBDenom)
	if err != nil {
		return sdk.Coin{}, PoolAsset{}, PoolAsset{}, err
	}
	return tokensA[0], Aasset, Basset, nil
}

func (p Pool) parsePoolAssetsCoins(tokensA sdk.Coins, tokensB sdk.Coins) (
	Aasset PoolAsset, Basset PoolAsset, err error,
) {
	if len(tokensB) != 1 {
		return Aasset, Basset, errors.New("expected tokensA to be of length one")
	}
	_, Aasset, Basset, err = p.parsePoolAssets(tokensA, tokensB[0].Denom)
	return Aasset, Basset, err
}

func (p Pool) parsePoolAssetsByDenoms(tokenADenom, tokenBDenom string) (
	Aasset PoolAsset, Basset PoolAsset, err error,
) {
	Aasset, found1 := GetPoolAssetByDenom(p.PoolAssets, tokenADenom)
	Basset, found2 := GetPoolAssetByDenom(p.PoolAssets, tokenBDenom)

	if !found1 {
		return PoolAsset{}, PoolAsset{}, fmt.Errorf("(%s) does not exist in the pool", tokenADenom)
	}
	if !found2 {
		return PoolAsset{}, PoolAsset{}, fmt.Errorf("(%s) does not exist in the pool", tokenBDenom)
	}
	return Aasset, Basset, nil
}

// setInitialPoolParams
func (p *Pool) setInitialPoolParams(params PoolParams, sortedAssets []PoolAsset, curBlockTime time.Time) error {
	p.PoolParams = params

	return nil
}

func (p *Pool) applySwap(ctx sdk.Context, tokensIn sdk.Coins, tokensOut sdk.Coins, swapFeeIn, swapFeeOut sdk.Dec, accPoolKeeper AccountedPoolKeeper) error {
	// Fixed gas consumption per swap to prevent spam
	ctx.GasMeter().ConsumeGas(BalancerGasFeeForSwap, "balancer swap computation")
	// Also ensures that len(tokensIn) = 1 = len(tokensOut)
	inPoolAsset, outPoolAsset, err := p.parsePoolAssetsCoins(tokensIn, tokensOut)
	if err != nil {
		return err
	}
	inTokensAfterFee := sdk.NewDecFromInt(tokensIn[0].Amount).Mul(sdk.OneDec().Sub(swapFeeIn)).TruncateInt()
	outTokensAfterFee := sdk.NewDecFromInt(tokensOut[0].Amount).Mul(sdk.OneDec().Sub(swapFeeOut)).TruncateInt()
	inPoolAsset.Token.Amount = inPoolAsset.Token.Amount.Add(inTokensAfterFee)
	outPoolAsset.Token.Amount = outPoolAsset.Token.Amount.Sub(outTokensAfterFee)

	return p.UpdatePoolAssetBalances(sdk.NewCoins(
		inPoolAsset.Token,
		outPoolAsset.Token,
	))
}

// SetInitialPoolAssets sets the PoolAssets in the pool. It is only designed to
// be called at the pool's creation. If the same denom's PoolAsset exists, will
// return error.
//
// The list of PoolAssets must be sorted. This is done to enable fast searching
// for a PoolAsset by denomination.
// TODO: Unify story for validation of []PoolAsset, some is here, some is in
// CreatePool.ValidateBasic()
func (p *Pool) SetInitialPoolAssets(PoolAssets []PoolAsset) error {
	exists := make(map[string]bool)
	for _, asset := range p.PoolAssets {
		exists[asset.Token.Denom] = true
	}

	newTotalWeight := p.TotalWeight
	scaledPoolAssets := make([]PoolAsset, 0, len(PoolAssets))

	// TODO: Refactor this into PoolAsset.validate()
	for _, asset := range PoolAssets {
		if asset.Token.Amount.LTE(sdk.ZeroInt()) {
			return fmt.Errorf("can't add the zero or negative balance of token")
		}

		err := asset.validateWeight()
		if err != nil {
			return err
		}

		if exists[asset.Token.Denom] {
			return fmt.Errorf("same PoolAsset already exists")
		}
		exists[asset.Token.Denom] = true

		// Scale weight from the user provided weight to the correct internal weight
		asset.Weight = asset.Weight.MulRaw(GuaranteedWeightPrecision)
		scaledPoolAssets = append(scaledPoolAssets, asset)
		newTotalWeight = newTotalWeight.Add(asset.Weight)
	}

	// TODO: Change this to a more efficient sorted insert algorithm.
	// Furthermore, consider changing the underlying data type to allow in-place modification if the
	// number of PoolAssets is expected to be large.
	p.PoolAssets = append(p.PoolAssets, scaledPoolAssets...)

	sortPoolAssetsByDenom(p.PoolAssets)

	p.TotalWeight = newTotalWeight

	return nil
}

func (p *Pool) AddTotalShares(amt math.Int) {
	p.TotalShares.Amount = p.TotalShares.Amount.Add(amt)
}

func (p *Pool) IncreaseLiquidity(sharesOut math.Int, coinsIn sdk.Coins) {
	err := p.addToPoolAssetBalances(coinsIn)
	if err != nil {
		panic(err)
	}
	p.AddTotalShares(sharesOut)
}

func (p *Pool) UpdatePoolAssetBalance(coin sdk.Coin) error {
	// Check that PoolAsset exists.
	assetIndex, existingAsset, err := p.GetPoolAssetAndIndex(coin.Denom)
	if err != nil {
		return err
	}

	if coin.Amount.LTE(sdk.ZeroInt()) {
		return fmt.Errorf("can't set the pool's balance of a token to be zero or negative")
	}

	// Update the supply of the asset
	existingAsset.Token = coin
	p.PoolAssets[assetIndex] = existingAsset
	return nil
}

func (p *Pool) UpdatePoolAssetBalances(coins sdk.Coins) error {
	// Ensures that there are no duplicate denoms, all denom's are valid,
	// and amount is > 0
	err := coins.Validate()
	if err != nil {
		return fmt.Errorf("provided coins are invalid, %v", err)
	}

	for _, coin := range coins {
		// TODO: We may be able to make this log(|coins|) faster in how it
		// looks up denom -> Coin by doing a multi-search,
		// but as we don't anticipate |coins| to be large, we omit this.
		err = p.UpdatePoolAssetBalance(coin)
		if err != nil {
			return err
		}
	}

	return nil
}

func (p Pool) GetAllPoolAssets() []PoolAsset {
	copyslice := make([]PoolAsset, len(p.PoolAssets))
	copy(copyslice, p.PoolAssets)
	return copyslice
}

func (p Pool) GetTotalPoolLiquidity() sdk.Coins {
	return poolAssetsCoins(p.PoolAssets)
}

// Returns a pool asset, and its index. If err != nil, then the index will be valid.
func (p Pool) GetPoolAssetAndIndex(denom string) (int, PoolAsset, error) {
	if denom == "" {
		return -1, PoolAsset{}, fmt.Errorf("you tried to find the PoolAsset with empty denom")
	}

	if len(p.PoolAssets) == 0 {
		return -1, PoolAsset{}, errorsmod.Wrapf(ErrDenomNotFoundInPool, fmt.Sprintf(FormatNoPoolAssetFoundErrFormat, denom))
	}

	i := sort.Search(len(p.PoolAssets), func(i int) bool {
		PoolAssetA := p.PoolAssets[i]

		compare := strings.Compare(PoolAssetA.Token.Denom, denom)
		return compare >= 0
	})

	if i < 0 || i >= len(p.PoolAssets) {
		return -1, PoolAsset{}, errorsmod.Wrapf(ErrDenomNotFoundInPool, fmt.Sprintf(FormatNoPoolAssetFoundErrFormat, denom))
	}

	if p.PoolAssets[i].Token.Denom != denom {
		return -1, PoolAsset{}, errorsmod.Wrapf(ErrDenomNotFoundInPool, fmt.Sprintf(FormatNoPoolAssetFoundErrFormat, denom))
	}

	return i, p.PoolAssets[i], nil
}

// getMaximalNoSwapLPAmount returns the coins(lp liquidity) needed to get the specified amount of shares in the pool.
// Steps to getting the needed lp liquidity coins needed for the share of the pools are
// 1. calculate how much percent of the pool does given share account for(# of input shares / # of current total shares)
// 2. since we know how much % of the pool we want, iterate through all pool liquidity to calculate how much coins we need for
// each pool asset.
func (pool Pool) GetMaximalNoSwapLPAmount(shareOutAmount math.Int) (neededLpLiquidity sdk.Coins, err error) {
	totalSharesAmount := pool.GetTotalShares()
	// shareRatio is the desired number of shares, divided by the total number of
	// shares currently in the pool. It is intended to be used in scenarios where you want
	shareRatio := sdk.NewDecFromBigInt(shareOutAmount.BigInt()).QuoInt(totalSharesAmount.Amount)
	if shareRatio.LTE(sdk.ZeroDec()) {
		return sdk.Coins{}, errorsmod.Wrapf(ErrInvalidMathApprox, "Too few shares out wanted. "+
			"(debug: getMaximalNoSwapLPAmount share ratio is zero or negative)")
	}

	poolLiquidity := pool.GetTotalPoolLiquidity()
	neededLpLiquidity = sdk.Coins{}

	for _, coin := range poolLiquidity {
		// (coin.Amt * shareRatio).Ceil()
		neededAmt := sdk.NewDecFromBigInt(coin.Amount.BigInt()).Mul(shareRatio).Ceil().RoundInt()
		if neededAmt.LTE(sdk.ZeroInt()) {
			return sdk.Coins{}, errorsmod.Wrapf(ErrInvalidMathApprox, "Too few shares out wanted")
		}
		neededCoin := sdk.Coin{Denom: coin.Denom, Amount: neededAmt}
		neededLpLiquidity = neededLpLiquidity.Add(neededCoin)
	}
	return neededLpLiquidity, nil
}

func (p *Pool) CalcExitPoolCoinsFromShares(
	ctx sdk.Context,
	oracleKeeper OracleKeeper,
	accountedPoolKeeper AccountedPoolKeeper,
	exitingShares math.Int,
	tokenOutDenom string,
) (exitedCoins sdk.Coins, weightBalanceBonus math.LegacyDec, err error) {
	return CalcExitPool(ctx, oracleKeeper, *p, accountedPoolKeeper, exitingShares, tokenOutDenom)
}

func (p *Pool) TVL(ctx sdk.Context, oracleKeeper OracleKeeper) (sdk.Dec, error) {
	// OracleAssetsTVL * TotalWeight / OracleAssetsWeight
	// E.g. JUNO / USDT / USDC (30:30:30)
	// TVL = USDC_USDT_liquidity * 90 / 60

	oracleAssetsTVL := sdk.ZeroDec()
	totalWeight := sdk.ZeroInt()
	oracleAssetsWeight := sdk.ZeroInt()
	for _, asset := range p.PoolAssets {
		tokenPrice := oracleKeeper.GetAssetPriceFromDenom(ctx, asset.Token.Denom)
		totalWeight = totalWeight.Add(asset.Weight)
		if tokenPrice.IsZero() {
			if p.PoolParams.UseOracle {
				return sdk.ZeroDec(), fmt.Errorf("token price not set: %s", asset.Token.Denom)
			}
		} else {
			v := tokenPrice.Mul(sdk.NewDecFromInt(asset.Token.Amount))
			oracleAssetsTVL = oracleAssetsTVL.Add(v)
			oracleAssetsWeight = oracleAssetsWeight.Add(asset.Weight)
		}
	}

	if oracleAssetsWeight.IsZero() {
		return sdk.ZeroDec(), nil
	}

	return oracleAssetsTVL.Mul(sdk.NewDecFromInt(totalWeight)).Quo(sdk.NewDecFromInt(oracleAssetsWeight)), nil
}

func (p *Pool) LpTokenPrice(ctx sdk.Context, oracleKeeper OracleKeeper) (sdk.Dec, error) {
	ammPoolTvl, err := p.TVL(ctx, oracleKeeper)
	if err != nil {
		return sdk.ZeroDec(), err
	}
	// Ensure ammPool.TotalShares is not zero to avoid division by zero
	if p.TotalShares.IsZero() {
		return sdk.OneDec(), nil
	}
	lpTokenPrice := ammPoolTvl.MulInt(OneShare).QuoInt(p.TotalShares.Amount)
	return lpTokenPrice, nil
}

func (pool Pool) Validate(poolId uint64) error {
	if pool.GetPoolId() != poolId {
		return errorsmod.Wrapf(ErrInvalidPool,
			"Pool was attempted to be created with incorrect pool ID.")
	}
	address, err := sdk.AccAddressFromBech32(pool.GetAddress())
	if err != nil {
		return errorsmod.Wrapf(ErrInvalidPool,
			"Pool was attempted to be created with invalid pool address.")
	}
	if !address.Equals(NewPoolAddress(poolId)) {
		return errorsmod.Wrapf(ErrInvalidPool,
			"Pool was attempted to be created with incorrect pool address.")
	}
	return nil
}
