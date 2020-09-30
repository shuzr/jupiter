package consul

import (
	"context"
	"testing"
	"time"

	"github.com/douyu/jupiter/pkg/server"
	"github.com/douyu/jupiter/pkg/xlog"
)

func TestRunConsulRegistry(t *testing.T) {
	conf := DefaultConfig()
	conf.Endpoints = []string{"192.168.88.206:8500"}
	registry := newConsulRegistry(&Config{
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

// type Engine struct {
// 	jupiter.Application
// }

// // NewEngine ...
// func NewEngine() *Engine {
// 	eng := &Engine{}
// 	if err := eng.Startup(
// 		eng.serveHTTP,
// 		// eng.setLogger,
// 	); err != nil {
// 		xlog.Panic("startup", xlog.Any("err", err))
// 	}
// 	return eng
// }

// // HTTP地址
// func (eng *Engine) serveHTTP() error {
// 	server := xecho.StdConfig("http").Build()
// 	server.GET("/hello", func(ctx echo.Context) error {
// 		return ctx.JSON(200, "Gopher Wuhan")
// 	})
// 	return eng.Serve(server)
// }

// // 治理地址
// func (eng *Engine) setLogger() error {
// 	xlog.DefaultLogger = xlog.StdConfig("default").Build()
// 	return nil
// }

// func runHttpSrvc() {
// 	eng := NewEngine()
// 	// 注册中心
// 	eng.SetRegistry(compound.New(
// 		StdConfig("test").Build(),
// 	),
// 	)
// 	if err := eng.Run(); err != nil {
// 		xlog.Panic(err.Error())
// 	}
// }

// func TestRunHttpSrvc(t *testing.T) {
// 	runHttpSrvc()
// }
