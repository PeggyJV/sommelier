# Sommelier

Sommelier is a coprocessor blockchain for Ethereum DeFi.

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


## Join the testnet!

### Installation

```bash 
# Create an installation directory
mkdir install && cd install

# Install Orchestrator
wget https://github.com/PeggyJV/gravity-bridge/releases/download/v0.1.5/client https://github.com/PeggyJV/gravity-bridge/releases/download/v0.1.5/contract-deployer https://github.com/PeggyJV/gravity-bridge/releases/download/v0.1.5/orchestrator https://github.com/PeggyJV/gravity-bridge/releases/download/v0.1.5/relayer && chmod +x * && sudo mv * /usr/bin

# Install Geth
wget https://gethstore.blob.core.windows.net/builds/geth-linux-amd64-1.10.4-aa637fd3.tar.gz && tar -xvf geth-linux-amd64-1.10.4-aa637fd3.tar.gz && sudo mv geth-linux-amd64-1.10.4-aa637fd3/geth /usr/bin/geth && rm -rf geth-linux-amd64-1.10.4-aa637fd3*

# Install Sommelier
wget https://github.com/PeggyJV/sommelier/releases/download/v0.1.2/sommelier_0.1.2_linux_amd64.tar.gz && tar -xf sommelier_0.1.2_linux_amd64.tar.gz && sudo mv sommelier /usr/bin && rm -rf sommelier_0.1.2_linux_amd64* LICENSE README.md

# Fetch systemd unit files
wget https://raw.githubusercontent.com/PeggyJV/sommelier/main/contrib/systemd/geth.goerli.service https://raw.githubusercontent.com/PeggyJV/sommelier/main/contrib/systemd/orchestrator.service https://raw.githubusercontent.com/PeggyJV/sommelier/main/contrib/systemd/sommelier.service

# Modify the unit files to fit your environment
nano geth.goerli.service
nano orchestrator.service
nano sommelier.service

# And install them to systemd
sudo mv geth.goerli.service /etc/systemd/system/geth.service && sudo mv orchestrator.service /etc/systemd/system/ && sudo mv sommelier.service /etc/systemd/system/ && sudo systemctl daemon-reload

# Start geth
sudo systemctl start geth && sudo journalctl -u geth -f

# Initialize the validator files
sommelier init myval --chain-id sommtest-1

# create 2 cosmos keys and 1 ethereum key
# NOTE: be sure to save the mnemonics and eth private key as you will need 
# these as arguements for the orchestrator and the client software

sommelier keys add validator # --keyring-backend test
sommelier keys add orchestrator # --keyring-backend test
# use metamask to create a new account/use exsting (goreli network)
# there is an "export private key" function that returns the key for this
# note: i've saved these in a keys.json file for easy access

# go ask Jack for some testnet $$$$ for both cosmos addresses and for some goreli eth

# Add the peers from contrib/testnets/peers.txt to the ~/.sommelier/config/config.toml file
nano ~/.sommelier/config/config.toml

# pull the genesis file 
wget https://raw.githubusercontent.com/PeggyJV/sommelier/main/contrib/testnets/sommtest-1/genesis.json -O 
mv genesis.json ~/.sommelier/config

# start your sommelier node - note it may take a minute or two to sync all of the blocks
sudo systemctl start sommelier && sudo journalctl -u sommelier -f

# once your node is synced, create your validator 
sommelier tx staking create-validator \
  --amount=10000000stake \
  --pubkey=$(sommelier tendermint show-validator) \
  --moniker="mymoniker" \
  --chain-id="sommtest-1" \
  --commission-rate="0.10" \
  --commission-max-rate="0.20" \
  --commission-max-change-rate="0.01" \
  --min-self-delegation="1000000" \
  --gas="auto" \
  --gas-prices="0.025stake" \
  --from=validator

# register delegate keys for eth and orchestrator keys
sommelier tx gravity set-delegate-keys $(sommelier keys show validator --bech val -a) $(sommelier keys show orchestrator -a) 0x0000000000000000000000000000000000000000

# edit the orchestrator unit file to include private keys for cosmos and eth as well as the proper contract address
# then start it
sudo nano /etc/systemd/system/orchestrator.service && sudo systemctl daemon-reload && sudo systemctl start orchestrator && sudo journalctl -u orchestrator -f
```

### Actions

Now you can try the bridge!!

```bash
# send somm to ethereum
client cosmos-to-eth \
    --cosmos-phrase="$(jq -r '.orchestrator' ~/keys.json)" \
    --cosmos-grpc="http://localhost:9090" \
    --cosmos-denom="somm" \
    --amount="100000000" \
    --eth-destination="0x25E17020b70A46ee0A4F6cb7259E1d835e4310Bb" \
    --cosmos-prefix="cosmos"

# send goreli uniswap tokens to cosmos
client eth-to-cosmos \
    --ethereum-key="$(jq -r '.eth' ~/keys.json)" \
    --ethereum-rpc="http://localhost:8545" \
    --cosmos-prefix="cosmos" \
    --contract-address="$(jq -r '.gravity' ~/keys.json)" \
    --erc20-address="0x1f9840a85d5af5bf1d1762f925bdaddc4201f984" \
    --amount="1.3530000" \
    --cosmos-destination="$(sommelier keys show orchestrator -a)"
    
```

## Notes:

### Genesis File Changes Necessary

```bash
# denom metadata
jq '.app_state.bank.denom_metadata += [{"base": "somm", display: "somm", "description": "A non-staking test token", "denom_units": [{"denom": "somm", "exponent": 6}]}, {"base": "stake", display: "stake", "description": "A staking test token", "denom_units": [{"denom": "stake", "exponent": 6}]}]' ~/.sommelier/config/genesis.json > ~/.sommelier/config/edited-genesis.json
mv ~/.sommelier/config/edited-genesis.json ~/.sommelier/config/genesis.json

# gravity params
jq '.app_state.gravity.params.bridge_chain_id = "5"' ~/.sommelier/config/genesis.json > ~/.sommelier/config/edited-genesis.json
mv ~/.sommelier/config/edited-genesis.json ~/.sommelier/config/genesis.json
```

### Deploy Peggy Contract

```bash
wget https://github.com/PeggyJV/gravity-bridge/releases/download/v0.1.5/Gravity.json
contract-deployer \
    --cosmos-node="http://localhost:26657" \
    --eth-node="http://localhost:8545" \
    --eth-privkey="0x0000000000000000000000000000000000000000000000000000000000000000" \
    --contract=Gravity.json \
    --test-mode=false
```

### Deploy Somm ERC20 representation

```bash
client deploy-erc20-representation \
    --ethereum-key="0x0000000000000000000000000000000000000000000000000000000000000000"  
    --cosmos-grpc="http://localhost:9090" 
    --cosmos-prefix=cosmos 
    --cosmos-denom=somm 
    --ethereum-rpc=http://localhost:8545 
    --contract-address="0x0000000000000000000000000000000000000000" 
    --erc20-name=somm 
    --erc20-symbol=somm 
    --erc20-decimals=6
```
