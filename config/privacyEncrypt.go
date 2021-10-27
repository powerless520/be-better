package config

type EncryptPair struct{
	Iv        string `mapstructure:"iv" json:"iv" yaml:"iv"`
	Key       string `mapstructure:"key" json:"key" yaml:"key"`
}