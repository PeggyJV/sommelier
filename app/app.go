package app

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"

	"github.com/cosmos/cosmos-sdk/runtime"

	autocliv1 "cosmossdk.io/api/cosmos/autocli/v1"
	reflectionv1 "cosmossdk.io/api/cosmos/reflection/v1"
	dbm "github.com/cometbft/cometbft-db"
	abci "github.com/cometbft/cometbft/abci/types"
	tmjson "github.com/cometbft/cometbft/libs/json"
	"github.com/cometbft/cometbft/libs/log"
	tmos "github.com/cometbft/cometbft/libs/os"
	"github.com/cosmos/cosmos-sdk/baseapp"
	"github.com/cosmos/cosmos-sdk/client"
	nodeservice "github.com/cosmos/cosmos-sdk/client/grpc/node"
	"github.com/cosmos/cosmos-sdk/client/grpc/tmservice"
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/codec/types"
	runtimeservices "github.com/cosmos/cosmos-sdk/runtime/services"
	"github.com/cosmos/cosmos-sdk/server/api"
	"github.com/cosmos/cosmos-sdk/server/config"
	servertypes "github.com/cosmos/cosmos-sdk/server/types"
	storetypes "github.com/cosmos/cosmos-sdk/store/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/module"
	"github.com/cosmos/cosmos-sdk/version"
	"github.com/cosmos/cosmos-sdk/x/auth"
	"github.com/cosmos/cosmos-sdk/x/auth/ante"
	authkeeper "github.com/cosmos/cosmos-sdk/x/auth/keeper"
	authtx "github.com/cosmos/cosmos-sdk/x/auth/tx"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	"github.com/cosmos/cosmos-sdk/x/auth/vesting"
	vestingtypes "github.com/cosmos/cosmos-sdk/x/auth/vesting/types"
	"github.com/cosmos/cosmos-sdk/x/authz"
	authzkeeper "github.com/cosmos/cosmos-sdk/x/authz/keeper"
	authzmodule "github.com/cosmos/cosmos-sdk/x/authz/module"
	"github.com/cosmos/cosmos-sdk/x/bank"
	bankkeeper "github.com/cosmos/cosmos-sdk/x/bank/keeper"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	"github.com/cosmos/cosmos-sdk/x/capability"
	capabilitykeeper "github.com/cosmos/cosmos-sdk/x/capability/keeper"
	capabilitytypes "github.com/cosmos/cosmos-sdk/x/capability/types"
	consensusparamkeeper "github.com/cosmos/cosmos-sdk/x/consensus/keeper"
	consensusparamtypes "github.com/cosmos/cosmos-sdk/x/consensus/types"
	"github.com/cosmos/cosmos-sdk/x/crisis"
	crisiskeeper "github.com/cosmos/cosmos-sdk/x/crisis/keeper"
	crisistypes "github.com/cosmos/cosmos-sdk/x/crisis/types"
	distr "github.com/cosmos/cosmos-sdk/x/distribution"
	distrkeeper "github.com/cosmos/cosmos-sdk/x/distribution/keeper"
	distrtypes "github.com/cosmos/cosmos-sdk/x/distribution/types"
	"github.com/cosmos/cosmos-sdk/x/evidence"
	evidencekeeper "github.com/cosmos/cosmos-sdk/x/evidence/keeper"
	evidencetypes "github.com/cosmos/cosmos-sdk/x/evidence/types"
	"github.com/cosmos/cosmos-sdk/x/feegrant"
	feegrantkeeper "github.com/cosmos/cosmos-sdk/x/feegrant/keeper"
	feegrantmodule "github.com/cosmos/cosmos-sdk/x/feegrant/module"
	"github.com/cosmos/cosmos-sdk/x/genutil"
	genutiltypes "github.com/cosmos/cosmos-sdk/x/genutil/types"
	"github.com/cosmos/cosmos-sdk/x/gov"
	govclient "github.com/cosmos/cosmos-sdk/x/gov/client"
	govkeeper "github.com/cosmos/cosmos-sdk/x/gov/keeper"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"
	govtypesv1beta1 "github.com/cosmos/cosmos-sdk/x/gov/types/v1beta1"
	"github.com/cosmos/cosmos-sdk/x/mint"
	mintkeeper "github.com/cosmos/cosmos-sdk/x/mint/keeper"
	minttypes "github.com/cosmos/cosmos-sdk/x/mint/types"
	"github.com/cosmos/cosmos-sdk/x/params"
	paramsclient "github.com/cosmos/cosmos-sdk/x/params/client"
	paramskeeper "github.com/cosmos/cosmos-sdk/x/params/keeper"
	paramstypes "github.com/cosmos/cosmos-sdk/x/params/types"
	paramproposal "github.com/cosmos/cosmos-sdk/x/params/types/proposal"
	"github.com/cosmos/cosmos-sdk/x/slashing"
	slashingkeeper "github.com/cosmos/cosmos-sdk/x/slashing/keeper"
	slashingtypes "github.com/cosmos/cosmos-sdk/x/slashing/types"
	"github.com/cosmos/cosmos-sdk/x/staking"
	stakingkeeper "github.com/cosmos/cosmos-sdk/x/staking/keeper"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	"github.com/cosmos/cosmos-sdk/x/upgrade"
	upgradeclient "github.com/cosmos/cosmos-sdk/x/upgrade/client"
	upgradekeeper "github.com/cosmos/cosmos-sdk/x/upgrade/keeper"
	upgradetypes "github.com/cosmos/cosmos-sdk/x/upgrade/types"
	ica "github.com/cosmos/ibc-go/v7/modules/apps/27-interchain-accounts"
	icahost "github.com/cosmos/ibc-go/v7/modules/apps/27-interchain-accounts/host"
	icahostkeeper "github.com/cosmos/ibc-go/v7/modules/apps/27-interchain-accounts/host/keeper"
	icahosttypes "github.com/cosmos/ibc-go/v7/modules/apps/27-interchain-accounts/host/types"
	icatypes "github.com/cosmos/ibc-go/v7/modules/apps/27-interchain-accounts/types"
	ibctransfer "github.com/cosmos/ibc-go/v7/modules/apps/transfer"
	ibctransferkeeper "github.com/cosmos/ibc-go/v7/modules/apps/transfer/keeper"
	ibctransfertypes "github.com/cosmos/ibc-go/v7/modules/apps/transfer/types"
	ibc "github.com/cosmos/ibc-go/v7/modules/core"
	ibcclient "github.com/cosmos/ibc-go/v7/modules/core/02-client"
	ibcclientclient "github.com/cosmos/ibc-go/v7/modules/core/02-client/client"
	ibcclienttypes "github.com/cosmos/ibc-go/v7/modules/core/02-client/types"
	ibcporttypes "github.com/cosmos/ibc-go/v7/modules/core/05-port/types"
	icaexported "github.com/cosmos/ibc-go/v7/modules/core/exported"
	ibckeeper "github.com/cosmos/ibc-go/v7/modules/core/keeper"
	"github.com/gorilla/mux"
	"github.com/peggyjv/gravity-bridge/module/v4/x/gravity"
	gravityclient "github.com/peggyjv/gravity-bridge/module/v4/x/gravity/client"
	gravitykeeper "github.com/peggyjv/gravity-bridge/module/v4/x/gravity/keeper"
	gravitytypes "github.com/peggyjv/gravity-bridge/module/v4/x/gravity/types"
	appParams "github.com/peggyjv/sommelier/v7/app/params"
	"github.com/peggyjv/sommelier/v7/x/addresses"
	addresseskeeper "github.com/peggyjv/sommelier/v7/x/addresses/keeper"
	addressestypes "github.com/peggyjv/sommelier/v7/x/addresses/types"
	"github.com/peggyjv/sommelier/v7/x/auction"
	auctionclient "github.com/peggyjv/sommelier/v7/x/auction/client"
	auctionkeeper "github.com/peggyjv/sommelier/v7/x/auction/keeper"
	auctiontypes "github.com/peggyjv/sommelier/v7/x/auction/types"
	"github.com/peggyjv/sommelier/v7/x/axelarcork"
	axelarcorkclient "github.com/peggyjv/sommelier/v7/x/axelarcork/client"
	axelarcorkkeeper "github.com/peggyjv/sommelier/v7/x/axelarcork/keeper"
	axelarcorktypes "github.com/peggyjv/sommelier/v7/x/axelarcork/types"
	"github.com/peggyjv/sommelier/v7/x/cellarfees"
	cellarfeeskeeper "github.com/peggyjv/sommelier/v7/x/cellarfees/keeper"
	cellarfeestypes "github.com/peggyjv/sommelier/v7/x/cellarfees/types"
	"github.com/peggyjv/sommelier/v7/x/cork"
	corkclient "github.com/peggyjv/sommelier/v7/x/cork/client"
	corkkeeper "github.com/peggyjv/sommelier/v7/x/cork/keeper"
	corktypes "github.com/peggyjv/sommelier/v7/x/cork/types"
	"github.com/peggyjv/sommelier/v7/x/incentives"
	incentiveskeeper "github.com/peggyjv/sommelier/v7/x/incentives/keeper"
	incentivestypes "github.com/peggyjv/sommelier/v7/x/incentives/types"
	"github.com/peggyjv/sommelier/v7/x/pubsub"
	pubsubclient "github.com/peggyjv/sommelier/v7/x/pubsub/client"
	pubsubkeeper "github.com/peggyjv/sommelier/v7/x/pubsub/keeper"
	pubsubtypes "github.com/peggyjv/sommelier/v7/x/pubsub/types"
	"github.com/rakyll/statik/fs"

	// unnamed import of statik for swagger UI support
	_ "github.com/cosmos/cosmos-sdk/client/docs/statik"
)

const appName = "SommelierApp"

// These variables are commentted out to support future upgrades of the gravity module.

// const upgradeName = "CabernetFranc"
// const newGravityContractAddress = "0x0000000000000000000000000000000000000000"
// const newGravityContractDeployHeight = 1000

var (
	// DefaultNodeHome default home directories for the application daemon
	DefaultNodeHome string

	// ModuleBasics defines the module BasicManager is in charge of setting up basic,
	// non-dependant module elements, such as codec registration
	// and genesis verification.
	ModuleBasics = module.NewBasicManager(
		auth.AppModuleBasic{},
		genutil.AppModuleBasic{},
		bank.AppModuleBasic{},
		capability.AppModuleBasic{},
		staking.AppModuleBasic{},
		mint.AppModuleBasic{},
		distr.AppModuleBasic{},
		gov.NewAppModuleBasic(
			[]govclient.ProposalHandler{
				paramsclient.ProposalHandler,
				upgradeclient.LegacyProposalHandler,
				upgradeclient.LegacyCancelProposalHandler,
				gravityclient.ProposalHandler,
				corkclient.AddProposalHandler,
				corkclient.RemoveProposalHandler,
				axelarcorkclient.AddProposalHandler,
				axelarcorkclient.RemoveProposalHandler,
				axelarcorkclient.AddChainConfigurationHandler,
				axelarcorkclient.RemoveChainConfigurationHandler,
				axelarcorkclient.ScheduledCorkProposalHandler,
				axelarcorkclient.CommunityPoolEthereumSpendProposalHandler,
				axelarcorkclient.UpgradeAxelarProxyContractHandler,
				axelarcorkclient.CancelAxelarProxyContractUpgradeHandler,
				corkclient.ScheduledCorkProposalHandler,
				auctionclient.SetProposalHandler,
				pubsubclient.AddPublisherProposalHandler,
				pubsubclient.RemovePublisherProposalHandler,
				pubsubclient.AddDefaultSubscriptionProposalHandler,
				pubsubclient.RemoveDefaultSubscriptionProposalHandler,
				ibcclientclient.UpdateClientProposalHandler,
				ibcclientclient.UpgradeProposalHandler,
			},
		),
		params.AppModuleBasic{},
		crisis.AppModuleBasic{},
		slashing.AppModuleBasic{},
		ibc.AppModuleBasic{},
		ica.AppModuleBasic{},
		feegrantmodule.AppModuleBasic{},
		upgrade.AppModuleBasic{},
		evidence.AppModuleBasic{},
		ibctransfer.AppModuleBasic{},
		vesting.AppModuleBasic{},
		authzmodule.AppModuleBasic{},
		gravity.AppModuleBasic{},
		cork.AppModuleBasic{},
		axelarcork.AppModuleBasic{},
		cellarfees.AppModuleBasic{},
		incentives.AppModuleBasic{},
		auction.AppModuleBasic{},
		pubsub.AppModuleBasic{},
		addresses.AppModuleBasic{},
	)

	// module account permissions
	maccPerms = map[string][]string{
		authtypes.FeeCollectorName:     nil,
		distrtypes.ModuleName:          nil,
		minttypes.ModuleName:           {authtypes.Minter},
		stakingtypes.BondedPoolName:    {authtypes.Burner, authtypes.Staking},
		stakingtypes.NotBondedPoolName: {authtypes.Burner, authtypes.Staking},
		govtypes.ModuleName:            {authtypes.Burner},
		ibctransfertypes.ModuleName:    {authtypes.Minter, authtypes.Burner},
		icatypes.ModuleName:            nil,
		gravitytypes.ModuleName:        {authtypes.Minter, authtypes.Burner},
		cellarfeestypes.ModuleName:     nil,
		incentivestypes.ModuleName:     nil,
		axelarcorktypes.ModuleName:     nil,
		auctiontypes.ModuleName:        nil,
		pubsubtypes.ModuleName:         nil,
		addressestypes.ModuleName:      nil,
	}

	// module accounts that are allowed to receive tokens
	// incidentally this permission is also required to be able to send tokens from module accounts
	allowedReceivingModAcc = map[string]bool{
		axelarcorktypes.ModuleName: true,
		cellarfeestypes.ModuleName: true,
	}

	_ runtime.AppI            = (*SommelierApp)(nil)
	_ servertypes.Application = (*SommelierApp)(nil)
)

// SommelierApp extends an ABCI application, but with most of its parameters exported.
// They are exported for convenience in creating helper functions, as object
// capabilities aren't needed for testing.
type SommelierApp struct {
	*baseapp.BaseApp
	legacyAmino       *codec.LegacyAmino
	appCodec          codec.Codec
	interfaceRegistry types.InterfaceRegistry

	invCheckPeriod uint

	// keys to access the substores
	keys    map[string]*storetypes.KVStoreKey
	tkeys   map[string]*storetypes.TransientStoreKey
	memKeys map[string]*storetypes.MemoryStoreKey

	// SDK keepers
	AccountKeeper         authkeeper.AccountKeeper
	BankKeeper            bankkeeper.Keeper
	CapabilityKeeper      *capabilitykeeper.Keeper
	StakingKeeper         stakingkeeper.Keeper
	SlashingKeeper        slashingkeeper.Keeper
	MintKeeper            mintkeeper.Keeper
	DistrKeeper           distrkeeper.Keeper
	GovKeeper             govkeeper.Keeper
	ConsensusParamsKeeper consensusparamkeeper.Keeper
	CrisisKeeper          crisiskeeper.Keeper
	UpgradeKeeper         upgradekeeper.Keeper
	ParamsKeeper          paramskeeper.Keeper
	IBCKeeper             *ibckeeper.Keeper // IBC Keeper must be a pointer in the app, so we can SetRouter on it correctly
	ICAHostKeeper         icahostkeeper.Keeper
	EvidenceKeeper        evidencekeeper.Keeper
	TransferKeeper        ibctransferkeeper.Keeper
	GravityKeeper         gravitykeeper.Keeper
	AuthzKeeper           authzkeeper.Keeper
	FeeGrantKeeper        feegrantkeeper.Keeper

	// Sommelier keepers
	CorkKeeper       corkkeeper.Keeper
	AxelarCorkKeeper axelarcorkkeeper.Keeper
	CellarFeesKeeper cellarfeeskeeper.Keeper
	IncentivesKeeper incentiveskeeper.Keeper
	AuctionKeeper    auctionkeeper.Keeper
	PubsubKeeper     pubsubkeeper.Keeper
	AddressesKeeper  addresseskeeper.Keeper

	// make capability scoped keepers public for test purposes (IBC only)
	ScopedAxelarCorkKeeper capabilitykeeper.ScopedKeeper
	ScopedIBCKeeper        capabilitykeeper.ScopedKeeper
	ScopedTransferKeeper   capabilitykeeper.ScopedKeeper
	ScopedICAHostKeeper    capabilitykeeper.ScopedKeeper

	// the module manager
	mm *module.Manager

	configurator module.Configurator

	// simulation manager
	sm *module.SimulationManager
}

func init() {
	userHomeDir, err := os.UserHomeDir()
	if err != nil {
		panic(err)
	}

	DefaultNodeHome = filepath.Join(userHomeDir, ".sommelier")
}

// NewSommelierApp returns a reference to an initialized Sommelier.
func NewSommelierApp(
	logger log.Logger, db dbm.DB, traceStore io.Writer, loadLatest bool, skipUpgradeHeights map[int64]bool,
	homePath string, invCheckPeriod uint, encodingConfig appParams.EncodingConfig, appOpts servertypes.AppOptions, baseAppOptions ...func(*baseapp.BaseApp),
) *SommelierApp {

	appCodec := encodingConfig.Marshaler
	legacyAmino := encodingConfig.Amino
	interfaceRegistry := encodingConfig.InterfaceRegistry

	bApp := baseapp.NewBaseApp(appName, logger, db, encodingConfig.TxConfig.TxDecoder(), baseAppOptions...)
	bApp.SetCommitMultiStoreTracer(traceStore)
	bApp.SetVersion(version.Version)
	bApp.SetInterfaceRegistry(interfaceRegistry)

	keys := sdk.NewKVStoreKeys(
		authtypes.StoreKey,
		banktypes.StoreKey,
		stakingtypes.StoreKey,
		minttypes.StoreKey,
		distrtypes.StoreKey,
		slashingtypes.StoreKey,
		govtypes.StoreKey,
		consensusparamtypes.StoreKey,
		paramstypes.StoreKey,
		icaexported.StoreKey,
		upgradetypes.StoreKey,
		crisistypes.StoreKey,
		evidencetypes.StoreKey,
		ibctransfertypes.StoreKey,
		icahosttypes.StoreKey,
		//icaexported.StoreKey,
		capabilitytypes.StoreKey,
		gravitytypes.StoreKey,
		feegrant.StoreKey,
		authzkeeper.StoreKey,
		corktypes.StoreKey,
		axelarcorktypes.StoreKey,
		incentivestypes.StoreKey,
		auctiontypes.StoreKey,
		cellarfeestypes.StoreKey,
		pubsubtypes.StoreKey,
		addressestypes.StoreKey,
	)
	tkeys := sdk.NewTransientStoreKeys(paramstypes.TStoreKey)
	memKeys := sdk.NewMemoryStoreKeys(capabilitytypes.MemStoreKey)

	app := &SommelierApp{
		BaseApp:           bApp,
		legacyAmino:       legacyAmino,
		appCodec:          appCodec,
		interfaceRegistry: interfaceRegistry,
		invCheckPeriod:    invCheckPeriod,
		keys:              keys,
		tkeys:             tkeys,
		memKeys:           memKeys,
	}

	app.ParamsKeeper = initParamsKeeper(appCodec, legacyAmino, keys[paramstypes.StoreKey], tkeys[paramstypes.TStoreKey])

	authority := authtypes.NewModuleAddress(govtypes.ModuleName).String()

	// set the BaseApp's parameter store
	app.ConsensusParamsKeeper = consensusparamkeeper.NewKeeper(appCodec, keys[consensusparamtypes.StoreKey], authority)
	bApp.SetParamStore(&app.ConsensusParamsKeeper)

	// add capability keeper and ScopeToModule for ibc module
	app.CapabilityKeeper = capabilitykeeper.NewKeeper(appCodec, keys[capabilitytypes.StoreKey], memKeys[capabilitytypes.MemStoreKey])
	scopedIBCKeeper := app.CapabilityKeeper.ScopeToModule(icaexported.ModuleName)
	scopedTransferKeeper := app.CapabilityKeeper.ScopeToModule(ibctransfertypes.ModuleName)
	scopedICAHostKeeper := app.CapabilityKeeper.ScopeToModule(icahosttypes.SubModuleName)
	scopedAxelarCorkKeeper := app.CapabilityKeeper.ScopeToModule(axelarcorktypes.ModuleName)
	app.CapabilityKeeper.Seal()

	// add keepers
	app.AccountKeeper = authkeeper.NewAccountKeeper(
		appCodec, keys[authtypes.StoreKey], authtypes.ProtoBaseAccount, maccPerms, appParams.Bech32PrefixAccAddr, authority,
	)
	app.BankKeeper = bankkeeper.NewBaseKeeper(
		appCodec, keys[banktypes.StoreKey], app.AccountKeeper, app.BlockedAddrs(), authority,
	)
	stakingKeeper := stakingkeeper.NewKeeper(
		appCodec, keys[stakingtypes.StoreKey], app.AccountKeeper, app.BankKeeper, authority,
	)
	app.MintKeeper = mintkeeper.NewKeeper(
		appCodec, keys[minttypes.StoreKey], stakingKeeper,
		app.AccountKeeper, app.BankKeeper, authtypes.FeeCollectorName,
		authority,
	)
	app.DistrKeeper = distrkeeper.NewKeeper(
		appCodec, keys[distrtypes.StoreKey], app.AccountKeeper, app.BankKeeper,
		stakingKeeper, authtypes.FeeCollectorName,
		authority,
	)
	app.SlashingKeeper = slashingkeeper.NewKeeper(
		appCodec, legacyAmino, keys[slashingtypes.StoreKey], stakingKeeper, authority,
	)
	app.CrisisKeeper = *crisiskeeper.NewKeeper(
		appCodec, app.keys[crisistypes.StoreKey], invCheckPeriod, app.BankKeeper, authtypes.FeeCollectorName, authority,
	)

	app.FeeGrantKeeper = feegrantkeeper.NewKeeper(appCodec, keys[feegrant.StoreKey], app.AccountKeeper)
	app.UpgradeKeeper = *upgradekeeper.NewKeeper(skipUpgradeHeights, keys[upgradetypes.StoreKey], appCodec, homePath, app.BaseApp,
		authority,
	)

	// register the staking hooks
	// NOTE: stakingKeeper above is passed by reference, so that it will contain these hooks
	stakingKeeper.SetHooks(
		stakingtypes.NewMultiStakingHooks(
			app.DistrKeeper.Hooks(),
			app.SlashingKeeper.Hooks(),
		),
	)

	app.StakingKeeper = *stakingKeeper

	app.AuthzKeeper = authzkeeper.NewKeeper(keys[authzkeeper.StoreKey], appCodec, app.MsgServiceRouter(), app.AccountKeeper)

	// Create IBC Keeper
	app.IBCKeeper = ibckeeper.NewKeeper(
		appCodec,
		keys[icaexported.StoreKey],
		app.GetSubspace(icaexported.ModuleName),
		app.StakingKeeper,
		app.UpgradeKeeper,
		scopedIBCKeeper,
	)

	// todo: check if default power reduction is appropriate
	app.GravityKeeper = gravitykeeper.NewKeeper(
		appCodec, keys[gravitytypes.StoreKey], app.GetSubspace(gravitytypes.ModuleName),
		app.AccountKeeper, app.StakingKeeper, app.BankKeeper, app.SlashingKeeper,
		app.DistrKeeper, sdk.DefaultPowerReduction,
		app.ModuleAccountAddressesToNames([]string{cellarfeestypes.ModuleName}),
		app.ModuleAccountAddressesToNames([]string{distrtypes.ModuleName}),
	)

	app.PubsubKeeper = pubsubkeeper.NewKeeper(
		appCodec, keys[pubsubtypes.StoreKey], app.GetSubspace(pubsubtypes.ModuleName),
		app.StakingKeeper, app.GravityKeeper,
	)

	// create axelar cork keeper
	app.AxelarCorkKeeper = axelarcorkkeeper.NewKeeper(
		appCodec,
		keys[axelarcorktypes.StoreKey],
		app.GetSubspace(axelarcorktypes.ModuleName),
		app.AccountKeeper,
		app.BankKeeper,
		app.StakingKeeper,
		app.TransferKeeper, // will be nil, circular dependency avoided by calling SetTransferKeeper later
		app.DistrKeeper,
		app.IBCKeeper.ChannelKeeper,
		app.GravityKeeper,
		app.PubsubKeeper,
	)

	// Create Transfer Keepers
	app.TransferKeeper = ibctransferkeeper.NewKeeper(
		appCodec,
		keys[ibctransfertypes.StoreKey],
		app.GetSubspace(ibctransfertypes.ModuleName),
		// replacing channelkeeper as ICS4 Middleware for Axelar packet validation
		app.AxelarCorkKeeper,
		app.IBCKeeper.ChannelKeeper,
		&app.IBCKeeper.PortKeeper,
		app.AccountKeeper,
		app.BankKeeper,
		scopedTransferKeeper,
	)

	app.AxelarCorkKeeper.SetTransferKeeper(app.TransferKeeper)

	transferModule := ibctransfer.NewAppModule(app.TransferKeeper)
	transferIBCModule := ibctransfer.NewIBCModule(app.TransferKeeper)

	// Create the ICAHost Keeper
	app.ICAHostKeeper = icahostkeeper.NewKeeper(
		appCodec,
		app.keys[icahosttypes.StoreKey],
		app.GetSubspace(icahosttypes.SubModuleName),
		app.IBCKeeper.ChannelKeeper,
		app.IBCKeeper.ChannelKeeper,
		&app.IBCKeeper.PortKeeper,
		app.AccountKeeper,
		scopedICAHostKeeper,
		bApp.MsgServiceRouter(),
	)
	app.ICAHostKeeper.WithQueryRouter(app.GRPCQueryRouter())
	icaHostIBCModule := icahost.NewIBCModule(app.ICAHostKeeper)

	var transferStack ibcporttypes.IBCModule = transferIBCModule
	transferStack = axelarcork.NewIBCMiddleware(&app.AxelarCorkKeeper, transferStack)

	app.CorkKeeper = corkkeeper.NewKeeper(
		appCodec, keys[corktypes.StoreKey], app.GetSubspace(corktypes.ModuleName),
		app.StakingKeeper, app.GravityKeeper, app.PubsubKeeper,
	)

	app.AuctionKeeper = auctionkeeper.NewKeeper(
		appCodec, keys[auctiontypes.StoreKey], app.GetSubspace(auctiontypes.ModuleName),
		app.BankKeeper, app.AccountKeeper, map[string]bool{cellarfeestypes.ModuleName: true},
		map[string]bool{cellarfeestypes.ModuleName: true},
	)

	app.CellarFeesKeeper = cellarfeeskeeper.NewKeeper(
		appCodec, keys[cellarfeestypes.StoreKey], app.GetSubspace(cellarfeestypes.ModuleName),
		app.AccountKeeper, app.BankKeeper, app.MintKeeper, app.CorkKeeper, app.AuctionKeeper,
	)

	app.IncentivesKeeper = incentiveskeeper.NewKeeper(
		appCodec, keys[incentivestypes.StoreKey], app.GetSubspace(incentivestypes.ModuleName), app.DistrKeeper, app.BankKeeper, app.MintKeeper,
	)

	app.IncentivesKeeper = incentiveskeeper.NewKeeper(
		appCodec, keys[incentivestypes.StoreKey], app.GetSubspace(incentivestypes.ModuleName), app.DistrKeeper, app.BankKeeper, app.MintKeeper,
	)

	app.GravityKeeper = *app.GravityKeeper.SetHooks(
		gravitytypes.NewMultiGravityHooks(
			app.CorkKeeper.Hooks(),
		))

	app.AddressesKeeper = *addresseskeeper.NewKeeper(
		appCodec, keys[addressestypes.StoreKey], app.GetSubspace(addressestypes.ModuleName),
	)

	// register the proposal types
	govRouter := govtypesv1beta1.NewRouter()
	govRouter.AddRoute(govtypes.RouterKey, govtypesv1beta1.ProposalHandler).
		AddRoute(paramproposal.RouterKey, params.NewParamChangeProposalHandler(app.ParamsKeeper)).
		AddRoute(upgradetypes.RouterKey, upgrade.NewSoftwareUpgradeProposalHandler(&app.UpgradeKeeper)).
		AddRoute(ibcclienttypes.RouterKey, ibcclient.NewClientProposalHandler(app.IBCKeeper.ClientKeeper)).
		AddRoute(corktypes.RouterKey, cork.NewProposalHandler(app.CorkKeeper)).
		AddRoute(axelarcorktypes.RouterKey, axelarcork.NewProposalHandler(app.AxelarCorkKeeper)).
		AddRoute(gravitytypes.RouterKey, gravity.NewCommunityPoolEthereumSpendProposalHandler(app.GravityKeeper)).
		AddRoute(auctiontypes.RouterKey, auction.NewSetTokenPricesProposalHandler(app.AuctionKeeper)).
		AddRoute(pubsubtypes.RouterKey, pubsub.NewPubsubProposalHandler(app.PubsubKeeper))

	app.GovKeeper = *govkeeper.NewKeeper(
		appCodec,
		keys[govtypes.StoreKey],
		app.AccountKeeper,
		app.BankKeeper,
		stakingKeeper,
		app.MsgServiceRouter(),
		govtypes.DefaultConfig(),
		authority,
	)

	app.GovKeeper.SetLegacyRouter(govRouter)

	// Create static IBC router, add transfer route, then set and seal it
	ibcRouter := ibcporttypes.NewRouter()
	ibcRouter.
		AddRoute(ibctransfertypes.ModuleName, transferStack).
		AddRoute(icahosttypes.SubModuleName, icaHostIBCModule)
	app.IBCKeeper.SetRouter(ibcRouter)

	// create evidence keeper with router
	evidenceKeeper := evidencekeeper.NewKeeper(
		appCodec, keys[evidencetypes.StoreKey], &app.StakingKeeper, app.SlashingKeeper,
	)
	// If evidence needs to be handled for the app, set routes in router here and seal
	app.EvidenceKeeper = *evidenceKeeper

	app.setupUpgradeStoreLoaders()

	/****  Module Options ****/
	var skipGenesisInvariants = false
	opt := appOpts.Get(crisis.FlagSkipGenesisInvariants)
	if opt, ok := opt.(bool); ok {
		skipGenesisInvariants = opt
	}

	// NOTE: Any module instantiated in the module manager that is later modified
	// must be passed by reference here.
	app.mm = module.NewManager(
		genutil.NewAppModule(
			app.AccountKeeper, app.StakingKeeper, app.BaseApp.DeliverTx,
			encodingConfig.TxConfig,
		),
		auth.NewAppModule(appCodec, app.AccountKeeper, nil, app.GetSubspace(authtypes.ModuleName)),
		vesting.NewAppModule(app.AccountKeeper, app.BankKeeper),
		bank.NewAppModule(appCodec, app.BankKeeper, app.AccountKeeper, app.GetSubspace(banktypes.ModuleName)),
		capability.NewAppModule(appCodec, *app.CapabilityKeeper, false),
		crisis.NewAppModule(&app.CrisisKeeper, skipGenesisInvariants, app.GetSubspace(crisistypes.ModuleName)),
		gov.NewAppModule(appCodec, &app.GovKeeper, app.AccountKeeper, app.BankKeeper, app.GetSubspace(govtypes.ModuleName)),
		feegrantmodule.NewAppModule(appCodec, app.AccountKeeper, app.BankKeeper, app.FeeGrantKeeper, app.interfaceRegistry),
		mint.NewAppModule(appCodec, app.MintKeeper, app.AccountKeeper, nil, app.GetSubspace(minttypes.ModuleName)),
		slashing.NewAppModule(appCodec, app.SlashingKeeper, app.AccountKeeper, app.BankKeeper, app.StakingKeeper, app.GetSubspace(slashingtypes.ModuleName)),
		distr.NewAppModule(appCodec, app.DistrKeeper, app.AccountKeeper, app.BankKeeper, app.StakingKeeper, app.GetSubspace(distrtypes.ModuleName)),
		staking.NewAppModule(appCodec, &app.StakingKeeper, app.AccountKeeper, app.BankKeeper, app.GetSubspace(stakingtypes.ModuleName)),
		upgrade.NewAppModule(&app.UpgradeKeeper),
		evidence.NewAppModule(app.EvidenceKeeper),
		ibc.NewAppModule(app.IBCKeeper),
		ica.NewAppModule(nil, &app.ICAHostKeeper),
		params.NewAppModule(app.ParamsKeeper),
		transferModule,
		axelarcork.NewAppModule(app.AxelarCorkKeeper, appCodec),
		gravity.NewAppModule(app.GravityKeeper, app.BankKeeper),
		authzmodule.NewAppModule(appCodec, app.AuthzKeeper, app.AccountKeeper, app.BankKeeper, app.interfaceRegistry),
		cork.NewAppModule(app.CorkKeeper, appCodec),
		incentives.NewAppModule(app.IncentivesKeeper, app.DistrKeeper, app.BankKeeper, app.MintKeeper, appCodec),
		cellarfees.NewAppModule(app.CellarFeesKeeper, appCodec, app.AccountKeeper, app.BankKeeper, app.MintKeeper, app.CorkKeeper, app.AuctionKeeper),
		auction.NewAppModule(app.AuctionKeeper, app.BankKeeper, app.AccountKeeper, appCodec),
		pubsub.NewAppModule(appCodec, app.PubsubKeeper, app.StakingKeeper, app.GravityKeeper),
		addresses.NewAppModule(appCodec, app.AddressesKeeper),
	)

	// During begin block slashing happens after distr.BeginBlocker so that
	// there is nothing left over in the validator fee pool, so as to keep the
	// CanWithdrawInvariant invariant.
	// NOTE: staking module is required if HistoricalEntries param > 0
	// NOTE: capability module's beginblocker must come before any modules using capabilities (e.g. IBC)
	app.mm.SetOrderBeginBlockers(
		upgradetypes.ModuleName,
		capabilitytypes.ModuleName,
		minttypes.ModuleName,
		distrtypes.ModuleName,
		slashingtypes.ModuleName,
		evidencetypes.ModuleName,
		stakingtypes.ModuleName,
		authtypes.ModuleName,
		banktypes.ModuleName,
		govtypes.ModuleName,
		crisistypes.ModuleName,
		icaexported.ModuleName,
		ibctransfertypes.ModuleName,
		icatypes.ModuleName,
		genutiltypes.ModuleName,
		authz.ModuleName,
		feegrant.ModuleName,
		paramstypes.ModuleName,
		vestingtypes.ModuleName,
		incentivestypes.ModuleName,
		gravitytypes.ModuleName,
		corktypes.ModuleName,
		axelarcorktypes.ModuleName,
		cellarfeestypes.ModuleName,
		auctiontypes.ModuleName,
		pubsubtypes.ModuleName,
		addressestypes.ModuleName,
	)

	// NOTE gov must come before staking
	app.mm.SetOrderEndBlockers(
		crisistypes.ModuleName,
		govtypes.ModuleName,
		stakingtypes.ModuleName,
		icaexported.ModuleName,
		ibctransfertypes.ModuleName,
		icatypes.ModuleName,
		capabilitytypes.ModuleName,
		authtypes.ModuleName,
		banktypes.ModuleName,
		distrtypes.ModuleName,
		slashingtypes.ModuleName,
		minttypes.ModuleName,
		genutiltypes.ModuleName,
		evidencetypes.ModuleName,
		authz.ModuleName,
		feegrant.ModuleName,
		paramstypes.ModuleName,
		upgradetypes.ModuleName,
		vestingtypes.ModuleName,
		incentivestypes.ModuleName,
		gravitytypes.ModuleName,
		corktypes.ModuleName,
		axelarcorktypes.ModuleName,
		cellarfeestypes.ModuleName,
		auctiontypes.ModuleName,
		pubsubtypes.ModuleName,
		addressestypes.ModuleName,
	)

	// NOTE: The genutils module must occur after staking so that pools are
	// properly initialized with tokens from genesis accounts.
	// NOTE: Capability module must occur first so that it can initialize any capabilities
	// so that other modules that want to create or claim capabilities afterwards in InitChain
	// can do so safely.
	// NOTE: During cellarfees development, moved crisis to end of order because when its InitChain
	// gets run, it enforces all registered invariants before moving on to other modules. This means
	// if a module listed after crisis has invariants that depend on their own state being initialized,
	// their keeper will panic when a nil state is returned.
	app.mm.SetOrderInitGenesis(
		capabilitytypes.ModuleName,
		authtypes.ModuleName,
		banktypes.ModuleName,
		distrtypes.ModuleName,
		govtypes.ModuleName,
		stakingtypes.ModuleName,
		slashingtypes.ModuleName,
		minttypes.ModuleName,
		crisistypes.ModuleName,
		genutiltypes.ModuleName,
		icaexported.ModuleName,
		ibctransfertypes.ModuleName,
		icatypes.ModuleName,
		evidencetypes.ModuleName,
		authz.ModuleName,
		feegrant.ModuleName,
		paramstypes.ModuleName,
		upgradetypes.ModuleName,
		vestingtypes.ModuleName,
		gravitytypes.ModuleName,
		incentivestypes.ModuleName,
		corktypes.ModuleName,
		axelarcorktypes.ModuleName,
		cellarfeestypes.ModuleName,
		auctiontypes.ModuleName,
		pubsubtypes.ModuleName,
		addressestypes.ModuleName,
	)

	app.mm.RegisterInvariants(&app.CrisisKeeper)
	app.configurator = module.NewConfigurator(app.appCodec, app.MsgServiceRouter(), app.GRPCQueryRouter())
	app.mm.RegisterServices(app.configurator)

	app.setupUpgradeHandlers()

	// create the simulation manager and define the order of the modules for deterministic simulations
	//
	// NOTE: this is not required apps that don't use the simulator for fuzz testing
	// transactions
	app.sm = module.NewSimulationManager(
		auth.NewAppModule(appCodec, app.AccountKeeper, nil, app.GetSubspace(authtypes.ModuleName)),
		bank.NewAppModule(appCodec, app.BankKeeper, app.AccountKeeper, app.GetSubspace(banktypes.ModuleName)),
		capability.NewAppModule(appCodec, *app.CapabilityKeeper, false),
		gov.NewAppModule(appCodec, &app.GovKeeper, app.AccountKeeper, app.BankKeeper, app.GetSubspace(govtypes.ModuleName)),
		feegrantmodule.NewAppModule(appCodec, app.AccountKeeper, app.BankKeeper, app.FeeGrantKeeper, app.interfaceRegistry),
		mint.NewAppModule(appCodec, app.MintKeeper, app.AccountKeeper, nil, app.GetSubspace(minttypes.ModuleName)),
		slashing.NewAppModule(appCodec, app.SlashingKeeper, app.AccountKeeper, app.BankKeeper, app.StakingKeeper, app.GetSubspace(slashingtypes.ModuleName)),
		distr.NewAppModule(appCodec, app.DistrKeeper, app.AccountKeeper, app.BankKeeper, app.StakingKeeper, app.GetSubspace(distrtypes.ModuleName)),
		staking.NewAppModule(appCodec, &app.StakingKeeper, app.AccountKeeper, app.BankKeeper, app.GetSubspace(stakingtypes.ModuleName)),
		evidence.NewAppModule(app.EvidenceKeeper),
		ibc.NewAppModule(app.IBCKeeper),
		ica.NewAppModule(nil, &app.ICAHostKeeper),
		params.NewAppModule(app.ParamsKeeper),
		transferModule,
		axelarcork.NewAppModule(app.AxelarCorkKeeper, appCodec),
		gravity.NewAppModule(app.GravityKeeper, app.BankKeeper),
		authzmodule.NewAppModule(appCodec, app.AuthzKeeper, app.AccountKeeper, app.BankKeeper, app.interfaceRegistry),
		cork.NewAppModule(app.CorkKeeper, appCodec),
		incentives.NewAppModule(app.IncentivesKeeper, app.DistrKeeper, app.BankKeeper, app.MintKeeper, appCodec),
		cellarfees.NewAppModule(app.CellarFeesKeeper, appCodec, app.AccountKeeper, app.BankKeeper, app.MintKeeper, app.CorkKeeper, app.AuctionKeeper),
		auction.NewAppModule(app.AuctionKeeper, app.BankKeeper, app.AccountKeeper, appCodec),
		pubsub.NewAppModule(appCodec, app.PubsubKeeper, app.StakingKeeper, app.GravityKeeper),
		addresses.NewAppModule(appCodec, app.AddressesKeeper),
	)

	app.sm.RegisterStoreDecoders()

	// initialize stores
	app.MountKVStores(keys)
	app.MountTransientStores(tkeys)
	app.MountMemoryStores(memKeys)

	autocliv1.RegisterQueryServer(app.GRPCQueryRouter(), runtimeservices.NewAutoCLIQueryService(app.mm.Modules))

	reflectionSvc, err := runtimeservices.NewReflectionService()
	if err != nil {
		panic(err)
	}
	reflectionv1.RegisterReflectionServiceServer(app.GRPCQueryRouter(), reflectionSvc)

	// initialize BaseApp
	app.SetInitChainer(app.InitChainer)
	app.SetBeginBlocker(app.BeginBlocker)

	// TODO(bolten): is there any reason we might want to customize the ante handler to include
	// the redundant relay IBC decorator?
	anteHandler, err := ante.NewAnteHandler(
		ante.HandlerOptions{
			AccountKeeper:   app.AccountKeeper,
			BankKeeper:      app.BankKeeper,
			SignModeHandler: encodingConfig.TxConfig.SignModeHandler(),
			FeegrantKeeper:  app.FeeGrantKeeper,
			SigGasConsumer:  ante.DefaultSigVerificationGasConsumer,
		},
	)
	if err != nil {
		panic(fmt.Errorf("failed to create ante handler: %s", err))
	}

	app.SetAnteHandler(anteHandler)
	app.SetEndBlocker(app.EndBlocker)

	if loadLatest {
		if err := app.LoadLatestVersion(); err != nil {
			tmos.Exit(err.Error())
		}
	}
	app.ScopedIBCKeeper = scopedIBCKeeper
	app.ScopedTransferKeeper = scopedTransferKeeper
	app.ScopedICAHostKeeper = scopedICAHostKeeper
	app.ScopedAxelarCorkKeeper = scopedAxelarCorkKeeper

	return app
}

// MakeCodecs constructs the *std.Codec and *codec.LegacyAmino instances used by
// simapp. It is useful for tests and clients who do not want to construct the
// full simapp
func MakeCodecs() (codec.Codec, *codec.LegacyAmino) {
	conf := MakeEncodingConfig()
	return conf.Marshaler, conf.Amino
}

// Name returns the name of the App
func (app *SommelierApp) Name() string { return app.BaseApp.Name() }

// BeginBlocker application updates every begin block
func (app *SommelierApp) BeginBlocker(ctx sdk.Context, req abci.RequestBeginBlock) abci.ResponseBeginBlock {
	return app.mm.BeginBlock(ctx, req)
}

// EndBlocker application updates every end block
func (app *SommelierApp) EndBlocker(ctx sdk.Context, req abci.RequestEndBlock) abci.ResponseEndBlock {
	return app.mm.EndBlock(ctx, req)
}

// InitChainer application update at chain initialization
func (app *SommelierApp) InitChainer(ctx sdk.Context, req abci.RequestInitChain) abci.ResponseInitChain {
	var genesisState GenesisState
	if err := tmjson.Unmarshal(req.AppStateBytes, &genesisState); err != nil {
		panic(err)
	}

	app.UpgradeKeeper.SetModuleVersionMap(ctx, app.mm.GetVersionMap())

	return app.mm.InitGenesis(ctx, app.appCodec, genesisState)
}

// LoadHeight loads a particular height
func (app *SommelierApp) LoadHeight(height int64) error {
	return app.LoadVersion(height)
}

// ModuleAccountAddrs returns all the app's module account addresses.
func (app *SommelierApp) ModuleAccountAddrs() map[string]bool {
	modAccAddrs := make(map[string]bool)
	for acc := range maccPerms {
		modAccAddrs[authtypes.NewModuleAddress(acc).String()] = true
	}

	return modAccAddrs
}

// ModuleAccountAddressesToNames returns a map of module account address to module name
func (app *SommelierApp) ModuleAccountAddressesToNames(moduleAccounts []string) map[string]string {
	modAccNames := make(map[string]string)
	for _, acc := range moduleAccounts {
		modAccNames[authtypes.NewModuleAddress(acc).String()] = acc
	}

	return modAccNames
}

// BlockedAddrs returns all the app's module account addresses that are not
// allowed to receive external tokens.
func (app *SommelierApp) BlockedAddrs() map[string]bool {
	blockedAddrs := make(map[string]bool)
	for acc := range maccPerms {
		blockedAddrs[authtypes.NewModuleAddress(acc).String()] = !allowedReceivingModAcc[acc]
	}

	return blockedAddrs
}

// LegacyAmino returns SimApp's amino codec.
//
// NOTE: This is solely to be used for testing purposes as it may be desirable
// for modules to register their own custom testing types.
func (app *SommelierApp) LegacyAmino() *codec.LegacyAmino {
	return app.legacyAmino
}

// AppCodec returns Sommelier's app codec.
//
// NOTE: This is solely to be used for testing purposes as it may be desirable
// for modules to register their own custom testing types.
func (app *SommelierApp) AppCodec() codec.Codec {
	return app.appCodec
}

// InterfaceRegistry returns Sommelier's InterfaceRegistry
func (app *SommelierApp) InterfaceRegistry() types.InterfaceRegistry {
	return app.interfaceRegistry
}

// GetKey returns the KVStoreKey for the provided store key.
//
// NOTE: This is solely to be used for testing purposes.
func (app *SommelierApp) GetKey(storeKey string) *storetypes.KVStoreKey {
	return app.keys[storeKey]
}

// GetTKey returns the TransientStoreKey for the provided store key.
//
// NOTE: This is solely to be used for testing purposes.
func (app *SommelierApp) GetTKey(storeKey string) *storetypes.TransientStoreKey {
	return app.tkeys[storeKey]
}

// GetMemKey returns the MemStoreKey for the provided mem key.
//
// NOTE: This is solely used for testing purposes.
func (app *SommelierApp) GetMemKey(storeKey string) *storetypes.MemoryStoreKey {
	return app.memKeys[storeKey]
}

// GetSubspace returns a param subspace for a given module name.
//
// NOTE: This is solely to be used for testing purposes.
func (app *SommelierApp) GetSubspace(moduleName string) paramstypes.Subspace {
	subspace, _ := app.ParamsKeeper.GetSubspace(moduleName)
	return subspace
}

// SimulationManager implements the SimulationApp interface
func (app *SommelierApp) SimulationManager() *module.SimulationManager {
	return app.sm
}

// RegisterAPIRoutes registers all application module routes with the provided
// API server.
func (app *SommelierApp) RegisterAPIRoutes(apiSvr *api.Server, apiConfig config.APIConfig) {
	clientCtx := apiSvr.ClientCtx
	// Register new tx routes from grpc-gateway.
	authtx.RegisterGRPCGatewayRoutes(clientCtx, apiSvr.GRPCGatewayRouter)
	// Register new tendermint queries routes from grpc-gateway.
	tmservice.RegisterGRPCGatewayRoutes(clientCtx, apiSvr.GRPCGatewayRouter)
	// Register grpc-gateway routes for all modules.
	ModuleBasics.RegisterGRPCGatewayRoutes(clientCtx, apiSvr.GRPCGatewayRouter)
	nodeservice.RegisterGRPCGatewayRoutes(clientCtx, apiSvr.GRPCGatewayRouter)

	// register swagger API from root so that other applications can override easily
	if apiConfig.Swagger {
		RegisterSwaggerAPI(clientCtx, apiSvr.Router)
	}
}

// RegisterNodeService registers the node gRPC Query service.
func (app *SommelierApp) RegisterNodeService(clientCtx client.Context) {
	nodeservice.RegisterNodeService(clientCtx, app.GRPCQueryRouter())
}

// RegisterTxService implements the Application.RegisterTxService method.
func (app *SommelierApp) RegisterTxService(clientCtx client.Context) {
	authtx.RegisterTxService(app.BaseApp.GRPCQueryRouter(), clientCtx, app.BaseApp.Simulate, app.interfaceRegistry)
}

// RegisterTendermintService implements the application.RegisterTendermintService method
func (app *SommelierApp) RegisterTendermintService(clientCtx client.Context) {
	tmservice.RegisterTendermintService(clientCtx, app.BaseApp.GRPCQueryRouter(), app.interfaceRegistry, app.Query)
}

// RegisterSwaggerAPI registers swagger route with API Server
func RegisterSwaggerAPI(_ client.Context, rtr *mux.Router) {
	statikFS, err := fs.New()
	if err != nil {
		panic(err)
	}

	staticServer := http.FileServer(statikFS)
	rtr.PathPrefix("/swagger/").Handler(http.StripPrefix("/swagger/", staticServer))
}

// GetMaccPerms returns a copy of the module account permissions
func GetMaccPerms() map[string][]string {
	dupMaccPerms := make(map[string][]string)
	for k, v := range maccPerms {
		dupMaccPerms[k] = v
	}
	return dupMaccPerms
}

// initParamsKeeper init params keeper and its subspaces
func initParamsKeeper(appCodec codec.BinaryCodec, legacyAmino *codec.LegacyAmino, key, tkey storetypes.StoreKey) paramskeeper.Keeper {
	paramsKeeper := paramskeeper.NewKeeper(appCodec, legacyAmino, key, tkey)

	paramsKeeper.Subspace(authtypes.ModuleName)
	paramsKeeper.Subspace(banktypes.ModuleName)
	paramsKeeper.Subspace(stakingtypes.ModuleName)
	paramsKeeper.Subspace(minttypes.ModuleName)
	paramsKeeper.Subspace(distrtypes.ModuleName)
	paramsKeeper.Subspace(slashingtypes.ModuleName)
	paramsKeeper.Subspace(govtypes.ModuleName)
	paramsKeeper.Subspace(crisistypes.ModuleName)
	paramsKeeper.Subspace(ibctransfertypes.ModuleName)
	paramsKeeper.Subspace(icaexported.ModuleName)
	paramsKeeper.Subspace(icahosttypes.SubModuleName)
	paramsKeeper.Subspace(gravitytypes.ModuleName)
	paramsKeeper.Subspace(corktypes.ModuleName)
	paramsKeeper.Subspace(axelarcorktypes.ModuleName)
	paramsKeeper.Subspace(cellarfeestypes.ModuleName)
	paramsKeeper.Subspace(incentivestypes.ModuleName)
	paramsKeeper.Subspace(auctiontypes.ModuleName)
	paramsKeeper.Subspace(pubsubtypes.ModuleName)
	paramsKeeper.Subspace(addressestypes.ModuleName)

	return paramsKeeper
}

func (app *SommelierApp) setupUpgradeStoreLoaders() {
	upgradeInfo, err := app.UpgradeKeeper.ReadUpgradeInfoFromDisk()
	if err != nil {
		panic(fmt.Sprintf("failed to read upgrade info from disk %s", err))
	}

	if app.UpgradeKeeper.IsSkipHeight(upgradeInfo.Height) {
		return
	}

	var storeUpgrades *storetypes.StoreUpgrades = nil

	// TODO: Add v8 store loader when writing upgrade handler

	if storeUpgrades != nil {
		app.SetStoreLoader(upgradetypes.UpgradeStoreLoader(upgradeInfo.Height, storeUpgrades))
	}
}

func (app *SommelierApp) setupUpgradeHandlers() {
	// TODO: Add v8 upgrade handler
}
