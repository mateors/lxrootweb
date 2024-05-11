package main

import (
	"fmt"
	"log"
	"net/http"
	"text/template"
)

func GetBaseURL(r *http.Request) string {

	var baseurl, proto string
	fproto := r.Header.Get("X-Forwarded-Proto")
	proto = "http"
	if fproto == "https" {
		proto = "https"
	} else if r.TLS != nil {
		proto = "https"
	}
	baseurl = fmt.Sprintf("%s://%s", proto, r.Host)
	return baseurl
}

func home(w http.ResponseWriter, r *http.Request) {

	tmplt, err := template.New("base.gohtml").Funcs(nil).ParseFiles(
		"templates/base.gohtml",
		"templates/header.gohtml",
		"templates/footer.gohtml",
		"wpages/home.gohtml", //
	)
	if err != nil {
		log.Println(err)
		return
	}

	base := GetBaseURL(r)
	data := struct {
		Title string
		Base  string
	}{
		Title: "LxRoot Website",
		Base:  base,
	}

	err = tmplt.Execute(w, data)
	if err != nil {
		log.Println(err)
	}

}
