package keeper

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/ethereum/go-ethereum/common"
	"strings"

	"github.com/armon/go-metrics"
	"github.com/cosmos/cosmos-sdk/telemetry"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"

	"github.com/peggyjv/sommelier/x/allocation/types"
)

var _ types.MsgServer = Keeper{}

// DelegateAllocations implements types.MsgServer
func (k Keeper) DelegateAllocations(c context.Context, msg *types.MsgDelegateAllocations) (*types.MsgDelegateAllocationsResponse, error) {
	ctx := sdk.UnwrapSDKContext(c)

	val, del := msg.MustGetValidator(), msg.MustGetDelegate()

	// check that the signer is a bonded validator and is not jailed
	validator := k.stakingKeeper.Validator(ctx, sdk.ValAddress(val))
	if validator == nil {
		return nil, sdkerrors.Wrap(stakingtypes.ErrNoValidatorFound, sdk.ValAddress(val).String())
	}

	if validator.IsUnbonded() {
		return nil, sdkerrors.Wrapf(sdkerrors.ErrUnauthorized, "validator %s cannot be unbonded", validator.GetOperator())
	}

	if validator.IsJailed() {
		return nil, sdkerrors.Wrap(stakingtypes.ErrValidatorJailed, validator.GetOperator().String())
	}

	// check that the delegate feeder is not a validator, this prevents mirroring and freeloading
	// See https://medium.com/fabric-ventures/decentralised-oracles-a-comprehensive-overview-d3168b9a8841
	if k.stakingKeeper.Validator(ctx, sdk.ValAddress(del)) != nil {
		return nil, sdkerrors.Wrapf(sdkerrors.ErrUnauthorized, "feeder delegate %s cannot be a validator", del)
	}

	k.SetValidatorDelegateAddress(ctx, del, sdk.ValAddress(val))

	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			sdk.EventTypeMessage,
			sdk.NewAttribute(sdk.AttributeKeyModule, types.AttributeValueCategory),
			sdk.NewAttribute(sdk.AttributeKeyAction, types.EventTypeDelegateAllocations),
			sdk.NewAttribute(sdk.AttributeKeySender, msg.Validator),
			sdk.NewAttribute(types.AttributeKeyValidator, msg.Validator),
			sdk.NewAttribute(types.AttributeKeyDeleagate, msg.Delegate),
		),
	)

	return &types.MsgDelegateAllocationsResponse{}, nil
}

// AllocationPrecommit implements types.MsgServer
func (k Keeper) AllocationPrecommit(c context.Context, msg *types.MsgAllocationPrecommit) (*types.MsgAllocationPrecommitResponse, error) {
	ctx := sdk.UnwrapSDKContext(c)

	signer := msg.MustGetSigner()
	validatorAddr := k.GetValidatorAddressFromDelegate(ctx, signer)
	if validatorAddr == nil {
		validator := k.stakingKeeper.Validator(ctx, sdk.ValAddress(signer))
		if validator == nil {
			return nil, sdkerrors.Wrap(stakingtypes.ErrNoValidatorFound, sdk.ValAddress(signer).String())
		}

		validatorAddr = validator.GetOperator()
		// NOTE: we set the validator address so we don't have to call look up for the validator
		// everytime the a validator feeder submits oracle data
		k.SetValidatorDelegateAddress(ctx, signer, validatorAddr)
	}

	// TODO: set precommit for current voting period
	var hashList []string
	var cellarList []string
	for _, ap := range msg.Precommit {
		cellarList = append(cellarList, ap.CellarId)
		hashList = append(hashList, ap.Hash.String())
		k.SetAllocationPrecommit(ctx, validatorAddr, common.HexToAddress(ap.CellarId), *ap)
	}
	// set miss counter now that the validator committed the provote
	if !k.HasMissCounter(ctx, validatorAddr) {
		k.SetMissCounter(ctx, validatorAddr, 0)
	}

	ctx.EventManager().EmitEvents(
		sdk.Events{
			sdk.NewEvent(
				sdk.EventTypeMessage,
				sdk.NewAttribute(sdk.AttributeKeyModule, types.AttributeValueCategory),
			),
			sdk.NewEvent(
				types.EventTypeAllocationPrecommit,
				sdk.NewAttribute(types.AttributeKeySigner, signer.String()),
				sdk.NewAttribute(types.AttributeKeyValidator, validatorAddr.String()),
				sdk.NewAttribute(types.AttributeKeyCellar, strings.Join(cellarList, ",")),
				sdk.NewAttribute(types.AttributeKeyPrevoteHash, strings.Join(hashList, ",")),
			),
		},
	)

	defer func() {
		telemetry.IncrCounter(1, types.ModuleName, "prevote")
	}()

	return &types.MsgAllocationPrecommitResponse{}, nil
}

// AllocationCommit implements types.MsgServer
func (k Keeper) AllocationCommit(c context.Context, msg *types.MsgAllocationCommit) (*types.MsgAllocationCommitResponse, error) {
	ctx := sdk.UnwrapSDKContext(c)

	// Make sure that the message was properly signed
	signer := msg.MustGetSigner()
	val := k.GetValidatorAddressFromDelegate(ctx, signer)
	if val == nil {
		validator := k.stakingKeeper.Validator(ctx, sdk.ValAddress(signer))
		if validator == nil {
			return nil, sdkerrors.Wrap(stakingtypes.ErrNoValidatorFound, sdk.ValAddress(signer).String())
		}

		val = validator.GetOperator()
		// NOTE: we set the validator address so we don't have to call look up for the validator
		// everytime the a validator feeder submits oracle data
		k.SetValidatorDelegateAddress(ctx, signer, val)
	}

	// check that the validator is still bonded and is not jailed

	allocationEvents := sdk.Events{
		sdk.NewEvent(
			sdk.EventTypeMessage,
			sdk.NewAttribute(sdk.AttributeKeyModule, types.AttributeValueCategory),
			sdk.NewAttribute(sdk.AttributeKeyAction, types.EventTypeAllocationCommit),
			sdk.NewAttribute(sdk.AttributeKeySender, msg.Signer),
		),
	}

	for _, commit := range msg.Commit {
		cel := common.HexToAddress(commit.CellarId)

		// check if there's an existing vote for the current voting period start
		if k.HasAllocationCommit(ctx, val, cel) {
			return nil, sdkerrors.Wrap(types.ErrAlreadyCommitted, fmt.Sprint("validator: %s cellar: %s", val.String(), cel.String()))
		}

		// Get the precommit for that validator from the store
		precommit, found := k.GetAllocationPrecommit(ctx, val, cel)
		// check that there is a precommit
		if !found || len(precommit.Hash) == 0 {
			return nil, sdkerrors.Wrap(types.ErrNoPrecommit, val.String())
		}

		//allowedTypesMap := make(map[string]bool)
		//allowedDataTypes := k.GetParamSet(ctx).DataTypes
		//
		//for _, allowedDataType := range allowedDataTypes {
		//	allowedTypesMap[allowedDataType] = true
		//}


		// parse data to json in order to compute the vote hash and sort
		jsonBz, err := json.Marshal(commit.TickWeights)
		if err != nil {
			return nil, sdkerrors.Wrap(
				sdkerrors.ErrJSONMarshal, "failed to marshal json allocation tick weights",
			)
		}

		jsonBz = sdk.MustSortJSON(jsonBz)

		// calculate the vote hash on the server
		commitHash := types.DataHash(commit.Salt, string(jsonBz), val)

		// compare to precommit hash
		if !bytes.Equal(commitHash, precommit.Hash) {
			return nil, sdkerrors.Wrapf(
				types.ErrHashMismatch,
				"precommit %x ≠ commit %x", precommit.Hash, commitHash,
			)
		}


		for _, tickWeight := range commit.TickWeights.Weights {
			// unpack the oracle data one by one
			// tickWeight, err := types.UnpackOracleData(oracleDataAny)
			// if err != nil {
			// 	return nil, sdkerrors.Wrapf(types.ErrUnpackOracleData, "index %d", i)
			// }

			// if !allowedTypesMap[tickWeight.Type()] {
			// 	return nil, sdkerrors.Wrapf(
			// 		types.ErrUnsupportedDataType,
			// 		"%s, allowed %v", tickWeight.Type(), allowedDataTypes,
			// 	)
			// }

			if

			if !k.HasOracleDataType(ctx, tickWeight.GetID()) {
				k.SetOracleDataType(ctx, tickWeight.Type(), tickWeight.GetID())
			}

			allocationEvents = append(
				allocationEvents,
				sdk.NewEvent(
					types.EventTypeAllocationCommit,
					sdk.NewAttribute(types.AttributeKeyOracleDataID, tickWeight.GetID()),
					sdk.NewAttribute(types.AttributeKeyOracleDataType, tickWeight.Type()),
				),
			)

			labels := []metrics.Label{
				telemetry.NewLabel("data-type", tickWeight.Type()),
				telemetry.NewLabel("data-id", tickWeight.GetID()),
			}

			defer func() {
				telemetry.IncrCounterWithLabels(
					[]string{types.ModuleName, "feed"},
					1,
					labels,
				)
			}()
		}
	}




	// set the vote in the store
	// TODO: set data for the current voting period
	k.SetAllocationCommit(ctx, val, *msg.Vote)
	ctx.EventManager().EmitEvents(allocationEvents)

	defer func() {
		telemetry.IncrCounterWithLabels(
			[]string{types.ModuleName, "vote"},
			1,
			[]metrics.Label{
				telemetry.NewLabel(types.AttributeKeyValidator, val.String()),
			},
		)
	}()

	return &types.MsgAllocationCommitResponse{}, nil
}
