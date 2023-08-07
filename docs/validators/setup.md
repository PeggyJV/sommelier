# Sommelier Validator Documentation

## Introduction

This documentation covers the requirements and steps to participate in the Sommelier chain's validator set and provision strategists' management of Sommelier Cellars (vaults).

The Sommelier chain uses Proof of Stake consensus via Tendermint, built using the Cosmos SDK. It is the engine of the larger Sommelier protocol. Validating for Sommelier is somewhat complex compared with typical Cosmos chains as it requires running two sidecar processes, `steward` and `orchestrator` which facilitate Strategists' management of Cellars and an in-house version of the Gravity Bridge, respectively. Failure to have fully functioning sidecars can result in slashing and jailing. With this in mind, we do not recommend Sommelier validation for less experienced validators or those without sufficient monitoring and alerting capabilities.

### Joining Mainnet

Running a validator node on the Sommelier mainnet requires three processes:

1. The validator node (`sommelier`)
2. The Gravity Bridge Orchestrator (`orchestrator`)
3. [Steward](https://github.com/peggyjv/steward) (`steward`)

We also recommend running a local `geth` process. Though a public Ethereum API service can work, you may run into request limits because of the frequent requests made by the `orchestrator` process.

The Steward CLI now supports all of the same commands as `gorc` and is the recommended way to configure delegate keys for new validators and to run the Orchestrator. __The Steward CLI is used to run *both* the `steward` and `orchestrator` processes__. There are post-installation steps for the `steward` process outlined at the end of the installation steps below. These are required for your Steward to participate in the protocol. For more information on these setup steps for Steward, see [Validators Instructions for Setting Up Steward](https://github.com/PeggyJV/steward/blob/main/docs/02-StewardForValidators.md) in the Steward repository.

> NOTE: The Steward CLI and Steward itself are distinct concepts in this document. The Steward CLI is used to start both the `steward` and `orchestrator` processes, while "Steward" refers specifically to the `steward` process.

### Installation & Setup

#### 1. Download binaries and install service files

First, download and install all of the binaries needed to run the complete Sommelier validator system:

```bash 
# Create an installation directory
mkdir install && cd install

# Install Steward 
wget https://github.com/PeggyJV/steward/releases/download/v3.3.2/steward
chmod +x steward 
sudo mv steward /usr/bin/steward

# Install Geth 
wget https://gethstore.blob.core.windows.net/builds/geth-linux-amd64-1.12.0-e501b3b0.tar.gz
tar -xvf geth-linux-amd64-1.12.0-e501b3b0.tar.gz
sudo mv geth-linux-amd64-1.12.0-e501b3b0/geth /usr/bin/geth 
rm -rf geth-linux-amd64-1.12.0-e501b3b0*

# Install Sommelier 
wget https://github.com/PeggyJV/sommelier/releases/download/v6.0.0/sommelier_6.0.0_linux_amd64.tar.gz
tar -xf sommelier_6.0.0_linux_amd64.tar.gz
sudo mv sommelier /usr/bin
rm -rf sommelier_6.0.0_linux_amd64* LICENSE README.md 
```

Confirm they are installed properly:

```bash 
which sommelier 
which steward 
which geth 
```

If any of these return an error, make sure there were no errors when downloading or moving the files in the steps above.

As a convenience, we provide systemd service files for each of the processes required for a full validator system: 

```bash 
# Fetch systemd unit file examples 
wget https://raw.githubusercontent.com/PeggyJV/sommelier/main/contrib/systemd/geth.service \
    https://raw.githubusercontent.com/PeggyJV/sommelier/main/contrib/systemd/sommelier.service \
    https://raw.githubusercontent.com/PeggyJV/sommelier/main/contrib/systemd/orchestrator.service \
    https://raw.githubusercontent.com/PeggyJV/sommelier/main/contrib/systemd/steward.service
```

Modify the unit files to fit your environment using your prefered code editor. Specifically, make sure `User` is your username, and that `WorkingDirectory` and `ExecStart` contain valid paths for your installations:

```bash nano geth.service nano orchestrator.service nano steward.service nano
sommelier.service ```

Once your service files are configured, install them to systemd:

```bash # And install them to systemd sudo mv geth.service
/etc/systemd/system/geth.service \ && sudo mv orchestrator.service
/etc/systemd/system/ \ && sudo mv steward.service /etc/systemd/system/ \ &&
sudo mv sommelier.service /etc/systemd/system/ \ && sudo systemctl
daemon-reload ```

At this point, if you chose to use it, geth can be started. The other services need further setup. Start geth and ensure there are no errors in the log:

```bash 
sudo systemctl start geth && sudo journalctl -u geth -f
```

#### 2. Initial configuration

Now let's set up a configuration file for the `steward` and `orchestrator` processes. These applications have a lot of overlap in their configuration, so in this walkthrough we'll just create one `config.toml` file and share it between the two.

Create a directory to work in and download the provided example config file:

```bash 
mkdir -p $HOME/steward && cd $HOME/steward wget
https://raw.githubusercontent.com/PeggyJV/sommelier/main/contrib/mainnet/sommelier-3/config.toml
```

Modify the steward/orchestrator config for your environment as needed. In particular, change the `keystore` field to the path where you would like steward to store the keys created in the following steps.

See the [Steward configuration doc for more information](https://github.com/PeggyJV/steward/blob/main/docs/01-Configuration.md)

```bash
nano config.toml
```

Initialize the validator node files. Replace "<MONIKER>" with your desired validator moniker (name).

```bash
sommelier init <MONIKER> --chain-id sommelier-3
```

#### 3. Create keys

For an overview of the purpose of each key in the Sommelier system, see [Keys](./keys.md).

> NOTE: Before you attempt to generate keys, ensure your `config.toml` has the keystore and both key derivation paths set. If you downloaded the example config file above, the derivation paths are already set for you.
>
> ```toml
> keystore = "/some/path/keystore"

> [cosmos]
> key_derivation_path = "m/44'/118'/0'/0/0"
>
> [ethereum]
> key_derivation_path = "m/44'/60'/0'/0/0"
> ```

First we will create your *delegate keys*. Generate a cosmos key called "orchestrator" and an ethereum key called "signer". Note that this example assumes that your config file is in the current working directory:

```bash
steward -c config.toml keys cosmos add orchestrator
steward -c config.toml keys eth add signer
```

You should now see two key files in your configure keystore path named "orchestrator" and "signer" respectively. 

Before the addresses associated with your delegate keys can be used, they will need a small amount of funds. Send a small amount of SOMM to the "orchestrator" address and a small amount of ETH to the "signer" address.

Next you will need a validator key, which we will name "validator". This example generates a new key. 

> NOTE: If you wish to recover an existing key you may use the `--recover` flag, then enter the mnemonic phrase when prompted.

```bash 
sommelier keys add validator
```

You will also need to send a small amount of SOMM to the validator address.

At this point you should have: 
- A funded cosmos wallet with a key named "orchestrator" in the steward keystore 
- A funded ethereum wallet with a key named "signer" in the steward keystore
- A funded cosmos wallet with a key named "validator" in the sommelier keystore

#### 4. Setup the validator node

From [this text file](/contrib/mainnet/sommelierr-3/peers.txt), copy the listed peers to the `persistent_peers` field of the config file ~/.sommelier/config/config.toml: 

```bash
nano ~/.sommelier/config/config.toml
```

Download the genesis file:

```bash
wget https://raw.githubusercontent.com/PeggyJV/sommelier/main/contrib/mainnet/sommelier-3/genesis.json -O $HOME/.sommelier/config/genesis.json
```

Sync your node. This may take a few minutes.

```bash
sudo systemctl start sommelier && sudo journalctl -u sommelier -f
```

Once your node is synced, create your validator on chain. Replace "<MONIKER>" in the example with your desired validator name, and change any other values you wish to configure:

```bash
sommelier tx staking create-validator \ 
    --amount=1000000usomm \ 
    --pubkey=$(sommelier tendermint show-validator) \ 
    --moniker="<MONIKER>" \ 
    --chain-id="sommelier-3" \
    --commission-rate="0.10" \ 
    --commission-max-rate="0.20" \
    --commission-max-change-rate="0.01" \ 
    --min-self-delegation="1" \ 
    --gas 300000\ 
    --fees="0usomm" \ 
    --from=validator
```

Register your delegate keys. This is how Sommelier authenticates your steward and orchestrator.

> NOTE: It is critical that you register your delegate keys. If you do not, your node will eventually be slashed and jailed.

```bash
sommelier tx gravity set-delegate-keys \ 
    $(sommelier keys show validator --bech val -a) \ # validator address 
    $(steward --config $HOME/steward/config.toml keys cosmos show orchestrator) \ # orchestrator address (this must be run manually and address extracted) 
    $(steward --config $HOME/steward/config.toml keys eth show signer) \ # eth signer address 
    $(steward --config $HOME/steward/config.toml sign-delegate-keys signer 
    $(sommelier keys show validator --bech val -a)) \
    --chain-id sommelier-3 \ 
    --from validator \ 
    -y
```

Start the orchestrator (note that we are not yet starting steward):

```bash
sudo systemctl start orchestrator && sudo journalctl -u orchestrator -f
```

#### 5. Steward setup

At this point, you should have a running validator node and Orchestrator.

Please follow the detailed guide in [Validators Instructions for Setting Up Steward](https://github.com/PeggyJV/steward/blob/main/docs/02-StewardForValidators.md) and return here afterward.

At this point you should have a server CA and server certificate for Steward, and your `config.toml` should be configured with those values. Now we can start the Steward service that we created during the other installation steps.

```bash 
sudo systemctl start steward && sudo journalctl -u steward -f
```

Once your Steward is running, ensure that its server endpoint is reachable over the internet. Then, if you haven't already, follow the steps outlined in the [Steward Registry repository](https://github.com/PeggyJV/steward-registry) to register your steward instance.

Your installation is complete! If you have any problems, please reach out in the validator lobby channels in Discord or Telegram.

