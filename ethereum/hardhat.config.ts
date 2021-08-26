import "@nomiclabs/hardhat-waffle";
import { task } from 'hardhat/config';
import CellarArtifact from './artifacts/Cellar.json';

task('integration_test_setup', 'Sets up contracts for the integration test', async (args, hre) => {
  const ADDRESSES = {
    CELLAR_OWNER: '0xB6C951cf962977f123bF37de42945f7ca1cd2A52',
    CELLAR: '0x6ea5992aB4A78D5720bD12A089D13c073d04B55d',
    GRAVITY: '0xFbB0BCfed0c82043A7d5387C35Ad8450b44A4cde'
  };

  // Take over the cellar owner so we can transfer ownership
  await hre.network.provider.request({
    method: "hardhat_impersonateAccount",
    params: [ADDRESSES.CELLAR_OWNER]
  });

  const signer = await hre.ethers.getSigner(ADDRESSES.CELLAR_OWNER)

  // Transfer ownership to gravity
  const Cellar = new hre.ethers.ContractFactory(CellarArtifact.abi, CellarArtifact.bytecode, signer);
  const cellar = await Cellar.attach(ADDRESSES.CELLAR);

  const { hash } = await cellar.transferOwnership(ADDRESSES.GRAVITY);

  console.log(`Cellar contract at ${ADDRESSES.CELLAR} is now owned by Gravity contract at ${ADDRESSES.GRAVITY}`);
  console.log(`Tx hash: ${hash}`);
  console.log('='.repeat(80));

  await hre.run('node');
});

/**
 * @type import('hardhat/config').HardhatUserConfig
 */
module.exports = {
  networks: {
    hardhat: {
      forking: {
        url: "https://mainnet.infura.io/v3/d6f22be0f7fd447186086d2495779003"
      }
    },
  },
  solidity: {
    version: "0.7.3",
    settings: {
      optimizer: {
        enabled: true
      }
    }
  },
  // TODO: add forking configuration
  typechain: {
    outDir: "typechain",
    target: "ethers-v5",
    runOnCompile: true
  },
  gasReporter: {
    enabled: true
  },
};
