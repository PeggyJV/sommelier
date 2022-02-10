package keeper

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"strings"

	mapset "github.com/deckarep/golang-set"
	"github.com/ethereum/go-ethereum/common"

	"github.com/cosmos/cosmos-sdk/telemetry"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"

	"github.com/peggyjv/sommelier/v3/x/allocation/types"
)

var _ types.MsgServer = Keeper{}

func (k Keeper) signerToValAddr(ctx sdk.Context, signer sdk.AccAddress) (sdk.ValAddress, error) {
	validatorAddr := k.gravityKeeper.GetOrchestratorValidatorAddress(ctx, signer)
	if validatorAddr == nil {
		validator := k.stakingKeeper.Validator(ctx, sdk.ValAddress(signer))
		if validator == nil {
			return nil, sdkerrors.Wrap(stakingtypes.ErrNoValidatorFound, sdk.ValAddress(signer).String())
		}

		validatorAddr = validator.GetOperator()
		// NOTE: we set the validator address so we don't have to call look up for the validator
		// everytime a validator feeder submits oracle data
		k.gravityKeeper.SetOrchestratorValidatorAddress(ctx, validatorAddr, signer)
	}
	return validatorAddr, nil
}

// AllocationPrecommit implements types.MsgServer
func (k Keeper) AllocationPrecommit(c context.Context, msg *types.MsgAllocationPrecommit) (*types.MsgAllocationPrecommitResponse, error) {
	ctx := sdk.UnwrapSDKContext(c)

	signer := msg.MustGetSigner()
	validatorAddr, err := k.signerToValAddr(ctx, signer)
	if err != nil {
		return nil, err
	}

	// TODO: set precommit for current voting period
	hashList := make([]string, len(msg.Precommit))
	cellarList := make([]string, len(msg.Precommit))

	cellarSet := mapset.NewThreadUnsafeSet()
	for _, cellar := range k.GetCellars(ctx) {
		cellarSet.Add(cellar.Id)
	}
	for _, ap := range msg.Precommit {
		if !cellarSet.Contains(ap.CellarId) {
			return nil, fmt.Errorf("precommit for unknown cellar ID %s", ap.CellarId)
		}
		cellarList = append(cellarList, ap.CellarId)
		hashList = append(hashList, string(ap.Hash))
		k.SetAllocationPrecommit(ctx, validatorAddr, common.HexToAddress(ap.CellarId), *ap)
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
	val, err := k.signerToValAddr(ctx, signer)
	if err != nil {
		return nil, err
	}

	allocationEvents := sdk.Events{
		sdk.NewEvent(
			sdk.EventTypeMessage,
			sdk.NewAttribute(sdk.AttributeKeyModule, types.AttributeValueCategory),
			sdk.NewAttribute(sdk.AttributeKeyAction, types.EventTypeAllocationCommit),
			sdk.NewAttribute(sdk.AttributeKeySender, msg.Signer),
		),
	}

	cellarSet := mapset.NewThreadUnsafeSet()
	for _, cellar := range k.GetCellars(ctx) {
		cellarSet.Add(cellar.Id)
	}

	// check if there's an existing vote for the current voting period start
	if k.HasAllocationCommit(ctx, val) {
		return nil, sdkerrors.Wrap(types.ErrAlreadyCommitted, fmt.Sprintf("validator: %s", val.String()))
	}

	for _, commit := range msg.Commit {
		cel := common.HexToAddress(commit.Vote.Cellar.Id)

		if cellarSet.Contains(commit.Vote.Cellar.Id) {
			cellarSet.Remove(commit.Vote.Cellar.Id)
		} else {
			return nil, fmt.Errorf("commit for unknown cellar: %s", commit.Vote.Cellar.Id)
		}

		// Get the precommit for that validator from the store
		precommit, found := k.GetAllocationPrecommit(ctx, val, cel)
		// check that there is a precommit
		if !found || len(precommit.Hash) == 0 {
			return nil, sdkerrors.Wrap(types.ErrNoPrecommit, val.String())
		}

		// calculate the vote hash on the server
		commitHash, err := commit.Vote.Hash(commit.Salt, val)
		if err != nil {
			return nil, sdkerrors.Wrap(sdkerrors.ErrJSONMarshal, "failed to hash cellar allocation")
		}

		// compare to precommit hash
		cellarJSON, err := json.Marshal(commit.Vote)
		if err != nil {
			return nil, err
		}
		if !bytes.Equal(commitHash, precommit.Hash) {
			k.Logger(ctx).Error(
				"error with hash",
				"msg", msg.String(),
				"commit", cellarJSON,
				"signer", val.String(),
				"salt", commit.Salt,
				"precommit hash", string(precommit.Hash),
			)
			return nil, sdkerrors.Wrapf(
				types.ErrHashMismatch,
				"precommit %x ≠ commit %x. cellar json: %s, signer val %s", precommit.Hash, commitHash, string(cellarJSON), val,
			)
		}

		allocationEvents = append(
			allocationEvents,
			sdk.NewEvent(
				types.EventTypeAllocationCommit,
				sdk.NewAttribute(types.AttributeKeyCellar, cel.String()),
				sdk.NewAttribute(types.AttributeKeyValidator, val.String()),
			),
		)

		// set the vote in the store
		k.SetAllocationCommit(ctx, val, cel, *commit)
	}

	if cellarSet.Cardinality() > 0 {
		return nil, fmt.Errorf("commits not included for cellars: %s", cellarSet.String())
	}

	// TODO: set data for the current voting period
	ctx.EventManager().EmitEvents(allocationEvents)

	return &types.MsgAllocationCommitResponse{}, nil
}
