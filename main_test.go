package main

import (
	"api/config"
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/whatsauth/watoken"
)

func TestWatoken(t *testing.T) {
	//privateKey, publicKey := watoken.GenerateKey()
	//fmt.Println("privateKey : ", privateKey)
	//fmt.Println("publicKey : ", publicKey)
	userid := "6287752000300"

	tokenstring, err := watoken.EncodeforHours(userid, config.PrivateKey, 720) //30hari
	require.NoError(t, err)
	body, err := watoken.Decode(config.PublicKey, tokenstring)
	fmt.Println("signed : ", tokenstring)
	fmt.Println("isi : ", body.Id)
	require.NoError(t, err)
}
