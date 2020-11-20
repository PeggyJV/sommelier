package types

import (
	"sort"
	"testing"

	"github.com/cosmos/cosmos-sdk/x/auth"
	"github.com/tendermint/tendermint/crypto"
	"github.com/tendermint/tendermint/crypto/secp256k1"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"
)

func TestMsgExchangeRatePrevote(t *testing.T) {
	_, addrs, _, _ := CreateGenAccounts(1, sdk.Coins{})

	bz := GetVoteHash("1", sdk.OneDec(), MicroSDRDenom, sdk.ValAddress(addrs[0]))

	tests := []struct {
		hash       VoteHash
		denom      string
		voter      sdk.AccAddress
		expectPass bool
	}{
		{bz, "", addrs[0], false},
		{bz, MicroCNYDenom, addrs[0], true},
		{bz, MicroCNYDenom, addrs[0], true},
		{bz, MicroCNYDenom, sdk.AccAddress{}, false},
		{VoteHash{}, MicroCNYDenom, addrs[0], false},
	}

	for i, tc := range tests {
		msg := NewMsgExchangeRatePrevote(tc.hash, tc.denom, tc.voter, sdk.ValAddress(tc.voter))
		if tc.expectPass {
			require.Nil(t, msg.ValidateBasic(), "test: %v", i)
		} else {
			require.NotNil(t, msg.ValidateBasic(), "test: %v", i)
		}
	}
}

func TestMsgExchangeRateVote(t *testing.T) {
	_, addrs, _, _ := CreateGenAccounts(1, sdk.Coins{})

	overflowExchangeRate, _ := sdk.NewDecFromStr("100000000000000000000000000000000000000000000000000000000")

	tests := []struct {
		denom      string
		voter      sdk.AccAddress
		salt       string
		rate       sdk.Dec
		expectPass bool
	}{
		{"", addrs[0], "123", sdk.OneDec(), false},
		{MicroCNYDenom, addrs[0], "123", sdk.OneDec().MulInt64(MicroUnit), true},
		{MicroCNYDenom, addrs[0], "123", sdk.ZeroDec(), true},
		{MicroCNYDenom, addrs[0], "123", overflowExchangeRate, false},
		{MicroCNYDenom, sdk.AccAddress{}, "123", sdk.OneDec().MulInt64(MicroUnit), false},
		{MicroCNYDenom, addrs[0], "", sdk.OneDec().MulInt64(MicroUnit), false},
	}

	for i, tc := range tests {
		msg := NewMsgExchangeRateVote(tc.rate, tc.salt, tc.denom, tc.voter, sdk.ValAddress(tc.voter))
		if tc.expectPass {
			require.Nil(t, msg.ValidateBasic(), "test: %v", i)
		} else {
			require.NotNil(t, msg.ValidateBasic(), "test: %v", i)
		}
	}
}

func TestMsgFeederDelegation(t *testing.T) {
	_, addrs, _, _ := CreateGenAccounts(2, sdk.Coins{})

	tests := []struct {
		delegator  sdk.ValAddress
		delegate   sdk.AccAddress
		expectPass bool
	}{
		{sdk.ValAddress(addrs[0]), addrs[1], true},
		{sdk.ValAddress{}, addrs[1], false},
		{sdk.ValAddress(addrs[0]), sdk.AccAddress{}, false},
		{nil, nil, false},
	}

	for i, tc := range tests {
		msg := NewMsgDelegateFeedConsent(tc.delegator, tc.delegate)
		if tc.expectPass {
			require.Nil(t, msg.ValidateBasic(), "test: %v", i)
		} else {
			require.NotNil(t, msg.ValidateBasic(), "test: %v", i)
		}
	}
}

func TestMsgAggregateExchangeRatePrevote(t *testing.T) {
	_, addrs, _, _ := CreateGenAccounts(1, sdk.Coins{})

	exchangeRates := sdk.DecCoins{sdk.NewDecCoinFromDec(MicroSDRDenom, sdk.OneDec()), sdk.NewDecCoinFromDec(MicroKRWDenom, sdk.NewDecWithPrec(32121, 1))}
	bz := GetAggregateVoteHash("1", exchangeRates.String(), sdk.ValAddress(addrs[0]))

	tests := []struct {
		hash          AggregateVoteHash
		exchangeRates sdk.DecCoins
		voter         sdk.AccAddress
		expectPass    bool
	}{
		{bz, exchangeRates, addrs[0], true},
		{bz[1:], exchangeRates, addrs[0], false},
		{bz, exchangeRates, sdk.AccAddress{}, false},
		{AggregateVoteHash{}, exchangeRates, addrs[0], false},
	}

	for i, tc := range tests {
		msg := NewMsgAggregateExchangeRatePrevote(tc.hash, tc.voter, sdk.ValAddress(tc.voter))
		if tc.expectPass {
			require.NoError(t, msg.ValidateBasic(), "test: %v", i)
		} else {
			require.Error(t, msg.ValidateBasic(), "test: %v", i)
		}
	}
}

func TestMsgAggregateExchangeRateVote(t *testing.T) {
	_, addrs, _, _ := CreateGenAccounts(1, sdk.Coins{})

	invalidExchangeRates := "a,b"
	exchangeRates := "1.0foo,1232.132bar"
	abstainExchangeRates := "0.0foo,1232.132bar"
	overFlowExchangeRates := "100000000000000000000000000000000000000000000000000000000.0foo,1232.132bar"

	tests := []struct {
		voter         sdk.AccAddress
		salt          string
		exchangeRates string
		expectPass    bool
	}{
		{addrs[0], "123", exchangeRates, true},
		{addrs[0], "123", invalidExchangeRates, false},
		{addrs[0], "123", abstainExchangeRates, true},
		{addrs[0], "123", overFlowExchangeRates, false},
		{sdk.AccAddress{}, "123", exchangeRates, false},
		{addrs[0], "", exchangeRates, false},
	}

	for i, tc := range tests {
		msg := NewMsgAggregateExchangeRateVote(tc.salt, tc.exchangeRates, tc.voter, sdk.ValAddress(tc.voter))
		if tc.expectPass {
			require.Nil(t, msg.ValidateBasic(), "test: %v", i)
		} else {
			require.NotNil(t, msg.ValidateBasic(), "test: %v", i)
		}
	}
}

// CreateGenAccounts generates genesis accounts loaded with coins, and returns
// their addresses, pubkeys, and privkeys.
func CreateGenAccounts(numAccs int, genCoins sdk.Coins) (genAccs []authexported.Account,
	addrs []sdk.AccAddress, pubKeys []crypto.PubKey, privKeys []crypto.PrivKey) {

	addrKeysSlice := AddrKeysSlice{}

	for i := 0; i < numAccs; i++ {
		privKey := secp256k1.GenPrivKey()
		pubKey := privKey.PubKey()
		addr := sdk.AccAddress(pubKey.Address())

		addrKeysSlice = append(addrKeysSlice, NewAddrKeys(addr, pubKey, privKey))
	}

	sort.Sort(addrKeysSlice)

	for i := range addrKeysSlice {
		addrs = append(addrs, addrKeysSlice[i].Address)
		pubKeys = append(pubKeys, addrKeysSlice[i].PubKey)
		privKeys = append(privKeys, addrKeysSlice[i].PrivKey)
		genAccs = append(genAccs, &auth.BaseAccount{
			Address: addrKeysSlice[i].Address,
			Coins:   genCoins,
		})
	}

	return
}
