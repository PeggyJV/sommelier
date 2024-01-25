package keeper

import (
	"encoding/json"
	"fmt"
	"sort"
	"time"

	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	distributiontypes "github.com/cosmos/cosmos-sdk/x/distribution/types"
	transfertypes "github.com/cosmos/ibc-go/v6/modules/apps/transfer/types"
	clienttypes "github.com/cosmos/ibc-go/v6/modules/core/02-client/types"

	errorsmod "cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/ethereum/go-ethereum/common"
	"github.com/peggyjv/sommelier/v7/x/axelarcork/types"
)

// HandleAddManagedCellarsProposal is a handler for executing a passed community cellar addition proposal
func HandleAddManagedCellarsProposal(ctx sdk.Context, k Keeper, p types.AddAxelarManagedCellarIDsProposal) error {
	config, ok := k.GetChainConfigurationByID(ctx, p.ChainId)
	if !ok {
		return fmt.Errorf("chain by id %d not found", p.ChainId)
	}

	if err := p.CellarIds.ValidateBasic(); err != nil {
		return err
	}

	cellarIDs := k.GetCellarIDs(ctx, config.Id)

	for _, proposedCellarID := range p.CellarIds.Ids {
		found := false
		for _, id := range cellarIDs {
			if id == common.HexToAddress(proposedCellarID) {
				found = true
			}
		}
		if !found {
			cellarIDs = append(cellarIDs, common.HexToAddress(proposedCellarID))
		}
	}

	idStrings := make([]string, len(cellarIDs))
	for i, cid := range cellarIDs {
		idStrings[i] = cid.String()
	}

	sort.Strings(idStrings)
	k.SetCellarIDs(ctx, config.Id, types.CellarIDSet{Ids: idStrings})

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
			outputCellarIDs.Ids = append(outputCellarIDs.Ids, existingID.Hex())
		}
	}
	outputCellarIDs.ChainId = config.Id

	// unlike for adding an ID, we don't need to re-sort because we're removing elements from an already sorted list
	k.SetCellarIDs(ctx, config.Id, outputCellarIDs)

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
	feePool := k.distributionKeeper.GetFeePool(ctx)

	// NOTE the community pool isn't a module account, however its coins
	// are held in the distribution module account. Thus, the community pool
	// must be reduced separately from the Axelar IBC calls
	newPool, negative := feePool.CommunityPool.SafeSub(sdk.NewDecCoinsFromCoins(p.Amount))
	if negative {
		return distributiontypes.ErrBadDistribution
	}

	feePool.CommunityPool = newPool
	sender := authtypes.NewModuleAddress(distributiontypes.ModuleName)

	params := k.GetParamSet(ctx)
	config, ok := k.GetChainConfigurationByID(ctx, p.ChainId)
	if !ok {
		return fmt.Errorf("chain by id %d not found", p.ChainId)
	}

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
		p.Amount,
		sender.String(),
		p.Recipient,
		clienttypes.ZeroHeight(),
		uint64(ctx.BlockTime().Add(time.Duration(params.TimeoutDuration)).UnixNano()),
		memo,
	)
	resp, err := k.transferKeeper.Transfer(ctx.Context(), transferMsg)
	if err != nil {
		return err
	}

	k.distributionKeeper.SetFeePool(ctx, feePool)
	k.Logger(ctx).Info("transfer from the community pool issued to the axelar bridge",
		"ibc sequence", resp,
		"amount", p.Amount.String(),
		"recipient", p.Recipient,
		"chain", config.Name,
		"sender", sender.String(),
		"timeout duration", params.TimeoutDuration,
	)

	return nil
}

// HandleAddChainConfigurationProposal is a handler for executing a passed chain configuration addition proposal
func HandleAddChainConfigurationProposal(ctx sdk.Context, k Keeper, p types.AddChainConfigurationProposal) error {
	k.SetChainConfiguration(ctx, p.ChainConfiguration.Id, *p.ChainConfiguration)

	return nil
}

// HandleRemoveChainConfigurationProposal is a handler for executing a passed chain configuration removal proposal
func HandleRemoveChainConfigurationProposal(ctx sdk.Context, k Keeper, p types.RemoveChainConfigurationProposal) error {
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
		cellars = append(cellars, c.Hex())
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
