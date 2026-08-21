package keeper

import (
	"context"

	sdkerrors "cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/errors"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"

	"github.com/peggyjv/sommelier/v10/x/poa/types"
)

var _ types.QueryServer = Querier{}

// Querier implements the gRPC query service.
type Querier struct {
	Keeper
}

// NewQuerier returns a Querier wrapping the keeper.
func NewQuerier(k Keeper) Querier { return Querier{k} }

func (q Querier) Params(c context.Context, _ *types.QueryParamsRequest) (*types.QueryParamsResponse, error) {
	ctx := sdk.UnwrapSDKContext(c)
	return &types.QueryParamsResponse{Params: q.GetParams(ctx)}, nil
}

func (q Querier) AuthoritySet(c context.Context, _ *types.QueryAuthoritySetRequest) (*types.QueryAuthoritySetResponse, error) {
	ctx := sdk.UnwrapSDKContext(c)
	addrs := q.GetAuthoritySet(ctx)
	out := make([]string, len(addrs))
	for i, a := range addrs {
		out[i] = a.String()
	}
	return &types.QueryAuthoritySetResponse{Validators: out}, nil
}

func (q Querier) EffectivePower(c context.Context, req *types.QueryEffectivePowerRequest) (*types.QueryEffectivePowerResponse, error) {
	if req == nil || req.OperatorAddress == "" {
		return nil, status.Error(codes.InvalidArgument, "operator_address required")
	}
	op, err := sdk.ValAddressFromBech32(req.OperatorAddress)
	if err != nil {
		return nil, sdkerrors.Wrap(errors.ErrInvalidAddress, req.OperatorAddress)
	}
	ctx := sdk.UnwrapSDKContext(c)
	return &types.QueryEffectivePowerResponse{
		Power:       q.sk.GetLastValidatorPower(ctx, op),
		IsAuthority: q.IsAuthority(ctx, op),
	}, nil
}

// SafeMode reports the authority-empty freeze state. This is the module's
// primary operational health check: while active, gravity / cork / axelarcork
// are frozen, and previously the only way to observe that was to scrape events
// or node logs.
func (q Querier) SafeMode(c context.Context, _ *types.QuerySafeModeRequest) (*types.QuerySafeModeResponse, error) {
	ctx := sdk.UnwrapSDKContext(c)

	var bonded uint64
	for _, op := range q.GetAuthoritySet(ctx) {
		v, found := q.sk.GetValidator(ctx, op)
		if found && !v.Jailed && v.IsBonded() {
			bonded++
		}
	}

	thaw, _ := q.getThawHeight(ctx)
	return &types.QuerySafeModeResponse{
		Active:               q.SafeModeActive(ctx),
		ThawHeight:           thaw,
		BondedAuthorityCount: bonded,
	}, nil
}
