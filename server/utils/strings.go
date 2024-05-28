package utils

import (
	"encoding/base64"
	"encoding/binary"
	"encoding/hex"
	"errors"
	"github.com/yanyiwu/gojieba"
	"math"
	"os"
	"path/filepath"
	"regexp"
	"runtime"
	"unicode"
)

func BytesToHexStr(data []byte) string {
	return hex.EncodeToString(data)
}

func HexStrToBytes(hexStr string) ([]byte, error) {
	// 解码十六进制字符串为字节数组
	data, err := hex.DecodeString(hexStr)
	if err != nil {
		return nil, err
	}
	return data, nil
}

// GetAbsolutePath 获取相对路径的绝对路径
func GetAbsolutePath(relativePath string) (string, error) {
	// 检查输入路径是否是绝对路径
	if filepath.IsAbs(relativePath) {
		return relativePath, nil
	}

	// 获取当前可执行文件的路径
	exePath, err := os.Executable()
	if err != nil {
		return "", err
	}

	// 判断操作系统类型
	var absPath string
	if runtime.GOOS == "windows" {
		// 如果是 Windows，使用 filepath.Join 拼接路径
		absPath = filepath.Join(filepath.Dir(exePath), relativePath)
	} else {
		// 如果是 Linux，直接拼接路径
		absPath = filepath.Join(filepath.Dir(exePath), relativePath)
	}

	// 返回绝对路径
	return absPath, nil
}

func BytesToInt64(data []byte) (int64, error) {
	if len(data) < 8 {
		return 0, errors.New("insufficient bytes to convert to int64")
	}
	return int64(binary.LittleEndian.Uint64(data[:8])), nil
}

func Int64ToInt16(i int64) int16 {
	module := math.MaxInt16
	return int16(i % int64(module))
}

func DecodeBase64(str string) ([]byte, error) {
	// 解码base64字符串为字节切片
	decodedBytes, err := base64.StdEncoding.DecodeString(str)
	return decodedBytes, err
}

// go:inline
func EncodeBase64(data []byte) string {
	// 解码base64字符串为字节切片
	return base64.StdEncoding.EncodeToString(data)

}

// IsValidEmail 检查电子邮件地址是否合法
func IsValidEmail(email string) bool {
	// 使用正则表达式匹配电子邮件地址的模式
	pattern := `^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`

	// 编译正则表达式
	reg := regexp.MustCompile(pattern)

	// 使用正则表达式判断是否匹配
	return reg.MatchString(email)
}

// 合并表
// mergeMap 合并两个map，并将结果存回第一个map
func MergeMap[K comparable, V any](m1, m2 map[K]V) {
	for k, v := range m2 {
		m1[k] = v
	}
}

func MergeMapMask(m1, m2 map[int64]uint32) {
	for k, v := range m2 {
		v1, ok := m1[k]
		if ok {
			m1[k] = v1 | v
		} else {
			m1[k] = v
		}
	}
}

//func SegmentText(str string) []string {
//	// 载入词典
//	var segmenter sego.Segmenter
//	segmenter.LoadDictionary("../ttf/dictionary.txt")
//
//	// 进行分词
//	text := []byte(str)
//	segments := segmenter.Segment(text)
//
//	// 将分词结果转换为 []string
//	result := sego.SegmentsToSlice(segments, false)
//
//	return result
//}

// 对英文进行分词
// "LifeIsFullOfChancesAndChallenges"
func SplitCamelCase(str string) []string {
	// 使用正则表达式匹配大写字母前的小写字母
	re := regexp.MustCompile(`[a-z]+|[A-Z][a-z]*`)
	words := re.FindAllString(str, -1)
	return words
}

// 中文进行分词
func SegmentTextChinese(text string) []string {
	// 进行分词
	x := gojieba.NewJieba()
	defer x.Free()

	words := x.Cut(text, true)
	tags := []string{}
	for _, word := range words {
		b := IsAllPunctuationOrWhitespace(word)
		if !b {
			tags = append(tags, word)
		}
	}
	return tags
}

// 检测字符串是否全是标点或空白
func IsAllPunctuationOrWhitespace(str string) bool {
	for _, runeValue := range str {
		// 检查字符是否为标点或空白
		if !unicode.IsPunct(runeValue) && !unicode.IsSpace(runeValue) {
			return false
		}
	}
	return true
}

func DetectLanguage(s string) string {
	for _, r := range s {
		switch {
		case unicode.Is(unicode.Han, r): // 汉字
			return "Chinese"
		case unicode.Is(unicode.Katakana, r), unicode.Is(unicode.Hiragana, r): // 日文片假名和平假名
			return "Japanese"
		case unicode.Is(unicode.Arabic, r): // 阿拉伯文
			return "Arabic"
		case unicode.Is(unicode.Cyrillic, r): // 西里尔文
			return "Russian"
		case unicode.Is(unicode.Latin, r): // 拉丁文
			return "Latin"
		case unicode.Is(unicode.Greek, r): // 希腊文
			return "Greek"
		case unicode.Is(unicode.Hebrew, r): // 希伯来文
			return "Hebrew"
		case unicode.Is(unicode.Bengali, r): // 孟加拉文
			return "Bengali"
		case unicode.Is(unicode.Devanagari, r): // 天城文
			return "Hindi"
		case unicode.Is(unicode.Tamil, r): // 泰米尔文
			return "Tamil"
		case unicode.Is(unicode.Thai, r): // 泰文
			return "Thai"
		}
		if isFrench(r) {
			return "French"
		}

		if isGerman(r) {
			return "German"
		}
	}
	return "English"
}

func isFrench(r rune) bool {
	// 法语中的特定字符范围
	frenchRanges := []*unicode.RangeTable{
		&unicode.RangeTable{
			R16: []unicode.Range16{
				{0x00C0, 0x00FF, 1}, // Latin-1 Supplement
			},
			R32: []unicode.Range32{
				{0x0152, 0x0153, 1}, // Latin Extended-A
				{0x0178, 0x0178, 1}, // Latin Extended-A
				{0x0192, 0x0192, 1}, // Latin Extended-A
				{0x02C6, 0x02C6, 1}, // Latin Extended-A
				{0x02C9, 0x02CB, 1}, // Latin Extended-A
				{0x02D8, 0x02DD, 1}, // Latin Extended-A
				{0x02DA, 0x02DA, 1}, // Latin Extended-A
				{0x02DC, 0x02DC, 1}, // Latin Extended-A
				{0x02E0, 0x02E4, 1}, // Latin Extended-A
				{0x1E00, 0x1EFF, 1}, // Latin Extended-A
				{0x2013, 0x2014, 1}, // Latin Extended-A
				{0x2018, 0x2019, 1}, // Latin Extended-A
				{0x201C, 0x201D, 1}, // Latin Extended-A
				{0x2020, 0x2022, 1}, // Latin Extended-A
				{0x2026, 0x2026, 1}, // Latin Extended-A
				{0x2030, 0x2030, 1}, // Latin Extended-A
				{0x2039, 0x203A, 1}, // Latin Extended-A
				{0x20AC, 0x20AC, 1}, // Latin Extended-A
				{0x2122, 0x2122, 1}, // Latin Extended-A
			},
		},
	}

	for _, charRange := range frenchRanges {
		if unicode.In(r, charRange) {
			return true
		}
	}

	return false
}

func isGerman(r rune) bool {
	// 德语中的特定字符范围
	germanRanges := []*unicode.RangeTable{
		&unicode.RangeTable{
			R16: []unicode.Range16{
				{0x00C4, 0x00D6, 1}, // Latin-1 Supplement
				{0x00DC, 0x00DC, 1}, // Latin-1 Supplement
				{0x00DF, 0x00E4, 1}, // Latin-1 Supplement
				{0x00E4, 0x00F6, 1}, // Latin-1 Supplement
				{0x00F6, 0x00FC, 1}, // Latin-1 Supplement
				{0x00FC, 0x00FF, 1}, // Latin-1 Supplement
			},
		},
	}

	for _, charRange := range germanRanges {
		if unicode.In(r, charRange) {
			return true
		}
	}

	return false
}
