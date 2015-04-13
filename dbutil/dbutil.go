package dbutil

import (
	"database/sql"
	cn "golang.org/x/text/encoding/simplifiedchinese"
	//	"golang.org/x/text/encoding/unicode"
	"bytes"
	"fmt"
	"github.com/hcsoft/webHealth/util"
	"github.com/larspensjo/config"
	"golang.org/x/text/transform"
	"io/ioutil"
	//	"unicode/utf16"
	//	"unicode/utf8"
	//	"encoding/binary"
)

func GetDB(inicfg config.Config) *sql.DB {
	dbtype, _ := inicfg.String("db", "dbtype")
	maxConns, _ := inicfg.Int("db", "maxconns")
	maxidles, _ := inicfg.Int("db", "maxidles")
	if dbtype == "odbc" {
		dsn, _ := inicfg.String(dbtype, "DSN")
		dbname, _ := inicfg.String(dbtype, "DATABASE")
		uid, _ := inicfg.String(dbtype, "UID")
		pwd, _ := inicfg.String(dbtype, "PWD")
		connectionString := fmt.Sprintf("DSN=%s;DATABASE=%s;UID=%s;PWD=%s", dsn, dbname, uid, pwd)
		fmt.Println(connectionString)
		db, err := sql.Open(dbtype, connectionString)
		util.CheckErr(err)
		db.SetMaxOpenConns(maxConns)
		db.SetMaxIdleConns(maxidles);
		return db
	} else if dbtype == "adodb" {
		provider, _ := inicfg.String(dbtype, "Provider")
		datatype, _ := inicfg.String(dbtype, "DataTypeCompatibility")
		server, _ := inicfg.String(dbtype, "Server")
		dbname, _ := inicfg.String(dbtype, "DATABASE")
		uid, _ := inicfg.String(dbtype, "UID")
		pwd, _ := inicfg.String(dbtype, "PWD")

		connectionString := fmt.Sprintf("Provider=%s;DataTypeCompatibility=%s;Server=%s;UID=%s;PWD=%s;Database=%s;", provider, datatype, server, uid, pwd, dbname)
		fmt.Println(connectionString)
		db, err := sql.Open(dbtype, connectionString)
		util.CheckErr(err)
		db.SetMaxOpenConns(maxConns)
		db.SetMaxIdleConns(maxidles);
		return db
	} else {
		panic(("数据库配置[db]下的dbtype配置类型错误,类型只能为odbc或adodb"))
	}
}

/*获得数据库的map类型的array*/

func GetResultArrayLimit(rows *sql.Rows, max int) []map[string]interface{} {
	cols, _ := rows.Columns()
	count := len(cols)
	maxount := 0
	var ret []map[string]interface{}
	for rows.Next() && maxount < max {
		//		fmt.Println("-----------------------")
		maxount++
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
			if ok {
				data, _ := ioutil.ReadAll(transform.NewReader(bytes.NewReader(b), cn.GB18030.NewDecoder()))
				//					data, _ := ioutil.ReadAll(transform.NewReader(bytes.NewReader(b), unicode.UTF16(unicode.LittleEndian,unicode.ExpectBOM).NewDecoder()))
				v = string(data)
			} else {

				v = val
			}

			row[s] = v
		}
		ret = append(ret, row)
	}
	return ret
}

func GetResultArray(rows *sql.Rows) []map[string]interface{} {
	cols, _ := rows.Columns()
	count := len(cols)
	var ret []map[string]interface{}
	for rows.Next() {
		//		fmt.Println("-----------------------")
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
			if ok {
				data, _ := ioutil.ReadAll(transform.NewReader(bytes.NewReader(b), cn.GB18030.NewDecoder()))
				//					data, _ := ioutil.ReadAll(transform.NewReader(bytes.NewReader(b), unicode.UTF16(unicode.LittleEndian,unicode.ExpectBOM).NewDecoder()))
				v = string(data)
			} else {

				v = val
			}

			row[s] = v
		}
		ret = append(ret, row)
	}
	return ret
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
		if ok {
			data, _ := ioutil.ReadAll(transform.NewReader(bytes.NewReader(b), cn.GB18030.NewDecoder()))
			v = string(data)
		} else {
			v = val
		}

		row[s] = v
	}
	return row
}
