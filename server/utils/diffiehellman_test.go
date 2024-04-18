package utils

import (
	"bytes"
	"crypto/elliptic"
	"encoding/base64"
	"fmt"
	"testing"
)

func TestCompareBytes(t *testing.T) {
	tests := []struct {
		name     string
		array1   []byte
		array2   []byte
		expected bool
	}{
		{"EqualArrays", []byte{1, 2, 3, 4, 5}, []byte{1, 2, 3, 4, 5}, true},
		{"DifferentArrays", []byte{1, 2, 3, 4, 5}, []byte{1, 2, 3, 4, 6}, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if result := CompareBytes(tt.array1, tt.array2); result != tt.expected {
				t.Errorf("Test failed for %s: expected %t, got %t", tt.name, tt.expected, result)
			}
		})
	}
}

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

func TestAES256(t *testing.T) {
	key := []byte("32-byte-long-key-123456789012345")

	plaintext := []byte("Hello, AES!")

	ciphertext, err := EncryptAES_CTR(plaintext, key)
	if err != nil {
		fmt.Println("Error encrypting:", err)
		return
	}
	fmt.Println("Ciphertext:", base64.StdEncoding.EncodeToString(ciphertext))

	decryptedText, err := DecryptAES_CTR(ciphertext, key)
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

// go test -run TestTwofish
func TestTwofish(t *testing.T) {
	key := []byte("16-byte-long-key")

	plaintext := []byte("Hello, Twofish!123")

	// 使用Twofish和CBC模式加密数据
	ciphertext, err := EncryptTwofish(key, plaintext)
	if err != nil {
		fmt.Println("Error encrypting data:", err)
		return
	}

	decryptText, err := DecryptTwofish(key, ciphertext)
	fmt.Printf("plainText =%s\nCiphertext=%s\n", string(plaintext), string(decryptText))
	fmt.Printf("plainText =%d\nCiphertext=%d\n", len(plaintext), len(decryptText))

	// 检查解密后的消息是否与原始消息相同
	if !bytes.Equal(plaintext, decryptText) {
		t.Error("Decrypted text does not match plaintext")
	}
}

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
	aliceSharedSecret, err := sharedSecret1(curve, bobPublicKey, alicePrivateKey)
	if err != nil {
		panic(err)
	}
	bobSharedSecret, err := sharedSecret1(curve, alicePublicKey, bobPrivateKey)
	if err != nil {
		panic(err)
	}

	// 验证共享密钥是否相同
	if !CompareBytes(aliceSharedSecret, bobSharedSecret) {
		panic("Shared secrets do not match!")
	}
	fmt.Println("len = ", len(aliceSharedSecret))

	//fmt.Printf("Shared Secret: %s\n", aliceSharedSecret.String())
}

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
