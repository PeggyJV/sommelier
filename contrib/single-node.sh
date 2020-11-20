#!/bin/sh

CHAINID=$1
GENACCT=$2

if [ -z "$1" ]; then
  echo "Need to input chain id..."
  exit 1
fi

if [ -z "$2" ]; then
  echo "Need to input genesis account address..."
  exit 1
fi

# Build genesis file incl account for passed address
coins="10000000000stake,100000000000samoleans"
sommelier init --chain-id $CHAINID $CHAINID
sommelier keys add validator --keyring-backend="test"
sommelier add-genesis-account $(sommelier keys show validator -a --keyring-backend="test") $coins
sommelier add-genesis-account $GENACCT $coins
sommelier gentx validator --keyring-backend="test" --chain-id $CHAINID
sommelier collect-gentxs

# Set proper defaults and change ports
sed -i 's#"tcp://127.0.0.1:26657"#"tcp://0.0.0.0:26657"#g' ~/.sommelier/config/config.toml
sed -i 's/timeout_commit = "5s"/timeout_commit = "1s"/g' ~/.sommelier/config/config.toml
sed -i 's/timeout_propose = "3s"/timeout_propose = "1s"/g' ~/.sommelier/config/config.toml
sed -i 's/index_all_keys = false/index_all_keys = true/g' ~/.sommelier/config/config.toml

# Start the sommelier
sommelier start --pruning=nothing
