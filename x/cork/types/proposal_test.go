package types

import (
	"testing"

	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/stretchr/testify/require"
)

func TestNewAddManagedCellarIDsProposal(t *testing.T) {
	expectedMsg := &AddManagedCellarIDsProposal{
		Title:       "Add Cellar IDs",
		Description: "Adds Cellar IDs",
		CellarIds:   &CellarIDSet{Ids: []string{"1", "2"}},
	}

	createdMsg := NewAddManagedCellarIDsProposal("Add Cellar IDs", "Adds Cellar IDs", &CellarIDSet{Ids: []string{"1", "2"}})
	require.Equal(t, expectedMsg, createdMsg)
}

func TestNewRemoveManagedCellarIDsProposal(t *testing.T) {
	expectedMsg := &RemoveManagedCellarIDsProposal{
		Title:       "Remove Cellar IDs",
		Description: "Remove Cellar IDs",
		CellarIds:   &CellarIDSet{Ids: []string{"1", "2"}},
	}

	createdMsg := NewRemoveManagedCellarIDsProposal("Remove Cellar IDs", "Remove Cellar IDs", &CellarIDSet{Ids: []string{"1", "2"}})
	require.Equal(t, expectedMsg, createdMsg)
}

func TestScheduledCorkProposal(t *testing.T) {
	expectedMsg := &ScheduledCorkProposal{
		Title:                 "Scheduled Cork",
		Description:           "Schedules a cork via governance",
		BlockHeight:           1,
		ContractCallProtoJson: "{\"thing\":1}",
		TargetContractAddress: "0x0000000000000000000000000000000000000000",
	}

	createdMsg, err := NewScheduledCorkProposal("Scheduled Cork", "Schedules a cork via governance", 1, "0x0000000000000000000000000000000000000000", "{\"thing\":1}")
	require.Nil(t, err)
	require.Equal(t, expectedMsg, createdMsg)
}

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
			err:     sdkerrors.Wrapf(ErrInvalidEthereumAddress, "0x01"),
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
			err:     sdkerrors.Wrap(ErrInvalidJson, "cannot have empty contract call"),
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
			err:     sdkerrors.Wrapf(ErrInvalidJson, "[}"),
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
