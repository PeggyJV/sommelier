package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// Defines the prefix of each query path
const (
	QueryParameters       = "parameters"
	QueryExchangeRate     = "exchangeRate"
	QueryExchangeRates    = "exchangeRates"
	QueryActives          = "actives"
	QueryPrevotes         = "prevotes"
	QueryVotes            = "votes"
	QueryFeederDelegation = "feederDelegation"
	QueryMissCounter      = "missCounter"
	QueryAggregatePrevote = "aggregatePrevote"
	QueryAggregateVote    = "aggregateVote"
	QueryVoteTargets      = "voteTargets"
	QueryTobinTax         = "tobinTax"
	QueryTobinTaxes       = "tobinTaxes"
)

// QueryExchangeRateParams defines the params for the following queries:
// - 'custom/oracle/exchange_rate'
type QueryExchangeRateParams struct {
	Denom string `json:"denom"`
}

// NewQueryExchangeRateParams returns params for exchange_rate query
func NewQueryExchangeRateParams(denom string) QueryExchangeRateParams {
	return QueryExchangeRateParams{denom}
}

// QueryPrevotesParams defines the params for the following queries:
// - 'custom/oracle/prevotes'
type QueryPrevotesParams struct {
	Voter sdk.ValAddress `json:"voter"`
	Denom string         `json:"denom"`
}

// NewQueryPrevotesParams returns params for exchange_rate prevotes query
func NewQueryPrevotesParams(voter sdk.ValAddress, denom string) QueryPrevotesParams {
	return QueryPrevotesParams{voter, denom}
}

// QueryVotesParams defines the params for the following queries:
// - 'custom/oracle/votes'
type QueryVotesParams struct {
	Voter sdk.ValAddress `json:"voter"`
	Denom string         `json:"denom"`
}

// NewQueryVotesParams returns params for exchange_rate votes query
func NewQueryVotesParams(voter sdk.ValAddress, denom string) QueryVotesParams {
	return QueryVotesParams{voter, denom}
}

// QueryFeederDelegationParams defeins the params for the following queries:
// - 'custom/oracle/feederDelegation'
type QueryFeederDelegationParams struct {
	Validator sdk.ValAddress `json:"validator"`
}

// NewQueryFeederDelegationParams returns params for feeder delegation query
func NewQueryFeederDelegationParams(validator sdk.ValAddress) QueryFeederDelegationParams {
	return QueryFeederDelegationParams{validator}
}

// QueryMissCounterParams defines the params for the following queries:
// - 'custom/oracle/missCounter'
type QueryMissCounterParams struct {
	Validator sdk.ValAddress `json:"validator"`
}

// NewQueryMissCounterParams returns params for feeder delegation query
func NewQueryMissCounterParams(validator sdk.ValAddress) QueryMissCounterParams {
	return QueryMissCounterParams{validator}
}

// QueryAggregatePrevoteParams defines the params for the following queries:
// - 'custom/oracle/aggregatePrevote'
type QueryAggregatePrevoteParams struct {
	Validator sdk.ValAddress `json:"validator"`
}

// NewQueryAggregatePrevoteParams returns params for feeder delegation query
func NewQueryAggregatePrevoteParams(validator sdk.ValAddress) QueryAggregatePrevoteParams {
	return QueryAggregatePrevoteParams{validator}
}

// QueryAggregateVoteParams defines the params for the following queries:
// - 'custom/oracle/aggregateVote'
type QueryAggregateVoteParams struct {
	Validator sdk.ValAddress `json:"validator"`
}

// NewQueryAggregateVoteParams returns params for feeder delegation query
func NewQueryAggregateVoteParams(validator sdk.ValAddress) QueryAggregateVoteParams {
	return QueryAggregateVoteParams{validator}
}

// QueryTobinTaxParams defines the params for the following queries:
// - 'custom/oracle/tobinTax'
type QueryTobinTaxParams struct {
	Denom string `json:"denom"`
}

// NewQueryTobinTaxParams returns params for tobin tax query
func NewQueryTobinTaxParams(denom string) QueryTobinTaxParams {
	return QueryTobinTaxParams{denom}
}
