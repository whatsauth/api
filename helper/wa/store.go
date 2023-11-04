package wa

import "fmt"

var Client []*WaClient

func GetWaClient(phonenumber string, client []*WaClient) (waclient WaClient) {
	id := WithPhoneNumber(phonenumber, client)
	if id >= 0 {
		waclient = *client[id]
	} else {
		waclient = ClientDB(phonenumber)
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
