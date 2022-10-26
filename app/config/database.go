package config

import (
	"database/sql"
	errorhandler "golang-depedency-injection/app/error_handler"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/spf13/viper"
)

// var DB *sql.DB

func NewDB() *sql.DB {
	// logLevel := logger.Info
	viper.GetString("ENVIRONMENT")

	consString := "root@tcp(localhost:3306)/golang_test"
	// dsn := fmt.Sprintf("%s@tcp(%s:%s)/%s?sslmode=disable TimeZone=Asia/Jakarta",
	// 	viper.GetString("DB_USER"),
	// 	viper.GetString("DB_HOST"),
	// 	viper.GetString("DB_PORT"),
	// 	viper.GetString("DB_Name"),
	// )
	db, err := sql.Open("mysql", consString)

	errorhandler.PanicIfNeeded(err)
	// fmt.Println("Connecting to database...")
	// defer db.Close()

	db.SetMaxIdleConns(5)
	db.SetMaxOpenConns(20)
	db.SetConnMaxLifetime(60 * time.Minute)
	db.SetConnMaxIdleTime(10 * time.Minute)

	err = db.Ping()
	if err != nil {
		panic(err.Error()) // proper error handling instead of panic in your app
	}

	return db
}
