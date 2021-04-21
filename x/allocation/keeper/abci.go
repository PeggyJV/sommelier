package keeper

import (
	"fmt"
	"github.com/ethereum/go-ethereum/common"

	"github.com/peggyjv/sommelier/x/allocation/types"

	sdk "github.com/cosmos/cosmos-sdk/types"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
)

// BeginBlocker is called at the beginning of every block
func (k Keeper) BeginBlocker(ctx sdk.Context) {
	// if there is not a vote period set, initialize it with current block height
	// TODO: consider removing
	if !k.HasCommitPeriodStart(ctx) {
		k.SetCommitPeriodStart(ctx, ctx.BlockHeight())
	}

	// On begin block, if we are tallying, emit the new vote period data
	params := k.GetParamSet(ctx)
	votePeriodStart, found := k.GetCommitPeriodStart(ctx)
	if !found {
		panic("vote period not set")
	}

	if (ctx.BlockHeight() - votePeriodStart) < params.VotePeriod {
		return
	}

	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			types.EventTypeCommitPeriod,
			sdk.NewAttribute(types.AttributeKeyCommitPeriodStart, fmt.Sprintf("%d", votePeriodStart)),
			sdk.NewAttribute(types.AttributeKeyCommitPeriodEnd, fmt.Sprintf("%d", votePeriodStart+params.VotePeriod)),
		),
	)
}

// EndBlocker defines the oracle logic that executes at the end of every block:
//
// 0) Checks if the voting period is over and performs a no-op if it's not.
//
// 1) Checks all the votes submitted in the last period and groups them by ID.
// This step also deletes the oracle votes.
//
// 2) Increments the miss counter for validators that didn't vote
//
// 3) Aggregates the data by ID and type
//
// 4) Compares each submitted data with the aggregated result in order to check
// for reward eligibility
//
// 5) Slashes validators that haven't voted in a while
//
// 6) Deletes all prevotes
//
// 7) Sets the new voting period to the next block

type PowerWeight struct {
	validator sdk.ValAddress
	cellar    common.Address
	fee_level sdk.Dec
	power     int64
	tick      uint32
}

func (k Keeper) EndBlocker(ctx sdk.Context) {
	params := k.GetParamSet(ctx)
	votePeriodStart, found := k.GetCommitPeriodStart(ctx)
	if !found {
		panic("vote period start not set")
	}

	// if the vote period has ended, tally the votes
	periodEnded := (ctx.BlockHeight() - votePeriodStart) >= params.VotePeriod
	if !periodEnded {
		return
	}

	k.Logger(ctx).Info("tallying oracle votes", "height", fmt.Sprintf("%d", ctx.BlockHeight()))

	votedPower := int64(0)

	var (
		// validators who submitted their vote
		validatorsCommittedMap = make(map[string]bool)
		// cellarMap is a map of cellars by address
		cellarsMap = make(map[common.Address]*types.Cellar)
		// cellarCommitsMap is a list of commits by cellar
		cellarCommitsMap = make(map[common.Address][]types.Allocation)
		// cellarTickWeightPowerMap is a map of cellars to pools to ticks with power adjusted allocations
		cellarPoolTickPowerMap = make(map[common.Address]map[sdk.Dec]map[uint32]sdk.Dec)
	)

	for _, cellar := range params.Cellars {
		cellarsMap[common.HexToAddress(cellar.CellarId)] = cellar
	}

	// iterate over the data votes
	// TODO: only iterate on the last voting period
	k.IterateAllocationCommitValidators(ctx, func(validatorAddr sdk.ValAddress) bool {
		// NOTE: the commit might have been submitted by a delegate, so we have to check the
		// original validator address

		validator := k.stakingKeeper.Validator(ctx, sdk.ValAddress(validatorAddr))
		if validator == nil {
			// validator nor found, continue with next commit
			k.Logger(ctx).Debug(
				"validator not found for oracle commit tally",
				"validator-address", validatorAddr.String(),
			)
			return false
		}

		validatorPower := validator.GetConsensusPower()

		k.IterateValidatorAllocationCommits(ctx, validatorAddr, func(cellar common.Address, commit types.Allocation) bool {
			cellarCommitsMap[cellar] = append(cellarCommitsMap[cellar], commit)
			for _, pa := range commit.PoolAllocations.Allocations {
				for _, tw := range pa.TickWeights.Weights {
					cellarPoolTickPowerMap[cellar][pa.FeeLevel][tw.Tick] =
						cellarPoolTickPowerMap[cellar][pa.FeeLevel][tw.Tick].Add(tw.Weight.MulInt64(validatorPower))
				}
			}

			// delete the commit as it has already been processed
			// TODO: consider keeping the votes for a few voting windows in order to
			// be able to submit evidence of inaccurate data feed
			k.DeleteAllocationCommit(ctx, validatorAddr, cellar)

			return false
		})

		// mark the validator voting as true
		validatorsCommittedMap[validatorAddr.String()] = true

		// remove the miss counter in the current voting window for validators who have already voted
		k.DeleteMissCounter(ctx, validatorAddr)

		votedPower += validatorPower

		return false
	})

	// iterate over the full list of validators to increment miss counters if they didn't vote
	totalPower := int64(0)
	k.stakingKeeper.IterateBondedValidatorsByPower(ctx, func(_ int64, validator stakingtypes.ValidatorI) bool {
		totalPower += validator.GetConsensusPower()

		validatorAddr := validator.GetOperator()

		// only increment miss counter for validators who have previously submitted data
		if !validatorsCommittedMap[validatorAddr.String()] && k.HasMissCounter(ctx, validatorAddr) {
			// TODO: this is wrong because the feeder could submit a single oracle data uniswap pair instead of all the required
			// ones and still be counted as if they voted correctly. Maybe consider adding each uniswap pair id to the params?
			// TODO: we need to define what is the exact data that we want the validators to submit.
			k.IncrementMissCounter(ctx, validatorAddr)
		}

		return false
	})

	// if the voted_power/total_power > params.VoteThreshold then we store the averages in the store
	// TODO: check for sample size too as there could be only a few feeders with a large voting power that submitted the
	// data to the oracle.
	// storeAverages := sdk.NewDec(votedPower).Quo(sdk.NewDec(totalPower)).GT(params.VoteThreshold)

	// slash validators who have missed to many votes
	k.IterateMissCounters(ctx, func(validatorAddr sdk.ValAddress, counter int64) bool {
		missedVotesPerWindow := sdk.NewDec(counter).Quo(sdk.NewDec(params.SlashWindow))
		if params.MinValidPerWindow.GTE(missedVotesPerWindow) {
			// continue with next counter
			return false
		}

		// slash validator below the minimum vote threshold
		validator := k.stakingKeeper.Validator(ctx, validatorAddr)
		if validator == nil {
			// validator not found, continue with the next counter
			k.Logger(ctx).Debug(
				"validator not found for miss counter",
				"validator-address", validatorAddr.String(),
			)

			// delete miss counter as the validator is not registered to the staking store
			// TODO: maybe use defer? not sure if deleting while iterating is bad or not?
			k.DeleteMissCounter(ctx, validatorAddr)
			return false
		}

		consensusAddr, err := validator.GetConsAddr()
		if err != nil {
			panic(err)
		}

		// the validator cannot be unbonded, otherwise slashing will panic
		if validator.IsUnbonded() {
			k.Logger(ctx).Debug("validator not slashed due to being unbonded", "address", consensusAddr.String())
			return false
		}

		k.stakingKeeper.Slash(ctx, consensusAddr, ctx.BlockHeight(), validator.GetConsensusPower(), params.SlashFraction)
		return false
	})

	// NOTE: the reward amount should be less than the slashed amount

	// TODO: Setup module account for oracle module.

	// Reset state prior to next round
	k.DeleteAllPrecommits(ctx)

	// After the tallying is done, reset the vote period start to the next block
	votePeriodStart = ctx.BlockHeight() + 1
	k.SetCommitPeriodStart(ctx, votePeriodStart)

	k.Logger(ctx).Info("vote period set", "height", fmt.Sprintf("%d", votePeriodStart))
}
