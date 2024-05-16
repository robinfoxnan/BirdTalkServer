package core

import (
	"birdtalk/server/pbmodel"
	"go.uber.org/zap"
)

// 群组相关的基本操作
func handleGroupOp(msg *pbmodel.Msg, session *Session) {
	groupOpMsg := msg.GetPlainMsg().GetGroupOp()
	if groupOpMsg == nil {
		Globals.Logger.Debug("receive wrong group op msg",
			zap.Int64("sid", session.Sid),
			zap.Int64("uid", session.UserID))
		sendBackErrorMsg(int(pbmodel.ErrorMsgType_ErrTMsgContent), "group op  is null", nil, session)
		return
	}

	// 都需要验证是否登录与权限
	ok := checkUserLogin(session)
	if !ok {
		sendBackErrorMsg(int(pbmodel.ErrorMsgType_ErrTNotLogin), "should login first.", nil, session)
		return
	}

	opCode := groupOpMsg.Operation
	switch opCode {
	case pbmodel.GroupOperationType_GroupCreate: // 创建
		break
	case pbmodel.GroupOperationType_GroupDissolve: // 解散
		break
	case pbmodel.GroupOperationType_GroupSetInfo: // 设置信息
		break
	case pbmodel.GroupOperationType_GroupKickMember: // 踢人
		break
	case pbmodel.GroupOperationType_GroupInviteRequest: // 邀请
		break
	case pbmodel.GroupOperationType_GroupInviteAnswer: // 邀请的应答
		break // 邀请后处理结果
	case pbmodel.GroupOperationType_GroupJoinRequest:
		break // 加入请求
	case pbmodel.GroupOperationType_GroupJoinAnswer:
		break // 加入请求的处理，同意、拒绝、问题
	case pbmodel.GroupOperationType_GroupQuit:
		break // 退出群组
	case pbmodel.GroupOperationType_GroupAddAdmin:
		break // 增加管理员
	case pbmodel.GroupOperationType_GroupDelAdmin:
		break // 删除管理员
	case pbmodel.GroupOperationType_GroupTransferOwner:
		break // 转让群主
		// 可以根据需要添加其他群组操作 case pbmodel.GroupOperationType_GroupSetMemberInfo GroupOperationType = 13 // 设置自己在群中的信息
	case pbmodel.GroupOperationType_GroupSearch:
		break // 搜素群组
	case pbmodel.GroupOperationType_GroupSearchMember:
		break // 人员
	}

	return
}

// 用户会应答邀请，管理员会应答申请操作，这里需要处理并转发
func handleGroupOpRet(msg *pbmodel.Msg, session *Session) {

	groupOpMsgRet := msg.GetPlainMsg().GetGroupOpRet()
	if groupOpMsgRet == nil {
		Globals.Logger.Debug("receive wrong group op msg",
			zap.Int64("sid", session.Sid),
			zap.Int64("uid", session.UserID))
		sendBackErrorMsg(int(pbmodel.ErrorMsgType_ErrTMsgContent), "group op  is null", nil, session)
		return
	}

	// 都需要验证是否登录与权限
	ok := checkUserLogin(session)
	if !ok {
		sendBackErrorMsg(int(pbmodel.ErrorMsgType_ErrTNotLogin), "should login first.", nil, session)
		return
	}

	opCode := groupOpMsgRet.Operation
	switch opCode {
	case pbmodel.GroupOperationType_GroupInviteAnswer:
		break
	case pbmodel.GroupOperationType_GroupJoinAnswer:
		break
	default:
		sendBackErrorMsg(int(pbmodel.ErrorMsgType_ErrTWrongCode), "group op ret has a wrong opcode", nil, session)
	}
	return
}
