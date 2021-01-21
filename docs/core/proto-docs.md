<!-- This file is auto-generated. Please do not modify it yourself. -->
# Protobuf Documentation
<a name="top"></a>

## Table of Contents

- [il/v1/il.proto](#il/v1/il.proto)
    - [Params](#il.v1.Params)
    - [Stoploss](#il.v1.Stoploss)
  
- [il/v1/genesis.proto](#il/v1/genesis.proto)
    - [GenesisState](#il.v1.GenesisState)
    - [StoplossPosition](#il.v1.StoplossPosition)
  
- [il/v1/query.proto](#il/v1/query.proto)
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
    - [AggregateExchangeRatePrevote](#oracle.v1.AggregateExchangeRatePrevote)
    - [AggregateExchangeRateVote](#oracle.v1.AggregateExchangeRateVote)
    - [Claim](#oracle.v1.Claim)
    - [ExchangeRatePrevote](#oracle.v1.ExchangeRatePrevote)
    - [ExchangeRateVote](#oracle.v1.ExchangeRateVote)
    - [Params](#oracle.v1.Params)
  
- [oracle/v1/genesis.proto](#oracle/v1/genesis.proto)
    - [GenesisState](#oracle.v1.GenesisState)
    - [OracleDelegation](#oracle.v1.OracleDelegation)
    - [ValidatorMissCounter](#oracle.v1.ValidatorMissCounter)
  
- [oracle/v1/query.proto](#oracle/v1/query.proto)
    - [QueryActivesRequest](#oracle.v1.QueryActivesRequest)
    - [QueryActivesResponse](#oracle.v1.QueryActivesResponse)
    - [QueryAggregatePrevoteRequest](#oracle.v1.QueryAggregatePrevoteRequest)
    - [QueryAggregatePrevoteResponse](#oracle.v1.QueryAggregatePrevoteResponse)
    - [QueryAggregateVoteRequest](#oracle.v1.QueryAggregateVoteRequest)
    - [QueryAggregateVoteResponse](#oracle.v1.QueryAggregateVoteResponse)
    - [QueryExchangeRateRequest](#oracle.v1.QueryExchangeRateRequest)
    - [QueryExchangeRateResponse](#oracle.v1.QueryExchangeRateResponse)
    - [QueryExchangeRatesRequest](#oracle.v1.QueryExchangeRatesRequest)
    - [QueryExchangeRatesResponse](#oracle.v1.QueryExchangeRatesResponse)
    - [QueryFeederDelegationRequest](#oracle.v1.QueryFeederDelegationRequest)
    - [QueryFeederDelegationResponse](#oracle.v1.QueryFeederDelegationResponse)
    - [QueryMissCounterRequest](#oracle.v1.QueryMissCounterRequest)
    - [QueryMissCounterResponse](#oracle.v1.QueryMissCounterResponse)
    - [QueryParametersRequest](#oracle.v1.QueryParametersRequest)
    - [QueryParametersResponse](#oracle.v1.QueryParametersResponse)
    - [QueryTobinTaxRequest](#oracle.v1.QueryTobinTaxRequest)
    - [QueryTobinTaxResponse](#oracle.v1.QueryTobinTaxResponse)
    - [QueryTobinTaxesRequest](#oracle.v1.QueryTobinTaxesRequest)
    - [QueryTobinTaxesResponse](#oracle.v1.QueryTobinTaxesResponse)
    - [QueryVoteTargetsRequest](#oracle.v1.QueryVoteTargetsRequest)
    - [QueryVoteTargetsResponse](#oracle.v1.QueryVoteTargetsResponse)
  
    - [Query](#oracle.v1.Query)
  
- [oracle/v1/tx.proto](#oracle/v1/tx.proto)
    - [MsgAggregateExchangeRatePrevote](#oracle.v1.MsgAggregateExchangeRatePrevote)
    - [MsgAggregateExchangeRatePrevoteResponse](#oracle.v1.MsgAggregateExchangeRatePrevoteResponse)
    - [MsgAggregateExchangeRateVote](#oracle.v1.MsgAggregateExchangeRateVote)
    - [MsgAggregateExchangeRateVoteResponse](#oracle.v1.MsgAggregateExchangeRateVoteResponse)
    - [MsgDelegateFeedConsent](#oracle.v1.MsgDelegateFeedConsent)
    - [MsgDelegateFeedConsentResponse](#oracle.v1.MsgDelegateFeedConsentResponse)
  
    - [Msg](#oracle.v1.Msg)
  
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
| `max_slippage` | [cosmos.base.v1beta1.DecCoin](#cosmos.base.v1beta1.DecCoin) |  | max slippage allowed before the stoploss is triggered |
| `reference_pair_ratio` | [cosmos.base.v1beta1.DecCoin](#cosmos.base.v1beta1.DecCoin) |  | starting token pair ration of the uniswap pool |





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
| `stoploss_positions` | [StoplossPosition](#il.v1.StoplossPosition) | repeated |  |






<a name="il.v1.StoplossPosition"></a>

### StoplossPosition
StoplossPosition represents an impermanent loss stop position for a given address and uniswap pair.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `address` | [string](#string) |  | account address |
| `stoploss` | [Stoploss](#il.v1.Stoploss) |  | account delegate address |





 <!-- end messages -->

 <!-- end enums -->

 <!-- end HasExtensions -->

 <!-- end services -->



<a name="il/v1/query.proto"></a>
<p align="right"><a href="#top">Top</a></p>

## il/v1/query.proto



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
| `stoploss_positions` | [Stoploss](#il.v1.Stoploss) | repeated |  |
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
| `stoploss` | [Stoploss](#il.v1.Stoploss) |  | token exchange rate |





 <!-- end messages -->

 <!-- end enums -->

 <!-- end HasExtensions -->


<a name="il.v1.Query"></a>

### Query


| Method Name | Request Type | Response Type | Description | HTTP Verb | Endpoint |
| ----------- | ------------ | ------------- | ------------| ------- | -------- |
| `Stoploss` | [QueryStoplossRequest](#il.v1.QueryStoplossRequest) | [QueryStoplossResponse](#il.v1.QueryStoplossResponse) |  | GET|/il/v1/stoploss_positions/{address}/{uniswap_pair}|
| `StoplossPositions` | [QueryStoplossPositionsRequest](#il.v1.QueryStoplossPositionsRequest) | [QueryStoplossPositionsResponse](#il.v1.QueryStoplossPositionsResponse) |  | GET|/il/v1/stoploss_positions/{address}|

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



<a name="oracle.v1.AggregateExchangeRatePrevote"></a>

### AggregateExchangeRatePrevote
AggregateExchangeRatePrevote - struct to store a validator's aggregate prevote on the rate of Luna in the denom asset


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `hash` | [bytes](#bytes) |  |  |
| `voter` | [string](#string) |  |  |
| `submit_block` | [int64](#int64) |  |  |






<a name="oracle.v1.AggregateExchangeRateVote"></a>

### AggregateExchangeRateVote
AggregateExchangeRateVote - struct to store a validator's aggregate vote on the rate of Luna in the denom asset


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `voter` | [string](#string) |  |  |
| `exchange_rate_tuples` | [cosmos.base.v1beta1.DecCoin](#cosmos.base.v1beta1.DecCoin) | repeated |  |






<a name="oracle.v1.Claim"></a>

### Claim
Claim is an interface that directs its rewards to an attached bank account.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `weight` | [int64](#int64) |  |  |
| `recipient` | [string](#string) |  |  |






<a name="oracle.v1.ExchangeRatePrevote"></a>

### ExchangeRatePrevote
ExchangeRatePrevote - struct to store a validator's prevote on the rate of Luna in the denom asset


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `hash` | [bytes](#bytes) |  |  |
| `denom` | [string](#string) |  |  |
| `voter` | [string](#string) |  |  |
| `submit_block` | [int64](#int64) |  |  |






<a name="oracle.v1.ExchangeRateVote"></a>

### ExchangeRateVote
ExchangeRateVote - struct to store a validator's vote on the rate of Luna in the denom asset


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `exchange_rate` | [string](#string) |  |  |
| `denom` | [string](#string) |  |  |
| `voter` | [string](#string) |  |  |






<a name="oracle.v1.Params"></a>

### Params
Params oracle parameters


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `vote_period` | [int64](#int64) |  |  |
| `vote_threshold` | [string](#string) |  |  |
| `reward_band` | [string](#string) |  |  |
| `reward_distribution_window` | [int64](#int64) |  |  |
| `whitelist` | [cosmos.base.v1beta1.DecCoin](#cosmos.base.v1beta1.DecCoin) | repeated | NOTE: The amounts here indicate the tobin tax for each currency |
| `slash_fraction` | [string](#string) |  |  |
| `slash_window` | [int64](#int64) |  |  |
| `min_valid_per_window` | [string](#string) |  |  |





 <!-- end messages -->

 <!-- end enums -->

 <!-- end HasExtensions -->

 <!-- end services -->



<a name="oracle/v1/genesis.proto"></a>
<p align="right"><a href="#top">Top</a></p>

## oracle/v1/genesis.proto



<a name="oracle.v1.GenesisState"></a>

### GenesisState
GenesisState - all oracle state that must be provided at genesis


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `params` | [Params](#oracle.v1.Params) |  |  |
| `feeder_delegations` | [OracleDelegation](#oracle.v1.OracleDelegation) | repeated |  |
| `exchange_rates` | [cosmos.base.v1beta1.DecCoin](#cosmos.base.v1beta1.DecCoin) | repeated |  |
| `miss_counters` | [ValidatorMissCounter](#oracle.v1.ValidatorMissCounter) | repeated |  |
| `aggregate_exchange_rate_prevotes` | [AggregateExchangeRatePrevote](#oracle.v1.AggregateExchangeRatePrevote) | repeated |  |
| `aggregate_exchange_rate_votes` | [AggregateExchangeRateVote](#oracle.v1.AggregateExchangeRateVote) | repeated |  |
| `tobin_taxes` | [cosmos.base.v1beta1.DecCoin](#cosmos.base.v1beta1.DecCoin) | repeated | NOTE: the amounts here indicate the tobin tax for a given USD/{denom} pair |






<a name="oracle.v1.OracleDelegation"></a>

### OracleDelegation
OracleDelegation represents a delegator-delegate pair for an oracle delegation


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `delegator_address` | [string](#string) |  | validator delegator address |
| `delegate_address` | [string](#string) |  | account delegate address |






<a name="oracle.v1.ValidatorMissCounter"></a>

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



<a name="oracle/v1/query.proto"></a>
<p align="right"><a href="#top">Top</a></p>

## oracle/v1/query.proto



<a name="oracle.v1.QueryActivesRequest"></a>

### QueryActivesRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `pagination` | [cosmos.base.query.v1beta1.PageRequest](#cosmos.base.query.v1beta1.PageRequest) |  | pagination request |






<a name="oracle.v1.QueryActivesResponse"></a>

### QueryActivesResponse



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `denoms` | [string](#string) | repeated |  |
| `pagination` | [cosmos.base.query.v1beta1.PageResponse](#cosmos.base.query.v1beta1.PageResponse) |  | pagination response |






<a name="oracle.v1.QueryAggregatePrevoteRequest"></a>

### QueryAggregatePrevoteRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `validator` | [string](#string) |  |  |






<a name="oracle.v1.QueryAggregatePrevoteResponse"></a>

### QueryAggregatePrevoteResponse



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `prevote` | [AggregateExchangeRatePrevote](#oracle.v1.AggregateExchangeRatePrevote) |  |  |






<a name="oracle.v1.QueryAggregateVoteRequest"></a>

### QueryAggregateVoteRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `validator` | [string](#string) |  |  |






<a name="oracle.v1.QueryAggregateVoteResponse"></a>

### QueryAggregateVoteResponse



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `vote` | [AggregateExchangeRateVote](#oracle.v1.AggregateExchangeRateVote) |  |  |






<a name="oracle.v1.QueryExchangeRateRequest"></a>

### QueryExchangeRateRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `denom` | [string](#string) |  | token denomination |






<a name="oracle.v1.QueryExchangeRateResponse"></a>

### QueryExchangeRateResponse



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `rate` | [string](#string) |  | token exchange rate |






<a name="oracle.v1.QueryExchangeRatesRequest"></a>

### QueryExchangeRatesRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `pagination` | [cosmos.base.query.v1beta1.PageRequest](#cosmos.base.query.v1beta1.PageRequest) |  | pagination request |






<a name="oracle.v1.QueryExchangeRatesResponse"></a>

### QueryExchangeRatesResponse



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `rates` | [cosmos.base.v1beta1.DecCoin](#cosmos.base.v1beta1.DecCoin) | repeated |  |
| `pagination` | [cosmos.base.query.v1beta1.PageResponse](#cosmos.base.query.v1beta1.PageResponse) |  | pagination response |






<a name="oracle.v1.QueryFeederDelegationRequest"></a>

### QueryFeederDelegationRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `validator` | [string](#string) |  |  |






<a name="oracle.v1.QueryFeederDelegationResponse"></a>

### QueryFeederDelegationResponse



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `address` | [string](#string) |  |  |






<a name="oracle.v1.QueryMissCounterRequest"></a>

### QueryMissCounterRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `validator` | [string](#string) |  |  |






<a name="oracle.v1.QueryMissCounterResponse"></a>

### QueryMissCounterResponse



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `counter` | [int64](#int64) |  |  |






<a name="oracle.v1.QueryParametersRequest"></a>

### QueryParametersRequest







<a name="oracle.v1.QueryParametersResponse"></a>

### QueryParametersResponse



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `params` | [Params](#oracle.v1.Params) |  |  |






<a name="oracle.v1.QueryTobinTaxRequest"></a>

### QueryTobinTaxRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `denom` | [string](#string) |  |  |






<a name="oracle.v1.QueryTobinTaxResponse"></a>

### QueryTobinTaxResponse



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `rate` | [string](#string) |  |  |






<a name="oracle.v1.QueryTobinTaxesRequest"></a>

### QueryTobinTaxesRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `pagination` | [cosmos.base.query.v1beta1.PageRequest](#cosmos.base.query.v1beta1.PageRequest) |  | pagination request |






<a name="oracle.v1.QueryTobinTaxesResponse"></a>

### QueryTobinTaxesResponse



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `rates` | [cosmos.base.v1beta1.DecCoin](#cosmos.base.v1beta1.DecCoin) | repeated |  |
| `pagination` | [cosmos.base.query.v1beta1.PageResponse](#cosmos.base.query.v1beta1.PageResponse) |  | pagination response |






<a name="oracle.v1.QueryVoteTargetsRequest"></a>

### QueryVoteTargetsRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `pagination` | [cosmos.base.query.v1beta1.PageRequest](#cosmos.base.query.v1beta1.PageRequest) |  | pagination request |






<a name="oracle.v1.QueryVoteTargetsResponse"></a>

### QueryVoteTargetsResponse



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `targets` | [string](#string) | repeated |  |
| `pagination` | [cosmos.base.query.v1beta1.PageResponse](#cosmos.base.query.v1beta1.PageResponse) |  | pagination response |





 <!-- end messages -->

 <!-- end enums -->

 <!-- end HasExtensions -->


<a name="oracle.v1.Query"></a>

### Query


| Method Name | Request Type | Response Type | Description | HTTP Verb | Endpoint |
| ----------- | ------------ | ------------- | ------------| ------- | -------- |
| `ExchangeRate` | [QueryExchangeRateRequest](#oracle.v1.QueryExchangeRateRequest) | [QueryExchangeRateResponse](#oracle.v1.QueryExchangeRateResponse) |  | GET|/oracle/v1/exchange_rate/{denom}|
| `ExchangeRates` | [QueryExchangeRatesRequest](#oracle.v1.QueryExchangeRatesRequest) | [QueryExchangeRatesResponse](#oracle.v1.QueryExchangeRatesResponse) |  | GET|/oracle/v1/exchange_rates|
| `Actives` | [QueryActivesRequest](#oracle.v1.QueryActivesRequest) | [QueryActivesResponse](#oracle.v1.QueryActivesResponse) |  | GET|/oracle/v1/actives|
| `Parameters` | [QueryParametersRequest](#oracle.v1.QueryParametersRequest) | [QueryParametersResponse](#oracle.v1.QueryParametersResponse) |  | GET|/oracle/v1/parameters|
| `FeederDelegation` | [QueryFeederDelegationRequest](#oracle.v1.QueryFeederDelegationRequest) | [QueryFeederDelegationResponse](#oracle.v1.QueryFeederDelegationResponse) |  | GET|/oracle/v1/feeder_delegation/{validator}|
| `MissCounter` | [QueryMissCounterRequest](#oracle.v1.QueryMissCounterRequest) | [QueryMissCounterResponse](#oracle.v1.QueryMissCounterResponse) |  | GET|/oracle/v1/miss_counter/{validator}|
| `AggregatePrevote` | [QueryAggregatePrevoteRequest](#oracle.v1.QueryAggregatePrevoteRequest) | [QueryAggregatePrevoteResponse](#oracle.v1.QueryAggregatePrevoteResponse) |  | GET|/oracle/v1/aggregate_prevote/{validator}|
| `AggregateVote` | [QueryAggregateVoteRequest](#oracle.v1.QueryAggregateVoteRequest) | [QueryAggregateVoteResponse](#oracle.v1.QueryAggregateVoteResponse) |  | GET|/oracle/v1/aggregate_vote/{validator}|
| `VoteTargets` | [QueryVoteTargetsRequest](#oracle.v1.QueryVoteTargetsRequest) | [QueryVoteTargetsResponse](#oracle.v1.QueryVoteTargetsResponse) |  | GET|/oracle/v1/vote_targets|
| `TobinTax` | [QueryTobinTaxRequest](#oracle.v1.QueryTobinTaxRequest) | [QueryTobinTaxResponse](#oracle.v1.QueryTobinTaxResponse) |  | GET|/oracle/v1/tobin_tax/{denom}|
| `TobinTaxes` | [QueryTobinTaxesRequest](#oracle.v1.QueryTobinTaxesRequest) | [QueryTobinTaxesResponse](#oracle.v1.QueryTobinTaxesResponse) |  | GET|/oracle/v1/tobin_taxes|

 <!-- end services -->



<a name="oracle/v1/tx.proto"></a>
<p align="right"><a href="#top">Top</a></p>

## oracle/v1/tx.proto



<a name="oracle.v1.MsgAggregateExchangeRatePrevote"></a>

### MsgAggregateExchangeRatePrevote
MsgAggregateExchangeRatePrevote - struct for aggregate prevoting on the ExchangeRateVote.
The purpose of aggregate prevote is to hide vote exchange rates with hash
which is formatted as hex string in SHA256("{salt}:{exchange rate}{denom},...,{exchange rate}{denom}:{voter}")


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `hash` | [bytes](#bytes) |  |  |
| `feeder` | [string](#string) |  |  |
| `validator` | [string](#string) |  |  |






<a name="oracle.v1.MsgAggregateExchangeRatePrevoteResponse"></a>

### MsgAggregateExchangeRatePrevoteResponse







<a name="oracle.v1.MsgAggregateExchangeRateVote"></a>

### MsgAggregateExchangeRateVote
MsgAggregateExchangeRateVote - struct for voting on the exchange rates of Luna denominated in various Terra assets.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `salt` | [string](#string) |  |  |
| `exchange_rates` | [string](#string) |  | NOTE: this exchange rates string is a DecCoins.String() |
| `feeder` | [string](#string) |  |  |
| `validator` | [string](#string) |  |  |






<a name="oracle.v1.MsgAggregateExchangeRateVoteResponse"></a>

### MsgAggregateExchangeRateVoteResponse







<a name="oracle.v1.MsgDelegateFeedConsent"></a>

### MsgDelegateFeedConsent
MsgDelegateFeedConsent - struct for delegating oracle voting rights to another address.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `operator` | [string](#string) |  |  |
| `delegate` | [string](#string) |  |  |






<a name="oracle.v1.MsgDelegateFeedConsentResponse"></a>

### MsgDelegateFeedConsentResponse






 <!-- end messages -->

 <!-- end enums -->

 <!-- end HasExtensions -->


<a name="oracle.v1.Msg"></a>

### Msg
MsgService defines the msgs that the oracle module handles.

| Method Name | Request Type | Response Type | Description | HTTP Verb | Endpoint |
| ----------- | ------------ | ------------- | ------------| ------- | -------- |
| `DelegateFeedConsent` | [MsgDelegateFeedConsent](#oracle.v1.MsgDelegateFeedConsent) | [MsgDelegateFeedConsentResponse](#oracle.v1.MsgDelegateFeedConsentResponse) |  | |
| `AggregateExchangeRatePrevote` | [MsgAggregateExchangeRatePrevote](#oracle.v1.MsgAggregateExchangeRatePrevote) | [MsgAggregateExchangeRatePrevoteResponse](#oracle.v1.MsgAggregateExchangeRatePrevoteResponse) |  | |
| `AggregateExchangeRateVote` | [MsgAggregateExchangeRateVote](#oracle.v1.MsgAggregateExchangeRateVote) | [MsgAggregateExchangeRateVoteResponse](#oracle.v1.MsgAggregateExchangeRateVoteResponse) |  | |

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

