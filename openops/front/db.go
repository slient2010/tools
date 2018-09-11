package front

import (
	"fmt"
	"github.com/ziutek/mymysql/mysql"
	_ "github.com/ziutek/mymysql/native" // Native engine
	"time"
	// "reflect"
	// "strconv"
)

// 服务器设备信息数据结构
type ServersInfo struct {
	Id           int    `json:"id"`
	ServerName   string `json:"server_name"`
	ServerIP     string `json:"IP"`
	ManageIP     string `json:"manage_ip"`
	Vendor       string `json:"vendor"`
	ServerType   string `json:"server_type"`
	Memory       int    `json:"memory"`
	CPU          string `json:"cpu"`
	Interface    string `json:"interface"`
	Disk         string `json:"disk"`
	ServerDetail string `json:"server_detail"`
	ServerSystem string `json:"server_system"`
	OsType       int    `json:"os_type"`
	OsTypeId     int    `json:"os_type_id"`
	Status       string `json:"status"`
	Position     string `json:"position"`
	CreatedTime  string `json:"created_time"`
	UpdateTime   string `json:"update_time"`
	Notes        string `json:"notes"`
}

// 交换机设备信息数据结构
type DeviceInfo struct {
	Id             int    `json:"id"`
	DeviceName     string `json:"device_name"`
	DeviceIP       string `json:"device_ip"`
	DeviceType     string `json:"device_type"`
	Vendor         string `json:"vendor"`
	DeviceNo       string `json:"device_no"`
	Username       string `json:"username"`
	Password       string `json:"password"`
	DeviceDetail   string `json:"device_detail"`
	DeviceStatus   string `json:"device_status"`
	DeviceLocation string `json:"device_location"`
	DevicePosition string `json:"device_position"`
	CreatedTime    string `json:"created_time"`
	UpdateTime     string `json:"update_time"`
	Notes          string `json:"notes"`
}

type logInfo struct {
	Id       int    `json:"id"`
	Lid      string `json:"lid"`
	Logname  string `json:logname"`
	Game     string `json:"game"`
	Thread   string `json:"thread"`
	Loglevel string `json:"loglevel"`
	Logtime  string `json:"logtime"`
	Loginfo  string `json:"loginfo"`
}

// 用户信息数据结构
type userInfo struct {
	Id            int    `json:"id"`
	Realname      string `json:"relaname"`
	Email         string `json:"email"`
	Phone         string `json:"phone"`
	LastIP        string `json:"lastlogip"`
	LastLoginTime string `json:"lastlogintime"`
	Notes         string `json:"notes"`
}

// 服务器权限数据结构
type serverPrivileges struct {
	Name     string `json:"name"`
	Hosts    string `json:"hosts"`
	Rsa_key  string `json:"rsa_key"`
	LastIP   string `json:"lastlogip"`
	LastTime string `json:"lasttime"`
	Created  string `json:"created"`
	Time     int    `json:"time"`
}

const (
	defaultDbUrl  = "172.17.0.8:3306"
	defaultDbUser = "ops"
	defaultDbPwd  = "work@Ljkj"
)

//查询交换机设备信息
func getDeviceInfo(queryType string) []DeviceInfo {
	db := mysql.New("tcp", "", defaultDbUrl, defaultDbUser, defaultDbPwd, "")
	db.Register("set names latin1")
	err := db.Connect()
	if err != nil {
		fmt.Print("connect db 失败: \n", err)
	}
	defer db.Close()
	stmt, err := db.Prepare("select * from ljops.netdeviceinfo")
	if err != nil {
		fmt.Println("no such tables!")
	}
	devices := []DeviceInfo{}
	rows, _, err := stmt.Exec()
	for _, row := range rows {
		device := DeviceInfo{
			Id:             row.Int(0),
			DeviceName:     row.Str(1),
			DeviceIP:       row.Str(2),
			DeviceType:     row.Str(3),
			Vendor:         row.Str(4),
			DeviceNo:       row.Str(5),
			Username:       row.Str(6),
			Password:       row.Str(7),
			DeviceDetail:   row.Str(8),
			DeviceStatus:   row.Str(9),
			DeviceLocation: row.Str(10),
			DevicePosition: row.Str(11),
			CreatedTime:    row.Str(12),
			UpdateTime:     row.Str(13),
			Notes:          row.Str(14),
		}
		//CreatedTime:  time.Unix(row.Int64(16), 0).Format("2006-01-02 15:04:05"),
		//UpdateTime:   time.Unix(row.Int64(17), 0).Format("2006-01-02 15:04:05"),
		devices = append(devices, device)
	}
	return devices
}

//获取服务器信息函数
func getServerInfo(queryType string) []ServersInfo {
	db := mysql.New("tcp", "", defaultDbUrl, defaultDbUser, defaultDbPwd, "")
	db.Register("set names utf8")
	err := db.Connect()
	if err != nil {
		fmt.Print("connect db 失败: \n", err)
	}
	defer db.Close()
	servers := []ServersInfo{}
	//物理机
	if queryType == "phy_tables" {
		stmt, err := db.Prepare("select * from ljops.serverinfo where os_type=0")
		if err != nil {
		}
		rows, _, err := stmt.Exec()
		for _, row := range rows {
			server := ServersInfo{
				Id:           row.Int(0),
				ServerName:   row.Str(1),
				ServerIP:     row.Str(2),
				ManageIP:     row.Str(3),
				Vendor:       row.Str(4),
				ServerType:   row.Str(5),
				Memory:       row.Int(6),
				CPU:          row.Str(7),
				Interface:    row.Str(8),
				Disk:         row.Str(9),
				ServerDetail: row.Str(10),
				ServerSystem: row.Str(11),
				OsType:       row.Int(12),
				OsTypeId:     row.Int(13),
				Status:       row.Str(14),
				Position:     row.Str(15),
				CreatedTime:  row.Str(16),
				UpdateTime:   row.Str(17),
				Notes:        row.Str(18),
			}
			servers = append(servers, server)
		}
		return servers
		//虚拟机
	} else if queryType == "vir_tables" {
		stmt, err := db.Prepare("select * from ljops.serverinfo where os_type=1")
		if err != nil {
		}
		servers := []ServersInfo{}
		rows, _, err := stmt.Exec()
		for _, row := range rows {
			server := ServersInfo{
				Id:           row.Int(0),
				ServerName:   row.Str(1),
				ServerIP:     row.Str(2),
				ManageIP:     row.Str(3),
				Vendor:       row.Str(4),
				ServerType:   row.Str(5),
				Memory:       row.Int(6),
				CPU:          row.Str(7),
				Interface:    row.Str(8),
				Disk:         row.Str(9),
				ServerDetail: row.Str(10),
				ServerSystem: row.Str(11),
				OsType:       row.Int(12),
				OsTypeId:     row.Int(13),
				Status:       row.Str(14),
				Position:     row.Str(15),
				CreatedTime:  row.Str(16),
				UpdateTime:   row.Str(17),
				Notes:        row.Str(18),
			}
			servers = append(servers, server)
		}
		return servers
	}
	return servers
}

//查询服务器权限信息
func getServerPrivileges() []serverPrivileges {
	db := mysql.New("tcp", "", defaultDbUrl, defaultDbUser, defaultDbPwd, "")
	db.Register("set names latin1")
	err := db.Connect()
	if err != nil {
		fmt.Print("connect db 失败: \n", err)
	}
	defer db.Close()
	stmt, err := db.Prepare("select * from ljops.servers_manage")
	if err != nil {
		fmt.Println("no such tables!")
	}

	serverprivilege := []serverPrivileges{}
	rows, _, err := stmt.Exec()
	for _, row := range rows {
		s := serverPrivileges{
			Name:     row.Str(0),
			Hosts:    row.Str(1),
			Rsa_key:  row.Str(2),
			LastIP:   row.Str(3),
			LastTime: row.Str(4),
			Created:  time.Unix(int64(row.Int(5)), 0).Format("2006-01-02 15:04:05"),
			// Created:       row.Int(5),
			Time: row.Int(6),
		}
		serverprivilege = append(serverprivilege, s)
	}
	return serverprivilege
}

//获取日志信息
func getLogInfo() []logInfo {
	db := mysql.New("tcp", "", defaultDbUrl, defaultDbUser, defaultDbPwd, "")
	db.Register("set names utf8")
	err := db.Connect()
	if err != nil {
		//        log.Printf("connect db 失败: %s\n", err)
		fmt.Print("connect db 失败: \n", err)
		//        Println(err)
	}
	defer db.Close()
	stmt, err := db.Prepare("select * from mysite.polls_datasave limit 10")
	if err != nil {
		fmt.Println(err)
	}
	info := []logInfo{}
	rows, _, err := stmt.Exec()
	for _, row := range rows {
		inf := logInfo{
			Id:       row.Int(0),
			Lid:      row.Str(1),
			Logname:  row.Str(2),
			Game:     row.Str(3),
			Thread:   row.Str(4),
			Loglevel: row.Str(5),
			Logtime:  row.Str(6),
			Loginfo:  row.Str(7),
		}

		info = append(info, inf)
	}
	return info
}

// 用户验证
func verifyUser(username string, password string) bool {
	db := mysql.New("tcp", "", defaultDbUrl, defaultDbUser, defaultDbPwd, "")
	db.Register("set names utf8")
	err := db.Connect()
	if err != nil {
		//        log.Printf("connect db 失败: %s\n", err)
		fmt.Print("connect db 失败: \n", err)
		//        Println(err)
	}
	defer db.Close()
	stmt, err := db.Prepare("select count(*) from ljops.user where email=? and password=?")
	if err != nil {
		fmt.Println(err)
	}
	rows, _, _ := stmt.Exec(username, password)
	for _, row := range rows {
		for _, r := range row {
			count := int64(0)
			if r == count {
				return false
			}
		}
	}
	return true
}

//获取用户信息
func getUserInfo(username string) []userInfo {
	db := mysql.New("tcp", "", defaultDbUrl, defaultDbUser, defaultDbPwd, "")
	db.Register("set names utf8")
	err := db.Connect()
	if err != nil {
		//        log.Printf("connect db 失败: %s\n", err)
		fmt.Print("connect db 失败: \n", err)
		//        Println(err)
	}
	defer db.Close()
	stmt, err := db.Prepare("select id, realname, email, phone, lastloginip, lastlogin, notes from ljops.user where email=?")
	if err != nil {
		fmt.Println(err)
	}
	userinfo := []userInfo{}
	rows, _, err := stmt.Exec(username)
	for _, row := range rows {
		inf := userInfo{
			Id:            row.Int(0),
			Realname:      row.Str(1),
			Email:         row.Str(2),
			Phone:         row.Str(3),
			LastIP:        row.Str(4),
			LastLoginTime: row.Str(5),
			Notes:         row.Str(6),
		}

		userinfo = append(userinfo, inf)
	}
	return userinfo
}

//删除服务器登录用户
func serverLoginUserDelete(username string) error {
	db := mysql.New("tcp", "", defaultDbUrl, defaultDbUser, defaultDbPwd, "")
	db.Register("set names utf8")
	err := db.Connect()
	if err != nil {
		//        log.Printf("connect db 失败: %s\n", err)
		fmt.Print("connect db 失败: \n", err)
		//        Println(err)
	}
	defer db.Close()
	stmt, err := db.Prepare("delete from ljops.servers_manage where name=?")
	if err != nil {
		fmt.Println(err)
	}
	_, _, e := stmt.Exec(username)
	return e
}

//更新服务器登录用户服务器范围
func serverLoginUpdateInfo(username, hosts string) error {
	db := mysql.New("tcp", "", defaultDbUrl, defaultDbUser, defaultDbPwd, "")
	db.Register("set names utf8")
	err := db.Connect()
	if err != nil {
		fmt.Print("connect db 失败: \n", err)
	}
	defer db.Close()
	stmt, err := db.Prepare("update ljops.servers_manage set hosts=? where name=?")
	if err != nil {
		fmt.Println(err)
	}
	_, _, e := stmt.Exec(hosts, username)
	return e
}


func saveUserIpAndTime(username, clientip, clienttime string) error {
	db := mysql.New("tcp", "", defaultDbUrl, defaultDbUser, defaultDbPwd, "")
	db.Register("set names utf8")
	err := db.Connect()
	if err != nil {
		fmt.Print("connect db 失败: \n", err)
	}
	defer db.Close()
	stmt, err := db.Prepare("update ljops.user set lastloginip = ?, lastlogin = ? where email = ?")
	if err != nil {
		fmt.Println(err)
	}
	_, _, e := stmt.Exec(clientip, clienttime, username)
	return e
}
