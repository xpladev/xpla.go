package provider

// Reset XplaClient.
// Remove recorded all parameters.
func ResetXplac(xplac XplaClient) XplaClient {
	return ResetModuleAndMsgXplac(xplac).
		WithOptions(Options{}).
		WithErr(nil)
}

// Reset XplaClient with removing module name and message.
func ResetModuleAndMsgXplac(xplac XplaClient) XplaClient {
	return xplac.
		WithModule("").
		WithMsgType("").
		WithMsg(nil)
}
