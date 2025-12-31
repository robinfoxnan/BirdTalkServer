package core

import (
	"birdtalk/server/db"
	"birdtalk/server/model"
	"birdtalk/server/pbmodel"
	"birdtalk/server/utils"
	"errors"
	"go.uber.org/zap"
)

// 通知好友当前个人信息或者在线状态改变的消息
func createSetUserInfoMessage(user *model.User, status string) *pbmodel.Msg {
	params := map[string]string{
		"status": status,
	}

	msgUserOp := pbmodel.UserOpReq{
		Operation: pbmodel.UserOperationType_SetUserInfo,
		User:      user.GetUserInfo(),
		SendId:    Globals.snow.GenerateID(),
		MsgId:     Globals.snow.GenerateID(),
		Params:    params,
	}

	msgPlain := pbmodel.MsgPlain{
		Message: &pbmodel.MsgPlain_UserOp{
			UserOp: &msgUserOp,
		},
	}

	msg := pbmodel.Msg{
		Version:  int32(ProtocolVersion),
		KeyPrint: 0,
		Tm:       utils.GetTimeStamp(),
		MsgType:  pbmodel.ComMsgType_MsgTUserOp,
		SubType:  0,
		Message: &pbmodel.Msg_PlainMsg{
			PlainMsg: &msgPlain,
		},
	}
	return &msg
}

// 设置用户在线与不在线，使用Set
func setUserOnline(user *model.User, bOnLine bool) {
	lst := user.GetMutualFriends()

	if bOnLine {
		msg := createSetUserInfoMessage(user, "online")
		trySendMsgToUserList(lst, msg)

	} else {
		// 如果会话断开了，那么应该检测是否已经完全离线
		b, _ := checkUserOnline(user.UserId)
		if !b {
			msg := createSetUserInfoMessage(user, "offline")
			trySendMsgToUserList(lst, msg)
		}

	}

}

// 一个人员加入到群的时候，保存人员信息
// 类似 groupHandler.go line-1353   func saveNewGroup(){}
func saveUserJoinGroup(user *model.User, group *model.Group) error {
	if user == nil {
		return errors.New("")
	}

	member := model.GroupMemberStore{
		Pk:   db.ComputePk(group.GroupId),
		Gid:  group.GroupId,
		Uid:  user.UserId,
		Tm:   utils.GetTimeStamp(),
		Role: model.RoleGroupMember,
		Nick: user.NickName,
	}

	item := model.UserInGStore{
		Pk:  db.ComputePk(user.UserId),
		Uid: user.UserId,
		Gid: group.GroupId,
	}

	// 1) 保存数据库
	err := Globals.scyllaCli.InsertGroupMember(&member, &item)
	if err != nil {
		return err
	}

	// 2)redis, 用户所在群
	// 这里做2个动作，2.1 群里有用户 2.2 用户所在群列表
	err = Globals.redisCli.SetUserJoinGroup(user.UserId, group.GroupId, user.NickName)
	if err != nil {
		// 这里不用返回，因为下次反正会同步数据
		Globals.Logger.Error("")
	}

	// 添加到内存
	group.AddMember(user.UserId, user.NickName, user)

	// 内存数据，设置用户在这个群里
	user.SetJoinGroup(group.GroupId)
	return nil
}

// 这个函数简单的调用旁边那个
func findUserInfo(uid int64) (*pbmodel.UserInfo, bool, error) {
	user, ok, err := findUser(uid)
	if err != nil {
		return nil, false, err
	}
	return user.GetUserInfo(), ok, nil
}

// 后面会有各种地方需要找到好友的信息
// 这个函数会造成加载内存
func findUser(uid int64) (*model.User, bool, error) {
	user, ok := Globals.uc.GetUser(uid)
	if ok && user != nil {
		return user, true, nil
	}

	userInfo, err := loadUserByFriend(uid)
	if userInfo != nil {
		user = model.NewUserFromInfo(userInfo)
		Globals.uc.SetOrUpdateUser(uid, user)
		return user, false, nil
	} else {
		return nil, false, err
	}
}

// 由其他人搜索才造成的加载redis, 基本信息，并且只加载一条fid的权限
// 至于uid是否关注了fid, 则从fid的粉丝表中加载；
func loadUserByFriend(fid int64) (*pbmodel.UserInfo, error) {

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
