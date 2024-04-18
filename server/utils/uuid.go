package utils

import (
	"github.com/google/uuid"
	"math/rand"
	"sync"
	"time"
)

// GenerateUUID 生成UUID字符串
func GenerateUUID() string {
	u := uuid.New()
	return u.String()
}

func GenerateShortID() string {
	return generateRandomID(16)
}

var onceShort sync.Once

const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
const shortIDLength = 16
const checkCodeLength = 4
const numset = "0123456789"

func generateRandomID(length int) string {
	onceShort.Do(func() {
		rand.Seed(time.Now().UnixNano())
	})
	if length < shortIDLength {
		length = shortIDLength
	}

	shortID := make([]byte, length)
	for i := 0; i < length; i++ {
		index := rand.Intn(len(charset))
		shortID[i] = charset[index]
	}

	return string(shortID)
}

// 计算一个随机的5位验证码
func GenerateCheckCode(length int) string {
	onceShort.Do(func() {
		rand.Seed(time.Now().UnixNano())
	})
	if length < checkCodeLength {
		length = checkCodeLength
	}

	shortID := make([]byte, length)
	for i := 0; i < length; i++ {
		index := rand.Intn(len(numset))
		shortID[i] = numset[index]
	}
	return string(shortID)
}
