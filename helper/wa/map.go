package wa

import (
	"github.com/puzpuzpuz/xsync/v3"
)

type MapClient struct {
	*xsync.MapOf[string, *WaClient]
}

func NewMapClient(size ...int) MapClient {
	sizeHint := 10
	if len(size) > 0 {
		sizeHint = size[0]
	}

	syncer := xsync.NewMapOfPresized[string, *WaClient](sizeHint)
	return MapClient{syncer}
}

func (m *MapClient) GetClient(id string) (client *WaClient, ok bool) {
	client, ok = m.Load(id)
	return
}

func (m *MapClient) StoreClient(id string, client *WaClient) {
	m.Store(id, client)
}

func (m *MapClient) StoreOnlineClient(id string, client *WaClient) (ok bool) {
	if client.WAClient == nil {
		goto STORECLIENT
	}

	if client.WAClient.IsConnected() {
		goto STORECLIENT
	}
	if err := client.WAClient.Connect(); err != nil {
		return
	}

STORECLIENT:
	m.Store(id, client)
	ok = true
	return
}

func (m *MapClient) CheckClientOnline(id string) (ok bool) {
	client, ok := m.Load(id)
	if !ok {
		return
	}

	if client.WAClient == nil {
		return
	}

	ok = client.WAClient.IsConnected()
	return
}

func (m *MapClient) GetAllClient() (listCli []*WaClient) {

	m.Range(func(k string, v *WaClient) bool {
		listCli = append(listCli, v)
		return true
	})
	return

}

func (m *MapClient) StatusAllClient() (res map[string]bool) {
	res = make(map[string]bool, m.Size())

	m.Range(func(k string, v *WaClient) (ok bool) {
		ok = true

		if v.WAClient == nil {
			res[k] = false
			return
		}

		res[k] = v.WAClient.IsConnected()

		return
	})

	return
}

func (m *MapClient) OfflineClient() (res []string) {

	m.Range(func(k string, v *WaClient) (ok bool) {
		ok = true

		if v.WAClient == nil {
			res = append(res, k)
			return
		}

		if !v.WAClient.IsConnected() {
			res = append(res, k)
		}

		return
	})

	return
}

func (m *MapClient) StoreAllClient(listClient []*WaClient) (ok bool) {
	if len(listClient) < 1 {
		return
	}

	// (client.WAClient.Store.ID.User == phonenumber) && (client.WAClient.Store.ID.Device == user.DeviceID)
	for _, v := range listClient {
		if v.WAClient == nil {
			continue
		}

		m.Store(DefaultID(v), v)
	}
	ok = true
	return
}

func (m *MapClient) StoreAllClientCustomId(listClient []*WaClient, f func(*WaClient) string) (ok bool) {
	if len(listClient) < 1 {
		return
	}

	if f == nil {
		return
	}

	// (client.WAClient.Store.ID.User == phonenumber) && (client.WAClient.Store.ID.Device == user.DeviceID)
	for _, v := range listClient {
		m.Store(f(v), v)
	}
	ok = true
	return
}

func (m *MapClient) SetOnlineClient(id string) (ok bool) {
	cli, ok := m.Load(id)
	if !ok {
		return
	}

	if cli.WAClient == nil {
		return
	}

	if cli.WAClient.IsConnected() {
		ok = true
		return
	}

	err := cli.WAClient.Connect()
	if err != nil {
		return
	}

	ok = true
	return
}
