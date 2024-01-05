package volunteer

import (
	"github.com/gogo/protobuf/proto"
	"github.com/xpladev/xpla.go/core"
	"github.com/xpladev/xpla.go/types"
	"github.com/xpladev/xpla.go/util"

	volunteertypes "github.com/xpladev/xpla/x/volunteer/types"
)

var out []byte
var res proto.Message
var err error

// Query client for volunteer module.
func QueryVolunteer(i core.QueryClient) (string, error) {
	if i.QueryType == types.QueryGrpc {
		return queryByGrpcVolunteer(i)
	} else {
		return queryByLcdVolunteer(i)
	}
}

func queryByGrpcVolunteer(i core.QueryClient) (string, error) {
	queryClient := volunteertypes.NewQueryClient(i.Ixplac.GetGrpcClient())

	switch {
	// Volunteer validators
	case i.Ixplac.GetMsgType() == VolunteerQueryValidatorsMsgType:
		convertMsg := i.Ixplac.GetMsg().(volunteertypes.QueryVolunteerValidatorsRequest)
		res, err = queryClient.VolunteerValidators(
			i.Ixplac.GetContext(),
			&convertMsg,
		)
		if err != nil {
			return "", i.Ixplac.GetLogger().Err(types.ErrWrap(types.ErrGrpcRequest, err))
		}

	default:
		return "", i.Ixplac.GetLogger().Err(types.ErrWrap(types.ErrInvalidMsgType, i.Ixplac.GetMsgType()))
	}

	out, err = core.PrintProto(i, res)
	if err != nil {
		return "", i.Ixplac.GetLogger().Err(err)
	}

	return string(out), nil
}

const (
	volunteerQueryValidatorsLabel = "validators"
)

func queryByLcdVolunteer(i core.QueryClient) (string, error) {
	url := "/xpla/volunteer/v1beta1/"

	switch {
	// Skating validator
	case i.Ixplac.GetMsgType() == VolunteerQueryValidatorsMsgType:
		url = url + util.MakeQueryLabels(volunteerQueryValidatorsLabel)

	default:
		return "", i.Ixplac.GetLogger().Err(types.ErrWrap(types.ErrInvalidMsgType, i.Ixplac.GetMsgType()))
	}

	i.Ixplac.GetHttpMutex().Lock()
	out, err := util.CtxHttpClient("GET", i.Ixplac.GetLcdURL()+url, nil, i.Ixplac.GetContext())
	if err != nil {
		i.Ixplac.GetHttpMutex().Unlock()
		return "", i.Ixplac.GetLogger().Err(err)
	}
	i.Ixplac.GetHttpMutex().Unlock()

	return string(out), nil
}
