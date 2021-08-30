package integration_tests

//func (s *IntegrationTestSuite) TestPhotonTokenTransfers() {
//	// deploy photon ERC20 token contact
//	var photonERC20Addr string
//	s.Run("deploy_photon_erc20", func() {
//		photonERC20Addr = s.deployERC20Token("photon")
//	})
//
//	// send 100 photon tokens from Umee to Ethereum
//	s.Run("send_photon_tokens_to_eth", func() {
//		ethRecipient := s.chain.validators[1].ethereumKey.address
//		s.sendFromUmeeToEth(0, ethRecipient, "100photon", "10photon", "3photon")
//
//		umeeEndpoint := fmt.Sprintf("http://%s", s.valResources[0].GetHostPort("1317/tcp"))
//		fromAddr := s.chain.validators[0].keyInfo.GetAddress()
//
//		// require the sender's (validator) balance decreased
//		balance, err := queryUmeeDenomBalance(umeeEndpoint, fromAddr.String(), "photon")
//		s.Require().NoError(err)
//		s.Require().Equal(99999999887, balance)
//
//		expEthBalance := 100
//		ethEndpoint := fmt.Sprintf("http://%s", s.ethResource.GetHostPort("8545/tcp"))
//
//		// require the Ethereum recipient balance increased
//		s.Require().Eventually(
//			func() bool {
//				ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
//				defer cancel()
//
//				b, err := queryEthTokenBalance(ctx, ethEndpoint, photonERC20Addr, ethRecipient)
//				if err != nil {
//					return false
//				}
//
//				return b == expEthBalance
//			},
//			2*time.Minute,
//			5*time.Second,
//		)
//	})
//
//	// send 100 photon tokens from Ethereum back to Umee
//	s.Run("send_photon_tokens_from_eth", func() {
//		s.sendFromEthToUmee(1, photonERC20Addr, s.chain.validators[0].keyInfo.GetAddress().String(), "100")
//
//		umeeEndpoint := fmt.Sprintf("http://%s", s.valResources[0].GetHostPort("1317/tcp"))
//		toAddr := s.chain.validators[0].keyInfo.GetAddress()
//		expBalance := 99999999987
//
//		// require the original sender's (validator) balance increased
//		s.Require().Eventually(
//			func() bool {
//				b, err := queryUmeeDenomBalance(umeeEndpoint, toAddr.String(), "photon")
//				if err != nil {
//					return false
//				}
//
//				return b == expBalance
//			},
//			2*time.Minute,
//			5*time.Second,
//		)
//	})
//}

func (s *IntegrationTestSuite) TestRebalance() {
	s.Run("Bring up chain, and submit a re-balance", func() {

	})
}
