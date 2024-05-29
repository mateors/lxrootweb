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
		//fmt.Println("<>", ftoken, ctokenh, ftokenh)
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

	email := args["email"].(string)
	sql := fmt.Sprintf("SELECT count(*)as cnt FROM %s WHERE username='%s';", tableToBucket("login"), email)
	//fmt.Println(sql)
	rows, err := lxql.GetRows(sql, db)
	if err != nil {
		return "ERROR wrong query"
	}
	if len(rows) > 0 {
		//fmt.Println(rows, len(rows))
		//for _, row := range rows {
		//fmt.Printf("%v %T\n", row, row)
		//fmt.Println(row["cnt"])
		//for key, val := range row {
		//fmt.Printf("%v %v %T\n", key, val, val)
		//}
		//}
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
	prow := db.QueryRow(sql)
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
