package main

import (
	"database/sql"
	"fmt"
	"log"
	"lxrootweb/database"
	"lxrootweb/utility"
	"net/http"
	"os"
	"path/filepath"

	_ "github.com/mateors/lxcb"
	"github.com/mateors/lxql"

	"github.com/CAFxX/httpcompression"
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
		lxql.RegisterModel(DocKeeper{})
		lxql.RegisterModel(TransactionRecord{})
		lxql.RegisterModel(LedgerTransaction{})
		lxql.RegisterModel(DocPayShipInfo{})
		lxql.RegisterModel(Subscription{})
		lxql.RegisterModel(Item{})
		lxql.RegisterModel(Event{})          //StripeEvent
		lxql.RegisterModel(Ticket{})         //
		lxql.RegisterModel(TicketResponse{}) //
		lxql.RegisterModel(Department{})
		lxql.RegisterModel(FileStore{})
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
	//deleteAccount("billahmdmostain@gmail.com")
	//loginId := usernameToLoginId("bill.rassel@gmail.com")
	//fmt.Println(loginId)
	//count := lxql.CheckCount("login", fmt.Sprintf(`username="%s"`, "bill.rassel@gmail.com"), database.DB)
	//fmt.Println(count)
	//aesEncDecTest()

	//id, err := addItem("lxroot", "subscription", "lxroot-license-monthly", "20 cents per app", "", "", "20", "")
	//fmt.Println(err, id)
	//requestBalance("sk_test_51OjqyFJFUQv2NTJsitgDUhNX3CPbns3eE3IyxSdTc8yEhI5p24SDyn9lyEI4AqaMSRghw6V25XoStkYa8Zl7zEOg006vuF1cTQ")
	//listAllPrices(secretKey string) (map[string]interface{}, error)
	//row, err := listAllPrices("sk_test_51OjqyFJFUQv2NTJsitgDUhNX3CPbns3eE3IyxSdTc8yEhI5p24SDyn9lyEI4AqaMSRghw6V25XoStkYa8Zl7zEOg006vuF1cTQs")
	// rurl := "https://api.stripe.com/v1/prices"
	//stripeKey := "sk_test_51OjqyFJFUQv2NTJsitgDUhNX3CPbns3eE3IyxSdTc8yEhI5p24SDyn9lyEI4AqaMSRghw6V25XoStkYa8Zl7zEOg006vuF1cTQ"
	// row, err := apiGetRequest(rurl, stripeKey)
	// if err != nil {
	// 	log.Fatal(err)
	// }
	// var fmap = make(map[string]string)
	// fmap["currency"] = "usd"
	// fmap["unit_amount"] = "2000"
	// fmap["recurring[interval]"] = "month"
	// fmap["product_data[name]"] = "LxRoot License Monthly" //product.name
	// fmap["nickname"] = "LxRoot License Monthly"
	// row, err := apiPostRequest("https://api.stripe.com/v1/prices", stripeKey, fmap)
	// row, err := createSession(stripeKey, "doc_200", "alauddin@mateors.com", "price_1PPlulJFUQv2NTJsqGsPFpLa", "1")
	// if err != nil {
	// 	log.Fatal(err)
	// }

	//get url and redirect to that url for stripe checkout
	//status
	//id

	// for key, val := range row {
	// 	fmt.Printf("%v = %v, %T\n", key, val, val)
	// }

	// fmt.Println(unixToDateTime(1720542106))
	// fmt.Println(strToTime("1720542106").Format("2006-01-02 15:04:05"))
	//fmt.Println(strToTime("1720542106").Format("2006-01-02 15:04:05"))

	//var evt Event
	//row := make(map[string]interface{})
	//err = json.Unmarshal([]byte(jsonTxt), &evt)
	//fmt.Println("##", err, evt.Created, evt.Object, evt.ID)
	// fmt.Printf("%v %T\n", row["created"], row["created"])
	// fmt.Printf("%v\n", row["object"])
	// fmt.Printf("%v %T\n", row["data"], row["data"])
	//fmt.Println(">", toTime(row["created"]).Format("2006-01-02 15:04:05"))
	//fmt.Println("<", toTime(int64(1717990849)))

	// ns := struct {
	// 	Name string
	// 	Age  int
	// 	Lang []string
	// 	Info map[string]interface{}
	// }{
	// 	Name: "Mostain",
	// 	Age:  40,
	// 	Lang: []string{"go", "rust", "python"},
	// 	Info: row,
	// }
	// row, err := structFieldValMap(&ns)
	// fmt.Println(err, row)
	//err = lxql.InsertUpdateObject("event", evt.ID, &evt, database.DB)
	//fmt.Println(err)

	// rurl := "https://pay.stripe.com/receipts/invoices/CAcaFwoVYWNjdF8xT2pxeUZKRlVRdjJOVEpzKLWOm7MGMgZ9ouzNg6s6LBbPoyN2PS0meir7wTTJAHXYzGPo0iTp-0u4GdFOsmFoLJi7IAVevpykoFy7?s=ap"
	// durl, err := stripeReceiptToPdfUrl(rurl)
	// fmt.Println(err, durl)

	// filename, err := DownloadFile("data/invoice", durl)
	// fmt.Println(err, filename)
	//docCheckoutProcess("cpjrjlq2r9et0ao3vuqg")

	//fmt.Println(uuid.NewV1())
	//fmt.Println(uuid.NewV4())
	// addDepartment("General", "GEN", "", "ticket")
	// addDepartment("Billing", "BIL", "", "ticket")
	// addDepartment("Sales", "SAL", "", "ticket")
	// addDepartment("Technical", "TEC", "", "ticket")
	// addDepartment("Bugs", "BUG", "", "ticket")
	// doc, err := emailToDocNumber("billahmdmostain@gmail.com")
	// fmt.Println(err, doc)

	// sql := `SELECT d.doc_number FROM lxroot._default.login l LEFT JOIN lxroot._default.doc_keeper d ON d.login_id=l.id WHERE d.doc_status='complete' AND l.username="billahmdmostain@gmail.com" ORDER BY d.id DESC LIMIT 1;`
	// row, err := singleRow(sql)

	//iurl := stripeInvoiceReceiptUrl("in_1PR7ZOJFUQv2NTJsHQ1dHZS0")
	//inumber := stripeInvoiceToNumber("in_1PR7ZOJFUQv2NTJsHQ1dHZS0")
	//start, end := subscriptionStartEnd()
	//fmt.Println(start, end)

	//rurl := invoiceToReceiptUrl("C90FC4CC-0001")
	//fmt.Println(rurl)

	//utility.LicenseEncDec()
	//aesEncDecTest()
	//key, err := genKey()
	//fmt.Println(err, key)
	//row, err := getInterfaces()
	//fmt.Println(err, row)
	//dataClean()
	//os.Exit(1)
}

func main() {

	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Use(middleware.RealIP)
	r.Use(middleware.Recoverer)

	assetPath := filepath.Join(workingDirectory, "assets")
	r.Handle("/resources/*", http.StripPrefix("/resources/", http.FileServer(http.Dir(assetPath))))
	//r.Handle("/vdata/*", http.StripPrefix("/vdata/", http.FileServer(http.Dir(filepath.Join(workingDirectory, "data")))))
	//fmt.Println("Allahuakbar", utility.WPORT)

	compress, _ := httpcompression.DefaultAdapter()

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
	r.HandleFunc("/complete", complete)             //complete shopping
	r.HandleFunc("/getstarted", product)            //product_details
	r.HandleFunc("/checkout", checkout)             //checkout + ***
	r.HandleFunc("/faqs", faqs)                     //
	r.HandleFunc("/about", about)                   //
	r.HandleFunc("/contact", contact)               //contact us
	r.HandleFunc("/join-waitlist", joinWaitlist)    //WaitList
	r.HandleFunc("/signup", signup)                 //signup + checkout
	r.HandleFunc("/verify", verify)                 //verify
	r.HandleFunc("/signin", signin)                 //login ***
	r.HandleFunc("/resetpass", resetpass)           //reset
	r.HandleFunc("/reset-pass-form", resetPassForm) //reset-pass
	r.HandleFunc("/dashboard", dashboard)           //dashboard OK
	r.HandleFunc("/profile", profile)               //profile **
	r.HandleFunc("/security", security)             //security **
	r.HandleFunc("/ticket", tickets)                //ticket **
	r.HandleFunc("/ticket/{tid}", ticketDetails)    //ticket/details
	r.HandleFunc("/orders", orders)                 //Billing > My orders OK
	r.HandleFunc("/orders/{oid}", orderDetails)     //orderDetails OK
	r.HandleFunc("/invoices", invoices)             //invoice ** OK
	r.HandleFunc("/license", licenseKey)            //license OK
	r.HandleFunc("/ticketnew", ticketNew)           //ticketnew
	r.HandleFunc("/logout", logout)                 //logout
	r.HandleFunc("/invoice/{order}", invoice)       //

	r.HandleFunc("/payhook", paymentHook)
	//r.HandleFunc("/webhook", webhookHandler)

	addr := fmt.Sprintf(":%s", utility.WPORT)
	err := http.ListenAndServe(addr, compress(r))
	fmt.Println(err)
}
