<!--
order: 1
-->

# Concepts

The `estaking` module in the Elys Network extends basic staking capabilities, providing additional functionalities such as managing and distributing rewards from multiple validators, updating staking parameters, and handling Eden and EdenB token mechanics in relation to staking rewards. This module aims to enhance the staking experience and efficiency within the network.

## Flow

`estaking` is a wrapper module of `staking` module to consider `Eden/EdenB` committed tokens in commitment module as virtual delegation to reuse `distribution` module for reward distribution to `Eden/EdenB` commit users.

`estaking` module inherits `staking` module functions and add `Eden` and `EdenB` virtual validators and `Eden/EdenB` virtual delegations.

`distribution` module considers `estaking` as `staking` module and distribute rewards based on `estaking` module's `delegation`.
