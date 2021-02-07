#!/bin/sh
# USAGE: ./two-node-net skip

# stop processes

# Constants
CWD=$(pwd)
CHAINID="somm"
CHAINDIR="$CWD/data"
SOM=sommelier
FED=oracle-feeder
hdir="$CHAINDIR/$CHAINID"

# Folders for nodes
n0dir="$hdir/n0"
n1dir="$hdir/n1"
n2dir="$hdir/n2"

# Home flag for folder
home0="--home $n0dir"
home1="--home $n1dir"
home2="--home $n2dir"

# Config directories for nodes
n0cfgDir="$n0dir/config"
n1cfgDir="$n1dir/config"
n2cfgDir="$n2dir/config"

# Config files for nodes
n0cfg="$n0cfgDir/config.toml"
n1cfg="$n1cfgDir/config.toml"
n2cfg="$n2cfgDir/config.toml"

# Config files for feeders
fd0cfg="$n0dir/config.yaml"
fd1cfg="$n1dir/config.yaml"
fd2cfg="$n2dir/config.yaml"

# App config files for nodes
n0app="$n0cfgDir/app.toml"
n1app="$n1cfgDir/app.toml"
n2app="$n2cfgDir/app.toml"

# Common flags
kbt="--keyring-backend test"
cid="--chain-id $CHAINID"

# Ensure user understands what will be deleted
if [[ -d $SIGNER_DATA ]] && [[ ! "$1" == "skip" ]]; then
  read -p "$0 will delete \$(pwd)/data folder. Do you wish to continue? (y/n): " -n 1 -r
  echo 
  if [[ ! $REPLY =~ ^[Yy]$ ]]; then
      exit 1
  fi
fi

echo "Creating 3x $SOM validators with chain-id=$CHAINID..."

# Build genesis file incl account for passed address
coins="100000000000stake,100000000000samoleans"

# Initialize the 2 home directories and add some keys
$SOM $home0 $cid init n0 &>/dev/null
$SOM $home0 keys add val $kbt &>/dev/null
$SOM $home1 $cid init n1 &>/dev/null
$SOM $home1 keys add val $kbt &>/dev/null
$SOM $home2 $cid init n2 &>/dev/null
$SOM $home2 keys add val $kbt &>/dev/null

# Add some keys and init feeder configs
$FED $home0 config init &>/dev/null
$FED $home0 keys add feeder &>/dev/null 
$FED $home1 config init &>/dev/null
$FED $home1 keys add feeder &>/dev/null
$FED $home2 config init &>/dev/null
$FED $home2 keys add feeder &>/dev/null

# Add addresses to genesis
$SOM $home0 add-genesis-account $($SOM $home0 keys show val -a $kbt) $coins &>/dev/null
$SOM $home0 add-genesis-account $($FED $home0 keys show feeder) $coins &>/dev/null
$SOM $home0 add-genesis-account $($SOM $home1 keys show val -a $kbt) $coins &>/dev/null
$SOM $home0 add-genesis-account $($FED $home1 keys show feeder) $coins &>/dev/null
$SOM $home0 add-genesis-account $($SOM $home2 keys show val -a $kbt) $coins &>/dev/null
$SOM $home0 add-genesis-account $($FED $home2 keys show feeder) $coins &>/dev/null

# Copy genesis around to sign
cp $n0cfgDir/genesis.json $n1cfgDir/genesis.json
cp $n0cfgDir/genesis.json $n2cfgDir/genesis.json

# Create gentxs and collect them in n0
$SOM $home0 gentx val 100000000000stake $kbt $cid &>/dev/null
$SOM $home1 gentx val 100000000000stake $kbt $cid &>/dev/null
$SOM $home2 gentx val 100000000000stake $kbt $cid &>/dev/null
cp $n1cfgDir/gentx/*.json $n0cfgDir/gentx/
cp $n2cfgDir/gentx/*.json $n0cfgDir/gentx/
$SOM $home0 collect-gentxs &>/dev/null

# Copy genesis file into n1 and n2s
cp $n0cfgDir/genesis.json $n1cfgDir/genesis.json
cp $n0cfgDir/genesis.json $n2cfgDir/genesis.json

# Switch sed command in the case of linux
SED="sed -i ''"
if [ `uname` = 'Linux' ]; then
   SED="sed -i"
fi

# Change ports on n0 val
$SED 's#"tcp://127.0.0.1:26657"#"tcp://0.0.0.0:26657"#g' $n0cfg
$SED 's#addr_book_strict = true#addr_book_strict = false#g' $n0cfg
$SED 's#allow_duplicate_ip = false#allow_duplicate_ip = true#g' $n0cfg
$SED 's#external_address = ""#external_address = "tcp://127.0.0.1:26657"#g' $n0cfg

# Change ports on n1 val
$SED 's#"tcp://127.0.0.1:26657"#"tcp://0.0.0.0:26667"#g' $n1cfg
$SED 's#"tcp://0.0.0.0:26656"#"tcp://0.0.0.0:26666"#g' $n1cfg
$SED 's#"localhost:6060"#"localhost:6061"#g' $n1cfg
$SED 's#"0.0.0.0:9090"#"0.0.0.0:9091"#g' $n1app
$SED 's#log_level = "main:info,state:info,statesync:info,*:error"#log_level = "info"#g' $n1cfg
$SED 's#addr_book_strict = true#addr_book_strict = false#g' $n1cfg
$SED 's#external_address = ""#external_address = "tcp://127.0.0.1:26667"#g' $n1cfg
$SED 's#allow_duplicate_ip = false#allow_duplicate_ip = true#g' $n1cfg

# Change ports on n1 feeder
$SED 's#http://localhost:9090#http://localhost:9091#g' $fd1cfg
$SED 's#http://http://localhost:26657#http://http://localhost:26667#g' $fd1cfg

# Change ports on n2 val
$SED 's#addr_book_strict = true#addr_book_strict = false#g' $n2cfg
$SED 's#external_address = ""#external_address = "tcp://127.0.0.1:26677"#g' $n2cfg
$SED 's#"tcp://127.0.0.1:26657"#"tcp://0.0.0.0:26677"#g' $n2cfg
$SED 's#"tcp://0.0.0.0:26656"#"tcp://0.0.0.0:26676"#g' $n2cfg
$SED 's#"localhost:6060"#"localhost:6062"#g' $n2cfg
$SED 's#"0.0.0.0:9090"#"0.0.0.0:9092"#g' $n2app
$SED 's#allow_duplicate_ip = false#allow_duplicate_ip = true#g' $n2cfg
$SED 's#log_level = "main:info,state:info,statesync:info,*:error"#log_level = "info"#g' $n2cfg

# Change ports on n2 feeder
$SED 's#http://localhost:9090#http://localhost:9092#g' $fd1cfg
$SED 's#http://http://localhost:26657#http://http://localhost:26677#g' $fd1cfg

# Set peers for all three nodes
peer0="$($SOM $home0 tendermint show-node-id)@127.0.0.1:26656"
peer1="$($SOM $home1 tendermint show-node-id)@127.0.0.1:26666"
peer2="$($SOM $home2 tendermint show-node-id)@127.0.0.1:26676"
sed -i '' 's#persistent_peers = ""#persistent_peers = "'$peer1','$peer2'"#g' $n0cfg
sed -i '' 's#persistent_peers = ""#persistent_peers = "'$peer0','$peer2'"#g' $n1cfg
sed -i '' 's#persistent_peers = ""#persistent_peers = "'$peer0','$peer1'"#g' $n2cfg

# Start the sommelier instances
echo "Starting nodes..."
$SOM $home0 start --pruning=nothing --grpc.address="0.0.0.0:9090" > $hdir.n0.log 2>&1 &
$SOM $home1 start --pruning=nothing --grpc.address="0.0.0.0:9091" > $hdir.n1.log 2>&1 &
$SOM $home2 start --pruning=nothing --grpc.address="0.0.0.0:9092" > $hdir.n2.log 2>&1 &

# Wait for chains to start
echo "Waiting for chains to start..."
sleep 8

# Delegate keys to the feeders
echo "Delegating feeder permissions..."
$SOM $home0 tx oracle delegate-feeder $($FED $home0 keys show feeder) $kbt --from val $cid -y &>/dev/null
$SOM $home1 tx oracle delegate-feeder $($FED $home1 keys show feeder) $kbt --from val $cid -y &>/dev/null
$SOM $home2 tx oracle delegate-feeder $($FED $home2 keys show feeder) $kbt --from val $cid -y &>/dev/null

# Start the oracle feeders
echo "Starting the oracle feeders..."
$FED $home0 start --log-level debug > $hdir.fed0.log 2>&1 &
$FED $home1 start --log-level debug > $hdir.fed1.log 2>&1 &
$FED $home2 start --log-level debug > $hdir.fed2.log 2>&1 &
echo
echo "Logs:"
echo "  - n0 'tail -f ./data/somm.n0.log'"
echo "  - n1 'tail -f ./data/somm.n1.log'"
echo "  - n2 'tail -f ./data/somm.n2.log'"
echo "  - f0 'tail -f ./data/somm.fed0.log'"
echo "  - f1 'tail -f ./data/somm.fed1.log'"
echo "  - f2 'tail -f ./data/somm.fed2.log'"
echo 
echo "Env for easy access:"
echo "export H1='--home ./data/somm/n0/'"
echo "export H2='--home ./data/somm/n1/'"
echo "export H3='--home ./data/somm/n2/'"
echo 
echo "Command Line Access:"
echo "  - n0 'sommelier --home ./data/somm/n0/ status'"
echo "  - n1 'sommelier --home ./data/somm/n1/ status'"
echo "  - n2 'sommelier --home ./data/somm/n2/ status'"
echo "  - e0 'oracle-feeder --home ./data/somm/n0/ q params'"
echo "  - e1 'oracle-feeder --home ./data/somm/n1/ q params'"
echo "  - e2 'oracle-feeder --home ./data/somm/n2/ q params'"
