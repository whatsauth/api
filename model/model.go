package model

type WebHook struct {
	URL    string `bson:"url"`
	Secret string `bson:"secret"`
}

type User struct {
	PhoneNumber string  `bson:"phonenumber"`
	WebHook     WebHook `bson:"webhook"`
	Token       string  `bson:"token"`
}
