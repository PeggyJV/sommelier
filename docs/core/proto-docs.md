<!-- This file is auto-generated. Please do not modify it yourself. -->
# Protobuf Documentation
<a name="top"></a>

## Table of Contents

- [uniswap-oracle/v1/oracle.proto](#uniswap-oracle/v1/oracle.proto)
    - [Claim](#uniswap_oracle.v1.Claim)
    - [Pair](#uniswap_oracle.v1.Pair)
    - [Params](#uniswap_oracle.v1.Params)
    - [Token](#uniswap_oracle.v1.Token)
    - [UniswapPrevote](#uniswap_oracle.v1.UniswapPrevote)
    - [UniswapVote](#uniswap_oracle.v1.UniswapVote)
  
- [uniswap-oracle/v1/genesis.proto](#uniswap-oracle/v1/genesis.proto)
    - [GenesisState](#uniswap_oracle.v1.GenesisState)
    - [OracleDelegation](#uniswap_oracle.v1.OracleDelegation)
    - [ValidatorMissCounter](#uniswap_oracle.v1.ValidatorMissCounter)
  
- [uniswap-oracle/v1/query.proto](#uniswap-oracle/v1/query.proto)
    - [QueryUniswapDataRequest](#uniswap_oracle.v1.QueryUniswapDataRequest)
    - [QueryUniswapDataResponse](#uniswap_oracle.v1.QueryUniswapDataResponse)
  
    - [Query](#uniswap_oracle.v1.Query)
  
- [uniswap-oracle/v1/tx.proto](#uniswap-oracle/v1/tx.proto)
    - [MsgDelegateFeedConsent](#uniswap_oracle.v1.MsgDelegateFeedConsent)
    - [MsgDelegateFeedConsentResponse](#uniswap_oracle.v1.MsgDelegateFeedConsentResponse)
    - [MsgUniswapDataPrevote](#uniswap_oracle.v1.MsgUniswapDataPrevote)
    - [MsgUniswapDataPrevoteResponse](#uniswap_oracle.v1.MsgUniswapDataPrevoteResponse)
    - [MsgUniswapDataVote](#uniswap_oracle.v1.MsgUniswapDataVote)
    - [MsgUniswapDataVoteResponse](#uniswap_oracle.v1.MsgUniswapDataVoteResponse)
  
    - [Msg](#uniswap_oracle.v1.Msg)
  
- [Scalar Value Types](#scalar-value-types)



<a name="uniswap-oracle/v1/oracle.proto"></a>
<p align="right"><a href="#top">Top</a></p>

## uniswap-oracle/v1/oracle.proto



<a name="uniswap_oracle.v1.Claim"></a>

### Claim
Claim is an interface that directs its rewards to an attached bank account.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `weight` | [int64](#int64) |  |  |
| `recipient` | [string](#string) |  |  |






<a name="uniswap_oracle.v1.Pair"></a>

### Pair



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `id` | [string](#string) |  |  |
| `reserve0` | [string](#string) |  |  |
| `reserve1` | [string](#string) |  |  |
| `reserve_usd` | [string](#string) |  |  |
| `token0` | [Token](#uniswap_oracle.v1.Token) |  |  |
| `token1` | [Token](#uniswap_oracle.v1.Token) |  |  |
| `token0_price` | [string](#string) |  |  |
| `token1_price` | [string](#string) |  |  |
| `total_supply` | [string](#string) |  |  |






<a name="uniswap_oracle.v1.Params"></a>

### Params
Params oracle parameters






<a name="uniswap_oracle.v1.Token"></a>

### Token
Token is the returned uniswap token representation


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `decimals` | [string](#string) |  |  |
| `id` | [string](#string) |  |  |






<a name="uniswap_oracle.v1.UniswapPrevote"></a>

### UniswapPrevote
UniswapPrevote - struct to store a validator's prevote on uniswap data


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `hash` | [bytes](#bytes) |  |  |
| `voter` | [string](#string) |  |  |
| `submit_block` | [int64](#int64) |  |  |






<a name="uniswap_oracle.v1.UniswapVote"></a>

### UniswapVote
UniswapVote - struct to store a validator's vote on uniswap data


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `voter` | [string](#string) |  |  |
| `pairs` | [Pair](#uniswap_oracle.v1.Pair) | repeated |  |





 <!-- end messages -->

 <!-- end enums -->

 <!-- end HasExtensions -->

 <!-- end services -->



<a name="uniswap-oracle/v1/genesis.proto"></a>
<p align="right"><a href="#top">Top</a></p>

## uniswap-oracle/v1/genesis.proto



<a name="uniswap_oracle.v1.GenesisState"></a>

### GenesisState
GenesisState - all oracle state that must be provided at genesis
GenesisState - all oracle state that must be provided at genesis


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `params` | [Params](#uniswap_oracle.v1.Params) |  |  |
| `feeder_delegations` | [OracleDelegation](#uniswap_oracle.v1.OracleDelegation) | repeated |  |
| `miss_counters` | [ValidatorMissCounter](#uniswap_oracle.v1.ValidatorMissCounter) | repeated |  |
| `uniswap_prevote` | [UniswapPrevote](#uniswap_oracle.v1.UniswapPrevote) | repeated |  |
| `uniswap_vote` | [UniswapVote](#uniswap_oracle.v1.UniswapVote) | repeated |  |






<a name="uniswap_oracle.v1.OracleDelegation"></a>

### OracleDelegation
OracleDelegation represents a delegator-delegate pair for an oracle delegation


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `validator_address` | [string](#string) |  | validator address which is delegating it's authority to vote |
| `delegate_address` | [string](#string) |  | account delegate address |






<a name="uniswap_oracle.v1.ValidatorMissCounter"></a>

### ValidatorMissCounter
ValidatorMissCounter represents a per validator miss counter


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `val_address` | [string](#string) |  | validator operator address |
| `missed_counter` | [int64](#int64) |  | missed oracle counter |





 <!-- end messages -->

 <!-- end enums -->

 <!-- end HasExtensions -->

 <!-- end services -->



<a name="uniswap-oracle/v1/query.proto"></a>
<p align="right"><a href="#top">Top</a></p>

## uniswap-oracle/v1/query.proto



<a name="uniswap_oracle.v1.QueryUniswapDataRequest"></a>

### QueryUniswapDataRequest







<a name="uniswap_oracle.v1.QueryUniswapDataResponse"></a>

### QueryUniswapDataResponse






 <!-- end messages -->

 <!-- end enums -->

 <!-- end HasExtensions -->


<a name="uniswap_oracle.v1.Query"></a>

### Query


| Method Name | Request Type | Response Type | Description | HTTP Verb | Endpoint |
| ----------- | ------------ | ------------- | ------------| ------- | -------- |
| `UniswapData` | [QueryUniswapDataRequest](#uniswap_oracle.v1.QueryUniswapDataRequest) | [QueryUniswapDataResponse](#uniswap_oracle.v1.QueryUniswapDataResponse) |  | GET|/uniswap-oracle/v1/uniswap-data|

 <!-- end services -->



<a name="uniswap-oracle/v1/tx.proto"></a>
<p align="right"><a href="#top">Top</a></p>

## uniswap-oracle/v1/tx.proto



<a name="uniswap_oracle.v1.MsgDelegateFeedConsent"></a>

### MsgDelegateFeedConsent
MsgDelegateFeedConsent - struct for delegating oracle voting rights to another address.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `operator` | [string](#string) |  |  |
| `delegate` | [string](#string) |  |  |






<a name="uniswap_oracle.v1.MsgDelegateFeedConsentResponse"></a>

### MsgDelegateFeedConsentResponse







<a name="uniswap_oracle.v1.MsgUniswapDataPrevote"></a>

### MsgUniswapDataPrevote
MsgUniswapDataPrevote - struct for aggregate prevoting on the ExchangeRateVote.
The purpose of aggregate prevote is to hide vote exchange rates with hash
which is formatted as hex string in SHA256("{salt}:{exchange rate}{denom},...,{exchange rate}{denom}:{voter}")


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `hash` | [bytes](#bytes) |  |  |
| `feeder` | [string](#string) |  |  |
| `validator` | [string](#string) |  |  |






<a name="uniswap_oracle.v1.MsgUniswapDataPrevoteResponse"></a>

### MsgUniswapDataPrevoteResponse







<a name="uniswap_oracle.v1.MsgUniswapDataVote"></a>

### MsgUniswapDataVote
MsgUniswapDataVote - struct for voting on the exchange rates of Luna denominated in various Terra assets.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `salt` | [string](#string) |  |  |
| `pairs` | [Pair](#uniswap_oracle.v1.Pair) | repeated |  |
| `feeder` | [string](#string) |  |  |
| `validator` | [string](#string) |  |  |






<a name="uniswap_oracle.v1.MsgUniswapDataVoteResponse"></a>

### MsgUniswapDataVoteResponse






 <!-- end messages -->

 <!-- end enums -->

 <!-- end HasExtensions -->


<a name="uniswap_oracle.v1.Msg"></a>

### Msg
MsgService defines the msgs that the oracle module handles.

| Method Name | Request Type | Response Type | Description | HTTP Verb | Endpoint |
| ----------- | ------------ | ------------- | ------------| ------- | -------- |
| `DelegateFeedConsent` | [MsgDelegateFeedConsent](#uniswap_oracle.v1.MsgDelegateFeedConsent) | [MsgDelegateFeedConsentResponse](#uniswap_oracle.v1.MsgDelegateFeedConsentResponse) |  | |
| `UniswapDataPrevote` | [MsgUniswapDataPrevote](#uniswap_oracle.v1.MsgUniswapDataPrevote) | [MsgUniswapDataPrevoteResponse](#uniswap_oracle.v1.MsgUniswapDataPrevoteResponse) |  | |
| `UniswapDataVote` | [MsgUniswapDataVote](#uniswap_oracle.v1.MsgUniswapDataVote) | [MsgUniswapDataVoteResponse](#uniswap_oracle.v1.MsgUniswapDataVoteResponse) |  | |

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

