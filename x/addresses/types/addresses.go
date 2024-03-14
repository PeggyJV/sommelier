package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/ethereum/go-ethereum/common"
)

func (am AddressMapping) ValidateBasic() error {
	if _, err := sdk.AccAddressFromBech32(am.CosmosAddress); err != nil {
		return sdkerrors.Wrapf(ErrInvalidCosmosAddress, "%s is not a valid cosmos address", am.CosmosAddress)
	}

	if !common.IsHexAddress(am.EvmAddress) {
		return sdkerrors.Wrapf(ErrInvalidEvmAddress, "%s is not a valid EVM address", am.EvmAddress)
	}

	return nil
}
