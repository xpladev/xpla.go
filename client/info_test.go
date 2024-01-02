package client_test

import (
	"github.com/xpladev/xpla.go/provider"
	"github.com/xpladev/xpla.go/types"
	"github.com/xpladev/xpla.go/util"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/authz"
)

func (s *ClientTestSuite) TestLoadAccount() {
	val := s.network.Validators[0].Address

	for i, api := range s.apis {
		if i == 0 {
			s.xplac.WithURL(api)
		} else {
			s.xplac.WithGrpc(api)
		}

		res, err := s.xplac.LoadAccount(val)
		s.Require().NoError(err)
		s.Require().Equal(val.String(), res.GetAddress().String())
	}
	s.xplac = provider.ResetXplac(s.xplac)
}

func (s *ClientTestSuite) TestSimulate() {
	val1 := s.network.Validators[0].Address
	s.xplac.WithPrivateKey(s.network.Validators[4].AdditionalAccount.PrivKey)

	for i, api := range s.apis {
		if i == 0 {
			s.xplac.WithURL(api)
		} else {
			s.xplac.WithGrpc(api)
		}

		xplac := s.xplac
		account, err := xplac.LoadAccount(sdk.AccAddress(xplac.GetPrivateKey().PubKey().Address()))
		s.Require().NoError(err)

		xplac.WithAccountNumber(util.FromUint64ToString(account.GetAccountNumber()))
		xplac.WithSequence(util.FromUint64ToString(account.GetSequence()))

		authzGrantMsg := types.AuthzGrantMsg{
			Granter:           s.network.Validators[4].AdditionalAccount.Address.String(),
			Grantee:           val1.String(),
			AuthorizationType: "send",
			SpendLimit:        "1000",
		}

		xplac = s.xplac.AuthzGrant(authzGrantMsg)
		s.Require().NoError(xplac.GetErr())

		builder := xplac.GetEncoding().TxConfig.NewTxBuilder()

		convertMsg, ok := xplac.GetMsg().(authz.MsgGrant)
		s.Require().True(ok)
		builder.SetMsgs(&convertMsg)

		_, err = xplac.Simulate(builder)
		s.Require().NoError(err)

	}
	s.xplac = provider.ResetXplac(s.xplac)
}
