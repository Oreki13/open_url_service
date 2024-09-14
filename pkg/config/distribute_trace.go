package config

type DistributeTraceConfig struct {
	TempoConfig `mapstructure:",squash"`
}

type TempoConfig struct {
	TempoHost string `mapstructure:"tempo_host"`
	TempoPort string `mapstructure:"tempo_port"`
}
