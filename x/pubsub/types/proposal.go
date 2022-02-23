package types

import govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"

const (
	ProposalTypeAddPublisher    = "AddPublisher"
	ProposalTypeRemovePublisher = "RemovePublisher"
)

var _ govtypes.Content = &AddPublisherProposal{}
var _ govtypes.Content = &RemovePublisherProposal{}

// TODO(bolten): fill out proposal boilerplate

//////////////////////////
// AddPublisherProposal //
//////////////////////////

/////////////////////////////
// RemovePublisherProposal //
/////////////////////////////
