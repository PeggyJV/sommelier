package rest

import (
	"net/http"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/tx"
	"github.com/cosmos/cosmos-sdk/types/rest"
	govrest "github.com/cosmos/cosmos-sdk/x/gov/client/rest"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"
	"github.com/peggyjv/sommelier/v4/x/pubsub/types"
)

// AddPublisherProposalRESTHandler returns a ProposalRESTHandler that exposes an add publisher REST handler with a given sub-route.
func AddPublisherProposalRESTHandler(clientCtx client.Context) govrest.ProposalRESTHandler {
	return govrest.ProposalRESTHandler{
		SubRoute: "add_publisher",
		Handler:  postAddPublisherProposalHandlerFn(clientCtx),
	}
}

// RemovePublisherProposalRESTHandler returns a ProposalRESTHandler that exposes a remove publisher REST handler with a given sub-route.
func RemovePublisherProposalRESTHandler(clientCtx client.Context) govrest.ProposalRESTHandler {
	return govrest.ProposalRESTHandler{
		SubRoute: "remove_publisher",
		Handler:  postRemovePublisherProposalHandlerFn(clientCtx),
	}
}

// AddDefaultSubscriptionProposalRESTHandler returns a ProposalRESTHandler that exposes an add default subscription REST handler with a given sub-route.
func AddDefaultSubscriptionProposalRESTHandler(clientCtx client.Context) govrest.ProposalRESTHandler {
	return govrest.ProposalRESTHandler{
		SubRoute: "add_default_subscription",
		Handler:  postAddDefaultSubscriptionProposalHandlerFn(clientCtx),
	}
}

// RemoveDefaultSubscriptionProposalRESTHandler returns a ProposalRESTHandler that exposes a remove default subscription REST handler with a given sub-route.
func RemoveDefaultSubscriptionProposalRESTHandler(clientCtx client.Context) govrest.ProposalRESTHandler {
	return govrest.ProposalRESTHandler{
		SubRoute: "remove_default_subscription",
		Handler:  postRemoveDefaultSubscriptionProposalHandlerFn(clientCtx),
	}
}

func postAddPublisherProposalHandlerFn(clientCtx client.Context) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req AddPublisherProposalReq
		if !rest.ReadRESTReq(w, r, clientCtx.LegacyAmino, &req) {
			return
		}

		req.BaseReq = req.BaseReq.Sanitize()
		if !req.BaseReq.ValidateBasic(w) {
			return
		}

		content := types.NewAddPublisherProposal(
			req.Title,
			req.Description,
			req.Domain,
			req.Address,
			req.ProofUrl,
			req.CaCert,
		)

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

func postRemovePublisherProposalHandlerFn(clientCtx client.Context) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req RemovePublisherProposalReq
		if !rest.ReadRESTReq(w, r, clientCtx.LegacyAmino, &req) {
			return
		}

		req.BaseReq = req.BaseReq.Sanitize()
		if !req.BaseReq.ValidateBasic(w) {
			return
		}

		content := types.NewRemovePublisherProposal(
			req.Title,
			req.Description,
			req.Domain,
		)

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

func postAddDefaultSubscriptionProposalHandlerFn(clientCtx client.Context) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req AddDefaultSubscriptionProposalReq
		if !rest.ReadRESTReq(w, r, clientCtx.LegacyAmino, &req) {
			return
		}

		req.BaseReq = req.BaseReq.Sanitize()
		if !req.BaseReq.ValidateBasic(w) {
			return
		}

		content := types.NewAddDefaultSubscriptionProposal(
			req.Title,
			req.Description,
			req.SubscriptionId,
			req.PublisherDomain,
		)

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

func postRemoveDefaultSubscriptionProposalHandlerFn(clientCtx client.Context) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req RemoveDefaultSubscriptionProposalReq
		if !rest.ReadRESTReq(w, r, clientCtx.LegacyAmino, &req) {
			return
		}

		req.BaseReq = req.BaseReq.Sanitize()
		if !req.BaseReq.ValidateBasic(w) {
			return
		}

		content := types.NewRemoveDefaultSubscriptionProposal(
			req.Title,
			req.Description,
			req.SubscriptionId,
		)

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
