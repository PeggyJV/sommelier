module github.com/peggyjv/sommelier/v4

go 1.15

require (
	github.com/Microsoft/go-winio v0.5.1 // indirect
	github.com/StackExchange/wmi v1.2.1 // indirect
	github.com/allegro/bigcache v1.2.1 // indirect
	github.com/cosmos/cosmos-sdk v0.44.5-patch
	github.com/cosmos/ibc-go/v2 v2.0.0
	github.com/ethereum/go-ethereum v1.10.11
	github.com/gogo/protobuf v1.3.3
	github.com/golang/mock v1.6.0
	github.com/golang/protobuf v1.5.2
	github.com/google/uuid v1.2.0 // indirect
	github.com/gorilla/mux v1.8.0
	github.com/grpc-ecosystem/grpc-gateway v1.16.0
	github.com/opencontainers/image-spec v1.0.2 // indirect
	github.com/peggyjv/gravity-bridge/module/v2 v2.0.2
	github.com/rakyll/statik v0.1.7
	github.com/regen-network/cosmos-proto v0.3.1
	github.com/spf13/cast v1.4.1
	github.com/spf13/cobra v1.2.1
	github.com/stretchr/testify v1.8.0
	github.com/tendermint/tendermint v0.34.14
	github.com/tendermint/tm-db v0.6.4
	github.com/tklauser/go-sysconf v0.3.7 // indirect
	google.golang.org/genproto v0.0.0-20211208223120-3a66f561d7aa
	google.golang.org/grpc v1.43.0
	gopkg.in/yaml.v2 v2.4.0
)

replace google.golang.org/grpc => google.golang.org/grpc v1.33.2

replace github.com/gogo/protobuf => github.com/regen-network/protobuf v1.3.3-alpha.regen.1

replace github.com/confio/ics23/go => github.com/cosmos/cosmos-sdk/ics23/go v0.8.0
