package core

import (
	"birdtalk/server/model"
	"errors"
	"fmt"
	"log"
)

func GetAllUsers(t string) ([]*model.User, error) {
	switch t {
	case "mem":
		// 从内存中获取用户
		users := Globals.uc.GetAllUsers()
		return users, nil

	case "db":
		// 从数据库获取用户（假设 core.DB 或 core.UserRepo 可用）
		users, err := Globals.mongoCli.LoadUsersFromId(10000, 100)
		if err != nil {
			log.Println("load users error:", err)
			return nil, err
		}

		//for _, u := range users {
		//	fmt.Println("user:", u.UserId, u.UserName)
		//}
		return users, nil
	case "redis":
		users, err := Globals.redisCli.GetUsersByRange("bsui_", 10000, 100)
		if err != nil {
			fmt.Println("Redis error:", err)
			return nil, err
		}

		for key, fields := range users {
			fmt.Printf("User: %s: %v\n", key, fields)
		}
		return users, nil

	default:

	}

	return nil, errors.New("type error")
}
