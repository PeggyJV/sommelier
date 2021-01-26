package types

// import (
// 	"fmt"

// 	yaml "gopkg.in/yaml.v2"

// 	sdk "github.com/cosmos/cosmos-sdk/types"
// 	paramtypes "github.com/cosmos/cosmos-sdk/x/params/types"
// )

// // Parameter keys
// var (
// 	ParamStoreKeyCommunityTax        = []byte("communitytax")
// )

// // ParamKeyTable returns the parameter key table.
// func ParamKeyTable() paramtypes.KeyTable {
// 	return paramtypes.NewKeyTable().RegisterParamSet(&Params{})
// }

// // DefaultParams returns default distribution parameters
// func DefaultParams() Params {
// 	return Params{
// 		CommunityTax:        sdk.NewDecWithPrec(2, 2), // 2%
// 	}
// }

// // ParamSetPairs returns the parameter set pairs.
// func (p *Params) ParamSetPairs() paramtypes.ParamSetPairs {
// 	return paramtypes.ParamSetPairs{
// 		paramtypes.NewParamSetPair(ParamStoreKeyCommunityTax, &p.CommunityTax, validateCommunityTax),
// 	}
// }

// // ValidateBasic performs basic validation on distribution parameters.
// func (p Params) ValidateBasic() error {
// 	if p.CommunityTax.IsNegative() || p.CommunityTax.GT(sdk.OneDec()) {
// 		return fmt.Errorf(
// 			"community tax should non-negative and less than one: %s", p.CommunityTax,
// 		)
// 	}

// 	return nil
// }

// func validateCommunityTax(i interface{}) error {
// 	v, ok := i.(sdk.Dec)
// 	if !ok {
// 		return fmt.Errorf("invalid parameter type: %T", i)
// 	}

// 	if v.IsNil() {
// 		return fmt.Errorf("community tax must be not nil")
// 	}
// 	if v.IsNegative() {
// 		return fmt.Errorf("community tax must be positive: %s", v)
// 	}
// 	if v.GT(sdk.OneDec()) {
// 		return fmt.Errorf("community tax too large: %s", v)
// 	}

// 	return nil
// }
