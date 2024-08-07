package types

import (
	errorsmod "cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/ethereum/go-ethereum/common"
)

func (am AddressMapping) ValidateBasic() error {
	if _, err := sdk.AccAddressFromBech32(am.CosmosAddress); err != nil {
		return errorsmod.Wrapf(ErrInvalidCosmosAddress, "%s is not a valid cosmos address", am.CosmosAddress)
	}

	if !common.IsHexAddress(am.EvmAddress) {
		return errorsmod.Wrapf(ErrInvalidEvmAddress, "%s is not a valid EVM address", am.EvmAddress)
	}

	return nil
}
