package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	"github.com/ethereum/go-ethereum/common"
	"github.com/peggyjv/gravity-bridge/module/v3/x/gravity/types"
	tmbytes "github.com/tendermint/tendermint/libs/bytes"
)

// StakingKeeper defines the expected staking keeper methods
type StakingKeeper interface {
	GetBondedValidatorsByPower(ctx sdk.Context) []stakingtypes.Validator
	GetLastValidatorPower(ctx sdk.Context, operator sdk.ValAddress) int64
	GetLastTotalPower(ctx sdk.Context) (power sdk.Int)
	IterateValidators(sdk.Context, func(index int64, validator stakingtypes.ValidatorI) (stop bool))
	IterateBondedValidatorsByPower(sdk.Context, func(index int64, validator stakingtypes.ValidatorI) (stop bool))
	IterateLastValidators(sdk.Context, func(index int64, validator stakingtypes.ValidatorI) (stop bool))
	Validator(sdk.Context, sdk.ValAddress) stakingtypes.ValidatorI
	ValidatorByConsAddr(sdk.Context, sdk.ConsAddress) stakingtypes.ValidatorI
	Slash(sdk.Context, sdk.ConsAddress, int64, int64, sdk.Dec)
	Jail(sdk.Context, sdk.ConsAddress)
	PowerReduction(ctx sdk.Context) sdk.Int
}

// GravityKeeper defines the expected gravity keeper methods
type GravityKeeper interface {
	SetOutgoingTx(ctx sdk.Context, outgoing types.OutgoingTx)
	CreateContractCallTx(
		ctx sdk.Context,
		invalidationNonce uint64,
		invalidationScope tmbytes.HexBytes,
		address common.Address,
		payload []byte,
		tokens []types.ERC20Token,
		fees []types.ERC20Token) *types.ContractCallTx
	GetOrchestratorValidatorAddress(ctx sdk.Context, orchAddr sdk.AccAddress) sdk.ValAddress
	GetValidatorEthereumAddress(ctx sdk.Context, valAddr sdk.ValAddress) common.Address
	GetEthereumOrchestratorAddress(ctx sdk.Context, ethAddr common.Address) sdk.AccAddress
	SetOrchestratorValidatorAddress(ctx sdk.Context, val sdk.ValAddress, orchAddr sdk.AccAddress)
}
