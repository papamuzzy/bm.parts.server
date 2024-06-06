package dbok

import (
	"bm.parts.server/config"
	"database/sql"
	"fmt"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

func NewClient() (*sql.DB, error) {
	//client, err := sql.Open("mysql", "ok_prod:Vh057|CS~e<M@tcp(84.247.137.2:3306)/ok_prod")
	client, err := sql.Open(config.OcDriver, config.OcUser+":"+config.OcPassw+"@tcp("+config.OcServer+":"+config.OcPort+")/"+config.OcDataBase)
	if err != nil {
		fmt.Println(err)
	}
	// See "Important settings" section.
	client.SetConnMaxLifetime(time.Minute * 3)
	client.SetMaxOpenConns(100)
	client.SetMaxIdleConns(100)

	return client, err
}
