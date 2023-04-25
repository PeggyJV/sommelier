package interchain_tests

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/icza/dyno"
	"go.uber.org/zap/zaptest"
	"os"
	"strings"
	"testing"
	"time"

	ibctest "github.com/strangelove-ventures/interchaintest/v3"
	"github.com/strangelove-ventures/interchaintest/v3/chain/cosmos"
	"github.com/strangelove-ventures/interchaintest/v3/ibc"
	"github.com/strangelove-ventures/interchaintest/v3/testreporter"
	"github.com/strangelove-ventures/interchaintest/v3/testutil"
	"github.com/stretchr/testify/require"
)

const (
	votingPeriod     = "10s"
	maxDepositPeriod = "10s"
)

func TestAxelarCorkSubmissions(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping in short mode")
	}

	t.Parallel()

	ctx := context.Background()

	var numFullNodes = 1
	var numValidators = 3

	// pulling image from env to foster local dev
	imageTag := os.Getenv("SOMMELIER_IMAGE")
	imageTagComponents := strings.Split(imageTag, ":")

	// disabling seeds in axelar because it causes intermittent test failures
	axlConfigFileOverrides := make(map[string]any)
	axlConfigTomlOverrides := make(testutil.Toml)

	axlP2POverrides := make(testutil.Toml)
	axlP2POverrides["seeds"] = ""
	axlConfigTomlOverrides["p2p"] = axlP2POverrides

	axlConfigFileOverrides["config/config.toml"] = axlConfigTomlOverrides

	// Chain factory
	cf := ibctest.NewBuiltinChainFactory(zaptest.NewLogger(t), []*ibctest.ChainSpec{
		{
			Name:    "axelar",
			Version: "v14.0.0",
			ChainConfig: ibc.ChainConfig{
				Images: []ibc.DockerImage{
					{
						Repository: "ghcr.io/strangelove-ventures/heighliner/axelar",
						Version:    "v14.0.0",
						UidGid:     "1025:1025",
					},
				},
				Type:                "cosmos",
				Bin:                 "axelard",
				Bech32Prefix:        "axl",
				Denom:               "uaxl",
				GasPrices:           "0.0uaxl",
				GasAdjustment:       1.3,
				TrustingPeriod:      "336h",
				NoHostMount:         false,
				ConfigFileOverrides: axlConfigFileOverrides,
			},
			NumValidators: &numValidators,
			NumFullNodes:  &numFullNodes,
		},
		{
			Name:    imageTagComponents[0],
			Version: imageTagComponents[1],
			ChainConfig: ibc.ChainConfig{
				Images: []ibc.DockerImage{
					{
						Repository: imageTagComponents[0],
						Version:    imageTagComponents[1],
						UidGid:     "1025:1025",
					},
				},
				GasPrices:      "0.0usomm",
				GasAdjustment:  1.3,
				Type:           "cosmos",
				ChainID:        "sommelier-1",
				Bin:            "sommelierd",
				Bech32Prefix:   "somm",
				Denom:          "usomm",
				TrustingPeriod: "336h",
				ModifyGenesis:  modifyGenesisShortProposals(votingPeriod, maxDepositPeriod),
			},
			NumValidators: &numValidators,
			NumFullNodes:  &numFullNodes,
		},
	})

	chains, err := cf.Chains(t.Name())
	require.NoError(t, err)

	axelar, sommelier := chains[0].(*cosmos.CosmosChain), chains[1].(*cosmos.CosmosChain)

	// Relayer Factory
	client, network := ibctest.DockerSetup(t)
	relayer := ibctest.NewBuiltinRelayerFactory(ibc.CosmosRly, zaptest.NewLogger(t)).Build(
		t, client, network)

	// Prep Interchain
	const ibcPath = "sommelier-axl-cork-test"
	ic := ibctest.NewInterchain().
		AddChain(sommelier).
		AddChain(axelar).
		AddRelayer(relayer, "relayer").
		AddLink(ibctest.InterchainLink{
			Chain1:  sommelier,
			Chain2:  axelar,
			Relayer: relayer,
			Path:    ibcPath,
		})

	// Log location
	f, err := ibctest.CreateLogFile(fmt.Sprintf("%d.json", time.Now().Unix()))
	require.NoError(t, err)
	// Reporter/logs
	rep := testreporter.NewReporter(f)
	eRep := rep.RelayerExecReporter(t)

	// Build Interchain
	require.NoError(t, ic.Build(ctx, eRep, ibctest.InterchainBuildOptions{
		TestName:          t.Name(),
		Client:            client,
		NetworkID:         network,
		BlockDatabaseFile: ibctest.DefaultBlockDatabaseFilepath(),

		SkipPathCreation: false},
	),
	)

	// Create and Fund User Wallets
	t.Log("creating and funding user accounts")
	fundAmount := int64(10_000_000)
	users := ibctest.GetAndFundTestUsers(t, ctx, "default", fundAmount, sommelier, axelar)
	sommelierUser := users[0]
	axelarUser := users[1]
	t.Logf("created sommelier user %s", sommelierUser.Address)
	t.Logf("created axelar user %s", axelarUser.Address)

	sommelierUserBalInitial, err := sommelier.GetBalance(ctx, sommelierUser.Address, sommelier.Config().Denom)
	require.NoError(t, err)
	require.Equal(t, fundAmount, sommelierUserBalInitial)

	// Get Channel ID
	t.Log("getting IBC channel IDs")
	sommelierChannelInfo, err := relayer.GetChannels(ctx, eRep, sommelier.Config().ChainID)
	require.NoError(t, err)
	sommelierChannelID := sommelierChannelInfo[0].ChannelID

	axlChannelInfo, err := relayer.GetChannels(ctx, eRep, axelar.Config().ChainID)
	require.NoError(t, err)
	axlChannelID := axlChannelInfo[0].ChannelID

}

func modifyGenesisShortProposals(votingPeriod string, maxDepositPeriod string) func(ibc.ChainConfig, []byte) ([]byte, error) {
	return func(chainConfig ibc.ChainConfig, genbz []byte) ([]byte, error) {
		g := make(map[string]interface{})
		if err := json.Unmarshal(genbz, &g); err != nil {
			return nil, fmt.Errorf("failed to unmarshal genesis file: %w", err)
		}
		if err := dyno.Set(g, votingPeriod, "app_state", "gov", "params", "voting_period"); err != nil {
			return nil, fmt.Errorf("failed to set voting period in genesis json: %w", err)
		}
		if err := dyno.Set(g, maxDepositPeriod, "app_state", "gov", "params", "max_deposit_period"); err != nil {
			return nil, fmt.Errorf("failed to set voting period in genesis json: %w", err)
		}
		if err := dyno.Set(g, chainConfig.Denom, "app_state", "gov", "params", "min_deposit", 0, "denom"); err != nil {
			return nil, fmt.Errorf("failed to set voting period in genesis json: %w", err)
		}
		if err := dyno.Set(g, "100", "app_state", "gov", "params", "min_deposit", 0, "amount"); err != nil {
			return nil, fmt.Errorf("failed to set voting period in genesis json: %w", err)
		}
		out, err := json.Marshal(g)
		if err != nil {
			return nil, fmt.Errorf("failed to marshal genesis bytes to json: %w", err)
		}
		return out, nil
	}
}
