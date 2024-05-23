package main

import (
	"fmt"
	_ "lxrootweb/lxql"
	"lxrootweb/utility"
	"net/http"
	"os"
	"path/filepath"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

var workingDirectory string

func init() {

	workingDirectory, _ = os.Getwd()
	//couchbaseConnTest()
	//lxqlCon()

}

func main() {

	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Use(middleware.RealIP)
	r.Use(middleware.Recoverer)

	assetPath := filepath.Join(workingDirectory, "assets")
	r.Handle("/resources/*", http.StripPrefix("/resources/", http.FileServer(http.Dir(assetPath))))
	//fmt.Println("Allahuakbar", utility.WPORT)

	r.HandleFunc("/", homePage)
	r.HandleFunc("/support", support)
	r.HandleFunc("/features/{slug}", features)
	r.HandleFunc("/technology", technology)
	r.HandleFunc("/apphosting", apphosting)
	r.HandleFunc("/roadmap", roadmap)
	r.HandleFunc("/pricing", pricing)
	r.HandleFunc("/terms", terms)
	r.HandleFunc("/privacy", privacy)
	r.HandleFunc("/shop", shop)         //store page
	r.HandleFunc("/product", product)   //product_details
	r.HandleFunc("/checkout", checkout) //checkout + signup
	r.HandleFunc("/signup", signup)     //signup + checkout
	r.HandleFunc("/faqs", faqs)         //
	r.HandleFunc("/about", about)       //
	r.HandleFunc("/contact", contact)   //contact us
	r.HandleFunc("/join-waitlist", joinWaitlist)

	//r.HandleFunc("/webhook", webhookHandler)

	addr := fmt.Sprintf(":%s", utility.WPORT)
	err := http.ListenAndServe(addr, r)
	fmt.Println(err)
}
