<!-- This file is auto-generated. Please do not modify it yourself. -->
# Protobuf Documentation
<a name="top"></a>

## Table of Contents

- [allocation/v1/allocation.proto](#allocation/v1/allocation.proto)
    - [Allocation](#allocation.v1.Allocation)
    - [AllocationPrecommit](#allocation.v1.AllocationPrecommit)
    - [Cellar](#allocation.v1.Cellar)
    - [CreateCellarsProposal](#allocation.v1.CreateCellarsProposal)
    - [Pool](#allocation.v1.Pool)
    - [Tick](#allocation.v1.Tick)
    - [TickWeight](#allocation.v1.TickWeight)
  
- [allocation/v1/tx.proto](#allocation/v1/tx.proto)
    - [MsgAllocationCommit](#allocation.v1.MsgAllocationCommit)
    - [MsgAllocationCommitResponse](#allocation.v1.MsgAllocationCommitResponse)
    - [MsgAllocationPrecommit](#allocation.v1.MsgAllocationPrecommit)
    - [MsgAllocationPrecommitResponse](#allocation.v1.MsgAllocationPrecommitResponse)
    - [MsgDelegateAllocations](#allocation.v1.MsgDelegateAllocations)
    - [MsgDelegateAllocationsResponse](#allocation.v1.MsgDelegateAllocationsResponse)
  
    - [Msg](#allocation.v1.Msg)
  
- [allocation/v1/genesis.proto](#allocation/v1/genesis.proto)
    - [AggregatedAllocationData](#allocation.v1.AggregatedAllocationData)
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
| `fee_level` | [string](#string) |  | sdk.Dec |
| `tick_weights` | [TickWeight](#allocation.v1.TickWeight) | repeated |  |






<a name="allocation.v1.AllocationPrecommit"></a>

### AllocationPrecommit
AllocationPrecommit defines an array of hashed decision data
that is used for the precommit phase of allocation


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `hash` | [bytes](#bytes) |  |  |






<a name="allocation.v1.Cellar"></a>

### Cellar
Cellar is the XXX


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `cellar_id` | [string](#string) |  |  |
| `token0` | [string](#string) |  |  |
| `token1` | [string](#string) |  |  |
| `pool` | [Pool](#allocation.v1.Pool) | repeated |  |






<a name="allocation.v1.CreateCellarsProposal"></a>

### CreateCellarsProposal
CreateCellarsProposal is a governance proposal content type


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `title` | [string](#string) |  |  |
| `description` | [string](#string) |  |  |
| `cellars` | [Cellar](#allocation.v1.Cellar) | repeated |  |






<a name="allocation.v1.Pool"></a>

### Pool
Pool is the XXX


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `fee_level` | [string](#string) |  | sdk.Dec |
| `tick_ranges` | [Tick](#allocation.v1.Tick) | repeated |  |






<a name="allocation.v1.Tick"></a>

### Tick
Tick is the XXX


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `min` | [uint64](#uint64) |  |  |
| `max` | [uint64](#uint64) |  |  |






<a name="allocation.v1.TickWeight"></a>

### TickWeight
TickWeight is the XXX


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `tick` | [Tick](#allocation.v1.Tick) |  |  |
| `weight` | [string](#string) |  | sdk.Dec |





 <!-- end messages -->

 <!-- end enums -->

 <!-- end HasExtensions -->

 <!-- end services -->



<a name="allocation/v1/tx.proto"></a>
<p align="right"><a href="#top">Top</a></p>

## allocation/v1/tx.proto



<a name="allocation.v1.MsgAllocationCommit"></a>

### MsgAllocationCommit
MsgAllocationCommit is the request type for AllocationCommit


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `allocations` | [Allocation](#allocation.v1.Allocation) | repeated |  |
| `salt` | [string](#string) |  |  |
| `signer` | [string](#string) |  |  |






<a name="allocation.v1.MsgAllocationCommitResponse"></a>

### MsgAllocationCommitResponse
MsgAllocationCommitResponse is the response type for the Msg/AllocationCommit gRPC method.






<a name="allocation.v1.MsgAllocationPrecommit"></a>

### MsgAllocationPrecommit
MsgAllocationPrecommit is the request type for AllocationPrecommit


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `precommit` | [AllocationPrecommit](#allocation.v1.AllocationPrecommit) |  |  |
| `signer` | [string](#string) |  |  |






<a name="allocation.v1.MsgAllocationPrecommitResponse"></a>

### MsgAllocationPrecommitResponse
MsgAllocationPrecommitResponse is the response type for MsgAllocationPrecommit






<a name="allocation.v1.MsgDelegateAllocations"></a>

### MsgDelegateAllocations
MsgDelegateAllocations is the request type for DelegateAllocations


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `delegate` | [string](#string) |  |  |
| `validator` | [string](#string) |  |  |






<a name="allocation.v1.MsgDelegateAllocationsResponse"></a>

### MsgDelegateAllocationsResponse
MsgDelegateAllocationsResponse is the response type for DelegateAllocations





 <!-- end messages -->

 <!-- end enums -->

 <!-- end HasExtensions -->


<a name="allocation.v1.Msg"></a>

### Msg
MsgService defines the messages the allocation module handles

| Method Name | Request Type | Response Type | Description | HTTP Verb | Endpoint |
| ----------- | ------------ | ------------- | ------------| ------- | -------- |
| `DelegateAllocations` | [MsgDelegateAllocations](#allocation.v1.MsgDelegateAllocations) | [MsgDelegateAllocationsResponse](#allocation.v1.MsgDelegateAllocationsResponse) | DelegateAllocations creates an index in the store linking the validator and the delegate key need to be able to query both the delegate and the validator given the other one | |
| `AllocationPrecommit` | [MsgAllocationPrecommit](#allocation.v1.MsgAllocationPrecommit) | [MsgAllocationPrecommitResponse](#allocation.v1.MsgAllocationPrecommitResponse) | AllocationPrecommit stores the precommit hash indexed by validator address | |
| `AllocationCommit` | [MsgAllocationCommit](#allocation.v1.MsgAllocationCommit) | [MsgAllocationCommitResponse](#allocation.v1.MsgAllocationCommitResponse) | AllocationCommit checks the precommit hash against the data, rejects the message if it doesn't match then records the commitment in the store indexed by validator address | |

 <!-- end services -->



<a name="allocation/v1/genesis.proto"></a>
<p align="right"><a href="#top">Top</a></p>

## allocation/v1/genesis.proto



<a name="allocation.v1.AggregatedAllocationData"></a>

### AggregatedAllocationData
AggregatedallocationData defines the aggregated allocation data at a given block height


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `height` | [int64](#int64) |  | block height in which the data was committed |
| `allocation` | [Allocation](#allocation.v1.Allocation) |  | allocation data |






<a name="allocation.v1.GenesisState"></a>

### GenesisState
GenesisState - all allocation state that must be provided at genesis


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `params` | [Params](#allocation.v1.Params) |  |  |
| `delegations` | [MsgDelegateAllocations](#allocation.v1.MsgDelegateAllocations) | repeated |  |
| `miss_counters` | [MissCounter](#allocation.v1.MissCounter) | repeated |  |
| `aggregates` | [AggregatedAllocationData](#allocation.v1.AggregatedAllocationData) | repeated |  |






<a name="allocation.v1.MissCounter"></a>

### MissCounter
MissCounter stores the validator address and the number of associated misses


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `validator` | [string](#string) |  | validator operator address |
| `misses` | [int64](#int64) |  | number of misses |






<a name="allocation.v1.Params"></a>

### Params
Params allocation parameters


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `vote_period` | [int64](#int64) |  | VotePeriod defines the number of blocks to wait for votes before attempting to tally |
| `vote_threshold` | [string](#string) |  | VoteThreshold defines the percentage of bonded stake required to vote each period |
| `slash_window` | [int64](#int64) |  | SlashWindow defines the number of blocks for the slashing window |
| `min_valid_per_window` | [string](#string) |  | MinValidPerWindow defines the number of misses a validator is allowed during each SlashWindow |
| `slash_fraction` | [string](#string) |  | SlashFraction defines the percentage of slash that a validator will suffer if it fails to send a vote |
| `target_threshold` | [string](#string) |  | TargetThreshold defines the max percentage difference that a given allocation data needs to have with the aggregated data in order for the feeder to be elegible for rewards. |
| `data_types` | [string](#string) | repeated | DataTypes defines which data types validators must submit each voting period |





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
| `type` | [string](#string) |  | allocation data type |
| `id` | [string](#string) |  | allocation data identifier |






<a name="allocation.v1.QueryAggregateDataResponse"></a>

### QueryAggregateDataResponse
QueryAggregateDataRequest is the response type for the Query/AggregateData gRPC method.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `allocation_data` | [Allocation](#allocation.v1.Allocation) |  | allocation data associated with the id and type from the request |
| `height` | [int64](#int64) |  | height at which the aggregated allocation data was stored |






<a name="allocation.v1.QueryAllocationCommitRequest"></a>

### QueryAllocationCommitRequest
QueryAllocationCommitRequest is the request type for the Query/QueryAllocationCommitRequest gRPC method.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `validator` | [string](#string) |  | validator operator address |






<a name="allocation.v1.QueryAllocationCommitResponse"></a>

### QueryAllocationCommitResponse
QueryAllocationCommitResponse is the response type for the Query/QueryAllocationCommitResponse gRPC method.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `allocation` | [Allocation](#allocation.v1.Allocation) |  | allocation containing the allocation feed submitted within the latest voting period |






<a name="allocation.v1.QueryAllocationPrecommitRequest"></a>

### QueryAllocationPrecommitRequest
QueryAllocationPrecommitRequest is the request type for the Query/QueryAllocationPrecommitRequest gRPC method.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `validator` | [string](#string) |  | validator operator address |






<a name="allocation.v1.QueryAllocationPrecommitResponse"></a>

### QueryAllocationPrecommitResponse
QueryAllocationPrecommitResponse is the response type for the Query/QueryAllocationPrecommitResponse gRPC method.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `precommit` | [AllocationPrecommit](#allocation.v1.AllocationPrecommit) |  | prevote submitted within the latest voting period |






<a name="allocation.v1.QueryCommitPeriodRequest"></a>

### QueryCommitPeriodRequest
QueryCommitPeriodRequest is the request type for the Query/VotePeriod gRPC method.






<a name="allocation.v1.QueryCommitPeriodResponse"></a>

### QueryCommitPeriodResponse
QueryCommitPeriodResponse is the response type for the Query/VotePeriod gRPC method.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `current_height` | [int64](#int64) |  | block height at which the query was processed |
| `commit_period_start` | [int64](#int64) |  | latest commit period start block height |
| `commit_period_end` | [int64](#int64) |  | block height at which the current commit period ends |






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
| `allocation_data` | [Allocation](#allocation.v1.Allocation) | repeated | allocation data associated with the |
| `height` | [int64](#int64) |  | height at which the aggregated allocation data was stored |
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
| `miss_counter` | [int64](#int64) |  | number of allocation feed votes missed since the last counter reset |






<a name="allocation.v1.QueryParamsRequest"></a>

### QueryParamsRequest
QueryParamsRequest is the request type for the Query/Params gRPC method.






<a name="allocation.v1.QueryParamsResponse"></a>

### QueryParamsResponse
QueryParamsRequest is the response type for the Query/Params gRPC method.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `params` | [Params](#allocation.v1.Params) |  | allocation parameters |






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
Query defines the gRPC querier service for the allocation module.

| Method Name | Request Type | Response Type | Description | HTTP Verb | Endpoint |
| ----------- | ------------ | ------------- | ------------| ------- | -------- |
| `QueryParams` | [QueryParamsRequest](#allocation.v1.QueryParamsRequest) | [QueryParamsResponse](#allocation.v1.QueryParamsResponse) | QueryParams queries the allocation module parameters. | GET|/sommelier/allocation/v1/params|
| `QueryDelegateAddress` | [QueryDelegateAddressRequest](#allocation.v1.QueryDelegateAddressRequest) | [QueryDelegateAddressResponse](#allocation.v1.QueryDelegateAddressResponse) | QueryDelegateAddress queries the delegate account address of a validator | GET|/sommelier/allocation/v1/delegates/{validator}|
| `QueryValidatorAddress` | [QueryValidatorAddressRequest](#allocation.v1.QueryValidatorAddressRequest) | [QueryValidatorAddressResponse](#allocation.v1.QueryValidatorAddressResponse) | QueryValidatorAddress returns the validator address of a given delegate | GET|/sommelier/allocation/v1/validators/{delegate}|
| `QueryAllocationPrecommit` | [QueryAllocationPrecommitRequest](#allocation.v1.QueryAllocationPrecommitRequest) | [QueryAllocationPrecommitResponse](#allocation.v1.QueryAllocationPrecommitResponse) | QueryAllocationPrecommit queries the validator precommit in the current voting period | GET|/sommelier/allocation/v1/prevotes/{validator}|
| `QueryAllocationCommit` | [QueryAllocationCommitRequest](#allocation.v1.QueryAllocationCommitRequest) | [QueryAllocationCommitResponse](#allocation.v1.QueryAllocationCommitResponse) | QueryAllocationCommit queries the validator allocation in the current voting period | GET|/sommelier/allocation/v1/votes/{validator}|
| `QueryCommitPeriod` | [QueryCommitPeriodRequest](#allocation.v1.QueryCommitPeriodRequest) | [QueryCommitPeriodResponse](#allocation.v1.QueryCommitPeriodResponse) | QueryCommitPeriod queries the heights for the current commit period (current, start and end) | GET|/sommelier/allocation/v1/vote_period|
| `QueryMissCounter` | [QueryMissCounterRequest](#allocation.v1.QueryMissCounterRequest) | [QueryMissCounterResponse](#allocation.v1.QueryMissCounterResponse) | QueryMissCounter queries the missed number of allocation data feed periods | GET|/sommelier/allocation/v1/miss_counters/{validator}|
| `QueryAggregateData` | [QueryAggregateDataRequest](#allocation.v1.QueryAggregateDataRequest) | [QueryAggregateDataResponse](#allocation.v1.QueryAggregateDataResponse) | QueryAggregateData returns the latest aggregated data value for a given type and identifioer | GET|/sommelier/allocation/v1/aggregate_data/{id}/{type}|
| `QueryLatestPeriodAggregateData` | [QueryLatestPeriodAggregateDataRequest](#allocation.v1.QueryLatestPeriodAggregateDataRequest) | [QueryLatestPeriodAggregateDataResponse](#allocation.v1.QueryLatestPeriodAggregateDataResponse) | QueryLatestPeriodAggregateData returns the aggregated data for a given pair an identifioer | GET|/sommelier/allocation/v1/aggregate_data|

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

