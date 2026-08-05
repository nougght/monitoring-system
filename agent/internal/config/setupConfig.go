package config

type SetupConfig struct {
	EnrollmentKey     string `yaml:"enrollment_key"`
	EnrollmentAddress string `yaml:"enrollment_address"`
	ServerAddress     string `yaml:"server_address"`
	CaPath            string `yaml:"ca_path"`
	KeyPath           string `yaml:"key_path"`
	CertPath          string `yaml:"cert_path"`
	// TODO: add agent mode
}
