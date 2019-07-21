package utils

import (
	"AutoTrading/config"
	"AutoTrading/models"
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
)

var dbsource string
var db *sql.DB
var err error

func init() {
	//dbsource = api.DBConf.UserName + ":" + api.DBConf.Password + "@tcp(" + api.DBConf.Addr + ":" + api.DBConf.Port + ")/" + api.DBConf.Schema
	dbsource = config.DBConf.UserName + ":" + config.DBConf.Password + "@/" + config.DBConf.Schema
	db, err = sql.Open("mysql", dbsource)
	if err != nil {
		panic(err.Error()) // Just for example purpose. You should use proper error handling instead of panic
		err = nil
	}
}

func CreateAccount(id string) *models.Account {
	stmtOut, err := db.Prepare("INSERT INTO `account` (`id`) values (?)")
	if err != nil {
		panic(err.Error())
		err = nil
	}
	defer stmtOut.Close()
	_, err = stmtOut.Exec(id)
	if err != nil {
		panic(err.Error())
		err = nil
	}
	return QueryAccountById(id)
}

func DeleteAccount(id string) error {
	stmtOut, err := db.Prepare("DELETE FROM `account` WHERE id = ?")
	if err != nil {
		panic(err.Error())
		err = nil
	}

	result, err := stmtOut.Exec(id)
	if err != nil {
		panic(err.Error())
		err = nil
	}

	if rowNum, err := result.RowsAffected(); err != nil || rowNum != int64(1) {
		panic("delete error")
		err = nil
	}
	return nil
}

func UpdateAccountById(id string, usdt float64, btc float64, eth float64, bnb float64, eos float64, xrp float64) *models.Account {
	stmtOut, err := db.Prepare("UPDATE `account` SET `usdt`=?, `btc`=?, `eth`=?, `bnb`=?, `eos`=?, `xrp`=? WHERE `id`=?")
	if err != nil {
		panic(err.Error())
	}
	defer stmtOut.Close()

	_, err = stmtOut.Exec(usdt, btc, eth, bnb, eos, xrp, id)
	if err != nil {
		panic(err.Error())
	}

	return QueryAccountById(id)
}

func QueryAccountById(id string) *models.Account {
	stmtOut, err := db.Prepare("SELECT * FROM `account` WHERE id = ?")
	if err != nil {
		panic(err.Error())
		err = nil
	}
	defer stmtOut.Close()
	rows := stmtOut.QueryRow(id)
	account := new(models.Account)
	err = rows.Scan(&account.Id, &account.Usdt, &account.Btc, &account.Eth, &account.Bnb, &account.Eos, &account.Xrp, &account.CreateTime, &account.UpdateTime)
	if err != nil {
		panic(err.Error())
		err = nil
	}
	return account
}

func CloseDB() error {
	return db.Close()
}

func Query(sql string) *sql.Row {
	defer db.Close()
	// Prepare statement for reading data
	stmtOut, err := db.Prepare(sql)
	if err != nil {
		panic(err.Error()) // proper error handling instead of panic in your app
		err = nil
	}
	defer stmtOut.Close()
	// Query the square-number of 13
	row := stmtOut.QueryRow() // WHERE number = 13
	if err != nil {
		panic(err.Error()) // proper error handling instead of panic in your app
		err = nil
	}
	return row
}
