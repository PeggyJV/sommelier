package types

import (
	"fmt"

	errorsmod "cosmossdk.io/errors"
	"github.com/ethereum/go-ethereum/common"
)

func (cc ChainConfiguration) ValidateBasic() error {
	if cc.ProxyAddress == "" {
		return fmt.Errorf("proxy address cannot be empty")
	}

	if !common.IsHexAddress(cc.ProxyAddress) {
		return errorsmod.Wrapf(ErrInvalidEVMAddress, "%s", cc.ProxyAddress)
	}

	if cc.Id == 0 {
		return fmt.Errorf("chain ID cannot be zero")
	}

	if cc.Name == "" {
		return fmt.Errorf("chain name cannot be empty")
	}

	return nil
}
