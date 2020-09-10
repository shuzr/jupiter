package consul

import (
	"testing"

	"github.com/douyu/jupiter"
	"github.com/douyu/jupiter/pkg/registry/compound"
	"github.com/douyu/jupiter/pkg/server/xecho"
	"github.com/douyu/jupiter/pkg/xlog"
	"github.com/labstack/echo/v4"
)

func TestRunConsulRegistry(t *testing.T) {
	// RunConsulRegistry()
}

type Engine struct {
	jupiter.Application
}

// NewEngine ...
func NewEngine() *Engine {
	eng := &Engine{}
	if err := eng.Startup(
		eng.serveHTTP,
		// eng.setLogger,
	); err != nil {
		xlog.Panic("startup", xlog.Any("err", err))
	}
	return eng
}

// HTTP地址
func (eng *Engine) serveHTTP() error {
	server := xecho.StdConfig("http").Build()
	server.GET("/hello", func(ctx echo.Context) error {
		return ctx.JSON(200, "Gopher Wuhan")
	})
	return eng.Serve(server)
}

// 治理地址
func (eng *Engine) setLogger() error {
	xlog.DefaultLogger = xlog.StdConfig("default").Build()
	return nil
}

func runHttpSrvc() {
	eng := NewEngine()
	// 注册中心
	eng.SetRegistry(compound.New(
		StdConfig("test").Build(),
	),
	)
	if err := eng.Run(); err != nil {
		xlog.Panic(err.Error())
	}
}

func TestRunHttpSrvc(t *testing.T) {
	runHttpSrvc()
}
