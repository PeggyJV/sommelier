package keeper

import (
	"encoding/json"
	"fmt"
	"time"

	distributiontypes "github.com/cosmos/cosmos-sdk/x/distribution/types"
	transfertypes "github.com/cosmos/ibc-go/v7/modules/apps/transfer/types"
	clienttypes "github.com/cosmos/ibc-go/v7/modules/core/02-client/types"

	errorsmod "cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/ethereum/go-ethereum/common"
	"github.com/peggyjv/sommelier/v9/x/axelarcork/types"
	pubsubtypes "github.com/peggyjv/sommelier/v9/x/pubsub/types"
)

func NewAxelarSubscriptionID(chainID uint64, address common.Address) string {
	return fmt.Sprintf("%d:%s", chainID, address.String())
}

// HandleAddManagedCellarsProposal is a handler for executing a passed community cellar addition proposal
func HandleAddManagedCellarsProposal(ctx sdk.Context, k Keeper, p types.AddAxelarManagedCellarIDsProposal) error {
	config, ok := k.GetChainConfigurationByID(ctx, p.ChainId)
	if !ok {
		return fmt.Errorf("chain by id %d not found", p.ChainId)
	}

	_, publisherFound := k.pubsubKeeper.GetPublisher(ctx, p.PublisherDomain)
	if !publisherFound {
		return fmt.Errorf("not an approved publisher: %s", p.PublisherDomain)
	}

	if err := p.CellarIds.ValidateBasic(); err != nil {
		return err
	}

	cellarAddresses := k.GetCellarIDs(ctx, config.Id)

	for _, proposedCellarID := range p.CellarIds.Ids {
		proposedCellarAddress := common.HexToAddress(proposedCellarID)
		found := false
		for _, id := range cellarAddresses {
			if id == proposedCellarAddress {
				found = true
			}
		}
		if !found {
			cellarAddresses = append(cellarAddresses, proposedCellarAddress)
			subscriptionID := NewAxelarSubscriptionID(p.ChainId, proposedCellarAddress)
			defaultSubscription := pubsubtypes.DefaultSubscription{
				SubscriptionId:  subscriptionID,
				PublisherDomain: p.PublisherDomain,
			}
			k.pubsubKeeper.SetDefaultSubscription(ctx, defaultSubscription)
		}
	}

	idStrings := make([]string, len(cellarAddresses))
	for i, cid := range cellarAddresses {
		idStrings[i] = cid.String()
	}

	k.SetCellarIDs(ctx, config.Id, types.CellarIDSet{ChainId: config.Id, Ids: idStrings})

	return nil
}

// HandleRemoveManagedCellarsProposal is a handler for executing a passed community cellar removal proposal
func HandleRemoveManagedCellarsProposal(ctx sdk.Context, k Keeper, p types.RemoveAxelarManagedCellarIDsProposal) error {
	config, ok := k.GetChainConfigurationByID(ctx, p.ChainId)
	if !ok {
		return fmt.Errorf("chain by id %d not found", p.ChainId)
	}

	if err := p.CellarIds.ValidateBasic(); err != nil {
		return err
	}

	var outputCellarIDs types.CellarIDSet

	for _, existingID := range k.GetCellarIDs(ctx, config.Id) {
		found := false
		for _, inputID := range p.CellarIds.Ids {
			if existingID == common.HexToAddress(inputID) {
				found = true
			}
		}

		if !found {
			outputCellarIDs.Ids = append(outputCellarIDs.Ids, existingID.String())
		}
	}
	outputCellarIDs.ChainId = config.Id

	// unlike for adding an ID, we don't need to re-sort because we're removing elements from an already sorted list
	k.SetCellarIDs(ctx, config.Id, outputCellarIDs)

	for _, cellarToDelete := range p.CellarIds.Ids {
		subscriptionID := NewAxelarSubscriptionID(p.ChainId, common.HexToAddress(cellarToDelete))
		k.pubsubKeeper.DeleteDefaultSubscription(ctx, subscriptionID)
	}

	return nil
}

// HandleScheduledCorkProposal is a handler for executing a passed scheduled cork proposal
func HandleScheduledCorkProposal(ctx sdk.Context, k Keeper, p types.AxelarScheduledCorkProposal) error {
	config, ok := k.GetChainConfigurationByID(ctx, p.ChainId)
	if !ok {
		return fmt.Errorf("chain by id %d not found", p.ChainId)
	}

	if !k.HasCellarID(ctx, config.Id, common.HexToAddress(p.TargetContractAddress)) {
		return errorsmod.Wrapf(types.ErrUnmanagedCellarAddress, "id: %s", p.TargetContractAddress)
	}

	return nil
}

func HandleCommunityPoolSpendProposal(ctx sdk.Context, k Keeper, p types.AxelarCommunityPoolSpendProposal) error {
	params := k.GetParamSet(ctx)
	config, ok := k.GetChainConfigurationByID(ctx, p.ChainId)
	if !ok {
		return fmt.Errorf("chain by id %d not found", p.ChainId)
	}

	feeFound, feeCoin := config.BridgeFees.Find(p.Amount.Denom)
	if !feeFound {
		return fmt.Errorf("no matching bridge fee for denom %s", p.Amount.Denom)
	}

	coinWithBridgeFee := p.Amount.Add(feeCoin)
	feePool := k.distributionKeeper.GetFeePool(ctx)

	// NOTE the community pool isn't a module account, however its coins
	// are held in the distribution module account. Thus, the community pool
	// must be reduced separately from the Axelar IBC calls
	newPool, negative := feePool.CommunityPool.SafeSub(sdk.NewDecCoinsFromCoins(coinWithBridgeFee))
	if negative {
		return distributiontypes.ErrBadDistribution
	}

	feePool.CommunityPool = newPool

	// since distribution is not an authorized sender, put them in the axelarcork module account
	if err := k.bankKeeper.SendCoinsFromModuleToModule(ctx, distributiontypes.ModuleName, types.ModuleName, sdk.NewCoins(coinWithBridgeFee)); err != nil {
		panic(err)
	}

	sender := k.GetSenderAccount(ctx).GetAddress().String()

	axelarMemo := types.AxelarBody{
		DestinationChain:   config.Name,
		DestinationAddress: p.Recipient,
		Payload:            nil,
		Type:               types.PureTokenTransfer,
	}
	memoBz, err := json.Marshal(axelarMemo)
	if err != nil {
		return err
	}

	memo := string(memoBz)
	transferMsg := transfertypes.NewMsgTransfer(
		params.IbcPort,
		params.IbcChannel,
		coinWithBridgeFee,
		sender,
		params.GmpAccount,
		clienttypes.ZeroHeight(),
		uint64(ctx.BlockTime().Add(time.Duration(params.TimeoutDuration)).UnixNano()),
		memo,
	)
	resp, err := k.transferKeeper.Transfer(sdk.WrapSDKContext(ctx), transferMsg)
	if err != nil {
		return err
	}

	k.distributionKeeper.SetFeePool(ctx, feePool)
	k.Logger(ctx).Info("transfer from the community pool issued to the axelar bridge",
		"ibc sequence", resp,
		"amount", coinWithBridgeFee.Amount.String(),
		"recipient", p.Recipient,
		"chain", config.Name,
		"sender", sender,
		"timeout duration", params.TimeoutDuration,
	)

	return nil
}

// HandleAddChainConfigurationProposal is a handler for executing a passed chain configuration addition proposal
func HandleAddChainConfigurationProposal(ctx sdk.Context, k Keeper, p types.AddChainConfigurationProposal) error {
	err := p.ChainConfiguration.ValidateBasic()
	if err != nil {
		return err
	}

	k.SetChainConfiguration(ctx, p.ChainConfiguration.Id, *p.ChainConfiguration)

	return nil
}

// HandleRemoveChainConfigurationProposal is a handler for executing a passed chain configuration removal proposal
func HandleRemoveChainConfigurationProposal(ctx sdk.Context, k Keeper, p types.RemoveChainConfigurationProposal) error {
	_, ok := k.GetChainConfigurationByID(ctx, p.ChainId)
	if !ok {
		return fmt.Errorf("chain by id %d not found", p.ChainId)
	}

	k.DeleteChainConfigurationByID(ctx, p.ChainId)

	return nil
}

// HandleUpgradeAxelarProxyContractProposal is a handler for executing a passed axelar proxy contract upgrade proposal
func HandleUpgradeAxelarProxyContractProposal(ctx sdk.Context, k Keeper, p types.UpgradeAxelarProxyContractProposal) error {
	_, ok := k.GetChainConfigurationByID(ctx, p.ChainId)
	if !ok {
		return fmt.Errorf("chain by id %d not found", p.ChainId)
	}

	cellars := []string{}
	for _, c := range k.GetCellarIDs(ctx, p.ChainId) {
		cellars = append(cellars, c.String())
	}

	payload, err := types.EncodeUpgradeArgs(p.NewProxyAddress, cellars)
	if err != nil {
		return err
	}

	upgradeData := types.AxelarUpgradeData{
		ChainId:                   p.ChainId,
		Payload:                   payload,
		ExecutableHeightThreshold: ctx.BlockHeight() + int64(types.DefaultExecutableHeightThreshold),
	}

	k.SetAxelarProxyUpgradeData(ctx, p.ChainId, upgradeData)

	return nil
}

// HandleCancelAxelarProxyContractUpgradeProposal is a handler for deleting axelar proxy contract upgrade data before the
// executable height threshold is reached.
func HandleCancelAxelarProxyContractUpgradeProposal(ctx sdk.Context, k Keeper, p types.CancelAxelarProxyContractUpgradeProposal) error {
	_, ok := k.GetChainConfigurationByID(ctx, p.ChainId)
	if !ok {
		return fmt.Errorf("chain id %d not found", p.ChainId)
	}

	k.DeleteAxelarProxyUpgradeData(ctx, p.ChainId)

	return nil
}
