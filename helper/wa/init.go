package wa

import (
	"context"

	"github.com/aiteung/atdb"
	_ "github.com/mattn/go-sqlite3"
	"go.mau.fi/whatsmeow"
	"go.mau.fi/whatsmeow/store"
	"go.mau.fi/whatsmeow/store/sqlstore"
	"go.mau.fi/whatsmeow/types/events"
	waLog "go.mau.fi/whatsmeow/util/log"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

func (mycli *WaClient) register() {
	mycli.eventHandlerID = mycli.WAClient.AddEventHandler(mycli.EventHandler)
}

func (mycli *WaClient) EventHandler(evt interface{}) {
	switch v := evt.(type) {
	case *events.Message:
		go HandlingMessage(&v.Info, v.Message, mycli)
	}
}

func ClientDB(phonenumber string, mongoconn *mongo.Database) (client WaClient) {
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
		deviceStore = container.NewDevice()
	}
	//deviceStore, err := container.GetAllDevices()
	client.PhoneNumber = phonenumber
	client.WAClient = whatsmeow.NewClient(deviceStore, waLog.Stdout("Client", "ERROR", true))
	client.Mongoconn = mongoconn
	client.register()
	return

}

func QRConnect(client WaClient, qr chan QRStatus) {
	if client.WAClient.Store.ID == nil {
		//client.PairPhone(PhoneNumber, true, whatsmeow.PairClientUnknown, "whatsauth.my.id")
		qrChan, _ := client.WAClient.GetQRChannel(context.Background())
		err := client.WAClient.Connect()
		if err != nil {
			panic(err)
		}
		// No ID stored, new login
		for evt := range qrChan {
			if evt.Event == "code" {
				qr <- QRStatus{client.PhoneNumber, true, evt.Code, evt.Event}
			} else {
				qr <- QRStatus{client.PhoneNumber, true, evt.Code, evt.Event}
			}
		}
	} else {
		message := "already login"
		err := client.WAClient.Connect()
		if err != nil {
			message = err.Error()
		}
		qr <- QRStatus{client.PhoneNumber, false, "", message}
	}

}

func PairConnect(client WaClient, qr chan QRStatus) {
	if client.WAClient.Store.ID == nil {
		err := client.WAClient.Connect()
		if err != nil {
			panic(err)
		}
		// No ID stored, new login
		message := "Pair Code Device"
		code, err := client.WAClient.PairPhone(client.PhoneNumber, true, whatsmeow.PairClientUnknown, "Chrome (Mac OS)")
		if err != nil {
			message = err.Error()
		}
		qr <- QRStatus{client.PhoneNumber, true, code, message}
	} else {
		message := "already login"
		if !client.WAClient.IsConnected() {
			message = "Melakukan Koneksi Ulang"
			err := client.WAClient.Connect()
			if err != nil {
				message = err.Error()
			}
		}
		qr <- QRStatus{client.PhoneNumber, false, "", message}
	}

}

func ConnectAllClient(mongoconn *mongo.Database) (clients []*WaClient) {
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
	for _, deviceStore := range deviceStores {
		if deviceStore.ID.User != nosebelumnya {
			client := whatsmeow.NewClient(deviceStore, waLog.Stdout("Client", "ERROR", true))
			//client.AddEventHandler(EventHandler)
			filter := bson.M{"phonenumber": deviceStore.ID.User}
			_, err := atdb.GetOneLatestDoc[User](mongoconn, "user", filter)
			if (client.Store.ID != nil) && (err == nil) {
				var mycli WaClient
				mycli.WAClient = client
				mycli.PhoneNumber = deviceStore.ID.User
				mycli.Mongoconn = mongoconn
				mycli.register()
				client.Connect()
				clients = append(clients, &mycli)

			}
			nosebelumnya = deviceStore.ID.User
		}
	}

	return

}
