package config

type Kafka struct {
	Username  string    `mapstructure:"username" json:"username" yaml:"username"`
	Password  string    `mapstructure:"password" json:"password" yaml:"password"`
	Url       []string  `mapstructure:"url" json:"Url" yaml:"Url"`
	Consumers Consumers `mapstructure:"consumers" json:"consumers" yaml:"consumers"`
}

type Consumers struct {
	Facm      ConsumerTopic `mapstructure:"facm" json:"facm" yaml:"facm"`
	FacmEvent ConsumerTopic `mapstructure:"facm-event" json:"face-event" yaml:"facm-event"`
}

type ConsumerTopic struct {
	Topic   string `mapstructure:"topic" json:"topic" yaml:"topic"`
	GroupId string `mapstructure:"group-id" json:"group-id" yaml:"group-id"`
}
