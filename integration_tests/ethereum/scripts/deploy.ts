import {ethers} from "hardhat";
import {BigNumberish, Signer} from "ethers";
import {SignerWithAddress} from "@nomiclabs/hardhat-ethers/signers";


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

    let abiEncoded = ethers.utils.defaultAbiCoder.encode(
        ["bytes32", "bytes32", "uint256", "address[]", "uint256[]"],
        [gravityId, methodName, valsetNonce, validators, powers]
    );

    return ethers.utils.keccak256(abiEncoded);
}

export async function deployTestnetContract() {
    let valAddresses: string[] = [
        '0xd312f0f1B39D54Db2829537595fC1167B14d4b34',
        '0x7bE2a04df4b9C3227928147461e19158eB2B11d1',
        '0xb8c6886FDDa38adaa0F416722dd5554886C43055',
        '0x14fdAC734De10065093C4Ed4a83C41638378005A'
    ];
    let powers: number[] = [1073741824,1073741824,1073741824,1073741824];
    let powerThreshold: number = 6666;

    let {gravity, checkpoint} = deployContracts("gravitytest", valAddresses, powers, powerThreshold);

 }

export async function deployContracts(
    gravityId: string = "foo",
    valAddresses: string[],
    powers: number[],
    powerThreshold: number
) {
    const Gravity = await ethers.getContractFactory("Gravity");
    const checkpoint = makeCheckpoint(valAddresses, powers, 0, gravityId);
    const gravity = (await Gravity.deploy(
        gravityId,
        powerThreshold,
        valAddresses,
        powers
    ));

    await gravity.deployed();
    console.log(`gravity contract deployed at ${gravity.address}`)

    return { gravity, checkpoint };
}
