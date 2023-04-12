package types

import "github.com/ethereum/go-ethereum/common"

func (cc *ChainConfiguration) HasCellarID(cellarID common.Address) bool {
	for _, cellarIDStr := range cc.CellarIds.Ids {
		addr := common.HexToAddress(cellarIDStr)
		if addr == cellarID {
			return true
		}
	}

	return false
}

func (cc *ChainConfiguration) GetCellarIDs() (cellarIDs []common.Address) {
	for _, cellarIDStr := range cc.CellarIds.Ids {
		addr := common.HexToAddress(cellarIDStr)
		cellarIDs = append(cellarIDs, addr)
	}

	return cellarIDs
}