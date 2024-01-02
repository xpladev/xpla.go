package volunteer_test

import (
	"math/rand"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/xpladev/xpla.go/core/volunteer"
	"github.com/xpladev/xpla.go/provider"
	"github.com/xpladev/xpla.go/types"
	"github.com/xpladev/xpla.go/util/testutil"
)

func (s *IntegrationTestSuite) TestCoreModule() {
	src := rand.NewSource(1)
	r := rand.New(src)
	accounts := testutil.RandomAccounts(r, 2)
	s.xplac.WithPrivateKey(accounts[0].PrivKey)

	c := volunteer.NewCoreModule()

	// test get name
	s.Require().Equal(volunteer.VolunteerModule, c.Name())

	// test tx
	var testMsg interface{}
	txBuilder := s.xplac.GetEncoding().TxConfig.NewTxBuilder()

	// register volunteer validator
	registerVolunteerValidatorMsg := types.RegisterVolunteerValidatorMsg{
		Title:       "register volunteer validator",
		Description: "register description",
		Deposit:     "1000",
		Amount:      "10000",
		ValPubKey:   `{"@type": "/cosmos.crypto.ed25519.PubKey", "key": "2z2yttKfEsLQyQnHYdgKEuky9zB3gscxapn9IyexxWk="}`,
		Moniker:     "volun moniker",
		Identity:    "volun identity",
		Website:     "volun website",
		Security:    "volun security",
		Details:     "volun details",
	}
	s.xplac.RegisterVolunteerValidator(registerVolunteerValidatorMsg)

	makeRegisterVolunteerValidatorMsg, err := volunteer.MakeRegisterVolunteerValidatorMsg(registerVolunteerValidatorMsg, s.xplac.GetEncoding(), s.xplac.GetFromAddress())
	s.Require().NoError(err)

	testMsg = makeRegisterVolunteerValidatorMsg
	txBuilder, err = c.NewTxRouter(txBuilder, volunteer.VolunteerRegisterVolunteerValidatorMsgType, testMsg)
	s.Require().NoError(err)
	s.Require().Equal(&makeRegisterVolunteerValidatorMsg, txBuilder.GetTx().GetMsgs()[0])

	// unregister volunteer validator
	tmpVal := sdk.ValAddress(accounts[0].Address)
	unregisterVolunteerValidatorMsg := types.UnregisterVolunteerValidatorMsg{
		Title:       "register volunteer validator",
		Description: "register description",
		Deposit:     "1000",
		ValAddress:  tmpVal.String(),
	}

	makeUnregisterVolunteerValidatorMsg, err := volunteer.MakeUnregisterVolunteerValidatorMsg(unregisterVolunteerValidatorMsg, s.xplac.GetEncoding(), s.xplac.GetFromAddress())
	s.Require().NoError(err)

	testMsg = makeUnregisterVolunteerValidatorMsg
	txBuilder, err = c.NewTxRouter(txBuilder, volunteer.VolunteerUnregisterVolunteerValidatorMsgType, testMsg)
	s.Require().NoError(err)
	s.Require().Equal(&makeUnregisterVolunteerValidatorMsg, txBuilder.GetTx().GetMsgs()[0])

	s.xplac = provider.ResetXplac(s.xplac)
}
