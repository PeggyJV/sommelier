package integration_tests

import (
	"github.com/peggyjv/sommelier/x/allocation/types"
)

func (s *IntegrationTestSuite) TestRebalance() {
	s.Run("Bring up chain, and submit a re-balance", func() {
		s.T().Logf("sending pre-commits")

		salt := "testsalt"
		commit := types.Allocation{
			Cellar: &types.Cellar{
				Id:         "0x6ea5992aB4A78D5720bD12A089D13c073d04B55d",
				TickRanges: []*types.TickRange{},
			},
		}

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
			s.Require().NoError(err, "unable to sign precommit")
			s.Require().NotZerof(response.Code, "non-zero response from rpc call for msg", commitMsg)
		}
	})
}
