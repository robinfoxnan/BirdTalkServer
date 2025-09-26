// 生成密钥对，并保存到全局变量中
async function createDHPair() {
    // 生成新的ECDH密钥对
    const keyPair = await window.crypto.subtle.generateKey(
        {
            name: "ECDH",
            namedCurve: "P-256"
        },
        true, // 允许导出私钥
        ["deriveKey", "deriveBits"] // 导出私钥的权限
    );
    return keyPair;
};

function bufferToBase64(buffer) {
    const binary = String.fromCharCode.apply(null, new Uint8Array(buffer));
    return btoa(binary);
}

async function getPEMPublicKey(keyPair) {
    const publicKeyBuffer = await crypto.subtle.exportKey("spki", keyPair.publicKey);
    const publicKeyBase64 = bufferToBase64(publicKeyBuffer);
    const pemKey = `-----BEGIN PUBLIC KEY-----\n${publicKeyBase64}\n-----END PUBLIC KEY-----\n`;
    return pemKey;
};

function stringToBytes(str) {
    const encoder = new TextEncoder();
    return encoder.encode(str);
};

function bytesToString(bytes) {
    const decoder = new TextDecoder();
    return decoder.decode(bytes);
}

function getCurrentTimestamp() {
    const timestamp = BigInt(Date.now());
    return timestamp.toString();
}


async function calculateSharedSecret(privateKey, publicKey) {
    // 执行Diffie-Hellman密钥交换以计算共享密钥
    const sharedSecret = await window.crypto.subtle.deriveBits(
        {
            name: "ECDH",
            namedCurve: "P-256",
            public: publicKey // 使用对方的公钥
        },
        privateKey, // 使用自己的私钥
        256,
    );
    return new Uint8Array(sharedSecret);
}

function stringToArrayBuffer(str) {
    const encoder = new TextEncoder();
    return encoder.encode(str);
}

async function importPublicKey(pemKey) {

    // 去除 PEM 头尾，并将字符串解码为 ArrayBuffer
    const pemContents = pemKey.replace(/-----BEGIN PUBLIC KEY-----/, '').replace(/-----END PUBLIC KEY-----/, '');
    const keyStr = pemContents.trim();
    console.log("remote der key:", keyStr);
    const pemArrayBuffer = Uint8Array.from(atob(keyStr), c => c.charCodeAt(0));
    console.log("remote der key:", pemArrayBuffer)

    // 将原始公钥字节数组导入为 CryptoKey 对象
    const publicKey = await crypto.subtle.importKey(
        "spki", // der
        pemArrayBuffer, // 原始公钥字节数组
        {
            name: "ECDH", // 算法名称
            namedCurve: "P-256" // 曲线名称
        },
        true, // 是否为公钥
        [] // 公钥的用途
    );
    return publicKey;
}

// 保存共享密钥到本地存储
function saveSharedKey(sharedKey) {
    // 将 Uint8Array 转换为字符串
    const sharedKeyString = sharedKey.join(',');

    // 存储共享密钥字符串到本地存储
    localStorage.setItem("birdSharedKey", sharedKeyString);
}

// 保存字符串
function saveShareKeyPrint(print){
    localStorage.setItem("birdSharedKeyPrint", print);
}

// 从本地存储加载共享密钥
function loadSharedKey() {
    // 从本地存储中获取共享密钥字符串
    const sharedKeyString = localStorage.getItem("birdSharedKey");

    if (sharedKeyString) {
        // 将共享密钥字符串转换为 Uint8Array
        const sharedKeyArray = sharedKeyString.split(',').map(Number);
        const sharedKey = new Uint8Array(sharedKeyArray);
        return sharedKey;
    } else {
        // 如果本地存储中没有共享密钥，则返回 null
        return null;
    }
}

function loadSharedKeyPrint() {
    return localStorage.getItem("birdSharedKeyPrint");
}

function deleteShareKey() {
    localStorage.removeItem("birdSharedKeyPrint");
    localStorage.removeItem("birdSharedKey");
}
//////////////////////////////////////////////////////////////
// 将字节数组转换为 int64, 取指纹
function bytesToInt64(data) {
    // 检查字节数组长度是否足够
    if (data.length < 8) {
        throw new Error("Insufficient bytes to convert to int64");
    }

    // 将字节数组转换为 int64
    const view = new DataView(new ArrayBuffer(8));
    for (let i = 0; i < 8; i++) {
        view.setUint8(i, data[i]);
    }
    const int64Value = view.getBigInt64(0, true); // 使用 little-endian 格式

    return int64Value.toString();
}



// async function encryptBytes(data, key) {
//     try {
//         // 生成随机的初始化向量（IV）
//         const iv = crypto.getRandomValues(new Uint8Array(12));
//
//         // 使用 CryptoKey 对象创建 ChaCha20 算法对象
//         const algorithm = { name: "AES-GCM", iv: iv };
//
//         // 使用 CryptoKey 对象加密数据
//         const encryptedData = await crypto.subtle.encrypt(algorithm, key, dataBuffer);
//
//         // 将 IV 和加密后的数据拼接在一起
//         const result = new Uint8Array(iv.length + encryptedData.byteLength);
//         result.set(iv, 0);
//         result.set(new Uint8Array(encryptedData), iv.length);
//
//         // 返回拼接后的数据
//         return result;
//     } catch (error) {
//         console.error("Encryption failed:", error);
//         throw error;
//     }
// };

// 加密整数
async function encryptInt64(int64Value, key) {
    // 将 int64 值转换为字节数组
    const dataBuffer = new ArrayBuffer(8); // 64位整数字节长度为8
    const view = new DataView(dataBuffer);
    view.setBigInt64(0, BigInt(int64Value));
    return encryptBytes(dataBuffer, key);
};

// 解密整数
async function decryptInt64AES(encryptedData, key) {
    try {
        // 使用 CryptoKey 对象创建 AES-GCM 算法对象
        const algorithm = { name: "AES-GCM", iv: encryptedData.slice(0, 12) };

        // 使用 CryptoKey 对象解密数据
        const decryptedData = await crypto.subtle.decrypt(algorithm, key, encryptedData.slice(12));

        // 将解密后的字节数组转换回 int64 值
        const view = new DataView(decryptedData);
        const decryptedInt64Value = view.getBigInt64(0);

        // 返回解密后的 int64 值
        return decryptedInt64Value.toString();
    } catch (error) {
        console.error("Decryption failed:", error);
        throw error;
    }
};

async function decryptInt64(encryptedData, key) {
    try {
        // 使用 CryptoKey 对象创建 ChaCha20 算法对象
        const algorithm = {
            name: "ChaCha20",
            iv: encryptedData.slice(0, 12) // ChaCha20 的初始向量长度为 12 字节
        };

        // 将原始密钥导入为 CryptoKey 对象
        // 将原始密钥导入为 CryptoKey 对象
        const importedKey = await crypto.subtle.importKey(
            "raw", // 密钥格式为原始数据
            key,   // 原始密钥字节数组
            { name: "ChaCha20" }, // 使用 ChaCha20 算法
            true, // 是否允许导出密钥
            ["decrypt"] // 密钥用途为解密
        );

        // 使用 CryptoKey 对象解密数据
        const decryptedData = await crypto.subtle.decrypt(algorithm, importedKey, encryptedData.slice(12));

        // 将解密后的字节数组转换为 int64 值
        const view = new DataView(decryptedData);
        const decryptedInt64Value = view.getBigInt64(0);

        // 返回解密后的 int64 值
        return decryptedInt64Value.toString();
    } catch (error) {
        console.error("Decryption failed:", error);
        throw error;
    }
}

// 将ArrayBuffer转换为Base64字符串
function arrayBufferToBase64(buffer) {
    var binary = '';
    var bytes = new Uint8Array(buffer);
    var len = bytes.byteLength;
    for (var i = 0; i < len; i++) {
        binary += String.fromCharCode(bytes[i]);
    }
    return window.btoa(binary);
};


// 预计算的CRC32表
const crcTable = (function () {
    let table = new Array(256);
    let c;
    for (let n = 0; n < 256; n++) {
        c = n;
        for (let k = 0; k < 8; k++) {
            c = ((c & 1) ? (0xEDB88320 ^ (c >>> 1)) : (c >>> 1));
        }
        table[n] = c;
    }
    return table;
})();

function crc32(input) {
    let crc = -1;
    for (let i = 0; i < input.length; i++) {
        const byte = input[i];
        crc = (crc >>> 8) ^ crcTable[(crc ^ byte) & 0xFF];
    }
    // 返回无符号32位整数
    return (crc ^ -1) >>> 0;
}



// 计算 SHA256 哈希
async function calculateSHA256(data) {
    try {
        // 将输入数据转换为 ArrayBuffer
        const buffer = new TextEncoder().encode(data);

        // 计算 SHA-256 哈希
        const hashBuffer = await crypto.subtle.digest('SHA-256', buffer);

        // 将 ArrayBuffer 转换为十六进制字符串
        const hashArray = Array.from(new Uint8Array(hashBuffer));
        const hashHex = hashArray.map(byte => byte.toString(16).padStart(2, '0')).join('');

        return hashHex;
    } catch (error) {
        console.error('SHA256 calculation error:', error);
        throw error;
    }
}

// Uint8Array
async function calculateSHA256Raw(data) {
    try {
        // 计算 SHA-256 哈希
        const hashBuffer = await crypto.subtle.digest('SHA-256', data);
        return hashBuffer;
    } catch (error) {
        console.error('SHA256 calculation error:', error);
        throw error;
    }
}


// 计算 MD5 哈希
async function calculateMD5(data) {
    try {
        // 将输入数据转换为 ArrayBuffer
        const buffer = new TextEncoder().encode(data);

        // 计算 MD5 哈希
        const hashBuffer = await crypto.subtle.digest('MD5', buffer);

        // 将 ArrayBuffer 转换为十六进制字符串
        const hashArray = Array.from(new Uint8Array(hashBuffer));
        const hashHex = hashArray.map(byte => byte.toString(16).padStart(2, '0')).join('');

        return hashHex;
    } catch (error) {
        console.error('MD5 calculation error:', error);
        throw error;
    }
};

async function encryptString(str, key){
    // 将文本编码为 ArrayBuffer
    const dataBuffer = new TextEncoder().encode(str);
    const temp = await encryptAES_CTR(dataBuffer, key);
    return arrayBufferToBase64(temp);
}

async function encryptAES_CTR(plaintext, key) {
    // 生成随机的初始化向量（IV）
    const iv = crypto.getRandomValues(new Uint8Array(16)); // 假设初始化向量长度为 16 字节

    // 使用 CryptoKey 对象创建 AES-CTR 算法对象
    const algorithm = {
        name: "AES-CTR",
        counter: iv, // 使用生成的随机初始化向量作为计数器
        length: 64 // 可选参数，指定计数器长度，默认为 128
    };

    // 将原始密钥导入为 CryptoKey 对象
    const importedKey = await crypto.subtle.importKey(
        "raw", // 密钥格式为原始数据
        key,   // 原始密钥字节数组
        { name: "AES-CTR", length: 256 }, // AES 密钥长度为 256 位
        true, // 是否允许导出密钥
        ["encrypt", "decrypt"] // 密钥用途为加密和解密
    );

    // 使用 CryptoKey 对象加密数据
    const encryptedData = await crypto.subtle.encrypt(algorithm, importedKey, plaintext);

    // 将随机初始化向量和加密后的数据拼接在一起
    const ciphertext = new Uint8Array(iv.length + encryptedData.byteLength);
    ciphertext.set(iv, 0); // 将初始化向量放在密文前面
    ciphertext.set(new Uint8Array(encryptedData), iv.length); // 将加密后的数据放在初始化向量后面

    return ciphertext;
}

async function decryptAES_CTR(ciphertext, key) {
    // 从密文中提取初始化向量（IV）
    const iv = ciphertext.slice(0, 16); // 假设初始化向量长度为 16 字节

    // 使用 CryptoKey 对象创建 AES-CTR 算法对象
    const algorithm = {
        name: "AES-CTR",
        counter: iv, // 使用提取的初始化向量作为计数器
        length: 64 // 可选参数，指定计数器长度，默认为 128
    };

    // 将原始密钥导入为 CryptoKey 对象
    const importedKey = await crypto.subtle.importKey(
        "raw", // 密钥格式为原始数据
        key,   // 原始密钥字节数组
        { name: "AES-CTR", length: 256 }, // AES 密钥长度为 256 位
        true, // 是否允许导出密钥
        ["encrypt", "decrypt"] // 密钥用途为加密和解密
    );

    // 使用 CryptoKey 对象解密数据
    const plaintext = await crypto.subtle.decrypt(algorithm, importedKey, ciphertext.slice(16)); // 去掉头部的 IV

    return plaintext;
}
///////////////////////////////////////////////////////////////////////////////////////
class FileUploader{
    constructor(ws, file, cz){
        this.ws = ws;
        this.file = file;
        this.chunkSize = cz;
        this.currentIndex = 0;
        this.totalChunks = Math.ceil(file.size / this.chunkSize);
        this.hashCode = "";
    }

    // 计算MD5
    async calculateFileMd5Hash(file) {
        return new Promise((resolve, reject) => {
            if (file) {
                const reader = new FileReader();

                reader.onload = function(e) {
                    const arrayBuffer = e.target.result;
                    const spark = new SparkMD5.ArrayBuffer();
                    spark.append(arrayBuffer);
                    const md5Hash = spark.end();
                    resolve(md5Hash);
                };

                reader.onerror = function() {
                    reject(new Error("File reading has failed"));
                };

                reader.readAsArrayBuffer(file);
            } else {
                reject(new Error("No file provided"));
            }
        });
    }

    async readFileChunkAsArrayBuffer(chunk) {
        return new Promise((resolve, reject) => {
            const reader = new FileReader();
            reader.onload = function(e) {
                resolve(e.target.result);
            };
            reader.onerror = function() {
                reject(new Error("Chunk reading has failed"));
            };
            reader.readAsArrayBuffer(chunk);
        });
    }
    // 开始传递第一块，看看能够秒传
    async uploadFirstTrunk(){
        this.hashCode = await this.calculateFileMd5Hash(this.file);
        const start = 0;
        const end = Math.min(start + this.chunkSize, this.file.size);
        const chunk = this.file.slice(start, end);
        try {
            const chunkBuffer = await this.readFileChunkAsArrayBuffer(chunk);
            // Convert ArrayBuffer to Uint8Array
            const chunkUint8Array = new Uint8Array(chunkBuffer);
            this.ws.sendUploadMessage(this.file.name, "file", this.file.size, this.chunkSize,
                this.totalChunks, 0,  chunkUint8Array, this.hashCode);
        } catch (error) {
            console.error('Error processing chunk:', error);

            return;
        }
    }

    // 如果不能秒传，传递其他的部分
    async uploadOther(index){
        this.currentIndex = index;
        const start = index * this.chunkSize;
        const end = Math.min(start + this.chunkSize, this.file.size);
        const chunk = this.file.slice(start, end);
        try {
            const chunkBuffer = await this.readFileChunkAsArrayBuffer(chunk);
            // Convert ArrayBuffer to Uint8Array
            const chunkUint8Array = new Uint8Array(chunkBuffer);
            this.ws.sendUploadMessage(this.file.name, "file", this.file.size, this.chunkSize,
                this.totalChunks, index,  chunkUint8Array, this.hashCode);
        } catch (error) {
            console.error('Error processing chunk:', error);

            return;
        }
    }
}
///////////////////////////////////////////////////////////////////////////////////////
//
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
        this.shareKeyPrint = "";
        this.shareKey = null;
        this.fileMap = new Map();  // 上传的开始
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
            //console.log(`Sent message: ${message}`);
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
        //console.warn("recv message type:", msgType);
        var str = "";
        switch (msgType) {
            case proto.model.ComMsgType.MSGTHELLO:
                this.onHello(msg);
                break;
            case proto.model.ComMsgType.MSGTERROR:
                this.onReplyError(msg);
                break;
            case proto.model.ComMsgType.MSGTKEYEXCHANGE:
                this.onKeyExchangeMsg(msg);
                break;

            case proto.model.ComMsgType.MSGTUSEROPRET:  // 用户
                this.onUserOpResult(msg);
                break;

            case proto.model.ComMsgType.MSGTFRIENDOP:  // 好友申请
                break;
            case proto.model.ComMsgType.MSGTFRIENDOPRET: //好友操作应答
                this.onFriendOpResult(msg);
                break;
            case proto.model.ComMsgType.MSGTGROUPOP:  // 需要管理员确认
                break;
            case proto.model.ComMsgType.MSGTGROUPOPRET: // 群操作应答
                this.onGroupOpResult(msg);
                break;

            case proto.model.ComMsgType.MSGTUPLOADREPLY:
                this.onUploadResult(msg);
                break;

            case proto.model.ComMsgType.MSGTCHATMSG:
                this.onChatMsg(msg);
                break;
            case proto.model.ComMsgType.MSGTCHATREPLY:
                this.onChatMsgReply(msg);
                break;

            case proto.model.ComMsgType.MSGTQUERYRESULT:
                this.onQueryResult(msg)
                break;
            default:
                // 其他类型的消息处理
                console.warn("Received unknown message type:", msgType);
                break;
        }
    }

    sendObject(msg){
        this.progressCallback("发送消息") ;// 字符串表示
        const binMsg = msg.serializeBinary();
        //this.progressCallback(binMsg)
        //var jsonStr =  JSON.stringify(msg);
        //showMessage(jsonStr);
        this.sendWs(binMsg);
    }

    // 服务器应答错误
    onReplyError(msg){
        // 如果是 Error 消息
        let str = "";
        const errorMsg = msg.getPlainmsg().getErrormsg();
        // 将 Error 消息内容显示在页面上
        str += "Received Error message:\n" +
            "Code: " + errorMsg.getCode() + "\n" +
            "Detail: " + errorMsg.getDetail() + "\n";
        //this.progressCallback(str);
        switch (errorMsg.getCode()) {
            case proto.model.ErrorMsgType.ERRTKEYPRINT:
                console.log("delete key");
                deleteShareKey();
                break;
            case proto.model.ErrorMsgType.ERRTREDIRECT:
                break;
            case proto.model.ErrorMsgType.ERRTWRONGPWD:
                break;

        }
        this.errorCallback(str)
    }
    // 收到消息应答
    async onHello(msg){
        // 如果是 Hello 消息
        const helloMsg = msg.getPlainmsg().getHello();

        // 将 Hello 消息内容显示在页面上
        let str = "";

        str += "Received Hello message:\n" +
            "Stage: " + helloMsg.getStage() + "\n" +
            "Server Version: " + helloMsg.getVersion() + "\n" +
            "Platform: " + helloMsg.getPlatform() + "\n";
        this.progressCallback(str);
    }

    // 通用消息
    async onKeyExchangeMsg(msg){
        const tm  = msg.getTm();
        console.log("tmstr=", tm);

        const keyMsg = msg.getPlainmsg().getKeyex();
        const stage = keyMsg.getStage();
        if (stage === 2){
            const remoteKeyPrint = keyMsg.getKeyprint();
            console.log("remote key print=", remoteKeyPrint)

            const remotePublicKeyData = keyMsg.getPubkey();
            const remotePublicKeyStr = bytesToString(remotePublicKeyData);
            console.log("remote public key=", remotePublicKeyStr)

            const remotePublicKey = await importPublicKey(remotePublicKeyStr);
            console.log("remote public key=", remotePublicKey)

            const sharedSecretLocal = await calculateSharedSecret(this.keyPairEx.privateKey, remotePublicKey);
            console.log("local share key=", sharedSecretLocal)

            const keyPrint = bytesToInt64(sharedSecretLocal);
            console.log("local key print is ", keyPrint)

            const checkData = keyMsg.getTempkey();
            const checkDataPlain = await decryptAES_CTR(checkData, sharedSecretLocal);
            const checkDataStr = bytesToString(checkDataPlain);
            console.log("decrypt data tm=", checkDataStr);
            if (keyPrint === remoteKeyPrint){
                if (checkDataStr === tm){
                    this.progressCallback("calculate share key ok, check data ok");
                    this.shareKey = sharedSecretLocal;
                    this.shareKeyPrint = keyPrint;
                    saveSharedKey(sharedSecretLocal);  // 保存共享密钥
                    saveShareKeyPrint(keyPrint);
                    this.sendExchange3();              // 发送验证结果

                }else{
                    this.progressCallback("calculate share key ok, check data fail");
                    this.errorCallback("calculate share key error!!");
                }

            }else{
                this.progressCallback("calculate share key error!!");
                // 发送错误应答；
                this.errorCallback("calculate share key error!!");
            }
        }else if (stage === 4){
            // 等待服务器的应答
            const status = keyMsg.getStatus();

            if (status === "needlogin"){
                this.progressCallback("server said that share-key is ok, but need login first.");
            }else if (status === "waitdata") {
                this.progressCallback("server said that share-key is ok, login is not needed.");
            }
        }

    }

    // 所有用户操作的应答
    async onUserOpResult(msg){

        const userOpRet = msg.getPlainmsg().getUseropret();
        const op = userOpRet.getOperation();
        switch (op){
            case proto.model.UserOperationType.REGISTERUSER:
                break;
        }

        let str = "";
        //str = userOpRet.toLocaleString();

        str += "Received user operation result  message:\n" +
            "OP: " + op.toLocaleString() + "\n" +
            "Status: " + userOpRet.getStatus() + "\n" +
            "Result: " + userOpRet.getResult() + "\n" +
            "Users info: " + userOpRet.getUsersList() + "\n";
        const users = userOpRet.getUsersList();
        const user = users[0];
        str += "user id: " + user.getUserid() + "\n" ;
        this.progressCallback(str);
    }

    // 朋友相关消息应答
    async onFriendOpResult(msg){
        const friendOpRet = msg.getPlainmsg().getFriendopret();
        const op = friendOpRet.getOperation();
        switch (op){
            case proto.model.UserOperationType.FINDUSER:
                break;
        }

        let str = "";
        str += "Received friend operation result  message:\n" +
            "OP: " + op.toLocaleString() + "\n" +
            // "Status: " + friendOpRet.getStatus() + "\n" +
            "Result: " + friendOpRet.getResult() + "\n" +
            "User info: " + friendOpRet.getUsersList() + "\n";
        const users = friendOpRet.getUsersList();
        // const user = users[0];
        // str += "user id: " + user.getUserid() + "\n" ;
        this.progressCallback(str);
    }

    // 群聊消息应答
    async onGroupOpResult(msg){
        const grpOpRet = msg.getPlainmsg().getGroupopret();
        const op = grpOpRet.getOperation();
        switch (op){
            case proto.model.GroupOperationType.GROUPCREATE:
                break;
        }

        let str = "";
        str += "Received group operation result  message:\n" +
            "OP: " + op.toString() + "\n" +
            // "Status: " + friendOpRet.getStatus() + "\n" +
            "Result: " + grpOpRet.getResult() + "\n" +
            "Detail: " + grpOpRet.getDetail() + "\n" +
            "group info: " + grpOpRet.getGroup()+ "\n" +
            "req Mem: " + grpOpRet.getReqmem() + "\n" +
            "group members: " + grpOpRet.getMembersList() +"\n" +
            "group list: " + grpOpRet.getGroupsList()+ "\n";

        //const users = grpOpRet.getUsersList();
        // const user = users[0];
        // str += "user id: " + user.getUserid() + "\n" ;
        this.progressCallback(str);
    }

    onUploadResult(msg){
        const OpRet = msg.getPlainmsg().getUploadreply();

        let str = "";
        str += "Received upload  result  message:\n" +
            // "Status: " + friendOpRet.getStatus() + "\n" +
            "Result: " + OpRet.getResult() + "\n" +
            "Detail: " + OpRet.getDetail() + "\n" +
            "chunk index: " + OpRet.getChunkindex()+ "\n" +
            "name " + OpRet.getUuidname()+ "\n";

        if (OpRet.getResult() == "sameok") {
            this.progressCallback(OpRet.getFilename() + " 秒传完成 ");
        }else if (OpRet.getResult() == "fail"){
            this.progressCallback(OpRet.getFilename() + " 失败: " + OpRet.getDetail());
        }else if (OpRet.getResult() == "chunkok"){
            this.progressCallback(OpRet.getFilename() + " 分片完毕: " + OpRet.getChunkindex());
            let uploader = this.fileMap.get(OpRet.getFilename());
            if (uploader) {
                uploader.uploadOther(OpRet.getChunkindex()+1);
            }
        }else if (OpRet.getResult() == "fileok"){
            this.progressCallback(OpRet.getFilename() + "上传完毕:" + OpRet.getDetail());
        }


        //const users = grpOpRet.getUsersList();
        // const user = users[0];
        // str += "user id: " + user.getUserid() + "\n" ;
        this.progressCallback(str);
    }

    onChatMsg(msg){
        const chatdata = msg.getPlainmsg().getChatdata();

        let str = "";
        const data = chatdata.getData();
        const dataStr = bytesToString(data);
        str += "Received chat message:\n" +
            // "Status: " + friendOpRet.getStatus() + "\n" +
            "from: " + chatdata.getFromid() + "\n" +
            "to: " + chatdata.getToid()  + "\n" +
            "data: " + dataStr + "\n" +
            "id " + chatdata.getMsgid() + "\n";

        this.progressCallback(str);

        if (chatdata.getChattype() == proto.model.ChatType.CHATTYPEP2P ){
            this.progressCallback("应答回执");

            const uid = chatdata.getToid();
            const fid = chatdata.getFromid();
            this.sendP2pMessageReply(uid, fid, chatdata.getSendid(),chatdata.getMsgid());
        }


    }

    onChatMsgReply(msg){
        const reply = msg.getPlainmsg().getChatreply();

        let str = "";


        str += "Received chatreply message:\n" +
            // "Status: " + friendOpRet.getStatus() + "\n" +
            "tm: " + reply.getSendok() + "\n" +
            "tm1: " + reply.getRecvok() + "\n" +
            "tm2 " + reply.getReadok() + "\n" +
            "id "  + reply.getMsgid() + "\n";

        this.progressCallback(str);

    }

    onQueryResult(msg){
        const reply = msg.getPlainmsg().getCommonqueryret();
        switch (reply.getQuerytype()){
            case proto.model.QueryDataType.QUERYDATATYPECHATDATA:
                this.onQueryChatData(reply);
                break;
        }

    }

    // 查询用户聊天的数据的应答
    onQueryChatData(reply){
        if (reply.getChattype() == proto.model.ChatType.CHATTYPEP2P){
            this.onQueryP2pChatData(reply);
        }else{
            this.onQueryGroupChatData(reply);
        }
    }

    formatTimestamp(timestamp) {
        // 确保时间戳是数值类型
        if (typeof timestamp !== 'number') {
            return 'Invalid timestamp';
        }

        // 创建一个新的Date对象
        const date = new Date(timestamp);

        // 格式化日期和时间
        const year = date.getFullYear();
        const month = String(date.getMonth() + 1).padStart(2, '0'); // 月份从0开始，所以要加1
        const day = String(date.getDate()).padStart(2, '0');
        const hours = String(date.getHours()).padStart(2, '0');
        const minutes = String(date.getMinutes()).padStart(2, '0');
        const seconds = String(date.getSeconds()).padStart(2, '0');

        // 拼接成时间字符串
        const formattedTime = `${year}-${month}-${day} ${hours}:${minutes}:${seconds}`;

        // 返回时间字符串
        return formattedTime;
    }

    // 同步私聊数据的结果
    onQueryP2pChatData(reply){
        let str = "";
        str += "Received query chat data p2p message:\n";

        const lst = reply.getChatdatalistList();
        str += "count=" + lst.length +"\n-----------------------------\n";
        for (let i = 0; i < lst.length; i++) {
            const chatdata = lst[i];
            const data = chatdata.getData();
            const dataStr = bytesToString(data);
            const tm = Number(chatdata.getTm());
            str +=
                "from: " + chatdata.getFromid() + "\n" +
                "to: " + chatdata.getToid()  + "\n" +
                "data: " + dataStr + "\n" +
                "id " + chatdata.getMsgid() + "\n" +
                "tm" + this.formatTimestamp(tm) + "\n\n";
        }


        this.progressCallback(str);
    }

    // 同步群聊数据的结果
    onQueryGroupChatData(reply){
        let str = "";
        str += "Received query chat data p2p message:\n";

        const lst = reply.getChatdatalistList();
        str += "count=" + lst.length +"\n-----------------------------\n";
        for (let i = 0; i < lst.length; i++) {
            const chatdata = lst[i];
            const data = chatdata.getData();
            const dataStr = bytesToString(data);
            const tm = Number(chatdata.getTm());
            str +=
                "from: " + chatdata.getFromid() + "\n" +
                "to: " + chatdata.getToid()  + "\n" +
                "data: " + dataStr + "\n" +
                "id " + chatdata.getMsgid() + "\n" +
                "tm" + this.formatTimestamp(tm) + "\n\n";
        }


        this.progressCallback(str);
    }
    //////////////////////////////////////////////////////////////////////////////

    // 1.0 发送hello包
    async sendHello(){
        const timestamp = getCurrentTimestamp();
        // 子消息
        const hello = new proto.model.MsgHello();
        hello.setClientid("uuid");  //js权限低，这里可以使用一个UUID，执行本地存储，每次都带着，用于服务端区分设备
        hello.setVersion("1.0");
        hello.setPlatform("web");
        hello.setStage("clienthello");

        const keyPrint = loadSharedKeyPrint();
        if (keyPrint) {
            // 如果 sharedKeyPrint 存在，则执行相应的操作
            console.log("sharedKeyPrint 存在");
            console.log(keyPrint);
            this.shareKeyPrint = keyPrint;

            const shareKey = loadSharedKey();
            if (shareKey){
                this.shareKey = shareKey;

                console.log("时间戳=", timestamp);
                const checkData = await encryptString(timestamp, this.shareKey);
                console.log(checkData);
                hello.setKeyprint(this.shareKeyPrint);
                let paramsMap = hello.getParamsMap();
                // 设置键值对
                paramsMap.set('checkTokenData', checkData);
            }
        } else {
            // 如果 sharedKeyPrint 不存在，则执行其他操作
            console.log("sharedKeyPrint 不存在");
        }



        // 将 MsgHello 消息设置为 Msg 消息的 plainMsg 字段
        const plainMsg = new proto.model.MsgPlain();
        plainMsg.setHello(hello);


        // 封装为通用消息
        const msg = new proto.model.Msg();
        msg.setMsgtype(proto.model.ComMsgType.MSGTHELLO);
        msg.setVersion(1);
        msg.setPlainmsg(plainMsg);

        msg.setTm(timestamp);

        this.sendObject(msg);
    }

    // 秘钥交换阶段1: 生成公私钥对，发送公钥过去
    async sendExchange1(){
        this.keyPairEx = await createDHPair();
        console.log("local key pair:", this.keyPairEx );

        const publicKey = await getPEMPublicKey(this.keyPairEx);
        console.log("local public key:", publicKey);

        const publicKeyRawData = stringToBytes(publicKey);
        console.log("local public key data:", publicKeyRawData);

        const exMsg = new proto.model.MsgKeyExchange();
        exMsg.setKeyprint("0");    // 不使用临时秘钥，也没有旧秘钥
        exMsg.setRsaprint("0");    // 不使用RSA加密
        exMsg.setStage(1);       // 第一次握手
        //exMsg.setTempkey(null);     // 不使用RSA加密,也就没有临时秘钥
        exMsg.setPubkey(publicKeyRawData); // 公钥的明文，这里不加密
        exMsg.setEnctype("AES-CTR");    // status和detail不设置


        const plainMsg = new proto.model.MsgPlain();
        plainMsg.setKeyex(exMsg);

        // 封装为通用消息
        const msg = new proto.model.Msg();
        msg.setMsgtype(proto.model.ComMsgType.MSGTKEYEXCHANGE);
        msg.setVersion(1);
        msg.setPlainmsg(plainMsg);
        const timestamp = getCurrentTimestamp();
        msg.setTm(timestamp);

        this.sendObject(msg);
    }

    // 秘钥交换阶段3：把自己的计算结果告诉服务器，也将时间戳加密，告诉服务器，让服务器验证；
    async sendExchange3(){

        const tmStr = getCurrentTimestamp();
        console.log("时间戳=", tmStr);
        const tmData = stringToBytes(tmStr);
        const checkData = await encryptAES_CTR(tmData, this.shareKey);


        const exMsg = new proto.model.MsgKeyExchange();
        exMsg.setKeyprint(this.shareKeyPrint);    // 计算过的秘钥
        exMsg.setStage(3);       // 第3次握手
        exMsg.setTempkey(checkData);     // 不使用RSA加密,也就没有临时秘钥
        exMsg.setEnctype("AES-CTR");    //
        exMsg.setStatus("ready");


        const plainMsg = new proto.model.MsgPlain();
        plainMsg.setKeyex(exMsg);

        // 封装为通用消息
        const msg = new proto.model.Msg();
        msg.setMsgtype(proto.model.ComMsgType.MSGTKEYEXCHANGE);
        msg.setVersion(1);
        msg.setPlainmsg(plainMsg);
        msg.setTm(tmStr);

        this.sendObject(msg);
    }

    // 2
    sendHeart(){
        const heart = new proto.model.MsgHeartBeat();
        heart.setUserid(1);
        const timestamp = getCurrentTimestamp();
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
    // "anonymous" "email"
    sendRegisterMessage(name, pwd, email, type) {
        showMessage("注册申请的消息")
        const userInfo = new proto.model.UserInfo();
        userInfo.setUsername(name);
        userInfo.setUserid(0);
        userInfo.setEmail(email);
        const paramsMap = userInfo.getParamsMap();
        paramsMap.set("pwd", pwd);

        // 注册用户
        const regOpReq = new proto.model.UserOpReq();
        regOpReq.setOperation(proto.model.UserOperationType.REGISTERUSER);
        regOpReq.setUser(userInfo)
        const params1 = regOpReq.getParamsMap();
        params1.set("regmode", type);   // pwd, email,phone

        const plainMsg = new proto.model.MsgPlain();
        plainMsg.setUserop(regOpReq);

        // 封装为通用消息
        const tmStr = getCurrentTimestamp();
        const msg = new proto.model.Msg();
        msg.setMsgtype(proto.model.ComMsgType.MSGTUSEROP);
        msg.setVersion(1);
        msg.setPlainmsg(plainMsg);
        msg.setTm(tmStr);

        this.sendObject(msg);
    }

    // 4. 发送验证码，所有发送验证码的都是一样的，服务端跟踪当前
    sendCodeMessage(email) {
        showMessage("发送验证消息")
        const userInfo = new proto.model.UserInfo();
        userInfo.setUserid(0);
        userInfo.setEmail(email);

        // 验证码
        const regOpReq = new proto.model.UserOpReq();
        regOpReq.setOperation(proto.model.UserOperationType.REALNAMEVERIFICATION);
        regOpReq.setUser(userInfo);
        regOpReq.getParamsMap()

        const paramsMap =  regOpReq.getParamsMap();
        paramsMap.set("code", "12345");
        paramsMap.set("regmode", "email");   // pwd, email,phone

        const plainMsg = new proto.model.MsgPlain();
        plainMsg.setUserop(regOpReq);


        // 封装为通用消息
        const msg = new proto.model.Msg();
        msg.setMsgtype(proto.model.ComMsgType.MSGTUSEROP);
        msg.setVersion(1);
        msg.setPlainmsg(plainMsg);
        msg.setTm(getCurrentTimestamp());

        this.sendObject(msg);

    }

    // 5.1 用户密码登录申请
    sendLoginMessage(mode, id, pwd){
        showMessage("发送登录消息")
        const userInfo = new proto.model.UserInfo();
        if (mode == "phone"){
            userInfo.setPhone(id)
        }else if (mode == "email"){
            userInfo.setEmail(id);
        }else{
            userInfo.setUserid(id);
            const paramsMap = userInfo.getParamsMap();
            paramsMap.set("pwd", pwd);
        }

        // 登录
        const regOpReq = new proto.model.UserOpReq();
        regOpReq.setOperation(proto.model.UserOperationType.LOGIN);
        regOpReq.setUser(userInfo);
        const params1 = regOpReq.getParamsMap();
        params1.set("loginmode", mode);   // uidpwd, email,phone


        const plainMsg = new proto.model.MsgPlain();
        plainMsg.setUserop(regOpReq);


        // 封装为通用消息
        const msg = new proto.model.Msg();
        msg.setMsgtype(proto.model.ComMsgType.MSGTUSEROP);
        msg.setVersion(1);
        msg.setPlainmsg(plainMsg);
        msg.setTm(getCurrentTimestamp());

        this.sendObject(msg);
    }


    // 5.3 邮箱登录验证码
    sendLoginCodeMessage(mode, id, code){
        showMessage("发送验证码消息")
        const userInfo = new proto.model.UserInfo();
        if (mode == "phone"){
            userInfo.setPhone(id)
        }else if (mode == "email"){
            userInfo.setEmail(id);
        }else{
            return
        }

        // 登录
        const regOpReq = new proto.model.UserOpReq();
        regOpReq.setOperation(proto.model.UserOperationType.REALNAMEVERIFICATION);
        regOpReq.setUser(userInfo);
        const params1 = regOpReq.getParamsMap();
        params1.set("loginmode", mode);   // uidpwd, email,phone
        params1.set("code", code);


        const plainMsg = new proto.model.MsgPlain();
        plainMsg.setUserop(regOpReq);


        // 封装为通用消息
        const msg = new proto.model.Msg();
        msg.setMsgtype(proto.model.ComMsgType.MSGTUSEROP);
        msg.setVersion(1);
        msg.setPlainmsg(plainMsg);
        msg.setTm(getCurrentTimestamp());

        this.sendObject(msg);
    }

    // 6.1 设置用户的各种信息，如果重新设置邮箱以及手机号，则需要验证码
    sendUserInfoMessage(){
        showMessage("发送用户详细消息")
        const userInfo = new proto.model.UserInfo();
        userInfo.setUserid(10003);

        // 设置信息
        const regOpReq = new proto.model.UserOpReq();
        regOpReq.setOperation(proto.model.UserOperationType.SETUSERINFO);
        regOpReq.setUser(userInfo);

        const paramsMap =  regOpReq.getParamsMap();
        paramsMap.set("UserName", "Robin.fox");
        paramsMap.set("NickName", "飞鸟真人");
        paramsMap.set("Age", "35");
        paramsMap.set("Intro", "我是一个爱运动的博主>_<...");
        paramsMap.set("Gender", "男");
        paramsMap.set("Region", "北京");
        paramsMap.set("Icon", "飞鸟真人");
        paramsMap.set("Params.title", "经理")
        ///paramsMap.set("Email", "robin-fox@sohu.com");


        const plainMsg = new proto.model.MsgPlain();
        plainMsg.setUserop(regOpReq);


        // 封装为通用消息
        const msg = new proto.model.Msg();
        msg.setMsgtype(proto.model.ComMsgType.MSGTUSEROP);
        msg.setVersion(1);
        msg.setPlainmsg(plainMsg);
        msg.setTm(getCurrentTimestamp());

        this.sendObject(msg);
    }

    sendUserInfoMessage1(){
        showMessage("发送用户详细消息")
        // 设置信息
        const regOpReq = new proto.model.UserOpReq();
        regOpReq.setOperation(proto.model.UserOperationType.SETUSERINFO);
        // regOpReq.setUser(userInfo);

        const paramsMap =  regOpReq.getParamsMap();
        paramsMap.set("Email", "robin-fox@sohu.com");


        const plainMsg = new proto.model.MsgPlain();
        plainMsg.setUserop(regOpReq);


        // 封装为通用消息
        const msg = new proto.model.Msg();
        msg.setMsgtype(proto.model.ComMsgType.MSGTUSEROP);
        msg.setVersion(1);
        msg.setPlainmsg(plainMsg);
        msg.setTm(getCurrentTimestamp());

        this.sendObject(msg);
    }

    // 6.2
    sendUserInfoCodeMessage(code){
        showMessage("发送更改信息验证码")

        // 发送验证码
        const regOpReq = new proto.model.UserOpReq();
        regOpReq.setOperation(proto.model.UserOperationType.REALNAMEVERIFICATION);
        const params1 = regOpReq.getParamsMap();
        params1.set("code", code);


        const plainMsg = new proto.model.MsgPlain();
        plainMsg.setUserop(regOpReq);


        // 封装为通用消息
        const msg = new proto.model.Msg();
        msg.setMsgtype(proto.model.ComMsgType.MSGTUSEROP);
        msg.setVersion(1);
        msg.setPlainmsg(plainMsg);
        msg.setTm(getCurrentTimestamp());

        this.sendObject(msg);
    }

    // 7.1 禁用
    sendUserDisableMessage(){
        showMessage("禁用用户")
        const userInfo = new proto.model.UserInfo();
        userInfo.setUserid(10003);

        // 禁用
        const regOpReq = new proto.model.UserOpReq();
        regOpReq.setOperation(proto.model.UserOperationType.DISABLEUSER);
        regOpReq.setUser(userInfo);


        const plainMsg = new proto.model.MsgPlain();
        plainMsg.setUserop(regOpReq);

        // 封装为通用消息
        const msg = new proto.model.Msg();
        msg.setMsgtype(proto.model.ComMsgType.MSGTUSEROP);
        msg.setVersion(1);
        msg.setPlainmsg(plainMsg);

        this.sendObject(msg);
    }

    // 7.2 解禁
    sendUserEnableMessage(){
        showMessage("解禁用户")
        const userInfo = new proto.model.UserInfo();
        userInfo.setUserid(10003);

        // 解禁
        const regOpReq = new proto.model.UserOpReq();
        regOpReq.setOperation(proto.model.UserOperationType.RECOVERUSER);
        regOpReq.setUser(userInfo);


        // 封装为通用消息
        const plainMsg = new proto.model.MsgPlain();
        plainMsg.setUserop(regOpReq);

        // 封装为通用消息
        const msg = new proto.model.Msg();
        msg.setMsgtype(proto.model.ComMsgType.MSGTUSEROP);
        msg.setVersion(1);
        msg.setPlainmsg(plainMsg);

        this.sendObject(msg);
    }

    // 8  注销
    sendUserUnregMessage(){
        showMessage("注销用户")
        const userInfo = new proto.model.UserInfo();
        userInfo.setUserid(10003);

        //
        const regOpReq = new proto.model.UserOpReq();
        regOpReq.setOperation(proto.model.UserOperationType.UNREGISTERUSER);
        regOpReq.setUser(userInfo);


        const plainMsg = new proto.model.MsgPlain();
        plainMsg.setUserop(regOpReq);

        // 封装为通用消息
        const msg = new proto.model.Msg();
        msg.setMsgtype(proto.model.ComMsgType.MSGTUSEROP);
        msg.setVersion(1);
        msg.setPlainmsg(plainMsg);

        this.sendObject(msg);
    }

    // 9 退出登录
    sendUserLogoutMessage(){
        showMessage("登出用户")
        const userInfo = new proto.model.UserInfo();
        userInfo.setUserid(10003);

        //
        const regOpReq = new proto.model.UserOpReq();
        regOpReq.setOperation(proto.model.UserOperationType.LOGOUT);
        regOpReq.setUser(userInfo);


        const plainMsg = new proto.model.MsgPlain();
        plainMsg.setUserop(regOpReq);

        // 封装为通用消息
        const msg = new proto.model.Msg();
        msg.setMsgtype(proto.model.ComMsgType.MSGTUSEROP);
        msg.setVersion(1);
        msg.setPlainmsg(plainMsg);

        this.sendObject(msg);
    }

    // 10 查询好友，通过params设置查询的信息
    sendFriendFindMessage(mode, id){
        showMessage("查找好友")
        const userInfo = new proto.model.UserInfo();



        //
        const opReq = new proto.model.FriendOpReq();
        opReq.setOperation(proto.model.UserOperationType.FINDUSER);
        opReq.setUser(userInfo);
        const params1 = opReq.getParamsMap();
        params1.set("mode", mode);
        params1.set("value", id);


        const plainMsg = new proto.model.MsgPlain();
        plainMsg.setFriendop(opReq);

        // 封装为通用消息
        const msg = new proto.model.Msg();
        msg.setMsgtype(proto.model.ComMsgType.MSGTFRIENDOP);
        msg.setVersion(1);
        msg.setPlainmsg(plainMsg);

        this.sendObject(msg);
    }
    // 11 请求添加好友
    sendFriendAddMessage(id){
        showMessage("请求添加好友")
        const userInfo = new proto.model.UserInfo();
        userInfo.setUserid(id)

        //
        const opReq = new proto.model.FriendOpReq();
        opReq.setOperation(proto.model.UserOperationType.ADDFRIEND);
        opReq.setUser(userInfo);

        const plainMsg = new proto.model.MsgPlain();
        plainMsg.setFriendop(opReq);

        // 封装为通用消息
        const msg = new proto.model.Msg();
        msg.setMsgtype(proto.model.ComMsgType.MSGTFRIENDOP);
        msg.setVersion(1);
        msg.setPlainmsg(plainMsg);

        this.sendObject(msg);
    }

    // 同意好友

    //13 删除关注
    sendRemoveFollow(id){
        showMessage("删除好友")
        const userInfo = new proto.model.UserInfo();
        userInfo.setUserid(id)

        //
        const opReq = new proto.model.FriendOpReq();
        opReq.setOperation(proto.model.UserOperationType.REMOVEFRIEND);
        opReq.setUser(userInfo);

        const plainMsg = new proto.model.MsgPlain();
        plainMsg.setFriendop(opReq);

        // 封装为通用消息
        const msg = new proto.model.Msg();
        msg.setMsgtype(proto.model.ComMsgType.MSGTFRIENDOP);
        msg.setVersion(1);
        msg.setPlainmsg(plainMsg);
        this.sendObject(msg);
    }

    // 拉黑
    sendBlockFollow(id){
        showMessage("拉黑好友")
        const userInfo = new proto.model.UserInfo();
        userInfo.setUserid(id)

        //
        const opReq = new proto.model.FriendOpReq();
        opReq.setOperation(proto.model.UserOperationType.BLOCKFRIEND);
        opReq.setUser(userInfo);

        const plainMsg = new proto.model.MsgPlain();
        plainMsg.setFriendop(opReq);

        // 封装为通用消息
        const msg = new proto.model.Msg();
        msg.setMsgtype(proto.model.ComMsgType.MSGTFRIENDOP);
        msg.setVersion(1);
        msg.setPlainmsg(plainMsg);
        this.sendObject(msg);
    }

    // 解除拉黑
    sendUnBlockFollow(id){
        showMessage("解除拉黑好友")
        const userInfo = new proto.model.UserInfo();
        userInfo.setUserid(id)

        //
        const opReq = new proto.model.FriendOpReq();
        opReq.setOperation(proto.model.UserOperationType.UNBLOCKFRIEND);
        opReq.setUser(userInfo);

        const plainMsg = new proto.model.MsgPlain();
        plainMsg.setFriendop(opReq);

        // 封装为通用消息
        const msg = new proto.model.Msg();
        msg.setMsgtype(proto.model.ComMsgType.MSGTFRIENDOP);
        msg.setVersion(1);
        msg.setPlainmsg(plainMsg);
        this.sendObject(msg);
    }

    sendSetPermission(id, mask){
        showMessage("设置权限")
        const userInfo = new proto.model.UserInfo();
        userInfo.setUserid(id)

        //
        const opReq = new proto.model.FriendOpReq();
        opReq.setOperation(proto.model.UserOperationType.SETFRIENDPERMISSION);
        opReq.setUser(userInfo);
        const params = opReq.getParamsMap();
        params["permission"]= mask

        const plainMsg = new proto.model.MsgPlain();
        plainMsg.setFriendop(opReq);

        // 封装为通用消息
        const msg = new proto.model.Msg();
        msg.setMsgtype(proto.model.ComMsgType.MSGTFRIENDOP);
        msg.setVersion(1);
        msg.setPlainmsg(plainMsg);
        this.sendObject(msg);
    }

    sendSetFriendMemo(id, mode, nick){
        showMessage("设置备注")
        const userInfo = new proto.model.UserInfo();
        userInfo.setUserid(id);
        userInfo.setNickname(nick);

        //
        const opReq = new proto.model.FriendOpReq();
        opReq.setOperation(proto.model.UserOperationType.SETFRIENDMEMO);
        opReq.setUser(userInfo);
        const params = opReq.getParamsMap();
        params.set("mode", mode)

        const plainMsg = new proto.model.MsgPlain();
        plainMsg.setFriendop(opReq);

        // 封装为通用消息
        const msg = new proto.model.Msg();
        msg.setMsgtype(proto.model.ComMsgType.MSGTFRIENDOP);
        msg.setVersion(1);
        msg.setPlainmsg(plainMsg);
        this.sendObject(msg);
    }

    sendListFriend(mode){
        showMessage("列出好友")
        //
        // const userInfo = new proto.model.UserInfo();
        // userInfo.setUserid(10003);
        //

        const opReq = new proto.model.FriendOpReq();
        opReq.setOperation(proto.model.UserOperationType.LISTFRIENDS);
        const params1 = opReq.getParamsMap();
        params1.set("mode", mode);
        //opReq.setUser(userInfo);


        const plainMsg = new proto.model.MsgPlain();
        plainMsg.setFriendop(opReq);


        // 封装为通用消息
        const msg = new proto.model.Msg();
        msg.setMsgtype(proto.model.ComMsgType.MSGTFRIENDOP);
        msg.setVersion(1);
        msg.setPlainmsg(plainMsg);
        this.sendObject(msg);
    }

    ////////////////////////////////////////////////////////
    sendCrateGroupMessage(name, visibility){
        showMessage("创建群");

        const group = new proto.model.GroupInfo();
        group.setGroupname(name);
        const params = group.getParamsMap();
        params.set("visibility", visibility)

        const opReq = new proto.model.GroupOpReq();
        opReq.setOperation(proto.model.GroupOperationType.GROUPCREATE);
        opReq.setGroup(group);

        const plainMsg = new proto.model.MsgPlain();
        plainMsg.setGroupop(opReq);


        // 封装为通用消息
        const msg = new proto.model.Msg();
        msg.setMsgtype(proto.model.ComMsgType.MSGTGROUPOP);
        msg.setVersion(1);
        msg.setPlainmsg(plainMsg);
        this.sendObject(msg);
    }

    // 查找群
    sendFindGroupMessage(keyword){
        showMessage("搜索群");
        const opReq = new proto.model.GroupOpReq();
        opReq.setOperation(proto.model.GroupOperationType.GROUPSEARCH);
        const params = opReq.getParamsMap();
        params.set("keyword", keyword);

        const plainMsg = new proto.model.MsgPlain();
        plainMsg.setGroupop(opReq);


        // 封装为通用消息
        const msg = new proto.model.Msg();
        msg.setMsgtype(proto.model.ComMsgType.MSGTGROUPOP);
        msg.setVersion(1);
        msg.setPlainmsg(plainMsg);
        this.sendObject(msg);
    }

    // 设置群信息
    sendSetgroupMemo(id, tags, brief, icon){
        showMessage("更新群");
        const group = new proto.model.GroupInfo();
        group.setGroupid(id);
        group.setGroupname("飞鸟和鱼鱼");
        group.setTagsList(tags);
        const params = group.getParamsMap();
        params.set("icon", icon);
        params.set("brief", brief);
        params.set("jointype", "any")
        // params.set("jointype", "admin")
        // params.set("jointype", "invite")
        // params.set("jointype", "question")

        //params["jointype"] = "any" | "invite" | "admin" | "question"



        const opReq = new proto.model.GroupOpReq();
        opReq.setOperation(proto.model.GroupOperationType.GROUPSETINFO);
        opReq.setGroup(group);

        const plainMsg = new proto.model.MsgPlain();
        plainMsg.setGroupop(opReq);


        // 封装为通用消息
        const msg = new proto.model.Msg();
        msg.setMsgtype(proto.model.ComMsgType.MSGTGROUPOP);
        msg.setVersion(1);
        msg.setPlainmsg(plainMsg);
        this.sendObject(msg);
    }
    // 请求加入群
    sendJoinGroupReq(id) {
        showMessage("尝试加入群");
        const group = new proto.model.GroupInfo();
        group.setGroupid(id);

        const opReq = new proto.model.GroupOpReq();
        opReq.setOperation(proto.model.GroupOperationType.GROUPJOINREQUEST);
        opReq.setGroup(group);

        const plainMsg = new proto.model.MsgPlain();
        plainMsg.setGroupop(opReq);


        // 封装为通用消息
        const msg = new proto.model.Msg();
        msg.setMsgtype(proto.model.ComMsgType.MSGTGROUPOP);
        msg.setVersion(1);
        msg.setPlainmsg(plainMsg);
        this.sendObject(msg);

    }

    // 从某一个UID开始查询群列表，最多返回100个，下次再分页
    sendListGroupMemberMessage(id){
        showMessage("查询群成员列表群");
        const group = new proto.model.GroupInfo();
        group.setGroupid(id);

        const opReq = new proto.model.GroupOpReq();
        opReq.setOperation(proto.model.GroupOperationType.GROUPSEARCHMEMBER);
        opReq.setGroup(group);
        const params = opReq.getParamsMap();
        params.set("uid", "100");

        const plainMsg = new proto.model.MsgPlain();
        plainMsg.setGroupop(opReq);


        // 封装为通用消息
        const msg = new proto.model.Msg();
        msg.setMsgtype(proto.model.ComMsgType.MSGTGROUPOP);
        msg.setVersion(1);
        msg.setPlainmsg(plainMsg);
        this.sendObject(msg);
    }

    // 添加管理员
    sendAddGroupAdminMessage(id){
        showMessage("添加管理员");
        const group = new proto.model.GroupInfo();
        group.setGroupid(id);

        const opReq = new proto.model.GroupOpReq();
        opReq.setOperation(proto.model.GroupOperationType.GROUPADDADMIN);
        opReq.setGroup(group);

        const mem = new proto.model.GroupMember();
        mem.setUserid(10005);
        const memList = opReq.getMembersList();
        memList.push(mem);



        const plainMsg = new proto.model.MsgPlain();
        plainMsg.setGroupop(opReq);


        // 封装为通用消息
        const msg = new proto.model.Msg();
        msg.setMsgtype(proto.model.ComMsgType.MSGTGROUPOP);
        msg.setVersion(1);
        msg.setPlainmsg(plainMsg);
        this.sendObject(msg);

    }
    // 移除管理员
    sendRemoveGroupAdminMessage(id){
        showMessage("添加管理员");
        const group = new proto.model.GroupInfo();
        group.setGroupid(id);

        const opReq = new proto.model.GroupOpReq();
        opReq.setOperation(proto.model.GroupOperationType.GROUPDELADMIN);
        opReq.setGroup(group);

        const mem = new proto.model.GroupMember();
        mem.setUserid(10005);
        const memList = opReq.getMembersList();
        memList.push(mem);



        const plainMsg = new proto.model.MsgPlain();
        plainMsg.setGroupop(opReq);


        // 封装为通用消息
        const msg = new proto.model.Msg();
        msg.setMsgtype(proto.model.ComMsgType.MSGTGROUPOP);
        msg.setVersion(1);
        msg.setPlainmsg(plainMsg);
        this.sendObject(msg);
    }
    // 转让群主
    sendTranferGroupOwnerMessage(){

    }
    // 设置自己的昵称
    sendSetGroupMemberNickMessage(id, nick){
        showMessage("设置自己的群昵称");
        const group = new proto.model.GroupInfo();
        group.setGroupid(id);

        const opReq = new proto.model.GroupOpReq();
        opReq.setOperation(proto.model.GroupOperationType.GROUPSETMEMBERINFO);
        opReq.setGroup(group);
        const params = opReq.getParamsMap();
        params.set("nick", nick);

        const plainMsg = new proto.model.MsgPlain();
        plainMsg.setGroupop(opReq);


        // 封装为通用消息
        const msg = new proto.model.Msg();
        msg.setMsgtype(proto.model.ComMsgType.MSGTGROUPOP);
        msg.setVersion(1);
        msg.setPlainmsg(plainMsg);
        this.sendObject(msg);
    }

    sendQuitGroupMessage(id){
        showMessage("退群");
        const group = new proto.model.GroupInfo();
        group.setGroupid(id);

        const opReq = new proto.model.GroupOpReq();
        opReq.setOperation(proto.model.GroupOperationType.GROUPQUIT);
        opReq.setGroup(group);


        const plainMsg = new proto.model.MsgPlain();
        plainMsg.setGroupop(opReq);


        // 封装为通用消息
        const msg = new proto.model.Msg();
        msg.setMsgtype(proto.model.ComMsgType.MSGTGROUPOP);
        msg.setVersion(1);
        msg.setPlainmsg(plainMsg);
        this.sendObject(msg);
    }

    sendUingMessage(){
        showMessage("查询自己加入了多少群");

        const opReq = new proto.model.GroupOpReq();
        opReq.setOperation(proto.model.GroupOperationType.GROUPLISTIN);

        const plainMsg = new proto.model.MsgPlain();
        plainMsg.setGroupop(opReq);

        // 封装为通用消息
        const msg = new proto.model.Msg();
        msg.setMsgtype(proto.model.ComMsgType.MSGTGROUPOP);
        msg.setVersion(1);
        msg.setPlainmsg(plainMsg);
        this.sendObject(msg);
    }

    ///////////////////////////////////////////
    sendUploadMessage( fileName, fileType, fileSize, chunkSize, chunkCount, chunkIndex,  data,  hashCode){

        const uploadMsg = new proto.model.MsgUploadReq();
        uploadMsg.setGroupid(0);
        uploadMsg.setFilename(fileName);
        uploadMsg.setFilesize(fileSize);
        uploadMsg.setFiletype(fileType);
        uploadMsg.setChunkindex(chunkIndex);
        uploadMsg.setChunkcount(chunkCount);
        uploadMsg.setChunksize(chunkSize);
        uploadMsg.setFiledata(data);
        uploadMsg.setHashtype("md5");
        uploadMsg.setHashcode(hashCode);
        uploadMsg.setSendid(1);

        const plainMsg = new proto.model.MsgPlain();
        plainMsg.setUploadreq(uploadMsg);

        // 封装为通用消息
        const msg = new proto.model.Msg();
        msg.setMsgtype(proto.model.ComMsgType.MSGTUPLOAD);
        msg.setVersion(1);
        msg.setPlainmsg(plainMsg);
        this.sendObject(msg);
    }

    // 开始传递文件，读取第一块并发送
    async startUploadFile(file, tz){
        const uploader = new FileUploader(this, file, tz)
        this.fileMap.set(file.name, uploader);
        uploader.uploadFirstTrunk();
    }

    // 发送私聊的消息
    sendP2pMessage(uid, fid, txt){
        const txtMsg = new proto.model.MsgChat();
        txtMsg.setUserid(uid);
        txtMsg.setFromid(uid);
        txtMsg.setToid(fid);
        txtMsg.setTm(getCurrentTimestamp());
        txtMsg.setDevid("");
        txtMsg.setSendid(1);
        txtMsg.setMsgtype(proto.model.ChatMsgType.TEXT);

        const data = stringToBytes(txt);
        txtMsg.setData(data);
        txtMsg.setRefmessageid(0);
        txtMsg.setChattype(proto.model.ChatType.CHATTYPEP2P);

        const plainMsg = new proto.model.MsgPlain();
        plainMsg.setChatdata(txtMsg);

        // 封装为通用消息
        const msg = new proto.model.Msg();
        msg.setMsgtype(proto.model.ComMsgType.MSGTCHATMSG);
        msg.setVersion(1);
        msg.setPlainmsg(plainMsg);
        this.sendObject(msg);
    }

    sendP2pMessageReply(uid, fid, sendId, msgId){
        const reply = new proto.model.MsgChatReply();
        reply.setUserid(fid);
        reply.setFromid(uid);
        reply.setSendid(sendId);
        reply.setMsgid(msgId);
        reply.setRecvok(getCurrentTimestamp());

        const plainMsg = new proto.model.MsgPlain();
        plainMsg.setChatreply(reply);

        // 封装为通用消息
        const msg = new proto.model.Msg();
        msg.setMsgtype(proto.model.ComMsgType.MSGTCHATREPLY);
        msg.setVersion(1);
        msg.setPlainmsg(plainMsg);
        this.sendObject(msg);
    }

    sendGroupMessage(uid, gid, txt){
        const txtMsg = new proto.model.MsgChat();
        txtMsg.setUserid(uid);
        txtMsg.setFromid(uid);
        txtMsg.setToid(gid);
        txtMsg.setTm(getCurrentTimestamp());
        txtMsg.setDevid("");
        txtMsg.setSendid(1);
        txtMsg.setMsgtype(proto.model.ChatMsgType.TEXT);

        const data = stringToBytes(txt);
        txtMsg.setData(data);
        txtMsg.setRefmessageid(0);
        txtMsg.setChattype(proto.model.ChatType.CHATTYPEGROUP);

        const plainMsg = new proto.model.MsgPlain();
        plainMsg.setChatdata(txtMsg);

        // 封装为通用消息
        const msg = new proto.model.Msg();
        msg.setMsgtype(proto.model.ComMsgType.MSGTCHATMSG);
        msg.setVersion(1);
        msg.setPlainmsg(plainMsg);
        this.sendObject(msg);
    }

    // 同步当前的
    sendSynP2pChatHistory(){
        const query = new proto.model.MsgQuery();
        query.setUserid(10004);
        query.setGroupid(0);
        query.setLittleid(0);
        query.setSyntype(proto.model.SynType.SYNTYPEFORWARD);
        query.setChattype(proto.model.ChatType.CHATTYPEP2P);
        query.setQuerytype(proto.model.QueryDataType.QUERYDATATYPECHATDATA);


        const plainMsg = new proto.model.MsgPlain();
        plainMsg.setCommonquery(query);

        // 封装为通用消息
        const msg = new proto.model.Msg();
        msg.setMsgtype(proto.model.ComMsgType.MSGTQUERY);
        msg.setVersion(1);
        msg.setPlainmsg(plainMsg);
        this.sendObject(msg);
    }

    // 同步群聊历史数据
    sendSynGroupChatHistory(gid){
        const query = new proto.model.MsgQuery();
        query.setUserid(10004);
        query.setGroupid(gid);
        query.setLittleid(0);
        query.setSyntype(proto.model.SynType.SYNTYPEBACKWARD);
        query.setChattype(proto.model.ChatType.CHATTYPEGROUP);
        query.setQuerytype(proto.model.QueryDataType.QUERYDATATYPECHATDATA);


        const plainMsg = new proto.model.MsgPlain();
        plainMsg.setCommonquery(query);

        // 封装为通用消息
        const msg = new proto.model.Msg();
        msg.setMsgtype(proto.model.ComMsgType.MSGTQUERY);
        msg.setVersion(1);
        msg.setPlainmsg(plainMsg);
        this.sendObject(msg);
    }


}


