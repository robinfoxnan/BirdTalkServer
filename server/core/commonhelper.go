package core

import (
	"birdtalk/server/model"
	"birdtalk/server/pbmodel"
	"birdtalk/server/utils"
	"errors"
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

func checkUserLogin(session *Session) bool {
	return false
}

func checkUserPermission(session *Session) bool {
	return true
}

// 先放到会话中，此时
func createTempUser(userInfo *pbmodel.UserInfo, session *Session) error {
	// 保存临时信息
	session.TempUserInfo = userInfo
	// 创建临时验证码
	code := utils.GenerateCheckCode(5)
	Globals.Logger.Info("user reg validate info", zap.String("code", code))
	session.SetKeyValue("validateCode", code)
	session.SetStatus(model.UserStatusRegister | model.UserStatusValidate)

	// 发送验证码
	//

	return nil
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

	// 2) 更新redis
	err = Globals.redisCli.SetUserInfo(&userInfos[0])

	// 3) 内存缓存, 多终端登录，所以需要
	user, ok := Globals.uc.GetUser(session.UserID)
	if ok {
		user.AddSessionID(session.Sid)
	} else {
		// 保存用户信息
		user = model.NewUserFromInfo(&userInfos[0])
		user.AddSessionID(session.Sid)
	}

	// 4) 绑定指纹
	if session.KeyEx != nil {
		Globals.redisCli.SaveToken(session.UserID, session.KeyEx)
	}

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
func justLoadUserFromRedis(sess *Session) (*model.User, uint32, error) {
	var mask uint32 = 1 | 2 | 4 | 8 | 16

	// 1
	userInfo, err := Globals.redisCli.FindUserById(sess.UserID)
	if err != nil {
		return nil, mask, err
	}
	user := model.NewUserFromInfo(userInfo)
	user.AddSessionID(sess.Sid)
	mask = 2 | 4 | 8 | 16

	// 4 好友，没有必要全部加载内存

	// 设置或者更新
	Globals.uc.SetOrUpdateUser(sess.UserID, user)
	return user, mask, err
}

func loadUserFromDb(session *Session, mask uint32) error {
	return nil
}

// 用户登录时候加载
func LoadUserLogin(session *Session) error {
	// 如果内存里有，就不用继续了
	user, ok := Globals.uc.GetUser(session.UserID)
	if ok {
		user.AddSessionID(session.Sid)
		return nil
	}
	user, mask, err := justLoadUserFromRedis(session)
	if err != nil {
		return err
	}

	if mask > 0 {
		err = loadUserFromDb(session, mask)
	}

	return err

}

func SendBackUserOp(opCode pbmodel.UserOperationType, userInfo *pbmodel.UserInfo, ret bool, status string, session *Session) error {
	return nil
}
