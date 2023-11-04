package wa

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

func (mycli *WaClient) register() {
	mycli.eventHandlerID = mycli.WAClient.AddEventHandler(mycli.EventHandler)
}

func (mycli *WaClient) EventHandler(evt interface{}) {
	switch v := evt.(type) {
	case *events.Message:
		go HandlingMessage(&v.Info, v.Message, mycli.WAClient)
		fmt.Println("Received a message!", v.Message.GetConversation())
	}
	// Handle event and access mycli.WAClient
}

func GetClient(phonenumber string) (client *whatsmeow.Client) {
	dbLog := waLog.Stdout("Database", "ERROR", true)
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
	if deviceStore == nil {
		fmt.Println("buat device baru")
		deviceStore = container.NewDevice()
	}
	//deviceStore, err := container.GetAllDevices()
	clientLog := waLog.Stdout("Client", "ERROR", true)
	client = whatsmeow.NewClient(deviceStore, clientLog)
	var mycli WaClient
	mycli.WAClient = client
	mycli.register()
	return

}

func Connect(PhoneNumber string, qr chan QRStatus) {
	client := GetClient(PhoneNumber)
	if client.Store.ID == nil {
		//client.PairPhone(PhoneNumber, true, whatsmeow.PairClientUnknown, "wa.my.id")
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
				qr <- QRStatus{PhoneNumber, true, evt.Code, evt.Event}
			} else {
				fmt.Println("Login event:", evt.Event)
				qr <- QRStatus{PhoneNumber, true, evt.Code, evt.Event}
			}
		}
	} else {
		message := "login"
		err := client.Connect()
		if err != nil {
			message = err.Error()
		}
		qr <- QRStatus{PhoneNumber, false, "", message}
	}

}

func Start(client *whatsmeow.Client) {
	if client.Store.ID != nil {
		err := client.Connect()
		if err != nil {
			fmt.Println(err)
		}
	}

}

func ConnectAllClient() {
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
	nosebelumnya := ""
	for i, deviceStore := range deviceStores {
		if deviceStore.ID.User != nosebelumnya {
			fmt.Printf("%d. %s", i, deviceStore.ID.User)
			clientLog := waLog.Stdout("Client", "ERROR", true)
			client := whatsmeow.NewClient(deviceStore, clientLog)
			var mycli WaClient
			mycli.WAClient = client
			mycli.register()
			//client.AddEventHandler(EventHandler)
			if client.Store.ID != nil {
				err := client.Connect()
				if err != nil {
					fmt.Println(err)
				}

			}
			nosebelumnya = deviceStore.ID.User
		}
	}

	return

}