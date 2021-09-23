package main


import (
	"encoding/json"
	"flag"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	_ "github.com/denisenkom/go-mssqldb"
	"github.com/go-xorm/xorm"
	"strconv"
	"github.com/sijms/go-ora/v2"
	"database/sql"
)

func mysql(username string,password string,port string,host string,protocol string,dbname string,sql string)  {
	var engine *xorm.Engine

	connString := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8",username,password,host,port,dbname)
	engine, err := xorm.NewEngine(protocol, connString)
	if err != nil {
		fmt.Println(err)
		return
	}
	//连接测试
	engine.ShowSQL(true)
	if err := engine.Ping(); err != nil {
		fmt.Println(err)
		return
	}
	defer engine.Close() //延迟关闭数据库
	fmt.Println("mysql数据库链接成功")

	//gsql := "select * from liushui limit 10;"
	gres, gerr := engine.QueryString(sql)
	if gerr != nil{
		panic(gerr)
	}
	//fmt.Println(gres)
	for _,v := range gres{
		//fmt.Printf("%v\n",v)

		childJson, _ := json.Marshal(v)
		childString := string(childJson)
		fmt.Printf("%s\n",childString)
	}
}
func mssql(username string,password string,port string,host string,protocol string,dbname string,sql string)  {
	var engine2 *xorm.Engine
	//server=127.0.0.1;user id=sa;password=123456;database=dbname
	//connString := fmt.Sprintf("server=%s:%s;user id=%s;password=%s;database=%s",host,port,username,password,dbname)
	//connString := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8",username,password,host,port,dbname)
	//engine2, err := xorm.NewEngine(protocol, connString)
	//engine,err = xorm.NewEngine("mssql",fmt.Sprintf("sqlserver://%s:%s@%s:%d?database=%s",Arg.User,Arg.Password,Arg.Host,Arg.Port,Arg.Db))
	//fmt.Sprintf("%s://%s:%s@%s:%s?database=%s",protocol,username,password,host,port,dbname)
	atoi, err2 := strconv.Atoi(port)
	if err2 != nil{
		panic(err2)
	}
	connString := fmt.Sprintf("%s://%s:%s@%s:%d?database=%s",protocol,username,password,host,atoi,dbname)
	fmt.Println(connString)
	engine2,err := xorm.NewEngine("mssql",connString)
	if err != nil {
		fmt.Println(err)
		return
	}
	//连接测试
	engine2.ShowSQL(true)
	if err := engine2.Ping(); err != nil {
		fmt.Println(err)
		return
	}
	defer engine2.Close() //延迟关闭数据库
	fmt.Println("mssql数据库链接成功")
	gres, gerr := engine2.QueryString(sql)
	if gerr != nil{
		panic(gerr)
	}
	//fmt.Println(gres)
	for _,v := range gres{
		//fmt.Printf("%v\n",v)

		childJson, _ := json.Marshal(v)
		childString := string(childJson)
		fmt.Printf("%s\n",childString)
	}



}
func oracle(username string,password string,port int,host string,protocol string,dbname string,sql1 string){
	//databaseURL := go_ora.BuildUrl("192.168.1.9", 1521, "orcl", "system", "123456", nil)
	databaseURL := go_ora.BuildUrl(host, port, dbname, username, password, nil)
	conn, _ := sql.Open(protocol, databaseURL)
	// check for err
	_ = conn.Ping()
	// check for err
	defer func() {
		err := conn.Close()
		// check for err
		if err != nil{
			panic(err)
		}
	}()
	rows, err := conn.Query(sql1)
	defer rows.Close()
	columnTypes, err := rows.ColumnTypes()
	if err != nil {
		fmt.Println(err)
	}
	count := len(columnTypes)
	finalRows := []interface{}{};
	for rows.Next() {
		scanArgs := make([]interface{}, count)
		for i, v := range columnTypes {
			switch v.DatabaseTypeName() {
			case "VARCHAR", "TEXT", "UUID", "TIMESTAMP":
				scanArgs[i] = new(sql.NullString)
				break;
			case "BOOL":
				scanArgs[i] = new(sql.NullBool)
				break;
			case "INT4":
				scanArgs[i] = new(sql.NullInt64)
				break;
			default:
				scanArgs[i] = new(sql.NullString)
			}
		}
		err := rows.Scan(scanArgs...)
		if err != nil {
			fmt.Println(err)
		}
		masterData := map[string]interface{}{}
		for i, v := range columnTypes {
			if z, ok := (scanArgs[i]).(*sql.NullBool); ok  {
				masterData[v.Name()] = z.Bool
				continue;
			}
			if z, ok := (scanArgs[i]).(*sql.NullString); ok  {
				masterData[v.Name()] = z.String
				continue;
			}
			if z, ok := (scanArgs[i]).(*sql.NullInt64); ok  {
				masterData[v.Name()] = z.Int64
				continue;
			}
			if z, ok := (scanArgs[i]).(*sql.NullFloat64); ok  {
				masterData[v.Name()] = z.Float64
				continue;
			}
			if z, ok := (scanArgs[i]).(*sql.NullInt32); ok  {
				masterData[v.Name()] = z.Int32
				continue;
			}
			masterData[v.Name()] = scanArgs[i]
		}
		finalRows = append(finalRows, masterData)
	}
	z, err := json.Marshal(finalRows)
	if err != nil{
		fmt.Println(err)
	}
	fmt.Printf("%s\n",z)

}

var db_protocol = flag.String("type", "mysql,mssql,oracle", "Input Your databases type")
var db_host = flag.String("host", "127.0.0.1", "Input Your host")
var db_port = flag.String("port", "3306", "Input Your port")
var db_user = flag.String("user", "root", "Input Your user")
var db_password = flag.String("password", "root", "Input Your password")
var db_name = flag.String("dbname", "mysql,orcl", "Input Your dbname")
var db_sql = flag.String("sql", "select version()", "Input Your sql")

var cliFlag int
func Init() {
	flag.IntVar(&cliFlag, "flagname", 1234, "Just for demo")
}


func main()  {
	//fmt.Println("sssss")
	Init()
	flag.Parse()
	//fmt.Println(*db_protocol)
	//fmt.Println(*db_host)
	//fmt.Println(*db_port)
	//fmt.Println(*db_user)
	//fmt.Println(*db_password)
	//fmt.Println(*db_name)
	//fmt.Println(*db_sql)
	//(username string,password string,port string,host string,protocol string,dbname string,sql st
	//mysql(*db_user,*db_password,*db_port,*db_host,*db_protocol,*db_name,*db_sql)
	//mssql(*db_user,*db_password,*db_port,*db_host,*db_protocol,*db_name,*db_sql)
	//port, err := strconv.Atoi(*db_port)
	//if err != nil{
	//	fmt.Println("port change failed")
	//}
	//oracle(*db_user,*db_password,port,*db_host,*db_protocol,*db_name,*db_sql)
	switch {
	case *db_protocol == "mysql":
		mysql(*db_user,*db_password,*db_port,*db_host,*db_protocol,*db_name,*db_sql)
	case *db_protocol == "mssql":
		mssql(*db_user,*db_password,*db_port,*db_host,*db_protocol,*db_name,*db_sql)
	case *db_protocol == "oracle":
		port, err := strconv.Atoi(*db_port)
		if err != nil{
			fmt.Println("端口参数字符串转int失败")
		}
		oracle(*db_user,*db_password,port,*db_host,*db_protocol,*db_name,*db_sql)

	default:
		fmt.Println("没有输入要连接的数据库类型,缺乏-type参数")
	}





	//gsql := "select * from liushui limit 10;"
	//mysql(gsql)

}
