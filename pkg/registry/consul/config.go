package consul

import (
	"time"

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
	// if err := conf.UnmarshalKey(key, &config); err != nil {
	// 	xlog.Panic("unmarshal key", xlog.FieldMod("registry.consul"), xlog.FieldErrKind(ecode.ErrKindUnmarshalConfigErr), xlog.FieldErr(err), xlog.String("key", key), xlog.Any("config", config))
	// }
	return config
}

// DefaultConfig ...
func DefaultConfig() *Config {
	return &Config{
		ReadTimeout: time.Second * 3,
		Prefix:      "jupiter",
		logger:      xlog.JupiterLogger,
	}
}

// Config ...
type Config struct {
	Endpoints   []string `json:"endpoints"`
	ReadTimeout time.Duration
	ConfigKey   string
	Prefix      string
	logger      *xlog.Logger
}

// Build ...
func (config *Config) Build() registry.Registry {
	if config.logger == nil {
		config.logger = xlog.JupiterLogger
	}
	config.logger = config.logger.With(xlog.FieldMod( /*ecode.ModRegistryConsul*/ "registry.consul"), xlog.FieldAddrAny(config.Endpoints))
	return newConsulRegistry(config)
}
