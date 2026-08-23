package cli

import (
	"os"
	"strings"
	"testing"

	"github.com/peggyjv/sommelier/v10/app/params"
	types "github.com/peggyjv/sommelier/v10/x/cork/types/v2"

	"github.com/cosmos/cosmos-sdk/testutil"
	moduletestutil "github.com/cosmos/cosmos-sdk/types/module/testutil"
	"github.com/stretchr/testify/require"
)

// The SDK config defaults to the "cosmos" bech32 prefix. The app calls
// SetAddressPrefixes at init; this package does not, so without it every
// somm1... address fails ValidateBasic with a misleading "expected cosmos".
func TestMain(m *testing.M) {
	params.SetAddressPrefixes()
	os.Exit(m.Run())
}

func TestParseAddManagedCellarsProposal(t *testing.T) {
	encodingConfig := moduletestutil.MakeTestEncodingConfig()

	okJSON := testutil.WriteToNewTempFile(t, `
{
  "title": "Dollary-doos LP Cellar Proposal",
  "description": "I have a hunch",
  "cellar_ids": ["0x123801a7D398351b8bE11C439e05C5B3259aeC9B", "0x456801a7D398351b8bE11C439e05C5B3259aeC9B"],
  "publisher_domain": "example.com",
  "deposit": "1000stake"
}
`)

	proposal := types.AddManagedCellarIDsProposalWithDeposit{}
	contents, err := os.ReadFile(okJSON.Name())
	require.NoError(t, err)

	err = encodingConfig.Codec.UnmarshalJSON(contents, &proposal)
	require.NoError(t, err)

	require.Equal(t, "Dollary-doos LP Cellar Proposal", proposal.Title)
	require.Equal(t, "I have a hunch", proposal.Description)
	require.Equal(t, "0x123801a7D398351b8bE11C439e05C5B3259aeC9B", proposal.CellarIds[0])
	require.Equal(t, "0x456801a7D398351b8bE11C439e05C5B3259aeC9B", proposal.CellarIds[1])
	require.Equal(t, "example.com", proposal.PublisherDomain)
	require.Equal(t, "1000stake", proposal.Deposit)
}

func TestParseRemoveManagedCellarsProposal(t *testing.T) {
	encodingConfig := moduletestutil.MakeTestEncodingConfig()

	okJSON := testutil.WriteToNewTempFile(t, `
{
  "title": "Dollary-doos LP Cellar Proposal",
  "description": "I have a hunch",
  "cellar_ids": ["0x123801a7D398351b8bE11C439e05C5B3259aeC9B", "0x456801a7D398351b8bE11C439e05C5B3259aeC9B"],
  "deposit": "1000stake"
}
`)

	proposal := types.RemoveManagedCellarIDsProposalWithDeposit{}
	contents, err := os.ReadFile(okJSON.Name())
	require.NoError(t, err)

	err = encodingConfig.Codec.UnmarshalJSON(contents, &proposal)
	require.NoError(t, err)

	require.Equal(t, "Dollary-doos LP Cellar Proposal", proposal.Title)
	require.Equal(t, "I have a hunch", proposal.Description)
	require.Equal(t, "0x123801a7D398351b8bE11C439e05C5B3259aeC9B", proposal.CellarIds[0])
	require.Equal(t, "0x456801a7D398351b8bE11C439e05C5B3259aeC9B", proposal.CellarIds[1])
	require.Equal(t, "1000stake", proposal.Deposit)
}

func TestParseSubmitScheduledCorkProposal(t *testing.T) {
	encodingConfig := moduletestutil.MakeTestEncodingConfig()

	okJSON := testutil.WriteToNewTempFile(t, `
{
  "title": "Scheduled cork proposal",
  "description": "I have a hunch",
  "contract_call_proto_json": "{\"cellar_id\":\"0x123801a7D398351b8bE11C439e05C5B3259aeC9B\",\"cellar_v1\":{\"some_fuction\":{\"function_args\":{}},\"block_height\":12345}}",
  "deposit": "1000stake"
}
`)

	proposal := types.ScheduledCorkProposalWithDeposit{}
	contents, err := os.ReadFile(okJSON.Name())
	require.NoError(t, err)

	err = encodingConfig.Codec.UnmarshalJSON(contents, &proposal)
	require.NoError(t, err)

	require.Equal(t, "Scheduled cork proposal", proposal.Title)
	require.Equal(t, "I have a hunch", proposal.Description)
	require.Equal(t, "{\"cellar_id\":\"0x123801a7D398351b8bE11C439e05C5B3259aeC9B\",\"cellar_v1\":{\"some_fuction\":{\"function_args\":{}},\"block_height\":12345}}", proposal.ContractCallProtoJson)
	require.Equal(t, "1000stake", proposal.Deposit)
}

// The cork authority needs a direct path to schedule a cork. Before v10 this
// was steward's job via the validator-delegate path, which v10 retires: the
// msg server now requires signer == params.CorkAuthority, so steward's corks
// are rejected and there was no CLI able to send one.
func TestGetTxCmdRegistersScheduleCork(t *testing.T) {
	cmd := GetTxCmd()
	var found bool
	for _, c := range cmd.Commands() {
		if strings.HasPrefix(c.Use, "schedule-cork") {
			found = true
		}
	}
	require.True(t, found,
		"tx cork must expose schedule-cork; without it the cork authority has no way "+
			"to schedule an Ethereum cork and the cellar wind-down cannot be driven")
}

func TestBuildScheduleCorkMsg(t *testing.T) {
	const (
		cellar = "0x123801a7D398351b8bE11C439e05C5B3259aeC9B"
		signer = "somm1lcsjy2d5s33h0sddd8lpuqvwyz5ruz7ju4aeqa"
	)

	t.Run("hex call is decoded to bytes, not passed through as ASCII", func(t *testing.T) {
		msg, err := buildScheduleCorkMsg(signer, cellar, 1000, "0xdeadbeef")
		require.NoError(t, err)
		// The chain expects raw ABI-encoded bytes. Passing []byte("0xdeadbeef")
		// would hand the cellar the ASCII of the hex string and the call would
		// revert on chain, with nothing locally to indicate why.
		require.Equal(t, []byte{0xde, 0xad, 0xbe, 0xef}, msg.Cork.EncodedContractCall)
		require.Equal(t, cellar, msg.Cork.TargetContractAddress)
		require.Equal(t, uint64(1000), msg.BlockHeight)
		require.Equal(t, signer, msg.Signer)
	})

	t.Run("0x prefix is optional", func(t *testing.T) {
		msg, err := buildScheduleCorkMsg(signer, cellar, 1000, "deadbeef")
		require.NoError(t, err)
		require.Equal(t, []byte{0xde, 0xad, 0xbe, 0xef}, msg.Cork.EncodedContractCall)
	})

	t.Run("rejects a non-hex call", func(t *testing.T) {
		_, err := buildScheduleCorkMsg(signer, cellar, 1000, "not-hex")
		require.Error(t, err)
	})

	t.Run("rejects an invalid cellar address", func(t *testing.T) {
		_, err := buildScheduleCorkMsg(signer, "0xnothex", 1000, "0xdeadbeef")
		require.Error(t, err)
	})

	t.Run("rejects an empty call", func(t *testing.T) {
		_, err := buildScheduleCorkMsg(signer, cellar, 1000, "0x")
		require.Error(t, err, "an empty call body is never a valid cellar instruction")
	})
}
