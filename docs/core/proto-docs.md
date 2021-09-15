<!-- This file is auto-generated. Please do not modify it yourself. -->
# Protobuf Documentation
<a name="top"></a>

## Table of Contents

- [allocation/v1/allocation.proto](#allocation/v1/allocation.proto)
    - [AddManagedCellarsProposal](#allocation.v1.AddManagedCellarsProposal)
    - [Allocation](#allocation.v1.Allocation)
    - [AllocationPrecommit](#allocation.v1.AllocationPrecommit)
    - [Cellar](#allocation.v1.Cellar)
    - [CellarUpdate](#allocation.v1.CellarUpdate)
    - [RemoveManagedCellarsProposal](#allocation.v1.RemoveManagedCellarsProposal)
    - [TickRange](#allocation.v1.TickRange)
  
- [allocation/v1/tx.proto](#allocation/v1/tx.proto)
    - [MsgAllocationCommit](#allocation.v1.MsgAllocationCommit)
    - [MsgAllocationCommitResponse](#allocation.v1.MsgAllocationCommitResponse)
    - [MsgAllocationPrecommit](#allocation.v1.MsgAllocationPrecommit)
    - [MsgAllocationPrecommitResponse](#allocation.v1.MsgAllocationPrecommitResponse)
  
    - [Msg](#allocation.v1.Msg)
  
- [allocation/v1/genesis.proto](#allocation/v1/genesis.proto)
    - [GenesisState](#allocation.v1.GenesisState)
    - [Params](#allocation.v1.Params)
  
- [allocation/v1/query.proto](#allocation/v1/query.proto)
    - [QueryAllocationCommitRequest](#allocation.v1.QueryAllocationCommitRequest)
    - [QueryAllocationCommitResponse](#allocation.v1.QueryAllocationCommitResponse)
    - [QueryAllocationPrecommitRequest](#allocation.v1.QueryAllocationPrecommitRequest)
    - [QueryAllocationPrecommitResponse](#allocation.v1.QueryAllocationPrecommitResponse)
    - [QueryCommitPeriodRequest](#allocation.v1.QueryCommitPeriodRequest)
    - [QueryCommitPeriodResponse](#allocation.v1.QueryCommitPeriodResponse)
    - [QueryParamsRequest](#allocation.v1.QueryParamsRequest)
    - [QueryParamsResponse](#allocation.v1.QueryParamsResponse)
  
    - [Query](#allocation.v1.Query)
  
- [Scalar Value Types](#scalar-value-types)



<a name="allocation/v1/allocation.proto"></a>
<p align="right"><a href="#top">Top</a></p>

## allocation/v1/allocation.proto



<a name="allocation.v1.AddManagedCellarsProposal"></a>

### AddManagedCellarsProposal



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `title` | [string](#string) |  |  |
| `description` | [string](#string) |  |  |
| `cellar_ids` | [string](#string) | repeated |  |






<a name="allocation.v1.Allocation"></a>

### Allocation
Allocation is the commit for all allocations for a cellar by a validator


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `cellar` | [Cellar](#allocation.v1.Cellar) |  |  |
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
| `id` | [string](#string) |  |  |
| `tick_ranges` | [TickRange](#allocation.v1.TickRange) | repeated |  |






<a name="allocation.v1.CellarUpdate"></a>

### CellarUpdate



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `invalidation_nonce` | [uint64](#uint64) |  |  |
| `cellar` | [Cellar](#allocation.v1.Cellar) |  |  |






<a name="allocation.v1.RemoveManagedCellarsProposal"></a>

### RemoveManagedCellarsProposal



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `title` | [string](#string) |  |  |
| `description` | [string](#string) |  |  |
| `cellar_ids` | [string](#string) | repeated |  |






<a name="allocation.v1.TickRange"></a>

### TickRange



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `upper` | [int32](#int32) |  |  |
| `lower` | [int32](#int32) |  |  |
| `weight` | [uint32](#uint32) |  |  |





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





 <!-- end messages -->

 <!-- end enums -->

 <!-- end HasExtensions -->


<a name="allocation.v1.Msg"></a>

### Msg
MsgService defines the msgs that the oracle module handles.

| Method Name | Request Type | Response Type | Description | HTTP Verb | Endpoint |
| ----------- | ------------ | ------------- | ------------| ------- | -------- |
| `AllocationPrecommit` | [MsgAllocationPrecommit](#allocation.v1.MsgAllocationPrecommit) | [MsgAllocationPrecommitResponse](#allocation.v1.MsgAllocationPrecommitResponse) | AllocationPrecommit defines a message that commits a hash of allocation data feed before the data is actually submitted. | |
| `AllocationCommit` | [MsgAllocationCommit](#allocation.v1.MsgAllocationCommit) | [MsgAllocationCommitResponse](#allocation.v1.MsgAllocationCommitResponse) | AllocationCommit defines a message to submit the actual allocation data that was committed by the feeder through the pre-commit. | |

 <!-- end services -->



<a name="allocation/v1/genesis.proto"></a>
<p align="right"><a href="#top">Top</a></p>

## allocation/v1/genesis.proto



<a name="allocation.v1.GenesisState"></a>

### GenesisState
GenesisState - all allocation state that must be provided at genesis


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `params` | [Params](#allocation.v1.Params) |  |  |
| `cellars` | [Cellar](#allocation.v1.Cellar) | repeated |  |






<a name="allocation.v1.Params"></a>

### Params
Params allocation parameters


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `vote_period` | [int64](#int64) |  | VotePeriod defines the number of blocks to wait for votes before attempting to tally |
| `vote_threshold` | [string](#string) |  | VoteThreshold defines the percentage of bonded stake required to vote each period |





 <!-- end messages -->

 <!-- end enums -->

 <!-- end HasExtensions -->

 <!-- end services -->



<a name="allocation/v1/query.proto"></a>
<p align="right"><a href="#top">Top</a></p>

## allocation/v1/query.proto



<a name="allocation.v1.QueryAllocationCommitRequest"></a>

### QueryAllocationCommitRequest
QueryAllocationCommitRequest is the request type for the Query/QueryallocationDataVote gRPC method.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `validator` | [string](#string) |  | validator operator address |
| `cellar` | [string](#string) |  | cellar contract address |






<a name="allocation.v1.QueryAllocationCommitResponse"></a>

### QueryAllocationCommitResponse
QueryAllocationCommitResponse is the response type for the Query/QueryallocationDataVote gRPC method.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `commit` | [Allocation](#allocation.v1.Allocation) |  | vote containing the allocation feed submitted within the latest voting period |






<a name="allocation.v1.QueryAllocationPrecommitRequest"></a>

### QueryAllocationPrecommitRequest
QueryAllocationPrecommitRequest is the request type for the Query/AllocationPrecommit gRPC method.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `validator` | [string](#string) |  | validator operator address |
| `cellar` | [string](#string) |  | cellar contract address |






<a name="allocation.v1.QueryAllocationPrecommitResponse"></a>

### QueryAllocationPrecommitResponse
QueryAllocationPrecommitResponse is the response type for the Query/QueryallocationDataPrevote gRPC method.


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
| `vote_period_start` | [int64](#int64) |  | latest vote period start block height |
| `vote_period_end` | [int64](#int64) |  | block height at which the current voting period ends |






<a name="allocation.v1.QueryParamsRequest"></a>

### QueryParamsRequest
QueryParamsRequest is the request type for the Query/Params gRPC method.






<a name="allocation.v1.QueryParamsResponse"></a>

### QueryParamsResponse
QueryParamsRequest is the response type for the Query/Params gRPC method.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `params` | [Params](#allocation.v1.Params) |  | allocation parameters |





 <!-- end messages -->

 <!-- end enums -->

 <!-- end HasExtensions -->


<a name="allocation.v1.Query"></a>

### Query
Query defines the gRPC querier service for the allocation module.

| Method Name | Request Type | Response Type | Description | HTTP Verb | Endpoint |
| ----------- | ------------ | ------------- | ------------| ------- | -------- |
| `QueryParams` | [QueryParamsRequest](#allocation.v1.QueryParamsRequest) | [QueryParamsResponse](#allocation.v1.QueryParamsResponse) | QueryParams queries the allocation module parameters. | GET|/sommelier/allocation/v1/params|
| `QueryAllocationPrecommit` | [QueryAllocationPrecommitRequest](#allocation.v1.QueryAllocationPrecommitRequest) | [QueryAllocationPrecommitResponse](#allocation.v1.QueryAllocationPrecommitResponse) | QueryAllocationPrecommit queries the validator prevote in the current voting period | GET|/sommelier/allocation/v1/precommits/{validator}|
| `QueryAllocationCommit` | [QueryAllocationCommitRequest](#allocation.v1.QueryAllocationCommitRequest) | [QueryAllocationCommitResponse](#allocation.v1.QueryAllocationCommitResponse) | QueryAllocationCommit queries the validator vote in the current voting period | GET|/sommelier/allocation/v1/commits/{validator}|
| `QueryCommitPeriod` | [QueryCommitPeriodRequest](#allocation.v1.QueryCommitPeriodRequest) | [QueryCommitPeriodResponse](#allocation.v1.QueryCommitPeriodResponse) | QueryVotePeriod queries the heights for the current voting period (current, start and end) | GET|/sommelier/allocation/v1/commit_period|

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

