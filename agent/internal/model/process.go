package model

type Process struct {
	Pid        int32    // id процесса
	Name       string   // название
	ParentPid  *int32   // id родительского процесса
	CPUPercent *float64 // процент использования CPU
	MemoryUsed *uint64  // объем использованного ОЗУ в байтах
}
