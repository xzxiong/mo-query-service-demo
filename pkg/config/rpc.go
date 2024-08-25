package config

import (
	"github.com/matrixorigin/matrixone/pkg/common/morpc"
)

func GetRpcConfig() *morpc.Config {
	var cfg morpc.Config
	cfg.Adjust()
	return &cfg
}
