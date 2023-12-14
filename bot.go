package zlp

import (
	"net/http"
)

type Bot struct {
	Config
	Client Doer
}

func WithRCFile(fileName string) ConfigFunction {
	return func(c *Config) {
		rc, err := loadRC(fileName)
		if err != nil {
			panic(err)
		}

		c.ZulipRC = *rc
	}
}

func WithRCEnv() ConfigFunction {
	return func(c *Config) {
		rc, _ := loadRCFromEnv()
		c.ZulipRC = *rc
	}
}

func WithAPIVersion(version APIVersion) ConfigFunction {
	return func(c *Config) {
		c.ApiVersion = version
	}
}

func WithHTTPClient(client Doer) ConfigFunction {
    return func(c *Config) {
        c.Client = client
    }
}

func NewBot(cfgs ...ConfigFunction) *Bot {
	cfg := defaultConfig()
	for _, fn := range cfgs {
		fn(&cfg)
	}

	return &Bot{
		Config: cfg,
        Client: cfg.Client,
	}
}

type Doer interface {
	Do(*http.Request) (*http.Response, error)
}
