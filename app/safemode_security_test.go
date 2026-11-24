package app

import (
	"bytes"
	"testing"

	dbm "github.com/cometbft/cometbft-db"
	abci "github.com/cometbft/cometbft/abci/types"
	tmjson "github.com/cometbft/cometbft/libs/json"
	"github.com/cometbft/cometbft/libs/log"
	tmproto "github.com/cometbft/cometbft/proto/tendermint/types"
	tmtypes "github.com/cometbft/cometbft/types"
	"github.com/cosmos/cosmos-sdk/crypto/keys/secp256k1"
	"github.com/cosmos/cosmos-sdk/testutil/mock"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkmodule "github.com/cosmos/cosmos-sdk/types/module"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	"github.com/cosmos/cosmos-sdk/x/authz"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	distrtypes "github.com/cosmos/cosmos-sdk/x/distribution/types"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	gogogrpc "github.com/cosmos/gogoproto/grpc"
	gravitytypes "github.com/peggyjv/gravity-bridge/module/v6/x/gravity/types"
	poatypes "github.com/peggyjv/sommelier/v10/x/poa/types"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc"
)

func newSafeModeSecurityApp(t *testing.T) (*SommelierApp, sdk.Context) {
	t.Helper()

	privVal := mock.NewPV()
	pubKey, err := privVal.GetPubKey()
	require.NoError(t, err)
	valSet := tmtypes.NewValidatorSet([]*tmtypes.Validator{tmtypes.NewValidator(pubKey, 1)})
	accountKey := secp256k1.GenPrivKeyFromSecret([]byte("safe-mode-security-account"))
	account := authtypes.NewBaseAccount(
		accountKey.PubKey().Address().Bytes(), accountKey.PubKey(), 0, 0,
	)
	balance := banktypes.Balance{
		Address: account.GetAddress().String(),
		Coins:   sdk.NewCoins(sdk.NewInt64Coin(sdk.DefaultBondDenom, 1)),
	}

	app := NewSommelierApp(
		log.NewNopLogger(), dbm.NewMemDB(), nil, true, map[int64]bool{},
		t.TempDir(), 0, MakeEncodingConfig(), EmptyAppOptions{},
	)
	genesis := genesisStateWithValSet(
		t, app, NewDefaultGenesisState(), valSet,
		[]authtypes.GenesisAccount{account}, balance,
	)
	stateBytes, err := tmjson.Marshal(genesis)
	require.NoError(t, err)
	app.InitChain(abci.RequestInitChain{
		ConsensusParams: DefaultConsensusParams,
		AppStateBytes:   stateBytes,
	})
	ctx := app.NewContext(false, tmproto.Header{Height: 1})
	app.PoaKeeper.SetSafeMode(ctx, true)
	return app, ctx
}

func TestEndBlockOrderMakesPoaPowerVisibleToGravity(t *testing.T) {
	app, _ := newSafeModeSecurityApp(t)

	index := func(moduleName string) int {
		for i, name := range app.mm.OrderEndBlockers {
			if name == moduleName {
				return i
			}
		}
		return -1
	}

	gravityIndex := index(gravitytypes.ModuleName)
	poaIndex := index(poatypes.ModuleName)
	require.NotEqual(t, -1, gravityIndex)
	require.NotEqual(t, -1, poaIndex)
	require.Less(t, poaIndex, gravityIndex,
		"gravity must observe PoA's current multiplier snapshot and validator powers")
}

func TestGravitySafeModeModuleKeepsThawBlockFrozen(t *testing.T) {
	app, thawCtx := newSafeModeSecurityApp(t)
	gravityModule, ok := app.mm.Modules[gravitytypes.ModuleName].(*gravitySafeModeModule)
	require.True(t, ok, "the app module manager must use the safe-mode Gravity wrapper")
	require.NotNil(t, gravityModule)

	// Exercise the exact wrapper instance installed in the app without running
	// unrelated modules. Gravity BeginBlock creates the first signer-set tx, so
	// its absence/presence is a durable observation of whether real processing
	// ran rather than an assertion on the wrapper's private flag.
	app.mm.OrderBeginBlockers = []string{gravitytypes.ModuleName}
	app.mm.OrderEndBlockers = []string{gravitytypes.ModuleName}
	require.Nil(t, app.GravityKeeper.GetLatestSignerSetTx(thawCtx))

	// BeginBlock observes the freeze. PoA clears its persisted flag later in
	// the same block, but gravity EndBlock must still no-op because validators
	// had no unfrozen DeliverTx phase in which to submit confirmations.
	app.mm.BeginBlock(thawCtx, abci.RequestBeginBlock{})
	app.PoaKeeper.SetSafeMode(thawCtx, false)
	require.NotPanics(t, func() {
		app.mm.EndBlock(thawCtx, abci.RequestEndBlock{})
	})
	require.Nil(t, app.GravityKeeper.GetLatestSignerSetTx(thawCtx),
		"Gravity processing must remain frozen through the thaw block")

	// The following block has an unfrozen DeliverTx phase. Its real Gravity
	// BeginBlock must resume and create the first signer-set tx.
	nextCtx := thawCtx.WithBlockHeight(thawCtx.BlockHeight() + 1)
	app.mm.BeginBlock(nextCtx, abci.RequestBeginBlock{})
	require.NotNil(t, app.GravityKeeper.GetLatestSignerSetTx(nextCtx),
		"Gravity processing must resume in the block after thaw")
}

func TestSafeModeMessageRouterRejectsCommunityPoolSpend(t *testing.T) {
	app, ctx := newSafeModeSecurityApp(t)
	govAddr := authtypes.NewModuleAddress(govtypes.ModuleName)
	msg := &distrtypes.MsgCommunityPoolSpend{
		Authority: govAddr.String(),
		Recipient: sdk.AccAddress(bytes.Repeat([]byte{0x41}, 20)).String(),
		Amount:    sdk.NewCoins(),
	}

	handler := app.MsgServiceRouter().Handler(msg)
	require.NotNil(t, handler)
	_, err := handler(ctx, msg)
	require.ErrorContains(t, err, "safe mode active")

	app.PoaKeeper.SetSafeMode(ctx, false)
	_, err = handler(ctx, msg)
	require.NoError(t, err, "community-pool spends must resume after safe mode")
}

func TestSafeModeMessageRouterRejectsDistributionParams(t *testing.T) {
	app, ctx := newSafeModeSecurityApp(t)
	msg := &distrtypes.MsgUpdateParams{
		Authority: authtypes.NewModuleAddress(govtypes.ModuleName).String(),
		Params:    app.DistrKeeper.GetParams(ctx),
	}

	handler := app.MsgServiceRouter().Handler(msg)
	require.NotNil(t, handler)
	_, err := handler(ctx, msg)
	require.ErrorContains(t, err, "safe mode active")

	app.PoaKeeper.SetSafeMode(ctx, false)
	_, err = handler(ctx, msg)
	require.NoError(t, err, "distribution parameter updates must resume after safe mode")
}

func TestSafeModeMessageRouterRejectsStakingParams(t *testing.T) {
	app, ctx := newSafeModeSecurityApp(t)
	msg := &stakingtypes.MsgUpdateParams{
		Authority: authtypes.NewModuleAddress(govtypes.ModuleName).String(),
		Params:    app.StakingKeeper.GetParams(ctx),
	}

	handler := app.MsgServiceRouter().Handler(msg)
	require.NotNil(t, handler)
	_, err := handler(ctx, msg)
	require.ErrorContains(t, err, "safe mode active")

	app.PoaKeeper.SetSafeMode(ctx, false)
	_, err = handler(ctx, msg)
	require.NoError(t, err, "staking parameter updates must resume after safe mode")
}

func TestSafeModeMessageRouterRejectsAuthzGrant(t *testing.T) {
	app, ctx := newSafeModeSecurityApp(t)
	govAddr := authtypes.NewModuleAddress(govtypes.ModuleName)
	grantee := sdk.AccAddress(bytes.Repeat([]byte{0x42}, 20))
	msg, err := authz.NewMsgGrant(
		govAddr,
		grantee,
		authz.NewGenericAuthorization(sdk.MsgTypeURL(&poatypes.MsgUpdateAuthoritySet{})),
		nil,
	)
	require.NoError(t, err)

	handler := app.MsgServiceRouter().Handler(msg)
	require.NotNil(t, handler)
	_, err = handler(ctx, msg)
	require.ErrorContains(t, err, "safe mode active")
	stored, _ := app.AuthzKeeper.GetAuthorization(
		ctx, grantee, govAddr, sdk.MsgTypeURL(&poatypes.MsgUpdateAuthoritySet{}),
	)
	require.Nil(t, stored)
}

type recordingGRPCServer struct {
	services []string
}

func (s *recordingGRPCServer) RegisterService(desc *grpc.ServiceDesc, _ interface{}) {
	s.services = append(s.services, desc.ServiceName)
}

type migrationRegistration struct {
	moduleName  string
	fromVersion uint64
	hasHandler  bool
}

type recordingConfigurator struct {
	msgServer   recordingGRPCServer
	queryServer recordingGRPCServer
	migrations  []migrationRegistration
}

func (c *recordingConfigurator) MsgServer() gogogrpc.Server   { return &c.msgServer }
func (c *recordingConfigurator) QueryServer() gogogrpc.Server { return &c.queryServer }
func (c *recordingConfigurator) RegisterMigration(
	moduleName string,
	fromVersion uint64,
	handler sdkmodule.MigrationHandler,
) error {
	c.migrations = append(c.migrations, migrationRegistration{
		moduleName: moduleName, fromVersion: fromVersion, hasHandler: handler != nil,
	})
	return nil
}

func TestSafeModeModulesPreserveSDKServiceAndMigrationRegistrations(t *testing.T) {
	tests := []struct {
		name           string
		register       func(sdkmodule.Configurator)
		wantMigrations []migrationRegistration
	}{
		{
			name:     authz.ModuleName,
			register: safeModeAuthzModule{}.RegisterServices,
			wantMigrations: []migrationRegistration{
				{moduleName: authz.ModuleName, fromVersion: 1, hasHandler: true},
			},
		},
		{
			name:     distrtypes.ModuleName,
			register: safeModeDistributionModule{}.RegisterServices,
			wantMigrations: []migrationRegistration{
				{moduleName: distrtypes.ModuleName, fromVersion: 1, hasHandler: true},
				{moduleName: distrtypes.ModuleName, fromVersion: 2, hasHandler: true},
			},
		},
		{
			name:     stakingtypes.ModuleName,
			register: safeModeStakingModule{}.RegisterServices,
			wantMigrations: []migrationRegistration{
				{moduleName: stakingtypes.ModuleName, fromVersion: 1, hasHandler: true},
				{moduleName: stakingtypes.ModuleName, fromVersion: 2, hasHandler: true},
				{moduleName: stakingtypes.ModuleName, fromVersion: 3, hasHandler: true},
			},
		},
	}

	for _, tc := range tests {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			cfg := &recordingConfigurator{}
			tc.register(cfg)

			require.Len(t, cfg.msgServer.services, 1, "SDK message service registration must be preserved")
			require.Len(t, cfg.queryServer.services, 1, "SDK query service registration must be preserved")
			require.Equal(t, tc.wantMigrations, cfg.migrations,
				"SDK migration registrations must be preserved exactly")
		})
	}
}
