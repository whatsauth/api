package config

import (
	"os"

	"github.com/aiteung/atdb"
)

var MongoString string = os.Getenv("MONGOSTRING")

var DBUlbimongoinfo = atdb.DBInfo{
	DBString: MongoString,
	DBName:   "apiwa",
}

var Mongoconn = atdb.MongoConnect(DBUlbimongoinfo)
