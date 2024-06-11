package main

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"html/template"
	"lxrootweb/utility"
	"mime/multipart"
	"net/http"
	"net/smtp"
	"os"
	"path/filepath"
	"strings"
)

var auth smtp.Auth

func send(toEmail, subject, body string) error {

	from := "info@lxroot.com" //youthictbd@gmail.com
	pass := "test4321"        //gmfwwdjyfqtusprj
	toEmails := []string{toEmail}
	toHeader := strings.Join(toEmails, ",")
	msg := "From: " + "LxRoot<info@lxroot.com>" + "\n" +
		"To: " + toHeader + "\n" +
		"Subject: " + subject + "\n\n" + body

	return smtp.SendMail("mail.lxroot.com:587", smtp.PlainAuth("", from, pass, "mail.lxroot.com"), from, toEmails, []byte(msg))
}

func SendEmail(toEmails []string, subject, body string) error {

	emailserver := settingsValue("emailserver")
	emailport := settingsValue("emailport")
	emailuser := settingsValue("emailuser")
	emailpass := settingsValue("emailpass")
	from := emailuser
	pass := emailpass
	toHeader := strings.Join(toEmails, ",")

	msg := "From: LxRoot<" + from + ">\n" +
		"To: " + toHeader + "\n" +
		"Subject: " + subject + "\n\n" + body

	err := smtp.SendMail(fmt.Sprintf("%s:%s", emailserver, emailport), smtp.PlainAuth("", from, pass, emailserver), from, toEmails, []byte(msg))
	if err != nil {
		return err
	}
	return nil
}

func signupEmailTemplateParser(name, location, verifyUrl string) (string, error) {

	filename := "templates/email.gohtml"
	templateName := filepath.Base(filename)
	//dateTime := time.Now().Format("January 2, 2006 at 3:04 PM")

	var tplOutput bytes.Buffer
	tpl, err := template.New(templateName).ParseFiles(filename)
	if err != nil {
		return "", err
	}

	tracker := "https://lxroot.com/resources/email/open.gif?u=30345700zz&id=fd127fffaa2f407d9aa80a6c1b77a964zzz"

	data := struct {
		Name       string
		Location   string
		VerifyURL  string
		TrackerUrl string
	}{
		Name:       name,
		Location:   location,
		VerifyURL:  verifyUrl,
		TrackerUrl: tracker,
	}
	err = tpl.Execute(&tplOutput, data)
	if err != nil {
		return "", err
	}
	markupText := tplOutput.String()
	return markupText, nil
}

func signupEmail(email, name, location, verifyUrl string) error {

	markup, err := signupEmailTemplateParser(name, location, verifyUrl)
	if err != nil {
		return err
	}
	subject := "🎉 Welcome to LxRoot - action required"
	err = htmlEmailer([]string{email}, subject, markup)
	return err
}

// EMAIL CONFIG
func htmlEmailer(toEmails []string, subject, body string) error {

	toHeader := strings.Join(toEmails, ",") //;
	from := utility.EMAILUSER
	mime := "MIME-version: 1.0;\nContent-Type: text/html; charset=\"UTF-8\";\n\n"
	//body := "<html><body><h1>Hello Mostain! can you read my image?</h1><img src='https://test.youthict.org/data/file_store/content/32.png'></body></html>"
	msg := "FROM: " + from + "\n" +
		"TO: " + toHeader + "\n" +
		"SUBJECT: " + subject + "\n" +
		mime + body
	//auth = smtp.PlainAuth("", "youthictbd@gmail.com", "gmfwwdjyfqtusprj", "smtp.gmail.com")
	auth = smtp.PlainAuth("", utility.EMAILUSER, utility.EMAILPASS, utility.EMAILSERVER)
	err = smtp.SendMail(fmt.Sprintf("%s:%s", utility.EMAILSERVER, utility.EMAILPORT), auth, from, toEmails, []byte(msg))
	return err
}

// var (
// 	host       = os.Getenv("EMAIL_HOST")
// 	username   = os.Getenv("EMAIL_USERNAME")
// 	password   = os.Getenv("EMAIL_PASSWORD")
// 	portNumber = os.Getenv("EMAIL_PORT")
// )

type Sender struct {
	auth smtp.Auth
}

type Message struct {
	To          []string
	CC          []string
	BCC         []string
	Subject     string
	Body        string
	Attachments map[string][]byte
}

func New() *Sender {
	auth := smtp.PlainAuth("", utility.EMAILUSER, utility.EMAILPASS, utility.EMAILSERVER)
	return &Sender{auth}
}

func (s *Sender) Send(m *Message) error {
	return smtp.SendMail(fmt.Sprintf("%s:%s", utility.EMAILSERVER, utility.EMAILPORT), s.auth, utility.EMAILUSER, m.To, m.ToBytes())
}

func NewMessage(s, b string) *Message {
	return &Message{Subject: s, Body: b, Attachments: make(map[string][]byte)}
}

func (m *Message) AttachFile(src string) error {
	b, err := os.ReadFile(src)
	if err != nil {
		return err
	}

	_, fileName := filepath.Split(src)
	m.Attachments[fileName] = b
	return nil
}

func (m *Message) ToBytes() []byte {

	buf := bytes.NewBuffer(nil)
	withAttachments := len(m.Attachments) > 0
	buf.WriteString(fmt.Sprintf("Subject: %s\r\n", m.Subject))
	buf.WriteString(fmt.Sprintf("To: %s\r\n", strings.Join(m.To, ",")))
	if len(m.CC) > 0 {
		buf.WriteString(fmt.Sprintf("Cc: %s\r\n", strings.Join(m.CC, ",")))
	}

	if len(m.BCC) > 0 {
		buf.WriteString(fmt.Sprintf("Bcc: %s\r\n", strings.Join(m.BCC, ",")))
	}

	buf.WriteString("MIME-Version: 1.0\n")
	writer := multipart.NewWriter(buf)
	boundary := writer.Boundary()
	if withAttachments {
		buf.WriteString(fmt.Sprintf("Content-Type: multipart/mixed; boundary=%s\r\n", boundary))
		buf.WriteString(fmt.Sprintf("--%s\n", boundary))
	} else {
		buf.WriteString("Content-Type: text/plain; charset=utf-8\r\n")
	}

	buf.WriteString(m.Body)
	if withAttachments {
		for k, v := range m.Attachments {
			buf.WriteString(fmt.Sprintf("\n\n--%s\n", boundary))
			buf.WriteString(fmt.Sprintf("Content-Type: %s\r\n", http.DetectContentType(v)))
			buf.WriteString("Content-Transfer-Encoding: base64\r\n")
			buf.WriteString(fmt.Sprintf("Content-Disposition: attachment; filename=%s\r\n", k))

			b := make([]byte, base64.StdEncoding.EncodedLen(len(v)))
			base64.StdEncoding.Encode(b, v)
			buf.Write(b)
			buf.WriteString(fmt.Sprintf("\n--%s", boundary))
		}

		buf.WriteString("--")
	}

	return buf.Bytes()
}

func emailWithAttchment() {

	sender := New()
	m := NewMessage("Invoice", "check the attached invoice")
	m.To = []string{"billahmdmostain@gmail.com"}
	//m.CC = []string{"bill.rassel@gmail.com"}
	//m.BCC = []string{"bc@gmail.com"}
	m.AttachFile("data/export.json")
	fmt.Println(sender.Send(m))
}
