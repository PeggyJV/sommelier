package types

import (
	"github.com/ethereum/go-ethereum/accounts/abi"
	"strings"
)

const CellarABI = `
[
    {
      "inputs": [
        {
          "internalType": "string",
          "name": "name_",
          "type": "string"
        },
        {
          "internalType": "string",
          "name": "symbol_",
          "type": "string"
        },
        {
          "internalType": "address",
          "name": "_token0",
          "type": "address"
        },
        {
          "internalType": "address",
          "name": "_token1",
          "type": "address"
        },
        {
          "internalType": "uint24",
          "name": "_feeLevel",
          "type": "uint24"
        },
        {
          "components": [
            {
              "internalType": "uint184",
              "name": "tokenId",
              "type": "uint184"
            },
            {
              "internalType": "int24",
              "name": "tickUpper",
              "type": "int24"
            },
            {
              "internalType": "int24",
              "name": "tickLower",
              "type": "int24"
            },
            {
              "internalType": "uint24",
              "name": "weight",
              "type": "uint24"
            }
          ],
          "internalType": "struct ICellarPoolShare.CellarTickInfo[]",
          "name": "_cellarTickInfo",
          "type": "tuple[]"
        }
      ],
      "stateMutability": "nonpayable",
      "type": "constructor"
    },
    {
      "anonymous": false,
      "inputs": [
        {
          "indexed": true,
          "internalType": "address",
          "name": "token0",
          "type": "address"
        },
        {
          "indexed": true,
          "internalType": "address",
          "name": "token1",
          "type": "address"
        },
        {
          "indexed": false,
          "internalType": "uint128",
          "name": "liquidity",
          "type": "uint128"
        },
        {
          "indexed": false,
          "internalType": "uint256",
          "name": "amount0",
          "type": "uint256"
        },
        {
          "indexed": false,
          "internalType": "uint256",
          "name": "amount1",
          "type": "uint256"
        }
      ],
      "name": "AddedLiquidity",
      "type": "event"
    },
    {
      "anonymous": false,
      "inputs": [
        {
          "indexed": true,
          "internalType": "address",
          "name": "owner",
          "type": "address"
        },
        {
          "indexed": true,
          "internalType": "address",
          "name": "spender",
          "type": "address"
        },
        {
          "indexed": false,
          "internalType": "uint256",
          "name": "value",
          "type": "uint256"
        }
      ],
      "name": "Approval",
      "type": "event"
    },
    {
      "anonymous": false,
      "inputs": [
        {
          "indexed": false,
          "internalType": "uint256",
          "name": "fees0",
          "type": "uint256"
        },
        {
          "indexed": false,
          "internalType": "uint256",
          "name": "fees1",
          "type": "uint256"
        },
        {
          "indexed": false,
          "internalType": "uint256",
          "name": "managementFee0",
          "type": "uint256"
        },
        {
          "indexed": false,
          "internalType": "uint256",
          "name": "managementFee1",
          "type": "uint256"
        },
        {
          "indexed": false,
          "internalType": "uint256",
          "name": "performanceFee0",
          "type": "uint256"
        },
        {
          "indexed": false,
          "internalType": "uint256",
          "name": "performanceFee1",
          "type": "uint256"
        },
        {
          "indexed": false,
          "internalType": "uint256",
          "name": "amount0",
          "type": "uint256"
        },
        {
          "indexed": false,
          "internalType": "uint256",
          "name": "amount1",
          "type": "uint256"
        }
      ],
      "name": "Rebalance",
      "type": "event"
    },
    {
      "anonymous": false,
      "inputs": [
        {
          "indexed": false,
          "internalType": "uint256",
          "name": "fees0",
          "type": "uint256"
        },
        {
          "indexed": false,
          "internalType": "uint256",
          "name": "fees1",
          "type": "uint256"
        },
        {
          "indexed": false,
          "internalType": "uint256",
          "name": "managementFee0",
          "type": "uint256"
        },
        {
          "indexed": false,
          "internalType": "uint256",
          "name": "managementFee1",
          "type": "uint256"
        },
        {
          "indexed": false,
          "internalType": "uint256",
          "name": "performanceFee0",
          "type": "uint256"
        },
        {
          "indexed": false,
          "internalType": "uint256",
          "name": "performanceFee1",
          "type": "uint256"
        },
        {
          "indexed": false,
          "internalType": "uint256",
          "name": "amount0",
          "type": "uint256"
        },
        {
          "indexed": false,
          "internalType": "uint256",
          "name": "amount1",
          "type": "uint256"
        }
      ],
      "name": "Reinvest",
      "type": "event"
    },
    {
      "anonymous": false,
      "inputs": [
        {
          "indexed": true,
          "internalType": "address",
          "name": "token0",
          "type": "address"
        },
        {
          "indexed": true,
          "internalType": "address",
          "name": "token1",
          "type": "address"
        },
        {
          "indexed": false,
          "internalType": "uint128",
          "name": "liquidity",
          "type": "uint128"
        },
        {
          "indexed": false,
          "internalType": "uint256",
          "name": "amount0",
          "type": "uint256"
        },
        {
          "indexed": false,
          "internalType": "uint256",
          "name": "amount1",
          "type": "uint256"
        }
      ],
      "name": "RemovedLiquidity",
      "type": "event"
    },
    {
      "anonymous": false,
      "inputs": [
        {
          "indexed": false,
          "internalType": "address",
          "name": "adjuster",
          "type": "address"
        },
        {
          "indexed": false,
          "internalType": "bool",
          "name": "value",
          "type": "bool"
        }
      ],
      "name": "SetAdjuster",
      "type": "event"
    },
    {
      "anonymous": false,
      "inputs": [
        {
          "indexed": false,
          "internalType": "uint256",
          "name": "newFee",
          "type": "uint256"
        }
      ],
      "name": "SetManagementFee",
      "type": "event"
    },
    {
      "anonymous": false,
      "inputs": [
        {
          "indexed": false,
          "internalType": "uint256",
          "name": "newFee",
          "type": "uint256"
        }
      ],
      "name": "SetPerformanceFee",
      "type": "event"
    },
    {
      "anonymous": false,
      "inputs": [
        {
          "indexed": false,
          "internalType": "address",
          "name": "validator",
          "type": "address"
        },
        {
          "indexed": false,
          "internalType": "bool",
          "name": "value",
          "type": "bool"
        }
      ],
      "name": "SetValidator",
      "type": "event"
    },
    {
      "anonymous": false,
      "inputs": [
        {
          "indexed": true,
          "internalType": "address",
          "name": "from",
          "type": "address"
        },
        {
          "indexed": true,
          "internalType": "address",
          "name": "to",
          "type": "address"
        },
        {
          "indexed": false,
          "internalType": "uint256",
          "name": "value",
          "type": "uint256"
        }
      ],
      "name": "Transfer",
      "type": "event"
    },
    {
      "anonymous": false,
      "inputs": [
        {
          "indexed": false,
          "internalType": "address",
          "name": "newOwner",
          "type": "address"
        }
      ],
      "name": "TransferOwnership",
      "type": "event"
    },
    {
      "inputs": [
        {
          "components": [
            {
              "internalType": "uint256",
              "name": "amount0Desired",
              "type": "uint256"
            },
            {
              "internalType": "uint256",
              "name": "amount1Desired",
              "type": "uint256"
            },
            {
              "internalType": "uint256",
              "name": "amount0Min",
              "type": "uint256"
            },
            {
              "internalType": "uint256",
              "name": "amount1Min",
              "type": "uint256"
            },
            {
              "internalType": "address",
              "name": "recipient",
              "type": "address"
            },
            {
              "internalType": "uint256",
              "name": "deadline",
              "type": "uint256"
            }
          ],
          "internalType": "struct ICellarPoolShare.CellarAddParams",
          "name": "cellarParams",
          "type": "tuple"
        }
      ],
      "name": "addLiquidityForUniV3",
      "outputs": [],
      "stateMutability": "payable",
      "type": "function"
    },
    {
      "inputs": [
        {
          "internalType": "address",
          "name": "",
          "type": "address"
        }
      ],
      "name": "adjuster",
      "outputs": [
        {
          "internalType": "bool",
          "name": "",
          "type": "bool"
        }
      ],
      "stateMutability": "view",
      "type": "function"
    },
    {
      "inputs": [
        {
          "internalType": "address",
          "name": "owner_",
          "type": "address"
        },
        {
          "internalType": "address",
          "name": "spender",
          "type": "address"
        }
      ],
      "name": "allowance",
      "outputs": [
        {
          "internalType": "uint256",
          "name": "",
          "type": "uint256"
        }
      ],
      "stateMutability": "view",
      "type": "function"
    },
    {
      "inputs": [
        {
          "internalType": "address",
          "name": "spender",
          "type": "address"
        },
        {
          "internalType": "uint256",
          "name": "amount",
          "type": "uint256"
        }
      ],
      "name": "approve",
      "outputs": [
        {
          "internalType": "bool",
          "name": "",
          "type": "bool"
        }
      ],
      "stateMutability": "nonpayable",
      "type": "function"
    },
    {
      "inputs": [
        {
          "internalType": "address",
          "name": "account",
          "type": "address"
        }
      ],
      "name": "balanceOf",
      "outputs": [
        {
          "internalType": "uint256",
          "name": "",
          "type": "uint256"
        }
      ],
      "stateMutability": "view",
      "type": "function"
    },
    {
      "inputs": [
        {
          "internalType": "uint256",
          "name": "",
          "type": "uint256"
        }
      ],
      "name": "cellarTickInfo",
      "outputs": [
        {
          "internalType": "uint184",
          "name": "tokenId",
          "type": "uint184"
        },
        {
          "internalType": "int24",
          "name": "tickUpper",
          "type": "int24"
        },
        {
          "internalType": "int24",
          "name": "tickLower",
          "type": "int24"
        },
        {
          "internalType": "uint24",
          "name": "weight",
          "type": "uint24"
        }
      ],
      "stateMutability": "view",
      "type": "function"
    },
    {
      "inputs": [],
      "name": "decimals",
      "outputs": [
        {
          "internalType": "uint8",
          "name": "",
          "type": "uint8"
        }
      ],
      "stateMutability": "pure",
      "type": "function"
    },
    {
      "inputs": [],
      "name": "feeLevel",
      "outputs": [
        {
          "internalType": "uint24",
          "name": "",
          "type": "uint24"
        }
      ],
      "stateMutability": "view",
      "type": "function"
    },
    {
      "inputs": [],
      "name": "getCellarTickInfo",
      "outputs": [
        {
          "components": [
            {
              "internalType": "uint184",
              "name": "tokenId",
              "type": "uint184"
            },
            {
              "internalType": "int24",
              "name": "tickUpper",
              "type": "int24"
            },
            {
              "internalType": "int24",
              "name": "tickLower",
              "type": "int24"
            },
            {
              "internalType": "uint24",
              "name": "weight",
              "type": "uint24"
            }
          ],
          "internalType": "struct ICellarPoolShare.CellarTickInfo[]",
          "name": "",
          "type": "tuple[]"
        }
      ],
      "stateMutability": "view",
      "type": "function"
    },
    {
      "inputs": [
        {
          "internalType": "address",
          "name": "",
          "type": "address"
        }
      ],
      "name": "lastLockedBlock",
      "outputs": [
        {
          "internalType": "uint256",
          "name": "",
          "type": "uint256"
        }
      ],
      "stateMutability": "view",
      "type": "function"
    },
    {
      "inputs": [],
      "name": "managementFee",
      "outputs": [
        {
          "internalType": "uint256",
          "name": "",
          "type": "uint256"
        }
      ],
      "stateMutability": "view",
      "type": "function"
    },
    {
      "inputs": [],
      "name": "name",
      "outputs": [
        {
          "internalType": "string",
          "name": "",
          "type": "string"
        }
      ],
      "stateMutability": "view",
      "type": "function"
    },
    {
      "inputs": [],
      "name": "owner",
      "outputs": [
        {
          "internalType": "address",
          "name": "",
          "type": "address"
        }
      ],
      "stateMutability": "view",
      "type": "function"
    },
    {
      "inputs": [],
      "name": "performanceFee",
      "outputs": [
        {
          "internalType": "uint256",
          "name": "",
          "type": "uint256"
        }
      ],
      "stateMutability": "view",
      "type": "function"
    },
    {
      "inputs": [
        {
          "components": [
            {
              "internalType": "uint184",
              "name": "tokenId",
              "type": "uint184"
            },
            {
              "internalType": "int24",
              "name": "tickUpper",
              "type": "int24"
            },
            {
              "internalType": "int24",
              "name": "tickLower",
              "type": "int24"
            },
            {
              "internalType": "uint24",
              "name": "weight",
              "type": "uint24"
            }
          ],
          "internalType": "struct ICellarPoolShare.CellarTickInfo[]",
          "name": "_cellarTickInfo",
          "type": "tuple[]"
        },
        {
          "internalType": "uint256",
          "name": "currentPriceX96",
          "type": "uint256"
        }
      ],
      "name": "rebalance",
      "outputs": [],
      "stateMutability": "nonpayable",
      "type": "function"
    },
    {
      "inputs": [
        {
          "internalType": "uint256",
          "name": "currentPriceX96",
          "type": "uint256"
        }
      ],
      "name": "reinvest",
      "outputs": [],
      "stateMutability": "nonpayable",
      "type": "function"
    },
    {
      "inputs": [
        {
          "components": [
            {
              "internalType": "uint256",
              "name": "tokenAmount",
              "type": "uint256"
            },
            {
              "internalType": "uint256",
              "name": "amount0Min",
              "type": "uint256"
            },
            {
              "internalType": "uint256",
              "name": "amount1Min",
              "type": "uint256"
            },
            {
              "internalType": "address",
              "name": "recipient",
              "type": "address"
            },
            {
              "internalType": "uint256",
              "name": "deadline",
              "type": "uint256"
            }
          ],
          "internalType": "struct ICellarPoolShare.CellarRemoveParams",
          "name": "cellarParams",
          "type": "tuple"
        }
      ],
      "name": "removeLiquidityFromUniV3",
      "outputs": [],
      "stateMutability": "nonpayable",
      "type": "function"
    },
    {
      "inputs": [
        {
          "internalType": "address",
          "name": "_adjuster",
          "type": "address"
        },
        {
          "internalType": "bool",
          "name": "value",
          "type": "bool"
        }
      ],
      "name": "setAdjuster",
      "outputs": [],
      "stateMutability": "nonpayable",
      "type": "function"
    },
    {
      "inputs": [
        {
          "internalType": "uint256",
          "name": "newFee",
          "type": "uint256"
        }
      ],
      "name": "setManagementFee",
      "outputs": [],
      "stateMutability": "nonpayable",
      "type": "function"
    },
    {
      "inputs": [
        {
          "internalType": "uint256",
          "name": "newFee",
          "type": "uint256"
        }
      ],
      "name": "setPerformanceFee",
      "outputs": [],
      "stateMutability": "nonpayable",
      "type": "function"
    },
    {
      "inputs": [
        {
          "internalType": "address",
          "name": "_validator",
          "type": "address"
        },
        {
          "internalType": "bool",
          "name": "value",
          "type": "bool"
        }
      ],
      "name": "setValidator",
      "outputs": [],
      "stateMutability": "nonpayable",
      "type": "function"
    },
    {
      "inputs": [],
      "name": "symbol",
      "outputs": [
        {
          "internalType": "string",
          "name": "",
          "type": "string"
        }
      ],
      "stateMutability": "view",
      "type": "function"
    },
    {
      "inputs": [],
      "name": "token0",
      "outputs": [
        {
          "internalType": "address",
          "name": "",
          "type": "address"
        }
      ],
      "stateMutability": "view",
      "type": "function"
    },
    {
      "inputs": [],
      "name": "token1",
      "outputs": [
        {
          "internalType": "address",
          "name": "",
          "type": "address"
        }
      ],
      "stateMutability": "view",
      "type": "function"
    },
    {
      "inputs": [],
      "name": "totalSupply",
      "outputs": [
        {
          "internalType": "uint256",
          "name": "",
          "type": "uint256"
        }
      ],
      "stateMutability": "view",
      "type": "function"
    },
    {
      "inputs": [
        {
          "internalType": "address",
          "name": "recipient",
          "type": "address"
        },
        {
          "internalType": "uint256",
          "name": "amount",
          "type": "uint256"
        }
      ],
      "name": "transfer",
      "outputs": [
        {
          "internalType": "bool",
          "name": "",
          "type": "bool"
        }
      ],
      "stateMutability": "nonpayable",
      "type": "function"
    },
    {
      "inputs": [
        {
          "internalType": "address",
          "name": "sender",
          "type": "address"
        },
        {
          "internalType": "address",
          "name": "recipient",
          "type": "address"
        },
        {
          "internalType": "uint256",
          "name": "amount",
          "type": "uint256"
        }
      ],
      "name": "transferFrom",
      "outputs": [
        {
          "internalType": "bool",
          "name": "",
          "type": "bool"
        }
      ],
      "stateMutability": "nonpayable",
      "type": "function"
    },
    {
      "inputs": [
        {
          "internalType": "address",
          "name": "newOwner",
          "type": "address"
        }
      ],
      "name": "transferOwnership",
      "outputs": [],
      "stateMutability": "nonpayable",
      "type": "function"
    },
    {
      "inputs": [
        {
          "internalType": "address",
          "name": "",
          "type": "address"
        }
      ],
      "name": "validator",
      "outputs": [
        {
          "internalType": "bool",
          "name": "",
          "type": "bool"
        }
      ],
      "stateMutability": "view",
      "type": "function"
    },
    {
      "stateMutability": "payable",
      "type": "receive"
    }
  ]`

func abiJSON() (abi.ABI, error) {
	return abi.JSON(strings.NewReader(CellarABI))
}