# v4 upgrade

This upgrade moves Sommelier to major version 4.

## Summary of changes

* Switch to use of the upgrade module
* Delete the allocation module, which was unused
* Add the cork module for receiving arbitrary logic calls
* Add the cellarfees module to provide a module account for cellar fees to be bridged to
* Support sending to specified module accounts over the bridge
* Community spend governance proposal for sending funds over the bridge to an Ethereum address
* Fix a bug affecting the capitalization of ERC20 addresses in denominations
* Fix a bug incorrectly setting the timeouts of ContractCallTxs