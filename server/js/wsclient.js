class WsClient {
    constructor(name, url, messageCallback, openCallback, closeCallback, errorCallback, progressCallback) {
        this.name = name;
        this.url = url;
        this.webSocket = null; // WebSocket 对象
        this.keyPairEx = null; // 密钥对对象
        this.messageCallback = messageCallback; // 消息回调函数
        this.openCallback = openCallback; // 连接成功回调函数
        this.closeCallback = closeCallback; // 连接关闭回调函数
        this.errorCallback = errorCallback; // 连接错误回调函数
        this.progressCallback = progressCallback;
    }

    // 连接 WebSocket
    connect() {
        this.webSocket = new WebSocket(this.url);

        // WebSocket 连接成功时的回调函数
        this.webSocket.onopen = () => {
            console.log(`WebSocket connected to ${this.url}`);
            if (this.openCallback) {
                this.openCallback();
            }
        };

        // WebSocket 接收到消息时的回调函数
        this.webSocket.onmessage = (event) => {
            console.log(`Received message: ${event.data}`);

            //showMessage("收到数据");
            if (typeof event.data === 'string') {
                this.progressCallback(event.data);
            }
            if (event.data instanceof Blob) {
                // 如果接收到的是二进制数据
                const reader = new FileReader();
                const self = this;
                reader.onload = function () {
                    // 读取数据完成后的回调函数
                    const binaryData = reader.result; // 二进制数据
                    // 在这里对二进制数据进行处理，可以将其解析为特定的格式
                    self.parseData(binaryData);
                };
                reader.readAsArrayBuffer(event.data); // 以ArrayBuffer格式读取二进制数据
            }

            // if (this.messageCallback) {
            //     this.messageCallback(event.data);
            // }
        };

        // WebSocket 连接关闭时的回调函数
        this.webSocket.onclose = () => {
            console.log(`WebSocket disconnected from ${this.url}`);
            if (this.closeCallback) {
                this.closeCallback();
            }
        };

        // WebSocket 连接发生错误时的回调函数
        this.webSocket.onerror = (error) => {
            console.error(`WebSocket error: ${error}`);
            if (this.errorCallback) {
                this.errorCallback(error);
            }
        };
    }

    // 发送消息
    sendWs(message) {
        if (this.webSocket && this.webSocket.readyState === WebSocket.OPEN) {
            this.webSocket.send(message);
            console.log(`Sent message: ${message}`);
        } else {
            console.error("WebSocket is not connected or ready to send messages.");
        }
    }

    isOpen(){
        return this.webSocket && this.webSocket.readyState === WebSocket.OPEN;
    }

    // 关闭连接
    disconnect() {
        if (this.webSocket) {
            this.webSocket.close();
        }
    }

    // 分发数据的核心函数
    parseData(binData){
        // 获取收到的二进制数据并转换为 Uint8Array
        const data = new Uint8Array(binData);
        this.progressCallback(data.length)
        printByteArray(data)

        // 解析 Msg 消息
        const msg = proto.model.Msg.deserializeBinary(data);

        // 获取 Msg 中的消息类型
        const msgType = msg.getMsgtype();
        const version = msg.getVersion();
        // this.progressCallback(msgType);

        // 根据消息类型进行相应处理
        var str = "";
        switch (msgType) {
            case proto.model.ComMsgType.MSGTHELLO:
                // 如果是 Hello 消息
                const helloMsg = msg.getPlainmsg().getHello();

                // 将 Hello 消息内容显示在页面上


                str += "Received Hello message:\n" +
                    "Stage: " + helloMsg.getStage() + "\n" +
                    "Server Version: " + helloMsg.getVersion() + "\n" +
                    "Platform: " + helloMsg.getPlatform() + "\n";
                this.progressCallback(str);
                break;
            case proto.model.ComMsgType.MSGTERROR:
                // 如果是 Error 消息
                var errorMsg = msg.getPlainmsg().getErrormsg();
                // 将 Error 消息内容显示在页面上
                str += "Received Error message:\n" +
                    "Code: " + errorMsg.getCode() + "\n" +
                    "Detail: " + errorMsg.getDetail() + "\n";
               //this.progressCallback(str);
                this.errorCallback(str)
                break;
            case proto.model.ComMsgType.MSGTKEYEXCHANGE:
                break;

            default:
                // 其他类型的消息处理
                console.warn("Received unknown message type:", msgType);
                break;
        }
    }

    sendObject(msg){
        this.progressCallback(msg.toString()) // 字符串表示
        const binMsg = msg.serializeBinary();
        //this.progressCallback(binMsg)
        //var jsonStr =  JSON.stringify(msg);
        //showMessage(jsonStr);
        this.sendWs(binMsg);
    }

    // 1.0 发送hello包
    sendHello(){
        // 子消息
        const hello = new proto.model.MsgHello();
        hello.setClientid("uuid")  //js权限低，这里可以使用一个UUID，执行本地存储，每次都带着，用于服务端区分设备
        hello.setVersion("1.0")
        hello.setPlatform("web")
        hello.setStage("clienthello")

        // 将 MsgHello 消息设置为 Msg 消息的 plainMsg 字段
        const plainMsg = new proto.model.MsgPlain();
        plainMsg.setHello(hello);


        // 封装为通用消息
        const msg = new proto.model.Msg();
        msg.setMsgtype(proto.model.ComMsgType.MSGTHELLO);
        msg.setVersion(1);
        msg.setPlainmsg(plainMsg);

        this.sendObject(msg);
    }

    // 秘钥交换阶段1: 生成公私钥对，发送公钥过去
    async sendExchange1(){
        this.keyPairEx = await createDHPair();
        const publicKeyRaw = await getPublicKeyRaw(this.keyPairEx);
        console.log("public key:", publicKey);

        const exMsg = new proto.model.MsgKeyExchange();
        exMsg.setKeyprint(0);    // 不使用临时秘钥，也没有旧秘钥
        exMsg.setRsaprint(0);    // 不使用RSA加密
        exMsg.setStage(1);       // 第一次握手
        exMsg.setTempkey(0);     // 不使用RSA加密,也就没有临时秘钥
        exMsg.setPubkey(new Uint8Array(publicKeyRaw)); // 公钥的明文，这里不加密
        exMsg.setEnctype("chacha20");    // status和detail不设置


        const plainMsg = new proto.model.MsgPlain();
        plainMsg.setKeyex(exMsg);

        // 封装为通用消息
        const msg = new proto.model.Msg();
        msg.setMsgtype(proto.model.ComMsgType.MSGTKEYEXCHANGE);
        msg.setVersion(1);
        msg.setPlainmsg(plainMsg);

        this.sendObject(msg);
    }

    // 秘钥交换阶段2
    sendExchange2(){

    }

    // 2
    sendHeart(){
        const heart = new proto.model.MsgHeartBeat();
        heart.setUserid(1);
        const timestamp = new Date().getTime();
        heart.setTm(timestamp);


        // const plainMsg = new proto.model.MsgPlain()
        // plainMsg.setHeartbeat(heart)

        // 封装为通用消息
        const msg = new proto.model.Msg();
        msg.setVersion(1);
        msg.setMsgtype(proto.model.ComMsgType.MSGTHEARTBEAT);

        msg.setPlainmsg(heart)

        this.sendObject(msg);
    }

    // 3注册申请
    sendRegisterMessage() {
        showMessage("注册申请的消息")
        const userInfo = new proto.model.UserInfo();
        userInfo.setUsername("robin");
        userInfo.setUserid(0);
        userInfo.setEmail("390017268@qq.com");
        const paramsMap = userInfo.getParamsMap();
        paramsMap.set("pwd", "12345678");
        paramsMap.set("regtype", "pwd");   // pwd, email,phone


        // 注册用户
        const regOpReq = new proto.model.UserOpReq();
        regOpReq.setOperation(proto.model.UserOperationType.REGISTERUSER);
        regOpReq.setUser(userInfo)


        // 封装为通用消息
        const msg = new proto.model.Msg();
        msg.setMsgtype(proto.model.ComMsgType.MSGTUSEROP);
        msg.setVersion(1);
        msg.setUserop(regOpReq)

        this.sendObject(msg);
    }

    // 4. 发送验证码，所有发送验证码的都是一样的，服务端跟踪当前
    sendcodeMessage() {
        showMessage("发送验证消息")
        const userInfo = new proto.model.UserInfo();
        userInfo.setUserid(1001);
        const paramsMap = userInfo.getParamsMap();
        paramsMap.set("code", "12345");
        paramsMap.set("regtype", "email");   // pwd, email,phone


        // 验证码
        const regOpReq = new proto.model.UserOpReq();
        regOpReq.setOperation(proto.model.UserOperationType.REALNAMEVERIFICATION);
        regOpReq.setUser(userInfo);


        // 封装为通用消息
        const msg = new proto.model.Msg();
        msg.setMsgtype(proto.model.ComMsgType.MSGTUSEROP);
        msg.setVersion(1);
        msg.setUserop(regOpReq)

        this.sendObject(msg);

    }

    // 5.1 用户密码登录申请
    sendLoginPwdMessage(){
        showMessage("发送登录消息")
        const userInfo = new proto.model.UserInfo();
        userInfo.setUserid(1001);
        const paramsMap = userInfo.getParamsMap();
        paramsMap.set("pwd", "12345");
        paramsMap.set("logintype", "pwd");   // uidpwd, email,phone


        // 登录
        const regOpReq = new proto.model.UserOpReq();
        regOpReq.setOperation(proto.model.UserOperationType.LOGIN);
        regOpReq.setUser(userInfo);


        // 封装为通用消息
        const msg = new proto.model.Msg();
        msg.setMsgtype(proto.model.ComMsgType.MSGTUSEROP);
        msg.setVersion(1);
        msg.setUserop(regOpReq)

        this.sendObject(msg);

    }

    // 5.2 用户邮箱登录申请
    sendLoginEmailMessage(){
        showMessage("发送登录消息")
        const userInfo = new proto.model.UserInfo();
        //userInfo.setUserid(1001);
        userInfo.setEmail("390017268@qq.com")
        const paramsMap = userInfo.getParamsMap();
        paramsMap.set("logintype", "email");   // uid, email,phone


        // 登录
        const regOpReq = new proto.model.UserOpReq();
        regOpReq.setOperation(proto.model.UserOperationType.LOGIN);
        regOpReq.setUser(userInfo);


        // 封装为通用消息
        const msg = new proto.model.Msg();
        msg.setMsgtype(proto.model.ComMsgType.MSGTUSEROP);
        msg.setVersion(1);
        msg.setUserop(regOpReq)

        this.sendObject(msg);
    }

    // 5.3 邮箱登录验证码
    sendLoginCodeMessage(){
        showMessage("发送验证码消息")
        const userInfo = new proto.model.UserInfo();
        //userInfo.setUserid(1001);
        userInfo.setEmail("390017268@qq.com")
        const paramsMap = userInfo.getParamsMap();
        paramsMap.set("logintype", "email");   // uid, email,phone
        paramsMap.set("code", "12345");


        // 登录
        const regOpReq = new proto.model.UserOpReq();
        regOpReq.setOperation(proto.model.UserOperationType.REALNAMEVERIFICATION);
        regOpReq.setUser(userInfo);


        // 封装为通用消息
        const msg = new proto.model.Msg();
        msg.setMsgtype(proto.model.ComMsgType.MSGTUSEROP);
        msg.setVersion(1);
        msg.setUserop(regOpReq)

        this.sendObject(msg);
    }

    // 6.1 设置用户的各种信息，如果重新设置邮箱以及手机号，则需要验证码
    sendUserInfoMessage(){
        showMessage("发送用户详细消息")
        const userInfo = new proto.model.UserInfo();
        userInfo.setUserid(1001);
        userInfo.setEmail("");
        userInfo.setGender("male");
        userInfo.setAge(5);
        userInfo.setIcon("sys://1");
        userInfo.setNickname("robinfox");
        userInfo.setUsername("robinfoxnan");
        const paramsMap = userInfo.getParamsMap();
        paramsMap.set("title", "manager");
        paramsMap.set("country", "china");


        // 登录
        const regOpReq = new proto.model.UserOpReq();
        regOpReq.setOperation(proto.model.UserOperationType.SETUSERINFO);
        regOpReq.setUser(userInfo);


        // 封装为通用消息
        const msg = new proto.model.Msg();
        msg.setMsgtype(proto.model.ComMsgType.MSGTUSEROP);
        msg.setVersion(1);
        msg.setUserop(regOpReq)

        this.sendObject(msg);
    }

    // 6.2
    sendUserInfoCodeMessage(){
        showMessage("发送用户详细消息")
        const userInfo = new proto.model.UserInfo();
        userInfo.setUserid(1001);
        userInfo.setEmail("");
        userInfo.setGender("male");
        userInfo.setAge(5);
        userInfo.setIcon("sys://1");
        userInfo.setNickname("robinfox");
        userInfo.setUsername("robinfoxnan");
        const paramsMap = userInfo.getParamsMap();
        paramsMap.set("title", "manager");
        paramsMap.set("country", "china");


        // 登录
        const regOpReq = new proto.model.UserOpReq();
        regOpReq.setOperation(proto.model.UserOperationType.REALNAMEVERIFICATION);
        regOpReq.setUser(userInfo);


        // 封装为通用消息
        const msg = new proto.model.Msg();
        msg.setMsgtype(proto.model.ComMsgType.MSGTUSEROP);
        msg.setVersion(1);
        msg.setUserop(regOpReq)

        this.sendObject(msg);
    }

    // 7.1 禁用
    sendUserDisableMessage(){
        showMessage("禁用用户")
        const userInfo = new proto.model.UserInfo();
        userInfo.setUserid(1001);

        // 禁用
        const regOpReq = new proto.model.UserOpReq();
        regOpReq.setOperation(proto.model.UserOperationType.DISABLEUSER);
        regOpReq.setUser(userInfo);


        // 封装为通用消息
        const msg = new proto.model.Msg();
        msg.setMsgtype(proto.model.ComMsgType.MSGTUSEROP);
        msg.setVersion(1);
        msg.setUserop(regOpReq)

        this.sendObject(msg);
    }

    // 7.2 解禁
    sendUserEnableMessage(){
        showMessage("禁用用户")
        const userInfo = new proto.model.UserInfo();
        userInfo.setUserid(1001);

        // 解禁
        const regOpReq = new proto.model.UserOpReq();
        regOpReq.setOperation(proto.model.UserOperationType.RECOVERUSER);
        regOpReq.setUser(userInfo);


        // 封装为通用消息
        const msg = new proto.model.Msg();
        msg.setMsgtype(proto.model.ComMsgType.MSGTUSEROP);
        msg.setVersion(1);
        msg.setUserop(regOpReq)

        this.sendObject(msg);
    }

    // 8  注销
    sendUserUnregMessage(){
        showMessage("注销用户")
        const userInfo = new proto.model.UserInfo();
        userInfo.setUserid(1001);

        //
        const regOpReq = new proto.model.UserOpReq();
        regOpReq.setOperation(proto.model.UserOperationType.UNREGISTERUSER);
        regOpReq.setUser(userInfo);


        // 封装为通用消息
        const msg = new proto.model.Msg();
        msg.setMsgtype(proto.model.ComMsgType.MSGTUSEROP);
        msg.setVersion(1);
        msg.setUserop(regOpReq)

        this.sendObject(msg);
    }

    // 9 退出登录
    sendUserLogoutMessage(){
        showMessage("登出用户")
        const userInfo = new proto.model.UserInfo();
        userInfo.setUserid(1001);

        //
        const regOpReq = new proto.model.UserOpReq();
        regOpReq.setOperation(proto.model.UserOperationType.LOGOUT);
        regOpReq.setUser(userInfo);


        // 封装为通用消息
        const msg = new proto.model.Msg();
        msg.setMsgtype(proto.model.ComMsgType.MSGTUSEROP);
        msg.setVersion(1);
        msg.setUserop(regOpReq)

        this.sendObject(msg);
    }

}

// 创建 WsClient 实例并连接 WebSocket，并传递回调函数
// const client = new WsClient(
//     "MyClient",
//     "ws://localhost:8080",
//     (message) => console.log("Received message callback:", message),
//     () => console.log("WebSocket connected"),
//     () => console.log("WebSocket disconnected"),
//     (error) => console.error("WebSocket error:", error)
// );
// client.connect();
//
// // 发送消息示例
// client.sendMessage("Hello, WebSocket!");

// 关闭连接示例
// client.disconnect();
