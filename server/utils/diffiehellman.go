package utils

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/md5"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/subtle"
	"crypto/x509"
	"encoding/hex"
	"encoding/pem"
	"errors"
	"fmt"
	"golang.org/x/crypto/chacha20"
	"golang.org/x/crypto/twofish"
	"io"
	"math/big"
	"os"
	"strings"
)

// 1. 发起端(客户端使用对方公钥加密临时秘钥给对方)
// 2. 服务端将ga作为公钥发给客户,（128位），使用临时秘钥加密公钥，算法为ChaCha20，256位（32字节）
// 3. 客户将gb （128位）作为公钥发给服务端， （使用临时对称密钥加密）；此时双方已经可以计算出共享密钥；用共享密钥加密临时秘钥，并计算MD5，
// 4. 服务端计算共享密钥，用共享密钥加密临时秘钥，并计算MD5，校验；
// 5， 如果校验成功，应答OK，如果失败应答FAIL

const (
	//PrimeLength    = 128
	//PrimeG         = 2 // 2、3或5等较小的素数。
	ChachaNonceLen = 12

	KeyExchangeStageNonce     = 1
	KeyExchangeStagePublicKey = 2
	KeyExchangeStageShareKey  = 3
	KeyExchangeStageValidate  = 4
	KeyExchangeStageFinished  = 5
)

var algorithmList = []string{"chacha20", "twofish128", "aes-ctr"}
var curve elliptic.Curve = elliptic.P256()

// 使用椭圆曲线计算共享密钥
type KeyExchange struct {
	Role    string // 发起方或者被动方
	TempKey []byte // 临时对称密钥，使用RSA加密
	Stage   int16  // 当前状态
	EncType string // 对称加密类型，三种

	PrivateKey      *big.Int
	PublicKey       []byte
	PublicKeyRemote []byte

	SharedKey      []byte // 秘钥交换的结果
	SharedKeyHash  []byte // 共享密钥通过SHA256得到一个哈希，用于对称秘钥
	SharedKeyPrint int64  // md hash 的前8个字节转为int64
}

func checkAlgorithm(enc string) bool {
	littleEnc := strings.ToLower(enc)
	for _, item := range algorithmList {
		if item == littleEnc {
			return true
		}
	}
	return false
}

// "s","c"
func NewKeyExchange(enc string) (*KeyExchange, error) {
	keyEx := KeyExchange{Role: "c",
		TempKey:       nil,
		SharedKey:     nil,
		SharedKeyHash: nil,
		Stage:         KeyExchangeStageNonce,
		EncType:       enc,
	}

	ok := checkAlgorithm(enc)
	if !ok {
		return nil, errors.New("encrypt algorithm is unknown")
	}

	var err error
	// 这里就生成了密钥对
	keyEx.PrivateKey, keyEx.PublicKey, err = generateDHKeyPair(curve)
	if err != nil {
		return nil, err
	}
	return &keyEx, err
}

// 发起方需要生成临时对称密钥
// "chacha20" "twofish128"  "aes256"
func (k *KeyExchange) GenTempKey(t string) error {
	k.Role = "c"
	var err error
	k.EncType = t
	if t == "chacha20" {
		k.TempKey, err = GenerateRandomKey(32) // 32 字节密钥
		return err
	} else if "aes256" == t {
		k.TempKey, err = GenerateRandomKey(32) // 32 字节密钥
		return err
	} else if "twofish128" == t {
		k.TempKey, err = GenerateRandomKey(16) // 16 字节密钥

		return err
	}

	return errors.New("error encrypt type.")
}

// 被动方需要接收临时秘钥，这里是被动方，使用自己的RSA私钥解码后的
func (k *KeyExchange) SetTempKeyAndNonce(key, nonce []byte, enc string) error {
	k.Role = "s"
	if enc != "chacha20" && enc != "twofish128" && enc != "aes256" {
		return errors.New("enc type unknow")
	}

	k.EncType = enc
	k.TempKey = key
	k.Stage = KeyExchangeStageNonce
	return nil
}

// 收到对方的公钥后，通过设置公钥就可以生成秘钥了
func (k *KeyExchange) GenShareKey(remotePublicKey []byte) (int64, error) {
	if k.Stage != KeyExchangeStagePublicKey {
		return 0, errors.New("stage is not having public key")
	}

	var err error
	k.PublicKeyRemote = remotePublicKey
	k.SharedKey, err = sharedSecretSPKI(curve, remotePublicKey, k.PrivateKey)
	if err != nil {
		return 0, err
	}

	if len(k.SharedKey) < 32 {
		return 0, errors.New("share key less than 32 bytes")
	}
	// 计算新的对称密钥

	switch strings.ToLower(k.EncType) {
	case "chacha20":
		k.SharedKeyHash = k.SharedKey
		//calculateSHA256(k.SharedKey)

	case "aes-ctr":
		k.SharedKeyHash = k.SharedKey
		//calculateSHA256(k.SharedKey)
	case "twofish128":
		k.SharedKeyHash = calculateMD5(k.SharedKey)
	default:
		k.SharedKeyHash = k.SharedKey
	}

	// 计算共享密钥指纹，用户后续直接使用，可以用于做SessionID
	k.SharedKeyPrint, err = BytesToInt64(k.SharedKeyHash)

	if err == nil {
		k.Stage = KeyExchangeStageShareKey
	}
	return k.SharedKeyPrint, err
}

// ////////////////////////////////////////////////////////////////////////
// 计算MD5，得到16字节
func calculateMD5(data []byte) []byte {
	hash := md5.New()
	hash.Write(data)
	return hash.Sum(nil)
}

// 计算SHA256,得到32字节
func calculateSHA256(data []byte) []byte {
	hash := sha256.Sum256(data)
	hashString := hex.EncodeToString(hash[:])
	fmt.Println("SHA-256 Hash:", hashString)
	return hash[:]
}

func GenerateRandomKey(length int) ([]byte, error) {

	key := make([]byte, length)
	_, err := rand.Read(key)
	if err != nil {
		return nil, err
	}
	return key, nil
}

// compareBytes 比较两个字节数组是否完全一致
func CompareBytes(a, b []byte) bool {
	// 使用crypto/subtle包中的ConstantTimeCompare函数进行比较
	return subtle.ConstantTimeCompare(a, b) == 1
}

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

// 使用了CTR模式进行加密，它需要一个初始化向量（IV）。
// 在这里，我们随机生成一个IV并将其作为密文的前16个字节。
// 在解密时，我们将前16个字节解释为IV，并将其与密文一起传递给解密函数。
func EncryptAES_CTR(plaintext, key []byte) ([]byte, error) {
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

func DecryptAES_CTR(ciphertext, key []byte) ([]byte, error) {
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

// ///////////////////////////////////////////////////////////////////////
func encodeSPKIPublicKey(publicKey *ecdsa.PublicKey) ([]byte, error) {

	// 将公钥编码为 DER 格式
	publicKeyDER, err := x509.MarshalPKIXPublicKey(publicKey)
	if err != nil {
		return nil, err
	}

	fmt.Println("DER public key:", publicKeyDER)

	// 创建一个 PEM 格式的数据块
	publicKeyPEM := &pem.Block{
		Type:  "PUBLIC KEY",
		Bytes: publicKeyDER,
	}

	// 将 PEM 格式的数据块转换为字节数组
	pemBytes := pem.EncodeToMemory(publicKeyPEM)

	return pemBytes, nil
}

// generateDHKeyPair 生成椭圆曲线的密钥对
func generateDHKeyPair(curve elliptic.Curve) (*big.Int, []byte, error) {
	// 生成私钥
	privateKey, err := rand.Int(rand.Reader, curve.Params().N)
	if err != nil {
		return nil, nil, err
	}

	// 计算公钥
	publicKeyX, publicKeyY := curve.ScalarBaseMult(privateKey.Bytes())
	//publicKey := elliptic.MarshalCompressed(curve, publicKeyX, publicKeyY) // 使用压缩格式
	publicKeyCurve := ecdsa.PublicKey{Curve: curve, X: publicKeyX, Y: publicKeyY}
	publicKeySPKI, err := encodeSPKIPublicKey(&publicKeyCurve)

	return privateKey, publicKeySPKI, err
}

// sharedSecret 计算共享密钥
//func sharedSecret(curve elliptic.Curve, publicKey []byte, privateKey *big.Int) ([]byte, error) {
//	// 解码公钥
//	x, y := elliptic.UnmarshalCompressed(curve, publicKey)
//	if x == nil {
//		return nil, fmt.Errorf("invalid public key")
//	}
//
//	// 计算共享密钥
//	sharedKeyX, sharedKeyY := curve.ScalarMult(x, y, privateKey.Bytes())
//	sharedKey := elliptic.Marshal(curve, sharedKeyX, sharedKeyY) // 返回整个点的字节表示
//
//	return sharedKey, nil
//}

func sharedSecret1(curve elliptic.Curve, publicKey []byte, privateKey *big.Int) ([]byte, error) {
	// 解码公钥
	x, y := elliptic.UnmarshalCompressed(curve, publicKey)
	if x == nil {
		return nil, fmt.Errorf("invalid public key")
	}

	// 计算共享密钥
	sharedKeyX, _ := curve.ScalarMult(x, y, privateKey.Bytes())
	sharedKey := sharedKeyX.Bytes()

	// 如果共享密钥长度超过32字节，截取前32个字节
	if len(sharedKey) > 32 {
		sharedKey = sharedKey[:32]
	}

	return sharedKey, nil
}

func sharedSecretSPKI(curve elliptic.Curve, publicKey []byte, privateKey *big.Int) ([]byte, error) {
	// 解码公钥
	publicKeySPKI, err := parsePublicKey(publicKey)
	if err != nil {
		return nil, err
	}

	// 计算共享密钥
	sharedKeyX, _ := curve.ScalarMult(publicKeySPKI.X, publicKeySPKI.Y, privateKey.Bytes())
	sharedKey := sharedKeyX.Bytes()

	// 如果共享密钥长度超过32字节，截取前32个字节
	//if len(sharedKey) > 32 {
	//	sharedKey = sharedKey[:32]
	//}

	return sharedKey, nil
}

/*
JavaScript 中使用 SPKI 格式导出的公钥可以在 Go 中解码和使用。
了在 Go 中解码 JavaScript 导出的公钥，需要使用 crypto/x509 包中的函数来解析公钥。
*/
func parsePublicKey(publicKeyPEM []byte) (*ecdsa.PublicKey, error) {
	// 解码 PEM 格式的公钥
	block, _ := pem.Decode(publicKeyPEM)
	if block == nil {
		return nil, fmt.Errorf("failed to decode PEM block containing public key")
	}
	// 解析公钥
	publicKey, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		return nil, fmt.Errorf("failed to parse public key: %v", err)
	}
	// 转换为 ECDSA 公钥
	ecdsaPublicKey, ok := publicKey.(*ecdsa.PublicKey)
	if !ok {
		return nil, fmt.Errorf("failed to convert to ECDSA public key")
	}
	return ecdsaPublicKey, nil
}

// generateLargePrime 生成一个指定位数的大素数
//func generateLargePrime(bits int) (*big.Int, error) {
//	// 创建一个指定位数的随机整数
//	n, err := rand.Int(rand.Reader, big.NewInt(1).Lsh(big.NewInt(1), uint(bits)))
//	if err != nil {
//		return nil, err
//	}
//
//	// 使用Miller-Rabin素性检验算法测试随机整数是否为素数
//	// 多次测试以增加准确率
//	prime := big.NewInt(0).Set(n)
//	isProbablePrime := prime.ProbablyPrime(20) // 20次迭代通常足够安全
//	if !isProbablePrime {
//		return generateLargePrime(bits) // 如果不是素数，递归调用自身以生成新的随机数
//	}
//
//	return prime, nil
//}

// （对方的公钥）做幂运算，指数为自己的私钥，对素数取mod
//func GenDhShare(remotePublicKey, privateKey, prime *big.Int) *big.Int {
//	sharedKey := new(big.Int).Exp(remotePublicKey, privateKey, prime)
//	return sharedKey
//}
