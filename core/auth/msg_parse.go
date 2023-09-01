package auth

import (
	"fmt"
	"strings"

	"github.com/xpladev/xpla.go/types"
	"github.com/xpladev/xpla.go/types/errors"
	"github.com/xpladev/xpla.go/util"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/rest"
	tmtypes "github.com/tendermint/tendermint/types"
)

// Parsing - transaction by evnets
func parseTxsByEventsArgs(txsByEventsMsg types.QueryTxsByEventsMsg) (QueryTxsByEventParseMsg, error) {
	eventFormat := "{eventType}.{eventAttribute}={value}"
	eventsRaw := txsByEventsMsg.Events
	eventsStr := strings.Trim(eventsRaw, "'")

	if txsByEventsMsg.Page == "" {
		txsByEventsMsg.Page = util.FromIntToString(rest.DefaultPage)
	}
	if txsByEventsMsg.Limit == "" {
		txsByEventsMsg.Limit = util.FromIntToString(rest.DefaultLimit)
	}

	var events []string
	if strings.Contains(eventsStr, "&") {
		events = strings.Split(eventsStr, "&")
	} else {
		events = append(events, eventsStr)
	}

	var tmEvents []string

	for _, event := range events {
		if !strings.Contains(event, "=") {
			return QueryTxsByEventParseMsg{}, util.LogErr(errors.ErrInvalidRequest, "invalid event; event", event, "should be of the format:", eventFormat)
		} else if strings.Count(event, "=") > 1 {
			return QueryTxsByEventParseMsg{}, util.LogErr(errors.ErrInvalidRequest, "invalid event; event", event, "should be of the format:", eventFormat)
		}

		tokens := strings.Split(event, "=")
		if tokens[0] == tmtypes.TxHeightKey {
			event = fmt.Sprintf("%s=%s", tokens[0], tokens[1])
		} else {
			event = fmt.Sprintf("%s='%s'", tokens[0], tokens[1])
		}

		tmEvents = append(tmEvents, event)
	}

	pageInt, err := util.FromStringToInt(txsByEventsMsg.Page)
	if err != nil {
		return QueryTxsByEventParseMsg{}, err
	}
	limitInt, err := util.FromStringToInt(txsByEventsMsg.Limit)
	if err != nil {
		return QueryTxsByEventParseMsg{}, err
	}

	queryTxsByEventParseMsg := QueryTxsByEventParseMsg{
		TmEvents: tmEvents,
		Page:     pageInt,
		Limit:    limitInt,
	}

	return queryTxsByEventParseMsg, nil
}

// Parsing - transaction
func parseQueryTxArgs(queryTxMsg types.QueryTxMsg) (QueryTxParseMsg, error) {
	var queryTxParseMsg QueryTxParseMsg

	if queryTxMsg.Type == "" || queryTxMsg.Type == "hash" {
		if queryTxMsg.Value == "" {
			return QueryTxParseMsg{}, util.LogErr(errors.ErrInvalidRequest, "argument should be a tx hash")
		}

		queryTxParseMsg.TmEvents = []string{queryTxMsg.Value}
		queryTxParseMsg.TxType = "hash"

		return queryTxParseMsg, nil

	} else if queryTxMsg.Type == "signature" {
		if queryTxMsg.Value == "" {
			return QueryTxParseMsg{}, fmt.Errorf("argument should be comma-separated signatures")
		}
		sigParts := strings.Split(queryTxMsg.Value, ",")

		tmEvents := make([]string, len(sigParts))
		for i, sig := range sigParts {
			tmEvents[i] = fmt.Sprintf("%s.%s='%s'", sdk.EventTypeTx, sdk.AttributeKeySignature, sig)
		}

		queryTxParseMsg.TmEvents = tmEvents
		queryTxParseMsg.TxType = queryTxMsg.Type

		return queryTxParseMsg, nil

	} else if queryTxMsg.Type == "acc_seq" {
		if queryTxMsg.Value == "" {
			return QueryTxParseMsg{}, util.LogErr(errors.ErrInvalidRequest, "`acc_seq` type takes an argument '<addr>/<seq>'")
		}

		tmEvents := []string{
			fmt.Sprintf("%s.%s='%s'", sdk.EventTypeTx, sdk.AttributeKeyAccountSequence, queryTxMsg.Value),
		}

		queryTxParseMsg.TmEvents = tmEvents
		queryTxParseMsg.TxType = queryTxMsg.Type

		return queryTxParseMsg, nil

	} else {
		return QueryTxParseMsg{}, util.LogErr(errors.ErrInvalidMsgType, "unknown type (hash|signature|acc_seq)")
	}
}
