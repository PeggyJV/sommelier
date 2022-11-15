pragma solidity ^0.8.10;

import "@openzeppelin/contracts/token/ERC20/ERC20.sol";

contract AlphaERC20 is ERC20 {
    constructor() public ERC20("Alpha Cellar Fees Test ERC20", "ALPHA") {
        // Mint to the validator/orchestrator that will be used for cellarfees testing
        _mint(0x14fdAC734De10065093C4Ed4a83C41638378005A, 1000000000000000);
    }
}
