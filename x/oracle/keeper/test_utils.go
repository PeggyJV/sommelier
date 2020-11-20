// nolint:deadcode unused noalias
package keeper

import (
	"testing"

	"github.com/peggyjv/sommelier/x/oracle/types"
	"github.com/stretchr/testify/require"

	"time"

	abci "github.com/tendermint/tendermint/abci/types"
	"github.com/tendermint/tendermint/crypto"
	"github.com/tendermint/tendermint/crypto/secp256k1"
	"github.com/tendermint/tendermint/libs/log"
	dbm "github.com/tendermint/tm-db"

	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/store"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/auth"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	"github.com/cosmos/cosmos-sdk/x/bank"
	distr "github.com/cosmos/cosmos-sdk/x/distribution"
	"github.com/cosmos/cosmos-sdk/x/params"
	"github.com/cosmos/cosmos-sdk/x/staking"
)

var (
	PubKeys = []crypto.PubKey{
		secp256k1.GenPrivKey().PubKey(),
		secp256k1.GenPrivKey().PubKey(),
		secp256k1.GenPrivKey().PubKey(),
		secp256k1.GenPrivKey().PubKey(),
		secp256k1.GenPrivKey().PubKey(),
	}

	Addrs = []sdk.AccAddress{
		sdk.AccAddress(PubKeys[0].Address()),
		sdk.AccAddress(PubKeys[1].Address()),
		sdk.AccAddress(PubKeys[2].Address()),
		sdk.AccAddress(PubKeys[3].Address()),
		sdk.AccAddress(PubKeys[4].Address()),
	}

	ValAddrs = []sdk.ValAddress{
		sdk.ValAddress(PubKeys[0].Address()),
		sdk.ValAddress(PubKeys[1].Address()),
		sdk.ValAddress(PubKeys[2].Address()),
		sdk.ValAddress(PubKeys[3].Address()),
		sdk.ValAddress(PubKeys[4].Address()),
	}

	InitTokens = sdk.TokensFromConsensusPower(200)
	InitCoins  = sdk.NewCoins(sdk.NewCoin(types.MicroLunaDenom, InitTokens))

	OracleDecPrecision = 8
)

// TestInput nolint
type TestInput struct {
	Ctx           sdk.Context
	Cdc           *codec.Codec
	AccKeeper     auth.AccountKeeper
	BankKeeper    bank.Keeper
	OracleKeeper  Keeper
	StakingKeeper staking.Keeper
	DistrKeeper   distr.Keeper
}

func newTestCodec() *codec.Codec {
	cdc := codec.New()

	types.RegisterCodec(cdc)
	auth.RegisterCodec(cdc)
	sdk.RegisterCodec(cdc)
	codec.RegisterCrypto(cdc)
	staking.RegisterCodec(cdc)
	distr.RegisterCodec(cdc)
	params.RegisterCodec(cdc)

	return cdc
}

// CreateTestInput nolint
func CreateTestInput(t *testing.T) TestInput {
	keyAcc := sdk.NewKVStoreKey(auth.StoreKey)
	keyParams := sdk.NewKVStoreKey(params.StoreKey)
	tKeyParams := sdk.NewTransientStoreKey(params.TStoreKey)
	keyOracle := sdk.NewKVStoreKey(types.StoreKey)
	keyStaking := sdk.NewKVStoreKey(staking.StoreKey)
	keyDistr := sdk.NewKVStoreKey(distr.StoreKey)

	cdc := newTestCodec()
	db := dbm.NewMemDB()
	ms := store.NewCommitMultiStore(db)
	ctx := sdk.NewContext(ms, abci.Header{Time: time.Now().UTC()}, false, log.NewNopLogger())

	ms.MountStoreWithDB(keyAcc, sdk.StoreTypeIAVL, db)
	ms.MountStoreWithDB(tKeyParams, sdk.StoreTypeTransient, db)
	ms.MountStoreWithDB(keyParams, sdk.StoreTypeIAVL, db)
	ms.MountStoreWithDB(keyOracle, sdk.StoreTypeIAVL, db)
	ms.MountStoreWithDB(keyStaking, sdk.StoreTypeIAVL, db)
	ms.MountStoreWithDB(keyDistr, sdk.StoreTypeIAVL, db)

	require.NoError(t, ms.LoadLatestVersion())

	blackListAddrs := map[string]bool{
		auth.FeeCollectorName:     true,
		staking.NotBondedPoolName: true,
		staking.BondedPoolName:    true,
		distr.ModuleName:          true,
		types.ModuleName:          true,
	}

	paramsKeeper := params.NewKeeper(cdc, keyParams, tKeyParams)
	accountKeeper := auth.NewAccountKeeper(cdc, keyAcc, paramsKeeper.Subspace(auth.DefaultParamspace), auth.ProtoBaseAccount)
	bankKeeper := bank.NewBaseKeeper(accountKeeper, paramsKeeper.Subspace(bank.DefaultParamspace), blackListAddrs)

	maccPerms := map[string][]string{
		auth.FeeCollectorName:     nil,
		staking.NotBondedPoolName: {authtypes.Burner, authtypes.Staking},
		staking.BondedPoolName:    {authtypes.Burner, authtypes.Staking},
		distr.ModuleName:          nil,
		types.ModuleName:          nil,
	}

	totalSupply := sdk.NewCoins(sdk.NewCoin(types.MicroLunaDenom, InitTokens.MulRaw(int64(len(Addrs)))))
	bankKeeper.SetSupply(ctx, banktypes.NewSupply(totalSupply))

	stakingKeeper := staking.NewKeeper(
		cdc,
		keyStaking,
		accountKeeper, bankKeeper, paramsKeeper.Subspace(staking.DefaultParamspace),
	)

	distrKeeper := distr.NewKeeper(
		cdc,
		keyDistr, paramsKeeper.Subspace(distr.DefaultParamspace),
		accountKeeper, bankKeeper, stakingKeeper, auth.FeeCollectorName, blackListAddrs)

	distrKeeper.SetFeePool(ctx, distr.InitialFeePool())
	distrParams := distr.DefaultParams()
	distrParams.CommunityTax = sdk.NewDecWithPrec(2, 2)
	distrParams.BaseProposerReward = sdk.NewDecWithPrec(1, 2)
	distrParams.BonusProposerReward = sdk.NewDecWithPrec(4, 2)
	distrKeeper.SetParams(ctx, distrParams)

	feeCollectorAcc := authtypes.NewEmptyModuleAccount(auth.FeeCollectorName)
	notBondedPool := authtypes.NewEmptyModuleAccount(staking.NotBondedPoolName, authtypes.Burner, authtypes.Staking)
	bondPool := authtypes.NewEmptyModuleAccount(staking.BondedPoolName, authtypes.Burner, authtypes.Staking)
	distrAcc := authtypes.NewEmptyModuleAccount(distr.ModuleName)
	oracleAcc := authtypes.NewEmptyModuleAccount(types.ModuleName, authtypes.Minter)

	notBondedPool.SetCoins(sdk.NewCoins(sdk.NewCoin(types.MicroLunaDenom, InitTokens.MulRaw(int64(len(Addrs))))))

	accountKeeper.SetModuleAccount(ctx, feeCollectorAcc)
	accountKeeper.SetModuleAccount(ctx, bondPool)
	accountKeeper.SetModuleAccount(ctx, notBondedPool)
	accountKeeper.SetModuleAccount(ctx, distrAcc)
	accountKeeper.SetModuleAccount(ctx, oracleAcc)

	genesis := staking.DefaultGenesisState()
	genesis.Params.BondDenom = types.MicroLunaDenom
	_ = staking.InitGenesis(ctx, genesis)

	for _, addr := range Addrs {
		_, err := bankKeeper.AddCoins(ctx, sdk.AccAddress(addr), InitCoins)
		require.NoError(t, err)
	}

	keeper := NewKeeper(cdc, keyOracle, paramsKeeper.Subspace(types.DefaultParamspace), distrKeeper, stakingKeeper, bankKeeper, distr.ModuleName)

	defaults := types.DefaultParams()
	keeper.SetParams(ctx, defaults)

	for _, denom := range defaults.Whitelist {
		keeper.SetTobinTax(ctx, denom.Name, denom.TobinTax)
	}

	stakingKeeper.SetHooks(staking.NewMultiStakingHooks(distrKeeper.Hooks()))

	return TestInput{ctx, cdc, accountKeeper, bankKeeper, keeper, stakingKeeper, distrKeeper}
}

func NewTestMsgCreateValidator(address sdk.ValAddress, pubKey crypto.PubKey, amt sdk.Int) staking.MsgCreateValidator {
	commission := staking.NewCommissionRates(sdk.ZeroDec(), sdk.ZeroDec(), sdk.ZeroDec())
	return staking.NewMsgCreateValidator(
		address, pubKey, sdk.NewCoin(types.MicroLunaDenom, amt),
		staking.Description{}, commission, sdk.OneInt(),
	)
}
