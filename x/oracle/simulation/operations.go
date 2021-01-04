package simulation

// DONTCOVER

import (
	"fmt"
	"math/rand"

	"github.com/cosmos/cosmos-sdk/baseapp"
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/simapp/helpers"
	simappparams "github.com/cosmos/cosmos-sdk/simapp/params"
	sdk "github.com/cosmos/cosmos-sdk/types"
	simtypes "github.com/cosmos/cosmos-sdk/types/simulation"
	authkeeper "github.com/cosmos/cosmos-sdk/x/auth/keeper"
	bankkeeper "github.com/cosmos/cosmos-sdk/x/bank/keeper"
	"github.com/cosmos/cosmos-sdk/x/simulation"

	"github.com/peggyjv/sommelier/x/oracle/keeper"
	"github.com/peggyjv/sommelier/x/oracle/types"
)

// Simulation operation weights constants
const (
	OpWeightMsgAggregateExchangeRatePrevote = "op_weight_msg_aggregate_exchange_rate_prevote"
	OpWeightMsgAggregateExchangeRateVote    = "op_weight_msg_aggregate_exchange_rate_vote"
	OpWeightMsgDelegateFeedConsent          = "op_weight_msg_exchange_feed_consent"

	salt = "1234"
)

var (
	whitelist                      = []string{types.MicroKRWDenom, types.MicroUSDDenom, types.MicroSDRDenom, types.MicroMNTDenom}
	voteHashMap map[string]sdk.Dec = make(map[string]sdk.Dec)
)

// WeightedOperations returns all the operations from the module with their respective weights
func WeightedOperations(
	appParams simtypes.AppParams,
	cdc codec.JSONMarshaler,
	ak authkeeper.AccountKeeper,
	bk bankkeeper.Keeper,
	k keeper.Keeper) simulation.WeightedOperations {
	var (
		weightMsgAggregateExchangeRatePrevote int
		weightMsgAggregateExchangeRateVote    int
		weightMsgDelegateFeedConsent          int
	)
	appParams.GetOrGenerate(cdc, OpWeightMsgAggregateExchangeRatePrevote, &weightMsgAggregateExchangeRatePrevote, nil,
		func(_ *rand.Rand) {
			weightMsgAggregateExchangeRatePrevote = simappparams.DefaultWeightMsgSend * 2
		},
	)

	appParams.GetOrGenerate(cdc, OpWeightMsgAggregateExchangeRateVote, &weightMsgAggregateExchangeRateVote, nil,
		func(_ *rand.Rand) {
			weightMsgAggregateExchangeRateVote = simappparams.DefaultWeightMsgSend * 2
		},
	)

	appParams.GetOrGenerate(cdc, OpWeightMsgDelegateFeedConsent, &weightMsgDelegateFeedConsent, nil,
		func(_ *rand.Rand) {
			weightMsgDelegateFeedConsent = simappparams.DefaultWeightMsgSetWithdrawAddress
		},
	)

	return simulation.WeightedOperations{
		simulation.NewWeightedOperation(
			weightMsgAggregateExchangeRatePrevote,
			SimulateMsgAggregateExchangeRatePrevote(ak, bk, k),
		),
		simulation.NewWeightedOperation(
			weightMsgAggregateExchangeRateVote,
			SimulateMsgExchangeRateVote(ak, bk, k),
		),
		simulation.NewWeightedOperation(
			weightMsgDelegateFeedConsent,
			SimulateMsgDelegateFeedConsent(ak, bk, k),
		),
	}
}

// SimulateMsgAggregateExchangeRatePrevote generates a MsgExchangeRatePrevote with random values.
// nolint: funlen
func SimulateMsgAggregateExchangeRatePrevote(ak authkeeper.AccountKeeper, bk bankkeeper.Keeper, k keeper.Keeper) simtypes.Operation {
	return func(
		r *rand.Rand, app *baseapp.BaseApp, ctx sdk.Context, accs []simtypes.Account, chainID string,
	) (simtypes.OperationMsg, []simtypes.FutureOperation, error) {
		fmt.Println("EXCHANGE RATE PREVOTE")
		txcfg := keeper.TestTxConfig()
		simAccount, _ := simtypes.RandomAcc(r, accs)
		address := sdk.ValAddress(simAccount.Address)

		// ensure the validator exists
		val := k.StakingKeeper.Validator(ctx, address)
		if val == nil {
			return simtypes.NoOpMsg(types.ModuleName, "", ""), nil, nil
		}

		denom := whitelist[simtypes.RandIntBetween(r, 0, len(whitelist))]
		price := sdk.NewDecWithPrec(int64(simtypes.RandIntBetween(r, 1, 10000)), int64(1))
		voteHash := types.GetAggregateVoteHash(salt, fmt.Sprintf("%d%s", price, denom), address)

		feederAddr := k.GetOracleDelegate(ctx, address)
		feederSimAccount, _ := simtypes.FindAccount(accs, feederAddr)
		feederAccount := ak.GetAccount(ctx, feederAddr)

		fees, err := simtypes.RandomFees(r, ctx, bk.SpendableCoins(ctx, feederAccount.GetAddress()))
		if err != nil {
			return simtypes.NoOpMsg(types.ModuleName, "", ""), nil, err
		}

		msg := types.NewMsgAggregateExchangeRatePrevote(voteHash, feederAddr, address)

		tx, err := helpers.GenTx(
			txcfg,
			[]sdk.Msg{msg},
			fees,
			helpers.DefaultGenTxGas,
			chainID,
			[]uint64{feederAccount.GetAccountNumber()},
			[]uint64{feederAccount.GetSequence()},
			feederSimAccount.PrivKey,
		)
		if err != nil {
			return simtypes.NoOpMsg(types.ModuleName, "", ""), nil, err
		}

		_, _, err = app.Deliver(txcfg.TxEncoder(), tx)
		if err != nil {
			return simtypes.NoOpMsg(types.ModuleName, "", ""), nil, err
		}

		voteHashMap[denom+address.String()] = price

		return simtypes.NewOperationMsg(msg, true, ""), nil, nil
	}
}

// SimulateMsgExchangeRateVote generates a MsgExchangeRateVote with random values.
// nolint: funlen
func SimulateMsgExchangeRateVote(ak authkeeper.AccountKeeper, bk bankkeeper.Keeper, k keeper.Keeper) simtypes.Operation {
	return func(
		r *rand.Rand, app *baseapp.BaseApp, ctx sdk.Context, accs []simtypes.Account, chainID string,
	) (simtypes.OperationMsg, []simtypes.FutureOperation, error) {
		fmt.Println("EXCHANGE RATE VOTE")
		txcfg := keeper.TestTxConfig()
		simAccount, _ := simtypes.RandomAcc(r, accs)
		address := sdk.ValAddress(simAccount.Address)

		// ensure the validator exists
		val := k.StakingKeeper.Validator(ctx, address)
		if val == nil {
			return simtypes.NoOpMsg(types.ModuleName, "", ""), nil, nil
		}

		// ensure vote hash exists
		denom := whitelist[simtypes.RandIntBetween(r, 0, len(whitelist))]
		price, ok := voteHashMap[denom+address.String()]
		if !ok {
			return simtypes.NoOpMsg(types.ModuleName, "", ""), nil, nil
		}

		// get prevote
		prevote, err := k.GetExchangeRatePrevote(ctx, denom, address)
		if err != nil {
			return simtypes.NoOpMsg(types.ModuleName, "", ""), nil, nil
		}

		params := k.GetParams(ctx)
		if (ctx.BlockHeight()/params.VotePeriod)-(prevote.SubmitBlock/params.VotePeriod) != 1 {
			return simtypes.NoOpMsg(types.ModuleName, "", ""), nil, nil
		}

		feederAddr := k.GetOracleDelegate(ctx, address)
		feederSimAccount, _ := simtypes.FindAccount(accs, feederAddr)
		feederAccount := ak.GetAccount(ctx, feederAddr)

		fees, err := simtypes.RandomFees(r, ctx, bk.SpendableCoins(ctx, feederAccount.GetAddress()))
		if err != nil {
			return simtypes.NoOpMsg(types.ModuleName, "", ""), nil, err
		}

		er := &sdk.DecCoin{Denom: denom, Amount: price}
		msg := types.NewMsgAggregateExchangeRateVote(salt, er.String(), feederAddr, address)

		tx, err := helpers.GenTx(
			txcfg,
			[]sdk.Msg{msg},
			fees,
			helpers.DefaultGenTxGas,
			chainID,
			[]uint64{feederAccount.GetAccountNumber()},
			[]uint64{feederAccount.GetSequence()},
			feederSimAccount.PrivKey,
		)
		if err != nil {
			return simtypes.NoOpMsg(types.ModuleName, "", ""), nil, err
		}

		_, _, err = app.Deliver(txcfg.TxEncoder(), tx)
		if err != nil {
			return simtypes.NoOpMsg(types.ModuleName, "", ""), nil, err
		}

		return simtypes.NewOperationMsg(msg, true, ""), nil, nil
	}
}

// SimulateMsgDelegateFeedConsent generates a MsgDelegateFeedConsent with random values.
// nolint: funlen
func SimulateMsgDelegateFeedConsent(ak authkeeper.AccountKeeper, bk bankkeeper.Keeper, k keeper.Keeper) simtypes.Operation {
	return func(
		r *rand.Rand, app *baseapp.BaseApp, ctx sdk.Context, accs []simtypes.Account, chainID string,
	) (simtypes.OperationMsg, []simtypes.FutureOperation, error) {
		fmt.Println("Delegate Feed Consent")
		txcfg := keeper.TestTxConfig()
		simAccount, _ := simtypes.RandomAcc(r, accs)
		fmt.Println("  - Validator", simAccount.Address)
		delegateAccount, _ := simtypes.RandomAcc(r, accs)
		valAddress := sdk.ValAddress(simAccount.Address)
		delegateValAddress := sdk.ValAddress(delegateAccount.Address)
		fmt.Println("  - Delegate", delegateAccount.Address)
		account := ak.GetAccount(ctx, simAccount.Address)
		// ensure the validator exists
		val := k.StakingKeeper.Validator(ctx, valAddress)
		if val == nil {
			return simtypes.NoOpMsg(types.ModuleName, "", ""), nil, nil
		}

		// ensure the target address is not a validator
		val2 := k.StakingKeeper.Validator(ctx, delegateValAddress)
		if val2 != nil {
			return simtypes.NoOpMsg(types.ModuleName, "", ""), nil, nil
		}

		fees, err := simtypes.RandomFees(r, ctx, bk.SpendableCoins(ctx, account.GetAddress()))
		if err != nil {
			return simtypes.NoOpMsg(types.ModuleName, "", ""), nil, err
		}

		msg := types.NewMsgDelegateFeedConsent(valAddress, delegateAccount.Address)

		tx, err := helpers.GenTx(
			txcfg,
			[]sdk.Msg{msg},
			fees,
			helpers.DefaultGenTxGas,
			chainID,
			[]uint64{account.GetAccountNumber()},
			[]uint64{account.GetSequence()},
			delegateAccount.PrivKey,
		)
		if err != nil {
			return simtypes.NoOpMsg(types.ModuleName, "", ""), nil, err
		}

		_, _, err = app.Deliver(txcfg.TxEncoder(), tx)
		if err != nil {
			return simtypes.NoOpMsg(types.ModuleName, "", ""), nil, err
		}

		return simtypes.NewOperationMsg(msg, true, ""), nil, nil
	}
}
