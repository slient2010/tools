package master

import (
	"github.com/jinzhu/gorm"
)

// server basic infomation
type ServerInfo struct {
	gorm.Model
	Name      string
	Ipaddress string
	CPU       int
	Mem       int
	Disk      int
	Business  string
	UUID      string
}
