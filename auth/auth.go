package auth

import ("database/sql"
	"github.com/go-martini/martini"
	"github.com/martini-contrib/sessions"
	"github.com/martini-contrib/render"
	"net/http"
	"fmt"
	"github.com/hcsoft/webHealth/dbutil"
	"github.com/hcsoft/webHealth/erutil"
	"encoding/json"
)

type JsonRet struct {
	Type string
	Status int
	Msg  string
	Data interface{}
}
func jsonp(obj interface{},callback string)string{
	b,_:=json.Marshal(obj)
	return callback+"("+string(b)+")"
}
func Login(session sessions.Session, db *sql.DB, r render.Render, req *http.Request ,writer http.ResponseWriter) string {
	writer.Header().Set("Content-Type", "text/javascript")
	userid := req.FormValue("userid")
	callback := req.FormValue("callback")
	fmt.Println("userid==" + userid)
	fmt.Println("callback==" + callback)
	password := req.FormValue("password")
	fmt.Println("password==" + password)
	if userid == "" {
//		r.JSON(401, JsonRet{false,"请登录", "abc"})
		return jsonp(JsonRet{"login",401,"请登录", "abc"},callback)
	} else {
		rows, err := db.Query("select * from sam_taxempcode where loginname= ? ", userid)
		erutil.CheckErr(err)
		if rows.Next() {
			values := dbutil.GetOneResult(rows)
			fmt.Println(values["password"]);
			if values["password"] == password {
				session.Set("userid", values["loginname"])
				session.Set("username", values["username"])
//				r.JSON(200, JsonRet{true,"登录成功", values["username"]})
				fmt.Println("登录成功!")
				return jsonp( JsonRet{"login",200,"登录成功", values["username"]},callback)
				//			r.HTML(200, "admin/index", nil)

			}else {
//				r.JSON(404, JsonRet{false,"登录失败!密码错误!", nil})
				return jsonp( JsonRet{"login",401,"登录失败!密码错误!", nil},callback)
			}
		}else {
//			r.JSON(401, JsonRet{false,"登录失败!用户名错误!", nil})
			return jsonp( JsonRet{"login",401,"登录失败!用户名错误!", nil},callback)
		}
	}
}

func Logout(session sessions.Session, r render.Render) {
	session.Delete("userid")
	r.HTML(200, "login", "登出成功")
}

func Auth(session sessions.Session, c martini.Context, r render.Render) {
	v := session.Get("userid")
	if v == nil {
		r.Redirect("/login")
	}else {
		c.Next();
	}
}

