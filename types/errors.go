package types

import (
	"fmt"
	"strings"
)

type XGoError struct {
	errCode uint64
	desc    string
}

// Return error code and message generating on the xpla.go.
var (
	ErrInvalidMsgType      = new(1, "invalid msg type")
	ErrInvalidRequest      = new(2, "invalid request")
	ErrNotSatisfiedOptions = new(3, "not satisfied option parameter")
	ErrFailedToMarshal     = new(4, "failed to marshal")
	ErrFailedToUnmarshal   = new(5, "failed to unmarshal")
	ErrNotSupport          = new(6, "not support xpla.go")
	ErrNotFound            = new(7, "not found")
	ErrTxFailed            = new(8, "tx failed")
	ErrInsufficientParams  = new(9, "insufficient parameters")
	ErrKeyNotFound         = new(10, "key not found")
	ErrAccountNotMatch     = new(11, "account not match")
	ErrHttpRequest         = new(12, "HTTP request error")
	ErrGrpcRequest         = new(13, "gRPC request error")
	ErrEvmRpcRequest       = new(14, "EVM RPC request error")
	ErrRpcRequest          = new(15, "RPC request error")
	ErrConvert             = new(16, "convert type error")
	ErrParse               = new(17, "parse error")
	ErrSdkClient           = new(18, "cosmos sdk client set error")
	ErrAlreadyExist        = new(19, "already exist")
	ErrCannotRead          = new(20, "cannot read")
)

func new(errCode uint64, desc string) XGoError {
	var xErr XGoError
	xErr.errCode = errCode
	xErr.desc = desc

	return xErr
}

func (x XGoError) ErrCode() uint64 {
	return x.errCode
}

func (x XGoError) Desc() string {
	return x.desc
}

func ErrWrap(errType XGoError, errDesc ...interface{}) error {
	return wrapper(
		"code",
		errType.ErrCode(), ":",
		errType.Desc(), "-", errDesc,
	)
}

func wrapper(log ...interface{}) error {
	return fmt.Errorf(ToStringTrim(log, ""))
}

func ToStringTrim(value interface{}, defaultValue string) string {
	s := fmt.Sprintf("%v", value)
	s = s[1 : len(s)-1]
	str := strings.TrimSpace(s)
	if str == "" {
		return defaultValue
	} else {
		return str
	}
}
