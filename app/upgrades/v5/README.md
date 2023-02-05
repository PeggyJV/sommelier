# v5 upgrade

This upgrade moves Sommelier to major version 5.

## Summary of changes

* Update the cork module to only use scheduled corks and attach an ID and status to cork execution
* Add the auction module to provide the capability of auctioning received cellar fee tokens in exchange for SOMM
* Update the cellarfees module to trigger auctions for balances stored in its module account based on amount of times fees have been sent
* Add the pubsub module to provide a method for coordinating publishers (strategy providers) and subscribers (validators) via subscription IDs (cellar addresses) to enable multi-strategist authenticaton and authorization to Steward
* Add the incentives module to provide a way of distributing community pool funds on a per-block, pro-rata basis to incentivize stakers while cellar TVL is built over time into a meaningful staking return