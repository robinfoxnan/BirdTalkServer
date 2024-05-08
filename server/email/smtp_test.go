package email

import (
	"fmt"
	"testing"
)

func TestSmtp(t *testing.T) {

	generator, err := NewEmailGenerator(`D:\GBuild\BirdTalkServer\server\emailtemp\email-validation-zh.templ`)

	// 创建 SMTP 客户端实例
	smtpClient := NewMailValidator("smtp.sina.com", "25", "sina.com", false,
		"tinodemaster@sina.com", "7ee60df2d8b390fd")

	data := EmailData{
		HostUrl: "http://birdtalk.com",
		Code:    "12345",
		Session: "123333333333",
		Server:  "1",
	}

	subject, txt, err := generator.GeneratePlainEmail(&data)
	//fmt.Println(subject, txt)

	err = smtpClient.SendMail([]string{"robin-fox@sohu.com"}, subject, txt)
	fmt.Println(err)

	data.Code = "45678"
	subject, txt, err = generator.GeneratePlainEmail(&data)
	err = smtpClient.SendMail([]string{"robin-fox@sohu.com"}, subject, txt)
	fmt.Println(err)
}

func TestSmtpText(t *testing.T) {

	data := EmailData{
		HostUrl: "http://birdtalk.com",
		Code:    "12345",
		Session: "123333333333",
		Server:  "1",
	}

	generator, err := NewEmailGenerator(`D:\GBuild\BirdTalkServer\server\emailtemp\email-validation-zh.templ`)
	fmt.Println(err)
	subject, txt, err := generator.GeneratePlainEmail(&data)
	fmt.Println(subject, txt)
	//txt, err := GeneratePlainEmail(template, &data)
	//fmt.Println(txt, err)
}
