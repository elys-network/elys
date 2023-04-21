<!--
order: 1
-->

# Concepts

The Burner Module is a Cosmos SDK module that allows for the automatic burning of native tokens on a regular basis. It depends on the Epochs module, which triggers the burning mechanism at the specified interval (e.g., daily, weekly, etc.).

In this scenario, a user sends native tokens to the zero address `elys1qqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqnrec2l`. Afterward, the module wallet retrieves the available native tokens from the zero address and transfers them to the module address. Finally, at an arbitrary interval, the module wallet burns the tokens that were transferred to the module address.
