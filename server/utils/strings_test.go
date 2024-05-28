package utils

import (
	"fmt"
	"strings"
	"testing"
)

func TestSegment(t *testing.T) {
	str := "初一英文下册5单元单词练习.docx"
	//str = "ALifeWithout AFriend IsLikealife without the sun"
	ret := SegmentTextChinese(str)
	fmt.Println(strings.Join(ret, "/"))

	str = "LifeIsFullOfChancesAndChallenges"
	words := SplitCamelCase(str)
	fmt.Println(strings.Join(words, "/"))

}
