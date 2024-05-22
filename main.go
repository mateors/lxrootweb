package main

import (
	"database/sql"
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

	//lxql.OpenConn()
	//lxql.ParseDSN("username:password@tcp(localhost:8309)/lxrootdb")
	//n1ql2, err := sql.Open("n1ql", "lxrtestusr:Test54321$@(172.93.55.179:8309)/lxrootdb")
	n1ql2, err := sql.Open("n1ql", "http://lxrtestusr:Test54321$@172.93.55.179:8093")
	fmt.Println(err)

	//ac := []byte(`[{"user": "admin:lxrtestusr", "pass": "Test54321$"}]`)
	//lxql.SetQueryParams("creds", string(ac))

	err = n1ql2.Ping()
	fmt.Println("ping..", err)

	/*
		rows, err := n1ql2.Query("select id,name,age from lxroot;")
		if err != nil {
			return
		}
		defer rows.Close()

		for rows.Next() {

			var id, name, age string
			if err := rows.Scan(&id, &name, &age); err != nil {
				log.Fatal(err)
			}
			log.Printf("Row returned -> %s,%s,%s : \n", id, name, age)
		}
	*/
	os.Exit(1)

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
