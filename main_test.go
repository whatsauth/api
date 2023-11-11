package main

import (
	"api/config"
	"fmt"
	"log"
	"testing"

	"github.com/whatsauth/wa"
	"go.mau.fi/whatsmeow"
	"go.mau.fi/whatsmeow/store"
	"go.mau.fi/whatsmeow/types"
	waLog "go.mau.fi/whatsmeow/util/log"
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
}

func TestNewClient(t *testing.T) {
	WAClient := whatsmeow.NewClient(deviceStore, waLog.Stdout("Client", "ERROR", true))
	var txt wa.TextMessage
	txt.IsGroup = false
	txt.Messages = "main test func gol"
	txt.To = "6281312000300"
	resp, err := wa.SendTextMessage(txt, WAClient)
	log.Println(resp)
	log.Println(err)

}
