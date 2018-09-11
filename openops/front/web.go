package front

import (
	"net/http"
	// "strings"
	"fmt"
	// "io"
	"github.com/gorilla/sessions"
	"log"
	"mime"
	"path/filepath"
	"time"
	// "reflect"
	"github.com/gorilla/context"
	// "os/exec"
	"strings"
	"crypto/sha1"
	"sort"
	"io"
)

func httpServer(addr string) {
	initTmpls()
	http.HandleFunc("/ping", ping)
	http.HandleFunc("/static/", static)
	http.HandleFunc("/", index)
	http.HandleFunc("/index", index)
	http.HandleFunc("/login", login)
	http.HandleFunc("/login.html", login)
	http.HandleFunc("/logout", logout)
	http.HandleFunc("/tables.html", tables)
	http.HandleFunc("/userinfo", user)
	http.HandleFunc("/server", server)
	http.HandleFunc("/upgrade.html", upgrade)
	http.HandleFunc("/testing.html", testing)
	http.HandleFunc("/blank.html", blank)
	http.HandleFunc("/wxplat", wxplat)
	http.HandleFunc("/notifications.html", notifications)

	fmt.Println("master http start!", addr)
	log.Fatal(http.ListenAndServe(addr, context.ClearHandler(http.DefaultServeMux)))

}


// 微信公众平台验证部分
const (
	token = "TheToken"
)

func makeSignature(timestamp, nonce string) string {
        sl := []string{token, timestamp, nonce}
        sort.Strings(sl)
        s := sha1.New()
        io.WriteString(s, strings.Join(sl, ""))
        return fmt.Sprintf("%x", s.Sum(nil))
}

func validateUrl(w http.ResponseWriter, r *http.Request) bool {
        timestamp := strings.Join(r.Form["timestamp"], "")
        nonce := strings.Join(r.Form["nonce"], "")
        signatureGen := makeSignature(timestamp, nonce)

        signatureIn := strings.Join(r.Form["signature"], "")
        if signatureGen != signatureIn {
                return false
        }
        echostr := strings.Join(r.Form["echostr"], "")
        fmt.Fprintf(w, echostr)
        return true


}

func wxplat(w http.ResponseWriter, r *http.Request) {
        r.ParseForm()
        if !validateUrl(w, r) {
                log.Println("Wechat Service: this http request is not from Wechat platform!")
                return
        }
        log.Println("Wechat Service: validateUrl Ok!")
}

// 创建白名单，限制访问
var (
	whiteIps = []string{
		"172.16.41.9",
		"192.168.4.208",
		"192.168.4.168",
		"127.0.0.1",
	}
)

func checkIP(r *http.Request) bool {
	getRemoteAddr := r.Header.Get("x-forwarded-for")
	var remoteAddr string
	 // nginx redirect
	if getRemoteAddr == "" {
		remoteAddr = strings.Split(string(r.RemoteAddr), ":")[0]
	} else {
		remoteAddr = string(getRemoteAddr)
	}

	isSafeIp := false
	for _, ip := range whiteIps {
		if ip == remoteAddr {
			isSafeIp = true
			return isSafeIp
		}
	}
	log.Println(remoteAddr, "is not allowed to visit")
	return isSafeIp
}

// cookie创建
var store = sessions.NewCookieStore([]byte("something-very-secret"))

// 检测web服务器状态
func ping(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("pong"))
}

type dataInfo struct {
	UserInfo    string
	TrafficTime string
	TrafficIn   string
	TrafficOut  string
}

//静态服务器资源
func static(w http.ResponseWriter, r *http.Request) {
	// 验证白名单，拒绝不合法请求
	if checkIP(r) == false {
		w.Write([]byte("不合法的请求!"))
		return
	}
	urlpath := r.URL.Path
	data, err := Asset("html" + urlpath)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	ctype := mime.TypeByExtension(filepath.Ext(urlpath))
	w.Header().Set("Content-Type", ctype)
	w.Write(data)
}

//默认访问请求
func index(w http.ResponseWriter, r *http.Request) {
	// 验证白名单，拒绝不合法请求
	if checkIP(r) == false {
		w.Write([]byte("不合法的请求!"))
		return
	}
	//从cookie中获取用户名，检验用户是否登录，未登录的话，就跳转到登录页面
	Sget(w, r, "username")
	session, _ := store.Get(r, "lju")
	Sessionusername, _ := session.Values["username"]
	var sname string
	if Sessionusername == nil {
		sname = ""
	} else {
		sname = Sessionusername.(string)
	}

	// for _, v := range Sessionusername {
	// 	user := v
	// }
	// 测试下
	getData()
	trafficTime, trafficIn, trafficOut := getChangYuanTrafficData()

	datainfo := dataInfo{
		UserInfo:    sname,
		TrafficTime: trafficTime,
		TrafficIn:   trafficIn,
		TrafficOut:  trafficOut,
	}

	err := templates.ExecuteTemplate(w, "index.html", datainfo)
	// err := templates.ExecuteTemplate(w, "index.html", Sessionusername)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

//获取cookie函数
func Sget(w http.ResponseWriter, r *http.Request, k string) {
	session, _ := store.Get(r, "lju")
	name, _ := session.Values[k].(string)
	if session.Values[k] == nil || name == "" {
		http.Redirect(w, r, "login.html", http.StatusFound)
	}
}

// 登录函数
func login(w http.ResponseWriter, r *http.Request) {
	// 验证白名单，拒绝不合法请求
	if checkIP(r) == false {
		w.Write([]byte("不合法的请求!"))
		return
	}
	if r.Method == "GET" {
		results := ""
		err := templates.ExecuteTemplate(w, "login.html", results)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	} else {
		// 获取用户提交参数并验证
		r.ParseForm()
		username := r.Form["username"]
		password := r.Form["password"]

		name := ""
		passwd := ""
		for _, v := range username {
			name = v
		}
		for _, val := range password {
			passwd = val
		}
		//验证用户
		results := verifyUser(name, passwd)
		if results {
			session, _ := store.Get(r, "lju")
			session.Values["username"] = name
			timeString := time.Now().Format("2006-01-02 15:04:05")
			session.Values["login_time"] = timeString
			session.Save(r, w)
			http.Redirect(w, r, "index", http.StatusFound)
		} else { // 验证不成功，则需重新登录。
			check_result := "用户名或密码不正确！"
			err := templates.ExecuteTemplate(w, "login.html", check_result)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
			}
		}
	}
}

//登出，删除cookie
func logout(w http.ResponseWriter, r *http.Request) {
	// 验证白名单，拒绝不合法请求
	if checkIP(r) == false {
		w.Write([]byte("不合法的请求!"))
		return
	}	
	// 记录最后登录IP和时间
	// 获取当前的客户端IP
	getRemoteAddr := r.Header.Get("x-forwarded-for")
	var remoteAddr string
	if getRemoteAddr == "" {
		remoteAddr = strings.Split(string(r.RemoteAddr), ":")[0]
	}
	// 获取当前时间
	timestamp := time.Now().Unix()
	//格式化为字符串,tm为Time类型
	tm := time.Unix(timestamp, 0)
	currentTime := (tm.Format("2006-01-02 15:04:05"))
	session, _ := store.Get(r, "lju")
	Sessionusername, _ := session.Values["username"].(string)
	saveResult := saveUserIpAndTime(Sessionusername,remoteAddr, currentTime)
	if saveResult != nil {
		w.Write([]byte("保存信息失败!"))
		return
	}
	// clear session
	for key, _ := range session.Values {
		delete(session.Values, key)
	}
	//删除session内属性也需要save
	session.Save(r, w)
	http.Redirect(w, r, "login.html", http.StatusFound)
}

// 设备列表
func tables(w http.ResponseWriter, r *http.Request) {
	// 验证白名单，拒绝不合法请求
	if checkIP(r) == false {
		w.Write([]byte("不合法的请求!"))
		return
	}		
	//从cookie中获取用户名，检验用户是否登录，未登录的话，就跳转到登录页面
	Sget(w, r, "username")
	server_type := r.FormValue("type")
	// 如果设备是交换机，从数据库中查询交换机信息
	if server_type == "switches" {
		results := getDeviceInfo(server_type)
		err := templates.ExecuteTemplate(w, "switches.html", results)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		// 如果设备是物理机，从数据库中查询物理机信息
	} else if server_type == "phy_tables" {
		results := getServerInfo(server_type)
		err := templates.ExecuteTemplate(w, "tables.html", results)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		// 如果设备是虚拟机，从数据库中查询虚拟机信息
	} else if server_type == "vir_tables" {
		results := getServerInfo(server_type)
		err := templates.ExecuteTemplate(w, "tables.html", results)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	} else {
		// default
		server_type := ""
		results := getServerInfo(server_type)
		err := templates.ExecuteTemplate(w, "tables.html", results)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	}
}

//获取用户信息
func user(w http.ResponseWriter, r *http.Request) {
	// 验证白名单，拒绝不合法请求
	if checkIP(r) == false {
		w.Write([]byte("不合法的请求!"))
		return
	}		
	Sget(w, r, "username")
	session, _ := store.Get(r, "lju")
	Sessionusername, _ := session.Values["username"].(string)

	userinfo := getUserInfo(Sessionusername)
	err := templates.ExecuteTemplate(w, "user.html", userinfo)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

//服务器管理
func server(w http.ResponseWriter, r *http.Request) {
	// 验证白名单，拒绝不合法请求
	if checkIP(r) == false {
		w.Write([]byte("不合法的请求!"))
		return
	}		
	Sget(w, r, "username")
	if r.Method == "GET" {
		action := r.FormValue("action")
		if action == "delete" {
			// 用户删除
			account := r.FormValue("account")
			err := serverLoginUserDelete(account)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
			}
			// 删除完成后，跳转到之前页面
			http.Redirect(w, r, "server", http.StatusFound)
		} else if action == "add_account" { // 添加用户
			fmt.Println(action)

		} else { // 响应默认请求
			serverinfo := getServerPrivileges()
			err := templates.ExecuteTemplate(w, "servers.html", serverinfo)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
			}
		}
	} else if r.Method == "POST" {
		username := r.FormValue("name")
		hosts := r.FormValue("hosts")
		fmt.Println(username)
		fmt.Println(hosts)
		err := serverLoginUpdateInfo(username, hosts)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		http.Redirect(w, r, "server", http.StatusFound)
	}
}

func blank(w http.ResponseWriter, r *http.Request) {
	err := templates.ExecuteTemplate(w, "blank.html", "")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

//升级函数
func upgrade(w http.ResponseWriter, r *http.Request) {
	// 验证白名单，拒绝不合法请求
	if checkIP(r) == false {
		w.Write([]byte("不合法的请求!"))
		return
	}		
	Sget(w, r, "username")
	if r.Method == "GET" {
		// 获取环境
		env_type := r.FormValue("type")
		// 获取操作
		action := r.FormValue("action")
		// 获取版本
		version := r.FormValue("version")
		// 获取项目
		project := r.FormValue("project")

		// 参数判断
		if env_type == "" || project == "" {
			// if env_type != "" || action != "" || version != "" || project != "" {
			http.Redirect(w, r, "index.html", http.StatusFound)
		}

		// 打包or升级
		if action == "package" {
			// 打包项目，环境
			result := doPackage(project, env_type)
			//打包结果
			if result {
				redirectUrl := "upgrade.html?type=" + env_type + "&project=" + project
				http.Redirect(w, r, redirectUrl, http.StatusFound)
			}
		} else if action == "upgrade" {
			if action != "" && version != "" && env_type != "" && project != "" {
				result := doUpgrade(action, version, env_type, project)
				err := templates.ExecuteTemplate(w, "upgrade_result.html", result)
				if err != nil {
					http.Error(w, err.Error(), http.StatusInternalServerError)
				}
			}
		} else {
			// 显示当前系统版本
			versions := getVersions(string(env_type), string(project))
			if versions != nil {
				err := templates.ExecuteTemplate(w, "upgrade.html", versions)
				if err != nil {
					http.Error(w, err.Error(), http.StatusInternalServerError)
				}
			}
		}
	}

	////////err := templates.ExecuteTemplate(w, "upgrade.html", "")
	////////if err != nil {
	////////	http.Error(w, err.Error(), http.StatusInternalServerError)
	////////}
}

func testing(w http.ResponseWriter, r *http.Request) {
	versions := getVersions("testing", "teaching")
	status := getServerStatus("testing")
	if r.Method == "GET" {
		env_type := r.FormValue("type")
		if env_type == "testing" {

			fmt.Println(status)
			err := templates.ExecuteTemplate(w, "testing.html", status)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
			}

		} else {
			// default
			err := templates.ExecuteTemplate(w, "testing.html", versions)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
			}
		}
	}
}

func notifications(w http.ResponseWriter, r *http.Request) {
	// 验证白名单，拒绝不合法请求
	if checkIP(r) == false {
		w.Write([]byte("不合法的请求!"))
		return
	}		
	Sget(w, r, "username")
	err := templates.ExecuteTemplate(w, "notifications.html", "")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
