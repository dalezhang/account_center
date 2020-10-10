package session

import (
	"fmt"
	"testing"

	"golang.org/x/crypto/bcrypt"
)

func TestBcrypt(t *testing.T) {

	passwordOK := "123456"
	passwordERR := "adminxx"

	hash, err := bcrypt.GenerateFromPassword([]byte(passwordOK), bcrypt.DefaultCost)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(hash)

	encodePW := string(hash) // 保存在数据库的密码，虽然每次生成都不同，只需保存一份即可
	fmt.Println(encodePW)

	// 正确密码验证
	err = bcrypt.CompareHashAndPassword([]byte(encodePW), []byte(passwordOK))
	if err != nil {
		fmt.Println("pw wrong")
	} else {
		fmt.Println("pw ok")
	}

	// 错误密码验证
	err = bcrypt.CompareHashAndPassword([]byte(encodePW), []byte(passwordERR))
	if err != nil {
		fmt.Println("pw wrong")
	} else {
		fmt.Println("pw ok")
	}
	encryptedPass := "$2a$10$iKuxkXemWludw3dgWvKea.RFQWGOMJSkGN/3VbTr5L6weW13wDk2O"
	pass := "123456"
	err = bcrypt.CompareHashAndPassword([]byte(encryptedPass), []byte(pass))
	if err != nil {
		t.Errorf("pw from phone: 19876543211 wrong")
	} else {
		t.Errorf("pw from phone: 19876543211 ok")
	}
}
