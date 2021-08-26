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