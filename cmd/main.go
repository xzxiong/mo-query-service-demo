package main

import (
	"context"
	"flag"
	"os"

	"github.com/google/uuid"
	"github.com/matrixorigin/matrixone/pkg/common/runtime"
	"github.com/matrixorigin/matrixone/pkg/defines"
	"github.com/matrixorigin/matrixone/pkg/pb/query"
	qclient "github.com/matrixorigin/matrixone/pkg/queryservice/client"
	rawzap "go.uber.org/zap"
	"sigs.k8s.io/controller-runtime/pkg/log/zap"

	"github.com/xzxiong/mo-query-service-demo/cmd/setup"
	"github.com/xzxiong/mo-query-service-demo/pkg/config"
)

func main() {
	var ctx, cancel = context.WithCancel(context.Background())
	var cfgPath string
	var logCaller bool
	flag.StringVar(&cfgPath, "cfg", "", "The config filepath")
	flag.BoolVar(&logCaller, "log-caller", true, "Enable log caller.")
	opts := zap.Options{
		Development:          true,
		EncoderConfigOptions: []zap.EncoderConfigOption{setup.CallerName()},
		ZapOpts:              []rawzap.Option{},
	}
	opts.BindFlags(flag.CommandLine)
	flag.Parse()
	// END> cmd-line parsed.
	defer cancel()

	if logCaller {
		opts.ZapOpts = append(opts.ZapOpts, rawzap.WithCaller(logCaller))
	}
	logger := zap.New(zap.UseFlagOptions(&opts))
	config.SetLogger(logger)

	// read config
	if cfgPath != "" {
		if _, err := config.InitConfiguration(logger, cfgPath); err != nil {
			os.Exit(1)
		}
	}
	cfg := config.GetConfiguration()
	if err := cfg.Validate(); err != nil {
		logger.Error(err, "failed to load configuration")
		os.Exit(1)
	}

	uid := uuid.Must(uuid.NewRandom())
	rpcCfg := config.GetRpcConfig()
	logger.Info("init QueryService client")
	client, err := qclient.NewQueryClient(uid.String(), *rpcCfg)
	if err != nil {
		logger.Error(err, "failed to init QueryService client")
		os.Exit(1)
	}

	// init for rpc env
	rt := runtime.DefaultRuntimeWithLevel(rawzap.InfoLevel)
	rt.SetGlobalVariables(runtime.MOProtocolVersion, defines.MORPCLatestVersion)
	runtime.SetupProcessLevelRuntime(rt)

	addr := cfg.App.GetRpcAddr("")
	req := client.NewRequest(query.CmdMethod_GetProtocolVersion)
	req.GetProtocolVersion = &query.GetProtocolVersionRequest{}
	deadlineCtx, dcCancel := context.WithTimeout(ctx, cfg.App.RpcTimeout)
	defer dcCancel()
	resp, err := client.SendMessage(deadlineCtx, addr, req)
	if err != nil {
		logger.Error(err, "failed to request QueryService")
		os.Exit(1)
	}
	logger.Info("GetProtocolVersion", "addr", addr, "version", resp.GetProtocolVersion.Version)
	client.Release(resp)
}
