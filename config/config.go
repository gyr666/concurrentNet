package config

type ConfigStrategy interface {
	Fill( *Config) error
}