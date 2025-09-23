package livekit

type Config struct {
	Host         string
	ApiKey       string
	ApiSecret    string
	EgressKey    string
	EgressSecret string
}

func NewConfig(
	host, apikey, apisecret string,
	egresskey string,
	egresssecret string,
) *Config {
	return &Config{
		Host:         host,
		ApiKey:       apikey,
		ApiSecret:    apisecret,
		EgressKey:    egresskey,
		EgressSecret: egresssecret,
	}
}
