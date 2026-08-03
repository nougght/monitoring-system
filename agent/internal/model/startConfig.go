package model

type StartConfig struct {
	EnrollmentKey string `yaml:"enrollment_key"`
	CaPath        string `yaml:"ca_path"`
	KeyPath       string `yaml:"key_path"`
	CertPath      string `yaml:"cert_path"`
}
