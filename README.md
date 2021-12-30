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


## Join the mainnet!

### Installation

```bash 
# Create an installation directory
mkdir install && cd install

# Install Orchestrator
wget https://github.com/PeggyJV/gravity-bridge/releases/download/v0.1.21/contract-deployer https://github.com/PeggyJV/gravity-bridge/releases/download/v0.1.21/orchestrator https://github.com/PeggyJV/gravity-bridge/releases/download/v0.1.21/relayer https://github.com/PeggyJV/gravity-bridge/releases/download/v0.1.21/gorc && chmod +x * && sudo mv * /usr/bin

# Install Geth
wget https://gethstore.blob.core.windows.net/builds/geth-linux-amd64-1.10.4-aa637fd3.tar.gz && tar -xvf geth-linux-amd64-1.10.4-aa637fd3.tar.gz && sudo mv geth-linux-amd64-1.10.4-aa637fd3/geth /usr/bin/geth && rm -rf geth-linux-amd64-1.10.4-aa637fd3*

# Install Sommelier
wget https://github.com/PeggyJV/sommelier/releases/download/v2.0.0/sommelier_2.0.0_linux_amd64.tar.gz && tar -xf sommelier_2.0.0_linux_amd64.tar.gz && sudo mv sommelier /usr/bin && rm -rf sommelier_2.0.0_linux_amd64* LICENSE README.md

# Fetch systemd unit files
wget https://raw.githubusercontent.com/PeggyJV/sommelier/main/contrib/systemd/geth.service https://raw.githubusercontent.com/PeggyJV/sommelier/main/contrib/systemd/gorc.service https://raw.githubusercontent.com/PeggyJV/sommelier/main/contrib/systemd/sommelier.service

# Modify the unit files to fit your environment
nano geth.service
nano gorc.service
nano sommelier.service

# And install them to systemd
sudo mv geth.service /etc/systemd/system/geth.service && sudo mv gorc.service /etc/systemd/system/ && sudo mv sommelier.service /etc/systemd/system/ && sudo systemctl daemon-reload

# Start geth
sudo systemctl start geth && sudo journalctl -u geth -f

# Init gorc configuration
mkdir -p $HOME/gorc && cd $HOME/gorc
wget https://raw.githubusercontent.com/PeggyJV/sommelier/main/contrib/mainnet/sommelier-2/config.toml

# modify gorc config for your environment
nano config.toml

# Initialize the validator files
sommelier init myval --chain-id sommelier-2

# create/restore 2 cosmos keys and 1 ethereum key
# NOTE: be sure to save the mnemonics and eth private key

# restore orchestrator key with gorc 
gorc --config $HOME/gorc/config.toml keys cosmos recover orchestrator "{menmonic}"

# OR: create orchestrator key with gorc 
gorc --config $HOME/gorc/config.toml keys cosmos add orchestrator

# restore eth key 

# EITHER: restore eth priv key from metamask with gorc 
gorc --config $HOME/gorc/config.toml keys eth import signer "0x0000..."

# OR: restore eth mnemonic with gorc
gorc --config $HOME/gorc/config.toml keys eth recover signer "{menomonic}"

# OR: create eth key with gorc 
gorc --config $HOME/gorc/config.toml keys eth add signer

# restore your validator mnemonic to the sommelier binary
sommelier keys add validator --recover

# OR: create your validator mnemonic to the sommelier binary
sommelier keys add validator

# NOTE: at the end of this process you need to have:
# - a key named "orchestrator" with funds on the cosmos chain in the gorc keystore
# - a key named "signer" with funds on connected ETH chain in the gorc keystore
# - a key named "validator" with funds on the cosmos chain in the sommelier keystore

# Add the peers from contrib/mainnet/sommelier-2/peers.txt to the ~/.sommelier/config/config.toml file
nano ~/.sommelier/config/config.toml

# pull the genesis file 
wget https://raw.githubusercontent.com/PeggyJV/sommelier/main/contrib/mainnet/sommelier-2/genesis.json -O $HOME/.sommelier/config/genesis.json

# start your sommelier node - note it may take a minute or two to sync all of the blocks
sudo systemctl start sommelier && sudo journalctl -u sommelier -f

# once your node is synced, create your validator 
sommelier tx staking create-validator \
  --amount=1000000usomm \
  --pubkey=$(sommelier tendermint show-validator) \
  --moniker="MYMONIKER" \
  --chain-id="sommelier-2" \
  --commission-rate="0.10" \
  --commission-max-rate="0.20" \
  --commission-max-change-rate="0.01" \
  --min-self-delegation="1" \
  --gas 300000 \
  --fees="0usomm" \
  --from=validator

# register delegate keys for eth and orchestrator keys
sommelier tx gravity set-delegate-keys \
    $(sommelier keys show validator --bech val -a) \ # validator address
    $(gorc --config $HOME/gorc/config.toml keys cosmos show orchestrator) \ # orchestrator address (this must be run manually and address extracted)
    $(gorc --config $HOME/gorc/config.toml keys eth show signer) \ # eth signer address
    $(gorc --config $HOME/gorc/config.toml sign-delegate-keys signer $(sommelier keys show validator --bech val -a)) \
    --chain-id sommelier-2 \
    --from validator \
    -y

# start the orchestrator
sudo systemctl start gorc && sudo journalctl -u gorc -f
```

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
