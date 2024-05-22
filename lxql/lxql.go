package lxql

import (
	"bytes"
	"database/sql"
	"database/sql/driver"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"sync"
)

var params map[string]string

// HTTPClient to use for REST and view operations.
var MaxIdleConnsPerHost = 10
var HTTPTransport = &http.Transport{MaxIdleConnsPerHost: MaxIdleConnsPerHost}
var HTTPClient = &http.Client{Transport: HTTPTransport}

func SetQueryParams(key string, value string) error {

	if key == "" {
		return fmt.Errorf("N1QL: Key not specified")
	}
	params[key] = value
	return nil
}

// implements Driver interface
type n1qlDriver struct{}

func init() {
	sql.Register("n1ql", &n1qlDriver{})
	params = make(map[string]string)
	fmt.Println("params...")
}

//username:password@protocol(address)/dbname?param=value
func (n *n1qlDriver) Open(dataSourceName string) (driver.Conn, error) {
	fmt.Println("open...", dataSourceName)
	return openConn2(dataSourceName)
}

func settQueryParams(v *url.Values) {

	fmt.Println("paramsLen:", params, len(params))
	for key, value := range params {
		v.Set(key, value)
		fmt.Println(key, "==>", value)
	}
}

// prepare a http request for the query
func prepareRequest(query string, queryAPI string, args []driver.Value) (*http.Request, error) {

	postData := url.Values{}
	postData.Set("statement", query)

	// if len(args) > 0 {
	// 	paStr := buildPositionalArgList(args)
	// 	if len(paStr) > 0 {
	// 		postData.Set("args", paStr)
	// 	}
	// }

	//settQueryParams(&postData)
	fmt.Println("prepareRequest...", postData)
	request, _ := http.NewRequest("POST", queryAPI, bytes.NewBufferString(postData.Encode()))
	request.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	return request, nil
}

func openConn2(dataSourceName string) (driver.Conn, error) {

	queryAPI := fmt.Sprintf("%s/query/service", dataSourceName)
	conn := &n1qlConn{client: HTTPClient, queryAPI: queryAPI}

	request, err := prepareRequest("SELECT 1", queryAPI, nil)
	if err != nil {
		return nil, err
	}

	resp, err := conn.client.Do(request)
	//fmt.Println(params, dataSourceName, resp.Status)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		bod, _ := io.ReadAll(io.LimitReader(resp.Body, 512))
		return nil, fmt.Errorf("N1QL: Connection failure %s", bod)
	}
	return conn, err
}

// implements driver.Conn interface
type n1qlConn struct {
	//clusterAddress string
	queryAPI string
	client   *http.Client
	lock     sync.RWMutex
}

func (conn *n1qlConn) Prepare(query string) (driver.Stmt, error) {

	//n1qlStmt
	return nil, nil
}

func (conn *n1qlConn) Close() error {

	return nil
}

func (conn *n1qlConn) Begin() (driver.Tx, error) {

	return nil, nil
}

func (conn *n1qlConn) Query(query string, args []driver.Value) (driver.Rows, error) {

	fmt.Println("query...", query)
	return nil, nil
}

func (conn *n1qlConn) Exec(query string, args []driver.Value) (driver.Result, error) {

	return nil, nil
}

// implements driver.Stmt interface
type n1qlStmt struct {
	conn      *n1qlConn
	prepared  string
	signature string
	argCount  int
	name      string
}

func (stmt *n1qlStmt) Close() error {

	return nil
}

func (stmt *n1qlStmt) NumInput() int {

	return 0
}

func (stmt *n1qlStmt) Exec(args []driver.Value) (driver.Result, error) {

	return nil, nil
}

func (stmt *n1qlStmt) Query(args []driver.Value) (driver.Rows, error) {

	return nil, nil
}
