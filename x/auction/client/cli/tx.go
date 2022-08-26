package cli

import (
	"fmt"
	"io/ioutil"
	"strings"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/tx"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/version"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"
	types "github.com/peggyjv/sommelier/v4/x/auction/types"
	"github.com/spf13/cobra"
)

// GetTxCmd returns the transaction commands for this module
func GetTxCmd() *cobra.Command {
	auctionTxCmd := &cobra.Command{
		Use:                        "auction",
		Short:                      "Auction transaction subcommands",
		DisableFlagParsing:         true,
		SuggestionsMinimumDistance: 2,
		RunE:                       client.ValidateCmd,
	}

	return auctionTxCmd
}

// GetCmdSubmitAddProposal implements the command to submit a token update proposal
func GetCmdSubmitAddProposal() *cobra.Command {

	cmd := &cobra.Command{
		Use:   "update-token-prices [proposal-file]",
		Args:  cobra.ExactArgs(1),
		Short: "Submit a token update proposal",
		Long: strings.TrimSpace(
			fmt.Sprintf(`Submit a token update proposal along with an initial deposit.
The proposal details must be supplied via a JSON file.

Example:
$ %s tx gov submit-proposal update-token-prices <path/to/proposal.json> --from=<key_or_address>

Where proposal.json contains:

{
  "title": "Best token price proposal evar",
  "description": "Add the guac",
  "token_prices": [ { denom: "usomm", usd_price: 1000000 }, {denom: "gwei", usd_price: 0.12501 }],
  "deposit": "1000stake"
}
`,
				version.AppName,
			),
		),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			proposal := types.SetTokenPricesProposalWithDeposit{}
			contents, err := ioutil.ReadFile(args[0])
			if err != nil {
				return err
			}

			if err = clientCtx.Codec.UnmarshalJSON(contents, &proposal); err != nil {
				return err
			}

			deposit, err := sdk.ParseCoinsNormalized(proposal.Deposit)
			if err != nil {
				return err
			}

			for _, tokenPrice := range proposal.GetTokenPrices() {
				err := tokenPrice.ValidateBasic()
				if err != nil {
					return err
				}
			}

			content := types.NewAddSetTokenPricesProposal(
				proposal.Title,
				proposal.Description,
				proposal.TokenPrices,
			)

			from := clientCtx.GetFromAddress()
			msg, err := govtypes.NewMsgSubmitProposal(content, deposit, from)
			if err != nil {
				return err
			}

			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}

	return cmd
}
