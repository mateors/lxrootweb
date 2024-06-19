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

	//dataClean()
	// subscriptionStart, subscriptionEnd := subscriptionStartEnd()
	// licenseKey := uuid.NewV1().String()
	// id, err := addSubscription("accountId", "billahmdmostain@gmail.com", "stripeCustomer", licenseKey, "monthly", "20", "due", subscriptionStart, subscriptionEnd, "trial")
	// fmt.Println(err, id)
	// os.Exit(1)
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
	r.HandleFunc("/signup", signup)                 //signup + checkout
	r.HandleFunc("/verify", verify)                 //verify
	r.HandleFunc("/signin", signin)                 //login ***
	r.HandleFunc("/tfauth", tfAuth)                 //tfAuth ***
	r.HandleFunc("/resetpass", resetpass)           //reset
	r.HandleFunc("/reset-pass-form", resetPassForm) //reset-pass
	r.HandleFunc("/dashboard", dashboard)           //dashboard OK
	r.HandleFunc("/profile", profile)               //profile **
	r.HandleFunc("/security", security)             //security **
	r.HandleFunc("/ticket", tickets)                //ticket **
	r.HandleFunc("/ticket/{tid}", ticketDetails)    //ticket/details
	r.HandleFunc("/orders", orders)                 //Billing > My orders OK
	r.HandleFunc("/paymethods", payMethods)         //Billing > My orders OK
	r.HandleFunc("/activity", activityLog)          //Profile > Activity
	r.HandleFunc("/orders/{oid}", orderDetails)     //orderDetails OK
	r.HandleFunc("/invoices", invoices)             //invoice ** OK
	r.HandleFunc("/license", licenseKey)            //license OK
	r.HandleFunc("/ticketnew", ticketNew)           //ticketnew
	r.HandleFunc("/logout", logout)                 //logout
	r.HandleFunc("/invoice/{order}", invoice)       //
	r.HandleFunc("/qrcode/{id}", qrcodeimg)         //

	r.HandleFunc("/payhook", paymentHook)
	//r.HandleFunc("/join-waitlist", joinWaitlist)    //WaitList
	//r.HandleFunc("/webhook", webhookHandler)

	addr := fmt.Sprintf(":%s", utility.WPORT)
	err := http.ListenAndServe(addr, compress(r))
	fmt.Println(err)
}
