package consul

import (
	"context"
	"fmt"
	"net"
	"time"

	"github.com/douyu/jupiter/pkg/client/consul"
	"github.com/douyu/jupiter/pkg/ecode"
	"github.com/douyu/jupiter/pkg/registry"
	"github.com/douyu/jupiter/pkg/server"
	"github.com/douyu/jupiter/pkg/xlog"
	"github.com/hashicorp/consul/api"
	uuid "github.com/satori/go.uuid"
)

type consulRegistry struct {
	id string
	*Config
	client *consul.Client
	ctx    context.Context
}

func newConsulRegistry(config *Config) *consulRegistry {
	if config.logger == nil {
		config.logger = xlog.JupiterLogger
	}
	config.logger = config.logger.With(xlog.FieldMod(ecode.ModRegistryConsul), xlog.FieldAddrAny(config.Config.Endpoints))
	return &consulRegistry{client: config.Config.Build(), Config: config}
}

func (reg *consulRegistry) RegisterService(ctx context.Context, info *server.ServiceInfo) error {

	addr, err := net.ResolveTCPAddr("", info.Address)
	if err != nil {
		return err
	}
	reg.ctx = ctx
	// 服务注册配置
	registration := new(api.AgentServiceRegistration)
	registration.Name = info.Name
	registration.Address = addr.IP.String()
	registration.Port = addr.Port
	uuid, _ := uuid.NewV4()
	reg.id = info.Name + "." + uuid.String()
	registration.ID = reg.id

	// 健康检查配置
	check := new(api.AgentServiceCheck)
	check.HTTP = fmt.Sprintf("http://%s", info.Address)
	check.Timeout = "5s"
	check.Interval = "5s"
	check.DeregisterCriticalServiceAfter = "30s"
	registration.Check = check

	return reg.client.Agent().ServiceRegister(registration)
}

func (reg *consulRegistry) UnregisterService(ctx context.Context, info *server.ServiceInfo) error {
	return reg.client.Agent().ServiceDeregister(info.AppID)
}

func (reg *consulRegistry) ListServices(ctx context.Context, name string, id string) ([]*server.ServiceInfo, error) {
	return nil, nil
}

func (reg *consulRegistry) WatchServices(ctx context.Context, name string, id string) (chan registry.Endpoints, error) {
	return nil, nil
}

func (reg *consulRegistry) Close() error {
	return nil
}

// RunConsulRegistry ...
func RunConsulRegistry() {
	conf := consul.DefaultConfig()
	conf.Endpoints = []string{"192.168.88.206:8500"}
	registry := newConsulRegistry(&Config{
		Config:      conf,
		ReadTimeout: time.Second * 10,
		Prefix:      "jupiter",
		logger:      xlog.DefaultLogger,
	})

	registry.RegisterService(context.Background(), &server.ServiceInfo{
		Name:     "go.service.consul1.test",
		AppID:    "",
		Scheme:   "http",
		Address:  "192.168.88.57:37219",
		Weight:   0,
		Enable:   true,
		Healthy:  true,
		Metadata: map[string]string{},
		Region:   "default",
		Zone:     "default",
		// Kind:       constant.ServiceProvider,
		Deployment: "default",
		Group:      "",
	})
}
