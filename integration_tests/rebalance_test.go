package integration_tests

import (
	"context"
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/peggyjv/sommelier/x/allocation/types"
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

		s.T().Logf("checking that test cellar exists in the chain")
		val := s.chain.validators[0]
		s.Require().Eventuallyf(func() bool {
			kb, err := val.keyring()
			s.Require().NoError(err)
			clientCtx, err := s.chain.clientContext("tcp://localhost:26657", &kb, "val", val.keyInfo.GetAddress())
			s.Require().NoError(err)

			queryClient := types.NewQueryClient(clientCtx)
			res, err := queryClient.QueryCellars(context.Background(), &types.QueryCellarsRequest{})
			if err != nil {
				return false
			}
			if res == nil {
				return false
			}
			for _, c := range res.Cellars {
				if c.Id == commit.Cellar.Id {
					return true
				}
			}
			return false
		},
			30*time.Second,
			2*time.Second,
			"hardhat cellar not found in chain")

		s.T().Logf("wait for new vote period start")
		val = s.chain.validators[0]
		s.Require().Eventuallyf(func() bool {
			kb, err := val.keyring()
			s.Require().NoError(err)
			clientCtx, err := s.chain.clientContext("tcp://localhost:26657", &kb, "val", val.keyInfo.GetAddress())
			s.Require().NoError(err)

			queryClient := types.NewQueryClient(clientCtx)
			res, err := queryClient.QueryCommitPeriod(context.Background(), &types.QueryCommitPeriodRequest{})
			if err != nil {
				return false
			}
			if res.VotePeriodStart != res.CurrentHeight {
				return false
			}

			return true
		}, 65*time.Second, 1*time.Second,"new vote period never seen")

		s.T().Logf("sending pre-commits")
		for i, orch := range s.chain.orchestrators {
			s.Require().Eventuallyf(func() bool {
				clientCtx, err := s.chain.clientContext("tcp://localhost:26657", orch.keyring, "orch", orch.keyInfo.GetAddress())
				s.Require().NoError(err)

				precommitMsg, err := types.NewMsgAllocationPrecommit(*commit.Cellar, salt, orch.keyInfo.GetAddress())
				s.Require().NoError(err, "unable to create precommit")

				response, err := s.chain.sendMsgs(*clientCtx, precommitMsg)
				if err != nil {
					s.T().Logf("error: %s", err)
					return false
				}
				if response.Code != 0 {
					return false
				}
				return true
			}, 10*time.Second, 500*time.Millisecond, "unable to deploy precommit for node %d", i)
			s.T().Logf("precommit for %d node sent successfully", i)
		}

		s.T().Logf("checking pre-commit for first validator")
		val = s.chain.validators[0]
		s.Require().Eventuallyf(func() bool {
			kb, err := val.keyring()
			s.Require().NoError(err)
			clientCtx, err := s.chain.clientContext("tcp://localhost:26657", &kb, "val", val.keyInfo.GetAddress())
			s.Require().NoError(err)

			queryClient := types.NewQueryClient(clientCtx)

			all, err := queryClient.QueryAllocationPrecommits(context.Background(), &types.QueryAllocationPrecommitsRequest{})
			s.T().Logf("all %s", all.Precommits)

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
			expectedPrecommit, err := types.NewMsgAllocationPrecommit(*commit.Cellar, salt, s.chain.orchestrators[0].keyInfo.GetAddress())
			s.Require().NoError(err, "unable to create precommit")
			s.Require().Equal(res.Precommit.CellarId, commit.Cellar.Id, "cellar ids unequal")
			s.Require().Equal(res.Precommit.Hash, expectedPrecommit.Precommit[0].Hash, "commit hashes unequal")

			return true
		},
			30*time.Second,
			2*time.Second,
			"pre-commit not found for validator %s",
			val.keyInfo.GetAddress().String())

		s.T().Logf("sending commits")
		for i, orch := range s.chain.orchestrators {
			s.Require().Eventuallyf(func() bool {
				clientCtx, err := s.chain.clientContext("tcp://localhost:26657", orch.keyring, "orch", orch.keyInfo.GetAddress())
				s.Require().NoError(err)

				commitMsg := types.NewMsgAllocationCommit([]*types.Allocation{&commit}, orch.keyInfo.GetAddress())

				response, err := s.chain.sendMsgs(*clientCtx, commitMsg)
				if err != nil {
					return false
				}
				if response.Code != 0 {
					return false
				}

				return true
			}, 10*time.Second, 500*time.Millisecond, "unable to deploy commit for node %d", i)
		}

		s.T().Logf("checking for updated tick ranges in cellar")
		trs, err = s.getTickRanges()
		s.Require().NoError(err)
		s.Require().Len(trs, 4)
		for i, tr := range trs {
			s.Require().Equal((i+2)*100, tr.Upper)
			s.Require().Equal((i+1)*100, tr.Lower)
			s.Require().Equal((i+1)*10, tr.Weight)
		}
	})
}
