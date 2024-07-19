package types

//go:generate mockgen --source=x/cork/types/expected_keepers.go --destination=x/cork/testutil/expected_keepers_mocks.go --package=mock_types

import (
	"cosmossdk.io/math"

	tmbytes "github.com/cometbft/cometbft/libs/bytes"
	sdk "github.com/cosmos/cosmos-sdk/types"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	"github.com/ethereum/go-ethereum/common"
	gravitytypes "github.com/peggyjv/gravity-bridge/module/v4/x/gravity/types"
	pubsubtypes "github.com/peggyjv/sommelier/v7/x/pubsub/types"
)

// StakingKeeper defines the expected staking keeper methods
type StakingKeeper interface {
	GetBondedValidatorsByPower(ctx sdk.Context) []stakingtypes.Validator
	GetLastValidatorPower(ctx sdk.Context, operator sdk.ValAddress) int64
	GetLastTotalPower(ctx sdk.Context) (power math.Int)
	IterateValidators(sdk.Context, func(index int64, validator stakingtypes.ValidatorI) (stop bool))
	IterateBondedValidatorsByPower(sdk.Context, func(index int64, validator stakingtypes.ValidatorI) (stop bool))
	IterateLastValidators(sdk.Context, func(index int64, validator stakingtypes.ValidatorI) (stop bool))
	Validator(sdk.Context, sdk.ValAddress) stakingtypes.ValidatorI
	ValidatorByConsAddr(sdk.Context, sdk.ConsAddress) stakingtypes.ValidatorI
	Slash(sdk.Context, sdk.ConsAddress, int64, int64, sdk.Dec) math.Int
	Jail(sdk.Context, sdk.ConsAddress)
	PowerReduction(ctx sdk.Context) math.Int
}

// GravityKeeper defines the expected gravity keeper methods
type GravityKeeper interface {
	SetOutgoingTx(ctx sdk.Context, outgoing gravitytypes.OutgoingTx)
	CreateContractCallTx(
		ctx sdk.Context,
		invalidationNonce uint64,
		invalidationScope tmbytes.HexBytes,
		address common.Address,
		payload []byte,
		tokens []gravitytypes.ERC20Token,
		fees []gravitytypes.ERC20Token) *gravitytypes.ContractCallTx
	GetOrchestratorValidatorAddress(ctx sdk.Context, orchAddr sdk.AccAddress) sdk.ValAddress
	GetValidatorEthereumAddress(ctx sdk.Context, valAddr sdk.ValAddress) common.Address
	GetEthereumOrchestratorAddress(ctx sdk.Context, ethAddr common.Address) sdk.AccAddress
	SetOrchestratorValidatorAddress(ctx sdk.Context, val sdk.ValAddress, orchAddr sdk.AccAddress)
}

type PubsubKeeper interface {
	GetPublisher(ctx sdk.Context, publisherDomain string) (publisher pubsubtypes.Publisher, found bool)
	SetDefaultSubscription(ctx sdk.Context, defaultSubscription pubsubtypes.DefaultSubscription)
	DeleteDefaultSubscription(ctx sdk.Context, subscriptionID string)
}
