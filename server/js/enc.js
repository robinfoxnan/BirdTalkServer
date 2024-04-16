

// 使用示例
var remotePublicKey = [/* 远程公钥的字节数组 */];
var curve = "P-256";
var privateKey = [/* 本地私钥的字节数组 */];
var publicKey = [];
var encType = "chacha20"; // 或者 "aes256" 或 "twofish128"


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

    async function getPublicKey(keyPair) {
        const publicKey = await crypto.subtle.exportKey("spki", keyPair.publicKey);
        return publicKey;
    }

    async function getPublicKeyRaw(keyPair) {
        const publicKey = await crypto.subtle.exportKey("raw", keyPair.publicKey);
        return publicKey;
    }


    // async function getPrivateKey(keyPair) {
    //     const publicKey = await crypto.subtle.exportKey("pkcs8", keyPair.publicKey);
    //     return publicKey;
    // }


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

    async function importPublicKey(rawPublicKey) {
        // 将原始公钥字节数组导入为 CryptoKey 对象
        const publicKey = await crypto.subtle.importKey(
            "raw", // 原始数据格式
            rawPublicKey, // 原始公钥字节数组
            {
                name: "ECDH", // 算法名称
                namedCurve: "P-256" // 曲线名称
            },
            true, // 是否为公钥
            [] // 公钥的用途
        );
        return publicKey;
    }

    async function importPublicKeySPKI(PublicKey) {
        // 将原始公钥字节数组导入为 CryptoKey 对象
        const publicKey = await crypto.subtle.importKey(
            "spki", // 使用 SPKI 格式
            PublicKey, // 原始公钥字节数组
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

        return Number(int64Value);
    }

    async function encryptString(str, key){
        // 将文本编码为 ArrayBuffer
        const dataBuffer = new TextEncoder().encode(str);
        return encryptBytes(dataBuffer, key)
    }

    async function encryptBytes(data, key) {
        try {
            // 生成随机的初始化向量（IV）
            const iv = crypto.getRandomValues(new Uint8Array(12));

            // 使用 CryptoKey 对象创建 ChaCha20 算法对象
            const algorithm = { name: "AES-GCM", iv: iv };

            // 使用 CryptoKey 对象加密数据
            const encryptedData = await crypto.subtle.encrypt(algorithm, key, dataBuffer);

            // 将 IV 和加密后的数据拼接在一起
            const result = new Uint8Array(iv.length + encryptedData.byteLength);
            result.set(iv, 0);
            result.set(new Uint8Array(encryptedData), iv.length);

            // 返回拼接后的数据
            return result;
        } catch (error) {
            console.error("Encryption failed:", error);
            throw error;
        }
    }

    // 加密整数
    async function encryptInt64(int64Value, key) {
        // 将 int64 值转换为字节数组
        const dataBuffer = new ArrayBuffer(8); // 64位整数字节长度为8
        const view = new DataView(dataBuffer);
        view.setBigInt64(0, BigInt(int64Value));
        return encryptBytes(dataBuffer, key);
    }

    // 解密整数
    async function decryptInt64(encryptedData, key) {
        try {
            // 使用 CryptoKey 对象创建 AES-GCM 算法对象
            const algorithm = { name: "AES-GCM", iv: encryptedData.slice(0, 12) };

            // 使用 CryptoKey 对象解密数据
            const decryptedData = await crypto.subtle.decrypt(algorithm, key, encryptedData.slice(12));

            // 将解密后的字节数组转换回 int64 值
            const view = new DataView(decryptedData);
            const decryptedInt64Value = view.getBigInt64(0);

            // 返回解密后的 int64 值
            return decryptedInt64Value;
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
    }


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
}






