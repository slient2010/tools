package client

import (
	"github.com/cloudfoundry/gosigar"
	"strconv"
)

func Cpu() (string, error) {
	cpu := &sigar.Cpu{}
	cpu.Get()

	total := cpu.Total()
	idle := cpu.Idle

	freePer := (total - idle) * 100 / total
	return strconv.FormatUint(freePer, 10) + "%", nil
}

func Mem() (string, error) {
	mem := &sigar.Mem{}
	mem.Get()

	total := mem.Total
	free := mem.ActualFree

	freePer := (total - free) * 100 / total
	return strconv.FormatUint(freePer, 10) + "%", nil
}

func Disk() (string, error) {
	disk := &sigar.FileSystemUsage{}
	baseDir := "/data/tehang/apps"
	disk.Get(baseDir)

	total := disk.Total
	free := disk.Free

	freePer := (total - free) * 100 / total

	return strconv.FormatUint(freePer, 10) + "%", nil
}
