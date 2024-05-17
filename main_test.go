package main

import (
	"api/config"
	"context"
	"database/sql"
	"log"
	"testing"
	"time"

	"api/helper/wa"

	"github.com/lib/pq"
	waProto "go.mau.fi/whatsmeow/binary/proto"
	"go.mau.fi/whatsmeow/types"
	"google.golang.org/protobuf/proto"
)

func TestDB(t *testing.T) {
	log.Println("Test Koneksi DB")

	pgUrl, err := pq.ParseURL(config.Postgrestring)
	if err != nil {
		panic(err)
	}
	db, err := sql.Open("postgres", pgUrl)
	db.SetConnMaxLifetime(time.Second * 15)
	if err != nil {
		panic(err)
	}
	defer db.Close()
	err = db.Ping()
	if err != nil {
		log.Println(err)
	}
	log.Println("Kirim pesan")
	phonenumber := "6287752000300"
	waclient, err := wa.CreateClientfromContainer(phonenumber, config.Mongoconn, config.ContainerDB)
	log.Println(err)
	var wamsg waProto.Message
	wamsg.Conversation = proto.String("asd sa dsa dsad sa ")
	log.Println("Panggilan send message")
	waclient.WAClient.SendMessage(context.Background(), types.NewJID("6281312000300", "s.whatsapp.net"), &wamsg)
	waclient.WAClient.Disconnect()
	log.Println("Selesai Kirim pesan")

}
