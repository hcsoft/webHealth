package main

import (
	_ "code.google.com/p/odbc"
	"database/sql"
	"flag"
	"github.com/go-martini/martini"
	"github.com/hcsoft/webHealth/auth"
	"github.com/martini-contrib/render"
	"github.com/martini-contrib/sessions"
//	_ "github.com/mattn/go-adodb"
	"net/http"
	//	erutil "github.com/hcsoft/webHealth/error"
	"github.com/hcsoft/webHealth/admin"
	"github.com/hcsoft/webHealth/dbutil"
	"github.com/hcsoft/webHealth/helpmaker"
	"github.com/larspensjo/config"
	//	"log"
	//	"fmt"
	//	"runtime"
)

func main() {
	m := martini.Classic()
	store := sessions.NewCookieStore([]byte("secret123"))
	m.Use(sessions.Sessions("webhealth_session", store))
	m.Use(render.Renderer())
	//配置文件
	configFile := flag.String("configfile", "config.ini", "配置文件")
	inicfg, _ := config.ReadDefault(*configFile)
	m.Map(inicfg)
	//数据库
	db := dbutil.GetDB(*inicfg)
	m.Map(db)
	//缓存
	cache := make(map[string]interface{})
	district, _ := inicfg.String("base", "district")
	cache["district"] = auth.GetDistrict(db, district)
	m.Map(cache)

	m.Any("/login", auth.Login)
	m.Get("/logout", auth.Logout)
	m.Get("/", index)
	m.Get("/cats/:catid", helpmaker.Cats)
	m.Get("/pages/:id", helpmaker.Pages)

	//静态内容
	m.Use(martini.Static("static"))
	//需要权限的内容
	m.Group("/admin", admin.Router, auth.Auth)
//	m.Run()
	m.RunOnAddr(":3333")
}

func index(db *sql.DB, r render.Render, req *http.Request, inicfg *config.Config) {
	ret := make(map[string]interface{})
	//	catid := req.FormValue("catid")
	//	ret["cats"] = helpmaker.GetCats(catid, db)
	r.HTML(200, "index", ret)
}
