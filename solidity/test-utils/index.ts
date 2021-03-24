import { Gravity } from "../typechain/Gravity";
import { TestERC20A } from "../typechain/TestERC20A";
import { ethers } from "hardhat";
import { makeCheckpoint, signHash, getSignerAddresses } from "./pure";
import { Signer } from "ethers";

type DeployContractsOptions = {
  corruptSig?: boolean;
};

export async function deployContracts(
  gravityId: string = "foo",
  validators: Signer[],
  powers: number[],
  powerThreshold: number,
  opts?: DeployContractsOptions
) {
  const TestERC20 = await ethers.getContractFactory("TestERC20A");
  const testERC20 = (await TestERC20.deploy()) as TestERC20A;

  const Gravity = await ethers.getContractFactory("Gravity");

  const valAddresses = await getSignerAddresses(validators);

  const checkpoint = makeCheckpoint(valAddresses, powers, 0, gravityId);

  const gravity = (await Gravity.deploy(
    gravityId,
    powerThreshold,
    valAddresses,
    powers
  )) as Gravity;

  await gravity.deployed();

  return { gravity, testERC20, checkpoint };
}
