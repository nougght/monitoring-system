package agent

type AgentSetupConfig struct {
	EnrollmentKey     string `yaml:"enrollment_key"`
	EnrollmentAddress string `yaml:"enrollment_address"`
	ServerAddress     string `yaml:"server_address"`
}
