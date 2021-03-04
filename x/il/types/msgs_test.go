package types

import (
	"testing"

	"github.com/cosmos/cosmos-sdk/testutil/testdata"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"
)

func TestMsgCreateStoplossValidate(t *testing.T) {
	_, _, addr1 := testdata.KeyTestPubAddr()

	stoploss := &Stoploss{
		UniswapPairID:       "0x3041cbd36888becc7bbcbc0045e3b1f144466f5f",
		LiquidityPoolShares: 10,
		MaxSlippage:         sdk.MustNewDecFromStr("0.05"),
		ReferencePairRatio:  sdk.MustNewDecFromStr("0.1"),
		ReceiverAddress:     "0x98950be0984d7cf7f5a098a6d8e53fc9c956d4bc",
	}

	testCases := []struct {
		name    string
		msg     sdk.Msg
		expPass bool
	}{
		{
			"valid MsgCreateStoploss",
			NewMsgCreateStoploss(addr1, stoploss),
			true,
		},
		{
			"invalid address",
			NewMsgCreateStoploss(sdk.AccAddress{}, stoploss),
			false,
		},
		{
			"nil stoploss",
			NewMsgCreateStoploss(addr1, nil),
			false,
		},
		{
			"invalid stoploss",
			NewMsgCreateStoploss(addr1, &Stoploss{}),
			false,
		},
	}

	for _, tc := range testCases {

		err := tc.msg.ValidateBasic()

		if tc.expPass {
			require.NoError(t, err, tc.name)
		} else {
			require.Error(t, err, tc.name)
		}
	}
}
