package main

import (
	"api/config"
	"database/sql"
	"testing"

	"github.com/lib/pq"
)

/* func TestWatoken(t *testing.T) {
	//privateKey, publicKey := watoken.GenerateKey()
	//fmt.Println("privateKey : ", privateKey)
	//fmt.Println("publicKey : ", publicKey)
	userid := "6283131895000"

	tokenstring, err := watoken.EncodeforHours(userid, config.PrivateKey, 720) //30hari
	require.NoError(t, err)
	body, err := watoken.Decode(config.PublicKey, tokenstring)
	fmt.Println("signed : ", tokenstring)
	fmt.Println("isi : ", body.Id)
	require.NoError(t, err)
} */

func TestInsertDB(t *testing.T) {
	pgUrl, err := pq.ParseURL(config.Postgrestring)
	db, err := sql.Open("postgres", pgUrl)
	if err != nil {
		panic(err)
	}
	defer db.Close()

	// var user = wa.User{
	// 	PhoneNumber: "62831",
	// 	WebHook:     wa.WebHook{URL: "https://eov6tgpfbhsve67.m.pipedream.net", Secret: "sajdisandsa"},
	// 	Token:       "v4.public.eyJleHAiOiIyMDIzLTEyLTA0VDA5OjE1OjE0KzA3OjAwIiwiaWF0IjoiMjAyMy0xMS0wNFQwOToxNToxNCswNzowMCIsImlkIjoiNjI4MzEzMTg5NTAwMCIsIm5iZiI6IjIwMjMtMTEtMDRUMDk6MTU6MTQrMDc6MDAifSqR5kBfQhwRfrtrMiOxXNoPP0syIUPpEbtOMqdPOMEfXbOC6boO6NDFKCKKSqjY8WfTcDBXAHtC9N7NHjrvmwM",
	// }
	// idinsert := atdb.InsertOneDoc(config.Mongoconn, "user", user)
	// fmt.Println(idinsert)
	// a, err := atdb.GetOneLatestDoc[wa.User](config.Mongoconn, "user", bson.M{"phonenumber": "62831"})
	// if err == nil {
	// 	fmt.Println("ada isinya hapus dulu")
	// }
	// fmt.Println(a)
	// fmt.Println(err)
	// anu := atdb.ReplaceOneDoc(config.Mongoconn, "user", bson.M{"phonenumber": "628310"}, user)
	// fmt.Println(anu.ModifiedCount)
}
