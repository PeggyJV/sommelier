package keeper

import (
	"fmt"
	"sort"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/ethereum/go-ethereum/common"
	"github.com/peggyjv/sommelier/v9/x/axelarcork/types"
)

// InitGenesis initialize default parameters
// and the keeper's address to pubkey map
func InitGenesis(ctx sdk.Context, k Keeper, gs types.GenesisState) {
	k.SetParams(ctx, *gs.Params)

	senderAccount := k.GetSenderAccount(ctx)
	if senderAccount == nil {
		panic(fmt.Sprintf("%s module account has not been set", types.ModuleName))
	}
	k.accountKeeper.SetModuleAccount(ctx, senderAccount)

	for _, config := range gs.ChainConfigurations.Configurations {
		k.SetChainConfiguration(ctx, config.Id, *config)
	}

	for _, cellarIDSet := range gs.CellarIds {
		k.SetCellarIDs(ctx, cellarIDSet.ChainId, *cellarIDSet)
	}

	for _, corkResult := range gs.CorkResults.CorkResults {
		k.SetAxelarCorkResult(
			ctx,
			corkResult.Cork.ChainId,
			corkResult.Cork.IDHash(corkResult.BlockHeight),
			*corkResult,
		)
	}

	for _, scheduledCork := range gs.ScheduledCorks.ScheduledCorks {
		valAddr, err := sdk.ValAddressFromBech32(scheduledCork.Validator)
		if err != nil {
			panic(err)
		}

		k.SetScheduledAxelarCork(ctx, scheduledCork.Cork.ChainId, scheduledCork.BlockHeight, valAddr, *scheduledCork.Cork)
	}

	// TODO(bolten): not a huge risk since they can be re-sent, but the genesis state is missing WinningAxelarCorks
	for _, n := range gs.AxelarContractCallNonces {
		if _, found := k.GetChainConfigurationByID(ctx, n.ChainId); !found {
			panic(fmt.Sprintf("chain configuration %d not found", n.ChainId))
		}

		if !common.IsHexAddress(n.ContractAddress) {
			panic(fmt.Sprintf("invalid contract address %s", n.ContractAddress))
		}

		k.SetAxelarContractCallNonce(ctx, n.ChainId, n.ContractAddress, n.Nonce)
	}

	for _, ud := range gs.AxelarUpgradeData {
		if _, found := k.GetChainConfigurationByID(ctx, ud.ChainId); !found {
			panic(fmt.Sprintf("chain configuration %d not found for upgrade data", ud.ChainId))
		}

		k.SetAxelarProxyUpgradeData(ctx, ud.ChainId, *ud)
	}
}

// ExportGenesis writes the current store values
// to a genesis file, which can be imported again
// with InitGenesis
func ExportGenesis(ctx sdk.Context, k Keeper) types.GenesisState {
	var gs types.GenesisState

	// These fields are nil by default, so we need to initialize them
	ps := k.GetParamSet(ctx)
	gs.Params = &ps
	gs.CorkResults = &types.AxelarCorkResults{}
	gs.ScheduledCorks = &types.ScheduledAxelarCorks{}

	k.IterateChainConfigurations(ctx, func(config types.ChainConfiguration) (stop bool) {
		gs.ChainConfigurations.Configurations = append(gs.ChainConfigurations.Configurations, &config)

		cellarIDs := k.GetCellarIDs(ctx, config.Id)
		var cellarIDSet types.CellarIDSet
		cellarIDSet.ChainId = config.Id
		cellarIDSetIDs := make([]string, 0, len(cellarIDs))
		for _, id := range cellarIDs {
			cellarIDSetIDs = append(cellarIDSetIDs, id.String())
		}
		sort.Strings(cellarIDSetIDs)
		cellarIDSet.Ids = cellarIDSetIDs
		gs.CellarIds = append(gs.CellarIds, &cellarIDSet)

		gs.ScheduledCorks.ScheduledCorks = append(gs.ScheduledCorks.ScheduledCorks, k.GetScheduledAxelarCorks(ctx, config.Id)...)
		gs.CorkResults.CorkResults = append(gs.CorkResults.CorkResults, k.GetAxelarCorkResults(ctx, config.Id)...)

		return false
	})

	k.IterateAxelarContractCallNonces(ctx, func(chainID uint64, contractAddress common.Address, nonce uint64) (stop bool) {
		accn := &types.AxelarContractCallNonce{
			ChainId:         chainID,
			ContractAddress: contractAddress.String(),
			Nonce:           nonce,
		}

		gs.AxelarContractCallNonces = append(gs.AxelarContractCallNonces, accn)

		return false
	})

	k.IterateAxelarProxyUpgradeData(ctx, func(chainID uint64, upgradeData types.AxelarUpgradeData) (stop bool) {
		gs.AxelarUpgradeData = append(gs.AxelarUpgradeData, &upgradeData)

		return false
	})

	return gs
}
