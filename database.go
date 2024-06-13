package main

import (
	"database/sql"
	"errors"
	"fmt"
	"log"
	"lxrootweb/database"
	"lxrootweb/lxql"
	"lxrootweb/utility"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/mateors/mtool"
	"github.com/rs/xid"
)

var (
	COMPANY_TABLE = "Company"
)

// func checkCollectionExist(collectionName string) bool {

// 	qs := `SELECT aid,count,name FROM %s USE KEYS["%s"];`
// 	sql := fmt.Sprintf(qs, tableToBucket("tableRowCount"), collectionName)
// 	pres := db.Query(sql)
// 	rows := pres.GetRows()
// 	return len(rows) != 0
// }

// func checkFileExist(filePath string) error {
// 	var err error
// 	if _, err = os.Stat(filePath); err == nil {
// 		return nil
// 	} else if errors.Is(err, os.ErrNotExist) {
// 		return fmt.Errorf("file %s does not exist!", filePath)
// 	}
// 	return err
// }

func createCollection(collectionName string, db *sql.DB) error {

	_, err = db.Exec("CREATE COLLECTION ?;", tableToBucket(collectionName))
	//sql := fmt.Sprintf("CREATE COLLECTION %s;", tableToBucket(collectionName))
	//pres := db.Query(sql)
	// if pres.Status != "success" {
	// 	var emsg string
	// 	for _, e := range pres.Errors {
	// 		emsg += fmt.Sprintf("%s-%s,", e.Message, e.Message)
	// 	}
	// 	log.Println("ERRsql:", sql)
	// 	return fmt.Errorf("ERROR %s", strings.TrimRight(emsg, ","))
	// }
	//time.Sleep(time.Millisecond * 200)
	return err
}

// func createSecondaryIndex(collectionName string) error {

// 	sql := fmt.Sprintf("CREATE INDEX `status%s` ON %s(status);", collectionName, tableToBucket(collectionName))
// 	//fmt.Println(sql)
// 	pres := db.Query(sql)
// 	if pres.Status != "success" {
// 		log.Println("ERRsql:", sql)
// 		return fmt.Errorf("ERROR %s", singleMsg(pres))
// 	}
// 	time.Sleep(time.Millisecond * 200)
// 	return nil
// }

// func createPrimaryIndex(collectionName string) error {
// 	sql := fmt.Sprintf("CREATE PRIMARY INDEX ON %s;", tableToBucket(collectionName))
// 	//fmt.Println(sql)
// 	pres := db.Query(sql)
// 	if pres.Status != "success" {
// 		log.Println("ERRsql:", sql)
// 		return fmt.Errorf("ERROR %s", singleMsg(pres))
// 	}
// 	time.Sleep(time.Millisecond * 200)
// 	return nil
// }

// func collectionInstance(collectionName string) error {

// 	var vmap = make(map[string]interface{})
// 	vmap["aid"] = collectionName
// 	vmap["count"] = 0
// 	vmap["name"] = collectionName
// 	vmap["type"] = "tableRowCount"
// 	insertStatm := `INSERT INTO %s (KEY,VALUE) VALUES("%s",%s);`
// 	bs, err := json.Marshal(&vmap)
// 	if err != nil {
// 		return err
// 	}
// 	sql := fmt.Sprintf(insertStatm, tableToBucket("tableRowCount"), collectionName, string(bs))
// 	pres := db.Query(sql)
// 	if pres.Status != "success" {
// 		return fmt.Errorf("ERROR %s", singleMsg(pres))
// 	}
// 	time.Sleep(time.Millisecond * 200)
// 	return nil
// }

// func nextSerial(tableName string) int {
// 	return lxql.CheckCount(tableName, fmt.Sprintf("type='%s'", tableName), db) + 1
// }

// func GetColumnNamesFromQuery(query string) []string {

// 	query = strings.ToLower(query)
// 	// Remove leading/trailing whitespaces and semicolon (if any)
// 	query = strings.TrimSpace(query)
// 	if strings.HasSuffix(query, ";") {
// 		query = query[:len(query)-1]
// 	}

// 	// Split the query by spaces
// 	parts := strings.Fields(query)

// 	// Find the index of "SELECT" and "FROM" keywords
// 	selectIndex := -1
// 	fromIndex := -1
// 	for i, part := range parts {
// 		if strings.EqualFold(part, "select") {
// 			selectIndex = i
// 		} else if strings.EqualFold(part, "from") {
// 			fromIndex = i
// 			break
// 		}
// 	}

// 	// Extract column names between "SELECT" and "FROM"
// 	if selectIndex != -1 && fromIndex != -1 && fromIndex > selectIndex+1 {

//			columns := parts[selectIndex+1 : fromIndex]
//			// Remove commas from column names
//			for i, col := range columns {
//				columns[i] = strings.TrimSuffix(col, ",")
//			}
//			var rcols []string
//			for _, col := range columns {
//				slc := strings.Split(col, ",")
//				rcols = slc
//			}
//			return rcols
//		}
//		return nil
//	}
func basicForm() {

	var form = make(map[string]interface{})
	form["id"] = "id::1"
	form["company_name"] = "Mostain"
	form["age"] = "40"
	form["lang"] = []string{"golang", "rust"}
	form["table"] = "Company"

	err = lxql.InsertUpdateMap(form, database.DB)
	fmt.Println(err)
}

func addCompany(companyName string) error {

	table := customTableName(COMPANY_TABLE)
	var form = make(map[string]interface{})
	id := xid.New().String()
	form["id"] = id
	form["company_name"] = companyName
	form["table"] = COMPANY_TABLE //model
	form["type"] = table
	form["status"] = 1
	err = lxql.InsertUpdateMap(form, database.DB)
	return err
}

func modelUpsert(modelName string, form url.Values) error {

	var mForm = make(map[string]interface{})
	table := customTableName(modelName) //database table
	mForm["table"] = modelName
	mForm["id"] = xid.New().String()
	mForm["type"] = table

	for key := range form {
		mForm[key] = form.Get(key)
	}
	err = lxql.InsertUpdateMap(mForm, database.DB)
	return err
}

func addSettings(fieldName, fieldValue, purpose string) error {

	modelName := structName(Settings{})
	table := customTableName(modelName)
	var form = make(map[string]interface{})
	id := xid.New().String()
	form["id"] = id
	form["type"] = table
	form["cid"] = COMPANY_ID
	form["table"] = modelName
	form["field_name"] = fieldName
	form["field_value"] = fieldValue
	form["purpose"] = purpose
	form["status"] = 1
	return lxql.InsertUpdateMap(form, database.DB)
}

func settingsValue(fieldName string) string {
	return lxql.FieldByValue("settings", "field_value", fmt.Sprintf("field_name='%s'", fieldName), database.DB)
}

func addCountry(name, isoCode, countryCode string) error {

	modelName := structName(Country{})
	table := customTableName(modelName)
	var form = make(map[string]interface{})
	id := xid.New().String()
	form["id"] = id
	form["type"] = table
	form["cid"] = COMPANY_ID
	form["table"] = modelName
	form["name"] = name
	form["iso_code"] = isoCode
	form["country_code"] = countryCode
	form["status"] = 1
	return lxql.InsertUpdateMap(form, database.DB)
}

func addAccess(accessName string) error {

	modelName := structName(Access{})
	table := customTableName(modelName)
	var form = make(map[string]interface{})
	id := xid.New().String()
	form["id"] = id
	form["type"] = table
	form["cid"] = COMPANY_ID
	form["table"] = modelName
	form["access_name"] = accessName
	form["status"] = 1
	return lxql.InsertUpdateMap(form, database.DB)
}

func addAccount(parentId, accountType, email, accountName, firstName, lastName string) (id string, err error) {

	modelName := structName(Account{})
	table := customTableName(modelName)
	var form = make(map[string]interface{})
	id = xid.New().String()
	form["id"] = id
	form["type"] = table
	form["cid"] = COMPANY_ID
	form["table"] = modelName
	form["parent_id"] = parentId
	form["account_type"] = accountType //vendor,customer
	form["account_name"] = accountName
	form["code"] = id
	form["first_name"] = firstName
	form["last_name"] = lastName
	form["email"] = email
	form["create_date"] = mtool.TimeNow()
	form["status"] = 0 //in active by default
	err = lxql.InsertUpdateMap(form, database.DB)
	return id, err
}

func addAddress(accountId, addressType, country, state, city, address1, address2, zip string) (id string, err error) {

	modelName := structName(Address{})
	table := customTableName(modelName)
	var form = make(map[string]interface{})
	id = xid.New().String()
	form["id"] = id
	form["type"] = table
	form["cid"] = COMPANY_ID
	form["table"] = modelName
	form["account_id"] = accountId
	form["address_type"] = addressType //billing
	form["country"] = country
	form["state"] = state
	form["city"] = city
	form["address1"] = address1
	form["address2"] = address2
	form["zip"] = zip
	form["status"] = 1
	err = lxql.InsertUpdateMap(form, database.DB)
	return id, err
}

func addLogin(accountId, accessId, accessName, username, plainPassword string) (id string, err error) {

	modelName := structName(Login{})
	table := customTableName(modelName)
	var form = make(map[string]interface{})
	id = xid.New().String()
	form["id"] = id
	form["type"] = table
	form["cid"] = COMPANY_ID
	form["table"] = modelName
	form["account_id"] = accountId
	form["access_id"] = accessId     //billing
	form["access_name"] = accessName //superadmin,admin,client,partner
	form["username"] = username
	form["passw"] = mtool.HashBcrypt(plainPassword)
	form["tfa_status"] = 0
	form["create_date"] = mtool.TimeNow()
	form["status"] = 0 //inactive by default
	err = lxql.InsertUpdateMap(form, database.DB)
	return id, err
}

func addVerification(username, purpose, code, messageId string) (id string, err error) {

	modelName := structName(Verification{})
	table := customTableName(modelName)
	var form = make(map[string]interface{})
	id = xid.New().String()
	form["id"] = id
	form["type"] = table
	form["cid"] = COMPANY_ID
	form["table"] = modelName
	form["username"] = username
	form["verification_purpose"] = purpose                               //signup
	form["verification_code"] = mtool.EncodeStr(code, utility.JWTSECRET) //
	form["message_id"] = messageId
	form["create_date"] = mtool.TimeNow()
	form["status"] = 0 //inactive by default
	err = lxql.InsertUpdateMap(form, database.DB)
	return id, err
}

func addActiviyLog(loginId, activityType, ownerTable, parameter, logDetails, ipAddress string) (id string, err error) {

	modelName := structName(ActivityLog{})
	table := customTableName(modelName)
	var form = make(map[string]interface{})
	id = xid.New().String()
	form["id"] = id
	form["type"] = table
	form["cid"] = COMPANY_ID
	form["table"] = modelName
	form["activity_type"] = activityType //UPDATE|INSERT|DELETE|CREATE
	form["table_name"] = ownerTable
	form["parameter"] = parameter
	form["log_details"] = logDetails
	form["ip_address"] = ipAddress
	form["login_id"] = loginId
	form["create_date"] = mtool.TimeNow()
	form["status"] = 1
	err = lxql.InsertUpdateMap(form, database.DB)
	return id, err
}

func addLoginSession(loginId, visitorSessionID, ipAddress, city, country, userAgent string) (id string, err error) {

	loginTime := mtool.TimeNow()
	modelName := structName(LoginSession{})
	table := customTableName(modelName)
	var form = make(map[string]interface{})
	id = xid.New().String()
	form["id"] = id
	form["type"] = table
	form["cid"] = COMPANY_ID
	form["table"] = modelName
	form["session_code"] = visitorSessionID
	form["login_id"] = loginId
	form["ip_address"] = ipAddress
	form["city"] = city
	form["country"] = country
	form["user_agent"] = userAgent
	form["login_time"] = loginTime
	form["create_date"] = loginTime
	form["status"] = 1
	err = lxql.InsertUpdateMap(form, database.DB)
	return id, err
}

// crow, err := TokenToClaim(token)
// 	if err != nil {
// 		return err
// 	}
// 	loginId := crow["login_id"].(string)
// 	exp := crow["exp"].(float64)
// 	expireDate := dateFormat("", exp)

func addAuthc(loginId, token, ipAddress, expireDate string) (id string, err error) {

	modelName := structName(Authc{})
	table := customTableName(modelName)
	var form = make(map[string]interface{})
	id = xid.New().String()
	form["id"] = id
	form["type"] = table
	form["cid"] = COMPANY_ID
	form["table"] = modelName
	form["login_id"] = loginId
	form["token"] = token
	form["ip_address"] = ipAddress
	form["create_date"] = mtool.TimeNow()
	form["expire_date"] = expireDate
	form["update_date"] = ""
	form["status"] = 1
	err = lxql.InsertUpdateMap(form, database.DB)
	return id, err
}

func addDocKeeper(docName, docType, docRef, docNumber, postingDate, docStatus, totalDiscount, totalTax, totalPayable, loginId, accountId, ipAddress string) (id string, err error) {

	modelName := structName(DocKeeper{})
	table := customTableName(modelName)
	var form = make(map[string]interface{})
	id = xid.New().String()
	if postingDate == "" {
		postingDate = time.Now().Format("2006-01-02")
	}
	if docNumber == "" {
		docNumber = id
	}
	form["id"] = id
	form["type"] = table
	form["cid"] = COMPANY_ID
	form["table"] = modelName
	form["doc_name"] = docName
	form["doc_type"] = docType
	form["doc_ref"] = docRef //visitorSession
	form["doc_number"] = docNumber
	form["posting_date"] = postingDate
	form["login_id"] = loginId
	form["account_id"] = accountId
	form["total_discount"] = totalDiscount
	form["total_tax"] = totalTax
	form["total_payable"] = totalPayable
	form["doc_status"] = docStatus
	form["ip_address"] = ipAddress
	form["create_date"] = mtool.TimeNow()
	form["status"] = 1
	err = lxql.InsertUpdateMap(form, database.DB)
	return id, err
}

func addDepartment(name, code, description, owner string) (id string, err error) {

	modelName := structName(Department{})
	table := customTableName(modelName)
	var form = make(map[string]interface{})
	id = xid.New().String()
	form["id"] = id
	form["type"] = table
	form["cid"] = COMPANY_ID
	form["table"] = modelName
	form["name"] = name
	form["code"] = code
	form["description"] = description
	form["owner"] = owner //ownerTable
	form["status"] = 1
	err = lxql.InsertUpdateMap(form, database.DB)
	return id, err
}

func addTicket(loginId, department, subject, message, reference, ipAddress string) (id string, err error) {

	modelName := structName(Ticket{})
	table := customTableName(modelName)
	var form = make(map[string]interface{})
	id = xid.New().String()
	if reference == "" {
		reference = id
	}
	form["id"] = id
	form["type"] = table
	form["cid"] = COMPANY_ID
	form["table"] = modelName
	form["department"] = department
	form["subject"] = subject
	form["message"] = message
	form["reference"] = reference //ticketNumber
	form["login_id"] = loginId
	form["ticket_status"] = "open"
	form["ip_address"] = ipAddress
	form["create_date"] = mtool.TimeNow()
	form["status"] = 1
	err = lxql.InsertUpdateMap(form, database.DB)
	return id, err
}

func str2int(val string) int {
	ival, err := strconv.Atoi(val)
	if err != nil {
		return 0
	}
	return ival
}

func addTransactionRecord(trxType, docNumber, itemId, itemInfo, itemSerial, stripePriceId, qty, price string) (id string, err error) {

	modelName := structName(TransactionRecord{})
	table := customTableName(modelName)
	var form = make(map[string]interface{})
	id = xid.New().String()
	form["id"] = id
	form["type"] = table
	form["cid"] = COMPANY_ID
	form["table"] = modelName
	form["trx_type"] = trxType
	form["doc_number"] = docNumber
	form["item_id"] = itemId
	form["item_info"] = itemInfo
	form["item_serial"] = itemSerial
	form["stock_info"] = stripePriceId
	form["quantity"] = qty
	form["rate"] = 1
	form["price"] = price
	form["payable_amount"] = str2int(qty) * str2int(price)
	form["create_date"] = mtool.TimeNow()
	form["status"] = 1
	err = lxql.InsertUpdateMap(form, database.DB)
	return id, err
}

// Meta Data
func addItem(itemName, itemCategory, itemCode, itemDesc, image, buyPrice, salePrice, supplier string) (id string, err error) {

	modelName := structName(Item{})
	table := customTableName(modelName)
	var form = make(map[string]interface{})
	id = xid.New().String()
	form["id"] = id
	form["type"] = table
	form["cid"] = COMPANY_ID
	form["table"] = modelName
	form["item_code"] = itemCode
	form["category_id"] = itemCategory //license,domain,hosting
	form["item_type"] = "subscription" //service -> subscription
	form["item_name"] = itemName
	form["item_description"] = itemDesc
	form["item_image"] = image
	form["buy_price"] = buyPrice
	form["sale_price"] = salePrice
	form["tags"] = "lxroot"
	form["supplier"] = supplier
	form["uom"] = "unit"
	form["tracking"] = "unique_serial"
	form["availability"] = "available"
	form["status"] = 1
	err = lxql.InsertUpdateMap(form, database.DB)
	return id, err
}

func addEvent(id string) (string, error) {

	modelName := structName(Event{})
	//table := customTableName(modelName)
	var form = make(map[string]interface{})
	if id == "" {
		id = xid.New().String()
	}
	form["id"] = id
	//form["type"] = table
	//form["cid"] = COMPANY_ID
	form["table"] = modelName
	//form["item_code"] = itemCode
	//form["category_id"] = itemCategory //license,domain,hosting
	//form["item_type"] = "subscription" //service -> subscription
	//form["status"] = 1
	err = lxql.InsertUpdateMap(form, database.DB)
	return id, err
}

func addSubscription(accountId, stripeCustomer, licenseKey, billing, price, paymentStatus, subscriptionStart, subscriptionEnd, remarks string) (id string, err error) {

	modelName := structName(Subscription{})
	table := customTableName(modelName)
	var form = make(map[string]interface{})
	id = xid.New().String()
	form["id"] = id
	form["type"] = table
	form["cid"] = COMPANY_ID
	form["table"] = modelName
	form["account_id"] = accountId
	form["subscriber"] = stripeCustomer //stripe.customer
	form["license_key"] = licenseKey    //
	//form["domain"] = domain
	form["billing"] = billing
	form["price"] = price
	form["payment_status"] = paymentStatus
	form["subscription_start"] = subscriptionStart
	form["subscription_end"] = subscriptionEnd
	form["create_date"] = mtool.TimeNow()
	form["remarks"] = remarks
	form["status"] = 1
	err = lxql.InsertUpdateMap(form, database.DB)
	return id, err
}

func addFileStore(ownerTable, reference, fileType, filepath, remarks string) (id string, err error) {

	modelName := structName(FileStore{})
	table := customTableName(modelName)
	var form = make(map[string]interface{})
	id = xid.New().String()
	form["id"] = id
	form["type"] = table
	form["cid"] = COMPANY_ID
	form["table"] = modelName
	form["owner_table"] = ownerTable //doc_keeper
	form["reference"] = reference    //doc_number
	form["file_type"] = fileType     //pdf
	form["filepath"] = filepath      //
	form["remarks"] = remarks        //stripe.invoice.id => evt_3PQ9PsJFUQv2NTJs0BEhJIyn
	form["create_date"] = mtool.TimeNow()
	form["status"] = 1
	err = lxql.InsertUpdateMap(form, database.DB)
	return id, err
}

func addTicketResponse(ticketId, message, loginId, ipAddress string) (id string, err error) {

	modelName := structName(TicketResponse{})
	table := customTableName(modelName)
	var form = make(map[string]interface{})
	id = xid.New().String()
	form["id"] = id
	form["type"] = table
	form["cid"] = COMPANY_ID
	form["table"] = modelName
	form["ticket_id"] = ticketId
	form["respond_by"] = loginId
	form["message"] = message
	form["ip_address"] = ipAddress
	form["create_date"] = mtool.TimeNow()
	form["status"] = 1
	err = lxql.InsertUpdateMap(form, database.DB)
	return id, err
}

func accessIdByName(accessName string) string {

	sql := fmt.Sprintf("SELECT id,status FROM %s WHERE access_name='%s';", tableToBucket("access"), accessName)
	row, err := singleRow(sql)
	if err != nil {
		log.Println("accessIdByName:", err, sql)
		return ""
	}
	return row["id"].(string)
}

func verifySignup(email, token string) error {

	sql := fmt.Sprintf("SELECT id,verification_code FROM %s WHERE username='%s' AND verification_purpose='signup' AND status=0;", tableToBucket("verification"), email)
	row, err := singleRow(sql)
	if err != nil {
		return errors.New("invalid link") //The provided link appears to be invalid.
	}
	hexCode := row["verification_code"].(string) //encoded with jwtsecret
	plainTxt := mtool.DecodeStr(hexCode, utility.JWTSECRET)
	if token == plainTxt {
		id := row["id"].(string)
		lxql.RawSQL(fmt.Sprintf("UPDATE %s SET status=1,update_date=%q WHERE id=%q;", tableToBucket("verification"), mtool.TimeNow(), id), database.DB)
		return nil
	}
	return errors.New("invalid token")
}

func singleRow(sql string) (map[string]interface{}, error) {

	rows, err := lxql.GetRows(sql, database.DB)
	if err != nil {
		return nil, err
	}
	for _, row := range rows {
		return row, nil
	}
	return nil, errors.New("no record found")
}

func cleanText(input string) string {
	return strings.Join(strings.Fields(input), " ")
}

func usernameToAccounInfo(username string) (map[string]interface{}, error) {

	sql := fmt.Sprintf(`SELECT 
	a.cid,
	a.account_type,
	a.id as account_id,
	l.id as login_id,
	a.first_name,
	a.last_name,
	a.code,
	a.customid,
	a.phone,
	a.email,
	a.photo,
	a.referral_url,
	a.status,
	l.access_name,
	l.last_login,
	l.passw,
	l.tfa_medium,
	l.tfa_setupkey,
	l.tfa_status 
	FROM lxroot._default.login l 
	LEFT JOIN lxroot._default.account a ON a.id=l.account_id
	WHERE l.username="%s";`, username)
	sql = cleanText(sql)
	return singleRow(sql)
}

// func bytesToStr(slc []uint8) string {

// 	var str string
// 	for _, c := range slc {
// 		if c != 34 { //remove "
// 			str += fmt.Sprintf("%c", c)
// 		}
// 	}
// 	return str
// }

// func getRows(sql string, db *sql.DB) ([]map[string]interface{}, error) {

// 	rows, err := db.Query(sql)
// 	if err != nil {
// 		return nil, err
// 	}

// 	defer rows.Close()
// 	columns, err := rows.Columns()
// 	if err != nil {
// 		return nil, err
// 	}

// 	count := len(columns)
// 	values := make([]interface{}, count)
// 	valuePtrs := make([]interface{}, count)
// 	//fmt.Println(columns, len(columns))
// 	var isStarFound bool
// 	var colCount int

// 	var orows = make([]map[string]interface{}, 0)

// 	for rows.Next() {

// 		for i := range columns {
// 			valuePtrs[i] = &values[i]
// 		}

// 		rows.Scan(valuePtrs...)
// 		var orow = make(map[string]interface{})
// 		for i, col := range columns {
// 			colCount++
// 			if col == "*" {
// 				isStarFound = true
// 			}
// 			val := values[i]
// 			orow[col] = val
// 		}
// 		orows = append(orows, orow)

// 	} //

// 	//process
// 	var nrows = make([]map[string]interface{}, 0)

// 	if isStarFound {

// 		//fmt.Println("* found...")
// 		for _, row := range orows {

// 			//fmt.Println(row["*"])
// 			for _, col := range columns {
// 				//fmt.Printf("%v = %v %T\n", col, bytesToStr(row[col].([]uint8)), row[col])
// 				var vmap = make(map[string]interface{})
// 				json.Unmarshal(row[col].([]uint8), &vmap)
// 				for key := range vmap {
// 					vrow, isOk := vmap[key].(map[string]interface{})
// 					if isOk {
// 						nrows = append(nrows, vrow)
// 					} else {
// 						fmt.Printf("%v %T\n", vmap[key], vmap[key])
// 					}
// 				}
// 			}
// 		}

// 	} else if colCount == 1 {

// 		//fmt.Println("1 col count...", columns, len(columns), columns[0])
// 		for _, row := range orows {
// 			for _, val := range row {
// 				json.Unmarshal(val.([]uint8), &row)
// 			}
// 			nrows = append(nrows, row)
// 		}

// 	} else if colCount > 1 {

// 		for _, row := range orows {
// 			var srow = make(map[string]interface{})
// 			for key, val := range row {
// 				srow[key] = bytesToStr(val.([]uint8))
// 				//fmt.Printf("%v => %v %T\n", key, val, val)
// 			}
// 			nrows = append(nrows, srow)
// 		}

// 	} else {
// 		fmt.Println("else:>", orows, len(orows))
// 	}
// 	return nrows, nil
// }

func emailToDocNumber(email string) (docNumber string, err error) {

	start, end := docFindDateTime()
	sql := fmt.Sprintf(`SELECT d.doc_number FROM %s l LEFT JOIN %s d ON d.login_id=l.id WHERE d.doc_type='cart' AND l.username="%s" AND d.create_date BETWEEN '%s' AND '%s' ORDER BY d.id DESC LIMIT 1;`, tableToBucket("login"), tableToBucket("doc_keeper"), email, start, end)
	row, err := singleRow(sql)
	if err != nil {
		log.Println(sql)
		return "", err
	}
	fmt.Println(sql)
	docNumber, _ = row["doc_number"].(string)
	return
}

func docToCartCount(docNumber string) (count int) {

	qs := `SELECT count(*)as cnt FROM %s d LEFT JOIN %s t ON t.doc_number=d.doc_number WHERE d.doc_number="%s" AND d.doc_status="pending";`
	sql := fmt.Sprintf(qs, tableToBucket("doc_keeper"), tableToBucket("transaction_record"), docNumber)
	row, err := singleRow(sql)
	if err == nil {
		cnt := fmt.Sprint(row["cnt"]) //float64
		count = str2int(cnt)
	}
	return
}

func docWiseItems(docNumber string) ([]map[string]interface{}, error) {

	sql := fmt.Sprintf("SELECT stock_info,quantity FROM %s WHERE doc_number=%q;", tableToBucket("transaction_record"), docNumber)
	return lxql.GetRows(sql, database.DB)
}

func groupItemPrice(rows []map[string]interface{}) map[string]interface{} {

	var rmap = make(map[string]interface{})
	for _, row := range rows {
		key, _ := row["stock_info"].(string) //stripePriceId
		qty, _ := row["quantity"].(string)
		val, isExist := rmap[key]
		if isExist {
			rmap[key] = val.(int) + str2int(qty)
		} else {
			rmap[key] = str2int(qty)
		}
	}
	return rmap
}

func docCheckoutProcess(docNumber string) map[string]interface{} {

	//check and make sure doc_kepper.doc_status is pending
	rows, err := docWiseItems(docNumber)
	if err != nil {
		return nil
	}
	return groupItemPrice(rows)
}

func totalOrdersByAccount(accountId string) int {
	return lxql.CheckCount("doc_keeper", fmt.Sprintf("account_id=%q AND doc_type='cart' AND doc_status='complete' AND status=1", accountId), database.DB)
}

func totalInvoicesByAccount(accountId string) int {
	return lxql.CheckCount("doc_keeper", fmt.Sprintf("account_id=%q AND doc_type='sales' AND doc_status IN ['complete','pending'] AND status=1", accountId), database.DB)
}

func totalUnpaidInvoicesByAccount(accountId string) int {
	return lxql.CheckCount("doc_keeper", fmt.Sprintf("account_id=%q AND doc_type='sales' AND doc_status='pending' AND status=1", accountId), database.DB)
}

func totalActiveTicketByUser(loginId string) int {
	return lxql.CheckCount("ticket", fmt.Sprintf("login_id=%q AND status=1", loginId), database.DB)
}

func subscriptionDetailsByAccount(accountId string) (map[string]interface{}, error) {

	sql := fmt.Sprintf("SELECT license_key,billing,price,payment_status,subscription_end,create_date FROM %s WHERE account_id=%q;", tableToBucket("subscription"), accountId)
	return singleRow(sql)
}

func myOrders(accountId string) ([]map[string]interface{}, error) {

	sql := fmt.Sprintf(`SELECT d.id, d.doc_status,d.doc_number,d.posting_date,d.receipt_url,d.payment_status,d.doc_name,d.doc_description,d.doc_ref,d.create_date,d.total_payable,t.item_info,t.price 
FROM lxroot._default.doc_keeper d 
LEFT JOIN  lxroot._default.transaction_record t ON d.doc_number=t.doc_number
WHERE d.account_id="%s" AND d.doc_type='cart' AND d.doc_status='complete' AND d.status=1;`, accountId)
	return lxql.GetRows(sql, database.DB)
}

func profileLastLogin(loginId string) (lastLoginDate, clientSince string) {

	sql := fmt.Sprintf("SELECT id,create_date,last_login FROM %s WHERE id=%q;", tableToBucket("login"), loginId)
	row, err := singleRow(sql)
	if err != nil {
		return
	}
	lastLoginDate = row["last_login"].(string)
	joinDate := row["create_date"].(string)
	clientSince = mtool.DateTimeParser(joinDate, "2006-01-02 15:04:05", "January 02, 2006")
	return
}

func profileInfo(accountId string) map[string]interface{} {

	sql := fmt.Sprintf("SELECT a.first_name,a.last_name,a.email,l.create_date,l.last_login,l.username FROM lxroot._default.account a LEFT JOIN lxroot._default.login l ON l.account_id=a.id WHERE a.id=%q;", accountId)
	//fmt.Println(sql)
	row, err := singleRow(sql)
	if err != nil {
		return nil
	}
	row["client_since"] = mtool.DateTimeParser(row["create_date"].(string), "2006-01-02 15:04:05", "January 02, 2006")
	return row
}

func addressList(accountId string) []map[string]interface{} {

	sql := fmt.Sprintf("SELECT address1,city,state,zip,country,address_type FROM lxroot._default.address WHERE account_id=%q;", accountId)
	rows, err := lxql.GetRows(sql, database.DB)
	if err != nil {
		log.Println(err)
		return nil
	}
	return rows
}

func docNumberToAccountInfo(docNumber string) map[string]interface{} {

	qs := `SELECT a.account_name,a.email,d.receipt_url,d.total_payable FROM lxroot._default.doc_keeper d LEFT JOIN lxroot._default.account a ON d.account_id=a.id WHERE d.doc_number=%q;`
	sql := fmt.Sprintf(qs, docNumber)
	row, err := singleRow(sql)
	if err != nil {
		log.Println(err)
		return nil
	}
	return row
}

func dataClean() {

	lxql.RawSQL("DELETE FROM lxroot._default.doc_keeper;", database.DB)
	lxql.RawSQL("DELETE FROM lxroot._default.transaction_record;", database.DB)
	lxql.RawSQL("DELETE FROM lxroot._default.event;", database.DB)
	lxql.RawSQL("DELETE FROM lxroot._default.subscription;", database.DB)
	lxql.RawSQL("DELETE FROM lxroot._default.file_store;", database.DB)
	lxql.RawSQL("DELETE FROM lxroot._default.ticket;", database.DB)
	lxql.RawSQL("DELETE FROM lxroot._default.ticket_response;", database.DB)
}
