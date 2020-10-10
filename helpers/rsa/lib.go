package rsa

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/base64"
	"encoding/pem"
	"fmt"
	"io/ioutil"
	"os"

	"gitee.com/dalezhang/account_center/util"
)

var GOPATH string
var AppPath string

func init() {
	GOPATH = os.Getenv("GOPATH")
	AppPath = "src/gitee.com/dalezhang/account_center"
}

func GetPrivateKey(path string) (privateKey *rsa.PrivateKey) {
	// mpkf := util.Config.GetString("epsp.myPrivateKeyFile")
	keyPath := fmt.Sprintf("%s/%s/%s", GOPATH, AppPath, path)
	fmt.Printf("PrivateKeyPath: %s", keyPath)
	mpk, err := ioutil.ReadFile(keyPath)
	if err != nil {
		panic(err)
	}
	block, _ := pem.Decode(mpk)
	privateKey, _ = x509.ParsePKCS1PrivateKey(block.Bytes)
	return

}

func GetPubKey(path string) (pubKey *rsa.PublicKey) {
	// epkf := util.Config.GetString("epsp.epspPublicKeyFile")
	keyPath := fmt.Sprintf("%s/%s/%s", GOPATH, AppPath, path)
	fmt.Printf("PubKeyPath: %s", keyPath)
	epk, err := ioutil.ReadFile(keyPath)
	if err != nil {
		panic(err)
	}
	block, _ := pem.Decode(epk)
	cert, _ := x509.ParseCertificate(block.Bytes)
	pubKey = cert.PublicKey.(*rsa.PublicKey)
	return
}

func Decrypt(encryptedData string) (err error, result []byte) {
	out, err := base64.StdEncoding.DecodeString(encryptedData)
	if err != nil {
		return
	}
	epkf := util.Config.GetString("rsa.private_key_path")
	privateKey := GetPrivateKey(epkf)
	result, err = rsa.DecryptPKCS1v15(rand.Reader, privateKey, out)
	if err != nil {
		return
	}
	return
}
