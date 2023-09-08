package feegrant_test

import (
	mfeegrant "github.com/xpladev/xpla.go/core/feegrant"
	"github.com/xpladev/xpla.go/types"
	"github.com/xpladev/xpla.go/util/testutil"
)

func (s *IntegrationTestSuite) TestFeegrantTx() {
	s.xplac.WithPrivateKey(s.accounts[0].PrivKey)
	// feegrant
	feeGrantMsg := types.FeeGrantMsg{
		Granter:    s.accounts[0].Address.String(),
		Grantee:    s.accounts[1].Address.String(),
		SpendLimit: "1000",
		// Period:      "3600",
		// PeriodLimit: "10",
		Expiration: "2100-01-01T23:59:59+00:00",
	}
	s.xplac.FeeGrant(feeGrantMsg)

	makeFeeGrantMsg, err := mfeegrant.MakeFeeGrantMsg(feeGrantMsg, s.xplac.GetPrivateKey())
	s.Require().NoError(err)

	s.Require().Equal(makeFeeGrantMsg, s.xplac.GetMsg())
	s.Require().Equal(mfeegrant.FeegrantModule, s.xplac.GetModule())
	s.Require().Equal(mfeegrant.FeegrantGrantMsgType, s.xplac.GetMsgType())

	feegrantFeegrantTxbytes, err := s.xplac.FeeGrant(feeGrantMsg).CreateAndSignTx()
	s.Require().NoError(err)

	feegrantFeegrantJsonTxbytes, err := s.xplac.EncodedTxbytesToJsonTx(feegrantFeegrantTxbytes)
	s.Require().NoError(err)
	s.Require().Equal(testutil.FeegrantFeegrantTxTemplates, string(feegrantFeegrantJsonTxbytes))

	// revoke feegrant
	revokeFeeGrantMsg := types.RevokeFeeGrantMsg{
		Granter: s.accounts[0].Address.String(),
		Grantee: s.accounts[1].Address.String(),
	}
	s.xplac.RevokeFeeGrant(revokeFeeGrantMsg)

	makeRevokeFeeGrantMsg, err := mfeegrant.MakeRevokeFeeGrantMsg(revokeFeeGrantMsg, s.xplac.GetPrivateKey())
	s.Require().NoError(err)

	s.Require().Equal(makeRevokeFeeGrantMsg, s.xplac.GetMsg())
	s.Require().Equal(mfeegrant.FeegrantModule, s.xplac.GetModule())
	s.Require().Equal(mfeegrant.FeegrantRevokeGrantMsgType, s.xplac.GetMsgType())

	feegrantRevokeFeegrantTxbytes, err := s.xplac.RevokeFeeGrant(revokeFeeGrantMsg).CreateAndSignTx()
	s.Require().NoError(err)

	feegrantRevokeFeegrantJsonTxbytes, err := s.xplac.EncodedTxbytesToJsonTx(feegrantRevokeFeegrantTxbytes)
	s.Require().NoError(err)
	s.Require().Equal(testutil.FeegrantRevokeFeegrantTxTemplates, string(feegrantRevokeFeegrantJsonTxbytes))
}

func (s *IntegrationTestSuite) TestFeegrant() {
	// feegrant
	queryFeeGrantMsg := types.QueryFeeGrantMsg{
		Grantee: s.accounts[0].Address.String(),
		Granter: s.accounts[1].Address.String(),
	}
	s.xplac.QueryFeeGrants(queryFeeGrantMsg)

	makeQueryFeeGrantMsg, err := mfeegrant.MakeQueryFeeGrantMsg(queryFeeGrantMsg)
	s.Require().NoError(err)

	s.Require().Equal(makeQueryFeeGrantMsg, s.xplac.GetMsg())
	s.Require().Equal(mfeegrant.FeegrantModule, s.xplac.GetModule())
	s.Require().Equal(mfeegrant.FeegrantQueryGrantMsgType, s.xplac.GetMsgType())

	// feegrant by grantee
	queryFeeGrantMsg = types.QueryFeeGrantMsg{
		Grantee: s.accounts[0].Address.String(),
	}
	s.xplac.QueryFeeGrants(queryFeeGrantMsg)

	makeQueryFeeGrantsByGranteeMsg, err := mfeegrant.MakeQueryFeeGrantsByGranteeMsg(queryFeeGrantMsg)
	s.Require().NoError(err)

	s.Require().Equal(makeQueryFeeGrantsByGranteeMsg, s.xplac.GetMsg())
	s.Require().Equal(mfeegrant.FeegrantModule, s.xplac.GetModule())
	s.Require().Equal(mfeegrant.FeegrantQueryGrantsByGranteeMsgType, s.xplac.GetMsgType())

	// feegrant by granter
	queryFeeGrantMsg = types.QueryFeeGrantMsg{
		Granter: s.accounts[1].Address.String(),
	}
	s.xplac.QueryFeeGrants(queryFeeGrantMsg)

	makeQueryFeeGrantsByGranterMsg, err := mfeegrant.MakeQueryFeeGrantsByGranterMsg(queryFeeGrantMsg)
	s.Require().NoError(err)

	s.Require().Equal(makeQueryFeeGrantsByGranterMsg, s.xplac.GetMsg())
	s.Require().Equal(mfeegrant.FeegrantModule, s.xplac.GetModule())
	s.Require().Equal(mfeegrant.FeegrantQueryGrantsByGranterMsgType, s.xplac.GetMsgType())
}
