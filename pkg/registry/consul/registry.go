package consul

import (
	"context"
	"fmt"
	"net"

	"github.com/douyu/jupiter/pkg/client/consul"
	"github.com/douyu/jupiter/pkg/constant"
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
}

func newConsulRegistry(config *Config) *consulRegistry {
	if config.logger == nil {
		config.logger = xlog.JupiterLogger
	}
	config.logger = config.logger.With(xlog.FieldMod( /*ecode.ModRegistryConsul*/ "registry.consul"), xlog.FieldAddrAny(config.Config.Endpoints))
	return &consulRegistry{client: config.Config.Build(), Config: config}
}

func (reg *consulRegistry) RegisterService(ctx context.Context, info *server.ServiceInfo) error {

	addr, err := net.ResolveTCPAddr("", info.Address)
	if err != nil {
		return err
	}
	// 服务注册配置
	registration := new(api.AgentServiceRegistration)
	registration.Name = info.Name
	registration.Address = addr.IP.String()
	registration.Port = addr.Port
	uuid := uuid.NewV4()
	reg.id = info.Name + "." + uuid.String()
	registration.ID = reg.id

	// 健康检查配置
	check := new(api.AgentServiceCheck)
	check.HTTP = fmt.Sprintf("http://%s", info.Address)
	check.Timeout = "3s"
	check.Interval = "5s"
	check.DeregisterCriticalServiceAfter = "30s"
	registration.Check = check

	return reg.client.Agent().ServiceRegister(registration)
}

func (reg *consulRegistry) UnregisterService(ctx context.Context, info *server.ServiceInfo) error {
	return reg.client.Agent().ServiceDeregister(reg.id)
}

func (reg *consulRegistry) ListServices(ctx context.Context, name string, scheme string) (services []*server.ServiceInfo, err error) {
	srvs, err := reg.client.Agent().Services()
	if err != nil {
		return
	}
	for _, v := range srvs {
		if v.Service == name {
			services = append(services, makeService(v))
		}
	}
	return
}

func (reg *consulRegistry) WatchServices(ctx context.Context, name string, scheme string) (chan registry.Endpoints, error) {
	return nil, nil
}

func (reg *consulRegistry) Close() error {
	return reg.client.Agent().ServiceDeregister(reg.id)
}

func makeService(as *api.AgentService) *server.ServiceInfo {
	return &server.ServiceInfo{
		Name:       as.Service,
		AppID:      "",
		Scheme:     "http",
		Address:    fmt.Sprintf("%s:%d", as.Address, as.Port),
		Weight:     0,
		Enable:     true,
		Healthy:    true,
		Metadata:   as.Meta,
		Region:     "default",
		Zone:       "default",
		Kind:       constant.ServiceProvider,
		Deployment: "default",
		Group:      "",
	}
}
