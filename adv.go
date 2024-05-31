package main

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/json"
	"errors"
	"fmt"
	"html/template"
	"log"
	"lxrootweb/database"
	"lxrootweb/lxql"
	"net/http"
	"net/url"
	"reflect"
	"strconv"
	"strings"
	"unicode"

	"github.com/mateors/mtool"
	uuid "github.com/satori/go.uuid"
)

type charInfo struct {
	Index int
	Char  rune
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

func validSignupField(args map[string]interface{}) string {

	firstName := args["first_name"].(string)
	if firstName == "" {
		return "ERROR first_name is required"
	}
	lastName := args["last_name"].(string)
	if lastName == "" {
		return "ERROR last_name is required"
	}
	passwd := args["passwd"].(string)
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

func bytesToStr(slc []uint8) string {

	var str string
	for _, c := range slc {
		if c != 34 { //remove "
			str += fmt.Sprintf("%c", c)
		}
	}
	return str
}

func colsToRowMap(cols []string, orow map[string]interface{}) map[string]interface{} {

	//var orow = make(map[string]interface{})
	if len(cols) > 1 {
		for key := range orow {
			slc := orow[key].([]uint8)
			orow[key] = bytesToStr(slc)
		}

	} else if len(cols) == 1 {

		colname := cols[0]
		//fmt.Println(colname, orow)
		if colname == "*" {
			var row = make(map[string]map[string]interface{})
			json.Unmarshal(orow[colname].([]uint8), &row)
			for key := range row {
				orow = row[key]
			}
		} else {
			var row = make(map[string]interface{})
			json.Unmarshal(orow[colname].([]uint8), &row)
			orow = row
		}
	}
	return orow
}

func singleRow(sql string) (map[string]interface{}, error) {

	var orow = make(map[string]interface{})

	sql = strings.ToLower(sql)
	cols := GetColumnNamesFromQuery(sql) //[]string{"id", "status"}
	row := make([][]byte, len(cols))
	rowPtr := make([]any, len(cols))

	for i := range row {
		rowPtr[i] = &row[i]
	}

	srow := database.DB.QueryRow(sql)
	err := srow.Scan(rowPtr...)
	if err != nil {
		return nil, err
	}

	for i, rowp := range rowPtr {

		switch val := rowp.(type) {

		case *[]uint8:
			orow[cols[i]] = *val

		default:
			fmt.Println("GetRows() Type is unknown!")
		}
	}

	orow = colsToRowMap(cols, orow)
	return orow, nil
}

func GetColumnNamesFromQuery(query string) []string {

	query = strings.ToLower(query)
	// Remove leading/trailing whitespaces and semicolon (if any)
	query = strings.TrimSpace(query)
	if strings.HasSuffix(query, ";") {
		query = query[:len(query)-1]
	}

	// Split the query by spaces
	parts := strings.Fields(query)

	// Find the index of "SELECT" and "FROM" keywords
	selectIndex := -1
	fromIndex := -1
	for i, part := range parts {
		if strings.EqualFold(part, "select") {
			selectIndex = i
		} else if strings.EqualFold(part, "from") {
			fromIndex = i
			break
		}
	}

	// Extract column names between "SELECT" and "FROM"
	if selectIndex != -1 && fromIndex != -1 && fromIndex > selectIndex+1 {

		columns := parts[selectIndex+1 : fromIndex]
		// Remove commas from column names
		for i, col := range columns {
			columns[i] = strings.TrimSuffix(col, ",")
		}
		var rcols []string
		for _, col := range columns {
			slc := strings.Split(col, ",")
			rcols = slc
		}
		return rcols
	}
	return nil
}
