package types

import (
	"crypto/sha256"
	"fmt"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// DataHash returns the hash for a precommit given the proper args
func DataHash(salt, jsonData string, signer sdk.ValAddress) []byte {
	h := sha256.New()
	h.Write([]byte(fmt.Sprintf("%s:%s:%s", salt, jsonData, signer.String())))
	return h.Sum(nil)
}

func (ac *Allocation) ValidateBasic() error {
	if ac.CellarId == "" {
		return fmt.Errorf("no cellar id provided for allocation")
	}

	totalWeight := sdk.ZeroDec()
	for _, pa := range ac.PoolAllocations.Allocations {
		if err := pa.ValidateBasic(); err != nil {
			return err
		}
		totalWeight = totalWeight.Add(pa.TotalWeight())
	}
	if !totalWeight.Equal(sdk.OneDec()) {
		return fmt.Errorf("tick weights do not total to 1.0")
	}

	return nil
}

func (pa *PoolAllocation) TotalWeight() sdk.Dec {
	totalWeight := sdk.ZeroDec()
	for _, tw := range pa.TickWeights.Weights {
		totalWeight = totalWeight.Add(tw.Weight)
	}
	return totalWeight
}

func (pa *PoolAllocation) ValidateBasic() error {
	if pa.FeeLevel.LTE(sdk.ZeroDec()) {
		return fmt.Errorf("invalid fee level provided")
	}

	for _, tw := range pa.TickWeights.Weights {
		if tw.Tick <= 0 {
			return fmt.Errorf("no tick provided for tick weight")
		}
		if tw.Weight.LTE(sdk.ZeroDec()) {
			return fmt.Errorf("invalid tick weight provided")
		}
	}

	return nil
}
