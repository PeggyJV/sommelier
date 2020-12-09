package simulation

import (
	"encoding/binary"
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"

	tmkv "github.com/cosmos/cosmos-sdk/types/kv"
	"github.com/tendermint/tendermint/crypto/ed25519"

	"github.com/cosmos/cosmos-sdk/codec"
	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	ccodec "github.com/cosmos/cosmos-sdk/crypto/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	distrtypes "github.com/cosmos/cosmos-sdk/x/distribution/types"
	paramsproposal "github.com/cosmos/cosmos-sdk/x/params/types/proposal"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	"github.com/peggyjv/sommelier/x/oracle/types"
)

var (
	delPk      = ed25519.GenPrivKey().PubKey()
	feederAddr = sdk.AccAddress(delPk.Address())
	valAddr    = sdk.ValAddress(delPk.Address())
)

func makeTestCodec() (cdc *codec.LegacyAmino) {
	cdc = codec.NewLegacyAmino()
	sdk.RegisterLegacyAminoCodec(cdc)
	ccodec.RegisterCrypto(cdc)
	types.RegisterLegacyAminoCodec(cdc)
	return
}

func newTestMarshaler() codec.ProtoCodecMarshaler {
	ir := codectypes.NewInterfaceRegistry()
	types.RegisterInterfaces(ir)
	authtypes.RegisterInterfaces(ir)
	sdk.RegisterInterfaces(ir)
	ccodec.RegisterInterfaces(ir)
	stakingtypes.RegisterInterfaces(ir)
	distrtypes.RegisterInterfaces(ir)
	paramsproposal.RegisterInterfaces(ir)
	return codec.NewProtoCodec(ir)
}

func TestDecodeDistributionStore(t *testing.T) {
	cdc := newTestMarshaler()

	exchangeRate := sdk.NewDecWithPrec(1234, 1)
	missCounter := 123

	aggregatePrevote := types.NewAggregateExchangeRatePrevote(types.AggregateVoteHash([]byte("12345")), valAddr, 123)
	aggregateVote := types.NewAggregateExchangeRateVote(sdk.DecCoins{
		{types.MicroKRWDenom, sdk.NewDecWithPrec(1234, 1)},
		{types.MicroKRWDenom, sdk.NewDecWithPrec(4321, 1)},
	}, valAddr)

	tobinTax := sdk.NewDecWithPrec(2, 2)
	missCounterBz := make([]byte, 8)
	binary.BigEndian.PutUint64(missCounterBz, uint64(missCounter))
	marEr, err := exchangeRate.Marshal()
	require.NoError(t, err)
	marTt, err := tobinTax.Marshal()
	require.NoError(t, err)

	kvPairs := []tmkv.Pair{
		{types.ExchangeRateKey, marEr},
		{types.FeederDelegationKey, []byte(feederAddr.String())},
		{types.MissCounterKey, missCounterBz},
		{types.AggregateExchangeRatePrevoteKey, cdc.MustMarshalBinaryLengthPrefixed(&aggregatePrevote)},
		{types.AggregateExchangeRateVoteKey, cdc.MustMarshalBinaryLengthPrefixed(&aggregateVote)},
		{types.TobinTaxKey, marTt},
		{[]byte{0x99}, []byte{0x99}},
	}

	tests := []struct {
		name        string
		expectedLog string
	}{
		{"ExchangeRate", fmt.Sprintf("%v\n%v", exchangeRate, exchangeRate)},
		{"FeederDelegation", fmt.Sprintf("%v\n%v", feederAddr, feederAddr)},
		{"MissCounter", fmt.Sprintf("%v\n%v", missCounter, missCounter)},
		{"AggregatePrevote", fmt.Sprintf("%v\n%v", aggregatePrevote, aggregatePrevote)},
		{"AggregateVote", fmt.Sprintf("%v\n%v", aggregateVote, aggregateVote)},
		{"TobinTax", fmt.Sprintf("%v\n%v", tobinTax, tobinTax)},
		{"other", ""},
	}

	for i, tt := range tests {
		i, tt := i, tt
		t.Run(tt.name, func(t *testing.T) {
			switch i {
			case len(tests) - 1:
				require.Panics(t, func() { DecodeStore(cdc)(kvPairs[i], kvPairs[i]) }, tt.name)
			default:
				require.Equal(t, tt.expectedLog, DecodeStore(cdc)(kvPairs[i], kvPairs[i]), tt.name)
			}
		})
	}
}
