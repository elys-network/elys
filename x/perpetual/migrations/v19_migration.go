package migrations

import (
	"cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/elys-network/elys/v7/x/perpetual/types"
)

func (m Migrator) V19Migration(ctx sdk.Context) error {

	legacyParams := m.keeper.GetLegacyParams(ctx)

	params := types.Params{
		LeverageMax:                         legacyParams.LeverageMax,
		BorrowInterestRateMax:               legacyParams.BorrowInterestRateMax,
		BorrowInterestRateMin:               legacyParams.BorrowInterestRateMin,
		BorrowInterestRateIncrease:          legacyParams.BorrowInterestRateIncrease,
		BorrowInterestRateDecrease:          legacyParams.BorrowInterestRateDecrease,
		HealthGainFactor:                    legacyParams.HealthGainFactor,
		MaxOpenPositions:                    legacyParams.MaxOpenPositions,
		PoolMaxLiabilitiesThreshold:         math.LegacyMustNewDecFromStr("0.3"),
		BorrowInterestPaymentFundPercentage: legacyParams.BorrowInterestPaymentFundPercentage,
		SafetyFactor:                        math.LegacyMustNewDecFromStr("1.035"),
		BorrowInterestPaymentEnabled:        legacyParams.BorrowInterestPaymentEnabled,
		WhitelistingEnabled:                 true,
		PerpetualSwapFee:                    math.LegacyMustNewDecFromStr("0.0005"),
		MaxLimitOrder:                       legacyParams.MaxLimitOrder,
		FixedFundingRate:                    math.LegacyMustNewDecFromStr("0.3"),
		MinimumLongTakeProfitPriceRatio:     math.LegacyMustNewDecFromStr("1.01"),
		MaximumLongTakeProfitPriceRatio:     legacyParams.MaximumLongTakeProfitPriceRatio,
		MaximumShortTakeProfitPriceRatio:    legacyParams.MaximumShortTakeProfitPriceRatio,
		WeightBreakingFeeFactor:             legacyParams.WeightBreakingFeeFactor,
		EnabledPools:                        []uint64{1, 13},
		MinimumNotionalValue:                math.LegacyNewDec(10),
		LongMinimumLiabilityAmount:          math.NewInt(1_000_000),
	}

	err := m.keeper.SetParams(ctx, &params)
	if err != nil {
		return err
	}

	allLegacyPools := m.keeper.GetAllLegacyPools(ctx)
	for _, legacyPool := range allLegacyPools {
		var poolAssetLong []types.PoolAsset
		for _, legacyPoolAsset := range legacyPool.PoolAssetsLong {
			poolAssetLong = append(poolAssetLong, types.PoolAsset{
				AssetDenom:  legacyPoolAsset.AssetDenom,
				Liabilities: legacyPoolAsset.Liabilities,
				Custody:     legacyPoolAsset.Custody,
				Collateral:  legacyPoolAsset.Collateral,
			})
		}

		var poolAssetShort []types.PoolAsset
		for _, legacyPoolAsset := range legacyPool.PoolAssetsShort {
			poolAssetShort = append(poolAssetShort, types.PoolAsset{
				AssetDenom:  legacyPoolAsset.AssetDenom,
				Liabilities: legacyPoolAsset.Liabilities,
				Custody:     legacyPoolAsset.Custody,
				Collateral:  legacyPoolAsset.Collateral,
			})
		}
		pool := types.Pool{
			AmmPoolId:                            legacyPool.AmmPoolId,
			BaseAssetLiabilitiesRatio:            legacyPool.BaseAssetLiabilitiesRatio,
			QuoteAssetLiabilitiesRatio:           legacyPool.QuoteAssetLiabilitiesRatio,
			BorrowInterestRate:                   legacyPool.BorrowInterestRate,
			PoolAssetsLong:                       poolAssetLong,
			PoolAssetsShort:                      poolAssetShort,
			LastHeightBorrowInterestRateComputed: legacyPool.LastHeightBorrowInterestRateComputed,
			FundingRate:                          legacyPool.FundingRate,
			FeesCollected:                        legacyPool.FeesCollected,
			LeverageMax:                          legacyPool.LeverageMax,
		}
		m.keeper.SetPool(ctx, pool)
	}
	err = m.keeper.ResetStore(ctx)
	if err != nil {
		return err
	}

	addresses := []string{
		"elys1tnk45uvxvdsemhd0pnl5y72jl43lwwm0wp0sls",
		"elys1pg2l3len5ueu2anchzchwnlm2msus4u2ng2ayh",
		"elys1rqv99msxgkyww7nm8p785p7cj2tc2upuzkq86t",
		"elys1kqh8fjrpj2cqmunfx7fmz0rxfyx06l3f829e3k",
		"elys1uuvvh83hpd3jndpt9l8tw63r7k63z9as4pw56k",
		"elys14mwhegj74cs90nlgdfmrgjlgd7rpq0yxcmpk5y",
		"elys1r5muk8te0z2url8kntgjst7wygnfmr8xyvk6hh",
		"elys1jr47lp2smyrj05f0w978pdw3thrqlsnv28h4qq",
		"elys106ps2h8y6mpnvl4f5lpjna87pv0ajasasllk96",
		"elys1yyv8gvfjg2x9d6zqrpr849mt87k49swehlmglj",
		"elys1tapvntnssnm4vqqqh8jxnv0fswc87nh9gncvrn",
		"elys1fth2nyd6g9q03n90ddrljfcxa8vja0wc8negut",
		"elys1c4u0h8h4wcwf3vq05ey3fee04tq445cz4xt2z4",
		"elys1k0zdx92k6d60vszdhm45hk4tpn06f6kglr5tu9",
		"elys1gsj0jh9z9mynemm4vjyc88wldw2trtg68hs7q3",
		"elys1mvmdf9vhw4seyfsgyn03czk4kzts7xj9sxkxzg",
		"elys1x3ylxxpjtfxvtrvgg5rh80feh02rymmvu2jl0g",
		"elys1ms0vtn0580p26rdlvcz8gulykxc0dl0qq9twj3",
		"elys190k0y8m6pktew6xs58q2dgfarqyxj3zsahpkpg",
		"elys1vm469w820ug8cxln03wlnsnwmt2pwx828w3ajq",
		"elys1s7ek72d84jvrhnugcte57rxgnf8vjwtlyuj587",
		"elys10xh5sna7f7j6ptsysp8d7c86rv42ddq4f7h003",
		"elys1u870ksd35rctnae70kl602hzrqcnwe8tqula3z",
		"elys1u93z584nyvufs647ut6re6dy4z9q3kr9gy7c6p",
		"elys133df0rz42d49jl3hzh6gtv6gf42x8vqam4he8p",
		"elys17qacjuwz58vupckhyz0g20f742khmfxs57yvjp",
		"elys1tcfusja8ekyxjkmguvfyfm5jl2jl20pxmyznaj",
		"elys1ur7tfac3y6m5v4xlvx66avpyv602na4v5safsr",
		"elys1un3ttecl2223hlep88qke0els9eaw3vqvzr5vt",
		"elys1c69yrjgrka9yafg7wcx4a0kwru9jwvwq8kf7ek",
		"elys1nsqef88pu5n2qa3gxppsm6ud0mua93xrqg3cm0",
	}

	for _, address := range addresses {
		m.keeper.WhitelistAddress(ctx, sdk.MustAccAddressFromBech32(address))
	}

	atomPool, found := m.keeper.GetPool(ctx, 1)
	if found {
		atomPool.LeverageMax = math.LegacyMustNewDecFromStr("5.0")
		m.keeper.SetPool(ctx, atomPool)
	}

	btcPool, found := m.keeper.GetPool(ctx, 13)
	if found {
		btcPool.LeverageMax = math.LegacyMustNewDecFromStr("10.0")
		m.keeper.SetPool(ctx, btcPool)
	}
	return nil
}
