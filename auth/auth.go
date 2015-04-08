package auth

import ("database/sql"
	"github.com/go-martini/martini"
	"github.com/martini-contrib/sessions"
	"github.com/martini-contrib/render"
	"net/http"
	"fmt"
	"github.com/hcsoft/webHealth/dbutil"
	"github.com/hcsoft/webHealth/util"

)

type JsonRet struct {
	Type string
	Status int
	Msg  string
	Data interface{}
}

func Login(session sessions.Session, db *sql.DB, r render.Render, req *http.Request ,writer http.ResponseWriter) string {
	writer.Header().Set("Content-Type", "text/javascript")
	userid := req.FormValue("userid")
	callback := req.FormValue("callback")
	password := req.FormValue("password")
	if userid == "" {
		return util.Jsonp(JsonRet{"login",401,"请登录", "abc"},callback)
	} else {
		rows, err := db.Query("select * from sam_taxempcode where loginname= ? ", userid)
		defer rows.Close()
		util.CheckErr(err)
		if rows.Next() {
			values := dbutil.GetOneResult(rows)
			fmt.Println(values);
			if values["password"] == password {
				session.Set("userid", values["loginname"])
				session.Set("username", values["username"])
				fmt.Println("登录成功!")
				return util.Jsonp( JsonRet{"login",200,"登录成功", values["username"]},callback)

			}else {
				return util.Jsonp( JsonRet{"login",401,"登录失败!密码错误!", nil},callback)
			}
		}else {
			return util.Jsonp( JsonRet{"login",401,"登录失败!用户名错误!", nil},callback)
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

func GetDistrictById(id string , cache map[string]interface{}) interface{}{
	var districts map[string]interface{}
	districts = cache["district"].(map[string]interface{})
	count := (len(id)-6) /3
	root := districts
	for i := 1 ; i <count ; i++{
		key := id[:6+i*3]
		child := root["child"].([]map[string]interface{})
		idindex := root["idindex"].( map[string]int)
		root = child[idindex[key]]
	}
//	fmt.Println(root);
	return root
}

func GetDistrict(db *sql.DB,dist string)map[string]interface{}{
	root, _ := db.Query("select * from district where id= ? ",dist)
	defer root.Close()
	root.Next()
	rootdata := dbutil.GetOneResult(root)
	child,idindex := getDistrictChild(db,rootdata["ID"].(string))
	rootdata["child"] = child
	rootdata["idindex"] = idindex
	if len(child)>0{
		rootdata["haschild"] = true
	}else{
		rootdata["haschild"] = false
	}
	return rootdata
}

func getDistrictChild(db *sql.DB,dist string) ([]map[string]interface{} , map[string]int){
	child, _  := db.Query("select * from district where parentid= ? ",dist)
	defer child.Close()
	childdata := dbutil.GetResultArray(child)
	idindex  := make(map[string]int)
	for i,v := range childdata{
		id :=v["ID"].(string)
		idindex[id] = i
		child,childidindex := getDistrictChild(db,id)
		v["child"] = child
		v["idindex"] = childidindex
		if len(child)>0{
			v["haschild"] = true
		}else{
			v["haschild"] = false
		}
	}
	return childdata ,idindex
}
