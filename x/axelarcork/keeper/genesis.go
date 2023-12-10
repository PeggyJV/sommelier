package keeper

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/ethereum/go-ethereum/common"
	"github.com/peggyjv/sommelier/v7/x/axelarcork/types"
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

	for i, config := range gs.ChainConfigurations.Configurations {
		k.SetChainConfiguration(ctx, config.Id, *config)
		k.SetCellarIDs(ctx, config.Id, *gs.CellarIds[i])
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

	ps := k.GetParamSet(ctx)
	gs.Params = &ps

	k.IterateChainConfigurations(ctx, func(config types.ChainConfiguration) (stop bool) {
		gs.ChainConfigurations.Configurations = append(gs.ChainConfigurations.Configurations, &config)

		cellarIDs := k.GetCellarIDs(ctx, config.Id)
		var cellarIDSet types.CellarIDSet
		for _, id := range cellarIDs {
			cellarIDSet.Ids = append(cellarIDSet.Ids, id.String())
		}
		gs.CellarIds = append(gs.CellarIds, &cellarIDSet)

		gs.ScheduledCorks.ScheduledCorks = append(gs.ScheduledCorks.ScheduledCorks, k.GetScheduledAxelarCorks(ctx, config.Id)...)
		gs.CorkResults.CorkResults = append(gs.CorkResults.CorkResults, k.GetAxelarCorkResults(ctx, config.Id)...)

		return false
	})

	k.IterateAxelarContractCallNonces(ctx, func(chainID uint64, contractAddress common.Address, nonce uint64) (stop bool) {
		accn := &types.AxelarContractCallNonce{
			ChainId:         chainID,
			ContractAddress: contractAddress.Hex(),
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
