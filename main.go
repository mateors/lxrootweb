package main

import (
	"database/sql"
	"fmt"
	"log"
	"lxrootweb/database"
	_ "lxrootweb/lxcb"
	"lxrootweb/lxql"
	"lxrootweb/utility"
	"net/http"
	"os"
	"path/filepath"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

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

var workingDirectory string
var err error
var COMPANY_ID string

func init() {

	workingDirectory, _ = os.Getwd()

	if DRIVER_NAME == "n1ql" {

		lxql.BUCKET = BUCKET_NAME
		lxql.SCOPE = SCOPE_NAME
		lxql.RegisterModel(Company{})
		lxql.RegisterModel(WaitList{})
		lxql.RegisterModel(Settings{})
		lxql.RegisterModel(Country{})
		lxql.RegisterModel(Access{})
		lxql.RegisterModel(Account{})
		lxql.RegisterModel(Address{})
		lxql.RegisterModel(Login{})
		lxql.RegisterModel(Verification{})
		lxql.RegisterModel(VisitorSession{})
		lxql.RegisterModel(LoginSession{})
		lxql.RegisterModel(ActivityLog{})
		lxql.RegisterModel(Authc{})
	}

	dataSourceName := fmt.Sprintf("http://%s:%s@%s:%s", DBUSER, DBPASS, SERVERIP, DBPORT)

	database.DB, err = sql.Open(DRIVER_NAME, dataSourceName)
	if err != nil {
		log.Fatal(err)
	}
	err = database.DB.Ping()
	if err != nil {
		log.Fatal(err)
	}
	log.Println("db ping successfull")
	COMPANY_ID = companyId("lxroot.com") //company id need to be inserted before proceed

	//err = SendEmail([]string{"billahmdmostain@gmail.com"}, "WELCOME", "Welcome to LxRoot")
	//fmt.Println(err)

	//fmt.Println(settingsValue("emailuser"))
	// addSettings("emailuser", "info@lxroot.com", "sysemail")
	// addSettings("emailpass", "test4321", "sysemail")
	// addSettings("emailserver", "mail.lxroot.com", "sysemail")
	// addSettings("emailport", "587", "sysemail")

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

	//name := structName(Settings{})
	//name2 := structName(&WaitList{})
	//fmt.Println(customTableName(structName(Settings{})))
	//fmt.Println(COMPANY_ID)

	// sql := "SELECT * FROM lxroot._default.country"
	// rows, err := lxql.GetRows(sql, db)
	// if err != nil {
	// 	return
	// }
	// fmt.Println("len:", len(rows))
	// for i, row := range rows {
	// 	fmt.Println(i, row)
	// }
	// stime := time.Now()
	// countryImportFromExcel("data/country.xlsx")
	// timeTaken := time.Since(stime).Seconds()
	// fmt.Println("timeTaken:", timeTaken, "s")

	// addAccess("superadmin")
	// addAccess("admin")
	// addAccess("client")
	// addAccess("partner")
	//slc := GetColumnNamesFromQuery("SELECT id,status FROM lxroot._default.access WHERE access_name='client';")
	//id := accessIdByName("client")
	//fmt.Println(id)
	// sql := fmt.Sprintf("SELECT id,name,iso_code,status FROM %s LIMIT 10;", tableToBucket("country"))
	// rows, err := lxql.GetRows(sql, db)
	// if err != nil {
	// 	log.Println(err)
	// 	return
	// }
	// for i, row := range rows {
	// 	fmt.Println(i, row, len(row))
	// }
	// query := "SELECT * FROM lxroot._default.access WHERE access_name='admin';"
	// row, err := singleRow(query)
	// fmt.Println(err, row)
	// email := ""
	// count := lxql.CheckCount("login", fmt.Sprintf(`username="%s"`, email), db)
	// fmt.Println(count)
	// var name, location, verifyUrl string = "SANZIDA YASMIN", "Rangpur, Bangladesh", "https://lxroot.com/verify?email=bill.rassel@gmail.com&token=112121212"
	// err = signupEmail("mostain@lxroot.com", name, location, verifyUrl)
	// fmt.Println(err)
	//code := uuid.NewV4().String()
	//fmt.Println(code)

	// email := "bill.rassel@gmail.com"
	// sql := fmt.Sprintf("SELECT verification_code FROM %s WHERE username='%s' AND verification_purpose='signup';", tableToBucket("verification"), email)
	// row, err := singleRow(sql)
	// if err != nil {
	// 	return
	// }
	// hexCode := row["verification_code"].(string) //
	// fmt.Println(hexCode)
	//_, err = database.DB.Exec("UPDATE ? SET status=1 WHERE id=?", tableToBucket("account"), "cpcaq4i2r9emrfgpk940") //not work
	//sql := fmt.Sprintf("UPDATE %s SET status=0 WHERE id=%q", tableToBucket("account"), "cpcaq4i2r9emrfgpk940")
	//_, err = database.DB.Exec(sql) //this works
	// fmt.Println(sql)
	// err = lxql.RawSQL(sql, database.DB)
	//fmt.Println(err)
	//row, err := usernameToAccounInfo("billahmdmostain@gmail.com")
	//fmt.Println(err, row)
	//aid := accessIdByName("client")
	//fmt.Println(aid)

	//sql := "SELECT a.cid,a.account_type,a.id as account_id FROM lxroot._default.login l LEFT JOIN lxroot._default.account a ON a.id=l.account_id WHERE l.username='billahmdmostain@gmail.com'"
	//sql := "SELECT a.cid FROM lxroot._default.login l LEFT JOIN lxroot._default.account a ON a.id=l.account_id WHERE l.username='billahmdmostain@gmail.com'"
	//sql := "SELECT * FROM lxroot._default.access;"
	//sql := "SELECT a.id as account_id, l.id as login_id FROM lxroot._default.login l LEFT JOIN lxroot._default.account a ON a.id=l.account_id WHERE l.username='billahmdmostain@gmail.com';"
	// sql := "SELECT id,cid,account_id,access_name,username,passw,label,tfa_status,tfa_medium,tfa_setupkey FROM lxroot._default.login WHERE username='billahmdmostain@gmail.com' AND status IN[1,6];"
	// rows, err := lxql.GetRows(sql, database.DB)
	// if err != nil {
	// 	return
	// }
	// for i, row := range rows {
	// 	fmt.Println(i, row)
	// }
	//getLocationWithin(ipAddress string)
	//txt := getLocationWithinSec("103.124.226.98")
	//fmt.Println(txt)
	//txt, err = getLocation("2602:ff16:4:0:1:127:0:1")
	//fmt.Println(err, txt)
	//fmt.Println(IsIPv4("103.124.226.98"), IsIPv6("103.124.226.98"))
	//fmt.Println(IsIPv4("2602:ff16:4:0:1:127:0:1"), IsIPv6("2602:ff16:4:0:1:127:0:1"))

	//fmt.Println(mtool.HashBcrypt("test321")) //$2a$14$LlXWMQxVBhW91WuJqjbCbuO5craaprMyM9tNYiOZeGjJ0mCy4Uiz2

	//err = addSettings("login_email", "-", "template")
	//err = addSettings("resetpass_email", "-", "template")
	//err = send("bill.rassel@gmail.com", "hello", "hi there")
	//fmt.Println(err)
	//os.Exit(1)
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
	r.HandleFunc("/shop", shop)                     //store page
	r.HandleFunc("/product", product)               //product_details
	r.HandleFunc("/checkout", checkout)             //checkout + signup
	r.HandleFunc("/faqs", faqs)                     //
	r.HandleFunc("/about", about)                   //
	r.HandleFunc("/contact", contact)               //contact us
	r.HandleFunc("/join-waitlist", joinWaitlist)    //WaitList
	r.HandleFunc("/signup", signup)                 //signup + checkout
	r.HandleFunc("/verify", verify)                 //verify
	r.HandleFunc("/signin", signin)                 //login
	r.HandleFunc("/resetpass", resetpass)           //reset
	r.HandleFunc("/reset-pass-form", resetPassForm) //reset-pass

	//r.HandleFunc("/webhook", webhookHandler)

	addr := fmt.Sprintf(":%s", utility.WPORT)
	err := http.ListenAndServe(addr, r)
	fmt.Println(err)
}
