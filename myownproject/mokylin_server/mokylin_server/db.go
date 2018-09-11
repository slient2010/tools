package mokylin_server

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/ziutek/mymysql/mysql"
	_ "github.com/ziutek/mymysql/native" // Native engine
	"log"
	"strconv"
	"strings"
	"time"
)

//func CheckClient(clientId string) (int, error) {
//	// connect to database and query db.
//	defaultDbUrl := databaseip + ":" + strconv.Itoa(databaseport)
//	defaultDbUser := databaseuser
//	defaultDbPwd := databasepass
//	db := mysql.New("tcp", "", defaultDbUrl, defaultDbUser, defaultDbPwd, "")
//	db.Register("set names utf8")
//	err := db.Connect()
//	if err != nil {
//		fmt.Print("connect db 失败: \n", err)
//	}
//	defer db.Close()
//	rows, _, err := db.Query("select * from mops.serverinfo where serverid like '%s'", clientId)
//	if err != nil {
//		log.Fatal(err)
//	}
//
//	return len(rows), nil
//}

func CheckClient(clientId, domain string) (int, error) {
	// connect to database and query db.
	defaultDbUrl := databaseip + ":" + strconv.Itoa(databaseport)
	defaultDbUser := databaseuser
	defaultDbPwd := databasepass
	db := mysql.New("tcp", "", defaultDbUrl, defaultDbUser, defaultDbPwd, "")
	db.Register("set names utf8")
	err := db.Connect()
	if err != nil {
		fmt.Print("connect db 失败: \n", err)
	}
	defer db.Close()
	// check before save
	rows, _, err := db.Query("select * from gameops.serverinfo where clientid like '%s' and domain like '%s'", clientId, domain)
	if err != nil {
		log.Fatal(err)
	}
	return len(rows), nil
}

func SaveClientInfo(clientId, ipaddress, domain, gamepath string) error {
	// connect to database and query db.
	defaultDbUrl := databaseip + ":" + strconv.Itoa(databaseport)
	defaultDbUser := databaseuser
	defaultDbPwd := databasepass
	db := mysql.New("tcp", "", defaultDbUrl, defaultDbUser, defaultDbPwd, "")
	db.Register("set names utf8")
	err := db.Connect()
	if err != nil {
		fmt.Print("connect db 失败: \n", err)
	}
	defer db.Close()

	var dat map[string]interface{}
	json.Unmarshal([]byte(ipaddress), &dat)
	////////if err := json.Unmarshal([]byte(ipaddress), &dat); err == nil {
	////////	fmt.Println(dat["ip"])
	////////}

	// check before save
	////////rows, _, err := db.Query("select * from gameops.serverinfo where clientid like '%s' and domain like '%s'", clientId, domain)
	////////if err != nil {
	////////	log.Fatal(err)
	////////}

	rows, _ := CheckClient(clientId, domain)

	// save or update
	if rows == 0 {
		// save client into project
		c := strings.Split(clientId, "-")
		if len(c) == 4 {
			proxy := c[0]
			agent_alias := c[1]
			gamename := c[2]

			// get gameid
			rs, _, err := db.Query("select game_id from gameops.game where proxy = '%s' and gamename = '%s'", proxy, gamename)
			if err != nil {
				log.Fatal(err)
			}

			game_id := rs[0].Int(0)

			// get agentid
			res, _, err := db.Query("select agent_id from gameops.proxy where proxy = '%s' and agent_alias = '%s'", proxy, agent_alias)
			if err != nil {
				log.Fatal(err)
			}

			agent_id := res[0].Int(0)

			// insert into db
			ins, err := db.Prepare("insert into gameops.serverinfo values (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)")
			if err != nil {
				log.Fatal(err)
			}
			_, err = ins.Run("", agent_id, game_id, clientId, domain, dat["ip"], gamepath, "gameserver", 1, time.Now().Format("2006-01-02 15:04:05"), time.Now().Format("2006-01-02 15:04:05"))
			if err != nil {
				log.Fatal(err)
			}

			////////// save client into project
			////////c := strings.Split(clientId, "-")
			////////proxy := c[0]
			////////gamename := c[1]

			////////rs, _, err := db.Query("select game_id from gameops.game where proxy = '%s' and gamename = '%s'", proxy, gamename)
			////////if err != nil {
			////////	log.Fatal(err)
			////////}
			////////	if len(rs) == 1 {
			////////		ins, err := db.Prepare("insert into gameops.serverlists values (?, ?, ?, ?, ?, ?, ?, ?, ?)")
			////////		if err != nil {
			////////			log.Fatal(err)
			////////		}
			////////		_, err = ins.Run("", "", game_id, domain, dat["ip"], "gameserver", 1, time.Now().Format("2006-01-02 15:04:05"), time.Now().Format("2006-01-02 15:04:05"))
			////////		if err != nil {
			////////			log.Fatal(err)
			////////		}
			////////	}

			return nil
		} else {
			return errors.New("client id name format error!")
		}

	} else {
		// update
		ins, err := db.Prepare("update gameops.serverinfo set ipaddress=?, status=?, updatetime=? where clientid = ? and domain = ?")
		if err != nil {
			log.Fatal(err)
		}
		_, err = ins.Run(dat["ip"], 1, time.Now().Format("2006-01-02 15:04:05"), clientId, domain)
		if err != nil {
			log.Fatal(err)
		}

		return nil
	}
	return errors.New("unknow error")
}

func UpdateClientUnRegister(clientId, domain string) error {
	// connect to database and query db.
	defaultDbUrl := databaseip + ":" + strconv.Itoa(databaseport)
	defaultDbUser := databaseuser
	defaultDbPwd := databasepass
	db := mysql.New("tcp", "", defaultDbUrl, defaultDbUser, defaultDbPwd, "")
	db.Register("set names utf8")
	err := db.Connect()
	if err != nil {
		fmt.Print("connect db 失败: \n", err)
	}
	defer db.Close()

	// client offline
	ins, err := db.Prepare("update gameops.serverinfo set status=0, updatetime=? where clientid = ? and domain = ?")
	if err != nil {
		log.Fatal(err)
	}
	_, err = ins.Run(time.Now().Format("2006-01-02 15:04:05"), clientId, domain)
	if err != nil {
		log.Fatal(err)
		return err
	}
	return nil
}
