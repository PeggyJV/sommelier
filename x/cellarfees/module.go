package cellarfees

import (
	"context"
	"encoding/json"
	"fmt"

	abci "github.com/cometbft/cometbft/abci/types"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/codec"
	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/module"
	sim "github.com/cosmos/cosmos-sdk/types/simulation"
	paramtypes "github.com/cosmos/cosmos-sdk/x/params/types"
	"github.com/grpc-ecosystem/grpc-gateway/runtime"
	"github.com/peggyjv/sommelier/v9/x/cellarfees/client/cli"
	"github.com/peggyjv/sommelier/v9/x/cellarfees/keeper"
	"github.com/peggyjv/sommelier/v9/x/cellarfees/types"
	typesv2 "github.com/peggyjv/sommelier/v9/x/cellarfees/types/v2"
	"github.com/spf13/cobra"
)

var (
	_ module.AppModule      = AppModule{}
	_ module.AppModuleBasic = AppModuleBasic{}
)

// AppModuleBasic defines the basic application module used by the cellarfees module.
type AppModuleBasic struct{}

// Name returns the cellarfees module's name
func (AppModuleBasic) Name() string {
	return types.ModuleName
}

// RegisterLegacyAminoCodec doesn't support amino
func (AppModuleBasic) RegisterLegacyAminoCodec(cdc *codec.LegacyAmino) {}

// DefaultGenesis returns default genesis state as raw bytes for the cellarfees
// module.
func (AppModuleBasic) DefaultGenesis(cdc codec.JSONCodec) json.RawMessage {
	gs := typesv2.DefaultGenesisState()
	return cdc.MustMarshalJSON(&gs)
}

// ValidateGenesis performs genesis state validation for the cellarfees module.
func (AppModuleBasic) ValidateGenesis(cdc codec.JSONCodec, _ client.TxEncodingConfig, bz json.RawMessage) error {
	var gs typesv2.GenesisState
	if err := cdc.UnmarshalJSON(bz, &gs); err != nil {
		return err
	}
	return gs.Validate()
}

// GetTxCmd returns the root tx command for the cellarfees module.
func (AppModuleBasic) GetTxCmd() *cobra.Command {
	return cli.GetTxCmd()
}

// GetQueryCmd returns the root query command for the cellarfees module.
func (AppModuleBasic) GetQueryCmd() *cobra.Command {
	return cli.GetQueryCmd(types.StoreKey)
}

// RegisterGRPCGatewayRoutes registers the gRPC Gateway routes for the cellarfees module.
// also implements AppModuleBasic
func (AppModuleBasic) RegisterGRPCGatewayRoutes(clientCtx client.Context, mux *runtime.ServeMux) {
	typesv2.RegisterQueryHandlerClient(context.Background(), mux, typesv2.NewQueryClient(clientCtx))
}

// RegisterInterfaces implements AppModuleBasic
func (b AppModuleBasic) RegisterInterfaces(registry codectypes.InterfaceRegistry) {
	typesv2.RegisterInterfaces(registry)
}

// AppModule implements an application module for the cellarfees module.
type AppModule struct {
	AppModuleBasic
	keeper        keeper.Keeper
	cdc           codec.Codec
	accountKeeper types.AccountKeeper
	bankKeeper    types.BankKeeper
	mintKeeper    types.MintKeeper
	corkKeeper    types.CorkKeeper
	auctionKeeper types.AuctionKeeper

	// legacy subspace used for v1 -> v2 migration
	legacySubspace paramtypes.Subspace
}

// NewAppModule creates a new AppModule object
func NewAppModule(keeper keeper.Keeper, cdc codec.Codec, accountKeeper types.AccountKeeper, bankKeeper types.BankKeeper,
	mintKeeper types.MintKeeper, corkKeeper types.CorkKeeper, auctionKeeper types.AuctionKeeper, legacySubspace paramtypes.Subspace) AppModule {
	return AppModule{
		AppModuleBasic: AppModuleBasic{},
		keeper:         keeper,
		cdc:            cdc,
		accountKeeper:  accountKeeper,
		bankKeeper:     bankKeeper,
		mintKeeper:     mintKeeper,
		corkKeeper:     corkKeeper,
		auctionKeeper:  auctionKeeper,
		legacySubspace: legacySubspace,
	}
}

// Name returns the cellarfees module's name.
func (AppModule) Name() string { return types.ModuleName }

// RegisterInvariants performs a no-op.
func (am AppModule) RegisterInvariants(ir sdk.InvariantRegistry) {}

// QuerierRoute returns the cellarfees module's querier route name.
func (AppModule) QuerierRoute() string { return types.QuerierRoute }

// ConsensusVersion implements AppModule/ConsensusVersion.
func (AppModule) ConsensusVersion() uint64 {
	return 3
}

func (am AppModule) WeightedOperations(simState module.SimulationState) []sim.WeightedOperation {
	return nil
}

// RegisterServices registers module services.
func (am AppModule) RegisterServices(cfg module.Configurator) {
	//types.RegisterMsgServer(cfg.MsgServer(), am.keeper)
	typesv2.RegisterQueryServer(cfg.QueryServer(), am.keeper)

	migrator := keeper.NewMigrator(am.keeper, am.legacySubspace)
	if err := cfg.RegisterMigration(types.ModuleName, 1, migrator.Migrate1to2); err != nil {
		panic(fmt.Sprintf("failed to migrate x/cellarfees from version 1 to 2: %v", err))
	}
	if err := cfg.RegisterMigration(types.ModuleName, 2, migrator.Migrate2to3); err != nil {
		panic(fmt.Sprintf("failed to migrate x/cellarfees from version 2 to 3: %v", err))
	}
}

// InitGenesis performs genesis initialization for the cellarfees module.
func (am AppModule) InitGenesis(ctx sdk.Context, cdc codec.JSONCodec, data json.RawMessage) []abci.ValidatorUpdate {
	var genesisState typesv2.GenesisState
	cdc.MustUnmarshalJSON(data, &genesisState)
	am.keeper.InitGenesis(ctx, genesisState)

	return []abci.ValidatorUpdate{}
}

// ExportGenesis returns the exported genesis state as raw bytes for the cellarfees
// module.
func (am AppModule) ExportGenesis(ctx sdk.Context, cdc codec.JSONCodec) json.RawMessage {
	genesisState := am.keeper.ExportGenesis(ctx)
	return cdc.MustMarshalJSON(&genesisState)
}

// BeginBlock returns the begin blocker for the cellarfees module.
func (am AppModule) BeginBlock(ctx sdk.Context, _ abci.RequestBeginBlock) {
	am.keeper.BeginBlocker(ctx)
}

// EndBlock returns the end blocker for the cellarfees module.
func (am AppModule) EndBlock(ctx sdk.Context, _ abci.RequestEndBlock) []abci.ValidatorUpdate {
	am.keeper.EndBlocker(ctx)
	return []abci.ValidatorUpdate{}
}

// AppModuleSimulation functions

// GenerateGenesisState creates a randomized GenState of the cellarfees module.
func (AppModule) GenerateGenesisState(simState *module.SimulationState) {}

// ProposalContents returns all the cellarfees content functions used to
// simulate governance proposals.
func (am AppModule) ProposalContents(_ module.SimulationState) []sim.WeightedProposalMsg {
	return nil
}

// RegisterStoreDecoder registers a decoder for cellarfees module's types
func (am AppModule) RegisterStoreDecoder(sdr sdk.StoreDecoderRegistry) {}
