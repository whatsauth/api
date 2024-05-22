package helper

import (
	"api/config"
	"api/model"
	"math/rand"

	"api/helper/atdb"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

func PickPantun(filter primitive.M) string {
	pantun, _ := atdb.GetOneDoc[model.Pantun](config.Mongoconn, "pantun", filter)
	pantuns := pantun.Pantun
	randomIndex := rand.Intn(len(pantuns))
	return pantuns[randomIndex]
}
