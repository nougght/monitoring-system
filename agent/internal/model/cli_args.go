package model

var (
	Args = []ArgInfo{
		ArgConfig,
	}
	ArgConfig = ArgInfo{
		Name:         "config",
		DefaultValue: "",
		Description:  "agent config path",
		Required:     true,
	}
)

type ArgInfo struct {
	Name         string
	DefaultValue string
	Description  string
	Required     bool
}
