package core

import (
	"birdtalk/server/pbmodel"
	"birdtalk/server/utils"
	"strconv"
	"strings"
)

func handleFileUpload(msg *pbmodel.Msg, session *Session) {
	uploadMsg := msg.GetPlainMsg().GetUploadReq()
	hashType := strings.ToLower(uploadMsg.GetHashType())
	if hashType != "md5" && hashType != "sha1" {
		sendBackErrorMsg(int(pbmodel.ErrorMsgType_ErrTMsgContent), "hash type is not accepted", nil, session)
		return
	}

	if uploadMsg.HashCode == nil {
		sendBackErrorMsg(int(pbmodel.ErrorMsgType_ErrTMsgContent), "hash code can not be nil", nil, session)
		return
	}

	hashStr := utils.BytesToHexStr(uploadMsg.HashCode)
	Globals.mongoCli.FindFileByHash(hashStr)

	return
}

func handleFileDownload(msg *pbmodel.Msg, session *Session) {

}

// 计算流水号文件名
func nextFileName() string {
	id := uint64(Globals.snow.GenerateID())
	filename := strconv.FormatUint(id, 36)
	return filename
}
