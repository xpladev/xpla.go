package mint_test

import mmint "github.com/xpladev/xpla.go/core/mint"

func (s *IntegrationTestSuite) TestMint() {
	// mint params
	s.xplac.MintParams()

	makeQueryMintParamsMsg, err := mmint.MakeQueryMintParamsMsg()
	s.Require().NoError(err)

	s.Require().Equal(makeQueryMintParamsMsg, s.xplac.GetMsg())
	s.Require().Equal(mmint.MintModule, s.xplac.GetModule())
	s.Require().Equal(mmint.MintQueryMintParamsMsgType, s.xplac.GetMsgType())

	// inflation
	s.xplac.Inflation()

	makeQueryInflationMsg, err := mmint.MakeQueryInflationMsg()
	s.Require().NoError(err)

	s.Require().Equal(makeQueryInflationMsg, s.xplac.GetMsg())
	s.Require().Equal(mmint.MintModule, s.xplac.GetModule())
	s.Require().Equal(mmint.MintQueryInflationMsgType, s.xplac.GetMsgType())

	// annual provisions
	s.xplac.AnnualProvisions()

	makeQueryAnnualProvisionsMsg, err := mmint.MakeQueryAnnualProvisionsMsg()
	s.Require().NoError(err)

	s.Require().Equal(makeQueryAnnualProvisionsMsg, s.xplac.GetMsg())
	s.Require().Equal(mmint.MintModule, s.xplac.GetModule())
	s.Require().Equal(mmint.MintQueryAnnualProvisionsMsgType, s.xplac.GetMsgType())
}
