package watoken

import (
	"crypto/md5"
	"encoding/base64"
	"encoding/hex"
	"math/rand"
	"strings"
	"time"

	"golang.org/x/crypto/bcrypt"
)

func GetAppUrl(wuid string) string {
	ss := strings.Split(wuid, ".")
	s := ss[len(ss)-1]
	base64Text := make([]byte, base64.StdEncoding.DecodedLen(len(s)))
	l, _ := base64.StdEncoding.Decode(base64Text, []byte(s))
	return string(base64Text[:l])
}

func GetAppInfo(wuid string) (protocol, hostname, pathname string) {
	appurl := GetAppUrl(wuid)
	var hostandpath string
	protocol, hostandpath, _ = strings.Cut(appurl, "://")
	hostname, pathname, _ = strings.Cut(hostandpath, "/")
	return
}

func GetAppSubDomain(wuid string) (subdomain string) {
	_, hostname, _ := GetAppInfo(wuid)
	subdomain, _, _ = strings.Cut(hostname, ".")
	return
}

func RandomLowerCaseStringwithSpecialCharacter(length int) string {
	rand.Seed(time.Now().UnixNano())
	chars := []rune("abcdefghijklmnopqrstuvwxyz0123456789!@#$%^&*()_+-[]{}<>/?|=,.~`")
	var b strings.Builder
	for i := 0; i < length; i++ {
		b.WriteRune(chars[rand.Intn(len(chars))])
	}
	return b.String()
}

func RandomLowerCaseString(length int) string {
	rand.Seed(time.Now().UnixNano())
	chars := []rune("abcdefghijklmnopqrstuvwxyz0123456789")
	var b strings.Builder
	for i := 0; i < length; i++ {
		b.WriteRune(chars[rand.Intn(len(chars))])
	}
	return b.String()
}

func RandomString(length int) string {
	rand.Seed(time.Now().UnixNano())
	chars := []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789")
	var b strings.Builder
	for i := 0; i < length; i++ {
		b.WriteRune(chars[rand.Intn(len(chars))])
	}
	return b.String()
}

func GetMD5Hash(text string) string {
	hash := md5.Sum([]byte(text))
	return hex.EncodeToString(hash[:])
}

func GetBcryptHash(text string) string {
	bytes, _ := bcrypt.GenerateFromPassword([]byte(text), 14)
	return string(bytes)
}
