package main

import (
	"database/sql"
	"errors"
	"fmt"
	"io"
	"lxrootweb/database"
	"lxrootweb/utility"
	"math"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"reflect"
	"strconv"
	"strings"
	"time"

	"github.com/mateors/lxql"

	"github.com/mateors/mtool"
	"github.com/mileusna/useragent"
	"github.com/rs/xid"
	uuid "github.com/satori/go.uuid"
)

func featureBySlug(slugName string, rows []map[string]interface{}) map[string]interface{} {

	for _, row := range rows {
		slug := row["slug"].(string)
		if slug == slugName {
			return row
		}
	}
	return nil
}

func findFeaturesById(featureId int, fRows []map[string]interface{}) []map[string]interface{} {

	var rows = make([]map[string]interface{}, 0)
	for _, row := range fRows {
		fid := row["feature_id"].(int)
		if fid == featureId {
			rows = append(rows, row)
		}
	}
	return rows
}

// Round float to 2 decimal places
func roundFloat(val float64, precision uint) float64 {
	ratio := math.Pow(10, float64(precision))
	return math.Round(val*ratio) / ratio
}

func calculateAppPrice(numApps int) float64 {

	costPerAppRanges := []struct {
		min  int
		max  int
		cost float64
	}{
		{1, 5000, 0.15},
		{5001, 25000, 0.10},
		{25001, 100000, 0.075},
	}

	var totalCost float64
	remainingApps := numApps

	for _, rangeItem := range costPerAppRanges {

		appsInRange := int(math.Min(float64(remainingApps), float64(rangeItem.max-rangeItem.min+1)))
		totalCost += float64(appsInRange) * rangeItem.cost
		remainingApps -= appsInRange

		if remainingApps <= 0 {
			break
		}
	}

	// Ensure minimum billing threshold
	//fmt.Sprintf("%.2f", number)
	totalCost = math.Max(totalCost, 15)
	return totalCost
}

// func couchbaseConnTest() {

// 	//couchbase
// 	//"user:password@/dbname"
// 	n1ql, err := sql.Open("n1ql", "http://172.93.55.179:8093")
// 	if err != nil {
// 		fmt.Println("***", err)
// 		//log.Fatal(err)
// 		return
// 	}

// 	ac := []byte(`[{"user": "admin:Administrator", "pass": "Mostain321$"}]`)
// 	go_n1ql.SetQueryParams("creds", string(ac))

// 	//go_n1ql.SetQueryParams("timeout", "10s")
// 	//go_n1ql.SetQueryParams("scan_consistency", "request_plus")

// 	// err = n1ql.Ping()
// 	// if err != nil {
// 	// 	fmt.Println("###")
// 	// 	log.Fatal(err)
// 	// }

// 	// fmt.Println("ping success...")

// 	// Set query parameters
// 	//ac := []byte(`[{"user": "admin:Administrator", "pass": "asdasd"}]`)
// 	//go_n1ql.SetQueryParams("creds", string(ac))
// 	//go_n1ql.SetQueryParams("timeout", "10s")
// 	// go_n1ql.SetQueryParams("scan_consistency", "request_plus")

// 	rows, err := n1ql.Query("select id,name,age from lxroot;")
// 	if err != nil {
// 		return
// 	}
// 	defer rows.Close()

// 	columnNames, err := rows.Columns()
// 	if err != nil {
// 		return
// 	}

// 	rc := newMapStringScan(columnNames)
// 	tableData := make([]map[string]interface{}, 0)

// 	for rows.Next() {

// 		//var id, name, age string
// 		//var age float64

// 		// if err := rows.Scan(&id, &name, &age); err != nil {
// 		// 	log.Fatal(err)
// 		// }
// 		//log.Printf("Row returned -> %s,%s,%s : \n", id, name, age)

// 		err = rc.Update(rows)
// 		if err != nil {
// 			break
// 		}
// 		cv := rc.Get()
// 		dd := make(map[string]interface{})
// 		for _, col := range columnNames {
// 			dd[col] = cv[col]
// 		}
// 		tableData = append(tableData, dd)
// 	}
// 	fmt.Printf("%v %T\n", tableData, tableData)

// 	// if err := rows.Err(); err != nil {
// 	// 	log.Fatal(err)
// 	// }

// 	// stt, err := n1ql.Prepare(`INSERT INTO lxroot (KEY, VALUE) VALUES ("doc2", {"name":"sanzida","age":32,"id":"doc2"}) RETURNING *`)
// 	// stt.Exec()
// 	// fmt.Println("insert:", err)

// 	// stt, err := n1ql.Prepare(`INSERT INTO lxroot (KEY, VALUE) VALUES ("doc3", {"name":"Wania","age":3,"id":"doc3"})`)
// 	// stt.Exec()
// 	// fmt.Println("insert:", err)

// 	// stt, err := n1ql.Prepare(`UPSERT INTO lxroot (KEY, VALUE) VALUES ("doc2", {"name":"Sanzida","age":33,"id":"doc2"})`)
// 	// stt.Exec()
// 	// fmt.Println("insert:", err)

// 	//SELECT id,name,age FROM lxroot USE KEYS "id::1"
// 	//res, err := n1ql.Exec("DELETE FROM lxroot USE KEYS ?", "id::1")
// 	//fmt.Println("delete:", res, err)

// }

func ParseDSN(dsn string) (err error) {

	// [user[:password]@][net[(addr)]]/dbname[?param1=value1&paramN=valueN]
	// Find the last '/' (since the password or the net addr might contain a '/')
	var user, passwd, addr, net, dbname string
	foundSlash := false
	for i := len(dsn) - 1; i >= 0; i-- {
		if dsn[i] == '/' {
			foundSlash = true
			var j, k int

			// left part is empty if i <= 0
			if i > 0 {
				// [username[:password]@][protocol[(address)]]
				// Find the last '@' in dsn[:i]
				for j = i; j >= 0; j-- {
					if dsn[j] == '@' {
						// username[:password]
						// Find the first ':' in dsn[:j]
						for k = 0; k < j; k++ {
							if dsn[k] == ':' {
								passwd = dsn[k+1 : j]
								break
							}
						}
						user = dsn[:k]

						break
					}
				}

				// [protocol[(address)]]
				// Find the first '(' in dsn[j+1:i]
				for k = j + 1; k < i; k++ {
					if dsn[k] == '(' {
						// dsn[i-1] must be == ')' if an address is specified
						if dsn[i-1] != ')' {
							if strings.ContainsRune(dsn[k+1:i], ')') {
								return errors.New("invalid DSN: did you forget to escape a param value")
							}
							return errors.New("invalid DSN: network address not terminated (missing closing brace)")
						}
						addr = dsn[k+1 : i-1]
						break
					}
				}
				net = dsn[j+1 : k]
			}

			// dbname[?param1=value1&...&paramN=valueN]
			// Find the first '?' in dsn[i+1:]
			for j = i + 1; j < len(dsn); j++ {
				if dsn[j] == '?' {
					//if err = parseDSNParams(cfg, dsn[j+1:]); err != nil {
					//return
					//}
					break
				}
			}

			dbname = dsn[i+1 : j]
			// if cfg.DBName, err = url.PathUnescape(dbname); err != nil {
			// 	return fmt.Errorf("invalid dbname %q: %w", dbname, err)
			// }
			break
		}
	}

	if !foundSlash && len(dsn) > 0 {
		return errors.New("invalid DSN: missing the slash separating the database name")
	}

	// if err = cfg.normalize(); err != nil {
	// 	return  err
	// }
	fmt.Println(user, passwd, addr, net, dbname)
	// Qparams["user"] = user
	// Qparams["passwd"] = passwd
	// Qparams["address"] = addr
	// Qparams["protocol"] = net
	// Qparams["dbname"] = dbname
	return
}

func lxqlCon() {

	//go_n1ql.OpenN1QLConnection()
	//lxql.OpenConn()
	//lxql.ParseDSN("username:password@tcp(localhost:8309)/lxrootdb")
	//n1ql2, err := sql.Open("n1ql", "lxrtestusr:Test54321$@(172.93.55.179:8309)/lxrootdb")
	n1ql2, err := sql.Open("n1ql", "http://lxrtestusr:Test54321$@172.93.55.179:8093")
	//n1ql2, err := sql.Open("n1ql", "http://172.93.55.179:8093")
	fmt.Println(err)

	//ac := []byte(`[{"user": "admin:lxrtestusr", "pass": "Test54321$"}]`)
	//lxql.SetQueryParams("creds", string(ac))

	err = n1ql2.Ping()
	fmt.Println("ping..", err)

	/*
		stt, err := n1ql2.Prepare("SELECT id,name,age FROM lxroot WHERE name=? AND age=?")
		if err != nil {
			log.Println(err)
			return
		}

		rows, err := stt.Query("Sanzida", 33)
		if err != nil {
			log.Fatal(err)
		}

		for rows.Next() {

			var id, name, age string
			if err := rows.Scan(&id, &name, &age); err != nil {
				log.Fatal(err)
			}
			log.Printf("-> %s,%s,%s : \n", id, name, age)
		}
	*/

	/*
		stmt, err := n1ql2.Prepare("Upsert INTO lxroot values (?,?)")
		if err != nil {
			log.Fatal(err)
		}

		// Map Values need to be marshaled
		value, _ := json.Marshal(map[string]interface{}{"name": "irish", "type": "contact"})
		result, err := stmt.Exec("irish4", value)
		if err != nil {
			log.Fatal(err)
		}

		rowsAffected, err := result.RowsAffected()
		if err != nil {
			log.Fatal(err)
		}
		log.Printf("Rows affected %d", rowsAffected)

		for i := 0; i < 20; i++ {

			key := fmt.Sprintf("irish%d", i)
			result, err = stmt.Exec(key, value)
			if err != nil {
				log.Fatal(err)
			}

			ra, err := result.RowsAffected()
			if err != nil {
				log.Fatal(err)
			}
			rowsAffected += ra
		}
		log.Printf("Total Rows Affected %d", rowsAffected)
		stmt.Close()
	*/

	// lid, err := res.LastInsertId()
	// fmt.Println("lastInsert:", lid, err)

	//res, err := n1ql2.Exec("DELETE FROM lxroot USE KEYS ?", "doc4")
	//fmt.Println("delete:", res, err)

	//os.Exit(1)
}

func dateFormat(layout string, intTime int64) string {
	//intTime := int64(d)
	t := time.Unix(intTime, 0)
	if layout == "" {
		layout = "2006-01-02 15:04:05"
	}
	return t.Format(layout)
}

func vAuthToken(loginId, accountId, username, role, ipAddress string) (string, error) {

	expireDuration := time.Now().Add(time.Hour * 24 * 30).Unix() //30 days long
	var row = make(map[string]interface{})
	row["cid"] = 1
	row["login_id"] = loginId
	row["parent_acc"] = accountId //
	row["email"] = username
	row["role"] = role
	row["session_code"] = xid.New().String()
	row["exp"] = expireDuration
	token, err := utility.JWTEncode(row, utility.JWTSECRET)
	if err == nil {
		sql := fmt.Sprintf("UPDATE %s SET update_date='%s',status=0 WHERE login_id='%s' AND status=1;", tableToBucket("authc"), mtool.TimeNow(), loginId)
		database.DB.Exec(sql)
		expireDate := dateFormat("", expireDuration)
		addAuthc(loginId, token, ipAddress, expireDate)
	}
	return token, err
}

func getSessionInfo(r *http.Request) (map[string]interface{}, error) {

	tokenStr, err := getCookie("login_session", r)
	if err != nil {
		return nil, err
	}
	return utility.JWTDecode(tokenStr, utility.JWTSECRET)
}

func userAgntDetails(uagentStr string) (device, osVersion, browserVersion string) {

	//uagentStr := r.UserAgent()
	//var device, osVersion, browserVersion string
	ua := useragent.Parse(uagentStr)
	if ua.Desktop {
		device = "Desktop"
	} else if ua.Mobile {
		device = "Mobile"
		if ua.IsAndroid() {
			device = "Android_Mobile"
		} else if ua.IsIOS() {
			device = "IOS_Mobile"
		}

	} else if ua.Tablet {
		device = "Tablet"
	} else if ua.Bot {
		device = "Bot"
		if ua.IsFacebookbot() {
			device = "FacebookBot"
		}
		if ua.IsGooglebot() {
			device = "GooglBot"
		}
		if ua.IsTwitterbot() {
			device = "TwitterBot"
		}
	} else {
		device = "Unknown"
	}

	if ua.IsLinux() {
		osVersion = "Linux-" + ua.OSVersion
	} else if ua.IsWindows() {
		osVersion = "Windows-" + ua.OSVersion
	} else if ua.IsMacOS() {
		osVersion = "MacOS-" + ua.OSVersion
	} else if ua.IsChromeOS() {
		osVersion = "ChromeOS-" + ua.OSVersion
	} else {
		osVersion = ua.OS + "*" + ua.OSVersion
	}

	if ua.IsChrome() {
		browserVersion = "chrome-" + ua.Version
	} else if ua.IsEdge() {
		browserVersion = "edge-" + ua.Version
	} else if ua.IsOpera() {
		browserVersion = "opera-" + ua.Version
	} else if ua.IsFirefox() {
		browserVersion = "firefox-" + ua.Version
	} else if ua.IsSafari() {
		browserVersion = "safari-" + ua.Version
	} else if ua.IsInternetExplorer() {
		browserVersion = "internetExplorer-" + ua.Version
	} else {
		browserVersion = "Unknown-0"
	}
	return
}

func visitorInfo(r *http.Request, w http.ResponseWriter) (sessionCode string) {

	ip := cleanIp(r.RemoteAddr)
	charger := r.FormValue("charger") //
	screen := r.FormValue("screen")
	uagentStr := r.UserAgent()

	if ip == "" {
		ip = mtool.ReadUserIP(r)
	}
	var todo, device, osVersion, browserVersion string
	device, osVersion, browserVersion = userAgntDetails(uagentStr)
	if charger == "Yes" {
		device = "Laptop"
	}
	var row = make(map[string]interface{})
	vcook, err := r.Cookie("visitor_session")
	if err == nil {
		sessionCode = vcook.Value
		todo = "update"
	} else {
		todo = "insert"
		sessionCode = xid.New().String()
	}

	//fmt.Println("todo:", todo)
	if todo == "insert" {

		var city, country string
		location := getLocationWithinSec(ip)
		slc := strings.Split(location, ",")
		if len(slc) == 2 {
			city = slc[0]
			country = slc[1]
		}
		row["id"] = xid.New().String()
		row["cid"] = COMPANY_ID
		row["city"] = city
		row["country"] = country
		row["session_code"] = sessionCode
		row["device"] = device
		row["screen_size"] = screen
		row["browser_version"] = browserVersion
		row["os_version"] = osVersion
		row["ip_address"] = ip
		row["vcount"] = 0
		row["create_date"] = mtool.TimeNow()
		row["status"] = 1
		row["table"] = structName(VisitorSession{})
		row["todo"] = todo
		//row["pkfield"] = "id"
		err = lxql.InsertUpdateMap(row, database.DB)
		logError("visitor_session", err)
		setCookie("visitor_session", sessionCode, 86400*365, w) //1 year

	} else if todo == "update" {

		sql := fmt.Sprintf("UPDATE %s SET ip_address=%q, vcount=vcount+1, update_date=%q WHERE session_code=%q;", tableToBucket("visitor_session"), ip, mtool.TimeNow(), sessionCode)
		database.DB.Exec(sql)
		if screen != "" {
			sql := fmt.Sprintf("UPDATE %s SET screen_size=%q, update_date=%q WHERE session_code=%q;", tableToBucket("visitor_session"), screen, mtool.TimeNow(), sessionCode)
			database.DB.Exec(sql)
		}
		if ip != "" {
			sql := fmt.Sprintf("UPDATE %s SET ip_address=%q, update_date=%q WHERE session_code =%q;", tableToBucket("visitor_session"), ip, mtool.TimeNow(), sessionCode)
			database.DB.Exec(sql)
		}
	}
	return
}

func loginToAccountRow(loginID string) (map[string]interface{}, error) {
	sql := fmt.Sprintf("SELECT l.account_id,a.parent_id,a.account_type,a.account_name FROM %s a LEFT JOIN %s l ON l.account_id=a.id WHERE l.id='%s';", tableToBucket("account"), tableToBucket("login"), loginID)
	return singleRow(sql)
}
func usernameToLoginId(username string) string {
	return lxql.FieldByValue("login", "id", fmt.Sprintf("username='%s'", username), database.DB)
}

func loginNotificationEmail(email, ipAddress, browser string) {

	subject := "LxRoot | Login Notification"
	emailTemplate := settingsValue("login_email")

	dmap := make(map[string]interface{})
	dmap["username"] = email
	dmap["ip"] = ipAddress
	dmap["browser"] = strings.ToUpper(browser)
	dmap["time"] = time.Now().Format("Mon, 02 Jan 2006 15:04:05 MST")
	emailBody, _ := templatePrepare(emailTemplate, dmap)
	err = SendEmail([]string{email}, subject, emailBody)
	logError("loginNotificationEmail", err)
}

func vcode() string {
	return uuid.NewV4().String()
}

func resetPassNotificationEmail(email, ipAddress, browser string) error {

	subject := "Reset your LxRoot password"
	emailTemplate := settingsValue("resetpass_email")

	code := utility.GeneratePassword(6, true, false)
	id, err := addVerification(email, "RESET-PASS", code, "")
	if err != nil {
		return err
	}

	var row = make(map[string]interface{})
	row["email"] = email
	row["id"] = id
	row["code"] = code
	row["exp"] = time.Now().Add(time.Minute * 30).Unix() //***
	token, err := utility.JWTEncode(row, utility.JWTSECRET)
	if err != nil {
		return err
	}

	dmap := make(map[string]interface{})
	dmap["username"] = email
	dmap["ip"] = ipAddress
	dmap["browser"] = strings.ToUpper(browser)
	dmap["reset_link"] = fmt.Sprintf("https://lxroot.com/reset-pass-form?token=%s", token)
	emailBody, _ := templatePrepare(emailTemplate, dmap)
	err = SendEmail([]string{email}, subject, emailBody)
	logError("resetPassNotificationEmail", err)
	return err
}

func strToTime(expStr string) time.Time {

	expUnix, _ := strconv.ParseInt(expStr, 10, 64)
	expTime := time.Unix(expUnix, 0)
	return expTime
}

func float64totime(val float64) time.Time {
	intTimestamp := int64(val)
	return time.Unix(intTimestamp, 0)
}

func toTime(val interface{}) time.Time {

	switch s := val.(type) {

	case string:
		sint64, _ := strconv.ParseInt(s, 10, 64)
		return time.Unix(sint64, 0)

	case float64:
		return time.Unix(int64(s), 0)

	case int:
		return time.Unix(int64(s), 0)

	case int64:
		return time.Unix(s, 0)

	default:
		fmt.Printf("toTime() -> type of val is %T", s)
	}
	return time.Time{}
}

// true = not expire, false = expired
func checkUnExpired(expStr string) bool {

	// Convert Unix timestamp to time.Time
	expTime := strToTime(expStr)

	if expTime.After(time.Now()) {
		return true
	}
	// Check if the expiration time is in the future
	return false
}

// true = not expire, false = expired
func checkUnExpired2(expStr float64) bool {

	// Convert float64 to time.Time
	expTime := float64totime(expStr)

	if expTime.After(time.Now()) {
		return true
	}
	// Check if the expiration time is in the future
	return false
}

func tokenInfo(token string) (map[string]interface{}, error) {
	return utility.JWTDecode(token, utility.JWTSECRET)
}

func checkTokenCodeValid(row map[string]interface{}) bool {

	expStr, isOk := row["exp"].(float64)
	if isOk {
		//intTimestamp := int64(expStr)
		//tm := time.Unix(intTimestamp, 0)
		//fmt.Println(tm.Format("2006-01-02 15:04:05"))
		//fmt.Println("isValid:", isValid)
		return checkUnExpired2(expStr)
	}
	return false
}

func deleteAccount(email string) error {

	//loginId := lxql.FieldByValue(tableToBucket("login"), "id", fmt.Sprintf("account_id='%s'", accountId), database.DB)
	//fmt.Println(loginId)
	sql := fmt.Sprintf(`SELECT a.id as accountId, l.id as loginId,l.username, i.id as addressId FROM %s a
	LEFT JOIN %s l ON l.account_id=a.id
	LEFT JOIN %s i ON i.account_id=a.id
	WHERE a.email="%s";`, tableToBucket("account"), tableToBucket("login"), tableToBucket("address"), email)
	row, err := singleRow(sql)
	if err != nil {
		return err
	}
	accountId := row["accountId"].(string)
	loginId := row["loginId"].(string)
	addressId := row["addressId"].(string)
	username := row["username"].(string)
	fmt.Println(loginId, addressId)

	sql = fmt.Sprintf("DELETE FROM %s WHERE id=%q;", tableToBucket("address"), addressId)
	err = lxql.RawSQL(sql, database.DB)
	logError("address", err)

	sql = fmt.Sprintf("DELETE FROM %s WHERE login_id=%q;", tableToBucket("activity_log"), loginId)
	err = lxql.RawSQL(sql, database.DB)
	logError("activity_log", err)

	sql = fmt.Sprintf("DELETE FROM %s WHERE login_id=%q;", tableToBucket("login_session"), loginId)
	err = lxql.RawSQL(sql, database.DB)
	logError("login_session", err)

	sql = fmt.Sprintf("DELETE FROM %s WHERE login_id=%q;", tableToBucket("authc"), loginId)
	err = lxql.RawSQL(sql, database.DB)
	logError("authc", err)

	//indirect relation tables -> verification[username], message[sender|receiver]
	sql = fmt.Sprintf("DELETE FROM %s WHERE username=%q;", tableToBucket("verification"), username)
	err = lxql.RawSQL(sql, database.DB)
	logError("verification", err)

	// sql = fmt.Sprintf("DELETE FROM %s WHERE login_id=%q;", tableToBucket("device_log"), loginId)
	// err = lxql.RawSQL(sql, database.DB)
	// logError("device_log", err)

	// sql = fmt.Sprintf("DELETE FROM %s WHERE login_id=%q;", tableToBucket("doc_keeper"), loginId)
	// err = lxql.RawSQL(sql, database.DB)
	// logError("doc_keeper", err)

	//sql = fmt.Sprintf("DELETE FROM %s WHERE sender=%q OR receiver=%q;", tableToBucket("message"), username, username)
	//fmt.Println(sql)
	//lxql.RawSQL(sql, database.DB)

	sql = fmt.Sprintf("DELETE FROM %s WHERE id=%q;", tableToBucket("account"), accountId)
	err = lxql.RawSQL(sql, database.DB)
	logError("account", err)

	sql = fmt.Sprintf("DELETE FROM %s WHERE id=%q;", tableToBucket("login"), loginId)
	err = lxql.RawSQL(sql, database.DB)
	logError("login", err)

	return nil
}

func addToCart(itemId, qty, docRef, docNumber, loginId, accountId, ipAddress string) (docId string, err error) {

	docName := "shopping cart"
	docType := "cart"
	postingDate := ""
	docStatus := "pending"
	docId = docNumber

	sql := fmt.Sprintf("SELECT item_name,item_code,sale_price FROM %s WHERE id='%s';", tableToBucket("item"), itemId)
	row, err := singleRow(sql)
	logError("addToCart", err)
	price, _ := row["sale_price"].(string)
	itemInfo, _ := row["item_name"].(string)
	stripePriceId, _ := row["item_code"].(string) //
	itemSerial := ""

	totalDiscount := ""
	totalTax := ""
	totalPayable := fmt.Sprint(str2int(qty) * str2int(price))

	var docUpdate bool = true
	if docNumber == "" {
		docUpdate = false
		docId, err = addDocKeeper(docName, docType, docRef, docNumber, postingDate, docStatus, totalDiscount, totalTax, totalPayable, loginId, accountId, ipAddress)
		if err != nil {
			return
		}
	}

	addTransactionRecord(docType, docId, itemId, itemInfo, itemSerial, stripePriceId, qty, price)
	if docUpdate {
		sql = fmt.Sprintf("UPDATE %s SET total_payable=total_payable+%s WHERE doc_number=%q;", tableToBucket("doc_keeper"), totalPayable, docId)
		lxql.RawSQL(sql, database.DB)
	}
	return
}

func structFieldValMap(anyStructToPointer interface{}) (map[string]interface{}, error) {

	empv := reflect.ValueOf(anyStructToPointer) //must be a pointer
	empt := reflect.TypeOf(anyStructToPointer)  //

	if empv.Kind() != reflect.Pointer {
		return nil, errors.New("anyStructToPointer must be a pointer")
	}

	var row = make(map[string]interface{})

	for i := 0; i < empv.Elem().NumField(); i++ {

		//vField := empv.Elem().Field(i)
		var key string
		kfield := empt.Elem().Field(i)
		key = kfield.Name
		tag := kfield.Tag.Get("json")
		if tag != "" {
			key = tag
		}
		//fmt.Println(empt.Elem().Field(i), empv.Elem().Field(i))
		row[key] = empv.Elem().Field(i).Interface()
	}
	//fmt.Println(row, len(row))
	return row, nil

}

func xidToNumber(id string) int32 {
	mid, err := xid.FromString(id)
	if err != nil {
		return 0
	}
	return mid.Counter()
}

func saveFile(fh *multipart.FileHeader, fileAbsPath string) error {

	file, err := fh.Open()
	if err != nil {
		return err
	}
	defer file.Close()
	buff := make([]byte, 512)
	_, err = file.Read(buff)
	if err != nil {
		return err
	}
	_, err = file.Seek(0, io.SeekStart)
	if err != nil {
		return err
	}
	os.MkdirAll(filepath.Dir(fileAbsPath), 0755)
	f, err := os.Create(fileAbsPath)
	if err != nil {
		return err
	}
	defer f.Close()
	_, err = io.Copy(f, file)
	if err != nil {
		return err
	}
	return nil
}

func parseMultipartTodo(r *http.Request) (todo string) {

	r.ParseForm()
	contentType := r.Header.Get("Content-Type")
	if strings.Contains(contentType, "multipart/form-data") {
		r.ParseMultipartForm(r.ContentLength)
	}
	todo = r.FormValue("todo")
	if todo == "" {
		todo = "search"
	}
	return
}

func docFindDateTime() (startDate, endDate string) {

	tnow := time.Now()
	diff := 1 * time.Minute
	then := tnow.Add(-diff)
	startDate = then.Format(DATE_TIME_FORMAT)
	endDate = tnow.Format(DATE_TIME_FORMAT)
	return
}

func subscriptionStartEnd() (startDate, endDate string) {

	tnow := time.Now()
	then := tnow.AddDate(0, 1, 0)
	startDate = tnow.Format(DATE_TIME_FORMAT)
	endDate = then.Format(DATE_TIME_FORMAT)
	return
}

func nameLabel(firstName, lastName string) (label string) {

	if firstName != "" && lastName != "" {
		label = fmt.Sprintf("%c%c", firstName[0], lastName[0])
	}
	return
}
