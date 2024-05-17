package wa

type StoreClient interface {
	StoreClient(id string, client *WaClient)
	StoreOnlineClient(id string, client *WaClient) (ok bool)
}

type GetClient interface {
	GetClient(id string) (client *WaClient, ok bool)
}

type GetStoreClient interface {
	StoreClient
	GetClient
}
