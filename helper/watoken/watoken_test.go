package watoken

import (
	"os"
	"testing"
)

func TestUpdateSenderRecoiver(t *testing.T) {
	privkey := os.Getenv("PRIVATEKEY")
	res, _ := EncodeforHours("628888", "123123", privkey, 43800)
	println(res)
	println("=====================")

	//pyl, _ := Decode("", res)
	//println(pyl.Id)

}
