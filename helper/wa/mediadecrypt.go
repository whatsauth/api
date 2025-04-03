package wa

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"

	"github.com/Rhymen/go-whatsapp/crypto/cbc"
	"github.com/Rhymen/go-whatsapp/crypto/hkdf"
	"github.com/rs/xid"
)

type MediaType int

const (
	_              = iota
	MediaTypeImage = MediaType(iota)
	MediaTypeVideo
	MediaTypeAudio
	MediaTypeDocument
)

var (
	AppInfo = map[MediaType]string{
		MediaTypeImage:    "WhatsApp Image Keys",
		MediaTypeVideo:    "WhatsApp Video Keys",
		MediaTypeAudio:    "WhatsApp Audio Keys",
		MediaTypeDocument: "WhatsApp Document Keys",
	}
)

func GetBase64Filedata(urlenc *string, MediaKey []byte) string {
	encFilePath := xid.New().String()
	err := DownloadFile(encFilePath, *urlenc)
	if err != nil {
		log.Println("GetBase64Filedata,DownloadFile error ", err)
		return ""
	}
	encFileData, err := os.ReadFile(encFilePath)
	if err != nil {
		log.Println("GetBase64Filedata,ReadFile error ", err)
		return ""
	}
	data, err := decryptMedia(encFileData, MediaKey, MediaType(4))
	if err != nil {
		log.Println("GetBase64Filedata,decryptMedia error ", err)
		return ""
	}
	err = os.Remove(encFilePath)
	if err != nil {
		log.Println("GetBase64Filedata,os.Remove error: ", err)
		return ""
	}
	return base64.StdEncoding.EncodeToString(data)
}

func decryptMedia(encFileData []byte, mediaKey []byte, mt MediaType) (
	[]byte,
	error,
) {
	//
	// Implement reverse engineered media decryption algorithm from:
	// https://github.com/sigalor/whatsapp-web-reveng#decryption
	//

	// mediaKey should be 32 bytes
	if len(mediaKey) != 32 {
		return nil, fmt.Errorf("mediaKey length %d != 32",
			len(mediaKey))
	}

	mediaKeyExpanded, err := hkdf.Expand(mediaKey, 112, AppInfo[mt])
	if err != nil {
		return nil, err
	}

	iv := mediaKeyExpanded[0:16]
	cipherKey := mediaKeyExpanded[16:48]
	macKey := mediaKeyExpanded[48:80]
	//refKey := mediaKeyExpanded[80:]

	fileLen := len(encFileData) - 10
	file := encFileData[:fileLen]
	mac := encFileData[fileLen:]

	err = validateMedia(iv, file, macKey, mac)
	if err != nil {
		return nil, err
	}

	data, err := cbc.Decrypt(cipherKey, iv, file)
	if err != nil {
		return nil, err
	}

	return data, nil
}
func validateMedia(iv []byte, file []byte, macKey []byte, mac []byte) error {
	h := hmac.New(sha256.New, macKey)
	n, err := h.Write(append(iv, file...))
	if err != nil {
		return err
	}
	if n < 10 {
		return fmt.Errorf("hash to short")
	}
	if !hmac.Equal(h.Sum(nil)[:10], mac) {
		return fmt.Errorf("invalid media hmac")
	}
	return nil
}

func DownloadFile(filepath string, url string) error {

	// Get the data
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// Create the file
	out, err := os.Create(filepath)
	if err != nil {
		return err
	}
	defer out.Close()

	// Write the body to file
	_, err = io.Copy(out, resp.Body)
	return err
}
