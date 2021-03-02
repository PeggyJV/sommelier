<!-- This file is auto-generated. Please do not modify it yourself. -->
# Protobuf Documentation
<a name="top"></a>

## Table of Contents

- [oracle/v1/oracle.proto](#oracle/v1/oracle.proto)
    - [OracleFeed](#oracle.v1.OracleFeed)
    - [OraclePrevote](#oracle.v1.OraclePrevote)
    - [OracleVote](#oracle.v1.OracleVote)
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
    - [AggregatedOracleData](#oracle.v1.AggregatedOracleData)
    - [GenesisState](#oracle.v1.GenesisState)
    - [MissCounter](#oracle.v1.MissCounter)
    - [Params](#oracle.v1.Params)
  
- [oracle/v1/query.proto](#oracle/v1/query.proto)
    - [QueryAggregateDataRequest](#oracle.v1.QueryAggregateDataRequest)
    - [QueryAggregateDataResponse](#oracle.v1.QueryAggregateDataResponse)
    - [QueryDelegateAddressRequest](#oracle.v1.QueryDelegateAddressRequest)
    - [QueryDelegateAddressResponse](#oracle.v1.QueryDelegateAddressResponse)
    - [QueryMissCounterRequest](#oracle.v1.QueryMissCounterRequest)
    - [QueryMissCounterResponse](#oracle.v1.QueryMissCounterResponse)
    - [QueryOracleDataPrevoteRequest](#oracle.v1.QueryOracleDataPrevoteRequest)
    - [QueryOracleDataPrevoteResponse](#oracle.v1.QueryOracleDataPrevoteResponse)
    - [QueryOracleDataVoteRequest](#oracle.v1.QueryOracleDataVoteRequest)
    - [QueryOracleDataVoteResponse](#oracle.v1.QueryOracleDataVoteResponse)
    - [QueryParamsRequest](#oracle.v1.QueryParamsRequest)
    - [QueryParamsResponse](#oracle.v1.QueryParamsResponse)
    - [QueryValidatorAddressRequest](#oracle.v1.QueryValidatorAddressRequest)
    - [QueryValidatorAddressResponse](#oracle.v1.QueryValidatorAddressResponse)
    - [QueryVotePeriodRequest](#oracle.v1.QueryVotePeriodRequest)
    - [QueryVotePeriodResponse](#oracle.v1.QueryVotePeriodResponse)
  
    - [Query](#oracle.v1.Query)
  
- [Scalar Value Types](#scalar-value-types)



<a name="oracle/v1/oracle.proto"></a>
<p align="right"><a href="#top">Top</a></p>

## oracle/v1/oracle.proto



<a name="oracle.v1.OracleFeed"></a>

### OracleFeed
OracleFeed represents an array of oracle data that is


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `data` | [UniswapPair](#oracle.v1.UniswapPair) | repeated |  |






<a name="oracle.v1.OraclePrevote"></a>

### OraclePrevote
OraclePrevote defines an array of hashed from oracle data that are used
for the prevote phase of the oracle data feeding.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `hash` | [bytes](#bytes) |  | hex formated hash of an oracle feed |






<a name="oracle.v1.OracleVote"></a>

### OracleVote
UniswapToken is the returned uniswap token representation


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `salt` | [string](#string) |  |  |
| `feed` | [OracleFeed](#oracle.v1.OracleFeed) |  |  |






<a name="oracle.v1.UniswapPair"></a>

### UniswapPair
UniswapPair represents an SDK compatible uniswap pair info fetched from The Graph.


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
| `id` | [string](#string) |  | token address |
| `decimals` | [uint64](#uint64) |  | number of decimal positions of the pair token |





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






<a name="oracle.v1.MsgDelegateFeedConsentResponse"></a>

### MsgDelegateFeedConsentResponse







<a name="oracle.v1.MsgOracleDataPrevote"></a>

### MsgOracleDataPrevote
MsgOracleDataPrevote - sdk.Msg for prevoting on an array of oracle data types.
The purpose of the prevote is to hide vote for data with hashes formatted as hex string: 
SHA256("{salt}:{data_cannonical_json}:{voter}")


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `prevote` | [OraclePrevote](#oracle.v1.OraclePrevote) |  |  |
| `signer` | [string](#string) |  |  |






<a name="oracle.v1.MsgOracleDataPrevoteResponse"></a>

### MsgOracleDataPrevoteResponse







<a name="oracle.v1.MsgOracleDataVote"></a>

### MsgOracleDataVote
MsgOracleDataVote - sdk.Msg for submitting arbitrary oracle data that has been prevoted on


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `vote` | [OracleVote](#oracle.v1.OracleVote) |  |  |
| `signer` | [string](#string) |  |  |






<a name="oracle.v1.MsgOracleDataVoteResponse"></a>

### MsgOracleDataVoteResponse






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



<a name="oracle.v1.AggregatedOracleData"></a>

### AggregatedOracleData
AggregatedOracleData defines the aggregated oracle data at a given block height


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `height` | [int64](#int64) |  | block height in which the data was committed |
| `data` | [UniswapPair](#oracle.v1.UniswapPair) |  | oracle data |






<a name="oracle.v1.GenesisState"></a>

### GenesisState
GenesisState - all oracle state that must be provided at genesis


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `params` | [Params](#oracle.v1.Params) |  |  |
| `feeder_delegations` | [MsgDelegateFeedConsent](#oracle.v1.MsgDelegateFeedConsent) | repeated |  |
| `miss_counters` | [MissCounter](#oracle.v1.MissCounter) | repeated |  |
| `aggregates` | [AggregatedOracleData](#oracle.v1.AggregatedOracleData) | repeated |  |






<a name="oracle.v1.MissCounter"></a>

### MissCounter
MissCounter stores the validator address and the number of associated misses


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `validator` | [string](#string) |  |  |
| `misses` | [int64](#int64) |  |  |






<a name="oracle.v1.Params"></a>

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





 <!-- end messages -->

 <!-- end enums -->

 <!-- end HasExtensions -->

 <!-- end services -->



<a name="oracle/v1/query.proto"></a>
<p align="right"><a href="#top">Top</a></p>

## oracle/v1/query.proto



<a name="oracle.v1.QueryAggregateDataRequest"></a>

### QueryAggregateDataRequest
QueryAggregateDataRequest is the request type for the Query/AggregateData RPC method.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `type` | [string](#string) |  |  |
| `id` | [string](#string) |  |  |






<a name="oracle.v1.QueryAggregateDataResponse"></a>

### QueryAggregateDataResponse



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `oracle_data` | [UniswapPair](#oracle.v1.UniswapPair) |  |  |






<a name="oracle.v1.QueryDelegateAddressRequest"></a>

### QueryDelegateAddressRequest
QueryParamsRequest is the request type for the Query/Params RPC method.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `validator` | [string](#string) |  |  |






<a name="oracle.v1.QueryDelegateAddressResponse"></a>

### QueryDelegateAddressResponse



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `delegate` | [string](#string) |  |  |






<a name="oracle.v1.QueryMissCounterRequest"></a>

### QueryMissCounterRequest
QueryMissCounterRequest is the request type for the Query/MissCounter RPC method.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `validator` | [string](#string) |  |  |






<a name="oracle.v1.QueryMissCounterResponse"></a>

### QueryMissCounterResponse



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `miss_counter` | [int64](#int64) |  |  |






<a name="oracle.v1.QueryOracleDataPrevoteRequest"></a>

### QueryOracleDataPrevoteRequest
QueryParamsRequest is the request type for the Query/Params RPC method.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `validator` | [string](#string) |  |  |






<a name="oracle.v1.QueryOracleDataPrevoteResponse"></a>

### QueryOracleDataPrevoteResponse



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `prevote` | [OraclePrevote](#oracle.v1.OraclePrevote) |  |  |






<a name="oracle.v1.QueryOracleDataVoteRequest"></a>

### QueryOracleDataVoteRequest
QueryOracleDataVoteRequest is the request type for the Query/Odra RPC method.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `validator` | [string](#string) |  |  |






<a name="oracle.v1.QueryOracleDataVoteResponse"></a>

### QueryOracleDataVoteResponse



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `vote` | [OracleVote](#oracle.v1.OracleVote) |  |  |






<a name="oracle.v1.QueryParamsRequest"></a>

### QueryParamsRequest
QueryParamsRequest is the request type for the Query/Params RPC method.






<a name="oracle.v1.QueryParamsResponse"></a>

### QueryParamsResponse



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `params` | [Params](#oracle.v1.Params) |  |  |






<a name="oracle.v1.QueryValidatorAddressRequest"></a>

### QueryValidatorAddressRequest
QueryParamsRequest is the request type for the Query/Params RPC method.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `delegate` | [string](#string) |  |  |






<a name="oracle.v1.QueryValidatorAddressResponse"></a>

### QueryValidatorAddressResponse



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `validator` | [string](#string) |  |  |






<a name="oracle.v1.QueryVotePeriodRequest"></a>

### QueryVotePeriodRequest
QueryVotePeriodRequest is the request type for the Query/VotePeriod RPC method.






<a name="oracle.v1.QueryVotePeriodResponse"></a>

### QueryVotePeriodResponse



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
Query defines the gRPC querier service for the oracle module.

| Method Name | Request Type | Response Type | Description | HTTP Verb | Endpoint |
| ----------- | ------------ | ------------- | ------------| ------- | -------- |
| `QueryParams` | [QueryParamsRequest](#oracle.v1.QueryParamsRequest) | [QueryParamsResponse](#oracle.v1.QueryParamsResponse) | Params queries the oracle module parameters. | GET|/sommelier/oracle/v1/params|
| `QueryDelegateAddress` | [QueryDelegateAddressRequest](#oracle.v1.QueryDelegateAddressRequest) | [QueryDelegateAddressResponse](#oracle.v1.QueryDelegateAddressResponse) | QueryDelegateAddress queries the delegate account address of a validator | GET|/sommelier/oracle/v1/delegates/{validator}|
| `QueryValidatorAddress` | [QueryValidatorAddressRequest](#oracle.v1.QueryValidatorAddressRequest) | [QueryValidatorAddressResponse](#oracle.v1.QueryValidatorAddressResponse) | QueryValidatorAddress returns the validator address of a given delegate | GET|/sommelier/oracle/v1/validators/{delegate}|
| `QueryOracleDataPrevote` | [QueryOracleDataPrevoteRequest](#oracle.v1.QueryOracleDataPrevoteRequest) | [QueryOracleDataPrevoteResponse](#oracle.v1.QueryOracleDataPrevoteResponse) | QueryOracleDataPrevote queries the validator prevote in the current voting period | GET|/sommelier/oracle/v1/prevotes/{validator}|
| `QueryOracleDataVote` | [QueryOracleDataVoteRequest](#oracle.v1.QueryOracleDataVoteRequest) | [QueryOracleDataVoteResponse](#oracle.v1.QueryOracleDataVoteResponse) | QueryOracleDataVote queries the validator vote in the current voting period | GET|/sommelier/oracle/v1/votes/{validator}|
| `QueryVotePeriod` | [QueryVotePeriodRequest](#oracle.v1.QueryVotePeriodRequest) | [QueryVotePeriodResponse](#oracle.v1.QueryVotePeriodResponse) | QueryVotePeriod queries the heights for the current voting period (current, start and end) | GET|/sommelier/oracle/v1/vote_period|
| `QueryMissCounter` | [QueryMissCounterRequest](#oracle.v1.QueryMissCounterRequest) | [QueryMissCounterResponse](#oracle.v1.QueryMissCounterResponse) | QueryMissCounter queries the missed number of oracle data feed periods | GET|/sommelier/oracle/v1/miss_counters/{validator}|
| `QueryAggregateData` | [QueryAggregateDataRequest](#oracle.v1.QueryAggregateDataRequest) | [QueryAggregateDataResponse](#oracle.v1.QueryAggregateDataResponse) | QueryAggregateData returns the aggregated data for a given pair an identifioer | GET|/sommelier/oracle/v1/aggregate_data/{type}/{id}|

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

