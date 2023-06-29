package cli

import (
	"io/ioutil"
	"testing"

	"github.com/cosmos/cosmos-sdk/simapp/params"
	"github.com/peggyjv/sommelier/v6/x/axelarcork/types"

	"github.com/cosmos/cosmos-sdk/testutil"
	"github.com/stretchr/testify/require"
)

func TestParseAddManagedCellarsProposal(t *testing.T) {
	encodingConfig := params.MakeTestEncodingConfig()

	okJSON := testutil.WriteToNewTempFile(t, `
{
  "title": "Dollary-doos LP Cellar Proposal",
  "description": "I have a hunch",
  "cellar_ids": {"ids": ["0x123801a7D398351b8bE11C439e05C5B3259aeC9B", "0x456801a7D398351b8bE11C439e05C5B3259aeC9B"]},
  "deposit": "1000stake"
}
`)

	proposal := types.AddAxelarManagedCellarIDsProposalWithDeposit{}
	contents, err := ioutil.ReadFile(okJSON.Name())
	require.NoError(t, err)

	err = encodingConfig.Marshaler.UnmarshalJSON(contents, &proposal)
	require.NoError(t, err)

	require.Equal(t, "Dollary-doos LP Cellar Proposal", proposal.Title)
	require.Equal(t, "I have a hunch", proposal.Description)
	require.Equal(t, "0x123801a7D398351b8bE11C439e05C5B3259aeC9B", proposal.CellarIds.Ids[0])
	require.Equal(t, "0x456801a7D398351b8bE11C439e05C5B3259aeC9B", proposal.CellarIds.Ids[1])
	require.Equal(t, "1000stake", proposal.Deposit)
}

func TestParseRemoveManagedCellarsProposal(t *testing.T) {
	encodingConfig := params.MakeTestEncodingConfig()

	okJSON := testutil.WriteToNewTempFile(t, `
{
  "title": "Dollary-doos LP Cellar Proposal",
  "description": "I have a hunch",
  "cellar_ids": {"ids": ["0x123801a7D398351b8bE11C439e05C5B3259aeC9B", "0x456801a7D398351b8bE11C439e05C5B3259aeC9B"]},
  "deposit": "1000stake"
}
`)

	proposal := types.RemoveAxelarManagedCellarIDsProposalWithDeposit{}
	contents, err := ioutil.ReadFile(okJSON.Name())
	require.NoError(t, err)

	err = encodingConfig.Marshaler.UnmarshalJSON(contents, &proposal)
	require.NoError(t, err)

	require.Equal(t, "Dollary-doos LP Cellar Proposal", proposal.Title)
	require.Equal(t, "I have a hunch", proposal.Description)
	require.Equal(t, "0x123801a7D398351b8bE11C439e05C5B3259aeC9B", proposal.CellarIds.Ids[0])
	require.Equal(t, "0x456801a7D398351b8bE11C439e05C5B3259aeC9B", proposal.CellarIds.Ids[1])
	require.Equal(t, "1000stake", proposal.Deposit)
}

func TestParseSubmitScheduledCorkProposal(t *testing.T) {
	encodingConfig := params.MakeTestEncodingConfig()

	okJSON := testutil.WriteToNewTempFile(t, `
{
  "title": "Scheduled cork proposal",
  "description": "I have a hunch",
  "contract_call_proto_json": "{\"cellar_id\":\"0x123801a7D398351b8bE11C439e05C5B3259aeC9B\",\"cellar_v1\":{\"some_fuction\":{\"function_args\":{}},\"block_height\":12345}}",
  "deposit": "1000stake"
}
`)

	proposal := types.AxelarScheduledCorkProposalWithDeposit{}
	contents, err := ioutil.ReadFile(okJSON.Name())
	require.NoError(t, err)

	err = encodingConfig.Marshaler.UnmarshalJSON(contents, &proposal)
	require.NoError(t, err)

	require.Equal(t, "Scheduled cork proposal", proposal.Title)
	require.Equal(t, "I have a hunch", proposal.Description)
	require.Equal(t, "{\"cellar_id\":\"0x123801a7D398351b8bE11C439e05C5B3259aeC9B\",\"cellar_v1\":{\"some_fuction\":{\"function_args\":{}},\"block_height\":12345}}", proposal.ContractCallProtoJson)
	require.Equal(t, "1000stake", proposal.Deposit)
}
