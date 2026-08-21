package keeper

import (
	"context"

	sdkerrors "cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/peggyjv/sommelier/v10/x/poa/types"
)

type msgServer struct {
	Keeper
}

// NewMsgServerImpl returns a types.MsgServer implementation backed by the
// PoA Keeper. Both messages are gov-only.
func NewMsgServerImpl(k Keeper) types.MsgServer {
	return &msgServer{Keeper: k}
}

var _ types.MsgServer = msgServer{}

// UpdateAuthoritySet replaces the authority allowlist.
//
// Two gates beyond the gov-authority check:
//
//   - SAFE MODE. While the chain is in authority-empty safe mode, the only
//     bonded stake is community stake, so governance is decided entirely by the
//     very set the freeze exists to distrust. Allowing an authority-set
//     replacement there would let an attacker who induced the outage vote
//     itself the authority set, wait out the thaw delay, and own the bridge —
//     the exact escalation Option A claims to prevent. The supported recovery
//     from safe mode is for the EXISTING authority validators to unjail and
//     re-bond, which needs no governance at all. Permanent loss of the
//     authority set is a social-consensus event requiring a coordinated
//     restart, the same as halt mode.
//
//   - LIVENESS. At least one proposed validator must be bonded and unjailed,
//     otherwise the message would drop the chain into safe mode on the very
//     next EndBlocker (freezing gravity/cork/axelarcork) the moment it lands.
func (s msgServer) UpdateAuthoritySet(goCtx context.Context, msg *types.MsgUpdateAuthoritySet) (*types.MsgUpdateAuthoritySetResponse, error) {
	if msg.Authority != s.authority {
		return nil, sdkerrors.Wrapf(types.ErrUnauthorized, "expected %s, got %s", s.authority, msg.Authority)
	}
	if err := msg.ValidateBasic(); err != nil {
		return nil, err
	}
	ctx := sdk.UnwrapSDKContext(goCtx)
	if s.SafeModeActive(ctx) {
		return nil, types.ErrSafeModeGovFrozen
	}
	addrs := make([]sdk.ValAddress, len(msg.Validators))
	for i, v := range msg.Validators {
		addr, err := sdk.ValAddressFromBech32(v)
		if err != nil {
			return nil, err
		}
		addrs[i] = addr
	}
	if !s.anyBondedAndUnjailed(ctx, addrs) {
		return nil, sdkerrors.Wrap(types.ErrNoBondedAuthority,
			"proposed authority set contains no bonded, unjailed validator; applying it would enter safe mode immediately")
	}
	s.SetAuthoritySet(ctx, addrs)

	ctx.EventManager().EmitEvent(sdk.NewEvent(
		types.EventTypeAuthoritySetUpdated,
		sdk.NewAttribute("size", sdk.NewInt(int64(len(addrs))).String()),
	))
	return &types.MsgUpdateAuthoritySetResponse{}, nil
}

// anyBondedAndUnjailed reports whether at least one of `addrs` is a validator
// that currently counts toward authority power in the EndBlocker.
func (s msgServer) anyBondedAndUnjailed(ctx sdk.Context, addrs []sdk.ValAddress) bool {
	for _, a := range addrs {
		v, found := s.sk.GetValidator(ctx, a)
		if found && !v.Jailed && v.IsBonded() {
			return true
		}
	}
	return false
}

// UpdateParams updates the PoA params.
//
// Frozen in safe mode for the same reason as UpdateAuthoritySet: governance is
// community-only there, and `Enabled=false` clears safe mode IMMEDIATELY with
// no thaw delay (see EndBlocker), which would unfreeze gravity/cork/axelarcork
// under the untrusted set in a single block. HaltWhenAuthorityEmpty is likewise
// not something the distrusted set should be able to flip mid-incident.
func (s msgServer) UpdateParams(goCtx context.Context, msg *types.MsgUpdateParams) (*types.MsgUpdateParamsResponse, error) {
	if msg.Authority != s.authority {
		return nil, sdkerrors.Wrapf(types.ErrUnauthorized, "expected %s, got %s", s.authority, msg.Authority)
	}
	if err := msg.ValidateBasic(); err != nil {
		return nil, err
	}
	ctx := sdk.UnwrapSDKContext(goCtx)
	if s.SafeModeActive(ctx) {
		return nil, types.ErrSafeModeGovFrozen
	}
	s.SetParams(ctx, msg.Params)

	ctx.EventManager().EmitEvent(sdk.NewEvent(
		types.EventTypeParamsUpdated,
		sdk.NewAttribute("floor_fraction", msg.Params.FloorFraction.String()),
	))
	return &types.MsgUpdateParamsResponse{}, nil
}
