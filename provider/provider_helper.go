package provider

func ResetXplac(xplac XplaClient) XplaClient {
	return ResetModuleAndMsgXplac(xplac).
		WithOptions(Options{}).
		WithErr(nil)
}

func ResetModuleAndMsgXplac(xplac XplaClient) XplaClient {
	return xplac.
		WithModule("").
		WithMsgType("").
		WithMsg(nil)
}
