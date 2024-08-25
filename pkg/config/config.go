// Copyright 2024 Matrix Origin.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package config

import (
	"fmt"
	"os"
	"time"

	"github.com/go-logr/logr"
	"github.com/spf13/viper"
)

var gCfg *Configuration

func InitConfiguration(l logr.Logger, cfgPath string) (*Configuration, error) {
	var cfg = NewConfiguration()
	viper.SetConfigFile(cfgPath)
	if err := viper.ReadInConfig(); err != nil {
		l.Error(err, "failed to read config file")
		return nil, err
	}
	for _, key := range viper.AllKeys() {
		if value := os.Getenv(key); value != "" {
			viper.Set(key, value)
		}
	}
	if err := viper.Unmarshal(cfg); err != nil {
		l.Error(err, "failed to unmarshal config file")
		return nil, err
	}
	setConfiguration(cfg)
	return cfg, nil
}

type Configuration struct {
	App AppConfig `yaml:"app"`
}

func GetConfiguration() *Configuration {
	if gCfg == nil {
		gCfg = NewConfiguration()
	}
	return gCfg
}

// setConfiguration should only do ONCE in main.go
func setConfiguration(cfg *Configuration) {
	gCfg = cfg
}

func (c *Configuration) Validate() error {
	if err := c.App.Validate(); err != nil {
		return err
	}
	return nil
}

func NewConfiguration() *Configuration {
	return &Configuration{
		App: *NewAppConfig(),
	}
}

type AppConfig struct {
	RpcPort    int           `yaml:"rpcPort"`
	RpcTimeout time.Duration `yaml:"rpcTimeout"`
}

func NewAppConfig() *AppConfig {
	return &AppConfig{
		RpcPort:    18002,
		RpcTimeout: 3 * time.Second,
	}
}

func (c *AppConfig) Validate() error {
	var err error
	return err
}

func (c *AppConfig) GetRpcAddr(instance string) string {
	if instance == "" {
		instance = "127.0.0.1"
	}
	return fmt.Sprintf("%s:%d", instance, c.RpcPort)
}
