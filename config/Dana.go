package config

type Dana struct {
	Facm          DanaTopic    `mapstructure:"facm" json:"facm" yaml:"facm"`
	FacmEvent     DanaTopic    `mapstructure:"facm-event" json:"facm-event" yaml:"facm-event"`
}

type DanaTopic struct{
	Url         string `mapstructure:"url" json:"url" yaml:"url"`
	Topic       string `mapstructure:"topic" json:"topic" yaml:"topic"`
	Token       string `mapstructure:"token" json:"token" yaml:"token"`
}