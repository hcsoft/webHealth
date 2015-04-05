package  db

import(
	"database/sql"
	cn "golang.org/x/text/encoding/simplifiedchinese"
//	"golang.org/x/text/encoding/unicode"
	"golang.org/x/text/transform"
	"io/ioutil"
	"bytes"
//	"fmt"
)

func GetDB(){

	//ado sql.Open("adodb", "Provider=SQLNCLI11;DataTypeCompatibility=80;Server=127.0.0.1;UID=sa;PWD=11111111;Database=helpsystem;")

	//odbc sql.Open("odbc","DSN=mssql;DATABASE=helpsystem;UID=sa;PWD=11111111")
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
