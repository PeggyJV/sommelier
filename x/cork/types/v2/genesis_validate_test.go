package v2

import (
	"testing"

	"github.com/stretchr/testify/require"

	_ "github.com/peggyjv/sommelier/v10/app/params"
)

// An exported genesis containing a queued authority cork must validate.
//
// Authority corks have no scheduling validator, so ScheduledCork.Validator is
// empty. If ValidateBasic still demands a bech32 validator, every exported
// genesis with a queued cork fails validate-genesis and cannot be imported.
func TestScheduledCorkValidateBasicAcceptsAuthorityCork(t *testing.T) {
	cork := Cork{
		TargetContractAddress: "0x1111111111111111111111111111111111111111",
		EncodedContractCall:   []byte{0x01, 0x02},
	}
	height := uint64(42)

	sc := ScheduledCork{
		Cork:        &cork,
		BlockHeight: height,
		Id:          cork.IDHash(height), // 32-byte keccak
		Validator:   "",                  // authority corks have no validator
	}
	require.NoError(t, sc.ValidateBasic())
}

// A genesis entry with a null cork must be rejected, not panic. An absent
// embedded message decodes to a nil pointer, so this is reachable from a
// malformed genesis file or a decoded MsgScheduleCorkRequest -- and both
// modules' InitGenesis dereference Cork directly afterwards.
func TestScheduledCorkValidateBasicRejectsNilCork(t *testing.T) {
	sc := ScheduledCork{
		Cork:        nil,
		BlockHeight: 42,
		Id:          make([]byte, 32),
	}
	require.NotPanics(t, func() {
		require.Error(t, sc.ValidateBasic())
	})
}
