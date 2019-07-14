package utils

import (
	"AutoTrading/api"
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
)

var dbsource string

func init() {
	//dbsource = api.DBConf.UserName + ":" + api.DBConf.Password + "@tcp(" + api.DBConf.Addr + ":" + api.DBConf.Port + ")/" + api.DBConf.Schema
	dbsource = api.DBConf.UserName + ":" + api.DBConf.Password + "@/" + api.DBConf.Schema
}

func GetMySQL() *sql.DB {
	//[username[:password]@][protocol[(address)]]/dbname[?param1=value1&...&paramN=valueN]
	db, err := sql.Open("mysql", dbsource)
	if err != nil {
		panic(err.Error()) // Just for example purpose. You should use proper error handling instead of panic
	}

	// Open doesn't open a connection. Validate DSN data:
	err = db.Ping()
	if err != nil {
		panic(err.Error()) // proper error handling instead of panic in your app
	}
	return db
}

func Query(sql string) *sql.Row {
	db := GetMySQL()
	defer db.Close()
	// Prepare statement for reading data
	stmtOut, err := db.Prepare(sql)
	if err != nil {
		panic(err.Error()) // proper error handling instead of panic in your app
	}
	defer stmtOut.Close()
	// Query the square-number of 13
	row := stmtOut.QueryRow() // WHERE number = 13
	if err != nil {
		panic(err.Error()) // proper error handling instead of panic in your app
	}
	return row
}

func Exec(sql string) {
	db := GetMySQL()
	defer db.Close()
	// Prepare statement for inserting data
	stmtIns, err := db.Prepare(sql) // ? = placeholder
	if err != nil {
		panic(err.Error()) // proper error handling instead of panic in your app
	}
	defer stmtIns.Close() // Close the statement when we leave main() / the program terminates
	stmtIns.Exec()
}
