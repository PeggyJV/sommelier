# Delegate Keys Overview

In addition to the typical validator key, Sommelier validators are required to generate *delegate keys* to fully participate in the Sommelier protocol. These delegate keys are used to authenticate each validator's Steward and Orchestrator. They are registered in the `gravity` module to create a relationship between the validator's address, the signer (Ethereum) address, and the orchestrator (Cosmos) address.

## Signer key

The Signer is the delegate Ethereum identity key and is used by Orchestrator to submit confirmations to the Sommelier chain for any transactions meant to be sent across the Gravity bridge to Ethereum. These confirmations contain signatures that will be used to prove to the Gravity smart contract that the validator set has agreed to bridge the transaction. A quorum (2/3) of validator signing power must be represented by the submitted confirmations before a transaction can be successfully relayed to Ethereum. 

Failing to confirm bridge transactions for a certain number of blocks is a slashable offense.

## Orchestrator key

The Orchestrator key represents the delegate Cosmos identity. It is shared by the `steward` and `orchestrator` processes to sign transactions submitted to the Sommelier chain (including confirmation transactions).

This key is critical for Steward to be able to submit transactions containing strategists' Cellar contract calls, and for Orchestrator to submit confirmations and oracle events, therefore a misconfigured (unregistered) orchestrator key will result in the validator being slashed and jailed.

## Registering delegate keys

Example prerequisites:
- [Steward]() and a [config.toml]() with the keystore path set
- A signer key (in this example named "signer")
- An orchestrator key (in this example named "orchestrator")
- A validator key (in this example named "validator")

See the [Validator Setup](/docs/validators/setup.md) doc for more details.

The gravity module defines a `set-delegate-key` command for the Sommelier CLI. It requires that you pass in the `signer` key signature over your validator operator address:

```bash
# Get addresses
ORCHESTRATOR_ADDR=$(steward --config config.toml keys cosmos show orchestrator)
SIGNER_ADDR=$(steward --config config.toml keys eth show signer)
VALOPER_ADDR=$(sommelier keys show validator --bech val -a)

# Sign the validator operator address with the signer key
SIGNATURE=$(steward --config config.toml sign-delegate-keys signer $VALOPER_ADDR) 

# Submit delegate keys tx
sommelier tx gravity set-delegate-keys \
    $VALOPERADDR \
    $ORCHESTRATOR_ADDR \
    $SIGNER_ADDR \
    $SIGNATURE \
    --chain-id sommelier-3 \
    --from validator \
    -y
```

