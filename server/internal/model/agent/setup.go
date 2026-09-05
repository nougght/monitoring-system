package agent_model

var (
	AgentFilePaths = map[string]string{
		"agent.exe":    "../bin/agent.exe",
		"creds/ca.crt": "creds/root-ca.crt",
	}
)

const SetupConfigFileName = "setupconfig.yaml"

type AgentSetupConfig struct {
	EnrollmentKey     string `yaml:"enrollment_key"`
	EnrollmentAddress string `yaml:"enrollment_address"`
	ServerAddress     string `yaml:"server_address"`
}
