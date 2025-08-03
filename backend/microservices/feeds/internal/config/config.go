package config

type (
	Config struct {
		DevMode      bool `yaml:"devMode"      default:"False"`
		DebugLogging bool `yaml:"debugLogging" default:"False"`

		ServerPort int `yaml:"serverPort" default:"4000" validate:"gt=0"`

		Redis RedisConfig `yaml:"redis" validate:"required"`
	}

	RedisConfig struct {
		NodeAddresses []string `yaml:"nodeAddresses" validate:"required,gt=0,dive,notblank"`

		Username string `yaml:"username" validate:"notblank"`
		Password string `yaml:"password" validate:"notblank"`
	}
)
