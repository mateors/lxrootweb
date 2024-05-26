package main

import (
	"bytes"
	"fmt"
	"html/template"
	"lxrootweb/utility"
	"net/smtp"
	"path/filepath"
	"strings"
	"time"
)

var auth smtp.Auth

func send(toEmail, subject, body string) error {

	from := "info@lxroot.com" //youthictbd@gmail.com
	pass := "test321"         //gmfwwdjyfqtusprj
	toEmails := []string{toEmail}
	toHeader := strings.Join(toEmails, ",")
	msg := "From: " + from + "\n" +
		"To: " + toHeader + "\n" +
		"Subject: " + subject + "\n\n" + body

	return smtp.SendMail("mx.lxroot.com:587", smtp.PlainAuth("", from, pass, "mx.lxroot.com"), from, toEmails, []byte(msg))
}

// func send(body string) {

// 	from := "no-reply@bag-n-brand.com"
// 	pass := "Te$st123!"
// 	//to := "ash@complimention.com"
// 	toEmails := []string{"bill.rassel@gmail.com", "ash@complimention.com"}

// 	msg := "From: " + from + "\n" +
// 		"To: " + toEmails[1] + "\n" +
// 		"Subject: Hello there\n\n" +
// 		body

// 		//smtp.zoho.com
// 		//smtp.gmail.com:587
// 	err := smtp.SendMail("smtp.zoho.com:587",
// 		smtp.PlainAuth("", from, pass, "smtp.zoho.com"),
// 		from, toEmails, []byte(msg))

// 	if err != nil {
// 		log.Printf("smtp error: %s", err)
// 		return
// 	}

// 	//log.Print("sent, visit http://foobarbazz.mailinator.com")
// }

func SendEmail(toEmails []string, subject, body string) error {

	emailserver := settingsValue("emailserver")
	emailport := settingsValue("emailport")
	emailuser := settingsValue("emailuser")
	emailpass := settingsValue("emailpass")
	from := emailuser
	pass := emailpass
	toHeader := strings.Join(toEmails, ",")

	msg := "From: " + from + "\n" +
		"To: " + toHeader + "\n" +
		"Subject: " + subject + "\n\n" + body

	err := smtp.SendMail(fmt.Sprintf("%s:%s", emailserver, emailport), smtp.PlainAuth("", from, pass, emailserver), from, toEmails, []byte(msg))
	if err != nil {
		return err
	}
	return nil
}

func studentEmailTemplateParser(base, name, regno, courseName, courseMode, session, username, password string) (string, error) {

	filename := "template/email/student.gohtml"
	templateName := filepath.Base(filename)

	dateTime := time.Now().Format("January 2, 2006 at 3:04 PM")
	var tplOutput bytes.Buffer

	tpl, err := template.New(templateName).ParseFiles(filename)
	if err != nil {
		return "", err
	}

	data := struct {
		Base       string
		DateTime   string
		Name       string
		RegNo      string
		CourseName string
		CourseMode string //online|offline
		Session    string
		Username   string
		Password   string
	}{
		Base:       base,
		DateTime:   dateTime,
		Name:       name,
		RegNo:      regno,
		CourseName: courseName,
		CourseMode: courseMode,
		Session:    session,
		Username:   username,
		Password:   password,
	}
	err = tpl.Execute(&tplOutput, data)
	if err != nil {
		return "", err
	}
	markupText := tplOutput.String()
	return markupText, nil
}

func branchEmailTemplateParser(base, name, director, code, upazila, district, username, password string) (string, error) {

	filename := "template/email/branch.gohtml"
	templateName := filepath.Base(filename)

	dateTime := time.Now().Format("January 2, 2006 at 3:04 PM")
	var tplOutput bytes.Buffer

	tpl, err := template.New(templateName).ParseFiles(filename)
	if err != nil {
		return "", err
	}

	data := struct {
		Base     string
		DateTime string
		Name     string //Institue name
		Director string //owner name
		Code     string //branch code
		Upazila  string //thana
		District string
		Username string
		Password string
	}{
		Base:     base,
		DateTime: dateTime,
		Name:     name,
		Director: director,
		Code:     code,
		Upazila:  upazila,
		District: district,
		Username: username,
		Password: password,
	}
	err = tpl.Execute(&tplOutput, data)
	if err != nil {
		return "", err
	}
	markupText := tplOutput.String()
	return markupText, nil
}

func branchSignupEmail(base, name, director, code, district, upazila, username, password string) error {

	// base := ""
	// name := "High Tech Training Center"
	// director := "Kamrul Islam"
	// code := "1559"
	// upazila := "Araihazar"
	// district := "Narayanganj"
	// username := "username@gmail.com"
	// password := generatePassword(12, true, true)
	markup, err := branchEmailTemplateParser(base, name, director, code, upazila, district, username, password)
	if err != nil {
		return err
	}
	toEmails := []string{username}
	err = htmlEmailer(toEmails, markup)
	return err
}

func studentSignupEmail(base, name, regno, courseName, courseMode, session, username, password string) error {

	// base := ""
	// name := "MOSTAIN BILLAH"
	// regno := "300001"
	// courseName := "Golang Programming"
	// courseMode := "Online"
	// session := "January 2023"
	// username := "mateors"
	//password := generatePassword(12, true, true)
	markup, err := studentEmailTemplateParser(base, name, regno, courseName, courseMode, session, username, password)
	if err != nil {
		return err
	}
	//toEmails := []string{"billahmdmostain@gmail.com"}
	err = htmlEmailer([]string{username}, markup)
	return err
}

// EMAIL CONFI
func htmlEmailer(toEmails []string, body string) error {

	//toEmails := []string{"billahmdmostain@gmail.com"} //"nasarulhasan@gmail.com"
	toHeader := strings.Join(toEmails, ";")
	//ccEmail := []string{"admin@mateors.com"}
	//ccHeader := strings.Join(ccEmail, ";")
	//bccEmail := []string{"bill.rassel@gmail.com"}
	//bccHeader := strings.Join(bccEmail, ";")
	subject := "REGISTRATION EMAIL"
	from := "youthictbd@gmail.com"
	mime := "MIME-version: 1.0;\nContent-Type: text/html; charset=\"UTF-8\";\n\n"
	//body := "<html><body><h1>Hello Mostain! can you read my image?</h1><img src='https://test.youthict.org/data/file_store/content/32.png'></body></html>"
	msg := "FROM: " + from + "\n" +
		"TO: " + toHeader + "\n" +
		//"CC: " + ccHeader + "\n" +
		//"BCC: " + bccHeader + "\n" +
		"SUBJECT: " + subject + "\n" +
		mime + body
	//auth = smtp.PlainAuth("", "youthictbd@gmail.com", "gmfwwdjyfqtusprj", "smtp.gmail.com")
	auth = smtp.PlainAuth("", utility.EMAILUSER, utility.EMAILPASS, utility.EMAILSERVER)
	err = smtp.SendMail(fmt.Sprintf("%s:%s", utility.EMAILSERVER, utility.EMAILPORT), auth, from, toEmails, []byte(msg))
	return err
}

// func htmlEmail() {

// 	toEmails := []string{"billahmdmostain@gmail.com"} //"nasarulhasan@gmail.com"
// 	toHeader := strings.Join(toEmails, ";")
// 	ccEmail := []string{"admin@mateors.com"}
// 	ccHeader := strings.Join(ccEmail, ";")
// 	bccEmail := []string{"bill.rassel@gmail.com"}
// 	bccHeader := strings.Join(bccEmail, ";")
// 	subject := "REGISTRATION EMAIL"
// 	//toHeader:="To: " + toHeader + "\n" +
// 	//subject := "Subject: Email from Go!\n"
// 	from := "youthictbd@gmail.com"
// 	mime := "MIME-version: 1.0;\nContent-Type: text/html; charset=\"UTF-8\";\n\n"
// 	body := "<html><body><h1>Hello Mostain! can you read my image?</h1><img src='https://test.youthict.org/data/file_store/content/32.png'></body></html>"
// 	msg := "FROM: " + from + "\n" +
// 		"TO: " + toHeader + "\n" +
// 		"CC: " + ccHeader + "\n" +
// 		"BCC: " + bccHeader + "\n" +
// 		"SUBJECT: " + subject + "\n" +
// 		mime + body

// 	//msg := []byte(subject + mime + body)
// 	auth = smtp.PlainAuth("", "youthictbd@gmail.com", "gmfwwdjyfqtusprj", "smtp.gmail.com")
// 	err = smtp.SendMail("smtp.gmail.com:587", auth, from, toEmails, []byte(msg))

// }

// func init() {
// 	fmt.Println("loading...email.go")
// }

// Request struct
// type Request struct {
// 	from    string
// 	to      []string
// 	subject string
// 	body    string
// }

// func NewRequest(to []string, subject, body string) *Request {
// 	return &Request{
// 		to:      to,
// 		subject: subject,
// 		body:    body,
// 	}
// }

// func templateEmail() {

// 	auth = smtp.PlainAuth("", "youthictbd@gmail.com", "gmfwwdjyfqtusprj", "smtp.gmail.com")
// 	templateData := struct {
// 		Name string
// 		URL  string
// 	}{
// 		Name: "YOUTHICT",
// 		URL:  "http://youthict.org",
// 	}
// 	r := NewRequest([]string{"billahmdmostain@gmail.com"}, "Hello Mostain!", "Assalamualikum, hope you are fine with the grace of almighty.")
// 	err = r.ParseTemplate("template.html", templateData)
// 	fmt.Println(err)
// 	if err = r.ParseTemplate("template.html", templateData); err == nil {
// 		ok, _ := r.SendEmail()
// 		fmt.Println(ok)
// 	}
// }

// func (r *Request) SendEmail() (bool, error) {
// 	mime := "MIME-version: 1.0;\nContent-Type: text/plain; charset=\"UTF-8\";\n\n"
// 	subject := "Subject: " + r.subject + "!\n"
// 	msg := []byte(subject + mime + "\n" + r.body)
// 	addr := "smtp.gmail.com:587"

// 	if err := smtp.SendMail(addr, auth, "billahmdmostain@gmail.com", r.to, msg); err != nil {
// 		return false, err
// 	}
// 	return true, nil
// }

// func (r *Request) ParseTemplate(templateFileName string, data interface{}) error {
// 	t, err := template.ParseFiles(templateFileName)
// 	if err != nil {
// 		return err
// 	}
// 	buf := new(bytes.Buffer)
// 	if err = t.Execute(buf, data); err != nil {
// 		return err
// 	}
// 	r.body = buf.String()
// 	return nil
// }
