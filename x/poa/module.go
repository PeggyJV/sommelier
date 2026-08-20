package poa

import (
	"context"
	"encoding/json"

	abci "github.com/cometbft/cometbft/abci/types"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/codec"
	cdctypes "github.com/cosmos/cosmos-sdk/codec/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/module"
	sim "github.com/cosmos/cosmos-sdk/types/simulation"
	"github.com/grpc-ecosystem/grpc-gateway/runtime"
	"github.com/spf13/cobra"

	"github.com/peggyjv/sommelier/v10/x/poa/client/cli"
	"github.com/peggyjv/sommelier/v10/x/poa/keeper"
	"github.com/peggyjv/sommelier/v10/x/poa/types"
)

var (
	_ module.AppModule      = AppModule{}
	_ module.AppModuleBasic = AppModuleBasic{}
)

// AppModuleBasic implements module.AppModuleBasic for x/poa.
type AppModuleBasic struct{}

func (AppModuleBasic) Name() string { return types.ModuleName }

func (AppModuleBasic) RegisterLegacyAminoCodec(cdc *codec.LegacyAmino) {
	types.RegisterLegacyAminoCodec(cdc)
}

func (AppModuleBasic) RegisterInterfaces(registry cdctypes.InterfaceRegistry) {
	types.RegisterInterfaces(registry)
}

func (AppModuleBasic) DefaultGenesis(cdc codec.JSONCodec) json.RawMessage {
	return cdc.MustMarshalJSON(types.DefaultGenesis())
}

func (AppModuleBasic) ValidateGenesis(cdc codec.JSONCodec, _ client.TxEncodingConfig, bz json.RawMessage) error {
	var gs types.GenesisState
	if err := cdc.UnmarshalJSON(bz, &gs); err != nil {
		return err
	}
	return gs.Validate()
}

func (AppModuleBasic) RegisterGRPCGatewayRoutes(clientCtx client.Context, mux *runtime.ServeMux) {
	if err := types.RegisterQueryHandlerClient(context.Background(), mux, types.NewQueryClient(clientCtx)); err != nil {
		panic(err)
	}
}

// GetTxCmd returns nil: both PoA messages are gov-only, so there is no
// user-submittable tx. They are proposed through `tx gov submit-proposal`.
func (AppModuleBasic) GetTxCmd() *cobra.Command { return nil }

func (AppModuleBasic) GetQueryCmd() *cobra.Command { return cli.GetQueryCmd() }

// AppModule wires the PoA keeper into the SDK module manager.
type AppModule struct {
	AppModuleBasic
	keeper            keeper.Keeper
	cdc               codec.Codec
	stakingEndBlocker keeper.StakingEndBlockerFn
}

// NewAppModule constructs the PoA AppModule. `stakingEndBlocker` MUST be a
// closure over the production *stakingkeeper.Keeper so PoA can drive
// staking's EndBlocker exactly once per block; see x/poa/keeper.EndBlocker.
func NewAppModule(cdc codec.Codec, k keeper.Keeper, stakingEndBlocker keeper.StakingEndBlockerFn) AppModule {
	return AppModule{
		AppModuleBasic:    AppModuleBasic{},
		keeper:            k,
		cdc:               cdc,
		stakingEndBlocker: stakingEndBlocker,
	}
}

func (AppModule) Name() string                               { return types.ModuleName }
func (AppModule) RegisterInvariants(_ sdk.InvariantRegistry) {}
func (AppModule) QuerierRoute() string                       { return types.QuerierRoute }
func (AppModule) ConsensusVersion() uint64                   { return 1 }

// RegisterServices wires the msg/query servers.
func (am AppModule) RegisterServices(cfg module.Configurator) {
	types.RegisterMsgServer(cfg.MsgServer(), keeper.NewMsgServerImpl(am.keeper))
	types.RegisterQueryServer(cfg.QueryServer(), keeper.NewQuerier(am.keeper))
}

func (am AppModule) InitGenesis(ctx sdk.Context, cdc codec.JSONCodec, data json.RawMessage) []abci.ValidatorUpdate {
	var gs types.GenesisState
	cdc.MustUnmarshalJSON(data, &gs)
	keeper.InitGenesis(ctx, am.keeper, gs)
	return nil
}

func (am AppModule) ExportGenesis(ctx sdk.Context, cdc codec.JSONCodec) json.RawMessage {
	gs := keeper.ExportGenesis(ctx, am.keeper)
	return cdc.MustMarshalJSON(&gs)
}

func (am AppModule) BeginBlock(_ sdk.Context, _ abci.RequestBeginBlock) {}

// EndBlock runs the PoA EndBlocker, which itself drives staking's EndBlocker
// (see keeper.EndBlocker for ordering rationale).
func (am AppModule) EndBlock(ctx sdk.Context, _ abci.RequestEndBlock) []abci.ValidatorUpdate {
	return keeper.EndBlocker(ctx, am.keeper, am.stakingEndBlocker)
}

// AppModuleSimulation: PoA does not currently participate in the simulator.

func (AppModule) GenerateGenesisState(_ *module.SimulationState)  {}
func (AppModule) RegisterStoreDecoder(_ sdk.StoreDecoderRegistry) {}
func (AppModule) WeightedOperations(_ module.SimulationState) []sim.WeightedOperation {
	return nil
}
