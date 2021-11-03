from brownie import accounts, interface

VALIDATOR_MNEMONICS = [
    "say monitor orient heart super local purse cricket caution primary bring insane road expect rather help two extend own execute throw nation plunge subject",
    "march carpet enact kiss tribe plastic wash enter index lift topic riot try juice replace supreme original shift hover adapt mutual holiday manual nut",
    "assault section bleak gadget venture ship oblige pave fabric more initial april dutch scene parade shallow educate gesture lunar match patch hawk member problem",
    "receive roof marine sure lady hundred sea enact exist place bean wagon kingdom betray science photo loop funny bargain floor suspect only strike endless",
]


def main():
    whale = accounts[0]
    gravity_contract = accounts[1]
    # gravity_owner = accounts[2]
    cellar_contract = accounts[3]
    cellar_owner = accounts[4]

    sommelier0 = accounts.from_mnemonic(VALIDATOR_MNEMONICS[0])
    sommelier1 = accounts.from_mnemonic(VALIDATOR_MNEMONICS[1])
    sommelier2 = accounts.from_mnemonic(VALIDATOR_MNEMONICS[2])
    sommelier3 = accounts.from_mnemonic(VALIDATOR_MNEMONICS[3])

    whale.transfer(sommelier0, 100 * 1000000000000000000)
    whale.transfer(sommelier1, 100 * 1000000000000000000)
    whale.transfer(sommelier2, 100 * 1000000000000000000)
    whale.transfer(sommelier3, 100 * 1000000000000000000)
    whale.transfer(cellar_contract, 100 * 1000000000000000000)
    whale.transfer(gravity_contract, 100 * 1000000000000000000)

    cellar = interface.CellarPoolShareLimitUSDCETH('0x08c0a0B8D2eDB1d040d4f2C00A1d2f9d9b9F2677')
    response = cellar.transferOwnership.call(gravity_contract, {
        'from': cellar_owner,
        'gas_limit': 99916001694,
    })
    print(response)


