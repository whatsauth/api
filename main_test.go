package main

import (
	"api/config"
	"fmt"
	"testing"

	"github.com/whatsauth/wa"

	"github.com/aiteung/atdb"
	"github.com/stretchr/testify/require"
	"github.com/whatsauth/watoken"
)

func TestWatoken(t *testing.T) {
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
}

func TestInsertDB(t *testing.T) {
	var user = wa.User{
		PhoneNumber: "6283131895000",
		WebHook:     wa.WebHook{URL: "https://eov6tgpfbhsve67.m.pipedream.net", Secret: "sajdisandsa"},
		Token:       "v4.public.eyJleHAiOiIyMDIzLTEyLTA0VDA5OjE1OjE0KzA3OjAwIiwiaWF0IjoiMjAyMy0xMS0wNFQwOToxNToxNCswNzowMCIsImlkIjoiNjI4MzEzMTg5NTAwMCIsIm5iZiI6IjIwMjMtMTEtMDRUMDk6MTU6MTQrMDc6MDAifSqR5kBfQhwRfrtrMiOxXNoPP0syIUPpEbtOMqdPOMEfXbOC6boO6NDFKCKKSqjY8WfTcDBXAHtC9N7NHjrvmwM",
	}
	atdb.InsertOneDoc(config.Mongoconn, "user", user)

}
