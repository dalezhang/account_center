package rsa

import (
	"crypto/rand"
	"crypto/rsa"
	"encoding/base64"
	"fmt"
	"testing"

	"gitee.com/dalezhang/account_center/util"
)

func TestRsa(t *testing.T) {
	result1 := "{\"timestamp\":\"2019-08-08 19:52:09\",\"type\":\"topic\",\"id\":8038}"
	encryptedData := "qlOfYacbR5wGPil2MGADezFqtr3ugd6bTmlp2XYE0EzLaZU13ic9uJcydxEHjw+uvZYbZIA4F8XQjUb/04WkMZF3m1A/WjyBwrCYOElLUy5W4jEW4656xd/P9oFm/DMtPgd/SOOntTsyF/XDq6h+gpfKCrYShjg4no9SroEKHxc="
	out, err := base64.StdEncoding.DecodeString(encryptedData)
	if err != nil {
		t.Error(err)
	}
	epkf := util.Config.GetString("rsa.private_key_path")
	privateKey := GetPrivateKey(epkf)
	result2, err := rsa.DecryptPKCS1v15(rand.Reader, privateKey, out)
	if err != nil {
		t.Error(err)
	}
	fmt.Printf("\n\n result: %+v", string(result1))
	fmt.Printf("\n result2: %+v \n", string(result2))
	if string(result1) != string(result2) {
		t.Error("not equal")
	}
}
