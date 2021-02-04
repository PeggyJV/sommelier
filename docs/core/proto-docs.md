<!-- This file is auto-generated. Please do not modify it yourself. -->
# Protobuf Documentation
<a name="top"></a>

## Table of Contents

- [il/v1/il.proto](#il/v1/il.proto)
    - [Params](#il.v1.Params)
    - [Stoploss](#il.v1.Stoploss)
  
- [il/v1/genesis.proto](#il/v1/genesis.proto)
    - [GenesisState](#il.v1.GenesisState)
    - [StoplossPositions](#il.v1.StoplossPositions)
  
- [il/v1/query.proto](#il/v1/query.proto)
    - [QueryParametersRequest](#il.v1.QueryParametersRequest)
    - [QueryParametersResponse](#il.v1.QueryParametersResponse)
    - [QueryStoplossPositionsRequest](#il.v1.QueryStoplossPositionsRequest)
    - [QueryStoplossPositionsResponse](#il.v1.QueryStoplossPositionsResponse)
    - [QueryStoplossRequest](#il.v1.QueryStoplossRequest)
    - [QueryStoplossResponse](#il.v1.QueryStoplossResponse)
  
    - [Query](#il.v1.Query)
  
- [il/v1/tx.proto](#il/v1/tx.proto)
    - [MsgStoploss](#il.v1.MsgStoploss)
    - [MsgStoplossResponse](#il.v1.MsgStoplossResponse)
  
    - [Msg](#il.v1.Msg)
  
- [oracle/v1/oracle.proto](#oracle/v1/oracle.proto)
    - [UniswapData](#oracle.v1.UniswapData)
    - [UniswapPair](#oracle.v1.UniswapPair)
    - [UniswapToken](#oracle.v1.UniswapToken)
  
- [oracle/v1/tx.proto](#oracle/v1/tx.proto)
    - [MsgDelegateFeedConsent](#oracle.v1.MsgDelegateFeedConsent)
    - [MsgDelegateFeedConsentResponse](#oracle.v1.MsgDelegateFeedConsentResponse)
    - [MsgOracleDataPrevote](#oracle.v1.MsgOracleDataPrevote)
    - [MsgOracleDataPrevoteResponse](#oracle.v1.MsgOracleDataPrevoteResponse)
    - [MsgOracleDataVote](#oracle.v1.MsgOracleDataVote)
    - [MsgOracleDataVoteResponse](#oracle.v1.MsgOracleDataVoteResponse)
  
    - [Msg](#oracle.v1.Msg)
  
- [oracle/v1/genesis.proto](#oracle/v1/genesis.proto)
    - [GenesisState](#oracle.v1.GenesisState)
    - [MissCounter](#oracle.v1.MissCounter)
    - [Params](#oracle.v1.Params)
  
- [oracle/v1/query.proto](#oracle/v1/query.proto)
    - [QueryDelegeateAddressRequest](#oracle.v1.QueryDelegeateAddressRequest)
    - [QueryDelegeateAddressResponse](#oracle.v1.QueryDelegeateAddressResponse)
    - [QueryMissCounterRequest](#oracle.v1.QueryMissCounterRequest)
    - [QueryMissCounterResponse](#oracle.v1.QueryMissCounterResponse)
    - [QueryOracleDataPrevoteRequest](#oracle.v1.QueryOracleDataPrevoteRequest)
    - [QueryOracleDataPrevoteResponse](#oracle.v1.QueryOracleDataPrevoteResponse)
    - [QueryOracleDataRequest](#oracle.v1.QueryOracleDataRequest)
    - [QueryOracleDataResponse](#oracle.v1.QueryOracleDataResponse)
    - [QueryOracleDataVoteRequest](#oracle.v1.QueryOracleDataVoteRequest)
    - [QueryParamsRequest](#oracle.v1.QueryParamsRequest)
    - [QueryParamsResponse](#oracle.v1.QueryParamsResponse)
    - [QueryValidatorAddressRequest](#oracle.v1.QueryValidatorAddressRequest)
    - [QueryValidatorAddressResponse](#oracle.v1.QueryValidatorAddressResponse)
    - [QueryVotePeriodRequest](#oracle.v1.QueryVotePeriodRequest)
    - [VotePeriod](#oracle.v1.VotePeriod)
  
    - [Query](#oracle.v1.Query)
  
- [Scalar Value Types](#scalar-value-types)



<a name="il/v1/il.proto"></a>
<p align="right"><a href="#top">Top</a></p>

## il/v1/il.proto



<a name="il.v1.Params"></a>

### Params
Params define the impermanent loss module parameters


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `contract_address` | [string](#string) |  | contract address for impermanent loss handling on ethereum |






<a name="il.v1.Stoploss"></a>

### Stoploss
Stoploss defines a set of parameters that together trigger a stoploss withdrawal.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `uniswap_pair_id` | [string](#string) |  | uniswap pair identifier |
| `liquidity_pool_shares` | [int64](#int64) |  | amount of shares from the liquidity pool to redeem if current slippage > max slipage |
| `max_slippage` | [string](#string) |  | max slippage allowed before the stoploss is triggered |
| `reference_pair_ratio` | [string](#string) |  | starting token pair ration of the uniswap pool |





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






<a name="il.v1.StoplossPositions"></a>

### StoplossPositions
StoplossPosition represents all the impermanent loss stop positions for a given LP address and uniswap pair.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `address` | [string](#string) |  | LP account address |
| `stoploss_positions` | [Stoploss](#il.v1.Stoploss) | repeated | set of possitions owned by the address |





 <!-- end messages -->

 <!-- end enums -->

 <!-- end HasExtensions -->

 <!-- end services -->



<a name="il/v1/query.proto"></a>
<p align="right"><a href="#top">Top</a></p>

## il/v1/query.proto



<a name="il.v1.QueryParametersRequest"></a>

### QueryParametersRequest
QueryParametersRequest is an empty request to query for the impermanent loss params






<a name="il.v1.QueryParametersResponse"></a>

### QueryParametersResponse
QueryParametersResponse


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `params` | [Params](#il.v1.Params) |  | impermanent loss parameters |






<a name="il.v1.QueryStoplossPositionsRequest"></a>

### QueryStoplossPositionsRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `address` | [string](#string) |  | account address that owns the stoploss positions |
| `pagination` | [cosmos.base.query.v1beta1.PageRequest](#cosmos.base.query.v1beta1.PageRequest) |  | pagination request |






<a name="il.v1.QueryStoplossPositionsResponse"></a>

### QueryStoplossPositionsResponse



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `stoploss_positions` | [Stoploss](#il.v1.Stoploss) | repeated | set of possitions owned by the given address |
| `pagination` | [cosmos.base.query.v1beta1.PageResponse](#cosmos.base.query.v1beta1.PageResponse) |  | pagination response |






<a name="il.v1.QueryStoplossRequest"></a>

### QueryStoplossRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `address` | [string](#string) |  | account address that owns the stoploss position |
| `uniswap_pair` | [string](#string) |  | uniswap pair of the position |






<a name="il.v1.QueryStoplossResponse"></a>

### QueryStoplossResponse



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
| `Stoploss` | [QueryStoplossRequest](#il.v1.QueryStoplossRequest) | [QueryStoplossResponse](#il.v1.QueryStoplossResponse) |  | GET|/il/v1/stoploss_positions/{address}/{uniswap_pair}|
| `StoplossPositions` | [QueryStoplossPositionsRequest](#il.v1.QueryStoplossPositionsRequest) | [QueryStoplossPositionsResponse](#il.v1.QueryStoplossPositionsResponse) |  | GET|/il/v1/stoploss_positions/{address}|
| `Parameters` | [QueryParametersRequest](#il.v1.QueryParametersRequest) | [QueryParametersResponse](#il.v1.QueryParametersResponse) |  | GET|/il/v1/parameters|

 <!-- end services -->



<a name="il/v1/tx.proto"></a>
<p align="right"><a href="#top">Top</a></p>

## il/v1/tx.proto



<a name="il.v1.MsgStoploss"></a>

### MsgStoploss
MsgStoploss defines a stoploss position


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `address` | [string](#string) |  |  |
| `stoploss` | [Stoploss](#il.v1.Stoploss) |  |  |






<a name="il.v1.MsgStoplossResponse"></a>

### MsgStoplossResponse






 <!-- end messages -->

 <!-- end enums -->

 <!-- end HasExtensions -->


<a name="il.v1.Msg"></a>

### Msg
MsgService defines the msgs that the oracle module handles.

| Method Name | Request Type | Response Type | Description | HTTP Verb | Endpoint |
| ----------- | ------------ | ------------- | ------------| ------- | -------- |
| `CreateStoploss` | [MsgStoploss](#il.v1.MsgStoploss) | [MsgStoplossResponse](#il.v1.MsgStoplossResponse) |  | |

 <!-- end services -->



<a name="oracle/v1/oracle.proto"></a>
<p align="right"><a href="#top">Top</a></p>

## oracle/v1/oracle.proto



<a name="oracle.v1.UniswapData"></a>
<<<<<<< HEAD
<<<<<<< HEAD

### UniswapData
UniswapData is an implementation of OracleData


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `pairs` | [UniswapPair](#oracle.v1.UniswapPair) | repeated |  |






<a name="oracle.v1.UniswapPair"></a>

### UniswapPair
UniswapPair represents the necessary data for a uniswap pair


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `id` | [string](#string) |  |  |
| `reserve0` | [string](#string) |  |  |
| `reserve1` | [string](#string) |  |  |
| `reserve_usd` | [string](#string) |  |  |
| `token0` | [UniswapToken](#oracle.v1.UniswapToken) |  |  |
| `token1` | [UniswapToken](#oracle.v1.UniswapToken) |  |  |
| `token0_price` | [string](#string) |  |  |
| `token1_price` | [string](#string) |  |  |
| `total_supply` | [string](#string) |  |  |






<a name="oracle.v1.UniswapToken"></a>

### UniswapToken
UniswapToken is the returned uniswap token representation
=======

### UniswapData
UniswapData is an implementation of OracleData
>>>>>>> 80d76df335316d61841ce9bca513aaf530ec40d3
=======

### UniswapData
UniswapData is an implementation of OracleData
>>>>>>> 80d76df335316d61841ce9bca513aaf530ec40d3


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
<<<<<<< HEAD
<<<<<<< HEAD
=======
=======
>>>>>>> 80d76df335316d61841ce9bca513aaf530ec40d3
| `pairs` | [UniswapPair](#oracle.v1.UniswapPair) | repeated |  |






<a name="oracle.v1.UniswapPair"></a>

### UniswapPair
UniswapPair represents the necessary data for a uniswap pair


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `id` | [string](#string) |  |  |
| `reserve0` | [string](#string) |  |  |
| `reserve1` | [string](#string) |  |  |
| `reserve_usd` | [string](#string) |  |  |
| `token0` | [UniswapToken](#oracle.v1.UniswapToken) |  |  |
| `token1` | [UniswapToken](#oracle.v1.UniswapToken) |  |  |
| `token0_price` | [string](#string) |  |  |
| `token1_price` | [string](#string) |  |  |
| `total_supply` | [string](#string) |  |  |






<a name="oracle.v1.UniswapToken"></a>

### UniswapToken
UniswapToken is the returned uniswap token representation


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
<<<<<<< HEAD
>>>>>>> 80d76df335316d61841ce9bca513aaf530ec40d3
| `decimals` | [string](#string) |  |  |
| `id` | [string](#string) |  |  |





 <!-- end messages -->

 <!-- end enums -->

 <!-- end HasExtensions -->

 <!-- end services -->



<a name="oracle/v1/tx.proto"></a>
<p align="right"><a href="#top">Top</a></p>

## oracle/v1/tx.proto



<a name="oracle.v1.MsgDelegateFeedConsent"></a>

### MsgDelegateFeedConsent
MsgDelegateFeedConsent - sdk.Msg for delegating oracle voting rights from a validator
to another address, must be signed by an active validator


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `delegate` | [string](#string) |  |  |
| `validator` | [string](#string) |  |  |
=======
| `decimals` | [string](#string) |  |  |
| `id` | [string](#string) |  |  |
>>>>>>> 80d76df335316d61841ce9bca513aaf530ec40d3






<a name="oracle.v1.MsgDelegateFeedConsentResponse"></a>

### MsgDelegateFeedConsentResponse




<<<<<<< HEAD



<a name="oracle.v1.MsgOracleDataPrevote"></a>

### MsgOracleDataPrevote
MsgOracleDataPrevote - sdk.Msg for prevoting on an array of oracle data types.
The purpose of the prevote is to hide vote for data with hashes formatted as hex string: 
SHA256("{salt}:{data_cannonical_json}:{voter}")


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `hashes` | [bytes](#bytes) | repeated |  |
| `signer` | [string](#string) |  |  |





<<<<<<< HEAD

<a name="oracle.v1.MsgOracleDataPrevoteResponse"></a>

### MsgOracleDataPrevoteResponse





## oracle/v1/query.proto
=======
>>>>>>> 80d76df335316d61841ce9bca513aaf530ec40d3

<a name="oracle.v1.MsgOracleDataPrevoteResponse"></a>

### MsgOracleDataPrevoteResponse



=======
<a name="oracle/v1/tx.proto"></a>
<p align="right"><a href="#top">Top</a></p>

## oracle/v1/tx.proto



<a name="oracle.v1.MsgDelegateFeedConsent"></a>

### MsgDelegateFeedConsent
MsgDelegateFeedConsent - sdk.Msg for delegating oracle voting rights from a validator
to another address, must be signed by an active validator
>>>>>>> 80d76df335316d61841ce9bca513aaf530ec40d3


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
<<<<<<< HEAD
| `pagination` | [cosmos.base.query.v1beta1.PageRequest](#cosmos.base.query.v1beta1.PageRequest) |  | pagination request |


=======
| `delegate` | [string](#string) |  |  |
| `validator` | [string](#string) |  |  |
>>>>>>> 80d76df335316d61841ce9bca513aaf530ec40d3


<a name="oracle.v1.MsgOracleDataVote"></a>

<<<<<<< HEAD
<<<<<<< HEAD

<a name="oracle.v1.QueryActivesResponse"></a>

<<<<<<< HEAD
### QueryActivesResponse

=======

<a name="oracle.v1.MsgDelegateFeedConsentResponse"></a>
>>>>>>> 80d76df335316d61841ce9bca513aaf530ec40d3

### MsgDelegateFeedConsentResponse

<<<<<<< HEAD
| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `denoms` | [string](#string) | repeated |  |
| `pagination` | [cosmos.base.query.v1beta1.PageResponse](#cosmos.base.query.v1beta1.PageResponse) |  | pagination response |
=======
>>>>>>> 80d76df335316d61841ce9bca513aaf530ec40d3






<<<<<<< HEAD
<a name="oracle.v1.QueryAggregatePrevoteRequest"></a>

### QueryAggregatePrevoteRequest

=======
### MsgOracleDataVote
MsgOracleDataVote - sdk.Msg for submitting arbitrary oracle data that has been prevoted on
>>>>>>> 80d76df335316d61841ce9bca513aaf530ec40d3
=======
<a name="oracle.v1.MsgOracleDataPrevote"></a>

### MsgOracleDataPrevote
MsgOracleDataPrevote - sdk.Msg for prevoting on an array of oracle data types.
The purpose of the prevote is to hide vote for data with hashes formatted as hex string: 
SHA256("{salt}:{data_cannonical_json}:{voter}")
>>>>>>> 80d76df335316d61841ce9bca513aaf530ec40d3
=======
### MsgOracleDataVote
MsgOracleDataVote - sdk.Msg for submitting arbitrary oracle data that has been prevoted on
>>>>>>> 5b665fcbd01d71b7bbcfcceba5a2d70aa2299577


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
<<<<<<< HEAD
| `salt` | [string](#string) | repeated |  |
| `oracle_data` | [google.protobuf.Any](#google.protobuf.Any) | repeated |  |
| `signer` | [string](#string) |  |  |






<a name="oracle.v1.MsgOracleDataVoteResponse"></a>

### MsgOracleDataVoteResponse

=======
| `hashes` | [bytes](#bytes) | repeated |  |
| `signer` | [string](#string) |  |  |
>>>>>>> 80d76df335316d61841ce9bca513aaf530ec40d3





<<<<<<< HEAD
 <!-- end messages -->

 <!-- end enums -->

 <!-- end HasExtensions -->


<a name="oracle.v1.Msg"></a>

### Msg
MsgService defines the msgs that the oracle module handles.

| Method Name | Request Type | Response Type | Description | HTTP Verb | Endpoint |
| ----------- | ------------ | ------------- | ------------| ------- | -------- |
| `DelegateFeedConsent` | [MsgDelegateFeedConsent](#oracle.v1.MsgDelegateFeedConsent) | [MsgDelegateFeedConsentResponse](#oracle.v1.MsgDelegateFeedConsentResponse) |  | |
| `OracleDataPrevote` | [MsgOracleDataPrevote](#oracle.v1.MsgOracleDataPrevote) | [MsgOracleDataPrevoteResponse](#oracle.v1.MsgOracleDataPrevoteResponse) |  | |
| `OracleDataVote` | [MsgOracleDataVote](#oracle.v1.MsgOracleDataVote) | [MsgOracleDataVoteResponse](#oracle.v1.MsgOracleDataVoteResponse) |  | |

 <!-- end services -->



<a name="oracle/v1/genesis.proto"></a>
<p align="right"><a href="#top">Top</a></p>

## oracle/v1/genesis.proto
=======
>>>>>>> 80d76df335316d61841ce9bca513aaf530ec40d3

<a name="oracle.v1.MsgOracleDataPrevoteResponse"></a>

### MsgOracleDataPrevoteResponse

<<<<<<< HEAD
<a name="oracle.v1.GenesisState"></a>

### GenesisState
GenesisState - all oracle state that must be provided at genesis
=======

>>>>>>> 80d76df335316d61841ce9bca513aaf530ec40d3


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `params` | [Params](#oracle.v1.Params) |  |  |
| `feeder_delegations` | [MsgDelegateFeedConsent](#oracle.v1.MsgDelegateFeedConsent) | repeated |  |
| `miss_counters` | [MissCounter](#oracle.v1.MissCounter) | repeated |  |



<a name="oracle.v1.MsgOracleDataVote"></a>

<<<<<<< HEAD


<a name="oracle.v1.MissCounter"></a>

### MissCounter
MissCounter stores the validator address and the number of associated misses


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
<<<<<<< HEAD
<<<<<<< HEAD
| `denom` | [string](#string) |  | token denomination |

=======
| `validator` | [string](#string) |  |  |
| `misses` | [int64](#int64) |  |  |
>>>>>>> 80d76df335316d61841ce9bca513aaf530ec40d3
=======
| `validator` | [string](#string) |  |  |
| `misses` | [int64](#int64) |  |  |
>>>>>>> 5b665fcbd01d71b7bbcfcceba5a2d70aa2299577






<a name="oracle.v1.Params"></a>

### Params
Params oracle parameters
=======
### MsgOracleDataVote
MsgOracleDataVote - sdk.Msg for submitting arbitrary oracle data that has been prevoted on
>>>>>>> 80d76df335316d61841ce9bca513aaf530ec40d3


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
<<<<<<< HEAD
<<<<<<< HEAD
<<<<<<< HEAD
| `rate` | [string](#string) |  | token exchange rate |
=======
| `salt` | [string](#string) | repeated |  |
| `oracle_data` | [google.protobuf.Any](#google.protobuf.Any) | repeated |  |
| `signer` | [string](#string) |  |  |
>>>>>>> 80d76df335316d61841ce9bca513aaf530ec40d3






<<<<<<< HEAD
<a name="oracle.v1.QueryExchangeRatesRequest"></a>

### QueryExchangeRatesRequest
=======
<a name="oracle.v1.MsgOracleDataVoteResponse"></a>

### MsgOracleDataVoteResponse
>>>>>>> 80d76df335316d61841ce9bca513aaf530ec40d3


=======
=======
>>>>>>> 5b665fcbd01d71b7bbcfcceba5a2d70aa2299577
| `vote_period` | [int64](#int64) |  | VotePeriod defines the number of blocks to wait for votes before attempting to tally |
| `vote_threshold` | [string](#string) |  | VoteThreshold defines the percentage of bonded stake required to vote each period |
| `slash_window` | [int64](#int64) |  | SlashWindow defines the number of blocks for the slashing window |
| `min_valid_per_window` | [string](#string) |  | MinValidPerWindow defines the number of misses a validator is allowed during each SlashWindow |
| `slash_fraction` | [string](#string) |  | SlashFraction defines the percentage of slash that a validator will suffer if it fails to send a vote |
| `data_types` | [string](#string) | repeated | DataTypes defines which data types validators must submit each voting period |
<<<<<<< HEAD
>>>>>>> 80d76df335316d61841ce9bca513aaf530ec40d3
=======
>>>>>>> 5b665fcbd01d71b7bbcfcceba5a2d70aa2299577

<<<<<<< HEAD
| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `pagination` | [cosmos.base.query.v1beta1.PageRequest](#cosmos.base.query.v1beta1.PageRequest) |  | pagination request |
=======
>>>>>>> 80d76df335316d61841ce9bca513aaf530ec40d3



 <!-- end messages -->

 <!-- end enums -->

 <!-- end HasExtensions -->

<<<<<<< HEAD
 <!-- end messages -->

 <!-- end enums -->

 <!-- end HasExtensions -->

 <!-- end services -->

<<<<<<< HEAD
<<<<<<< HEAD
| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `rates` | [cosmos.base.v1beta1.DecCoin](#cosmos.base.v1beta1.DecCoin) | repeated |  |
| `pagination` | [cosmos.base.query.v1beta1.PageResponse](#cosmos.base.query.v1beta1.PageResponse) |  | pagination response |
=======

<a name="oracle.v1.Msg"></a>

### Msg
MsgService defines the msgs that the oracle module handles.

| Method Name | Request Type | Response Type | Description | HTTP Verb | Endpoint |
| ----------- | ------------ | ------------- | ------------| ------- | -------- |
| `DelegateFeedConsent` | [MsgDelegateFeedConsent](#oracle.v1.MsgDelegateFeedConsent) | [MsgDelegateFeedConsentResponse](#oracle.v1.MsgDelegateFeedConsentResponse) |  | |
| `OracleDataPrevote` | [MsgOracleDataPrevote](#oracle.v1.MsgOracleDataPrevote) | [MsgOracleDataPrevoteResponse](#oracle.v1.MsgOracleDataPrevoteResponse) |  | |
| `OracleDataVote` | [MsgOracleDataVote](#oracle.v1.MsgOracleDataVote) | [MsgOracleDataVoteResponse](#oracle.v1.MsgOracleDataVoteResponse) |  | |

 <!-- end services -->
>>>>>>> 80d76df335316d61841ce9bca513aaf530ec40d3



<a name="oracle/v1/genesis.proto"></a>
<p align="right"><a href="#top">Top</a></p>

<<<<<<< HEAD
<a name="oracle.v1.MsgOracleDataVote"></a>
=======
## oracle/v1/genesis.proto
>>>>>>> 80d76df335316d61841ce9bca513aaf530ec40d3

### MsgOracleDataVote
MsgOracleDataVote - sdk.Msg for submitting arbitrary oracle data that has been prevoted on


<<<<<<< HEAD
| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `salt` | [string](#string) | repeated |  |
| `oracle_data` | [google.protobuf.Any](#google.protobuf.Any) | repeated |  |
| `signer` | [string](#string) |  |  |
=======
>>>>>>> 80d76df335316d61841ce9bca513aaf530ec40d3
=======
<a name="oracle.v1.GenesisState"></a>
>>>>>>> 80d76df335316d61841ce9bca513aaf530ec40d3
=======
>>>>>>> 5b665fcbd01d71b7bbcfcceba5a2d70aa2299577

### GenesisState
GenesisState - all oracle state that must be provided at genesis

<a name="oracle/v1/query.proto"></a>
<p align="right"><a href="#top">Top</a></p>

<<<<<<< HEAD
## oracle/v1/query.proto

=======
| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `params` | [Params](#oracle.v1.Params) |  |  |
| `feeder_delegations` | [MsgDelegateFeedConsent](#oracle.v1.MsgDelegateFeedConsent) | repeated |  |
| `miss_counters` | [MissCounter](#oracle.v1.MissCounter) | repeated |  |
>>>>>>> 80d76df335316d61841ce9bca513aaf530ec40d3


<<<<<<< HEAD
<a name="oracle.v1.MsgOracleDataVoteResponse"></a>

### MsgOracleDataVoteResponse



<<<<<<< HEAD
<a name="oracle.v1.QueryTobinTaxResponse"></a>

### QueryTobinTaxResponse
=======
<a name="oracle.v1.QueryDelegeateAddressRequest"></a>

### QueryDelegeateAddressRequest
>>>>>>> 80d76df335316d61841ce9bca513aaf530ec40d3
=======
<a name="oracle.v1.MissCounter"></a>
>>>>>>> 80d76df335316d61841ce9bca513aaf530ec40d3

### MissCounter
MissCounter stores the validator address and the number of associated misses


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
<<<<<<< HEAD
| `rate` | [string](#string) |  |  |

=======
| `validator` | [string](#string) |  |  |
| `misses` | [int64](#int64) |  |  |
>>>>>>> 80d76df335316d61841ce9bca513aaf530ec40d3





<<<<<<< HEAD
<<<<<<< HEAD
<a name="oracle.v1.QueryTobinTaxesRequest"></a>

### QueryTobinTaxesRequest
=======
<a name="oracle.v1.QueryDelegeateAddressResponse"></a>

### QueryDelegeateAddressResponse
>>>>>>> 80d76df335316d61841ce9bca513aaf530ec40d3
=======

<a name="oracle.v1.Params"></a>
>>>>>>> 80d76df335316d61841ce9bca513aaf530ec40d3

### Params
Params oracle parameters


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
<<<<<<< HEAD
<<<<<<< HEAD
| `pagination` | [cosmos.base.query.v1beta1.PageRequest](#cosmos.base.query.v1beta1.PageRequest) |  | pagination request |
=======
| `delegate` | [string](#string) |  |  |
>>>>>>> 80d76df335316d61841ce9bca513aaf530ec40d3






<a name="oracle.v1.QueryTobinTaxesResponse"></a>

### QueryTobinTaxesResponse


=======
| `vote_period` | [int64](#int64) |  | VotePeriod defines the number of blocks to wait for votes before attempting to tally |
| `vote_threshold` | [string](#string) |  | VoteThreshold defines the percentage of bonded stake required to vote each period |
| `slash_window` | [int64](#int64) |  | SlashWindow defines the number of blocks for the slashing window |
| `min_valid_per_window` | [string](#string) |  | MinValidPerWindow defines the number of misses a validator is allowed during each SlashWindow |
| `slash_fraction` | [string](#string) |  | SlashFraction defines the percentage of slash that a validator will suffer if it fails to send a vote |
| `data_types` | [string](#string) | repeated | DataTypes defines which data types validators must submit each voting period |
>>>>>>> 80d76df335316d61841ce9bca513aaf530ec40d3

| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `rates` | [cosmos.base.v1beta1.DecCoin](#cosmos.base.v1beta1.DecCoin) | repeated |  |
| `pagination` | [cosmos.base.query.v1beta1.PageResponse](#cosmos.base.query.v1beta1.PageResponse) |  | pagination response |






<<<<<<< HEAD
<a name="oracle.v1.QueryVoteTargetsRequest"></a>

### QueryVoteTargetsRequest
=======
 <!-- end messages -->

 <!-- end enums -->
>>>>>>> 80d76df335316d61841ce9bca513aaf530ec40d3

 <!-- end HasExtensions -->

 <!-- end services -->

<<<<<<< HEAD
| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
<<<<<<< HEAD
| `pagination` | [cosmos.base.query.v1beta1.PageRequest](#cosmos.base.query.v1beta1.PageRequest) |  | pagination request |
=======
>>>>>>> 80d76df335316d61841ce9bca513aaf530ec40d3


<a name="oracle/v1/query.proto"></a>
<p align="right"><a href="#top">Top</a></p>

## oracle/v1/query.proto



<<<<<<< HEAD
<a name="oracle.v1.QueryVoteTargetsResponse"></a>

### QueryVoteTargetsResponse
=======
<a name="oracle.v1.QueryDelegeateAddressRequest"></a>

### QueryDelegeateAddressRequest
>>>>>>> 80d76df335316d61841ce9bca513aaf530ec40d3



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `targets` | [string](#string) | repeated |  |
| `pagination` | [cosmos.base.query.v1beta1.PageResponse](#cosmos.base.query.v1beta1.PageResponse) |  | pagination response |

=======
| `miss_counter` | [int64](#int64) |  |  |
>>>>>>> 80d76df335316d61841ce9bca513aaf530ec40d3




 <!-- end messages -->

 <!-- end enums -->

<<<<<<< HEAD
 <!-- end HasExtensions -->


<a name="oracle.v1.Msg"></a>

### Msg
MsgService defines the msgs that the oracle module handles.

| Method Name | Request Type | Response Type | Description | HTTP Verb | Endpoint |
| ----------- | ------------ | ------------- | ------------| ------- | -------- |
| `DelegateFeedConsent` | [MsgDelegateFeedConsent](#oracle.v1.MsgDelegateFeedConsent) | [MsgDelegateFeedConsentResponse](#oracle.v1.MsgDelegateFeedConsentResponse) |  | |
| `OracleDataPrevote` | [MsgOracleDataPrevote](#oracle.v1.MsgOracleDataPrevote) | [MsgOracleDataPrevoteResponse](#oracle.v1.MsgOracleDataPrevoteResponse) |  | |
| `OracleDataVote` | [MsgOracleDataVote](#oracle.v1.MsgOracleDataVote) | [MsgOracleDataVoteResponse](#oracle.v1.MsgOracleDataVoteResponse) |  | |

 <!-- end services -->

=======
<a name="oracle.v1.QueryOracleDataPrevoteRequest"></a>

<<<<<<< HEAD
### QueryOracleDataPrevoteRequest
>>>>>>> 80d76df335316d61841ce9bca513aaf530ec40d3

=======
<a name="oracle.v1.QueryDelegeateAddressResponse"></a>

### QueryDelegeateAddressResponse
>>>>>>> 80d76df335316d61841ce9bca513aaf530ec40d3

<a name="oracle/v1/genesis.proto"></a>
<p align="right"><a href="#top">Top</a></p>

## oracle/v1/genesis.proto



<a name="oracle.v1.GenesisState"></a>

### GenesisState
GenesisState - all oracle state that must be provided at genesis


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
<<<<<<< HEAD
<<<<<<< HEAD
| `params` | [Params](#oracle.v1.Params) |  |  |
| `feeder_delegations` | [MsgDelegateFeedConsent](#oracle.v1.MsgDelegateFeedConsent) | repeated |  |
| `miss_counters` | [MissCounter](#oracle.v1.MissCounter) | repeated |  |
=======
| `delegate` | [string](#string) |  |  |

>>>>>>> 80d76df335316d61841ce9bca513aaf530ec40d3






<a name="oracle.v1.MissCounter"></a>

### MissCounter
MissCounter stores the validator address and the number of associated misses

=======
| `validator` | [string](#string) |  |  |
>>>>>>> 80d76df335316d61841ce9bca513aaf530ec40d3

| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `validator` | [string](#string) |  |  |
| `misses` | [int64](#int64) |  |  |





<<<<<<< HEAD

<a name="oracle.v1.Params"></a>
=======
<a name="oracle.v1.QueryOracleDataPrevoteResponse"></a>

### QueryOracleDataPrevoteResponse
>>>>>>> 80d76df335316d61841ce9bca513aaf530ec40d3

### Params
Params oracle parameters


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
<<<<<<< HEAD
<<<<<<< HEAD
| `vote_period` | [int64](#int64) |  | VotePeriod defines the number of blocks to wait for votes before attempting to tally |
| `vote_threshold` | [string](#string) |  | VoteThreshold defines the percentage of bonded stake required to vote each period |
| `slash_window` | [int64](#int64) |  | SlashWindow defines the number of blocks for the slashing window |
| `min_valid_per_window` | [string](#string) |  | MinValidPerWindow defines the number of misses a validator is allowed during each SlashWindow |
| `slash_fraction` | [string](#string) |  | SlashFraction defines the percentage of slash that a validator will suffer if it fails to send a vote |
| `data_types` | [string](#string) | repeated | DataTypes defines which data types validators must submit each voting period |
=======
| `hashes` | [bytes](#bytes) | repeated |  |
>>>>>>> 80d76df335316d61841ce9bca513aaf530ec40d3




=======
| `miss_counter` | [int64](#int64) |  |  |
>>>>>>> 80d76df335316d61841ce9bca513aaf530ec40d3

 <!-- end messages -->

<<<<<<< HEAD
 <!-- end enums -->

 <!-- end HasExtensions -->
=======

>>>>>>> 80d76df335316d61841ce9bca513aaf530ec40d3

 <!-- end services -->


<a name="oracle.v1.QueryOracleDataPrevoteRequest"></a>

<<<<<<< HEAD
<a name="oracle/v1/query.proto"></a>
<p align="right"><a href="#top">Top</a></p>

## oracle/v1/query.proto

<<<<<<< HEAD


<a name="oracle.v1.QueryDelegeateAddressRequest"></a>

### QueryDelegeateAddressRequest
=======
<a name="oracle.v1.QueryOracleDataRequest"></a>

### QueryOracleDataRequest
>>>>>>> 80d76df335316d61841ce9bca513aaf530ec40d3
=======
### QueryOracleDataPrevoteRequest
>>>>>>> 80d76df335316d61841ce9bca513aaf530ec40d3



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
<<<<<<< HEAD
<<<<<<< HEAD
=======
>>>>>>> 80d76df335316d61841ce9bca513aaf530ec40d3
| `validator` | [string](#string) |  |  |






<<<<<<< HEAD
<a name="oracle.v1.QueryDelegeateAddressResponse"></a>

### QueryDelegeateAddressResponse
=======
| `type` | [string](#string) |  |  |
>>>>>>> 80d76df335316d61841ce9bca513aaf530ec40d3
=======
<a name="oracle.v1.QueryOracleDataPrevoteResponse"></a>

### QueryOracleDataPrevoteResponse
>>>>>>> 80d76df335316d61841ce9bca513aaf530ec40d3



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
<<<<<<< HEAD
| `delegate` | [string](#string) |  |  |
=======
| `hashes` | [bytes](#bytes) | repeated |  |
>>>>>>> 80d76df335316d61841ce9bca513aaf530ec40d3




<a name="oracle.v1.QueryOracleDataResponse"></a>

<<<<<<< HEAD

<<<<<<< HEAD
<a name="oracle.v1.QueryMissCounterRequest"></a>

### QueryMissCounterRequest
=======
### QueryOracleDataResponse
>>>>>>> 80d76df335316d61841ce9bca513aaf530ec40d3
=======
<a name="oracle.v1.QueryOracleDataRequest"></a>

### QueryOracleDataRequest
>>>>>>> 80d76df335316d61841ce9bca513aaf530ec40d3



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
<<<<<<< HEAD
<<<<<<< HEAD
| `validator` | [string](#string) |  |  |






<a name="oracle.v1.QueryMissCounterResponse"></a>

### QueryMissCounterResponse
=======
| `oracle_data` | [google.protobuf.Any](#google.protobuf.Any) |  |  |
=======
| `type` | [string](#string) |  |  |
>>>>>>> 80d76df335316d61841ce9bca513aaf530ec40d3

>>>>>>> 80d76df335316d61841ce9bca513aaf530ec40d3



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `miss_counter` | [int64](#int64) |  |  |


<<<<<<< HEAD
<a name="oracle.v1.QueryOracleDataVoteRequest"></a>

<<<<<<< HEAD



<a name="oracle.v1.QueryOracleDataPrevoteRequest"></a>

### QueryOracleDataPrevoteRequest
=======
### QueryOracleDataVoteRequest
>>>>>>> 80d76df335316d61841ce9bca513aaf530ec40d3
=======
<a name="oracle.v1.QueryOracleDataResponse"></a>

### QueryOracleDataResponse
>>>>>>> 80d76df335316d61841ce9bca513aaf530ec40d3



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
<<<<<<< HEAD
<<<<<<< HEAD
| `validator` | [string](#string) |  |  |






<<<<<<< HEAD
<a name="oracle.v1.QueryOracleDataPrevoteResponse"></a>

### QueryOracleDataPrevoteResponse
=======
=======
>>>>>>> 5b665fcbd01d71b7bbcfcceba5a2d70aa2299577
| `oracle_data` | [google.protobuf.Any](#google.protobuf.Any) |  |  |

>>>>>>> 80d76df335316d61841ce9bca513aaf530ec40d3



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `hashes` | [bytes](#bytes) | repeated |  |
=======
<a name="oracle.v1.QueryParamsRequest"></a>

### QueryParamsRequest

<a name="oracle.v1.QueryOracleDataVoteRequest"></a>

<<<<<<< HEAD

>>>>>>> 80d76df335316d61841ce9bca513aaf530ec40d3




<a name="oracle.v1.QueryParamsResponse"></a>
=======
### QueryOracleDataVoteRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `validator` | [string](#string) |  |  |
>>>>>>> 80d76df335316d61841ce9bca513aaf530ec40d3

### QueryParamsResponse

<<<<<<< HEAD
<a name="oracle.v1.QueryOracleDataRequest"></a>

### QueryOracleDataRequest



<<<<<<< HEAD
| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `type` | [string](#string) |  |  |

=======
<a name="oracle.v1.QueryParamsRequest"></a>

### QueryParamsRequest
>>>>>>> 80d76df335316d61841ce9bca513aaf530ec40d3





<<<<<<< HEAD
<a name="oracle.v1.QueryOracleDataResponse"></a>

### QueryOracleDataResponse
=======

>>>>>>> 80d76df335316d61841ce9bca513aaf530ec40d3

<a name="oracle.v1.QueryParamsResponse"></a>

### QueryParamsResponse

<<<<<<< HEAD
| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `oracle_data` | [google.protobuf.Any](#google.protobuf.Any) |  |  |
=======
>>>>>>> 80d76df335316d61841ce9bca513aaf530ec40d3


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `params` | [Params](#oracle.v1.Params) |  |  |

=======
>>>>>>> 80d76df335316d61841ce9bca513aaf530ec40d3


<<<<<<< HEAD
| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `params` | [Params](#oracle.v1.Params) |  |  |

<a name="oracle.v1.QueryOracleDataVoteRequest"></a>

<<<<<<< HEAD
### QueryOracleDataVoteRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `validator` | [string](#string) |  |  |
<<<<<<< HEAD
=======

=======


>>>>>>> 80d76df335316d61841ce9bca513aaf530ec40d3

<a name="oracle.v1.QueryValidatorAddressRequest"></a>
=======
>>>>>>> 5b665fcbd01d71b7bbcfcceba5a2d70aa2299577

<<<<<<< HEAD
>>>>>>> 80d76df335316d61841ce9bca513aaf530ec40d3
=======
### QueryValidatorAddressRequest
>>>>>>> 80d76df335316d61841ce9bca513aaf530ec40d3

<a name="oracle.v1.QueryValidatorAddressRequest"></a>

### QueryValidatorAddressRequest

<<<<<<< HEAD


<<<<<<< HEAD
=======
| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `delegate` | [string](#string) |  |  |
>>>>>>> 80d76df335316d61841ce9bca513aaf530ec40d3

<a name="oracle.v1.QueryParamsRequest"></a>

### QueryParamsRequest
=======
| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `delegate` | [string](#string) |  |  |
>>>>>>> 80d76df335316d61841ce9bca513aaf530ec40d3





<a name="oracle.v1.QueryValidatorAddressResponse"></a>

<<<<<<< HEAD
<a name="oracle.v1.QueryValidatorAddressResponse"></a>

<<<<<<< HEAD
<a name="oracle.v1.QueryParamsResponse"></a>

### QueryParamsResponse
=======
### QueryValidatorAddressResponse
>>>>>>> 80d76df335316d61841ce9bca513aaf530ec40d3
=======
### QueryValidatorAddressResponse
>>>>>>> 80d76df335316d61841ce9bca513aaf530ec40d3



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
<<<<<<< HEAD
<<<<<<< HEAD
| `params` | [Params](#oracle.v1.Params) |  |  |
=======
=======
>>>>>>> 80d76df335316d61841ce9bca513aaf530ec40d3
| `validator` | [string](#string) |  |  |
>>>>>>> 80d76df335316d61841ce9bca513aaf530ec40d3






<<<<<<< HEAD
<a name="oracle.v1.QueryValidatorAddressRequest"></a>

### QueryValidatorAddressRequest
=======
<a name="oracle.v1.QueryVotePeriodRequest"></a>

<<<<<<< HEAD
### QueryVotePeriodRequest
>>>>>>> 80d76df335316d61841ce9bca513aaf530ec40d3

=======
<a name="oracle.v1.QueryVotePeriodRequest"></a>

### QueryVotePeriodRequest
>>>>>>> 80d76df335316d61841ce9bca513aaf530ec40d3


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `delegate` | [string](#string) |  |  |




<<<<<<< HEAD

<<<<<<< HEAD

<a name="oracle.v1.QueryValidatorAddressResponse"></a>

### QueryValidatorAddressResponse

=======
<a name="oracle.v1.VotePeriod"></a>

### VotePeriod

>>>>>>> 80d76df335316d61841ce9bca513aaf530ec40d3
=======
<a name="oracle.v1.VotePeriod"></a>

### VotePeriod

>>>>>>> 80d76df335316d61841ce9bca513aaf530ec40d3


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
<<<<<<< HEAD
<<<<<<< HEAD
| `validator` | [string](#string) |  |  |






<a name="oracle.v1.QueryVotePeriodRequest"></a>

### QueryVotePeriodRequest


=======
| `current_height` | [int64](#int64) |  |  |
| `vote_period_start` | [int64](#int64) |  |  |
| `vote_period_end` | [int64](#int64) |  |  |
>>>>>>> 80d76df335316d61841ce9bca513aaf530ec40d3

=======
| `current_height` | [int64](#int64) |  |  |
| `vote_period_start` | [int64](#int64) |  |  |
| `vote_period_end` | [int64](#int64) |  |  |
>>>>>>> 80d76df335316d61841ce9bca513aaf530ec40d3




<a name="oracle.v1.VotePeriod"></a>

### VotePeriod



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `current_height` | [int64](#int64) |  |  |
| `vote_period_start` | [int64](#int64) |  |  |
| `vote_period_end` | [int64](#int64) |  |  |





 <!-- end messages -->

 <!-- end enums -->

 <!-- end HasExtensions -->


<a name="oracle.v1.Query"></a>

### Query


| Method Name | Request Type | Response Type | Description | HTTP Verb | Endpoint |
| ----------- | ------------ | ------------- | ------------| ------- | -------- |
| `QueryParams` | [QueryParamsRequest](#oracle.v1.QueryParamsRequest) | [QueryParamsResponse](#oracle.v1.QueryParamsResponse) |  | |
| `QueryDelegeateAddress` | [QueryDelegeateAddressRequest](#oracle.v1.QueryDelegeateAddressRequest) | [QueryDelegeateAddressResponse](#oracle.v1.QueryDelegeateAddressResponse) |  | |
| `QueryValidatorAddress` | [QueryValidatorAddressRequest](#oracle.v1.QueryValidatorAddressRequest) | [QueryValidatorAddressResponse](#oracle.v1.QueryValidatorAddressResponse) |  | |
| `QueryOracleDataPrevote` | [QueryOracleDataPrevoteRequest](#oracle.v1.QueryOracleDataPrevoteRequest) | [QueryOracleDataPrevoteResponse](#oracle.v1.QueryOracleDataPrevoteResponse) |  | |
| `QueryOracleDataVote` | [QueryOracleDataVoteRequest](#oracle.v1.QueryOracleDataVoteRequest) | [MsgOracleDataVote](#oracle.v1.MsgOracleDataVote) |  | |
| `QueryVotePeriod` | [QueryVotePeriodRequest](#oracle.v1.QueryVotePeriodRequest) | [VotePeriod](#oracle.v1.VotePeriod) |  | |
| `QueryMissCounter` | [QueryMissCounterRequest](#oracle.v1.QueryMissCounterRequest) | [QueryMissCounterResponse](#oracle.v1.QueryMissCounterResponse) |  | |
| `OracleData` | [QueryOracleDataRequest](#oracle.v1.QueryOracleDataRequest) | [QueryOracleDataResponse](#oracle.v1.QueryOracleDataResponse) |  | |

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

