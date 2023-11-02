package config

import (
	"github.com/aiteung/atdb"
)

var MongoString string = "mongodb+srv://awangga:8uPiRHynbtRuHv6X@potp.x8hnwy3.mongodb.net/waapi"

var DBUlbimongoinfo = atdb.DBInfo{
	DBString: MongoString,
	DBName:   "waapi",
}

var Mongoconn = atdb.MongoConnect(DBUlbimongoinfo)
