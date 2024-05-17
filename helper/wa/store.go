package wa

import (
	"fmt"

	"github.com/aiteung/atdb"
	"go.mau.fi/whatsmeow/store/sqlstore"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

func GetWaClient(phonenumber string, client []*WaClient, mongoconn *mongo.Database, container *sqlstore.Container) (waclient *WaClient, IsCreateNewClient bool, err error) {
	id, err := WithPhoneNumber(phonenumber, client, mongoconn)
	if id >= 0 {
		waclient = client[id]
	} else {
		waclient, err = CreateClientfromContainer(phonenumber, mongoconn, container)
		if err != nil {
			return
		}
		IsCreateNewClient = true
		//if IsNewClient{config.Client = append(config.Client, client)}
	}
	return
}

func GetWaClientMap(phonenumber string, client GetStoreClient, mongoconn *mongo.Database, container *sqlstore.Container) (waclient *WaClient, err error) {
	id, err := FindByPhoneNum(phonenumber, client, mongoconn)
	if id == "" {
		waclient, err = CreateClientfromContainer(phonenumber, mongoconn, container)
		if err != nil {
			return
		}
		client.StoreOnlineClient(id, waclient)
		return
	}

	waclient, ok := client.GetClient(id)
	if !ok {
		waclient, err = CreateClientfromContainer(phonenumber, mongoconn, container)
		if err != nil {
			return
		}
		client.StoreOnlineClient(id, waclient)
		return
	}
	return
}

func SetWaClient(phonenumber string, client *Clients, mongoconn *mongo.Database, container *sqlstore.Container) (waclient *WaClient, err error) {
	id, err := WithPhoneNumber(phonenumber, client.List, mongoconn)
	if id >= 0 {
		waclient = client.List[id]
	} else {
		waclient, err = CreateClientfromContainer(phonenumber, mongoconn, container)
		client.List = append(client.List, waclient)
	}
	return
}

func WithPhoneNumber(phonenumber string, clients []*WaClient, mongoconn *mongo.Database) (idx int, err error) {
	user, err := atdb.GetOneLatestDoc[User](mongoconn, "user", bson.M{"phonenumber": phonenumber})
	idx = -1
	for i, client := range clients {
		if client.WAClient.Store.ID != nil {
			if (client.WAClient.Store.ID.User == phonenumber) && (client.WAClient.Store.ID.Device == user.DeviceID) {
				idx = i
			}
		}
	}
	return
}

func FindByPhoneNum(phonenumber string, clients GetStoreClient, mongoconn *mongo.Database) (idMap string, err error) {
	user, err := atdb.GetOneLatestDoc[User](mongoconn, "user", bson.M{"phonenumber": phonenumber})
	if err != nil {
		return
	}
	idMap = fmt.Sprintf("%s-%d", phonenumber, user.DeviceID)
	return
}
