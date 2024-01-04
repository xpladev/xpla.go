package client

import (
	"github.com/xpladev/xpla.go/types"
	"github.com/xpladev/xpla.go/util"

	cmclient "github.com/cosmos/cosmos-sdk/client"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdktx "github.com/cosmos/cosmos-sdk/types/tx"
	"github.com/cosmos/cosmos-sdk/types/tx/signing"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
)

const (
	userInfoUrl  = "/cosmos/auth/v1beta1/accounts/"
	simulateUrl  = "/cosmos/tx/v1beta1/simulate"
	broadcastUrl = "/cosmos/tx/v1beta1/txs"
)

// LoadAccount gets the account info by AccAddress
// If xpla client has gRPC client, query account information by using gRPC
func (xplac *xplaClient) LoadAccount(address sdk.AccAddress) (res authtypes.AccountI, err error) {

	if xplac.GetGrpcUrl() == "" {
		xplac.GetHttpMutex().Lock()
		out, err := util.CtxHttpClient("GET", xplac.GetLcdURL()+userInfoUrl+address.String(), nil, xplac.GetContext())
		if err != nil {
			xplac.GetHttpMutex().Unlock()
			return nil, xplac.GetLogger().Err(err)
		}
		xplac.GetHttpMutex().Unlock()

		var response authtypes.QueryAccountResponse
		err = xplac.GetEncoding().Codec.UnmarshalJSON(out, &response)
		if err != nil {
			return nil, xplac.GetLogger().Err(types.ErrWrap(types.ErrFailedToUnmarshal, err))
		}
		return response.Account.GetCachedValue().(authtypes.AccountI), nil

	} else {
		queryClient := authtypes.NewQueryClient(xplac.GetGrpcClient())
		queryAccountRequest := authtypes.QueryAccountRequest{
			Address: address.String(),
		}
		response, err := queryClient.Account(xplac.GetContext(), &queryAccountRequest)
		if err != nil {
			return nil, xplac.GetLogger().Err(types.ErrWrap(types.ErrGrpcRequest, err))
		}

		var newAccount authtypes.AccountI
		err = xplac.GetEncoding().InterfaceRegistry.UnpackAny(response.Account, &newAccount)
		if err != nil {
			return nil, xplac.GetLogger().Err(types.ErrWrap(types.ErrParse, err))
		}

		return newAccount, nil
	}
}

// Simulate tx and get response
// If xpla client has gRPC client, query simulation by using gRPC
func (xplac *xplaClient) Simulate(txbuilder cmclient.TxBuilder) (*sdktx.SimulateResponse, error) {
	seq, err := util.FromStringToUint64(xplac.GetSequence())
	if err != nil {
		return nil, xplac.GetLogger().Err(types.ErrWrap(types.ErrConvert, err))
	}

	pubKey := xplac.GetPublicKey()
	if xplac.GetPublicKey() == nil {
		if xplac.GetFromAddress() != nil {
			accountInfo, err := xplac.LoadAccount(xplac.GetFromAddress())
			if err != nil {
				return nil, err
			}

			pubKey = accountInfo.GetPubKey()

		} else {
			return nil, xplac.GetLogger().Err(types.ErrWrap(types.ErrInvalidRequest, "cannot be simulated without the public key."))
		}
	}

	sig := signing.SignatureV2{
		PubKey: pubKey,
		Data: &signing.SingleSignatureData{
			SignMode: xplac.GetSignMode(),
		},
		Sequence: seq,
	}

	if err := txbuilder.SetSignatures(sig); err != nil {
		return nil, xplac.GetLogger().Err(types.ErrWrap(types.ErrParse, err))
	}

	sdkTx := txbuilder.GetTx()
	txBytes, err := xplac.GetEncoding().TxConfig.TxEncoder()(sdkTx)
	if err != nil {
		return nil, xplac.GetLogger().Err(types.ErrWrap(types.ErrParse, err))
	}

	if xplac.GetGrpcUrl() == "" {
		reqBytes, err := xplac.GetEncoding().Codec.MarshalJSON(&sdktx.SimulateRequest{
			TxBytes: txBytes,
		})
		if err != nil {
			return nil, xplac.GetLogger().Err(types.ErrWrap(types.ErrFailedToMarshal, err))
		}

		xplac.GetHttpMutex().Lock()
		out, err := util.CtxHttpClient("POST", xplac.GetLcdURL()+simulateUrl, reqBytes, xplac.GetContext())
		if err != nil {
			xplac.GetHttpMutex().Unlock()
			return nil, xplac.GetLogger().Err(err)
		}
		xplac.GetHttpMutex().Unlock()

		var response sdktx.SimulateResponse
		err = xplac.GetEncoding().Codec.UnmarshalJSON(out, &response)
		if err != nil {
			return nil, xplac.GetLogger().Err(types.ErrWrap(types.ErrFailedToUnmarshal, err))
		}

		return &response, nil
	} else {
		serviceClient := sdktx.NewServiceClient(xplac.GetGrpcClient())
		simulateRequest := sdktx.SimulateRequest{
			TxBytes: txBytes,
		}

		response, err := serviceClient.Simulate(xplac.GetContext(), &simulateRequest)
		if err != nil {
			return nil, xplac.GetLogger().Err(types.ErrWrap(types.ErrGrpcRequest, err))
		}

		return response, nil
	}
}
