package admin

import (
	"database/sql"
	"fmt"
	"github.com/go-martini/martini"
	"github.com/hcsoft/webHealth/admin/query"
	"github.com/hcsoft/webHealth/auth"
	"github.com/hcsoft/webHealth/dbutil"
	"github.com/hcsoft/webHealth/helpmaker"
	"github.com/hcsoft/webHealth/util"
	"github.com/martini-contrib/render"
	"github.com/martini-contrib/sessions"
	"github.com/satori/go.uuid"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

func Router(router martini.Router) {
	router.Get("/district", func(req *http.Request, session sessions.Session, db *sql.DB, cache map[string]interface{}, writer http.ResponseWriter) string {
		writer.Header().Set("Content-Type", "text/javascript")
		callback := req.FormValue("callback")
		userid := session.Get("userid").(string)
		data, _ := db.Query("select * from sam_taxempcode where loginname= ? ", userid)
		defer data.Close()
		data.Next()
		userinfo := dbutil.GetOneResult(data)
		districtid := userinfo["district_id"].(string)
		fmt.Println(cache)
		return util.Jsonp(auth.GetDistrictById(districtid, cache), callback)
	})
	router.Group("/query", query.Router)

	router.Get("", func(r render.Render, session sessions.Session) {
		r.HTML(200, "admin/index", session.Get("username"))
	})
	router.Get("/index", func(r render.Render, session sessions.Session) {
		r.HTML(200, "admin/index", session.Get("username"))
	})
	router.Get("/EditPages/:id", helpmaker.EditPages)
	router.Get("/helpmanager", func(r render.Render, session sessions.Session, db *sql.DB) {
		r.HTML(200, "admin/helpmanager", helpmaker.GetCats("0", db))
	})
	router.Get("/helpcatsave/:id/:parentid/:ord/:name", func(r render.Render, params martini.Params, db *sql.DB) string {
		id := params["id"]
		parentid := params["parentid"]
		name := params["name"]
		ord := params["ord"]
		rows, err := db.Query("select * from help_cat where id= ? ", id)
		defer rows.Close()
		util.CheckErr(err)
		if rows.Next() {
			_, err := db.Exec("update help_cat set parentid=? ,name=? ,ord=? where id= ? ", parentid, name, ord, id)
			util.CheckErr(err)
			return "保存成功"
		} else {
			_, err := db.Exec("insert into help_cat (id,name,parentid,ord)values(?,?,?,?) ", id, name, parentid, ord)
			util.CheckErr(err)
			return "保存成功"
		}
		return "保存失败"
	})
	router.Get("/helpcatdel/:id", func(r render.Render, params martini.Params, db *sql.DB) string {
		id := params["id"]
		_, err := db.Exec("delete help_cat  where id = ? ", id)
		util.CheckErr(err)
		return "删除成功"
	})

	router.Post("/pic/delete", func(r render.Render, params martini.Params, db *sql.DB, req *http.Request) {
		req.ParseForm()
		ret := map[string]string{"msg": "删除成功"}
		r.JSON(200, ret)
	})
	router.Post("/pic/upload", func(r render.Render, params martini.Params, db *sql.DB, req *http.Request, w http.ResponseWriter) {
		dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
		util.CheckErr(err)
		err = req.ParseMultipartForm(100000)
		util.CheckErr(err)

		//get a ref to the parsed multipart form
		m := req.MultipartForm

		//get the *fileheaders
		files := m.File["file_data"]
		filenames := []string{}
		for i, _ := range files {
			//for each fileheader, get a handle to the actual file
			file, err := files[i].Open()
			defer file.Close()
			util.CheckErr(err)
			//create destination file making sure the path is writeable.
			ext := filepath.Ext(files[i].Filename)
			newfilename := uuid.NewV4().String() + ext
			if _, err := os.Stat(newfilename); err == nil {
				newfilename = uuid.NewV4().String() + ext
			}
			fmt.Println(dir + "/static/upload/" + newfilename)
			dst, err := os.Create(dir + "/static/upload/" + newfilename)
			filenames = append(filenames, "/upload/"+newfilename)
			defer dst.Close()
			util.CheckErr(err)
			if _, err := io.Copy(dst, file); err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
		}
		ret := map[string]string{"urls": strings.Join(filenames, ",")}
		r.JSON(200, ret)
	})
	router.Post("/page/save/:id", func(r render.Render, db *sql.DB, req *http.Request, params martini.Params) {
		req.ParseForm()
		fmt.Println(req.Form["urls[]"])
		id := params["id"]
		_, err := db.Exec("delete help_pages  where catid = ? ", id)
		util.CheckErr(err)
		fmt.Println(len(req.Form["urls[]"]))
		for key, value := range req.Form["urls[]"] {
			strs := strings.Split(value, ",")
			strarry := []string{}
			for _, str := range strs {
				if strings.TrimSpace(str) != "" {
					strarry = append(strarry, strings.TrimSpace(str))
				}
			}
			fmt.Println(value)
			fmt.Println(id, key, strings.Join(strarry, ","))
			_, err := db.Exec("insert into  help_pages (catid,idx,url) values(?,?,?) ", id, key+1, strings.Join(strarry, ","))
			util.CheckErr(err)
		}
		ret := map[string]string{"msg": "保存成功"}
		r.JSON(200, ret)
	})
}
