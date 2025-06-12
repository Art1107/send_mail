package config

type Config struct {
	LarkWebhookURL string
	ServerPort     string
}

func NewConfig() *Config {
	return &Config{
		LarkWebhookURL: "https://open.larksuite.com/open-apis/bot/v2/hook/662a59e9-d3b5-41db-81c0-911e083525d5",
		ServerPort:     ":8080",
	}
}
