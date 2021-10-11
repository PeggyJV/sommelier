module github.com/peggyjv/sommelier

go 1.15

require (
	github.com/CosmWasm/wasmvm v1.0.0-soon2 // indirect
	github.com/cenkalti/backoff/v4 v4.1.1 // indirect
	github.com/cosmos/cosmos-sdk v0.44.1
	github.com/cosmos/go-bip39 v1.0.0
	github.com/cosmos/ibc-go v1.1.0
	github.com/deckarep/golang-set v0.0.0-20180603214616-504e848d77ea
	github.com/ethereum/go-ethereum v1.10.8
	github.com/gogo/protobuf v1.3.3
	github.com/golang/protobuf v1.5.2
	github.com/google/gofuzz v1.2.0 // indirect
	github.com/gorilla/mux v1.8.0
	github.com/grpc-ecosystem/grpc-gateway v1.16.0
	github.com/miguelmota/go-ethereum-hdwallet v0.1.1
	github.com/moby/term v0.0.0-20210619224110-3f7ff695adc6 // indirect
	github.com/ory/dockertest/v3 v3.7.0
	github.com/peggyjv/gravity-bridge/module v0.2.14-0.20210923174501-aea3338bb368
	//github.com/peggyjv/gravity-bridge/module v0.2.0
	github.com/rakyll/statik v0.1.7
	github.com/regen-network/cosmos-proto v0.3.1
	github.com/rs/zerolog v1.25.0 // indirect
	github.com/spf13/cast v1.4.1
	github.com/spf13/cobra v1.2.1
	github.com/spf13/viper v1.9.0
	github.com/stretchr/testify v1.7.0
	github.com/tendermint/tendermint v0.34.13
	github.com/tendermint/tm-db v0.6.4
	github.com/terra-money/core v0.5.6
	github.com/xeipuuv/gojsonpointer v0.0.0-20190905194746-02993c407bfb // indirect
	google.golang.org/genproto v0.0.0-20210828152312-66f60bf46e71
	google.golang.org/grpc v1.40.0
)

replace google.golang.org/grpc => google.golang.org/grpc v1.33.2

replace github.com/gogo/protobuf => github.com/regen-network/protobuf v1.3.3-alpha.regen.1
