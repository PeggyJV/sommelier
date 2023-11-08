package rest

import (
	"net/http"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/tx"
	"github.com/cosmos/cosmos-sdk/types/rest"
	govrest "github.com/cosmos/cosmos-sdk/x/gov/client/rest"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"
	"github.com/peggyjv/sommelier/v7/x/cork/types"
)

// AddProposalRESTHandler returns a ProposalRESTHandler that exposes add managed cellar IDs REST handler with a given sub-route.
func AddProposalRESTHandler(clientCtx client.Context) govrest.ProposalRESTHandler {
	return govrest.ProposalRESTHandler{
		SubRoute: "add_managed_cellar_ids",
		Handler:  postAddProposalHandlerFn(clientCtx),
	}
}

// RemoveProposalRESTHandler returns a ProposalRESTHandler that exposes remove managed cellar IDs REST handler with a given sub-route.
func RemoveProposalRESTHandler(clientCtx client.Context) govrest.ProposalRESTHandler {
	return govrest.ProposalRESTHandler{
		SubRoute: "remove_managed_cellar_ids",
		Handler:  postRemoveProposalHandlerFn(clientCtx),
	}
}

// ScheduledCorkProposalRESTHandler returns a ProposalRESTHandler that exposes the scheduled cork REST handler with a given sub-route.
func ScheduledCorkProposalRESTHandler(clientCtx client.Context) govrest.ProposalRESTHandler {
	return govrest.ProposalRESTHandler{
		SubRoute: "scheduled_cork",
		Handler:  postScheduledCorkProposalHandlerFn(clientCtx),
	}
}

func postAddProposalHandlerFn(clientCtx client.Context) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req AddManagedCellarIDsProposalReq
		if !rest.ReadRESTReq(w, r, clientCtx.LegacyAmino, &req) {
			return
		}

		req.BaseReq = req.BaseReq.Sanitize()
		if !req.BaseReq.ValidateBasic(w) {
			return
		}

		content := types.NewAddManagedCellarIDsProposal(
			req.Title,
			req.Description,
			&types.CellarIDSet{
				Ids: req.CellarIDs,
			})

		msg, err := govtypes.NewMsgSubmitProposal(content, req.Deposit, req.Proposer)
		if rest.CheckBadRequestError(w, err) {
			return
		}
		if rest.CheckBadRequestError(w, msg.ValidateBasic()) {
			return
		}

		tx.WriteGeneratedTxResponse(clientCtx, w, req.BaseReq, msg)
	}
}

func postRemoveProposalHandlerFn(clientCtx client.Context) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req RemoveManagedCellarIDsProposalReq
		if !rest.ReadRESTReq(w, r, clientCtx.LegacyAmino, &req) {
			return
		}

		req.BaseReq = req.BaseReq.Sanitize()
		if !req.BaseReq.ValidateBasic(w) {
			return
		}

		content := types.NewRemoveManagedCellarIDsProposal(
			req.Title,
			req.Description,
			&types.CellarIDSet{
				Ids: req.CellarIDs,
			})

		msg, err := govtypes.NewMsgSubmitProposal(content, req.Deposit, req.Proposer)
		if rest.CheckBadRequestError(w, err) {
			return
		}
		if rest.CheckBadRequestError(w, msg.ValidateBasic()) {
			return
		}

		tx.WriteGeneratedTxResponse(clientCtx, w, req.BaseReq, msg)
	}
}

func postScheduledCorkProposalHandlerFn(clientCtx client.Context) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req ScheduledCorkProposalReq
		if !rest.ReadRESTReq(w, r, clientCtx.LegacyAmino, &req) {
			return
		}

		req.BaseReq = req.BaseReq.Sanitize()
		if !req.BaseReq.ValidateBasic(w) {
			return
		}

		content := types.NewScheduledCorkProposal(
			req.Title,
			req.Description,
			req.BlockHeight,
			req.TargetContractAddress,
			req.ContractCallProtoJSON,
		)
		if rest.CheckBadRequestError(w, content.ValidateBasic()) {
			return
		}
		msg, err := govtypes.NewMsgSubmitProposal(content, req.Deposit, req.Proposer)
		if rest.CheckBadRequestError(w, err) {
			return
		}
		if rest.CheckBadRequestError(w, msg.ValidateBasic()) {
			return
		}

		tx.WriteGeneratedTxResponse(clientCtx, w, req.BaseReq, msg)
	}
}
