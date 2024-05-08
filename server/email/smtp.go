package email

import (
	"bytes"
	"crypto/tls"
	"encoding/base64"
	"errors"
	"fmt"
	"math/rand"
	"mime"
	qp "mime/quotedprintable"
	"net/smtp"
	"strings"
)

type MailValidator struct {
	SMTPAddr              string
	SMTPPort              string
	SMTPHeloHost          string
	UserName              string
	TLSInsecureSkipVerify bool
	auth                  smtp.Auth
	conn                  *smtp.Client
}

type LoginAuth struct {
	username, password []byte
}

func (a *LoginAuth) Start(server *smtp.ServerInfo) (string, []byte, error) {
	return "LOGIN", []byte(a.username), nil
}

// Next continues the authentication. Exported only to satisfy the interface definition.
func (a *LoginAuth) Next(fromServer []byte, more bool) ([]byte, error) {
	if more {
		switch strings.ToLower(string(fromServer)) {
		case "username:":
			return a.username, nil
		case "password:":
			return a.password, nil
		default:
			return nil, fmt.Errorf("LOGIN AUTH unknown server response '%s'", string(fromServer))
		}
	}
	return nil, nil
}

func NewMailValidator(smtpAddr, smtpPort, smtpHeloHost string, tlsInsecureSkipVerify bool,
	userName, passWord string) *MailValidator {

	auth := &LoginAuth{[]byte(userName), []byte(passWord)}
	return &MailValidator{
		SMTPAddr:              smtpAddr,
		SMTPPort:              smtpPort,
		SMTPHeloHost:          smtpHeloHost,
		UserName:              userName,
		TLSInsecureSkipVerify: tlsInsecureSkipVerify,
		auth:                  auth,
	}
}

func (v *MailValidator) SendMail(rcpt []string, subject, msg string) error {
	if v.conn == nil {
		err := v.Connect()
		if err != nil {
			return err
		}
	}

	if err := v.conn.Noop(); err != nil {
		err = v.Connect()
		if err != nil {
			return err
		}
	}

	// 重新封装数据
	msgData, _ := v.GetMessage(v.UserName, rcpt[0], subject, "plain", msg)

	return v.doSendMail(rcpt, msgData)
}

func (v *MailValidator) Close() {
	if nil != v.conn {
		v.conn.Close()
		v.conn = nil
	}
}

func (v *MailValidator) Connect() error {
	client, err := smtp.Dial(v.SMTPAddr + ":" + v.SMTPPort)
	if err != nil {
		v.conn = nil
		return err
	}

	if err = client.Hello(v.SMTPHeloHost); err != nil {
		v.conn = nil
		return err
	}

	useTls, _ := client.Extension("STARTTLS")
	if useTls {
		tlsConfig := &tls.Config{
			InsecureSkipVerify: v.TLSInsecureSkipVerify,
			ServerName:         v.SMTPAddr,
		}

		if err = client.StartTLS(tlsConfig); err != nil {
			v.conn = nil
			return err
		}
	}

	if v.auth != nil {
		if isAuth, _ := client.Extension("AUTH"); isAuth {
			err = client.Auth(v.auth)
			if err != nil {
				v.conn = nil
				return err
			}
		}
	}

	v.conn = client
	return nil
}

func (v *MailValidator) doSendMail(rcpt []string, msg []byte) error {
	if v.conn == nil {
		return errors.New("smtp client is nil")
	}
	client := v.conn
	if err := client.Mail(v.UserName); err != nil {
		return err
	}

	for _, to := range rcpt {
		if err := client.Rcpt(to); err != nil {
			return err
		}
	}

	w, err := client.Data()
	if err != nil {
		return err
	}

	if _, err = w.Write(msg); err != nil {
		return err
	}

	if err := w.Close(); err != nil {
		return err
	}

	return nil
}

func (v *MailValidator) randomBoundary() string {
	var buf [24]byte
	rand.Read(buf[:])
	return fmt.Sprintf("tinode--%x", buf[:])
}

func (v *MailValidator) GetMessage(SendFrom, to, title, dataType, content string) ([]byte, error) {
	message := &bytes.Buffer{}

	// Common headers.
	fmt.Fprintf(message, "From: %s\r\n", SendFrom)
	fmt.Fprintf(message, "To: %s\r\n", to)
	message.WriteString("Subject: ")
	// Old email clients may barf on UTF-8 strings.
	// Encode as quoted printable with 75-char strings separated by spaces, split by spaces, reassemble.
	message.WriteString(strings.Join(strings.Split(mime.QEncoding.Encode("utf-8", title), " "), "\r\n    "))
	message.WriteString("\r\n")
	message.WriteString("MIME-version: 1.0;\r\n")

	if dataType == "plain" {
		// Plain text message
		message.WriteString("Content-Type: text/plain; charset=\"UTF-8\"; format=flowed; delsp=yes\r\n")
		message.WriteString("Content-Transfer-Encoding: base64\r\n\r\n")
		b64w := base64.NewEncoder(base64.StdEncoding, message)
		b64w.Write([]byte(content))
		b64w.Close()
	} else if dataType == "html" {
		// HTML-formatted message
		message.WriteString("Content-Type: text/html; charset=\"UTF-8\"\r\n")
		message.WriteString("Content-Transfer-Encoding: quoted-printable\r\n\r\n")
		qpw := qp.NewWriter(message)
		qpw.Write([]byte(content))
		qpw.Close()
	} else {
		// Multipart-alternative message includes both HTML and plain text components.
		boundary := v.randomBoundary()
		message.WriteString("Content-Type: multipart/alternative; boundary=\"" + boundary + "\"\r\n\r\n")

		message.WriteString("--" + boundary + "\r\n")
		message.WriteString("Content-Type: text/plain; charset=\"UTF-8\"; format=flowed; delsp=yes\r\n")
		message.WriteString("Content-Transfer-Encoding: base64\r\n\r\n")
		b64w := base64.NewEncoder(base64.StdEncoding, message)
		b64w.Write([]byte(content))
		b64w.Close()

		message.WriteString("\r\n")

		message.WriteString("--" + boundary + "\r\n")
		message.WriteString("Content-Type: text/html; charset=\"UTF-8\"\r\n")
		message.WriteString("Content-Transfer-Encoding: quoted-printable\r\n\r\n")
		qpw := qp.NewWriter(message)
		qpw.Write([]byte(content))
		qpw.Close()

		message.WriteString("\r\n--" + boundary + "--")
	}

	message.WriteString("\r\n")
	return message.Bytes(), nil
}
