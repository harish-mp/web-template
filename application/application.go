package application

import (
	"log"
	"net/http"
	"os"

	"github.com/harish-mp/web-template/db"
	"github.com/harish-mp/web-template/router"
)

func AppInit() {

	dbVendor := os.Getenv("APPDB")
	dbhost := os.Getenv("DBHOST")
	dbname := os.Getenv("DBNAME")
	if dbVendor == "MongoDB" {
		var (
			dbVendor db.MongoDb
		)
		db.DbInstance = &dbVendor
		db.Dbh = db.DbInstance.Connect(dbhost, dbname)
	} else {
		log.Fatal("DB Vendor " + dbVendor + " not supported")
	}

	router.AddRoutes()
	log.Fatal(http.ListenAndServe(":8080", nil))
}
