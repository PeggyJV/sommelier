# Sommelier

Sommelier is a platform for running DeFi strategies managed by off-chain computation. It's a blockchain built with the [Cosmos SDK](https://github.com/cosmos/cosmos-sdk), and uses its own fork of the [Gravity Bridge](https://github.com/peggyjv/gravity-bridge) to enable cross-chain execution.

[![codecov](https://codecov.io/gh/peggyjv/sommelier/branch/main/graph/badge.svg)](https://codecov.io/gh/peggyjv/sommelier)
[![Go Report Card](https://goreportcard.com/badge/github.com/peggyjv/sommelier)](https://goreportcard.com/report/github.com/peggyjv/sommelier)
[![license](https://img.shields.io/github/license/peggyjv/sommelier.svg)](https://github.com/peggyjv/sommelier/blob/main/LICENSE)
[![LoC](https://tokei.rs/b1/github/peggyjv/sommelier)](https://github.com/peggyjv/sommelier)
[![GolangCI](https://golangci.com/badges/github.com/peggyjv/sommelier.svg)](https://golangci.com/r/github.com/peggyjv/sommelier)

## Talk to us!

We have active, helpful communities on Twitter, Discord, and Telegram.

* [Twitter](https://twitter.com/sommfinance)
* [Discord](https://discord.gg/gZzaPmDzUq)
* [Telegram](https://t.me/peggyvaults)

## Sommelier

The initial release of the Sommelier blockchain will consist of a standard cosmos-sdk chain and the recently completed [Gravity Bridge refactor](https://github.com/peggyjv/gravity-bridge).

### Gravity Bridge

The Gravity Bridge requires some additional pieces to be deployed to support it:

- [ ] [Ethereum Contract](https://github.com/PeggyJV/gravity-bridge/tree/main/solidity) and associated tooling
- [ ] Orchestrator/Relayer binaries built from the `go.mod` commit

### Steward

[Steward](https://github.com/peggyjv/steward) is a middleware between the Strategy Provider and the protocol that facilitates Cellar function calls. It's also a CLI that subsumes the functionality of `gorc`, and is used in this document to configure and run the orchestrator.

## Join the mainnet!

Running a validator node on the Sommelier mainnet requires three processes:

1. The validator node
2. The Gravity Bridge Orchestrator
3. [Steward](https://github.com/peggyjv/steward)

We also recommend running a local `geth` process. Though a public Ethereum API service can work, you may run into request limits because of the frequent requests made by the `orchestrator` process.

The Steward CLI now supports all of the same commands as `gorc` and is the recommended way to configure delegate keys for new validators and to run the Orchestrator. __The Steward CLI is used to run *both* the `steward` and `orchestrator` processes__. There are post-installation steps for the `steward` process outlined at the end of the installation steps below. These are required for your Steward to participate in the protocol. For more information on these setup steps for Steward, see [Validators Instructions for Setting Up Steward](https://github.com/PeggyJV/steward/blob/main/docs/02-StewardForValidators.md) in the Steward repository.

> NOTE: The Steward CLI and Steward itself are distinct concepts in this document. The Steward CLI is used to start both the `steward` and `orchestrator` processes, while "Steward" refers specifically to the `steward` process.

### Installation

```bash
# Create an installation directory
mkdir install && cd install

# Install Steward
wget https://github.com/PeggyJV/steward/releases/download/v2.0.2/steward && chmod +x * && sudo mv * /usr/bin

# Install Geth
wget https://gethstore.blob.core.windows.net/builds/geth-linux-amd64-1.10.23-d901d853.tar.gz && tar -xvf geth-linux-amd64-1.10.23-d901d853.tar.gz && sudo mv geth-linux-amd64-1.10.23-d901d853/geth /usr/bin/geth && rm -rf geth-linux-amd64-1.10.23-d901d853*

# Install Sommelier
wget https://github.com/PeggyJV/sommelier/releases/download/v4.0.1/sommelier_4.0.1_linux_amd64.tar.gz && tar -xf sommelier_4.0.1_linux_amd64.tar.gz && sudo mv sommelier /usr/bin && rm -rf sommelier_4.0.1_linux_amd64* LICENSE README.md

# Fetch systemd unit file examples
wget https://raw.githubusercontent.com/PeggyJV/sommelier/main/contrib/systemd/geth.service https://raw.githubusercontent.com/PeggyJV/sommelier/main/contrib/systemd/sommelier.service https://raw.githubusercontent.com/PeggyJV/sommelier/main/contrib/systemd/orchestrator.service https://raw.githubusercontent.com/PeggyJV/sommelier/main/contrib/systemd/steward.service

# Modify the unit files to fit your environment
nano geth.service
nano orchestrator.service
nano steward.service
nano sommelier.service

# And install them to systemd
sudo mv geth.service /etc/systemd/system/geth.service \
    && sudo mv orchestrator.service /etc/systemd/system/ \
    && sudo mv steward.service /etc/systemd/system/ \
    && sudo mv sommelier.service /etc/systemd/system/ \
    && sudo systemctl daemon-reload

# Start geth
sudo systemctl start geth && sudo journalctl -u geth -f

# Init steward/orchestrator configuration. Note that the steward and orchestrator processes share
# much of the same configuration fields, so we share the config.toml for convenience.
mkdir -p $HOME/steward && cd $HOME/steward
wget https://raw.githubusercontent.com/PeggyJV/sommelier/main/contrib/mainnet/sommelier-3/config.toml

# modify steward/orchestrator config for your environment
# You can use alchemy or infura endpoints as your RPC and don't necessarily need to sync the
# blockchain eth with geth
nano config.toml

# Initialize the validator files
sommelier init myval --chain-id sommelier-3
```

At this point you need to create orchestrator keys OR restore them if you already created them. __Please follow [these instructions](https://github.com/PeggyJV/steward/blob/main/docs/03-TheOrchestrator.md#setup) to create or restore these keys with the Steward CLI__, then return to this doc for steps to add them to your validator.

```bash
# restore your validator mnemonic to the sommelier binary
sommelier keys add validator --recover

# OR: create your validator mnemonic to the sommelier binary
sommelier keys add validator

# NOTE: at the end of this process you need to have:
# - a key named "orchestrator" with funds on the cosmos chain in the steward keystore
# - a key named "signer" with funds on connected ETH chain in the steward keystore
# - a key named "validator" with funds on the cosmos chain in the sommelier keystore

# Add the peers from contrib/mainnet/sommelier-3/peers.txt to the ~/.sommelier/config/config.toml file
nano ~/.sommelier/config/config.toml

# pull the genesis file
wget https://raw.githubusercontent.com/PeggyJV/sommelier/main/contrib/mainnet/sommelier-3/genesis.json -O $HOME/.sommelier/config/genesis.json

# start your sommelier node - note it may take a minute or two to sync all of the blocks
sudo systemctl start sommelier && sudo journalctl -u sommelier -f

# once your node is synced, create your validator
sommelier tx staking create-validator \
  --amount=1000000usomm \
  --pubkey=$(sommelier tendermint show-validator) \
  --moniker="MYMONIKER" \
  --chain-id="sommelier-3" \
  --commission-rate="0.10" \
  --commission-max-rate="0.20" \
  --commission-max-change-rate="0.01" \
  --min-self-delegation="1" \
  --gas 300000
  --fees="0usomm"
  --from=validator

# register delegate keys for eth and orchestrator keys
sommelier tx gravity set-delegate-keys \
    $(sommelier keys show validator --bech val -a) \ # validator address
    $(steward --config $HOME/steward/config.toml keys cosmos show orchestrator) \ # orchestrator address (this must be run manually and address extracted)
    $(steward --config $HOME/steward/config.toml keys eth show signer) \ # eth signer address
    $(steward --config $HOME/steward/config.toml sign-delegate-keys signer $(sommelier keys show validator --bech val -a)) \
    --chain-id sommelier-3 \
    --from validator \
    -y

# start the orchestrator. note that we are not yet starting steward
sudo systemctl start orchestrator && sudo journalctl -u orchestrator -f
```

At this point, you should have a running validator node and Orchestrator.

Now it's time to complete the setup for Steward. Please follow the detailed guide in [Validators Instructions for Setting Up Steward](https://github.com/PeggyJV/steward/blob/main/docs/02-StewardForValidators.md) and return here.

At this point you should have a server CA and server certificate for Steward, and your `config.toml` should be configured with those values. Now we can start the Steward service that we created during the other installation steps.

```bash
# start steward
sudo systemctl start steward && sudo journalctl -u steward -f
```

Once your Steward is running, ensure that its server endpoint is reachable over the internet. Then, if you haven't already, follow the steps outlined in the [Steward Registry repository](https://github.com/PeggyJV/steward-registry) to register your steward instance.

Your installation is complete! If you have any problems, please reach out in the validator lobby channels in Discord or Telegram.

### Actions

Now you can try the bridge!!

```bash
# send somm to ethereum
gorc cosmos-to-eth \
    --cosmos-phrase="$(jq -r '.orchestrator' ~/keys.json)" \
    --cosmos-grpc="http://localhost:9090" \
    --cosmos-denom="somm" \
    --amount="100000000" \
    --eth-destination=$(gorc --config $HOME/gorc/config.toml keys eth show signer) \
    --cosmos-prefix="cosmos"

# send goreli uniswap tokens to cosmos
gorc eth-to-cosmos \
    --ethereum-key="$(jq -r '.eth' ~/keys.json)" \
    --ethereum-rpc="http://localhost:8545" \
    --cosmos-prefix="cosmos" \
    --contract-address="$(jq -r '.gravity' ~/keys.json)" \
    --erc20-address="0x0000000000000000000000000000000000000000" \
    --amount="1.3530000" \
    --cosmos-destination="$(sommelier keys show orchestrator -a)"

```

## Notes:

### Genesis File Changes Necessary

```bash
# change stake to usomm
sed -i 's/stake/usomm/g' ~/.sommelier/config/genesis.json

# denom metadata
# TODO: add name and symbol here
jq -rMc '.app_state.bank.denom_metadata += [{"base": "usomm", display: "somm", "description": "A staking test token", "denom_units": [{"denom": "usomm", "exponent": 0}, {"denom": "somm", "exponent": 6}]}]' ~/.sommelier/config/genesis.json > ~/.sommelier/config/genesis.json

# gravity params
jq -rMc '.app_state.gravity.params.bridge_chain_id = "5"' ~/.sommelier/config/genesis.json > ~/.sommelier/config/genesis.json
```

### Deploy Peggy Contract

```bash
wget https://github.com/PeggyJV/gravity-bridge/releases/download/v0.1.21/Gravity.json
contract-deployer \
    --cosmos-node="http://localhost:26657" \
    --eth-node="http://localhost:8545" \
    --eth-privkey="0x0000000000000000000000000000000000000000000000000000000000000000" \
    --contract=Gravity.json \
    --test-mode=false
```

### Deploy Somm ERC20 representation

```bash
gorc deploy erc20 \
    --ethereum-key="0x0000000000000000000000000000000000000000000000000000000000000000" \
    --cosmos-grpc="http://localhost:9090" \
    --cosmos-prefix=cosmos \
    --cosmos-denom=usomm \
    --ethereum-rpc=http://localhost:8545 \
    --contract-address="0x8887F26882a3F920e40A91969D1A40D1Ef7efe10" \
    --erc20-name=usomm \
    --erc20-symbol=usomm \
    --erc20-decimals=6
```
