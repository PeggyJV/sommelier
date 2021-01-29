module github.com/peggyjv/sommelier

go 1.15

require (
	github.com/althea-net/peggy/module v0.0.0-20210128235107-1d1f78f131b6
	github.com/cosmos/cosmos-sdk v0.40.0
	github.com/gogo/protobuf v1.3.1
	github.com/golang/protobuf v1.4.3
	github.com/google/gofuzz v1.1.1-0.20200604201612-c04b05f3adfa // indirect
	github.com/gorilla/mux v1.8.0
	github.com/grpc-ecosystem/grpc-gateway v1.16.0
	github.com/pkg/errors v0.9.1
	github.com/rakyll/statik v0.1.7
	github.com/spf13/cast v1.3.1
	github.com/spf13/cobra v1.1.1
	github.com/stretchr/testify v1.6.1
	github.com/tendermint/tendermint v0.34.1
	github.com/tendermint/tm-db v0.6.3
	google.golang.org/genproto v0.0.0-20210108203827-ffc7fda8c3d7
	google.golang.org/grpc v1.33.2
	gopkg.in/yaml.v2 v2.4.0
)

replace github.com/gogo/protobuf => github.com/regen-network/protobuf v1.3.2-alpha.regen.4
