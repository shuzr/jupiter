package consul

import "github.com/hashicorp/consul/api"

// Client ...
type Client struct {
	*api.Client
	config *Config
}

func newClient(config *Config) *Client {
	conf := api.DefaultConfig()
	conf.Address = config.Endpoints[0]

	client, err := api.NewClient(conf)
	if err != nil {
		config.logger.Panic(err.Error())
	}
	return &Client{Client: client, config: config}
}
