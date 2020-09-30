package consul

import (
	"context"
	"fmt"
	"net"

	"github.com/douyu/jupiter/pkg/constant"
	"github.com/douyu/jupiter/pkg/registry"
	"github.com/douyu/jupiter/pkg/server"
	uuid "github.com/google/uuid"
	"github.com/hashicorp/consul/api"
)

type consulRegistry struct {
	*Config
	client *api.Client
}

func newConsulRegistry(config *Config) *consulRegistry {
	conf := api.DefaultConfig()
	conf.Address = config.Endpoints[0]
	c, err := api.NewClient(conf)
	if err != nil {
		config.logger.Panic(err.Error())
	}
	return &consulRegistry{client: c, Config: config}
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
	uuid := uuid.New().String()
	registration.ID = info.Name + "." + uuid

	check := new(api.AgentServiceCheck)
	// 健康检查配置 TTL
	// check.TTL = "10s"
	// check.Notes = "Web app does a curl internally every 10 seconds"
	// registration.Check = check

	// 健康检查配置 HTTP
	check.HTTP = fmt.Sprintf("http://%s/health", info.Address)
	check.Timeout = "2s"
	check.Interval = "5s"
	check.DeregisterCriticalServiceAfter = "15s"
	registration.Check = check

	return reg.client.Agent().ServiceRegister(registration)
}

func (reg *consulRegistry) UnregisterService(ctx context.Context, info *server.ServiceInfo) error {
	fmt.Println("UnregisterService---------------", ctx.Value("serviceid"))
	return reg.client.Agent().ServiceDeregister(info.AppID)
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
	return nil
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
