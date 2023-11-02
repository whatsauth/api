package main

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/whatsauth/watoken"
)

func TestWatoken(t *testing.T) {
	privateKey, publicKey := watoken.GenerateKey()
	fmt.Println("privateKey : ", privateKey)
	fmt.Println("publicKey : ", publicKey)
	userid := "awangga"

	tokenstring, err := watoken.Encode(userid, privateKey)
	require.NoError(t, err)
	body, err := watoken.Decode(publicKey, tokenstring)
	fmt.Println("signed : ", tokenstring)
	fmt.Println("isi : ", body)
	require.NoError(t, err)
}
