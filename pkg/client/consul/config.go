package consul

import (
	"time"

	"github.com/douyu/jupiter/pkg/conf"
	"github.com/douyu/jupiter/pkg/ecode"
	"github.com/douyu/jupiter/pkg/util/xtime"
	"github.com/douyu/jupiter/pkg/xlog"
)

// Config ...
type Config struct {
	Endpoints      []string      `json:"endpoints"`
	ConnectTimeout time.Duration `json:"connectTimeout"`
	logger         *xlog.Logger
}

// DefaultConfig 返回默认配置
func DefaultConfig() *Config {
	return &Config{
		ConnectTimeout: xtime.Duration("5s"),
		logger:         xlog.JupiterLogger.With(xlog.FieldMod("client.consul")),
	}
}

// StdConfig ...
func StdConfig(name string) *Config {
	return RawConfig("jupiter.consul." + name)
}

// RawConfig ...
func RawConfig(key string) *Config {
	var config = DefaultConfig()
	if err := conf.UnmarshalKey(key, config); err != nil {
		config.logger.Panic("client etcd parse config panic", xlog.FieldErrKind(ecode.ErrKindUnmarshalConfigErr), xlog.FieldErr(err), xlog.FieldKey(key), xlog.FieldValueAny(config))
	}
	return config
}

// WithLogger ...
func (config *Config) WithLogger(logger *xlog.Logger) *Config {
	config.logger = logger
	return config
}

// Build ...
func (config *Config) Build() *Client {
	cc := newClient(config)
	return cc
}
