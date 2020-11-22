package simulation

// DONTCOVER

import (
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
	core "github.com/peggyjv/sommelier/x/oracle/types"
)

// Simulation operation weights constants
const (
	OpWeightMsgExchangeRatePrevote = "op_weight_msg_exchange_rate_prevote"
	OpWeightMsgExchangeRateVote    = "op_weight_msg_exchange_rate_vote"
	OpWeightMsgDelegateFeedConsent = "op_weight_msg_exchange_feed_consent"

	salt = "1234"
)

var (
	whitelist                      = []string{core.MicroKRWDenom, core.MicroUSDDenom, core.MicroSDRDenom, core.MicroMNTDenom}
	voteHashMap map[string]sdk.Dec = make(map[string]sdk.Dec)
)

// WeightedOperations returns all the operations from the module with their respective weights
func WeightedOperations(
	appParams simtypes.AppParams,
	cdc codec.JSONMarshaler,
	ak authkeeper.AccountKeeper,
	bk bankkeeper.BaseKeeper,
	k keeper.Keeper) simulation.WeightedOperations {
	var (
		weightMsgExchangeRatePrevote int
		weightMsgExchangeRateVote    int
		weightMsgDelegateFeedConsent int
	)
	appParams.GetOrGenerate(cdc, OpWeightMsgExchangeRatePrevote, &weightMsgExchangeRatePrevote, nil,
		func(_ *rand.Rand) {
			weightMsgExchangeRatePrevote = simappparams.DefaultWeightMsgSend * 2
		},
	)

	appParams.GetOrGenerate(cdc, OpWeightMsgExchangeRateVote, &weightMsgExchangeRateVote, nil,
		func(_ *rand.Rand) {
			weightMsgExchangeRateVote = simappparams.DefaultWeightMsgSend * 2
		},
	)

	appParams.GetOrGenerate(cdc, OpWeightMsgDelegateFeedConsent, &weightMsgDelegateFeedConsent, nil,
		func(_ *rand.Rand) {
			weightMsgDelegateFeedConsent = simappparams.DefaultWeightMsgSetWithdrawAddress
		},
	)

	return simulation.WeightedOperations{
		// TODO: simulate aggregatePrevote
		// TODO: simulate aggregateVote
		simulation.NewWeightedOperation(
			weightMsgDelegateFeedConsent,
			SimulateMsgDelegateFeedConsent(ak, bk, k, cdc),
		),
	}
}

// SimulateMsgDelegateFeedConsent generates a MsgDelegateFeedConsent with random values.
// nolint: funlen
func SimulateMsgDelegateFeedConsent(ak authkeeper.AccountKeeper, bk bankkeeper.BaseKeeper, k keeper.Keeper, cdc codec.JSONMarshaler) simtypes.Operation {
	return func(
		r *rand.Rand, app *baseapp.BaseApp, ctx sdk.Context, accs []simtypes.Account, chainID string,
	) (simtypes.OperationMsg, []simtypes.FutureOperation, error) {
		txcfg := keeper.TestTxConfig()
		simAccount, _ := simtypes.RandomAcc(r, accs)
		delegateAccount, _ := simtypes.RandomAcc(r, accs)
		valAddress := sdk.ValAddress(simAccount.Address)
		delegateValAddress := sdk.ValAddress(delegateAccount.Address)
		account := ak.GetAccount(ctx, simAccount.Address)

		// ensure the validator exists
		val := k.StakingKeeper.Validator(ctx, valAddress)
		if val == nil {
			// TODO: maybe this needs to be something else?
			return simtypes.NoOpMsg(types.ModuleName, "", ""), nil, nil
		}

		// ensure the target address is not a validator
		val2 := k.StakingKeeper.Validator(ctx, delegateValAddress)
		if val2 != nil {
			// TODO: maybe this needs to be something else?
			return simtypes.NoOpMsg(types.ModuleName, "", ""), nil, nil
		}
		fees, err := simtypes.RandomFees(r, ctx, bk.SpendableCoins(ctx, account.GetAddress()))
		if err != nil {
			// TODO: maybe this needs to be something else?
			return simtypes.NoOpMsg(types.ModuleName, "", ""), nil, err
		}

		msg := types.NewMsgDelegateFeedConsent(valAddress, delegateAccount.Address)

		// TODO: create txConfig
		tx, err := helpers.GenTx(
			txcfg,
			[]sdk.Msg{&msg},
			fees,
			helpers.DefaultGenTxGas,
			chainID,
			[]uint64{account.GetAccountNumber()},
			[]uint64{account.GetSequence()},
			simAccount.PrivKey,
		)
		if err != nil {
			return simtypes.NoOpMsg(types.ModuleName, "", ""), nil, err
		}

		_, _, err = app.Deliver(txcfg.TxEncoder(), tx)
		if err != nil {
			return simtypes.NoOpMsg(types.ModuleName, "", ""), nil, err
		}

		return simtypes.NewOperationMsg(&msg, true, ""), nil, nil
	}
}
