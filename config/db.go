package config

import (
	"context"
	"database/sql"
	"log"
	"os"

	"api/helper/atdb"

	"go.mau.fi/whatsmeow/store/sqlstore"
)

var MongoString string = os.Getenv("MONGOSTRING")

var DBUlbimongoinfo = atdb.DBInfo{
	DBString: MongoString,
	DBName:   "waapi",
}

var Mongoconn, ErrorConnect = atdb.MongoConnect(DBUlbimongoinfo)

var Postgrestring = os.Getenv("PGSTRING")

var ContainerDB *sqlstore.Container

func init() {
	db, err := sql.Open("postgres", Postgrestring)
	if err != nil {
		log.Fatal(err)
	}
	ContainerDB = sqlstore.NewWithDB(db, "postgres", nil)
	err = ContainerDB.Upgrade(context.TODO())
	if err != nil {
		log.Fatal(err)
	}
}
