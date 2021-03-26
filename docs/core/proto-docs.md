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
  
- [nft/v1/nft.proto](#nft/v1/nft.proto)
    - [BaseNFT](#nft.v1.BaseNFT)
    - [Collection](#nft.v1.Collection)
    - [Denom](#nft.v1.Denom)
    - [IDCollection](#nft.v1.IDCollection)
    - [Owner](#nft.v1.Owner)
  
- [nft/v1/genesis.proto](#nft/v1/genesis.proto)
    - [GenesisState](#nft.v1.GenesisState)
  
- [nft/v1/query.proto](#nft/v1/query.proto)
    - [QueryCollectionRequest](#nft.v1.QueryCollectionRequest)
    - [QueryCollectionResponse](#nft.v1.QueryCollectionResponse)
    - [QueryDenomRequest](#nft.v1.QueryDenomRequest)
    - [QueryDenomResponse](#nft.v1.QueryDenomResponse)
    - [QueryDenomsRequest](#nft.v1.QueryDenomsRequest)
    - [QueryDenomsResponse](#nft.v1.QueryDenomsResponse)
    - [QueryNFTRequest](#nft.v1.QueryNFTRequest)
    - [QueryNFTResponse](#nft.v1.QueryNFTResponse)
    - [QueryOwnerRequest](#nft.v1.QueryOwnerRequest)
    - [QueryOwnerResponse](#nft.v1.QueryOwnerResponse)
    - [QuerySupplyRequest](#nft.v1.QuerySupplyRequest)
    - [QuerySupplyResponse](#nft.v1.QuerySupplyResponse)
  
    - [Query](#nft.v1.Query)
  
- [nft/v1/tx.proto](#nft/v1/tx.proto)
    - [MsgBurnNFT](#nft.v1.MsgBurnNFT)
    - [MsgBurnNFTResponse](#nft.v1.MsgBurnNFTResponse)
    - [MsgEditNFT](#nft.v1.MsgEditNFT)
    - [MsgEditNFTResponse](#nft.v1.MsgEditNFTResponse)
    - [MsgIssueDenom](#nft.v1.MsgIssueDenom)
    - [MsgIssueDenomResponse](#nft.v1.MsgIssueDenomResponse)
    - [MsgMintNFT](#nft.v1.MsgMintNFT)
    - [MsgMintNFTResponse](#nft.v1.MsgMintNFTResponse)
    - [MsgTransferNFT](#nft.v1.MsgTransferNFT)
    - [MsgTransferNFTResponse](#nft.v1.MsgTransferNFTResponse)
  
    - [Msg](#nft.v1.Msg)
  
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
    - [QueryLatestPeriodAggregateDataRequest](#oracle.v1.QueryLatestPeriodAggregateDataRequest)
    - [QueryLatestPeriodAggregateDataResponse](#oracle.v1.QueryLatestPeriodAggregateDataResponse)
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



<a name="nft/v1/nft.proto"></a>
<p align="right"><a href="#top">Top</a></p>

## nft/v1/nft.proto



<a name="nft.v1.BaseNFT"></a>

### BaseNFT
BaseNFT defines a non-fungible token


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `id` | [string](#string) |  |  |
| `name` | [string](#string) |  |  |
| `uri` | [string](#string) |  |  |
| `data` | [string](#string) |  |  |
| `owner` | [string](#string) |  |  |






<a name="nft.v1.Collection"></a>

### Collection
Collection defines a type of collection


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `denom` | [Denom](#nft.v1.Denom) |  |  |
| `nfts` | [BaseNFT](#nft.v1.BaseNFT) | repeated |  |






<a name="nft.v1.Denom"></a>

### Denom
Denom defines a type of NFT


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `id` | [string](#string) |  |  |
| `name` | [string](#string) |  |  |
| `schema` | [string](#string) |  |  |
| `creator` | [string](#string) |  |  |






<a name="nft.v1.IDCollection"></a>

### IDCollection
IDCollection defines a type of collection with specified ID


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `denom_id` | [string](#string) |  |  |
| `token_ids` | [string](#string) | repeated |  |






<a name="nft.v1.Owner"></a>

### Owner
Owner defines a type of owner


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `address` | [string](#string) |  |  |
| `id_collections` | [IDCollection](#nft.v1.IDCollection) | repeated |  |





 <!-- end messages -->

 <!-- end enums -->

 <!-- end HasExtensions -->

 <!-- end services -->



<a name="nft/v1/genesis.proto"></a>
<p align="right"><a href="#top">Top</a></p>

## nft/v1/genesis.proto



<a name="nft.v1.GenesisState"></a>

### GenesisState
GenesisState defines the NFT module's genesis state


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `collections` | [Collection](#nft.v1.Collection) | repeated |  |





 <!-- end messages -->

 <!-- end enums -->

 <!-- end HasExtensions -->

 <!-- end services -->



<a name="nft/v1/query.proto"></a>
<p align="right"><a href="#top">Top</a></p>

## nft/v1/query.proto



<a name="nft.v1.QueryCollectionRequest"></a>

### QueryCollectionRequest
QueryCollectionRequest is the request type for the Query/Collection RPC method


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `denom_id` | [string](#string) |  |  |
| `pagination` | [cosmos.base.query.v1beta1.PageRequest](#cosmos.base.query.v1beta1.PageRequest) |  | pagination defines an optional pagination for the request. |






<a name="nft.v1.QueryCollectionResponse"></a>

### QueryCollectionResponse
QueryCollectionResponse is the response type for the Query/Collection RPC method


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `collection` | [Collection](#nft.v1.Collection) |  |  |
| `pagination` | [cosmos.base.query.v1beta1.PageResponse](#cosmos.base.query.v1beta1.PageResponse) |  |  |






<a name="nft.v1.QueryDenomRequest"></a>

### QueryDenomRequest
QueryDenomRequest is the request type for the Query/Denom RPC method


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `denom_id` | [string](#string) |  |  |






<a name="nft.v1.QueryDenomResponse"></a>

### QueryDenomResponse
QueryDenomResponse is the response type for the Query/Denom RPC method


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `denom` | [Denom](#nft.v1.Denom) |  |  |






<a name="nft.v1.QueryDenomsRequest"></a>

### QueryDenomsRequest
QueryDenomsRequest is the request type for the Query/Denoms RPC method


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `pagination` | [cosmos.base.query.v1beta1.PageRequest](#cosmos.base.query.v1beta1.PageRequest) |  | pagination defines an optional pagination for the request. |






<a name="nft.v1.QueryDenomsResponse"></a>

### QueryDenomsResponse
QueryDenomsResponse is the response type for the Query/Denoms RPC method


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `denoms` | [Denom](#nft.v1.Denom) | repeated |  |
| `pagination` | [cosmos.base.query.v1beta1.PageResponse](#cosmos.base.query.v1beta1.PageResponse) |  |  |






<a name="nft.v1.QueryNFTRequest"></a>

### QueryNFTRequest
QueryNFTRequest is the request type for the Query/NFT RPC method


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `denom_id` | [string](#string) |  |  |
| `token_id` | [string](#string) |  |  |






<a name="nft.v1.QueryNFTResponse"></a>

### QueryNFTResponse
QueryNFTResponse is the response type for the Query/NFT RPC method


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `nft` | [BaseNFT](#nft.v1.BaseNFT) |  |  |






<a name="nft.v1.QueryOwnerRequest"></a>

### QueryOwnerRequest
QueryOwnerRequest is the request type for the Query/Owner RPC method


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `denom_id` | [string](#string) |  |  |
| `owner` | [string](#string) |  |  |
| `pagination` | [cosmos.base.query.v1beta1.PageRequest](#cosmos.base.query.v1beta1.PageRequest) |  | pagination defines an optional pagination for the request. |






<a name="nft.v1.QueryOwnerResponse"></a>

### QueryOwnerResponse
QueryOwnerResponse is the response type for the Query/Owner RPC method


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `owner` | [Owner](#nft.v1.Owner) |  |  |
| `pagination` | [cosmos.base.query.v1beta1.PageResponse](#cosmos.base.query.v1beta1.PageResponse) |  |  |






<a name="nft.v1.QuerySupplyRequest"></a>

### QuerySupplyRequest
QuerySupplyRequest is the request type for the Query/HTLC RPC method


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `denom_id` | [string](#string) |  |  |
| `owner` | [string](#string) |  |  |






<a name="nft.v1.QuerySupplyResponse"></a>

### QuerySupplyResponse
QuerySupplyResponse is the response type for the Query/Supply RPC method


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `amount` | [uint64](#uint64) |  |  |





 <!-- end messages -->

 <!-- end enums -->

 <!-- end HasExtensions -->


<a name="nft.v1.Query"></a>

### Query
Query defines the gRPC querier service for NFT module

| Method Name | Request Type | Response Type | Description | HTTP Verb | Endpoint |
| ----------- | ------------ | ------------- | ------------| ------- | -------- |
| `Supply` | [QuerySupplyRequest](#nft.v1.QuerySupplyRequest) | [QuerySupplyResponse](#nft.v1.QuerySupplyResponse) | Supply queries the total supply of a given denom or owner | GET|/sommelier/v1/nft/collections/{denom_id}/supply|
| `Owner` | [QueryOwnerRequest](#nft.v1.QueryOwnerRequest) | [QueryOwnerResponse](#nft.v1.QueryOwnerResponse) | Owner queries the NFTs of the specified owner | GET|/sommelier/v1/nft/nfts|
| `Collection` | [QueryCollectionRequest](#nft.v1.QueryCollectionRequest) | [QueryCollectionResponse](#nft.v1.QueryCollectionResponse) | Collection queries the NFTs of the specified denom | GET|/sommelier/v1/nft/collections/{denom_id}|
| `Denom` | [QueryDenomRequest](#nft.v1.QueryDenomRequest) | [QueryDenomResponse](#nft.v1.QueryDenomResponse) | Denom queries the definition of a given denom | GET|/sommelier/v1/nft/denoms/{denom_id}|
| `Denoms` | [QueryDenomsRequest](#nft.v1.QueryDenomsRequest) | [QueryDenomsResponse](#nft.v1.QueryDenomsResponse) | Denoms queries all the denoms | GET|/sommelier/v1/nft/denoms|
| `NFT` | [QueryNFTRequest](#nft.v1.QueryNFTRequest) | [QueryNFTResponse](#nft.v1.QueryNFTResponse) | NFT queries the NFT for the given denom and token ID | GET|/sommelier/v1/nft/nfts/{denom_id}/{token_id}|

 <!-- end services -->



<a name="nft/v1/tx.proto"></a>
<p align="right"><a href="#top">Top</a></p>

## nft/v1/tx.proto



<a name="nft.v1.MsgBurnNFT"></a>

### MsgBurnNFT
MsgBurnNFT defines an SDK message for burning a NFT.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `id` | [string](#string) |  |  |
| `denom_id` | [string](#string) |  |  |
| `sender` | [string](#string) |  |  |






<a name="nft.v1.MsgBurnNFTResponse"></a>

### MsgBurnNFTResponse
MsgBurnNFTResponse defines the Msg/BurnNFT response type.






<a name="nft.v1.MsgEditNFT"></a>

### MsgEditNFT
MsgEditNFT defines an SDK message for editing a nft.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `id` | [string](#string) |  |  |
| `denom_id` | [string](#string) |  |  |
| `name` | [string](#string) |  |  |
| `uri` | [string](#string) |  |  |
| `data` | [string](#string) |  |  |
| `sender` | [string](#string) |  |  |






<a name="nft.v1.MsgEditNFTResponse"></a>

### MsgEditNFTResponse
MsgEditNFTResponse defines the Msg/EditNFT response type.






<a name="nft.v1.MsgIssueDenom"></a>

### MsgIssueDenom
MsgIssueDenom defines an SDK message for creating a new denom.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `id` | [string](#string) |  |  |
| `name` | [string](#string) |  |  |
| `schema` | [string](#string) |  |  |
| `sender` | [string](#string) |  |  |






<a name="nft.v1.MsgIssueDenomResponse"></a>

### MsgIssueDenomResponse
MsgIssueDenomResponse defines the Msg/IssueDenom response type.






<a name="nft.v1.MsgMintNFT"></a>

### MsgMintNFT
MsgMintNFT defines an SDK message for creating a new NFT.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `id` | [string](#string) |  |  |
| `denom_id` | [string](#string) |  |  |
| `name` | [string](#string) |  |  |
| `uri` | [string](#string) |  |  |
| `data` | [string](#string) |  |  |
| `sender` | [string](#string) |  |  |
| `recipient` | [string](#string) |  |  |






<a name="nft.v1.MsgMintNFTResponse"></a>

### MsgMintNFTResponse
MsgMintNFTResponse defines the Msg/MintNFT response type.






<a name="nft.v1.MsgTransferNFT"></a>

### MsgTransferNFT
MsgTransferNFT defines an SDK message for transferring an NFT to recipient.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `id` | [string](#string) |  |  |
| `denom_id` | [string](#string) |  |  |
| `name` | [string](#string) |  |  |
| `uri` | [string](#string) |  |  |
| `data` | [string](#string) |  |  |
| `sender` | [string](#string) |  |  |
| `recipient` | [string](#string) |  |  |






<a name="nft.v1.MsgTransferNFTResponse"></a>

### MsgTransferNFTResponse
MsgTransferNFTResponse defines the Msg/TransferNFT response type.





 <!-- end messages -->

 <!-- end enums -->

 <!-- end HasExtensions -->


<a name="nft.v1.Msg"></a>

### Msg
Msg defines the htlc Msg service.

| Method Name | Request Type | Response Type | Description | HTTP Verb | Endpoint |
| ----------- | ------------ | ------------- | ------------| ------- | -------- |
| `IssueDenom` | [MsgIssueDenom](#nft.v1.MsgIssueDenom) | [MsgIssueDenomResponse](#nft.v1.MsgIssueDenomResponse) | IssueDenom defines a method for issue a denom. | |
| `MintNFT` | [MsgMintNFT](#nft.v1.MsgMintNFT) | [MsgMintNFTResponse](#nft.v1.MsgMintNFTResponse) | MintNFT defines a method for mint a new nft | |
| `EditNFT` | [MsgEditNFT](#nft.v1.MsgEditNFT) | [MsgEditNFTResponse](#nft.v1.MsgEditNFTResponse) | RefundHTLC defines a method for editing a nft. | |
| `TransferNFT` | [MsgTransferNFT](#nft.v1.MsgTransferNFT) | [MsgTransferNFTResponse](#nft.v1.MsgTransferNFTResponse) | TransferNFT defines a method for transferring a nft. | |
| `BurnNFT` | [MsgBurnNFT](#nft.v1.MsgBurnNFT) | [MsgBurnNFTResponse](#nft.v1.MsgBurnNFTResponse) | BurnNFT defines a method for burning a nft. | |

 <!-- end services -->



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
MsgDelegateFeedConsent defines sdk.Msg for delegating oracle voting rights from a validator
to another address, must be signed by an active validator


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `delegate` | [string](#string) |  | delegate account address |
| `validator` | [string](#string) |  | validator operator address |






<a name="oracle.v1.MsgDelegateFeedConsentResponse"></a>

### MsgDelegateFeedConsentResponse
MsgDelegateFeedConsentResponse is the response type for the Msg/DelegateFeedConsent gRPC method.






<a name="oracle.v1.MsgOracleDataPrevote"></a>

### MsgOracleDataPrevote
MsgOracleDataPrevote - sdk.Msg for prevoting on an array of oracle data types.
The purpose of the prevote is to hide vote for data with hashes formatted as hex string: 
SHA256("{salt}:{data_cannonical_json}:{voter}")


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `prevote` | [OraclePrevote](#oracle.v1.OraclePrevote) |  | prevote containing the hash of the oracle feed vote contents |
| `signer` | [string](#string) |  | signer (i.e feeder) account address |






<a name="oracle.v1.MsgOracleDataPrevoteResponse"></a>

### MsgOracleDataPrevoteResponse
MsgOracleDataPrevoteResponse is the response type for the Msg/OracleDataPrevote gRPC method.






<a name="oracle.v1.MsgOracleDataVote"></a>

### MsgOracleDataVote
MsgOracleDataVote - sdk.Msg for submitting arbitrary oracle data that has been prevoted on


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `vote` | [OracleVote](#oracle.v1.OracleVote) |  | vote containing the oracle data feed |
| `signer` | [string](#string) |  | signer (i.e feeder) account address |






<a name="oracle.v1.MsgOracleDataVoteResponse"></a>

### MsgOracleDataVoteResponse
MsgOracleDataVoteResponse is the response type for the Msg/OracleDataVote gRPC method.





 <!-- end messages -->

 <!-- end enums -->

 <!-- end HasExtensions -->


<a name="oracle.v1.Msg"></a>

### Msg
MsgService defines the msgs that the oracle module handles.

| Method Name | Request Type | Response Type | Description | HTTP Verb | Endpoint |
| ----------- | ------------ | ------------- | ------------| ------- | -------- |
| `DelegateFeedConsent` | [MsgDelegateFeedConsent](#oracle.v1.MsgDelegateFeedConsent) | [MsgDelegateFeedConsentResponse](#oracle.v1.MsgDelegateFeedConsentResponse) | DelegateFeedConsent defines a message that delegates the oracle feeding to an account address. | |
| `OracleDataPrevote` | [MsgOracleDataPrevote](#oracle.v1.MsgOracleDataPrevote) | [MsgOracleDataPrevoteResponse](#oracle.v1.MsgOracleDataPrevoteResponse) | OracleDataPrevote defines a message that commits a hash of a oracle data feed before the data is actually submitted. | |
| `OracleDataVote` | [MsgOracleDataVote](#oracle.v1.MsgOracleDataVote) | [MsgOracleDataVoteResponse](#oracle.v1.MsgOracleDataVoteResponse) | OracleDataVote defines a message to submit the actual oracle data that was committed by the feeder through the prevote. | |

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
| `validator` | [string](#string) |  | validator operator address |
| `misses` | [int64](#int64) |  | number of misses |






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
QueryAggregateDataRequest is the request type for the Query/AggregateData gRPC method.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `type` | [string](#string) |  | oracle data type |
| `id` | [string](#string) |  | oracle data identifier |






<a name="oracle.v1.QueryAggregateDataResponse"></a>

### QueryAggregateDataResponse
QueryAggregateDataRequest is the response type for the Query/AggregateData gRPC method.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `oracle_data` | [UniswapPair](#oracle.v1.UniswapPair) |  | oracle data associated with the id and type from the request |
| `height` | [int64](#int64) |  | height at which the aggregated oracle data was stored |






<a name="oracle.v1.QueryDelegateAddressRequest"></a>

### QueryDelegateAddressRequest
QueryDelegateAddressRequest is the request type for the Query/QueryDelegateAddress gRPC method.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `validator` | [string](#string) |  | validator operator address |






<a name="oracle.v1.QueryDelegateAddressResponse"></a>

### QueryDelegateAddressResponse
QueryDelegateAddressResponse is the response type for the Query/QueryDelegateAddress gRPC method.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `delegate` | [string](#string) |  | delegate account address |






<a name="oracle.v1.QueryLatestPeriodAggregateDataRequest"></a>

### QueryLatestPeriodAggregateDataRequest
QueryLatestPeriodAggregateDataRequest is the request type for the Query/QueryLatestPeriodAggregateData gRPC method.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `pagination` | [cosmos.base.query.v1beta1.PageRequest](#cosmos.base.query.v1beta1.PageRequest) |  | pagination defines an optional pagination for the request. |






<a name="oracle.v1.QueryLatestPeriodAggregateDataResponse"></a>

### QueryLatestPeriodAggregateDataResponse
QueryLatestPeriodAggregateDataResponse is the response type for the Query/QueryLatestPeriodAggregateData gRPC method.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `oracle_data` | [UniswapPair](#oracle.v1.UniswapPair) | repeated | oracle data associated with the |
| `height` | [int64](#int64) |  | height at which the aggregated oracle data was stored |
| `pagination` | [cosmos.base.query.v1beta1.PageResponse](#cosmos.base.query.v1beta1.PageResponse) |  | pagination defines the pagination in the response. |






<a name="oracle.v1.QueryMissCounterRequest"></a>

### QueryMissCounterRequest
QueryMissCounterRequest is the request type for the Query/MissCounter gRPC method.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `validator` | [string](#string) |  | validator operator address |






<a name="oracle.v1.QueryMissCounterResponse"></a>

### QueryMissCounterResponse
QueryMissCounterResponse is the response type for the Query/MissCounter gRPC method.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `miss_counter` | [int64](#int64) |  | number of oracle feed votes missed since the last counter reset |






<a name="oracle.v1.QueryOracleDataPrevoteRequest"></a>

### QueryOracleDataPrevoteRequest
QueryOracleDataPrevoteRequest is the request type for the Query/QueryOracleDataPrevote gRPC method.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `validator` | [string](#string) |  | validator operator address |






<a name="oracle.v1.QueryOracleDataPrevoteResponse"></a>

### QueryOracleDataPrevoteResponse
QueryOracleDataPrevoteResponse is the response type for the Query/QueryOracleDataPrevote gRPC method.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `prevote` | [OraclePrevote](#oracle.v1.OraclePrevote) |  | prevote submitted within the latest voting period |






<a name="oracle.v1.QueryOracleDataVoteRequest"></a>

### QueryOracleDataVoteRequest
QueryOracleDataVoteRequest is the request type for the Query/QueryOracleDataVote gRPC method.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `validator` | [string](#string) |  | validator operator address |






<a name="oracle.v1.QueryOracleDataVoteResponse"></a>

### QueryOracleDataVoteResponse
QueryOracleDataVoteResponse is the response type for the Query/QueryOracleDataVote gRPC method.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `vote` | [OracleVote](#oracle.v1.OracleVote) |  | vote containing the oracle feed submitted within the latest voting period |






<a name="oracle.v1.QueryParamsRequest"></a>

### QueryParamsRequest
QueryParamsRequest is the request type for the Query/Params gRPC method.






<a name="oracle.v1.QueryParamsResponse"></a>

### QueryParamsResponse
QueryParamsRequest is the response type for the Query/Params gRPC method.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `params` | [Params](#oracle.v1.Params) |  | oracle parameters |






<a name="oracle.v1.QueryValidatorAddressRequest"></a>

### QueryValidatorAddressRequest
QueryValidatorAddressRequest is the request type for the Query/Params gRPC method.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `delegate` | [string](#string) |  | delegate account address |






<a name="oracle.v1.QueryValidatorAddressResponse"></a>

### QueryValidatorAddressResponse
QueryValidatorAddressResponse is the response type for the Query/Params gRPC method.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `validator` | [string](#string) |  | validator operator address |






<a name="oracle.v1.QueryVotePeriodRequest"></a>

### QueryVotePeriodRequest
QueryVotePeriodRequest is the request type for the Query/VotePeriod gRPC method.






<a name="oracle.v1.QueryVotePeriodResponse"></a>

### QueryVotePeriodResponse
QueryVotePeriodResponse is the response type for the Query/VotePeriod gRPC method.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `current_height` | [int64](#int64) |  | block height at which the query was processed |
| `vote_period_start` | [int64](#int64) |  | latest vote period start block height |
| `vote_period_end` | [int64](#int64) |  | block height at which the current voting period ends |





 <!-- end messages -->

 <!-- end enums -->

 <!-- end HasExtensions -->


<a name="oracle.v1.Query"></a>

### Query
Query defines the gRPC querier service for the oracle module.

| Method Name | Request Type | Response Type | Description | HTTP Verb | Endpoint |
| ----------- | ------------ | ------------- | ------------| ------- | -------- |
| `QueryParams` | [QueryParamsRequest](#oracle.v1.QueryParamsRequest) | [QueryParamsResponse](#oracle.v1.QueryParamsResponse) | QueryParams queries the oracle module parameters. | GET|/sommelier/oracle/v1/params|
| `QueryDelegateAddress` | [QueryDelegateAddressRequest](#oracle.v1.QueryDelegateAddressRequest) | [QueryDelegateAddressResponse](#oracle.v1.QueryDelegateAddressResponse) | QueryDelegateAddress queries the delegate account address of a validator | GET|/sommelier/oracle/v1/delegates/{validator}|
| `QueryValidatorAddress` | [QueryValidatorAddressRequest](#oracle.v1.QueryValidatorAddressRequest) | [QueryValidatorAddressResponse](#oracle.v1.QueryValidatorAddressResponse) | QueryValidatorAddress returns the validator address of a given delegate | GET|/sommelier/oracle/v1/validators/{delegate}|
| `QueryOracleDataPrevote` | [QueryOracleDataPrevoteRequest](#oracle.v1.QueryOracleDataPrevoteRequest) | [QueryOracleDataPrevoteResponse](#oracle.v1.QueryOracleDataPrevoteResponse) | QueryOracleDataPrevote queries the validator prevote in the current voting period | GET|/sommelier/oracle/v1/prevotes/{validator}|
| `QueryOracleDataVote` | [QueryOracleDataVoteRequest](#oracle.v1.QueryOracleDataVoteRequest) | [QueryOracleDataVoteResponse](#oracle.v1.QueryOracleDataVoteResponse) | QueryOracleDataVote queries the validator vote in the current voting period | GET|/sommelier/oracle/v1/votes/{validator}|
| `QueryVotePeriod` | [QueryVotePeriodRequest](#oracle.v1.QueryVotePeriodRequest) | [QueryVotePeriodResponse](#oracle.v1.QueryVotePeriodResponse) | QueryVotePeriod queries the heights for the current voting period (current, start and end) | GET|/sommelier/oracle/v1/vote_period|
| `QueryMissCounter` | [QueryMissCounterRequest](#oracle.v1.QueryMissCounterRequest) | [QueryMissCounterResponse](#oracle.v1.QueryMissCounterResponse) | QueryMissCounter queries the missed number of oracle data feed periods | GET|/sommelier/oracle/v1/miss_counters/{validator}|
| `QueryAggregateData` | [QueryAggregateDataRequest](#oracle.v1.QueryAggregateDataRequest) | [QueryAggregateDataResponse](#oracle.v1.QueryAggregateDataResponse) | QueryAggregateData returns the latest aggregated data value for a given type and identifioer | GET|/sommelier/oracle/v1/aggregate_data/{id}/{type}|
| `QueryLatestPeriodAggregateData` | [QueryLatestPeriodAggregateDataRequest](#oracle.v1.QueryLatestPeriodAggregateDataRequest) | [QueryLatestPeriodAggregateDataResponse](#oracle.v1.QueryLatestPeriodAggregateDataResponse) | QueryLatestPeriodAggregateData returns the aggregated data for a given pair an identifioer | GET|/sommelier/oracle/v1/aggregate_data|

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

