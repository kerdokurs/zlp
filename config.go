package zlp

import "net/http"

type APIVersion = string

const (
	ApiV1 APIVersion = "v1"
)

type Config struct {
	ZulipRC
	ApiVersion string
    Client Doer
    UserAgent string
}

type ConfigFunction func(*Config)

func defaultConfig() Config {
	return Config{
		ApiVersion: ApiV1,
        Client: http.DefaultClient,
	}
}
