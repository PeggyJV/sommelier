package cli

import (
	"os"
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/peggyjv/sommelier/v10/app/params"
	"github.com/peggyjv/sommelier/v10/x/axelarcork/types"

	"github.com/cosmos/cosmos-sdk/testutil"
	moduletestutil "github.com/cosmos/cosmos-sdk/types/module/testutil"
	"github.com/stretchr/testify/require"
)

// The SDK config defaults to the "cosmos" bech32 prefix; the app sets "somm" at
// init. Without this, somm1... addresses fail ValidateBasic misleadingly.
func TestMain(m *testing.M) {
	params.SetAddressPrefixes()
	os.Exit(m.Run())
}

func TestParseAddManagedCellarsProposal(t *testing.T) {
	encodingConfig := moduletestutil.MakeTestEncodingConfig()

	okJSON := testutil.WriteToNewTempFile(t, `
{
  "title": "Dollary-doos LP Arbitrum Cellar Proposal",
  "description": "I have a hunch",
  "chain_id": 42161,
  "cellar_ids": ["0x123801a7D398351b8bE11C439e05C5B3259aeC9B", "0x456801a7D398351b8bE11C439e05C5B3259aeC9B"],
  "publisher_domain": "example.com",
  "deposit": "10000usomm"
}
`)

	proposal := types.AddAxelarManagedCellarIDsProposalWithDeposit{}
	contents, err := os.ReadFile(okJSON.Name())
	require.NoError(t, err)

	err = encodingConfig.Codec.UnmarshalJSON(contents, &proposal)
	require.NoError(t, err)

	require.Equal(t, "Dollary-doos LP Arbitrum Cellar Proposal", proposal.Title)
	require.Equal(t, "I have a hunch", proposal.Description)
	require.Equal(t, uint64(42161), proposal.ChainId)
	require.Equal(t, "0x123801a7D398351b8bE11C439e05C5B3259aeC9B", proposal.CellarIds[0])
	require.Equal(t, "0x456801a7D398351b8bE11C439e05C5B3259aeC9B", proposal.CellarIds[1])
	require.Equal(t, "example.com", proposal.PublisherDomain)
	require.Equal(t, "10000usomm", proposal.Deposit)
}

func TestParseRemoveManagedCellarsProposal(t *testing.T) {
	encodingConfig := moduletestutil.MakeTestEncodingConfig()

	okJSON := testutil.WriteToNewTempFile(t, `
{
  "title": "Dollary-doos LP Arbitrum Cellar Proposal",
  "description": "I have a hunch",
  "chain_id": 42161,
  "cellar_ids": ["0x123801a7D398351b8bE11C439e05C5B3259aeC9B", "0x456801a7D398351b8bE11C439e05C5B3259aeC9B"],
  "deposit": "10000usomm"
}
`)

	proposal := types.RemoveAxelarManagedCellarIDsProposalWithDeposit{}
	contents, err := os.ReadFile(okJSON.Name())
	require.NoError(t, err)

	err = encodingConfig.Codec.UnmarshalJSON(contents, &proposal)
	require.NoError(t, err)

	require.Equal(t, "Dollary-doos LP Arbitrum Cellar Proposal", proposal.Title)
	require.Equal(t, "I have a hunch", proposal.Description)
	require.Equal(t, uint64(42161), proposal.ChainId)
	require.Equal(t, "0x123801a7D398351b8bE11C439e05C5B3259aeC9B", proposal.CellarIds[0])
	require.Equal(t, "0x456801a7D398351b8bE11C439e05C5B3259aeC9B", proposal.CellarIds[1])
	require.Equal(t, "10000usomm", proposal.Deposit)
}

func TestParseSubmitScheduledCorkProposal(t *testing.T) {
	encodingConfig := moduletestutil.MakeTestEncodingConfig()

	okJSON := testutil.WriteToNewTempFile(t, `
{
  "title": "Scheduled Axelar Arbitrum cork proposal",
  "description": "I have a hunch",
  "block_height": 100000,
  "chain_id": 42161,
  "target_contract_address": "0x123801a7D398351b8bE11C439e05C5B3259aeC9B",
  "contract_call_proto_json": "{\"cellar_id\":\"0x123801a7D398351b8bE11C439e05C5B3259aeC9B\",\"cellar_v1\":{\"some_fuction\":{\"function_args\":{}},\"block_height\":12345}}",
  "deadline": 1706225320,
  "deposit": "10000usomm"
}
`)

	proposal := types.AxelarScheduledCorkProposalWithDeposit{}
	contents, err := os.ReadFile(okJSON.Name())
	require.NoError(t, err)

	err = encodingConfig.Codec.UnmarshalJSON(contents, &proposal)
	require.NoError(t, err)

	require.Equal(t, "Scheduled Axelar Arbitrum cork proposal", proposal.Title)
	require.Equal(t, "I have a hunch", proposal.Description)
	require.Equal(t, uint64(100000), proposal.BlockHeight)
	require.Equal(t, uint64(42161), proposal.ChainId)
	require.Equal(t, "0x123801a7D398351b8bE11C439e05C5B3259aeC9B", proposal.TargetContractAddress)
	require.Equal(t, "{\"cellar_id\":\"0x123801a7D398351b8bE11C439e05C5B3259aeC9B\",\"cellar_v1\":{\"some_fuction\":{\"function_args\":{}},\"block_height\":12345}}", proposal.ContractCallProtoJson)
	require.Equal(t, uint64(1706225320), proposal.Deadline)
	require.Equal(t, "10000usomm", proposal.Deposit)
}

func TestParseAxelarCommunityPoolSpendProposal(t *testing.T) {
	encodingConfig := moduletestutil.MakeTestEncodingConfig()

	okJSON := testutil.WriteToNewTempFile(t, `
{
  "title": "Community Pool Axelar Spend",
  "description": "Bridge me some tokens to Arbitrum!",
  "recipient": "0x123801a7D398351b8bE11C439e05C5B3259aeC9B",
  "chain_id": 42161,
  "amount": "20000usomm",
  "deposit": "10000usomm"
}
`)

	proposal := types.AxelarCommunityPoolSpendProposalForCLI{}
	contents, err := os.ReadFile(okJSON.Name())
	require.NoError(t, err)

	err = encodingConfig.Codec.UnmarshalJSON(contents, &proposal)
	require.NoError(t, err)

	require.Equal(t, "Community Pool Axelar Spend", proposal.Title)
	require.Equal(t, "Bridge me some tokens to Arbitrum!", proposal.Description)
	require.Equal(t, "0x123801a7D398351b8bE11C439e05C5B3259aeC9B", proposal.Recipient)
	require.Equal(t, uint64(42161), proposal.ChainId)
	require.Equal(t, "20000usomm", proposal.Amount)
	require.Equal(t, "10000usomm", proposal.Deposit)
}

func TestParseAddChainConfigurationProposal(t *testing.T) {
	encodingConfig := moduletestutil.MakeTestEncodingConfig()

	okJSON := testutil.WriteToNewTempFile(t, `
{
  "title": "Enable Arbitrum proposal",
  "description": "allow cellars to be used on Arbitrum",
  "chain_configuration": {
	"name": "arbitrum",
	"id": 42161,
	"proxy_address": "0x0000000000000000000000000000000000000000",
	"bridge_fees": [
		{
			"denom": "usomm",
			"amount": "100000"
		}
	]
  },
  "deposit": "10000usomm"
}
`)

	proposal := types.AddChainConfigurationProposalWithDeposit{}
	contents, err := os.ReadFile(okJSON.Name())
	require.NoError(t, err)

	err = encodingConfig.Codec.UnmarshalJSON(contents, &proposal)
	require.NoError(t, err)

	require.Equal(t, "Enable Arbitrum proposal", proposal.Title)
	require.Equal(t, "allow cellars to be used on Arbitrum", proposal.Description)
	require.Equal(t, "arbitrum", proposal.ChainConfiguration.Name)
	require.Equal(t, uint64(42161), proposal.ChainConfiguration.Id)
	require.Equal(t, "0x0000000000000000000000000000000000000000", proposal.ChainConfiguration.ProxyAddress)
	require.Equal(t, "usomm", proposal.ChainConfiguration.BridgeFees[0].Denom)
	require.Equal(t, sdk.NewInt(100000), proposal.ChainConfiguration.BridgeFees[0].Amount)
	require.Equal(t, "10000usomm", proposal.Deposit)
}

func TestParseRemoveChainConfigurationProposal(t *testing.T) {
	encodingConfig := moduletestutil.MakeTestEncodingConfig()

	okJSON := testutil.WriteToNewTempFile(t, `
{
  "title": "Remove Arbitrum chain config",
  "description": "not using Arbitrum any more",
  "chain_id": 42161,
  "deposit": "10000usomm"
}
`)

	proposal := types.RemoveChainConfigurationProposalWithDeposit{}
	contents, err := os.ReadFile(okJSON.Name())
	require.NoError(t, err)

	err = encodingConfig.Codec.UnmarshalJSON(contents, &proposal)
	require.NoError(t, err)

	require.Equal(t, "Remove Arbitrum chain config", proposal.Title)
	require.Equal(t, "not using Arbitrum any more", proposal.Description)
	require.Equal(t, uint64(42161), proposal.ChainId)
	require.Equal(t, "10000usomm", proposal.Deposit)
}

func TestParseUpgradeAxelarProxyContractProposal(t *testing.T) {
	encodingConfig := moduletestutil.MakeTestEncodingConfig()

	okJSON := testutil.WriteToNewTempFile(t, `
{
  "title": "Upgrade Axelar proxy contract proposal",
  "description": "I have a hunch",
  "chain_id": 42161,
  "new_proxy_address": "0x123801a7D398351b8bE11C439e05C5B3259aeC9B",
  "deposit": "10000usomm"
}
`)

	proposal := types.UpgradeAxelarProxyContractProposalWithDeposit{}
	contents, err := os.ReadFile(okJSON.Name())
	require.NoError(t, err)

	err = encodingConfig.Codec.UnmarshalJSON(contents, &proposal)
	require.NoError(t, err)

	require.Equal(t, "Upgrade Axelar proxy contract proposal", proposal.Title)
	require.Equal(t, "I have a hunch", proposal.Description)
	require.Equal(t, uint64(42161), proposal.ChainId)
	require.Equal(t, "0x123801a7D398351b8bE11C439e05C5B3259aeC9B", proposal.NewProxyAddress)
	require.Equal(t, "10000usomm", proposal.Deposit)
}

func TestParseCancelAxelarProxyContractUpgradeProposal(t *testing.T) {
	encodingConfig := moduletestutil.MakeTestEncodingConfig()

	okJSON := testutil.WriteToNewTempFile(t, `
{
  "title": "Cancel Axelar proxy contract upgrade proposal",
  "description": "I have a hunch",
  "chain_id": 42161,
  "deposit": "10000usomm"
}
`)

	proposal := types.CancelAxelarProxyContractUpgradeProposalWithDeposit{}
	contents, err := os.ReadFile(okJSON.Name())
	require.NoError(t, err)

	err = encodingConfig.Codec.UnmarshalJSON(contents, &proposal)
	require.NoError(t, err)

	require.Equal(t, "Cancel Axelar proxy contract upgrade proposal", proposal.Title)
	require.Equal(t, "I have a hunch", proposal.Description)
	require.Equal(t, uint64(42161), proposal.ChainId)
	require.Equal(t, "10000usomm", proposal.Deposit)
}

// schedule-axelar-cork previously did []byte(args[3]) under a
// "todo: how are contract calls submitted?" comment, which handed the cellar the
// ASCII of the hex string rather than the ABI-encoded call. The transaction
// succeeds on Sommelier and the call then reverts on the destination chain, with
// nothing locally to indicate why -- a bad failure mode for a cellar wind-down.
func TestBuildScheduleAxelarCorkMsg(t *testing.T) {
	const (
		cellar  = "0x123801a7D398351b8bE11C439e05C5B3259aeC9B"
		signer  = "somm1lcsjy2d5s33h0sddd8lpuqvwyz5ruz7ju4aeqa"
		chainID = uint64(42161)
	)

	t.Run("hex call is decoded to bytes", func(t *testing.T) {
		msg, err := buildScheduleAxelarCorkMsg(signer, chainID, cellar, 1000, 1800000000, "0xdeadbeef")
		require.NoError(t, err)
		require.Equal(t, []byte{0xde, 0xad, 0xbe, 0xef}, msg.Cork.EncodedContractCall)
		require.Equal(t, cellar, msg.Cork.TargetContractAddress)
		require.Equal(t, chainID, msg.Cork.ChainId)
		require.Equal(t, chainID, msg.ChainId)
		require.Equal(t, uint64(1000), msg.BlockHeight)
		require.Equal(t, signer, msg.Signer)
	})

	t.Run("0x prefix is optional", func(t *testing.T) {
		msg, err := buildScheduleAxelarCorkMsg(signer, chainID, cellar, 1000, 1800000000, "deadbeef")
		require.NoError(t, err)
		require.Equal(t, []byte{0xde, 0xad, 0xbe, 0xef}, msg.Cork.EncodedContractCall)
	})

	t.Run("rejects a non-hex call", func(t *testing.T) {
		_, err := buildScheduleAxelarCorkMsg(signer, chainID, cellar, 1000, 1800000000, "not-hex")
		require.Error(t, err)
	})

	t.Run("rejects an invalid cellar address", func(t *testing.T) {
		_, err := buildScheduleAxelarCorkMsg(signer, chainID, "0xnothex", 1000, 1800000000, "0xdeadbeef")
		require.Error(t, err)
	})

	t.Run("rejects an empty call", func(t *testing.T) {
		_, err := buildScheduleAxelarCorkMsg(signer, chainID, cellar, 1000, 1800000000, "0x")
		require.Error(t, err)
	})
}

// The deadline is what made this command unusable: it was never set, so
// ValidateBasic rejected every invocation. Pin that it is required.
func TestBuildScheduleAxelarCorkMsgRequiresDeadline(t *testing.T) {
	_, err := buildScheduleAxelarCorkMsg(
		"somm1lcsjy2d5s33h0sddd8lpuqvwyz5ruz7ju4aeqa", 42161,
		"0x123801a7D398351b8bE11C439e05C5B3259aeC9B", 1000, 0, "0xdeadbeef")
	require.ErrorContains(t, err, "deadline")
}

// The flag must actually be wired to the command, or the deadline silently
// stays zero and every invocation fails again.
func TestScheduleAxelarCorkCmdExposesDeadlineFlag(t *testing.T) {
	cmd := CmdScheduleAxelarCork()
	require.NotNil(t, cmd.Flags().Lookup(FlagDeadline),
		"schedule-axelar-cork must expose --"+FlagDeadline)
}
