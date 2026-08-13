package model

type AgentCommand string

const (
	AgentCommandRun    AgentCommand = "run"
	AgentCommandEnroll AgentCommand = "enroll"
)
const (
	EnrollmentKeyUsed string = "enrollment_key_used"
)

// TODO: replace with normal paths
const (
	DefaultCAPath        string = "./creds/ca.crt"
	DefaultAgentCertPath string = "./creds/agent.crt"
	DefaultAgentKeyPath  string = "./creds/agent.key"
)
