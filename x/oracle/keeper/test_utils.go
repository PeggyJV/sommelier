// nolint:deadcode unused noalias
package keeper

import (
	"testing"

	"github.com/peggyjv/sommelier/x/oracle/types"
	"github.com/stretchr/testify/require"

	"time"

	"github.com/cosmos/cosmos-sdk/crypto/keys/ed25519"
	"github.com/cosmos/cosmos-sdk/crypto/keys/secp256k1"
	crypto "github.com/cosmos/cosmos-sdk/crypto/types"
	"github.com/tendermint/tendermint/libs/log"
	tmproto "github.com/tendermint/tendermint/proto/tendermint/types"
	dbm "github.com/tendermint/tm-db"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/codec"
	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	ccodec "github.com/cosmos/cosmos-sdk/crypto/codec"
	"github.com/cosmos/cosmos-sdk/store"
	sdk "github.com/cosmos/cosmos-sdk/types"
	authkeeper "github.com/cosmos/cosmos-sdk/x/auth/keeper"
	"github.com/cosmos/cosmos-sdk/x/auth/tx"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	bankkeeper "github.com/cosmos/cosmos-sdk/x/bank/keeper"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	distrkeeper "github.com/cosmos/cosmos-sdk/x/distribution/keeper"
	distrtypes "github.com/cosmos/cosmos-sdk/x/distribution/types"
	paramskeeper "github.com/cosmos/cosmos-sdk/x/params/keeper"
	paramstypes "github.com/cosmos/cosmos-sdk/x/params/types"
	paramsproposal "github.com/cosmos/cosmos-sdk/x/params/types/proposal"
	"github.com/cosmos/cosmos-sdk/x/staking"
	stakingkeeper "github.com/cosmos/cosmos-sdk/x/staking/keeper"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
)

var (
	PubKeys = []crypto.PubKey{
		ed25519.GenPrivKey().PubKey(),
		ed25519.GenPrivKey().PubKey(),
		ed25519.GenPrivKey().PubKey(),
		ed25519.GenPrivKey().PubKey(),
		ed25519.GenPrivKey().PubKey(),
	}

	AccPubKeys = []crypto.PubKey{
		secp256k1.GenPrivKey().PubKey(),
		secp256k1.GenPrivKey().PubKey(),
		secp256k1.GenPrivKey().PubKey(),
		secp256k1.GenPrivKey().PubKey(),
		secp256k1.GenPrivKey().PubKey(),
	}

	Addrs = []sdk.AccAddress{
		sdk.AccAddress(AccPubKeys[0].Address()),
		sdk.AccAddress(AccPubKeys[1].Address()),
		sdk.AccAddress(AccPubKeys[2].Address()),
		sdk.AccAddress(AccPubKeys[3].Address()),
		sdk.AccAddress(AccPubKeys[4].Address()),
	}

	ValAddrs = []sdk.ValAddress{
		sdk.ValAddress(AccPubKeys[0].Address()),
		sdk.ValAddress(AccPubKeys[1].Address()),
		sdk.ValAddress(AccPubKeys[2].Address()),
		sdk.ValAddress(AccPubKeys[3].Address()),
		sdk.ValAddress(AccPubKeys[4].Address()),
	}

	InitTokens = sdk.TokensFromConsensusPower(200)
	InitCoins  = sdk.NewCoins(sdk.NewCoin(types.MicroLunaDenom, InitTokens))

	OracleDecPrecision = 8
)

// TestInput nolint
type TestInput struct {
	Ctx           sdk.Context
	Cdc           codec.BinaryMarshaler
	AccKeeper     authkeeper.AccountKeeper
	BankKeeper    bankkeeper.Keeper
	OracleKeeper  Keeper
	StakingKeeper stakingkeeper.Keeper
	DistrKeeper   distrkeeper.Keeper
}

func newTestCodec() *codec.LegacyAmino {
	cdc := codec.NewLegacyAmino()

	authtypes.RegisterLegacyAminoCodec(cdc)
	banktypes.RegisterLegacyAminoCodec(cdc)
	sdk.RegisterLegacyAminoCodec(cdc)
	ccodec.RegisterCrypto(cdc)
	stakingtypes.RegisterLegacyAminoCodec(cdc)
	distrtypes.RegisterLegacyAminoCodec(cdc)
	paramsproposal.RegisterLegacyAminoCodec(cdc)

	return cdc
}

// NewTestMarshaler returns a new ProtoCodecMarshaler
func NewTestMarshaler() codec.ProtoCodecMarshaler {
	ir := codectypes.NewInterfaceRegistry()
	types.RegisterInterfaces(ir)
	authtypes.RegisterInterfaces(ir)
	banktypes.RegisterInterfaces(ir)
	sdk.RegisterInterfaces(ir)
	ccodec.RegisterInterfaces(ir)
	stakingtypes.RegisterInterfaces(ir)
	distrtypes.RegisterInterfaces(ir)
	paramsproposal.RegisterInterfaces(ir)
	return codec.NewProtoCodec(ir)
}

// TestTxConfig returns a tx config for testing
func TestTxConfig() client.TxConfig {
	return tx.NewTxConfig(NewTestMarshaler(), tx.DefaultSignModes)
}

// CreateTestInput nolint
func CreateTestInput(t *testing.T) TestInput {
	keyAcc := sdk.NewKVStoreKey(authtypes.StoreKey)
	keyBank := sdk.NewKVStoreKey(banktypes.StoreKey)
	keyParams := sdk.NewKVStoreKey(paramstypes.StoreKey)
	tKeyParams := sdk.NewTransientStoreKey(paramstypes.TStoreKey)
	keyOracle := sdk.NewKVStoreKey(types.StoreKey)
	keyStaking := sdk.NewKVStoreKey(stakingtypes.StoreKey)
	keyDistr := sdk.NewKVStoreKey(distrtypes.StoreKey)

	cdc := NewTestMarshaler()
	amino := newTestCodec()
	db := dbm.NewMemDB()
	ms := store.NewCommitMultiStore(db)

	ctx := sdk.NewContext(ms, tmproto.Header{Time: time.Now().UTC()}, false, log.NewNopLogger())

	ms.MountStoreWithDB(keyAcc, sdk.StoreTypeIAVL, db)
	ms.MountStoreWithDB(keyBank, sdk.StoreTypeIAVL, db)
	ms.MountStoreWithDB(keyParams, sdk.StoreTypeIAVL, db)
	ms.MountStoreWithDB(tKeyParams, sdk.StoreTypeTransient, db)
	ms.MountStoreWithDB(keyOracle, sdk.StoreTypeIAVL, db)
	ms.MountStoreWithDB(keyStaking, sdk.StoreTypeIAVL, db)
	ms.MountStoreWithDB(keyDistr, sdk.StoreTypeIAVL, db)

	require.NoError(t, ms.LoadLatestVersion())

	blackListAddrs := map[string]bool{
		authtypes.FeeCollectorName:     true,
		stakingtypes.NotBondedPoolName: true,
		stakingtypes.BondedPoolName:    true,
		distrtypes.ModuleName:          true,
		types.ModuleName:               true,
	}

	maccPerms := map[string][]string{
		authtypes.FeeCollectorName:     nil,
		stakingtypes.NotBondedPoolName: {authtypes.Burner, authtypes.Staking},
		stakingtypes.BondedPoolName:    {authtypes.Burner, authtypes.Staking},
		distrtypes.ModuleName:          nil,
		types.ModuleName:               nil,
	}

	paramsKeeper := paramskeeper.NewKeeper(cdc, amino, keyParams, tKeyParams)
	paramsKeeper.Subspace(authtypes.ModuleName)
	paramsKeeper.Subspace(banktypes.ModuleName)
	paramsKeeper.Subspace(stakingtypes.ModuleName)
	paramsKeeper.Subspace(distrtypes.ModuleName)
	paramsKeeper.Subspace(types.ModuleName)

	authsub, ok := paramsKeeper.GetSubspace(authtypes.ModuleName)
	if !ok {
		panic("NOT OK")
	}
	banksub, ok := paramsKeeper.GetSubspace(banktypes.ModuleName)
	if !ok {
		panic("NOT OK")
	}
	stakingsub, ok := paramsKeeper.GetSubspace(stakingtypes.ModuleName)
	if !ok {
		panic("NOT OK")
	}
	distrsub, ok := paramsKeeper.GetSubspace(distrtypes.ModuleName)
	if !ok {
		panic("NOT OK")
	}
	oraclesub, ok := paramsKeeper.GetSubspace(types.ModuleName)
	if !ok {
		panic("NOT OK")
	}

	accountKeeper := authkeeper.NewAccountKeeper(cdc, keyAcc, authsub, authtypes.ProtoBaseAccount, maccPerms)
	bankKeeper := bankkeeper.NewBaseKeeper(cdc, keyBank, accountKeeper, banksub, blackListAddrs)

	ctx = ctx.WithMultiStore(ms)
	totalSupply := sdk.NewCoins(sdk.NewCoin(types.MicroLunaDenom, InitTokens.MulRaw(int64(len(Addrs)))))
	bankKeeper.SetSupply(ctx.WithMultiStore(ms), banktypes.NewSupply(totalSupply))

	stakingKeeper := stakingkeeper.NewKeeper(
		cdc,
		keyStaking,
		accountKeeper, bankKeeper, stakingsub,
	)

	distrKeeper := distrkeeper.NewKeeper(
		cdc,
		keyDistr, distrsub,
		accountKeeper, bankKeeper, stakingKeeper, authtypes.FeeCollectorName, blackListAddrs)

	distrKeeper.SetFeePool(ctx, distrtypes.InitialFeePool())
	distrParams := distrtypes.DefaultParams()
	distrParams.CommunityTax = sdk.NewDecWithPrec(2, 2)
	distrParams.BaseProposerReward = sdk.NewDecWithPrec(1, 2)
	distrParams.BonusProposerReward = sdk.NewDecWithPrec(4, 2)
	distrKeeper.SetParams(ctx, distrParams)

	feeCollectorAcc := authtypes.NewEmptyModuleAccount(authtypes.FeeCollectorName)
	notBondedPool := authtypes.NewEmptyModuleAccount(stakingtypes.NotBondedPoolName, authtypes.Burner, authtypes.Staking)
	bondPool := authtypes.NewEmptyModuleAccount(stakingtypes.BondedPoolName, authtypes.Burner, authtypes.Staking)
	distrAcc := authtypes.NewEmptyModuleAccount(distrtypes.ModuleName)
	oracleAcc := authtypes.NewEmptyModuleAccount(types.ModuleName, authtypes.Minter)

	bankKeeper.SetBalances(ctx, notBondedPool.GetAddress(), sdk.NewCoins(sdk.NewCoin(types.MicroLunaDenom, InitTokens.MulRaw(int64(len(Addrs))))))

	accountKeeper.SetModuleAccount(ctx, feeCollectorAcc)
	accountKeeper.SetModuleAccount(ctx, bondPool)
	accountKeeper.SetModuleAccount(ctx, notBondedPool)
	accountKeeper.SetModuleAccount(ctx, distrAcc)
	accountKeeper.SetModuleAccount(ctx, oracleAcc)

	genesis := stakingtypes.DefaultGenesisState()
	genesis.Params.BondDenom = types.MicroLunaDenom
	_ = staking.InitGenesis(ctx, stakingKeeper, accountKeeper, bankKeeper, genesis)

	for _, addr := range Addrs {
		require.NoError(t, bankKeeper.AddCoins(ctx, sdk.AccAddress(addr), InitCoins))
	}

	keeper := NewKeeper(cdc, keyOracle, oraclesub, distrKeeper, stakingKeeper, accountKeeper, bankKeeper, distrtypes.ModuleName)

	defaults := types.DefaultParams()
	keeper.SetParams(ctx, defaults)

	for _, denom := range defaults.Whitelist {
		keeper.SetTobinTax(ctx, denom.Denom, denom.Amount)
	}

	stakingKeeper.SetHooks(stakingtypes.NewMultiStakingHooks(distrKeeper.Hooks()))

	return TestInput{ctx, cdc, accountKeeper, bankKeeper, keeper, stakingKeeper, distrKeeper}
}

func NewTestMsgCreateValidator(address sdk.ValAddress, pubKey crypto.PubKey, amt sdk.Int) *stakingtypes.MsgCreateValidator {
	commission := stakingtypes.NewCommissionRates(sdk.ZeroDec(), sdk.ZeroDec(), sdk.ZeroDec())
	out, err := stakingtypes.NewMsgCreateValidator(
		address, pubKey, sdk.NewCoin(types.MicroLunaDenom, amt),
		stakingtypes.Description{}, commission, sdk.OneInt(),
	)
	if err != nil {
		panic(err)
	}
	return out
}
