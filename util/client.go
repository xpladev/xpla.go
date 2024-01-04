package util

import (
	"bytes"
	"context"
	"io"
	"net/http"
	"time"

	cmclient "github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/tx"
	"github.com/cosmos/cosmos-sdk/crypto/keyring"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	"github.com/ethereum/go-ethereum/ethclient"
	erpc "github.com/ethereum/go-ethereum/rpc"
	"github.com/evmos/ethermint/crypto/hd"
	"github.com/xpladev/xpla.go/types"
	"golang.org/x/net/context/ctxhttp"
)

const (
	BackendFile   = "file"
	BackendMemory = "memory"
	BackendTest   = "test"
)

// Provide cosmos sdk client.
func NewClient() (cmclient.Context, error) {
	clientCtx := cmclient.Context{}
	encodingConfig := MakeEncodingConfig()
	clientKeyring, err := NewKeyring(BackendMemory, "")
	if err != nil {
		return cmclient.Context{}, types.ErrWrap(types.ErrKeyNotFound, err)
	}

	clientCtx = clientCtx.
		WithTxConfig(encodingConfig.TxConfig).
		WithCodec(encodingConfig.Codec).
		WithLegacyAmino(encodingConfig.Amino).
		WithInterfaceRegistry(encodingConfig.InterfaceRegistry).
		WithKeyringOptions(hd.EthSecp256k1Option()).
		WithKeyring(clientKeyring).
		WithAccountRetriever(authtypes.AccountRetriever{})

	return clientCtx, nil
}

const (
	DefaultEvmGasLimit         = "100000"
	DefaultEvmQueryGasLimit    = "200000" // Gas is not consumed when querying
	DefaultSolidityValue       = "0"
	DefaultEvmTxReceiptTimeout = 100
)

type EvmClient struct {
	Ctx       context.Context
	Client    *ethclient.Client
	RpcClient *erpc.Client
}

// Make new evm client using RPC URL which normally TCP port number is 8545.
// It supports that sending transaction, contract deployment, executing/querying contract and etc.
func NewEvmClient(evmRpcUrl string, ctx context.Context) (*EvmClient, error) {
	// Target blockchain node URL
	httpDefaultTransport := http.DefaultTransport
	defaultTransport, ok := httpDefaultTransport.(*http.Transport)
	if !ok {
		return nil, types.ErrWrap(types.ErrInvalidRequest, "default transport pointer err")
	}
	defaultTransport.DisableKeepAlives = true

	httpClient := &http.Client{Transport: defaultTransport}
	rpcClient, err := erpc.DialHTTPWithClient(evmRpcUrl, httpClient)
	if err != nil {
		return nil, types.ErrWrap(types.ErrEvmRpcRequest, err)
	}

	ethClient := ethclient.NewClient(rpcClient)

	return &EvmClient{ctx, ethClient, rpcClient}, nil
}

// Provide cosmos sdk keyring
func NewKeyring(backendType string, keyringPath string) (keyring.Keyring, error) {
	switch {
	case backendType == BackendMemory:
		k, err := keyring.New(
			types.XplaToolDefaultName,
			keyring.BackendMemory,
			"",
			nil,
			hd.EthSecp256k1Option(),
		)
		if err != nil {
			return nil, types.ErrWrap(types.ErrKeyNotFound, err)
		}

		return k, nil

	case backendType == BackendFile:
		k, err := keyring.New(
			types.XplaToolDefaultName,
			keyring.BackendFile,
			keyringPath,
			nil,
			hd.EthSecp256k1Option(),
		)
		if err != nil {
			return nil, types.ErrWrap(types.ErrKeyNotFound, err)
		}

		return k, nil

	case backendType == BackendTest:
		k, err := keyring.New(
			types.XplaToolDefaultName,
			keyring.BackendTest,
			keyringPath,
			nil,
			hd.EthSecp256k1Option(),
		)
		if err != nil {
			return nil, types.ErrWrap(types.ErrKeyNotFound, err)
		}

		return k, nil

	default:
		return nil, types.ErrWrap(types.ErrInvalidMsgType, "invalid keyring backend type")
	}
}

// Provide cosmos sdk tx factory.
func NewFactory(clientCtx cmclient.Context) tx.Factory {
	txFactory := tx.Factory{}.
		WithTxConfig(clientCtx.TxConfig).
		WithKeybase(clientCtx.Keyring).
		WithAccountRetriever(clientCtx.AccountRetriever)

	return txFactory
}

// Make new http client for inquiring several information.
func CtxHttpClient(methodType string, url string, reqBody []byte, ctx context.Context) ([]byte, error) {
	var resp *http.Response
	var err error

	httpClient := &http.Client{Timeout: 30 * time.Second}

	if methodType == "GET" {
		resp, err = ctxhttp.Get(ctx, httpClient, url)
		if err != nil {
			return nil, types.ErrWrap(types.ErrHttpRequest, "failed GET method", err)
		}
	} else if methodType == "POST" {
		resp, err = ctxhttp.Post(ctx, httpClient, url, "application/json", bytes.NewBuffer(reqBody))
		if err != nil {
			return nil, types.ErrWrap(types.ErrHttpRequest, "failed POST method", err)
		}
	} else {
		return nil, types.ErrWrap(types.ErrHttpRequest, "not correct method", err)
	}

	defer resp.Body.Close()

	out, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, types.ErrWrap(types.ErrHttpRequest, "failed to read response", err)
	}

	if resp.StatusCode != 200 {
		return nil, types.ErrWrap(types.ErrHttpRequest, resp.StatusCode, ":", string(out))
	}

	return out, nil
}
