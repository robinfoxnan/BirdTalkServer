package core

import (
	"birdtalk/server/model"
	"birdtalk/server/pbmodel"
	"birdtalk/server/utils"
	"crypto/md5"
	"encoding/binary"
	"encoding/hex"
	"errors"
	"fmt"
	"go.uber.org/zap"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"sync"
)

func handleFileUpload(msg *pbmodel.Msg, session *Session) {
	uploadMsg := msg.GetPlainMsg().GetUploadReq()

	Globals.Logger.Debug("File Upload msg:", zap.Int64("user id", session.UserID),
		zap.String("filename", uploadMsg.FileName), zap.String("hash", uploadMsg.HashCode), zap.Int64("send id", uploadMsg.SendId))

	if uploadMsg == nil {
		sendBackErrorMsg(int(pbmodel.ErrorMsgType_ErrTMsgContent), "upload request is null", nil, session)
		return
	}
	hashType := strings.ToLower(uploadMsg.GetHashType())
	if hashType != "md5" && hashType != "sha1" {
		sendBackErrorMsg(int(pbmodel.ErrorMsgType_ErrTMsgContent), "hash type is not accepted", nil, session)
		return
	}

	if len(uploadMsg.HashCode) < 16 {
		sendBackErrorMsg(int(pbmodel.ErrorMsgType_ErrTMsgContent), "hash code can not be nil", nil, session)
		return
	}

	if uploadMsg.FileData == nil {
		{
			sendBackErrorMsg(int(pbmodel.ErrorMsgType_ErrTMsgContent), "data can not be nil", nil, session)
			return
		}
	}

	// 第一片
	if uploadMsg.ChunkIndex == 0 {
		onHandleUploadTrunkFirst(uploadMsg, session)
	} else {
		onHandleUploadTrunkOther(uploadMsg, session)
	}

	return
}

// 下载的请求
func handleFileDownload(msg *pbmodel.Msg, session *Session) {
	downLoadReq := msg.GetPlainMsg().GetDownloadReq()
	if downLoadReq == nil {
		sendBackErrorMsg(int(pbmodel.ErrorMsgType_ErrTMsgContent), "download request is null", nil, session)
		return
	}

	fileFullPath, _ := utils.FileName2FilePath(Globals.Config.Server.FileBasePath, downLoadReq.GetFileName(), false)

	st, err := os.Stat(fileFullPath)
	if os.IsNotExist(err) {

		msgRet := createDownloadRetMsg(downLoadReq, "", "fail", "not exist",
			0, 0, 0, 0, nil, "")
		session.SendMessage(msgRet)
		return
	}
	sz := st.Size()
	if sz < (1<<20)*2 {
		sendbackFile(downLoadReq, sz, fileFullPath, session)
	} else {
		go sendbackFile(downLoadReq, sz, fileFullPath, session)
	}

}

// 应答文件数据
func sendbackFile(downLoadReq *pbmodel.MsgDownloadReq, sz int64, fileFullPath string, session *Session) {

	file, err := os.Open(fileFullPath)
	if err != nil {
		msgRet := createDownloadRetMsg(downLoadReq, "", "fail", "open file error",
			0, 0, 0, 0, nil, "")
		session.SendMessage(msgRet)
		return
	}
	defer file.Close()

	chSize := int64(1 << 20)
	chCount := (sz + chSize) / int64(chSize)
	chIndex := 0

	buffer := make([]byte, chSize)

	for {
		n, err := file.Read(buffer)
		if err != nil {
			msgRet := createDownloadRetMsg(downLoadReq, "", "fail", "open file error",
				0, 0, 0, 0, nil, "")
			session.SendMessage(msgRet)
			return
		}

		if n > 0 {
			msgRet := createDownloadRetMsg(downLoadReq, "", "trunk", "",
				sz, int32(chIndex), int32(chCount), int32(chSize), buffer, "")
			session.SendMessage(msgRet)
		}

		if int64(n) < chSize {
			break
		}

		chIndex++
	}
	msgRet := createDownloadRetMsg(downLoadReq, "", "finish", "",
		sz, int32(chIndex), int32(chCount), int32(chSize), buffer, "")
	session.SendMessage(msgRet)
}

// 上传的消息的应答
func createUpLoadRetMsg(uniqName string, uploadMsg *pbmodel.MsgUploadReq, result, detail string) *pbmodel.Msg {

	msgUploadRet := pbmodel.MsgUploadReply{
		Result:     result,
		FileName:   uploadMsg.FileName,
		SendId:     uploadMsg.SendId,
		UuidName:   uniqName,
		ChunkIndex: uploadMsg.ChunkIndex,
		Detail:     detail,
	}

	msgPlain := pbmodel.MsgPlain{
		Message: &pbmodel.MsgPlain_UploadReply{
			UploadReply: &msgUploadRet,
		},
	}

	msg := pbmodel.Msg{
		Version:  int32(ProtocolVersion),
		KeyPrint: 0,
		Tm:       utils.GetTimeStamp(),
		MsgType:  pbmodel.ComMsgType_MsgTUploadReply,
		SubType:  0,
		Message: &pbmodel.Msg_PlainMsg{
			PlainMsg: &msgPlain,
		},
	}

	return &msg
}

// 创建下载回复的消息
func createDownloadRetMsg(downloadRet *pbmodel.MsgDownloadReq, fileName, result, detail string,
	sz int64, index, chCount, chSize int32, data []byte, hashCode string) *pbmodel.Msg {

	msgDownloadRet := pbmodel.MsgDownloadReply{
		Result:     result,
		Detail:     detail,
		FileName:   downloadRet.FileName,
		RealName:   fileName,
		Offset:     0,
		Size:       sz,
		ChunkIndex: index,
		ChunkCount: chCount,
		ChunkSize:  chSize,
		Data:       data,
		HashCode:   hashCode,
		HashType:   "md5",
	}

	msgPlain := pbmodel.MsgPlain{
		Message: &pbmodel.MsgPlain_DownloadReply{
			DownloadReply: &msgDownloadRet,
		},
	}

	msg := pbmodel.Msg{
		Version:  int32(ProtocolVersion),
		KeyPrint: 0,
		Tm:       utils.GetTimeStamp(),
		MsgType:  pbmodel.ComMsgType_MsgTDownloadReply,
		SubType:  0,
		Message: &pbmodel.Msg_PlainMsg{
			PlainMsg: &msgPlain,
		},
	}

	return &msg
}

// 计算流水号文件名
// 示例哈希字符串
// hashString := "b58f6e861b4f82a6f8a86b0d1b646216be2f5489d70595cb91c4a2c97a18ff67"
func getUUIDFileName(fileName string, fileHash string) (string, error) {
	ext := filepath.Ext(fileName)
	// 将哈希字符串转换为字节数组
	id := uint64(0)
	hashBytes, err := hex.DecodeString(fileHash)
	if err != nil {
		id = uint64(Globals.snow.GenerateID())
	} else {
		// 使用 binary.LittleEndian 解析字节数组为 int64
		id = uint64(binary.LittleEndian.Uint64(hashBytes[:8])) // 取前8字节转换为 int64
	}

	idName := strconv.FormatUint(id, 36)

	return idName + ext, err
}

func sendBackFileUploadErr(uniqName string, uploadMsg *pbmodel.MsgUploadReq, result string, detail string, session *Session) {
	Globals.Logger.Error(detail)
	msgRet := createUpLoadRetMsg(uniqName, uploadMsg, result, detail)
	session.SendMessage(msgRet)
}

// 文件名，存在否，一样否
func checkFileExist(fileName string) (bool, int64) {
	filePath, err := utils.FileName2FilePath(Globals.Config.Server.FileBasePath, fileName, true)
	if err != nil {
		return false, 0
	}
	// 检测文件是否存在
	st, err1 := os.Stat(filePath)
	if os.IsNotExist(err1) {
		return false, 0
	}

	return true, st.Size()
}

// 处理第一个分片
// 创建文件，但是不写库
// todo: 将文件放在临时列表中，会话如果出错了，则需要删除
func onHandleUploadTrunkFirst(uploadMsg *pbmodel.MsgUploadReq, session *Session) {

	// 检查是否能秒传
	hashStr := uploadMsg.HashCode
	fileStore, _ := Globals.mongoCli.FindFileByHashCode(hashStr)
	if fileStore != nil {
		// 防止哈希冲突，以及防止那个文件被清除了
		b, sz := checkFileExist(fileStore.UniqName)
		if b {
			if sz == fileStore.FileSize {
				// 添加记录
				fileInfo := model.FileInfo{
					HashCode:  hashStr,
					StoreType: "dir",
					FileName:  uploadMsg.FileName,
					UniqName:  fileStore.UniqName,
					Gid:       uploadMsg.GroupId,
					Status:    "",
					Tm:        utils.GetTimeStamp(),
					FileSize:  uploadMsg.FileSize,
					UserId:    session.UserID,
					Tags:      Globals.segment.Cut(filepath.Base(uploadMsg.FileName)), // 对主文件名分词
				}
				err := Globals.mongoCli.SaveNewFile(&fileInfo)
				if err != nil {
					sendBackErrorMsg(int(pbmodel.ErrorMsgType_ErrTServerInside), "save file info err", nil, session)
					return

				}

				msgRet := createUpLoadRetMsg(fileStore.UniqName, uploadMsg, "sameok", "find same one")
				session.SendMessage(msgRet)
				return
			} else {
				Globals.Logger.Fatal("file hashcode same ,size not same")
			}
		}
	}

	idName, err := getUUIDFileName(uploadMsg.FileName, uploadMsg.HashCode)
	if err != nil {
		Globals.Logger.Error("hash is invalid")
	}
	// 会话重建立一个结构
	sFile := &SessionFile{
		File:       nil,
		isUpload:   true,
		FileName:   uploadMsg.FileName,
		FullPath:   "",
		Gid:        uploadMsg.GroupId,
		UniqName:   idName,
		HashCode:   uploadMsg.HashCode,
		FileSize:   uploadMsg.FileSize,
		Hash:       md5.New(),
		Lock:       sync.Mutex{},
		ChunkSize:  uploadMsg.ChunkSize,
		ChunkIndex: 0,
		ChunkCount: uploadMsg.ChunkCount,
	}

	fullPath, err := utils.FileName2FilePath(Globals.Config.Server.FileBasePath, sFile.UniqName, true)
	if err != nil {
		sendBackFileUploadErr("", uploadMsg, "fail", "save file info err", session)
		return
	}
	sFile.FullPath = fullPath
	fmt.Println(fullPath)

	sFile.File, err = os.Create(fullPath)
	if err != nil {
		cleanSFile(sFile, session)
		sendBackFileUploadErr("", uploadMsg, "fail", "create file  err", session)
		return
	}

	// 写文件
	err = writeToFile(uploadMsg.FileData, sFile)
	if err != nil {
		cleanSFile(sFile, session)
		sendBackFileUploadErr("", uploadMsg, "fail", "write file data err", session)
		return
	}

	// 最后一片
	if uploadMsg.GetChunkCount() == 1 {
		// 计算MD5
		hashInBytes := sFile.Hash.Sum(nil)
		hashString := hex.EncodeToString(hashInBytes)
		fmt.Printf("本地计算的md5 = %s, 对方的=%s \n", hashString, sFile.HashCode)
		if hashString != sFile.HashCode {
			cleanSFile(sFile, session)
			Globals.Logger.Error("md5 not same", zap.String("file", hashString), zap.String("remote", sFile.HashCode))
			sendBackFileUploadErr("", uploadMsg, "fail", "md5 hash code is not same", session)
			return
		}

		closeSFile(sFile, session)
		msgRet := createUpLoadRetMsg(sFile.UniqName, uploadMsg, "fileok", "finish")
		session.SendMessage(msgRet)

	} else {
		// 保存到
		session.SetFile(uploadMsg.FileName, sFile)
		msgRet := createUpLoadRetMsg("", uploadMsg, "chunkok", "wait next truck")
		session.SendMessage(msgRet)
	}

}

func closeSFile(sFile *SessionFile, session *Session) {
	if sFile.File != nil {
		sFile.File.Close()
		sFile.File = nil
	}
	// 保存到数据库

	fileInfo := model.FileInfo{
		HashCode:  sFile.HashCode,
		StoreType: "dir",
		FileName:  sFile.FileName,
		UniqName:  sFile.UniqName,
		Gid:       sFile.Gid,
		Status:    "",
		Tm:        utils.GetTimeStamp(),
		FileSize:  sFile.FileSize,
		UserId:    session.UserID,
		Tags:      Globals.segment.Cut(filepath.Base(sFile.FileName)), // 对主文件名分词
	}
	err := Globals.mongoCli.SaveNewFile(&fileInfo)
	if err != nil {
		//sendBackErrorMsg(int(pbmodel.ErrorMsgType_ErrTServerInside), "save file info err", nil, session)
	}

	session.RemoveFile(sFile.FileName)
}

// 如果中间过程出错，则删除文件，删除缓存
func cleanSFile(sFile *SessionFile, session *Session) {
	if sFile.File != nil {
		sFile.File.Close()
		sFile.File = nil
	}
	// 尝试删除文件
	err := os.Remove(sFile.FullPath)
	if err != nil {
		fmt.Printf("无法删除文件：%s\n", err)
		return
	}
	session.RemoveFile(sFile.FileName)
}
func writeToFile(data []byte, sFile *SessionFile) error {
	// 将数据写入文件
	if sFile.File == nil {
		return errors.New("nil file")
	}
	_, err := sFile.File.Write(data)
	if err != nil {
		return err
	}

	//fmt.Println("数据写入文件成功")
	if _, err := sFile.Hash.Write(data); err != nil {
		fmt.Printf("Error writing to hash: %v\n", err)

	}

	return nil
}

// 如果文件过小，可能只有1片
func onHandleUploadTrunkOther(uploadMsg *pbmodel.MsgUploadReq, session *Session) {
	sFile := session.GetFile(uploadMsg.FileName)
	if sFile == nil {
		sendBackFileUploadErr("", uploadMsg, "fail", "should send truck 0 first", session)
		return
	}

	bLast := false
	if uploadMsg.ChunkIndex == (uploadMsg.ChunkCount - 1) {
		bLast = true
	}

	// 写文件
	err := writeToFile(uploadMsg.FileData, sFile)
	if err != nil {
		cleanSFile(sFile, session)
		sendBackFileUploadErr("", uploadMsg, "fail", "write file data err", session)
		return
	}

	if bLast {
		// 计算MD5
		hashInBytes := sFile.Hash.Sum(nil)
		hashString := hex.EncodeToString(hashInBytes)
		fmt.Printf("接收的md5 = %s \n", hashString)
		if hashString != sFile.HashCode {
			cleanSFile(sFile, session)
			Globals.Logger.Error("md5 not same", zap.String("file", hashString), zap.String("remote", sFile.HashCode))
			sendBackFileUploadErr("", uploadMsg, "fail", "md5 hash code is not same", session)
			return
		}

		closeSFile(sFile, session)
		msgRet := createUpLoadRetMsg(sFile.UniqName, uploadMsg, "fileok", "finish")
		session.SendMessage(msgRet)
	} else {
		msgRet := createUpLoadRetMsg("", uploadMsg, "chunkok", "wait next")
		session.SendMessage(msgRet)
	}

}
