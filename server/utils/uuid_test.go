package utils

import (
	"fmt"
	"testing"
)

func TestGetCheckCode(t *testing.T) {

	for i := 0; i < 100; i++ {
		str := GenerateCheckCode(5)
		fmt.Println(str)
	}
}
