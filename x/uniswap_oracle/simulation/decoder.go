package simulation

import (
	"bytes"
	"encoding/binary"
	"fmt"

	"github.com/cosmos/cosmos-sdk/types/kv"

	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/peggyjv/sommelier/x/uniswap_oracle/types"
)

// DecodeStore unmarshals the KVPair's Value to the corresponding distribution type
func DecodeStore(cdc codec.Marshaler) func(kvA, kvB kv.Pair) string {
	return func(kvA, kvB kv.Pair) string {
		switch {
		case bytes.Equal(kvA.Key[:1], types.ExchangeRateKey):
			var exchangeRateA, exchangeRateB sdk.Dec
			_ = exchangeRateA.Unmarshal(kvA.Value)
			_ = exchangeRateB.Unmarshal(kvB.Value)
			return fmt.Sprintf("%v\n%v", exchangeRateA, exchangeRateB)
		case bytes.Equal(kvA.Key[:1], types.FeederDelegationKey):
			addressA, _ := sdk.AccAddressFromBech32(string(kvA.Value))
			addressB, _ := sdk.AccAddressFromBech32(string(kvB.Value))
			return fmt.Sprintf("%v\n%v", addressA, addressB)
		case bytes.Equal(kvA.Key[:1], types.MissCounterKey):
			counterA := int64(binary.BigEndian.Uint64(kvA.Value))
			counterB := int64(binary.BigEndian.Uint64(kvB.Value))
			return fmt.Sprintf("%v\n%v", counterA, counterB)
		case bytes.Equal(kvA.Key[:1], types.AggregateExchangeRatePrevoteKey):
			var prevoteA, prevoteB types.AggregateExchangeRatePrevote
			cdc.MustUnmarshalBinaryBare(kvA.Value, &prevoteA)
			cdc.MustUnmarshalBinaryBare(kvB.Value, &prevoteB)
			return fmt.Sprintf("%v\n%v", prevoteA, prevoteB)
		case bytes.Equal(kvA.Key[:1], types.AggregateExchangeRateVoteKey):
			var voteA, voteB types.AggregateExchangeRateVote
			cdc.MustUnmarshalBinaryBare(kvA.Value, &voteA)
			cdc.MustUnmarshalBinaryBare(kvB.Value, &voteB)
			return fmt.Sprintf("%v\n%v", voteA, voteB)
		case bytes.Equal(kvA.Key[:1], types.TobinTaxKey):
			var tobinTaxA, tobinTaxB sdk.Dec
			_ = tobinTaxA.Unmarshal(kvA.Value)
			_ = tobinTaxB.Unmarshal(kvB.Value)
			return fmt.Sprintf("%v\n%v", tobinTaxA, tobinTaxB)
		default:
			panic(fmt.Sprintf("invalid oracle key prefix %X", kvA.Key[:1]))
		}
	}
}
