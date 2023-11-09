package config

import (
	"github.com/aiteung/atdb"
	"github.com/whatsauth/wa"
)

var MongoString string = "mongodb+srv://awangga:8uPiRHynbtRuHv6X@potp.x8hnwy3.mongodb.net/waapi"

var DBUlbimongoinfo = atdb.DBInfo{
	DBString: MongoString,
	DBName:   "waapi",
}

var Mongoconn = atdb.MongoConnect(DBUlbimongoinfo)

var Postgrestring = "postgres://obruyswq:ZPHsdZ9LYSujKDoHEIehA5uJ3LYkDbv0@satao.db.elephantsql.com/obruyswq"

var ContainerDB, _ = wa.CreateContainerDB(Postgrestring)
