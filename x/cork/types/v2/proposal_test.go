package v2

import (
	"testing"

	errorsmod "cosmossdk.io/errors"
	types "github.com/peggyjv/sommelier/v7/x/cork/types"
	"github.com/stretchr/testify/require"
)

func TestScheduledCorkProposalValidation(t *testing.T) {
	testCases := []struct {
		name                  string
		scheduledCorkProposal ScheduledCorkProposal
		expPass               bool
		err                   error
	}{
		{
			name: "Happy path",
			scheduledCorkProposal: ScheduledCorkProposal{
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
			scheduledCorkProposal: ScheduledCorkProposal{
				Title:                 "Scheduled Cork",
				Description:           "Schedules a cork via governance",
				BlockHeight:           1,
				ContractCallProtoJson: "{\"thing\":1}",
				TargetContractAddress: "0x01",
			},
			expPass: false,
			err:     errorsmod.Wrapf(types.ErrInvalidEthereumAddress, "0x01"),
		},
		{
			name: "Empty proto JSON",
			scheduledCorkProposal: ScheduledCorkProposal{
				Title:                 "Scheduled Cork",
				Description:           "Schedules a cork via governance",
				BlockHeight:           1,
				ContractCallProtoJson: "",
				TargetContractAddress: "0x0000000000000000000000000000000000000000",
			},
			expPass: false,
			err:     errorsmod.Wrap(types.ErrInvalidJSON, "cannot have empty contract call"),
		},
		{
			name: "Invalid JSON",
			scheduledCorkProposal: ScheduledCorkProposal{
				Title:                 "Scheduled Cork",
				Description:           "Schedules a cork via governance",
				BlockHeight:           1,
				ContractCallProtoJson: "[}",
				TargetContractAddress: "0x0000000000000000000000000000000000000000",
			},
			expPass: false,
			err:     errorsmod.Wrapf(types.ErrInvalidJSON, "[}"),
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
