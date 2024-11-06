package keeper

import (
	"os"

	"cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/peggyjv/sommelier/v7/x/axelarcork/types"
	"github.com/stretchr/testify/require"
	"github.com/tendermint/tendermint/libs/log"
)

func (suite *KeeperTestSuite) TestExportGenesis() {
	logger := log.NewTMLogger(log.NewSyncWriter(os.Stdout))
	k, ctx := suite.axelarcorkKeeper, suite.ctx
	ctx = ctx.WithLogger(logger)
	t := suite.T()

	params := types.DefaultParams()
	k.SetParams(ctx, params)

	suite.Run("exports default genesis state", func() {
		// Export with default state
		genesis := ExportGenesis(ctx, k)

		require.NotNil(t, genesis)
		require.NotNil(t, genesis.Params, "Params should not be nil")

		require.NotNil(t, genesis.ScheduledCorks, "ScheduledCorks should not be nil")
		require.NotNil(t, genesis.CorkResults, "CorkResults should not be nil")

		require.NoError(t, genesis.Validate())
		require.Equal(t, types.DefaultParams(), *genesis.Params)
		require.Empty(t, genesis.ChainConfigurations.Configurations)
		require.Empty(t, genesis.CellarIds)
		require.Empty(t, genesis.ScheduledCorks.ScheduledCorks)
		require.Empty(t, genesis.CorkResults.CorkResults)
		require.Empty(t, genesis.AxelarContractCallNonces)
		require.Empty(t, genesis.AxelarUpgradeData)
	})

	suite.Run("exports state with data", func() {
		// Set up test data
		params := types.DefaultParams()
		k.SetParams(ctx, params)

		// Set multiple chain configurations
		chainConfigs := []*types.ChainConfiguration{
			{
				Name:         "test-chain-1",
				Id:           1,
				ProxyAddress: "0x0000000000000000000000000000000000000000",
				BridgeFees:   sdk.NewCoins(sdk.NewCoin("uatom", math.NewInt(1000000))),
			},
			{
				Name:         "test-chain-2",
				Id:           2,
				ProxyAddress: "0x0000000000000000000000000000000000000000",
				BridgeFees:   sdk.NewCoins(sdk.NewCoin("uatom", math.NewInt(2000000))),
			},
		}
		for _, cc := range chainConfigs {
			k.SetChainConfiguration(ctx, cc.Id, *cc)
		}

		// Set multiple cellar ID sets
		cellarIds := []*types.CellarIDSet{
			{
				ChainId: 1,
				Ids:     []string{"0x1111111111111111111111111111111111111111", "0x2222222222222222222222222222222222222222"},
			},
			{
				ChainId: 2,
				Ids:     []string{"0x3333333333333333333333333333333333333333", "0x4444444444444444444444444444444444444444"},
			},
		}
		for _, cid := range cellarIds {
			k.SetCellarIDs(ctx, cid.ChainId, *cid)
		}

		// Set multiple scheduled corks
		scheduledCorks := []*types.ScheduledAxelarCork{
			{
				Cork: &types.AxelarCork{
					EncodedContractCall:   []byte{0x1},
					ChainId:               1,
					TargetContractAddress: "0x0000000000000000000000000000000000000000",
					Deadline:              100,
				},
				BlockHeight: 100,
				Validator:   "cosmosvaloper1weskc6tyv96x7u3dxyqqqqqqqqqqqqqqr0fqts",
				Id:          "15549d518be7ad8d29e2af15adfb38e24c50d76e5540e1efb5603d17c9959518",
			},
			{
				Cork: &types.AxelarCork{
					EncodedContractCall:   []byte{0x2},
					ChainId:               2,
					TargetContractAddress: "0x0000000000000000000000000000000000000000",
					Deadline:              200,
				},
				BlockHeight: 200,
				Validator:   "cosmosvaloper1weskc6tyv96x7u3dxgqqqqqqqqqqqqqqqpd58s",
				Id:          "eb8eb8bbef7351946a6716e81e41beeb207470e0c97109056810aa5c539f06cf",
			},
		}
		for _, sc := range scheduledCorks {
			valAddr, err := sdk.ValAddressFromBech32(sc.Validator)
			require.NoError(t, err)
			k.SetScheduledAxelarCork(ctx, sc.Cork.ChainId, sc.BlockHeight, valAddr, *sc.Cork)
		}

		// Set multiple cork results
		corkResults := []*types.AxelarCorkResult{
			{
				Cork: &types.AxelarCork{
					EncodedContractCall:   []byte{0x1},
					ChainId:               1,
					TargetContractAddress: "0x0000000000000000000000000000000000000000",
					Deadline:              100,
				},
				BlockHeight:        100,
				Approved:           true,
				ApprovalPercentage: "100.0",
			},
			{
				Cork: &types.AxelarCork{
					EncodedContractCall:   []byte{0x2},
					ChainId:               2,
					TargetContractAddress: "0x0000000000000000000000000000000000000000",
					Deadline:              200,
				},
				BlockHeight:        200,
				Approved:           false,
				ApprovalPercentage: "0.0",
			},
		}
		for _, cr := range corkResults {
			k.SetAxelarCorkResult(ctx, cr.Cork.ChainId, cr.Cork.EncodedContractCall, *cr)
		}

		// Set multiple contract call nonces
		nonces := []*types.AxelarContractCallNonce{
			{
				ChainId:         1,
				Nonce:           1,
				ContractAddress: "0x0000000000000000000000000000000000000000",
			},
			{
				ChainId:         2,
				Nonce:           2,
				ContractAddress: "0x0000000000000000000000000000000000000000",
			},
		}
		for _, n := range nonces {
			k.SetAxelarContractCallNonce(ctx, n.ChainId, n.ContractAddress, n.Nonce)
		}

		// Set multiple upgrade data entries
		upgradeData := []*types.AxelarUpgradeData{
			{
				ChainId:                   1,
				Payload:                   []byte{0x1},
				ExecutableHeightThreshold: 100,
			},
			{
				ChainId:                   2,
				Payload:                   []byte{0x2},
				ExecutableHeightThreshold: 200,
			},
		}
		for _, ud := range upgradeData {
			k.SetAxelarProxyUpgradeData(ctx, ud.ChainId, *ud)
		}

		// Export genesis
		genesis := ExportGenesis(ctx, k)

		// Verify exported state
		require.NotNil(t, genesis)
		require.NotNil(t, genesis.Params, "Params should not be nil")
		require.NotNil(t, genesis.ScheduledCorks, "ScheduledCorks should not be nil")
		require.NotNil(t, genesis.CorkResults, "CorkResults should not be nil")

		require.NoError(t, genesis.Validate())
		require.Equal(t, params, *genesis.Params)
		require.Len(t, genesis.ChainConfigurations.Configurations, 2)
		require.ElementsMatch(t, chainConfigs, genesis.ChainConfigurations.Configurations)
		require.Len(t, genesis.CellarIds, 2)
		require.ElementsMatch(t, cellarIds, genesis.CellarIds)
		require.Len(t, genesis.ScheduledCorks.ScheduledCorks, 2)
		require.ElementsMatch(t, scheduledCorks, genesis.ScheduledCorks.ScheduledCorks)
		require.Len(t, genesis.CorkResults.CorkResults, 2)
		require.ElementsMatch(t, corkResults, genesis.CorkResults.CorkResults)
		require.Len(t, genesis.AxelarContractCallNonces, 2)
		require.ElementsMatch(t, nonces, genesis.AxelarContractCallNonces)
		require.Len(t, genesis.AxelarUpgradeData, 2)
		require.ElementsMatch(t, upgradeData, genesis.AxelarUpgradeData)
	})
}
