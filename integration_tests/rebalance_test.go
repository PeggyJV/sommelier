package integration_tests

import (
	"context"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/peggyjv/sommelier/x/allocation/types"
	"time"
)

func (s *IntegrationTestSuite) TestRebalance() {
	s.Run("Bring up chain, and submit a re-balance", func() {

		trs, err := s.getTickRanges()
		s.Require().NoError(err)
		s.Require().Len(trs, 3)

		salt := "testsalt"
		commit := types.Allocation{
			Cellar: &types.Cellar{
				Id: hardhatCellar.String(),
				TickRanges: []*types.TickRange{
					{200, 100, 10},
					{300, 200, 20},
					{400, 300, 30},
					{500, 400, 40},
				},
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
			s.Require().NoError(err, "unable to send precommit")
			s.Require().Zerof(response.Code, "non-zero response from rpc call for msg %s, response %s", precommitMsg, response)
		}

		s.T().Logf("checking pre-commit for first validator")
		val := s.chain.validators[0]
		s.Require().Eventuallyf(func() bool {
			clientCtx, err := s.chain.clientContext("tcp://localhost:26657", *val)
			s.Require().NoError(err)

			queryClient := types.NewQueryClient(clientCtx)
			res, err := queryClient.QueryAllocationPrecommit(context.Background(), &types.QueryAllocationPrecommitRequest{
				Validator: sdk.ValAddress(val.keyInfo.GetAddress()).String(),
				Cellar:    hardhatCellar.String(),
			})
			if err != nil {
				return false
			}
			if res == nil {
				return false
			}

			return true
		},
		30 * time.Second,
		2 * time.Second,
		"pre-commit not found for validator %s",
		val.keyInfo.GetAddress().String())

		s.T().Logf("sending commits")
		for _, val := range s.chain.validators {
			clientCtx, err := s.chain.clientContext("tcp://localhost:26657", *val)
			s.Require().NoError(err)

			commitMsg := types.MsgAllocationCommit{
				Commit: []*types.Allocation{&commit},
				Signer: val.keyInfo.GetAddress().String(),
			}

			response, err := s.chain.sendMsgs(*clientCtx, &commitMsg)
			s.Require().NoError(err, "unable to send commit")
			s.Require().Zerof(response.Code, "non-zero response from rpc call for msg %s, response %s", commitMsg, response)
		}

		s.T().Logf("checking for updated tick ranges in cellar")
		trs, err = s.getTickRanges()
		s.Require().NoError(err)
		s.Require().Len(trs, 4)
		for i, tr := range trs {
			s.Require().Equal((i + 2) * 100, tr.Upper)
			s.Require().Equal((i + 1) * 100, tr.Lower)
			s.Require().Equal((i + 1) * 10, tr.Weight)
		}
	})
}
