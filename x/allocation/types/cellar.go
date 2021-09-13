package types

import (
	"encoding/json"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/ethereum/go-ethereum/common"
	tmbytes "github.com/tendermint/tendermint/libs/bytes"
)

func (c *Cellar) Address() common.Address {
	return common.HexToAddress(c.Id)
}

func (c *Cellar) InvalidationScope() tmbytes.HexBytes {
	return c.GetCheckpoint()
}

func (c *Cellar) Equals(other Cellar) bool {
	if c.Id != other.Id {
		return false
	}

	if len(c.TickRanges) != len(other.TickRanges) {
		return false
	}

	for _, tr := range c.TickRanges {
		found := false
		for _, otr := range other.TickRanges {
			if tr.Equals(*otr) {
				found = true
			}
		}

		if !found {
			return false
		}
	}

	return true
}

func (c *Cellar) Hash(salt string, val sdk.ValAddress) ([]byte, error) {
	jsonBz, err := json.Marshal(c)
	if err != nil {
		return nil, sdkerrors.Wrap(
			sdkerrors.ErrJSONMarshal, "failed to marshal cellar",
		)
	}

	jsonBz = sdk.MustSortJSON(jsonBz)
	commitHash := DataHash(salt, string(jsonBz), val)
	return commitHash, nil
}
