package query

import (
	"database/sql"
	"fmt"
	"github.com/go-martini/martini"
	"github.com/hcsoft/webHealth/dbutil"
	"github.com/hcsoft/webHealth/util"
	"github.com/martini-contrib/sessions"
	"net/http"
	"strings"
)

func Router(router martini.Router) {
	router.Get("/file", func(req *http.Request, session sessions.Session, db *sql.DB, cache map[string]interface{}, writer http.ResponseWriter) string {
		writer.Header().Set("Content-Type", "text/javascript")
		district := req.FormValue("district")
		callback := req.FormValue("callback")
		querystring := req.FormValue("querystring")
		querytype := req.FormValue("querytype")
		fmt.Println(district, querytype, querystring)
		sql := "select top 10 a.fileno,a.name,a.address,a.tel,b.birthday from healthfile a, personalinfo b  where a.fileno = b.fileno and a.fileno like ?+'%'  "
		params := []interface{}{district, querystring}
		switch querytype {
		case "name":
			sql += " and  a.name like '%'+?+'%' "
		case "idnumber":
			sql += " and  b.idnumber = ? "
		case "address":
			sql += " and  a.address like '%'+?+'%'  "
		case "birthday":
			if strings.Index(querystring, "-") > 0 {
				split := strings.Split(querystring, "-")
				params = append(params[:1], split[0], split[1])
				sql += " and b.birthday >=?  and b.birthday < dateadd(day,1,?) "
			} else {
				sql += " and  b.birthday = ? "
			}
		}
		fmt.Println(sql, params)
		rows, _ := db.Query(sql, params...)
		//select a.fileno,a.name,a.address,a.tel,b.birthday from healthfile a, personalinfo b  where a.fileno = b.fileno and a.districtnumber like ?+'%'   and  a.name like '%'+?+'%'  [530521 123]
		//			sql = "select a.* from healthfile a, personalinfo b  where a.fileno = b.fileno and a.DistrictNumber like ?+'%'   and  a.name = ? "

		//			abc, _ := db.Query("select a.* from healthfile a, personalinfo b  where a.fileno = b.fileno and a.fileno like ?+'%'  and  a.name = ? ","530521",  "zjz" )
		data := dbutil.GetResultArrayLimit(rows, 10)
		defer rows.Close()
		fmt.Println(data)
		return util.Jsonp(data, callback)
	})
}
