package types

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestFeeAccrualCounters(t *testing.T) {
	expected := FeeAccrualCounters{
		Counters: []FeeAccrualCounter{
			{
				Denom: "uatom",
				Count: 1,
			},
			{
				Denom: "uist",
				Count: 0,
			},
			{
				Denom: "uusdc",
				Count: 2,
			},
		},
	}

	actual := FeeAccrualCounters{
		Counters: make([]FeeAccrualCounter, 0),
	}

	require.Equal(t, 0, len(actual.Counters))

	// uist: 0
	actual.ResetCounter("uist")
	require.Equal(t, 1, len(actual.Counters))
	require.Equal(t, uint64(0), actual.Counters[0].Count)

	// uist: 0
	// uusdc: 1
	actual.IncrementCounter("uusdc")
	require.Equal(t, len(actual.Counters), 2)
	require.Equal(t, "uusdc", actual.Counters[1].Denom)
	require.Equal(t, uint64(1), actual.Counters[1].Count)

	// uist: 0
	// uusdc: 2
	actual.IncrementCounter("uusdc")

	// uatom: 1
	// uist: 0
	// uusdc: 2
	actual.IncrementCounter("uatom")
	require.Equal(t, len(actual.Counters), 3)
	require.Equal(t, "uatom", actual.Counters[0].Denom)
	require.Equal(t, uint64(1), actual.Counters[0].Count)
	require.Equal(t, expected, actual)

	// uatom: 1
	// uist: 0
	// uusdc: 0
	actual.ResetCounter("uusdc")

	// uatom: 1
	// uist: 0
	// uusdc: 1
	require.Equal(t, uint64(1), actual.IncrementCounter("uusdc"))
}
