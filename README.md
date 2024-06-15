# LxRoot New Website

> go clean -cache

> go clean -modcache


```
SELECT OBJECT_NAMES(b) FROM lxroot._default.company b;

SELECT OBJECT_NAMES(b)as fields FROM lxroot._default.company b;

SELECT DISTINCT OBJECT_NAMES(b)as fields FROM lxroot._default.company b;

CREATE INDEX `idx_obj_name_db` ON `db`((object_names(db)))  USING GSI;
```


### System info query
* `SELECT * FROM system:datastores`
* `SELECT * FROM system:namespaces`
* `SELECT * FROM system:buckets`
* `SELECT * FROM system:scopes`
* `SELECT * FROM system:keyspaces` | list collection
* `SELECT * FROM system:indexes`


### How do i add Collection?
Add collection manually in Community Edition 7.6.1 build 3200 (as REST API request only allowed in enterprise edition)

## Bucket Export
> cbexport json -c couchbase://127.0.0.1 -u Administrator -p Mostain321$  -b lxroot -o data.json -f lines -t 4 --scope-field scope --collection-field collection

## Bucket Import
> cbimport json -c couchbase://127.0.0.1 -u Administrator -p Mostain321$ -b lxerp -d file://data.json -f lines -g %id% --scope-collection-exp %scope%.%collection%


```
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
```
### Reference
* https://docs.couchbase.com/server/current/tools/cbexport-json.html
* https://docs.couchbase.com/server/current/tools/cbimport-json.html
* https://docs.couchbase.com/server/current/n1ql/n1ql-intro/sysinfo.html