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
	if dbVendor == "MongoDB" {
		var (
			dbVendor db.MongoDb
		)
		db.DbInstance = &dbVendor
		db.Dbh = db.DbInstance.Connect("localhost:27017", "demoDb")
	} else {
		log.Fatal("DB Vendor " + dbVendor + " not supported")
	}

	router.AddRoutes()
	log.Fatal(http.ListenAndServe(":8080", nil))
}
