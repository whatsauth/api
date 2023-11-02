package helper

import (
	"context"
	"fmt"

	_ "github.com/mattn/go-sqlite3"
	"go.mau.fi/whatsmeow"
	"go.mau.fi/whatsmeow/store"
	"go.mau.fi/whatsmeow/store/sqlstore"
	"go.mau.fi/whatsmeow/types/events"
	waLog "go.mau.fi/whatsmeow/util/log"
)

func EventHandler(evt interface{}) {
	switch v := evt.(type) {
	case *events.Message:
		fmt.Println("Received a message!", v.Message.GetConversation())
	}
}

func GetClient(phonenumber string) (client *whatsmeow.Client) {
	dbLog := waLog.Stdout("Database", "DEBUG", true)
	// Make sure you add appropriate DB connector imports, e.g. github.com/mattn/go-sqlite3 for SQLite
	container, err := sqlstore.New("sqlite3", "file:wa.db?_foreign_keys=on", dbLog)
	if err != nil {
		panic(err)
	}
	// If you want multiple sessions, remember their JIDs and use .GetDevice(jid) or .GetAllDevices() instead.
	//deviceStore, err := container.GetFirstDevice()
	deviceStores, err := container.GetAllDevices()
	//deviceStore, err := container.GetDevice(jid)
	var deviceStore *store.Device
	for _, dv := range deviceStores {
		if dv.ID.User == phonenumber {
			deviceStore = dv
		}
	}
	if err != nil {
		fmt.Println(err)
	}
	if deviceStore == nil {
		fmt.Println("buat device baru")
		deviceStore = container.NewDevice()
	}
	//deviceStore, err := container.GetAllDevices()
	clientLog := waLog.Stdout("Client", "ERROR", true)
	client = whatsmeow.NewClient(deviceStore, clientLog)
	client.AddEventHandler(EventHandler)
	return

}

func Connect(client *whatsmeow.Client, qr chan QRStatus) {
	if client.Store.ID == nil {
		qrChan, _ := client.GetQRChannel(context.Background())
		err := client.Connect()
		if err != nil {
			panic(err)
		}
		// No ID stored, new login
		for evt := range qrChan {
			if evt.Event == "code" {
				// Render the QR code here
				// e.g. qrterminal.GenerateHalfBlock(evt.Code, qrterminal.L, os.Stdout)
				// or just manually `echo 2@... | qrencode -t ansiutf8` in a terminal
				fmt.Println("QR code:", evt.Code)
				qr <- QRStatus{true, evt.Code, evt.Event}
			} else {
				fmt.Println("Login event:", evt.Event)
				qr <- QRStatus{true, evt.Code, evt.Event}
			}
		}
	} else {
		message := "login"
		err := client.Connect()
		if err != nil {
			message = err.Error()
		}
		qr <- QRStatus{false, "", message}
	}

}
