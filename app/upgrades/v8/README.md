# v8 upgrade

This upgrade moves Sommelier to major version 8.

## Summary of changes

* Add the addresses module for mapping cosmos/evm addresses
* Update the cellarfees module to start fee auctions based on the accrued USD value of a particular denom
* Update the auction module to allow a portion of SOMM proceeds earned by auctions to be burned
* Upgrade the gravity module to v5
* Upgrade incentives module the support validator-specific rewards subsidized by the community pool
* Add events to cork and axelarcork modules for easier status tracking of requests
