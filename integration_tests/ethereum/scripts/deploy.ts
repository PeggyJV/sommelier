import {ethers, network} from "hardhat";
import {BigNumberish, Signer} from "ethers";
import * as constants from "../addresses";

const hre = require("hardhat");

export async function getSignerAddresses(signers: Signer[]) {
    return await Promise.all(signers.map(signer => signer.getAddress()));
}

export function makeCheckpoint(
    validators: string[],
    powers: BigNumberish[],
    valsetNonce: BigNumberish,
    gravityId: string
) {
    const methodName = ethers.utils.formatBytes32String("checkpoint");
    const gravityIDBytes = ethers.utils.formatBytes32String(gravityId);

    let abiEncoded = ethers.utils.defaultAbiCoder.encode(
        ["bytes32", "bytes32", "uint256", "address[]", "uint256[]"],
        [gravityIDBytes, methodName, valsetNonce, validators, powers]
    );

    return ethers.utils.keccak256(abiEncoded);
}

export async function deployTestnetContract() {
    let powers: number[] = [1073741824,1073741824,1073741824,1073741824];
    let powerThreshold: number = 6666;
    let {gravity} = await deployContracts("gravitytest", constants.VALIDATORS, powers, powerThreshold);

    console.log('taking over cellar owner');
    // Take over the cellar owner so we can transfer ownership
    await network.provider.request({
        method: 'hardhat_impersonateAccount',
        params: [constants.CELLAR_OWNER],
    });
    const cellarSigner = await ethers.getSigner(constants.CELLAR_OWNER);
    const Cellar = ethers.getContractAt(
        'CellarPoolShare',
        constants.CELLAR,
        cellarSigner,
    );
    const cellar = await Cellar;

    let { hash } = await cellar.transferOwnership(gravity.address, {
        gasPrice: ethers.BigNumber.from('99916001694'),
    });
    console.log(
        `Cellar contract at ${cellar.address} is now owned by Gravity contract at ${gravity.address}`,
    );

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

    await hre.run('node');

}

export async function deployContracts(
    gravityId: string = "foo",
    valAddresses: string[],
    powers: number[],
    powerThreshold: number
) {
    console.log(`creating gravity contract`)
    const Gravity = await ethers.getContractFactory("Gravity");
    console.log(`creating checkpoint`)
    const checkpoint = makeCheckpoint(valAddresses, powers, 0, gravityId);
    console.log(`deploying gravity contract`)
    const gravity = (await Gravity.deploy(
        ethers.utils.formatBytes32String(gravityId),
        powerThreshold,
        valAddresses,
        powers
    ));

    await gravity.deployed();
    console.log(`gravity contract deployed at - ${gravity.address}`)

    return { gravity, checkpoint };
}

deployTestnetContract().then(() => process.exit(0))
    .catch((error) => {
        console.error(error);
        process.exit(1);
    });