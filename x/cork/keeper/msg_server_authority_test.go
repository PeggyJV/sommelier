package keeper

import (
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/require"

	// The SDK's bech32 prefix config is process-global and defaults to "cosmos";
	// this installs the "somm" prefixes the authority address is encoded with.
	_ "github.com/peggyjv/sommelier/v10/app/params"
	corktypes "github.com/peggyjv/sommelier/v10/x/cork/types"
	v2types "github.com/peggyjv/sommelier/v10/x/cork/types/v2"
)

const testCorkAuthority = "somm1lcsjy2d5s33h0sddd8lpuqvwyz5ruz7ju4aeqa"

// scheduleCorkFixture wires a keeper with an allowlisted cellar and returns a
// well-formed request scheduling a cork one block in the future.
func scheduleCorkFixture(t *testing.T, authority string) (Keeper, sdk.Context, *v2types.MsgScheduleCorkRequest, common.Address) {
	t.Helper()

	k, ctx, _, ctrl := setupCorkKeeper(t)
	t.Cleanup(ctrl.Finish)

	params := v2types.DefaultParams()
	params.CorkAuthority = authority
	k.SetParams(ctx, params)

	contract := common.HexToAddress("0x1111111111111111111111111111111111111111")
	k.SetCellarIDs(ctx, v2types.CellarIDSet{Ids: []string{contract.String()}})

	msg := &v2types.MsgScheduleCorkRequest{
		Cork: &v2types.Cork{
			TargetContractAddress: contract.String(),
			EncodedContractCall:   []byte{0xde, 0xad, 0xbe, 0xef},
		},
		BlockHeight: uint64(ctx.BlockHeight()) + 1,
		Signer:      testCorkAuthority,
	}

	return k, ctx, msg, contract
}

func TestScheduleCorkAcceptsAuthority(t *testing.T) {
	k, ctx, msg, contract := scheduleCorkFixture(t, testCorkAuthority)

	resp, err := k.ScheduleCork(sdk.WrapSDKContext(ctx), msg)
	require.NoError(t, err)
	require.NotEmpty(t, resp.Id)

	// The cork must land in the authority queue, not the legacy validator queue.
	var found int
	k.IterateAuthorityCorksByBlockHeight(ctx, msg.BlockHeight, func(_ uint64, _ []byte, gotContract common.Address, c v2types.Cork) bool {
		found++
		require.Equal(t, contract, gotContract)
		require.Equal(t, msg.Cork.EncodedContractCall, c.EncodedContractCall)
		return false
	})
	require.Equal(t, 1, found)

	// The typed legacy accessors are deleted; check the raw prefix directly.
	legacy := 0
	it := sdk.KVStorePrefixIterator(ctx.KVStore(k.storeKey), []byte{corktypes.ScheduledCorkKeyPrefix})
	for ; it.Valid(); it.Next() {
		legacy++
	}
	it.Close()
	require.Zero(t, legacy, "legacy validator queue must not be written")
}

func TestScheduleCorkRejectsNonAuthority(t *testing.T) {
	// Authority is configured, but the signer is a different address.
	k, ctx, msg, _ := scheduleCorkFixture(t, "somm1fcl08ymkl70dhyg3vmx4hjsqvxym7dawnp0zfp")

	_, err := k.ScheduleCork(sdk.WrapSDKContext(ctx), msg)
	require.Error(t, err)
	require.Contains(t, err.Error(), "is not the cork authority")
}

func TestScheduleCorkRejectsEmptyAuthority(t *testing.T) {
	// Fail-closed: an unset authority means nobody may schedule, including a
	// well-formed signer.
	k, ctx, msg, _ := scheduleCorkFixture(t, "")

	_, err := k.ScheduleCork(sdk.WrapSDKContext(ctx), msg)
	require.Error(t, err)
	require.Contains(t, err.Error(), "is not the cork authority")
}

func TestScheduleCorkRejectsUnmanagedCellar(t *testing.T) {
	k, ctx, msg, _ := scheduleCorkFixture(t, testCorkAuthority)

	// Retarget at a cellar that is not on the allowlist.
	msg.Cork.TargetContractAddress = common.HexToAddress("0x2222222222222222222222222222222222222222").String()

	_, err := k.ScheduleCork(sdk.WrapSDKContext(ctx), msg)
	require.ErrorIs(t, err, corktypes.ErrUnmanagedCellarAddress)
}

func TestScheduleCorkRejectsPastHeight(t *testing.T) {
	k, ctx, msg, _ := scheduleCorkFixture(t, testCorkAuthority)

	msg.BlockHeight = uint64(ctx.BlockHeight())

	_, err := k.ScheduleCork(sdk.WrapSDKContext(ctx), msg)
	require.ErrorIs(t, err, corktypes.ErrSchedulingInThePast)
}
