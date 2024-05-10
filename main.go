package main

import (
	"fmt"
	"lxrootweb/utility"
	"net/http"
	"os"
	"path/filepath"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

var workingDirectory string

func init() {

	//hmacTest()
	workingDirectory, _ = os.Getwd()
	//fmt.Println("Bismillah...")
}

func main() {

	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Use(middleware.RealIP)
	r.Use(middleware.Recoverer)

	assetPath := filepath.Join(workingDirectory, "assets")
	r.Handle("/resources/*", http.StripPrefix("/resources/", http.FileServer(http.Dir(assetPath))))
	//fmt.Println("Allahuakbar", utility.WPORT)

	//r.HandleFunc("/", home)
	//r.HandleFunc("/webhook", webhookHandler)

	addr := fmt.Sprintf(":%s", utility.WPORT)
	err := http.ListenAndServe(addr, r)
	fmt.Println(err)
}
