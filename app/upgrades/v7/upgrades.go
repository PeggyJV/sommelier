package v7

import (
	"fmt"
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/module"
	upgradetypes "github.com/cosmos/cosmos-sdk/x/upgrade/types"
	icahostkeeper "github.com/cosmos/ibc-go/v6/modules/apps/27-interchain-accounts/host/keeper"
	icahosttypes "github.com/cosmos/ibc-go/v6/modules/apps/27-interchain-accounts/host/types"
	ibctransfertypes "github.com/cosmos/ibc-go/v6/modules/apps/transfer/types"
	auctionkeeper "github.com/peggyjv/sommelier/v7/x/auction/keeper"
	auctiontypes "github.com/peggyjv/sommelier/v7/x/auction/types"
	axelarcorkkeeper "github.com/peggyjv/sommelier/v7/x/axelarcork/keeper"
	axelarcorktypes "github.com/peggyjv/sommelier/v7/x/axelarcork/types"
	cellarfeeskeeper "github.com/peggyjv/sommelier/v7/x/cellarfees/keeper"
	cellarfeestypes "github.com/peggyjv/sommelier/v7/x/cellarfees/types"
	corkkeeper "github.com/peggyjv/sommelier/v7/x/cork/keeper"
	corktypes "github.com/peggyjv/sommelier/v7/x/cork/types"
	pubsubkeeper "github.com/peggyjv/sommelier/v7/x/pubsub/keeper"
	pubsubtypes "github.com/peggyjv/sommelier/v7/x/pubsub/types"
)

func CreateUpgradeHandler(
	mm *module.Manager,
	configurator module.Configurator,
	auctionKeeper auctionkeeper.Keeper,
	axelarcorkKeeper axelarcorkkeeper.Keeper,
	cellarfeesKeeper cellarfeeskeeper.Keeper,
	corkKeeper corkkeeper.Keeper,
	icaHostKeeper icahostkeeper.Keeper,
	pubsubKeeper pubsubkeeper.Keeper,
) upgradetypes.UpgradeHandler {
	return func(ctx sdk.Context, plan upgradetypes.Plan, fromVM module.VersionMap) (module.VersionMap, error) {
		ctx.Logger().Info("v7 upgrade: entering handler")

		// Now that we're on IBC V6, we can update the ICA host module to allow all message types rather than
		// the list we specified in the v6 upgrade -- a default of HostEnabled: true and the string "*" for messages
		ctx.Logger().Info("v7 upgrade: setting ICA host params to allow all messages")
		icaParams := icahosttypes.DefaultParams()
		icaHostKeeper.SetParams(ctx, icaParams)

		// Given that we've removed and added params in cork v2, we'll also set the cork v2 params to their defaults
		ctx.Logger().Info("v7 upgrade: setting cork v2 params to defaults")
		corkParams := corktypes.DefaultParams()
		corkKeeper.SetParams(ctx, corkParams)

		// We must manually run InitGenesis for auction, axelarcork, and pubsub so we can adjust their values
		// during the upgrade process. RunMigrations will migrate to the new cork version. Setting the consensus
		// version to 1 prevents RunMigrations from running InitGenesis itself.
		fromVM[auctiontypes.ModuleName] = mm.Modules[auctiontypes.ModuleName].ConsensusVersion()
		fromVM[axelarcorktypes.ModuleName] = mm.Modules[axelarcorktypes.ModuleName].ConsensusVersion()
		fromVM[pubsubtypes.ModuleName] = mm.Modules[pubsubtypes.ModuleName].ConsensusVersion()

		// Params values were introduced in this upgrade but no migration was necessary, so we initialize the
		// new values to their defaults
		ctx.Logger().Info("v7 upgrading: setting cellarfees default params")
		cellarfeesKeeper.SetParams(ctx, cellarfeestypes.DefaultParams())

		ctx.Logger().Info("v7 upgrade: initializing auction genesis state")
		auctionInitGenesis(ctx, auctionKeeper)

		ctx.Logger().Info("v7 upgrade: intializing axelarcork genesis state")
		axelarcorkInitGenesis(ctx, axelarcorkKeeper)

		ctx.Logger().Info("v7 upgrade: initializing pubsub genesis state")
		pubsubInitGenesis(ctx, pubsubKeeper)

		ctx.Logger().Info("v7 upgrade: running migrations and exiting handler")
		return mm.RunMigrations(ctx, configurator, fromVM)
	}
}

// Initialize the auction module with prices for some stablecoins and SOMM.
func auctionInitGenesis(ctx sdk.Context, auctionKeeper auctionkeeper.Keeper) {
	genesisState := auctiontypes.DefaultGenesisState()

	genesisState.Params.MinimumAuctionHeight = 13020000 // roughly 2024-02-08 09:30 UTC

	usomm52WeekLow := sdk.MustNewDecFromStr("0.079151")
	eth52WeekHigh := sdk.MustNewDecFromStr("2618.33")
	btc52WeekHigh := sdk.MustNewDecFromStr("46936.19")
	oneDollar := sdk.MustNewDecFromStr("1.0")

	var lastUpdatedBlock uint64 = uint64(ctx.BlockHeight())

	usommPrice := auctiontypes.TokenPrice{
		Denom:            "usomm",
		UsdPrice:         usomm52WeekLow,
		Exponent:         6,
		LastUpdatedBlock: lastUpdatedBlock,
	}

	// setting stables to 1 dollar
	usdcPrice := auctiontypes.TokenPrice{
		Denom:            "gravity0xA0b86991c6218b36c1d19D4a2e9Eb0cE3606eB48",
		UsdPrice:         oneDollar,
		Exponent:         6,
		LastUpdatedBlock: lastUpdatedBlock,
	}

	usdtPrice := auctiontypes.TokenPrice{
		Denom:            "gravity0xdAC17F958D2ee523a2206206994597C13D831ec7",
		UsdPrice:         oneDollar,
		Exponent:         6,
		LastUpdatedBlock: lastUpdatedBlock,
	}

	daiPrice := auctiontypes.TokenPrice{
		Denom:            "gravity0x6B175474E89094C44Da98b954EedeAC495271d0F",
		UsdPrice:         oneDollar,
		Exponent:         18,
		LastUpdatedBlock: lastUpdatedBlock,
	}

	fraxPrice := auctiontypes.TokenPrice{
		Denom:            "gravity0x853d955aCEf822Db058eb8505911ED77F175b99e",
		UsdPrice:         oneDollar,
		Exponent:         18,
		LastUpdatedBlock: lastUpdatedBlock,
	}

	// setting non-stables
	wethPrice := auctiontypes.TokenPrice{
		Denom:            "gravity0xC02aaA39b223FE8D0A0e5C4F27eAD9083C756Cc2",
		UsdPrice:         eth52WeekHigh,
		Exponent:         18,
		LastUpdatedBlock: lastUpdatedBlock,
	}

	wbtcPrice := auctiontypes.TokenPrice{
		Denom:            "gravity0x2260FAC5E5542a773Aa44fBCfeDf7C193bc2C599",
		UsdPrice:         btc52WeekHigh,
		Exponent:         8,
		LastUpdatedBlock: lastUpdatedBlock,
	}

	genesisState.TokenPrices = []*auctiontypes.TokenPrice{
		&usommPrice,
		&usdcPrice,
		&usdtPrice,
		&daiPrice,
		&fraxPrice,
		&wethPrice,
		&wbtcPrice,
	}

	if err := genesisState.Validate(); err != nil {
		panic(fmt.Errorf("auction genesis state invalid: %s", err))
	}

	auctionkeeper.InitGenesis(ctx, auctionKeeper, genesisState)
}

// Initialize the Axelar cork module with the correct parameters.
func axelarcorkInitGenesis(ctx sdk.Context, axelarcorkKeeper axelarcorkkeeper.Keeper) {
	genesisState := axelarcorktypes.DefaultGenesisState()

	genesisState.Params.Enabled = true
	genesisState.Params.TimeoutDuration = uint64(6 * time.Hour)
	genesisState.Params.IbcChannel = "channel-5"
	genesisState.Params.IbcPort = ibctransfertypes.PortID
	genesisState.Params.GmpAccount = "axelar1dv4u5k73pzqrxlzujxg3qp8kvc3pje7jtdvu72npnt5zhq05ejcsn5qme5s"
	genesisState.Params.ExecutorAccount = "axelar1aythygn6z5thymj6tmzfwekzh05ewg3l7d6y89"
	genesisState.Params.CorkTimeoutBlocks = 5000

	// pure token transfers have a fixed fee deducted from the amount sent in the ICS-20 message depending
	// on the asset and destination chain
	// they can be calculated here: https://docs.axelar.dev/resources/mainnet#cross-chain-relayer-gas-fee
	genesisState.ChainConfigurations = axelarcorktypes.ChainConfigurations{
		Configurations: []*axelarcorktypes.ChainConfiguration{
			{
				Name:         "arbitrum",
				Id:           42161,
				ProxyAddress: "0xEe75bA2C81C04DcA4b0ED6d1B7077c188FEde4d2",
				BridgeFees:   sdk.NewCoins(sdk.NewCoin("usomm", sdk.NewIntFromUint64(33670000))),
			},
			{
				Name:         "Avalanche",
				Id:           43114,
				ProxyAddress: "0xEe75bA2C81C04DcA4b0ED6d1B7077c188FEde4d2",
				BridgeFees:   sdk.NewCoins(sdk.NewCoin("usomm", sdk.NewIntFromUint64(10670000))),
			},
		},
	}

	genesisState.CellarIds = []*axelarcorktypes.CellarIDSet{
		{
			ChainId: 42161, // arbitrum with test cellar
			Ids:     []string{"0x438087f7c226A89762a791F187d7c3D4a0e95ae6"},
		},
		{
			ChainId: 43114, // Avalanche
			Ids:     []string{},
		},
	}

	if err := genesisState.Validate(); err != nil {
		panic(fmt.Errorf("axelarcork genesis state invalid: %s", err))
	}

	axelarcorkkeeper.InitGenesis(ctx, axelarcorkKeeper, genesisState)
}

// Set up the initial pubsub state to mirror what is currently used in practice already, with 7seas as the
// first publisher using its existing CA certificate, its default subscriptions as the already launched cellars,
// and the subscribers as reflected in the steward-registry repo.
func pubsubInitGenesis(ctx sdk.Context, pubsubKeeper pubsubkeeper.Keeper) {
	genesisState := pubsubtypes.DefaultGenesisState()

	// Initialize the 7seas publisher.
	publisher := pubsubtypes.Publisher{
		Address: "somm14zsm5frvjuqxk3f9377altc6xq368dglhmkxmp",
		Domain:  SevenSeasDomain,
		CaCert:  SevenSeasPublisherCA,
	}
	publishers := []*pubsubtypes.Publisher{&publisher}

	cellars := []string{
		"1:0x7bAD5DF5E11151Dc5Ee1a648800057C5c934c0d5",     // Aave V2
		"1:0x03df2A53Cbed19B824347D6a45d09016C2D1676a",     // DeFi Stars
		"1:0x6c51041A91C91C86f3F08a72cB4D3F67f1208897",     // ETH Trend Growth
		"1:0x6b7f87279982d919Bbf85182DDeAB179B366D8f2",     // ETH-BTC Trend
		"1:0x6E2dAc3b9E9ADc0CbbaE2D0B9Fd81952a8D33872",     // ETH-BTC Momentum
		"1:0xDBe19d1c3F21b1bB250ca7BDaE0687A97B5f77e6",     // Fraximal
		"1:0xcf4B531b4Cde95BD35d71926e09B2b54c564F5b6",     // Max Morpho ETH
		"1:0xC7b69E15D86C5c1581dacce3caCaF5b68cd6596F",     // Real Yield 1INCH
		"1:0x0274a704a6D9129F90A62dDC6f6024b33EcDad36",     // Real Yield BTC
		"1:0x18ea937aba6053bC232d9Ae2C42abE7a8a2Be440",     // Real Yield ENS
		"1:0xb5b29320d2Dde5BA5BAFA1EbcD270052070483ec",     // Real Yield ETH
		"1:0x4068BDD217a45F8F668EF19F1E3A1f043e4c4934",     // Real Yield LINK
		"1:0xcBf2250F33c4161e18D4A2FA47464520Af5216b5",     // Real Yield SNX
		"1:0x6A6AF5393DC23D7e3dB28D28Ef422DB7c40932B6",     // Real Yield UNI
		"1:0x97e6E0a40a3D02F12d1cEC30ebfbAE04e37C119E",     // Real Yield USD
		"1:0x3F07A84eCdf494310D397d24c1C78B041D2fa622",     // Steady ETH
		"1:0x4986fD36b6b16f49b43282Ee2e24C5cF90ed166d",     // Steady BTC
		"1:0x05641a27C82799AaF22b436F20A3110410f29652",     // Steady MATIC
		"1:0x6F069F711281618467dAe7873541EcC082761B33",     // Steady UNI
		"1:0x6a6731f1754e3088ea35c8ddfd55cFF4d0FA5052",     // Turbo CRVUSD
		"1:0x9a7b4980C6F0FCaa50CD5f288Ad7038f434c692e",     // Turbo EETH
		"1:0x0C190DEd9Be5f512Bd72827bdaD4003e9Cc7975C",     // Turbo GHO
		"1:0x5195222f69c5821f8095ec565E71e18aB6A2298f",     // Turbo SOMM
		"1:0xc7372Ab5dd315606dB799246E8aA112405abAeFf",     // Turbo stETH (stETH deposit)
		"1:0xfd6db5011b171B05E1Ea3b92f9EAcaEEb055e971",     // Turbo stETH (WETH deposit)
		"1:0xd33dAd974b938744dAC81fE00ac67cb5AA13958E",     // Turbo swETH
		"1:0x991Fc0B9f027A80E2d642Eb260a41FfC14b2f450",     // Yield Maxi USD
		"42161:0x438087f7c226A89762a791F187d7c3D4a0e95ae6", // Arbitrum test cellar
	}

	// Set 7seas publisher intents for existing cellars
	publisherIntents := make([]*pubsubtypes.PublisherIntent, 0, 25)
	for _, cellar := range cellars {
		publisherIntents = append(publisherIntents, &pubsubtypes.PublisherIntent{
			SubscriptionId:     cellar,
			PublisherDomain:    SevenSeasDomain,
			Method:             pubsubtypes.PublishMethod_PUSH,
			AllowedSubscribers: pubsubtypes.AllowedSubscribers_VALIDATORS,
		})
	}

	// Set default subscriptions for 7seas as the publisher for existing cellars
	defaultSubscriptions := make([]*pubsubtypes.DefaultSubscription, 0, 25)
	for _, cellar := range cellars {
		defaultSubscriptions = append(defaultSubscriptions, &pubsubtypes.DefaultSubscription{
			SubscriptionId:  cellar,
			PublisherDomain: SevenSeasDomain,
		})
	}

	// Create subscribers and intents for existing validators
	subscribers := createSubscribers()
	subscriberIntents := make([]*pubsubtypes.SubscriberIntent, 0, 875)
	for _, subscriber := range subscribers {
		for _, cellar := range cellars {
			subscriberIntents = append(subscriberIntents, &pubsubtypes.SubscriberIntent{
				SubscriptionId:    cellar,
				SubscriberAddress: subscriber.Address,
				PublisherDomain:   SevenSeasDomain,
			})
		}
	}

	genesisState.Publishers = publishers
	genesisState.PublisherIntents = publisherIntents
	genesisState.DefaultSubscriptions = defaultSubscriptions
	genesisState.Subscribers = subscribers
	genesisState.SubscriberIntents = subscriberIntents

	if err := genesisState.Validate(); err != nil {
		panic(fmt.Errorf("pubsub genesis state invalid: %s", err))
	}

	pubsubkeeper.InitGenesis(ctx, pubsubKeeper, genesisState)
}

// Addresses are from the validator's delegate orchestrator key and certs/URLs captured from the
// steward-registry repo.
// query to get orchestrator key: sommelier query gravity delegate-keys-by-validator sommvaloper<rest_of_val_address>
// See source data at: https://github.com/PeggyJV/steward-registry
// data captured at commit ecdb7f386e7e573edb5d8f6ad22a1a67cfa21863
// leaving out made_in_block because I can't find their validator on-chain
// blockhunters hadn't been merged, but verified and added here
func createSubscribers() []*pubsubtypes.Subscriber {
	subscribers := make([]*pubsubtypes.Subscriber, 0, 35)

	subscribers = append(subscribers, &pubsubtypes.Subscriber{
		Address: "somm1s2q8avjykkztudpl8k60f0ns4v5mvnjp5t366c",
		CaCert:  FigmentSubscriberCA,
		PushUrl: "sommelier-steward.staking.production.figment.io:5734",
	})

	subscribers = append(subscribers, &pubsubtypes.Subscriber{
		Address: "somm148a27t9usz9u5xzzjnkt2u8fergs48935dzdnt",
		CaCert:  StandardCryptoSubscriberCA,
		PushUrl: "steward.sommelier.standardcryptovc.com:5734",
	})

	subscribers = append(subscribers, &pubsubtypes.Subscriber{
		Address: "somm1kfzuvra3ym8nxffwdlyj0xvkky87qc0ywh9d42",
		CaCert:  RockawaySubscriberCA,
		PushUrl: "steward-01.rbf.systems:5734",
	})

	subscribers = append(subscribers, &pubsubtypes.Subscriber{
		Address: "somm1tznv6agw8pdzv34ykdpau243kdwyvf9lz4dedh",
		CaCert:  BlockscapeSubscriberCA,
		PushUrl: "steward.blockscape.network:5734",
	})

	subscribers = append(subscribers, &pubsubtypes.Subscriber{
		Address: "somm1p7vn9hajt44fxwn4ecfxjs2r469l0tgmjlqzmp",
		CaCert:  SimplySubscriberCA,
		PushUrl: "sommelier-steward.simply-vc.com.mt:5734",
	})

	subscribers = append(subscribers, &pubsubtypes.Subscriber{
		Address: "somm1lxktamk5tw30cksdlafyzr47vc5cdm76u4tkjm",
		CaCert:  PupmosSubscriberCA,
		PushUrl: "steward.sommelier.pupmos.network:5734",
	})

	subscribers = append(subscribers, &pubsubtypes.Subscriber{
		Address: "somm1yrcqnv7xvfztuh0020vxrnuhgc6dghv3kxvvnk",
		CaCert:  LavenderFiveSubscriberCA,
		PushUrl: "steward.lavenderfive.com:5734",
	})

	subscribers = append(subscribers, &pubsubtypes.Subscriber{
		Address: "somm1y0yvuwnk7g3at6yvl6ctgsvzuxaeqjkw53tduu",
		CaCert:  PolkachuSubscriberCA,
		PushUrl: "steward.polkachu.com:5734",
	})

	subscribers = append(subscribers, &pubsubtypes.Subscriber{
		Address: "somm10qfl0f55vruhcuqwnqg00uykkdpnl4g3fzx2m7",
		CaCert:  StakecitoSubscriberCA,
		PushUrl: "steward.stakesandstone.com:5734",
	})

	subscribers = append(subscribers, &pubsubtypes.Subscriber{
		Address: "somm19ts4umurauutumqdcu5n8x73fn9dfwwshhf8a4",
		CaCert:  ChorusOneSubscriberCA,
		PushUrl: "sommelier-steward.chorus.one:443",
	})

	subscribers = append(subscribers, &pubsubtypes.Subscriber{
		Address: "somm1p7tskps8hya4ldeu8qfghxwq72g5fzp6aekap7",
		CaCert:  ImperatorSubscriberCA,
		PushUrl: "steward.imperator.co:30812",
	})

	subscribers = append(subscribers, &pubsubtypes.Subscriber{
		Address: "somm1h0df0k7wlzzg53wnglayftwnv6du74ggdu28fz",
		CaCert:  TekuSubscriberCA,
		PushUrl: "sommelier.teku.network:5734",
	})

	subscribers = append(subscribers, &pubsubtypes.Subscriber{
		Address: "somm1tj8echy75u4z0f04z4vda3jzgx4x02de9umhv9",
		CaCert:  ForboleSubscriberCA,
		PushUrl: "steward.sommelier.forbole.com:5734",
	})

	subscribers = append(subscribers, &pubsubtypes.Subscriber{
		Address: "somm15h8ls6mwt8k709wc8f48ycxa80vcu690j6rnwy",
		CaCert:  BoubouSubscriberCA,
		PushUrl: "sommelier-steward.boubounode.com:5734",
	})

	subscribers = append(subscribers, &pubsubtypes.Subscriber{
		Address: "somm1pajq6rx0vgxjzdccukyh3h403rqprfhcsvhrat",
		CaCert:  SleepyKittenSubscriberCA,
		PushUrl: "steward.sommelier.sleepykitten.info:5734",
	})

	subscribers = append(subscribers, &pubsubtypes.Subscriber{
		Address: "somm17qlejm42fz4re8cskz5hlah4hh3s8w9y2yxgu6",
		CaCert:  EverstakeSubscriberCA,
		PushUrl: "sommelier-steward.everstake.one:5734",
	})

	subscribers = append(subscribers, &pubsubtypes.Subscriber{
		Address: "somm1zqklnhsp0q0rew352akcg35a45ruq3vn2c7fym",
		CaCert:  TesselatedSubscriberCA,
		PushUrl: "sommelier.tessageo.com:5734",
	})

	subscribers = append(subscribers, &pubsubtypes.Subscriber{
		Address: "somm1z3jzvtxplxh2c7qn8j4teq63vsngkf0488w3uj",
		CaCert:  ZtakeSubscriberCA,
		PushUrl: "sommelier.ztake.org:5734",
	})

	subscribers = append(subscribers, &pubsubtypes.Subscriber{
		Address: "somm1ydtvj3ruqqq7zxkz9w5lze5ecylh82v03h5udg",
		CaCert:  TwoBuckChuckSubscriberCA,
		PushUrl: "tastings.two-buck-chuck.xyz:5734",
	})

	subscribers = append(subscribers, &pubsubtypes.Subscriber{
		Address: "somm1xvfdclzyw03ye5mfskw6lvvmervexx8hx58823",
		CaCert:  CosmostationSubscriberCA,
		PushUrl: "steward.sommelier.cosmostation.io:15413",
	})

	subscribers = append(subscribers, &pubsubtypes.Subscriber{
		Address: "somm1ge3qxg5ydd25huhf4v8nge8kjsgyps83qvw775",
		CaCert:  MCBSubscriberCA,
		PushUrl: "sommelier.mcbnode.com:5734",
	})

	subscribers = append(subscribers, &pubsubtypes.Subscriber{
		Address: "somm19e4mx2geplpm6ksexxex8dg4dr4a5p7utl4y8z",
		CaCert:  PolychainSubscriberCA,
		PushUrl: "steward.sommelier.unit410.com:5734",
	})

	subscribers = append(subscribers, &pubsubtypes.Subscriber{
		Address: "somm1v78rxsl5ycptaq2dq5mu6ftzqcgk2aqtfx4ryr",
		CaCert:  KingSuperSubscriberCA,
		PushUrl: "sommelier-steward.kingsuper.services:5734",
	})

	subscribers = append(subscribers, &pubsubtypes.Subscriber{
		Address: "somm1rs4hcjr0jgw9ah8ml5p84cvc0yxcvf8krer8wu",
		CaCert:  ChillValidationSubscriberCA,
		PushUrl: "steward0.chillvalidation.com:5734",
	})

	subscribers = append(subscribers, &pubsubtypes.Subscriber{
		Address: "somm1g7kwsxw5khxg2zftpfcla8x4pz7zjukzc6luqy",
		CaCert:  ChainnodesSubscriberCA,
		PushUrl: "steward.chainnodes.net:5734",
	})

	subscribers = append(subscribers, &pubsubtypes.Subscriber{
		Address: "somm13d6vkp03nelzu7aq4v6n88nw0tye2ht7j9xern",
		CaCert:  SevenSeasSubscriberCA,
		PushUrl: "steward.sommelier.7seas.capital:5734",
	})

	subscribers = append(subscribers, &pubsubtypes.Subscriber{
		Address: "somm1xlg9tu8nwyratwhppkmen62putwf3dltkeqvl9",
		CaCert:  GoldenRatioSubscriberCA,
		PushUrl: "sommelier.goldenratiostaking.net:5734",
	})

	subscribers = append(subscribers, &pubsubtypes.Subscriber{
		Address: "somm1kdq2pwdnn5225y0fjk7nzd93errzxmj2ncp60z",
		CaCert:  CryptoCrewSubscriberCA,
		PushUrl: "steward-somm.ccvalidators.com:5734",
	})

	subscribers = append(subscribers, &pubsubtypes.Subscriber{
		Address: "somm1cfdkpueekdxgax0gu5fwq30nfwd2h0mg3kwtqq",
		CaCert:  DoraFactorySubscriberCA,
		PushUrl: "sommelier.dorafactory.org:5734",
	})

	subscribers = append(subscribers, &pubsubtypes.Subscriber{
		Address: "somm1p0d4cg70pk9x49xzrg9dllvj6wxkvtqxfc8490",
		CaCert:  FrenchChocolatineSubscriberCA,
		PushUrl: "sommelier.frenchchocolatine.com:5734",
	})

	subscribers = append(subscribers, &pubsubtypes.Subscriber{
		Address: "somm1ye5qdw92yjj0a2fvpgwgmxh9yymrcmaxn8ed3u",
		CaCert:  FreshStakingSubscriberCA,
		PushUrl: "somm-steward.mitera.net:5734",
	})

	subscribers = append(subscribers, &pubsubtypes.Subscriber{
		Address: "somm1fmw3y7heca7qhfkt5uu3u65mk8gj5tx24k9x68",
		CaCert:  KleomedesSubscriberCA,
		PushUrl: "steward.kleomedes.network:5734",
	})

	subscribers = append(subscribers, &pubsubtypes.Subscriber{
		Address: "somm1q9k53u4fu2v0ksgs84ek4c0xrh269haykxuqrk",
		CaCert:  MeriaSubscriberCA,
		PushUrl: "sommelier.meria.com:5734",
	})

	subscribers = append(subscribers, &pubsubtypes.Subscriber{
		Address: "somm1s70pr2uyct7jtjc69kpkwm3ajysmfgzpwl32vj",
		CaCert:  RorcualSubscriberCA,
		PushUrl: "steward.rorcualnodes.com:5734",
	})

	subscribers = append(subscribers, &pubsubtypes.Subscriber{
		Address: "somm1u7n35gtu85qrtu92ws5fsgs6ea4ay32nach7q7",
		CaCert:  BlockHuntersSubscriberCA,
		PushUrl: "somm.blockhunters.sbs:5734",
	})

	return subscribers
}
