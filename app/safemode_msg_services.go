package app

import (
	"context"
	"fmt"

	errorsmod "cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/cosmos/cosmos-sdk/types/module"
	"github.com/cosmos/cosmos-sdk/x/authz"
	authzkeeper "github.com/cosmos/cosmos-sdk/x/authz/keeper"
	authzmodule "github.com/cosmos/cosmos-sdk/x/authz/module"
	distr "github.com/cosmos/cosmos-sdk/x/distribution"
	distrexported "github.com/cosmos/cosmos-sdk/x/distribution/exported"
	distrkeeper "github.com/cosmos/cosmos-sdk/x/distribution/keeper"
	distrtypes "github.com/cosmos/cosmos-sdk/x/distribution/types"
	staking "github.com/cosmos/cosmos-sdk/x/staking"
	stakingexported "github.com/cosmos/cosmos-sdk/x/staking/exported"
	stakingkeeper "github.com/cosmos/cosmos-sdk/x/staking/keeper"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
)

func rejectGovernanceActionInSafeMode(goCtx context.Context, poa PoaSafeModeReader) error {
	if poa != nil && poa.SafeModeActive(sdk.UnwrapSDKContext(goCtx)) {
		return errorsmod.Wrap(
			sdkerrors.ErrInvalidRequest,
			"x/poa safe mode active: governance-authority action is frozen until the authority set is restored",
		)
	}
	return nil
}

// safeModeAuthzMsgServer enforces the grant freeze at the message boundary.
// Governance and ICA both call the message router without passing through ante.
type safeModeAuthzMsgServer struct {
	authz.MsgServer
	poa PoaSafeModeReader
}

func (s safeModeAuthzMsgServer) Grant(ctx context.Context, msg *authz.MsgGrant) (*authz.MsgGrantResponse, error) {
	if err := rejectGovernanceActionInSafeMode(ctx, s.poa); err != nil {
		return nil, err
	}
	return s.MsgServer.Grant(ctx, msg)
}

// safeModeAuthzModule retains the SDK module's genesis, query, and migration
// behavior while registering the protected message server.
type safeModeAuthzModule struct {
	authzmodule.AppModule
	keeper authzkeeper.Keeper
	poa    PoaSafeModeReader
}

func (m safeModeAuthzModule) RegisterServices(cfg module.Configurator) {
	authz.RegisterQueryServer(cfg.QueryServer(), m.keeper)
	authz.RegisterMsgServer(cfg.MsgServer(), safeModeAuthzMsgServer{MsgServer: m.keeper, poa: m.poa})
	if err := cfg.RegisterMigration(authz.ModuleName, 1, authzkeeper.NewMigrator(m.keeper).Migrate1to2); err != nil {
		panic(fmt.Sprintf("failed to migrate x/%s from version 1 to 2: %v", authz.ModuleName, err))
	}
}

type safeModeDistributionMsgServer struct {
	distrtypes.MsgServer
	poa PoaSafeModeReader
}

func (s safeModeDistributionMsgServer) UpdateParams(ctx context.Context, msg *distrtypes.MsgUpdateParams) (*distrtypes.MsgUpdateParamsResponse, error) {
	if err := rejectGovernanceActionInSafeMode(ctx, s.poa); err != nil {
		return nil, err
	}
	return s.MsgServer.UpdateParams(ctx, msg)
}

func (s safeModeDistributionMsgServer) CommunityPoolSpend(ctx context.Context, msg *distrtypes.MsgCommunityPoolSpend) (*distrtypes.MsgCommunityPoolSpendResponse, error) {
	if err := rejectGovernanceActionInSafeMode(ctx, s.poa); err != nil {
		return nil, err
	}
	return s.MsgServer.CommunityPoolSpend(ctx, msg)
}

type safeModeDistributionModule struct {
	distr.AppModule
	keeper         distrkeeper.Keeper
	legacySubspace distrexported.Subspace
	poa            PoaSafeModeReader
}

func (m safeModeDistributionModule) RegisterServices(cfg module.Configurator) {
	distrtypes.RegisterMsgServer(cfg.MsgServer(), safeModeDistributionMsgServer{
		MsgServer: distrkeeper.NewMsgServerImpl(m.keeper),
		poa:       m.poa,
	})
	distrtypes.RegisterQueryServer(cfg.QueryServer(), distrkeeper.NewQuerier(m.keeper))

	migrator := distrkeeper.NewMigrator(m.keeper, m.legacySubspace)
	if err := cfg.RegisterMigration(distrtypes.ModuleName, 1, migrator.Migrate1to2); err != nil {
		panic(fmt.Sprintf("failed to migrate x/%s from version 1 to 2: %v", distrtypes.ModuleName, err))
	}
	if err := cfg.RegisterMigration(distrtypes.ModuleName, 2, migrator.Migrate2to3); err != nil {
		panic(fmt.Sprintf("failed to migrate x/%s from version 2 to 3: %v", distrtypes.ModuleName, err))
	}
}

type safeModeStakingMsgServer struct {
	stakingtypes.MsgServer
	poa PoaSafeModeReader
}

func (s safeModeStakingMsgServer) UpdateParams(ctx context.Context, msg *stakingtypes.MsgUpdateParams) (*stakingtypes.MsgUpdateParamsResponse, error) {
	if err := rejectGovernanceActionInSafeMode(ctx, s.poa); err != nil {
		return nil, err
	}
	return s.MsgServer.UpdateParams(ctx, msg)
}

type safeModeStakingModule struct {
	staking.AppModule
	keeper         *stakingkeeper.Keeper
	legacySubspace stakingexported.Subspace
	poa            PoaSafeModeReader
}

func (m safeModeStakingModule) RegisterServices(cfg module.Configurator) {
	stakingtypes.RegisterMsgServer(cfg.MsgServer(), safeModeStakingMsgServer{
		MsgServer: stakingkeeper.NewMsgServerImpl(m.keeper),
		poa:       m.poa,
	})
	stakingtypes.RegisterQueryServer(cfg.QueryServer(), stakingkeeper.Querier{Keeper: m.keeper})

	migrator := stakingkeeper.NewMigrator(m.keeper, m.legacySubspace)
	if err := cfg.RegisterMigration(stakingtypes.ModuleName, 1, migrator.Migrate1to2); err != nil {
		panic(fmt.Sprintf("failed to migrate x/%s from version 1 to 2: %v", stakingtypes.ModuleName, err))
	}
	if err := cfg.RegisterMigration(stakingtypes.ModuleName, 2, migrator.Migrate2to3); err != nil {
		panic(fmt.Sprintf("failed to migrate x/%s from version 2 to 3: %v", stakingtypes.ModuleName, err))
	}
	if err := cfg.RegisterMigration(stakingtypes.ModuleName, 3, migrator.Migrate3to4); err != nil {
		panic(fmt.Sprintf("failed to migrate x/%s from version 3 to 4: %v", stakingtypes.ModuleName, err))
	}
}
