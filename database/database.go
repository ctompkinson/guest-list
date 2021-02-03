package database

import (
	"database/sql"
	"fmt"
	"github.com/ctompkinson/guest-list/model"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"os"
)

var (
	db *gorm.DB

	username     = ""
	password     = ""
	address      = ""
	port         = ""
	databaseName = ""
)

// Get returns an instance of the Gorm DB
func Get() *gorm.DB {
	return db
}

// Init initialises the database if it has not already been setup
func Init() error {
	if db != nil {
		return nil
	}

	username = os.Getenv("GUESTLIST_DB_USERNAME")
	if username == "" {
		username = "root"
	}

	password = os.Getenv("GUESTLIST_DB_PASSWORD")
	if password == "" {
		password = "foo"
	}

	databaseName = os.Getenv("GUESTLIST_DB_NAME")
	if databaseName == "" {
		databaseName = "guestlist"
	}

	address = os.Getenv("GUESTLIST_DB_ADDRESS")
	if address == "" {
		address = "localhost"
	}

	port = os.Getenv("GUESTLIST_DB_PORT")
	if port == "" {
		port = "3306"
	}

	Create()
	var err error
	db, err = gorm.Open(mysql.Open(
		fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local", username, password, address, port, databaseName)), &gorm.Config{})
	if err != nil {
		return err
	}

	if err := Migrate(); err != nil {
		return err
	}

	return nil
}

// Migrate sets up database tables using gorm models
func Migrate() error {
	if err := db.AutoMigrate(&model.Table{}); err != nil {
		return err
	}
	if err := db.AutoMigrate(&model.Reservation{}); err != nil {
		return err
	}
	return nil
}

// Create creates the database to be used for Gorm
func Create() {
	db, err := sql.Open("mysql", fmt.Sprintf("%s:%s@tcp(%s:%s)/", username, password, address, port))
	if err != nil {
		panic(err)
	}
	defer db.Close()
	_, err = db.Exec("CREATE DATABASE IF NOT EXISTS guestlist")
	if err != nil {
		panic(err)
	}
}

// ClearAndCreate deletes the database and recreates it for testing purposes
func ClearAndCreate() {
	db, err := sql.Open("mysql", fmt.Sprintf("%s:%s@tcp(%s:%s)/", username, password, address, port))
	if err != nil {
		panic(err)
	}
	defer db.Close()
	_, err = db.Exec("DROP DATABASE IF EXISTS guestlist")
	if err != nil {
		panic(err)
	}
	_, err = db.Exec("CREATE DATABASE IF NOT EXISTS guestlist")
	if err != nil {
		panic(err)
	}

	if err := Migrate(); err != nil {
		panic(err)
	}
}
