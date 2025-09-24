package livekit

type Config struct {
	Host         string
	ApiKey       string
	ApiSecret    string
	EgressKey    string
	EgressSecret string
	WebhookKey   string
}

func NewConfig(
	host, apikey, apisecret string,
	egresskey string,
	egresssecret string,
	webhook string,
) *Config {
	return &Config{
		Host:         host,
		ApiKey:       apikey,
		ApiSecret:    apisecret,
		EgressKey:    egresskey,
		EgressSecret: egresssecret,
		WebhookKey:   webhook,
	}
}
