package server

type Config struct {
	BindAddr string
	Loglevel string
}

func NewCofig() *Config {
	return &Config{
		BindAddr: ":8080",
		Loglevel: "Debug",
	}
}
