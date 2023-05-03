package types

import (
	"fmt"
)

func (cc ChainConfiguration) ValidateBasic() error {
	if cc.ProxyAddress == "" {
		return fmt.Errorf("proxy address cannot be empty")
	}
	
	if cc.Id == 0 {
		return fmt.Errorf("chain ID cannot be zero")
	}

	if cc.Name == "" {
		return fmt.Errorf("chain name cannot be empty")
	}

	if cc.VoteThreshold.IsZero() {
		return fmt.Errorf("chain vote threshold cannot be zero")
	}

	return nil
}
