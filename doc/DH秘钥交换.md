



# 1. 秘钥交换协议

椭圆曲线加密（Elliptic Curve Cryptography，ECC）和迪菲-赫尔曼密钥交换（Diffie-Hellman Key Exchange，DH）都是现代密码学中常用的技术，它们可以用于安全地交换密钥和加密通信。

1. **Diffie-Hellman密钥交换**（DH）是一种密钥协商协议，允许双方在没有事先共享密钥的情况下协商出一个共享密钥。DH协议的安全性基于数论中的离散对数问题。
2. **椭圆曲线加密**（ECC）则是一种基于椭圆曲线数学结构的加密算法。ECC提供了与RSA等传统加密算法相似的安全性，但使用更短的密钥长度。这意味着在相同的安全级别下，ECC需要更少的计算资源，更适合于资源受限的环境，比如移动设备和物联网设备。

虽然它们是两种不同的技术，但它们之间存在联系：

- **ECDH（椭圆曲线迪菲-赫尔曼密钥交换）**是将椭圆曲线加密与Diffie-Hellman密钥交换相结合的协议。它使用椭圆曲线上的点来执行密钥交换，从而提供了一种更高效的方法来协商共享密钥。ECDH协议的安全性基于椭圆曲线离散对数问题，与传统的DH协议相比，它提供了更高的安全性和更短的密钥长度。
- 椭圆曲线加密还可以用于数字签名和身份验证，而DH主要用于密钥交换。但是，两者都基于数论和离散数学的原理，具有相似的数学基础。

综上所述，椭圆曲线加密和Diffie-Hellman密钥交换是现代密码学中常用的技术，它们可以单独使用，也可以结合在一起提供更高级别的安全性。

# 2.使用GO和椭圆曲线实现



在telegram中使用的是DL秘钥交换，毕竟这个软件诞生时间比较早了。

协议实现：https://blog.csdn.net/robinfoxnan/article/details/127322483

国人有人做了telegram的开源服务端。

DL秘钥交换如果自己动手实现时候有几点需要注意：

1）g通常取小素数，2，3，5等；

2）p应该取一个大的素数，推荐2048位；如果随机数不能保证是素数，计算结果可能不正确；

在 Go 语言中，`crypto/dh` 包是 Go 标准库的一部分。

需要注意的是，从 Go 1.17 版本开始，`crypto/dh` 包已经被标记为废弃（deprecated），并且不建议在新的代码中使用它。

一个替代方案是使用 `crypto/elliptic` 包来手动实现 Diffie-Hellman 算法。

```go
// generateDHKeyPair 生成椭圆曲线的密钥对
func generateDHKeyPair(curve elliptic.Curve) (*big.Int, []byte, error) {
	// 生成私钥
	privateKey, err := rand.Int(rand.Reader, curve.Params().N)
	if err != nil {
		return nil, nil, err
	}

	// 计算公钥
	publicKeyX, publicKeyY := curve.ScalarBaseMult(privateKey.Bytes())
	publicKey := elliptic.MarshalCompressed(curve, publicKeyX, publicKeyY) // 使用压缩格式

	return privateKey, publicKey, nil
}

// sharedSecret 计算共享密钥
func sharedSecret(curve elliptic.Curve, publicKey []byte, privateKey *big.Int) ([]byte, error) {
	// 解码公钥
	x, y := elliptic.UnmarshalCompressed(curve, publicKey)
	if x == nil {
		return nil, fmt.Errorf("invalid public key")
	}

	// 计算共享密钥
	sharedKeyX, sharedKeyY := curve.ScalarMult(x, y, privateKey.Bytes())
	sharedKey := elliptic.Marshal(curve, sharedKeyX, sharedKeyY) // 返回整个点的字节表示

	return sharedKey, nil
}
```

测试代码如下：

```go
// go test -run TestGenDHkey
func TestGenDHkey(t *testing.T) {

	// 选择椭圆曲线，这里使用P-256曲线作为示例
	curve := elliptic.P256()

	// Alice生成密钥对
	alicePrivateKey, alicePublicKey, err := generateDHKeyPair(curve)
	if err != nil {
		panic(err)
	}
	fmt.Printf("Alice's Public Key: %v\n", alicePublicKey)
	fmt.Printf("Alice's Private Key: %v\n", alicePrivateKey)

	// Bob生成密钥对
	bobPrivateKey, bobPublicKey, err := generateDHKeyPair(curve)
	if err != nil {
		panic(err)
	}
	fmt.Printf("Bob's Public Key: %v\n", bobPublicKey)
	fmt.Printf("Bob's Private Key: %v\n", bobPrivateKey)

	// Alice和Bob交换公钥，并计算共享密钥
	aliceSharedSecret, err := sharedSecret(curve, bobPublicKey, alicePrivateKey)
	if err != nil {
		panic(err)
	}
	bobSharedSecret, err := sharedSecret(curve, alicePublicKey, bobPrivateKey)
	if err != nil {
		panic(err)
	}

	// 验证共享密钥是否相同
	if !CompareBytes(aliceSharedSecret, bobSharedSecret) {
		panic("Shared secrets do not match!")
	}

	//fmt.Printf("Shared Secret: %s\n", aliceSharedSecret.String())
}
```

可以看出使用椭圆曲线交互更加方便，因为算法内置了相关的参数，不需要传递或者协商P和Q参数；双方只需要交换自己的公钥就可以协商出合适的共享密钥；



**备注**：

生成的共享密钥仅仅是一个字节数组，那么在作为共享密钥时候，可以使用MD5和SHA256做一个哈希，那么就能得到128或者256位的秘钥字节数组。

# 3.用RSA防止中间人攻击

为了防止中间人攻击，需要在客户端内置服务端的RSA公钥，在协商的时候，只有服务端才能解开传递的公钥部分；

这样攻击在HTTPS中常见，为了截获数据，在中间充当一个双向代理，分别与客户端与服务器建立链接并协商秘钥，那么这个中间代理就能看到明文；

telegram 方案就是在客户端内置了服务端的几个RSA公钥，这样可以防止中间人监听明文；在协商对称密钥的时候，使用RSA加密。



RSA实现

```go
// len = 2048
// SA 密钥对是由公钥和私钥组成的，其中私钥包含了公钥的信息。因此，从私钥可以推算出公钥。
func GenRSAKey(privateKeyFile, publicKeyFile string) (*rsa.PrivateKey, error) {
	privateKey, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	err = SaveRSAPrivateKeyToFile(privateKey, &privateKey.PublicKey, privateKeyFile, publicKeyFile)
	return privateKey, err

}

func SaveRSAPrivateKeyToFile(privateKey *rsa.PrivateKey, publicKey *rsa.PublicKey,
	privateKeyFileName, publicKeyFileName string) error {
	// 创建一个文件用于保存私钥
	// 保存私钥到文件
	privateKeyFile, err := os.Create(privateKeyFileName)
	if err != nil {
		return err
	}
	defer privateKeyFile.Close()

	privateKeyPEM := &pem.Block{
		Type:  "RSA PRIVATE KEY",
		Bytes: x509.MarshalPKCS1PrivateKey(privateKey),
	}

	if err := pem.Encode(privateKeyFile, privateKeyPEM); err != nil {
		return err
	}

	// 保存公钥到文件
	publicKeyFile, err := os.Create(publicKeyFileName)
	if err != nil {
		return err
	}
	defer publicKeyFile.Close()

	publicKeyBytes, err := x509.MarshalPKIXPublicKey(publicKey)
	if err != nil {
		return err
	}

	publicKeyPEM := &pem.Block{
		Type:  "PUBLIC KEY",
		Bytes: publicKeyBytes,
	}

	if err := pem.Encode(publicKeyFile, publicKeyPEM); err != nil {
		return err
	}

	return nil
}

// EncryptRSA 使用公钥加密消息
func EncryptRSA(publicKey *rsa.PublicKey, message []byte) ([]byte, error) {
	return rsa.EncryptPKCS1v15(rand.Reader, publicKey, message)
}

// DecryptRSA 使用私钥解密消息
func DecryptRSA(privateKey *rsa.PrivateKey, encryptedMessage []byte) ([]byte, error) {
	return rsa.DecryptPKCS1v15(rand.Reader, privateKey, encryptedMessage)
}
```

`rsa.EncryptPKCS1v15` 是 RSA 加密中使用的一种填充方案，称为 PKCS#1 v1.5 填充。在这种填充方案中，要加密的消息被转换为长度等于密钥长度减去一些固定字节的数据块，然后再对这个数据块进行加密。这种填充方案是 RSA 加密的标准之一，广泛应用于加密通信和数字签名等领域。

测试

```go
func TestRsa(t *testing.T) {
	// 示例：生成 RSA 密钥对
	privateKey, err := GenRSAKey("chat_server_key.pem", "chat_server_public.pem")
	if err != nil {
		fmt.Println(err)
		return
	}

	msg := []byte("Hello, World!")
	// 使用公钥加密消息
	encryptedMessage, err := EncryptRSA(&privateKey.PublicKey, msg)
	if err != nil {
		fmt.Println(err)
		return
	}

	// 使用私钥解密消息
	decryptedMessage, err := DecryptRSA(privateKey, encryptedMessage)
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println("Decrypted message:", string(decryptedMessage)) // 输出：Decrypted message: Hello, World!
	if !CompareBytes(msg, decryptedMessage) {
		panic("Shared secrets do not match!")
	}
}
```

在使用 RSA 公钥时，可以遇到两种不同的 PEM 格式：`RSA PUBLIC KEY` 和 `PUBLIC KEY`。

1. **RSA PUBLIC KEY**: 这种格式通常包含了 RSA 公钥的 ASN.1 DER 编码，使用了 PKCS#1 标准。在这种情况下，您可以使用 `x509.ParsePKCS1PublicKey` 函数来解析这种类型的 PEM 数据块。
2. **PUBLIC KEY**: 这种格式通常包含了公钥的 X.509 DER 编码，可能包含了除 RSA 之外的其他公钥类型。在这种情况下，您可以使用 `x509.ParsePKIXPublicKey` 函数来解析这种类型的 PEM 数据块。

因此，区别在于这两种格式所包含的公钥编码的不同。

这里的公钥文件为：

```
-----BEGIN PUBLIC KEY-----
MIIBIjANBgkqhkiG9w0BAQEFAAOCAQ8AMIIBCgKCAQEA1DF287O19UDIg5ejz6aS
sLqiQ+8Ge3ExDrrsOi0TCU7qht/IUVDyv+UYPftCqSpSX0570rBTuwSam+6vUcBv
N/IAv4kg52b8n9IHu8eh+SxefEovEV3/LPiEDIAI8+qcolsorxwv3T0rATDxdMmu
gCmGJb6DAxOhmsuvTpOGrgovwaQ6svFCZ9LV1OsTjNVEVtm7utjKsYGtc6xunKA6
c5mzcqlOqdWhXXpPsDL46/224kDivv5lS1+wVgyACOadRQXtOyeE2ld4CcA0WzsC
HFzzsYeCALJv8D8EHQfm+8FP6Ur1ANByP/sHfBS4C2pdexk14obaqvOyDSCdvBqp
mwIDAQAB
-----END PUBLIC KEY-----

```



下面分别从文件加载公钥和私钥：

```go
// ExtractPublicKey 从 PEM 格式的公钥字符串中提取 RSA 公钥
func ExtractPublicKey(publicKeyPEM string) (*rsa.PublicKey, error) {
	// 去除 PEM 标记
	block, _ := pem.Decode([]byte(publicKeyPEM))
	if block == nil {
		return nil, errors.New("failed to decode PEM block containing public key")
	}

	var err error
	// 根据数据块类型解析RSA公钥
	var pubKey interface{}
	switch block.Type {
	case "RSA PUBLIC KEY":
		pubKey, err = x509.ParsePKCS1PublicKey(block.Bytes)
	case "PUBLIC KEY":
		pubKey, err = x509.ParsePKIXPublicKey(block.Bytes)
	default:
		return nil, fmt.Errorf("unsupported PEM block type: %s", block.Type)
	}
	if err != nil {
		return nil, err
	}

	// 转换为 *rsa.PublicKey 类型
	publicKey, ok := pubKey.(*rsa.PublicKey)
	if !ok {
		return nil, errors.New("failed to convert public key to RSA public key")
	}

	return publicKey, nil
}

// LoadRSAPublicKeyFromPEM 从 PEM 文件中加载 RSA 公钥
func LoadRSAPublicKeyFromPEM(filePath string) (*rsa.PublicKey, error) {
	// 如果是相对路径，则需要转换
	filePath, err := GetAbsolutePath(filePath)
	if err != nil {
		return nil, err
	}
	// 读取PEM格式的公钥文件
	pemBytes, err := os.ReadFile(filePath)
	if err != nil {
		return nil, err
	}

	// 解码PEM数据块
	block, _ := pem.Decode(pemBytes)
	if block == nil {
		return nil, fmt.Errorf("failed to decode PEM block")
	}

	// 根据数据块类型解析RSA公钥
	var pubKey interface{}
	switch block.Type {
	case "RSA PUBLIC KEY":
		pubKey, err = x509.ParsePKCS1PublicKey(block.Bytes)
	case "PUBLIC KEY":
		pubKey, err = x509.ParsePKIXPublicKey(block.Bytes)
	default:
		return nil, fmt.Errorf("unsupported PEM block type: %s", block.Type)
	}
	if err != nil {
		return nil, err
	}

	// 转换为RSA公钥类型
	rsaPubKey, ok := pubKey.(*rsa.PublicKey)
	if !ok {
		return nil, fmt.Errorf("failed to convert to RSA public key")
	}

	return rsaPubKey, nil
}

// LoadRSAPrivateKeyFromPEM 从 PEM 文件中加载 RSA 私钥
func LoadRSAPrivateKeyFromPEM(filePath string) (*rsa.PrivateKey, error) {
	filePath, err := GetAbsolutePath(filePath)
	if err != nil {
		return nil, err
	}
	// 读取 PEM 文件内容
	pemData, err := os.ReadFile(filePath)
	if err != nil {
		return nil, err
	}

	// 解码PEM数据块
	block, _ := pem.Decode(pemData)
	if block == nil {
		return nil, fmt.Errorf("failed to decode PEM block")
	}

	// 根据数据块类型解析RSA私钥
	var privKey interface{}
	switch block.Type {
	case "RSA PRIVATE KEY":
		privKey, err = x509.ParsePKCS1PrivateKey(block.Bytes)
	case "PRIVATE KEY":
		privKey, err = x509.ParsePKCS8PrivateKey(block.Bytes)
	default:
		return nil, fmt.Errorf("unsupported PEM block type: %s", block.Type)
	}
	if err != nil {
		return nil, err
	}

	// 转换为RSA私钥类型
	rsaPrivKey, ok := privKey.(*rsa.PrivateKey)
	if !ok {
		return nil, fmt.Errorf("failed to convert to RSA private key")
	}

	return rsaPrivKey, nil
}
```

测试函数：

```go
func TestLoadFileRSA(t *testing.T) {

	// 加载公钥
	publicKey, err := LoadRSAPublicKeyFromPEM(`D:\GBuild\BirdTalkServer\server\utils\chat_server_public.pem`)
	if err != nil {
		fmt.Println("Error loading public key:", err)
		return
	}
	//fmt.Println("Public key:", publicKey)

	// 加载私钥
	privateKey, err := LoadRSAPrivateKeyFromPEM(`D:\GBuild\BirdTalkServer\server\utils\chat_server_key.pem`)
	if err != nil {
		fmt.Println("Error loading private key:", err)
		return
	}
	//fmt.Println("Private key:", privateKey)

	msg := []byte("Hello, World!")
	// 使用公钥加密消息
	encryptedMessage, err := EncryptRSA(publicKey, msg)
	if err != nil {
		fmt.Println(err)
		return
	}

	// 使用私钥解密消息
	decryptedMessage, err := DecryptRSA(privateKey, encryptedMessage)
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println("Decrypted message:", string(decryptedMessage)) // 输出：Decrypted message: Hello, World!
	if !CompareBytes(msg, decryptedMessage) {
		panic("Shared secrets do not match!")
	}
}

func TestRSAString(t *testing.T) {
	// 公钥字符串
	publicKeyPEM := `-----BEGIN PUBLIC KEY-----
MIIBIjANBgkqhkiG9w0BAQEFAAOCAQ8AMIIBCgKCAQEA1DF287O19UDIg5ejz6aS
sLqiQ+8Ge3ExDrrsOi0TCU7qht/IUVDyv+UYPftCqSpSX0570rBTuwSam+6vUcBv
N/IAv4kg52b8n9IHu8eh+SxefEovEV3/LPiEDIAI8+qcolsorxwv3T0rATDxdMmu
gCmGJb6DAxOhmsuvTpOGrgovwaQ6svFCZ9LV1OsTjNVEVtm7utjKsYGtc6xunKA6
c5mzcqlOqdWhXXpPsDL46/224kDivv5lS1+wVgyACOadRQXtOyeE2ld4CcA0WzsC
HFzzsYeCALJv8D8EHQfm+8FP6Ur1ANByP/sHfBS4C2pdexk14obaqvOyDSCdvBqp
mwIDAQAB
-----END PUBLIC KEY-----`

	// 提取公钥
	publicKey, err := ExtractPublicKey(publicKeyPEM)
	if err != nil {
		fmt.Println("Error extracting public key:", err)
		return
	}

	fmt.Println("Extracted public key:", publicKey)

	// 加载私钥
	privateKey, err := LoadRSAPrivateKeyFromPEM(`D:\GBuild\BirdTalkServer\server\utils\chat_server_key.pem`)
	if err != nil {
		fmt.Println("Error loading private key:", err)
		return
	}
	//fmt.Println("Private key:", privateKey)

	msg := []byte("Hello, World!")
	// 使用公钥加密消息
	encryptedMessage, err := EncryptRSA(publicKey, msg)
	if err != nil {
		fmt.Println(err)
		return
	}

	// 使用私钥解密消息
	decryptedMessage, err := DecryptRSA(privateKey, encryptedMessage)
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println("Decrypted message:", string(decryptedMessage)) // 输出：Decrypted message: Hello, World!
	if !CompareBytes(msg, decryptedMessage) {
		panic("Shared secrets do not match!")
	}
}
```



# 4.对称加密加密算法

1. **AES（Advanced Encryption Standard）**：目前广泛使用的对称加密算法之一，支持多种密钥长度，包括128位、192位和256位。

2. **DES（Data Encryption Standard）**：是一种比较古老的对称加密算法，使用56位密钥，由于密钥较短，安全性受到了质疑。

3. **3DES（Triple Data Encryption Standard）**：是DES的一种变种，通过对数据使用三次DES加密来增强安全性。

4. **Blowfish**：是一种对称密钥加密算法，可用于加密各种长度的数据，其密钥长度可变，通常为32到448位。

5. **Twofish**：是一种高级的对称密钥分组加密算法，与Blowfish类似，但是Twofish支持更长的密钥长度，可以达到256位。

6. **RC4（Rivest Cipher 4）**：是一种流密码算法，用于将明文转换为密文，曾经广泛用于SSL/TLS协议，但由于存在一些安全问题，已经不再推荐使用。

7. **IDEA（International Data Encryption Algorithm）**：是一种对称加密算法，其设计简单，但在一定程度上提供了较高的安全性。

8. **Chacha20**：Chacha20 是一种流密码算法，由丹尼尔·J·伯恩斯坦（Daniel J. Bernstein）设计，被认为是一种安全性很高且性能优异的加密算法。Chacha20 提供了较高的安全性，并且在许多现代加密协议中得到了广泛应用，如TLS 1.3。

   

   


**这里考虑三种实现：Chacha20, Twofish-128，AES-256。**通常情况下，Chacha20由于其简单的设计和高效的实现，在软件和硬件上都能获得很好的性能。Twofish-128是一种块密码算法，在某些情况下可能会比AES-256慢，因为Twofish-128的轮数较多，但是在其他情况下，两者的性能可能会很接近。而AES-256由于其广泛应用和优化实现，在许多情况下都能提供很高的性能。

总的来说，这些算法都被认为是安全的，但在选择算法时应根据具体的需求和应用场景来进行评估。AES-256 由于其广泛的应用和被广泛认可的安全性，是许多应用的首选。Chacha20 和 Twofish-128 则提供了备选选择，适用于特定的安全性和性能需求。

几点实现的说明：

-    我都是把随机数放在密文前面传递，随机数向量长度有可能会不同；

-    Twofish-128需要手动将文本对齐到块16大小；另2种算法不需要填充；

## 4.1 chacha20

```go
// 使用 ChaCha20 加密消息
func EncryptChaCha20(plaintext, key []byte) ([]byte, error) {

	nonce, err := GenerateRandomKey(ChachaNonceLen) // 12 字节 nonce
	if err != nil {
		return nil, err
	}

	cipher, err := chacha20.NewUnauthenticatedCipher(key, nonce)
	if err != nil {
		return nil, err
	}

	ciphertext := make([]byte, ChachaNonceLen+len(plaintext))
	copy(ciphertext, nonce)
	cipher.XORKeyStream(ciphertext[ChachaNonceLen:], plaintext)
	return ciphertext, nil
}

// 使用 ChaCha20 解密消息
func DecryptChaCha20(ciphertext, key []byte) ([]byte, error) {
	if len(ciphertext) < (ChachaNonceLen + 1) {
		return nil, errors.New("too short, at least 13 bytes")
	}

	nonce := ciphertext[0:ChachaNonceLen]
	cipher, err := chacha20.NewUnauthenticatedCipher(key, nonce)
	if err != nil {
		return nil, err
	}
	plainText := make([]byte, len(ciphertext)-ChachaNonceLen)
	cipher.XORKeyStream(plainText, ciphertext[ChachaNonceLen:])
	return plainText, nil
}
```

测试代码

```go
// go test -run TestEncryptDecrypt
func TestEncryptChaCha20(t *testing.T) {
	// 生成随机密钥和随机 nonce
	key, _ := GenerateRandomKey(32) // 32 字节密钥
	fmt.Printf("%x \n", key)
	// 明文消息
	plaintext := []byte("Hello, ChaCha20!")

	// 加密
	ciphertext, err := EncryptChaCha20(plaintext, key)
	if err != nil {
		t.Errorf("Encryption failed: %v", err)
	}

	// 解密
	decryptedText, err := DecryptChaCha20(ciphertext, key)
	if err != nil {
		t.Errorf("Decryption failed: %v", err)
	}

	// 检查解密后的消息是否与原始消息相同
	if !bytes.Equal(plaintext, decryptedText) {
		t.Error("Decrypted text does not match plaintext")
	}
}
```



## 4.2 aes256

```go
// 使用了CTR模式进行加密，它需要一个初始化向量（IV）。
// 在这里，我们随机生成一个IV并将其作为密文的前16个字节。
// 在解密时，我们将前16个字节解释为IV，并将其与密文一起传递给解密函数。
func EncryptAES(plaintext, key []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	ciphertext := make([]byte, aes.BlockSize+len(plaintext))
	iv := ciphertext[:aes.BlockSize]
	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
		return nil, err
	}

	stream := cipher.NewCTR(block, iv)
	stream.XORKeyStream(ciphertext[aes.BlockSize:], plaintext)

	return ciphertext, nil
}

func DecryptAES(ciphertext, key []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	if len(ciphertext) < aes.BlockSize {
		return nil, fmt.Errorf("ciphertext too short")
	}

	iv := ciphertext[:aes.BlockSize]
	ciphertext = ciphertext[aes.BlockSize:]

	stream := cipher.NewCTR(block, iv)
	stream.XORKeyStream(ciphertext, ciphertext)

	return ciphertext, nil
}
```
测试代码
```go
func TestAES256(t *testing.T) {
	key := []byte("32-byte-long-key-123456789012345")

	plaintext := []byte("Hello, AES!")

	ciphertext, err := EncryptAES(plaintext, key)
	if err != nil {
		fmt.Println("Error encrypting:", err)
		return
	}
	fmt.Println("Ciphertext:", base64.StdEncoding.EncodeToString(ciphertext))

	decryptedText, err := DecryptAES(ciphertext, key)
	if err != nil {
		fmt.Println("Error decrypting:", err)
		return
	}
	fmt.Println("Decrypted text:", string(decryptedText))

	// 检查解密后的消息是否与原始消息相同
	if !bytes.Equal(plaintext, decryptedText) {
		t.Error("Decrypted text does not match plaintext")
	}
}
```



## 4.3 twofish128

```go
// PKCS7Padding 对明文进行填充以满足块大小，如果对齐了也需要填充一个块，否则无法识别是否填充了；
func PKCS7Padding(plaintext []byte, blockSize int) []byte {
	padding := blockSize - len(plaintext)%blockSize
	padtext := bytes.Repeat([]byte{byte(padding)}, padding)
	return append(plaintext, padtext...)
}

// PKCS7UnPadding removes PKCS#7 padding from the decrypted data
func PKCS7UnPadding(data []byte) ([]byte, error) {
	length := len(data)
	if length == 0 {
		return nil, errors.New("invalid data")
	}
	padding := int(data[length-1])
	if padding > length || padding == 0 {
		return nil, errors.New("invalid padding")
	}
	return data[:length-padding], nil
}

// EncryptTwofish 使用Twofish算法对数据进行加密，将12字节的IV放置密文前边
func EncryptTwofish(key, plaintext []byte) ([]byte, error) {
	// 创建Twofish块
	block, err := twofish.NewCipher(key)
	if err != nil {
		return nil, err
	}

	iv, err := GenerateRandomKey(block.BlockSize())
	if err != nil {
		return nil, err
	}

	// 填充明文以满足块大小
	plaintext1 := PKCS7Padding(plaintext, block.BlockSize())
	//fmt.Printf("after padding len = %v \n", len(plaintext1))

	// 使用CBC模式进行加密
	ciphertext := make([]byte, len(iv)+len(plaintext1))
	// 将IV放在密文最前面
	copy(ciphertext, iv)

	//fmt.Printf("ciptertext len = %v, \n", len(ciphertext))
	mode := cipher.NewCBCEncrypter(block, iv)
	mode.CryptBlocks(ciphertext[len(iv):], plaintext1)

	return ciphertext, nil
}

// DecryptTwofish decrypts the given ciphertext using Twofish-128 in CBC mode
func DecryptTwofish(key, ciphertext []byte) ([]byte, error) {
	block, err := twofish.NewCipher(key)
	if err != nil {
		return nil, err
	}

	// 至少有1个数据块，并且有一个IV 16字节
	if len(ciphertext) < (block.BlockSize() + block.BlockSize()) {
		return nil, errors.New("two short, at least 32bytes")
	}

	iv := ciphertext[0:block.BlockSize()]
	cipherTextReal := ciphertext[block.BlockSize():]
	// CBC mode decryption
	mode := cipher.NewCBCDecrypter(block, iv)
	plaintext := make([]byte, len(cipherTextReal))
	mode.CryptBlocks(plaintext, cipherTextReal)

	// Remove PKCS#7 padding
	plaintext, err = PKCS7UnPadding(plaintext)
	if err != nil {
		return nil, err
	}

	return plaintext, nil
}
```

测试代码

```go
func TestTwofish(t *testing.T) {
	key := []byte("16-byte-long-key")
	// 初始化向量（IV），与密钥一样长
	iv, _ := GenerateRandomKey(16)

	plaintext := []byte("Hello, Twofish!123")

	// 使用Twofish和CBC模式加密数据
	ciphertext, err := EncryptTwofish(key, iv, plaintext)
	if err != nil {
		fmt.Println("Error encrypting data:", err)
		return
	}

	decryptText, err := DecryptTwofish(key, iv, ciphertext)
	fmt.Println("Ciphertext:", string(decryptText))

	// 检查解密后的消息是否与原始消息相同
	if !bytes.Equal(plaintext, decryptText) {
		t.Error("Decrypted text does not match plaintext")
	}
}
```

# 5. 协议封装

```protobuf
message MsgKeyExchange {
  int64 keyPrint = 1;
  int64 rsaPrint = 2;
  int32 stage = 3;   // 当前处于状态机

  bytes tempKey = 4;  // 临时秘钥，需要RSA加密
  bytes pubKey = 5;   // 临时公钥，需要RSA加密

  string encType = 6;      // plain, rsa加密，对阵加密类型
  string status = 7;       // ok, fail
  string detail = 8;       // 错误信息
}
```

|                       | 发起端                                                       | 服务端                                                       |
| --------------------- | ------------------------------------------------------------ | ------------------------------------------------------------ |
| 1.1直接使用旧的秘钥； | MsgHello 中包含keyPrint，说明执行过协商;                     | 服务端检查缓存以及数据库是否有该信息；如果有则应该中添加keyPrint，如果找不到，则设置为0； |
| 1.2完成               | 客户端检查应答的MsgHello存在同样的keyPrint；后续使用cipher字段发送消息，需要手动序列化对象并加密，之后设置到该字段； | 服务端根据keyPrint字段提示了解密使用的秘钥；                 |
| 2.1新发起协商         | MsgHello 中包含是否需要加密；                                | 服务端需要清理旧的缓存，这可能是因为客户端丢失了会话秘钥；服务端需要应答RSA公钥指纹，让用户使用公钥协商； |
| 2.2交换公钥           | 使用MsgKeyExchange设置使用RSA指纹，加密临时秘钥，加密公钥，设置算法，发给对方； | 如果自己没有指纹对应的RSA，则说明证书不匹配；无法协商，告诉客户错误；<br/>如果算法不支持，也返回错误；<br>检查如果可以解密，则设置临时秘钥，使用临时秘钥以及算法加密自己的公钥，以及计算的共享密钥指纹； |
| 2.3                   | 客户收到服务端的公钥后，计算共享密钥和共享密钥指纹并比对，如果正确，需要发送协商完毕的消息； | 收到协商完毕的消息，就可以进入加密通信模式了；               |
| 3. 用户会话级别协商   | 用户设置信息时，会自动生成RSA公私钥对，会一起将公钥发送到服务器上；<br>客户之一可以发起进入点对点加密模式， | 接收方的处理与服务端类似；                                   |
|                       | 协商完毕后，后续使用key对后续的MsgChat的内容字段加密；服务端无法解密；协商的秘钥也不存在于服务器端； |                                                              |
|                       |                                                              |                                                              |
|                       |                                                              |                                                              |

