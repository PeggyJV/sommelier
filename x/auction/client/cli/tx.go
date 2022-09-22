package cli

import (
	"fmt"
	"io/ioutil"
	"strconv"
	"strings"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
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

// GetCmdSubmitSetTokenPricesProposal implements the command to submit a token price set proposal
func GetCmdSubmitSetTokenPricesProposal() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "set-token-prices [proposal-file]",
		Args:  cobra.ExactArgs(1),
		Short: "Submit a set token prices proposal",
		Long: strings.TrimSpace(
			fmt.Sprintf(`Submit a set token prices proposal along with an initial deposit.
The proposal details must be supplied via a JSON file.

Example:
$ %s tx gov submit-proposal set-token-prices <path/to/proposal.json> --from=<key_or_address>

Where proposal.json contains:

{
  "title": "Best token price proposal evar",
  "description": "Add the guac",
  "token_prices": [ { denom: "usomm", usd_price: 1000000 }, { denom: "gravity0xdac17f958d2ee523a2206206994597c13d831ec7", usd_price: 0.12501 } ],
  "deposit": "10000usommm"
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

			content := types.NewSetTokenPricesProposal(
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

// GetCmdSubmitBid implements the command to submit a bid
func GetCmdSubmitBid() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "submit-bid [auction-id] [max-bid-in-usomm] [sale-token-minimum-amount]",
		Aliases: []string{"b", "bid"},
		Args:    cobra.ExactArgs(3),
		Short:   "Submit a bid for an auction",
		Long: strings.TrimSpace(
			fmt.Sprintf(`Example:
$ %s tx auction submit-bid 1 10000usomm 50000gravity0xdac17f958d2ee523a2206206994597c13d831ec7 --from=<key_or_address>
`, version.AppName),
		),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			auctionID, err := strconv.ParseUint(args[0], 10, 32)
			if err != nil {
				return err
			}

			maxBidInUsomm, err := sdk.ParseCoinNormalized(args[1])
			if err != nil {
				return err
			}

			saleTokenMinimumAmount, err := sdk.ParseCoinNormalized(args[2])
			if err != nil {
				return err
			}

			bidder := clientCtx.GetFromAddress()
			if bidder == nil {
				return fmt.Errorf("must include `--from` flag")
			}

			msg, err := types.NewMsgSubmitBidRequest(uint32(auctionID), maxBidInUsomm, saleTokenMinimumAmount, bidder)
			if err != nil {
				return err
			}

			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}

	flags.AddTxFlagsToCmd(cmd)
	return cmd
}
