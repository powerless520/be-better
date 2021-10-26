package config

type Server struct {
	System System `mapstructure:"system" json:"system" yaml:"system"`
	Mysql  Mysql  `mapstructure:"mysql" json:"mysql" yaml:"mysql"`
	Kafka  Kafka  `mapstructue:"kafka" json:"kafka" yaml:"kafka"`
}
