package keeper

import (
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"
)

func TestGetApportionments(t *testing.T) {
	tests := []struct {
		name        string
		numPortions uint64
		value       sdk.Dec
		maxPortion  sdk.Dec
		want        []sdk.Dec
		wantErr     bool
	}{
		{
			name:        "Negative value",
			numPortions: 10,
			value:       sdk.NewDec(-100),
			maxPortion:  sdk.MustNewDecFromStr("0.1"),
			wantErr:     true,
		},
		{
			name:        "Zero length slice",
			numPortions: 0,
			value:       sdk.NewDec(100),
			maxPortion:  sdk.MustNewDecFromStr("0.1"),
			want:        make([]sdk.Dec, 0),
		},
		{
			name:        "Zero value",
			numPortions: 10,
			value:       sdk.NewDec(0),
			maxPortion:  sdk.MustNewDecFromStr("0.1"),
			want:        []sdk.Dec{sdk.ZeroDec(), sdk.ZeroDec(), sdk.ZeroDec(), sdk.ZeroDec(), sdk.ZeroDec(), sdk.ZeroDec(), sdk.ZeroDec(), sdk.ZeroDec(), sdk.ZeroDec(), sdk.ZeroDec()},
		},
		{
			name:        "Max portion",
			numPortions: 10,
			value:       sdk.NewDec(100),
			maxPortion:  sdk.OneDec(),
			want:        []sdk.Dec{sdk.NewDec(100), sdk.ZeroDec(), sdk.ZeroDec(), sdk.ZeroDec(), sdk.ZeroDec(), sdk.ZeroDec(), sdk.ZeroDec(), sdk.ZeroDec(), sdk.ZeroDec(), sdk.ZeroDec()},
		},
		{
			name:        "Max portion is zero",
			numPortions: 10,
			value:       sdk.NewDec(100),
			maxPortion:  sdk.ZeroDec(),
			want:        []sdk.Dec{sdk.ZeroDec(), sdk.ZeroDec(), sdk.ZeroDec(), sdk.ZeroDec(), sdk.ZeroDec(), sdk.ZeroDec(), sdk.ZeroDec(), sdk.ZeroDec(), sdk.ZeroDec(), sdk.ZeroDec()},
		},
		{
			name:        "Max portion is negative",
			numPortions: 10,
			value:       sdk.NewDec(100),
			maxPortion:  sdk.MustNewDecFromStr("-0.1"),
			wantErr:     true,
		},
		{
			name:        "Max portion greater than 1",
			numPortions: 10,
			value:       sdk.NewDec(100),
			maxPortion:  sdk.NewDec(2),
			wantErr:     true,
		},
		{
			name:        "Specific distribution with 50 portions with 5% max portion",
			numPortions: 50,
			value:       sdk.NewDec(100),
			maxPortion:  sdk.MustNewDecFromStr("0.05"),
			want: []sdk.Dec{
				sdk.MustNewDecFromStr("5"),
				sdk.MustNewDecFromStr("4.75"),
				sdk.MustNewDecFromStr("4.5125"),
				sdk.MustNewDecFromStr("4.286875"),
				sdk.MustNewDecFromStr("4.07253125"),
				sdk.MustNewDecFromStr("3.868904688"),
				sdk.MustNewDecFromStr("3.675459453"),
				sdk.MustNewDecFromStr("3.49168648"),
				sdk.MustNewDecFromStr("3.317102156"),
				sdk.MustNewDecFromStr("3.151247049"),
				sdk.MustNewDecFromStr("2.993684696"),
				sdk.MustNewDecFromStr("2.844000461"),
				sdk.MustNewDecFromStr("2.701800438"),
				sdk.MustNewDecFromStr("2.566710416"),
				sdk.MustNewDecFromStr("2.438374896"),
				sdk.MustNewDecFromStr("2.316456151"),
				sdk.MustNewDecFromStr("2.200633343"),
				sdk.MustNewDecFromStr("2.090601676"),
				sdk.MustNewDecFromStr("1.986071592"),
				sdk.MustNewDecFromStr("1.886768013"),
				sdk.MustNewDecFromStr("1.792429612"),
				sdk.MustNewDecFromStr("1.702808131"),
				sdk.MustNewDecFromStr("1.617667725"),
				sdk.MustNewDecFromStr("1.536784339"),
				sdk.MustNewDecFromStr("1.459945122"),
				sdk.MustNewDecFromStr("1.386947866"),
				sdk.MustNewDecFromStr("1.317600472"),
				sdk.MustNewDecFromStr("1.251720449"),
				sdk.MustNewDecFromStr("1.189134426"),
				sdk.MustNewDecFromStr("1.129677705"),
				sdk.MustNewDecFromStr("1.07319382"),
				sdk.MustNewDecFromStr("1.019534129"),
				sdk.MustNewDecFromStr("0.9685574223"),
				sdk.MustNewDecFromStr("0.9201295512"),
				sdk.MustNewDecFromStr("0.8741230736"),
				sdk.MustNewDecFromStr("0.8304169199"),
				sdk.MustNewDecFromStr("0.7888960739"),
				sdk.MustNewDecFromStr("0.7494512702"),
				sdk.MustNewDecFromStr("0.7119787067"),
				sdk.MustNewDecFromStr("0.6763797714"),
				sdk.MustNewDecFromStr("0.6425607828"),
				sdk.MustNewDecFromStr("0.6104327437"),
				sdk.MustNewDecFromStr("0.5799111065"),
				sdk.MustNewDecFromStr("0.5509155512"),
				sdk.MustNewDecFromStr("0.5233697736"),
				sdk.MustNewDecFromStr("0.4972012849"),
				sdk.MustNewDecFromStr("0.4723412207"),
				sdk.MustNewDecFromStr("0.4487241597"),
				sdk.MustNewDecFromStr("0.4262879517"),
				sdk.MustNewDecFromStr("0.4049735541"),
			},
		},
		{
			name:        "Specific distribution with 50 portions with 10% max portion",
			numPortions: 50,
			value:       sdk.NewDec(100),
			maxPortion:  sdk.MustNewDecFromStr("0.1"),
			want: []sdk.Dec{
				sdk.MustNewDecFromStr("10"),
				sdk.MustNewDecFromStr("9"),
				sdk.MustNewDecFromStr("8.1"),
				sdk.MustNewDecFromStr("7.29"),
				sdk.MustNewDecFromStr("6.561"),
				sdk.MustNewDecFromStr("5.9049"),
				sdk.MustNewDecFromStr("5.31441"),
				sdk.MustNewDecFromStr("4.782969"),
				sdk.MustNewDecFromStr("4.3046721"),
				sdk.MustNewDecFromStr("3.87420489"),
				sdk.MustNewDecFromStr("3.486784401"),
				sdk.MustNewDecFromStr("3.138105961"),
				sdk.MustNewDecFromStr("2.824295365"),
				sdk.MustNewDecFromStr("2.541865828"),
				sdk.MustNewDecFromStr("2.287679245"),
				sdk.MustNewDecFromStr("2.058911321"),
				sdk.MustNewDecFromStr("1.853020189"),
				sdk.MustNewDecFromStr("1.66771817"),
				sdk.MustNewDecFromStr("1.500946353"),
				sdk.MustNewDecFromStr("1.350851718"),
				sdk.MustNewDecFromStr("1.215766546"),
				sdk.MustNewDecFromStr("1.094189891"),
				sdk.MustNewDecFromStr("0.9847709022"),
				sdk.MustNewDecFromStr("0.886293812"),
				sdk.MustNewDecFromStr("0.7976644308"),
				sdk.MustNewDecFromStr("0.7178979877"),
				sdk.MustNewDecFromStr("0.6461081889"),
				sdk.MustNewDecFromStr("0.58149737"),
				sdk.MustNewDecFromStr("0.523347633"),
				sdk.MustNewDecFromStr("0.4710128697"),
				sdk.MustNewDecFromStr("0.4239115828"),
				sdk.MustNewDecFromStr("0.3815204245"),
				sdk.MustNewDecFromStr("0.343368382"),
				sdk.MustNewDecFromStr("0.3090315438"),
				sdk.MustNewDecFromStr("0.2781283894"),
				sdk.MustNewDecFromStr("0.2503155505"),
				sdk.MustNewDecFromStr("0.2252839954"),
				sdk.MustNewDecFromStr("0.2027555959"),
				sdk.MustNewDecFromStr("0.1824800363"),
				sdk.MustNewDecFromStr("0.1642320327"),
				sdk.MustNewDecFromStr("0.1478088294"),
				sdk.MustNewDecFromStr("0.1330279465"),
				sdk.MustNewDecFromStr("0.1197251518"),
				sdk.MustNewDecFromStr("0.1077526366"),
				sdk.MustNewDecFromStr("0.09697737298"),
				sdk.MustNewDecFromStr("0.08727963568"),
				sdk.MustNewDecFromStr("0.07855167211"),
				sdk.MustNewDecFromStr("0.0706965049"),
				sdk.MustNewDecFromStr("0.06362685441"),
				sdk.MustNewDecFromStr("0.05726416897"),
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt := tt
			got, remaining, err := getApportionments(tt.numPortions, tt.value, tt.maxPortion)
			if tt.wantErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)

				// Correct number of portions
				require.Equal(t, tt.numPortions, uint(len(got)))

				// We round before comparing to avoid floating point precision issues.
				// Check that the values returned in the slice line up with the expected values.
				for i, value := range got {
					// Correct apportionment values
					require.Equal(t, roundToPrecision(tt.want[i], 6), roundToPrecision(value, 6))

					// No negative apportionments
					require.True(t, value.GTE(sdk.ZeroDec()))
				}

				// Non-negative remaining value
				require.True(t, remaining.GTE(sdk.ZeroDec()))

				// Sum of apportionments plus remaining equals the original value
				gotSum := sdk.ZeroDec()
				for _, g := range got {
					gotSum = gotSum.Add(g)
				}
				require.True(t, gotSum.LTE(tt.value))
				require.Equal(t, tt.value, gotSum.Add(remaining))

				// Test that each value is greater than or equal to the next
				for i := 0; i < len(got)-1; i++ {
					if i == len(got)-1 {
						break
					}

					require.True(t, got[i].GTE(got[i+1]), "Value at index %d should be greater than or equal to value at index %d", i, i+1)
				}
			}
		})
	}
}

func roundToPrecision(d sdk.Dec, precision uint64) sdk.Dec {
	multiplier := sdk.NewDec(10).Power(precision)
	return d.Mul(multiplier).RoundInt().ToLegacyDec().Quo(multiplier)
}
