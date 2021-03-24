module github.com/peggyjv/sommelier

go 1.15

require (
	github.com/armon/go-metrics v0.3.6
	github.com/cosmos/cosmos-sdk v0.42.2
	github.com/cosmos/go-bip39 v1.0.0
	github.com/cosmos/gravity-bridge/module v0.0.0-20210324205831-32074c2eab98
	github.com/ethereum/go-ethereum v1.10.1
	github.com/gogo/protobuf v1.3.3
	github.com/golang/protobuf v1.5.1
	github.com/gorilla/mux v1.8.0
	github.com/grpc-ecosystem/grpc-gateway v1.16.0
	github.com/machinebox/graphql v0.2.2
	github.com/matryer/is v1.4.0 // indirect
	github.com/rakyll/statik v0.1.7
	github.com/regen-network/cosmos-proto v0.3.1
	github.com/spf13/cast v1.3.1
	github.com/spf13/cobra v1.1.3
	github.com/spf13/viper v1.7.1
	github.com/stretchr/testify v1.7.0
	github.com/tendermint/tendermint v0.34.8
	github.com/tendermint/tm-db v0.6.4
	golang.org/x/crypto v0.0.0-20210322153248-0c34fe9e7dc2 // indirect
	golang.org/x/sync v0.0.0-20210220032951-036812b2e83c
	golang.org/x/sys v0.0.0-20210324051608-47abb6519492 // indirect
	google.golang.org/genproto v0.0.0-20210324141432-3032e8ff099e
	google.golang.org/grpc v1.36.0
	gopkg.in/yaml.v2 v2.4.0
)

replace google.golang.org/grpc => google.golang.org/grpc v1.33.2

replace github.com/gogo/protobuf => github.com/regen-network/protobuf v1.3.3-alpha.regen.1
