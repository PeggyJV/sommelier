package app

import (
	"testing"
	"time"

	dbm "github.com/cometbft/cometbft-db"
	abci "github.com/cometbft/cometbft/abci/types"
	tmjson "github.com/cometbft/cometbft/libs/json"
	"github.com/cometbft/cometbft/libs/log"
	tmproto "github.com/cometbft/cometbft/proto/tendermint/types"
	tmtypes "github.com/cometbft/cometbft/types"
	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	"github.com/cosmos/cosmos-sdk/crypto/keys/ed25519"
	"github.com/cosmos/cosmos-sdk/crypto/keys/secp256k1"
	"github.com/cosmos/cosmos-sdk/testutil/mock"
	sdk "github.com/cosmos/cosmos-sdk/types"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	"github.com/cosmos/cosmos-sdk/x/authz"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	govmodule "github.com/cosmos/cosmos-sdk/x/gov"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"
	govtypesv1 "github.com/cosmos/cosmos-sdk/x/gov/types/v1"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	"github.com/cosmos/gogoproto/proto"
	icahost "github.com/cosmos/ibc-go/v7/modules/apps/27-interchain-accounts/host"
	icahosttypes "github.com/cosmos/ibc-go/v7/modules/apps/27-interchain-accounts/host/types"
	icatypes "github.com/cosmos/ibc-go/v7/modules/apps/27-interchain-accounts/types"
	clienttypes "github.com/cosmos/ibc-go/v7/modules/core/02-client/types"
	channeltypes "github.com/cosmos/ibc-go/v7/modules/core/04-channel/types"
	poakeeper "github.com/peggyjv/sommelier/v10/x/poa/keeper"
	poatypes "github.com/peggyjv/sommelier/v10/x/poa/types"
	"github.com/stretchr/testify/require"
)

const (
	pocControllerPort    = "icacontroller-attacker"
	pocControllerChannel = "channel-77"
	pocHostConnection    = "connection-75"
	pocHostChannel       = "channel-999"
)

// Adapted from the reported ICA PoC. ICA Host calls the message router without
// running ante, so the grant must be rejected when governance executes it.
func TestICASafeModeRejectsGovernanceGrant(t *testing.T) {
	app, ctx, icaAddr, attackerVal, _ := setupICASafeModeRegression(t)
	govAddr := authtypes.NewModuleAddress(govtypes.ModuleName)
	grant, err := authz.NewMsgGrant(govAddr, sdk.AccAddress(attackerVal),
		authz.NewGenericAuthorization(sdk.MsgTypeURL(&poatypes.MsgUpdateAuthoritySet{})), nil)
	require.NoError(t, err)
	proposal, err := govtypesv1.NewMsgSubmitProposal([]sdk.Msg{grant},
		sdk.NewCoins(app.GovKeeper.GetParams(ctx).MinDeposit...), icaAddr.String(), "",
		"safe-mode persistence", "install standing governance authority")
	require.NoError(t, err)

	// Direct delivery and ICA delivery must enforce the same freeze.
	handler := app.MsgServiceRouter().Handler(grant)
	require.NotNil(t, handler)
	_, err = handler(ctx, grant)
	require.ErrorContains(t, err, "safe mode active")
	require.True(t, sendICASafeModeRegressionPacket(t, app, ctx, proposal, 1))
	stored, found := app.GovKeeper.GetProposal(ctx, govtypesv1.DefaultStartingProposalID)
	require.True(t, found)
	vote := govtypesv1.NewMsgVote(icaAddr, stored.Id, govtypesv1.OptionYes, "")
	require.True(t, sendICASafeModeRegressionPacket(t, app, ctx, vote, 2))
	require.NotNil(t, stored.VotingEndTime)
	endCtx := ctx.WithBlockHeight(2).WithBlockTime(stored.VotingEndTime.Add(time.Second))
	govmodule.EndBlocker(endCtx, &app.GovKeeper)
	stored, found = app.GovKeeper.GetProposal(endCtx, stored.Id)
	require.True(t, found)
	require.Equal(t, govtypesv1.StatusFailed, stored.Status)
	installed, _ := app.AuthzKeeper.GetAuthorization(endCtx, sdk.AccAddress(attackerVal), govAddr, sdk.MsgTypeURL(&poatypes.MsgUpdateAuthoritySet{}))
	require.Nil(t, installed)

	// Once safe mode ends, the ICA may again submit an ordinary proposal.
	app.PoaKeeper.SetSafeMode(endCtx, false)
	require.True(t, sendICASafeModeRegressionPacket(t, app, endCtx, proposal, 3))
	_, found = app.GovKeeper.GetProposal(endCtx, stored.Id+1)
	require.True(t, found)
	_, err = handler(endCtx, grant)
	require.NoError(t, err)
	installed, _ = app.AuthzKeeper.GetAuthorization(endCtx, sdk.AccAddress(attackerVal), govAddr, sdk.MsgTypeURL(&poatypes.MsgUpdateAuthoritySet{}))
	require.NotNil(t, installed)
}

func setupICASafeModeRegression(t *testing.T) (*SommelierApp, sdk.Context, sdk.AccAddress, sdk.ValAddress, sdk.ValAddress) {
	t.Helper()

	privVal := mock.PV{PrivKey: ed25519.GenPrivKeyFromSecret([]byte("somm-ica-safe-mode-validator"))}
	pubKey, err := privVal.GetPubKey()
	require.NoError(t, err)
	validator := tmtypes.NewValidator(pubKey, 1)
	valSet := tmtypes.NewValidatorSet([]*tmtypes.Validator{validator})

	icaKey := secp256k1.GenPrivKeyFromSecret([]byte("somm-ica-safe-mode-account"))
	icaAccount := authtypes.NewBaseAccount(icaKey.PubKey().Address().Bytes(), icaKey.PubKey(), 0, 0)
	icaBalance := banktypes.Balance{
		Address: icaAccount.GetAddress().String(),
		Coins: sdk.NewCoins(
			sdk.NewCoin(sdk.DefaultBondDenom, sdk.NewInt(1_000_000_000_000_000)),
		),
	}

	app := NewSommelierApp(log.NewNopLogger(), dbm.NewMemDB(), nil, true, map[int64]bool{}, DefaultNodeHome, 0, MakeEncodingConfig(), EmptyAppOptions{})
	genesis := genesisStateWithValSet(t, app, NewDefaultGenesisState(), valSet, []authtypes.GenesisAccount{icaAccount}, icaBalance)
	stateBytes, err := tmjson.Marshal(genesis)
	require.NoError(t, err)

	blockTime := time.Date(2026, time.August, 28, 0, 0, 0, 0, time.UTC)
	app.InitChain(abci.RequestInitChain{
		Time:            blockTime,
		Validators:      []abci.ValidatorUpdate{},
		ConsensusParams: DefaultConsensusParams,
		AppStateBytes:   stateBytes,
	})
	ctx := app.NewContext(false, tmproto.Header{Height: 1, Time: blockTime})

	// Start the fixture in community-only safe mode, with a jailed authority.
	attackerVal := sdk.ValAddress(validator.Address)
	authorityKey := ed25519.GenPrivKeyFromSecret([]byte("somm-ica-safe-mode-returning-authority"))
	authorityPubkey, err := codectypes.NewAnyWithValue(authorityKey.PubKey())
	require.NoError(t, err)
	offlineAuthority := sdk.ValAddress(authorityKey.PubKey().Address())
	communityVal, found := app.StakingKeeper.GetValidator(ctx, attackerVal)
	require.True(t, found)
	offlineVal := communityVal
	offlineVal.OperatorAddress = offlineAuthority.String()
	offlineVal.ConsensusPubkey = authorityPubkey
	offlineVal.Jailed = true
	app.StakingKeeper.SetValidator(ctx, offlineVal)
	require.NoError(t, app.StakingKeeper.SetValidatorByConsAddr(ctx, offlineVal))

	app.PoaKeeper.SetParams(ctx, poatypes.DefaultParams())
	app.PoaKeeper.SetAuthoritySet(ctx, []sdk.ValAddress{offlineAuthority})
	poakeeper.EndBlocker(ctx, app.PoaKeeper, func(sdk.Context) []abci.ValidatorUpdate { return nil })
	require.True(t, app.PoaKeeper.SafeModeActive(ctx))

	// Enable ICA Host with the reported v10 unrestricted message allowlist.
	app.ICAHostKeeper.SetParams(ctx, icahosttypes.NewParams(true, []string{"*"}))
	app.ICAHostKeeper.SetInterchainAccountAddress(ctx, pocHostConnection, pocControllerPort, icaAccount.GetAddress().String())

	metadata := icatypes.NewDefaultMetadata("connection-0", pocHostConnection)
	metadata.Address = icaAccount.GetAddress().String()
	versionBytes, err := icatypes.ModuleCdc.MarshalJSON(&metadata)
	require.NoError(t, err)
	channel := channeltypes.NewChannel(
		channeltypes.OPEN,
		channeltypes.ORDERED,
		channeltypes.NewCounterparty(pocControllerPort, pocControllerChannel),
		[]string{pocHostConnection},
		string(versionBytes),
	)
	app.IBCKeeper.ChannelKeeper.SetChannel(ctx, icatypes.HostPortID, pocHostChannel, channel)

	// Ensure the proposed replacement is a real bonded, unjailed validator, as
	// required by x/poa's own validation.
	storedVal, found := app.StakingKeeper.GetValidator(ctx, attackerVal)
	require.True(t, found)
	require.Equal(t, stakingtypes.Bonded, storedVal.Status)
	require.False(t, storedVal.Jailed)

	return app, ctx, icaAccount.GetAddress(), attackerVal, offlineAuthority
}

func sendICASafeModeRegressionPacket(t *testing.T, app *SommelierApp, ctx sdk.Context, msg sdk.Msg, sequence uint64) bool {
	t.Helper()

	data, err := icatypes.SerializeCosmosTx(app.AppCodec(), []proto.Message{msg})
	require.NoError(t, err)
	packetData := icatypes.InterchainAccountPacketData{Type: icatypes.EXECUTE_TX, Data: data}
	packet := channeltypes.NewPacket(
		packetData.GetBytes(),
		sequence,
		pocControllerPort,
		pocControllerChannel,
		icatypes.HostPortID,
		pocHostChannel,
		clienttypes.NewHeight(0, 1000),
		0,
	)

	ack := icahost.NewIBCModule(app.ICAHostKeeper).OnRecvPacket(ctx, packet, nil)
	return ack.Success()
}
