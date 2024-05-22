package lxql

import (
	"bytes"
	"database/sql"
	"database/sql/driver"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"sync"
)

var params map[string]string
var N1QL_PASSTHROUGH_MODE = false

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

	settQueryParams(&postData)
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
	fmt.Println(params, dataSourceName, queryAPI, resp.Status)
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

func (conn *n1qlConn) doClientRequest(query string, requestValues *url.Values) (*http.Response, error) {

	ok := false
	for !ok {

		var request *http.Request
		var err error

		// select query API
		//rand.Seed(time.Now().Unix())
		numNodes := 1 //len(conn.queryAPIs)

		//selectedNode := rand.Intn(numNodes)
		conn.lock.RLock()
		queryAPI := conn.queryAPI //conn.queryAPIs[selectedNode]
		conn.lock.RUnlock()

		if query != "" {
			request, err = prepareRequest(query, queryAPI, nil)
			if err != nil {
				return nil, err
			}
		} else {
			if requestValues != nil {
				request, _ = http.NewRequest("POST", queryAPI, bytes.NewBufferString(requestValues.Encode()))
			} else {
				request, _ = http.NewRequest("POST", queryAPI, nil)
			}
			request.Header.Add("Content-Type", "application/x-www-form-urlencoded")
		}

		resp, err := conn.client.Do(request)
		if err != nil {
			// if this is the last node return with error
			if numNodes == 1 {
				break
			}
			// remove the node that failed from the list of query nodes
			//conn.lock.Lock()
			//conn.queryAPIs = append(conn.queryAPIs[:selectedNode], conn.queryAPIs[selectedNode+1:]...)
			//conn.lock.Unlock()
			continue
		} else {
			return resp, nil
		}
	}
	return nil, fmt.Errorf("N1QL: Query nodes not responding")
}

func decodeSignature(signature *json.RawMessage) interface{} {

	var sign interface{}
	var rows map[string]interface{}
	json.Unmarshal(*signature, &sign)

	switch s := sign.(type) {
	case map[string]interface{}:
		return s
	case string:
		return s
	default:
		fmt.Printf(" Cannot decode signature. Type of this signature is %T", s)
		return map[string]interface{}{"*": "*"}
	}
	return rows
}

func (conn *n1qlConn) performQuery(query string, requestValues *url.Values) (driver.Rows, error) {

	resp, err := conn.doClientRequest(query, requestValues)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != 200 {
		bod, _ := io.ReadAll(io.LimitReader(resp.Body, 512))
		return nil, fmt.Errorf("%s", bod)
	}

	var resultMap map[string]*json.RawMessage
	decoder := json.NewDecoder(resp.Body)

	err = decoder.Decode(&resultMap)
	if err != nil {
		return nil, fmt.Errorf(" N1QL: Failed to decode result %v", err)
	}

	var signature interface{}
	var resultRows *json.RawMessage
	var metrics interface{}
	var status interface{}
	var requestId interface{}
	var errs interface{}

	for name, results := range resultMap {
		switch name {
		case "errors":
			_ = json.Unmarshal(*results, &errs)
		case "signature":
			if results != nil {
				signature = decodeSignature(results)
			} else if N1QL_PASSTHROUGH_MODE {
				// for certain types of DML queries, the returned signature could be null
				// however in passthrough mode we always return the metrics, status etc as
				// rows therefore we need to ensure that there is a default signature.
				signature = map[string]interface{}{"*": "*"}
			}
		case "results":
			resultRows = results
		case "metrics":
			if N1QL_PASSTHROUGH_MODE {
				_ = json.Unmarshal(*results, &metrics)
			}
		case "status":
			if N1QL_PASSTHROUGH_MODE {
				_ = json.Unmarshal(*results, &status)
			}
		case "requestID":
			if N1QL_PASSTHROUGH_MODE {
				_ = json.Unmarshal(*results, &requestId)
			}
		}
	}

	if N1QL_PASSTHROUGH_MODE {
		extraVals := map[string]interface{}{"requestID": requestId,
			"status":    status,
			"signature": signature,
		}

		// in passthrough mode last line will always be en error line
		errors := map[string]interface{}{"errors": errs}
		return resultToRows(bytes.NewReader(*resultRows), resp, signature, metrics, errors, extraVals)
	}

	// we return the errors with the rows because we can have scenarios where there are valid
	// results returned along with the error and this interface doesn't allow for both to be
	// returned and hence this workaround.
	return resultToRows(bytes.NewReader(*resultRows), resp, signature, nil, errs, nil)
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
	//conn.performQuery(query, nil)
	return conn.performQuery(query, nil)
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
