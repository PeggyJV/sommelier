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
)

// safeModeAuthzMsgServer enforces the grant freeze at the message boundary.
// ICA and governance both call the message router without passing through ante.
type safeModeAuthzMsgServer struct {
	authz.MsgServer
	poa PoaSafeModeReader
}

func (s safeModeAuthzMsgServer) Grant(ctx context.Context, msg *authz.MsgGrant) (*authz.MsgGrantResponse, error) {
	if s.poa != nil && s.poa.SafeModeActive(sdk.UnwrapSDKContext(ctx)) {
		return nil, errorsmod.Wrap(sdkerrors.ErrInvalidRequest,
			"x/poa safe mode active: authz MsgGrant is frozen until the authority set is restored")
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
