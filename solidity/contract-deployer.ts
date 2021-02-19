import { Peggy } from "./typechain/Peggy";
import { TestERC20A } from "./typechain/TestERC20A";
import { TestERC20B } from "./typechain/TestERC20B";
import { TestERC20C } from "./typechain/TestERC20C";
import { ethers } from "ethers";
import fs from "fs";
import commandLineArgs from "command-line-args";
import axios, { AxiosError, AxiosRequestConfig, AxiosResponse } from "axios";
import { exit } from "process";

const args = commandLineArgs([
  // the ethernum node used to deploy the contract
  { name: "eth-node", type: String },
  // the cosmos node that will be used to grab the validator set via RPC (TODO),
  { name: "cosmos-node", type: String },
  // the Ethereum private key that will contain the gas required to pay for the contact deployment
  { name: "eth-privkey", type: String },
  // the peggy contract .json file
  { name: "contract", type: String },
  // test mode, if enabled this script deploys three ERC20 contracts for testing
  { name: "test-mode", type: String },
]);

// 4. Now, the deployer script hits a full node api, gets the Eth signatures of the valset from the latest block, and deploys the Ethereum contract.
//     - We will consider the scenario that many deployers deploy many valid peggy eth contracts.
// 5. The deployer submits the address of the peggy contract that it deployed to Ethereum.
//     - The peggy module checks the Ethereum chain for each submitted address, and makes sure that the peggy contract at that address is using the correct source code, and has the correct validator set.
type Validator = {
  power: number;
  ethereum_address: string;
};
type ValsetTypeWrapper = {
  type: string;
  value: Valset;
}
type Valset = {
  members: Validator[];
  nonce: number;
};
type ABCIWrapper = {
  jsonrpc: string;
  id: string;
  result: ABCIResponse;
};
type ABCIResponse = {
  response: ABCIResult
}
type ABCIResult = {
  code: number
  log: string,
  info: string,
  index: string,
  value: string,
  height: string,
  codespace: string,
};
type StatusWrapper = {
  jsonrpc: string,
  id: string,
  result: NodeStatus
};
type NodeInfo = {
  protocol_version: JSON,
  id: string,
  listen_addr: string,
  network: string,
  version: string,
  channels: string,
  moniker: string,
  other: JSON,
};
type SyncInfo = {
  latest_block_hash: string,
  latest_app_hash: string,
  latest_block_height: Number
  latest_block_time: string,
  earliest_block_hash: string,
  earliest_app_hash: string,
  earliest_block_height: Number,
  earliest_block_time: string,
  catching_up: boolean,
}
type NodeStatus = {
  node_info: NodeInfo,
  sync_info: SyncInfo,
  validator_info: JSON,
};

async function deploy() {
  const provider = await new ethers.providers.JsonRpcProvider(args["eth-node"]);
  let wallet = new ethers.Wallet(args["eth-privkey"], provider);

  if (args["test-mode"] == "True" || args["test-mode"] == "true") {
    console.log("Test mode, deploying ERC20 contracts");
    const { abi, bytecode } = getContractArtifacts("/peggy/solidity/artifacts/contracts/TestERC20A.sol/TestERC20A.json");
    const erc20Factory = new ethers.ContractFactory(abi, bytecode, wallet);
    const testERC20 = (await erc20Factory.deploy()) as TestERC20A;
    await testERC20.deployed();
    const erc20TestAddress = testERC20.address;
    console.log("ERC20 deployed at Address - ", erc20TestAddress);
    const { abi: abi1, bytecode: bytecode1 } = getContractArtifacts("/peggy/solidity/artifacts/contracts/TestERC20B.sol/TestERC20B.json");
    const erc20Factory1 = new ethers.ContractFactory(abi1, bytecode1, wallet);
    const testERC201 = (await erc20Factory1.deploy()) as TestERC20B;
    await testERC201.deployed();
    const erc20TestAddress1 = testERC201.address;
    console.log("ERC20 deployed at Address - ", erc20TestAddress1);
    const { abi: abi2, bytecode: bytecode2 } = getContractArtifacts("/peggy/solidity/artifacts/contracts/TestERC20C.sol/TestERC20C.json");
    const erc20Factory2 = new ethers.ContractFactory(abi2, bytecode2, wallet);
    const testERC202 = (await erc20Factory2.deploy()) as TestERC20C;
    await testERC202.deployed();
    const erc20TestAddress2 = testERC202.address;
    console.log("ERC20 deployed at Address - ", erc20TestAddress2);
  }
  const peggyIdString = await getPeggyId();
  const peggyId = ethers.utils.formatBytes32String(peggyIdString);

  console.log("Starting Peggy contract deploy");
  const { abi, bytecode } = getContractArtifacts(args["contract"]);
  const factory = new ethers.ContractFactory(abi, bytecode, wallet);

  console.log("About to get latest Peggy valset");
  const latestValset = await getLatestValset();

  let eth_addresses = [];
  let powers = [];
  let powers_sum = 0;
  // this MUST be sorted uniformly across all components of Peggy in this
  // case we perform the sorting in module/x/peggy/keeper/types.go to the
  // output of the endpoint should always be sorted correctly. If you're
  // having strange problems with updating the validator set you should go
  // look there.
  for (let i = 0; i < latestValset.members.length; i++) {
    if (latestValset.members[i].ethereum_address == null) {
      continue;
    }
    eth_addresses.push(latestValset.members[i].ethereum_address);
    powers.push(latestValset.members[i].power);
    powers_sum += latestValset.members[i].power;
  }

  // 66% of uint32_max
  let vote_power = 2834678415;
  if (powers_sum < vote_power) {
    console.log("Refusing to deploy! Incorrect power! Please inspect the validator set below")
    console.log("If less than 66% of the current voting power has unset Ethereum Addresses we refuse to deploy")
    console.log(latestValset)
    exit(1)
  }

  const peggy = (await factory.deploy(
    // todo generate this randomly at deployment time that way we can avoid
    // anything but intentional conflicts
    peggyId,
    vote_power,
    eth_addresses,
    powers
  )) as Peggy;

  await peggy.deployed();
  console.log("Peggy deployed at Address - ", peggy.address);
  await submitPeggyAddress(peggy.address);
}

function getContractArtifacts(path: string): { bytecode: string; abi: string } {
  var { bytecode, abi } = JSON.parse(fs.readFileSync(path, "utf8").toString());
  return { bytecode, abi };
}
const decode = (str: string):string => Buffer.from(str, 'base64').toString('binary');

async function getLatestValset(): Promise<Valset> {
  let block_height_request_string = args["cosmos-node"] + '/status';
  let block_height_response = await axios.get(block_height_request_string);
  let info: StatusWrapper = await block_height_response.data;
  let block_height = info.result.sync_info.latest_block_height;
  if (info.result.sync_info.catching_up) {
    console.log("This node is still syncing! You can not deploy using this validator set!");
    exit(1);
  }
  let request_string = args["cosmos-node"] + "/abci_query"
  let response = await axios.get(request_string, {params: {
    path: "\"/custom/peggy/currentValset/\"",
    height: block_height,
    prove: "false",
  }});
  let valsets: ABCIWrapper = await response.data;
  console.log(decode(valsets.result.response.value));
  let valset: ValsetTypeWrapper = JSON.parse(decode(valsets.result.response.value))
  return valset.value;
}
async function getPeggyId(): Promise<string> {
  let block_height_request_string = args["cosmos-node"] + '/status';
  let block_height_response = await axios.get(block_height_request_string);
  let info: StatusWrapper = await block_height_response.data;
  let block_height = info.result.sync_info.latest_block_height;
  if (info.result.sync_info.catching_up) {
    console.log("This node is still syncing! You can not deploy using this peggyID!");
    exit(1);
  }
  let request_string = args["cosmos-node"] + "/abci_query"
  let response = await axios.get(request_string, {params: {
    path: "\"/custom/peggy/peggyID/\"",
    height: block_height,
    prove: "false",
  }});
  let peggyIDABCIResponse: ABCIWrapper = await response.data;
  let peggyID: string = JSON.parse(decode(peggyIDABCIResponse.result.response.value))
  return peggyID;

}

async function submitPeggyAddress(address: string) {}

async function main() {
  await deploy();
}

main();
