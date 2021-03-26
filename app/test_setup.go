package app

import (
	"encoding/json"

	abci "github.com/tendermint/tendermint/abci/types"
	"github.com/tendermint/tendermint/libs/log"
	dbm "github.com/tendermint/tm-db"

	"github.com/cosmos/cosmos-sdk/baseapp"
	servertypes "github.com/cosmos/cosmos-sdk/server/types"
	"github.com/cosmos/cosmos-sdk/simapp"
	storetypes "github.com/cosmos/cosmos-sdk/store/types"
	"github.com/cosmos/cosmos-sdk/testutil/network"
)

// Setup initializes a new SimApp. A Nop logger is set in SimApp.
func Setup(isCheckTx bool) *SommelierApp {
	app, genesisState := setup(!isCheckTx, 5)
	if !isCheckTx {
		// init chain must be called to stop deliverState from being nil
		stateBytes, err := json.MarshalIndent(genesisState, "", " ")
		if err != nil {
			panic(err)
		}

		// Initialize the chain
		app.InitChain(
			abci.RequestInitChain{
				Validators:      []abci.ValidatorUpdate{},
				ConsensusParams: simapp.DefaultConsensusParams,
				AppStateBytes:   stateBytes,
			},
		)
	}

	return app
}

func setup(withGenesis bool, invCheckPeriod uint) (*SommelierApp, GenesisState) {
	db := dbm.NewMemDB()
	encCdc := MakeEncodingConfig()
	app := NewSommelierApp(log.NewNopLogger(), db, nil, true, map[int64]bool{}, DefaultNodeHome, invCheckPeriod, encCdc, simapp.EmptyAppOptions{})
	if withGenesis {
		return app, NewDefaultGenesisState()
	}
	return app, GenesisState{}
}

func NewConfig() network.Config {
	cfg := network.DefaultConfig()
	encCfg := MakeEncodingConfig()
	cfg.Codec = encCfg.Marshaler
	cfg.TxConfig = encCfg.TxConfig
	cfg.LegacyAmino = encCfg.Amino
	cfg.InterfaceRegistry = encCfg.InterfaceRegistry
	cfg.AppConstructor = AppConstructor
	cfg.GenesisState = NewDefaultGenesisState()
	return cfg
}

func AppConstructor(val network.Validator) servertypes.Application {
	return NewSommelierApp(
		val.Ctx.Logger, dbm.NewMemDB(), nil, true, make(map[int64]bool),
		val.Ctx.Config.RootDir, 0, MakeEncodingConfig(), simapp.EmptyAppOptions{},
		baseapp.SetPruning(storetypes.NewPruningOptionsFromString(val.AppConfig.Pruning)),
		baseapp.SetMinGasPrices(val.AppConfig.MinGasPrices),
	)
}
