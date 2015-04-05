package main

import (
	"flag"
	"database/sql"
	"github.com/go-martini/martini"
	_ "code.google.com/p/odbc"
	"github.com/martini-contrib/sessions"
	"github.com/martini-contrib/render"
	"net/http"
	"github.com/hcsoft/helpsystem/auth"
	erutil "github.com/hcsoft/helpsystem/error"
	"github.com/hcsoft/helpsystem/helpmaker"
	"github.com/hcsoft/helpsystem/admin"
	"github.com/larspensjo/config"
	"log"
	"fmt"
	"runtime"
)
var (
	configFile = flag.String("configfile", "config.ini", "配置文件")
)
var cfg = make(map[string]string)

func loadCfg(){
	runtime.GOMAXPROCS(runtime.NumCPU())
	flag.Parse()
	//set config file std
	cfgrs, err := config.ReadDefault(*configFile)
	if err != nil {
		log.Fatalf("Fail to find", *configFile, err)
	}
	//set config file std End

	//Initialized topic from the configuration
	if cfgrs.HasSection("topicArr") {
		section, err := cfgrs.SectionOptions("topicArr")
		if err == nil {
			for _, v := range section {
				options, err := cfgrs.String("topicArr", v)
				if err == nil {
					cfg[v] = options
				}
			}
		}
	}
	//Initialized topic from the configuration END
	fmt.Println(cfg)
}
func main() {
	m := martini.Classic()
	store := sessions.NewCookieStore([]byte("secret123"))
	m.Use(sessions.Sessions("helpsystem_session", store))
	m.Use(render.Renderer())
	flag.Parse()
	//数据库
//	db, err := sql.Open("adodb", "Provider=SQLNCLI11;DataTypeCompatibility=80;Server=127.0.0.1;UID=sa;PWD=11111111;Database=helpsystem;")
	db, err := sql.Open("odbc","DSN=mssql;DATABASE=helpsystem;UID=sa;PWD=11111111")
	erutil.CheckErr(err)
	db.SetMaxOpenConns(100)
	m.Map(db)
	m.Any("/login", auth.Login)
	m.Get("/logout", auth.Logout)
	m.Get("/", index)
	m.Get("/cats/:catid", helpmaker.Cats)
	m.Get("/pages/:id", helpmaker.Pages)

	//静态内容
	m.Use(martini.Static("static"))
	//需要权限的内容

	m.Group("/admin", admin.Router , auth.Auth)
	m.Run()
}


func index(db *sql.DB , r render.Render, req *http.Request) {
	ret := make(map[string]interface{})
	catid := req.FormValue("catid")
	ret["cats"] = helpmaker.GetCats(catid, db)
	r.HTML(200, "index", ret)
}
