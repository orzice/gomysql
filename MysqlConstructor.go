// +----------------------------------------------------------------------
// | Author: Orzice(小涛)
// +----------------------------------------------------------------------
// | 联系我: https://i.orzice.com
// +----------------------------------------------------------------------
// | gitee: https://gitee.com/orzice
// +----------------------------------------------------------------------
// | github: https://github.com/orzice
// +----------------------------------------------------------------------
// | DateTime: 2021-09-29 15:43:14
// +----------------------------------------------------------------------
package extend

import (
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
	"os"
	"strconv"
	"strings"
)

type DB_Mysql struct {
	Connection *sql.DB
	Prefix string	
	Options map[string]interface{}	
	Bind interface{}	
	Config map[string]string 
	Sql_num string 
}

func (u *DB_Mysql) init() {
	u.Options = make(map[string]interface{})
	u.Options["where"] = make(map[int]interface{})
	u.Options["order"] = make(map[int]interface{})
	u.Options["set"] = make(map[int]interface{})
	u.Options["limit"] = ""
	u.Options["field"] = ""
	u.Options["table"] = ""
	u.Options["exp"] = ""
}

func (u *DB_Mysql) Construct() *DB_Mysql {
	u.init()

	S_config := map[string]string{
		// Mysql
		"type": "mysql", 
		"prefix": "",
		"hostname": "localhost",
		"database": "",
		"username": "root",
		"password": "root",
		"hostport": "3306",
		"charset": "utf8",
	}
	u.Config = S_config
	u.connect()
	u.Prefix = u.Config["prefix"]
	return u
}

func (u *DB_Mysql) connect() bool{
	
	db, err := sql.Open(u.Config["type"], u.Config["username"]+":"+u.Config["password"]+"@("+u.Config["hostname"]+":"+u.Config["hostport"]+")/"+u.Config["database"]+"?charset="+u.Config["charset"])
	
	err = db.Ping()
	if err != nil{
		config.TOOL_log(err)
		os.Exit(10001)
		return false
	}
	u.Connection = db
	return true
}
func (u *DB_Mysql) Close()  {
	u.init()
	u.Connection.Close()
}
func (u *DB_Mysql) Getconnect() *sql.DB {
	return u.Connection
}

func (u *DB_Mysql) Table(table string) *DB_Mysql {
	//初始化
	u.init()

	u.Options["table"] =  u.Prefix + table

	return u
}

func (u *DB_Mysql) Field(table string) *DB_Mysql {
	u.Options["field"] =  table

	return u
}

func (u *DB_Mysql) EXP(table string,data string) *DB_Mysql {
	ls := table + "("
	if data == "" {
		ls += "*"
	}else{
		ls +=  "`" + data + "`"
	}
	ls += ")"
	u.Options["exp"] =  ls

	return u
}
func (u *DB_Mysql) Tables(table string) *DB_Mysql {
	u.Options["table"] =  table

	return u
}
func (u *DB_Mysql) Where(field string, op string, condition string) *DB_Mysql {
	u.parseWhereExp("AND", field, op, condition);

	return u
}

func (u *DB_Mysql) WhereOr(field string, op string, condition string) *DB_Mysql {
	u.parseWhereExp("OR", field, op, condition);

	return u
}

func (u *DB_Mysql) Set(field string, order string)  *DB_Mysql {

	u.Options["set"].(map[int]interface{})[len(u.Options["set"].(map[int]interface{}))] =  map[string]string{
		"field" :field,
		"order" :order,
	}

	return u
}
func (u *DB_Mysql) Limit(offset int64, length int64) *DB_Mysql {
	sqll  := strconv.FormatInt(offset, 10)

	if length != 0 {
		sqll += ","+strconv.FormatInt(length, 10)
	}
	u.Options["limit"] = sqll

	return u
}
func (u *DB_Mysql) Order(field string, order string)  *DB_Mysql {

	u.Options["order"].(map[int]interface{})[len(u.Options["order"].(map[int]interface{}))] =  map[string]string{
		"field" :field,
		"order" :order,
	}

	return u

}

func (u *DB_Mysql) Get() map[int]map[string]string {
	sql := u.Constructor("select",u.Options,"")

	rows2, err := u.Connection.Query(sql);
	if err != nil{
		return map[int]map[string]string{}
	}
	cols, err := rows2.Columns();
	if err != nil{
		return map[int]map[string]string{}
	}
	vals := make([][]byte, len(cols));
	scans := make([]interface{}, len(cols));
	for k, _ := range vals {
		scans[k] = &vals[k];
	}
	i := 0;
	result := make(map[int]map[string]string);
	for rows2.Next() {
		rows2.Scan(scans...);
		row := make(map[string]string);
		for k, v := range vals {
			key := cols[k];
			row[key] = string(v);
		}
		result[i] = row;
		i++;
	}
	u.init()
	return result
}


func (u *DB_Mysql) Insert(data map[string]string) (int64) {
	value := make(map[int]interface{})

	for k, v := range data {
		value[len(value)] = map[string]string{
			"name" : k,
			"value" : v,
		}
	}
	sql := u.Constructor("insert",u.Options,value)

	s,err := u.Connection.Exec(sql)
	if err == nil {
		u.init()
		return  0
	}
	t,err := s.LastInsertId()
	if err != nil {
		t = 0
	}
	u.init()

	return t
}
func (u *DB_Mysql) Update(data map[string]string)  {
	value := make(map[int]interface{})

	for k, v := range data {
		value[len(value)] = map[string]string{
			"name" : k,
			"value" : v,
		}
	}
	sql := u.Constructor("update",u.Options,value)
	u.Connection.Exec(sql)

	u.init()
}
func (u *DB_Mysql) Delete() {
	sql := u.Constructor("delete",u.Options,"")
	u.Connection.Exec(sql)

	u.init()
}


func (u *DB_Mysql) parseWhereExp(logic string, field string, op string, condition string) {

		u.Options["where"].(map[int]interface{})[len(u.Options["where"].(map[int]interface{}))] = map[string]string{
			"logic" : logic,
			"field" : field,
			"op" : op,
			"condition" : condition,
		}
}


var selectSql string = "SELECT%DISTINCT% %FIELD% FROM %TABLE%%FORCE%%JOIN%%WHERE%%GROUP%%HAVING%%UNION%%ORDER%%LIMIT%%LOCK%%COMMENT%"
var insertSql string = "%INSERT% INTO %TABLE% (%FIELD%) VALUES (%DATA%) %COMMENT%"
var updateSql string = "UPDATE %TABLE% SET %SET% %JOIN% %WHERE% %ORDER%%LIMIT% %LOCK%%COMMENT%"
var deleteSql string = "DELETE FROM %TABLE% %USING% %JOIN% %WHERE% %ORDER%%LIMIT% %LOCK%%COMMENT%"

var exp = map[string]string{
	"eq" : "=",
	"neq" : "<>",
	"like" : ">",
	"egt" : ">=",
}
func (u *DB_Mysql) Constructor(types string , options  map[string]interface{}, data interface{}) string {

	switch types {
		case "select":
			return u.const_select(options)
		case "insert":
			return u.const_insert(data.(map[int]interface{}),options)
		case "update":
			return u.const_update(data.(map[int]interface{}),options)
		case "delete":
			return u.const_delete(options)
		default:
			return ""
	}
	return ""
}

func (u *DB_Mysql) const_select(options  map[string]interface{}) string{
	sqls := selectSql
	//sqls = strings.Replace(sqls, "%TABLE%", "", -1)
	sqls = strings.Replace(sqls, "%JOIN%", "", -1)
	sqls = strings.Replace(sqls, "%GROUP%", "", -1)
	sqls = strings.Replace(sqls, "%HAVING%", "", -1)
	sqls = strings.Replace(sqls, "%UNION%", "", -1)
	sqls = strings.Replace(sqls, "%LOCK%", "", -1)
	sqls = strings.Replace(sqls, "%COMMENT%", "", -1)
	sqls = strings.Replace(sqls, "%FORCE%", "", -1)
	sqls = strings.Replace(sqls, "%DISTINCT%", "", -1)

	sqls = strings.Replace(sqls, "%FIELD%", u.parseField(options["field"].(string),options["exp"].(string)), -1)
	sqls = strings.Replace(sqls, "%TABLE%", u.parseTable(options["table"].(string)), -1)
	sqls = strings.Replace(sqls, "%WHERE%", u.parseWhere(options["where"].(map[int]interface{})), -1)
	sqls = strings.Replace(sqls, "%ORDER%", u.parseOrder(options["order"].(map[int]interface{})), -1)
	sqls = strings.Replace(sqls, "%LIMIT%", u.parselimit(options["limit"].(string)), -1)

	return sqls
}

func (u *DB_Mysql) const_insert(data map[int]interface{} ,options  map[string]interface{}) string{
	sqls := insertSql
	fields := ""
	values := ""
	for i:=0;i< len(data) ; i++ {
		if i != 0{
			fields += " , "
			values += " , "
		}
		fields += "`"+data[i].(map[string]string)["name"]+"`"
		values += "'"+data[i].(map[string]string)["value"]+"'"
	}

	//sqls = strings.Replace(sqls, "%TABLE%", "", -1)
	sqls = strings.Replace(sqls, "%INSERT%", "INSERT", -1)
	sqls = strings.Replace(sqls, "%FIELD%", fields, -1)
	sqls = strings.Replace(sqls, "%DATA%", values, -1)
	sqls = strings.Replace(sqls, "%COMMENT%", "", -1)//注释

	sqls = strings.Replace(sqls, "%TABLE%", u.parseTable(options["table"].(string)), -1)

	return sqls
}


func (u *DB_Mysql) const_update(data map[int]interface{} ,options  map[string]interface{}) string{
	sqls := updateSql

	SET := ""
	for i:=0;i< len(data) ; i++ {
		if i != 0{
			SET += " , "
		}
		SET += data[i].(map[string]string)["name"] + " = '"+data[i].(map[string]string)["value"]+"'"
	}

	//sqls = strings.Replace(sqls, "%SET%", "", -1)
	sqls = strings.Replace(sqls, "%JOIN%", SET, -1)
	sqls = strings.Replace(sqls, "%LOCK%", "", -1)
	sqls = strings.Replace(sqls, "%COMMENT%", "", -1)

	sqls = strings.Replace(sqls, "%SET%", u.parseSet(options["set"].(map[int]interface{})), -1)
	sqls = strings.Replace(sqls, "%TABLE%", u.parseTable(options["table"].(string)), -1)
	sqls = strings.Replace(sqls, "%WHERE%", u.parseWhere(options["where"].(map[int]interface{})), -1)
	sqls = strings.Replace(sqls, "%ORDER%", u.parseOrder(options["order"].(map[int]interface{})), -1)
	sqls = strings.Replace(sqls, "%LIMIT%", u.parselimit(options["limit"].(string)), -1)

	return sqls
}

func (u *DB_Mysql) const_delete(options  map[string]interface{}) string{
	sqls := deleteSql

	sqls = strings.Replace(sqls, "%USING%", "", -1)
	sqls = strings.Replace(sqls, "%COMMENT%", "", -1)
	sqls = strings.Replace(sqls, "%JOIN%", "", -1)
	sqls = strings.Replace(sqls, "%LOCK%", "", -1)

	sqls = strings.Replace(sqls, "%TABLE%", u.parseTable(options["table"].(string)), -1)
	sqls = strings.Replace(sqls, "%WHERE%", u.parseWhere(options["where"].(map[int]interface{})), -1)
	sqls = strings.Replace(sqls, "%ORDER%", u.parseOrder(options["order"].(map[int]interface{})), -1)
	sqls = strings.Replace(sqls, "%LIMIT%", u.parselimit(options["limit"].(string)), -1)

	return sqls
}


func (u *DB_Mysql) parseWhere(options map[int]interface{}) string{
	whereStr := ""

	for i := 0; i < len(options) ; i++  {
		lsdata := options[i].(map[string]string);
		if i==0 {
			whereStr = " WHERE "
		}else{
			whereStr +=  lsdata["logic"]
		}
		whereStr += " `" + lsdata["field"] + "` " + lsdata["op"] + " '" + lsdata["condition"] + "' "
	}
	return  whereStr
}

func (u *DB_Mysql) parseSet(options map[int]interface{}) string{
	whereStr := ""

	for i := 0; i < len(options) ; i++  {
		lsdata := options[i].(map[string]string);

		whereStr += lsdata["field"] + "=" + lsdata["order"]

		if i +1 != len(options) && len(options) > 1 {
			whereStr +=  ","
		}
	}
	return  whereStr
}


func (u *DB_Mysql) parseOrder(options map[int]interface{}) string{
	whereStr := ""

	for i := 0; i < len(options) ; i++  {
		if i==0 {
			whereStr = "  ORDER BY "
		}
		lsdata := options[i].(map[string]string);
		whereStr += "`" + lsdata["field"] + "` "
		if lsdata["order"] != ""{
			whereStr += "DESC"
		}

		if i +1 != len(options) && len(options) > 1 {
			whereStr +=  ","
		}
	}
	return  whereStr
}
func (u *DB_Mysql) parselimit(data string) string{
	if data == "" {
		return ""
	}
	return " LIMIT " + data
}

func (u *DB_Mysql) parseTable(data string) string{
	if data == "" {
		return "*"
	}

	return "`" + data + "`"
}

func (u *DB_Mysql) parseField(data string,exp string) string{
	ls := ""
	if exp != "" {
		return exp
	}else{
		if data == "" {
			return "*"
		}
		s := strings.Split(data, ",")
		for i:=0;i < len(s)  ;i++  {
			ls += "`"+s[i] + "`"

			if i +1 != len(s) && len(s) > 1 {
				ls +=  ","
			}
		}
	}
	return ls
}

