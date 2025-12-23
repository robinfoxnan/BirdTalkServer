package main

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"io"
	"strconv"
	"strings"
	"time"
)

// EncryptFileToken 使用 AES-GCM 加密文件名和过期时间，生成 token
func EncryptFileToken(secretKey []byte, filename string, expires int64) (string, error) {
	// 拼接要加密的数据
	// 拼接数据：expires|filename
	data := fmt.Sprintf("%d|%s", expires, filename)
	plaintext := []byte(data)

	block, err := aes.NewCipher(secretKey)
	if err != nil {
		return "", err
	}

	aesGCM, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}

	// 生成随机 nonce
	nonce := make([]byte, aesGCM.NonceSize())
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return "", err
	}

	// 加密
	ciphertext := aesGCM.Seal(nil, nonce, plaintext, nil)

	// 将 nonce + 密文一起编码成 base64
	//token := base64.URLEncoding.EncodeToString(append(nonce, ciphertext...))
	token := base64.RawURLEncoding.EncodeToString(append(nonce, ciphertext...))
	return token, nil
}

// 解密 token 获取文件名和过期时间
func DecryptFileToken(secretKey []byte, token string) (expires int64, filename string, err error) {
	//raw, err := base64.URLEncoding.DecodeString(token)
	raw, err := base64.RawURLEncoding.DecodeString(token)
	if err != nil {
		return 0, "", err
	}

	block, err := aes.NewCipher(secretKey)
	if err != nil {
		return 0, "", err
	}

	aesGCM, err := cipher.NewGCM(block)
	if err != nil {
		return 0, "", err
	}

	nonceSize := aesGCM.NonceSize()
	if len(raw) < nonceSize {
		return 0, "", fmt.Errorf("invalid token")
	}

	nonce, ciphertext := raw[:nonceSize], raw[nonceSize:]

	plaintext, err := aesGCM.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		return 0, "", err
	}

	// 使用 strings.SplitN，防止文件名里有 | 时出错
	parts := strings.SplitN(string(plaintext), "|", 2)
	if len(parts) != 2 {
		return 0, "", fmt.Errorf("invalid plaintext format")
	}

	// 先解析 expires
	expires, err = strconv.ParseInt(parts[0], 10, 64)
	if err != nil {
		return 0, "", err
	}

	// 后面就是文件名
	filename = parts[1]

	return expires, filename, nil
}

func testAll() {
	key := make([]byte, 32) // 32 bytes = 256 bit
	_, err := rand.Read(key)
	if err != nil {
		panic(err)
	}
	fmt.Println("生成的密钥:", hex.EncodeToString(key))

	filename := "o_abc.jpg"
	expires := time.Now().Add(5 * time.Minute).Unix()

	token, err := EncryptFileToken(key, filename, expires)
	if err != nil {
		panic(err)
	}
	fmt.Println("生成 token:", token)

	// 模拟 CDN 解密
	exp, name, err := DecryptFileToken(key, token)
	if err != nil {
		fmt.Println(err)
		return
	}

	// 转换为可读时间字符串
	t := time.Unix(exp, 0) // 秒级时间戳转换
	//fmt.Println("可读时间:", t.Format(time.RFC3339)) // 输出格式：2006-01-02T15:04:05Z07:00

	// 也可以自定义格式
	fmt.Println("自定义格式:", t.Format("2006-01-02 15:04:05"))
	fmt.Printf("解密后 filename=%s, expires=%d\n", name, exp)
}

func test1(keyHex, tokenStr string) {

	// 反向解码成 []byte
	key, err := hex.DecodeString(keyHex)
	if err != nil {
		panic(err)
	}

	// 模拟 CDN 解密
	exp, name, err := DecryptFileToken(key, tokenStr)
	if err != nil {
		fmt.Println(err)
		return
	}

	// 转换为可读时间字符串
	t := time.Unix(exp, 0) // 秒级时间戳转换
	//fmt.Println("可读时间:", t.Format(time.RFC3339)) // 输出格式：2006-01-02T15:04:05Z07:00

	// 也可以自定义格式
	fmt.Println("自定义格式:", t.Format("2006-01-02 15:04:05"))
	fmt.Printf("解密后 filename=%s, expires=%d\n", name, exp)
}

// openssl rand -hex 32
func main() {
	//test1("adffbde3eac40de08ab49d6f48cea72ee78df550479c11caea0e173e9fea7794",
	//	"KRf9VJBSZp8bUknV0xP2Xn5LeldpCC6N92osHZJ3kT16DHsq8eCEvIpLkpyBmIq1")

	// 创建一个能用的

	key, err := hex.DecodeString("278b3e9f1c18f488314c7991e163c8bac880d5f09f4a15dc7589823dc6b43264")
	if err != nil {
		panic(err)
	}
	sign, err := EncryptFileToken(key, "3zi18u3l5inm.jpg", time.Now().Unix())
	print(sign)
	// 7MwoETBf4VzHrcUeS9HDlWjJTAcPKpdafXQSEBSbjJ7-M4NMBS8mpjvsp-GwnOVGqqVaO2tKlA

	test1("278b3e9f1c18f488314c7991e163c8bac880d5f09f4a15dc7589823dc6b43264",
		"7MwoETBf4VzHrcUeS9HDlWjJTAcPKpdafXQSEBSbjJ7-M4NMBS8mpjvsp-GwnOVGqqVaO2tKlA")

}
