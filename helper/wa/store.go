package wa

import (
	"fmt"

	"go.mongodb.org/mongo-driver/mongo"
)

var Client []*WaClient

func GetWaClient(phonenumber string, client []*WaClient, mongoconn *mongo.Database) (waclient WaClient) {
	id := WithPhoneNumber(phonenumber, client)
	if id >= 0 {
		waclient = *client[id]
	} else {
		waclient = ClientDB(phonenumber, mongoconn)
		client = append(client, &waclient)
	}
	return
}

func WithPhoneNumber(phonenumber string, clients []*WaClient) int {
	for i, client := range clients {
		fmt.Println(client)
		if client.PhoneNumber == phonenumber {
			return i
		}
	}
	return -1
}
