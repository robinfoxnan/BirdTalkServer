package email

import (
	"html/template"
	"strings"
)

type EmailGenerator struct {
	tmpl *template.Template
}

type EmailData struct {
	HostUrl string
	Code    string
	Session string
	Server  string
	Email   string
}

func NewEmailGenerator(filepath string) (*EmailGenerator, error) {
	tmpl, err := template.ParseFiles(filepath)
	if err != nil {
		return nil, err
	}
	g := &EmailGenerator{
		tmpl: tmpl,
	}

	return g, nil
}

func (g *EmailGenerator) GeneratePlainEmail(data *EmailData) (string, string, error) {
	// 解析模板
	var plainBody strings.Builder
	err := g.tmpl.ExecuteTemplate(&plainBody, "body_plain", data)

	var subject strings.Builder
	g.tmpl.ExecuteTemplate(&subject, "subject", data)
	return subject.String(), plainBody.String(), err
}
