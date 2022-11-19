package exchange

type Config struct {
	APIKey string `envconfig:"APP_LAYER_KEY" default:""`
}
