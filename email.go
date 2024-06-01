package main

import (
	"bytes"
	"fmt"
	"html/template"
	"lxrootweb/utility"
	"net/smtp"
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
