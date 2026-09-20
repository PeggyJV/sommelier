package app

import (
	"encoding/json"
	"testing"
	"time"

	dbm "github.com/cometbft/cometbft-db"
	abci "github.com/cometbft/cometbft/abci/types"
	"github.com/cometbft/cometbft/libs/log"
	tmproto "github.com/cometbft/cometbft/proto/tendermint/types"
	tmtypes "github.com/cometbft/cometbft/types"
	"github.com/cosmos/cosmos-sdk/testutil/mock"
	sdk "github.com/cosmos/cosmos-sdk/types"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	evidencetypes "github.com/cosmos/cosmos-sdk/x/evidence/types"
	"github.com/cosmos/cosmos-sdk/x/slashing"
	slashingtypes "github.com/cosmos/cosmos-sdk/x/slashing/types"
	"github.com/cosmos/cosmos-sdk/x/staking"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	poakeeper "github.com/peggyjv/sommelier/v10/x/poa/keeper"
	poatypes "github.com/peggyjv/sommelier/v10/x/poa/types"
	"github.com/stretchr/testify/require"
)

// Exercise real SDK handlers, staking and bank state. The consensus power is
// taken from the real PoA EndBlocker output and delivered two blocks later,
// matching CometBFT 0.37's validator-update pipeline.
func TestPoASlashMultiplierTransition(t *testing.T) {
	for _, path := range []string{"downtime", "equivocation", "gravity"} {
		for _, transition := range []string{"boost_removed", "boost_added", "stable"} {
			t.Run(path+"/"+transition, func(t *testing.T) {
				app, ctx, auth, community := newPoASlashApp(t)
				boostedPower := int64(8122)
				if path == "gravity" {
					// Isolate snapshot timing from the separate concrete-token
					// rounding in Gravity's ConsensusPower read.
					params := app.PoaKeeper.GetParams(ctx)
					params.FloorFraction = sdk.MustNewDecFromStr("0.75")
					app.PoaKeeper.SetParams(ctx, params)
					boostedPower = 12000
				}
				comCons, err := community.GetConsAddr()
				require.NoError(t, err)
				if transition == "boost_added" {
					app.StakingKeeper.Jail(ctx, comCons)
				}
				runEnd := func(height int64) []abci.ValidatorUpdate {
					ctx = ctx.WithBlockHeight(height)
					return poakeeper.EndBlocker(ctx, app.PoaKeeper, func(ctx sdk.Context) []abci.ValidatorUpdate {
						return staking.EndBlocker(ctx, &app.StakingKeeper)
					})
				}
				updates := runEnd(97)
				authKey, err := auth.ConsPubKey()
				require.NoError(t, err)
				// An unchanged raw validator can be absent from the update list.
				power := int64(1000)
				for _, update := range updates {
					if string(update.PubKey.GetEd25519()) == string(authKey.Bytes()) {
						power = update.Power
					}
				}
				if transition == "boost_added" {
					require.Equal(t, int64(1000), power)
				} else {
					require.Equal(t, boostedPower, power)
				}
				if transition == "boost_removed" {
					app.StakingKeeper.Jail(ctx, comCons)
				}
				if transition == "boost_added" {
					app.StakingKeeper.Unjail(ctx, comCons)
				}
				runEnd(98)
				authCons, err := auth.GetConsAddr()
				require.NoError(t, err)
				// At height 100, LastCommit is for 99, whose power came from EndBlock 97.
				// Evidence for height 99 uses the same power and distribution height 98.
				ctx = ctx.WithBlockHeight(100)
				supplyBefore := app.BankKeeper.GetSupply(ctx, sdk.DefaultBondDenom).Amount
				switch path {
				case "downtime":
					slashing.BeginBlocker(ctx, abci.RequestBeginBlock{LastCommitInfo: abci.CommitInfo{Votes: []abci.VoteInfo{{Validator: abci.Validator{Address: authCons, Power: power}, SignedLastBlock: false}}}}, app.SlashingKeeper)
				case "equivocation":
					app.EvidenceKeeper.HandleEquivocationEvidence(ctx, &evidencetypes.Equivocation{Height: 99, Time: ctx.BlockTime(), Power: power, ConsensusAddress: authCons.String()})
					require.True(t, app.SlashingKeeper.IsTombstoned(ctx, authCons))
				case "gravity":
					// Gravity runs after PoA's EndBlocker and uses current, adapted power.
					runEnd(100)
					app.GravityKeeper.SlashAndJail(ctx, auth, "test")
				}
				after, found := app.StakingKeeper.GetValidator(ctx, auth.GetOperator())
				require.True(t, found)
				require.True(t, after.Jailed)
				// 5% of 1000 raw power * 1e6 tokens/power = 50,000,000 tokens.
				require.Equal(t, "950000000", after.Tokens.String())
				require.Equal(t, "50000000", supplyBefore.Sub(app.BankKeeper.GetSupply(ctx, sdk.DefaultBondDenom).Amount).String())
			})
		}
	}
}

func newPoASlashApp(t *testing.T) (*SommelierApp, sdk.Context, stakingtypes.Validator, stakingtypes.Validator) {
	t.Helper()
	app := NewSommelierApp(log.NewNopLogger(), dbm.NewMemDB(), nil, true, map[int64]bool{}, t.TempDir(), 0, MakeEncodingConfig(), EmptyAppOptions{})
	vals := make([]*tmtypes.Validator, 2)
	for i := range vals {
		pv := mock.NewPV()
		key, err := pv.GetPubKey()
		require.NoError(t, err)
		vals[i] = tmtypes.NewValidator(key, 1)
	}
	acc := authtypes.NewBaseAccountWithAddress(sdk.AccAddress(vals[0].Address))
	genesis := genesisStateWithValSet(t, app, NewDefaultGenesisState(), tmtypes.NewValidatorSet(vals), []authtypes.GenesisAccount{acc})
	var sg stakingtypes.GenesisState
	app.AppCodec().MustUnmarshalJSON(genesis[stakingtypes.ModuleName], &sg)
	for i := range sg.Validators {
		tokens := sdk.NewInt(1000000000)
		if i == 1 {
			tokens = sdk.NewInt(4000000000)
		}
		sg.Validators[i].Tokens = tokens
		sg.Validators[i].DelegatorShares = sdk.NewDecFromInt(tokens)
		sg.Delegations[i].Shares = sdk.NewDecFromInt(tokens)
	}
	genesis[stakingtypes.ModuleName] = app.AppCodec().MustMarshalJSON(&sg)
	bg := banktypes.DefaultGenesisState()
	bg.Balances = []banktypes.Balance{{Address: authtypes.NewModuleAddress(stakingtypes.BondedPoolName).String(), Coins: sdk.NewCoins(sdk.NewInt64Coin(sdk.DefaultBondDenom, 5000000000))}}
	bg.Supply = sdk.NewCoins(sdk.NewInt64Coin(sdk.DefaultBondDenom, 5000000000))
	genesis[banktypes.ModuleName] = app.AppCodec().MustMarshalJSON(bg)
	pg := poatypes.DefaultGenesis()
	pg.Params = poatypes.DefaultParams()
	pg.AuthoritySet = []poatypes.AuthorityValidator{{OperatorAddress: sg.Validators[0].OperatorAddress}}
	genesis[poatypes.ModuleName] = app.AppCodec().MustMarshalJSON(pg)
	state, err := json.Marshal(genesis)
	require.NoError(t, err)
	now := time.Unix(1700000000, 0).UTC()
	app.InitChain(abci.RequestInitChain{Time: now, ConsensusParams: DefaultConsensusParams, AppStateBytes: state})
	app.Commit()
	ctx := app.BaseApp.NewContext(true, tmproto.Header{Height: 97, Time: now.Add(time.Minute)})
	for _, val := range sg.Validators {
		pubkey, err := val.ConsPubKey()
		require.NoError(t, err)
		require.NoError(t, app.SlashingKeeper.AddPubkey(ctx, pubkey))
		cons, err := val.GetConsAddr()
		require.NoError(t, err)
		app.SlashingKeeper.SetValidatorSigningInfo(ctx, cons, slashingtypes.ValidatorSigningInfo{Address: cons.String(), StartHeight: 0})
	}
	sp := slashingtypes.DefaultParams()
	sp.SignedBlocksWindow = 1
	sp.MinSignedPerWindow = sdk.OneDec()
	sp.SlashFractionDowntime = sdk.MustNewDecFromStr("0.05")
	sp.SlashFractionDoubleSign = sdk.MustNewDecFromStr("0.05")
	require.NoError(t, app.SlashingKeeper.SetParams(ctx, sp))
	gp := app.GravityKeeper.GetParams(ctx)
	gp.SlashFractionBatch = sdk.MustNewDecFromStr("0.05")
	app.GravityKeeper.SetParams(ctx, gp)
	return app, ctx, sg.Validators[0], sg.Validators[1]
}

// Moving the height forwarded to staking back with the snapshot would make
// stake unbonded at 97 slashable even though the SDK cutoff is 98.
func TestPoASlashPreservesUnbondingCutoff(t *testing.T) {
	app, ctx, auth, _ := newPoASlashApp(t)
	poakeeper.EndBlocker(ctx, app.PoaKeeper, func(ctx sdk.Context) []abci.ValidatorUpdate {
		return staking.EndBlocker(ctx, &app.StakingKeeper)
	})
	delegations := app.StakingKeeper.GetValidatorDelegations(ctx, auth.GetOperator())
	require.Len(t, delegations, 1)
	delegator := delegations[0].GetDelegatorAddr()
	_, err := app.StakingKeeper.Undelegate(ctx, delegator, auth.GetOperator(), sdk.NewDec(100000000))
	require.NoError(t, err)
	ctx = ctx.WithBlockHeight(98)
	poakeeper.EndBlocker(ctx, app.PoaKeeper, func(ctx sdk.Context) []abci.ValidatorUpdate {
		return staking.EndBlocker(ctx, &app.StakingKeeper)
	})
	ctx = ctx.WithBlockHeight(100)
	cons, err := auth.GetConsAddr()
	require.NoError(t, err)
	supplyBefore := app.BankKeeper.GetSupply(ctx, sdk.DefaultBondDenom).Amount
	app.EvidenceKeeper.HandleEquivocationEvidence(ctx, &evidencetypes.Equivocation{
		Height: 99, Time: ctx.BlockTime(), Power: 8122, ConsensusAddress: cons.String(),
	})
	unbonding, found := app.StakingKeeper.GetUnbondingDelegation(ctx, delegator, auth.GetOperator())
	require.True(t, found)
	require.Len(t, unbonding.Entries, 1)
	require.Equal(t, int64(97), unbonding.Entries[0].CreationHeight)
	require.Equal(t, "100000000", unbonding.Entries[0].Balance.String())
	after, found := app.StakingKeeper.GetValidator(ctx, auth.GetOperator())
	require.True(t, found)
	require.Equal(t, "850000000", after.Tokens.String())
	require.Equal(t, "50000000", supplyBefore.Sub(app.BankKeeper.GetSupply(ctx, sdk.DefaultBondDenom).Amount).String())
}
