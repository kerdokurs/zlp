package zlp

type APIVersion = string

const (
	ApiV1 APIVersion = "v1"
)

type Config struct {
	ZulipRC
	ApiVersion string
}

type ConfigFunction func(*Config)

func defaultConfig() Config {
	return Config{
		ApiVersion: ApiV1,
	}
}
