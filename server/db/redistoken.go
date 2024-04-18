package db

import (
	"birdtalk/server/utils"
	"errors"
	"strconv"
)

// 保存秘钥相关内容
func (cli *RedisClient) SaveToken(uid int64, keyEx *utils.KeyExchange) error {
	key := GetUserTokenKey(keyEx.SharedKeyPrint)

	data := make(map[string]interface{})
	data["key"] = utils.EncodeBase64(keyEx.SharedKey)
	data["keyh"] = utils.EncodeBase64(keyEx.SharedKeyHash)
	data["uid"] = strconv.FormatInt(uid, 10)
	data["enc"] = keyEx.EncType

	_, err := cli.Cmd.HMSet(key, data).Result()
	return err
}

func (cli *RedisClient) LoadToken(id int64) (int64, *utils.KeyExchange, error) {
	key := GetUserTokenKey(id)
	data, err := cli.Cmd.HGetAll(key).Result()
	if err != nil {
		return 0, nil, err
	}

	KeyEx := utils.KeyExchange{
		SharedKeyPrint: id,
	}

	var ok bool
	KeyEx.EncType, ok = data["enc"]
	tempStr, ok := data["key"]
	if !ok {
		return 0, nil, errors.New("not has field share key")
	}
	KeyEx.SharedKey, err = utils.DecodeBase64(tempStr)
	if err != nil {
		if !ok {
			return 0, nil, errors.New("field share key is not correct")
		}
	}

	tempStr, ok = data["keyh"]
	if !ok {
		return 0, nil, errors.New("not has field share key hash")
	}
	KeyEx.SharedKeyHash, err = utils.DecodeBase64(tempStr)
	if err != nil {
		if !ok {
			return 0, nil, errors.New("field share key hash is not correct")
		}
	}

	uidStr, ok := data["uid"]
	if !ok {
		return 0, &KeyEx, err
	}
	uid, err := strconv.ParseInt(uidStr, 10, 64)

	return uid, &KeyEx, err
}
