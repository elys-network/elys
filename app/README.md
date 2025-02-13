# Dependencies map

## Keepers dependencies

`commitment` -> `staking`
`tokenomics` -> `commitment`
`estaking` -> `staking`, `commitment`, `distr`, `tokenomics`
`slashing` -> `estaking`
`amm` -> `commitment`
`stablestake` -> `commitment`
`distr` -> `commitment`, `estaking`
`masterchef` -> `tokenomics`
`leveragelp` -> `stablestake`, `commitment`

## Hooks dependencies

`staking` -> `slashing`, `distr`, `estaking`
`amm` -> `perpetual`, `leveragelp`, `masterchef`
`epochs` -> `commitment`, `perpetual`
`perpetual` -> `accountedpool`
`commitment` -> `estaking`
`stablestake` -> `masterchef`

## The deepest stack actions

### Eden uncommit

`commitment` hooks `estaking` with `EdenUncommitted`
`estaking` calls `commitment` with `BurnEdenBoost`
`commitment` hooks `BeforeEdenBCommitChange`, `BeforeEdenCommitChange`, `CommitmentChanged`
`estaking` hooks `staking` with `BeforeDelegationCreated`, `BeforeDelegationSharesModified`, `BeforeDelegationRemoved`, `AfterDelegationModified`
`staking` hooks `estaking` with `AfterDelegationModified`
`staking` hooks `distr` with `BeforeDelegationCreated`, `BeforeDelegationSharesModified`, `BeforeDelegationRemoved`, `AfterDelegationModified`
`distr` executes rewards claim on `BeforeDelegationSharesModified`

### Elys unstake

`staking` hooks `AfterDelegationModified`
`estaking` calls `commitment` with `BurnEdenBoost`
`commitment` hooks `BeforeEdenBCommitChange`, `BeforeEdenCommitChange`, `CommitmentChanged`
`estaking` hooks `staking` with `BeforeDelegationCreated`, `BeforeDelegationSharesModified`, `BeforeDelegationRemoved`, `AfterDelegationModified`
`staking` hooks `estaking` with `AfterDelegationModified`
`staking` hooks `distr` with `BeforeDelegationCreated`, `BeforeDelegationSharesModified`, `BeforeDelegationRemoved`, `AfterDelegationModified`
`distr` executes rewards claim on `BeforeDelegationSharesModified`
