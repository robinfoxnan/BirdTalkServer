package model

data class MsgPriority(override val value: Int) : pbandk.Message.Enum {
    companion object : pbandk.Message.Enum.Companion<MsgPriority> {
        val LOW = MsgPriority(0)
        val NORMAL = MsgPriority(1)
        val HIGH = MsgPriority(2)

        override fun fromValue(value: Int) = when (value) {
            0 -> LOW
            1 -> NORMAL
            2 -> HIGH
            else -> MsgPriority(value)
        }
    }
}

data class ChatMsgStatus(override val value: Int) : pbandk.Message.Enum {
    companion object : pbandk.Message.Enum.Companion<ChatMsgStatus> {
        val SENDING = ChatMsgStatus(0)
        val SENT = ChatMsgStatus(1)
        val FAILED = ChatMsgStatus(2)
        val DELIVERED = ChatMsgStatus(3)
        val READ = ChatMsgStatus(4)
        val DELETED = ChatMsgStatus(5)

        override fun fromValue(value: Int) = when (value) {
            0 -> SENDING
            1 -> SENT
            2 -> FAILED
            3 -> DELIVERED
            4 -> READ
            5 -> DELETED
            else -> ChatMsgStatus(value)
        }
    }
}

data class ChatMsgType(override val value: Int) : pbandk.Message.Enum {
    companion object : pbandk.Message.Enum.Companion<ChatMsgType> {
        val TEXT = ChatMsgType(0)
        val IMAGE = ChatMsgType(1)
        val VOICE = ChatMsgType(2)
        val VIDEO = ChatMsgType(3)
        val FILE = ChatMsgType(4)
        val DELETE = ChatMsgType(5)
        val KEY = ChatMsgType(6)
        val PLUGIN = ChatMsgType(100)

        override fun fromValue(value: Int) = when (value) {
            0 -> TEXT
            1 -> IMAGE
            2 -> VOICE
            3 -> VIDEO
            4 -> FILE
            5 -> DELETE
            6 -> KEY
            100 -> PLUGIN
            else -> ChatMsgType(value)
        }
    }
}

data class EncryptType(override val value: Int) : pbandk.Message.Enum {
    companion object : pbandk.Message.Enum.Companion<EncryptType> {
        val PLAIN = EncryptType(0)
        val CUSTOM = EncryptType(1)
        val CHACHA20 = EncryptType(2)
        val TWOFISH = EncryptType(3)
        val AES = EncryptType(4)

        override fun fromValue(value: Int) = when (value) {
            0 -> PLAIN
            1 -> CUSTOM
            2 -> CHACHA20
            3 -> TWOFISH
            4 -> AES
            else -> EncryptType(value)
        }
    }
}

data class ChatType(override val value: Int) : pbandk.Message.Enum {
    companion object : pbandk.Message.Enum.Companion<ChatType> {
        val PRIVATECHATTYPE = ChatType(0)
        val GROUPCHATTYPE = ChatType(1)

        override fun fromValue(value: Int) = when (value) {
            0 -> PRIVATECHATTYPE
            1 -> GROUPCHATTYPE
            else -> ChatType(value)
        }
    }
}

data class QueryDataType(override val value: Int) : pbandk.Message.Enum {
    companion object : pbandk.Message.Enum.Companion<QueryDataType> {
        val QUERYDATATYPEDATA = QueryDataType(0)
        val QUERYDATATYPEREPLY = QueryDataType(1)
        val QUERYDATATYPEBYTAG = QueryDataType(2)

        override fun fromValue(value: Int) = when (value) {
            0 -> QUERYDATATYPEDATA
            1 -> QUERYDATATYPEREPLY
            2 -> QUERYDATATYPEBYTAG
            else -> QueryDataType(value)
        }
    }
}

data class ErrorMsgType(override val value: Int) : pbandk.Message.Enum {
    companion object : pbandk.Message.Enum.Companion<ErrorMsgType> {
        val ERRTNONE = ErrorMsgType(0)
        val ERRTVERSION = ErrorMsgType(1)
        val ERRTKEYPRINT = ErrorMsgType(2)
        val ERRTREDIRECT = ErrorMsgType(3)
        val ERRTWRONGPWD = ErrorMsgType(4)
        val ERRTWRONGCODE = ErrorMsgType(5)

        override fun fromValue(value: Int) = when (value) {
            0 -> ERRTNONE
            1 -> ERRTVERSION
            2 -> ERRTKEYPRINT
            3 -> ERRTREDIRECT
            4 -> ERRTWRONGPWD
            5 -> ERRTWRONGCODE
            else -> ErrorMsgType(value)
        }
    }
}

data class ComMsgType(override val value: Int) : pbandk.Message.Enum {
    companion object : pbandk.Message.Enum.Companion<ComMsgType> {
        val MSGTUNUSED = ComMsgType(0)
        val MSGTHELLO = ComMsgType(1)
        val MSGTHEARTBEAT = ComMsgType(2)
        val MSGTERROR = ComMsgType(3)
        val MSGTKEYEXCHANGE = ComMsgType(4)
        val MSGTCHATMSG = ComMsgType(11)
        val MSGTCHATREPLY = ComMsgType(12)
        val MSGTQUERY = ComMsgType(13)
        val MSGTQUERYRESULT = ComMsgType(14)
        val MSGTUPLOAD = ComMsgType(21)
        val MSGTDOWNLOAD = ComMsgType(22)
        val MSGTUPLOADREPLY = ComMsgType(23)
        val MSGTDOWNLOADREPLY = ComMsgType(24)
        val MSGTUSEROP = ComMsgType(31)
        val MSGTUSEROPRET = ComMsgType(32)
        val MSGTFRIENDOP = ComMsgType(33)
        val MSGTFRIENDOPRET = ComMsgType(34)
        val MSGTGROUPOP = ComMsgType(35)
        val MSGTGROUPOPRET = ComMsgType(36)
        val MSGTOTHER = ComMsgType(100)

        override fun fromValue(value: Int) = when (value) {
            0 -> MSGTUNUSED
            1 -> MSGTHELLO
            2 -> MSGTHEARTBEAT
            3 -> MSGTERROR
            4 -> MSGTKEYEXCHANGE
            11 -> MSGTCHATMSG
            12 -> MSGTCHATREPLY
            13 -> MSGTQUERY
            14 -> MSGTQUERYRESULT
            21 -> MSGTUPLOAD
            22 -> MSGTDOWNLOAD
            23 -> MSGTUPLOADREPLY
            24 -> MSGTDOWNLOADREPLY
            31 -> MSGTUSEROP
            32 -> MSGTUSEROPRET
            33 -> MSGTFRIENDOP
            34 -> MSGTFRIENDOPRET
            35 -> MSGTGROUPOP
            36 -> MSGTGROUPOPRET
            100 -> MSGTOTHER
            else -> ComMsgType(value)
        }
    }
}

data class MsgHello(
    val clientId: String = "",
    val version: String = "",
    val platform: String = "",
    val stage: String = "",
    val keyPrint: Long = 0L,
    val rsaPrint: Long = 0L,
    val params: Map<String, String> = emptyMap(),
    val unknownFields: Map<Int, pbandk.UnknownField> = emptyMap()
) : pbandk.Message<MsgHello> {
    override operator fun plus(other: MsgHello?) = protoMergeImpl(other)
    override val protoSize by lazy { protoSizeImpl() }
    override fun protoMarshal(m: pbandk.Marshaller) = protoMarshalImpl(m)
    companion object : pbandk.Message.Companion<MsgHello> {
        override fun protoUnmarshal(u: pbandk.Unmarshaller) = MsgHello.protoUnmarshalImpl(u)
    }

    data class ParamsEntry(
        override val key: String = "",
        override val value: String = "",
        val unknownFields: Map<Int, pbandk.UnknownField> = emptyMap()
    ) : pbandk.Message<ParamsEntry>, Map.Entry<String, String> {
        override operator fun plus(other: ParamsEntry?) = protoMergeImpl(other)
        override val protoSize by lazy { protoSizeImpl() }
        override fun protoMarshal(m: pbandk.Marshaller) = protoMarshalImpl(m)
        companion object : pbandk.Message.Companion<ParamsEntry> {
            override fun protoUnmarshal(u: pbandk.Unmarshaller) = ParamsEntry.protoUnmarshalImpl(u)
        }
    }
}

data class MsgKeyExchange(
    val keyPrint: Long = 0L,
    val rsaPrint: Long = 0L,
    val stage: Int = 0,
    val tempKey: pbandk.ByteArr = pbandk.ByteArr.empty,
    val pubKey: pbandk.ByteArr = pbandk.ByteArr.empty,
    val encType: String = "",
    val status: String = "",
    val detail: String = "",
    val unknownFields: Map<Int, pbandk.UnknownField> = emptyMap()
) : pbandk.Message<MsgKeyExchange> {
    override operator fun plus(other: MsgKeyExchange?) = protoMergeImpl(other)
    override val protoSize by lazy { protoSizeImpl() }
    override fun protoMarshal(m: pbandk.Marshaller) = protoMarshalImpl(m)
    companion object : pbandk.Message.Companion<MsgKeyExchange> {
        override fun protoUnmarshal(u: pbandk.Unmarshaller) = MsgKeyExchange.protoUnmarshalImpl(u)
    }
}

data class MsgHeartBeat(
    val tm: Long = 0L,
    val userId: Long = 0L,
    val unknownFields: Map<Int, pbandk.UnknownField> = emptyMap()
) : pbandk.Message<MsgHeartBeat> {
    override operator fun plus(other: MsgHeartBeat?) = protoMergeImpl(other)
    override val protoSize by lazy { protoSizeImpl() }
    override fun protoMarshal(m: pbandk.Marshaller) = protoMarshalImpl(m)
    companion object : pbandk.Message.Companion<MsgHeartBeat> {
        override fun protoUnmarshal(u: pbandk.Unmarshaller) = MsgHeartBeat.protoUnmarshalImpl(u)
    }
}

data class MsgChat(
    val msgId: Long = 0L,
    val userId: Long = 0L,
    val fromId: Long = 0L,
    val toId: Long = 0L,
    val tm: Long = 0L,
    val devId: String = "",
    val sendId: Long = 0L,
    val msgType: model.ChatMsgType = model.ChatMsgType.fromValue(0),
    val data: pbandk.ByteArr = pbandk.ByteArr.empty,
    val priority: model.MsgPriority = model.MsgPriority.fromValue(0),
    val refMessageId: Long = 0L,
    val status: model.ChatMsgStatus = model.ChatMsgStatus.fromValue(0),
    val sendReply: Long = 0L,
    val recvReply: Long = 0L,
    val readReply: Long = 0L,
    val encType: model.EncryptType = model.EncryptType.fromValue(0),
    val chatType: model.ChatType = model.ChatType.fromValue(0),
    val subMsgType: Int = 0,
    val keyPrint: Long = 0L,
    val params: Map<String, String> = emptyMap(),
    val unknownFields: Map<Int, pbandk.UnknownField> = emptyMap()
) : pbandk.Message<MsgChat> {
    override operator fun plus(other: MsgChat?) = protoMergeImpl(other)
    override val protoSize by lazy { protoSizeImpl() }
    override fun protoMarshal(m: pbandk.Marshaller) = protoMarshalImpl(m)
    companion object : pbandk.Message.Companion<MsgChat> {
        override fun protoUnmarshal(u: pbandk.Unmarshaller) = MsgChat.protoUnmarshalImpl(u)
    }

    data class ParamsEntry(
        override val key: String = "",
        override val value: String = "",
        val unknownFields: Map<Int, pbandk.UnknownField> = emptyMap()
    ) : pbandk.Message<ParamsEntry>, Map.Entry<String, String> {
        override operator fun plus(other: ParamsEntry?) = protoMergeImpl(other)
        override val protoSize by lazy { protoSizeImpl() }
        override fun protoMarshal(m: pbandk.Marshaller) = protoMarshalImpl(m)
        companion object : pbandk.Message.Companion<ParamsEntry> {
            override fun protoUnmarshal(u: pbandk.Unmarshaller) = ParamsEntry.protoUnmarshalImpl(u)
        }
    }
}

data class MsgChatReply(
    val msgId: Long = 0L,
    val sendId: Long = 0L,
    val sendOk: Boolean = false,
    val recvOk: Boolean = false,
    val readOk: Boolean = false,
    val extraMsg: String = "",
    val userId: Long = 0L,
    val fromId: Long = 0L,
    val params: Map<String, String> = emptyMap(),
    val unknownFields: Map<Int, pbandk.UnknownField> = emptyMap()
) : pbandk.Message<MsgChatReply> {
    override operator fun plus(other: MsgChatReply?) = protoMergeImpl(other)
    override val protoSize by lazy { protoSizeImpl() }
    override fun protoMarshal(m: pbandk.Marshaller) = protoMarshalImpl(m)
    companion object : pbandk.Message.Companion<MsgChatReply> {
        override fun protoUnmarshal(u: pbandk.Unmarshaller) = MsgChatReply.protoUnmarshalImpl(u)
    }

    data class ParamsEntry(
        override val key: String = "",
        override val value: String = "",
        val unknownFields: Map<Int, pbandk.UnknownField> = emptyMap()
    ) : pbandk.Message<ParamsEntry>, Map.Entry<String, String> {
        override operator fun plus(other: ParamsEntry?) = protoMergeImpl(other)
        override val protoSize by lazy { protoSizeImpl() }
        override fun protoMarshal(m: pbandk.Marshaller) = protoMarshalImpl(m)
        companion object : pbandk.Message.Companion<ParamsEntry> {
            override fun protoUnmarshal(u: pbandk.Unmarshaller) = ParamsEntry.protoUnmarshalImpl(u)
        }
    }
}

data class MsgQuery(
    val userId: Long = 0L,
    val groupId: Long = 0L,
    val bigId: Long = 0L,
    val littleId: Long = 0L,
    val synType: Int = 0,
    val tm: Long = 0L,
    val chatType: model.ChatType = model.ChatType.fromValue(0),
    val queryType: model.QueryDataType = model.QueryDataType.fromValue(0),
    val params: Map<String, String> = emptyMap(),
    val unknownFields: Map<Int, pbandk.UnknownField> = emptyMap()
) : pbandk.Message<MsgQuery> {
    override operator fun plus(other: MsgQuery?) = protoMergeImpl(other)
    override val protoSize by lazy { protoSizeImpl() }
    override fun protoMarshal(m: pbandk.Marshaller) = protoMarshalImpl(m)
    companion object : pbandk.Message.Companion<MsgQuery> {
        override fun protoUnmarshal(u: pbandk.Unmarshaller) = MsgQuery.protoUnmarshalImpl(u)
    }

    data class ParamsEntry(
        override val key: String = "",
        override val value: String = "",
        val unknownFields: Map<Int, pbandk.UnknownField> = emptyMap()
    ) : pbandk.Message<ParamsEntry>, Map.Entry<String, String> {
        override operator fun plus(other: ParamsEntry?) = protoMergeImpl(other)
        override val protoSize by lazy { protoSizeImpl() }
        override fun protoMarshal(m: pbandk.Marshaller) = protoMarshalImpl(m)
        companion object : pbandk.Message.Companion<ParamsEntry> {
            override fun protoUnmarshal(u: pbandk.Unmarshaller) = ParamsEntry.protoUnmarshalImpl(u)
        }
    }
}

data class MsgQueryResult(
    val userId: Long = 0L,
    val anId: Long = 0L,
    val bigId: Long = 0L,
    val littleId: Long = 0L,
    val synType: Int = 0,
    val tm: Long = 0L,
    val chatType: model.ChatType = model.ChatType.fromValue(0),
    val queryType: model.QueryDataType = model.QueryDataType.fromValue(0),
    val chatDataList: List<model.MsgChat> = emptyList(),
    val chatReplyList: List<model.MsgChatReply> = emptyList(),
    val params: Map<String, String> = emptyMap(),
    val unknownFields: Map<Int, pbandk.UnknownField> = emptyMap()
) : pbandk.Message<MsgQueryResult> {
    override operator fun plus(other: MsgQueryResult?) = protoMergeImpl(other)
    override val protoSize by lazy { protoSizeImpl() }
    override fun protoMarshal(m: pbandk.Marshaller) = protoMarshalImpl(m)
    companion object : pbandk.Message.Companion<MsgQueryResult> {
        override fun protoUnmarshal(u: pbandk.Unmarshaller) = MsgQueryResult.protoUnmarshalImpl(u)
    }

    data class ParamsEntry(
        override val key: String = "",
        override val value: String = "",
        val unknownFields: Map<Int, pbandk.UnknownField> = emptyMap()
    ) : pbandk.Message<ParamsEntry>, Map.Entry<String, String> {
        override operator fun plus(other: ParamsEntry?) = protoMergeImpl(other)
        override val protoSize by lazy { protoSizeImpl() }
        override fun protoMarshal(m: pbandk.Marshaller) = protoMarshalImpl(m)
        companion object : pbandk.Message.Companion<ParamsEntry> {
            override fun protoUnmarshal(u: pbandk.Unmarshaller) = ParamsEntry.protoUnmarshalImpl(u)
        }
    }
}

data class MsgUploadReq(
    val fileName: String = "",
    val fileSize: Long = 0L,
    val fileData: pbandk.ByteArr = pbandk.ByteArr.empty,
    val hashType: String = "",
    val hashCode: pbandk.ByteArr = pbandk.ByteArr.empty,
    val fileType: String = "",
    val sendId: String = "",
    val unknownFields: Map<Int, pbandk.UnknownField> = emptyMap()
) : pbandk.Message<MsgUploadReq> {
    override operator fun plus(other: MsgUploadReq?) = protoMergeImpl(other)
    override val protoSize by lazy { protoSizeImpl() }
    override fun protoMarshal(m: pbandk.Marshaller) = protoMarshalImpl(m)
    companion object : pbandk.Message.Companion<MsgUploadReq> {
        override fun protoUnmarshal(u: pbandk.Unmarshaller) = MsgUploadReq.protoUnmarshalImpl(u)
    }
}

data class MsgUploadReply(
    val fileName: String = "",
    val sendId: String = "",
    val uuidName: String = "",
    val result: String = "",
    val unknownFields: Map<Int, pbandk.UnknownField> = emptyMap()
) : pbandk.Message<MsgUploadReply> {
    override operator fun plus(other: MsgUploadReply?) = protoMergeImpl(other)
    override val protoSize by lazy { protoSizeImpl() }
    override fun protoMarshal(m: pbandk.Marshaller) = protoMarshalImpl(m)
    companion object : pbandk.Message.Companion<MsgUploadReply> {
        override fun protoUnmarshal(u: pbandk.Unmarshaller) = MsgUploadReply.protoUnmarshalImpl(u)
    }
}

data class MsgDownloadReq(
    val sendId: String = "",
    val fileName: String = "",
    val offset: Long = 0L,
    val unknownFields: Map<Int, pbandk.UnknownField> = emptyMap()
) : pbandk.Message<MsgDownloadReq> {
    override operator fun plus(other: MsgDownloadReq?) = protoMergeImpl(other)
    override val protoSize by lazy { protoSizeImpl() }
    override fun protoMarshal(m: pbandk.Marshaller) = protoMarshalImpl(m)
    companion object : pbandk.Message.Companion<MsgDownloadReq> {
        override fun protoUnmarshal(u: pbandk.Unmarshaller) = MsgDownloadReq.protoUnmarshalImpl(u)
    }
}

data class MsgDownloadReply(
    val sendId: String = "",
    val fileName: String = "",
    val realName: String = "",
    val fileType: String = "",
    val hashType: String = "",
    val hashCode: pbandk.ByteArr = pbandk.ByteArr.empty,
    val data: pbandk.ByteArr = pbandk.ByteArr.empty,
    val size: Long = 0L,
    val offset: Long = 0L,
    val unknownFields: Map<Int, pbandk.UnknownField> = emptyMap()
) : pbandk.Message<MsgDownloadReply> {
    override operator fun plus(other: MsgDownloadReply?) = protoMergeImpl(other)
    override val protoSize by lazy { protoSizeImpl() }
    override fun protoMarshal(m: pbandk.Marshaller) = protoMarshalImpl(m)
    companion object : pbandk.Message.Companion<MsgDownloadReply> {
        override fun protoUnmarshal(u: pbandk.Unmarshaller) = MsgDownloadReply.protoUnmarshalImpl(u)
    }
}

data class MsgError(
    val code: Int = 0,
    val detail: String = "",
    val params: Map<String, String> = emptyMap(),
    val unknownFields: Map<Int, pbandk.UnknownField> = emptyMap()
) : pbandk.Message<MsgError> {
    override operator fun plus(other: MsgError?) = protoMergeImpl(other)
    override val protoSize by lazy { protoSizeImpl() }
    override fun protoMarshal(m: pbandk.Marshaller) = protoMarshalImpl(m)
    companion object : pbandk.Message.Companion<MsgError> {
        override fun protoUnmarshal(u: pbandk.Unmarshaller) = MsgError.protoUnmarshalImpl(u)
    }

    data class ParamsEntry(
        override val key: String = "",
        override val value: String = "",
        val unknownFields: Map<Int, pbandk.UnknownField> = emptyMap()
    ) : pbandk.Message<ParamsEntry>, Map.Entry<String, String> {
        override operator fun plus(other: ParamsEntry?) = protoMergeImpl(other)
        override val protoSize by lazy { protoSizeImpl() }
        override fun protoMarshal(m: pbandk.Marshaller) = protoMarshalImpl(m)
        companion object : pbandk.Message.Companion<ParamsEntry> {
            override fun protoUnmarshal(u: pbandk.Unmarshaller) = ParamsEntry.protoUnmarshalImpl(u)
        }
    }
}

data class MsgPlain(
    val message: Message? = null,
    val unknownFields: Map<Int, pbandk.UnknownField> = emptyMap()
) : pbandk.Message<MsgPlain> {
    sealed class Message {
        data class Hello(val hello: model.MsgHello) : Message()
        data class HeartBeat(val heartBeat: model.MsgHeartBeat) : Message()
        data class ErrorMsg(val errorMsg: model.MsgError) : Message()
        data class KeyEx(val keyEx: model.MsgKeyExchange) : Message()
        data class ChatData(val chatData: model.MsgChat) : Message()
        data class ChatReply(val chatReply: model.MsgChatReply) : Message()
        data class CommonQuery(val commonQuery: model.MsgQuery) : Message()
        data class CommonQueryRet(val commonQueryRet: model.MsgQueryResult) : Message()
        data class UploadReq(val uploadReq: model.MsgUploadReq) : Message()
        data class DownloadReq(val downloadReq: model.MsgDownloadReq) : Message()
        data class UploadReply(val uploadReply: model.MsgUploadReply) : Message()
        data class DownloadReply(val downloadReply: model.MsgDownloadReply) : Message()
        data class UserOp(val userOp: model.UserOpReq) : Message()
        data class UserOpRet(val userOpRet: model.UserOpResult) : Message()
        data class FriendOp(val friendOp: model.FriendOpReq) : Message()
        data class FriendOpRet(val friendOpRet: model.FriendOpResult) : Message()
        data class GroupOp(val groupOp: model.GroupOpReq) : Message()
        data class GroupOpRet(val groupOpRet: model.GroupOpResult) : Message()
        data class OtherTypeMsg(val otherTypeMsg: pbandk.ByteArr = pbandk.ByteArr.empty) : Message()
    }

    override operator fun plus(other: MsgPlain?) = protoMergeImpl(other)
    override val protoSize by lazy { protoSizeImpl() }
    override fun protoMarshal(m: pbandk.Marshaller) = protoMarshalImpl(m)
    companion object : pbandk.Message.Companion<MsgPlain> {
        override fun protoUnmarshal(u: pbandk.Unmarshaller) = MsgPlain.protoUnmarshalImpl(u)
    }
}

data class Msg(
    val version: Int = 0,
    val keyPrint: Long = 0L,
    val tm: Long = 0L,
    val msgType: model.ComMsgType = model.ComMsgType.fromValue(0),
    val subType: Int = 0,
    val message: Message_? = null,
    val unknownFields: Map<Int, pbandk.UnknownField> = emptyMap()
) : pbandk.Message<Msg> {
    sealed class Message_ {
        data class Cipher(val cipher: pbandk.ByteArr = pbandk.ByteArr.empty) : Message_()
        data class PlainMsg(val plainMsg: model.MsgPlain) : Message_()
    }

    override operator fun plus(other: Msg?) = protoMergeImpl(other)
    override val protoSize by lazy { protoSizeImpl() }
    override fun protoMarshal(m: pbandk.Marshaller) = protoMarshalImpl(m)
    companion object : pbandk.Message.Companion<Msg> {
        override fun protoUnmarshal(u: pbandk.Unmarshaller) = Msg.protoUnmarshalImpl(u)
    }
}

private fun MsgHello.protoMergeImpl(plus: MsgHello?): MsgHello = plus?.copy(
    params = params + plus.params,
    unknownFields = unknownFields + plus.unknownFields
) ?: this

private fun MsgHello.protoSizeImpl(): Int {
    var protoSize = 0
    if (clientId.isNotEmpty()) protoSize += pbandk.Sizer.tagSize(1) + pbandk.Sizer.stringSize(clientId)
    if (version.isNotEmpty()) protoSize += pbandk.Sizer.tagSize(2) + pbandk.Sizer.stringSize(version)
    if (platform.isNotEmpty()) protoSize += pbandk.Sizer.tagSize(3) + pbandk.Sizer.stringSize(platform)
    if (stage.isNotEmpty()) protoSize += pbandk.Sizer.tagSize(5) + pbandk.Sizer.stringSize(stage)
    if (keyPrint != 0L) protoSize += pbandk.Sizer.tagSize(6) + pbandk.Sizer.int64Size(keyPrint)
    if (rsaPrint != 0L) protoSize += pbandk.Sizer.tagSize(7) + pbandk.Sizer.int64Size(rsaPrint)
    if (params.isNotEmpty()) protoSize += pbandk.Sizer.mapSize(8, params, model.MsgHello::ParamsEntry)
    protoSize += unknownFields.entries.sumBy { it.value.size() }
    return protoSize
}

private fun MsgHello.protoMarshalImpl(protoMarshal: pbandk.Marshaller) {
    if (clientId.isNotEmpty()) protoMarshal.writeTag(10).writeString(clientId)
    if (version.isNotEmpty()) protoMarshal.writeTag(18).writeString(version)
    if (platform.isNotEmpty()) protoMarshal.writeTag(26).writeString(platform)
    if (stage.isNotEmpty()) protoMarshal.writeTag(42).writeString(stage)
    if (keyPrint != 0L) protoMarshal.writeTag(48).writeInt64(keyPrint)
    if (rsaPrint != 0L) protoMarshal.writeTag(56).writeInt64(rsaPrint)
    if (params.isNotEmpty()) protoMarshal.writeMap(66, params, model.MsgHello::ParamsEntry)
    if (unknownFields.isNotEmpty()) protoMarshal.writeUnknownFields(unknownFields)
}

private fun MsgHello.Companion.protoUnmarshalImpl(protoUnmarshal: pbandk.Unmarshaller): MsgHello {
    var clientId = ""
    var version = ""
    var platform = ""
    var stage = ""
    var keyPrint = 0L
    var rsaPrint = 0L
    var params: pbandk.MessageMap.Builder<String, String>? = null
    while (true) when (protoUnmarshal.readTag()) {
        0 -> return MsgHello(clientId, version, platform, stage,
            keyPrint, rsaPrint, pbandk.MessageMap.Builder.fixed(params), protoUnmarshal.unknownFields())
        10 -> clientId = protoUnmarshal.readString()
        18 -> version = protoUnmarshal.readString()
        26 -> platform = protoUnmarshal.readString()
        42 -> stage = protoUnmarshal.readString()
        48 -> keyPrint = protoUnmarshal.readInt64()
        56 -> rsaPrint = protoUnmarshal.readInt64()
        66 -> params = protoUnmarshal.readMap(params, model.MsgHello.ParamsEntry.Companion, true)
        else -> protoUnmarshal.unknownField()
    }
}

private fun MsgHello.ParamsEntry.protoMergeImpl(plus: MsgHello.ParamsEntry?): MsgHello.ParamsEntry = plus?.copy(
    unknownFields = unknownFields + plus.unknownFields
) ?: this

private fun MsgHello.ParamsEntry.protoSizeImpl(): Int {
    var protoSize = 0
    if (key.isNotEmpty()) protoSize += pbandk.Sizer.tagSize(1) + pbandk.Sizer.stringSize(key)
    if (value.isNotEmpty()) protoSize += pbandk.Sizer.tagSize(2) + pbandk.Sizer.stringSize(value)
    protoSize += unknownFields.entries.sumBy { it.value.size() }
    return protoSize
}

private fun MsgHello.ParamsEntry.protoMarshalImpl(protoMarshal: pbandk.Marshaller) {
    if (key.isNotEmpty()) protoMarshal.writeTag(10).writeString(key)
    if (value.isNotEmpty()) protoMarshal.writeTag(18).writeString(value)
    if (unknownFields.isNotEmpty()) protoMarshal.writeUnknownFields(unknownFields)
}

private fun MsgHello.ParamsEntry.Companion.protoUnmarshalImpl(protoUnmarshal: pbandk.Unmarshaller): MsgHello.ParamsEntry {
    var key = ""
    var value = ""
    while (true) when (protoUnmarshal.readTag()) {
        0 -> return MsgHello.ParamsEntry(key, value, protoUnmarshal.unknownFields())
        10 -> key = protoUnmarshal.readString()
        18 -> value = protoUnmarshal.readString()
        else -> protoUnmarshal.unknownField()
    }
}

private fun MsgKeyExchange.protoMergeImpl(plus: MsgKeyExchange?): MsgKeyExchange = plus?.copy(
    unknownFields = unknownFields + plus.unknownFields
) ?: this

private fun MsgKeyExchange.protoSizeImpl(): Int {
    var protoSize = 0
    if (keyPrint != 0L) protoSize += pbandk.Sizer.tagSize(1) + pbandk.Sizer.int64Size(keyPrint)
    if (rsaPrint != 0L) protoSize += pbandk.Sizer.tagSize(2) + pbandk.Sizer.int64Size(rsaPrint)
    if (stage != 0) protoSize += pbandk.Sizer.tagSize(3) + pbandk.Sizer.int32Size(stage)
    if (tempKey.array.isNotEmpty()) protoSize += pbandk.Sizer.tagSize(4) + pbandk.Sizer.bytesSize(tempKey)
    if (pubKey.array.isNotEmpty()) protoSize += pbandk.Sizer.tagSize(5) + pbandk.Sizer.bytesSize(pubKey)
    if (encType.isNotEmpty()) protoSize += pbandk.Sizer.tagSize(6) + pbandk.Sizer.stringSize(encType)
    if (status.isNotEmpty()) protoSize += pbandk.Sizer.tagSize(7) + pbandk.Sizer.stringSize(status)
    if (detail.isNotEmpty()) protoSize += pbandk.Sizer.tagSize(8) + pbandk.Sizer.stringSize(detail)
    protoSize += unknownFields.entries.sumBy { it.value.size() }
    return protoSize
}

private fun MsgKeyExchange.protoMarshalImpl(protoMarshal: pbandk.Marshaller) {
    if (keyPrint != 0L) protoMarshal.writeTag(8).writeInt64(keyPrint)
    if (rsaPrint != 0L) protoMarshal.writeTag(16).writeInt64(rsaPrint)
    if (stage != 0) protoMarshal.writeTag(24).writeInt32(stage)
    if (tempKey.array.isNotEmpty()) protoMarshal.writeTag(34).writeBytes(tempKey)
    if (pubKey.array.isNotEmpty()) protoMarshal.writeTag(42).writeBytes(pubKey)
    if (encType.isNotEmpty()) protoMarshal.writeTag(50).writeString(encType)
    if (status.isNotEmpty()) protoMarshal.writeTag(58).writeString(status)
    if (detail.isNotEmpty()) protoMarshal.writeTag(66).writeString(detail)
    if (unknownFields.isNotEmpty()) protoMarshal.writeUnknownFields(unknownFields)
}

private fun MsgKeyExchange.Companion.protoUnmarshalImpl(protoUnmarshal: pbandk.Unmarshaller): MsgKeyExchange {
    var keyPrint = 0L
    var rsaPrint = 0L
    var stage = 0
    var tempKey: pbandk.ByteArr = pbandk.ByteArr.empty
    var pubKey: pbandk.ByteArr = pbandk.ByteArr.empty
    var encType = ""
    var status = ""
    var detail = ""
    while (true) when (protoUnmarshal.readTag()) {
        0 -> return MsgKeyExchange(keyPrint, rsaPrint, stage, tempKey,
            pubKey, encType, status, detail, protoUnmarshal.unknownFields())
        8 -> keyPrint = protoUnmarshal.readInt64()
        16 -> rsaPrint = protoUnmarshal.readInt64()
        24 -> stage = protoUnmarshal.readInt32()
        34 -> tempKey = protoUnmarshal.readBytes()
        42 -> pubKey = protoUnmarshal.readBytes()
        50 -> encType = protoUnmarshal.readString()
        58 -> status = protoUnmarshal.readString()
        66 -> detail = protoUnmarshal.readString()
        else -> protoUnmarshal.unknownField()
    }
}

private fun MsgHeartBeat.protoMergeImpl(plus: MsgHeartBeat?): MsgHeartBeat = plus?.copy(
    unknownFields = unknownFields + plus.unknownFields
) ?: this

private fun MsgHeartBeat.protoSizeImpl(): Int {
    var protoSize = 0
    if (tm != 0L) protoSize += pbandk.Sizer.tagSize(1) + pbandk.Sizer.int64Size(tm)
    if (userId != 0L) protoSize += pbandk.Sizer.tagSize(2) + pbandk.Sizer.int64Size(userId)
    protoSize += unknownFields.entries.sumBy { it.value.size() }
    return protoSize
}

private fun MsgHeartBeat.protoMarshalImpl(protoMarshal: pbandk.Marshaller) {
    if (tm != 0L) protoMarshal.writeTag(8).writeInt64(tm)
    if (userId != 0L) protoMarshal.writeTag(16).writeInt64(userId)
    if (unknownFields.isNotEmpty()) protoMarshal.writeUnknownFields(unknownFields)
}

private fun MsgHeartBeat.Companion.protoUnmarshalImpl(protoUnmarshal: pbandk.Unmarshaller): MsgHeartBeat {
    var tm = 0L
    var userId = 0L
    while (true) when (protoUnmarshal.readTag()) {
        0 -> return MsgHeartBeat(tm, userId, protoUnmarshal.unknownFields())
        8 -> tm = protoUnmarshal.readInt64()
        16 -> userId = protoUnmarshal.readInt64()
        else -> protoUnmarshal.unknownField()
    }
}

private fun MsgChat.protoMergeImpl(plus: MsgChat?): MsgChat = plus?.copy(
    params = params + plus.params,
    unknownFields = unknownFields + plus.unknownFields
) ?: this

private fun MsgChat.protoSizeImpl(): Int {
    var protoSize = 0
    if (msgId != 0L) protoSize += pbandk.Sizer.tagSize(1) + pbandk.Sizer.int64Size(msgId)
    if (userId != 0L) protoSize += pbandk.Sizer.tagSize(2) + pbandk.Sizer.int64Size(userId)
    if (fromId != 0L) protoSize += pbandk.Sizer.tagSize(3) + pbandk.Sizer.int64Size(fromId)
    if (toId != 0L) protoSize += pbandk.Sizer.tagSize(4) + pbandk.Sizer.int64Size(toId)
    if (tm != 0L) protoSize += pbandk.Sizer.tagSize(5) + pbandk.Sizer.int64Size(tm)
    if (devId.isNotEmpty()) protoSize += pbandk.Sizer.tagSize(6) + pbandk.Sizer.stringSize(devId)
    if (sendId != 0L) protoSize += pbandk.Sizer.tagSize(7) + pbandk.Sizer.int64Size(sendId)
    if (msgType.value != 0) protoSize += pbandk.Sizer.tagSize(8) + pbandk.Sizer.enumSize(msgType)
    if (data.array.isNotEmpty()) protoSize += pbandk.Sizer.tagSize(9) + pbandk.Sizer.bytesSize(data)
    if (priority.value != 0) protoSize += pbandk.Sizer.tagSize(10) + pbandk.Sizer.enumSize(priority)
    if (refMessageId != 0L) protoSize += pbandk.Sizer.tagSize(11) + pbandk.Sizer.int64Size(refMessageId)
    if (status.value != 0) protoSize += pbandk.Sizer.tagSize(12) + pbandk.Sizer.enumSize(status)
    if (sendReply != 0L) protoSize += pbandk.Sizer.tagSize(13) + pbandk.Sizer.int64Size(sendReply)
    if (recvReply != 0L) protoSize += pbandk.Sizer.tagSize(14) + pbandk.Sizer.int64Size(recvReply)
    if (readReply != 0L) protoSize += pbandk.Sizer.tagSize(15) + pbandk.Sizer.int64Size(readReply)
    if (encType.value != 0) protoSize += pbandk.Sizer.tagSize(16) + pbandk.Sizer.enumSize(encType)
    if (chatType.value != 0) protoSize += pbandk.Sizer.tagSize(17) + pbandk.Sizer.enumSize(chatType)
    if (subMsgType != 0) protoSize += pbandk.Sizer.tagSize(18) + pbandk.Sizer.int32Size(subMsgType)
    if (keyPrint != 0L) protoSize += pbandk.Sizer.tagSize(19) + pbandk.Sizer.int64Size(keyPrint)
    if (params.isNotEmpty()) protoSize += pbandk.Sizer.mapSize(30, params, model.MsgChat::ParamsEntry)
    protoSize += unknownFields.entries.sumBy { it.value.size() }
    return protoSize
}

private fun MsgChat.protoMarshalImpl(protoMarshal: pbandk.Marshaller) {
    if (msgId != 0L) protoMarshal.writeTag(8).writeInt64(msgId)
    if (userId != 0L) protoMarshal.writeTag(16).writeInt64(userId)
    if (fromId != 0L) protoMarshal.writeTag(24).writeInt64(fromId)
    if (toId != 0L) protoMarshal.writeTag(32).writeInt64(toId)
    if (tm != 0L) protoMarshal.writeTag(40).writeInt64(tm)
    if (devId.isNotEmpty()) protoMarshal.writeTag(50).writeString(devId)
    if (sendId != 0L) protoMarshal.writeTag(56).writeInt64(sendId)
    if (msgType.value != 0) protoMarshal.writeTag(64).writeEnum(msgType)
    if (data.array.isNotEmpty()) protoMarshal.writeTag(74).writeBytes(data)
    if (priority.value != 0) protoMarshal.writeTag(80).writeEnum(priority)
    if (refMessageId != 0L) protoMarshal.writeTag(88).writeInt64(refMessageId)
    if (status.value != 0) protoMarshal.writeTag(96).writeEnum(status)
    if (sendReply != 0L) protoMarshal.writeTag(104).writeInt64(sendReply)
    if (recvReply != 0L) protoMarshal.writeTag(112).writeInt64(recvReply)
    if (readReply != 0L) protoMarshal.writeTag(120).writeInt64(readReply)
    if (encType.value != 0) protoMarshal.writeTag(128).writeEnum(encType)
    if (chatType.value != 0) protoMarshal.writeTag(136).writeEnum(chatType)
    if (subMsgType != 0) protoMarshal.writeTag(144).writeInt32(subMsgType)
    if (keyPrint != 0L) protoMarshal.writeTag(152).writeInt64(keyPrint)
    if (params.isNotEmpty()) protoMarshal.writeMap(242, params, model.MsgChat::ParamsEntry)
    if (unknownFields.isNotEmpty()) protoMarshal.writeUnknownFields(unknownFields)
}

private fun MsgChat.Companion.protoUnmarshalImpl(protoUnmarshal: pbandk.Unmarshaller): MsgChat {
    var msgId = 0L
    var userId = 0L
    var fromId = 0L
    var toId = 0L
    var tm = 0L
    var devId = ""
    var sendId = 0L
    var msgType: model.ChatMsgType = model.ChatMsgType.fromValue(0)
    var data: pbandk.ByteArr = pbandk.ByteArr.empty
    var priority: model.MsgPriority = model.MsgPriority.fromValue(0)
    var refMessageId = 0L
    var status: model.ChatMsgStatus = model.ChatMsgStatus.fromValue(0)
    var sendReply = 0L
    var recvReply = 0L
    var readReply = 0L
    var encType: model.EncryptType = model.EncryptType.fromValue(0)
    var chatType: model.ChatType = model.ChatType.fromValue(0)
    var subMsgType = 0
    var keyPrint = 0L
    var params: pbandk.MessageMap.Builder<String, String>? = null
    while (true) when (protoUnmarshal.readTag()) {
        0 -> return MsgChat(msgId, userId, fromId, toId,
            tm, devId, sendId, msgType,
            data, priority, refMessageId, status,
            sendReply, recvReply, readReply, encType,
            chatType, subMsgType, keyPrint, pbandk.MessageMap.Builder.fixed(params), protoUnmarshal.unknownFields())
        8 -> msgId = protoUnmarshal.readInt64()
        16 -> userId = protoUnmarshal.readInt64()
        24 -> fromId = protoUnmarshal.readInt64()
        32 -> toId = protoUnmarshal.readInt64()
        40 -> tm = protoUnmarshal.readInt64()
        50 -> devId = protoUnmarshal.readString()
        56 -> sendId = protoUnmarshal.readInt64()
        64 -> msgType = protoUnmarshal.readEnum(model.ChatMsgType.Companion)
        74 -> data = protoUnmarshal.readBytes()
        80 -> priority = protoUnmarshal.readEnum(model.MsgPriority.Companion)
        88 -> refMessageId = protoUnmarshal.readInt64()
        96 -> status = protoUnmarshal.readEnum(model.ChatMsgStatus.Companion)
        104 -> sendReply = protoUnmarshal.readInt64()
        112 -> recvReply = protoUnmarshal.readInt64()
        120 -> readReply = protoUnmarshal.readInt64()
        128 -> encType = protoUnmarshal.readEnum(model.EncryptType.Companion)
        136 -> chatType = protoUnmarshal.readEnum(model.ChatType.Companion)
        144 -> subMsgType = protoUnmarshal.readInt32()
        152 -> keyPrint = protoUnmarshal.readInt64()
        242 -> params = protoUnmarshal.readMap(params, model.MsgChat.ParamsEntry.Companion, true)
        else -> protoUnmarshal.unknownField()
    }
}

private fun MsgChat.ParamsEntry.protoMergeImpl(plus: MsgChat.ParamsEntry?): MsgChat.ParamsEntry = plus?.copy(
    unknownFields = unknownFields + plus.unknownFields
) ?: this

private fun MsgChat.ParamsEntry.protoSizeImpl(): Int {
    var protoSize = 0
    if (key.isNotEmpty()) protoSize += pbandk.Sizer.tagSize(1) + pbandk.Sizer.stringSize(key)
    if (value.isNotEmpty()) protoSize += pbandk.Sizer.tagSize(2) + pbandk.Sizer.stringSize(value)
    protoSize += unknownFields.entries.sumBy { it.value.size() }
    return protoSize
}

private fun MsgChat.ParamsEntry.protoMarshalImpl(protoMarshal: pbandk.Marshaller) {
    if (key.isNotEmpty()) protoMarshal.writeTag(10).writeString(key)
    if (value.isNotEmpty()) protoMarshal.writeTag(18).writeString(value)
    if (unknownFields.isNotEmpty()) protoMarshal.writeUnknownFields(unknownFields)
}

private fun MsgChat.ParamsEntry.Companion.protoUnmarshalImpl(protoUnmarshal: pbandk.Unmarshaller): MsgChat.ParamsEntry {
    var key = ""
    var value = ""
    while (true) when (protoUnmarshal.readTag()) {
        0 -> return MsgChat.ParamsEntry(key, value, protoUnmarshal.unknownFields())
        10 -> key = protoUnmarshal.readString()
        18 -> value = protoUnmarshal.readString()
        else -> protoUnmarshal.unknownField()
    }
}

private fun MsgChatReply.protoMergeImpl(plus: MsgChatReply?): MsgChatReply = plus?.copy(
    params = params + plus.params,
    unknownFields = unknownFields + plus.unknownFields
) ?: this

private fun MsgChatReply.protoSizeImpl(): Int {
    var protoSize = 0
    if (msgId != 0L) protoSize += pbandk.Sizer.tagSize(1) + pbandk.Sizer.int64Size(msgId)
    if (sendId != 0L) protoSize += pbandk.Sizer.tagSize(2) + pbandk.Sizer.int64Size(sendId)
    if (sendOk) protoSize += pbandk.Sizer.tagSize(3) + pbandk.Sizer.boolSize(sendOk)
    if (recvOk) protoSize += pbandk.Sizer.tagSize(4) + pbandk.Sizer.boolSize(recvOk)
    if (readOk) protoSize += pbandk.Sizer.tagSize(5) + pbandk.Sizer.boolSize(readOk)
    if (extraMsg.isNotEmpty()) protoSize += pbandk.Sizer.tagSize(6) + pbandk.Sizer.stringSize(extraMsg)
    if (userId != 0L) protoSize += pbandk.Sizer.tagSize(7) + pbandk.Sizer.int64Size(userId)
    if (fromId != 0L) protoSize += pbandk.Sizer.tagSize(8) + pbandk.Sizer.int64Size(fromId)
    if (params.isNotEmpty()) protoSize += pbandk.Sizer.mapSize(30, params, model.MsgChatReply::ParamsEntry)
    protoSize += unknownFields.entries.sumBy { it.value.size() }
    return protoSize
}

private fun MsgChatReply.protoMarshalImpl(protoMarshal: pbandk.Marshaller) {
    if (msgId != 0L) protoMarshal.writeTag(8).writeInt64(msgId)
    if (sendId != 0L) protoMarshal.writeTag(16).writeInt64(sendId)
    if (sendOk) protoMarshal.writeTag(24).writeBool(sendOk)
    if (recvOk) protoMarshal.writeTag(32).writeBool(recvOk)
    if (readOk) protoMarshal.writeTag(40).writeBool(readOk)
    if (extraMsg.isNotEmpty()) protoMarshal.writeTag(50).writeString(extraMsg)
    if (userId != 0L) protoMarshal.writeTag(56).writeInt64(userId)
    if (fromId != 0L) protoMarshal.writeTag(64).writeInt64(fromId)
    if (params.isNotEmpty()) protoMarshal.writeMap(242, params, model.MsgChatReply::ParamsEntry)
    if (unknownFields.isNotEmpty()) protoMarshal.writeUnknownFields(unknownFields)
}

private fun MsgChatReply.Companion.protoUnmarshalImpl(protoUnmarshal: pbandk.Unmarshaller): MsgChatReply {
    var msgId = 0L
    var sendId = 0L
    var sendOk = false
    var recvOk = false
    var readOk = false
    var extraMsg = ""
    var userId = 0L
    var fromId = 0L
    var params: pbandk.MessageMap.Builder<String, String>? = null
    while (true) when (protoUnmarshal.readTag()) {
        0 -> return MsgChatReply(msgId, sendId, sendOk, recvOk,
            readOk, extraMsg, userId, fromId,
            pbandk.MessageMap.Builder.fixed(params), protoUnmarshal.unknownFields())
        8 -> msgId = protoUnmarshal.readInt64()
        16 -> sendId = protoUnmarshal.readInt64()
        24 -> sendOk = protoUnmarshal.readBool()
        32 -> recvOk = protoUnmarshal.readBool()
        40 -> readOk = protoUnmarshal.readBool()
        50 -> extraMsg = protoUnmarshal.readString()
        56 -> userId = protoUnmarshal.readInt64()
        64 -> fromId = protoUnmarshal.readInt64()
        242 -> params = protoUnmarshal.readMap(params, model.MsgChatReply.ParamsEntry.Companion, true)
        else -> protoUnmarshal.unknownField()
    }
}

private fun MsgChatReply.ParamsEntry.protoMergeImpl(plus: MsgChatReply.ParamsEntry?): MsgChatReply.ParamsEntry = plus?.copy(
    unknownFields = unknownFields + plus.unknownFields
) ?: this

private fun MsgChatReply.ParamsEntry.protoSizeImpl(): Int {
    var protoSize = 0
    if (key.isNotEmpty()) protoSize += pbandk.Sizer.tagSize(1) + pbandk.Sizer.stringSize(key)
    if (value.isNotEmpty()) protoSize += pbandk.Sizer.tagSize(2) + pbandk.Sizer.stringSize(value)
    protoSize += unknownFields.entries.sumBy { it.value.size() }
    return protoSize
}

private fun MsgChatReply.ParamsEntry.protoMarshalImpl(protoMarshal: pbandk.Marshaller) {
    if (key.isNotEmpty()) protoMarshal.writeTag(10).writeString(key)
    if (value.isNotEmpty()) protoMarshal.writeTag(18).writeString(value)
    if (unknownFields.isNotEmpty()) protoMarshal.writeUnknownFields(unknownFields)
}

private fun MsgChatReply.ParamsEntry.Companion.protoUnmarshalImpl(protoUnmarshal: pbandk.Unmarshaller): MsgChatReply.ParamsEntry {
    var key = ""
    var value = ""
    while (true) when (protoUnmarshal.readTag()) {
        0 -> return MsgChatReply.ParamsEntry(key, value, protoUnmarshal.unknownFields())
        10 -> key = protoUnmarshal.readString()
        18 -> value = protoUnmarshal.readString()
        else -> protoUnmarshal.unknownField()
    }
}

private fun MsgQuery.protoMergeImpl(plus: MsgQuery?): MsgQuery = plus?.copy(
    params = params + plus.params,
    unknownFields = unknownFields + plus.unknownFields
) ?: this

private fun MsgQuery.protoSizeImpl(): Int {
    var protoSize = 0
    if (userId != 0L) protoSize += pbandk.Sizer.tagSize(1) + pbandk.Sizer.int64Size(userId)
    if (groupId != 0L) protoSize += pbandk.Sizer.tagSize(2) + pbandk.Sizer.int64Size(groupId)
    if (bigId != 0L) protoSize += pbandk.Sizer.tagSize(3) + pbandk.Sizer.int64Size(bigId)
    if (littleId != 0L) protoSize += pbandk.Sizer.tagSize(4) + pbandk.Sizer.int64Size(littleId)
    if (synType != 0) protoSize += pbandk.Sizer.tagSize(5) + pbandk.Sizer.int32Size(synType)
    if (tm != 0L) protoSize += pbandk.Sizer.tagSize(6) + pbandk.Sizer.int64Size(tm)
    if (chatType.value != 0) protoSize += pbandk.Sizer.tagSize(7) + pbandk.Sizer.enumSize(chatType)
    if (queryType.value != 0) protoSize += pbandk.Sizer.tagSize(8) + pbandk.Sizer.enumSize(queryType)
    if (params.isNotEmpty()) protoSize += pbandk.Sizer.mapSize(9, params, model.MsgQuery::ParamsEntry)
    protoSize += unknownFields.entries.sumBy { it.value.size() }
    return protoSize
}

private fun MsgQuery.protoMarshalImpl(protoMarshal: pbandk.Marshaller) {
    if (userId != 0L) protoMarshal.writeTag(8).writeInt64(userId)
    if (groupId != 0L) protoMarshal.writeTag(16).writeInt64(groupId)
    if (bigId != 0L) protoMarshal.writeTag(24).writeInt64(bigId)
    if (littleId != 0L) protoMarshal.writeTag(32).writeInt64(littleId)
    if (synType != 0) protoMarshal.writeTag(40).writeInt32(synType)
    if (tm != 0L) protoMarshal.writeTag(48).writeInt64(tm)
    if (chatType.value != 0) protoMarshal.writeTag(56).writeEnum(chatType)
    if (queryType.value != 0) protoMarshal.writeTag(64).writeEnum(queryType)
    if (params.isNotEmpty()) protoMarshal.writeMap(74, params, model.MsgQuery::ParamsEntry)
    if (unknownFields.isNotEmpty()) protoMarshal.writeUnknownFields(unknownFields)
}

private fun MsgQuery.Companion.protoUnmarshalImpl(protoUnmarshal: pbandk.Unmarshaller): MsgQuery {
    var userId = 0L
    var groupId = 0L
    var bigId = 0L
    var littleId = 0L
    var synType = 0
    var tm = 0L
    var chatType: model.ChatType = model.ChatType.fromValue(0)
    var queryType: model.QueryDataType = model.QueryDataType.fromValue(0)
    var params: pbandk.MessageMap.Builder<String, String>? = null
    while (true) when (protoUnmarshal.readTag()) {
        0 -> return MsgQuery(userId, groupId, bigId, littleId,
            synType, tm, chatType, queryType,
            pbandk.MessageMap.Builder.fixed(params), protoUnmarshal.unknownFields())
        8 -> userId = protoUnmarshal.readInt64()
        16 -> groupId = protoUnmarshal.readInt64()
        24 -> bigId = protoUnmarshal.readInt64()
        32 -> littleId = protoUnmarshal.readInt64()
        40 -> synType = protoUnmarshal.readInt32()
        48 -> tm = protoUnmarshal.readInt64()
        56 -> chatType = protoUnmarshal.readEnum(model.ChatType.Companion)
        64 -> queryType = protoUnmarshal.readEnum(model.QueryDataType.Companion)
        74 -> params = protoUnmarshal.readMap(params, model.MsgQuery.ParamsEntry.Companion, true)
        else -> protoUnmarshal.unknownField()
    }
}

private fun MsgQuery.ParamsEntry.protoMergeImpl(plus: MsgQuery.ParamsEntry?): MsgQuery.ParamsEntry = plus?.copy(
    unknownFields = unknownFields + plus.unknownFields
) ?: this

private fun MsgQuery.ParamsEntry.protoSizeImpl(): Int {
    var protoSize = 0
    if (key.isNotEmpty()) protoSize += pbandk.Sizer.tagSize(1) + pbandk.Sizer.stringSize(key)
    if (value.isNotEmpty()) protoSize += pbandk.Sizer.tagSize(2) + pbandk.Sizer.stringSize(value)
    protoSize += unknownFields.entries.sumBy { it.value.size() }
    return protoSize
}

private fun MsgQuery.ParamsEntry.protoMarshalImpl(protoMarshal: pbandk.Marshaller) {
    if (key.isNotEmpty()) protoMarshal.writeTag(10).writeString(key)
    if (value.isNotEmpty()) protoMarshal.writeTag(18).writeString(value)
    if (unknownFields.isNotEmpty()) protoMarshal.writeUnknownFields(unknownFields)
}

private fun MsgQuery.ParamsEntry.Companion.protoUnmarshalImpl(protoUnmarshal: pbandk.Unmarshaller): MsgQuery.ParamsEntry {
    var key = ""
    var value = ""
    while (true) when (protoUnmarshal.readTag()) {
        0 -> return MsgQuery.ParamsEntry(key, value, protoUnmarshal.unknownFields())
        10 -> key = protoUnmarshal.readString()
        18 -> value = protoUnmarshal.readString()
        else -> protoUnmarshal.unknownField()
    }
}

private fun MsgQueryResult.protoMergeImpl(plus: MsgQueryResult?): MsgQueryResult = plus?.copy(
    chatDataList = chatDataList + plus.chatDataList,
    chatReplyList = chatReplyList + plus.chatReplyList,
    params = params + plus.params,
    unknownFields = unknownFields + plus.unknownFields
) ?: this

private fun MsgQueryResult.protoSizeImpl(): Int {
    var protoSize = 0
    if (userId != 0L) protoSize += pbandk.Sizer.tagSize(1) + pbandk.Sizer.int64Size(userId)
    if (anId != 0L) protoSize += pbandk.Sizer.tagSize(2) + pbandk.Sizer.int64Size(anId)
    if (bigId != 0L) protoSize += pbandk.Sizer.tagSize(3) + pbandk.Sizer.int64Size(bigId)
    if (littleId != 0L) protoSize += pbandk.Sizer.tagSize(4) + pbandk.Sizer.int64Size(littleId)
    if (synType != 0) protoSize += pbandk.Sizer.tagSize(5) + pbandk.Sizer.int32Size(synType)
    if (tm != 0L) protoSize += pbandk.Sizer.tagSize(6) + pbandk.Sizer.int64Size(tm)
    if (chatType.value != 0) protoSize += pbandk.Sizer.tagSize(7) + pbandk.Sizer.enumSize(chatType)
    if (queryType.value != 0) protoSize += pbandk.Sizer.tagSize(8) + pbandk.Sizer.enumSize(queryType)
    if (chatDataList.isNotEmpty()) protoSize += (pbandk.Sizer.tagSize(11) * chatDataList.size) + chatDataList.sumBy(pbandk.Sizer::messageSize)
    if (chatReplyList.isNotEmpty()) protoSize += (pbandk.Sizer.tagSize(12) * chatReplyList.size) + chatReplyList.sumBy(pbandk.Sizer::messageSize)
    if (params.isNotEmpty()) protoSize += pbandk.Sizer.mapSize(13, params, model.MsgQueryResult::ParamsEntry)
    protoSize += unknownFields.entries.sumBy { it.value.size() }
    return protoSize
}

private fun MsgQueryResult.protoMarshalImpl(protoMarshal: pbandk.Marshaller) {
    if (userId != 0L) protoMarshal.writeTag(8).writeInt64(userId)
    if (anId != 0L) protoMarshal.writeTag(16).writeInt64(anId)
    if (bigId != 0L) protoMarshal.writeTag(24).writeInt64(bigId)
    if (littleId != 0L) protoMarshal.writeTag(32).writeInt64(littleId)
    if (synType != 0) protoMarshal.writeTag(40).writeInt32(synType)
    if (tm != 0L) protoMarshal.writeTag(48).writeInt64(tm)
    if (chatType.value != 0) protoMarshal.writeTag(56).writeEnum(chatType)
    if (queryType.value != 0) protoMarshal.writeTag(64).writeEnum(queryType)
    if (chatDataList.isNotEmpty()) chatDataList.forEach { protoMarshal.writeTag(90).writeMessage(it) }
    if (chatReplyList.isNotEmpty()) chatReplyList.forEach { protoMarshal.writeTag(98).writeMessage(it) }
    if (params.isNotEmpty()) protoMarshal.writeMap(106, params, model.MsgQueryResult::ParamsEntry)
    if (unknownFields.isNotEmpty()) protoMarshal.writeUnknownFields(unknownFields)
}

private fun MsgQueryResult.Companion.protoUnmarshalImpl(protoUnmarshal: pbandk.Unmarshaller): MsgQueryResult {
    var userId = 0L
    var anId = 0L
    var bigId = 0L
    var littleId = 0L
    var synType = 0
    var tm = 0L
    var chatType: model.ChatType = model.ChatType.fromValue(0)
    var queryType: model.QueryDataType = model.QueryDataType.fromValue(0)
    var chatDataList: pbandk.ListWithSize.Builder<model.MsgChat>? = null
    var chatReplyList: pbandk.ListWithSize.Builder<model.MsgChatReply>? = null
    var params: pbandk.MessageMap.Builder<String, String>? = null
    while (true) when (protoUnmarshal.readTag()) {
        0 -> return MsgQueryResult(userId, anId, bigId, littleId,
            synType, tm, chatType, queryType,
            pbandk.ListWithSize.Builder.fixed(chatDataList), pbandk.ListWithSize.Builder.fixed(chatReplyList), pbandk.MessageMap.Builder.fixed(params), protoUnmarshal.unknownFields())
        8 -> userId = protoUnmarshal.readInt64()
        16 -> anId = protoUnmarshal.readInt64()
        24 -> bigId = protoUnmarshal.readInt64()
        32 -> littleId = protoUnmarshal.readInt64()
        40 -> synType = protoUnmarshal.readInt32()
        48 -> tm = protoUnmarshal.readInt64()
        56 -> chatType = protoUnmarshal.readEnum(model.ChatType.Companion)
        64 -> queryType = protoUnmarshal.readEnum(model.QueryDataType.Companion)
        90 -> chatDataList = protoUnmarshal.readRepeatedMessage(chatDataList, model.MsgChat.Companion, true)
        98 -> chatReplyList = protoUnmarshal.readRepeatedMessage(chatReplyList, model.MsgChatReply.Companion, true)
        106 -> params = protoUnmarshal.readMap(params, model.MsgQueryResult.ParamsEntry.Companion, true)
        else -> protoUnmarshal.unknownField()
    }
}

private fun MsgQueryResult.ParamsEntry.protoMergeImpl(plus: MsgQueryResult.ParamsEntry?): MsgQueryResult.ParamsEntry = plus?.copy(
    unknownFields = unknownFields + plus.unknownFields
) ?: this

private fun MsgQueryResult.ParamsEntry.protoSizeImpl(): Int {
    var protoSize = 0
    if (key.isNotEmpty()) protoSize += pbandk.Sizer.tagSize(1) + pbandk.Sizer.stringSize(key)
    if (value.isNotEmpty()) protoSize += pbandk.Sizer.tagSize(2) + pbandk.Sizer.stringSize(value)
    protoSize += unknownFields.entries.sumBy { it.value.size() }
    return protoSize
}

private fun MsgQueryResult.ParamsEntry.protoMarshalImpl(protoMarshal: pbandk.Marshaller) {
    if (key.isNotEmpty()) protoMarshal.writeTag(10).writeString(key)
    if (value.isNotEmpty()) protoMarshal.writeTag(18).writeString(value)
    if (unknownFields.isNotEmpty()) protoMarshal.writeUnknownFields(unknownFields)
}

private fun MsgQueryResult.ParamsEntry.Companion.protoUnmarshalImpl(protoUnmarshal: pbandk.Unmarshaller): MsgQueryResult.ParamsEntry {
    var key = ""
    var value = ""
    while (true) when (protoUnmarshal.readTag()) {
        0 -> return MsgQueryResult.ParamsEntry(key, value, protoUnmarshal.unknownFields())
        10 -> key = protoUnmarshal.readString()
        18 -> value = protoUnmarshal.readString()
        else -> protoUnmarshal.unknownField()
    }
}

private fun MsgUploadReq.protoMergeImpl(plus: MsgUploadReq?): MsgUploadReq = plus?.copy(
    unknownFields = unknownFields + plus.unknownFields
) ?: this

private fun MsgUploadReq.protoSizeImpl(): Int {
    var protoSize = 0
    if (fileName.isNotEmpty()) protoSize += pbandk.Sizer.tagSize(1) + pbandk.Sizer.stringSize(fileName)
    if (fileSize != 0L) protoSize += pbandk.Sizer.tagSize(2) + pbandk.Sizer.int64Size(fileSize)
    if (fileData.array.isNotEmpty()) protoSize += pbandk.Sizer.tagSize(3) + pbandk.Sizer.bytesSize(fileData)
    if (hashType.isNotEmpty()) protoSize += pbandk.Sizer.tagSize(4) + pbandk.Sizer.stringSize(hashType)
    if (hashCode.array.isNotEmpty()) protoSize += pbandk.Sizer.tagSize(5) + pbandk.Sizer.bytesSize(hashCode)
    if (fileType.isNotEmpty()) protoSize += pbandk.Sizer.tagSize(6) + pbandk.Sizer.stringSize(fileType)
    if (sendId.isNotEmpty()) protoSize += pbandk.Sizer.tagSize(7) + pbandk.Sizer.stringSize(sendId)
    protoSize += unknownFields.entries.sumBy { it.value.size() }
    return protoSize
}

private fun MsgUploadReq.protoMarshalImpl(protoMarshal: pbandk.Marshaller) {
    if (fileName.isNotEmpty()) protoMarshal.writeTag(10).writeString(fileName)
    if (fileSize != 0L) protoMarshal.writeTag(16).writeInt64(fileSize)
    if (fileData.array.isNotEmpty()) protoMarshal.writeTag(26).writeBytes(fileData)
    if (hashType.isNotEmpty()) protoMarshal.writeTag(34).writeString(hashType)
    if (hashCode.array.isNotEmpty()) protoMarshal.writeTag(42).writeBytes(hashCode)
    if (fileType.isNotEmpty()) protoMarshal.writeTag(50).writeString(fileType)
    if (sendId.isNotEmpty()) protoMarshal.writeTag(58).writeString(sendId)
    if (unknownFields.isNotEmpty()) protoMarshal.writeUnknownFields(unknownFields)
}

private fun MsgUploadReq.Companion.protoUnmarshalImpl(protoUnmarshal: pbandk.Unmarshaller): MsgUploadReq {
    var fileName = ""
    var fileSize = 0L
    var fileData: pbandk.ByteArr = pbandk.ByteArr.empty
    var hashType = ""
    var hashCode: pbandk.ByteArr = pbandk.ByteArr.empty
    var fileType = ""
    var sendId = ""
    while (true) when (protoUnmarshal.readTag()) {
        0 -> return MsgUploadReq(fileName, fileSize, fileData, hashType,
            hashCode, fileType, sendId, protoUnmarshal.unknownFields())
        10 -> fileName = protoUnmarshal.readString()
        16 -> fileSize = protoUnmarshal.readInt64()
        26 -> fileData = protoUnmarshal.readBytes()
        34 -> hashType = protoUnmarshal.readString()
        42 -> hashCode = protoUnmarshal.readBytes()
        50 -> fileType = protoUnmarshal.readString()
        58 -> sendId = protoUnmarshal.readString()
        else -> protoUnmarshal.unknownField()
    }
}

private fun MsgUploadReply.protoMergeImpl(plus: MsgUploadReply?): MsgUploadReply = plus?.copy(
    unknownFields = unknownFields + plus.unknownFields
) ?: this

private fun MsgUploadReply.protoSizeImpl(): Int {
    var protoSize = 0
    if (fileName.isNotEmpty()) protoSize += pbandk.Sizer.tagSize(1) + pbandk.Sizer.stringSize(fileName)
    if (sendId.isNotEmpty()) protoSize += pbandk.Sizer.tagSize(2) + pbandk.Sizer.stringSize(sendId)
    if (uuidName.isNotEmpty()) protoSize += pbandk.Sizer.tagSize(3) + pbandk.Sizer.stringSize(uuidName)
    if (result.isNotEmpty()) protoSize += pbandk.Sizer.tagSize(4) + pbandk.Sizer.stringSize(result)
    protoSize += unknownFields.entries.sumBy { it.value.size() }
    return protoSize
}

private fun MsgUploadReply.protoMarshalImpl(protoMarshal: pbandk.Marshaller) {
    if (fileName.isNotEmpty()) protoMarshal.writeTag(10).writeString(fileName)
    if (sendId.isNotEmpty()) protoMarshal.writeTag(18).writeString(sendId)
    if (uuidName.isNotEmpty()) protoMarshal.writeTag(26).writeString(uuidName)
    if (result.isNotEmpty()) protoMarshal.writeTag(34).writeString(result)
    if (unknownFields.isNotEmpty()) protoMarshal.writeUnknownFields(unknownFields)
}

private fun MsgUploadReply.Companion.protoUnmarshalImpl(protoUnmarshal: pbandk.Unmarshaller): MsgUploadReply {
    var fileName = ""
    var sendId = ""
    var uuidName = ""
    var result = ""
    while (true) when (protoUnmarshal.readTag()) {
        0 -> return MsgUploadReply(fileName, sendId, uuidName, result, protoUnmarshal.unknownFields())
        10 -> fileName = protoUnmarshal.readString()
        18 -> sendId = protoUnmarshal.readString()
        26 -> uuidName = protoUnmarshal.readString()
        34 -> result = protoUnmarshal.readString()
        else -> protoUnmarshal.unknownField()
    }
}

private fun MsgDownloadReq.protoMergeImpl(plus: MsgDownloadReq?): MsgDownloadReq = plus?.copy(
    unknownFields = unknownFields + plus.unknownFields
) ?: this

private fun MsgDownloadReq.protoSizeImpl(): Int {
    var protoSize = 0
    if (sendId.isNotEmpty()) protoSize += pbandk.Sizer.tagSize(1) + pbandk.Sizer.stringSize(sendId)
    if (fileName.isNotEmpty()) protoSize += pbandk.Sizer.tagSize(2) + pbandk.Sizer.stringSize(fileName)
    if (offset != 0L) protoSize += pbandk.Sizer.tagSize(3) + pbandk.Sizer.int64Size(offset)
    protoSize += unknownFields.entries.sumBy { it.value.size() }
    return protoSize
}

private fun MsgDownloadReq.protoMarshalImpl(protoMarshal: pbandk.Marshaller) {
    if (sendId.isNotEmpty()) protoMarshal.writeTag(10).writeString(sendId)
    if (fileName.isNotEmpty()) protoMarshal.writeTag(18).writeString(fileName)
    if (offset != 0L) protoMarshal.writeTag(24).writeInt64(offset)
    if (unknownFields.isNotEmpty()) protoMarshal.writeUnknownFields(unknownFields)
}

private fun MsgDownloadReq.Companion.protoUnmarshalImpl(protoUnmarshal: pbandk.Unmarshaller): MsgDownloadReq {
    var sendId = ""
    var fileName = ""
    var offset = 0L
    while (true) when (protoUnmarshal.readTag()) {
        0 -> return MsgDownloadReq(sendId, fileName, offset, protoUnmarshal.unknownFields())
        10 -> sendId = protoUnmarshal.readString()
        18 -> fileName = protoUnmarshal.readString()
        24 -> offset = protoUnmarshal.readInt64()
        else -> protoUnmarshal.unknownField()
    }
}

private fun MsgDownloadReply.protoMergeImpl(plus: MsgDownloadReply?): MsgDownloadReply = plus?.copy(
    unknownFields = unknownFields + plus.unknownFields
) ?: this

private fun MsgDownloadReply.protoSizeImpl(): Int {
    var protoSize = 0
    if (sendId.isNotEmpty()) protoSize += pbandk.Sizer.tagSize(1) + pbandk.Sizer.stringSize(sendId)
    if (fileName.isNotEmpty()) protoSize += pbandk.Sizer.tagSize(2) + pbandk.Sizer.stringSize(fileName)
    if (realName.isNotEmpty()) protoSize += pbandk.Sizer.tagSize(3) + pbandk.Sizer.stringSize(realName)
    if (fileType.isNotEmpty()) protoSize += pbandk.Sizer.tagSize(4) + pbandk.Sizer.stringSize(fileType)
    if (hashType.isNotEmpty()) protoSize += pbandk.Sizer.tagSize(5) + pbandk.Sizer.stringSize(hashType)
    if (hashCode.array.isNotEmpty()) protoSize += pbandk.Sizer.tagSize(6) + pbandk.Sizer.bytesSize(hashCode)
    if (data.array.isNotEmpty()) protoSize += pbandk.Sizer.tagSize(7) + pbandk.Sizer.bytesSize(data)
    if (size != 0L) protoSize += pbandk.Sizer.tagSize(8) + pbandk.Sizer.int64Size(size)
    if (offset != 0L) protoSize += pbandk.Sizer.tagSize(9) + pbandk.Sizer.int64Size(offset)
    protoSize += unknownFields.entries.sumBy { it.value.size() }
    return protoSize
}

private fun MsgDownloadReply.protoMarshalImpl(protoMarshal: pbandk.Marshaller) {
    if (sendId.isNotEmpty()) protoMarshal.writeTag(10).writeString(sendId)
    if (fileName.isNotEmpty()) protoMarshal.writeTag(18).writeString(fileName)
    if (realName.isNotEmpty()) protoMarshal.writeTag(26).writeString(realName)
    if (fileType.isNotEmpty()) protoMarshal.writeTag(34).writeString(fileType)
    if (hashType.isNotEmpty()) protoMarshal.writeTag(42).writeString(hashType)
    if (hashCode.array.isNotEmpty()) protoMarshal.writeTag(50).writeBytes(hashCode)
    if (data.array.isNotEmpty()) protoMarshal.writeTag(58).writeBytes(data)
    if (size != 0L) protoMarshal.writeTag(64).writeInt64(size)
    if (offset != 0L) protoMarshal.writeTag(72).writeInt64(offset)
    if (unknownFields.isNotEmpty()) protoMarshal.writeUnknownFields(unknownFields)
}

private fun MsgDownloadReply.Companion.protoUnmarshalImpl(protoUnmarshal: pbandk.Unmarshaller): MsgDownloadReply {
    var sendId = ""
    var fileName = ""
    var realName = ""
    var fileType = ""
    var hashType = ""
    var hashCode: pbandk.ByteArr = pbandk.ByteArr.empty
    var data: pbandk.ByteArr = pbandk.ByteArr.empty
    var size = 0L
    var offset = 0L
    while (true) when (protoUnmarshal.readTag()) {
        0 -> return MsgDownloadReply(sendId, fileName, realName, fileType,
            hashType, hashCode, data, size,
            offset, protoUnmarshal.unknownFields())
        10 -> sendId = protoUnmarshal.readString()
        18 -> fileName = protoUnmarshal.readString()
        26 -> realName = protoUnmarshal.readString()
        34 -> fileType = protoUnmarshal.readString()
        42 -> hashType = protoUnmarshal.readString()
        50 -> hashCode = protoUnmarshal.readBytes()
        58 -> data = protoUnmarshal.readBytes()
        64 -> size = protoUnmarshal.readInt64()
        72 -> offset = protoUnmarshal.readInt64()
        else -> protoUnmarshal.unknownField()
    }
}

private fun MsgError.protoMergeImpl(plus: MsgError?): MsgError = plus?.copy(
    params = params + plus.params,
    unknownFields = unknownFields + plus.unknownFields
) ?: this

private fun MsgError.protoSizeImpl(): Int {
    var protoSize = 0
    if (code != 0) protoSize += pbandk.Sizer.tagSize(1) + pbandk.Sizer.int32Size(code)
    if (detail.isNotEmpty()) protoSize += pbandk.Sizer.tagSize(2) + pbandk.Sizer.stringSize(detail)
    if (params.isNotEmpty()) protoSize += pbandk.Sizer.mapSize(9, params, model.MsgError::ParamsEntry)
    protoSize += unknownFields.entries.sumBy { it.value.size() }
    return protoSize
}

private fun MsgError.protoMarshalImpl(protoMarshal: pbandk.Marshaller) {
    if (code != 0) protoMarshal.writeTag(8).writeInt32(code)
    if (detail.isNotEmpty()) protoMarshal.writeTag(18).writeString(detail)
    if (params.isNotEmpty()) protoMarshal.writeMap(74, params, model.MsgError::ParamsEntry)
    if (unknownFields.isNotEmpty()) protoMarshal.writeUnknownFields(unknownFields)
}

private fun MsgError.Companion.protoUnmarshalImpl(protoUnmarshal: pbandk.Unmarshaller): MsgError {
    var code = 0
    var detail = ""
    var params: pbandk.MessageMap.Builder<String, String>? = null
    while (true) when (protoUnmarshal.readTag()) {
        0 -> return MsgError(code, detail, pbandk.MessageMap.Builder.fixed(params), protoUnmarshal.unknownFields())
        8 -> code = protoUnmarshal.readInt32()
        18 -> detail = protoUnmarshal.readString()
        74 -> params = protoUnmarshal.readMap(params, model.MsgError.ParamsEntry.Companion, true)
        else -> protoUnmarshal.unknownField()
    }
}

private fun MsgError.ParamsEntry.protoMergeImpl(plus: MsgError.ParamsEntry?): MsgError.ParamsEntry = plus?.copy(
    unknownFields = unknownFields + plus.unknownFields
) ?: this

private fun MsgError.ParamsEntry.protoSizeImpl(): Int {
    var protoSize = 0
    if (key.isNotEmpty()) protoSize += pbandk.Sizer.tagSize(1) + pbandk.Sizer.stringSize(key)
    if (value.isNotEmpty()) protoSize += pbandk.Sizer.tagSize(2) + pbandk.Sizer.stringSize(value)
    protoSize += unknownFields.entries.sumBy { it.value.size() }
    return protoSize
}

private fun MsgError.ParamsEntry.protoMarshalImpl(protoMarshal: pbandk.Marshaller) {
    if (key.isNotEmpty()) protoMarshal.writeTag(10).writeString(key)
    if (value.isNotEmpty()) protoMarshal.writeTag(18).writeString(value)
    if (unknownFields.isNotEmpty()) protoMarshal.writeUnknownFields(unknownFields)
}

private fun MsgError.ParamsEntry.Companion.protoUnmarshalImpl(protoUnmarshal: pbandk.Unmarshaller): MsgError.ParamsEntry {
    var key = ""
    var value = ""
    while (true) when (protoUnmarshal.readTag()) {
        0 -> return MsgError.ParamsEntry(key, value, protoUnmarshal.unknownFields())
        10 -> key = protoUnmarshal.readString()
        18 -> value = protoUnmarshal.readString()
        else -> protoUnmarshal.unknownField()
    }
}

private fun MsgPlain.protoMergeImpl(plus: MsgPlain?): MsgPlain = plus?.copy(
    message = when {
        message is MsgPlain.Message.Hello && plus.message is MsgPlain.Message.Hello ->
            MsgPlain.Message.Hello(message.hello + plus.message.hello)
        message is MsgPlain.Message.HeartBeat && plus.message is MsgPlain.Message.HeartBeat ->
            MsgPlain.Message.HeartBeat(message.heartBeat + plus.message.heartBeat)
        message is MsgPlain.Message.ErrorMsg && plus.message is MsgPlain.Message.ErrorMsg ->
            MsgPlain.Message.ErrorMsg(message.errorMsg + plus.message.errorMsg)
        message is MsgPlain.Message.KeyEx && plus.message is MsgPlain.Message.KeyEx ->
            MsgPlain.Message.KeyEx(message.keyEx + plus.message.keyEx)
        message is MsgPlain.Message.ChatData && plus.message is MsgPlain.Message.ChatData ->
            MsgPlain.Message.ChatData(message.chatData + plus.message.chatData)
        message is MsgPlain.Message.ChatReply && plus.message is MsgPlain.Message.ChatReply ->
            MsgPlain.Message.ChatReply(message.chatReply + plus.message.chatReply)
        message is MsgPlain.Message.CommonQuery && plus.message is MsgPlain.Message.CommonQuery ->
            MsgPlain.Message.CommonQuery(message.commonQuery + plus.message.commonQuery)
        message is MsgPlain.Message.CommonQueryRet && plus.message is MsgPlain.Message.CommonQueryRet ->
            MsgPlain.Message.CommonQueryRet(message.commonQueryRet + plus.message.commonQueryRet)
        message is MsgPlain.Message.UploadReq && plus.message is MsgPlain.Message.UploadReq ->
            MsgPlain.Message.UploadReq(message.uploadReq + plus.message.uploadReq)
        message is MsgPlain.Message.DownloadReq && plus.message is MsgPlain.Message.DownloadReq ->
            MsgPlain.Message.DownloadReq(message.downloadReq + plus.message.downloadReq)
        message is MsgPlain.Message.UploadReply && plus.message is MsgPlain.Message.UploadReply ->
            MsgPlain.Message.UploadReply(message.uploadReply + plus.message.uploadReply)
        message is MsgPlain.Message.DownloadReply && plus.message is MsgPlain.Message.DownloadReply ->
            MsgPlain.Message.DownloadReply(message.downloadReply + plus.message.downloadReply)
        message is MsgPlain.Message.UserOp && plus.message is MsgPlain.Message.UserOp ->
            MsgPlain.Message.UserOp(message.userOp + plus.message.userOp)
        message is MsgPlain.Message.UserOpRet && plus.message is MsgPlain.Message.UserOpRet ->
            MsgPlain.Message.UserOpRet(message.userOpRet + plus.message.userOpRet)
        message is MsgPlain.Message.FriendOp && plus.message is MsgPlain.Message.FriendOp ->
            MsgPlain.Message.FriendOp(message.friendOp + plus.message.friendOp)
        message is MsgPlain.Message.FriendOpRet && plus.message is MsgPlain.Message.FriendOpRet ->
            MsgPlain.Message.FriendOpRet(message.friendOpRet + plus.message.friendOpRet)
        message is MsgPlain.Message.GroupOp && plus.message is MsgPlain.Message.GroupOp ->
            MsgPlain.Message.GroupOp(message.groupOp + plus.message.groupOp)
        message is MsgPlain.Message.GroupOpRet && plus.message is MsgPlain.Message.GroupOpRet ->
            MsgPlain.Message.GroupOpRet(message.groupOpRet + plus.message.groupOpRet)
        else ->
            plus.message ?: message
    },
    unknownFields = unknownFields + plus.unknownFields
) ?: this

private fun MsgPlain.protoSizeImpl(): Int {
    var protoSize = 0
    when (message) {
        is MsgPlain.Message.Hello -> protoSize += pbandk.Sizer.tagSize(7) + pbandk.Sizer.messageSize(message.hello)
        is MsgPlain.Message.HeartBeat -> protoSize += pbandk.Sizer.tagSize(8) + pbandk.Sizer.messageSize(message.heartBeat)
        is MsgPlain.Message.ErrorMsg -> protoSize += pbandk.Sizer.tagSize(9) + pbandk.Sizer.messageSize(message.errorMsg)
        is MsgPlain.Message.KeyEx -> protoSize += pbandk.Sizer.tagSize(10) + pbandk.Sizer.messageSize(message.keyEx)
        is MsgPlain.Message.ChatData -> protoSize += pbandk.Sizer.tagSize(11) + pbandk.Sizer.messageSize(message.chatData)
        is MsgPlain.Message.ChatReply -> protoSize += pbandk.Sizer.tagSize(12) + pbandk.Sizer.messageSize(message.chatReply)
        is MsgPlain.Message.CommonQuery -> protoSize += pbandk.Sizer.tagSize(13) + pbandk.Sizer.messageSize(message.commonQuery)
        is MsgPlain.Message.CommonQueryRet -> protoSize += pbandk.Sizer.tagSize(14) + pbandk.Sizer.messageSize(message.commonQueryRet)
        is MsgPlain.Message.UploadReq -> protoSize += pbandk.Sizer.tagSize(21) + pbandk.Sizer.messageSize(message.uploadReq)
        is MsgPlain.Message.DownloadReq -> protoSize += pbandk.Sizer.tagSize(22) + pbandk.Sizer.messageSize(message.downloadReq)
        is MsgPlain.Message.UploadReply -> protoSize += pbandk.Sizer.tagSize(23) + pbandk.Sizer.messageSize(message.uploadReply)
        is MsgPlain.Message.DownloadReply -> protoSize += pbandk.Sizer.tagSize(24) + pbandk.Sizer.messageSize(message.downloadReply)
        is MsgPlain.Message.UserOp -> protoSize += pbandk.Sizer.tagSize(31) + pbandk.Sizer.messageSize(message.userOp)
        is MsgPlain.Message.UserOpRet -> protoSize += pbandk.Sizer.tagSize(32) + pbandk.Sizer.messageSize(message.userOpRet)
        is MsgPlain.Message.FriendOp -> protoSize += pbandk.Sizer.tagSize(33) + pbandk.Sizer.messageSize(message.friendOp)
        is MsgPlain.Message.FriendOpRet -> protoSize += pbandk.Sizer.tagSize(34) + pbandk.Sizer.messageSize(message.friendOpRet)
        is MsgPlain.Message.GroupOp -> protoSize += pbandk.Sizer.tagSize(35) + pbandk.Sizer.messageSize(message.groupOp)
        is MsgPlain.Message.GroupOpRet -> protoSize += pbandk.Sizer.tagSize(36) + pbandk.Sizer.messageSize(message.groupOpRet)
        is MsgPlain.Message.OtherTypeMsg -> protoSize += pbandk.Sizer.tagSize(100) + pbandk.Sizer.bytesSize(message.otherTypeMsg)
    }
    protoSize += unknownFields.entries.sumBy { it.value.size() }
    return protoSize
}

private fun MsgPlain.protoMarshalImpl(protoMarshal: pbandk.Marshaller) {
    if (message is MsgPlain.Message.Hello) protoMarshal.writeTag(58).writeMessage(message.hello)
    if (message is MsgPlain.Message.HeartBeat) protoMarshal.writeTag(66).writeMessage(message.heartBeat)
    if (message is MsgPlain.Message.ErrorMsg) protoMarshal.writeTag(74).writeMessage(message.errorMsg)
    if (message is MsgPlain.Message.KeyEx) protoMarshal.writeTag(82).writeMessage(message.keyEx)
    if (message is MsgPlain.Message.ChatData) protoMarshal.writeTag(90).writeMessage(message.chatData)
    if (message is MsgPlain.Message.ChatReply) protoMarshal.writeTag(98).writeMessage(message.chatReply)
    if (message is MsgPlain.Message.CommonQuery) protoMarshal.writeTag(106).writeMessage(message.commonQuery)
    if (message is MsgPlain.Message.CommonQueryRet) protoMarshal.writeTag(114).writeMessage(message.commonQueryRet)
    if (message is MsgPlain.Message.UploadReq) protoMarshal.writeTag(170).writeMessage(message.uploadReq)
    if (message is MsgPlain.Message.DownloadReq) protoMarshal.writeTag(178).writeMessage(message.downloadReq)
    if (message is MsgPlain.Message.UploadReply) protoMarshal.writeTag(186).writeMessage(message.uploadReply)
    if (message is MsgPlain.Message.DownloadReply) protoMarshal.writeTag(194).writeMessage(message.downloadReply)
    if (message is MsgPlain.Message.UserOp) protoMarshal.writeTag(250).writeMessage(message.userOp)
    if (message is MsgPlain.Message.UserOpRet) protoMarshal.writeTag(258).writeMessage(message.userOpRet)
    if (message is MsgPlain.Message.FriendOp) protoMarshal.writeTag(266).writeMessage(message.friendOp)
    if (message is MsgPlain.Message.FriendOpRet) protoMarshal.writeTag(274).writeMessage(message.friendOpRet)
    if (message is MsgPlain.Message.GroupOp) protoMarshal.writeTag(282).writeMessage(message.groupOp)
    if (message is MsgPlain.Message.GroupOpRet) protoMarshal.writeTag(290).writeMessage(message.groupOpRet)
    if (message is MsgPlain.Message.OtherTypeMsg) protoMarshal.writeTag(802).writeBytes(message.otherTypeMsg)
    if (unknownFields.isNotEmpty()) protoMarshal.writeUnknownFields(unknownFields)
}

private fun MsgPlain.Companion.protoUnmarshalImpl(protoUnmarshal: pbandk.Unmarshaller): MsgPlain {
    var message: MsgPlain.Message? = null
    while (true) when (protoUnmarshal.readTag()) {
        0 -> return MsgPlain(message, protoUnmarshal.unknownFields())
        58 -> message = MsgPlain.Message.Hello(protoUnmarshal.readMessage(model.MsgHello.Companion))
        66 -> message = MsgPlain.Message.HeartBeat(protoUnmarshal.readMessage(model.MsgHeartBeat.Companion))
        74 -> message = MsgPlain.Message.ErrorMsg(protoUnmarshal.readMessage(model.MsgError.Companion))
        82 -> message = MsgPlain.Message.KeyEx(protoUnmarshal.readMessage(model.MsgKeyExchange.Companion))
        90 -> message = MsgPlain.Message.ChatData(protoUnmarshal.readMessage(model.MsgChat.Companion))
        98 -> message = MsgPlain.Message.ChatReply(protoUnmarshal.readMessage(model.MsgChatReply.Companion))
        106 -> message = MsgPlain.Message.CommonQuery(protoUnmarshal.readMessage(model.MsgQuery.Companion))
        114 -> message = MsgPlain.Message.CommonQueryRet(protoUnmarshal.readMessage(model.MsgQueryResult.Companion))
        170 -> message = MsgPlain.Message.UploadReq(protoUnmarshal.readMessage(model.MsgUploadReq.Companion))
        178 -> message = MsgPlain.Message.DownloadReq(protoUnmarshal.readMessage(model.MsgDownloadReq.Companion))
        186 -> message = MsgPlain.Message.UploadReply(protoUnmarshal.readMessage(model.MsgUploadReply.Companion))
        194 -> message = MsgPlain.Message.DownloadReply(protoUnmarshal.readMessage(model.MsgDownloadReply.Companion))
        250 -> message = MsgPlain.Message.UserOp(protoUnmarshal.readMessage(model.UserOpReq.Companion))
        258 -> message = MsgPlain.Message.UserOpRet(protoUnmarshal.readMessage(model.UserOpResult.Companion))
        266 -> message = MsgPlain.Message.FriendOp(protoUnmarshal.readMessage(model.FriendOpReq.Companion))
        274 -> message = MsgPlain.Message.FriendOpRet(protoUnmarshal.readMessage(model.FriendOpResult.Companion))
        282 -> message = MsgPlain.Message.GroupOp(protoUnmarshal.readMessage(model.GroupOpReq.Companion))
        290 -> message = MsgPlain.Message.GroupOpRet(protoUnmarshal.readMessage(model.GroupOpResult.Companion))
        802 -> message = MsgPlain.Message.OtherTypeMsg(protoUnmarshal.readBytes())
        else -> protoUnmarshal.unknownField()
    }
}

private fun Msg.protoMergeImpl(plus: Msg?): Msg = plus?.copy(
    message = when {
        message is Msg.Message_.PlainMsg && plus.message is Msg.Message_.PlainMsg ->
            Msg.Message_.PlainMsg(message.plainMsg + plus.message.plainMsg)
        else ->
            plus.message ?: message
    },
    unknownFields = unknownFields + plus.unknownFields
) ?: this

private fun Msg.protoSizeImpl(): Int {
    var protoSize = 0
    if (version != 0) protoSize += pbandk.Sizer.tagSize(1) + pbandk.Sizer.int32Size(version)
    if (keyPrint != 0L) protoSize += pbandk.Sizer.tagSize(2) + pbandk.Sizer.int64Size(keyPrint)
    if (tm != 0L) protoSize += pbandk.Sizer.tagSize(3) + pbandk.Sizer.int64Size(tm)
    if (msgType.value != 0) protoSize += pbandk.Sizer.tagSize(4) + pbandk.Sizer.enumSize(msgType)
    if (subType != 0) protoSize += pbandk.Sizer.tagSize(5) + pbandk.Sizer.int32Size(subType)
    when (message) {
        is Msg.Message_.Cipher -> protoSize += pbandk.Sizer.tagSize(11) + pbandk.Sizer.bytesSize(message.cipher)
        is Msg.Message_.PlainMsg -> protoSize += pbandk.Sizer.tagSize(12) + pbandk.Sizer.messageSize(message.plainMsg)
    }
    protoSize += unknownFields.entries.sumBy { it.value.size() }
    return protoSize
}

private fun Msg.protoMarshalImpl(protoMarshal: pbandk.Marshaller) {
    if (version != 0) protoMarshal.writeTag(8).writeInt32(version)
    if (keyPrint != 0L) protoMarshal.writeTag(16).writeInt64(keyPrint)
    if (tm != 0L) protoMarshal.writeTag(24).writeInt64(tm)
    if (msgType.value != 0) protoMarshal.writeTag(32).writeEnum(msgType)
    if (subType != 0) protoMarshal.writeTag(40).writeInt32(subType)
    if (message is Msg.Message_.Cipher) protoMarshal.writeTag(90).writeBytes(message.cipher)
    if (message is Msg.Message_.PlainMsg) protoMarshal.writeTag(98).writeMessage(message.plainMsg)
    if (unknownFields.isNotEmpty()) protoMarshal.writeUnknownFields(unknownFields)
}

private fun Msg.Companion.protoUnmarshalImpl(protoUnmarshal: pbandk.Unmarshaller): Msg {
    var version = 0
    var keyPrint = 0L
    var tm = 0L
    var msgType: model.ComMsgType = model.ComMsgType.fromValue(0)
    var subType = 0
    var message: Msg.Message_? = null
    while (true) when (protoUnmarshal.readTag()) {
        0 -> return Msg(version, keyPrint, tm, msgType,
            subType, message, protoUnmarshal.unknownFields())
        8 -> version = protoUnmarshal.readInt32()
        16 -> keyPrint = protoUnmarshal.readInt64()
        24 -> tm = protoUnmarshal.readInt64()
        32 -> msgType = protoUnmarshal.readEnum(model.ComMsgType.Companion)
        40 -> subType = protoUnmarshal.readInt32()
        90 -> message = Msg.Message_.Cipher(protoUnmarshal.readBytes())
        98 -> message = Msg.Message_.PlainMsg(protoUnmarshal.readMessage(model.MsgPlain.Companion))
        else -> protoUnmarshal.unknownField()
    }
}
