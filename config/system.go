package config

type System struct {
	Port          int    `mapstructure:"port" json:"port" yaml:"port"`
	LogPath       string `mapstructure:"log-path" json:"log-path" yaml:"log-path"`
}
