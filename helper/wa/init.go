package wa

import (
	"context"
	"errors"
	"fmt"

	"github.com/aiteung/atdb"
	"github.com/lib/pq"
	_ "github.com/mattn/go-sqlite3"
	"go.mau.fi/whatsmeow"
	"go.mau.fi/whatsmeow/store"
	"go.mau.fi/whatsmeow/store/sqlstore"
	"go.mau.fi/whatsmeow/types"
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

func CreateContainerDB(pgstring string) (container *sqlstore.Container, err error) {
	dbLog := waLog.Stdout("Database", "ERROR", true)
	pgUrl, err := pq.ParseURL(pgstring)
	if err != nil {
		return
	}
	container, err = sqlstore.New("postgres", pgUrl, dbLog)
	if err != nil {
		return
	}
	return
}

func ResetDeviceStore(mongoconn *mongo.Database, client *WaClient, container *sqlstore.Container) (err error) {
	if client.WAClient.Store.ID != nil {
		var id uint16
		id, err = GetDeviceIDFromContainer(client.PhoneNumber, container)
		if err != nil {
			return
		}
		filter := bson.M{"phonenumber": client.PhoneNumber}
		var user User
		user, err = atdb.GetOneLatestDoc[User](mongoconn, "user", filter)
		user.DeviceID = id
		client.WAClient.Store.ID.Device = id
		atdb.ReplaceOneDoc(mongoconn, "user", bson.M{"phonenumber": user.PhoneNumber}, user)

	}
	return
}

func CreateClientfromContainer(phonenumber string, mongoconn *mongo.Database, container *sqlstore.Container) (client *WaClient, err error) {
	var user User
	user, err = atdb.GetOneLatestDoc[User](mongoconn, "user", bson.M{"phonenumber": phonenumber})
	if err != nil {
		user.PhoneNumber = phonenumber
	}
	var deviceStore *store.Device
	if user.DeviceID == 0 {
		var deviceid uint16
		deviceid, err = GetDeviceIDFromContainer(phonenumber, container)
		if err != nil {
			return
		}
		deviceStore, err = container.GetDevice(types.JID{User: user.PhoneNumber, Device: deviceid, Server: "s.whatsapp.net"})
	} else {
		deviceStore, err = container.GetDevice(types.JID{User: user.PhoneNumber, Device: user.DeviceID, Server: "s.whatsapp.net"})
	}
	if deviceStore == nil {
		deviceStore = container.NewDevice()
	}
	var wc WaClient
	wc.PhoneNumber = phonenumber
	wc.WAClient = whatsmeow.NewClient(deviceStore, waLog.Stdout("Client", "ERROR", true))
	wc.Mongoconn = mongoconn
	wc.register()
	ConnectClient(wc.WAClient)
	if deviceStore.ID != nil { //tanda belum terkoneksi
		user.DeviceID = deviceStore.ID.Device
	} else {
		user.DeviceID = 0
	}
	atdb.ReplaceOneDoc(mongoconn, "user", bson.M{"phonenumber": phonenumber}, user)
	client = &wc
	return
}

func GetDeviceIDFromContainer(phonenumber string, container *sqlstore.Container) (deviceid uint16, err error) {
	deviceStores, err := container.GetAllDevices()
	fmt.Println(err)
	for _, dv := range deviceStores {
		if dv.ID.User == phonenumber {
			deviceid = dv.ID.Device
		}
	}
	return
}

func GetDeviceStoreFromContainer(phonenumber string, container *sqlstore.Container) (device *store.Device, err error) {
	deviceStores, err := container.GetAllDevices()
	fmt.Println(err)
	for _, dv := range deviceStores {
		if dv.ID.User == phonenumber {
			device = dv
		}
	}
	return
}

func QRConnect(client *WaClient, qr chan QRStatus) {
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

func PairConnect(client *WaClient, qr chan QRStatus) {
	if client.WAClient.Store.ID == nil {
		message := "Silahkan Masukkan Pair Code Device di Handphone kakak"
		err := client.WAClient.Connect()
		if err != nil {
			message = err.Error()
			//qr <- QRStatus{client.PhoneNumber, false, "", message}
		}
		// No ID stored, new login
		code, err := client.WAClient.PairPhone(client.PhoneNumber, true, whatsmeow.PairClientUnknown, "Chrome (Mac OS)")
		if err != nil {
			message = err.Error()
			qr <- QRStatus{client.PhoneNumber, false, "", message}
		}
		qr <- QRStatus{client.PhoneNumber, true, code, message}

	} else {
		message := "Sudah login kak"
		if !client.WAClient.IsConnected() {
			message = "Koneksi Gagal Mencoba Melakukan Koneksi Ulang"
			err := client.WAClient.Connect()
			if err != nil {
				message = err.Error()
			}
			qr <- QRStatus{client.PhoneNumber, false, "", message}
		}
		qr <- QRStatus{client.PhoneNumber, false, "", message}

	}

}
func PairConnectStore(client *WaClient, storeMap GetStoreClient, qr chan QRStatus) {
	if client.WAClient.Store.ID == nil {
		message := "Silahkan Masukkan Pair Code Device di Handphone kakak"

		err := ConnectClient(client.WAClient)
		if err != nil {
			message = err.Error()
			//qr <- QRStatus{client.PhoneNumber, false, "", message}
		}
		// No ID stored, new login
		code, err := client.WAClient.PairPhone(client.PhoneNumber, true, whatsmeow.PairClientUnknown, "Chrome (Mac OS)")
		if err != nil {
			message = err.Error()
			qr <- QRStatus{client.PhoneNumber, false, "", message}
		}
		qr <- QRStatus{client.PhoneNumber, true, code, message}
		storeMap.StoreOnlineClient(client.PhoneNumber, client)
	} else {
		message := "Sudah login kak"
		if !client.WAClient.IsConnected() {
			message = "Koneksi Gagal Mencoba Melakukan Koneksi Ulang"
			err := client.WAClient.Connect()
			if err != nil {
				message = err.Error()
			}
			qr <- QRStatus{client.PhoneNumber, false, "", message}
		}
		qr <- QRStatus{client.PhoneNumber, false, "", message}

	}
}

func PairConnectStoreMap(client *WaClient, storeMap GetStoreClient, qr chan QRStatus) (err error) {
	err = ConnectClient(client.WAClient)

	status := QRStatus{PhoneNumber: client.PhoneNumber}

	defer func(st QRStatus, stchan chan QRStatus) {
		stchan <- st
	}(status, qr)

	if client.WAClient.Store.ID == nil {
		status.Message = "Silahkan Masukkan Pair Code Device di Handphone kakak"

		if err != nil {
			status.Message = err.Error()
			return
		}
		var code string
		code, err = client.WAClient.PairPhone(client.PhoneNumber, true, whatsmeow.PairClientUnknown, "Chrome (Mac OS)")
		if err != nil {
			status.Message = err.Error()
			return
		}

		status.Status = true
		status.Code = code

		storemapsuccess := storeMap.StoreOnlineClient(client.PhoneNumber, client)
		if !storemapsuccess {
			status.Message = "storeMap.StoreOnlineClient(client.PhoneNumber, client) Failed"
			err = errors.New(status.Message)
			return
		}
		return
	}

	status.Message = "Sudah login kak"
	if !client.WAClient.IsConnected() {
		status.Message = "Koneksi Gagal Mencoba Melakukan Koneksi Ulang"
		err = client.WAClient.Connect()
		if err != nil {
			status.Message = err.Error()
		}
		return
	}
	return
}

func RePairConnect(client *WaClient) (qr QRStatus, err error) {
	err = ConnectClient(client.WAClient)
	qr.PhoneNumber = client.PhoneNumber
	qr.Message = "Silahkan Masukkan Pair Code Device di Handphone kakak"

	if err != nil {
		qr.Message = err.Error()
		return
	}

	if client.WAClient.Store.ID == nil {
		qr.Code, err = client.WAClient.PairPhone(client.PhoneNumber, true, whatsmeow.PairClientUnknown, "Chrome (Mac OS)")
		if err != nil {
			qr.Message = err.Error()
			return
		}

		qr.Status = true
		return
	}

	qr.Status = true
	return
}

func ConnectClient(client *whatsmeow.Client) error {
	if !client.IsConnected() {
		return client.Connect()
	}
	return nil
}

func ConnectAllClient(mongoconn *mongo.Database, container *sqlstore.Container) (clients []*WaClient, err error) {
	deviceStores, err := container.GetAllDevices()
	for _, deviceStore := range deviceStores {
		client := whatsmeow.NewClient(deviceStore, waLog.Stdout("Client", "ERROR", true))
		filter := bson.M{"phonenumber": deviceStore.ID.User}
		user, err := atdb.GetOneLatestDoc[User](mongoconn, "user", filter)
		if (client.Store.ID != nil) && (err == nil) {
			var mycli WaClient
			mycli.WAClient = client
			mycli.PhoneNumber = deviceStore.ID.User
			mycli.Mongoconn = mongoconn
			mycli.register()
			ConnectClient(client)
			clients = append(clients, &mycli)
			user.DeviceID = deviceStore.ID.Device
			atdb.ReplaceOneDoc(mongoconn, "user", bson.M{"phonenumber": user.PhoneNumber}, user)
		}

	}

	return

}
