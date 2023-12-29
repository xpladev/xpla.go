package volunteer_test

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	mvolunteer "github.com/xpladev/xpla.go/core/volunteer"
	"github.com/xpladev/xpla.go/types"
)

func (s *IntegrationTestSuite) TestVolunteerTx() {
	account0 := s.network.Validators[0].AdditionalAccount

	s.xplac.WithPrivateKey(account0.PrivKey)
	tmpVal := sdk.ValAddress(account0.Address)

	// register volunteer validator.
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

	makeRegisterVolunteerValidatorMsg, err := mvolunteer.MakeRegisterVolunteerValidatorMsg(registerVolunteerValidatorMsg, s.xplac.GetEncoding(), s.xplac.GetFromAddress())
	s.Require().NoError(err)

	s.Require().Equal(makeRegisterVolunteerValidatorMsg, s.xplac.GetMsg())
	s.Require().Equal(mvolunteer.VolunteerModule, s.xplac.GetModule())
	s.Require().Equal(mvolunteer.VolunteerRegisterVolunteerValidatorMsgType, s.xplac.GetMsgType())

	_, err = s.xplac.RegisterVolunteerValidator(registerVolunteerValidatorMsg).CreateAndSignTx()
	s.Require().NoError(err)

	// unregister volunteer validator.
	unregisterVolunteerValidatorMsg := types.UnregisterVolunteerValidatorMsg{
		Title:       "register volunteer validator",
		Description: "register description",
		Deposit:     "1000",
		ValAddress:  tmpVal.String(),
	}
	s.xplac.UnregisterVolunteerValidator(unregisterVolunteerValidatorMsg)

	makeUnregisterVolunteerValidatorMsg, err := mvolunteer.MakeUnregisterVolunteerValidatorMsg(unregisterVolunteerValidatorMsg, s.xplac.GetEncoding(), s.xplac.GetFromAddress())
	s.Require().NoError(err)

	s.Require().Equal(makeUnregisterVolunteerValidatorMsg, s.xplac.GetMsg())
	s.Require().Equal(mvolunteer.VolunteerModule, s.xplac.GetModule())
	s.Require().Equal(mvolunteer.VolunteerUnregisterVolunteerValidatorMsgType, s.xplac.GetMsgType())

	_, err = s.xplac.UnregisterVolunteerValidator(unregisterVolunteerValidatorMsg).CreateAndSignTx()
	s.Require().NoError(err)
}

func (s *IntegrationTestSuite) TestVolunteer() {
	// query validators
	s.xplac.QueryVolunteerValidators()

	makeQueryVolunteerValidatorsMsg, err := mvolunteer.MakeQueryVolunteerValidatorsMsg()
	s.Require().NoError(err)

	s.Require().Equal(makeQueryVolunteerValidatorsMsg, s.xplac.GetMsg())
	s.Require().Equal(mvolunteer.VolunteerModule, s.xplac.GetModule())
	s.Require().Equal(mvolunteer.VolunteerQueryValidatorsMsgType, s.xplac.GetMsgType())
}
