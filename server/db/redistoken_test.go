package db

import (
	"birdtalk/server/utils"
	"fmt"
	"testing"
)

func TestSaveToken(t *testing.T) {

	redisCli, err := NewRedisClient("127.0.0.1:6379", "")
	if err != nil {
		fmt.Println(err)
		return
	}

	token := utils.KeyExchange{
		SharedKeyPrint: int64(123456789),
		SharedKey:      []byte("123456780000000000000"),
		SharedKeyHash:  []byte("00000000123456780000000000000"),
		EncType:        "AES-CTR",
	}

	err = redisCli.SaveToken(10001, &token)
	fmt.Println(err)
}

func TestLoadToken(t *testing.T) {
	redisCli, err := NewRedisClient("127.0.0.1:6379", "")
	if err != nil {
		fmt.Println(err)
		return
	}

	uid, data, err := redisCli.LoadToken(123456789)
	fmt.Println(err)
	if err == nil {
		fmt.Println(uid, data.EncType, data.SharedKeyPrint, string(data.SharedKey), string(data.SharedKeyHash))
	}

}
