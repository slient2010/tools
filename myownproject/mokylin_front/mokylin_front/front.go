package mokylin_front

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"github.com/golang/protobuf/proto"
	"github.com/gorilla/sessions"
	"github.com/ziutek/mymysql/mysql"
	_ "github.com/ziutek/mymysql/native" // Native engine
	"log"
	"mime"
	"mokylin_common/msg/pb"
	"net"
	"net/http"
	"path/filepath"
	_ "reflect"
	"strconv"
	"strings"
)

var store = sessions.NewCookieStore([]byte("something-very-secret"))

func HttpServer(addr string) {
	initTmpls()
	http.HandleFunc("/ping", ping)
	http.HandleFunc("/static/", static)
	http.HandleFunc("/", index)
	http.HandleFunc("/login", login)
	http.HandleFunc("/logout", logout)
	http.HandleFunc("/loginPost", loginPost)
	http.HandleFunc("/api", api)
	http.HandleFunc("/serverlists", serverlist)
	http.HandleFunc("/gamemanage", gamemanage)
	http.HandleFunc("/userinfo", userInfoQuery)
	http.HandleFunc("/platandid", platAndId)
	http.HandleFunc("/notifications.html", notifications)
	http.HandleFunc("/index", index)
	//	http.HandleFunc("/allgame", allgames)
	http.HandleFunc("/serverCreate", serverCreate)
	http.HandleFunc("/serverCreatePost", serverCreatePost)

	fmt.Println("ms http start!", addr)
	log.Fatal(http.ListenAndServe(addr, nil))
}

// check the server status
func ping(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("pong"))
}

// static resources
func static(w http.ResponseWriter, r *http.Request) {
	urlpath := r.URL.Path
	data, err := Asset("mokylin_front/html" + urlpath)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	ctype := mime.TypeByExtension(filepath.Ext(urlpath))
	w.Header().Set("Content-Type", ctype)
	w.Write(data)
}

func index(w http.ResponseWriter, r *http.Request) {
	//	header := strings.Split(r.Header.Get("Referer"), "/")
	//	refer := header[len(header)-1 : len(header)]

	// check session
	session, _ := store.Get(r, "SESSID")
	username := session.Values["auth"]
	if username == nil {
		http.Redirect(w, r, "/login", 302)
		return
	}

	// get the login IP and update.
	tag := true
	ip, _, e := net.SplitHostPort(r.RemoteAddr)
	if e != nil {
		fmt.Println(e)
	}
	userIP := net.ParseIP(ip)
	if userIP == nil {
		fmt.Println("user ip is not  IP:PORT", r.RemoteAddr)
	}
	updateUserLoginIPAndTime(ip, username.(string), tag)

	// header
	headerr := templates.ExecuteTemplate(w, "header", nil)
	if headerr != nil {
		http.Error(w, headerr.Error(), http.StatusInternalServerError)
	}

	// navigation
	naverr := templates.ExecuteTemplate(w, "nav", map[string]interface{}{
		"Username": username,
	})
	if naverr != nil {
		http.Error(w, naverr.Error(), http.StatusInternalServerError)
	}

	// content
	err := templates.ExecuteTemplate(w, "index.html", map[string]interface{}{
		"Username": username,
	})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	// footer
	footerr := templates.ExecuteTemplate(w, "footer", nil)
	if footerr != nil {
		http.Error(w, footerr.Error(), http.StatusInternalServerError)
	}
}

// machine lists
func serverlist(w http.ResponseWriter, r *http.Request) {
	// game project
	var project string
	gameProject := map[string]string{"gjqt": "古剑奇谭", "rxtl": "热血屠龙", "qmrsy": "秦美人手游"}
	// get the GET params
	r.ParseForm()
	serverdl := r.Form["dl"]
	game := r.Form["game"]
	// check session
	session, _ := store.Get(r, "SESSID")
	username := session.Values["auth"]
	if username == nil {
		http.Redirect(w, r, "/login", 302)
		return
	}

	// header
	headerr := templates.ExecuteTemplate(w, "header", nil)

	if headerr != nil {
		http.Error(w, headerr.Error(), http.StatusInternalServerError)
	}

	// navigation
	naverr := templates.ExecuteTemplate(w, "nav", map[string]interface{}{
		"Username": username,
	})
	if naverr != nil {
		http.Error(w, naverr.Error(), http.StatusInternalServerError)
	}

	// content
	if len(game) == 1 {
		for i, _ := range game {
			project = gameProject[game[i]]
		}
		if project == "" {
			project = "暂无项目"
		}
	}

	results := getServersInfo(serverdl, game)
	err := templates.ExecuteTemplate(w, "serverlists.html", map[string]interface{}{
		"Project": project,
		"results": results,
	})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	// footer
	footerr := templates.ExecuteTemplate(w, "footer", nil)
	if footerr != nil {
		http.Error(w, footerr.Error(), http.StatusInternalServerError)
	}
}

// game manage
func gamemanage(w http.ResponseWriter, r *http.Request) {
	// game project
	var project string
	var operation string
	gameProject := map[string]string{"gjqt": "古剑奇谭", "rxtl": "热血屠龙", "qmrsy": "秦美人手游"}
	// get the GET params
	r.ParseForm()
	serverdl := r.Form["dl"]
	game := r.Form["game"]
	servertype := r.Form["type"]

	// check session
	session, _ := store.Get(r, "SESSID")
	username := session.Values["auth"]
	if username == nil {
		http.Redirect(w, r, "/login", 302)
		return
	}

	// header
	headerr := templates.ExecuteTemplate(w, "header", nil)

	if headerr != nil {
		http.Error(w, headerr.Error(), http.StatusInternalServerError)
	}

	// navigation
	naverr := templates.ExecuteTemplate(w, "nav", map[string]interface{}{
		"Username": username,
	})
	if naverr != nil {
		http.Error(w, naverr.Error(), http.StatusInternalServerError)
	}

	// content
	if len(game) == 1 {
		for i, _ := range game {
			project = gameProject[game[i]]
		}
	}

	switch servertype[0] {
	case "agent":
		results := getManageAgentInfo(serverdl, game, servertype)
		if len(servertype) == 1 {
			operation = servertype[0]
		} else {
			operation = ""
		}
		err := templates.ExecuteTemplate(w, "gamemanage.html", map[string]interface{}{
			"Project":   project,
			"results":   results,
			"Operation": operation,
		})
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	case "gameservers":
		platandids, _ := getSelectPlatAndId(serverdl, game)
		results := getManageGameInfo(serverdl, game, servertype)
		clientandgames := getSelectClientAndNums(serverdl, game)
		backendversions, _ := ResList(serverdl, game)
		if len(servertype) == 1 {
			operation = servertype[0]
		} else {
			operation = ""
		}
		err := templates.ExecuteTemplate(w, "gamemanage.html", map[string]interface{}{
			"Project":         project,
			"PlatformAndIds":  platandids,
			"results":         results,
			"ClientAndGames":  clientandgames,
			"BackendVersions": backendversions,
			"Operation":       operation,
		})
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	case "merge":
		results := getManageMergeInfo(serverdl, game, servertype)
		backendversion := []string{"server_2016-06-16-16-53-48@client_77-276_c.tar", "server_2016-06-14-16-53-48@client_77-275_c.tar"}
		if len(servertype) == 1 {
			operation = servertype[0]
		} else {
			operation = ""
		}
		err := templates.ExecuteTemplate(w, "gamemanage.html", map[string]interface{}{
			"Project":   project,
			"results":   results,
			"res":       backendversion,
			"Operation": operation,
		})
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	}

	// footer
	footerr := templates.ExecuteTemplate(w, "footer", nil)
	if footerr != nil {
		http.Error(w, footerr.Error(), http.StatusInternalServerError)
	}
}

// machine lists
func userInfoQuery(w http.ResponseWriter, r *http.Request) {
	// check session
	session, _ := store.Get(r, "SESSID")
	username := session.Values["auth"]
	if username == nil {
		http.Redirect(w, r, "/login", 302)
		return
	}

	// header
	headerr := templates.ExecuteTemplate(w, "header", nil)

	if headerr != nil {
		http.Error(w, headerr.Error(), http.StatusInternalServerError)
	}

	// navigation
	naverr := templates.ExecuteTemplate(w, "nav", map[string]interface{}{
		"Username": username,
	})
	if naverr != nil {
		http.Error(w, naverr.Error(), http.StatusInternalServerError)
	}

	// content
	str, _ := username.(string)
	results, _ := getUserInfo(str)

	err := templates.ExecuteTemplate(w, "userinfo.html", map[string]interface{}{
		"results": results,
	})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	// footer
	footerr := templates.ExecuteTemplate(w, "footer", nil)
	if footerr != nil {
		http.Error(w, footerr.Error(), http.StatusInternalServerError)
	}
}

func api(w http.ResponseWriter, r *http.Request) {
	// check session
	session, _ := store.Get(r, "SESSID")
	username := session.Values["auth"]
	if username == nil {
		http.Redirect(w, r, "/login", 302)
		return
	}

	// header
	headerr := templates.ExecuteTemplate(w, "header", nil)

	if headerr != nil {
		http.Error(w, headerr.Error(), http.StatusInternalServerError)
	}

	// navigation
	naverr := templates.ExecuteTemplate(w, "nav", map[string]interface{}{
		"Username": username,
	})
	if naverr != nil {
		http.Error(w, naverr.Error(), http.StatusInternalServerError)
	}

	// content
	results := getApi()
	err := templates.ExecuteTemplate(w, "api.html", map[string]interface{}{
		"results": results,
	})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	// footer
	footerr := templates.ExecuteTemplate(w, "footer", nil)
	if footerr != nil {
		http.Error(w, footerr.Error(), http.StatusInternalServerError)
	}
}

func platAndId(w http.ResponseWriter, r *http.Request) {
	// check session
	session, _ := store.Get(r, "SESSID")
	username := session.Values["auth"]
	if username == nil {
		http.Redirect(w, r, "/login", 302)
		return
	}

	// header
	headerr := templates.ExecuteTemplate(w, "header", nil)

	if headerr != nil {
		http.Error(w, headerr.Error(), http.StatusInternalServerError)
	}

	// navigation
	naverr := templates.ExecuteTemplate(w, "nav", map[string]interface{}{
		"Username": username,
	})
	if naverr != nil {
		http.Error(w, naverr.Error(), http.StatusInternalServerError)
	}

	// content
	results, _ := getPlatAndId()
	err := templates.ExecuteTemplate(w, "platandid.html", map[string]interface{}{
		"results": results,
	})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	// footer
	footerr := templates.ExecuteTemplate(w, "footer", nil)
	if footerr != nil {
		http.Error(w, footerr.Error(), http.StatusInternalServerError)
	}
}

func notifications(w http.ResponseWriter, r *http.Request) {
	err := templates.ExecuteTemplate(w, "notifications.html", "")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func clients(w http.ResponseWriter, r *http.Request) {
	// slice not use now, disable
	// htmlClients := []HtmlClient{}
	// connect to db and get clients info
	// temp disable clients := GetClientsIncludeNotConnected()
	clients := []string{"client1", "client2", "client3"}
	err := templates.ExecuteTemplate(w, "logs.html", clients)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func login(w http.ResponseWriter, r *http.Request) {
	session, _ := store.Get(r, "SESSID")
	username := session.Values["auth"]
	if username == nil {
		err := templates.ExecuteTemplate(w, "login.html", map[string]interface{}{
			"Failed":   "",
			"Username": username,
		})
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	} else {
		err := templates.ExecuteTemplate(w, "index.html", map[string]interface{}{
			"Username": username,
		})
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	}
}

func loginPost(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		http.Redirect(w, r, "/login", 302)
	case "POST":
		username := r.FormValue("email")
		if username == "" {
			err := templates.ExecuteTemplate(w, "login.html", map[string]interface{}{
				"Failed":   "用户名不能为空!",
				"Username": username,
			})
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
			}
			//			w.Write([]byte("用户名不能为空"))
			return
		}

		password := r.FormValue("password")
		if password == "" {
			err := templates.ExecuteTemplate(w, "login.html", map[string]interface{}{
				"Failed":   "密码不能为空!",
				"Username": username,
			})
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
			}
			//			w.Write([]byte("密码不能为空"))
			return
		}
		encdata, _ := Encrypt(string(username) + string("@@") + string(password))
		result := verifyUser(encdata)
		if !result {
			err := templates.ExecuteTemplate(w, "login.html", map[string]interface{}{
				"Failed":   "用户名或者密码错误!",
				"Username": username,
			})
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
			}

			//			w.Write([]byte("用户名或者密码错误"))
			return
		}

		session, _ := store.Get(r, "SESSID")
		session.Values["auth"] = username
		session.Save(r, w)
		http.Redirect(w, r, "/index", 302)
		return

		////////err := templates.ExecuteTemplate(w, "index.html", "")
		////////if err != nil {
		////////	http.Error(w, err.Error(), http.StatusInternalServerError)
		////////}
	}
}

func verifyUser(userdata []byte) bool {
	// decrypt
	decdata, _ := Decrypt(userdata)
	// connect to db
	db := mysql.New("tcp", "", defaultDbUrl, defaultDbUser, defaultDbPwd, "")
	db.Register("set names utf8")
	err := db.Connect()
	if err != nil {
		fmt.Print("connect db 失败: \n", err)
	}
	defer db.Close()
	passenc := md5.New()
	passenc.Write([]byte(strings.Split(decdata, "@@")[1]))
	encresult := hex.EncodeToString(passenc.Sum(nil))
	stmt, err := db.Prepare("select username, password from gameops.users where username=? and password = ?")
	if err != nil {
		return false
	}
	rows, _, _ := stmt.Exec(strings.Split(decdata, "@@")[0], encresult)
	if err != nil || len(rows) != 1 {
		return false
	}
	return true
}

func logout(w http.ResponseWriter, r *http.Request) {
	session, _ := store.Get(r, "SESSID")
	username := session.Values["auth"]

	// get the login IP and update.
	ip, _, e := net.SplitHostPort(r.RemoteAddr)
	if e != nil {
		fmt.Println(e)
	}
	userIP := net.ParseIP(ip)
	if userIP == nil {
		fmt.Println("user ip is not  IP:PORT", r.RemoteAddr)
	}
	tag := false
	updateUserLoginIPAndTime(ip, username.(string), tag)
	delete(session.Values, "auth")
	session.Options.MaxAge = -1
	session.Save(r, w)
	http.Redirect(w, r, "/login", 302)
}

func serverCreate(w http.ResponseWriter, r *http.Request) {
	// check session
	session, _ := store.Get(r, "SESSID")
	username := session.Values["auth"]
	if username == nil {
		http.Redirect(w, r, "/login", 302)
		return
	}

	// header
	////////headerr := templates.ExecuteTemplate(w, "header", nil)

	////////if headerr != nil {
	////////	http.Error(w, headerr.Error(), http.StatusInternalServerError)
	////////}

	////////// navigation
	////////naverr := templates.ExecuteTemplate(w, "nav", map[string]interface{}{
	////////	"Username": username,
	////////})
	////////if naverr != nil {
	////////	http.Error(w, naverr.Error(), http.StatusInternalServerError)
	////////}

	// content
	results, _ := getPlatAndId()
	err := templates.ExecuteTemplate(w, "servercreate.html", map[string]interface{}{
		"results": results,
	})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	////////// footer
	////////footerr := templates.ExecuteTemplate(w, "footer", nil)
	////////if footerr != nil {
	////////	http.Error(w, footerr.Error(), http.StatusInternalServerError)
	////////}
}

func serverCreatePost(w http.ResponseWriter, r *http.Request) {
	platformIdStr := r.FormValue("platformId")
	if platformIdStr == "" {
		w.Write([]byte("平台ID不能为空!"))
		return
	}

	serverIdsStr := r.FormValue("serverIds") // 1-2-4&&5&&9-10
	if serverIdsStr == "" {
		w.Write([]byte("服务器ID不能为空!"))
		return
	}

	serverCreateTimesStr := r.FormValue("serverCreateTimes")
	if serverCreateTimesStr == "" {
		w.Write([]byte("开服时间不能为空!"))
		return
	}

	resUrl := r.FormValue("resUrl")
	if resUrl == "" {
		w.Write([]byte("资源路径不能为空!"))
		return
	}

	clientId := r.FormValue("clientId")
	if clientId == "" {
		w.Write([]byte("clientId不能为空!"))
		return
	}

	serverIds := strings.Split(serverIdsStr, "&&")
	serverCreateTimes := strings.Split(serverCreateTimesStr, "&&")

	l := len(serverIds)
	if l != len(serverCreateTimes) {
		http.Error(w, "len(serverIds) != len(serverCreateTimes)", http.StatusInternalServerError)
		return
	}

	serverIdCheckMap := make(map[int]string)
	serverIdsInt := []int{}
	domainSidsInt := []int{}
	for _, sidsStr := range serverIds {
		for _, serverIDStr := range strings.Split(sidsStr, "-") {

			if platformIdStr != "199" {
				sidInt, err := strconv.Atoi(serverIDStr)
				if err != nil {
					w.Write([]byte("服务器id格式非法!"))
					return
				}

				if _, ok := serverIdCheckMap[sidInt]; ok {
					w.Write([]byte("服务器id重复!"))
					return
				}

				serverIdCheckMap[sidInt] = sidsStr
				serverIdsInt = append(serverIdsInt, sidInt)
				domainSidsInt = append(domainSidsInt, sidInt)
			} else {
				sidx := strings.Index(serverIDStr, "(")
				dsidStr := serverIDStr[:sidx]
				dsid, err := strconv.Atoi(dsidStr)
				if err != nil {
					w.Write([]byte("服务器id格式非法!"))
					return
				}

				domainSidsInt = append(domainSidsInt, dsid)

				eidx := strings.Index(serverIDStr, ")")
				if eidx != len(serverIDStr)-1 {
					w.Write([]byte("服务器id格式非法!"))
					return
				}

				sidStr := serverIDStr[sidx+1 : eidx]
				sidInt, err := strconv.Atoi(sidStr)
				if err != nil {
					w.Write([]byte("服务器id格式非法!"))
					return
				}

				if _, ok := serverIdCheckMap[sidInt]; ok {
					w.Write([]byte("服务器id重复!"))
					return
				}

				serverIdsInt = append(serverIdsInt, sidInt)
			}
		}
	}

	platformId, err := strconv.Atoi(platformIdStr)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	existingServers := getServers()
	idsMap := make(map[int][]int)
	idsMap[platformId] = serverIdsInt

	propServerVal := "server="
	for i := 0; i < l; i++ {
		propServerVal += platformIdStr + "@" + serverCreateTimes[i] + "#" + serverIds[i]
		if i < l-1 {
			// propServerVal += server.PropServerSep
			propServerVal += "@#@"
		}
	}

	for i := 0; i < len(serverIdsInt); i++ {
		combineID := combineOperatorIDAndServerID(platformId, serverIdsInt[i])
		if _, existed := existingServers[combineID]; existed {
			w.Write([]byte("服务器已存在"))
			return
		}
	}

	// platformDomain := platformIdToDomain(platformId)

	jobProto := &pb.CreateServerJobProto{
		ResUrl: proto.String(resUrl),
		Cfg:    proto.String(propServerVal),
		//		Domains: platformDomainMap,
	}

	replyData := submitCreateServerJob(clientId, jobProto)

	w.Write(replyData)
}

func combineOperatorIDAndServerID(operatorID, serverID int) int {
	return (operatorID << 16) | serverID
}
