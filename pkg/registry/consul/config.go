package consul

import (
	"time"

	"github.com/douyu/jupiter/pkg/client/consul"
	"github.com/douyu/jupiter/pkg/conf"
	"github.com/douyu/jupiter/pkg/ecode"
	"github.com/douyu/jupiter/pkg/registry"

	"github.com/douyu/jupiter/pkg/xlog"
)

// StdConfig ...
func StdConfig(name string) *Config {
	return RawConfig("jupiter.registry." + name)
}

// RawConfig ...
func RawConfig(key string) *Config {
	var config = DefaultConfig()
	// 解析最外层配置
	if err := conf.UnmarshalKey(key, &config); err != nil {
		xlog.Panic("unmarshal key", xlog.FieldMod("registry.consul"), xlog.FieldErrKind(ecode.ErrKindUnmarshalConfigErr), xlog.FieldErr(err), xlog.String("key", key), xlog.Any("config", config))
	}
	// 解析嵌套配置
	if err := conf.UnmarshalKey(key, &config.Config); err != nil {
		xlog.Panic("unmarshal key", xlog.FieldMod("registry.consul"), xlog.FieldErrKind(ecode.ErrKindUnmarshalConfigErr), xlog.FieldErr(err), xlog.String("key", key), xlog.Any("config", config))
	}
	return config
}

// DefaultConfig ...
func DefaultConfig() *Config {
	return &Config{
		Config:      consul.DefaultConfig(),
		ReadTimeout: time.Second * 3,
		Prefix:      "jupiter",
		logger:      xlog.JupiterLogger,
	}
}

// Config ...
type Config struct {
	*consul.Config
	ReadTimeout time.Duration
	ConfigKey   string
	Prefix      string
	logger      *xlog.Logger
}

// Build ...
func (config Config) Build() registry.Registry {
	if config.ConfigKey != "" {
		config.Config = consul.RawConfig(config.ConfigKey)
	}
	return newConsulRegistry(&config)
}
