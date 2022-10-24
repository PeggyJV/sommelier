package keeper

import (
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/peggyjv/sommelier/v4/x/cork/types"
	"github.com/stretchr/testify/assert"
)

var (
	sampleCellarHex  = "0xc0ffee254729296a45a3885639AC7E10F9d54979"
	sampleCellarAddr = common.HexToAddress(sampleCellarHex)
)

func TestCellarIDs_SetGetHas(t *testing.T) {
	t.Run("happy path", func(t *testing.T) {
		k, ctx, _, _ := setupCorkKeeper(t)

		cellarID := sampleCellarAddr
		cellarAddrs := k.GetCellarIDs(ctx)
		assert.Len(t, cellarAddrs, 0)
		assert.Equal(t, false, k.HasCellarID(ctx, cellarID))

		k.SetCellarIDs(ctx, types.CellarIDSet{
			Ids: []string{cellarID.String()},
		})
		cellarAddrs = k.GetCellarIDs(ctx)
		assert.Len(t, cellarAddrs, 1)
		assert.Contains(t, cellarAddrs[0].String(), cellarID.String())
		assert.Equal(t, true, k.HasCellarID(ctx, cellarID))
	})
}

func TestSetSecheduledCorkGetSecheduledCCork_Unit(t *testing.T) {

}

func TestGetWinningVotes_Unit(t *testing.T) {

}
