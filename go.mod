module github.com/peggyjv/sommelier

go 1.15

require (
	github.com/cosmos/cosmos-sdk v0.42.4-0.20210623214207-eb0fc466c99b
	github.com/ethereum/go-ethereum v1.10.1 // indirect
	github.com/gorilla/mux v1.8.0
	github.com/peggyjv/gravity-bridge/module v0.1.14
	github.com/rakyll/statik v0.1.7
	github.com/spf13/cast v1.3.1
	github.com/spf13/cobra v1.1.3
	github.com/stretchr/testify v1.7.0
	github.com/tendermint/tendermint v0.34.8
	github.com/tendermint/tm-db v0.6.4
	golang.org/x/crypto v0.0.0-20210322153248-0c34fe9e7dc2 // indirect
	google.golang.org/genproto v0.0.0-20210324141432-3032e8ff099e // indirect
)

replace google.golang.org/grpc => google.golang.org/grpc v1.33.2

replace github.com/gogo/protobuf => github.com/regen-network/protobuf v1.3.3-alpha.regen.1
