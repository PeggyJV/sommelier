package integration_tests

import (
	"encoding/json"
	"fmt"
	"os"

	tmtypes "github.com/cometbft/cometbft/types"
	"github.com/cosmos/cosmos-sdk/server"
	sdk "github.com/cosmos/cosmos-sdk/types"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	"github.com/cosmos/cosmos-sdk/x/genutil"
	genutiltypes "github.com/cosmos/cosmos-sdk/x/genutil/types"
	"github.com/peggyjv/sommelier/v10/app/params"
	incentivestypes "github.com/peggyjv/sommelier/v10/x/incentives/types"
	poatypes "github.com/peggyjv/sommelier/v10/x/poa/types"
)

func getGenDoc(path string) (*tmtypes.GenesisDoc, error) {
	serverCtx := server.NewDefaultContext()
	config := serverCtx.Config
	config.SetRoot(path)

	genFile := config.GenesisFile()
	doc := &tmtypes.GenesisDoc{}

	if _, err := os.Stat(genFile); err != nil {
		if !os.IsNotExist(err) {
			return nil, err
		}
	} else {
		var err error

		doc, err = tmtypes.GenesisDocFromFile(genFile)
		if err != nil {
			return nil, fmt.Errorf("failed to read genesis doc from file: %w", err)
		}
	}

	return doc, nil
}

func addGenesisAccount(path, moniker, amountStr string, accAddr sdk.AccAddress) error { //nolint:unparam
	serverCtx := server.NewDefaultContext()
	config := serverCtx.Config

	config.SetRoot(path)
	config.Moniker = moniker

	coins, err := sdk.ParseCoinsNormalized(amountStr)
	if err != nil {
		return fmt.Errorf("failed to parse coins: %w", err)
	}

	balances := banktypes.Balance{Address: accAddr.String(), Coins: coins.Sort()}
	genAccount := authtypes.NewBaseAccount(accAddr, nil, 0, 0)

	genFile := config.GenesisFile()
	appState, genDoc, err := genutiltypes.GenesisStateFromGenFile(genFile)
	if err != nil {
		return fmt.Errorf("failed to unmarshal genesis state: %w", err)
	}

	authGenState := authtypes.GetGenesisStateFromAppState(cdc, appState)

	accs, err := authtypes.UnpackAccounts(authGenState.Accounts)
	if err != nil {
		return fmt.Errorf("failed to get accounts from any: %w", err)
	}

	if accs.Contains(accAddr) {
		return fmt.Errorf("failed to add account to genesis state; account already exists: %s", accAddr)
	}

	// Add the new account to the set of genesis accounts and sanitize the
	// accounts afterwards.
	accs = append(accs, genAccount)
	accs = authtypes.SanitizeGenesisAccounts(accs)

	genAccs, err := authtypes.PackAccounts(accs)
	if err != nil {
		return fmt.Errorf("failed to convert accounts into any's: %w", err)
	}

	authGenState.Accounts = genAccs

	authGenStateBz, err := cdc.MarshalJSON(&authGenState)
	if err != nil {
		return fmt.Errorf("failed to marshal auth genesis state: %w", err)
	}

	appState[authtypes.ModuleName] = authGenStateBz

	bankGenState := banktypes.GetGenesisStateFromAppState(cdc, appState)
	bankGenState.Balances = append(bankGenState.Balances, balances)
	bankGenState.Balances = banktypes.SanitizeGenesisBalances(bankGenState.Balances)

	bankGenStateBz, err := cdc.MarshalJSON(bankGenState)
	if err != nil {
		return fmt.Errorf("failed to marshal bank genesis state: %w", err)
	}

	appState[banktypes.ModuleName] = bankGenStateBz

	appStateJSON, err := json.Marshal(appState)
	if err != nil {
		return fmt.Errorf("failed to marshal application genesis state: %w", err)
	}

	genDoc.AppState = appStateJSON
	return genutil.ExportGenesisFile(genDoc, genFile)
}

func (s *IntegrationTestSuite) setIncentivesGenState(appGenState map[string]json.RawMessage) error {
	incentivesGenState := incentivestypes.DefaultGenesisState()
	err := cdc.UnmarshalJSON(appGenState[incentivestypes.ModuleName], &incentivesGenState)
	if err != nil {
		return fmt.Errorf("failed to unmarshal incentives genesis state: %w", err)
	}

	incentivesGenState.Params.ValidatorIncentivesCutoffHeight = 0
	incentivesGenState.Params.ValidatorMaxDistributionPerBlock = sdk.NewCoin(params.BaseCoinUnit, sdk.NewInt(0))

	appGenState[incentivestypes.ModuleName] = cdc.MustMarshalJSON(&incentivesGenState)
	return nil
}

// setPoaGenState seeds the x/poa authority allowlist with every validator in
// the test chain.
//
// Without this the module's DefaultGenesis leaves the allowlist EMPTY, the
// first EndBlocker enters authority-empty safe mode, and the whole
// value-bearing surface freezes: gravity's BeginBlock/EndBlock become no-ops
// and MsgSendToEthereum / MsgSubmitEthereumEvent /
// MsgSubmitEthereumTxConfirmation / MsgScheduleCork are all rejected. Every
// bridge and cork scenario in this suite would fail for reasons that have
// nothing to do with what it is testing.
//
// All test validators are authority validators, so the boost multiplier clamps
// to 1 and consensus power is unchanged from a pre-PoA run. Tests that want to
// exercise the community bucket (or safe mode) should narrow this set.
func (s *IntegrationTestSuite) setPoaGenState(appGenState map[string]json.RawMessage) error {
	poaGenState := poatypes.DefaultGenesis()
	if err := cdc.UnmarshalJSON(appGenState[poatypes.ModuleName], poaGenState); err != nil {
		return fmt.Errorf("failed to unmarshal poa genesis state: %w", err)
	}

	authoritySet := make([]poatypes.AuthorityValidator, 0, len(s.chain.validators))
	for _, val := range s.chain.validators {
		authoritySet = append(authoritySet, poatypes.AuthorityValidator{
			OperatorAddress: val.validatorAddress().String(),
		})
	}
	poaGenState.AuthoritySet = authoritySet

	if err := poaGenState.Validate(); err != nil {
		return fmt.Errorf("invalid poa genesis state: %w", err)
	}

	appGenState[poatypes.ModuleName] = cdc.MustMarshalJSON(poaGenState)
	return nil
}
