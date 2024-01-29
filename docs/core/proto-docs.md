<!-- This file is auto-generated. Please do not modify it yourself. -->
# Protobuf Documentation
<a name="top"></a>

## Table of Contents

- [auction/v1/auction.proto](#auction/v1/auction.proto)
    - [Auction](#auction.v1.Auction)
    - [Bid](#auction.v1.Bid)
    - [ProposedTokenPrice](#auction.v1.ProposedTokenPrice)
    - [TokenPrice](#auction.v1.TokenPrice)
  
- [auction/v1/tx.proto](#auction/v1/tx.proto)
    - [MsgSubmitBidRequest](#auction.v1.MsgSubmitBidRequest)
    - [MsgSubmitBidResponse](#auction.v1.MsgSubmitBidResponse)
  
    - [Msg](#auction.v1.Msg)
  
- [auction/v1/genesis.proto](#auction/v1/genesis.proto)
    - [GenesisState](#auction.v1.GenesisState)
    - [Params](#auction.v1.Params)
  
- [auction/v1/proposal.proto](#auction/v1/proposal.proto)
    - [SetTokenPricesProposal](#auction.v1.SetTokenPricesProposal)
    - [SetTokenPricesProposalWithDeposit](#auction.v1.SetTokenPricesProposalWithDeposit)
  
- [auction/v1/query.proto](#auction/v1/query.proto)
    - [QueryActiveAuctionRequest](#auction.v1.QueryActiveAuctionRequest)
    - [QueryActiveAuctionResponse](#auction.v1.QueryActiveAuctionResponse)
    - [QueryActiveAuctionsRequest](#auction.v1.QueryActiveAuctionsRequest)
    - [QueryActiveAuctionsResponse](#auction.v1.QueryActiveAuctionsResponse)
    - [QueryBidRequest](#auction.v1.QueryBidRequest)
    - [QueryBidResponse](#auction.v1.QueryBidResponse)
    - [QueryBidsByAuctionRequest](#auction.v1.QueryBidsByAuctionRequest)
    - [QueryBidsByAuctionResponse](#auction.v1.QueryBidsByAuctionResponse)
    - [QueryEndedAuctionRequest](#auction.v1.QueryEndedAuctionRequest)
    - [QueryEndedAuctionResponse](#auction.v1.QueryEndedAuctionResponse)
    - [QueryEndedAuctionsRequest](#auction.v1.QueryEndedAuctionsRequest)
    - [QueryEndedAuctionsResponse](#auction.v1.QueryEndedAuctionsResponse)
    - [QueryParamsRequest](#auction.v1.QueryParamsRequest)
    - [QueryParamsResponse](#auction.v1.QueryParamsResponse)
    - [QueryTokenPriceRequest](#auction.v1.QueryTokenPriceRequest)
    - [QueryTokenPriceResponse](#auction.v1.QueryTokenPriceResponse)
    - [QueryTokenPricesRequest](#auction.v1.QueryTokenPricesRequest)
    - [QueryTokenPricesResponse](#auction.v1.QueryTokenPricesResponse)
  
    - [Query](#auction.v1.Query)
  
- [axelarcork/v1/axelarcork.proto](#axelarcork/v1/axelarcork.proto)
    - [AxelarContractCallNonce](#axelarcork.v1.AxelarContractCallNonce)
    - [AxelarCork](#axelarcork.v1.AxelarCork)
    - [AxelarCorkResult](#axelarcork.v1.AxelarCorkResult)
    - [AxelarCorkResults](#axelarcork.v1.AxelarCorkResults)
    - [AxelarUpgradeData](#axelarcork.v1.AxelarUpgradeData)
    - [CellarIDSet](#axelarcork.v1.CellarIDSet)
    - [ChainConfiguration](#axelarcork.v1.ChainConfiguration)
    - [ChainConfigurations](#axelarcork.v1.ChainConfigurations)
    - [ScheduledAxelarCork](#axelarcork.v1.ScheduledAxelarCork)
    - [ScheduledAxelarCorks](#axelarcork.v1.ScheduledAxelarCorks)
  
- [axelarcork/v1/event.proto](#axelarcork/v1/event.proto)
    - [ScheduleCorkEvent](#axelarcork.v1.ScheduleCorkEvent)
  
- [axelarcork/v1/genesis.proto](#axelarcork/v1/genesis.proto)
    - [GenesisState](#axelarcork.v1.GenesisState)
    - [Params](#axelarcork.v1.Params)
  
- [axelarcork/v1/proposal.proto](#axelarcork/v1/proposal.proto)
    - [AddAxelarManagedCellarIDsProposal](#axelarcork.v1.AddAxelarManagedCellarIDsProposal)
    - [AddAxelarManagedCellarIDsProposalWithDeposit](#axelarcork.v1.AddAxelarManagedCellarIDsProposalWithDeposit)
    - [AddChainConfigurationProposal](#axelarcork.v1.AddChainConfigurationProposal)
    - [AddChainConfigurationProposalWithDeposit](#axelarcork.v1.AddChainConfigurationProposalWithDeposit)
    - [AxelarCommunityPoolSpendProposal](#axelarcork.v1.AxelarCommunityPoolSpendProposal)
    - [AxelarCommunityPoolSpendProposalForCLI](#axelarcork.v1.AxelarCommunityPoolSpendProposalForCLI)
    - [AxelarScheduledCorkProposal](#axelarcork.v1.AxelarScheduledCorkProposal)
    - [AxelarScheduledCorkProposalWithDeposit](#axelarcork.v1.AxelarScheduledCorkProposalWithDeposit)
    - [CancelAxelarProxyContractUpgradeProposal](#axelarcork.v1.CancelAxelarProxyContractUpgradeProposal)
    - [CancelAxelarProxyContractUpgradeProposalWithDeposit](#axelarcork.v1.CancelAxelarProxyContractUpgradeProposalWithDeposit)
    - [RemoveAxelarManagedCellarIDsProposal](#axelarcork.v1.RemoveAxelarManagedCellarIDsProposal)
    - [RemoveAxelarManagedCellarIDsProposalWithDeposit](#axelarcork.v1.RemoveAxelarManagedCellarIDsProposalWithDeposit)
    - [RemoveChainConfigurationProposal](#axelarcork.v1.RemoveChainConfigurationProposal)
    - [RemoveChainConfigurationProposalWithDeposit](#axelarcork.v1.RemoveChainConfigurationProposalWithDeposit)
    - [UpgradeAxelarProxyContractProposal](#axelarcork.v1.UpgradeAxelarProxyContractProposal)
    - [UpgradeAxelarProxyContractProposalWithDeposit](#axelarcork.v1.UpgradeAxelarProxyContractProposalWithDeposit)
  
- [axelarcork/v1/query.proto](#axelarcork/v1/query.proto)
    - [QueryAxelarContractCallNoncesRequest](#axelarcork.v1.QueryAxelarContractCallNoncesRequest)
    - [QueryAxelarContractCallNoncesResponse](#axelarcork.v1.QueryAxelarContractCallNoncesResponse)
    - [QueryAxelarProxyUpgradeDataRequest](#axelarcork.v1.QueryAxelarProxyUpgradeDataRequest)
    - [QueryAxelarProxyUpgradeDataResponse](#axelarcork.v1.QueryAxelarProxyUpgradeDataResponse)
    - [QueryCellarIDsByChainIDRequest](#axelarcork.v1.QueryCellarIDsByChainIDRequest)
    - [QueryCellarIDsByChainIDResponse](#axelarcork.v1.QueryCellarIDsByChainIDResponse)
    - [QueryCellarIDsRequest](#axelarcork.v1.QueryCellarIDsRequest)
    - [QueryCellarIDsResponse](#axelarcork.v1.QueryCellarIDsResponse)
    - [QueryChainConfigurationsRequest](#axelarcork.v1.QueryChainConfigurationsRequest)
    - [QueryChainConfigurationsResponse](#axelarcork.v1.QueryChainConfigurationsResponse)
    - [QueryCorkResultRequest](#axelarcork.v1.QueryCorkResultRequest)
    - [QueryCorkResultResponse](#axelarcork.v1.QueryCorkResultResponse)
    - [QueryCorkResultsRequest](#axelarcork.v1.QueryCorkResultsRequest)
    - [QueryCorkResultsResponse](#axelarcork.v1.QueryCorkResultsResponse)
    - [QueryParamsRequest](#axelarcork.v1.QueryParamsRequest)
    - [QueryParamsResponse](#axelarcork.v1.QueryParamsResponse)
    - [QueryScheduledBlockHeightsRequest](#axelarcork.v1.QueryScheduledBlockHeightsRequest)
    - [QueryScheduledBlockHeightsResponse](#axelarcork.v1.QueryScheduledBlockHeightsResponse)
    - [QueryScheduledCorksByBlockHeightRequest](#axelarcork.v1.QueryScheduledCorksByBlockHeightRequest)
    - [QueryScheduledCorksByBlockHeightResponse](#axelarcork.v1.QueryScheduledCorksByBlockHeightResponse)
    - [QueryScheduledCorksByIDRequest](#axelarcork.v1.QueryScheduledCorksByIDRequest)
    - [QueryScheduledCorksByIDResponse](#axelarcork.v1.QueryScheduledCorksByIDResponse)
    - [QueryScheduledCorksRequest](#axelarcork.v1.QueryScheduledCorksRequest)
    - [QueryScheduledCorksResponse](#axelarcork.v1.QueryScheduledCorksResponse)
    - [QueryWinningAxelarCorkRequest](#axelarcork.v1.QueryWinningAxelarCorkRequest)
    - [QueryWinningAxelarCorkResponse](#axelarcork.v1.QueryWinningAxelarCorkResponse)
  
    - [Query](#axelarcork.v1.Query)
  
- [axelarcork/v1/tx.proto](#axelarcork/v1/tx.proto)
    - [MsgBumpAxelarCorkGasRequest](#axelarcork.v1.MsgBumpAxelarCorkGasRequest)
    - [MsgBumpAxelarCorkGasResponse](#axelarcork.v1.MsgBumpAxelarCorkGasResponse)
    - [MsgCancelAxelarCorkRequest](#axelarcork.v1.MsgCancelAxelarCorkRequest)
    - [MsgCancelAxelarCorkResponse](#axelarcork.v1.MsgCancelAxelarCorkResponse)
    - [MsgRelayAxelarCorkRequest](#axelarcork.v1.MsgRelayAxelarCorkRequest)
    - [MsgRelayAxelarCorkResponse](#axelarcork.v1.MsgRelayAxelarCorkResponse)
    - [MsgRelayAxelarProxyUpgradeRequest](#axelarcork.v1.MsgRelayAxelarProxyUpgradeRequest)
    - [MsgRelayAxelarProxyUpgradeResponse](#axelarcork.v1.MsgRelayAxelarProxyUpgradeResponse)
    - [MsgScheduleAxelarCorkRequest](#axelarcork.v1.MsgScheduleAxelarCorkRequest)
    - [MsgScheduleAxelarCorkResponse](#axelarcork.v1.MsgScheduleAxelarCorkResponse)
  
    - [Msg](#axelarcork.v1.Msg)
  
- [cellarfees/v1/cellarfees.proto](#cellarfees/v1/cellarfees.proto)
    - [FeeAccrualCounter](#cellarfees.v1.FeeAccrualCounter)
    - [FeeAccrualCounters](#cellarfees.v1.FeeAccrualCounters)
  
- [cellarfees/v1/params.proto](#cellarfees/v1/params.proto)
    - [Params](#cellarfees.v1.Params)
  
- [cellarfees/v1/genesis.proto](#cellarfees/v1/genesis.proto)
    - [GenesisState](#cellarfees.v1.GenesisState)
  
- [cellarfees/v1/query.proto](#cellarfees/v1/query.proto)
    - [QueryAPYRequest](#cellarfees.v1.QueryAPYRequest)
    - [QueryAPYResponse](#cellarfees.v1.QueryAPYResponse)
    - [QueryFeeAccrualCountersRequest](#cellarfees.v1.QueryFeeAccrualCountersRequest)
    - [QueryFeeAccrualCountersResponse](#cellarfees.v1.QueryFeeAccrualCountersResponse)
    - [QueryLastRewardSupplyPeakRequest](#cellarfees.v1.QueryLastRewardSupplyPeakRequest)
    - [QueryLastRewardSupplyPeakResponse](#cellarfees.v1.QueryLastRewardSupplyPeakResponse)
    - [QueryModuleAccountsRequest](#cellarfees.v1.QueryModuleAccountsRequest)
    - [QueryModuleAccountsResponse](#cellarfees.v1.QueryModuleAccountsResponse)
    - [QueryParamsRequest](#cellarfees.v1.QueryParamsRequest)
    - [QueryParamsResponse](#cellarfees.v1.QueryParamsResponse)
  
    - [Query](#cellarfees.v1.Query)
  
- [cork/v2/cork.proto](#cork/v2/cork.proto)
    - [CellarIDSet](#cork.v2.CellarIDSet)
    - [Cork](#cork.v2.Cork)
    - [CorkResult](#cork.v2.CorkResult)
    - [ScheduledCork](#cork.v2.ScheduledCork)
  
- [cork/v2/genesis.proto](#cork/v2/genesis.proto)
    - [GenesisState](#cork.v2.GenesisState)
    - [Params](#cork.v2.Params)
  
- [cork/v2/proposal.proto](#cork/v2/proposal.proto)
    - [AddManagedCellarIDsProposal](#cork.v2.AddManagedCellarIDsProposal)
    - [AddManagedCellarIDsProposalWithDeposit](#cork.v2.AddManagedCellarIDsProposalWithDeposit)
    - [RemoveManagedCellarIDsProposal](#cork.v2.RemoveManagedCellarIDsProposal)
    - [RemoveManagedCellarIDsProposalWithDeposit](#cork.v2.RemoveManagedCellarIDsProposalWithDeposit)
    - [ScheduledCorkProposal](#cork.v2.ScheduledCorkProposal)
    - [ScheduledCorkProposalWithDeposit](#cork.v2.ScheduledCorkProposalWithDeposit)
  
- [cork/v2/query.proto](#cork/v2/query.proto)
    - [QueryCellarIDsRequest](#cork.v2.QueryCellarIDsRequest)
    - [QueryCellarIDsResponse](#cork.v2.QueryCellarIDsResponse)
    - [QueryCorkResultRequest](#cork.v2.QueryCorkResultRequest)
    - [QueryCorkResultResponse](#cork.v2.QueryCorkResultResponse)
    - [QueryCorkResultsRequest](#cork.v2.QueryCorkResultsRequest)
    - [QueryCorkResultsResponse](#cork.v2.QueryCorkResultsResponse)
    - [QueryParamsRequest](#cork.v2.QueryParamsRequest)
    - [QueryParamsResponse](#cork.v2.QueryParamsResponse)
    - [QueryScheduledBlockHeightsRequest](#cork.v2.QueryScheduledBlockHeightsRequest)
    - [QueryScheduledBlockHeightsResponse](#cork.v2.QueryScheduledBlockHeightsResponse)
    - [QueryScheduledCorksByBlockHeightRequest](#cork.v2.QueryScheduledCorksByBlockHeightRequest)
    - [QueryScheduledCorksByBlockHeightResponse](#cork.v2.QueryScheduledCorksByBlockHeightResponse)
    - [QueryScheduledCorksByIDRequest](#cork.v2.QueryScheduledCorksByIDRequest)
    - [QueryScheduledCorksByIDResponse](#cork.v2.QueryScheduledCorksByIDResponse)
    - [QueryScheduledCorksRequest](#cork.v2.QueryScheduledCorksRequest)
    - [QueryScheduledCorksResponse](#cork.v2.QueryScheduledCorksResponse)
  
    - [Query](#cork.v2.Query)
  
- [cork/v2/tx.proto](#cork/v2/tx.proto)
    - [MsgScheduleCorkRequest](#cork.v2.MsgScheduleCorkRequest)
    - [MsgScheduleCorkResponse](#cork.v2.MsgScheduleCorkResponse)
  
    - [Msg](#cork.v2.Msg)
  
- [incentives/v1/genesis.proto](#incentives/v1/genesis.proto)
    - [GenesisState](#incentives.v1.GenesisState)
    - [Params](#incentives.v1.Params)
  
- [incentives/v1/query.proto](#incentives/v1/query.proto)
    - [QueryAPYRequest](#incentives.v1.QueryAPYRequest)
    - [QueryAPYResponse](#incentives.v1.QueryAPYResponse)
    - [QueryParamsRequest](#incentives.v1.QueryParamsRequest)
    - [QueryParamsResponse](#incentives.v1.QueryParamsResponse)
  
    - [Query](#incentives.v1.Query)
  
- [pubsub/v1/params.proto](#pubsub/v1/params.proto)
    - [Params](#pubsub.v1.Params)
  
- [pubsub/v1/pubsub.proto](#pubsub/v1/pubsub.proto)
    - [AddDefaultSubscriptionProposal](#pubsub.v1.AddDefaultSubscriptionProposal)
    - [AddDefaultSubscriptionProposalWithDeposit](#pubsub.v1.AddDefaultSubscriptionProposalWithDeposit)
    - [AddPublisherProposal](#pubsub.v1.AddPublisherProposal)
    - [AddPublisherProposalWithDeposit](#pubsub.v1.AddPublisherProposalWithDeposit)
    - [DefaultSubscription](#pubsub.v1.DefaultSubscription)
    - [Publisher](#pubsub.v1.Publisher)
    - [PublisherIntent](#pubsub.v1.PublisherIntent)
    - [RemoveDefaultSubscriptionProposal](#pubsub.v1.RemoveDefaultSubscriptionProposal)
    - [RemoveDefaultSubscriptionProposalWithDeposit](#pubsub.v1.RemoveDefaultSubscriptionProposalWithDeposit)
    - [RemovePublisherProposal](#pubsub.v1.RemovePublisherProposal)
    - [RemovePublisherProposalWithDeposit](#pubsub.v1.RemovePublisherProposalWithDeposit)
    - [Subscriber](#pubsub.v1.Subscriber)
    - [SubscriberIntent](#pubsub.v1.SubscriberIntent)
  
    - [AllowedSubscribers](#pubsub.v1.AllowedSubscribers)
    - [PublishMethod](#pubsub.v1.PublishMethod)
  
- [pubsub/v1/genesis.proto](#pubsub/v1/genesis.proto)
    - [GenesisState](#pubsub.v1.GenesisState)
  
- [pubsub/v1/query.proto](#pubsub/v1/query.proto)
    - [QueryDefaultSubscriptionRequest](#pubsub.v1.QueryDefaultSubscriptionRequest)
    - [QueryDefaultSubscriptionResponse](#pubsub.v1.QueryDefaultSubscriptionResponse)
    - [QueryDefaultSubscriptionsRequest](#pubsub.v1.QueryDefaultSubscriptionsRequest)
    - [QueryDefaultSubscriptionsResponse](#pubsub.v1.QueryDefaultSubscriptionsResponse)
    - [QueryParamsRequest](#pubsub.v1.QueryParamsRequest)
    - [QueryParamsResponse](#pubsub.v1.QueryParamsResponse)
    - [QueryPublisherIntentRequest](#pubsub.v1.QueryPublisherIntentRequest)
    - [QueryPublisherIntentResponse](#pubsub.v1.QueryPublisherIntentResponse)
    - [QueryPublisherIntentsByPublisherDomainRequest](#pubsub.v1.QueryPublisherIntentsByPublisherDomainRequest)
    - [QueryPublisherIntentsByPublisherDomainResponse](#pubsub.v1.QueryPublisherIntentsByPublisherDomainResponse)
    - [QueryPublisherIntentsBySubscriptionIDRequest](#pubsub.v1.QueryPublisherIntentsBySubscriptionIDRequest)
    - [QueryPublisherIntentsBySubscriptionIDResponse](#pubsub.v1.QueryPublisherIntentsBySubscriptionIDResponse)
    - [QueryPublisherIntentsRequest](#pubsub.v1.QueryPublisherIntentsRequest)
    - [QueryPublisherIntentsResponse](#pubsub.v1.QueryPublisherIntentsResponse)
    - [QueryPublisherRequest](#pubsub.v1.QueryPublisherRequest)
    - [QueryPublisherResponse](#pubsub.v1.QueryPublisherResponse)
    - [QueryPublishersRequest](#pubsub.v1.QueryPublishersRequest)
    - [QueryPublishersResponse](#pubsub.v1.QueryPublishersResponse)
    - [QuerySubscriberIntentRequest](#pubsub.v1.QuerySubscriberIntentRequest)
    - [QuerySubscriberIntentResponse](#pubsub.v1.QuerySubscriberIntentResponse)
    - [QuerySubscriberIntentsByPublisherDomainRequest](#pubsub.v1.QuerySubscriberIntentsByPublisherDomainRequest)
    - [QuerySubscriberIntentsByPublisherDomainResponse](#pubsub.v1.QuerySubscriberIntentsByPublisherDomainResponse)
    - [QuerySubscriberIntentsBySubscriberAddressRequest](#pubsub.v1.QuerySubscriberIntentsBySubscriberAddressRequest)
    - [QuerySubscriberIntentsBySubscriberAddressResponse](#pubsub.v1.QuerySubscriberIntentsBySubscriberAddressResponse)
    - [QuerySubscriberIntentsBySubscriptionIDRequest](#pubsub.v1.QuerySubscriberIntentsBySubscriptionIDRequest)
    - [QuerySubscriberIntentsBySubscriptionIDResponse](#pubsub.v1.QuerySubscriberIntentsBySubscriptionIDResponse)
    - [QuerySubscriberIntentsRequest](#pubsub.v1.QuerySubscriberIntentsRequest)
    - [QuerySubscriberIntentsResponse](#pubsub.v1.QuerySubscriberIntentsResponse)
    - [QuerySubscriberRequest](#pubsub.v1.QuerySubscriberRequest)
    - [QuerySubscriberResponse](#pubsub.v1.QuerySubscriberResponse)
    - [QuerySubscribersRequest](#pubsub.v1.QuerySubscribersRequest)
    - [QuerySubscribersResponse](#pubsub.v1.QuerySubscribersResponse)
  
    - [Query](#pubsub.v1.Query)
  
- [pubsub/v1/tx.proto](#pubsub/v1/tx.proto)
    - [MsgAddPublisherIntentRequest](#pubsub.v1.MsgAddPublisherIntentRequest)
    - [MsgAddPublisherIntentResponse](#pubsub.v1.MsgAddPublisherIntentResponse)
    - [MsgAddSubscriberIntentRequest](#pubsub.v1.MsgAddSubscriberIntentRequest)
    - [MsgAddSubscriberIntentResponse](#pubsub.v1.MsgAddSubscriberIntentResponse)
    - [MsgAddSubscriberRequest](#pubsub.v1.MsgAddSubscriberRequest)
    - [MsgAddSubscriberResponse](#pubsub.v1.MsgAddSubscriberResponse)
    - [MsgRemovePublisherIntentRequest](#pubsub.v1.MsgRemovePublisherIntentRequest)
    - [MsgRemovePublisherIntentResponse](#pubsub.v1.MsgRemovePublisherIntentResponse)
    - [MsgRemovePublisherRequest](#pubsub.v1.MsgRemovePublisherRequest)
    - [MsgRemovePublisherResponse](#pubsub.v1.MsgRemovePublisherResponse)
    - [MsgRemoveSubscriberIntentRequest](#pubsub.v1.MsgRemoveSubscriberIntentRequest)
    - [MsgRemoveSubscriberIntentResponse](#pubsub.v1.MsgRemoveSubscriberIntentResponse)
    - [MsgRemoveSubscriberRequest](#pubsub.v1.MsgRemoveSubscriberRequest)
    - [MsgRemoveSubscriberResponse](#pubsub.v1.MsgRemoveSubscriberResponse)
  
    - [Msg](#pubsub.v1.Msg)
  
- [Scalar Value Types](#scalar-value-types)



<a name="auction/v1/auction.proto"></a>
<p align="right"><a href="#top">Top</a></p>

## auction/v1/auction.proto



<a name="auction.v1.Auction"></a>

### Auction



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `id` | [uint32](#uint32) |  |  |
| `starting_tokens_for_sale` | [cosmos.base.v1beta1.Coin](#cosmos.base.v1beta1.Coin) |  |  |
| `start_block` | [uint64](#uint64) |  |  |
| `end_block` | [uint64](#uint64) |  |  |
| `initial_price_decrease_rate` | [string](#string) |  |  |
| `current_price_decrease_rate` | [string](#string) |  |  |
| `price_decrease_block_interval` | [uint64](#uint64) |  |  |
| `initial_unit_price_in_usomm` | [string](#string) |  |  |
| `current_unit_price_in_usomm` | [string](#string) |  |  |
| `remaining_tokens_for_sale` | [cosmos.base.v1beta1.Coin](#cosmos.base.v1beta1.Coin) |  |  |
| `funding_module_account` | [string](#string) |  |  |
| `proceeds_module_account` | [string](#string) |  |  |






<a name="auction.v1.Bid"></a>

### Bid



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `id` | [uint64](#uint64) |  |  |
| `auction_id` | [uint32](#uint32) |  |  |
| `bidder` | [string](#string) |  |  |
| `max_bid_in_usomm` | [cosmos.base.v1beta1.Coin](#cosmos.base.v1beta1.Coin) |  |  |
| `sale_token_minimum_amount` | [cosmos.base.v1beta1.Coin](#cosmos.base.v1beta1.Coin) |  |  |
| `total_fulfilled_sale_tokens` | [cosmos.base.v1beta1.Coin](#cosmos.base.v1beta1.Coin) |  |  |
| `sale_token_unit_price_in_usomm` | [string](#string) |  |  |
| `total_usomm_paid` | [cosmos.base.v1beta1.Coin](#cosmos.base.v1beta1.Coin) |  |  |
| `block_height` | [uint64](#uint64) |  |  |






<a name="auction.v1.ProposedTokenPrice"></a>

### ProposedTokenPrice



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `denom` | [string](#string) |  |  |
| `exponent` | [uint64](#uint64) |  |  |
| `usd_price` | [string](#string) |  |  |






<a name="auction.v1.TokenPrice"></a>

### TokenPrice
USD price is the value for one non-fractional token (smallest unit of the token * 10^exponent)


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `denom` | [string](#string) |  |  |
| `exponent` | [uint64](#uint64) |  |  |
| `usd_price` | [string](#string) |  |  |
| `last_updated_block` | [uint64](#uint64) |  |  |





 <!-- end messages -->

 <!-- end enums -->

 <!-- end HasExtensions -->

 <!-- end services -->



<a name="auction/v1/tx.proto"></a>
<p align="right"><a href="#top">Top</a></p>

## auction/v1/tx.proto



<a name="auction.v1.MsgSubmitBidRequest"></a>

### MsgSubmitBidRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `auction_id` | [uint32](#uint32) |  |  |
| `signer` | [string](#string) |  |  |
| `max_bid_in_usomm` | [cosmos.base.v1beta1.Coin](#cosmos.base.v1beta1.Coin) |  |  |
| `sale_token_minimum_amount` | [cosmos.base.v1beta1.Coin](#cosmos.base.v1beta1.Coin) |  |  |






<a name="auction.v1.MsgSubmitBidResponse"></a>

### MsgSubmitBidResponse



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `bid` | [Bid](#auction.v1.Bid) |  |  |





 <!-- end messages -->

 <!-- end enums -->

 <!-- end HasExtensions -->


<a name="auction.v1.Msg"></a>

### Msg


| Method Name | Request Type | Response Type | Description | HTTP Verb | Endpoint |
| ----------- | ------------ | ------------- | ------------| ------- | -------- |
| `SubmitBid` | [MsgSubmitBidRequest](#auction.v1.MsgSubmitBidRequest) | [MsgSubmitBidResponse](#auction.v1.MsgSubmitBidResponse) |  | |

 <!-- end services -->



<a name="auction/v1/genesis.proto"></a>
<p align="right"><a href="#top">Top</a></p>

## auction/v1/genesis.proto



<a name="auction.v1.GenesisState"></a>

### GenesisState



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `params` | [Params](#auction.v1.Params) |  |  |
| `auctions` | [Auction](#auction.v1.Auction) | repeated |  |
| `bids` | [Bid](#auction.v1.Bid) | repeated |  |
| `token_prices` | [TokenPrice](#auction.v1.TokenPrice) | repeated |  |
| `last_auction_id` | [uint32](#uint32) |  |  |
| `last_bid_id` | [uint64](#uint64) |  |  |






<a name="auction.v1.Params"></a>

### Params



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `price_max_block_age` | [uint64](#uint64) |  |  |
| `minimum_bid_in_usomm` | [uint64](#uint64) |  |  |
| `minimum_sale_tokens_usd_value` | [string](#string) |  |  |
| `auction_max_block_age` | [uint64](#uint64) |  |  |
| `auction_price_decrease_acceleration_rate` | [string](#string) |  |  |
| `minimum_auction_height` | [uint64](#uint64) |  |  |





 <!-- end messages -->

 <!-- end enums -->

 <!-- end HasExtensions -->

 <!-- end services -->



<a name="auction/v1/proposal.proto"></a>
<p align="right"><a href="#top">Top</a></p>

## auction/v1/proposal.proto



<a name="auction.v1.SetTokenPricesProposal"></a>

### SetTokenPricesProposal



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `title` | [string](#string) |  |  |
| `description` | [string](#string) |  |  |
| `token_prices` | [ProposedTokenPrice](#auction.v1.ProposedTokenPrice) | repeated |  |






<a name="auction.v1.SetTokenPricesProposalWithDeposit"></a>

### SetTokenPricesProposalWithDeposit



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `title` | [string](#string) |  |  |
| `description` | [string](#string) |  |  |
| `token_prices` | [ProposedTokenPrice](#auction.v1.ProposedTokenPrice) | repeated |  |
| `deposit` | [string](#string) |  |  |





 <!-- end messages -->

 <!-- end enums -->

 <!-- end HasExtensions -->

 <!-- end services -->



<a name="auction/v1/query.proto"></a>
<p align="right"><a href="#top">Top</a></p>

## auction/v1/query.proto



<a name="auction.v1.QueryActiveAuctionRequest"></a>

### QueryActiveAuctionRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `auction_id` | [uint32](#uint32) |  |  |






<a name="auction.v1.QueryActiveAuctionResponse"></a>

### QueryActiveAuctionResponse



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `auction` | [Auction](#auction.v1.Auction) |  |  |






<a name="auction.v1.QueryActiveAuctionsRequest"></a>

### QueryActiveAuctionsRequest







<a name="auction.v1.QueryActiveAuctionsResponse"></a>

### QueryActiveAuctionsResponse



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `auctions` | [Auction](#auction.v1.Auction) | repeated |  |






<a name="auction.v1.QueryBidRequest"></a>

### QueryBidRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `bid_id` | [uint64](#uint64) |  |  |
| `auction_id` | [uint32](#uint32) |  |  |






<a name="auction.v1.QueryBidResponse"></a>

### QueryBidResponse



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `bid` | [Bid](#auction.v1.Bid) |  |  |






<a name="auction.v1.QueryBidsByAuctionRequest"></a>

### QueryBidsByAuctionRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `auction_id` | [uint32](#uint32) |  |  |
| `pagination` | [cosmos.base.query.v1beta1.PageRequest](#cosmos.base.query.v1beta1.PageRequest) |  |  |






<a name="auction.v1.QueryBidsByAuctionResponse"></a>

### QueryBidsByAuctionResponse



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `bids` | [Bid](#auction.v1.Bid) | repeated |  |
| `pagination` | [cosmos.base.query.v1beta1.PageResponse](#cosmos.base.query.v1beta1.PageResponse) |  |  |






<a name="auction.v1.QueryEndedAuctionRequest"></a>

### QueryEndedAuctionRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `auction_id` | [uint32](#uint32) |  |  |






<a name="auction.v1.QueryEndedAuctionResponse"></a>

### QueryEndedAuctionResponse



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `auction` | [Auction](#auction.v1.Auction) |  |  |






<a name="auction.v1.QueryEndedAuctionsRequest"></a>

### QueryEndedAuctionsRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `pagination` | [cosmos.base.query.v1beta1.PageRequest](#cosmos.base.query.v1beta1.PageRequest) |  |  |






<a name="auction.v1.QueryEndedAuctionsResponse"></a>

### QueryEndedAuctionsResponse



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `auctions` | [Auction](#auction.v1.Auction) | repeated |  |
| `pagination` | [cosmos.base.query.v1beta1.PageResponse](#cosmos.base.query.v1beta1.PageResponse) |  |  |






<a name="auction.v1.QueryParamsRequest"></a>

### QueryParamsRequest







<a name="auction.v1.QueryParamsResponse"></a>

### QueryParamsResponse



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `params` | [Params](#auction.v1.Params) |  |  |






<a name="auction.v1.QueryTokenPriceRequest"></a>

### QueryTokenPriceRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `denom` | [string](#string) |  |  |






<a name="auction.v1.QueryTokenPriceResponse"></a>

### QueryTokenPriceResponse



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `token_price` | [TokenPrice](#auction.v1.TokenPrice) |  |  |






<a name="auction.v1.QueryTokenPricesRequest"></a>

### QueryTokenPricesRequest







<a name="auction.v1.QueryTokenPricesResponse"></a>

### QueryTokenPricesResponse



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `token_prices` | [TokenPrice](#auction.v1.TokenPrice) | repeated |  |





 <!-- end messages -->

 <!-- end enums -->

 <!-- end HasExtensions -->


<a name="auction.v1.Query"></a>

### Query


| Method Name | Request Type | Response Type | Description | HTTP Verb | Endpoint |
| ----------- | ------------ | ------------- | ------------| ------- | -------- |
| `QueryParams` | [QueryParamsRequest](#auction.v1.QueryParamsRequest) | [QueryParamsResponse](#auction.v1.QueryParamsResponse) |  | GET|/sommelier/auction/v1/params|
| `QueryActiveAuction` | [QueryActiveAuctionRequest](#auction.v1.QueryActiveAuctionRequest) | [QueryActiveAuctionResponse](#auction.v1.QueryActiveAuctionResponse) |  | GET|/sommelier/auction/v1/active_auctions/{auction_id}|
| `QueryEndedAuction` | [QueryEndedAuctionRequest](#auction.v1.QueryEndedAuctionRequest) | [QueryEndedAuctionResponse](#auction.v1.QueryEndedAuctionResponse) |  | GET|/sommelier/auction/v1/ended_auctions/{auction_id}|
| `QueryActiveAuctions` | [QueryActiveAuctionsRequest](#auction.v1.QueryActiveAuctionsRequest) | [QueryActiveAuctionsResponse](#auction.v1.QueryActiveAuctionsResponse) |  | GET|/sommelier/auction/v1/active_auctions|
| `QueryEndedAuctions` | [QueryEndedAuctionsRequest](#auction.v1.QueryEndedAuctionsRequest) | [QueryEndedAuctionsResponse](#auction.v1.QueryEndedAuctionsResponse) |  | GET|/sommelier/auction/v1/ended_auctions|
| `QueryBid` | [QueryBidRequest](#auction.v1.QueryBidRequest) | [QueryBidResponse](#auction.v1.QueryBidResponse) |  | GET|/sommelier/auction/v1/auctions/{auction_id}/bids/{bid_id}|
| `QueryBidsByAuction` | [QueryBidsByAuctionRequest](#auction.v1.QueryBidsByAuctionRequest) | [QueryBidsByAuctionResponse](#auction.v1.QueryBidsByAuctionResponse) |  | GET|/sommelier/auction/v1/auctions/{auction_id}/bids|
| `QueryTokenPrice` | [QueryTokenPriceRequest](#auction.v1.QueryTokenPriceRequest) | [QueryTokenPriceResponse](#auction.v1.QueryTokenPriceResponse) |  | GET|/sommelier/auction/v1/token_prices/{denom}|
| `QueryTokenPrices` | [QueryTokenPricesRequest](#auction.v1.QueryTokenPricesRequest) | [QueryTokenPricesResponse](#auction.v1.QueryTokenPricesResponse) |  | GET|/sommelier/auction/v1/token_prices|

 <!-- end services -->



<a name="axelarcork/v1/axelarcork.proto"></a>
<p align="right"><a href="#top">Top</a></p>

## axelarcork/v1/axelarcork.proto



<a name="axelarcork.v1.AxelarContractCallNonce"></a>

### AxelarContractCallNonce
Used to enforce strictly newer call ordering per contract


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `chain_id` | [uint64](#uint64) |  |  |
| `contract_address` | [string](#string) |  |  |
| `nonce` | [uint64](#uint64) |  |  |






<a name="axelarcork.v1.AxelarCork"></a>

### AxelarCork



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `encoded_contract_call` | [bytes](#bytes) |  | call body containing the ABI encoded bytes to send to the contract |
| `chain_id` | [uint64](#uint64) |  | the chain ID of the evm target chain |
| `target_contract_address` | [string](#string) |  | address of the contract to send the call |
| `deadline` | [uint64](#uint64) |  | unix timestamp before which the contract call must be executed. enforced by the proxy contract. |






<a name="axelarcork.v1.AxelarCorkResult"></a>

### AxelarCorkResult



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `cork` | [AxelarCork](#axelarcork.v1.AxelarCork) |  |  |
| `block_height` | [uint64](#uint64) |  |  |
| `approved` | [bool](#bool) |  |  |
| `approval_percentage` | [string](#string) |  |  |






<a name="axelarcork.v1.AxelarCorkResults"></a>

### AxelarCorkResults



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `cork_results` | [AxelarCorkResult](#axelarcork.v1.AxelarCorkResult) | repeated |  |






<a name="axelarcork.v1.AxelarUpgradeData"></a>

### AxelarUpgradeData
Represents a proxy contract upgrade approved by governance with a delay in
execution in case of an error.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `chain_id` | [uint64](#uint64) |  |  |
| `payload` | [bytes](#bytes) |  |  |
| `executable_height_threshold` | [int64](#int64) |  |  |






<a name="axelarcork.v1.CellarIDSet"></a>

### CellarIDSet



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `chain_id` | [uint64](#uint64) |  |  |
| `ids` | [string](#string) | repeated |  |






<a name="axelarcork.v1.ChainConfiguration"></a>

### ChainConfiguration



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `name` | [string](#string) |  |  |
| `id` | [uint64](#uint64) |  |  |
| `proxy_address` | [string](#string) |  |  |
| `bridge_fees` | [cosmos.base.v1beta1.Coin](#cosmos.base.v1beta1.Coin) | repeated | pure token transfers have a fixed fee deducted from the amount sent in the ICS-20 message depending on the asset and destination chain they can be calculated here: https://docs.axelar.dev/resources/mainnet#cross-chain-relayer-gas-fee |






<a name="axelarcork.v1.ChainConfigurations"></a>

### ChainConfigurations



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `configurations` | [ChainConfiguration](#axelarcork.v1.ChainConfiguration) | repeated |  |






<a name="axelarcork.v1.ScheduledAxelarCork"></a>

### ScheduledAxelarCork



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `cork` | [AxelarCork](#axelarcork.v1.AxelarCork) |  |  |
| `block_height` | [uint64](#uint64) |  |  |
| `validator` | [string](#string) |  |  |
| `id` | [string](#string) |  |  |






<a name="axelarcork.v1.ScheduledAxelarCorks"></a>

### ScheduledAxelarCorks



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `scheduled_corks` | [ScheduledAxelarCork](#axelarcork.v1.ScheduledAxelarCork) | repeated |  |





 <!-- end messages -->

 <!-- end enums -->

 <!-- end HasExtensions -->

 <!-- end services -->



<a name="axelarcork/v1/event.proto"></a>
<p align="right"><a href="#top">Top</a></p>

## axelarcork/v1/event.proto



<a name="axelarcork.v1.ScheduleCorkEvent"></a>

### ScheduleCorkEvent



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `signer` | [string](#string) |  |  |
| `validator` | [string](#string) |  |  |
| `cork` | [string](#string) |  |  |
| `block_height` | [uint64](#uint64) |  |  |
| `chain_id` | [uint64](#uint64) |  |  |





 <!-- end messages -->

 <!-- end enums -->

 <!-- end HasExtensions -->

 <!-- end services -->



<a name="axelarcork/v1/genesis.proto"></a>
<p align="right"><a href="#top">Top</a></p>

## axelarcork/v1/genesis.proto



<a name="axelarcork.v1.GenesisState"></a>

### GenesisState
GenesisState - all cork state that must be provided at genesis


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `params` | [Params](#axelarcork.v1.Params) |  |  |
| `chain_configurations` | [ChainConfigurations](#axelarcork.v1.ChainConfigurations) |  |  |
| `cellar_ids` | [CellarIDSet](#axelarcork.v1.CellarIDSet) | repeated |  |
| `scheduled_corks` | [ScheduledAxelarCorks](#axelarcork.v1.ScheduledAxelarCorks) |  |  |
| `cork_results` | [AxelarCorkResults](#axelarcork.v1.AxelarCorkResults) |  |  |
| `axelar_contract_call_nonces` | [AxelarContractCallNonce](#axelarcork.v1.AxelarContractCallNonce) | repeated |  |
| `axelar_upgrade_data` | [AxelarUpgradeData](#axelarcork.v1.AxelarUpgradeData) | repeated |  |






<a name="axelarcork.v1.Params"></a>

### Params



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `enabled` | [bool](#bool) |  |  |
| `ibc_channel` | [string](#string) |  |  |
| `ibc_port` | [string](#string) |  |  |
| `gmp_account` | [string](#string) |  |  |
| `executor_account` | [string](#string) |  |  |
| `timeout_duration` | [uint64](#uint64) |  |  |
| `cork_timeout_blocks` | [uint64](#uint64) |  |  |





 <!-- end messages -->

 <!-- end enums -->

 <!-- end HasExtensions -->

 <!-- end services -->



<a name="axelarcork/v1/proposal.proto"></a>
<p align="right"><a href="#top">Top</a></p>

## axelarcork/v1/proposal.proto



<a name="axelarcork.v1.AddAxelarManagedCellarIDsProposal"></a>

### AddAxelarManagedCellarIDsProposal



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `title` | [string](#string) |  |  |
| `description` | [string](#string) |  |  |
| `chain_id` | [uint64](#uint64) |  |  |
| `cellar_ids` | [CellarIDSet](#axelarcork.v1.CellarIDSet) |  |  |
| `publisher_domain` | [string](#string) |  |  |






<a name="axelarcork.v1.AddAxelarManagedCellarIDsProposalWithDeposit"></a>

### AddAxelarManagedCellarIDsProposalWithDeposit
AddAxelarManagedCellarIDsProposalWithDeposit is a specific definition for CLI commands


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `title` | [string](#string) |  |  |
| `description` | [string](#string) |  |  |
| `chain_id` | [uint64](#uint64) |  |  |
| `cellar_ids` | [string](#string) | repeated |  |
| `publisher_domain` | [string](#string) |  |  |
| `deposit` | [string](#string) |  |  |






<a name="axelarcork.v1.AddChainConfigurationProposal"></a>

### AddChainConfigurationProposal



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `title` | [string](#string) |  |  |
| `description` | [string](#string) |  |  |
| `chain_configuration` | [ChainConfiguration](#axelarcork.v1.ChainConfiguration) |  |  |






<a name="axelarcork.v1.AddChainConfigurationProposalWithDeposit"></a>

### AddChainConfigurationProposalWithDeposit



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `title` | [string](#string) |  |  |
| `description` | [string](#string) |  |  |
| `chain_configuration` | [ChainConfiguration](#axelarcork.v1.ChainConfiguration) |  |  |
| `deposit` | [string](#string) |  |  |






<a name="axelarcork.v1.AxelarCommunityPoolSpendProposal"></a>

### AxelarCommunityPoolSpendProposal



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `title` | [string](#string) |  |  |
| `description` | [string](#string) |  |  |
| `recipient` | [string](#string) |  |  |
| `chain_id` | [uint64](#uint64) |  |  |
| `amount` | [cosmos.base.v1beta1.Coin](#cosmos.base.v1beta1.Coin) |  |  |






<a name="axelarcork.v1.AxelarCommunityPoolSpendProposalForCLI"></a>

### AxelarCommunityPoolSpendProposalForCLI
This format of the Axelar community spend Ethereum proposal is specifically for
the CLI to allow simple text serialization.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `title` | [string](#string) |  |  |
| `description` | [string](#string) |  |  |
| `recipient` | [string](#string) |  |  |
| `chain_id` | [uint64](#uint64) |  |  |
| `amount` | [string](#string) |  |  |
| `deposit` | [string](#string) |  |  |






<a name="axelarcork.v1.AxelarScheduledCorkProposal"></a>

### AxelarScheduledCorkProposal



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `title` | [string](#string) |  |  |
| `description` | [string](#string) |  |  |
| `block_height` | [uint64](#uint64) |  |  |
| `chain_id` | [uint64](#uint64) |  |  |
| `target_contract_address` | [string](#string) |  |  |
| `contract_call_proto_json` | [string](#string) |  | The JSON representation of a ScheduleRequest defined in the Steward protos

Example: The following is the JSON form of a ScheduleRequest containing a steward.v2.cellar_v1.TrustPosition message, which maps to the `trustPosition(address)` function of the the V1 Cellar contract.

{ "cellar_id": "0x1234567890000000000000000000000000000000", "cellar_v1": { "trust_position": { "erc20_address": "0x1234567890000000000000000000000000000000" } }, "block_height": 1000000 }

You can use the Steward CLI to generate the required JSON rather than constructing it by hand https://github.com/peggyjv/steward |
| `deadline` | [uint64](#uint64) |  | unix timestamp before which the contract call must be executed. enforced by the Axelar proxy contract |






<a name="axelarcork.v1.AxelarScheduledCorkProposalWithDeposit"></a>

### AxelarScheduledCorkProposalWithDeposit
AxelarScheduledCorkProposalWithDeposit is a specific definition for CLI commands


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `title` | [string](#string) |  |  |
| `description` | [string](#string) |  |  |
| `block_height` | [uint64](#uint64) |  |  |
| `chain_id` | [uint64](#uint64) |  |  |
| `target_contract_address` | [string](#string) |  |  |
| `contract_call_proto_json` | [string](#string) |  |  |
| `deadline` | [uint64](#uint64) |  |  |
| `deposit` | [string](#string) |  |  |






<a name="axelarcork.v1.CancelAxelarProxyContractUpgradeProposal"></a>

### CancelAxelarProxyContractUpgradeProposal



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `title` | [string](#string) |  |  |
| `description` | [string](#string) |  |  |
| `chain_id` | [uint64](#uint64) |  |  |






<a name="axelarcork.v1.CancelAxelarProxyContractUpgradeProposalWithDeposit"></a>

### CancelAxelarProxyContractUpgradeProposalWithDeposit



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `title` | [string](#string) |  |  |
| `description` | [string](#string) |  |  |
| `chain_id` | [uint64](#uint64) |  |  |
| `deposit` | [string](#string) |  |  |






<a name="axelarcork.v1.RemoveAxelarManagedCellarIDsProposal"></a>

### RemoveAxelarManagedCellarIDsProposal



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `title` | [string](#string) |  |  |
| `description` | [string](#string) |  |  |
| `chain_id` | [uint64](#uint64) |  |  |
| `cellar_ids` | [CellarIDSet](#axelarcork.v1.CellarIDSet) |  |  |






<a name="axelarcork.v1.RemoveAxelarManagedCellarIDsProposalWithDeposit"></a>

### RemoveAxelarManagedCellarIDsProposalWithDeposit
RemoveAxelarManagedCellarIDsProposalWithDeposit is a specific definition for CLI commands


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `title` | [string](#string) |  |  |
| `description` | [string](#string) |  |  |
| `chain_id` | [uint64](#uint64) |  |  |
| `cellar_ids` | [string](#string) | repeated |  |
| `deposit` | [string](#string) |  |  |






<a name="axelarcork.v1.RemoveChainConfigurationProposal"></a>

### RemoveChainConfigurationProposal



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `title` | [string](#string) |  |  |
| `description` | [string](#string) |  |  |
| `chain_id` | [uint64](#uint64) |  |  |






<a name="axelarcork.v1.RemoveChainConfigurationProposalWithDeposit"></a>

### RemoveChainConfigurationProposalWithDeposit



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `title` | [string](#string) |  |  |
| `description` | [string](#string) |  |  |
| `chain_id` | [uint64](#uint64) |  |  |
| `deposit` | [string](#string) |  |  |






<a name="axelarcork.v1.UpgradeAxelarProxyContractProposal"></a>

### UpgradeAxelarProxyContractProposal



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `title` | [string](#string) |  |  |
| `description` | [string](#string) |  |  |
| `chain_id` | [uint64](#uint64) |  |  |
| `new_proxy_address` | [string](#string) |  |  |






<a name="axelarcork.v1.UpgradeAxelarProxyContractProposalWithDeposit"></a>

### UpgradeAxelarProxyContractProposalWithDeposit



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `title` | [string](#string) |  |  |
| `description` | [string](#string) |  |  |
| `chain_id` | [uint64](#uint64) |  |  |
| `new_proxy_address` | [string](#string) |  |  |
| `deposit` | [string](#string) |  |  |





 <!-- end messages -->

 <!-- end enums -->

 <!-- end HasExtensions -->

 <!-- end services -->



<a name="axelarcork/v1/query.proto"></a>
<p align="right"><a href="#top">Top</a></p>

## axelarcork/v1/query.proto



<a name="axelarcork.v1.QueryAxelarContractCallNoncesRequest"></a>

### QueryAxelarContractCallNoncesRequest







<a name="axelarcork.v1.QueryAxelarContractCallNoncesResponse"></a>

### QueryAxelarContractCallNoncesResponse



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `contract_call_nonces` | [AxelarContractCallNonce](#axelarcork.v1.AxelarContractCallNonce) | repeated |  |






<a name="axelarcork.v1.QueryAxelarProxyUpgradeDataRequest"></a>

### QueryAxelarProxyUpgradeDataRequest







<a name="axelarcork.v1.QueryAxelarProxyUpgradeDataResponse"></a>

### QueryAxelarProxyUpgradeDataResponse



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `proxy_upgrade_data` | [AxelarUpgradeData](#axelarcork.v1.AxelarUpgradeData) | repeated |  |






<a name="axelarcork.v1.QueryCellarIDsByChainIDRequest"></a>

### QueryCellarIDsByChainIDRequest
QueryCellarIDsByChainIDRequest is the request type for Query/QueryCellarIDsByChainID gRPC method.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `chain_id` | [uint64](#uint64) |  |  |






<a name="axelarcork.v1.QueryCellarIDsByChainIDResponse"></a>

### QueryCellarIDsByChainIDResponse
QueryCellarIDsByChainIDResponse is the response type for Query/QueryCellarIDsByChainID gRPC method.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `cellar_ids` | [string](#string) | repeated |  |






<a name="axelarcork.v1.QueryCellarIDsRequest"></a>

### QueryCellarIDsRequest
QueryCellarIDs is the request type for Query/QueryCellarIDs gRPC method.






<a name="axelarcork.v1.QueryCellarIDsResponse"></a>

### QueryCellarIDsResponse
QueryCellarIDsResponse is the response type for Query/QueryCellarIDs gRPC method.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `cellar_ids` | [CellarIDSet](#axelarcork.v1.CellarIDSet) | repeated |  |






<a name="axelarcork.v1.QueryChainConfigurationsRequest"></a>

### QueryChainConfigurationsRequest







<a name="axelarcork.v1.QueryChainConfigurationsResponse"></a>

### QueryChainConfigurationsResponse



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `configurations` | [ChainConfiguration](#axelarcork.v1.ChainConfiguration) | repeated |  |






<a name="axelarcork.v1.QueryCorkResultRequest"></a>

### QueryCorkResultRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `id` | [string](#string) |  |  |
| `chain_id` | [uint64](#uint64) |  |  |






<a name="axelarcork.v1.QueryCorkResultResponse"></a>

### QueryCorkResultResponse



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `corkResult` | [AxelarCorkResult](#axelarcork.v1.AxelarCorkResult) |  |  |






<a name="axelarcork.v1.QueryCorkResultsRequest"></a>

### QueryCorkResultsRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `chain_id` | [uint64](#uint64) |  |  |






<a name="axelarcork.v1.QueryCorkResultsResponse"></a>

### QueryCorkResultsResponse



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `corkResults` | [AxelarCorkResult](#axelarcork.v1.AxelarCorkResult) | repeated |  |






<a name="axelarcork.v1.QueryParamsRequest"></a>

### QueryParamsRequest
QueryParamsRequest is the request type for the Query/Params gRPC method.






<a name="axelarcork.v1.QueryParamsResponse"></a>

### QueryParamsResponse
QueryParamsRequest is the response type for the Query/Params gRPC method.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `params` | [Params](#axelarcork.v1.Params) |  | allocation parameters |






<a name="axelarcork.v1.QueryScheduledBlockHeightsRequest"></a>

### QueryScheduledBlockHeightsRequest
QueryScheduledBlockHeightsRequest


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `chain_id` | [uint64](#uint64) |  |  |






<a name="axelarcork.v1.QueryScheduledBlockHeightsResponse"></a>

### QueryScheduledBlockHeightsResponse
QueryScheduledBlockHeightsResponse


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `block_heights` | [uint64](#uint64) | repeated |  |






<a name="axelarcork.v1.QueryScheduledCorksByBlockHeightRequest"></a>

### QueryScheduledCorksByBlockHeightRequest
QueryScheduledCorksByBlockHeightRequest


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `block_height` | [uint64](#uint64) |  |  |
| `chain_id` | [uint64](#uint64) |  |  |






<a name="axelarcork.v1.QueryScheduledCorksByBlockHeightResponse"></a>

### QueryScheduledCorksByBlockHeightResponse
QueryScheduledCorksByBlockHeightResponse


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `corks` | [ScheduledAxelarCork](#axelarcork.v1.ScheduledAxelarCork) | repeated |  |






<a name="axelarcork.v1.QueryScheduledCorksByIDRequest"></a>

### QueryScheduledCorksByIDRequest
QueryScheduledCorksByIDRequest


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `id` | [string](#string) |  |  |
| `chain_id` | [uint64](#uint64) |  |  |






<a name="axelarcork.v1.QueryScheduledCorksByIDResponse"></a>

### QueryScheduledCorksByIDResponse
QueryScheduledCorksByIDResponse


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `corks` | [ScheduledAxelarCork](#axelarcork.v1.ScheduledAxelarCork) | repeated |  |






<a name="axelarcork.v1.QueryScheduledCorksRequest"></a>

### QueryScheduledCorksRequest
QueryScheduledCorksRequest


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `chain_id` | [uint64](#uint64) |  |  |






<a name="axelarcork.v1.QueryScheduledCorksResponse"></a>

### QueryScheduledCorksResponse
QueryScheduledCorksResponse


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `corks` | [ScheduledAxelarCork](#axelarcork.v1.ScheduledAxelarCork) | repeated |  |






<a name="axelarcork.v1.QueryWinningAxelarCorkRequest"></a>

### QueryWinningAxelarCorkRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `chain_id` | [uint64](#uint64) |  |  |
| `contract_address` | [string](#string) |  |  |






<a name="axelarcork.v1.QueryWinningAxelarCorkResponse"></a>

### QueryWinningAxelarCorkResponse



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `cork` | [AxelarCork](#axelarcork.v1.AxelarCork) |  |  |
| `block_height` | [uint64](#uint64) |  |  |





 <!-- end messages -->

 <!-- end enums -->

 <!-- end HasExtensions -->


<a name="axelarcork.v1.Query"></a>

### Query
Query defines the gRPC query service for the cork module.

| Method Name | Request Type | Response Type | Description | HTTP Verb | Endpoint |
| ----------- | ------------ | ------------- | ------------| ------- | -------- |
| `QueryParams` | [QueryParamsRequest](#axelarcork.v1.QueryParamsRequest) | [QueryParamsResponse](#axelarcork.v1.QueryParamsResponse) | QueryParams queries the axelar cork module parameters. | GET|/sommelier/cork/v1/params|
| `QueryCellarIDs` | [QueryCellarIDsRequest](#axelarcork.v1.QueryCellarIDsRequest) | [QueryCellarIDsResponse](#axelarcork.v1.QueryCellarIDsResponse) | QueryCellarIDs queries approved cellar ids of all supported chains | GET|/sommelier/axelarcork/v1/cellar_ids|
| `QueryCellarIDsByChainID` | [QueryCellarIDsByChainIDRequest](#axelarcork.v1.QueryCellarIDsByChainIDRequest) | [QueryCellarIDsByChainIDResponse](#axelarcork.v1.QueryCellarIDsByChainIDResponse) | QueryCellarIDsByChainID returns all cellars and current tick ranges | GET|/sommelier/axelarcork/v1/cellar_ids_by_chain_id|
| `QueryScheduledCorks` | [QueryScheduledCorksRequest](#axelarcork.v1.QueryScheduledCorksRequest) | [QueryScheduledCorksResponse](#axelarcork.v1.QueryScheduledCorksResponse) | QueryScheduledCorks returns all scheduled corks | GET|/sommelier/axelarcork/v1/scheduled_corks|
| `QueryScheduledBlockHeights` | [QueryScheduledBlockHeightsRequest](#axelarcork.v1.QueryScheduledBlockHeightsRequest) | [QueryScheduledBlockHeightsResponse](#axelarcork.v1.QueryScheduledBlockHeightsResponse) | QueryScheduledBlockHeights returns all scheduled block heights | GET|/sommelier/axelarcork/v1/scheduled_block_heights|
| `QueryScheduledCorksByBlockHeight` | [QueryScheduledCorksByBlockHeightRequest](#axelarcork.v1.QueryScheduledCorksByBlockHeightRequest) | [QueryScheduledCorksByBlockHeightResponse](#axelarcork.v1.QueryScheduledCorksByBlockHeightResponse) | QueryScheduledCorks returns all scheduled corks at a block height | GET|/sommelier/axelarcork/v1/scheduled_corks_by_block_height/{block_height}|
| `QueryScheduledCorksByID` | [QueryScheduledCorksByIDRequest](#axelarcork.v1.QueryScheduledCorksByIDRequest) | [QueryScheduledCorksByIDResponse](#axelarcork.v1.QueryScheduledCorksByIDResponse) | QueryScheduledCorks returns all scheduled corks with the specified ID | GET|/sommelier/axelarcork/v1/scheduled_corks_by_id/{id}|
| `QueryCorkResult` | [QueryCorkResultRequest](#axelarcork.v1.QueryCorkResultRequest) | [QueryCorkResultResponse](#axelarcork.v1.QueryCorkResultResponse) |  | GET|/sommelier/axelarcork/v1/cork_results/{id}|
| `QueryCorkResults` | [QueryCorkResultsRequest](#axelarcork.v1.QueryCorkResultsRequest) | [QueryCorkResultsResponse](#axelarcork.v1.QueryCorkResultsResponse) |  | GET|/sommelier/axelarcork/v1/cork_results|
| `QueryChainConfigurations` | [QueryChainConfigurationsRequest](#axelarcork.v1.QueryChainConfigurationsRequest) | [QueryChainConfigurationsResponse](#axelarcork.v1.QueryChainConfigurationsResponse) |  | GET|/sommelier/axelarcork/v1/chain_configurations|
| `QueryAxelarContractCallNonces` | [QueryAxelarContractCallNoncesRequest](#axelarcork.v1.QueryAxelarContractCallNoncesRequest) | [QueryAxelarContractCallNoncesResponse](#axelarcork.v1.QueryAxelarContractCallNoncesResponse) |  | GET|/sommelier/axelarcork/v1/contract_call_nonces|
| `QueryAxelarProxyUpgradeData` | [QueryAxelarProxyUpgradeDataRequest](#axelarcork.v1.QueryAxelarProxyUpgradeDataRequest) | [QueryAxelarProxyUpgradeDataResponse](#axelarcork.v1.QueryAxelarProxyUpgradeDataResponse) |  | GET|/sommelier/axelarcork/v1/proxy_upgrade_data|
| `QueryWinningAxelarCork` | [QueryWinningAxelarCorkRequest](#axelarcork.v1.QueryWinningAxelarCorkRequest) | [QueryWinningAxelarCorkResponse](#axelarcork.v1.QueryWinningAxelarCorkResponse) |  | GET|/sommelier/axelarcork/v1/winning_cork|

 <!-- end services -->



<a name="axelarcork/v1/tx.proto"></a>
<p align="right"><a href="#top">Top</a></p>

## axelarcork/v1/tx.proto



<a name="axelarcork.v1.MsgBumpAxelarCorkGasRequest"></a>

### MsgBumpAxelarCorkGasRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `signer` | [string](#string) |  |  |
| `token` | [cosmos.base.v1beta1.Coin](#cosmos.base.v1beta1.Coin) |  |  |
| `message_id` | [string](#string) |  |  |






<a name="axelarcork.v1.MsgBumpAxelarCorkGasResponse"></a>

### MsgBumpAxelarCorkGasResponse







<a name="axelarcork.v1.MsgCancelAxelarCorkRequest"></a>

### MsgCancelAxelarCorkRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `signer` | [string](#string) |  |  |
| `chain_id` | [uint64](#uint64) |  |  |
| `target_contract_address` | [string](#string) |  |  |






<a name="axelarcork.v1.MsgCancelAxelarCorkResponse"></a>

### MsgCancelAxelarCorkResponse







<a name="axelarcork.v1.MsgRelayAxelarCorkRequest"></a>

### MsgRelayAxelarCorkRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `signer` | [string](#string) |  |  |
| `token` | [cosmos.base.v1beta1.Coin](#cosmos.base.v1beta1.Coin) |  |  |
| `fee` | [uint64](#uint64) |  |  |
| `chain_id` | [uint64](#uint64) |  |  |
| `target_contract_address` | [string](#string) |  |  |






<a name="axelarcork.v1.MsgRelayAxelarCorkResponse"></a>

### MsgRelayAxelarCorkResponse







<a name="axelarcork.v1.MsgRelayAxelarProxyUpgradeRequest"></a>

### MsgRelayAxelarProxyUpgradeRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `signer` | [string](#string) |  |  |
| `token` | [cosmos.base.v1beta1.Coin](#cosmos.base.v1beta1.Coin) |  |  |
| `fee` | [uint64](#uint64) |  |  |
| `chain_id` | [uint64](#uint64) |  |  |






<a name="axelarcork.v1.MsgRelayAxelarProxyUpgradeResponse"></a>

### MsgRelayAxelarProxyUpgradeResponse







<a name="axelarcork.v1.MsgScheduleAxelarCorkRequest"></a>

### MsgScheduleAxelarCorkRequest
MsgScheduleCorkRequest - sdk.Msg for scheduling a cork request for on or after a specific block height


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `cork` | [AxelarCork](#axelarcork.v1.AxelarCork) |  | the scheduled cork |
| `chain_id` | [uint64](#uint64) |  | the chain id |
| `block_height` | [uint64](#uint64) |  | the block height that must be reached |
| `signer` | [string](#string) |  | signer account address |






<a name="axelarcork.v1.MsgScheduleAxelarCorkResponse"></a>

### MsgScheduleAxelarCorkResponse



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `id` | [string](#string) |  | cork ID |





 <!-- end messages -->

 <!-- end enums -->

 <!-- end HasExtensions -->


<a name="axelarcork.v1.Msg"></a>

### Msg
MsgService defines the msgs that the cork module handles

| Method Name | Request Type | Response Type | Description | HTTP Verb | Endpoint |
| ----------- | ------------ | ------------- | ------------| ------- | -------- |
| `ScheduleCork` | [MsgScheduleAxelarCorkRequest](#axelarcork.v1.MsgScheduleAxelarCorkRequest) | [MsgScheduleAxelarCorkResponse](#axelarcork.v1.MsgScheduleAxelarCorkResponse) |  | |
| `RelayCork` | [MsgRelayAxelarCorkRequest](#axelarcork.v1.MsgRelayAxelarCorkRequest) | [MsgRelayAxelarCorkResponse](#axelarcork.v1.MsgRelayAxelarCorkResponse) |  | |
| `BumpCorkGas` | [MsgBumpAxelarCorkGasRequest](#axelarcork.v1.MsgBumpAxelarCorkGasRequest) | [MsgBumpAxelarCorkGasResponse](#axelarcork.v1.MsgBumpAxelarCorkGasResponse) |  | |
| `CancelScheduledCork` | [MsgCancelAxelarCorkRequest](#axelarcork.v1.MsgCancelAxelarCorkRequest) | [MsgCancelAxelarCorkResponse](#axelarcork.v1.MsgCancelAxelarCorkResponse) |  | |

 <!-- end services -->



<a name="cellarfees/v1/cellarfees.proto"></a>
<p align="right"><a href="#top">Top</a></p>

## cellarfees/v1/cellarfees.proto



<a name="cellarfees.v1.FeeAccrualCounter"></a>

### FeeAccrualCounter



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `denom` | [string](#string) |  |  |
| `count` | [uint64](#uint64) |  |  |






<a name="cellarfees.v1.FeeAccrualCounters"></a>

### FeeAccrualCounters



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `counters` | [FeeAccrualCounter](#cellarfees.v1.FeeAccrualCounter) | repeated |  |





 <!-- end messages -->

 <!-- end enums -->

 <!-- end HasExtensions -->

 <!-- end services -->



<a name="cellarfees/v1/params.proto"></a>
<p align="right"><a href="#top">Top</a></p>

## cellarfees/v1/params.proto



<a name="cellarfees.v1.Params"></a>

### Params
Params defines the parameters for the module.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `fee_accrual_auction_threshold` | [uint64](#uint64) |  | The number of fee accruals after which an auction should be started |
| `reward_emission_period` | [uint64](#uint64) |  | Emission rate factor. Specifically, the number of blocks over which to distribute some amount of staking rewards. |
| `initial_price_decrease_rate` | [string](#string) |  | The initial rate at which auctions should decrease their denom's price in SOMM |
| `price_decrease_block_interval` | [uint64](#uint64) |  | Number of blocks between auction price decreases |
| `auction_interval` | [uint64](#uint64) |  | The interval between starting auctions |





 <!-- end messages -->

 <!-- end enums -->

 <!-- end HasExtensions -->

 <!-- end services -->



<a name="cellarfees/v1/genesis.proto"></a>
<p align="right"><a href="#top">Top</a></p>

## cellarfees/v1/genesis.proto



<a name="cellarfees.v1.GenesisState"></a>

### GenesisState
GenesisState defines the cellarfees module's genesis state.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `params` | [Params](#cellarfees.v1.Params) |  |  |
| `fee_accrual_counters` | [FeeAccrualCounters](#cellarfees.v1.FeeAccrualCounters) |  |  |
| `last_reward_supply_peak` | [string](#string) |  |  |





 <!-- end messages -->

 <!-- end enums -->

 <!-- end HasExtensions -->

 <!-- end services -->



<a name="cellarfees/v1/query.proto"></a>
<p align="right"><a href="#top">Top</a></p>

## cellarfees/v1/query.proto



<a name="cellarfees.v1.QueryAPYRequest"></a>

### QueryAPYRequest







<a name="cellarfees.v1.QueryAPYResponse"></a>

### QueryAPYResponse



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `apy` | [string](#string) |  |  |






<a name="cellarfees.v1.QueryFeeAccrualCountersRequest"></a>

### QueryFeeAccrualCountersRequest







<a name="cellarfees.v1.QueryFeeAccrualCountersResponse"></a>

### QueryFeeAccrualCountersResponse



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `fee_accrual_counters` | [FeeAccrualCounters](#cellarfees.v1.FeeAccrualCounters) |  |  |






<a name="cellarfees.v1.QueryLastRewardSupplyPeakRequest"></a>

### QueryLastRewardSupplyPeakRequest







<a name="cellarfees.v1.QueryLastRewardSupplyPeakResponse"></a>

### QueryLastRewardSupplyPeakResponse



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `last_reward_supply_peak` | [string](#string) |  |  |






<a name="cellarfees.v1.QueryModuleAccountsRequest"></a>

### QueryModuleAccountsRequest







<a name="cellarfees.v1.QueryModuleAccountsResponse"></a>

### QueryModuleAccountsResponse



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `fees_address` | [string](#string) |  |  |






<a name="cellarfees.v1.QueryParamsRequest"></a>

### QueryParamsRequest







<a name="cellarfees.v1.QueryParamsResponse"></a>

### QueryParamsResponse



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `params` | [Params](#cellarfees.v1.Params) |  |  |





 <!-- end messages -->

 <!-- end enums -->

 <!-- end HasExtensions -->


<a name="cellarfees.v1.Query"></a>

### Query


| Method Name | Request Type | Response Type | Description | HTTP Verb | Endpoint |
| ----------- | ------------ | ------------- | ------------| ------- | -------- |
| `QueryParams` | [QueryParamsRequest](#cellarfees.v1.QueryParamsRequest) | [QueryParamsResponse](#cellarfees.v1.QueryParamsResponse) |  | GET|/sommelier/cellarfees/v1/params|
| `QueryModuleAccounts` | [QueryModuleAccountsRequest](#cellarfees.v1.QueryModuleAccountsRequest) | [QueryModuleAccountsResponse](#cellarfees.v1.QueryModuleAccountsResponse) |  | GET|/sommelier/cellarfees/v1/module_accounts|
| `QueryLastRewardSupplyPeak` | [QueryLastRewardSupplyPeakRequest](#cellarfees.v1.QueryLastRewardSupplyPeakRequest) | [QueryLastRewardSupplyPeakResponse](#cellarfees.v1.QueryLastRewardSupplyPeakResponse) |  | GET|/sommelier/cellarfees/v1/last_reward_supply_peak|
| `QueryFeeAccrualCounters` | [QueryFeeAccrualCountersRequest](#cellarfees.v1.QueryFeeAccrualCountersRequest) | [QueryFeeAccrualCountersResponse](#cellarfees.v1.QueryFeeAccrualCountersResponse) |  | GET|/sommelier/cellarfees/v1/fee_accrual_counters|
| `QueryAPY` | [QueryAPYRequest](#cellarfees.v1.QueryAPYRequest) | [QueryAPYResponse](#cellarfees.v1.QueryAPYResponse) |  | GET|/sommelier/cellarfees/v1/apy|

 <!-- end services -->



<a name="cork/v2/cork.proto"></a>
<p align="right"><a href="#top">Top</a></p>

## cork/v2/cork.proto



<a name="cork.v2.CellarIDSet"></a>

### CellarIDSet



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `ids` | [string](#string) | repeated |  |






<a name="cork.v2.Cork"></a>

### Cork



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `encoded_contract_call` | [bytes](#bytes) |  | call body containing the ABI encoded bytes to send to the contract |
| `target_contract_address` | [string](#string) |  | address of the contract to send the call |






<a name="cork.v2.CorkResult"></a>

### CorkResult



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `cork` | [Cork](#cork.v2.Cork) |  |  |
| `block_height` | [uint64](#uint64) |  |  |
| `approved` | [bool](#bool) |  |  |
| `approval_percentage` | [string](#string) |  |  |






<a name="cork.v2.ScheduledCork"></a>

### ScheduledCork



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `cork` | [Cork](#cork.v2.Cork) |  |  |
| `block_height` | [uint64](#uint64) |  |  |
| `validator` | [string](#string) |  |  |
| `id` | [bytes](#bytes) |  |  |





 <!-- end messages -->

 <!-- end enums -->

 <!-- end HasExtensions -->

 <!-- end services -->



<a name="cork/v2/genesis.proto"></a>
<p align="right"><a href="#top">Top</a></p>

## cork/v2/genesis.proto



<a name="cork.v2.GenesisState"></a>

### GenesisState
GenesisState - all cork state that must be provided at genesis


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `params` | [Params](#cork.v2.Params) |  |  |
| `cellar_ids` | [CellarIDSet](#cork.v2.CellarIDSet) |  |  |
| `invalidation_nonce` | [uint64](#uint64) |  |  |
| `scheduled_corks` | [ScheduledCork](#cork.v2.ScheduledCork) | repeated |  |
| `cork_results` | [CorkResult](#cork.v2.CorkResult) | repeated |  |






<a name="cork.v2.Params"></a>

### Params
Params cork parameters


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `vote_threshold` | [string](#string) |  | Deprecated VoteThreshold defines the percentage of bonded stake required to vote for a scheduled cork to be approved |
| `max_corks_per_validator` | [uint64](#uint64) |  |  |





 <!-- end messages -->

 <!-- end enums -->

 <!-- end HasExtensions -->

 <!-- end services -->



<a name="cork/v2/proposal.proto"></a>
<p align="right"><a href="#top">Top</a></p>

## cork/v2/proposal.proto



<a name="cork.v2.AddManagedCellarIDsProposal"></a>

### AddManagedCellarIDsProposal



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `title` | [string](#string) |  |  |
| `description` | [string](#string) |  |  |
| `cellar_ids` | [CellarIDSet](#cork.v2.CellarIDSet) |  |  |
| `publisher_domain` | [string](#string) |  |  |






<a name="cork.v2.AddManagedCellarIDsProposalWithDeposit"></a>

### AddManagedCellarIDsProposalWithDeposit
AddManagedCellarIDsProposalWithDeposit is a specific definition for CLI commands


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `title` | [string](#string) |  |  |
| `description` | [string](#string) |  |  |
| `cellar_ids` | [string](#string) | repeated |  |
| `publisher_domain` | [string](#string) |  |  |
| `deposit` | [string](#string) |  |  |






<a name="cork.v2.RemoveManagedCellarIDsProposal"></a>

### RemoveManagedCellarIDsProposal



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `title` | [string](#string) |  |  |
| `description` | [string](#string) |  |  |
| `cellar_ids` | [CellarIDSet](#cork.v2.CellarIDSet) |  |  |






<a name="cork.v2.RemoveManagedCellarIDsProposalWithDeposit"></a>

### RemoveManagedCellarIDsProposalWithDeposit
RemoveManagedCellarIDsProposalWithDeposit is a specific definition for CLI commands


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `title` | [string](#string) |  |  |
| `description` | [string](#string) |  |  |
| `cellar_ids` | [string](#string) | repeated |  |
| `deposit` | [string](#string) |  |  |






<a name="cork.v2.ScheduledCorkProposal"></a>

### ScheduledCorkProposal



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `title` | [string](#string) |  |  |
| `description` | [string](#string) |  |  |
| `block_height` | [uint64](#uint64) |  |  |
| `target_contract_address` | [string](#string) |  |  |
| `contract_call_proto_json` | [string](#string) |  | The JSON representation of a ScheduleRequest defined in the Steward protos

Example: The following is the JSON form of a ScheduleRequest containing a steward.v2.cellar_v1.TrustPosition message, which maps to the `trustPosition(address)` function of the the V1 Cellar contract.

{ "cellar_id": "0x1234567890000000000000000000000000000000", "cellar_v1": { "trust_position": { "erc20_address": "0x1234567890000000000000000000000000000000" } }, "block_height": 1000000 }

You can use the Steward CLI to generate the required JSON rather than constructing it by hand https://github.com/peggyjv/steward |






<a name="cork.v2.ScheduledCorkProposalWithDeposit"></a>

### ScheduledCorkProposalWithDeposit
ScheduledCorkProposalWithDeposit is a specific definition for CLI commands


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `title` | [string](#string) |  |  |
| `description` | [string](#string) |  |  |
| `block_height` | [uint64](#uint64) |  |  |
| `target_contract_address` | [string](#string) |  |  |
| `contract_call_proto_json` | [string](#string) |  |  |
| `deposit` | [string](#string) |  |  |





 <!-- end messages -->

 <!-- end enums -->

 <!-- end HasExtensions -->

 <!-- end services -->



<a name="cork/v2/query.proto"></a>
<p align="right"><a href="#top">Top</a></p>

## cork/v2/query.proto



<a name="cork.v2.QueryCellarIDsRequest"></a>

### QueryCellarIDsRequest
QueryCellarIDsRequest is the request type for Query/QueryCellarIDs gRPC method.






<a name="cork.v2.QueryCellarIDsResponse"></a>

### QueryCellarIDsResponse
QueryCellarIDsResponse is the response type for Query/QueryCellars gRPC method.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `cellar_ids` | [string](#string) | repeated |  |






<a name="cork.v2.QueryCorkResultRequest"></a>

### QueryCorkResultRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `id` | [string](#string) |  |  |






<a name="cork.v2.QueryCorkResultResponse"></a>

### QueryCorkResultResponse



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `corkResult` | [CorkResult](#cork.v2.CorkResult) |  |  |






<a name="cork.v2.QueryCorkResultsRequest"></a>

### QueryCorkResultsRequest







<a name="cork.v2.QueryCorkResultsResponse"></a>

### QueryCorkResultsResponse



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `corkResults` | [CorkResult](#cork.v2.CorkResult) | repeated |  |






<a name="cork.v2.QueryParamsRequest"></a>

### QueryParamsRequest
QueryParamsRequest is the request type for the Query/Params gRPC method.






<a name="cork.v2.QueryParamsResponse"></a>

### QueryParamsResponse
QueryParamsRequest is the response type for the Query/Params gRPC method.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `params` | [Params](#cork.v2.Params) |  | allocation parameters |






<a name="cork.v2.QueryScheduledBlockHeightsRequest"></a>

### QueryScheduledBlockHeightsRequest
QueryScheduledBlockHeightsRequest






<a name="cork.v2.QueryScheduledBlockHeightsResponse"></a>

### QueryScheduledBlockHeightsResponse
QueryScheduledBlockHeightsResponse


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `block_heights` | [uint64](#uint64) | repeated |  |






<a name="cork.v2.QueryScheduledCorksByBlockHeightRequest"></a>

### QueryScheduledCorksByBlockHeightRequest
QueryScheduledCorksByBlockHeightRequest


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `block_height` | [uint64](#uint64) |  |  |






<a name="cork.v2.QueryScheduledCorksByBlockHeightResponse"></a>

### QueryScheduledCorksByBlockHeightResponse
QueryScheduledCorksByBlockHeightResponse


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `corks` | [ScheduledCork](#cork.v2.ScheduledCork) | repeated |  |






<a name="cork.v2.QueryScheduledCorksByIDRequest"></a>

### QueryScheduledCorksByIDRequest
QueryScheduledCorksByIDRequest


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `id` | [string](#string) |  |  |






<a name="cork.v2.QueryScheduledCorksByIDResponse"></a>

### QueryScheduledCorksByIDResponse
QueryScheduledCorksByIDResponse


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `corks` | [ScheduledCork](#cork.v2.ScheduledCork) | repeated |  |






<a name="cork.v2.QueryScheduledCorksRequest"></a>

### QueryScheduledCorksRequest
QueryScheduledCorksRequest






<a name="cork.v2.QueryScheduledCorksResponse"></a>

### QueryScheduledCorksResponse
QueryScheduledCorksResponse


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `corks` | [ScheduledCork](#cork.v2.ScheduledCork) | repeated |  |





 <!-- end messages -->

 <!-- end enums -->

 <!-- end HasExtensions -->


<a name="cork.v2.Query"></a>

### Query
Query defines the gRPC query service for the cork module.

| Method Name | Request Type | Response Type | Description | HTTP Verb | Endpoint |
| ----------- | ------------ | ------------- | ------------| ------- | -------- |
| `QueryParams` | [QueryParamsRequest](#cork.v2.QueryParamsRequest) | [QueryParamsResponse](#cork.v2.QueryParamsResponse) | QueryParams queries the allocation module parameters. | GET|/sommelier/cork/v2/params|
| `QueryCellarIDs` | [QueryCellarIDsRequest](#cork.v2.QueryCellarIDsRequest) | [QueryCellarIDsResponse](#cork.v2.QueryCellarIDsResponse) | QueryCellarIDs returns all cellars and current tick ranges | GET|/sommelier/cork/v2/cellar_ids|
| `QueryScheduledCorks` | [QueryScheduledCorksRequest](#cork.v2.QueryScheduledCorksRequest) | [QueryScheduledCorksResponse](#cork.v2.QueryScheduledCorksResponse) | QueryScheduledCorks returns all scheduled corks | GET|/sommelier/cork/v2/scheduled_corks|
| `QueryScheduledBlockHeights` | [QueryScheduledBlockHeightsRequest](#cork.v2.QueryScheduledBlockHeightsRequest) | [QueryScheduledBlockHeightsResponse](#cork.v2.QueryScheduledBlockHeightsResponse) | QueryScheduledBlockHeights returns all scheduled block heights | GET|/sommelier/cork/v2/scheduled_block_heights|
| `QueryScheduledCorksByBlockHeight` | [QueryScheduledCorksByBlockHeightRequest](#cork.v2.QueryScheduledCorksByBlockHeightRequest) | [QueryScheduledCorksByBlockHeightResponse](#cork.v2.QueryScheduledCorksByBlockHeightResponse) | QueryScheduledCorks returns all scheduled corks at a block height | GET|/sommelier/cork/v2/scheduled_corks_by_block_height/{block_height}|
| `QueryScheduledCorksByID` | [QueryScheduledCorksByIDRequest](#cork.v2.QueryScheduledCorksByIDRequest) | [QueryScheduledCorksByIDResponse](#cork.v2.QueryScheduledCorksByIDResponse) | QueryScheduledCorks returns all scheduled corks with the specified ID | GET|/sommelier/cork/v2/scheduled_corks_by_id/{id}|
| `QueryCorkResult` | [QueryCorkResultRequest](#cork.v2.QueryCorkResultRequest) | [QueryCorkResultResponse](#cork.v2.QueryCorkResultResponse) |  | GET|/sommelier/cork/v2/cork_results/{id}|
| `QueryCorkResults` | [QueryCorkResultsRequest](#cork.v2.QueryCorkResultsRequest) | [QueryCorkResultsResponse](#cork.v2.QueryCorkResultsResponse) |  | GET|/sommelier/cork/v2/cork_results|

 <!-- end services -->



<a name="cork/v2/tx.proto"></a>
<p align="right"><a href="#top">Top</a></p>

## cork/v2/tx.proto



<a name="cork.v2.MsgScheduleCorkRequest"></a>

### MsgScheduleCorkRequest
MsgScheduleCorkRequest - sdk.Msg for scheduling a cork request for on or after a specific block height


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `cork` | [Cork](#cork.v2.Cork) |  | the scheduled cork |
| `block_height` | [uint64](#uint64) |  | the block height that must be reached |
| `signer` | [string](#string) |  | signer account address |






<a name="cork.v2.MsgScheduleCorkResponse"></a>

### MsgScheduleCorkResponse



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `id` | [string](#string) |  | cork ID |





 <!-- end messages -->

 <!-- end enums -->

 <!-- end HasExtensions -->


<a name="cork.v2.Msg"></a>

### Msg
MsgService defines the msgs that the cork module handles

| Method Name | Request Type | Response Type | Description | HTTP Verb | Endpoint |
| ----------- | ------------ | ------------- | ------------| ------- | -------- |
| `ScheduleCork` | [MsgScheduleCorkRequest](#cork.v2.MsgScheduleCorkRequest) | [MsgScheduleCorkResponse](#cork.v2.MsgScheduleCorkResponse) |  | |

 <!-- end services -->



<a name="incentives/v1/genesis.proto"></a>
<p align="right"><a href="#top">Top</a></p>

## incentives/v1/genesis.proto



<a name="incentives.v1.GenesisState"></a>

### GenesisState



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `params` | [Params](#incentives.v1.Params) |  |  |






<a name="incentives.v1.Params"></a>

### Params
Params incentives parameters


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `distribution_per_block` | [cosmos.base.v1beta1.Coin](#cosmos.base.v1beta1.Coin) |  | DistributionPerBlock defines the coin to be sent to the distribution module from the community pool every block |
| `incentives_cutoff_height` | [uint64](#uint64) |  | IncentivesCutoffHeight defines the block height after which the incentives module will stop sending coins to the distribution module from the community pool |





 <!-- end messages -->

 <!-- end enums -->

 <!-- end HasExtensions -->

 <!-- end services -->



<a name="incentives/v1/query.proto"></a>
<p align="right"><a href="#top">Top</a></p>

## incentives/v1/query.proto



<a name="incentives.v1.QueryAPYRequest"></a>

### QueryAPYRequest
QueryAPYRequest is the request type for the QueryAPY gRPC method.






<a name="incentives.v1.QueryAPYResponse"></a>

### QueryAPYResponse
QueryAPYRequest is the response type for the QueryAPY gRPC method.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `apy` | [string](#string) |  |  |






<a name="incentives.v1.QueryParamsRequest"></a>

### QueryParamsRequest
QueryParamsRequest is the request type for the QueryParams gRPC method.






<a name="incentives.v1.QueryParamsResponse"></a>

### QueryParamsResponse
QueryParamsRequest is the response type for the QueryParams gRPC method.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `params` | [Params](#incentives.v1.Params) |  | allocation parameters |





 <!-- end messages -->

 <!-- end enums -->

 <!-- end HasExtensions -->


<a name="incentives.v1.Query"></a>

### Query
Query defines the gRPC query service for the cork module.

| Method Name | Request Type | Response Type | Description | HTTP Verb | Endpoint |
| ----------- | ------------ | ------------- | ------------| ------- | -------- |
| `QueryParams` | [QueryParamsRequest](#incentives.v1.QueryParamsRequest) | [QueryParamsResponse](#incentives.v1.QueryParamsResponse) | QueryParams queries the allocation module parameters. | GET|/sommelier/incentives/v1/params|
| `QueryAPY` | [QueryAPYRequest](#incentives.v1.QueryAPYRequest) | [QueryAPYResponse](#incentives.v1.QueryAPYResponse) | QueryAPY queries the APY returned from the incentives module. | GET|/sommelier/incentives/v1/apy|

 <!-- end services -->



<a name="pubsub/v1/params.proto"></a>
<p align="right"><a href="#top">Top</a></p>

## pubsub/v1/params.proto



<a name="pubsub.v1.Params"></a>

### Params
Params defines the parameters for the module.





 <!-- end messages -->

 <!-- end enums -->

 <!-- end HasExtensions -->

 <!-- end services -->



<a name="pubsub/v1/pubsub.proto"></a>
<p align="right"><a href="#top">Top</a></p>

## pubsub/v1/pubsub.proto



<a name="pubsub.v1.AddDefaultSubscriptionProposal"></a>

### AddDefaultSubscriptionProposal
set the default publisher for a given subscription ID
these can be overridden by the client


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `title` | [string](#string) |  |  |
| `description` | [string](#string) |  |  |
| `subscription_id` | [string](#string) |  |  |
| `publisher_domain` | [string](#string) |  |  |






<a name="pubsub.v1.AddDefaultSubscriptionProposalWithDeposit"></a>

### AddDefaultSubscriptionProposalWithDeposit



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `title` | [string](#string) |  |  |
| `description` | [string](#string) |  |  |
| `subscription_id` | [string](#string) |  |  |
| `publisher_domain` | [string](#string) |  |  |
| `deposit` | [string](#string) |  |  |






<a name="pubsub.v1.AddPublisherProposal"></a>

### AddPublisherProposal
governance proposal to add a publisher, with domain, adress, and ca_cert the same as the Publisher type
proof URL expected in the format: https://<domain>/<address>/cacert.pem and serving cacert.pem matching ca_cert


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `title` | [string](#string) |  |  |
| `description` | [string](#string) |  |  |
| `domain` | [string](#string) |  |  |
| `address` | [string](#string) |  |  |
| `proof_url` | [string](#string) |  |  |
| `ca_cert` | [string](#string) |  |  |






<a name="pubsub.v1.AddPublisherProposalWithDeposit"></a>

### AddPublisherProposalWithDeposit



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `title` | [string](#string) |  |  |
| `description` | [string](#string) |  |  |
| `domain` | [string](#string) |  |  |
| `address` | [string](#string) |  |  |
| `proof_url` | [string](#string) |  |  |
| `ca_cert` | [string](#string) |  |  |
| `deposit` | [string](#string) |  |  |






<a name="pubsub.v1.DefaultSubscription"></a>

### DefaultSubscription
represents a default subscription voted in by governance that can be overridden by a subscriber


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `subscription_id` | [string](#string) |  | arbitary string representing a subscription, max length of 128 |
| `publisher_domain` | [string](#string) |  | FQDN of the publisher, max length of 256 |






<a name="pubsub.v1.Publisher"></a>

### Publisher
represents a publisher, which are added via governance


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `address` | [string](#string) |  | account address of the publisher |
| `domain` | [string](#string) |  | unique key, FQDN of the publisher, max length of 256 |
| `ca_cert` | [string](#string) |  | the publisher's self-signed CA cert PEM file, expecting TLS 1.3 compatible ECDSA certificates, max length 4096 |






<a name="pubsub.v1.PublisherIntent"></a>

### PublisherIntent
represents a publisher committing to sending messages for a specific subscription ID


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `subscription_id` | [string](#string) |  | arbitary string representing a subscription, max length of 128 |
| `publisher_domain` | [string](#string) |  | FQDN of the publisher, max length of 256 |
| `method` | [PublishMethod](#pubsub.v1.PublishMethod) |  | either PULL or PUSH (see enum above for details) |
| `pull_url` | [string](#string) |  | optional, only needs to be set if using the PULL method, max length of 512 |
| `allowed_subscribers` | [AllowedSubscribers](#pubsub.v1.AllowedSubscribers) |  | either ANY, VALIDATORS, or LIST (see enum above for details) |
| `allowed_addresses` | [string](#string) | repeated | optional, must be provided if allowed_subscribers is LIST, list of account addresses, max length 256 |






<a name="pubsub.v1.RemoveDefaultSubscriptionProposal"></a>

### RemoveDefaultSubscriptionProposal
remove a default subscription


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `title` | [string](#string) |  |  |
| `description` | [string](#string) |  |  |
| `subscription_id` | [string](#string) |  |  |






<a name="pubsub.v1.RemoveDefaultSubscriptionProposalWithDeposit"></a>

### RemoveDefaultSubscriptionProposalWithDeposit



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `title` | [string](#string) |  |  |
| `description` | [string](#string) |  |  |
| `subscription_id` | [string](#string) |  |  |
| `deposit` | [string](#string) |  |  |






<a name="pubsub.v1.RemovePublisherProposal"></a>

### RemovePublisherProposal
governance proposal to remove a publisher (publishers can remove themselves, but this might be necessary in the
event of a malicious publisher or a key compromise), since Publishers are unique by domain, it's the only
necessary information to remove one


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `title` | [string](#string) |  |  |
| `description` | [string](#string) |  |  |
| `domain` | [string](#string) |  |  |






<a name="pubsub.v1.RemovePublisherProposalWithDeposit"></a>

### RemovePublisherProposalWithDeposit



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `title` | [string](#string) |  |  |
| `description` | [string](#string) |  |  |
| `domain` | [string](#string) |  |  |
| `deposit` | [string](#string) |  |  |






<a name="pubsub.v1.Subscriber"></a>

### Subscriber
represents a subscriber, can be set or modified by the owner of the subscriber address


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `address` | [string](#string) |  | unique key, account address representation of either an account or a validator |
| `ca_cert` | [string](#string) |  | the subscriber's self-signed CA cert PEM file, expecting TLS 1.3 compatible ECDSA certificates, max length 4096 |
| `push_url` | [string](#string) |  | max length of 512 |






<a name="pubsub.v1.SubscriberIntent"></a>

### SubscriberIntent
represents a subscriber requesting messages for a specific subscription ID and publisher


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `subscription_id` | [string](#string) |  | arbitary string representing a subscription, max length of 128 |
| `subscriber_address` | [string](#string) |  | account address of the subscriber |
| `publisher_domain` | [string](#string) |  | FQDN of the publisher, max length of 256 |





 <!-- end messages -->


<a name="pubsub.v1.AllowedSubscribers"></a>

### AllowedSubscribers
for a given PublisherIntent, determines what types of subscribers may subscribe

| Name | Number | Description |
| ---- | ------ | ----------- |
| ANY | 0 | any valid account address |
| VALIDATORS | 1 | account address must map to a validator in the active validator set |
| LIST | 2 | a specific list of account addresses |



<a name="pubsub.v1.PublishMethod"></a>

### PublishMethod
for a given PublisherIntent, whether or not it is pulled or pushed

| Name | Number | Description |
| ---- | ------ | ----------- |
| PULL | 0 | subscribers should pull from the provided URL |
| PUSH | 1 | subscribers must provide a URL to receive push messages |


 <!-- end enums -->

 <!-- end HasExtensions -->

 <!-- end services -->



<a name="pubsub/v1/genesis.proto"></a>
<p align="right"><a href="#top">Top</a></p>

## pubsub/v1/genesis.proto



<a name="pubsub.v1.GenesisState"></a>

### GenesisState
GenesisState defines the pubsub module's genesis state.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `params` | [Params](#pubsub.v1.Params) |  |  |
| `publishers` | [Publisher](#pubsub.v1.Publisher) | repeated |  |
| `subscribers` | [Subscriber](#pubsub.v1.Subscriber) | repeated |  |
| `publisher_intents` | [PublisherIntent](#pubsub.v1.PublisherIntent) | repeated |  |
| `subscriber_intents` | [SubscriberIntent](#pubsub.v1.SubscriberIntent) | repeated |  |
| `default_subscriptions` | [DefaultSubscription](#pubsub.v1.DefaultSubscription) | repeated |  |





 <!-- end messages -->

 <!-- end enums -->

 <!-- end HasExtensions -->

 <!-- end services -->



<a name="pubsub/v1/query.proto"></a>
<p align="right"><a href="#top">Top</a></p>

## pubsub/v1/query.proto



<a name="pubsub.v1.QueryDefaultSubscriptionRequest"></a>

### QueryDefaultSubscriptionRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `subscription_id` | [string](#string) |  |  |






<a name="pubsub.v1.QueryDefaultSubscriptionResponse"></a>

### QueryDefaultSubscriptionResponse



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `default_subscription` | [DefaultSubscription](#pubsub.v1.DefaultSubscription) |  |  |






<a name="pubsub.v1.QueryDefaultSubscriptionsRequest"></a>

### QueryDefaultSubscriptionsRequest







<a name="pubsub.v1.QueryDefaultSubscriptionsResponse"></a>

### QueryDefaultSubscriptionsResponse



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `default_subscriptions` | [DefaultSubscription](#pubsub.v1.DefaultSubscription) | repeated |  |






<a name="pubsub.v1.QueryParamsRequest"></a>

### QueryParamsRequest







<a name="pubsub.v1.QueryParamsResponse"></a>

### QueryParamsResponse



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `params` | [Params](#pubsub.v1.Params) |  |  |






<a name="pubsub.v1.QueryPublisherIntentRequest"></a>

### QueryPublisherIntentRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `publisher_domain` | [string](#string) |  |  |
| `subscription_id` | [string](#string) |  |  |






<a name="pubsub.v1.QueryPublisherIntentResponse"></a>

### QueryPublisherIntentResponse



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `publisher_intent` | [PublisherIntent](#pubsub.v1.PublisherIntent) |  |  |






<a name="pubsub.v1.QueryPublisherIntentsByPublisherDomainRequest"></a>

### QueryPublisherIntentsByPublisherDomainRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `publisher_domain` | [string](#string) |  |  |






<a name="pubsub.v1.QueryPublisherIntentsByPublisherDomainResponse"></a>

### QueryPublisherIntentsByPublisherDomainResponse



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `publisher_intents` | [PublisherIntent](#pubsub.v1.PublisherIntent) | repeated |  |






<a name="pubsub.v1.QueryPublisherIntentsBySubscriptionIDRequest"></a>

### QueryPublisherIntentsBySubscriptionIDRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `subscription_id` | [string](#string) |  |  |






<a name="pubsub.v1.QueryPublisherIntentsBySubscriptionIDResponse"></a>

### QueryPublisherIntentsBySubscriptionIDResponse



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `publisher_intents` | [PublisherIntent](#pubsub.v1.PublisherIntent) | repeated |  |






<a name="pubsub.v1.QueryPublisherIntentsRequest"></a>

### QueryPublisherIntentsRequest







<a name="pubsub.v1.QueryPublisherIntentsResponse"></a>

### QueryPublisherIntentsResponse



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `publisher_intents` | [PublisherIntent](#pubsub.v1.PublisherIntent) | repeated |  |






<a name="pubsub.v1.QueryPublisherRequest"></a>

### QueryPublisherRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `publisher_domain` | [string](#string) |  |  |






<a name="pubsub.v1.QueryPublisherResponse"></a>

### QueryPublisherResponse



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `publisher` | [Publisher](#pubsub.v1.Publisher) |  |  |






<a name="pubsub.v1.QueryPublishersRequest"></a>

### QueryPublishersRequest







<a name="pubsub.v1.QueryPublishersResponse"></a>

### QueryPublishersResponse



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `publishers` | [Publisher](#pubsub.v1.Publisher) | repeated |  |






<a name="pubsub.v1.QuerySubscriberIntentRequest"></a>

### QuerySubscriberIntentRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `subscriber_address` | [string](#string) |  |  |
| `subscription_id` | [string](#string) |  |  |






<a name="pubsub.v1.QuerySubscriberIntentResponse"></a>

### QuerySubscriberIntentResponse



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `subscriber_intent` | [SubscriberIntent](#pubsub.v1.SubscriberIntent) |  |  |






<a name="pubsub.v1.QuerySubscriberIntentsByPublisherDomainRequest"></a>

### QuerySubscriberIntentsByPublisherDomainRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `publisher_domain` | [string](#string) |  |  |






<a name="pubsub.v1.QuerySubscriberIntentsByPublisherDomainResponse"></a>

### QuerySubscriberIntentsByPublisherDomainResponse



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `subscriber_intents` | [SubscriberIntent](#pubsub.v1.SubscriberIntent) | repeated |  |






<a name="pubsub.v1.QuerySubscriberIntentsBySubscriberAddressRequest"></a>

### QuerySubscriberIntentsBySubscriberAddressRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `subscriber_address` | [string](#string) |  |  |






<a name="pubsub.v1.QuerySubscriberIntentsBySubscriberAddressResponse"></a>

### QuerySubscriberIntentsBySubscriberAddressResponse



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `subscriber_intents` | [SubscriberIntent](#pubsub.v1.SubscriberIntent) | repeated |  |






<a name="pubsub.v1.QuerySubscriberIntentsBySubscriptionIDRequest"></a>

### QuerySubscriberIntentsBySubscriptionIDRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `subscription_id` | [string](#string) |  |  |






<a name="pubsub.v1.QuerySubscriberIntentsBySubscriptionIDResponse"></a>

### QuerySubscriberIntentsBySubscriptionIDResponse



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `subscriber_intents` | [SubscriberIntent](#pubsub.v1.SubscriberIntent) | repeated |  |






<a name="pubsub.v1.QuerySubscriberIntentsRequest"></a>

### QuerySubscriberIntentsRequest







<a name="pubsub.v1.QuerySubscriberIntentsResponse"></a>

### QuerySubscriberIntentsResponse



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `subscriber_intents` | [SubscriberIntent](#pubsub.v1.SubscriberIntent) | repeated |  |






<a name="pubsub.v1.QuerySubscriberRequest"></a>

### QuerySubscriberRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `subscriber_address` | [string](#string) |  |  |






<a name="pubsub.v1.QuerySubscriberResponse"></a>

### QuerySubscriberResponse



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `subscriber` | [Subscriber](#pubsub.v1.Subscriber) |  |  |






<a name="pubsub.v1.QuerySubscribersRequest"></a>

### QuerySubscribersRequest







<a name="pubsub.v1.QuerySubscribersResponse"></a>

### QuerySubscribersResponse



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `subscribers` | [Subscriber](#pubsub.v1.Subscriber) | repeated |  |





 <!-- end messages -->

 <!-- end enums -->

 <!-- end HasExtensions -->


<a name="pubsub.v1.Query"></a>

### Query


| Method Name | Request Type | Response Type | Description | HTTP Verb | Endpoint |
| ----------- | ------------ | ------------- | ------------| ------- | -------- |
| `Params` | [QueryParamsRequest](#pubsub.v1.QueryParamsRequest) | [QueryParamsResponse](#pubsub.v1.QueryParamsResponse) |  | GET|/sommelier/pubsub/v1/params|
| `QueryPublisher` | [QueryPublisherRequest](#pubsub.v1.QueryPublisherRequest) | [QueryPublisherResponse](#pubsub.v1.QueryPublisherResponse) |  | GET|/sommelier/pubsub/v1/publishers/{publisher_domain}|
| `QueryPublishers` | [QueryPublishersRequest](#pubsub.v1.QueryPublishersRequest) | [QueryPublishersResponse](#pubsub.v1.QueryPublishersResponse) |  | GET|/sommelier/pubsub/v1/publishers|
| `QuerySubscriber` | [QuerySubscriberRequest](#pubsub.v1.QuerySubscriberRequest) | [QuerySubscriberResponse](#pubsub.v1.QuerySubscriberResponse) |  | GET|/sommelier/pubsub/v1/subscribers/{subscriber_address}|
| `QuerySubscribers` | [QuerySubscribersRequest](#pubsub.v1.QuerySubscribersRequest) | [QuerySubscribersResponse](#pubsub.v1.QuerySubscribersResponse) |  | GET|/sommelier/pubsub/v1/subscribers|
| `QueryPublisherIntent` | [QueryPublisherIntentRequest](#pubsub.v1.QueryPublisherIntentRequest) | [QueryPublisherIntentResponse](#pubsub.v1.QueryPublisherIntentResponse) |  | GET|/sommelier/pubsub/v1/publisher_intents/{publisher_domain}/{subscription_id}|
| `QueryPublisherIntents` | [QueryPublisherIntentsRequest](#pubsub.v1.QueryPublisherIntentsRequest) | [QueryPublisherIntentsResponse](#pubsub.v1.QueryPublisherIntentsResponse) |  | GET|/sommelier/pubsub/v1/publisher_intents|
| `QueryPublisherIntentsByPublisherDomain` | [QueryPublisherIntentsByPublisherDomainRequest](#pubsub.v1.QueryPublisherIntentsByPublisherDomainRequest) | [QueryPublisherIntentsByPublisherDomainResponse](#pubsub.v1.QueryPublisherIntentsByPublisherDomainResponse) |  | GET|/sommelier/pubsub/v1/publisher_intents/{publisher_domain}|
| `QueryPublisherIntentsBySubscriptionID` | [QueryPublisherIntentsBySubscriptionIDRequest](#pubsub.v1.QueryPublisherIntentsBySubscriptionIDRequest) | [QueryPublisherIntentsBySubscriptionIDResponse](#pubsub.v1.QueryPublisherIntentsBySubscriptionIDResponse) |  | GET|/sommelier/pubsub/v1/publisher_intents_by_subscription_id/{subscription_id}|
| `QuerySubscriberIntent` | [QuerySubscriberIntentRequest](#pubsub.v1.QuerySubscriberIntentRequest) | [QuerySubscriberIntentResponse](#pubsub.v1.QuerySubscriberIntentResponse) |  | GET|/sommelier/pubsub/v1/subscriber_intents/{subscriber_address}/{subscription_id}|
| `QuerySubscriberIntents` | [QuerySubscriberIntentsRequest](#pubsub.v1.QuerySubscriberIntentsRequest) | [QuerySubscriberIntentsResponse](#pubsub.v1.QuerySubscriberIntentsResponse) |  | GET|/sommelier/pubsub/v1/subscriber_intents|
| `QuerySubscriberIntentsBySubscriberAddress` | [QuerySubscriberIntentsBySubscriberAddressRequest](#pubsub.v1.QuerySubscriberIntentsBySubscriberAddressRequest) | [QuerySubscriberIntentsBySubscriberAddressResponse](#pubsub.v1.QuerySubscriberIntentsBySubscriberAddressResponse) |  | GET|/sommelier/pubsub/v1/subscriber_intents/{subscriber_address}|
| `QuerySubscriberIntentsBySubscriptionID` | [QuerySubscriberIntentsBySubscriptionIDRequest](#pubsub.v1.QuerySubscriberIntentsBySubscriptionIDRequest) | [QuerySubscriberIntentsBySubscriptionIDResponse](#pubsub.v1.QuerySubscriberIntentsBySubscriptionIDResponse) |  | GET|/sommelier/pubsub/v1/subscriber_intents_by_subscription_id/{subscription_id}|
| `QuerySubscriberIntentsByPublisherDomain` | [QuerySubscriberIntentsByPublisherDomainRequest](#pubsub.v1.QuerySubscriberIntentsByPublisherDomainRequest) | [QuerySubscriberIntentsByPublisherDomainResponse](#pubsub.v1.QuerySubscriberIntentsByPublisherDomainResponse) |  | GET|/sommelier/pubsub/v1/subscriber_intents_by_publisher_domain/{publisher_domain}|
| `QueryDefaultSubscription` | [QueryDefaultSubscriptionRequest](#pubsub.v1.QueryDefaultSubscriptionRequest) | [QueryDefaultSubscriptionResponse](#pubsub.v1.QueryDefaultSubscriptionResponse) |  | GET|/sommelier/pubsub/v1/default_subscriptions/{subscription_id}|
| `QueryDefaultSubscriptions` | [QueryDefaultSubscriptionsRequest](#pubsub.v1.QueryDefaultSubscriptionsRequest) | [QueryDefaultSubscriptionsResponse](#pubsub.v1.QueryDefaultSubscriptionsResponse) |  | GET|/sommelier/pubsub/v1/default_subscriptions|

 <!-- end services -->



<a name="pubsub/v1/tx.proto"></a>
<p align="right"><a href="#top">Top</a></p>

## pubsub/v1/tx.proto



<a name="pubsub.v1.MsgAddPublisherIntentRequest"></a>

### MsgAddPublisherIntentRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `publisher_intent` | [PublisherIntent](#pubsub.v1.PublisherIntent) |  |  |
| `signer` | [string](#string) |  |  |






<a name="pubsub.v1.MsgAddPublisherIntentResponse"></a>

### MsgAddPublisherIntentResponse







<a name="pubsub.v1.MsgAddSubscriberIntentRequest"></a>

### MsgAddSubscriberIntentRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `subscriber_intent` | [SubscriberIntent](#pubsub.v1.SubscriberIntent) |  |  |
| `signer` | [string](#string) |  |  |






<a name="pubsub.v1.MsgAddSubscriberIntentResponse"></a>

### MsgAddSubscriberIntentResponse







<a name="pubsub.v1.MsgAddSubscriberRequest"></a>

### MsgAddSubscriberRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `subscriber` | [Subscriber](#pubsub.v1.Subscriber) |  |  |
| `signer` | [string](#string) |  |  |






<a name="pubsub.v1.MsgAddSubscriberResponse"></a>

### MsgAddSubscriberResponse







<a name="pubsub.v1.MsgRemovePublisherIntentRequest"></a>

### MsgRemovePublisherIntentRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `subscription_id` | [string](#string) |  |  |
| `publisher_domain` | [string](#string) |  |  |
| `signer` | [string](#string) |  |  |






<a name="pubsub.v1.MsgRemovePublisherIntentResponse"></a>

### MsgRemovePublisherIntentResponse







<a name="pubsub.v1.MsgRemovePublisherRequest"></a>

### MsgRemovePublisherRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `publisher_domain` | [string](#string) |  |  |
| `signer` | [string](#string) |  |  |






<a name="pubsub.v1.MsgRemovePublisherResponse"></a>

### MsgRemovePublisherResponse







<a name="pubsub.v1.MsgRemoveSubscriberIntentRequest"></a>

### MsgRemoveSubscriberIntentRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `subscription_id` | [string](#string) |  |  |
| `subscriber_address` | [string](#string) |  |  |
| `signer` | [string](#string) |  |  |






<a name="pubsub.v1.MsgRemoveSubscriberIntentResponse"></a>

### MsgRemoveSubscriberIntentResponse







<a name="pubsub.v1.MsgRemoveSubscriberRequest"></a>

### MsgRemoveSubscriberRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `subscriber_address` | [string](#string) |  |  |
| `signer` | [string](#string) |  |  |






<a name="pubsub.v1.MsgRemoveSubscriberResponse"></a>

### MsgRemoveSubscriberResponse






 <!-- end messages -->

 <!-- end enums -->

 <!-- end HasExtensions -->


<a name="pubsub.v1.Msg"></a>

### Msg


| Method Name | Request Type | Response Type | Description | HTTP Verb | Endpoint |
| ----------- | ------------ | ------------- | ------------| ------- | -------- |
| `RemovePublisher` | [MsgRemovePublisherRequest](#pubsub.v1.MsgRemovePublisherRequest) | [MsgRemovePublisherResponse](#pubsub.v1.MsgRemovePublisherResponse) |  | |
| `AddSubscriber` | [MsgAddSubscriberRequest](#pubsub.v1.MsgAddSubscriberRequest) | [MsgAddSubscriberResponse](#pubsub.v1.MsgAddSubscriberResponse) |  | |
| `RemoveSubscriber` | [MsgRemoveSubscriberRequest](#pubsub.v1.MsgRemoveSubscriberRequest) | [MsgRemoveSubscriberResponse](#pubsub.v1.MsgRemoveSubscriberResponse) |  | |
| `AddPublisherIntent` | [MsgAddPublisherIntentRequest](#pubsub.v1.MsgAddPublisherIntentRequest) | [MsgAddPublisherIntentResponse](#pubsub.v1.MsgAddPublisherIntentResponse) |  | |
| `RemovePublisherIntent` | [MsgRemovePublisherIntentRequest](#pubsub.v1.MsgRemovePublisherIntentRequest) | [MsgRemovePublisherIntentResponse](#pubsub.v1.MsgRemovePublisherIntentResponse) |  | |
| `AddSubscriberIntent` | [MsgAddSubscriberIntentRequest](#pubsub.v1.MsgAddSubscriberIntentRequest) | [MsgAddSubscriberIntentResponse](#pubsub.v1.MsgAddSubscriberIntentResponse) |  | |
| `RemoveSubscriberIntent` | [MsgRemoveSubscriberIntentRequest](#pubsub.v1.MsgRemoveSubscriberIntentRequest) | [MsgRemoveSubscriberIntentResponse](#pubsub.v1.MsgRemoveSubscriberIntentResponse) |  | |

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

