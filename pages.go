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

func homePage(w http.ResponseWriter, r *http.Request) {

	if r.Method == http.MethodGet {

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
			Title        string
			Base         string
			BodyClass    string
			MainDivClass string
		}{
			Title:        "LxRoot Website",
			Base:         base,
			BodyClass:    "",
			MainDivClass: "main min-h-[calc(100vh-312px)]",
		}

		err = tmplt.Execute(w, data)
		if err != nil {
			log.Println(err)
		}
	}
}

func support(w http.ResponseWriter, r *http.Request) {

	if r.Method == http.MethodGet {

		tmplt, err := template.New("base.gohtml").Funcs(nil).ParseFiles(
			"templates/base.gohtml",
			"templates/header.gohtml",
			"templates/footer.gohtml",
			"wpages/support.gohtml", //
		)
		if err != nil {
			log.Println(err)
			return
		}

		base := GetBaseURL(r)
		data := struct {
			Title        string
			Base         string
			BodyClass    string
			MainDivClass string
		}{
			Title:        "Support | LxRoot",
			Base:         base,
			BodyClass:    "",
			MainDivClass: "main min-h-[calc(100vh-312px)]",
		}

		err = tmplt.Execute(w, data)
		if err != nil {
			log.Println(err)
		}
	}
}

func features(w http.ResponseWriter, r *http.Request) {

	if r.Method == http.MethodGet {

		tmplt, err := template.New("base.gohtml").Funcs(nil).ParseFiles(
			"templates/base.gohtml",
			"templates/header.gohtml",
			"templates/footer.gohtml",
			"wpages/features.gohtml", //
		)
		if err != nil {
			log.Println(err)
			return
		}

		base := GetBaseURL(r)
		data := struct {
			Title        string
			Base         string
			BodyClass    string
			MainDivClass string
		}{
			Title:        "Features | LxRoot",
			Base:         base,
			BodyClass:    "",
			MainDivClass: "main min-h-[calc(100vh-312px)]",
		}

		err = tmplt.Execute(w, data)
		if err != nil {
			log.Println(err)
		}
	}
}

func technology(w http.ResponseWriter, r *http.Request) {

	if r.Method == http.MethodGet {

		tmplt, err := template.New("base.gohtml").Funcs(nil).ParseFiles(
			"templates/base.gohtml",
			"templates/header.gohtml",
			"templates/footer.gohtml",
			"wpages/technology.gohtml", //
		)
		if err != nil {
			log.Println(err)
			return
		}

		base := GetBaseURL(r)
		data := struct {
			Title        string
			Base         string
			BodyClass    string
			MainDivClass string
		}{
			Title:        "Technology | LxRoot",
			Base:         base,
			BodyClass:    "",
			MainDivClass: "main min-h-[calc(100vh-312px)]",
		}

		err = tmplt.Execute(w, data)
		if err != nil {
			log.Println(err)
		}
	}
}

func apphosting(w http.ResponseWriter, r *http.Request) {

	if r.Method == http.MethodGet {

		tmplt, err := template.New("base.gohtml").Funcs(nil).ParseFiles(
			"templates/base.gohtml",
			"templates/header.gohtml",
			"templates/footer.gohtml",
			"wpages/apphosting.gohtml", //
		)
		if err != nil {
			log.Println(err)
			return
		}

		base := GetBaseURL(r)
		data := struct {
			Title        string
			Base         string
			BodyClass    string
			MainDivClass string
		}{
			Title:        "Application Hosting | LxRoot",
			Base:         base,
			BodyClass:    "",
			MainDivClass: "main min-h-[calc(100vh-312px)]",
		}

		err = tmplt.Execute(w, data)
		if err != nil {
			log.Println(err)
		}
	}
}

func roadmap(w http.ResponseWriter, r *http.Request) {

	if r.Method == http.MethodGet {

		tmplt, err := template.New("base.gohtml").Funcs(nil).ParseFiles(
			"templates/base.gohtml",
			"templates/header.gohtml",
			"templates/footer.gohtml",
			"wpages/roadmap.gohtml", //
		)
		if err != nil {
			log.Println(err)
			return
		}

		base := GetBaseURL(r)
		data := struct {
			Title        string
			Base         string
			BodyClass    string
			MainDivClass string
		}{
			Title:        "Product Roadmap | LxRoot",
			Base:         base,
			BodyClass:    "",
			MainDivClass: "main min-h-[calc(100vh-312px)]",
		}

		err = tmplt.Execute(w, data)
		if err != nil {
			log.Println(err)
		}
	}
}

func pricing(w http.ResponseWriter, r *http.Request) {

	if r.Method == http.MethodGet {

		tmplt, err := template.New("base.gohtml").Funcs(nil).ParseFiles(
			"templates/base.gohtml",
			"templates/header.gohtml",
			"templates/footer.gohtml",
			"wpages/pricing.gohtml", //
		)
		if err != nil {
			log.Println(err)
			return
		}

		base := GetBaseURL(r)
		data := struct {
			Title        string
			Base         string
			BodyClass    string
			MainDivClass string
		}{
			Title:        "Pricing | LxRoot",
			Base:         base,
			BodyClass:    "",
			MainDivClass: "main min-h-[calc(100vh-312px)]",
		}

		err = tmplt.Execute(w, data)
		if err != nil {
			log.Println(err)
		}
	}
}

func terms(w http.ResponseWriter, r *http.Request) {

	if r.Method == http.MethodGet {

		tmplt, err := template.New("base.gohtml").Funcs(nil).ParseFiles(
			"templates/base.gohtml",
			"templates/header.gohtml",
			"templates/footer.gohtml",
			"wpages/terms.gohtml", //
		)
		if err != nil {
			log.Println(err)
			return
		}

		base := GetBaseURL(r)
		data := struct {
			Title        string
			Base         string
			BodyClass    string
			MainDivClass string
		}{
			Title:        "Terms & Conditions | LxRoot",
			Base:         base,
			BodyClass:    "",
			MainDivClass: "main min-h-[calc(100vh-312px)]",
		}

		err = tmplt.Execute(w, data)
		if err != nil {
			log.Println(err)
		}
	}
}

func privacy(w http.ResponseWriter, r *http.Request) {

	if r.Method == http.MethodGet {

		tmplt, err := template.New("base.gohtml").Funcs(nil).ParseFiles(
			"templates/base.gohtml",
			"templates/header.gohtml",
			"templates/footer.gohtml",
			"wpages/privacy.gohtml", //
		)
		if err != nil {
			log.Println(err)
			return
		}

		base := GetBaseURL(r)
		data := struct {
			Title        string
			Base         string
			BodyClass    string
			MainDivClass string
		}{
			Title:        "Privacy Policy | LxRoot",
			Base:         base,
			BodyClass:    "",
			MainDivClass: "main min-h-[calc(100vh-312px)]",
		}

		err = tmplt.Execute(w, data)
		if err != nil {
			log.Println(err)
		}
	}
}

func shop(w http.ResponseWriter, r *http.Request) {

	if r.Method == http.MethodGet {

		tmplt, err := template.New("base.gohtml").Funcs(nil).ParseFiles(
			"templates/base.gohtml",
			"templates/header2.gohtml",
			"templates/footer2.gohtml",
			"wpages/shop.gohtml", //
		)
		if err != nil {
			log.Println(err)
			return
		}

		base := GetBaseURL(r)
		data := struct {
			Title        string
			Base         string
			BodyClass    string
			MainDivClass string
		}{
			Title:        "Shop | LxRoot",
			Base:         base,
			BodyClass:    "bg-white text-slate-700",
			MainDivClass: "main min-h-[calc(100vh-52px)]",
		}

		err = tmplt.Execute(w, data)
		if err != nil {
			log.Println(err)
		}
	}
}

func product(w http.ResponseWriter, r *http.Request) {

	if r.Method == http.MethodGet {

		tmplt, err := template.New("base.gohtml").Funcs(nil).ParseFiles(
			"templates/base.gohtml",
			"templates/header2.gohtml",
			"templates/footer2.gohtml",
			"wpages/product.gohtml", //
		)
		if err != nil {
			log.Println(err)
			return
		}

		base := GetBaseURL(r)
		data := struct {
			Title        string
			Base         string
			BodyClass    string
			MainDivClass string
		}{
			Title:        "Product Details | LxRoot",
			Base:         base,
			BodyClass:    "bg-white text-slate-700",
			MainDivClass: "main min-h-[calc(100vh-52px)]",
		}

		err = tmplt.Execute(w, data)
		if err != nil {
			log.Println(err)
		}
	}
}

func checkout(w http.ResponseWriter, r *http.Request) {

	if r.Method == http.MethodGet {

		tmplt, err := template.New("base.gohtml").Funcs(nil).ParseFiles(
			"templates/base.gohtml",
			"templates/header2.gohtml",
			"templates/footer2.gohtml",
			"wpages/checkout.gohtml", //
		)
		if err != nil {
			log.Println(err)
			return
		}

		base := GetBaseURL(r)
		data := struct {
			Title        string
			Base         string
			BodyClass    string
			MainDivClass string
		}{
			Title:        "Checkout | LxRoot",
			Base:         base,
			BodyClass:    "bg-white text-slate-700",
			MainDivClass: "main min-h-[calc(100vh-52px)]",
		}

		err = tmplt.Execute(w, data)
		if err != nil {
			log.Println(err)
		}
	}
}

func signup(w http.ResponseWriter, r *http.Request) {

	if r.Method == http.MethodGet {

		tmplt, err := template.New("base.gohtml").Funcs(nil).ParseFiles(
			"templates/base.gohtml",
			"templates/header2.gohtml",
			"templates/footer2.gohtml",
			"wpages/signup.gohtml", //
		)
		if err != nil {
			log.Println(err)
			return
		}

		base := GetBaseURL(r)
		data := struct {
			Title        string
			Base         string
			BodyClass    string
			MainDivClass string
		}{
			Title:        "Authenticate | LxRoot",
			Base:         base,
			BodyClass:    "",
			MainDivClass: "main min-h-[calc(100vh-52px)]",
		}

		err = tmplt.Execute(w, data)
		if err != nil {
			log.Println(err)
		}
	}
}
