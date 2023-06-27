package types

import (
	"testing"

	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/stretchr/testify/require"
)

func TestScheduledCorkProposalValidation(t *testing.T) {
	testCases := []struct {
		name                  string
		scheduledCorkProposal AxelarScheduledCorkProposal
		expPass               bool
		err                   error
	}{
		{
			name: "Happy path",
			scheduledCorkProposal: AxelarScheduledCorkProposal{
				Title:                 "Scheduled Cork",
				Description:           "Schedules a cork via governance",
				BlockHeight:           1,
				ContractCallProtoJson: "{\"thing\":1}",
				TargetContractAddress: "0x0000000000000000000000000000000000000000",
			},
			expPass: true,
			err:     nil,
		},
		{
			name: "Contract address invalid",
			scheduledCorkProposal: AxelarScheduledCorkProposal{
				Title:                 "Scheduled Cork",
				Description:           "Schedules a cork via governance",
				BlockHeight:           1,
				ContractCallProtoJson: "{\"thing\":1}",
				TargetContractAddress: "0x01",
			},
			expPass: false,
			err:     sdkerrors.Wrapf(ErrInvalidEVMAddress, "0x01"),
		},
		{
			name: "Empty proto JSON",
			scheduledCorkProposal: AxelarScheduledCorkProposal{
				Title:                 "Scheduled Cork",
				Description:           "Schedules a cork via governance",
				BlockHeight:           1,
				ContractCallProtoJson: "",
				TargetContractAddress: "0x0000000000000000000000000000000000000000",
			},
			expPass: false,
			err:     sdkerrors.Wrap(ErrInvalidJSON, "cannot have empty contract call"),
		},
		{
			name: "Invalid JSON",
			scheduledCorkProposal: AxelarScheduledCorkProposal{
				Title:                 "Scheduled Cork",
				Description:           "Schedules a cork via governance",
				BlockHeight:           1,
				ContractCallProtoJson: "[}",
				TargetContractAddress: "0x0000000000000000000000000000000000000000",
			},
			expPass: false,
			err:     sdkerrors.Wrapf(ErrInvalidJSON, "[}"),
		},
	}

	for _, tc := range testCases {
		err := tc.scheduledCorkProposal.ValidateBasic()
		if tc.expPass {
			require.NoError(t, err, tc.name)
			require.Nil(t, err)
		} else {
			require.Error(t, err, tc.name)
			require.Equal(t, tc.err.Error(), err.Error())
		}
	}
}
