package front

import (
	"fmt"
	"github.com/ziutek/mymysql/mysql"
	_ "github.com/ziutek/mymysql/native" // Native engine
	"os/exec"
	"time"
	"strings"
)

//
type Versions struct {
	Id           int    `json:"id"`
	VersionsName string `json:"versions_name"`
	Versions     string `json:"versions"`
	Environment  string `json:"enviroment"`
	Project      string `json:"project"`
	ProjectName  string `json:"project_name"`
	UpdateTime   string `json:"update_time"`
	Notes        string `json:"notes"`
}

//
type ServersStatus struct {
	Ip      string
	Status  string
	Version string
}

//
type UpgradeResult struct {
	Project     string
	Version     string
	Environment string
	Result      string
}

// func contains(s []string, e string) bool {
func contains(s []interface{}, e interface{}) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}
	return false
}

//获取版本信息
func getVersions(envtype, project string) []Versions {
	db := mysql.New("tcp", "", defaultDbUrl, defaultDbUser, defaultDbPwd, "")
	db.Register("set names utf8")
	err := db.Connect()
	if err != nil {
		fmt.Print("connect db 失败: \n", err)
	}
	defer db.Close()
	// env_type := []string{"dev", "testing", "online"}
	env_type := []interface{}{"dev", "testing", "online"}
	versions := []Versions{}
	if contains(env_type, envtype) {
		stmt, err := db.Prepare("select * from ljops.versions where enviroment = ? and project = ?")
		if err != nil {
			fmt.Println("no such tables!")
		}
		rows, _, err := stmt.Exec(envtype, project)
		for _, row := range rows {
			version := Versions{
				Id:           row.Int(0),
				VersionsName: row.Str(1),
				Versions:     row.Str(2),
				Environment:  row.Str(3),
				Project:      row.Str(4),
				ProjectName:  row.Str(5),
				UpdateTime:   row.Str(6),
				Notes:        row.Str(7),
			}
			versions = append(versions, version)
		}
		return versions
	} else {
		return versions
	}
}

//打包操作
func doPackage(project, env string) bool {
	fmt.Println(project)
	fmt.Println(env)
	return true
}

//升级操作
func doUpgrade(action, version, env, project string) []UpgradeResult {
	//获取操作，版本，环境情况进行操作。
	// fmt.Println(action)
	// fmt.Println(version)
	// fmt.Println(env)
	projectName := getProjectInfo(project, env, version)
	time.Sleep(5 * time.Second)
	upgradeResult := []UpgradeResult{}
	result := UpgradeResult{
		Project:     projectName,
		Version:     version,
		Environment: env,
		Result:      "success",
	}
	upgradeResult = append(upgradeResult, result)
	return upgradeResult
}

func getProjectInfo(project, envtype, version string) string {
	projectName := ""
	db := mysql.New("tcp", "", defaultDbUrl, defaultDbUser, defaultDbPwd, "")
	db.Register("set names utf8")
	err := db.Connect()
	if err != nil {
		fmt.Print("connect db 失败: \n", err)
	}
	defer db.Close()
	env_type := []interface{}{"dev", "testing", "online"}
	if contains(env_type, envtype) {
		stmt, err := db.Prepare("select distinct(project_name) from ljops.versions where project = ?  and  enviroment = ? and versions = ?")
		if err != nil {
			fmt.Println("no such tables!")
		}
		rows, _, err := stmt.Exec(project, envtype, version)
		for _, row := range rows {
			projectName = row.Str(0)
			fmt.Println(projectName)
		}
	}
	return projectName
}

//获取服务器状态，todo
func getServerStatus(envtype string) []ServersStatus {
	Status := []ServersStatus{}
	if envtype == "testing" {
		// get testing versions
		f, err := exec.Command("ls", "/").Output()
		if err != nil {
			fmt.Println(err.Error())
		}
		fmt.Println(string(f))
		return Status
	} else if envtype == "online" {

		// get onlie versions
		fmt.Println("test")
	}
	return Status
}

//执行python脚本，调用zabbix api获取长垣出口流量并返回 
func getChangYuanTrafficData() (string, string, string) {
	f, err := exec.Command("/usr/bin/python", "/data/openops/scripts/get_graph.py").Output()
	if err != nil {
		fmt.Println(err.Error())
	}
	// fmt.Println(reflect.TypeOf(string(f)))

	results := string(f)
	trafficTime := string(strings.Split(results, "+")[0])
	trafficIn := string(strings.Split(results, "+")[1])
	trafficOut := string(strings.Split(results, "+")[2])
	trafficOt := strings.Split(trafficOut, "\n")[0]
	return trafficTime, trafficIn, trafficOt

}


//执行python脚本，调用elasticsearch api获取长垣教学平台服务器访问日志（nginx）
func getData() {
	f, err := exec.Command("/usr/bin/python", "/data/openops/scripts/get_nginx_log_and_analysis.py").Output()
	if err != nil {
		fmt.Println(err.Error())
	}
	// fmt.Println(reflect.TypeOf(string(f)))
	results := string(f)
	fmt.Println(results)
}


//获取服务器可用状态
// 维度：
// 1.服务器端口正常
// 2.能登录
func serverAvailability() {
	fmt.Println("")
}
