package main

import (
	"database/sql"
	"fmt"
	"log"
	_ "lxrootweb/lxql"
	"lxrootweb/utility"
	"net/http"
	"os"
	"path/filepath"
	"reflect"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

var workingDirectory string
var typeRegistry = make(map[string]reflect.Type)

const (
	SERVERIP    = "172.93.55.179" //
	DBUSER      = "lxrtestusr"
	DBPASS      = "Test54321$" //
	DBPORT      = "8093"
	BUCKET_NAME = "lxroot"
	SCOPE_NAME  = "_default"
)

var db *sql.DB
var err error

func init() {

	workingDirectory, _ = os.Getwd()
	//couchbaseConnTest()
	//lxqlCon()
	registerType(Company{})

	dataSourceName := fmt.Sprintf("http://%s:%s@%s:%s", DBUSER, DBPASS, SERVERIP, DBPORT)
	db, err = sql.Open("n1ql", dataSourceName)
	if err != nil {
		log.Fatal(err)
	}
	err = db.Ping()
	if err != nil {
		log.Fatal(err)
	}
	log.Println("db ping successfull")

	//tableName := customTableName("AccountTable")
	//err = createCollection("company", db)
	//err = addCompany("LXROOT LLC")
	//count := CheckCount("company", "type='company'", db)
	//fmt.Println(count)

	// sql := fmt.Sprintf("SELECT id,company_name,serial,status FROM %s;", tableToBucket("company"))
	// rows, err := GetRows(sql, db)
	// if err != nil {
	// 	log.Println(err)
	// 	return
	// }
	// for _, row := range rows {
	// 	fmt.Println(row)
	// }

	//slc := strStructToFields("Test")
	//fmt.Println(slc)

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
