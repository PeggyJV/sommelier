<!-- This file is auto-generated. Please do not modify it yourself. -->
# Protobuf Documentation
<a name="top"></a>

## Table of Contents

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
  
- [cellarfees/v1/params.proto](#cellarfees/v1/params.proto)
    - [Params](#cellarfees.v1.Params)
  
- [cellarfees/v1/genesis.proto](#cellarfees/v1/genesis.proto)
    - [GenesisState](#cellarfees.v1.GenesisState)
  
- [cellarfees/v1/query.proto](#cellarfees/v1/query.proto)
    - [QueryModuleAccountsRequest](#cellarfees.v1.QueryModuleAccountsRequest)
    - [QueryModuleAccountsResponse](#cellarfees.v1.QueryModuleAccountsResponse)
    - [QueryParamsRequest](#cellarfees.v1.QueryParamsRequest)
    - [QueryParamsResponse](#cellarfees.v1.QueryParamsResponse)
  
    - [Query](#cellarfees.v1.Query)
  
- [cork/v1/cork.proto](#cork/v1/cork.proto)
    - [CellarIDSet](#cork.v1.CellarIDSet)
    - [Cork](#cork.v1.Cork)
    - [ScheduledCork](#cork.v1.ScheduledCork)
    - [ValidatorCork](#cork.v1.ValidatorCork)
  
- [cork/v1/tx.proto](#cork/v1/tx.proto)
    - [MsgScheduleCorkRequest](#cork.v1.MsgScheduleCorkRequest)
    - [MsgScheduleCorkResponse](#cork.v1.MsgScheduleCorkResponse)
    - [MsgSubmitCorkRequest](#cork.v1.MsgSubmitCorkRequest)
    - [MsgSubmitCorkResponse](#cork.v1.MsgSubmitCorkResponse)
  
    - [Msg](#cork.v1.Msg)
  
- [cork/v1/genesis.proto](#cork/v1/genesis.proto)
    - [GenesisState](#cork.v1.GenesisState)
    - [Params](#cork.v1.Params)
  
- [cork/v1/proposal.proto](#cork/v1/proposal.proto)
    - [AddManagedCellarIDsProposal](#cork.v1.AddManagedCellarIDsProposal)
    - [AddManagedCellarIDsProposalWithDeposit](#cork.v1.AddManagedCellarIDsProposalWithDeposit)
    - [RemoveManagedCellarIDsProposal](#cork.v1.RemoveManagedCellarIDsProposal)
    - [RemoveManagedCellarIDsProposalWithDeposit](#cork.v1.RemoveManagedCellarIDsProposalWithDeposit)
  
- [cork/v1/query.proto](#cork/v1/query.proto)
    - [QueryCellarIDsRequest](#cork.v1.QueryCellarIDsRequest)
    - [QueryCellarIDsResponse](#cork.v1.QueryCellarIDsResponse)
    - [QueryCommitPeriodRequest](#cork.v1.QueryCommitPeriodRequest)
    - [QueryCommitPeriodResponse](#cork.v1.QueryCommitPeriodResponse)
    - [QueryParamsRequest](#cork.v1.QueryParamsRequest)
    - [QueryParamsResponse](#cork.v1.QueryParamsResponse)
    - [QueryScheduledBlockHeightsRequest](#cork.v1.QueryScheduledBlockHeightsRequest)
    - [QueryScheduledBlockHeightsResponse](#cork.v1.QueryScheduledBlockHeightsResponse)
    - [QueryScheduledCorksByBlockHeightRequest](#cork.v1.QueryScheduledCorksByBlockHeightRequest)
    - [QueryScheduledCorksByBlockHeightResponse](#cork.v1.QueryScheduledCorksByBlockHeightResponse)
    - [QueryScheduledCorksRequest](#cork.v1.QueryScheduledCorksRequest)
    - [QueryScheduledCorksResponse](#cork.v1.QueryScheduledCorksResponse)
    - [QuerySubmittedCorksRequest](#cork.v1.QuerySubmittedCorksRequest)
    - [QuerySubmittedCorksResponse](#cork.v1.QuerySubmittedCorksResponse)
  
    - [Query](#cork.v1.Query)
  
- [incentives/v1/genesis.proto](#incentives/v1/genesis.proto)
    - [GenesisState](#incentives.v1.GenesisState)
    - [Params](#incentives.v1.Params)
  
- [incentives/v1/query.proto](#incentives/v1/query.proto)
    - [QueryAPYRequest](#incentives.v1.QueryAPYRequest)
    - [QueryAPYResponse](#incentives.v1.QueryAPYResponse)
    - [QueryParamsRequest](#incentives.v1.QueryParamsRequest)
    - [QueryParamsResponse](#incentives.v1.QueryParamsResponse)
  
    - [Query](#incentives.v1.Query)
  
- [Scalar Value Types](#scalar-value-types)



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
| `chain` | [ChainConfiguration](#axelarcork.v1.ChainConfiguration) |  |  |
| `ids` | [string](#string) | repeated |  |






<a name="axelarcork.v1.ChainConfiguration"></a>

### ChainConfiguration



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `name` | [string](#string) |  |  |
| `id` | [uint64](#uint64) |  |  |
| `proxy_address` | [string](#string) |  |  |






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






<a name="axelarcork.v1.AddAxelarManagedCellarIDsProposalWithDeposit"></a>

### AddAxelarManagedCellarIDsProposalWithDeposit
AddManagedCellarIDsProposalWithDeposit is a specific definition for CLI commands


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `title` | [string](#string) |  |  |
| `description` | [string](#string) |  |  |
| `chain_id` | [uint64](#uint64) |  |  |
| `cellar_ids` | [CellarIDSet](#axelarcork.v1.CellarIDSet) |  |  |
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
This format of the community spend Ethereum proposal is specifically for
the CLI to allow simple text serialization.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `title` | [string](#string) |  |  |
| `description` | [string](#string) |  |  |
| `recipient` | [string](#string) |  |  |
| `chain_id` | [uint64](#uint64) |  |  |
| `chain_name` | [string](#string) |  |  |
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
ScheduledCorkProposalWithDeposit is a specific definition for CLI commands


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
RemoveManagedCellarIDsProposalWithDeposit is a specific definition for CLI commands


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `title` | [string](#string) |  |  |
| `description` | [string](#string) |  |  |
| `chain_id` | [uint64](#uint64) |  |  |
| `cellar_ids` | [CellarIDSet](#axelarcork.v1.CellarIDSet) |  |  |
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



<a name="cellarfees/v1/params.proto"></a>
<p align="right"><a href="#top">Top</a></p>

## cellarfees/v1/params.proto



<a name="cellarfees.v1.Params"></a>

### Params
Params defines the parameters for the module.





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





 <!-- end messages -->

 <!-- end enums -->

 <!-- end HasExtensions -->

 <!-- end services -->



<a name="cellarfees/v1/query.proto"></a>
<p align="right"><a href="#top">Top</a></p>

## cellarfees/v1/query.proto



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
| `Params` | [QueryParamsRequest](#cellarfees.v1.QueryParamsRequest) | [QueryParamsResponse](#cellarfees.v1.QueryParamsResponse) |  | GET|/sommelier/cellarfees/v1/params|
| `ModuleAccounts` | [QueryModuleAccountsRequest](#cellarfees.v1.QueryModuleAccountsRequest) | [QueryModuleAccountsResponse](#cellarfees.v1.QueryModuleAccountsResponse) |  | GET|/sommeliers/cellarfees/v1/module_accounts|

 <!-- end services -->



<a name="cork/v1/cork.proto"></a>
<p align="right"><a href="#top">Top</a></p>

## cork/v1/cork.proto



<a name="cork.v1.CellarIDSet"></a>

### CellarIDSet



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `ids` | [string](#string) | repeated |  |






<a name="cork.v1.Cork"></a>

### Cork



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `encoded_contract_call` | [bytes](#bytes) |  | call body containing the ABI encoded bytes to send to the contract |
| `target_contract_address` | [string](#string) |  | address of the contract to send the call |






<a name="cork.v1.ScheduledCork"></a>

### ScheduledCork



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `cork` | [Cork](#cork.v1.Cork) |  |  |
| `block_height` | [uint64](#uint64) |  |  |
| `validator` | [string](#string) |  |  |






<a name="cork.v1.ValidatorCork"></a>

### ValidatorCork



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `cork` | [Cork](#cork.v1.Cork) |  |  |
| `validator` | [string](#string) |  |  |





 <!-- end messages -->

 <!-- end enums -->

 <!-- end HasExtensions -->

 <!-- end services -->



<a name="cork/v1/tx.proto"></a>
<p align="right"><a href="#top">Top</a></p>

## cork/v1/tx.proto



<a name="cork.v1.MsgScheduleCorkRequest"></a>

### MsgScheduleCorkRequest
MsgScheduleCorkRequest - sdk.Msg for scheduling a cork request for on or after a specific block height


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `cork` | [Cork](#cork.v1.Cork) |  | the scheduled cork |
| `block_height` | [uint64](#uint64) |  | the block height that must be reached |
| `signer` | [string](#string) |  | signer account address |






<a name="cork.v1.MsgScheduleCorkResponse"></a>

### MsgScheduleCorkResponse







<a name="cork.v1.MsgSubmitCorkRequest"></a>

### MsgSubmitCorkRequest
MsgSubmitCorkRequest - sdk.Msg for submitting calls to Ethereum through the gravity bridge contract


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `cork` | [Cork](#cork.v1.Cork) |  | the cork to send across the bridge |
| `signer` | [string](#string) |  | signer account address |






<a name="cork.v1.MsgSubmitCorkResponse"></a>

### MsgSubmitCorkResponse






 <!-- end messages -->

 <!-- end enums -->

 <!-- end HasExtensions -->


<a name="cork.v1.Msg"></a>

### Msg
MsgService defines the msgs that the cork module handles

| Method Name | Request Type | Response Type | Description | HTTP Verb | Endpoint |
| ----------- | ------------ | ------------- | ------------| ------- | -------- |
| `SubmitCork` | [MsgSubmitCorkRequest](#cork.v1.MsgSubmitCorkRequest) | [MsgSubmitCorkResponse](#cork.v1.MsgSubmitCorkResponse) |  | |
| `ScheduleCork` | [MsgScheduleCorkRequest](#cork.v1.MsgScheduleCorkRequest) | [MsgScheduleCorkResponse](#cork.v1.MsgScheduleCorkResponse) |  | |

 <!-- end services -->



<a name="cork/v1/genesis.proto"></a>
<p align="right"><a href="#top">Top</a></p>

## cork/v1/genesis.proto



<a name="cork.v1.GenesisState"></a>

### GenesisState
GenesisState - all cork state that must be provided at genesis


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `params` | [Params](#cork.v1.Params) |  |  |
| `cellar_ids` | [CellarIDSet](#cork.v1.CellarIDSet) |  |  |
| `invalidation_nonce` | [uint64](#uint64) |  |  |
| `corks` | [ValidatorCork](#cork.v1.ValidatorCork) | repeated |  |
| `scheduled_corks` | [ScheduledCork](#cork.v1.ScheduledCork) | repeated |  |






<a name="cork.v1.Params"></a>

### Params
Params cork parameters


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `vote_period` | [int64](#int64) |  | VotePeriod defines the number of blocks to wait for votes before attempting to tally |
| `vote_threshold` | [string](#string) |  | VoteThreshold defines the percentage of bonded stake required to vote each period |





 <!-- end messages -->

 <!-- end enums -->

 <!-- end HasExtensions -->

 <!-- end services -->



<a name="cork/v1/proposal.proto"></a>
<p align="right"><a href="#top">Top</a></p>

## cork/v1/proposal.proto



<a name="cork.v1.AddManagedCellarIDsProposal"></a>

### AddManagedCellarIDsProposal



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `title` | [string](#string) |  |  |
| `description` | [string](#string) |  |  |
| `cellar_ids` | [CellarIDSet](#cork.v1.CellarIDSet) |  |  |






<a name="cork.v1.AddManagedCellarIDsProposalWithDeposit"></a>

### AddManagedCellarIDsProposalWithDeposit
AddManagedCellarIDsProposalWithDeposit is a specific definition for CLI commands


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `title` | [string](#string) |  |  |
| `description` | [string](#string) |  |  |
| `cellar_ids` | [string](#string) | repeated |  |
| `deposit` | [string](#string) |  |  |






<a name="cork.v1.RemoveManagedCellarIDsProposal"></a>

### RemoveManagedCellarIDsProposal



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `title` | [string](#string) |  |  |
| `description` | [string](#string) |  |  |
| `cellar_ids` | [CellarIDSet](#cork.v1.CellarIDSet) |  |  |






<a name="cork.v1.RemoveManagedCellarIDsProposalWithDeposit"></a>

### RemoveManagedCellarIDsProposalWithDeposit
RemoveManagedCellarIDsProposalWithDeposit is a specific definition for CLI commands


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `title` | [string](#string) |  |  |
| `description` | [string](#string) |  |  |
| `cellar_ids` | [string](#string) | repeated |  |
| `deposit` | [string](#string) |  |  |





 <!-- end messages -->

 <!-- end enums -->

 <!-- end HasExtensions -->

 <!-- end services -->



<a name="cork/v1/query.proto"></a>
<p align="right"><a href="#top">Top</a></p>

## cork/v1/query.proto



<a name="cork.v1.QueryCellarIDsRequest"></a>

### QueryCellarIDsRequest
QueryCellarIDsRequest is the request type for Query/QueryCellarIDs gRPC method.






<a name="cork.v1.QueryCellarIDsResponse"></a>

### QueryCellarIDsResponse
QueryCellarIDsResponse is the response type for Query/QueryCellars gRPC method.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `cellar_ids` | [string](#string) | repeated |  |






<a name="cork.v1.QueryCommitPeriodRequest"></a>

### QueryCommitPeriodRequest
QueryCommitPeriodRequest is the request type for the Query/QueryCommitPeriod gRPC method.






<a name="cork.v1.QueryCommitPeriodResponse"></a>

### QueryCommitPeriodResponse
QueryCommitPeriodResponse is the response type for the Query/QueryCommitPeriod gRPC method.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `current_height` | [int64](#int64) |  | block height at which the query was processed |
| `vote_period_start` | [int64](#int64) |  | latest vote period start block height |
| `vote_period_end` | [int64](#int64) |  | block height at which the current voting period ends |






<a name="cork.v1.QueryParamsRequest"></a>

### QueryParamsRequest
QueryParamsRequest is the request type for the Query/Params gRPC method.






<a name="cork.v1.QueryParamsResponse"></a>

### QueryParamsResponse
QueryParamsRequest is the response type for the Query/Params gRPC method.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `params` | [Params](#cork.v1.Params) |  | allocation parameters |






<a name="cork.v1.QueryScheduledBlockHeightsRequest"></a>

### QueryScheduledBlockHeightsRequest
QueryScheduledBlockHeightsRequest






<a name="cork.v1.QueryScheduledBlockHeightsResponse"></a>

### QueryScheduledBlockHeightsResponse
QueryScheduledBlockHeightsResponse


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `block_heights` | [uint64](#uint64) | repeated |  |






<a name="cork.v1.QueryScheduledCorksByBlockHeightRequest"></a>

### QueryScheduledCorksByBlockHeightRequest
QueryScheduledCorksByBlockHeightRequest


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `block_height` | [uint64](#uint64) |  |  |






<a name="cork.v1.QueryScheduledCorksByBlockHeightResponse"></a>

### QueryScheduledCorksByBlockHeightResponse
QueryScheduledCorksByBlockHeightResponse


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `corks` | [ScheduledCork](#cork.v1.ScheduledCork) | repeated |  |






<a name="cork.v1.QueryScheduledCorksRequest"></a>

### QueryScheduledCorksRequest
QueryScheduledCorksRequest






<a name="cork.v1.QueryScheduledCorksResponse"></a>

### QueryScheduledCorksResponse
QueryScheduledCorksResponse


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `corks` | [ScheduledCork](#cork.v1.ScheduledCork) | repeated |  |






<a name="cork.v1.QuerySubmittedCorksRequest"></a>

### QuerySubmittedCorksRequest
QuerySubmittedCorksRequest is the request type for the Query/QuerySubmittedCorks gRPC query method.






<a name="cork.v1.QuerySubmittedCorksResponse"></a>

### QuerySubmittedCorksResponse
QuerySubmittedCorksResponse is the response type for the Query/QuerySubmittedCorks gRPC query method.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `corks` | [Cork](#cork.v1.Cork) | repeated | corks in keeper awaiting vote |





 <!-- end messages -->

 <!-- end enums -->

 <!-- end HasExtensions -->


<a name="cork.v1.Query"></a>

### Query
Query defines the gRPC query service for the cork module.

| Method Name | Request Type | Response Type | Description | HTTP Verb | Endpoint |
| ----------- | ------------ | ------------- | ------------| ------- | -------- |
| `QueryParams` | [QueryParamsRequest](#cork.v1.QueryParamsRequest) | [QueryParamsResponse](#cork.v1.QueryParamsResponse) | QueryParams queries the allocation module parameters. | GET|/sommelier/cork/v1/params|
| `QuerySubmittedCorks` | [QuerySubmittedCorksRequest](#cork.v1.QuerySubmittedCorksRequest) | [QuerySubmittedCorksResponse](#cork.v1.QuerySubmittedCorksResponse) | QuerySubmittedCorks queries the submitted corks awaiting vote | GET|/sommelier/cork/v1/submitted|
| `QueryCommitPeriod` | [QueryCommitPeriodRequest](#cork.v1.QueryCommitPeriodRequest) | [QueryCommitPeriodResponse](#cork.v1.QueryCommitPeriodResponse) | QueryCommitPeriod queries the heights for the current voting period (current, start and end) | GET|/sommelier/cork/v1/commit_period|
| `QueryCellarIDs` | [QueryCellarIDsRequest](#cork.v1.QueryCellarIDsRequest) | [QueryCellarIDsResponse](#cork.v1.QueryCellarIDsResponse) | QueryCellarIDs returns all cellars and current tick ranges | GET|/sommelier/cork/v1/cellar_ids|
| `QueryScheduledCorks` | [QueryScheduledCorksRequest](#cork.v1.QueryScheduledCorksRequest) | [QueryScheduledCorksResponse](#cork.v1.QueryScheduledCorksResponse) | QueryScheduledCorks returns all scheduled corks | GET|/sommelier/cork/v1/scheduled_corks|
| `QueryScheduledBlockHeights` | [QueryScheduledBlockHeightsRequest](#cork.v1.QueryScheduledBlockHeightsRequest) | [QueryScheduledBlockHeightsResponse](#cork.v1.QueryScheduledBlockHeightsResponse) | QueryScheduledBlockHeights returns all scheduled block heights | GET|/sommelier/cork/v1/scheduled_block_heights|
| `QueryScheduledCorksByBlockHeight` | [QueryScheduledCorksByBlockHeightRequest](#cork.v1.QueryScheduledCorksByBlockHeightRequest) | [QueryScheduledCorksByBlockHeightResponse](#cork.v1.QueryScheduledCorksByBlockHeightResponse) | QueryScheduledCorks returns all scheduled corks at a block height | GET|/sommelier/cork/v1/scheduled_corks_by_block_height/{block_height}|

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

