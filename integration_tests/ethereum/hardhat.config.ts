import '@nomiclabs/hardhat-waffle';
// import "@nomiclabs/hardhat-etherscan";
import { task } from 'hardhat/config';
import CellarArtifact from './artifacts/Cellar.json';

task(
    'integration_test_setup',
    'Sets up contracts for the integration test',
    async (args, hre) => {
        const ADDRESSES = {
            CELLAR_OWNER: '0xB6C951cf962977f123bF37de42945f7ca1cd2A52',
            CELLAR: '0x6ea5992aB4A78D5720bD12A089D13c073d04B55d',
            GRAVITY: '0xFbB0BCfed0c82043A7d5387C35Ad8450b44A4cde',
            WHALE: '0xd8da6bf26964af9d7eed9e03e53415d37aa96045',
        };

        // Take over the cellar owner so we can transfer ownership
        await hre.network.provider.request({
            method: 'hardhat_impersonateAccount',
            params: [ADDRESSES.CELLAR_OWNER],
        });

        const signer = await hre.ethers.getSigner(ADDRESSES.CELLAR_OWNER);

        // Attach Gravity contract
        const Gravity = hre.ethers.getContractAt('Gravity', ADDRESSES.GRAVITY);
        const gravity = await Gravity;

        // Transfer ownership to gravity
        const Cellar = new hre.ethers.ContractFactory(
            CellarArtifact.abi,
            CellarArtifact.bytecode,
            signer,
        );
        const cellar = await Cellar.attach(ADDRESSES.CELLAR);

        const { hash } = await cellar.transferOwnership(ADDRESSES.GRAVITY);

        // Send ETH to needed accounts

        console.log(
            `Cellar contract at ${ADDRESSES.CELLAR} is now owned by Gravity contract at ${ADDRESSES.GRAVITY}`,
        );
        console.log(`Tx hash: ${hash}`);
        console.log('='.repeat(80));

        // Take over vitalik.eth
        await hre.network.provider.request({
            method: 'hardhat_impersonateAccount',
            params: [ADDRESSES.WHALE],
        });

        // Send ETH to needed parties
        const whaleSigner = await hre.ethers.getSigner(ADDRESSES.WHALE);

        const recipients = [
            '0xd312f0f1B39D54Db2829537595fC1167B14d4b34',
            '0x7bE2a04df4b9C3227928147461e19158eB2B11d1',
            '0xb8c6886FDDa38adaa0F416722dd5554886C43055',
            '0x14fdAC734De10065093C4Ed4a83C41638378005A',
        ];

        for (let addr of recipients) {
            await whaleSigner.sendTransaction({
                to: addr,
                value: hre.ethers.utils.parseEther('100'),
            });
        }

        // start the ethereum node after all setup is complete
        await hre.run('node');
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
                blockNumber: 13103326,
            },
        },
    },
    solidity: {
        compilers: [
            {
                version: '0.6.6',
                settings: {
                    optimizer: {
                        enabled: true,
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
