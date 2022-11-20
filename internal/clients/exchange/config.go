package exchange

type Config struct {
	APIKey     string `envconfig:"APP_LAYER_KEY" default:""`
	APIBaseURL string `envconfig:"APP_API_BASE_URL" default:"https://api.exchangerate.host"`
}
