import '@nomiclabs/hardhat-waffle';
import {task} from 'hardhat/config';
import * as constants from "./addresses";

task(
    'integration_test_setup',
    'Sets up contracts for the integration test',
    async (args, hre) => {


        // console.log('retrieving gravity contract');
        // const gravitySigner = await hre.ethers.getSigner(ADDRESSES.GRAVITY_OWNER)
        // // Attach Gravity contract
        // const Gravity = hre.ethers.getContractAt(
        //     'Gravity',
        //     ADDRESSES.GRAVITY,
        //     gravitySigner
        // );
        // const gravity = await Gravity;
        // console.log(`attached to gravity: ${gravity}`)

        // console.log('taking over cellar owner');
        // // Take over the cellar owner so we can transfer ownership
        // await hre.network.provider.request({
        //     method: 'hardhat_impersonateAccount',
        //     params: [ADDRESSES.CELLAR_OWNER],
        // });

        // console.log('getting cellar signer');
        // const cellarSigner = await hre.ethers.getSigner(ADDRESSES.CELLAR_OWNER);
        //
        // console.log('retrieving cellar contract');
        // // Transfer ownership to gravity
        // const Cellar = hre.ethers.getContractAt(
        //     'CellarPoolShare',
        //     ADDRESSES.CELLAR,
        //     cellarSigner,
        // );
        // const cellar = await Cellar;
        //
        // let { hash } = await cellar.transferOwnership(ADDRESSES.GRAVITY, {
        //     gasPrice: hre.ethers.BigNumber.from('99916001694'),
        // });
        //
        // // Send ETH to needed accounts
        //
        // console.log(
        //     `Cellar contract at ${ADDRESSES.CELLAR} is now owned by Gravity contract at ${ADDRESSES.GRAVITY}`,
        // );
        // console.log(`Tx hash: ${hash}`);
        // console.log('='.repeat(80));

        // Take over vitalik.eth
        // await hre.network.provider.request({
        //     method: 'hardhat_impersonateAccount',
        //     params: [constants.WHALE],
        // });
        //
        // // Send ETH to needed parties
        // const whaleSigner = await hre.ethers.getSigner(constants.WHALE);
        //
        // for (let addr of constants.VALIDATORS) {
        //     await whaleSigner.sendTransaction({
        //         to: addr,
        //         value: hre.ethers.utils.parseEther('100'),
        //     });
        // }

        // start the ethereum node after all setup is complete
        // await hre.run('node');
    },
);

/**
 * @type import('hardhat/config').HardhatUserConfig
 */
module.exports = {
    networks: {
        hardhat: {
            forking: {
                url: 'https://mainnet.infura.io/v3/d6f22be0f7fd447186086d2495779003',
                blockNumber: 13357100,
            },
        },
        mainnet: {
            url: 'https://mainnet.infura.io/v3/d6f22be0f7fd447186086d2495779003',
        },
    },
    solidity: {
        compilers: [
            {
                version: '0.6.6',
                settings: {
                    optimizer: {
                        enabled: true,
                        runs: 200,
                    },
                },
            },
            {
                version: '0.7.6',
                settings: {
                    optimizer: {
                        enabled: true,
                        runs: 200,
                    },
                },
            },
            {
                version: '0.8.0',
                settings: {
                    optimizer: {
                        enabled: true,
                    },
                },
            },
        ],
    },
    // TODO: add forking configuration
    typechain: {
        outDir: 'typechain',
        target: 'ethers-v5',
        runOnCompile: true,
    },
    gasReporter: {
        enabled: true,
    },
};
