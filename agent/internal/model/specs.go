package model

type CpuSpecs struct {
	ModelName        string `json:"modelName"`
	CoreCount        int    `json:"coreCount"`
	LogicalCoreCount int    `json:"logicalCoreCount"`
}
