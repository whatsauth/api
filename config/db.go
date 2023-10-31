package config

import (
	"os"

	"github.com/aiteung/atdb"
)

var IteungIPAddress string = os.Getenv("ITEUNGBEV1")

var MongoString string = os.Getenv("MONGOSTRING")

var MariaStringAkademik string = os.Getenv("MARIASTRINGAKADEMIK")

var DBUlbimariainfo = atdb.DBInfo{
	DBString: MariaStringAkademik,
	DBName:   "xia3fhuwzm5wo0zo",
}

var DBUlbimongoinfo = atdb.DBInfo{
	DBString: MongoString,
	DBName:   "iteung",
}

var Ulbimariaconn = atdb.MariaConnect(DBUlbimariainfo)

var Ulbimongoconn = atdb.MongoConnect(DBUlbimongoinfo)
