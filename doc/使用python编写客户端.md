# 使用python编写客户端

# 一、websocket

# 二、hello

## 2.1 需要的基本功能

1. 操作系统
2. 协议版本号
3. SDK版本
4. 硬件的基础信息字符串
5. 时区，地区

## 2.2 基本信息



## 2.3 直接验证秘钥

如果有keyprint和共享密钥，那么应该直接使用秘钥验证身份，这样速度更快；

1）生成毫秒时间戳，并转为字符串；

2）使用AES_CTR方式加密时间戳，作为验证数据；验证数据BASE64编码放在了param["checkTokenData"]中；

3）在hello包中设置keyprint;

# 三、秘钥交换

## 3.1 需要的基本功能

1. 生成秘钥对
2. 导出公钥
3. 生成共享密钥
4. 生成秘钥签名
5. 使用AEC_CTR加解密
6. 保存和加载密钥与签名

## 3.2 阶段1

1）生成密钥对，放在内存；

2）将公钥按照PEM方式编码并转为byte[]写到publickey字段；

3）如果使用RSA公钥，则在RSA字段指定公钥的指纹，并加密用于交换的公钥；（目前不支持）

4）设置为阶段1；

```python
 # 生成阶段1的消息
    def create_keyex1(self) -> msg_pb2.Msg:
        self.keyEx.generate_key_pair()
        pub_key = self.keyEx.get_public_key()  

        exMsg = msg_pb2.MsgKeyExchange()
        exMsg.keyPrint = 0
        exMsg.rsaPrint = 0
        exMsg.stage = 1
        exMsg.pubKey = pub_key    # PEM格式编码的字节流UTF-8
        exMsg.encType = "AES-CTR" # 加密算法

        # Create an instance of MsgPlain and set its keyex field
        plainMsg = msg_pb2.MsgPlain()
        plainMsg.keyEx.CopyFrom(exMsg)

        # Create an instance of Msg and set its fields
        msg = msg_pb2.Msg()
        msg.msgType = msg_pb2.ComMsgType.MsgTKeyExchange
        msg.version = 1
        msg.plainMsg.CopyFrom(plainMsg)
        msg.tm = int(time.time())  # returns a Unix timestamp in seconds
        return msg
```



## 3.3 阶段3

收到服务器发来的阶段2的信息；里面包含了对方公钥，指纹，验证数据；

1）这里收到的公钥是PEM字符串的字节流形式；生成共享密钥；生成指纹；

2）验证指纹，tempkey字段是tm时间戳的字符串加密后的验证数据；这里是byte[]不编码；

3）验证出错，说明不兼容，无法继续；

4）验证通过，应该向服务器发送指纹和验证数据，服务器也验证一遍交换结果；

```python
# 生成阶段3的消息, 已经秘钥交换完成；
    def create_keyex3(self) -> msg_pb2.Msg:

        tm = int(time.time())
        tmStr = str(tm)
        exMsg = msg_pb2.MsgKeyExchange()
        exMsg.keyPrint = self.keyEx.get_key_print()
        exMsg.tempKey = self.keyEx.encrypt_aes_ctr_bytes_to_bytes(tmStr)
        exMsg.rsaPrint = 0
        exMsg.stage = 3
        exMsg.encType = "AES-CTR" # 加密算法
        exMsg.status = "ready"

        # Create an instance of MsgPlain and set its keyex field
        plainMsg = msg_pb2.MsgPlain()
        plainMsg.keyEx.CopyFrom(exMsg)

        # Create an instance of Msg and set its fields
        msg = msg_pb2.Msg()
        msg.msgType = msg_pb2.ComMsgType.MsgTKeyExchange
        msg.version = 1
        msg.plainMsg.CopyFrom(plainMsg)
        msg.tm = tm  # returns a Unix timestamp in seconds
        return msg
```

交换结束后，有可能需要登录，也有可能需要注册，或者直接登录成功；

```python
# 这里会是阶段2，或者阶段4的应答
    async def on_key_exchange(self, msg: msg_pb2.Msg):
        keyex = msg.plainMsg.keyEx
        if keyex.stage == 2:
            key_print = keyex.keyPrint
            pub_key = keyex.pubKey
            print(f"remote public key {pub_key}")
            self.keyEx.exchange_keys(pub_key)


            self.keyEx.save_key_print(self.printName)
            self.keyEx.save_shared_key(self.keyName)
            local_print = self.keyEx.get_int64_print()
            print(f"local print{local_print}")
            print(f"local share key{self.keyEx.get_shared_key()}")
            if local_print != key_print:
                print(f"local keyprint:{local_print} is not same with remote key:{key_print}")
                return
            
            data = self.create_keyex3()
            await self.send(data)
        
        elif keyex.stage == 4:  # 交换秘钥之后也需要登录，或者注册
            if keyex.status == "waitdata":
                self.client_state = ClientState.READY
            if keyex.status == "needlogin":
                self.client_state = ClientState.WAIT_LOGIN
```



# 四、登录与同步消息

这个库并没有提供注册的功能，因为主要用户就是机器人，机器人也需要使用客户端手动注册；





# 五、用户与好友操作



# 六、群组操作





