package main

import (
	"database/sql"
	"fmt"
	"log"
	_ "lxrootweb/lxcb"
	"lxrootweb/lxql"
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
	DRIVER_NAME = "n1ql"
	ENCDECPASS  = "MosT$sLxRoot"
)

var db *sql.DB
var err error
var COMPANY_ID string

func init() {

	workingDirectory, _ = os.Getwd()
	//couchbaseConnTest()
	//lxqlCon()
	if DRIVER_NAME == "n1ql" {

		lxql.BUCKET = BUCKET_NAME
		lxql.SCOPE = SCOPE_NAME
		lxql.RegisterModel(Company{})
		lxql.RegisterModel(WaitList{})
	}

	dataSourceName := fmt.Sprintf("http://%s:%s@%s:%s", DBUSER, DBPASS, SERVERIP, DBPORT)
	db, err = sql.Open(DRIVER_NAME, dataSourceName)
	if err != nil {
		log.Fatal(err)
	}
	err = db.Ping()
	if err != nil {
		log.Fatal(err)
	}
	log.Println("db ping successfull")
	COMPANY_ID = companyId("lxroot.com") //company id need to be inserted before proceed

	//tableName := customTableName("AccountTable")
	//err = createCollection("company", db)
	//err = addCompany("MATEORS DOT COM LLC")
	//count := CheckCount("company", "type='company'", db)

	//DELETE FROM lxroot._default.company USE KEYS ["cp8346a2r9eu4jj9mhjg","cp8858a2r9et68vde730"]
	//_, err = db.Exec("DELETE FROM lxroot._default.company WHERE id='cp8ba4i2r9eu5orkqda0';")
	//fmt.Println(err)

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
	//cols, err := ReadTable2Columns("Company", db)
	//fmt.Println(err, cols)
	//fmt.Println(customTableName("WaitList"))

	//name := structName(Company{})
	//name2 := structName(&WaitList{})
	//fmt.Println(name, name2)
	//fmt.Println(COMPANY_ID)

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
