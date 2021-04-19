<!-- This file is auto-generated. Please do not modify it yourself. -->
# Protobuf Documentation
<a name="top"></a>

## Table of Contents

- [allocation/v1/allocation.proto](#allocation/v1/allocation.proto)
    - [Allocation](#allocation.v1.Allocation)
    - [AllocationPrecommit](#allocation.v1.AllocationPrecommit)
    - [Cellar](#allocation.v1.Cellar)
    - [OracleFeed](#allocation.v1.OracleFeed)
    - [Pool](#allocation.v1.Pool)
    - [Tick](#allocation.v1.Tick)
    - [TickWeight](#allocation.v1.TickWeight)
    - [TickWeights](#allocation.v1.TickWeights)
    - [UniswapPair](#allocation.v1.UniswapPair)
    - [UniswapToken](#allocation.v1.UniswapToken)
  
- [allocation/v1/tx.proto](#allocation/v1/tx.proto)
    - [MsgAllocationCommit](#allocation.v1.MsgAllocationCommit)
    - [MsgAllocationCommitResponse](#allocation.v1.MsgAllocationCommitResponse)
    - [MsgAllocationPrecommit](#allocation.v1.MsgAllocationPrecommit)
    - [MsgAllocationPrecommitResponse](#allocation.v1.MsgAllocationPrecommitResponse)
    - [MsgDelegateAllocations](#allocation.v1.MsgDelegateAllocations)
    - [MsgDelegateAllocationsResponse](#allocation.v1.MsgDelegateAllocationsResponse)
  
    - [Msg](#allocation.v1.Msg)
  
- [allocation/v1/genesis.proto](#allocation/v1/genesis.proto)
    - [AggregatedOracleData](#allocation.v1.AggregatedOracleData)
    - [GenesisState](#allocation.v1.GenesisState)
    - [MissCounter](#allocation.v1.MissCounter)
    - [Params](#allocation.v1.Params)
  
- [allocation/v1/query.proto](#allocation/v1/query.proto)
    - [QueryAggregateDataRequest](#allocation.v1.QueryAggregateDataRequest)
    - [QueryAggregateDataResponse](#allocation.v1.QueryAggregateDataResponse)
    - [QueryAllocationCommitRequest](#allocation.v1.QueryAllocationCommitRequest)
    - [QueryAllocationCommitResponse](#allocation.v1.QueryAllocationCommitResponse)
    - [QueryAllocationPrecommitRequest](#allocation.v1.QueryAllocationPrecommitRequest)
    - [QueryAllocationPrecommitResponse](#allocation.v1.QueryAllocationPrecommitResponse)
    - [QueryCommitPeriodRequest](#allocation.v1.QueryCommitPeriodRequest)
    - [QueryCommitPeriodResponse](#allocation.v1.QueryCommitPeriodResponse)
    - [QueryDelegateAddressRequest](#allocation.v1.QueryDelegateAddressRequest)
    - [QueryDelegateAddressResponse](#allocation.v1.QueryDelegateAddressResponse)
    - [QueryLatestPeriodAggregateDataRequest](#allocation.v1.QueryLatestPeriodAggregateDataRequest)
    - [QueryLatestPeriodAggregateDataResponse](#allocation.v1.QueryLatestPeriodAggregateDataResponse)
    - [QueryMissCounterRequest](#allocation.v1.QueryMissCounterRequest)
    - [QueryMissCounterResponse](#allocation.v1.QueryMissCounterResponse)
    - [QueryParamsRequest](#allocation.v1.QueryParamsRequest)
    - [QueryParamsResponse](#allocation.v1.QueryParamsResponse)
    - [QueryValidatorAddressRequest](#allocation.v1.QueryValidatorAddressRequest)
    - [QueryValidatorAddressResponse](#allocation.v1.QueryValidatorAddressResponse)
  
    - [Query](#allocation.v1.Query)
  
- [il/v1/il.proto](#il/v1/il.proto)
    - [Params](#il.v1.Params)
    - [Stoploss](#il.v1.Stoploss)
  
- [il/v1/genesis.proto](#il/v1/genesis.proto)
    - [GenesisState](#il.v1.GenesisState)
    - [StoplossPositions](#il.v1.StoplossPositions)
    - [SubmittedPosition](#il.v1.SubmittedPosition)
  
- [il/v1/query.proto](#il/v1/query.proto)
    - [QueryParamsRequest](#il.v1.QueryParamsRequest)
    - [QueryParamsResponse](#il.v1.QueryParamsResponse)
    - [QueryStoplossPositionsRequest](#il.v1.QueryStoplossPositionsRequest)
    - [QueryStoplossPositionsResponse](#il.v1.QueryStoplossPositionsResponse)
    - [QueryStoplossRequest](#il.v1.QueryStoplossRequest)
    - [QueryStoplossResponse](#il.v1.QueryStoplossResponse)
  
    - [Query](#il.v1.Query)
  
- [il/v1/tx.proto](#il/v1/tx.proto)
    - [MsgCreateStoploss](#il.v1.MsgCreateStoploss)
    - [MsgCreateStoplossResponse](#il.v1.MsgCreateStoplossResponse)
    - [MsgDeleteStoploss](#il.v1.MsgDeleteStoploss)
    - [MsgDeleteStoplossResponse](#il.v1.MsgDeleteStoplossResponse)
  
    - [Msg](#il.v1.Msg)
  
- [Scalar Value Types](#scalar-value-types)



<a name="allocation/v1/allocation.proto"></a>
<p align="right"><a href="#top">Top</a></p>

## allocation/v1/allocation.proto



<a name="allocation.v1.Allocation"></a>

### Allocation
Allocation is the XXX


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `cellar_id` | [string](#string) |  |  |
| `fee_level` | [string](#string) |  |  |
| `tick_weights` | [TickWeights](#allocation.v1.TickWeights) |  |  |
| `salt` | [string](#string) |  |  |






<a name="allocation.v1.AllocationPrecommit"></a>

### AllocationPrecommit
AllocationPrecommit defines an array of hashed data to be used for the precommit phase
of allocation


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `hash` | [bytes](#bytes) |  |  |
| `cellar_id` | [string](#string) |  |  |






<a name="allocation.v1.Cellar"></a>

### Cellar
Cellar is a collection of pools for a token pair


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `cellar_id` | [string](#string) |  |  |
| `token0` | [string](#string) |  |  |
| `token1` | [string](#string) |  |  |
| `pool` | [Pool](#allocation.v1.Pool) | repeated |  |






<a name="allocation.v1.OracleFeed"></a>

### OracleFeed
OracleFeed represents an array of oracle data that is


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `data` | [UniswapPair](#allocation.v1.UniswapPair) | repeated |  |






<a name="allocation.v1.Pool"></a>

### Pool
Pool collects tick information


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `fee_level` | [string](#string) |  |  |
| `tick_ranges` | [Tick](#allocation.v1.Tick) | repeated |  |






<a name="allocation.v1.Tick"></a>

### Tick
Tick is the metadata for an instance of token data


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `min` | [uint64](#uint64) |  |  |
| `max` | [uint64](#uint64) |  |  |






<a name="allocation.v1.TickWeight"></a>

### TickWeight
TickWeight is the weighted Tick value


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `tick` | [Tick](#allocation.v1.Tick) |  |  |
| `weight` | [string](#string) |  |  |






<a name="allocation.v1.TickWeights"></a>

### TickWeights
TickWeights is a struct for holding a collection of TickWeight


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `weights` | [TickWeight](#allocation.v1.TickWeight) | repeated |  |






<a name="allocation.v1.UniswapPair"></a>

### UniswapPair
UniswapPair represents an SDK compatible uniswap pair info fetched from The Graph.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `id` | [string](#string) |  |  |
| `reserve0` | [string](#string) |  |  |
| `reserve1` | [string](#string) |  |  |
| `reserve_usd` | [string](#string) |  |  |
| `token0` | [UniswapToken](#allocation.v1.UniswapToken) |  |  |
| `token1` | [UniswapToken](#allocation.v1.UniswapToken) |  |  |
| `token0_price` | [string](#string) |  |  |
| `token1_price` | [string](#string) |  |  |
| `total_supply` | [string](#string) |  |  |






<a name="allocation.v1.UniswapToken"></a>

### UniswapToken
UniswapToken is the returned uniswap token representation


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `id` | [string](#string) |  | token address |
| `decimals` | [uint64](#uint64) |  | number of decimal positions of the pair token |





 <!-- end messages -->

 <!-- end enums -->

 <!-- end HasExtensions -->

 <!-- end services -->



<a name="allocation/v1/tx.proto"></a>
<p align="right"><a href="#top">Top</a></p>

## allocation/v1/tx.proto



<a name="allocation.v1.MsgAllocationCommit"></a>

### MsgAllocationCommit
MsgAllocationCommit - sdk.Msg for submitting arbitrary oracle data that has been prevoted on


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `commit` | [Allocation](#allocation.v1.Allocation) | repeated | vote containing the oracle data feed |
| `signer` | [string](#string) |  | signer (i.e feeder) account address |






<a name="allocation.v1.MsgAllocationCommitResponse"></a>

### MsgAllocationCommitResponse
MsgAllocationCommitResponse is the response type for the Msg/AllocationCommit gRPC method.






<a name="allocation.v1.MsgAllocationPrecommit"></a>

### MsgAllocationPrecommit
MsgAllocationPrecommit - sdk.Msg for prevoting on an array of oracle data types.
The purpose of the prevote is to hide vote for data with hashes formatted as hex string:
SHA256("{salt}:{data_cannonical_json}:{voter}")


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `precommit` | [AllocationPrecommit](#allocation.v1.AllocationPrecommit) | repeated | precommit containing the hash of the allocation precommit contents |
| `signer` | [string](#string) |  | signer (i.e feeder) account address |






<a name="allocation.v1.MsgAllocationPrecommitResponse"></a>

### MsgAllocationPrecommitResponse
MsgAllocationPrecommitResponse is the response type for the Msg/AllocationPrecommitResponse gRPC method.






<a name="allocation.v1.MsgDelegateAllocations"></a>

### MsgDelegateAllocations
MsgDelegateAllocations defines sdk.Msg for delegating allocation rights from a validator
to another address, must be signed by an active validator


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `delegate` | [string](#string) |  | delegate account address |
| `validator` | [string](#string) |  | validator operator address |






<a name="allocation.v1.MsgDelegateAllocationsResponse"></a>

### MsgDelegateAllocationsResponse
MsgDelegateAllocationsResponse is the response type for the Msg/DelegateAllocations gRPC method.





 <!-- end messages -->

 <!-- end enums -->

 <!-- end HasExtensions -->


<a name="allocation.v1.Msg"></a>

### Msg
MsgService defines the msgs that the oracle module handles.

| Method Name | Request Type | Response Type | Description | HTTP Verb | Endpoint |
| ----------- | ------------ | ------------- | ------------| ------- | -------- |
| `DelegateAllocations` | [MsgDelegateAllocations](#allocation.v1.MsgDelegateAllocations) | [MsgDelegateAllocationsResponse](#allocation.v1.MsgDelegateAllocationsResponse) | DelegateAllocations defines a message that delegates the allocating to an account address. | |
| `AllocationPrecommit` | [MsgAllocationPrecommit](#allocation.v1.MsgAllocationPrecommit) | [MsgAllocationPrecommitResponse](#allocation.v1.MsgAllocationPrecommitResponse) | OracleDataPrevote defines a message that commits a hash of a oracle data feed before the data is actually submitted. | |
| `AllocationCommit` | [MsgAllocationCommit](#allocation.v1.MsgAllocationCommit) | [MsgAllocationCommitResponse](#allocation.v1.MsgAllocationCommitResponse) | OracleDataVote defines a message to submit the actual oracle data that was committed by the feeder through the prevote. | |

 <!-- end services -->



<a name="allocation/v1/genesis.proto"></a>
<p align="right"><a href="#top">Top</a></p>

## allocation/v1/genesis.proto



<a name="allocation.v1.AggregatedOracleData"></a>

### AggregatedOracleData
AggregatedOracleData defines the aggregated oracle data at a given block height


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `height` | [int64](#int64) |  | block height in which the data was committed |
| `data` | [UniswapPair](#allocation.v1.UniswapPair) |  | oracle data |






<a name="allocation.v1.GenesisState"></a>

### GenesisState
GenesisState - all oracle state that must be provided at genesis


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `params` | [Params](#allocation.v1.Params) |  |  |
| `feeder_delegations` | [MsgDelegateAllocations](#allocation.v1.MsgDelegateAllocations) | repeated |  |
| `miss_counters` | [MissCounter](#allocation.v1.MissCounter) | repeated |  |
| `aggregates` | [AggregatedOracleData](#allocation.v1.AggregatedOracleData) | repeated |  |






<a name="allocation.v1.MissCounter"></a>

### MissCounter
MissCounter stores the validator address and the number of associated misses


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `validator` | [string](#string) |  | validator operator address |
| `misses` | [int64](#int64) |  | number of misses |






<a name="allocation.v1.Params"></a>

### Params
Params oracle parameters


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `vote_period` | [int64](#int64) |  | VotePeriod defines the number of blocks to wait for votes before attempting to tally |
| `vote_threshold` | [string](#string) |  | VoteThreshold defines the percentage of bonded stake required to vote each period |
| `slash_window` | [int64](#int64) |  | SlashWindow defines the number of blocks for the slashing window |
| `min_valid_per_window` | [string](#string) |  | MinValidPerWindow defines the number of misses a validator is allowed during each SlashWindow |
| `slash_fraction` | [string](#string) |  | SlashFraction defines the percentage of slash that a validator will suffer if it fails to send a vote |
| `target_threshold` | [string](#string) |  | TargetThreshold defines the max percentage difference that a given oracle data needs to have with the aggregated data in order for the feeder to be elegible for rewards. |
| `data_types` | [string](#string) | repeated | DataTypes defines which data types validators must submit each voting period |
| `cellars` | [Cellar](#allocation.v1.Cellar) | repeated | Cellars is the collection of pools that the allocation module references |





 <!-- end messages -->

 <!-- end enums -->

 <!-- end HasExtensions -->

 <!-- end services -->



<a name="allocation/v1/query.proto"></a>
<p align="right"><a href="#top">Top</a></p>

## allocation/v1/query.proto



<a name="allocation.v1.QueryAggregateDataRequest"></a>

### QueryAggregateDataRequest
QueryAggregateDataRequest is the request type for the Query/AggregateData gRPC method.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `type` | [string](#string) |  | oracle data type |
| `id` | [string](#string) |  | oracle data identifier |






<a name="allocation.v1.QueryAggregateDataResponse"></a>

### QueryAggregateDataResponse
QueryAggregateDataRequest is the response type for the Query/AggregateData gRPC method.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `oracle_data` | [UniswapPair](#allocation.v1.UniswapPair) |  | oracle data associated with the id and type from the request |
| `height` | [int64](#int64) |  | height at which the aggregated oracle data was stored |






<a name="allocation.v1.QueryAllocationCommitRequest"></a>

### QueryAllocationCommitRequest
QueryOracleDataVoteRequest is the request type for the Query/QueryOracleDataVote gRPC method.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `validator` | [string](#string) |  | validator operator address |






<a name="allocation.v1.QueryAllocationCommitResponse"></a>

### QueryAllocationCommitResponse
QueryOracleDataVoteResponse is the response type for the Query/QueryOracleDataVote gRPC method.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `commit` | [Allocation](#allocation.v1.Allocation) |  | vote containing the oracle feed submitted within the latest voting period |






<a name="allocation.v1.QueryAllocationPrecommitRequest"></a>

### QueryAllocationPrecommitRequest
QueryAllocationPrecommitRequest is the request type for the Query/AllocationPrecommit gRPC method.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `validator` | [string](#string) |  | validator operator address |






<a name="allocation.v1.QueryAllocationPrecommitResponse"></a>

### QueryAllocationPrecommitResponse
QueryOracleDataPrevoteResponse is the response type for the Query/QueryOracleDataPrevote gRPC method.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `precommit` | [AllocationPrecommit](#allocation.v1.AllocationPrecommit) |  | prevote submitted within the latest voting period |






<a name="allocation.v1.QueryCommitPeriodRequest"></a>

### QueryCommitPeriodRequest
QueryVotePeriodRequest is the request type for the Query/VotePeriod gRPC method.






<a name="allocation.v1.QueryCommitPeriodResponse"></a>

### QueryCommitPeriodResponse
QueryVotePeriodResponse is the response type for the Query/VotePeriod gRPC method.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `current_height` | [int64](#int64) |  | block height at which the query was processed |
| `vote_period_start` | [int64](#int64) |  | latest vote period start block height |
| `vote_period_end` | [int64](#int64) |  | block height at which the current voting period ends |






<a name="allocation.v1.QueryDelegateAddressRequest"></a>

### QueryDelegateAddressRequest
QueryDelegateAddressRequest is the request type for the Query/QueryDelegateAddress gRPC method.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `validator` | [string](#string) |  | validator operator address |






<a name="allocation.v1.QueryDelegateAddressResponse"></a>

### QueryDelegateAddressResponse
QueryDelegateAddressResponse is the response type for the Query/QueryDelegateAddress gRPC method.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `delegate` | [string](#string) |  | delegate account address |






<a name="allocation.v1.QueryLatestPeriodAggregateDataRequest"></a>

### QueryLatestPeriodAggregateDataRequest
QueryLatestPeriodAggregateDataRequest is the request type for the Query/QueryLatestPeriodAggregateData gRPC method.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `pagination` | [cosmos.base.query.v1beta1.PageRequest](#cosmos.base.query.v1beta1.PageRequest) |  | pagination defines an optional pagination for the request. |






<a name="allocation.v1.QueryLatestPeriodAggregateDataResponse"></a>

### QueryLatestPeriodAggregateDataResponse
QueryLatestPeriodAggregateDataResponse is the response type for the Query/QueryLatestPeriodAggregateData gRPC method.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `oracle_data` | [UniswapPair](#allocation.v1.UniswapPair) | repeated | oracle data associated with the |
| `height` | [int64](#int64) |  | height at which the aggregated oracle data was stored |
| `pagination` | [cosmos.base.query.v1beta1.PageResponse](#cosmos.base.query.v1beta1.PageResponse) |  | pagination defines the pagination in the response. |






<a name="allocation.v1.QueryMissCounterRequest"></a>

### QueryMissCounterRequest
QueryMissCounterRequest is the request type for the Query/MissCounter gRPC method.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `validator` | [string](#string) |  | validator operator address |






<a name="allocation.v1.QueryMissCounterResponse"></a>

### QueryMissCounterResponse
QueryMissCounterResponse is the response type for the Query/MissCounter gRPC method.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `miss_counter` | [int64](#int64) |  | number of oracle feed votes missed since the last counter reset |






<a name="allocation.v1.QueryParamsRequest"></a>

### QueryParamsRequest
QueryParamsRequest is the request type for the Query/Params gRPC method.






<a name="allocation.v1.QueryParamsResponse"></a>

### QueryParamsResponse
QueryParamsRequest is the response type for the Query/Params gRPC method.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `params` | [Params](#allocation.v1.Params) |  | oracle parameters |






<a name="allocation.v1.QueryValidatorAddressRequest"></a>

### QueryValidatorAddressRequest
QueryValidatorAddressRequest is the request type for the Query/Params gRPC method.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `delegate` | [string](#string) |  | delegate account address |






<a name="allocation.v1.QueryValidatorAddressResponse"></a>

### QueryValidatorAddressResponse
QueryValidatorAddressResponse is the response type for the Query/Params gRPC method.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `validator` | [string](#string) |  | validator operator address |





 <!-- end messages -->

 <!-- end enums -->

 <!-- end HasExtensions -->


<a name="allocation.v1.Query"></a>

### Query
Query defines the gRPC querier service for the oracle module.

| Method Name | Request Type | Response Type | Description | HTTP Verb | Endpoint |
| ----------- | ------------ | ------------- | ------------| ------- | -------- |
| `QueryParams` | [QueryParamsRequest](#allocation.v1.QueryParamsRequest) | [QueryParamsResponse](#allocation.v1.QueryParamsResponse) | QueryParams queries the oracle module parameters. | GET|/sommelier/allocation/v1/params|
| `QueryDelegateAddress` | [QueryDelegateAddressRequest](#allocation.v1.QueryDelegateAddressRequest) | [QueryDelegateAddressResponse](#allocation.v1.QueryDelegateAddressResponse) | QueryDelegateAddress queries the delegate account address of a validator | GET|/sommelier/allocation/v1/delegates/{validator}|
| `QueryValidatorAddress` | [QueryValidatorAddressRequest](#allocation.v1.QueryValidatorAddressRequest) | [QueryValidatorAddressResponse](#allocation.v1.QueryValidatorAddressResponse) | QueryValidatorAddress returns the validator address of a given delegate | GET|/sommelier/allocation/v1/validators/{delegate}|
| `QueryAllocationPrecommit` | [QueryAllocationPrecommitRequest](#allocation.v1.QueryAllocationPrecommitRequest) | [QueryAllocationPrecommitResponse](#allocation.v1.QueryAllocationPrecommitResponse) | QueryOracleDataPrevote queries the validator prevote in the current voting period | GET|/sommelier/allocation/v1/precommits/{validator}|
| `QueryAllocationCommit` | [QueryAllocationCommitRequest](#allocation.v1.QueryAllocationCommitRequest) | [QueryAllocationCommitResponse](#allocation.v1.QueryAllocationCommitResponse) | QueryOracleDataVote queries the validator vote in the current voting period | GET|/sommelier/allocation/v1/commits/{validator}|
| `QueryCommitPeriod` | [QueryCommitPeriodRequest](#allocation.v1.QueryCommitPeriodRequest) | [QueryCommitPeriodResponse](#allocation.v1.QueryCommitPeriodResponse) | QueryVotePeriod queries the heights for the current voting period (current, start and end) | GET|/sommelier/allocation/v1/commit_period|
| `QueryMissCounter` | [QueryMissCounterRequest](#allocation.v1.QueryMissCounterRequest) | [QueryMissCounterResponse](#allocation.v1.QueryMissCounterResponse) | QueryMissCounter queries the missed number of oracle data feed periods | GET|/sommelier/allocation/v1/miss_counters/{validator}|
| `QueryAggregateData` | [QueryAggregateDataRequest](#allocation.v1.QueryAggregateDataRequest) | [QueryAggregateDataResponse](#allocation.v1.QueryAggregateDataResponse) | QueryAggregateData returns the latest aggregated data value for a given type and identifioer | GET|/sommelier/allocation/v1/aggregate_data/{id}/{type}|
| `QueryLatestPeriodAggregateData` | [QueryLatestPeriodAggregateDataRequest](#allocation.v1.QueryLatestPeriodAggregateDataRequest) | [QueryLatestPeriodAggregateDataResponse](#allocation.v1.QueryLatestPeriodAggregateDataResponse) | QueryLatestPeriodAggregateData returns the aggregated data for a given pair an identifioer | GET|/sommelier/allocation/v1/aggregate_data|

 <!-- end services -->



<a name="il/v1/il.proto"></a>
<p align="right"><a href="#top">Top</a></p>

## il/v1/il.proto



<a name="il.v1.Params"></a>

### Params
Params define the impermanent loss module parameters


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `batch_contract_address` | [string](#string) |  | batch contract address for batching impermanent loss handling (i.e remove liquidity) on Ethereum. This contract calls the liquidity contract address for each position that has been batched. |
| `liquidity_contract_address` | [string](#string) |  | liquidity contract address for removing liquidity from each LP position |
| `eth_timeout_blocks` | [uint64](#uint64) |  | timeout block height value for the custom ethereum outgoing logic. This is value added to the last seen ethereum height when executing the stoploss logic on EndBlock. |
| `eth_timeout_timestamp` | [uint64](#uint64) |  | timeout timestamp second duration value for the redeemLiquidity deadline. This value is added to the block unix timestamp when executing the stoploss logic on EndBlock. |






<a name="il.v1.Stoploss"></a>

### Stoploss
Stoploss defines a set of parameters that together trigger a stoploss withdrawal.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `uniswap_pair_id` | [string](#string) |  | uniswap pair hex address |
| `liquidity_pool_shares` | [uint64](#uint64) |  | amount of shares from the liquidity pool to redeem if current slippage > max slipage |
| `max_slippage` | [string](#string) |  | max slippage allowed before the stoploss is triggered |
| `reference_pair_ratio` | [string](#string) |  | starting token pair ratio of the uniswap pool |
| `receiver_address` | [string](#string) |  | ethereum receiving address in hex format |
| `redeem_eth` | [bool](#bool) |  | redeem liquidity for eth or for the corresponding pair tokens once the stoploss position is executed |
| `submitted` | [bool](#bool) |  | track submission to the bridge in order to reenable the position if the tx timeouts |





 <!-- end messages -->

 <!-- end enums -->

 <!-- end HasExtensions -->

 <!-- end services -->



<a name="il/v1/genesis.proto"></a>
<p align="right"><a href="#top">Top</a></p>

## il/v1/genesis.proto



<a name="il.v1.GenesisState"></a>

### GenesisState
GenesisState all impermanent loss state that must be provided at genesis.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `params` | [Params](#il.v1.Params) |  |  |
| `lps_stoploss_positions` | [StoplossPositions](#il.v1.StoplossPositions) | repeated |  |
| `invalidation_id` | [uint64](#uint64) |  |  |
| `submitted_positions_queue` | [SubmittedPosition](#il.v1.SubmittedPosition) | repeated |  |






<a name="il.v1.StoplossPositions"></a>

### StoplossPositions
StoplossPosition represents all the impermanent loss stop positions for a given LP address and uniswap pair.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `address` | [string](#string) |  | LP account address |
| `stoploss_positions` | [Stoploss](#il.v1.Stoploss) | repeated | set of positions owned by the address |






<a name="il.v1.SubmittedPosition"></a>

### SubmittedPosition
SubmittedPosition contains an impermanent loss position owned by an account
that has already been submitted as outgoing bridge tx request to ethereum.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `address` | [string](#string) |  | account address |
| `timeout_height` | [uint64](#uint64) |  | ethereum block height at which the position tx timeouts |
| `pair_id` | [string](#string) |  | stoploss position's uniswap pair id |





 <!-- end messages -->

 <!-- end enums -->

 <!-- end HasExtensions -->

 <!-- end services -->



<a name="il/v1/query.proto"></a>
<p align="right"><a href="#top">Top</a></p>

## il/v1/query.proto



<a name="il.v1.QueryParamsRequest"></a>

### QueryParamsRequest
QueryParamsRequest is the request type for the Query/Params gRPC method.






<a name="il.v1.QueryParamsResponse"></a>

### QueryParamsResponse
QueryParamsResponse is the response type for the Query/Params gRPC method.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `params` | [Params](#il.v1.Params) |  | impermanent loss parameters |






<a name="il.v1.QueryStoplossPositionsRequest"></a>

### QueryStoplossPositionsRequest
QueryStoplossPisitionsRequest is the request type for the Query/StoplossPositions gRPC method.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `address` | [string](#string) |  | account address that owns the stoploss positions |
| `pagination` | [cosmos.base.query.v1beta1.PageRequest](#cosmos.base.query.v1beta1.PageRequest) |  | pagination request |






<a name="il.v1.QueryStoplossPositionsResponse"></a>

### QueryStoplossPositionsResponse
QueryStoplossPositionsResponse is the response type for the Query/StoplossPositions gRPC method.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `stoploss_positions` | [Stoploss](#il.v1.Stoploss) | repeated | set of possitions owned by the given address |
| `pagination` | [cosmos.base.query.v1beta1.PageResponse](#cosmos.base.query.v1beta1.PageResponse) |  | pagination response |






<a name="il.v1.QueryStoplossRequest"></a>

### QueryStoplossRequest
QueryStoplossRequest is the request type for the Query/Stoploss gRPC method.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `address` | [string](#string) |  | account address that owns the stoploss position |
| `uniswap_pair` | [string](#string) |  | uniswap pair of the position |






<a name="il.v1.QueryStoplossResponse"></a>

### QueryStoplossResponse
QueryStoplossResponse is the response type for the Query/Stoploss gRPC method.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `stoploss` | [Stoploss](#il.v1.Stoploss) |  | stoploss position for the given address and pair. |





 <!-- end messages -->

 <!-- end enums -->

 <!-- end HasExtensions -->


<a name="il.v1.Query"></a>

### Query
Query defines a gRPC query service for the impermanent loss module.

| Method Name | Request Type | Response Type | Description | HTTP Verb | Endpoint |
| ----------- | ------------ | ------------- | ------------| ------- | -------- |
| `Stoploss` | [QueryStoplossRequest](#il.v1.QueryStoplossRequest) | [QueryStoplossResponse](#il.v1.QueryStoplossResponse) | Stoploss queries a stoploss position for a given pair and account address. | GET|/il/v1/stoploss_positions/{address}/{uniswap_pair}|
| `StoplossPositions` | [QueryStoplossPositionsRequest](#il.v1.QueryStoplossPositionsRequest) | [QueryStoplossPositionsResponse](#il.v1.QueryStoplossPositionsResponse) | Stoploss returns all stoploss positions from an address. | GET|/il/v1/stoploss_positions/{address}|
| `Params` | [QueryParamsRequest](#il.v1.QueryParamsRequest) | [QueryParamsResponse](#il.v1.QueryParamsResponse) | Params queries the IL module parameters. | GET|/il/v1/params|

 <!-- end services -->



<a name="il/v1/tx.proto"></a>
<p align="right"><a href="#top">Top</a></p>

## il/v1/tx.proto



<a name="il.v1.MsgCreateStoploss"></a>

### MsgCreateStoploss
MsgStoploss defines a stoploss position


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `address` | [string](#string) |  | account address that owns the stoploss position |
| `stoploss` | [Stoploss](#il.v1.Stoploss) |  | stoploss position details |






<a name="il.v1.MsgCreateStoplossResponse"></a>

### MsgCreateStoplossResponse
MsgCreateStoplossResponse is the response type for the Msg/CreateStoploss gRPC method.






<a name="il.v1.MsgDeleteStoploss"></a>

### MsgDeleteStoploss
MsgDeleteStoploss removes a stoploss position


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `address` | [string](#string) |  | account address that owns the stoploss position |
| `uniswap_pair_id` | [string](#string) |  | uniswap pair hex address |






<a name="il.v1.MsgDeleteStoplossResponse"></a>

### MsgDeleteStoplossResponse
MsgDeleteStoplossResponse is the response type for the Msg/DeleteStoploss gRPC method.





 <!-- end messages -->

 <!-- end enums -->

 <!-- end HasExtensions -->


<a name="il.v1.Msg"></a>

### Msg
MsgService defines the msgs that the il module handles.

| Method Name | Request Type | Response Type | Description | HTTP Verb | Endpoint |
| ----------- | ------------ | ------------- | ------------| ------- | -------- |
| `CreateStoploss` | [MsgCreateStoploss](#il.v1.MsgCreateStoploss) | [MsgCreateStoplossResponse](#il.v1.MsgCreateStoplossResponse) | CreateStoploss sets a new tracking stoploss position for a uniswap pair | |
| `DeleteStoploss` | [MsgDeleteStoploss](#il.v1.MsgDeleteStoploss) | [MsgDeleteStoplossResponse](#il.v1.MsgDeleteStoplossResponse) | DeleteStoploss deletes an existing stoploss position | |

 <!-- end services -->



## Scalar Value Types

| .proto Type | Notes | C++ | Java | Python | Go | C# | PHP | Ruby |
| ----------- | ----- | --- | ---- | ------ | -- | -- | --- | ---- |
| <a name="double" /> double |  | double | double | float | float64 | double | float | Float |
| <a name="float" /> float |  | float | float | float | float32 | float | float | Float |
| <a name="int32" /> int32 | Uses variable-length encoding. Inefficient for encoding negative numbers – if your field is likely to have negative values, use sint32 instead. | int32 | int | int | int32 | int | integer | Bignum or Fixnum (as required) |
| <a name="int64" /> int64 | Uses variable-length encoding. Inefficient for encoding negative numbers – if your field is likely to have negative values, use sint64 instead. | int64 | long | int/long | int64 | long | integer/string | Bignum |
| <a name="uint32" /> uint32 | Uses variable-length encoding. | uint32 | int | int/long | uint32 | uint | integer | Bignum or Fixnum (as required) |
| <a name="uint64" /> uint64 | Uses variable-length encoding. | uint64 | long | int/long | uint64 | ulong | integer/string | Bignum or Fixnum (as required) |
| <a name="sint32" /> sint32 | Uses variable-length encoding. Signed int value. These more efficiently encode negative numbers than regular int32s. | int32 | int | int | int32 | int | integer | Bignum or Fixnum (as required) |
| <a name="sint64" /> sint64 | Uses variable-length encoding. Signed int value. These more efficiently encode negative numbers than regular int64s. | int64 | long | int/long | int64 | long | integer/string | Bignum |
| <a name="fixed32" /> fixed32 | Always four bytes. More efficient than uint32 if values are often greater than 2^28. | uint32 | int | int | uint32 | uint | integer | Bignum or Fixnum (as required) |
| <a name="fixed64" /> fixed64 | Always eight bytes. More efficient than uint64 if values are often greater than 2^56. | uint64 | long | int/long | uint64 | ulong | integer/string | Bignum |
| <a name="sfixed32" /> sfixed32 | Always four bytes. | int32 | int | int | int32 | int | integer | Bignum or Fixnum (as required) |
| <a name="sfixed64" /> sfixed64 | Always eight bytes. | int64 | long | int/long | int64 | long | integer/string | Bignum |
| <a name="bool" /> bool |  | bool | boolean | boolean | bool | bool | boolean | TrueClass/FalseClass |
| <a name="string" /> string | A string must always contain UTF-8 encoded or 7-bit ASCII text. | string | String | str/unicode | string | string | string | String (UTF-8) |
| <a name="bytes" /> bytes | May contain any arbitrary sequence of bytes. | string | ByteString | str | []byte | ByteString | string | String (ASCII-8BIT) |

