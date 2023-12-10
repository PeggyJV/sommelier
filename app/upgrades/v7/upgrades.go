package v7

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/module"
	upgradetypes "github.com/cosmos/cosmos-sdk/x/upgrade/types"
	auctionkeeper "github.com/peggyjv/sommelier/v7/x/auction/keeper"
	auctiontypes "github.com/peggyjv/sommelier/v7/x/auction/types"
	cellarfeeskeeper "github.com/peggyjv/sommelier/v7/x/cellarfees/keeper"
	cellarfeestypes "github.com/peggyjv/sommelier/v7/x/cellarfees/types"
	pubsubkeeper "github.com/peggyjv/sommelier/v7/x/pubsub/keeper"
	pubsubtypes "github.com/peggyjv/sommelier/v7/x/pubsub/types"
)

func CreateUpgradeHandler(
	mm *module.Manager,
	configurator module.Configurator,
	auctionKeeper auctionkeeper.Keeper,
	cellarfeesKeeper cellarfeeskeeper.Keeper,
	pubsubKeeper pubsubkeeper.Keeper,
) upgradetypes.UpgradeHandler {
	return func(ctx sdk.Context, plan upgradetypes.Plan, fromVM module.VersionMap) (module.VersionMap, error) {
		ctx.Logger().Info("v7 upgrade: entering handler")

		// We must manually run InitGenesis for incentives, pubsub, and auctions so we can adjust their values
		// during the upgrade process. RunMigrations will migrate to the new cork version. Setting the consensus
		// version to 1 prevents RunMigrations from running InitGenesis itself.
		fromVM[auctiontypes.ModuleName] = mm.Modules[auctiontypes.ModuleName].ConsensusVersion()
		fromVM[pubsubtypes.ModuleName] = mm.Modules[pubsubtypes.ModuleName].ConsensusVersion()

		// Params values were introduced in this upgrade but no migration was necessary, so we initialize the
		// new values to their defaults
		ctx.Logger().Info("v7 upgrading: setting cellarfees default params")
		cellarfeesKeeper.SetParams(ctx, cellarfeestypes.DefaultParams())

		ctx.Logger().Info("v7 upgrade: initializing auction genesis state")
		auctionInitGenesis(ctx, auctionKeeper)

		ctx.Logger().Info("v7 upgrade: initializing pubsub genesis state")
		pubsubInitGenesis(ctx, pubsubKeeper)

		ctx.Logger().Info("v7 upgrade: running migrations and exiting handler")
		return mm.RunMigrations(ctx, configurator, fromVM)
	}
}

// Initialize the auction module with prices for some stablecoins and SOMM.
func auctionInitGenesis(ctx sdk.Context, auctionKeeper auctionkeeper.Keeper) {
	genesisState := auctiontypes.DefaultGenesisState()

	usomm52WeekLow := sdk.MustNewDecFromStr("0.079151")
	oneDollar := sdk.MustNewDecFromStr("1.0")
	// TODO(bolten): update LastUpdatedBlock to the upgrade height when finalized
	usommPrice := auctiontypes.TokenPrice{
		Denom:            "usomm",
		UsdPrice:         usomm52WeekLow,
		Exponent:         6,
		LastUpdatedBlock: 1,
	}

	// setting stables to 1 dollar
	usdcPrice := auctiontypes.TokenPrice{
		Denom:            "gravity0xA0b86991c6218b36c1d19D4a2e9Eb0cE3606eB48",
		UsdPrice:         oneDollar,
		Exponent:         6,
		LastUpdatedBlock: 1,
	}

	usdtPrice := auctiontypes.TokenPrice{
		Denom:            "gravity0xdAC17F958D2ee523a2206206994597C13D831ec7",
		UsdPrice:         oneDollar,
		Exponent:         6,
		LastUpdatedBlock: 1,
	}

	daiPrice := auctiontypes.TokenPrice{
		Denom:            "gravity0x6B175474E89094C44Da98b954EedeAC495271d0F",
		UsdPrice:         oneDollar,
		Exponent:         18,
		LastUpdatedBlock: 1,
	}

	genesisState.TokenPrices = []*auctiontypes.TokenPrice{&usommPrice, &usdcPrice, &usdtPrice, &daiPrice}

	auctionkeeper.InitGenesis(ctx, auctionKeeper, genesisState)
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

	// TODO(bolten): update the cellars list with the current total cellars
	cellars := []string{
		"0x7bAD5DF5E11151Dc5Ee1a648800057C5c934c0d5", // Aave V2
		"0x6b7f87279982d919bbf85182ddeab179b366d8f2", // ETH-BTC Trend
		"0x6e2dac3b9e9adc0cbbae2d0b9fd81952a8d33872", // ETH-BTC Momentum
		"0x3F07A84eCdf494310D397d24c1C78B041D2fa622", // Steady ETH
		"0x4986fD36b6b16f49b43282Ee2e24C5cF90ed166d", // Steady BTC
		"0x05641a27C82799AaF22b436F20A3110410f29652", // Steady MATIC
		"0x6f069f711281618467dae7873541ecc082761b33", // Steady UNI
		"0x97e6E0a40a3D02F12d1cEC30ebfbAE04e37C119E", // Real Yield USD
	}

	// Set 7seas publisher intents for existing cellars
	publisherIntents := make([]*pubsubtypes.PublisherIntent, 8)
	for _, cellar := range cellars {
		publisherIntents = append(publisherIntents, &pubsubtypes.PublisherIntent{
			SubscriptionId:     cellar,
			PublisherDomain:    SevenSeasDomain,
			Method:             pubsubtypes.PublishMethod_PUSH,
			AllowedSubscribers: pubsubtypes.AllowedSubscribers_VALIDATORS,
		})
	}

	// Set default subscriptions for 7seas as the publisher for existing cellars
	defaultSubscriptions := make([]*pubsubtypes.DefaultSubscription, 8)
	for _, cellar := range cellars {
		defaultSubscriptions = append(defaultSubscriptions, &pubsubtypes.DefaultSubscription{
			SubscriptionId:  cellar,
			PublisherDomain: SevenSeasDomain,
		})
	}

	// Create subscribers and intents for existing validators
	subscribers := createSubscribers()
	subscriberIntents := make([]*pubsubtypes.SubscriberIntent, 208)
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

	pubsubkeeper.InitGenesis(ctx, pubsubKeeper, genesisState)
}

// Addresses are from the validator's delegate orchestrator key and certs/URLs captured from the
// steward-registry repo.
// query to get orchestrator key: sommelier query gravity delegate-keys-by-validator sommvaloper<rest_of_val_address>
// See source data at: https://github.com/PeggyJV/steward-registry
// data captured at commit ecdb7f386e7e573edb5d8f6ad22a1a67cfa21863
func createSubscribers() []*pubsubtypes.Subscriber {
	subscribers := make([]*pubsubtypes.Subscriber, 26)

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
