package http

type Config struct {
	Port int `envconfig:"APP_HTTP_PORT" default:"8888"`
}
