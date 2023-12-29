package network

import (
	"bufio"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"math/big"
	"math/rand"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/CosmWasm/wasmd/x/wasm"
	"github.com/cosmos/cosmos-sdk/crypto/keys/ed25519"
	"github.com/xpladev/xpla.go/types"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/stretchr/testify/require"
	tmcfg "github.com/tendermint/tendermint/config"
	tmflags "github.com/tendermint/tendermint/libs/cli/flags"
	"github.com/tendermint/tendermint/libs/log"
	"github.com/tendermint/tendermint/node"
	tmclient "github.com/tendermint/tendermint/rpc/client"
	dbm "github.com/tendermint/tm-db"
	"google.golang.org/grpc"

	"github.com/cosmos/cosmos-sdk/baseapp"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/tx"
	"github.com/cosmos/cosmos-sdk/codec"
	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	"github.com/cosmos/cosmos-sdk/crypto/keyring"
	cryptotypes "github.com/cosmos/cosmos-sdk/crypto/types"
	"github.com/cosmos/cosmos-sdk/server"
	"github.com/cosmos/cosmos-sdk/server/api"
	srvconfig "github.com/cosmos/cosmos-sdk/server/config"
	servertypes "github.com/cosmos/cosmos-sdk/server/types"
	"github.com/cosmos/cosmos-sdk/simapp"
	storetypes "github.com/cosmos/cosmos-sdk/store/types"
	"github.com/cosmos/cosmos-sdk/testutil"
	sdk "github.com/cosmos/cosmos-sdk/types"
	simtypes "github.com/cosmos/cosmos-sdk/types/simulation"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	"github.com/cosmos/cosmos-sdk/x/genutil"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	xtestutil "github.com/xpladev/xpla.go/util/testutil"

	"github.com/evmos/ethermint/crypto/hd"
	"github.com/evmos/ethermint/server/config"
	ethermint "github.com/evmos/ethermint/types"
	evmtypes "github.com/evmos/ethermint/x/evm/types"

	"github.com/xpladev/xpla/app"
	"github.com/xpladev/xpla/app/params"
)

const testChainID = "cube_47-5"

// package-wide network lock to only allow one test network at a time
var lock = new(sync.Mutex)

// AppConstructor defines a function which accepts a network configuration and
// creates an ABCI Application to provide to Tendermint.
type AppConstructor = func(val Validator) servertypes.Application

// NewAppConstructor returns a new simapp AppConstructor
func NewAppConstructor(encodingCfg params.EncodingConfig) AppConstructor {
	return func(val Validator) servertypes.Application {
		return app.NewXplaApp(
			val.Ctx.Logger, dbm.NewMemDB(), nil, true, make(map[int64]bool), val.Ctx.Config.RootDir, 0,
			encodingCfg,
			app.GetEnabledProposals(),
			simapp.EmptyAppOptions{},
			[]wasm.Option{},
			baseapp.SetPruning(storetypes.NewPruningOptionsFromString(val.AppConfig.Pruning)),
			baseapp.SetMinGasPrices(val.AppConfig.MinGasPrices),
		)
	}
}

// Config defines the necessary configuration used to bootstrap and start an
// in-process local testing network.
type Config struct {
	KeyringOptions    []keyring.Option // keyring configuration options
	Codec             codec.Codec
	LegacyAmino       *codec.LegacyAmino // TODO: Remove!
	InterfaceRegistry codectypes.InterfaceRegistry
	TxConfig          client.TxConfig
	AccountRetriever  client.AccountRetriever
	AppConstructor    AppConstructor      // the ABCI application constructor
	GenesisState      simapp.GenesisState // custom gensis state to provide
	TimeoutCommit     time.Duration       // the consensus commitment timeout
	AccountTokens     sdk.Int             // the amount of unique validator tokens (e.g. 1000node0)
	StakingTokens     sdk.Int             // the amount of tokens each validator has available to stake
	BondedTokens      sdk.Int             // the amount of tokens each validator stakes
	NumValidators     int                 // the total number of validators to create and bond
	ChainID           string              // the network chain-id
	BondDenom         string              // the staking bond denomination
	MinGasPrices      string              // the minimum gas prices each validator will accept
	PruningStrategy   string              // the pruning strategy each validator will have
	SigningAlgo       string              // signing algorithm for keys
	RPCAddress        string              // RPC listen address (including port)
	JSONRPCAddress    string              // JSON-RPC listen address (including port)
	APIAddress        string              // REST API listen address (including port)
	GRPCAddress       string              // GRPC server listen address (including port)
	EnableTMLogging   bool                // enable Tendermint logging to STDOUT
	CleanupDir        bool                // remove base temporary directory during cleanup
	AdditionalAccount bool
}

// DefaultConfig returns a sane default configuration suitable for nearly all
// testing requirements.
func DefaultConfig() Config {
	encCfg := app.MakeTestEncodingConfig()
	powerReduction := sdk.NewIntFromBigInt(new(big.Int).Exp(big.NewInt(10), big.NewInt(types.BaseDenomUnit), nil))

	return Config{
		Codec:             encCfg.Codec,
		TxConfig:          encCfg.TxConfig,
		LegacyAmino:       encCfg.Amino,
		InterfaceRegistry: encCfg.InterfaceRegistry,
		AccountRetriever:  authtypes.AccountRetriever{},
		AppConstructor:    NewAppConstructor(encCfg),
		GenesisState:      app.ModuleBasics.DefaultGenesis(encCfg.Codec),
		TimeoutCommit:     2 * time.Second,
		ChainID:           testChainID,
		NumValidators:     4,
		BondDenom:         types.XplaDenom,
		MinGasPrices:      fmt.Sprintf("0.000006%s", types.XplaDenom),
		AccountTokens:     sdk.TokensFromConsensusPower(1000, powerReduction),
		StakingTokens:     sdk.TokensFromConsensusPower(500, powerReduction),
		BondedTokens:      sdk.TokensFromConsensusPower(100, powerReduction),
		PruningStrategy:   storetypes.PruningOptionNothing,
		CleanupDir:        true,
		SigningAlgo:       string(hd.EthSecp256k1Type),
		KeyringOptions:    []keyring.Option{hd.EthSecp256k1Option()},
		AdditionalAccount: true,
	}
}

type (
	// Network defines a local in-process testing network using SimApp. It can be
	// configured to start any number of validators, each with its own RPC and API
	// clients. Typically, this test network would be used in client and integration
	// testing where user input is expected.
	//
	// Note, due to Tendermint constraints in regards to RPC functionality, there
	// may only be one test network running at a time. Thus, any caller must be
	// sure to Cleanup after testing is finished in order to allow other tests
	// to create networks. In addition, only the first validator will have a valid
	// RPC and API server/client.
	Network struct {
		T          *testing.T
		BaseDir    string
		Validators []*Validator

		Config Config
	}

	// Validator defines an in-process Tendermint validator node. Through this object,
	// a client can make RPC and API calls and interact with any client command
	// or handler.
	Validator struct {
		AppConfig         *config.Config
		ClientCtx         client.Context
		Ctx               *server.Context
		Dir               string
		NodeID            string
		PubKey            cryptotypes.PubKey
		Moniker           string
		APIAddress        string
		RPCAddress        string
		P2PAddress        string
		Address           sdk.AccAddress
		ValAddress        sdk.ValAddress
		RPCClient         tmclient.Client
		JSONRPCClient     *ethclient.Client
		AdditionalAccount simtypes.Account

		tmNode      *node.Node
		api         *api.Server
		grpc        *grpc.Server
		grpcWeb     *http.Server
		jsonrpc     *http.Server
		jsonrpcDone chan struct{}
	}
)

// New creates a new Network for integration tests or in-process testnets run via the CLI
func New(t *testing.T, cfg Config) Network {
	// only one caller/test can create and use a network at a time
	t.Log("acquiring test network lock")
	lock.Lock()

	baseDir, err := os.MkdirTemp(t.TempDir(), cfg.ChainID)
	require.NoError(t, err)
	t.Logf("created temporary directory: %s", baseDir)

	network := Network{
		T:          t,
		BaseDir:    baseDir,
		Validators: make([]*Validator, cfg.NumValidators),
		Config:     cfg,
	}

	t.Log("preparing test network...")

	monikers := make([]string, cfg.NumValidators)
	nodeIDs := make([]string, cfg.NumValidators)
	valPubKeys := make([]cryptotypes.PubKey, cfg.NumValidators)

	var (
		genAccounts []authtypes.GenesisAccount
		genBalances []banktypes.Balance
		genFiles    []string
	)

	buf := bufio.NewReader(os.Stdin)

	// generate private keys, node IDs, and initial transactions
	for i := 0; i < cfg.NumValidators; i++ {
		appCfg := config.DefaultConfig()
		appCfg.Pruning = cfg.PruningStrategy
		appCfg.MinGasPrices = cfg.MinGasPrices
		appCfg.API.Enable = true
		appCfg.API.Swagger = false
		appCfg.Telemetry.Enabled = false
		appCfg.Telemetry.GlobalLabels = [][]string{{"chain_id", cfg.ChainID}}

		ctx := server.NewDefaultContext()
		tmCfg := ctx.Config
		tmCfg.Consensus.TimeoutCommit = cfg.TimeoutCommit

		// Only allow the first validator to expose an RPC, API and gRPC
		// server/client due to Tendermint in-process constraints.
		apiAddr := ""
		tmCfg.RPC.ListenAddress = ""
		appCfg.GRPC.Enable = false
		appCfg.GRPCWeb.Enable = false
		apiListenAddr := ""
		if i == 0 {
			if cfg.APIAddress != "" {
				apiListenAddr = cfg.APIAddress
			} else {
				var err error
				apiListenAddr, _, err = server.FreeTCPAddr()
				require.NoError(t, err)
			}

			appCfg.API.Address = apiListenAddr
			apiURL, err := url.Parse(apiListenAddr)
			require.NoError(t, err)
			apiAddr = fmt.Sprintf("http://%s:%s", apiURL.Hostname(), apiURL.Port())

			if cfg.RPCAddress != "" {
				tmCfg.RPC.ListenAddress = cfg.RPCAddress
			} else {
				rpcAddr, _, err := server.FreeTCPAddr()
				require.NoError(t, err)
				tmCfg.RPC.ListenAddress = rpcAddr
			}

			if cfg.GRPCAddress != "" {
				appCfg.GRPC.Address = cfg.GRPCAddress
			} else {
				_, grpcPort, err := server.FreeTCPAddr()
				require.NoError(t, err)
				appCfg.GRPC.Address = fmt.Sprintf("0.0.0.0:%s", grpcPort)
			}
			appCfg.GRPC.Enable = true

			_, grpcWebPort, err := server.FreeTCPAddr()
			require.NoError(t, err)
			appCfg.GRPCWeb.Address = fmt.Sprintf("0.0.0.0:%s", grpcWebPort)
			appCfg.GRPCWeb.Enable = true

			if cfg.JSONRPCAddress != "" {
				appCfg.JSONRPC.Address = cfg.JSONRPCAddress
			} else {
				_, jsonRPCPort, err := server.FreeTCPAddr()
				require.NoError(t, err)
				appCfg.JSONRPC.Address = fmt.Sprintf("0.0.0.0:%s", jsonRPCPort)
			}
			appCfg.JSONRPC.Enable = true
			appCfg.JSONRPC.API = config.GetAPINamespaces()
		} else {
			appCfg.JSONRPC.Address = ""
			appCfg.JSONRPC.Enable = true
		}

		logger := log.NewNopLogger()
		if cfg.EnableTMLogging {
			logger = log.NewTMLogger(log.NewSyncWriter(os.Stdout))
			logger, _ = tmflags.ParseLogLevel("info", logger, tmcfg.DefaultLogLevel)
		}

		ctx.Logger = logger

		nodeDirName := fmt.Sprintf("node%d", i)
		nodeDir := filepath.Join(network.BaseDir, nodeDirName, "xplad")
		clientDir := filepath.Join(network.BaseDir, nodeDirName, "xplacli")
		gentxsDir := filepath.Join(network.BaseDir, "gentxs")

		require.NoError(t, os.MkdirAll(filepath.Join(nodeDir, "config"), 0o755))
		require.NoError(t, os.MkdirAll(clientDir, 0o755))

		tmCfg.SetRoot(nodeDir)
		tmCfg.Moniker = nodeDirName
		monikers[i] = nodeDirName

		proxyAddr, _, err := server.FreeTCPAddr()
		require.NoError(t, err)
		tmCfg.ProxyApp = proxyAddr

		p2pAddr, _, err := server.FreeTCPAddr()
		require.NoError(t, err)

		tmCfg.P2P.ListenAddress = p2pAddr
		tmCfg.P2P.AddrBookStrict = false
		tmCfg.P2P.AllowDuplicateIP = true

		nodeID, pubKey, err := genutil.InitializeNodeValidatorFiles(tmCfg)
		require.NoError(t, err)
		nodeIDs[i] = nodeID
		valPubKeys[i] = pubKey

		kb, err := keyring.New(sdk.KeyringServiceName(), keyring.BackendTest, clientDir, buf, cfg.KeyringOptions...)
		require.NoError(t, err)

		keyringAlgos, _ := kb.SupportedAlgorithms()
		algo, err := keyring.NewSigningAlgoFromString(cfg.SigningAlgo, keyringAlgos)
		require.NoError(t, err)

		addr, secret, err := testutil.GenerateSaveCoinKey(kb, nodeDirName, "", true, algo)
		require.NoError(t, err)

		info := map[string]string{"secret": secret}
		infoBz, err := json.Marshal(info)
		require.NoError(t, err)

		// save private key seed words
		require.NoError(t, writeFile(fmt.Sprintf("%v.json", "key_seed"), clientDir, infoBz))

		balances := sdk.NewCoins(
			sdk.NewCoin(fmt.Sprintf("%stoken", nodeDirName), cfg.AccountTokens),
			sdk.NewCoin(cfg.BondDenom, cfg.StakingTokens),
		)

		genFiles = append(genFiles, tmCfg.GenesisFile())
		genBalances = append(genBalances, banktypes.Balance{Address: addr.String(), Coins: balances.Sort()})
		genAccounts = append(genAccounts, &ethermint.EthAccount{
			BaseAccount: authtypes.NewBaseAccount(addr, nil, 0, 0),
			CodeHash:    common.BytesToHash(evmtypes.EmptyCodeHash).Hex(),
		})

		commission, err := sdk.NewDecFromStr("0.5")
		require.NoError(t, err)

		createValMsg, err := stakingtypes.NewMsgCreateValidator(
			sdk.ValAddress(addr),
			valPubKeys[i],
			sdk.NewCoin(cfg.BondDenom, cfg.BondedTokens),
			stakingtypes.NewDescription(nodeDirName, "", "", "", ""),
			stakingtypes.NewCommissionRates(commission, sdk.OneDec(), sdk.OneDec()),
			sdk.OneInt(),
		)
		require.NoError(t, err)

		p2pURL, err := url.Parse(p2pAddr)
		require.NoError(t, err)

		memo := fmt.Sprintf("%s@%s:%s", nodeIDs[i], p2pURL.Hostname(), p2pURL.Port())
		fee := sdk.NewCoins(sdk.NewCoin(cfg.BondDenom, sdk.NewInt(0)))
		txBuilder := cfg.TxConfig.NewTxBuilder()
		err = txBuilder.SetMsgs(createValMsg)
		require.NoError(t, err)

		txBuilder.SetFeeAmount(fee)    // Arbitrary fee
		txBuilder.SetGasLimit(1000000) // Need at least 100386
		txBuilder.SetMemo(memo)

		txFactory := tx.Factory{}
		txFactory = txFactory.
			WithChainID(cfg.ChainID).
			WithMemo(memo).
			WithKeybase(kb).
			WithTxConfig(cfg.TxConfig)

		err = tx.Sign(txFactory, nodeDirName, txBuilder, true)
		require.NoError(t, err)

		txBz, err := cfg.TxConfig.TxJSONEncoder()(txBuilder.GetTx())
		require.NoError(t, err)

		err = writeFile(fmt.Sprintf("%v.json", nodeDirName), gentxsDir, txBz)
		require.NoError(t, err)

		customAppTemplate, _ := config.AppConfig(types.XplaDenom)
		srvconfig.SetConfigTemplate(customAppTemplate)
		srvconfig.WriteConfigFile(filepath.Join(nodeDir, "config/app.toml"), appCfg)

		ctx.Viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_", "-", "_"))
		ctx.Viper.SetConfigFile(filepath.Join(nodeDir, "config/app.toml"))
		err = ctx.Viper.ReadInConfig()
		require.NoError(t, err)

		var accs simtypes.Account
		if cfg.AdditionalAccount {
			privkeySeed := make([]byte, 16)

			src := rand.NewSource(int64(i))
			r := rand.New(src)
			r.Read(privkeySeed)

			mnemonic, err := xtestutil.NewTestMnemonic(privkeySeed)
			require.NoError(t, err)

			keyringAlgos, _ := kb.SupportedAlgorithms()
			algo, err := keyring.NewSigningAlgoFromString(cfg.SigningAlgo, keyringAlgos)
			require.NoError(t, err)

			addr, secret, err := testutil.GenerateSaveCoinKey(kb, "add_account_"+nodeDirName, mnemonic, false, algo)
			require.NoError(t, err)

			info := map[string]string{"secret": secret}
			infoBz, err := json.Marshal(info)
			require.NoError(t, err)
			require.NoError(t, writeFile(fmt.Sprintf("%v.json", "add_key_seed"), clientDir, infoBz))

			balances := sdk.NewCoins(
				sdk.NewCoin(cfg.BondDenom, cfg.StakingTokens),
			)
			genBalances = append(genBalances, banktypes.Balance{Address: addr.String(), Coins: balances.Sort()})
			genAccounts = append(genAccounts, &ethermint.EthAccount{
				BaseAccount: authtypes.NewBaseAccount(addr, nil, 0, 0),
				CodeHash:    common.BytesToHash(evmtypes.EmptyCodeHash).Hex(),
			})

			ethsecpPrivKey, err := xtestutil.NewTestEthSecpPrivKey(secret)
			require.NoError(t, err)

			accs.PrivKey = ethsecpPrivKey
			accs.PubKey = accs.PrivKey.PubKey()
			accs.Address = sdk.AccAddress(accs.PubKey.Address())
			accs.ConsKey = ed25519.GenPrivKeyFromSecret(privkeySeed)
		}

		clientCtx := client.Context{}.
			WithKeyringDir(clientDir).
			WithKeyring(kb).
			WithHomeDir(tmCfg.RootDir).
			WithChainID(cfg.ChainID).
			WithInterfaceRegistry(cfg.InterfaceRegistry).
			WithCodec(cfg.Codec).
			WithLegacyAmino(cfg.LegacyAmino).
			WithTxConfig(cfg.TxConfig).
			WithAccountRetriever(cfg.AccountRetriever)

		network.Validators[i] = &Validator{
			AppConfig:         appCfg,
			ClientCtx:         clientCtx,
			Ctx:               ctx,
			Dir:               filepath.Join(network.BaseDir, nodeDirName),
			NodeID:            nodeID,
			PubKey:            pubKey,
			Moniker:           nodeDirName,
			RPCAddress:        tmCfg.RPC.ListenAddress,
			P2PAddress:        tmCfg.P2P.ListenAddress,
			APIAddress:        apiAddr,
			Address:           addr,
			ValAddress:        sdk.ValAddress(addr),
			AdditionalAccount: accs,
		}
	}

	err = initGenFiles(cfg, genAccounts, genBalances, genFiles)
	require.NoError(t, err)

	err = collectGenFiles(cfg, network.Validators, network.BaseDir)
	require.NoError(t, err)

	t.Log("starting test network...")
	for _, v := range network.Validators {
		err := startInProcess(cfg, v)
		require.NoError(t, err)
	}

	t.Log("started test network")

	// Ensure we cleanup incase any test was abruptly halted (e.g. SIGINT) as any
	// defer in a test would not be called.
	server.TrapSignal(network.Cleanup)

	return network
}

// LatestHeight returns the latest height of the network or an error if the
// query fails or no validators exist.
func (n *Network) LatestHeight() (int64, error) {
	if len(n.Validators) == 0 {
		return 0, errors.New("no validators available")
	}

	status, err := n.Validators[0].RPCClient.Status(context.Background())
	if err != nil {
		return 0, err
	}

	return status.SyncInfo.LatestBlockHeight, nil
}

// WaitForHeight performs a blocking check where it waits for a block to be
// committed after a given block. If that height is not reached within a timeout,
// an error is returned. Regardless, the latest height queried is returned.
func (n *Network) WaitForHeight(h int64) (int64, error) {
	return n.WaitForHeightWithTimeout(h, 10*time.Second)
}

// WaitForHeightWithTimeout is the same as WaitForHeight except the caller can
// provide a custom timeout.
func (n *Network) WaitForHeightWithTimeout(h int64, t time.Duration) (int64, error) {
	ticker := time.NewTicker(time.Second)
	timeout := time.After(t)

	if len(n.Validators) == 0 {
		return 0, errors.New("no validators available")
	}

	var latestHeight int64
	val := n.Validators[0]

	for {
		select {
		case <-timeout:
			ticker.Stop()
			return latestHeight, errors.New("timeout exceeded waiting for block")
		case <-ticker.C:
			status, err := val.RPCClient.Status(context.Background())
			if err == nil && status != nil {
				latestHeight = status.SyncInfo.LatestBlockHeight
				if latestHeight >= h {
					return latestHeight, nil
				}
			}
		}
	}
}

// WaitForNextBlock waits for the next block to be committed, returning an error
// upon failure.
func (n *Network) WaitForNextBlock() error {
	lastBlock, err := n.LatestHeight()
	if err != nil {
		return err
	}

	_, err = n.WaitForHeight(lastBlock + 1)
	if err != nil {
		return err
	}

	return err
}

// Cleanup removes the root testing (temporary) directory and stops both the
// Tendermint and API services. It allows other callers to create and start
// test networks. This method must be called when a test is finished, typically
// in a defer.
func (n *Network) Cleanup() {
	defer func() {
		lock.Unlock()
		n.T.Log("released test network lock")
	}()

	n.T.Log("cleaning up test network...")

	for _, v := range n.Validators {
		if v.tmNode != nil && v.tmNode.IsRunning() {
			_ = v.tmNode.Stop()
		}

		if v.api != nil {
			_ = v.api.Close()
		}

		if v.grpc != nil {
			v.grpc.Stop()
			if v.grpcWeb != nil {
				_ = v.grpcWeb.Close()
			}
		}

		if v.jsonrpc != nil {
			shutdownCtx, cancelFn := context.WithTimeout(context.Background(), 10*time.Second)
			defer cancelFn()

			if err := v.jsonrpc.Shutdown(shutdownCtx); err != nil {
				v.tmNode.Logger.Error("HTTP server shutdown produced a warning", "error", err.Error())
			} else {
				v.tmNode.Logger.Info("HTTP server shut down, waiting 5 sec")
				select {
				case <-time.Tick(2 * time.Second):
				case <-v.jsonrpcDone:
				}
			}
		}
	}

	if n.Config.CleanupDir {
		_ = os.RemoveAll(n.BaseDir)
	}

	n.T.Log("finished cleaning up test network")
}
