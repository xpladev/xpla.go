package util

import (
	"errors"
	"fmt"

	xgoerrors "github.com/xpladev/xpla.go/types/errors"
)

func LogInfo(log ...interface{}) {
	fmt.Println(ToStringTrim(log, ""))
}

func LogErr(errType xgoerrors.XGoError, errDesc ...interface{}) error {
	return logErrReturn("code", errType.ErrCode(), ":", errType.Desc(), "-", errDesc)
}

func logErrReturn(log ...interface{}) error {
	return errors.New(ToStringTrim(log, ""))
}
