package types

import codectypes "github.com/cosmos/cosmos-sdk/codec/types"

// BlocksTillNextPeriod helper
func (vp *VotePeriod) BlocksTillNextPeriod() int64 {
	return vp.VotePeriodEnd - vp.CurrentHeight
}

var (
	_ codectypes.UnpackInterfacesMessage = &OracleFeed{}
	_ codectypes.UnpackInterfacesMessage = &OracleVote{}
)

// UnpackInterfaces implements UnpackInterfacesMessage
func (ov *OracleVote) UnpackInterfaces(unpacker codectypes.AnyUnpacker) error {
	return ov.Feed.UnpackInterfaces(unpacker)
}

// UnpackInterfaces implements UnpackInterfacesMessage
func (of *OracleFeed) UnpackInterfaces(unpacker codectypes.AnyUnpacker) error {
	for _, oracleDataAny := range of.OracleData {
		var od OracleData
		if err := unpacker.UnpackAny(oracleDataAny, &od); err != nil {
			return err
		}
	}

	return nil
}
