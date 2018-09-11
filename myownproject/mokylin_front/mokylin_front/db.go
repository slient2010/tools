package mokylin_front

import (
	"fmt"
	"github.com/ziutek/mymysql/mysql"
	_ "github.com/ziutek/mymysql/native" // Native engine
	_ "log"
	_ "mokylin_common"
	"time"
)

type platAndIdInfo struct {
	Proxy     string `json:"proxy"`
	Alias     string `json:"agent_id"`
	Agent_id  int    `json:"alias"`
	Agent     string `json:"agent"`
	GameAlias string `json:"gamealias"`
	GameName  string `json:"gamename"`
}

type ServersInfo struct {
	Proxy       string `json:"proxy"`
	Agent       string `json:"agent"`
	Project     string `json:"project"`
	GameDomain  string `json:"gameDomain"`
	GameIP      string `json:"gameIP"`
	ServerType  string `json:"serverType"`
	Status      string `json:"status"`
	CreateTime  string `json:"createTime"`
	UpdateTime  string `json:"updateTime"`
	ServerState int    `json:"serverState"`
}

type gameInfo struct {
	Starttime      string
	Game           string
	BackendVersion string
	Clientid       string
	GameStatus     int
}

type agentInfo struct {
	ClientId     string
	GamePath     string
	GameNums     int
	ClientStatus int
}

type mergeInfo struct {
	ClientId string
	Game     string
}

type userInfo struct {
	Id         int    `json:"id"`
	Depart     string `json:"department"`
	Username   string `json:"username"`
	LoginIP    string `json:"loginip"`
	CreateTime string `json:"createtime"`
	UpdateTime string `json:"updatetime"`
}

const (
	// defaultDbUrl = "172.17.0.1:3306"
	// defaultDbUser = "gameops"
	// defaultDbPwd = "work@Mokylin"
	defaultDbUrl  = "127.0.0.1:3306"
	defaultDbUser = "root"
	defaultDbPwd  = ""
)

// get physical servers information
func getServersInfo(proxy, game []string) []ServersInfo {
	// just for debug here, actually its should be query db.
	// fmt.Println(daili)
	db := mysql.New("tcp", "", defaultDbUrl, defaultDbUser, defaultDbPwd, "gameops")
	db.Register("set names utf8")
	err := db.Connect()
	if err != nil {
		fmt.Print("connect db 失败: \n", err)
	}
	defer db.Close()
	stmt, err := db.Prepare("select distinct proxy.alias as proxy_name, proxy.agent, game.alias as project_name, domain, ipaddress, machine_type, status, createtime, updatetime from serverinfo, game, proxy where proxy.agent_id = serverinfo.agent_id and game.game_id = serverinfo.game_id and proxy.proxy=? and game.gamename = ?")
	if err != nil {
		fmt.Println(err)
	}

	servers := []ServersInfo{}
	rows, _, err := stmt.Exec(string(proxy[0]), string(game[0]))
	for _, row := range rows {
		server := ServersInfo{
			Proxy:       row.Str(0),
			Agent:       row.Str(1),
			Project:     row.Str(2),
			GameDomain:  row.Str(3),
			GameIP:      row.Str(4),
			ServerType:  row.Str(5),
			ServerState: row.Int(6),
			CreateTime:  row.Str(7),
			UpdateTime:  row.Str(8),
		}
		servers = append(servers, server)
	}
	return servers
}

func getManageGameInfo(daili, gamename, serverType []string) []gameInfo {
	db := mysql.New("tcp", "", defaultDbUrl, defaultDbUser, defaultDbPwd, "gameops")
	db.Register("set names utf8")
	err := db.Connect()
	if err != nil {
		fmt.Print("connect db 失败: \n", err)
	}
	defer db.Close()
	stmt, err := db.Prepare("select starttime, game, backendversion, clientid, status from gameserverlists where gamename = ? and proxy= ?")
	if err != nil {
		fmt.Println(err)
	}
	infs := []gameInfo{}
	rows, _, err := stmt.Exec(gamename[0], daili[0])
	for _, row := range rows {
		inf := gameInfo{
			Starttime:      row.Str(0),
			Game:           row.Str(1),
			BackendVersion: row.Str(2),
			Clientid:       row.Str(3),
			GameStatus:     row.Int(4),
		}
		infs = append(infs, inf)
	}
	return infs

}

func getManageAgentInfo(daili, gamename, serverType []string) []agentInfo {
	db := mysql.New("tcp", "", defaultDbUrl, defaultDbUser, defaultDbPwd, "gameops")
	db.Register("set names utf8")
	err := db.Connect()
	if err != nil {
		fmt.Print("connect db 失败: \n", err)
	}
	defer db.Close()
	infs := []agentInfo{}

	//part-1
	stmt1, err1 := db.Prepare("select serverinfo.clientid, gamepath, gamenums, serverinfo.status from serverinfo, (select distinct count(clientid) as gamenums, clientid, proxy, gamename from gameserverlists where clientid in (select clientid from serverinfo) group by clientid) as new_tab where serverinfo.clientid = new_tab.clientid and gamename = ? and proxy = ?")
	if err1 != nil {
		fmt.Println(err1)
	}
	rows1, _, err1 := stmt1.Exec(gamename[0], daili[0])

	for _, row1 := range rows1 {
		inf1 := agentInfo{
			ClientId:     row1.Str(0),
			GamePath:     row1.Str(1),
			GameNums:     row1.Int(2),
			ClientStatus: row1.Int(3),
		}
		infs = append(infs, inf1)
	}

	//part-2
	stmt2, err2 := db.Prepare("select clientid, gamepath, '0' as gamenums, status from  (select new_tab.*, gamename from (select serverinfo.*, proxy.proxy from serverinfo, proxy where serverinfo.agent_id = proxy.agent_id) as new_tab, game where new_tab.game_id= game.game_id and new_tab.proxy=game.proxy) as new_tab1 where clientid not in (select clientid from gameserverlists) and gamename = ? and proxy = ?")
	if err != nil {
		fmt.Println(err2)
	}
	rows2, _, err2 := stmt2.Exec(gamename[0], daili[0])
	for _, row2 := range rows2 {
		inf2 := agentInfo{
			ClientId:     row2.Str(0),
			GamePath:     row2.Str(1),
			GameNums:     row2.Int(2),
			ClientStatus: row2.Int(3),
		}
		infs = append(infs, inf2)
	}

	return infs
}

func getSelectClientAndNums(daili, gamename []string) []agentInfo {
	db := mysql.New("tcp", "", defaultDbUrl, defaultDbUser, defaultDbPwd, "gameops")
	db.Register("set names utf8")
	err := db.Connect()
	if err != nil {
		fmt.Print("connect db 失败: \n", err)
	}
	defer db.Close()
	infs := []agentInfo{}

	//part-1
	stmt1, err1 := db.Prepare("select serverinfo.clientid, gamepath, gamenums, serverinfo.status from serverinfo, (select distinct count(clientid) as gamenums, clientid, proxy, gamename from gameserverlists where clientid in (select clientid from serverinfo) group by clientid) as new_tab where serverinfo.clientid = new_tab.clientid and gamename = ? and proxy = ?")
	if err1 != nil {
		fmt.Println(err1)
	}
	rows1, _, err1 := stmt1.Exec(gamename[0], daili[0])

	for _, row1 := range rows1 {
		inf1 := agentInfo{
			ClientId:     row1.Str(0),
			GamePath:     row1.Str(1),
			GameNums:     row1.Int(2),
			ClientStatus: row1.Int(3),
		}
		infs = append(infs, inf1)
	}

	//part-2
	stmt2, err2 := db.Prepare("select clientid, gamepath, '0' as gamenums, status from  (select new_tab.*, gamename from (select serverinfo.*, proxy.proxy from serverinfo, proxy where serverinfo.agent_id = proxy.agent_id) as new_tab, game where new_tab.game_id= game.game_id and new_tab.proxy=game.proxy) as new_tab1 where clientid not in (select clientid from gameserverlists) and gamename = ? and proxy = ?")
	if err != nil {
		fmt.Println(err2)
	}
	rows2, _, err2 := stmt2.Exec(gamename[0], daili[0])
	for _, row2 := range rows2 {
		inf2 := agentInfo{
			ClientId:     row2.Str(0),
			GamePath:     row2.Str(1),
			GameNums:     row2.Int(2),
			ClientStatus: row2.Int(3),
		}
		infs = append(infs, inf2)
	}

	return infs
}

func getManageMergeInfo(daili, gamename, serverType []string) []mergeInfo {
	db := mysql.New("tcp", "", defaultDbUrl, defaultDbUser, defaultDbPwd, "gameops")
	db.Register("set names utf8")
	err := db.Connect()
	if err != nil {
		fmt.Print("connect db 失败: \n", err)
	}
	defer db.Close()
	stmt, err := db.Prepare("select game, clientid from gameserverlists where status = 1 and gamename=? and proxy = ?")
	if err != nil {
		fmt.Println(err)
	}
	infs := []mergeInfo{}
	rows, _, err := stmt.Exec(gamename[0], daili[0])
	for _, row := range rows {
		inf := mergeInfo{
			ClientId: row.Str(0),
			Game:     row.Str(1),
		}
		infs = append(infs, inf)
	}
	return infs
}

func getUserInfo(user string) ([]userInfo, error) {
	userinfo := []userInfo{}
	db := mysql.New("tcp", "", defaultDbUrl, defaultDbUser, defaultDbPwd, "")
	db.Register("set names utf8")
	err := db.Connect()
	if err != nil {
		fmt.Print("connect db 失败: \n", err)
	}
	defer db.Close()
	stmt, err := db.Prepare("select id, department, username, loginip, createtime, updatetime from gameops.users where username = ?")
	if err != nil {
		fmt.Println(err)
		return userinfo, err
	}
	rows, _, err := stmt.Exec(user)
	for _, row := range rows {
		userinf := userInfo{
			Id:         row.Int(0),
			Depart:     row.Str(1),
			Username:   row.Str(2),
			LoginIP:    row.Str(3),
			CreateTime: row.Str(4),
			UpdateTime: row.Str(5),
		}
		return append(userinfo, userinf), nil
	}
	return userinfo, nil
}

func getPlatAndId() ([]platAndIdInfo, error) {
	// connect to db and query.
	platandid := []platAndIdInfo{}
	db := mysql.New("tcp", "", defaultDbUrl, defaultDbUser, defaultDbPwd, "gameops")
	db.Register("set names utf8")
	err := db.Connect()
	if err != nil {
		fmt.Print("connect db 失败: \n", err)
	}
	defer db.Close()
	stmt, err := db.Prepare("select proxy.proxy, proxy.alias, agent_id, proxy.agent, game.alias, game.gamename from proxy, game where game.proxy = proxy.proxy")
	if err != nil {
		fmt.Println(err)
	}
	rows, _, err := stmt.Exec()
	for _, row := range rows {
		platinfo := platAndIdInfo{
			Proxy:     row.Str(0),
			Alias:     row.Str(1),
			Agent_id:  row.Int(2),
			Agent:     row.Str(3),
			GameAlias: row.Str(4),
			GameName:  row.Str(5),
		}
		platandid = append(platandid, platinfo)
	}
	return platandid, nil
}

func getSelectPlatAndId(daili, gamename []string) ([]platAndIdInfo, error) {
	// connect to db and query.
	platandid := []platAndIdInfo{}
	db := mysql.New("tcp", "", defaultDbUrl, defaultDbUser, defaultDbPwd, "gameops")
	db.Register("set names utf8")
	err := db.Connect()
	if err != nil {
		fmt.Print("connect db 失败: \n", err)
	}
	defer db.Close()
	stmt, err := db.Prepare("select * from (select proxy.proxy, proxy.alias as proxy_alias, agent_id, proxy.agent, game.alias as game_alias, game.gamename from proxy, game where game.proxy = proxy.proxy) as new_tab where new_tab.proxy=? and new_tab.gamename=?")
	if err != nil {
		fmt.Println(err)
	}
	rows, _, err := stmt.Exec(daili[0], gamename[0])
	for _, row := range rows {
		platinfo := platAndIdInfo{
			Proxy:     row.Str(0),
			Alias:     row.Str(1),
			Agent_id:  row.Int(2),
			Agent:     row.Str(3),
			GameAlias: row.Str(4),
			GameName:  row.Str(5),
		}
		platandid = append(platandid, platinfo)
	}
	return platandid, nil
}

// api docment
func getApi() string {
	return ""
}

// update user login ip and time
func updateUserLoginIPAndTime(ip, username string, tag bool) error {
	db := mysql.New("tcp", "", defaultDbUrl, defaultDbUser, defaultDbPwd, "")
	db.Register("set names utf8")
	err := db.Connect()
	if err != nil {
		fmt.Print("connect db 失败: \n", err)
	}
	defer db.Close()
	stmt, err := db.Prepare("select * from gameops.users where loginip=? and username = ?")
	if err != nil {
		fmt.Println(err)
		return err
	}
	row, _, _ := stmt.Exec(ip, username)
	nowtime := time.Now().Format("2006-01-02 15:04:05")
	if len(row) == 0 && tag {
		stmt, err := db.Prepare("update gameops.users set loginip=?, updatetime=? where username = ?")
		if err != nil {
			fmt.Println(err)
		}
		_, _, e := stmt.Exec(ip, nowtime, username)
		if e != nil {
			fmt.Println(e)
			return err
		}
	}
	if tag == false {
		stmt, err := db.Prepare("update gameops.users set updatetime=? where username = ?")
		if err != nil {
			fmt.Println(err)
		}
		_, _, e := stmt.Exec(nowtime, username)
		if e != nil {
			fmt.Println(e)
			return err
		}
	}
	return nil
}

func getServers() map[int]bool {
	result := make(map[int]bool)
	///for _, client := range clients {
	///    for _, serverStatProto := range client.ServerStatProtoMap {
	///                serverName := serverStatProto.GetName()

	///                ids, err := mokylin_common.NameToIds(serverName)
	///                if err != nil {
	///                        log.Printf("getServers.NameToIds出错: %s @ %s\n", err, serverName)
	///                } else {
	///                        for operatorID, serverIDs := range ids {
	///                                for _, sids := range serverIDs {
	///                                        combineID := combineOperatorIDAndServerID(operatorID, sids)
	///                                        result[combineID] = true
	///                                }
	///                        }
	///                }
	///        }
	///}
	return result
}

//func getClients() []*client {
// func getSelectClientAndNums(daili, gamename []string) []agentInfo {
func getClients(daili, gamename []string) []agentInfo{
	db := mysql.New("tcp", "", defaultDbUrl, defaultDbUser, defaultDbPwd, "gameops")
	db.Register("set names utf8")
	err := db.Connect()
	if err != nil {
		fmt.Print("connect db 失败: \n", err)
	}
	defer db.Close()
	infs := []agentInfo{}

	//part-1
	stmt1, err1 := db.Prepare("select serverinfo.clientid, gamepath, gamenums, serverinfo.status from serverinfo, (select distinct count(clientid) as gamenums, clientid, proxy, gamename from gameserverlists where clientid in (select clientid from serverinfo) group by clientid) as new_tab where serverinfo.clientid = new_tab.clientid and gamename = ? and proxy = ?")
	if err1 != nil {
		fmt.Println(err1)
	}
	rows1, _, err1 := stmt1.Exec(gamename[0], daili[0])

	for _, row1 := range rows1 {
		inf1 := agentInfo{
			ClientId:     row1.Str(0),
			GamePath:     row1.Str(1),
			GameNums:     row1.Int(2),
			ClientStatus: row1.Int(3),
		}
		infs = append(infs, inf1)
	}

	//part-2
	stmt2, err2 := db.Prepare("select clientid, gamepath, '0' as gamenums, status from  (select new_tab.*, gamename from (select serverinfo.*, proxy.proxy from serverinfo, proxy where serverinfo.agent_id = proxy.agent_id) as new_tab, game where new_tab.game_id= game.game_id and new_tab.proxy=game.proxy) as new_tab1 where clientid not in (select clientid from gameserverlists) and gamename = ? and proxy = ?")
	if err != nil {
		fmt.Println(err2)
	}
	rows2, _, err2 := stmt2.Exec(gamename[0], daili[0])
	for _, row2 := range rows2 {
		inf2 := agentInfo{
			ClientId:     row2.Str(0),
			GamePath:     row2.Str(1),
			GameNums:     row2.Int(2),
			ClientStatus: row2.Int(3),
		}
		infs = append(infs, inf2)
	}

	return infs

}
