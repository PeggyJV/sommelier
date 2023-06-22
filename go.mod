module github.com/peggyjv/sommelier/v4

go 1.15

require (
	github.com/StackExchange/wmi v1.2.1 // indirect
	github.com/cosmos/cosmos-sdk v0.47.3
	github.com/cosmos/go-bip39 v1.0.0
	github.com/cosmos/ibc-go/v2 v2.0.0
	github.com/ethereum/go-ethereum v1.10.11
	github.com/gogo/protobuf v1.3.3
	github.com/golang/mock v1.6.0
	github.com/golang/protobuf v1.5.3
	github.com/gorilla/mux v1.8.0
	github.com/grpc-ecosystem/grpc-gateway v1.16.0
	github.com/miguelmota/go-ethereum-hdwallet v0.1.1
	github.com/ory/dockertest/v3 v3.9.1
	github.com/peggyjv/gravity-bridge/module/v2 v2.0.2
	github.com/rakyll/statik v0.1.7
	github.com/regen-network/cosmos-proto v0.3.1
	github.com/spf13/cast v1.5.0
	github.com/spf13/cobra v1.6.1
	github.com/spf13/viper v1.14.0
	github.com/stretchr/testify v1.8.2
	github.com/tendermint/tendermint v0.37.0-rc2
	github.com/tendermint/tm-db v0.6.7
	google.golang.org/genproto v0.0.0-20230306155012-7f2fa6fef1f4
	google.golang.org/grpc v1.55.0
	gopkg.in/yaml.v2 v2.4.0
)

replace google.golang.org/grpc => google.golang.org/grpc v1.33.2

replace github.com/gogo/protobuf => github.com/regen-network/protobuf v1.3.3-alpha.regen.1

replace github.com/confio/ics23/go => github.com/cosmos/cosmos-sdk/ics23/go v0.8.0
