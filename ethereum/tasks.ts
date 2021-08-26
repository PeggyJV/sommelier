import '@nomiclabs/hardhat-ethers';
import { task } from 'hardhat/config';
import { ethers, network } from 'hardhat';

import CellarArtifact from './artifacts/Cellar.json';

task('accounts', 'Prints the list of accounts', async (args, hre) => {
    const accounts = await hre.ethers.getSigners();

    for (const account of accounts) {
        console.log(account.address);
    }
});

task('integration_test_setup', 'Sets up contracts for the integration test', async (args, hre) => {
    const ADDRESSES = {
        CELLAR_OWNER: '0xB6C951cf962977f123bF37de42945f7ca1cd2A52',
        CELLAR: '0x6ea5992aB4A78D5720bD12A089D13c073d04B55d',
        GRAVITY: '0xFbB0BCfed0c82043A7d5387C35Ad8450b44A4cde'
    };

    // Take over the cellar owner so we can transfer ownership
    await network.provider.request({
        method: "hardhat_impersonateAccount",
        params: [ADDRESSES.CELLAR_OWNER]
    });

    // Transfer ownership to gravity
    const Cellar = new ethers.ContractFactory(CellarArtifact.abi, CellarArtifact.bytecode);
    const cellar = await Cellar.attach(ADDRESSES.CELLAR);

    await cellar.transferOwnership(ADDRESSES.GRAVITY);

    console.log(`Cellar contract at ${ADDRESSES.CELLAR} is now owned by Gravity contract at ${ADDRESSES.GRAVITY}`);
});

