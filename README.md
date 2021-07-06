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


```bash 
# Create an installation directory
mkdir install && cd install

# Install Orchestrator
wget https://github.com/PeggyJV/gravity-bridge/releases/download/v0.1.4/client https://github.com/PeggyJV/gravity-bridge/releases/download/v0.1.4/contract-deployer https://github.com/PeggyJV/gravity-bridge/releases/download/v0.1.4/orchestrator https://github.com/PeggyJV/gravity-bridge/releases/download/v0.1.4/relayer && chmod +x * && sudo mv * /usr/bin

# Install Geth
wget https://gethstore.blob.core.windows.net/builds/geth-linux-amd64-1.10.4-aa637fd3.tar.gz && tar -xvf geth-linux-amd64-1.10.4-aa637fd3.tar.gz && sudo mv geth-linux-amd64-1.10.4-aa637fd3/geth /usr/bin/geth && rm -rf geth-linux-amd64-1.10.4-aa637fd3*

# Install Sommelier
wget https://github.com/PeggyJV/sommelier/releases/download/v0.1.1/sommelier_0.1.1_linux_amd64.tar.gz && tar -xf sommelier_0.1.1_linux_amd64.tar.gz && sudo mv sommelier /usr/bin && rm -rf sommelier_0.1.1_linux_amd64* LICENSE README.md
```

## Genesis File Fixes

```bash
# denom metadata
jq '.app_state.bank.denom_metadata += [{"base": "somm", display: "somm", "description": "A non-staking test token", "denom_units": [{"denom": "somm", "exponent": 6}]}, {"base": "stake", display: "stake", "description": "A staking test token", "denom_units": [{"denom": "stake", "exponent": 6}]}]' ~/.sommelier/config/genesis.json > ~/.sommelier/config/edited-genesis.json
mv ~/.sommelier/config/edited-genesis.json ~/.sommelier/config/genesis.json

# gravity params
jq '.app_state.gravity.params.bridge_chain_id = "5"' ~/.sommelier/config/genesis.json > ~/.sommelier/config/edited-genesis.json
mv ~/.sommelier/config/edited-genesis.json ~/.sommelier/config/genesis.json
```

## Deploy Peggy Contract
```bash
wget https://github.com/PeggyJV/gravity-bridge/releases/download/v0.1.4/Gravity.json
contract-deployer \
    --cosmos-node="http://localhost:26657" \
    --eth-node="http://localhost:8545" \
    --eth-privkey="0x0000000000000000000000000000000000000000000000000000000000000000" \
    --contract=Gravity.json \
    --test-mode=false
```

```bash
orchestrator \
    --cosmos-phrase="electric gesture amount absent someone pudding waste you pilot truth pioneer nose surprise flavor ask cost art grit ladder girl detect height pet primary" \
    --ethereum-key=0xe24bfff133d5f90f046147cb4d23a6e68ae277e6855d58c5b1f1208822a68dec \
    --cosmos-grpc=http://localhost:9090 \
    --address-prefix=cosmos \
    --ethereum-rpc=http://localhost:8545 \
    --fees=somm \
    --contract-address=0x0D08d893bf138D958a8e215880cBd307C1946dd5
```
    0x0D08d893bf138D958a8e215880cBd307C1946dd5