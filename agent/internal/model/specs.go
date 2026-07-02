package model

type CpuSpecs struct {
	ModelName        string `json:"modelName"`
	CoreCount        int    `json:"cores"`
	LogicalCoreCount int    `json:"logicalCores"`
}
