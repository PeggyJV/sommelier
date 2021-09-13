package integration_tests

import (
	"context"
	"fmt"
	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/peggyjv/sommelier/x/allocation/types"
	"math/big"
)

func (s *IntegrationTestSuite) TestRebalance() {
	s.Run("Bring up chain, and submit a re-balance", func() {

		ethClient, err := ethclient.Dial(fmt.Sprintf("http://%s", s.ethResource.GetHostPort("8545/tcp")))
		s.Require().NoError(err)

		suggestedGasPrice, err := ethClient.SuggestGasPrice(context.Background())
		s.Require().NoError(err)

		//tickInfos := [4]types.TickRange{}
		for i := 0; i < 4; i++ {
			bz, err := ethClient.CallContract(context.Background(), ethereum.CallMsg{
				From:       common.HexToAddress(s.chain.validators[0].ethereumKey.address),
				To:         &hardhatCellar,
				Gas:        0,
				GasPrice:   suggestedGasPrice,
				GasFeeCap:  big.NewInt(1),
				GasTipCap:  big.NewInt(1),
				Value:      nil,
				Data:       types.CellarTickInfo(uint(i)),
				AccessList: nil,
			}, nil)
			s.T().Logf("bytes received %b", bz)
			s.Require().NoError(err)

		}
		salt := "testsalt"
		commit := types.Allocation{
			Cellar: &types.Cellar{
				Id:         hardhatCellar.String(),
				TickRanges: []*types.TickRange{},
			},
		}

		s.T().Logf("sending pre-commits")
		for _, val := range s.chain.validators {
			clientCtx, err := s.chain.clientContext("tcp://localhost:26657", *val)
			s.Require().NoError(err)

			hash, err := val.hashCellar(*commit.Cellar, salt)
			s.Require().NoError(err, "unable to hash cellar")
			precommitMsg := types.NewMsgAllocationPrecommit(hash, val.keyInfo.GetAddress())

			response, err := s.chain.sendMsgs(*clientCtx, precommitMsg)
			s.Require().NoError(err, "unable to sign precommit")
			s.Require().NotZerof(response.Code, "non-zero response from rpc call for msg", precommitMsg)
		}

		s.T().Logf("sending commits")
		for _, val := range s.chain.validators {
			clientCtx, err := s.chain.clientContext("tcp://localhost:26657", *val)
			s.Require().NoError(err)

			commitMsg := types.MsgAllocationCommit{
				Commit: []*types.Allocation{&commit},
				Signer: val.keyInfo.GetAddress().String(),
			}

			response, err := s.chain.sendMsgs(*clientCtx, &commitMsg)
			s.Require().NoError(err, "unable to send precommit")
			s.Require().NotZerof(response.Code, "non-zero response from rpc call for msg", commitMsg)
		}

		s.T().Logf("checking for updated tick ranges in cellar")
		s.Require().Fail("UNIMPLEMENTED")
	})
}
