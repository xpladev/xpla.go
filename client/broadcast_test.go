package client_test

import (
	"strings"

	"github.com/xpladev/xpla.go/client"
	"github.com/xpladev/xpla.go/types"
	"github.com/xpladev/xpla.go/util"

	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	"github.com/gogo/protobuf/jsonpb"
)

func (s *ClientTestSuite) TestBroadcast() {
	from := s.accounts[0]
	to := s.accounts[1]
	s.xplac.WithPrivateKey(s.accounts[0].PrivKey)

	for i, api := range s.apis {
		if i == 0 {
			s.xplac.WithURL(api)
		} else {
			s.xplac.WithGrpc(api)
			nowSeq, err := util.FromStringToInt(s.xplac.GetSequence())
			s.Require().NoError(err)
			s.xplac.WithSequence(util.FromIntToString(nowSeq + 1))
		}

		// check before send
		bankBalancesMsg := types.BankBalancesMsg{
			Address: to.Address.String(),
		}
		beforeToRes, err := s.xplac.BankBalances(bankBalancesMsg).Query()
		s.Require().NoError(err)

		var beforeQueryAllBalancesResponse banktypes.QueryAllBalancesResponse
		jsonpb.Unmarshal(strings.NewReader(beforeToRes), &beforeQueryAllBalancesResponse)

		// broadcast transaction - bank send
		bankSendMsg := types.BankSendMsg{
			FromAddress: from.Address.String(),
			ToAddress:   to.Address.String(),
			Amount:      testSendAmount,
		}
		txbytes, err := s.xplac.BankSend(bankSendMsg).CreateAndSignTx()
		s.Require().NoError(err)

		_, err = s.xplac.Broadcast(txbytes)
		s.Require().NoError(err)
		s.Require().NoError(s.network.WaitForNextBlock())

		// check after send
		bankBalancesMsg = types.BankBalancesMsg{
			Address: to.Address.String(),
		}
		afterToRes, err := s.xplac.BankBalances(bankBalancesMsg).Query()
		s.Require().NoError(err)

		var afterQueryAllBalancesResponse banktypes.QueryAllBalancesResponse
		jsonpb.Unmarshal(strings.NewReader(afterToRes), &afterQueryAllBalancesResponse)

		s.Require().Equal(
			testSendAmount,
			afterQueryAllBalancesResponse.Balances[0].Amount.Sub(beforeQueryAllBalancesResponse.Balances[0].Amount).String(),
		)
	}
	s.xplac = client.ResetXplac(s.xplac)
}

func (s *ClientTestSuite) TestBroadcastMode() {
	invalidBroadcastMode := "invalid-mode"

	from := s.accounts[0]
	to := s.accounts[1]

	s.xplac.
		WithURL(s.apis[0]).
		WithPrivateKey(s.accounts[0].PrivKey)

	modes := []string{"block", "async", "sync", "", invalidBroadcastMode}

	for _, mode := range modes {
		s.xplac.WithBroadcastMode(mode)

		bankSendMsg := types.BankSendMsg{
			FromAddress: from.Address.String(),
			ToAddress:   to.Address.String(),
			Amount:      testSendAmount,
		}
		txbytes, err := s.xplac.BankSend(bankSendMsg).CreateAndSignTx()
		s.Require().NoError(err)

		// if empty mode or invalid mode is changed to "sync"
		_, err = s.xplac.Broadcast(txbytes)
		s.Require().NoError(err)
		s.Require().NoError(s.network.WaitForNextBlock())

		nowSeq, err := util.FromStringToInt(s.xplac.GetSequence())
		s.Require().NoError(err)
		s.xplac.WithSequence(util.FromIntToString(nowSeq + 1))
	}
	s.xplac = client.ResetXplac(s.xplac)
}

func (s *ClientTestSuite) TestBroadcastEVM() {
	from := s.accounts[0]
	to := s.accounts[1]
	s.xplac.WithPrivateKey(s.accounts[0].PrivKey).
		WithURL(s.apis[0]).
		WithEvmRpc("http://" + s.network.Validators[0].AppConfig.JSONRPC.Address)

	// check before send
	bankBalancesMsg := types.BankBalancesMsg{
		Address: to.Address.String(),
	}
	beforeToRes, err := s.xplac.BankBalances(bankBalancesMsg).Query()
	s.Require().NoError(err)

	var beforeQueryAllBalancesResponse banktypes.QueryAllBalancesResponse
	jsonpb.Unmarshal(strings.NewReader(beforeToRes), &beforeQueryAllBalancesResponse)

	// broadcast transaction - evm send coin
	sendCoinMsg := types.SendCoinMsg{
		FromAddress: from.PubKey.Address().String(),
		ToAddress:   to.PubKey.Address().String(),
		Amount:      testSendAmount,
	}
	txbytes, err := s.xplac.EvmSendCoin(sendCoinMsg).CreateAndSignTx()
	s.Require().NoError(err)

	_, err = s.xplac.Broadcast(txbytes)
	s.Require().NoError(err)
	s.Require().NoError(s.network.WaitForNextBlock())

	// check after send
	bankBalancesMsg = types.BankBalancesMsg{
		Address: to.Address.String(),
	}
	afterToRes, err := s.xplac.BankBalances(bankBalancesMsg).Query()
	s.Require().NoError(err)

	var afterQueryAllBalancesResponse banktypes.QueryAllBalancesResponse
	jsonpb.Unmarshal(strings.NewReader(afterToRes), &afterQueryAllBalancesResponse)

	s.Require().Equal(
		testSendAmount,
		afterQueryAllBalancesResponse.Balances[0].Amount.Sub(beforeQueryAllBalancesResponse.Balances[0].Amount).String(),
	)
	s.xplac = client.ResetXplac(s.xplac)
}
