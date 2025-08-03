package config

type (
	Config struct {
		DevMode      bool `yaml:"devMode"      default:"False"`
		DebugLogging bool `yaml:"debugLogging" default:"False"`

		ServerPort int `yaml:"serverPort" default:"4000" validate:"gt=0"`

		Kafka KafkaConfig `yaml:"kafka" validate:"required"`

		Postgres      PostgresConfig      `yaml:"postgres"      validate:"required"`
		Elasticsearch ElasticsearchConfig `yaml:"elasticsearch" validate:"required"`
	}

	KafkaConfig struct {
		SeedBrokerURLs []string `yaml:"seedBrokerURLs" validate:"required,gt=0,dive,notblank"`
	}

	PostgresConfig struct {
		URL string `yaml:"url" validate:"notblank"`
	}

	ElasticsearchConfig struct {
		NodeAddresses []string `yaml:"nodeAddresses" validate:"required,gt=0,dive,notblank"`

		Username string `yaml:"username" validate:"notblank"`
		Password string `yaml:"password" validate:"notblank"`
	}
)
