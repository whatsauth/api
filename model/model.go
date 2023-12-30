package model

type Header struct {
	Token string `reqHeader:"token"`
}

type Pantun struct {
	Notif  string   `bson:"notif"`
	Pantun []string `bson:"pantun"`
}
