package utils

import (
	"birdtalk/server/utils/myavatar"
	"fmt"
	"github.com/disintegration/imaging"
	"github.com/fogleman/gg"
	"github.com/golang/freetype/truetype"
	"image"
	"image/color"
	"image/draw"
	"image/png"
	"os"
	"unicode/utf8"
)

var defaultFont *truetype.Font = nil

func InitFont(filepath string) error {
	fmt.Println("font path", filepath)
	font, err := loadFontToMemory(filepath)
	if err != nil {
		return err
	}
	defaultFont = font
	return nil
}

func initFont() {
	if defaultFont != nil {
		return
	}
	strFilePath := "./ttf/SourceHanSans-VF.ttf"
	//strFilePath = "C:\\Windows\\Fonts\\simfang.ttf"
	font, err := loadFontToMemory(strFilePath)
	if err != nil {
		return
	}
	defaultFont = font
}

func GenerateAvatar(name string, gender int, saveFilePath string) error {

	onceShort.Do(func() {
		initFont()
	})

	// 生成圆形图标
	avatar, err := generateAvatar(name, 100, gender)
	if err != nil {
		return err
	}

	// 创建相同尺寸的圆形蒙版
	mask := createCircularMask(avatar.Bounds())

	// 创建一个新的透明背景图像
	bgImg := image.NewRGBA(avatar.Bounds())

	// 使用圆形蒙版将文字图标覆盖到透明背景图像上
	draw.DrawMask(bgImg, bgImg.Bounds(), avatar, image.Point{}, mask, image.Point{}, draw.Over)

	// 保存图像到文件
	file, err := os.Create(saveFilePath)
	if err != nil {
		return err
	}
	defer file.Close()

	err = png.Encode(file, bgImg)
	if err != nil {
		return err
	}

	//fmt.Println("Avatar saved as avatar.png")
	return nil
}

// LoadFontToMemory 加载字体文件到内存中，返回一个粗体
func loadFontToMemory(fontPath string) (*truetype.Font, error) {
	// Read the font file into memory
	fontBytes, err := os.ReadFile(fontPath)
	if err != nil {
		return nil, err
	}

	// Parse the font file
	font, err := truetype.Parse(fontBytes)
	if err != nil {
		return nil, err
	}

	return font, nil
}

// createCircularMask 创建一个圆形蒙版，与指定的边界大小相同
func createCircularMask(bounds image.Rectangle) *image.Alpha {
	dc := gg.NewContext(bounds.Dx(), bounds.Dy())
	dc.DrawCircle(float64(bounds.Dx())/2, float64(bounds.Dy())/2, float64(bounds.Dx())/2)
	dc.Fill()

	mask := image.NewAlpha(bounds)
	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			_, _, _, a := dc.Image().At(x, y).RGBA()
			mask.Set(x, y, color.Alpha{A: uint8(a >> 8)})
		}
	}
	return mask
}

func generateAvatar(name string, sz int, gender int) (image.Image, error) {

	if defaultFont == nil {
		// 汉语会出错
	}
	myColor := color.RGBA{0x7f, 0xd3, 0xfa, 0xff} // 蓝色
	if gender == 0 {
		myColor = color.RGBA{0xf6, 0xa1, 0xbe, 0xff} // 粉色
	}
	options := &myavatar.Options{
		Font: defaultFont,
		Palette: []color.Color{
			myColor,
		},
	}

	first6Letter := TakeFirstCharacters(name)
	fmt.Println("返回：", first6Letter)
	avatar, err := myavatar.Draw(sz, first6Letter, options)

	if err != nil {
		return nil, err
	}

	return avatar, nil
}

// TakeFirstNRunes 返回字符串的前N个字符，确保总字节数不超过4个字节
// TakeFirstCharacters 根据字符类型取字符串的前N个字符
func TakeFirstCharacters(str string) string {

	result := ""
	count := 0
	for len(str) > 0 {
		r, size := utf8.DecodeRuneInString(str)
		//fmt.Printf("字符: %c, 字节数: %d\n", r, size)
		if (count + size) > 6 {
			break
		} else {
			result += string(r)
			count += size
		}

		str = str[size:] // 将已解码的字符从字符串中去除
	}

	return result
}

func TakeFirstNRunes(str string, n int) string {
	var result string
	var bytesCount int
	for _, runeValue := range str {
		runeSize := utf8.RuneLen(runeValue)
		if bytesCount+runeSize > 4 {
			break
		}
		bytesCount += runeSize
		result += string(runeValue)
		if bytesCount == 4 {
			break
		}
	}
	return result
}

// CreateThumbImage 函数用于生成缩略图，使用邻居插值方法。
// 这个函数主要是社区的图片用，为了布局好看
func CreateThumbImage34(filePath, destFilePath string, maxWidth int) (bool, error) {
	// 打开图片文件
	src, err := imaging.Open(filePath)
	if err != nil {
		return false, fmt.Errorf("failed to open image: %v", err)
	}

	thumbWidth := maxWidth
	thumbHeight := thumbWidth * 4.0 / 3.0

	width := src.Bounds().Dx()
	height := src.Bounds().Dy()
	if width < thumbWidth || height < thumbHeight {
		return false, nil
	}

	// 检查原始图像的宽高比是否为 3:4，如果不是，则裁剪为 3:4 的比例，取中间部分
	aspectRatio := float64(width) / float64(height)
	if aspectRatio < (3.0 / 4.0) {
		newHeight := width * 4.0 / 3.0
		src = imaging.CropCenter(src, width, int(newHeight))
	} else if aspectRatio > (3.0 / 4.0) {
		newWidth := height * 3.0 / 4.0
		src = imaging.CropCenter(src, int(newWidth), height)
	}

	// 使用邻居插值方法创建缩略图
	thumbnail := imaging.Resize(src, thumbWidth, thumbHeight, imaging.NearestNeighbor)

	// 保存缩略图到目标文件
	err = imaging.Save(thumbnail, destFilePath)
	if err != nil {
		return false, fmt.Errorf("failed to save thumbnail image: %v", err)
	}

	return true, nil
}

// 聊天的缩略图，不剪裁，但是按照等比例缩放，比例是算出来的，不超过指定宽度和高度
func ScaleWithAspectRatio(filePath, destFilePath string, maxWidth int) (bool, error) {
	// 打开图片文件
	src, err := imaging.Open(filePath)
	if err != nil {
		return false, fmt.Errorf("failed to open image: %v", err)
	}

	thumbWidth := maxWidth
	thumbHeight := thumbWidth * 4.0 / 3.0

	width := src.Bounds().Dx()
	height := src.Bounds().Dy()
	if width < thumbWidth || height < thumbHeight {
		return false, nil
	}

	// 检查原始图像的宽高比是否为 3:4，如果不是，则裁剪为 3:4 的比例，取中间部分
	aspectRatio := float64(width) / float64(height)
	if aspectRatio < (3.0 / 4.0) {
		// 细长的图片，以高为准
		thumbWidth = int(float64(thumbHeight) * aspectRatio)
	} else if aspectRatio > (3.0 / 4.0) {
		// 宽的图片，以宽为准
		thumbHeight = int(float64(thumbWidth) / aspectRatio)
	}

	// 使用邻居插值方法创建缩略图
	thumbnail := imaging.Resize(src, thumbWidth, thumbHeight, imaging.NearestNeighbor)

	// 保存缩略图到目标文件
	err = imaging.Save(thumbnail, destFilePath)
	if err != nil {
		return false, fmt.Errorf("failed to save thumbnail image: %v", err)
	}

	return true, nil
}
