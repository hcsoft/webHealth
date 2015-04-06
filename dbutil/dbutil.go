package  dbutil

import(
	"database/sql"
	cn "golang.org/x/text/encoding/simplifiedchinese"
	"github.com/larspensjo/config"
	"golang.org/x/text/transform"
	"github.com/hcsoft/webHealth/erutil"
	"io/ioutil"
	"bytes"
	"fmt"
)

func GetDB(inicfg config.Config) *sql.DB{
	dbtype ,_:= inicfg.String("db","dbtype")
	maxConns,_ := inicfg.Int("db","maxconns")
	if dbtype == "odbc"{
		dsn,_ :=inicfg.String(dbtype,"DSN")
		dbname,_ :=inicfg.String(dbtype,"DATABASE")
		uid,_ :=inicfg.String(dbtype,"UID")
		pwd,_ :=inicfg.String(dbtype,"PWD")
		connectionString := fmt.Sprintf("DSN=%s;DATABASE=%s;UID=%s;PWD=%s",dsn,dbname,uid,pwd)
		fmt.Println(connectionString)
		db,err := sql.Open(dbtype,connectionString)
		erutil.CheckErr(err)
		db.SetMaxOpenConns(maxConns)
		return db
	}else if dbtype == "adodb"{
		provider,_ :=inicfg.String(dbtype,"Provider")
		datatype,_ :=inicfg.String(dbtype,"DataTypeCompatibility")
		server,_ :=inicfg.String(dbtype,"Server")
		dbname,_ :=inicfg.String(dbtype,"DATABASE")
		uid,_ :=inicfg.String(dbtype,"UID")
		pwd,_ :=inicfg.String(dbtype,"PWD")

		connectionString := fmt.Sprintf("Provider=%s;DataTypeCompatibility=%s;Server=%s;UID=%s;PWD=%s;Database=%s;",provider,datatype,server,uid,pwd,dbname)
		fmt.Println(connectionString)
		db,err := sql.Open(dbtype,connectionString)
		erutil.CheckErr(err)
		db.SetMaxOpenConns(maxConns)
		return db
	}else{
		panic( ("数据库配置[db]下的dbtype配置类型错误,类型只能为odbc或adodb"))
	}
}

/*获得数据库的map类型的array*/
func GetResultArray(rows *sql.Rows) []map[string]interface{} {
	cols, _ := rows.Columns()
	count := len(cols)
	var ret []map[string]interface{};
	for rows.Next() {
		row := make(map[string]interface{})
		values := make([]interface{}, count)
		valuePtrs := make([]interface{}, count)
		for i, _ := range cols {
			valuePtrs[i] = &values[i]
		}
		rows.Scan(valuePtrs...)
		for i, s := range cols {
			var v interface{}

			val := values[i]

			b, ok := val.([]byte)
			if (ok) {
				data, _ := ioutil.ReadAll(transform.NewReader(bytes.NewReader(b), cn.GB18030.NewDecoder()))
				v= string(data)
			} else {
				v = val
			}
			row[s] = v
		}
		ret = append(ret, row);
	}
	return ret;
}

/*获得数据库的map类型单一结果*/
func GetOneResult(rows *sql.Rows) map[string]interface{} {
	cols, _ := rows.Columns()
	count := len(cols)
	row := make(map[string]interface{})
	values := make([]interface{}, count)
	valuePtrs := make([]interface{}, count)
	for i, _ := range cols {
		valuePtrs[i] = &values[i]
	}
	rows.Scan(valuePtrs...)

	for i, s := range cols {
		var v interface{}

		val := values[i]

		b, ok := val.([]byte)
		if (ok) {
			data, _ := ioutil.ReadAll(transform.NewReader(bytes.NewReader(b), cn.GB18030.NewDecoder()))
			v= string(data)
		} else {
			v = val
		}
		row[s] = v
	}
	return row;
}
