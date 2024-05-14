package config

import "os"

var PublicKey string = os.Getenv("PUBLICKEY")
var PrivateKey string = os.Getenv("PRIVATEKEY")
