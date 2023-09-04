package client

import (
	"github.com/xpladev/xpla.go/util/testutil"
)

func NewTestXplaClient() *XplaClient {
	return NewXplaClient(testutil.TestChainId)
}

func ResetXplac(xplac *XplaClient) *XplaClient {
	return ResetModuleAndMsgXplac(xplac).
		WithOptions(Options{}).
		WithErr(nil)
}

func ResetModuleAndMsgXplac(xplac *XplaClient) *XplaClient {
	return xplac.
		WithModule("").
		WithMsgType("").
		WithMsg(nil)
}

func (xplac *XplaClient) EncodedTxbytesToJsonTx(txbytes []byte) ([]byte, error) {
	sdkTx, err := xplac.GetEncoding().TxConfig.TxDecoder()(txbytes)
	if err != nil {
		return nil, err
	}
	jsonTx, err := xplac.GetEncoding().TxConfig.TxJSONEncoder()(sdkTx)
	if err != nil {
		return nil, err
	}
	return jsonTx, nil
}
