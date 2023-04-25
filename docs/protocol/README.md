# What is the Sommelier protocol?

The Sommelier protocol is a Cosmos SDK-based proof-of-stake chain for coordinating updates to smart contracts implementing investment strategies (referred to as "cellars"). Stakers of the SOMM token use governance functionality to approve of cellar contracts that have been deployed. Once approved, strategists update the investment positions of cellars for depositors in response to market conditions by sending update recommendations to validators, who then vote to approve an update. When a consensus of staked power has agreed, these updates are applied to the cellars using Sommelier's fork of the Gravity Bridge, which is the mechanism by which the Sommelier protocol communicates with Ethereum. As compensation for providing these services, depending on the specific configuration of a given cellar, the strategist and the protocol itself receive fees as a percentage of total assets and/or performance. Fees received by the protocol are currently held in a module account and in the v6 upgrade will be auctioned for SOMM, with the proceeds being distributed pro rata to stakers of SOMM.

# Validators

Validators run Sommelier nodes and vote to achieve consensus for every block produced in the network. Becoming a validator is permissionless -- anyone can become a validator by running the appropriate software and registering themselves with the protocol. A validator's voting power towards achieving consensus is determined by the amount of SOMM that has been delegated to them by stakers of the SOMM token. Validators will receive a percentage of SOMM paid to stakers determined by the commission they have configured. Stakers generally prefer validators who are engaged with the community, participate in governance, and update their software in a timely fashion. Running properly configured services and being responsive as a validator is a key component in the success of the Sommelier protocol.

## What software must a validator run?

The Sommelier protocol requires three distinct software processes to run. Each serves a different purpose in enabling end-to-end functionality of the protocol. Reliability and uptime for these processes is critical.

### Sommelier node

The Sommelier node is the Cosmos SDK-based full node that is required for participation in the network and achieving protocol consensus. Like all other Cosmos proof-of-stake chains, running this software is the primary responsibility of validators in order for the network to continue producing new blocks and processing messages.

The node software and releases can be found here:

https://github.com/PeggyJV/sommelier

### Gravity Bridge orchestrator

The Sommelier protocol communicates with Ethereum using a fork of the Gravity Bridge. The componenets of this system include a smart contract on Ethereum called the Gravity contract, and a module in the Sommelier node called the "gravity" module which is responsible for coordinating state between Cosmos and Ethereum.

This coordination is achieved by each validator running an orchestrator process. When registering with the chain, a validator will set up two "delegate" keys: one being a Cosmos account that will be used for sending orchestration messages to the chain, and one being an Ethereum account for interacting with the Gravity contract. The orchestrator serves two main functions:

* Signer: Events originating on the Cosmos side (updates to cellars, updates to the validator set data on the Gravity contract, bridging tokens to Ethereum) are recorded in the chain state. The orchestrator process will observe these events and create a signature confirming them (using their Ethereum delegate key) and submit them to the chain (using their Cosmos delegate key). Once enough of these confirmations have been received based on a consensus of staked power by validator, these signatures will be sent along to a call to the Gravity contract in order for the event to execute on Ethereum.

* Oracle: Events originating on the Ethereum side (bridging tokens to Sommelier, events emitted from successful Gravity calls) are observed by the oracle. Validators confirming their observation of these events is what closes the loop on the end-to-end process of Cosmos-originated events or bridging calls from Ethereum, and allows the bridge to continue to the next nonce and process further events.

It is recommended to run the orchestrator by using one of the operating modes of the Steward process.

Steward can be found here:
https://github.com/PeggyJV/steward

With documentation found here:
https://github.com/PeggyJV/steward/blob/main/docs/00-TableOfContents.md

### Strategist server

Strategists do not need to run code on-chain -- they can use machine learning, big data, or any other number of systems, and then submit their update recommendations based on their proprietary technology and algorithms. In order to make this possible, each validator must run a strategist server by which a strategist can communicate their update recommendations for approval by the validator set. The reference implementation of this server is called Steward, which by default takes these recommendations at face value and executes them.

The strategist server serves a number of purposes:

* Authentication and authorization: strategists communicate with the strategist server using mutually-authenticated TLS. Currently, Steward is pre-populated with the client-certificate CA for Seven Seas, our launch partner during the proof-of-concept phase of the network, but will in the v6 upgrade retrieve these CAs from on-chain data. In addition, permissions to update specific cellars are assigned to specific strategists (again, currently only Seven Seas). This same v6 upgrade will also allow for mechanisms to customize these strategist update subscriptions for each cellar according to the validator's wishes.

* API specification: Cellar contracts may have many different functions, some of which are the province of the individual strategist, and some of which are more sensitive and may require the approval of SOMM holders via governance. The strategist server defines what functions can be called in this manner by specifying a GRPC API for strategists to call, and then translating those calls into executable bytes for the Ethereum contract. Currently, calls which are not specified in Steward may be executed through reaching consensus via a text governance proposal and coordinating manually, but the v6 upgrade will create a mechanism for specifying calls scoped to governance.

Steward is provided as a reference implementation of the strategist server. Validators wishing to use the default behavior can run Steward without modification. It is not required that a validator specifically run Steward, as they could fork and modift it or create their own implementation, but it is required that they will expose the same GRPC interface a strategist expects to call and that their process actively participates in the cellar update process. Future updates to the protocol will introduce slashing penalties for non-participation in this process, as reaching consensus on cellar updates is required for Sommelier to function.

Steward can be found here:
https://github.com/PeggyJV/steward

With documentation found here:
https://github.com/PeggyJV/steward/blob/main/docs/00-TableOfContents.md

## Slashing

Penalties for validators that fail to participate or misbehave in a Cosmos proof-of-stake network are referred to as slashing. These can range from token penalties as a percentage of delegated stake (all stakers with the validator will lose this percentage) to the validator being jailed (temporarily removed from the consensus set and rewards) or tombstoned (permanently removed).

During the testing, bootstrapping, and proof of concept phases of the Sommelier protocol, token penalties are not currently being issued when slashed. Failing to have your orchestrator process keep pace with signing events can result in being jailed. Future updates will introduce penalties for the standard Cosmos SDK slashing cases of downtime or double signing, and protocol-specific penalties for failing to participate in the cellar update process.

## Communication channels

Telegram channel: how to get an invite link here?
Discord server: https://discord.com/invite/ZcAYgSBxvY

# Governance

Protocol-wide decisions are made via governance, with the consent of delegated power. As with most standard Cosmos SDK chains, validators can vote with their delegated power, and stakers can vote their own preferences, overriding the vote of the validators to which they have delegated. If a majority of non-abstaining voting power chooses Yes, a governance proposal will pass, unless greater than 1/3 of the dissenting power has chosen No With Veto. Sommelier governance consists of the standard Cosmos SDK proposal types, with additional protocol-specific proposals.

## Parameters

### Quorum

In order for a governance proposal to not be outright rejected, a quorum of staked power must vote. Sommelier has a higher than normal quorum requirement of 50%. Validators and stakers are highly encouraged to participate in governance.

### Deposit

When a governance proposal is submitted, it first enters the deposit period. If a full deposit is sent with proposal submission, it will immediately enter the voting period -- otherwise, the deposit period will last up to a maximum of 48 hours, during which other token holders may contribute to the deposit to reach the required amount. To discourage proposal spam, the current deposit requirement is 5000 SOMM.

### Voting period

Sommelier's voting period is relatively short compared to other chains. Once a proposal has reached the required deposit, it will enter the voting period, which will last a guaranteed 48 hours. During this period, validators and stakers can vote, and are free to change their votes at any time until the voting period has concluded. If a quorum of delegated power is reached by the end of the voting period, the proposal will be considered based on the received votes, otherwise it will be rejected. Due to the short voting period and high quorum requirement, Sommelier relies on the active involvement and consent of validators and stakers.

## Proposals

### Standard Cosmos SDK

Sommelier supports the standard Cosmos SDK proposal types as of v0.45.

#### Text

Text proposals describe a desired action for the protocol to take by community agreement, but do not execute any code when voting completes. These typically include general policies or philosophies of how the chain should operate or signaling proposals to get community consent for an upcoming feature or decision.

#### Parameter change

Each module in a Cosmos SDK chain may choose to define a set of parameters -- values that affect the module which can be adjusted via governance without requiring an upgrade.

#### Community spend
#### Upgrade

### Sommelier

# Gravity Bridge

# Steward

# Auctions

# Strategy Providers
