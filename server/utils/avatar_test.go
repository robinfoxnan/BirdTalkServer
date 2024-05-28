package utils

import (
	"fmt"
	"testing"
)

func TestCreateAvatar(t *testing.T) {
	GenerateAvatar("Bird Fish", 0, "Bird.png")
	GenerateAvatar("飞鸟真人", 1, "飞鸟.png")
	GenerateAvatar("Жозефина", 1, "约瑟芬.png")
	GenerateAvatar("Ярослав", 1, "雅罗斯拉夫.png")
	GenerateAvatar("ちびまる子ちゃん", 1, "日语.png")

}

func TestCreateThumb(t *testing.T) {
	// 调用 CreateThumbImage 函数生成缩略图
	b, err := CreateThumbImage34("e:\\test\\截图.jpg", "e:\\test\\2.jpg", 300)
	if err != nil {
		fmt.Println(b, err)
		return
	}
	fmt.Println("Thumbnail image created successfully.")
}

func TestCreateThumb1(t *testing.T) {
	// 调用 CreateThumbImage 函数生成缩略图
	b, err := ScaleWithAspectRatio("e:\\test\\截图.jpg", "e:\\test\\2-1.jpg", 300)
	if err != nil {
		fmt.Println(b, err)
		return
	}
	fmt.Println("Thumbnail image created successfully.")
}
