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

	tokenstring := "v4.public.eyJleHAiOiIyMDIzLTEyLTEwVDIxOjI0OjEwKzA3OjAwIiwiaWF0IjoiMjAyMy0xMS0xMFQyMToyNDoxMCswNzowMCIsImlkIjoiNjI4ODEwODI3MjYzNzkiLCJuYmYiOiIyMDIzLTExLTEwVDIxOjI0OjEwKzA3OjAwIn3810IYPVFjPsoGFa2XVkzqcEyVoLmeZKp0UvwjWdmsKDFPaERJY2OzhMl6QCU33HMUOPjHBbKwDIOQmZN4hFoP"
	body, err := watoken.Decode(config.PublicKey, tokenstring)
	fmt.Println("signed : ", tokenstring)
	fmt.Println("isi : ", body.Id)
	require.NoError(t, err)
	dvcs := config.ContainerDB.NewDevice()
	fmt.Println(dvcs)
	//waclient, err := wa.GetWaClient(body.Id, config.Client, config.Mongoconn, config.ContainerDB)
	//qr := make(chan wa.QRStatus)
	//wa.PairConnect(waclient, qr)
	//resp := <-qr
	//fmt.Println(resp)
}

/*
var phonenumber = "6287752000300"
var deviceid uint16
var deviceStore *store.Device

func TestGetAllDevice(t *testing.T) {
	deviceStores, err := config.ContainerDB.GetAllDevices()
	var deviceStore *store.Device
	for _, dv := range deviceStores {
		if dv.ID.User == phonenumber {
			deviceStore = dv
			deviceid = deviceStore.ID.Device
			fmt.Println("device id:", deviceStore.ID.Device)
		}
	}
	//deviceStore, err := config.ContainerDB.GetDevice(types.JID{User: "6287752000300", Server: "s.whatsapp.net"})
	fmt.Println(deviceStore)
	log.Println(err)
	//clientLog := waLog.Stdout("Client", "DEBUG", true)
	//client := whatsmeow.NewClient(deviceStore, waLog.Stdout("Client", "ERROR", true))
}
func TestGetDevice(t *testing.T) {
	deviceStore, err := config.ContainerDB.GetDevice(types.JID{User: phonenumber, Device: deviceid, Server: "s.whatsapp.net"})
	fmt.Println(deviceStore)
	log.Println(err)
	//clientLog := waLog.Stdout("Client", "DEBUG", true)
	//client := whatsmeow.NewClient(deviceStore, waLog.Stdout("Client", "ERROR", true))
} */

/* func TestNewClient(t *testing.T) {
	WAClient := whatsmeow.NewClient(deviceStore, waLog.Stdout("Client", "ERROR", true))
	var txt wa.TextMessage
	txt.IsGroup = false
	txt.Messages = "main test func gol"
	txt.To = "6281312000300"
	resp, err := wa.SendTextMessage(txt, WAClient)
	log.Println(resp)
	log.Println(err)

} */
