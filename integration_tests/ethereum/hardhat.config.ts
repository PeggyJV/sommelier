import '@nomiclabs/hardhat-ethers';
// import '@nomiclabs/hardhat-waffle';
import {task} from "hardhat/config";
import * as constants from "./addresses";

task(
    'integration_test_setup',
    'Sets up contracts for the integration test',
    async (args, hre) => {

        // Take over vitalik.eth
        await hre.network.provider.request({
            method: 'hardhat_impersonateAccount',
            params: [constants.WHALE],
        });

        // Send ETH to needed parties
        const whaleSigner = await hre.ethers.getSigner(constants.WHALE);

        for (let addr of constants.VALIDATORS) {
            await whaleSigner.sendTransaction({
                to: addr,
                value: hre.ethers.utils.parseEther('100'),
            });
        }

        let powers: number[] = [1073741823,1073741823,1073741823,1073741823];
        let powerThreshold: number = 6666;

        const Gravity = await hre.ethers.getContractFactory("Gravity");
        const gravity = (await Gravity.deploy(
            hre.ethers.utils.formatBytes32String("gravitytest"),
            powerThreshold,
            constants.VALIDATORS,
            powers
        ));

        await gravity.deployed();
        console.log(`gravity contract deployed at - ${gravity.address}`)

        const Cellar = await hre.ethers.getContractFactory("CellarPoolShare");
        const cellar = (await Cellar.deploy(
            "mock cellar",
            "mock",
            "0xa0b86991c6218b36c1d19d4a2e9eb0ce3606eb48",
            "0xc02aaa39b223fe8d0a0e5c4f27ead9083c756cc2",
            3000,
            [[0, 600, 300, 900]],
        ));
        await cellar.deployed();
        console.log(`cellar contract deploy at - ${cellar.address}`);

        let cellarSignerAddress = await cellar.signer.getAddress()
        await hre.network.provider.request({
            method: 'hardhat_impersonateAccount',
            params: [cellarSignerAddress],
        });

        let { adjusterHash } = await cellar.setAdjuster(gravity.address, true, {
            gasPrice: hre.ethers.BigNumber.from('99916001694'),
        })
        console.log(
            `Cellar contract at ${cellar.address} now has adjuster Gravity contract ${gravity.address} with hash ${adjusterHash}`,
        )
        let { hash } = await cellar.transferOwnership(gravity.address, {
            gasPrice: hre.ethers.BigNumber.from('99916001694'),
            from: cellarSignerAddress
        });
        console.log(
            `Cellar contract at ${cellar.address} is now owned by Gravity contract at ${gravity.address} with hash ${hash}`,
        );

        const Counter = await hre.ethers.getContractFactory("Counter");
        const counter = await Counter.deploy();
        await counter.deployed();
        console.log(`counter contract deployed at - ${counter.address}`);

        await hre.network.provider.send("evm_setIntervalMining", [1000]);

        await hre.run('node');
    });


/**
 * @type import('hardhat/config').HardhatUserConfig
 */
const ARCHIVE_NODE_URL = process.env.ARCHIVE_NODE_URL;

module.exports = {
    networks: {
        hardhat: {
            forking: {
                url: ARCHIVE_NODE_URL,
                blockNumber: 13405367,
            },
        },
    },
    solidity: {
        compilers: [
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
                version: '0.8.10',
                settings: {
                    optimizer: {
                        enabled: true,
                    },
                },
            },
        ],
    },
    typechain: {
        outDir: 'typechain',
        target: 'ethers-v5',
        runOnCompile: true,
    },
    gasReporter: {
        enabled: true,
    },
};
