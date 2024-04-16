package model

data class UserOperationType(override val value: Int) : pbandk.Message.Enum {
    companion object : pbandk.Message.Enum.Companion<UserOperationType> {
        val USERNONEACTION = UserOperationType(0)
        val REGISTERUSER = UserOperationType(1)
        val UNREGISTERUSER = UserOperationType(2)
        val DISABLEUSER = UserOperationType(3)
        val RECOVERUSER = UserOperationType(4)
        val SETUSERINFO = UserOperationType(5)
        val REALNAMEVERIFICATION = UserOperationType(6)
        val LOGIN = UserOperationType(7)
        val LOGOUT = UserOperationType(8)
        val FINDUSER = UserOperationType(9)
        val ADDFRIEND = UserOperationType(10)
        val APPROVEFRIEND = UserOperationType(11)
        val REMOVEFRIEND = UserOperationType(12)
        val BLOCKFRIEND = UserOperationType(13)
        val UNBLOCKFRIEND = UserOperationType(14)
        val SETFRIENDPERMISSION = UserOperationType(15)
        val SETFRIENDMEMO = UserOperationType(16)

        override fun fromValue(value: Int) = when (value) {
            0 -> USERNONEACTION
            1 -> REGISTERUSER
            2 -> UNREGISTERUSER
            3 -> DISABLEUSER
            4 -> RECOVERUSER
            5 -> SETUSERINFO
            6 -> REALNAMEVERIFICATION
            7 -> LOGIN
            8 -> LOGOUT
            9 -> FINDUSER
            10 -> ADDFRIEND
            11 -> APPROVEFRIEND
            12 -> REMOVEFRIEND
            13 -> BLOCKFRIEND
            14 -> UNBLOCKFRIEND
            15 -> SETFRIENDPERMISSION
            16 -> SETFRIENDMEMO
            else -> UserOperationType(value)
        }
    }
}

data class GroupOperationType(override val value: Int) : pbandk.Message.Enum {
    companion object : pbandk.Message.Enum.Companion<GroupOperationType> {
        val GROUPNONEACTION = GroupOperationType(0)
        val GROUPCREATE = GroupOperationType(1)
        val GROUPDISSOLVE = GroupOperationType(2)
        val GROUPSETINFO = GroupOperationType(3)
        val GROUPKICKMEMBER = GroupOperationType(4)
        val GROUPINVITEREQUEST = GroupOperationType(5)
        val GROUPINVITEANSWER = GroupOperationType(6)
        val GROUPJOINREQUEST = GroupOperationType(7)
        val GROUPJOINANSWER = GroupOperationType(8)
        val GROUPQUIT = GroupOperationType(9)
        val GROUPADDADMIN = GroupOperationType(10)
        val GROUPDELADMIN = GroupOperationType(11)
        val GROUPTRANSFEROWNER = GroupOperationType(12)
        val GROUPSETMEMBERINFO = GroupOperationType(13)
        val GROUPSEARCH = GroupOperationType(14)
        val GROUPSEARCHMEMBER = GroupOperationType(15)

        override fun fromValue(value: Int) = when (value) {
            0 -> GROUPNONEACTION
            1 -> GROUPCREATE
            2 -> GROUPDISSOLVE
            3 -> GROUPSETINFO
            4 -> GROUPKICKMEMBER
            5 -> GROUPINVITEREQUEST
            6 -> GROUPINVITEANSWER
            7 -> GROUPJOINREQUEST
            8 -> GROUPJOINANSWER
            9 -> GROUPQUIT
            10 -> GROUPADDADMIN
            11 -> GROUPDELADMIN
            12 -> GROUPTRANSFEROWNER
            13 -> GROUPSETMEMBERINFO
            14 -> GROUPSEARCH
            15 -> GROUPSEARCHMEMBER
            else -> GroupOperationType(value)
        }
    }
}

data class GroupOperationResultType(override val value: Int) : pbandk.Message.Enum {
    companion object : pbandk.Message.Enum.Companion<GroupOperationResultType> {
        val GROUPOPERATIONRESULTNONE = GroupOperationResultType(0)
        val GROUPOPERATIONRESULTOK = GroupOperationResultType(1)
        val GROUPOPERATIONRESULTREFUSE = GroupOperationResultType(2)

        override fun fromValue(value: Int) = when (value) {
            0 -> GROUPOPERATIONRESULTNONE
            1 -> GROUPOPERATIONRESULTOK
            2 -> GROUPOPERATIONRESULTREFUSE
            else -> GroupOperationResultType(value)
        }
    }
}

data class UserInfo(
    val userId: Long = 0L,
    val userName: String = "",
    val nickName: String = "",
    val email: String = "",
    val phone: String = "",
    val gender: String = "",
    val age: Int = 0,
    val region: String = "",
    val icon: String = "",
    val params: Map<String, String> = emptyMap(),
    val unknownFields: Map<Int, pbandk.UnknownField> = emptyMap()
) : pbandk.Message<UserInfo> {
    override operator fun plus(other: UserInfo?) = protoMergeImpl(other)
    override val protoSize by lazy { protoSizeImpl() }
    override fun protoMarshal(m: pbandk.Marshaller) = protoMarshalImpl(m)
    companion object : pbandk.Message.Companion<UserInfo> {
        override fun protoUnmarshal(u: pbandk.Unmarshaller) = UserInfo.protoUnmarshalImpl(u)
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

data class UserOpReq(
    val operation: model.UserOperationType = model.UserOperationType.fromValue(0),
    val user: model.UserInfo? = null,
    val params: Map<String, String> = emptyMap(),
    val unknownFields: Map<Int, pbandk.UnknownField> = emptyMap()
) : pbandk.Message<UserOpReq> {
    override operator fun plus(other: UserOpReq?) = protoMergeImpl(other)
    override val protoSize by lazy { protoSizeImpl() }
    override fun protoMarshal(m: pbandk.Marshaller) = protoMarshalImpl(m)
    companion object : pbandk.Message.Companion<UserOpReq> {
        override fun protoUnmarshal(u: pbandk.Unmarshaller) = UserOpReq.protoUnmarshalImpl(u)
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

data class UserOpResult(
    val operation: model.UserOperationType = model.UserOperationType.fromValue(0),
    val result: String = "",
    val users: List<model.UserInfo> = emptyList(),
    val params: Map<String, String> = emptyMap(),
    val unknownFields: Map<Int, pbandk.UnknownField> = emptyMap()
) : pbandk.Message<UserOpResult> {
    override operator fun plus(other: UserOpResult?) = protoMergeImpl(other)
    override val protoSize by lazy { protoSizeImpl() }
    override fun protoMarshal(m: pbandk.Marshaller) = protoMarshalImpl(m)
    companion object : pbandk.Message.Companion<UserOpResult> {
        override fun protoUnmarshal(u: pbandk.Unmarshaller) = UserOpResult.protoUnmarshalImpl(u)
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

data class FriendOpReq(
    val operation: model.UserOperationType = model.UserOperationType.fromValue(0),
    val user: model.UserInfo? = null,
    val sendId: Long = 0L,
    val msgId: Long = 0L,
    val params: Map<String, String> = emptyMap(),
    val unknownFields: Map<Int, pbandk.UnknownField> = emptyMap()
) : pbandk.Message<FriendOpReq> {
    override operator fun plus(other: FriendOpReq?) = protoMergeImpl(other)
    override val protoSize by lazy { protoSizeImpl() }
    override fun protoMarshal(m: pbandk.Marshaller) = protoMarshalImpl(m)
    companion object : pbandk.Message.Companion<FriendOpReq> {
        override fun protoUnmarshal(u: pbandk.Unmarshaller) = FriendOpReq.protoUnmarshalImpl(u)
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

data class FriendOpResult(
    val operation: model.UserOperationType = model.UserOperationType.fromValue(0),
    val result: String = "",
    val user: model.UserInfo? = null,
    val users: model.UserInfo? = null,
    val sendId: Long = 0L,
    val msgId: Long = 0L,
    val params: Map<String, String> = emptyMap(),
    val unknownFields: Map<Int, pbandk.UnknownField> = emptyMap()
) : pbandk.Message<FriendOpResult> {
    override operator fun plus(other: FriendOpResult?) = protoMergeImpl(other)
    override val protoSize by lazy { protoSizeImpl() }
    override fun protoMarshal(m: pbandk.Marshaller) = protoMarshalImpl(m)
    companion object : pbandk.Message.Companion<FriendOpResult> {
        override fun protoUnmarshal(u: pbandk.Unmarshaller) = FriendOpResult.protoUnmarshalImpl(u)
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

data class GroupMember(
    val userId: Long = 0L,
    val nick: String = "",
    val icon: String = "",
    val role: String = "",
    val groupId: Long = 0L,
    val params: Map<String, String> = emptyMap(),
    val unknownFields: Map<Int, pbandk.UnknownField> = emptyMap()
) : pbandk.Message<GroupMember> {
    override operator fun plus(other: GroupMember?) = protoMergeImpl(other)
    override val protoSize by lazy { protoSizeImpl() }
    override fun protoMarshal(m: pbandk.Marshaller) = protoMarshalImpl(m)
    companion object : pbandk.Message.Companion<GroupMember> {
        override fun protoUnmarshal(u: pbandk.Unmarshaller) = GroupMember.protoUnmarshalImpl(u)
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

data class GroupInfo(
    val groupId: Long = 0L,
    val tags: List<String> = emptyList(),
    val groupName: String = "",
    val groupType: String = "",
    val params: Map<String, String> = emptyMap(),
    val unknownFields: Map<Int, pbandk.UnknownField> = emptyMap()
) : pbandk.Message<GroupInfo> {
    override operator fun plus(other: GroupInfo?) = protoMergeImpl(other)
    override val protoSize by lazy { protoSizeImpl() }
    override fun protoMarshal(m: pbandk.Marshaller) = protoMarshalImpl(m)
    companion object : pbandk.Message.Companion<GroupInfo> {
        override fun protoUnmarshal(u: pbandk.Unmarshaller) = GroupInfo.protoUnmarshalImpl(u)
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

data class GroupOpReq(
    val operation: model.GroupOperationType = model.GroupOperationType.fromValue(0),
    val reqMem: model.GroupMember? = null,
    val group: model.GroupInfo? = null,
    val members: List<model.GroupMember> = emptyList(),
    val sendId: Long = 0L,
    val msgId: Long = 0L,
    val params: Map<String, String> = emptyMap(),
    val unknownFields: Map<Int, pbandk.UnknownField> = emptyMap()
) : pbandk.Message<GroupOpReq> {
    override operator fun plus(other: GroupOpReq?) = protoMergeImpl(other)
    override val protoSize by lazy { protoSizeImpl() }
    override fun protoMarshal(m: pbandk.Marshaller) = protoMarshalImpl(m)
    companion object : pbandk.Message.Companion<GroupOpReq> {
        override fun protoUnmarshal(u: pbandk.Unmarshaller) = GroupOpReq.protoUnmarshalImpl(u)
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

data class GroupOpResult(
    val operation: model.GroupOperationType = model.GroupOperationType.fromValue(0),
    val reqMem: model.GroupMember? = null,
    val result: String = "",
    val detail: String = "",
    val group: model.GroupInfo? = null,
    val sendId: Long = 0L,
    val msgId: Long = 0L,
    val members: List<model.GroupMember> = emptyList(),
    val params: Map<String, String> = emptyMap(),
    val unknownFields: Map<Int, pbandk.UnknownField> = emptyMap()
) : pbandk.Message<GroupOpResult> {
    override operator fun plus(other: GroupOpResult?) = protoMergeImpl(other)
    override val protoSize by lazy { protoSizeImpl() }
    override fun protoMarshal(m: pbandk.Marshaller) = protoMarshalImpl(m)
    companion object : pbandk.Message.Companion<GroupOpResult> {
        override fun protoUnmarshal(u: pbandk.Unmarshaller) = GroupOpResult.protoUnmarshalImpl(u)
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

private fun UserInfo.protoMergeImpl(plus: UserInfo?): UserInfo = plus?.copy(
    params = params + plus.params,
    unknownFields = unknownFields + plus.unknownFields
) ?: this

private fun UserInfo.protoSizeImpl(): Int {
    var protoSize = 0
    if (userId != 0L) protoSize += pbandk.Sizer.tagSize(1) + pbandk.Sizer.int64Size(userId)
    if (userName.isNotEmpty()) protoSize += pbandk.Sizer.tagSize(2) + pbandk.Sizer.stringSize(userName)
    if (nickName.isNotEmpty()) protoSize += pbandk.Sizer.tagSize(3) + pbandk.Sizer.stringSize(nickName)
    if (email.isNotEmpty()) protoSize += pbandk.Sizer.tagSize(4) + pbandk.Sizer.stringSize(email)
    if (phone.isNotEmpty()) protoSize += pbandk.Sizer.tagSize(5) + pbandk.Sizer.stringSize(phone)
    if (gender.isNotEmpty()) protoSize += pbandk.Sizer.tagSize(6) + pbandk.Sizer.stringSize(gender)
    if (age != 0) protoSize += pbandk.Sizer.tagSize(7) + pbandk.Sizer.int32Size(age)
    if (region.isNotEmpty()) protoSize += pbandk.Sizer.tagSize(8) + pbandk.Sizer.stringSize(region)
    if (icon.isNotEmpty()) protoSize += pbandk.Sizer.tagSize(9) + pbandk.Sizer.stringSize(icon)
    if (params.isNotEmpty()) protoSize += pbandk.Sizer.mapSize(10, params, model.UserInfo::ParamsEntry)
    protoSize += unknownFields.entries.sumBy { it.value.size() }
    return protoSize
}

private fun UserInfo.protoMarshalImpl(protoMarshal: pbandk.Marshaller) {
    if (userId != 0L) protoMarshal.writeTag(8).writeInt64(userId)
    if (userName.isNotEmpty()) protoMarshal.writeTag(18).writeString(userName)
    if (nickName.isNotEmpty()) protoMarshal.writeTag(26).writeString(nickName)
    if (email.isNotEmpty()) protoMarshal.writeTag(34).writeString(email)
    if (phone.isNotEmpty()) protoMarshal.writeTag(42).writeString(phone)
    if (gender.isNotEmpty()) protoMarshal.writeTag(50).writeString(gender)
    if (age != 0) protoMarshal.writeTag(56).writeInt32(age)
    if (region.isNotEmpty()) protoMarshal.writeTag(66).writeString(region)
    if (icon.isNotEmpty()) protoMarshal.writeTag(74).writeString(icon)
    if (params.isNotEmpty()) protoMarshal.writeMap(82, params, model.UserInfo::ParamsEntry)
    if (unknownFields.isNotEmpty()) protoMarshal.writeUnknownFields(unknownFields)
}

private fun UserInfo.Companion.protoUnmarshalImpl(protoUnmarshal: pbandk.Unmarshaller): UserInfo {
    var userId = 0L
    var userName = ""
    var nickName = ""
    var email = ""
    var phone = ""
    var gender = ""
    var age = 0
    var region = ""
    var icon = ""
    var params: pbandk.MessageMap.Builder<String, String>? = null
    while (true) when (protoUnmarshal.readTag()) {
        0 -> return UserInfo(userId, userName, nickName, email,
            phone, gender, age, region,
            icon, pbandk.MessageMap.Builder.fixed(params), protoUnmarshal.unknownFields())
        8 -> userId = protoUnmarshal.readInt64()
        18 -> userName = protoUnmarshal.readString()
        26 -> nickName = protoUnmarshal.readString()
        34 -> email = protoUnmarshal.readString()
        42 -> phone = protoUnmarshal.readString()
        50 -> gender = protoUnmarshal.readString()
        56 -> age = protoUnmarshal.readInt32()
        66 -> region = protoUnmarshal.readString()
        74 -> icon = protoUnmarshal.readString()
        82 -> params = protoUnmarshal.readMap(params, model.UserInfo.ParamsEntry.Companion, true)
        else -> protoUnmarshal.unknownField()
    }
}

private fun UserInfo.ParamsEntry.protoMergeImpl(plus: UserInfo.ParamsEntry?): UserInfo.ParamsEntry = plus?.copy(
    unknownFields = unknownFields + plus.unknownFields
) ?: this

private fun UserInfo.ParamsEntry.protoSizeImpl(): Int {
    var protoSize = 0
    if (key.isNotEmpty()) protoSize += pbandk.Sizer.tagSize(1) + pbandk.Sizer.stringSize(key)
    if (value.isNotEmpty()) protoSize += pbandk.Sizer.tagSize(2) + pbandk.Sizer.stringSize(value)
    protoSize += unknownFields.entries.sumBy { it.value.size() }
    return protoSize
}

private fun UserInfo.ParamsEntry.protoMarshalImpl(protoMarshal: pbandk.Marshaller) {
    if (key.isNotEmpty()) protoMarshal.writeTag(10).writeString(key)
    if (value.isNotEmpty()) protoMarshal.writeTag(18).writeString(value)
    if (unknownFields.isNotEmpty()) protoMarshal.writeUnknownFields(unknownFields)
}

private fun UserInfo.ParamsEntry.Companion.protoUnmarshalImpl(protoUnmarshal: pbandk.Unmarshaller): UserInfo.ParamsEntry {
    var key = ""
    var value = ""
    while (true) when (protoUnmarshal.readTag()) {
        0 -> return UserInfo.ParamsEntry(key, value, protoUnmarshal.unknownFields())
        10 -> key = protoUnmarshal.readString()
        18 -> value = protoUnmarshal.readString()
        else -> protoUnmarshal.unknownField()
    }
}

private fun UserOpReq.protoMergeImpl(plus: UserOpReq?): UserOpReq = plus?.copy(
    user = user?.plus(plus.user) ?: plus.user,
    params = params + plus.params,
    unknownFields = unknownFields + plus.unknownFields
) ?: this

private fun UserOpReq.protoSizeImpl(): Int {
    var protoSize = 0
    if (operation.value != 0) protoSize += pbandk.Sizer.tagSize(1) + pbandk.Sizer.enumSize(operation)
    if (user != null) protoSize += pbandk.Sizer.tagSize(2) + pbandk.Sizer.messageSize(user)
    if (params.isNotEmpty()) protoSize += pbandk.Sizer.mapSize(3, params, model.UserOpReq::ParamsEntry)
    protoSize += unknownFields.entries.sumBy { it.value.size() }
    return protoSize
}

private fun UserOpReq.protoMarshalImpl(protoMarshal: pbandk.Marshaller) {
    if (operation.value != 0) protoMarshal.writeTag(8).writeEnum(operation)
    if (user != null) protoMarshal.writeTag(18).writeMessage(user)
    if (params.isNotEmpty()) protoMarshal.writeMap(26, params, model.UserOpReq::ParamsEntry)
    if (unknownFields.isNotEmpty()) protoMarshal.writeUnknownFields(unknownFields)
}

private fun UserOpReq.Companion.protoUnmarshalImpl(protoUnmarshal: pbandk.Unmarshaller): UserOpReq {
    var operation: model.UserOperationType = model.UserOperationType.fromValue(0)
    var user: model.UserInfo? = null
    var params: pbandk.MessageMap.Builder<String, String>? = null
    while (true) when (protoUnmarshal.readTag()) {
        0 -> return UserOpReq(operation, user, pbandk.MessageMap.Builder.fixed(params), protoUnmarshal.unknownFields())
        8 -> operation = protoUnmarshal.readEnum(model.UserOperationType.Companion)
        18 -> user = protoUnmarshal.readMessage(model.UserInfo.Companion)
        26 -> params = protoUnmarshal.readMap(params, model.UserOpReq.ParamsEntry.Companion, true)
        else -> protoUnmarshal.unknownField()
    }
}

private fun UserOpReq.ParamsEntry.protoMergeImpl(plus: UserOpReq.ParamsEntry?): UserOpReq.ParamsEntry = plus?.copy(
    unknownFields = unknownFields + plus.unknownFields
) ?: this

private fun UserOpReq.ParamsEntry.protoSizeImpl(): Int {
    var protoSize = 0
    if (key.isNotEmpty()) protoSize += pbandk.Sizer.tagSize(1) + pbandk.Sizer.stringSize(key)
    if (value.isNotEmpty()) protoSize += pbandk.Sizer.tagSize(2) + pbandk.Sizer.stringSize(value)
    protoSize += unknownFields.entries.sumBy { it.value.size() }
    return protoSize
}

private fun UserOpReq.ParamsEntry.protoMarshalImpl(protoMarshal: pbandk.Marshaller) {
    if (key.isNotEmpty()) protoMarshal.writeTag(10).writeString(key)
    if (value.isNotEmpty()) protoMarshal.writeTag(18).writeString(value)
    if (unknownFields.isNotEmpty()) protoMarshal.writeUnknownFields(unknownFields)
}

private fun UserOpReq.ParamsEntry.Companion.protoUnmarshalImpl(protoUnmarshal: pbandk.Unmarshaller): UserOpReq.ParamsEntry {
    var key = ""
    var value = ""
    while (true) when (protoUnmarshal.readTag()) {
        0 -> return UserOpReq.ParamsEntry(key, value, protoUnmarshal.unknownFields())
        10 -> key = protoUnmarshal.readString()
        18 -> value = protoUnmarshal.readString()
        else -> protoUnmarshal.unknownField()
    }
}

private fun UserOpResult.protoMergeImpl(plus: UserOpResult?): UserOpResult = plus?.copy(
    users = users + plus.users,
    params = params + plus.params,
    unknownFields = unknownFields + plus.unknownFields
) ?: this

private fun UserOpResult.protoSizeImpl(): Int {
    var protoSize = 0
    if (operation.value != 0) protoSize += pbandk.Sizer.tagSize(1) + pbandk.Sizer.enumSize(operation)
    if (result.isNotEmpty()) protoSize += pbandk.Sizer.tagSize(2) + pbandk.Sizer.stringSize(result)
    if (users.isNotEmpty()) protoSize += (pbandk.Sizer.tagSize(3) * users.size) + users.sumBy(pbandk.Sizer::messageSize)
    if (params.isNotEmpty()) protoSize += pbandk.Sizer.mapSize(4, params, model.UserOpResult::ParamsEntry)
    protoSize += unknownFields.entries.sumBy { it.value.size() }
    return protoSize
}

private fun UserOpResult.protoMarshalImpl(protoMarshal: pbandk.Marshaller) {
    if (operation.value != 0) protoMarshal.writeTag(8).writeEnum(operation)
    if (result.isNotEmpty()) protoMarshal.writeTag(18).writeString(result)
    if (users.isNotEmpty()) users.forEach { protoMarshal.writeTag(26).writeMessage(it) }
    if (params.isNotEmpty()) protoMarshal.writeMap(34, params, model.UserOpResult::ParamsEntry)
    if (unknownFields.isNotEmpty()) protoMarshal.writeUnknownFields(unknownFields)
}

private fun UserOpResult.Companion.protoUnmarshalImpl(protoUnmarshal: pbandk.Unmarshaller): UserOpResult {
    var operation: model.UserOperationType = model.UserOperationType.fromValue(0)
    var result = ""
    var users: pbandk.ListWithSize.Builder<model.UserInfo>? = null
    var params: pbandk.MessageMap.Builder<String, String>? = null
    while (true) when (protoUnmarshal.readTag()) {
        0 -> return UserOpResult(operation, result, pbandk.ListWithSize.Builder.fixed(users), pbandk.MessageMap.Builder.fixed(params), protoUnmarshal.unknownFields())
        8 -> operation = protoUnmarshal.readEnum(model.UserOperationType.Companion)
        18 -> result = protoUnmarshal.readString()
        26 -> users = protoUnmarshal.readRepeatedMessage(users, model.UserInfo.Companion, true)
        34 -> params = protoUnmarshal.readMap(params, model.UserOpResult.ParamsEntry.Companion, true)
        else -> protoUnmarshal.unknownField()
    }
}

private fun UserOpResult.ParamsEntry.protoMergeImpl(plus: UserOpResult.ParamsEntry?): UserOpResult.ParamsEntry = plus?.copy(
    unknownFields = unknownFields + plus.unknownFields
) ?: this

private fun UserOpResult.ParamsEntry.protoSizeImpl(): Int {
    var protoSize = 0
    if (key.isNotEmpty()) protoSize += pbandk.Sizer.tagSize(1) + pbandk.Sizer.stringSize(key)
    if (value.isNotEmpty()) protoSize += pbandk.Sizer.tagSize(2) + pbandk.Sizer.stringSize(value)
    protoSize += unknownFields.entries.sumBy { it.value.size() }
    return protoSize
}

private fun UserOpResult.ParamsEntry.protoMarshalImpl(protoMarshal: pbandk.Marshaller) {
    if (key.isNotEmpty()) protoMarshal.writeTag(10).writeString(key)
    if (value.isNotEmpty()) protoMarshal.writeTag(18).writeString(value)
    if (unknownFields.isNotEmpty()) protoMarshal.writeUnknownFields(unknownFields)
}

private fun UserOpResult.ParamsEntry.Companion.protoUnmarshalImpl(protoUnmarshal: pbandk.Unmarshaller): UserOpResult.ParamsEntry {
    var key = ""
    var value = ""
    while (true) when (protoUnmarshal.readTag()) {
        0 -> return UserOpResult.ParamsEntry(key, value, protoUnmarshal.unknownFields())
        10 -> key = protoUnmarshal.readString()
        18 -> value = protoUnmarshal.readString()
        else -> protoUnmarshal.unknownField()
    }
}

private fun FriendOpReq.protoMergeImpl(plus: FriendOpReq?): FriendOpReq = plus?.copy(
    user = user?.plus(plus.user) ?: plus.user,
    params = params + plus.params,
    unknownFields = unknownFields + plus.unknownFields
) ?: this

private fun FriendOpReq.protoSizeImpl(): Int {
    var protoSize = 0
    if (operation.value != 0) protoSize += pbandk.Sizer.tagSize(1) + pbandk.Sizer.enumSize(operation)
    if (user != null) protoSize += pbandk.Sizer.tagSize(2) + pbandk.Sizer.messageSize(user)
    if (sendId != 0L) protoSize += pbandk.Sizer.tagSize(3) + pbandk.Sizer.int64Size(sendId)
    if (msgId != 0L) protoSize += pbandk.Sizer.tagSize(4) + pbandk.Sizer.int64Size(msgId)
    if (params.isNotEmpty()) protoSize += pbandk.Sizer.mapSize(5, params, model.FriendOpReq::ParamsEntry)
    protoSize += unknownFields.entries.sumBy { it.value.size() }
    return protoSize
}

private fun FriendOpReq.protoMarshalImpl(protoMarshal: pbandk.Marshaller) {
    if (operation.value != 0) protoMarshal.writeTag(8).writeEnum(operation)
    if (user != null) protoMarshal.writeTag(18).writeMessage(user)
    if (sendId != 0L) protoMarshal.writeTag(24).writeInt64(sendId)
    if (msgId != 0L) protoMarshal.writeTag(32).writeInt64(msgId)
    if (params.isNotEmpty()) protoMarshal.writeMap(42, params, model.FriendOpReq::ParamsEntry)
    if (unknownFields.isNotEmpty()) protoMarshal.writeUnknownFields(unknownFields)
}

private fun FriendOpReq.Companion.protoUnmarshalImpl(protoUnmarshal: pbandk.Unmarshaller): FriendOpReq {
    var operation: model.UserOperationType = model.UserOperationType.fromValue(0)
    var user: model.UserInfo? = null
    var sendId = 0L
    var msgId = 0L
    var params: pbandk.MessageMap.Builder<String, String>? = null
    while (true) when (protoUnmarshal.readTag()) {
        0 -> return FriendOpReq(operation, user, sendId, msgId,
            pbandk.MessageMap.Builder.fixed(params), protoUnmarshal.unknownFields())
        8 -> operation = protoUnmarshal.readEnum(model.UserOperationType.Companion)
        18 -> user = protoUnmarshal.readMessage(model.UserInfo.Companion)
        24 -> sendId = protoUnmarshal.readInt64()
        32 -> msgId = protoUnmarshal.readInt64()
        42 -> params = protoUnmarshal.readMap(params, model.FriendOpReq.ParamsEntry.Companion, true)
        else -> protoUnmarshal.unknownField()
    }
}

private fun FriendOpReq.ParamsEntry.protoMergeImpl(plus: FriendOpReq.ParamsEntry?): FriendOpReq.ParamsEntry = plus?.copy(
    unknownFields = unknownFields + plus.unknownFields
) ?: this

private fun FriendOpReq.ParamsEntry.protoSizeImpl(): Int {
    var protoSize = 0
    if (key.isNotEmpty()) protoSize += pbandk.Sizer.tagSize(1) + pbandk.Sizer.stringSize(key)
    if (value.isNotEmpty()) protoSize += pbandk.Sizer.tagSize(2) + pbandk.Sizer.stringSize(value)
    protoSize += unknownFields.entries.sumBy { it.value.size() }
    return protoSize
}

private fun FriendOpReq.ParamsEntry.protoMarshalImpl(protoMarshal: pbandk.Marshaller) {
    if (key.isNotEmpty()) protoMarshal.writeTag(10).writeString(key)
    if (value.isNotEmpty()) protoMarshal.writeTag(18).writeString(value)
    if (unknownFields.isNotEmpty()) protoMarshal.writeUnknownFields(unknownFields)
}

private fun FriendOpReq.ParamsEntry.Companion.protoUnmarshalImpl(protoUnmarshal: pbandk.Unmarshaller): FriendOpReq.ParamsEntry {
    var key = ""
    var value = ""
    while (true) when (protoUnmarshal.readTag()) {
        0 -> return FriendOpReq.ParamsEntry(key, value, protoUnmarshal.unknownFields())
        10 -> key = protoUnmarshal.readString()
        18 -> value = protoUnmarshal.readString()
        else -> protoUnmarshal.unknownField()
    }
}

private fun FriendOpResult.protoMergeImpl(plus: FriendOpResult?): FriendOpResult = plus?.copy(
    user = user?.plus(plus.user) ?: plus.user,
    users = users?.plus(plus.users) ?: plus.users,
    params = params + plus.params,
    unknownFields = unknownFields + plus.unknownFields
) ?: this

private fun FriendOpResult.protoSizeImpl(): Int {
    var protoSize = 0
    if (operation.value != 0) protoSize += pbandk.Sizer.tagSize(1) + pbandk.Sizer.enumSize(operation)
    if (result.isNotEmpty()) protoSize += pbandk.Sizer.tagSize(2) + pbandk.Sizer.stringSize(result)
    if (user != null) protoSize += pbandk.Sizer.tagSize(3) + pbandk.Sizer.messageSize(user)
    if (users != null) protoSize += pbandk.Sizer.tagSize(4) + pbandk.Sizer.messageSize(users)
    if (sendId != 0L) protoSize += pbandk.Sizer.tagSize(5) + pbandk.Sizer.int64Size(sendId)
    if (msgId != 0L) protoSize += pbandk.Sizer.tagSize(6) + pbandk.Sizer.int64Size(msgId)
    if (params.isNotEmpty()) protoSize += pbandk.Sizer.mapSize(7, params, model.FriendOpResult::ParamsEntry)
    protoSize += unknownFields.entries.sumBy { it.value.size() }
    return protoSize
}

private fun FriendOpResult.protoMarshalImpl(protoMarshal: pbandk.Marshaller) {
    if (operation.value != 0) protoMarshal.writeTag(8).writeEnum(operation)
    if (result.isNotEmpty()) protoMarshal.writeTag(18).writeString(result)
    if (user != null) protoMarshal.writeTag(26).writeMessage(user)
    if (users != null) protoMarshal.writeTag(34).writeMessage(users)
    if (sendId != 0L) protoMarshal.writeTag(40).writeInt64(sendId)
    if (msgId != 0L) protoMarshal.writeTag(48).writeInt64(msgId)
    if (params.isNotEmpty()) protoMarshal.writeMap(58, params, model.FriendOpResult::ParamsEntry)
    if (unknownFields.isNotEmpty()) protoMarshal.writeUnknownFields(unknownFields)
}

private fun FriendOpResult.Companion.protoUnmarshalImpl(protoUnmarshal: pbandk.Unmarshaller): FriendOpResult {
    var operation: model.UserOperationType = model.UserOperationType.fromValue(0)
    var result = ""
    var user: model.UserInfo? = null
    var users: model.UserInfo? = null
    var sendId = 0L
    var msgId = 0L
    var params: pbandk.MessageMap.Builder<String, String>? = null
    while (true) when (protoUnmarshal.readTag()) {
        0 -> return FriendOpResult(operation, result, user, users,
            sendId, msgId, pbandk.MessageMap.Builder.fixed(params), protoUnmarshal.unknownFields())
        8 -> operation = protoUnmarshal.readEnum(model.UserOperationType.Companion)
        18 -> result = protoUnmarshal.readString()
        26 -> user = protoUnmarshal.readMessage(model.UserInfo.Companion)
        34 -> users = protoUnmarshal.readMessage(model.UserInfo.Companion)
        40 -> sendId = protoUnmarshal.readInt64()
        48 -> msgId = protoUnmarshal.readInt64()
        58 -> params = protoUnmarshal.readMap(params, model.FriendOpResult.ParamsEntry.Companion, true)
        else -> protoUnmarshal.unknownField()
    }
}

private fun FriendOpResult.ParamsEntry.protoMergeImpl(plus: FriendOpResult.ParamsEntry?): FriendOpResult.ParamsEntry = plus?.copy(
    unknownFields = unknownFields + plus.unknownFields
) ?: this

private fun FriendOpResult.ParamsEntry.protoSizeImpl(): Int {
    var protoSize = 0
    if (key.isNotEmpty()) protoSize += pbandk.Sizer.tagSize(1) + pbandk.Sizer.stringSize(key)
    if (value.isNotEmpty()) protoSize += pbandk.Sizer.tagSize(2) + pbandk.Sizer.stringSize(value)
    protoSize += unknownFields.entries.sumBy { it.value.size() }
    return protoSize
}

private fun FriendOpResult.ParamsEntry.protoMarshalImpl(protoMarshal: pbandk.Marshaller) {
    if (key.isNotEmpty()) protoMarshal.writeTag(10).writeString(key)
    if (value.isNotEmpty()) protoMarshal.writeTag(18).writeString(value)
    if (unknownFields.isNotEmpty()) protoMarshal.writeUnknownFields(unknownFields)
}

private fun FriendOpResult.ParamsEntry.Companion.protoUnmarshalImpl(protoUnmarshal: pbandk.Unmarshaller): FriendOpResult.ParamsEntry {
    var key = ""
    var value = ""
    while (true) when (protoUnmarshal.readTag()) {
        0 -> return FriendOpResult.ParamsEntry(key, value, protoUnmarshal.unknownFields())
        10 -> key = protoUnmarshal.readString()
        18 -> value = protoUnmarshal.readString()
        else -> protoUnmarshal.unknownField()
    }
}

private fun GroupMember.protoMergeImpl(plus: GroupMember?): GroupMember = plus?.copy(
    params = params + plus.params,
    unknownFields = unknownFields + plus.unknownFields
) ?: this

private fun GroupMember.protoSizeImpl(): Int {
    var protoSize = 0
    if (userId != 0L) protoSize += pbandk.Sizer.tagSize(1) + pbandk.Sizer.int64Size(userId)
    if (nick.isNotEmpty()) protoSize += pbandk.Sizer.tagSize(2) + pbandk.Sizer.stringSize(nick)
    if (icon.isNotEmpty()) protoSize += pbandk.Sizer.tagSize(3) + pbandk.Sizer.stringSize(icon)
    if (role.isNotEmpty()) protoSize += pbandk.Sizer.tagSize(4) + pbandk.Sizer.stringSize(role)
    if (groupId != 0L) protoSize += pbandk.Sizer.tagSize(5) + pbandk.Sizer.int64Size(groupId)
    if (params.isNotEmpty()) protoSize += pbandk.Sizer.mapSize(6, params, model.GroupMember::ParamsEntry)
    protoSize += unknownFields.entries.sumBy { it.value.size() }
    return protoSize
}

private fun GroupMember.protoMarshalImpl(protoMarshal: pbandk.Marshaller) {
    if (userId != 0L) protoMarshal.writeTag(8).writeInt64(userId)
    if (nick.isNotEmpty()) protoMarshal.writeTag(18).writeString(nick)
    if (icon.isNotEmpty()) protoMarshal.writeTag(26).writeString(icon)
    if (role.isNotEmpty()) protoMarshal.writeTag(34).writeString(role)
    if (groupId != 0L) protoMarshal.writeTag(40).writeInt64(groupId)
    if (params.isNotEmpty()) protoMarshal.writeMap(50, params, model.GroupMember::ParamsEntry)
    if (unknownFields.isNotEmpty()) protoMarshal.writeUnknownFields(unknownFields)
}

private fun GroupMember.Companion.protoUnmarshalImpl(protoUnmarshal: pbandk.Unmarshaller): GroupMember {
    var userId = 0L
    var nick = ""
    var icon = ""
    var role = ""
    var groupId = 0L
    var params: pbandk.MessageMap.Builder<String, String>? = null
    while (true) when (protoUnmarshal.readTag()) {
        0 -> return GroupMember(userId, nick, icon, role,
            groupId, pbandk.MessageMap.Builder.fixed(params), protoUnmarshal.unknownFields())
        8 -> userId = protoUnmarshal.readInt64()
        18 -> nick = protoUnmarshal.readString()
        26 -> icon = protoUnmarshal.readString()
        34 -> role = protoUnmarshal.readString()
        40 -> groupId = protoUnmarshal.readInt64()
        50 -> params = protoUnmarshal.readMap(params, model.GroupMember.ParamsEntry.Companion, true)
        else -> protoUnmarshal.unknownField()
    }
}

private fun GroupMember.ParamsEntry.protoMergeImpl(plus: GroupMember.ParamsEntry?): GroupMember.ParamsEntry = plus?.copy(
    unknownFields = unknownFields + plus.unknownFields
) ?: this

private fun GroupMember.ParamsEntry.protoSizeImpl(): Int {
    var protoSize = 0
    if (key.isNotEmpty()) protoSize += pbandk.Sizer.tagSize(1) + pbandk.Sizer.stringSize(key)
    if (value.isNotEmpty()) protoSize += pbandk.Sizer.tagSize(2) + pbandk.Sizer.stringSize(value)
    protoSize += unknownFields.entries.sumBy { it.value.size() }
    return protoSize
}

private fun GroupMember.ParamsEntry.protoMarshalImpl(protoMarshal: pbandk.Marshaller) {
    if (key.isNotEmpty()) protoMarshal.writeTag(10).writeString(key)
    if (value.isNotEmpty()) protoMarshal.writeTag(18).writeString(value)
    if (unknownFields.isNotEmpty()) protoMarshal.writeUnknownFields(unknownFields)
}

private fun GroupMember.ParamsEntry.Companion.protoUnmarshalImpl(protoUnmarshal: pbandk.Unmarshaller): GroupMember.ParamsEntry {
    var key = ""
    var value = ""
    while (true) when (protoUnmarshal.readTag()) {
        0 -> return GroupMember.ParamsEntry(key, value, protoUnmarshal.unknownFields())
        10 -> key = protoUnmarshal.readString()
        18 -> value = protoUnmarshal.readString()
        else -> protoUnmarshal.unknownField()
    }
}

private fun GroupInfo.protoMergeImpl(plus: GroupInfo?): GroupInfo = plus?.copy(
    tags = tags + plus.tags,
    params = params + plus.params,
    unknownFields = unknownFields + plus.unknownFields
) ?: this

private fun GroupInfo.protoSizeImpl(): Int {
    var protoSize = 0
    if (groupId != 0L) protoSize += pbandk.Sizer.tagSize(1) + pbandk.Sizer.int64Size(groupId)
    if (tags.isNotEmpty()) protoSize += (pbandk.Sizer.tagSize(2) * tags.size) + tags.sumBy(pbandk.Sizer::stringSize)
    if (groupName.isNotEmpty()) protoSize += pbandk.Sizer.tagSize(3) + pbandk.Sizer.stringSize(groupName)
    if (groupType.isNotEmpty()) protoSize += pbandk.Sizer.tagSize(4) + pbandk.Sizer.stringSize(groupType)
    if (params.isNotEmpty()) protoSize += pbandk.Sizer.mapSize(5, params, model.GroupInfo::ParamsEntry)
    protoSize += unknownFields.entries.sumBy { it.value.size() }
    return protoSize
}

private fun GroupInfo.protoMarshalImpl(protoMarshal: pbandk.Marshaller) {
    if (groupId != 0L) protoMarshal.writeTag(8).writeInt64(groupId)
    if (tags.isNotEmpty()) tags.forEach { protoMarshal.writeTag(18).writeString(it) }
    if (groupName.isNotEmpty()) protoMarshal.writeTag(26).writeString(groupName)
    if (groupType.isNotEmpty()) protoMarshal.writeTag(34).writeString(groupType)
    if (params.isNotEmpty()) protoMarshal.writeMap(42, params, model.GroupInfo::ParamsEntry)
    if (unknownFields.isNotEmpty()) protoMarshal.writeUnknownFields(unknownFields)
}

private fun GroupInfo.Companion.protoUnmarshalImpl(protoUnmarshal: pbandk.Unmarshaller): GroupInfo {
    var groupId = 0L
    var tags: pbandk.ListWithSize.Builder<String>? = null
    var groupName = ""
    var groupType = ""
    var params: pbandk.MessageMap.Builder<String, String>? = null
    while (true) when (protoUnmarshal.readTag()) {
        0 -> return GroupInfo(groupId, pbandk.ListWithSize.Builder.fixed(tags), groupName, groupType,
            pbandk.MessageMap.Builder.fixed(params), protoUnmarshal.unknownFields())
        8 -> groupId = protoUnmarshal.readInt64()
        18 -> tags = protoUnmarshal.readRepeated(tags, protoUnmarshal::readString, true)
        26 -> groupName = protoUnmarshal.readString()
        34 -> groupType = protoUnmarshal.readString()
        42 -> params = protoUnmarshal.readMap(params, model.GroupInfo.ParamsEntry.Companion, true)
        else -> protoUnmarshal.unknownField()
    }
}

private fun GroupInfo.ParamsEntry.protoMergeImpl(plus: GroupInfo.ParamsEntry?): GroupInfo.ParamsEntry = plus?.copy(
    unknownFields = unknownFields + plus.unknownFields
) ?: this

private fun GroupInfo.ParamsEntry.protoSizeImpl(): Int {
    var protoSize = 0
    if (key.isNotEmpty()) protoSize += pbandk.Sizer.tagSize(1) + pbandk.Sizer.stringSize(key)
    if (value.isNotEmpty()) protoSize += pbandk.Sizer.tagSize(2) + pbandk.Sizer.stringSize(value)
    protoSize += unknownFields.entries.sumBy { it.value.size() }
    return protoSize
}

private fun GroupInfo.ParamsEntry.protoMarshalImpl(protoMarshal: pbandk.Marshaller) {
    if (key.isNotEmpty()) protoMarshal.writeTag(10).writeString(key)
    if (value.isNotEmpty()) protoMarshal.writeTag(18).writeString(value)
    if (unknownFields.isNotEmpty()) protoMarshal.writeUnknownFields(unknownFields)
}

private fun GroupInfo.ParamsEntry.Companion.protoUnmarshalImpl(protoUnmarshal: pbandk.Unmarshaller): GroupInfo.ParamsEntry {
    var key = ""
    var value = ""
    while (true) when (protoUnmarshal.readTag()) {
        0 -> return GroupInfo.ParamsEntry(key, value, protoUnmarshal.unknownFields())
        10 -> key = protoUnmarshal.readString()
        18 -> value = protoUnmarshal.readString()
        else -> protoUnmarshal.unknownField()
    }
}

private fun GroupOpReq.protoMergeImpl(plus: GroupOpReq?): GroupOpReq = plus?.copy(
    reqMem = reqMem?.plus(plus.reqMem) ?: plus.reqMem,
    group = group?.plus(plus.group) ?: plus.group,
    members = members + plus.members,
    params = params + plus.params,
    unknownFields = unknownFields + plus.unknownFields
) ?: this

private fun GroupOpReq.protoSizeImpl(): Int {
    var protoSize = 0
    if (operation.value != 0) protoSize += pbandk.Sizer.tagSize(1) + pbandk.Sizer.enumSize(operation)
    if (reqMem != null) protoSize += pbandk.Sizer.tagSize(2) + pbandk.Sizer.messageSize(reqMem)
    if (group != null) protoSize += pbandk.Sizer.tagSize(3) + pbandk.Sizer.messageSize(group)
    if (members.isNotEmpty()) protoSize += (pbandk.Sizer.tagSize(4) * members.size) + members.sumBy(pbandk.Sizer::messageSize)
    if (sendId != 0L) protoSize += pbandk.Sizer.tagSize(6) + pbandk.Sizer.int64Size(sendId)
    if (msgId != 0L) protoSize += pbandk.Sizer.tagSize(7) + pbandk.Sizer.int64Size(msgId)
    if (params.isNotEmpty()) protoSize += pbandk.Sizer.mapSize(8, params, model.GroupOpReq::ParamsEntry)
    protoSize += unknownFields.entries.sumBy { it.value.size() }
    return protoSize
}

private fun GroupOpReq.protoMarshalImpl(protoMarshal: pbandk.Marshaller) {
    if (operation.value != 0) protoMarshal.writeTag(8).writeEnum(operation)
    if (reqMem != null) protoMarshal.writeTag(18).writeMessage(reqMem)
    if (group != null) protoMarshal.writeTag(26).writeMessage(group)
    if (members.isNotEmpty()) members.forEach { protoMarshal.writeTag(34).writeMessage(it) }
    if (sendId != 0L) protoMarshal.writeTag(48).writeInt64(sendId)
    if (msgId != 0L) protoMarshal.writeTag(56).writeInt64(msgId)
    if (params.isNotEmpty()) protoMarshal.writeMap(66, params, model.GroupOpReq::ParamsEntry)
    if (unknownFields.isNotEmpty()) protoMarshal.writeUnknownFields(unknownFields)
}

private fun GroupOpReq.Companion.protoUnmarshalImpl(protoUnmarshal: pbandk.Unmarshaller): GroupOpReq {
    var operation: model.GroupOperationType = model.GroupOperationType.fromValue(0)
    var reqMem: model.GroupMember? = null
    var group: model.GroupInfo? = null
    var members: pbandk.ListWithSize.Builder<model.GroupMember>? = null
    var sendId = 0L
    var msgId = 0L
    var params: pbandk.MessageMap.Builder<String, String>? = null
    while (true) when (protoUnmarshal.readTag()) {
        0 -> return GroupOpReq(operation, reqMem, group, pbandk.ListWithSize.Builder.fixed(members),
            sendId, msgId, pbandk.MessageMap.Builder.fixed(params), protoUnmarshal.unknownFields())
        8 -> operation = protoUnmarshal.readEnum(model.GroupOperationType.Companion)
        18 -> reqMem = protoUnmarshal.readMessage(model.GroupMember.Companion)
        26 -> group = protoUnmarshal.readMessage(model.GroupInfo.Companion)
        34 -> members = protoUnmarshal.readRepeatedMessage(members, model.GroupMember.Companion, true)
        48 -> sendId = protoUnmarshal.readInt64()
        56 -> msgId = protoUnmarshal.readInt64()
        66 -> params = protoUnmarshal.readMap(params, model.GroupOpReq.ParamsEntry.Companion, true)
        else -> protoUnmarshal.unknownField()
    }
}

private fun GroupOpReq.ParamsEntry.protoMergeImpl(plus: GroupOpReq.ParamsEntry?): GroupOpReq.ParamsEntry = plus?.copy(
    unknownFields = unknownFields + plus.unknownFields
) ?: this

private fun GroupOpReq.ParamsEntry.protoSizeImpl(): Int {
    var protoSize = 0
    if (key.isNotEmpty()) protoSize += pbandk.Sizer.tagSize(1) + pbandk.Sizer.stringSize(key)
    if (value.isNotEmpty()) protoSize += pbandk.Sizer.tagSize(2) + pbandk.Sizer.stringSize(value)
    protoSize += unknownFields.entries.sumBy { it.value.size() }
    return protoSize
}

private fun GroupOpReq.ParamsEntry.protoMarshalImpl(protoMarshal: pbandk.Marshaller) {
    if (key.isNotEmpty()) protoMarshal.writeTag(10).writeString(key)
    if (value.isNotEmpty()) protoMarshal.writeTag(18).writeString(value)
    if (unknownFields.isNotEmpty()) protoMarshal.writeUnknownFields(unknownFields)
}

private fun GroupOpReq.ParamsEntry.Companion.protoUnmarshalImpl(protoUnmarshal: pbandk.Unmarshaller): GroupOpReq.ParamsEntry {
    var key = ""
    var value = ""
    while (true) when (protoUnmarshal.readTag()) {
        0 -> return GroupOpReq.ParamsEntry(key, value, protoUnmarshal.unknownFields())
        10 -> key = protoUnmarshal.readString()
        18 -> value = protoUnmarshal.readString()
        else -> protoUnmarshal.unknownField()
    }
}

private fun GroupOpResult.protoMergeImpl(plus: GroupOpResult?): GroupOpResult = plus?.copy(
    reqMem = reqMem?.plus(plus.reqMem) ?: plus.reqMem,
    group = group?.plus(plus.group) ?: plus.group,
    members = members + plus.members,
    params = params + plus.params,
    unknownFields = unknownFields + plus.unknownFields
) ?: this

private fun GroupOpResult.protoSizeImpl(): Int {
    var protoSize = 0
    if (operation.value != 0) protoSize += pbandk.Sizer.tagSize(1) + pbandk.Sizer.enumSize(operation)
    if (reqMem != null) protoSize += pbandk.Sizer.tagSize(2) + pbandk.Sizer.messageSize(reqMem)
    if (result.isNotEmpty()) protoSize += pbandk.Sizer.tagSize(3) + pbandk.Sizer.stringSize(result)
    if (detail.isNotEmpty()) protoSize += pbandk.Sizer.tagSize(4) + pbandk.Sizer.stringSize(detail)
    if (group != null) protoSize += pbandk.Sizer.tagSize(5) + pbandk.Sizer.messageSize(group)
    if (sendId != 0L) protoSize += pbandk.Sizer.tagSize(6) + pbandk.Sizer.int64Size(sendId)
    if (msgId != 0L) protoSize += pbandk.Sizer.tagSize(7) + pbandk.Sizer.int64Size(msgId)
    if (members.isNotEmpty()) protoSize += (pbandk.Sizer.tagSize(8) * members.size) + members.sumBy(pbandk.Sizer::messageSize)
    if (params.isNotEmpty()) protoSize += pbandk.Sizer.mapSize(9, params, model.GroupOpResult::ParamsEntry)
    protoSize += unknownFields.entries.sumBy { it.value.size() }
    return protoSize
}

private fun GroupOpResult.protoMarshalImpl(protoMarshal: pbandk.Marshaller) {
    if (operation.value != 0) protoMarshal.writeTag(8).writeEnum(operation)
    if (reqMem != null) protoMarshal.writeTag(18).writeMessage(reqMem)
    if (result.isNotEmpty()) protoMarshal.writeTag(26).writeString(result)
    if (detail.isNotEmpty()) protoMarshal.writeTag(34).writeString(detail)
    if (group != null) protoMarshal.writeTag(42).writeMessage(group)
    if (sendId != 0L) protoMarshal.writeTag(48).writeInt64(sendId)
    if (msgId != 0L) protoMarshal.writeTag(56).writeInt64(msgId)
    if (members.isNotEmpty()) members.forEach { protoMarshal.writeTag(66).writeMessage(it) }
    if (params.isNotEmpty()) protoMarshal.writeMap(74, params, model.GroupOpResult::ParamsEntry)
    if (unknownFields.isNotEmpty()) protoMarshal.writeUnknownFields(unknownFields)
}

private fun GroupOpResult.Companion.protoUnmarshalImpl(protoUnmarshal: pbandk.Unmarshaller): GroupOpResult {
    var operation: model.GroupOperationType = model.GroupOperationType.fromValue(0)
    var reqMem: model.GroupMember? = null
    var result = ""
    var detail = ""
    var group: model.GroupInfo? = null
    var sendId = 0L
    var msgId = 0L
    var members: pbandk.ListWithSize.Builder<model.GroupMember>? = null
    var params: pbandk.MessageMap.Builder<String, String>? = null
    while (true) when (protoUnmarshal.readTag()) {
        0 -> return GroupOpResult(operation, reqMem, result, detail,
            group, sendId, msgId, pbandk.ListWithSize.Builder.fixed(members),
            pbandk.MessageMap.Builder.fixed(params), protoUnmarshal.unknownFields())
        8 -> operation = protoUnmarshal.readEnum(model.GroupOperationType.Companion)
        18 -> reqMem = protoUnmarshal.readMessage(model.GroupMember.Companion)
        26 -> result = protoUnmarshal.readString()
        34 -> detail = protoUnmarshal.readString()
        42 -> group = protoUnmarshal.readMessage(model.GroupInfo.Companion)
        48 -> sendId = protoUnmarshal.readInt64()
        56 -> msgId = protoUnmarshal.readInt64()
        66 -> members = protoUnmarshal.readRepeatedMessage(members, model.GroupMember.Companion, true)
        74 -> params = protoUnmarshal.readMap(params, model.GroupOpResult.ParamsEntry.Companion, true)
        else -> protoUnmarshal.unknownField()
    }
}

private fun GroupOpResult.ParamsEntry.protoMergeImpl(plus: GroupOpResult.ParamsEntry?): GroupOpResult.ParamsEntry = plus?.copy(
    unknownFields = unknownFields + plus.unknownFields
) ?: this

private fun GroupOpResult.ParamsEntry.protoSizeImpl(): Int {
    var protoSize = 0
    if (key.isNotEmpty()) protoSize += pbandk.Sizer.tagSize(1) + pbandk.Sizer.stringSize(key)
    if (value.isNotEmpty()) protoSize += pbandk.Sizer.tagSize(2) + pbandk.Sizer.stringSize(value)
    protoSize += unknownFields.entries.sumBy { it.value.size() }
    return protoSize
}

private fun GroupOpResult.ParamsEntry.protoMarshalImpl(protoMarshal: pbandk.Marshaller) {
    if (key.isNotEmpty()) protoMarshal.writeTag(10).writeString(key)
    if (value.isNotEmpty()) protoMarshal.writeTag(18).writeString(value)
    if (unknownFields.isNotEmpty()) protoMarshal.writeUnknownFields(unknownFields)
}

private fun GroupOpResult.ParamsEntry.Companion.protoUnmarshalImpl(protoUnmarshal: pbandk.Unmarshaller): GroupOpResult.ParamsEntry {
    var key = ""
    var value = ""
    while (true) when (protoUnmarshal.readTag()) {
        0 -> return GroupOpResult.ParamsEntry(key, value, protoUnmarshal.unknownFields())
        10 -> key = protoUnmarshal.readString()
        18 -> value = protoUnmarshal.readString()
        else -> protoUnmarshal.unknownField()
    }
}
