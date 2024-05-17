package wa

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"net/url"
)

func PostStructWithToken[T any](tokenkey string, tokenvalue string, structname interface{}, urltarget string) (result T, err error) {
	_, err = url.ParseRequestURI(urltarget)
	if err != nil {
		err = errors.New("> ⓘ URL WebHook tidak valid:*" + urltarget + "* ERROR: " + err.Error())
		return
	}
	client := http.Client{}
	mJson, err := json.Marshal(structname)
	if err != nil {
		err = errors.New("> ⓘ Struct input tidak valid. json.Marshal(structname) ERROR: " + err.Error())
		return
	}
	req, err := http.NewRequest("POST", urltarget, bytes.NewBuffer(mJson))
	if err != nil {
		err = errors.New("> ⓘ http.NewRequest: error inisiasi request " + urltarget + ". ERROR: " + err.Error())
		return
	}
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add(tokenkey, tokenvalue)
	resp, err := client.Do(req)
	if err != nil {
		err = errors.New("> ⓘ client.Do: gagal melakukan request ke " + urltarget + ". ERROR: " + err.Error())
		return
	}
	defer resp.Body.Close()
	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		err = errors.New("> ⓘ io.ReadAll: Gagal membaca body respon dari " + urltarget + ". ERROR: " + err.Error())
		return
	}
	if err = json.Unmarshal(respBody, &result); err != nil {
		rawstring := string(respBody)
		err = errors.New("> ⓘ Not A Valid JSON Response from " + urltarget + ". CONTENT: " + rawstring)
		return
	}
	return
}
