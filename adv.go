package main

import (
	"bytes"
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/json"
	"errors"
	"fmt"
	"html/template"
	"io"
	"log"
	"lxrootweb/database"
	"lxrootweb/lxql"
	"net"
	"net/http"
	"net/url"
	"reflect"
	"strconv"
	"strings"
	"time"
	"unicode"

	"github.com/mateors/money"
	"github.com/mateors/mtool"
	uuid "github.com/satori/go.uuid"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

type charInfo struct {
	Index int
	Char  rune
}

var FuncMap = template.FuncMap{

	"moneyFormat":      moneyFormat,
	"subTotal":         subTotal,
	"taxTotal":         taxTotal,
	"checkoutDisabled": checkoutDisabled,
	"toTitle":          toTitle,
	"toUpper":          toUpper,
	"paymentIcon":      paymentIcon,
	"niceDate":         niceDate,
}

const DATE_TIME_FORMAT = "2006-01-02 15:04:05"
const NICE_DATE_FORMAT = "January 02, 2006"

func niceDate(createDateTime string) string {
	return mtool.DateTimeParser(createDateTime, DATE_TIME_FORMAT, NICE_DATE_FORMAT)
}

func paymentIcon(paymentStatus string) string {

	if paymentStatus == "paid" {
		return "paid"
	} else if paymentStatus == "refunded" {
		return "restore"
	}

	return "info"
}

func toUpper(text string) string {
	return strings.ToUpper(text)
}

func toTitle(str string) string {
	caser := cases.Title(language.English)
	return caser.String(str)
}

func checkoutDisabled(loginRequired bool, cartCount int) bool {

	if !loginRequired && cartCount == 0 {
		return true
	} else if loginRequired && cartCount > 0 {
		return true
	}
	return false
}

func moneyFormat(amount interface{}) string {

	return money.CommaSeparatedMoneyFormat(amount)
}

func subTotal(data []map[string]interface{}) float64 {

	var total float64
	for _, row := range data {
		payableAmount, _ := strconv.ParseFloat(row["payable_amount"].(string), 64)
		total += payableAmount
	}
	return total
}
func taxTotal(data []map[string]interface{}) float64 {

	var total float64
	for _, row := range data {
		amount, _ := strconv.ParseFloat(row["tax_amount"].(string), 64)
		total += amount
	}
	return total
}

func structName(myvar interface{}) string {
	if t := reflect.TypeOf(myvar); t.Kind() == reflect.Ptr {
		return t.Elem().Name()
	} else {
		return t.Name()
	}
}

func vMapToJsonStr(vMap map[string]interface{}) string {

	bs, err := json.Marshal(&vMap)
	if err != nil {
		return ""
	}
	return string(bs)
}

func tableToBucket(table string) string {
	return fmt.Sprintf(`%s.%s.%s`, BUCKET_NAME, SCOPE_NAME, table)
}

// helper func
func upperCount(text string) []*charInfo {

	var list []*charInfo
	for i, ch := range text {
		if unicode.IsUpper(ch) {
			list = append(list, &charInfo{i, ch})
		}
	}
	return list
}

// helper func
func splitByUpperCase(text string) []string {

	var list []string
	ci := upperCount(text)
	splitIndex := 0
	for i, c := range ci {
		if i > 0 {
			//fmt.Println(i, text[splitIndex:c.Index]) //c.Index, c.Char, fmt.Sprintf(`%c`, c.Char)
			list = append(list, text[splitIndex:c.Index])
			splitIndex = c.Index
		}
	}
	if splitIndex < len(text) {
		//fmt.Println(name[splitIndex:])
		list = append(list, text[splitIndex:])
	}
	return list
}

func customTableName(structName string) string {

	list := splitByUpperCase(structName)
	for i, part := range list {
		list[i] = strings.ToLower(part)
	}
	return strings.Join(list, "_")
}

// Anti CSRF token
func csrfToken() string {
	return uuid.NewV4().String()
}

func hmacHash(message, secretKey string) string {
	mac := hmac.New(sha256.New, []byte(secretKey))
	mac.Write([]byte(message))
	return fmt.Sprintf("%x", mac.Sum(nil))
}

func getCookie(cookieName string, r *http.Request) (string, error) {

	ecookie, err := r.Cookie(cookieName)
	if err != nil {
		return "", err
	}
	return ecookie.Value, nil
}

func delCookie(cookieName string, r *http.Request, w http.ResponseWriter) error {

	pcookie, err := r.Cookie(cookieName)
	if err != http.ErrNoCookie {
		pcookie.Name = cookieName
		pcookie.MaxAge = -1
		pcookie.Value = ""
		pcookie.Path = "/"
		pcookie.HttpOnly = true
		http.SetCookie(w, pcookie)
	}
	return err
}

func setCookie(cookieName, value string, timeInSec int, w http.ResponseWriter) {

	c := &http.Cookie{
		Name:     cookieName,
		Value:    value,
		MaxAge:   timeInSec, //300, //5 minutes = 300 seconds
		SameSite: http.SameSiteLaxMode,
		HttpOnly: true,
	}
	http.SetCookie(w, c)
}

func fCall(m map[string]interface{}, name string, params ...interface{}) (result []reflect.Value, err error) {

	f := reflect.ValueOf(m[name])
	if len(params) != f.Type().NumIn() {
		err = errors.New("the number of params is not adapted")
		return
	}
	in := make([]reflect.Value, len(params))
	for k, param := range params {
		in[k] = reflect.ValueOf(param)
	}
	result = f.Call(in)
	return

}

// slice:=["ERROR invalid username", "valid"]
func errorInSlice(slice []string, val string) (int, bool) {
	for i, item := range slice {
		if strings.Contains(item, val) {
			return i, true
		}
	}
	return -1, false
}

func CheckMultipleConditionTrue(args map[string]interface{}, funcsMap map[string]interface{}) string {

	var response string
	resAray := make([]string, 0)
	for key := range funcsMap {
		result, err := fCall(funcsMap, key, args) //result is type of []reflect.Value
		if err != nil {
			log.Println(err)
		}
		res, isOk := result[0].Interface().(string) //Converting reflect.Value to string
		if !isOk {
			log.Println("error at CheckMultipleConditionTrue")
		}
		resAray = append(resAray, res)
	}
	i, errorExist := errorInSlice(resAray, "ERROR")
	if errorExist {
		response = fmt.Sprintf("%v", resAray[i]) //ERROR EXIST in CheckError =>>
	} else {
		response = "OKAY"
	}
	return response
}

func validCSRF(args map[string]interface{}) string {
	ctokenh, isOk := args["ctoken"].(string)
	if isOk {
		ftoken := fmt.Sprint(args["ftoken"])
		ftokenh := hmacHash(ftoken, ENCDECPASS)
		if ftokenh == ctokenh {
			return "valid"
		}
	}
	return "ERROR invalid ctoken"
}

func resetPassValidation(args map[string]interface{}) string {

	pass1, _ := args["pass1"].(string)
	pass2, _ := args["pass2"].(string)

	if pass1 != pass2 {
		return "ERROR both passwords must match"
	}
	if len(pass1) < 8 && len(pass2) < 8 {
		return "ERROR a minimum 8-character password is required"
	}
	return "valid"
}

func vcodeIdValidation(args map[string]interface{}) string {

	vid, isOk := args["vid"].(string)
	if isOk {
		count := lxql.CheckCount("verification", fmt.Sprintf("id=%q AND status=0", vid), database.DB)
		if count == 1 {
			return "valid"
		}
	}
	return "ERROR invalid vsession"

}

func validSignupField(args map[string]interface{}) string {

	firstName, isOk := args["first_name"].(string)
	if !isOk {
		return "ERROR first name field required"
	}
	if firstName == "" {
		return "ERROR first_name is required"
	}
	lastName, isOk := args["last_name"].(string)
	if !isOk {
		return "ERROR last name field required"
	}
	if lastName == "" {
		return "ERROR last_name is required"
	}
	passwd, isOk := args["passwd"].(string)
	if !isOk {
		return "ERROR password field required"
	}
	if passwd == "" {
		return "ERROR password is required"
	}
	return "valid"
}

func validEmail(args map[string]interface{}) string {

	email, isOk := args["email"].(string)
	if !isOk {
		return "ERROR email missing"
	}
	email = strings.ToLower(email)
	count := lxql.CheckCount("login", fmt.Sprintf(`username="%s"`, email), database.DB)
	if count > 0 {
		return "ERROR email already exist"
	}
	return "valid"
}

func validUserName(args map[string]interface{}) string {

	username, isOk := args["email"].(string)
	if !isOk {
		return "ERROR username missing"
	}
	count := lxql.CheckCount("login", fmt.Sprintf(`username="%s"`, username), database.DB)
	if count == 1 {
		return "valid"
	}
	return "ERROR invalid username or password"
}

func logError(prefix string, err error) {
	if err != nil {
		log.Printf("ERR_%s: %s", prefix, err.Error())
	}
}

// value parser
func vParser(vtype, key string, form url.Values) (output interface{}) {

	value := form.Get(key)
	if vtype == "int" {
		cval, _ := strconv.Atoi(value)
		output = cval

	} else if vtype == "int64" {
		cval, _ := strconv.ParseInt(value, 10, 64)
		output = cval

	} else if vtype == "float64" {
		cval, _ := strconv.ParseFloat(value, 64)
		output = cval

	} else if vtype == "slice" {
		output = form[key]

	} else {
		output = value
	}
	return
}

func tokenPullNSet(r *http.Request) error {

	ctoken, err := getCookie("ctoken", r) //cross check with form token
	if err != nil {
		log.Println("ERR_1tokenPullNSet:", err)
		return err
	}
	token, err := getCookie("token", r) //API_TOKEN
	if err != nil {
		log.Println("ERR_2tokenPullNSet:", err)
		return err
	}
	r.Form.Set("token", token)   //API
	r.Form.Set("ctoken", ctoken) //uuid.NewV4()
	return nil
}

func companyId(website string) string {

	sql := fmt.Sprintf("SELECT id FROM %s WHERE website='%s';", tableToBucket("company"), website)
	prow := database.DB.QueryRow(sql)
	var cmap = make(map[string]interface{}, 0)
	var jsonBytes []uint8
	err := prow.Scan(&jsonBytes)
	if err != nil {
		log.Fatal(err)
	}
	json.Unmarshal(jsonBytes, &cmap)
	return cmap["id"].(string)
}

func commonDataSet(r *http.Request) error {

	ctoken, err := getCookie("ctoken", r) //cross check with form token
	if err != nil {
		log.Println("ERR_token:", err)
		return err
	}
	ipAddress := cleanIp(mtool.ReadUserIP(r))
	r.Form.Set("ip_address", ipAddress)
	r.Form.Set("cid", COMPANY_ID)
	r.Form.Set("create_date", mtool.TimeNow())
	r.Form.Set("ctoken", ctoken)
	r.Form.Set("status", "1")
	return nil
}

func cleanIp(ipwithport string) string {
	slc := strings.Split(ipwithport, ":")
	if len(slc) == 2 {
		return slc[0]
	}
	return ipwithport
}

func templatePrepare(tmpltText string, dmap map[string]interface{}) (string, error) {

	var tplOutput bytes.Buffer
	tpl := template.New("email")
	eml, err := tpl.Parse(tmpltText)
	if err != nil {
		return "", err
	}
	err = eml.Execute(&tplOutput, dmap)
	if err != nil {
		return "", err
	}
	return tplOutput.String(), nil
}

func countryImportFromExcel(filePath string) error {

	rows, err := excelReader(filePath, "Sheet1", nil)
	if err != nil {
		return err
	}

	for _, row := range rows {
		//fmt.Println(row)
		name := row["name"].(string)
		isoCode := row["iso_code"].(string)
		countryCode := row["country_code"].(string)
		err := addCountry(name, isoCode, countryCode)
		fmt.Println(err, name, isoCode, countryCode)
	}
	return nil
}

func IsIPv4(str string) bool {
	ip := net.ParseIP(str)
	return ip != nil && ip.To4() != nil
}

func IsIPv6(str string) bool {
	ip := net.ParseIP(str)
	return ip != nil && ip.To4() == nil
}

func getCountryRegionFromIp(ipAddress string) (map[string]string, error) {

	var row = make(map[string]string)
	rurl := fmt.Sprintf("http://ip-api.com/json/%s?fields=country,regionName", ipAddress)
	resp, err := http.Get(rurl)
	if err != nil {
		return nil, err
	}

	remaining := resp.Header.Get("X-Rl")
	ttl := resp.Header.Get("X-Ttl")
	err = json.NewDecoder(resp.Body).Decode(&row)
	if err != nil {
		return nil, err
	}
	row["remaining"] = remaining //
	row["ttl"] = ttl             // check if it is 0
	row["status"] = resp.Status  // check if 429
	//fmt.Println(row)
	return row, nil
}

func getIpToCountry(ipv4Address string) (string, error) {

	if !IsIPv4(ipv4Address) {
		return "", errors.New("ipv4 allowed only")
	}
	url := fmt.Sprintf("http://ip2c.org/?ip=%s", ipv4Address)
	resp, err := http.Get(url)
	if err != nil {
		return "", err
	}
	bs, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	text := string(bs)
	//slc := strings.Split(text, ";")
	//fmt.Println(slc, len(slc))
	return text, nil
}

func getLocation(ipAddress string) (string, error) {

	var country, region string
	row, err := getCountryRegionFromIp(ipAddress)
	if err != nil {
		return "", err
	}
	country = row["country"]
	region = row["regionName"]
	status := row["status"]
	region = strings.TrimSpace(strings.Replace(region, "Division", "", -1))
	//fmt.Println(row["ttl"], row["remaining"])

	if strings.Contains(status, "429") {
		if IsIPv4(ipAddress) {
			txt, err := getIpToCountry(ipAddress)
			logError("", err)
			slc := strings.Split(txt, ";")
			if len(slc) == 4 {
				country = slc[3]
			}
		} else {
			region = "IP"
			country = ipAddress
			log.Println(status, ipAddress, "ipv6 address has no provider yet")
		}
	}
	//time.Sleep(time.Second * 1)
	return fmt.Sprintf("%s,%s", region, country), nil
}

func getLocationWithinSec(ipAddress string) string {

	var location string
	// type output struct {
	// 	out string
	// 	err error
	// }
	ch := make(chan string)
	go func() {
		res, err := getLocation(ipAddress)
		logError("getLocation", err)
		ch <- res
	}()

	select {

	case <-time.After(1 * time.Second):
		location = "" //"timed out"
		log.Println("getLocation() time out")

	case x := <-ch:
		location = x
	}
	return location
}

func performTask(ctx context.Context, output chan<- string) {

	time.Sleep(time.Second * 3)
	output <- "mostain"
}

func ctxReq() {

	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()

	ch := make(chan string)
	go performTask(ctx, ch)

	select {
	case <-ctx.Done():
		fmt.Println("Task timed out")

	case x := <-ch:
		fmt.Println(x)
	}
}
