package auth

const (
	AuthModule                  = "auth"
	AuthQueryParamsMsgType      = "query-auth-params"
	AuthQueryAccAddressMsgType  = "query-account"
	AuthQueryAccountsMsgType    = "query-accounts"
	AuthQueryTxsByEventsMsgType = "query-txs-by-events"
	AuthQueryTxMsgType          = "query-tx"
)

type QueryTxsByEventParseMsg struct {
	TmEvents []string
	Page     int
	Limit    int
}

type QueryTxParseMsg struct {
	TmEvents []string
	TxType   string
}
