package setup

import (
	"go.uber.org/zap/zapcore"
	"sigs.k8s.io/controller-runtime/pkg/log/zap"
)

func CallerName() zap.EncoderConfigOption {
	return func(cfg *zapcore.EncoderConfig) {
		cfg.CallerKey = "caller"
	}
}
