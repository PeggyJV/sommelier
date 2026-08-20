package keeper

import (
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/require"

	// Installs the chain's "somm" bech32 prefixes; the SDK default is "cosmos".
	_ "github.com/peggyjv/sommelier/v10/app/params"
	"github.com/peggyjv/sommelier/v10/x/axelarcork/types"
)

// stubPoaKeeper lets these tests drive the safe-mode gate directly.
type stubPoaKeeper struct{ active bool }

func (s stubPoaKeeper) SafeModeActive(sdk.Context) bool { return s.active }

const (
	testCorkAuthority = "somm1lcsjy2d5s33h0sddd8lpuqvwyz5ruz7ju4aeqa"
	testOtherSigner   = "somm1fcl08ymkl70dhyg3vmx4hjsqvxym7dawnp0zfp"
)

// authorityFixture wires an enabled module with a configured chain, an
// allowlisted cellar, and the given cork authority.
func authorityFixture(t *testing.T, authority string) (Keeper, sdk.Context, common.Address) {
	t.Helper()

	k, ctx, _, ctrl := setupCorkKeeper(t)
	t.Cleanup(ctrl.Finish)

	k.SetPoaKeeper(stubPoaKeeper{active: false})

	params := types.DefaultParams()
	params.Enabled = true
	params.CorkAuthority = authority
	k.SetParams(ctx, params)

	k.SetChainConfiguration(ctx, testChainArbitrum, types.ChainConfiguration{
		Name:         "arbitrum",
		Id:           testChainArbitrum,
		ProxyAddress: "0x9999999999999999999999999999999999999999",
	})

	contract := common.HexToAddress("0x1111111111111111111111111111111111111111")
	k.SetCellarIDs(ctx, testChainArbitrum, types.CellarIDSet{Ids: []string{contract.String()}})

	return k, ctx, contract
}

func TestScheduleCorkAcceptsAuthority(t *testing.T) {
	k, ctx, contract := authorityFixture(t, testCorkAuthority)

	msg := &types.MsgScheduleAxelarCorkRequest{
		Cork: &types.AxelarCork{
			TargetContractAddress: contract.String(),
			EncodedContractCall:   []byte{0xde, 0xad},
			ChainId:               testChainArbitrum,
		},
		ChainId:     testChainArbitrum,
		BlockHeight: uint64(ctx.BlockHeight()) + 1,
		Signer:      testCorkAuthority,
	}

	_, err := k.ScheduleCork(sdk.WrapSDKContext(ctx), msg)
	require.NoError(t, err)

	require.Equal(t, 1, countAt(k, ctx, testChainArbitrum, msg.BlockHeight),
		"cork must land in the authority queue")
	require.Empty(t, k.GetScheduledAxelarCorks(ctx, testChainArbitrum),
		"legacy validator queue must not be written")
}

// Every mutating handler must reject a signer that is not the cork authority,
// and must do so before performing any side effect. The reject path is asserted
// with no bank/transfer mocks registered: if a handler got as far as moving
// funds, the nil keeper would panic rather than return an error.
func TestAllHandlersRejectNonAuthority(t *testing.T) {
	cases := []struct {
		name string
		call func(k Keeper, ctx sdk.Context, contract common.Address) error
	}{
		{"ScheduleCork", func(k Keeper, ctx sdk.Context, contract common.Address) error {
			_, err := k.ScheduleCork(sdk.WrapSDKContext(ctx), &types.MsgScheduleAxelarCorkRequest{
				Cork: &types.AxelarCork{
					TargetContractAddress: contract.String(),
					EncodedContractCall:   []byte{0x01},
					ChainId:               testChainArbitrum,
				},
				ChainId:     testChainArbitrum,
				BlockHeight: uint64(ctx.BlockHeight()) + 1,
				Signer:      testOtherSigner,
			})
			return err
		}},
		{"RelayCork", func(k Keeper, ctx sdk.Context, contract common.Address) error {
			_, err := k.RelayCork(sdk.WrapSDKContext(ctx), &types.MsgRelayAxelarCorkRequest{
				TargetContractAddress: contract.String(),
				ChainId:               testChainArbitrum,
				Signer:                testOtherSigner,
			})
			return err
		}},
		{"RelayProxyUpgrade", func(k Keeper, ctx sdk.Context, _ common.Address) error {
			_, err := k.RelayProxyUpgrade(sdk.WrapSDKContext(ctx), &types.MsgRelayAxelarProxyUpgradeRequest{
				ChainId: testChainArbitrum,
				Signer:  testOtherSigner,
			})
			return err
		}},
		{"BumpCorkGas", func(k Keeper, ctx sdk.Context, _ common.Address) error {
			_, err := k.BumpCorkGas(sdk.WrapSDKContext(ctx), &types.MsgBumpAxelarCorkGasRequest{
				MessageId: "some-message-id",
				Signer:    testOtherSigner,
			})
			return err
		}},
		{"CancelScheduledCork", func(k Keeper, ctx sdk.Context, _ common.Address) error {
			_, err := k.CancelScheduledCork(sdk.WrapSDKContext(ctx), &types.MsgCancelAxelarCorkRequest{
				ChainId: testChainArbitrum,
				Signer:  testOtherSigner,
			})
			return err
		}},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			k, ctx, contract := authorityFixture(t, testCorkAuthority)
			err := tc.call(k, ctx, contract)
			require.Error(t, err)
			require.Contains(t, err.Error(), "is not the cork authority")
		})

		t.Run(tc.name+"/EmptyAuthorityFailsClosed", func(t *testing.T) {
			// Fail-closed: with no authority configured, nobody may act --
			// including an address that would otherwise be the authority.
			k, ctx, contract := authorityFixture(t, "")
			err := tc.call(k, ctx, contract)
			require.Error(t, err)
			require.Contains(t, err.Error(), "is not the cork authority")
		})
	}
}
