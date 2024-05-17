package wa

import "fmt"

func DefaultID(client *WaClient) string {
	return fmt.Sprintf("%s-%d", client.WAClient.Store.ID.User, client.WAClient.Store.ID.Device)
}
