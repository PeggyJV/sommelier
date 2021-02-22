package cmd

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/cosmos/cosmos-sdk/client"
	kr "github.com/cosmos/cosmos-sdk/crypto/keyring"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/machinebox/graphql"
	oracletypes "github.com/peggyjv/sommelier/x/oracle/types"
	"github.com/spf13/cobra"
	rpchttp "github.com/tendermint/tendermint/rpc/client/http"
)

// queryCmd represents the keys command
func queryCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "query",
		Aliases: []string{"q"},
		Short:   "queries that can be run by the oracle-feeder",
	}

	cmd.AddCommand(
		queryUniswapData(),
		queryParams(),
		queryDelegateAddress(),
		queryValidatorAddress(),
		queryOracleDataPrevote(),
		queryOracleDataVote(),
		queryVotePeriod(),
		queryMissCounter(),
		queryOracleData(),
	)

	return cmd
}

func queryUniswapData() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "uniswap-data",
		Aliases: []string{"ud"},
		Args:    cobra.NoArgs,
		Short:   "queries uniswap data for the transactions",
		RunE: func(cmd *cobra.Command, _ []string) error {
			n, err := cmd.Flags().GetInt("num-markets")
			if err != nil {
				return err
			}

			out, err := config.GetPairs(context.Background(), n, 0)
			if err != nil {
				return err
			}

			bz, err := json.Marshal(out)
			if err != nil {
				return err
			}

			fmt.Println(string(bz))
			return nil
		},
	}
	cmd.Flags().IntP("num-markets", "n", 5, "number of markets to query")
	return cmd
}

func queryParams() *cobra.Command {
	return &cobra.Command{
		Use:     "parameters",
		Aliases: []string{"params"},
		Args:    cobra.NoArgs,
		Short:   "query oracle params from the chain",
		RunE: func(cmd *cobra.Command, _ []string) error {
			ctx, err := config.GetClientContext(cmd)
			if err != nil {
				return err
			}

			params, err := GetParams(ctx)
			if err != nil {
				return err
			}

			bz, err := json.Marshal(params)
			if err != nil {
				return err
			}

			fmt.Println(string(bz))
			return nil
		},
	}
}

// GetParams helper
func GetParams(ctx client.Context) (oracletypes.Params, error) {
	res, err := oracletypes.NewQueryClient(ctx).QueryParams(context.Background(), &oracletypes.QueryParamsRequest{})
	return res.Params, err
}

func queryDelegateAddress() *cobra.Command {
	return &cobra.Command{
		Use:     "delegate-address [validator-address]",
		Aliases: []string{"del"},
		Args:    cobra.ExactArgs(1),
		Short:   "query delegate address from the chain given validators address",
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx, err := config.GetClientContext(cmd)
			if err != nil {
				return err
			}

			val, err := sdk.AccAddressFromBech32(args[0])
			if err != nil {
				return err
			}

			del, err := GetDelFromVal(ctx, val)
			if err != nil {
				return err
			}

			fmt.Println(del.String())
			return nil
		},
	}
}

// GetDelFromVal helper
func GetDelFromVal(ctx client.Context, val sdk.AccAddress) (sdk.AccAddress, error) {
	queryClient := oracletypes.NewQueryClient(ctx)

	req := &oracletypes.QueryDelegateAddressRequest{
		Validator: val.String(),
	}

	res, err := queryClient.QueryDelegateAddress(context.Background(), req)
	if err != nil {
		return nil, err
	}

	return sdk.AccAddressFromBech32(res.Delegate)
}

func queryValidatorAddress() *cobra.Command {
	return &cobra.Command{
		Use:     "validator-address [delegate-address]",
		Aliases: []string{"val"},
		Args:    cobra.ExactArgs(1),
		Short:   "query validator address from the chain given the address that validator delegated to",
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx, err := config.GetClientContext(cmd)
			if err != nil {
				return err
			}

			del, err := sdk.AccAddressFromBech32(args[0])
			if err != nil {
				return err
			}

			val, err := GetValFromDel(ctx, del)
			if err != nil {
				return err
			}

			fmt.Println(val.String())
			return nil
		},
	}
}

// GetValFromDel helper
func GetValFromDel(ctx client.Context, del sdk.AccAddress) (sdk.ValAddress, error) {
	queryClient := oracletypes.NewQueryClient(ctx)

	req := &oracletypes.QueryValidatorAddressRequest{
		Delegate: del.String(),
	}

	res, err := queryClient.QueryValidatorAddress(context.Background(), req)
	if err != nil {
		return nil, err
	}

	return sdk.ValAddressFromBech32(res.Validator)
}

func queryOracleDataPrevote() *cobra.Command {
	return &cobra.Command{
		Use:     "oracle-prevote [signer]",
		Aliases: []string{"prevote", "pv"},
		Args:    cobra.ExactArgs(1),
		Short:   "query oracle data prevote from the chain, by either validator or delegate address",
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx, err := config.GetClientContext(cmd)
			if err != nil {
				return err
			}

			val, err := sdk.AccAddressFromBech32(args[0])
			if err != nil {
				return err
			}

			hashes, err := GetPrevote(ctx, val)
			if err != nil {
				return err
			}

			bz, err := json.Marshal(hashes)
			if err != nil {
				return err
			}

			fmt.Println(string(bz))
			return nil
		},
	}
}

// GetPrevote helper
func GetPrevote(ctx client.Context, val sdk.AccAddress) (*oracletypes.OraclePrevote, error) {
	res, err := oracletypes.NewQueryClient(ctx).QueryOracleDataPrevote(
		context.Background(), &oracletypes.QueryOracleDataPrevoteRequest{Validator: val.String()})
	return res.Prevote, err
}

func queryOracleDataVote() *cobra.Command {
	return &cobra.Command{
		Use:     "oracle-vote [signer]",
		Aliases: []string{"vote"},
		Args:    cobra.ExactArgs(1),
		Short:   "query oracle data vote from the chain, by either validator or delegate address",
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx, err := config.GetClientContext(cmd)
			if err != nil {
				return err
			}

			val, err := sdk.AccAddressFromBech32(args[0])
			if err != nil {
				return err
			}

			vote, err := GetVote(ctx, val)
			if err != nil {
				return err
			}

			bz, err := ctx.JSONMarshaler.MarshalJSON(vote)
			if err != nil {
				return err
			}

			fmt.Println(string(bz))
			return nil
		},
	}
}

// GetVote helper
func GetVote(ctx client.Context, val sdk.AccAddress) (*oracletypes.QueryOracleDataVoteResponse, error) {
	queryClient := oracletypes.NewQueryClient(ctx)
	req := &oracletypes.QueryOracleDataVoteRequest{Validator: val.String()}

	return queryClient.QueryOracleDataVote(context.Background(), req)
}

func queryVotePeriod() *cobra.Command {
	return &cobra.Command{
		Use:     "vote-period",
		Aliases: []string{"vp"},
		Args:    cobra.NoArgs,
		Short:   "query vote period data from the chain",
		RunE: func(cmd *cobra.Command, _ []string) error {
			ctx, err := config.GetClientContext(cmd)
			if err != nil {
				return err
			}
			res, err := GetVotePeriod(ctx)
			if err != nil {
				return err
			}

			bz, err := json.Marshal(res)
			if err != nil {
				return err
			}

			fmt.Println(string(bz))
			return nil
		},
	}
}

// GetVotePeriod helper
func GetVotePeriod(ctx client.Context) (*oracletypes.VotePeriod, error) {
	queryClient := oracletypes.NewQueryClient(ctx)

	return queryClient.QueryVotePeriod(context.Background(), &oracletypes.QueryVotePeriodRequest{})
}

func queryMissCounter() *cobra.Command {
	return &cobra.Command{
		Use:     "miss-counter [signer]",
		Aliases: []string{"mc"},
		Args:    cobra.ExactArgs(1),
		Short:   "query miss counter for a validator from the chain given its address or the delegate address",
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx, err := config.GetClientContext(cmd)
			if err != nil {
				return err
			}

			val, err := sdk.AccAddressFromBech32(args[0])
			if err != nil {
				return err
			}

			misses, err := GetMissCounter(ctx, val)
			if err != nil {
				return err
			}

			fmt.Println(misses)
			return nil
		},
	}
}

// GetMissCounter helper
func GetMissCounter(ctx client.Context, val sdk.AccAddress) (int64, error) {
	queryClient := oracletypes.NewQueryClient(ctx)

	req := &oracletypes.QueryMissCounterRequest{Validator: val.String()}

	res, err := queryClient.QueryMissCounter(context.Background(), req)
	if err != nil {
		return 0, err
	}

	return res.MissCounter, nil
}

func queryOracleData() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "oracle-data",
		Aliases: []string{"od"},
		Args:    cobra.NoArgs,
		Short:   "query consensus oracle data from the chain given its type",
		RunE: func(cmd *cobra.Command, _ []string) error {
			ctx, err := config.GetClientContext(cmd)
			if err != nil {
				return err
			}

			typ, err := cmd.Flags().GetString("type")
			if err != nil {
				return err
			}

			res, err := oracletypes.NewQueryClient(ctx).OracleData(context.Background(), &oracletypes.QueryOracleDataRequest{Type: typ})
			if err != nil {
				return err
			}

			od, err := oracletypes.UnpackOracleData(res.OracleData)
			if err != nil {
				return err
			}

			jsonBz, err := json.Marshal(od)
			if err != nil {
				return err
			}

			fmt.Println(jsonBz)
			return nil
		},
	}
	cmd.Flags().StringP("type", "t", oracletypes.UniswapDataType, "type of oracle data to fetch")
	return cmd
}

// GetData helper
func GetData(ctx client.Context, typ string) (oracletypes.OracleData, error) {
	res, err := oracletypes.NewQueryClient(ctx).OracleData(context.Background(), &oracletypes.QueryOracleDataRequest{Type: typ})
	if err != nil {
		return nil, err
	}
	return oracletypes.UnpackOracleData(res.OracleData)
}

// GetPairs returns the top N pairs from the Uniswap Subgraph
func (c *Config) GetPairs(ctx context.Context, first, skip int) (*oracletypes.OracleFeed, error) {
	req := graphql.NewRequest(fmt.Sprintf(`{ 
		pairs(first: %d, skip: %d, orderBy: volumeUSD, orderDirection: desc) {
			id
			reserveUSD
			totalSupply
			reserve0
			reserve1
			token0Price
			token1Price
			token0 {
				id
				decimals
			}
			token1 {
				id
				decimals
			}
		}
	}`, first, skip))

	var out *oracletypes.OracleFeed

	err := c.graphClient.Run(ctx, req, out)
	if err != nil {
		return nil, err
	}

	return out, nil
}

// GetClientContext reads in values from the config
func (c Config) GetClientContext(cmd *cobra.Command) (client.Context, error) {
	ctx := client.GetClientContextFromCmd(cmd)
	home, err := cmd.Flags().GetString("home")
	if err != nil {
		return client.Context{}, err
	}

	cl, err := rpchttp.New(c.ChainRPC, "/websocket")
	if err != nil {
		return client.Context{}, err
	}

	keyring, err := kr.New(AppName, "test", home, os.Stdin)
	if err != nil {
		return client.Context{}, err
	}

	return ctx.WithClient(cl).
		WithChainID(c.ChainID).
		WithFromName(c.SigningKey).
		WithFrom(c.SigningKey).
		WithKeyring(keyring).
		WithOutput(ioutil.Discard).
		WithHomeDir(home), nil
}
