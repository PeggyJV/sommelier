package upgrade_test

import (
	"fmt"

	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/peggyjv/sommelier/v4/app/params"
)

const (
	keyringPassphrase = "testpassphrase"
	keyringAppName    = "testnet"
)

var (
	encodingConfig params.EncodingConfig
	cdc            codec.Codec
)

type chain struct {
	dataDir       string
	id            string
	validators    []*validator
	orchestrators []*orchestrator
}

func (c *chain) configDir() string {
	return fmt.Sprintf("%s/%s", c.dataDir, c.id)
}
