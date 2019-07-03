package util

import (
	"AutoTrading/api"
	"database/sql"
)

//const dbsource = api.DBConf.UserName + ":" + api.DBConf.Password + "@/" + api.DBConf.Addr

func GetMySQL() *sql.DB {
	db, err := sql.Open("mysql", api.DBConf.UserName+":"+api.DBConf.Password+"@/"+api.DBConf.Addr)
	if err != nil {
		panic(err.Error()) // Just for example purpose. You should use proper error handling instead of panic
	}
	defer db.Close()

	// Open doesn't open a connection. Validate DSN data:
	err = db.Ping()
	if err != nil {
		panic(err.Error()) // proper error handling instead of panic in your app
	}
	return db
}

func Select(sql string) {
	db := GetMySQL()
	// Prepare statement for reading data
	stmtOut, err := db.Prepare("SELECT squareNumber FROM squarenum WHERE number = ?")
	if err != nil {
		panic(err.Error()) // proper error handling instead of panic in your app
	}
	defer stmtOut.Close()
}
