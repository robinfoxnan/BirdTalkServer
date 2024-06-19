package core

import (
	"birdtalk/server/db"
	"birdtalk/server/model"
	"birdtalk/server/pbmodel"
	"birdtalk/server/utils"
	"errors"
	"fmt"
	"go.uber.org/zap"
	"strings"
)

// 检查当前客户端的协议版本号；
func checkProtoVersion(msg *pbmodel.Msg, session *Session) bool {
	ver := int(msg.GetVersion()) // 协议版本号
	if ver != ProtocolVersion {
		sendBackErrorMsg(int(pbmodel.ErrorMsgType_ErrTVersion), "version is too high", nil, session)
		return false
	}
	return true
}

// 检查秘钥指纹是否存在
func checkKeyPrint(key int64) (int, string) {
	return int(pbmodel.ErrorMsgType_ErrTKeyPrint), "should be 0"
}

// 检查是否需要重定向
func checkNeedRedirect(session *Session) (int, string) {
	//return int(pbmodel.ErrorMsgType_ErrTRedirect), "should redirect"

	return int(pbmodel.ErrorMsgType_ErrTNone), ""
}

// 解码对方的公钥，有可能是加密，根据rasPrint区分
func decodeRemotePublicKey(exMsg *pbmodel.MsgKeyExchange, session *Session) ([]byte, error) {
	rsaPrint := exMsg.GetRsaPrint()
	if rsaPrint == 0 {

		publicKey := exMsg.GetPubKey()
		//fmt.Println("remote public key=", string(publicKey))
		Globals.Logger.Debug("decode public key from msg:", zap.String("publicKey", string(publicKey)),
			zap.Int64("session", session.Sid))
		if len(publicKey) < 65 {
			sendBackErrorMsg(int(pbmodel.ErrorMsgType_ErrTPublicKey), "", nil, session)
			return nil, errors.New("public key len less bytes")
		}
		return publicKey, nil
	}
	// 使用RAS解码对称密钥，使用RSA解码对方公钥；

	sendBackErrorMsg(int(pbmodel.ErrorMsgType_ErrTRsaPrint), "", nil, session)
	return nil, errors.New("decode remote public key error")
}

// 自动选择算法执行加密
func encryptDataAuto(data []byte, session *Session) ([]byte, error) {
	encType := strings.ToLower(session.KeyEx.EncType)
	switch encType {
	case "chacha20":
		return utils.EncryptChaCha20(data, session.KeyEx.SharedKeyHash)
	case "aes-ctr":
		return utils.EncryptAES_CTR(data, session.KeyEx.SharedKeyHash)
	case "twofish":
		return utils.EncryptTwofish(data, session.KeyEx.SharedKeyHash)
	}

	return nil, errors.New("not supported encrypt algorithm")
}

func decryptDataAuto(data []byte, session *Session) ([]byte, error) {
	encType := strings.ToLower(session.KeyEx.EncType)
	switch encType {
	case "chacha20":
		return utils.DecryptChaCha20(data, session.KeyEx.SharedKeyHash)
	case "aes-ctr":
		return utils.DecryptAES_CTR(data, session.KeyEx.SharedKeyHash)
	case "twofish":
		return utils.DecryptTwofish(data, session.KeyEx.SharedKeyHash)
	}

	return nil, errors.New("not supported encrypt algorithm")
}

// 收到消息，先检查用户是否登录了
func checkUserLogin(sess *Session) bool {

	user, ok := Globals.uc.GetUser(sess.UserID)
	if ok {
		if user.Params != nil {
			status, b := user.Params["status"]
			if b {
				if status == "disabled" {
					sendBackErrorMsg(int(pbmodel.ErrorMsgType_ErrTDisabled), "user is disabled.", nil, sess)
					return false
				}
			}
		}

		if sess.HasStatus(model.UserStatusOk) {
			return true
		}
	}

	sendBackErrorMsg(int(pbmodel.ErrorMsgType_ErrTNotLogin), "should login first.", nil, sess)
	return false
}

// 是否用户已经注销
func checkUserIsUsing(sess *Session) bool {

	user, ok := Globals.uc.GetUser(sess.UserID)
	if ok {
		if user.Params != nil {
			status, b := user.Params["status"]
			if b {
				if status == "deleted" {
					sendBackErrorMsg(int(pbmodel.ErrorMsgType_ErrTDeleted), "user is deleted.", nil, sess)
					return false
				}
			}
		}

		//if sess.HasStatus(model.UserStatusOk) {
		//	return true
		//}
	}

	//sendBackErrorMsg(int(pbmodel.ErrorMsgType_ErrTNotLogin), "should login first.", nil, sess)
	return true
}

// 先放到会话中，此时
func createTempUser(userInfo *pbmodel.UserInfo, session *Session) error {
	// 保存临时信息
	session.TempUserInfo = userInfo
	params := userInfo.GetParams()
	if params == nil {
		userInfo.Params = map[string]string{
			"email": "checked",
		}
	} else {
		params["email"] = "checked"
	}
	// 创建临时验证码
	code := utils.GenerateCheckCode(5)

	Globals.Logger.Info("user reg validate info", zap.String("code", code))
	session.SetKeyValue("code", code)
	session.SetStatus(model.UserStatusRegister | model.UserStatusValidate)

	// 发送验证码
	//
	SendEmailCode(session, userInfo.Email, code)

	return nil
}

// 检查验证码
func checkValidateCode(code string, session *Session) (bool, error) {
	if !session.HasStatus(model.UserStatusValidate) {
		return false, errors.New("not with status model.UserStatusValidate")
	}

	value := session.GetKeyValue("validateCode")
	if len(value) == 0 {
		return false, errors.New("no code here")
	}

	if code == value {
		return true, nil
	}

	return false, nil
}

// 创建新用户
func createUser(userInfo *pbmodel.UserInfo, session *Session) error {
	uid, err := Globals.redisCli.GetNextUserId()
	if err != nil {
		sendBackErrorMsg(int(pbmodel.ErrorMsgType_ErrTServerInside),
			"email is not correct",
			map[string]string{"field": userInfo.Email},
			session)
		Globals.Logger.Error("redis get next user id error")
		return errors.New("redis get next user id error")
	}

	// 保存用户信息
	userInfo.UserId = uid
	err = Globals.mongoCli.CreateNewUser(userInfo)

	if err != nil {
		sendBackErrorMsg(int(pbmodel.ErrorMsgType_ErrTServerInside),
			"email is not correct",
			map[string]string{"field": userInfo.Email},
			session)
		Globals.Logger.Error("save new user error")
		return errors.New("redis get next user id error")
	}

	// 保存到回话中
	session.UserID = uid
	session.SetStatus(model.UserWaitLogin)
	return nil
}

// 登录后设置各种信息,完整加载
// 1) 从数据库加载
// 2) 保存到redis
// 3) 同步内存
// 4) 将用户绑定到指纹
// 5) 发送同步信息，登录
// 6) 加载粉丝
// 7) 加载关注
// 8) 加载权限
// 9) 加载所在的群组列表
func LoadUserNew(session *Session) error {
	// 1) 数据库
	userInfos, err := Globals.mongoCli.FindUserById(session.UserID)
	if err != nil {
		return err
	}
	if len(userInfos) != 1 {
		Globals.Logger.Error("load user err, find count of user", zap.Int("userCount", len(userInfos)))
		return errors.New("load user count is not 1")
	}
	//userinfo := &userInfos[0]

	session.TempUserInfo = userInfos[0]
	// 2) 更新redis
	err = Globals.redisCli.SetUserInfo(userInfos[0])

	// 3) 内存缓存, 多终端登录，所以需要，但是这里是新注册，所以冲突的概率几乎为0
	user, ok := Globals.uc.GetUser(session.UserID)
	if ok {
		user.AddSessionID(session.Sid)
	} else {
		// 保存用户信息
		user = model.NewUserFromInfo(userInfos[0])
		user.MaskLoad = model.UserLoadStatusInfo // 目前仅仅加载了基本的数据
		user.AddSessionID(session.Sid)
		Globals.uc.SetOrUpdateUser(user.UserId, user) // 插入或者合并
	}

	// 4) 绑定指纹
	if session.KeyEx != nil {
		Globals.redisCli.SaveToken(session.UserID, session.KeyEx)
	}
	session.Status = model.UserLoadStatusAll

	// 5) todo: 广播信息

	return nil
}

// 仅仅是从redis 加载到内存
// 加载涉及到7个表，哪个出错了就返回掩码，但是不是全部加载到内存
// 1：用户基础信息表，bsui_
// 2：用户好友权限表，bsufb_
// 4：用户关注表，bsufo_
// 8：用户粉丝表，bsufa_
// 16：用户所属群组，bsuing_
func tryLoadUserFromRedis(sess *Session) (*model.User, uint32, error) {
	var mask uint32 = model.UserLoadStatusAll

	// 加载后尝试添加到内存缓存中
	user, ok := Globals.uc.GetUser(sess.UserID)
	if ok {
		user.AddSessionID(sess.Sid)
	} else {

		// 1 查看是否有基本的信息
		userInfo, err := Globals.redisCli.FindUserById(sess.UserID)
		if err != nil {
			return nil, mask, err
		}
		if userInfo.UserId == 0 {
			fmt.Println("userinfo is default")
			return nil, mask, err
		}
		loadedUser := model.NewUserFromInfo(userInfo)
		loadedUser.AddSessionID(sess.Sid)
		loadedUser.SetLoadMask(model.UserLoadStatusInfo)

		user = Globals.uc.SetOrUpdateUser(sess.UserID, loadedUser) // 插入或者合并
	}

	mask = model.UserLoadStatusAll & (^model.UserLoadStatusInfo)

	// 2 好友仅仅加载权限,权限全都加载到内存
	off := uint64(0)
	var permissionMap map[int64]uint32
	count := 0
	var err error
	for {
		off, permissionMap, err = Globals.redisCli.GetUserBLocks(sess.UserID, off)
		if err == nil {
			user.AddPermission(permissionMap)
		} else {
			break
		}
		count += len(permissionMap)
		fmt.Printf("len= %d off= %d, count = %d\n", len(permissionMap), off, count)

		// 查询结束
		if off == 0 {
			if err == nil {
				user.SetLoadMask(model.UserLoadStatusInfo)
			}
			break
		}
	}

	if count > 0 {
		mask = model.UserLoadStatusAll & (^model.UserLoadStatusPermission)
	}

	// 3 & 4 关注和粉丝都是仅仅检查一下redis中是否有；
	//off, fMap, err := Globals.redisCli.GetUserFollowing(sess.UserID, 0)
	//if err == nil {
	//
	//	user.SetLoadMask(model.UserLoadStatusFollow)
	//	mask = model.UserLoadStatusAll & (^model.UserLoadStatusFollow)
	//}
	//
	//off, fMap, err = Globals.redisCli.GetUserFans(sess.UserID, 0)
	//if err == nil {
	//
	//	user.SetLoadMask(model.UserLoadStatusFans)
	//	mask = model.UserLoadStatusAll & (^model.UserLoadStatusFans)
	//}

	// 用户在组中的
	gList, err := Globals.redisCli.GetUserInGroupAll(sess.UserID)
	if err == nil {
		user.SetInGroup(gList)
		mask = model.UserLoadStatusAll & (^model.UserLoadStatusGroups)
	}

	return user, mask, err
}

// 根据掩码来决定设置加载哪些
func loadUserFromDb(sess *Session, mask uint32) error {
	if mask == 0 {
		return nil
	}

	// 加载后尝试添加到内存缓存中
	user, ok := Globals.uc.GetUser(sess.UserID)
	if ok {
		user.AddSessionID(sess.Sid)
	} else {
		// 基础信息
		userInfos, err := Globals.mongoCli.FindUserById(sess.UserID)
		if err != nil {
			return err
		}
		if len(userInfos) != 1 {
			Globals.Logger.Error("loadUserFromDb() load user err, find count of user", zap.Int64("userid", sess.UserID), zap.Int("userCount", len(userInfos)))
			return errors.New("load user count is not 1")
		}
		userInfo := userInfos[0]
		Globals.redisCli.SetUserInfo(userInfo) // 同步到redis

		loadedUser := model.NewUserFromInfo(userInfo)
		loadedUser.AddSessionID(sess.Sid)
		loadedUser.SetLoadMask(model.UserLoadStatusInfo)

		user = Globals.uc.SetOrUpdateUser(sess.UserID, loadedUser) // 插入或者合并
	}

	//numFollow := 0
	//numFans := 0
	if (mask & model.UserLoadStatusFans) > 0 {
		fList, err := Globals.scyllaCli.FindFans(db.ComputePk(sess.UserID), sess.UserID, 0, db.MaxFriendCacheSize)
		if err != nil {
			return err
		}
		//numFans = len(fList)

		Globals.redisCli.SetUserFans(sess.UserID, fList) // 同步到redis
	}

	if (mask & model.UserLoadStatusFollow) > 0 {
		fList, err := Globals.scyllaCli.FindFollowing(db.ComputePk(sess.UserID), sess.UserID, 0, db.MaxFriendCacheSize)
		if err != nil {
			return err
		}
		//numFollow = len(fList)

		Globals.redisCli.SetUserFollowing(sess.UserID, fList) // 同步到redis
	}

	// 改为记录在用户的mongo中
	// 如果加载全了，则同步到缓存，入股不全就不同步了，
	//if numFollow < db.MaxFriendCacheSize && numFans < db.MaxFriendCacheSize {
	//	Globals.redisCli.SetUserFriendNum(sess.UserID, int64(numFollow), int64(numFans))
	//	user.NumFans = int64(numFans)
	//	user.NumFollow = int64(numFollow)
	//} else {
	//	user.NumFollow, user.NumFans, _ = Globals.redisCli.GetUserFriendNum(sess.UserID)
	//}

	if (mask & model.UserLoadStatusPermission) > 0 {
		count := 0
		fromId := int64(0)
		for {
			perList, err := Globals.scyllaCli.FindBlocks(sess.UserID, sess.UserID, fromId, db.MaxFriendCacheSize)
			if err != nil {
				return err
			}
			if len(perList) < db.MaxFriendCacheSize {
				break
			}
			count += len(perList)

			fromId = perList[len(perList)-1].Uid2

			// 加载到resdis与内存
			Globals.redisCli.AddUserBlocks(sess.UserID, perList)
			user.AddPermissionFromDb(perList)
		}
		Globals.Logger.Debug("Load user permission list", zap.Int64("user", sess.UserID),
			zap.Int("count", count))
	}

	if (mask & model.UserLoadStatusGroups) > 0 {
		gList, err := Globals.scyllaCli.FindUserInGroups(sess.UserID, sess.UserID, 0, db.MaxFriendCacheSize)
		if err != nil {
			return err
		}

		groupSet := make([]int64, len(gList))
		for i, item := range gList {
			groupSet[i] = item.Gid
		}

		// 加载到resdis与内存
		Globals.redisCli.SetUserInGroup(sess.UserID, groupSet)
		user.SetInGroup(groupSet)
	}
	return nil
}

// 用户登录时候加载
// 这里既然登录了，加载的东西相对比较全
func LoadUserLogin(session *Session) error {
	// 如果内存里有，就不用继续了
	user, ok := Globals.uc.GetUser(session.UserID)
	if ok {
		user.AddSessionID(session.Sid)
		return nil
	}
	user, mask, err := tryLoadUserFromRedis(session)
	fmt.Println("tryLoadUserFromRedis() mask=", mask)
	if err != nil {
		//
	}

	if mask > 0 {
		err = loadUserFromDb(session, mask)
	}

	return err

}

// 这个函数是给查询好友使用的
func findUserMongoRedis(uid int64) ([]*pbmodel.UserInfo, error) {
	user, ok := Globals.uc.GetUser(uid)
	if ok && user != nil {
		return []*pbmodel.UserInfo{&user.UserInfo}, nil
	}

	userInfo, err := LoadUserByFriend(uid)
	if userInfo != nil {
		return []*pbmodel.UserInfo{userInfo}, nil
	} else {
		return []*pbmodel.UserInfo{}, err
	}
}

// 后面会有各种地方需要找到好友的信息
// 这个函数会造成加载内存
func findUser(uid int64) (*model.User, bool, error) {
	user, ok := Globals.uc.GetUser(uid)
	if ok && user != nil {
		return user, true, nil
	}

	userInfo, err := LoadUserByFriend(uid)
	if userInfo != nil {
		user = model.NewUserFromInfo(userInfo)
		Globals.uc.SetOrUpdateUser(uid, user)
		return user, false, nil
	} else {
		return nil, false, err
	}
}

// 这个函数不会加载内存
func findUserInfo(uid int64) (*pbmodel.UserInfo, bool, error) {
	user, ok := Globals.uc.GetUser(uid)
	if ok && user != nil {
		return &user.UserInfo, true, nil
	}

	userInfo, err := LoadUserByFriend(uid)
	if userInfo != nil {
		return userInfo, false, nil
	} else {
		return nil, false, err
	}
}

// 由其他人搜索才造成的加载redis, 基本信息，并且只加载一条fid的权限
// 至于uid是否关注了fid, 则从fid的粉丝表中加载；
func LoadUserByFriend(fid int64) (*pbmodel.UserInfo, error) {

	// 1 查看是否有基本的信息
	userInfo, err := Globals.redisCli.FindUserById(fid)
	if userInfo != nil {
		return userInfo, nil
	}

	// 基础信息
	userInfos, err := Globals.mongoCli.FindUserById(fid)
	if err != nil {
		return nil, err
	}
	if len(userInfos) > 1 {
		Globals.Logger.Error("loadUserFromDb() load user err, find count of user", zap.Int("userCount", len(userInfos)))
		return nil, errors.New("load user count is not 1")
	}

	if len(userInfos) == 0 {
		return nil, errors.New("no user in db")
	}
	userInfo = userInfos[0]
	Globals.redisCli.SetUserInfo(userInfo) // 同步到redis

	return userInfo, nil
}

// 先检查内存是否有，然后检查redis中是否有
func checkUserOnline(fid int64) (bool, map[int64]int32) {
	user, ok := Globals.uc.GetUser(fid)
	if ok {
		count := user.GetSessionCount()
		return count > 0, nil
	}

	sessionMap, err := Globals.redisCli.GetUserSessionOnServer(fid)
	if err != nil {
		return false, nil
	}

	if len(sessionMap) == 0 {
		return false, nil
	}

	return true, sessionMap
}

// 检查对方给自己设置的权限，
// 首先看内存，如果没有，查对方redis给设计的权限, 如果没有就查数据库，如果实在没有设置默认的
func checkFriendPermission(uid, fid int64, bFan bool, bits uint32) bool {
	user, ok := Globals.uc.GetUser(uid)
	bRet := false
	mask := uint32(0)
	if ok {
		// 未设置有2种情况
		bRet, ok = user.CheckFriendToMeMask(fid, bits)
		if ok {
			return bRet
		}
	}

	// 从redis查对方给自己设置的权限，如果有就添加到自己的内存中
	mask, err := Globals.redisCli.CheckUserPermission(fid, uid)
	if err == nil && mask != 0 {
		// 借助基础函数来识别
		user.SetFriendToMeMask(fid, mask|model.PermissionMaskExist)
		bRet, ok = user.CheckFriendToMeMask(fid, bits)
		return bRet
	}

	// 查看对方自己设置的权限，默认是啥也没有的
	permission, err := Globals.scyllaCli.FindBlocksExact(db.ComputePk(fid), fid, uid)

	// 数据库中也没有，实在没有设置返回默认的，
	if err != nil || permission == nil {

		if bFan {
			mask = model.DefaultPermissionP2P | model.PermissionMaskFriend | model.PermissionMaskExist
		} else {
			mask = model.DefaultPermissionStranger | model.PermissionMaskExist
		}
	} else { // 找到了

		Globals.redisCli.AddUserPermission(fid, uid, uint32(permission.Perm|model.PermissionMaskExist))
		mask = uint32(permission.Perm)
	}

	// 借助基础函数来识别
	user.SetFriendToMeMask(fid, mask)
	bRet, ok = user.CheckFriendToMeMask(fid, bits)
	return bRet
}

// 检查内存，如果没有就检查redis,再没有查数据库，redis非粉丝设置nick为“”
func checkFriendIsFan(fid, uid int64) (bool, error) {
	user, ok := Globals.uc.GetUser(uid)

	// 这个函数为通用函数，不一定是在线用户
	if ok {
		isFan, bHas := user.CheckFun(fid)
		if bHas {
			return isFan, nil
		}
	}

	// 内存未设置，则应该从redis开始查找
	bFan, err := Globals.redisCli.CheckUserFan(uid, fid)
	// 主找到了，如果没好有找到，则err = redis: nil
	if err == nil {
		if user != nil {
			user.SetFan(fid, bFan)
		}
		return bFan, nil
	}

	// 没有找到，则从数据库开始加载
	//friend, err := Globals.scyllaCli.FindFansExact(db.ComputePk(uid), uid, fid)
	// 粉丝可能几百万，关注永远不会那么多
	friend, err := Globals.scyllaCli.FindFollowingExact(db.ComputePk(fid), fid, uid)

	// 数据库也没有设置，那就对方不是自己的粉丝，自己不是对方的朋友
	if friend == nil {
		fan := model.FriendStore{
			Pk:   db.ComputePk(uid),
			Uid1: uid,
			Uid2: fid,
			Tm:   0,
			Nick: "##",
		}
		Globals.redisCli.AddUserFans(uid, []model.FriendStore{fan})
		if user != nil {
			user.SetFan(fid, false)
		}
		return false, nil
	}

	// 数据库中找到了，那对方是自己的粉丝，自己是对方的朋友
	Globals.redisCli.AddUserFans(uid, []model.FriendStore{*friend})
	if user != nil {
		user.SetFan(fid, true)
	}
	return true, nil
}

// 格式转换
func FriendStore2UserInfo(lst []model.FriendStore) []*pbmodel.UserInfo {
	if lst == nil {
		return nil
	}
	retLst := make([]*pbmodel.UserInfo, len(lst))
	for index, f := range lst {
		data := pbmodel.UserInfo{
			UserId:   f.Uid1,
			UserName: "",
			NickName: f.Nick,
			Email:    "",
			Phone:    "",
			Gender:   "",
			Age:      0,
			Region:   "",
			Icon:     "",
			Params:   nil,
		}
		retLst[index] = &data
	}

	return retLst
}

// 对应答的数据过滤掉多余的信息
func filterUserInfo(userList []*pbmodel.UserInfo, mode string) {
	for _, p := range userList {
		filterUserInfo1(p, mode)
	}
}

func filterUserInfo1(p *pbmodel.UserInfo, mode string) {
	if p.Params != nil {
		delete(p.Params, "pwd")
		delete(p.Params, "friendaddmode")
		delete(p.Params, "friendaddanswer")
	}
	if mode != "phone" {
		p.Phone = "*"
	}
	if mode != "email" {
		p.Email = "*"
	}
}
