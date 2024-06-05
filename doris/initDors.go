package doris

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
)

var (
	// 定义一个全局对象db
	db *sql.DB
	//连接Doris的用户名
	userName string = "root"
	//连接Doris的密码
	password string = ""
	//连接Doris的地址
	ipAddress string = "192.168.152.128"
	//连接Doris的端口号,默认是9030
	port int = 9030
	//连接Doris的具体数据库名称
	dbName string = "demo"
)

func InitDB() (err error) {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s", userName, password, ipAddress, port, dbName)
	//Open打开一个driverName指定的数据库，dataSourceName指定数据源
	//不会校验用户名和密码是否正确，只会对dsn的格式进行检测
	db, err = sql.Open("mysql", dsn)
	//dsn格式不正确的时候会报错
	if err != nil {
		return err
	}
	defer db.Close()
	//尝试与数据库连接，校验dsn是否正确
	err = db.Ping()
	if err != nil {
		fmt.Println("校验失败,err", err)
		return err
	}
	// 设置最大连接数
	db.SetMaxOpenConns(50)
	// 设置最大的空闲连接数
	// db.SetMaxIdleConns(20)
	fmt.Println("连接数据库成功！")
	return nil
}

// QueryRow 查询
func QueryRow() {
	rows, _ := db.Query("select * from t_cn_search where demo MATCH_ANY '1'") //获取所有数据
	var md5 int
	var line string
	for rows.Next() { //循环显示所有的数据
		rows.Scan(&md5, &line)
		fmt.Println(md5, "--", line)
	}
}

// Insert 数据库插入
func Insert() {

}
