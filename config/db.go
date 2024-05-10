package config

import (
	"database/sql"
	"log"

	"github.com/aiteung/atdb"
	"go.mau.fi/whatsmeow/store/sqlstore"
)

var MongoString string = "mongodb+srv://awangga:8uPiRHynbtRuHv6X@potp.x8hnwy3.mongodb.net/waapi"

var DBUlbimongoinfo = atdb.DBInfo{
	DBString: MongoString,
	DBName:   "waapi",
}

var Mongoconn = atdb.MongoConnect(DBUlbimongoinfo)

// elephantsql 20mb
//var Postgrestring = "postgres://obruyswq:ZPHsdZ9LYSujKDoHEIehA5uJ3LYkDbv0@satao.db.elephantsql.com/obruyswq"

// neon tect 100hours
// var Postgrestring = "postgresql://awangga:z9iNkyTFZOt5@ep-steep-pine-25929021-pooler.ap-southeast-1.aws.neon.tech/whatsauth?sslmode=require"
// fly.io
// var Postgrestring = "postgres://postgres:iMTFz957Ov9eTmh@127.0.0.1/whatsauth?sslmode=disable"
var Postgrestring = "postgres://postgres:iMTFz957Ov9eTmh@whatsauth.flycast/whatsauth?sslmode=disable"

var ContainerDB *sqlstore.Container

func init() {
	db, err := sql.Open("postgres", Postgrestring)
	if err != nil {
		log.Fatal(err)
	}
	ContainerDB = sqlstore.NewWithDB(db, "postgres", nil)
	err = ContainerDB.Upgrade()
	if err != nil {
		log.Fatal(err)
	}
}
