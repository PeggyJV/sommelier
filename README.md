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
wget https://github.com/PeggyJV/gravity-bridge/releases/download/v0.1.12/client https://github.com/PeggyJV/gravity-bridge/releases/download/v0.1.12/contract-deployer https://github.com/PeggyJV/gravity-bridge/releases/download/v0.1.12/orchestrator https://github.com/PeggyJV/gravity-bridge/releases/download/v0.1.12/relayer https://github.com/PeggyJV/gravity-bridge/releases/download/v0.1.12/gorc && chmod +x * && sudo mv * /usr/bin

# Install Geth
wget https://gethstore.blob.core.windows.net/builds/geth-linux-amd64-1.10.4-aa637fd3.tar.gz && tar -xvf geth-linux-amd64-1.10.4-aa637fd3.tar.gz && sudo mv geth-linux-amd64-1.10.4-aa637fd3/geth /usr/bin/geth && rm -rf geth-linux-amd64-1.10.4-aa637fd3*

# Install Sommelier
wget https://github.com/PeggyJV/sommelier/releases/download/v0.1.6/sommelier_0.1.6_linux_amd64.tar.gz && tar -xf sommelier_0.1.6_linux_amd64.tar.gz && sudo mv sommelier /usr/bin && rm -rf sommelier_0.1.6_linux_amd64* LICENSE README.md

# Fetch systemd unit files
wget https://raw.githubusercontent.com/PeggyJV/sommelier/main/contrib/systemd/geth.goerli.service https://raw.githubusercontent.com/PeggyJV/sommelier/main/contrib/systemd/gorc.service https://raw.githubusercontent.com/PeggyJV/sommelier/main/contrib/systemd/sommelier.service

# Modify the unit files to fit your environment
nano geth.goerli.service
nano gorc.service
nano sommelier.service

# And install them to systemd
sudo mv geth.goerli.service /etc/systemd/system/geth.service && sudo mv gorc.service /etc/systemd/system/ && sudo mv sommelier.service /etc/systemd/system/ && sudo systemctl daemon-reload

# Start geth
sudo systemctl start geth && sudo journalctl -u geth -f

# Init gorc configuration
mkdir -p $HOME/gorc && cd $HOME/gorc
wget https://raw.githubusercontent.com/PeggyJV/sommelier/main/contrib/testnets/sommtest-2/config.toml

# modify gorc config for your environment
nano config.toml

# Initialize the validator files
sommelier init myval --chain-id sommtest-2

# create/restore 2 cosmos keys and 1 ethereum key
# NOTE: be sure to save the mnemonics and eth private key

# restore orchestrator key with gorc 
gorc --config $HOME/gorc/config.toml keys cosmos restore orchestrator "{menmonic}"

# restore eth priv key from metamask with gorc 
gorc --config $HOME/gorc/config.toml keys eth import signer "0x0000..."

# restore eth mnemonic with gorc
gorc --config $HOME/gorc/config.toml keys eth restore signer "{menomonic}"

# restore your validator mnemonic to the sommelier binary
sommelier keys add validator --recover 

# NOTE: at the end of this process you need to have:
# - a key named "orchestrator" with funds on the cosmos chain in the gorc keystore
# - a key named "signer" with funds on connected ETH chain in the gorc keystore
# - a key named "validator" with funds on the cosmos chain in the sommelier keystore

# Add the peers from contrib/testnets/peers.txt to the ~/.sommelier/config/config.toml file
nano ~/.sommelier/config/config.toml

# pull the genesis file 
wget https://raw.githubusercontent.com/PeggyJV/sommelier/main/contrib/testnets/sommtest-2/genesis.json -O $HOME/.sommelier/config/genesis.json

# start your sommelier node - note it may take a minute or two to sync all of the blocks
sudo systemctl start sommelier && sudo journalctl -u sommelier -f

# once your node is synced, create your validator 
sommelier tx staking create-validator \
  --amount=1000000000000usomm \
  --pubkey=$(sommelier tendermint show-validator) \
  --moniker="mymoniker" \
  --chain-id="sommtest-2" \
  --commission-rate="0.10" \
  --commission-max-rate="0.20" \
  --commission-max-change-rate="0.01" \
  --min-self-delegation="1000000" \
  --gas="auto" \
  --gas-prices="0.025usomm" \
  --from=validator

# register delegate keys for eth and orchestrator keys
sommelier tx gravity set-delegate-keys \
    $(sommelier keys show validator --bech val -a) \               # validator address
    $(sommelier keys show orchestrator -a) \                       # orchestrator address
    $(gorc --config $HOME/gorc/config.toml keys eth show signer) \ # eth signer address
    $(gorc --config $HOME/gorc/config sign-delegate-keys signer $(sommelier keys show validator --bech val -a) 0) \ 
    --chain-id sommtest-2 \ 
    --from validator \ 
    --fees 25000usomm -y

# edit the orchestrator unit file to include private keys for cosmos and eth as well as the proper contract address
# then start it
sudo systemctl start gorc && journalctl -u gorc -f
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
    --eth-destination=$(gorc --config $HOME/gorc/config.toml keys eth show signer) \
    --cosmos-prefix="cosmos"

# send goreli uniswap tokens to cosmos
client eth-to-cosmos \
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
# denom metadata
jq '.app_state.bank.denom_metadata += [{"base": "usomm", display: "somm", "description": "A staking test token", "denom_units": [{"denom": "usomm", "exponent": 0}, {"denom": "somm", "exponent": 6}]}]' ~/.sommelier/config/genesis.json > ~/.sommelier/config/edited-genesis.json
mv ~/.sommelier/config/edited-genesis.json ~/.sommelier/config/genesis.json

# gravity params
jq '.app_state.gravity.params.bridge_chain_id = "5"' ~/.sommelier/config/genesis.json > ~/.sommelier/config/edited-genesis.json
mv ~/.sommelier/config/edited-genesis.json ~/.sommelier/config/genesis.json
```

### Deploy Peggy Contract

```bash
wget https://github.com/PeggyJV/gravity-bridge/releases/download/v0.1.12/Gravity.json
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
    --ethereum-key="0x0000000000000000000000000000000000000000000000000000000000000000" \
    --cosmos-grpc="http://localhost:9090" \
    --cosmos-prefix=cosmos \
    --cosmos-denom=usomm \
    --ethereum-rpc=http://localhost:8545 \
    --contract-address="0x0000000000000000000000000000000000000000" \
    --erc20-name=usomm \
    --erc20-symbol=usomm \ 
    --erc20-decimals=6 \
```
