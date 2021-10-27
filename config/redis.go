package config

type Redis struct {
	DB       int    `mapstructure:"dbutil" json:"dbutil" yaml:"dbutil"`
	Addr     string `mapstructure:"addr" json:"addr" yaml:"addr"`
	Password string `mapstructure:"password" json:"password" yaml:"password"`
}