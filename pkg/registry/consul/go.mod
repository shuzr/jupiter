module consul

go 1.15

require (
	github.com/douyu/jupiter v0.2.4
	github.com/douyu/jupiter/pkg/client/consul v0.0.0-00010101000000-000000000000
	// github.com/douyu/jupiter/pkg/client/consul v0.0.0-00010101000000-000000000000
	github.com/hashicorp/consul/api v1.6.0
	github.com/labstack/echo/v4 v4.1.16
	github.com/satori/go.uuid v1.2.0
)

replace github.com/douyu/jupiter/pkg/ecode => E:/workspace/Code/Go/works/jupiter/pkg/ecode

replace github.com/douyu/jupiter/pkg/client/consul => E:/workspace/Code/Go/works/jupiter/pkg/client/consul
