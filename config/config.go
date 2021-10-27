package config

type Server struct {
	System System `mapstructure:"system" json:"system" yaml:"system"`
	Mysql  Mysql  `mapstructure:"mysql" json:"mysql" yaml:"mysql"`
	Kafka  Kafka  `mapstructue:"kafka" json:"kafka" yaml:"kafka"`
	Dana    Dana `mapstructure:"dana" json:"dana" yaml:"dana"`
	PrivacyEncrypt EncryptPair  `mapstructure:"privacy-encrypt" json:"privacy-encrypt" yaml:"privacy-encrypt"`
	IgnoreSign bool  `mapstructure:"ignore-sign" json:"ignore-sign" yaml:"ignore-sign"`
	DataCenterId int64   `mapstructure:"data-center-id" json:"data-center-id" yaml:"data-center-id"`
}
