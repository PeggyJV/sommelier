package keeper

import (
	"fmt"

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
		validatorCommitsMap = make(map[string]bool)
	)

	// iterate over the data votes
	// TODO: only iterate on the last voting period
	k.IterateAllocationCommits(ctx, func(validatorAddr sdk.ValAddress, commit types.Allocation) bool {
		// NOTE: the commit might have been submitted by a delegate, so we have to check the
		// original validator address

		// mark the validator voting as true
		validatorCommitsMap[validatorAddr.String()] = true

		// remove the miss counter in the current voting window for validators who have already voted
		k.DeleteMissCounter(ctx, validatorAddr)

		validator := k.stakingKeeper.Validator(ctx, sdk.ValAddress(validatorAddr))
		if validator == nil {
			// validator nor found, continue with next commit
			k.Logger(ctx).Debug(
				"validator not found for oracle commit tally",
				"validator-address", validatorAddr.String(),
			)
			return false
		}

		votedPower += validator.GetConsensusPower()

		// Safety check. Already validated on msg validation
		if len(commit.TickWeights.Weights) == 0 {
			k.Logger(ctx).Debug("attempted to process empty tick weights in commit", "validator", validatorAddr.String())
			return false
		}

		// save the oracle data for later processing
		for _, oracleData := range commit.Feed.Data {
			// oracleData, err := types.UnpackOracleData(oracleDataAny)
			// if err != nil {
			// 	// NOTE: this should never panic as the oracle data had already been checked before
			// 	// setting it to store.
			// 	panic(err)
			// }

			// add oracle data to maps
			feederVote := types.FeederVote{
				Data:    oracleData,
				Address: validatorAddr,
			}

			dataID := oracleData.GetID()

			oracleFeeds, ok := submittedFeedMap[dataID]
			if ok {
				submittedFeedMap[dataID] = append(oracleFeeds, feederVote)

				datas := submittedDataMap[dataID]
				submittedDataMap[dataID] = append(datas, oracleData)
			} else {
				submittedFeedMap[dataID] = []types.FeederVote{feederVote}
				submittedDataMap[dataID] = []types.OracleData{oracleData}
				submittedDataIDs = append(submittedDataIDs, dataID)
			}
		}

		// delete the commit as it has already been processed
		// TODO: consider keeping the votes for a few voting windows in order to
		// be able to submit evidence of inaccurate data feed
		k.DeleteAllocationCommit(ctx, validatorAddr)
		return false
	})

	// iterate over the full list of validators to increment miss counters if they didn't vote
	totalPower := int64(0)
	k.stakingKeeper.IterateBondedValidatorsByPower(ctx, func(_ int64, validator stakingtypes.ValidatorI) bool {
		totalPower += validator.GetConsensusPower()

		validatorAddr := validator.GetOperator()

		// only increment miss counter for validators who have previously submitted data
		if !validatorCommitsMap[validatorAddr.String()] && k.HasMissCounter(ctx, validatorAddr) {
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
	storeAverages := sdk.NewDec(votedPower).Quo(sdk.NewDec(totalPower)).GT(params.VoteThreshold)

	// Now, compute the aggregated data (eg: avg, median, etc.) for each type of data tracked by the oracle
	for _, id := range submittedDataIDs {
		oracleDatas, ok := submittedDataMap[id]
		if !ok {
			// no oracle data for the current id was submitted
			k.Logger(ctx).Debug(
				"no oracle data submitted for existing identifier",
				"id", id, "voting-period-start", fmt.Sprintf("%d", votePeriodStart),
			)
			continue
		}

		// aggregate the data  using the handler set up during app initialization
		aggregatedData, err := k.oracleHandler(ctx, oracleDatas)
		if err != nil {
			// TODO: ensure correctness or consider logging instead?
			panic(err)
		}

		if aggregatedData == nil {
			k.Logger(ctx).Debug(
				"aggregated data is nil",
				"id", id, "voting-period-start", fmt.Sprintf("%d", votePeriodStart),
			)
			continue
		}

		// once we have the aggregated data for the data type, we set it in the store

		// only store averages if there's enough voting power
		if storeAverages {
			k.SetAggregatedOracleData(ctx, ctx.BlockHeight(), aggregatedData)
			k.SetOracleDataHeight(ctx, aggregatedData.GetID(), ctx.BlockHeight())
		}

		// store the "average" for scoring validators later
		aggregateMap[aggregatedData.GetID()] = aggregatedData
	}

	for _, id := range submittedDataIDs {
		aggregateData, ok := aggregateMap[id]
		if !ok {
			k.Logger(ctx).Debug(
				"aggregated data for existing identifier",
				"id", id, "voting-period-start", fmt.Sprintf("%d", votePeriodStart),
			)
			continue
		}

		feederVotes, ok := submittedFeedMap[id]
		if !ok {
			// no oracle feed vote for the current id was submitted
			k.Logger(ctx).Debug(
				"no feed vote submitted for existing identifier",
				"id", id, "voting-period-start", fmt.Sprintf("%d", votePeriodStart),
			)
			continue
		}

		// Compare each validators vote for each data type against the
		// averages to define which are eligable for rewards
		for _, feederVote := range feederVotes {
			if !feederVote.Data.Compare(aggregateData, params.TargetThreshold) {
				continue
			}

			// TODO: reward validators / delegates
			// TODO: maybe consider splitting the delegate / validator reward in equal parts because the validator has the stake and the feeder
			// submits the data.
		}
	}

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
