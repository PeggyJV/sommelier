package integration_tests

import (
	"context"
	"time"

	gravitytypes "github.com/peggyjv/gravity-bridge/module/x/gravity/types"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/peggyjv/sommelier/x/allocation/types"
)

func (s *IntegrationTestSuite) TestRebalance() {
	s.Run("Bring up chain, and submit a re-balance", func() {

		trs := s.getTickRanges()
		s.Require().Len(trs, 1)
		s.Require().Equal(int32(198780), trs[0].Upper)
		s.Require().Equal(int32(198120), trs[0].Lower)
		s.Require().Equal(uint32(100), trs[0].Weight)

		commit := types.Allocation{
			Cellar: &types.Cellar{
				Id: hardhatCellar.String(),
				TickRanges: []*types.TickRange{
					{Upper: 200000, Lower: 190000, Weight: 200},
				},
			},
			Salt: "testsalt",
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
		}, 30*time.Second, 2*time.Second, "hardhat cellar not found in chain")

		//s.T().Logf("Trigger valset update via redelegation")
		//val = s.chain.validators[0]
		//s.Require().Eventuallyf(func() bool {
		//	kb, err := val.keyring()
		//	s.Require().NoError(err)
		//
		//	clientCtx, err := s.chain.clientContext("tcp://localhost:26657", &kb, "val", val.keyInfo.GetAddress())
		//	s.Require().NoError(err)
		//
		//	redelegateMsg := stakingtypes.NewMsgBeginRedelegate(
		//		s.chain.validators[0].keyInfo.GetAddress(),
		//		sdk.ValAddress(s.chain.validators[0].keyInfo.GetAddress()),
		//		sdk.ValAddress(s.chain.validators[1].keyInfo.GetAddress()),
		//		sdk.Coin{
		//			Denom:  bondDenom,
		//			Amount: sdk.NewInt(1073741824/2),
		//		},
		//	)
		//	response, err := s.chain.sendMsgs(*clientCtx, redelegateMsg)
		//	if err != nil {
		//		s.T().Logf("error: %s", err)
		//		return false
		//	}
		//	if response.Code != 0 {
		//		if response.Code != 32 {
		//			s.T().Log(response)
		//		}
		//		return false
		//	}
		//
		//	return true
		//}, 10*time.Second, 500*time.Millisecond, "unable to send redelegation from first node to second", )


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
				if res.CurrentHeight%10 == 0 {
					s.T().Logf("current height: %d, period end: %d", res.CurrentHeight, res.VotePeriodEnd)
				}
				return false
			}

			return true
		}, 105*time.Second, 1*time.Second, "new vote period never seen")

		s.T().Logf("sending pre-commits")
		for i, orch := range s.chain.orchestrators {
			i := i
			orch := orch
			s.Require().Eventuallyf(func() bool {
				clientCtx, err := s.chain.clientContext("tcp://localhost:26657", orch.keyring, "orch", orch.keyInfo.GetAddress())
				s.Require().NoError(err)

				delegatedVal := s.chain.validators[i]

				precommitMsg, err := types.NewMsgAllocationPrecommit(*commit.Cellar, commit.Salt, orch.keyInfo.GetAddress(), sdk.ValAddress(delegatedVal.keyInfo.GetAddress()))
				s.Require().NoError(err, "unable to create precommit")

				response, err := s.chain.sendMsgs(*clientCtx, precommitMsg)
				if err != nil {
					s.T().Logf("error: %s", err)
					return false
				}
				if response.Code != 0 {
					if response.Code != 32 {
						s.T().Log(response)
					}
					return false
				}

				s.T().Logf("precommit for %d node with hash %x sent successfully", i, precommitMsg.Precommit[0].Hash)
				return true
			}, 10*time.Second, 500*time.Millisecond, "unable to deploy precommit for node %d", i)
		}

		s.T().Logf("checking pre-commits for validators")
		for i, val := range s.chain.validators {
			i := i
			val := val
			s.Require().Eventuallyf(func() bool {
				kb, err := val.keyring()
				s.Require().NoError(err)
				clientCtx, err := s.chain.clientContext("tcp://localhost:26657", &kb, "val", val.keyInfo.GetAddress())
				s.Require().NoError(err)

				queryClient := types.NewQueryClient(clientCtx)
				signerVal := sdk.ValAddress(val.keyInfo.GetAddress())
				res, err := queryClient.QueryAllocationPrecommit(context.Background(), &types.QueryAllocationPrecommitRequest{
					Validator: signerVal.String(),
					Cellar:    hardhatCellar.String(),
				})
				if err != nil {
					return false
				}
				if res == nil {
					return false
				}
				expectedPrecommit, err := types.NewMsgAllocationPrecommit(*commit.Cellar, commit.Salt, s.chain.orchestrators[i].keyInfo.GetAddress(), sdk.ValAddress(val.keyInfo.GetAddress()))
				s.Require().NoError(err, "unable to create precommit")
				s.Require().Equal(res.Precommit.CellarId, commit.Cellar.Id, "cellar ids unequal")
				s.Require().Equal(res.Precommit.Hash, expectedPrecommit.Precommit[0].Hash, "commit hashes unequal")

				return true
			},
				30*time.Second,
				2*time.Second,
				"pre-commit not found for validator %s",
				val.keyInfo.GetAddress().String())
		}

		//kb, err := val.keyring()

		//clientCtx, err := s.chain.clientContext("tcp://localhost:26657", &kb, "val", val.keyInfo.GetAddress())
		//s.Require().NoError(err)
		//queryClient := types.NewQueryClient(clientCtx)
		//res, err := queryClient.QueryAllocationPrecommits(context.Background(), &types.QueryAllocationPrecommitsRequest{})
		//for _, pc := range res.Precommits {
		//	s.T().Logf("precommit: hash %x", pc.Hash)
		//}

		s.T().Logf("sending commits")
		for i, orch := range s.chain.orchestrators {
			i := i
			orch := orch
			s.Require().Eventuallyf(func() bool {
				clientCtx, err := s.chain.clientContext("tcp://localhost:26657", orch.keyring, "orch", orch.keyInfo.GetAddress())
				s.Require().NoError(err)

				commitMsg := types.NewMsgAllocationCommit([]*types.Allocation{&commit}, orch.keyInfo.GetAddress())
				response, err := s.chain.sendMsgs(*clientCtx, commitMsg)
				if err != nil {
					s.T().Logf("error: %s", err)
					return false
				}
				if response.Code != 0 {
					if response.Code != 32 {
						s.T().Logf("response: %s", response)
						s.FailNow("failing")
					}
					return false
				}

				return true
			}, 10*time.Second, 500*time.Millisecond, "unable to deploy commit for node %d", i)
		}

		s.T().Logf("checking commits for validators")
		for _, val := range s.chain.validators {
			val := val
			s.Require().Eventuallyf(func() bool {
				kb, err := val.keyring()
				s.Require().NoError(err)
				clientCtx, err := s.chain.clientContext("tcp://localhost:26657", &kb, "val", val.keyInfo.GetAddress())
				s.Require().NoError(err)

				queryClient := types.NewQueryClient(clientCtx)
				res, err := queryClient.QueryAllocationCommit(context.Background(), &types.QueryAllocationCommitRequest{Validator: sdk.ValAddress(val.keyInfo.GetAddress()).String(), Cellar: hardhatCellar.String()})
				if err != nil {
					return false
				}
				if res == nil {
					return false
				}
				s.Require().Equal(res.Commit.Cellar.Id, commit.Cellar.Id, "cellar ids unequal")

				return true
			},
				30*time.Second,
				2*time.Second,
				"commit not found for validator %s",
				val.keyInfo.GetAddress().String())
		}

		s.T().Logf("waiting for end of vote period, endblocker to run")
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
				if res.CurrentHeight%10 == 0 {
					s.T().Logf("current height: %d, period end: %d", res.CurrentHeight, res.VotePeriodEnd)
				}
				return false
			}

			return true
		}, 105*time.Second, 1*time.Second, "new vote period never seen")

		s.T().Logf("checking for updated tick ranges in cellar")
		s.Require().Eventuallyf(func() bool {
			kb, err := val.keyring()
			s.Require().NoError(err)
			clientCtx, err := s.chain.clientContext("tcp://localhost:26657", &kb, "val", val.keyInfo.GetAddress())
			s.Require().NoError(err)
			gravityQueryClient := gravitytypes.NewQueryClient(clientCtx)
			res, err := gravityQueryClient.UnsignedContractCallTxs(context.Background(), &gravitytypes.UnsignedContractCallTxsRequest{
				Address: val.keyInfo.GetAddress().String(),
			})
			if err != nil {
				s.T().Logf("error: %s", err)
				if res != nil {
					s.T().Logf("response: %s", res)
				}
			}
			//s.T().Logf("unsigned contract call txs: %s", res.Calls)
			for _, call := range res.Calls {
				s.T().Logf("contract call; nonce: %d, scope: %x, store index: %x", call.InvalidationNonce, call.InvalidationScope, call.GetStoreIndex())
			}

			confirmsRes, err := gravityQueryClient.ContractCallTxConfirmations(context.Background(), &gravitytypes.ContractCallTxConfirmationsRequest{
				InvalidationScope: commit.Cellar.ABIEncodedRebalanceBytes(),
				InvalidationNonce: 1,
			})

			if err != nil {
				s.T().Logf("error: %s", err)
				if res != nil {
					s.T().Logf("response: %s", confirmsRes)
				}
			}
			s.T().Logf("contract call tx confirms: %s", confirmsRes.Signatures)

			trs = s.getTickRanges()
			if err != nil {
				s.T().Logf("error: %s", err)
				return false
			}
			if len(trs) != 4 {
				s.T().Logf("wrong length: %d", len(trs))
				return false
			}
			for i, tr := range trs {
				if int32((i+2)*100) != tr.Upper {
					s.T().Logf("wrong upper %s", tr.String())
					return false
				}
				if int32((i+1)*100) != tr.Lower {
					s.T().Logf("wrong lower %s", tr.String())
					return false
				}
				if uint32((i+1)*10) != tr.Weight {
					s.T().Logf("wrong weight %s", tr.String())
					return false
				}
			}

			return true
		}, 30*time.Second, 2*time.Second, "cellar ticks never updated")
	})
}
