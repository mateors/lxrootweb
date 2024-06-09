package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/mateors/mtool"
)

// unixToDateTime 1642614110, "2006-01-02 15:04:05"
func unixToDateTime(sec int64) string {
	tm := time.Unix(sec, 0)
	return tm.Format("2006-01-02 15:04:05") //dateTimeFormat
}

func requestBalance(secretKey string) (map[string]interface{}, error) {

	url := "https://api.stripe.com/v1/balance"
	method := "GET"

	//payload := strings.NewReader("email=sunzida%40gmail.com&name=Dr%20Sanzida&phone=01765110255&source=tok_1IU9EmEJdAixg8N3xAwmFVBP")
	client := &http.Client{}
	req, err := http.NewRequest(method, url, nil)

	if err != nil {
		return nil, err
	}
	//req.Header.Add("Stripe-Account", "acct_1Es5BxEJdAixg8N3")
	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", secretKey))
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	res, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}
	//fmt.Println(string(body))
	var result map[string]interface{}
	err = json.Unmarshal(body, &result)
	if err != nil {
		return nil, err
	}
	// fmt.Println("output")
	// for key, val := range result {
	// 	fmt.Println(key, val)
	// }
	return result, nil
}

func listAllPrices(secretKey string) (map[string]interface{}, error) {

	url := "https://api.stripe.com/v1/prices"
	method := "GET"

	//payload := strings.NewReader("email=sunzida%40gmail.com&name=Dr%20Sanzida&phone=01765110255&source=tok_1IU9EmEJdAixg8N3xAwmFVBP")
	client := &http.Client{}
	req, err := http.NewRequest(method, url, nil)

	if err != nil {
		return nil, err
	}
	//req.Header.Add("Stripe-Account", "acct_1Es5BxEJdAixg8N3")
	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", secretKey))
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	res, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}
	//fmt.Println(string(body))
	var result map[string]interface{}
	err = json.Unmarshal(body, &result)
	if err != nil {
		return nil, err
	}
	// fmt.Println("output")
	// for key, val := range result {
	// 	fmt.Println(key, val)
	// }
	return result, nil
}

func apiGetRequest(rurl, secretKey string) (map[string]interface{}, error) {

	req, err := http.NewRequest(http.MethodGet, rurl, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Add("Accept", "application/json")
	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", secretKey))
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()
	var rmap = make(map[string]interface{})
	json.NewDecoder(res.Body).Decode(&rmap)

	//fmt.Println("res.StatusCode:", res.StatusCode, res.StatusCode != http.StatusOK)
	//error handling
	if res.StatusCode != http.StatusOK {
		//fmt.Printf("%v %T\n", rmap, rmap)
		erow, isOk := rmap["error"].(map[string]interface{})
		if isOk {
			etype, isOk := erow["type"].(string)
			if isOk {
				errTypes := []string{"invalid_request_error", "idempotency_error", "card_error", "api_error"}
				if mtool.ArrayValueExist(errTypes, etype) {
					message, _ := erow["message"].(string)
					return nil, errors.New(message)
				}
			}
		}
	}
	return rmap, nil
}

func apiPostRequest(rurl, secretKey string, fmap map[string]string) (map[string]interface{}, error) {

	var fv = make(url.Values)
	for key, val := range fmap {
		fv.Set(key, val)
	}
	payload := strings.NewReader(fv.Encode())
	req, err := http.NewRequest("POST", rurl, payload)
	if err != nil {
		return nil, err
	}
	req.Header.Add("Accept", "application/json")
	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", secretKey))
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()
	var rmap = make(map[string]interface{})
	json.NewDecoder(res.Body).Decode(&rmap)

	//error handling
	if res.StatusCode != http.StatusOK {
		etype, isOk := rmap["type"].(string)
		if isOk {
			errTypes := []string{"invalid_request_error", "idempotency_error", "card_error", "api_error"}
			if mtool.ArrayValueExist(errTypes, etype) {
				message, _ := rmap["message"].(string)
				return nil, errors.New(message)
			}
		}
	}

	return rmap, nil
}

// when we create price it automatically added under a product
func createSubscriptionPrice(stripeSecretKey, productName, priceInCents string) (map[string]interface{}, error) {

	//type = recurring
	//stripeKey := "sk_test_51OjqyFJFUQv2NTJsitgDUhNX3CPbns3eE3IyxSdTc8yEhI5p24SDyn9lyEI4AqaMSRghw6V25XoStkYa8Zl7zEOg006vuF1cTQ"
	var fmap = make(map[string]string)
	fmap["currency"] = "usd"
	fmap["unit_amount"] = priceInCents
	fmap["recurring[interval]"] = "month"
	fmap["product_data[name]"] = productName //product.name
	fmap["nickname"] = productName
	return apiPostRequest("https://api.stripe.com/v1/prices", stripeSecretKey, fmap)
}

// stripe checkout session
func createSession(stripeSecretKey, docNumber, customerEmail, priceId, qty string) (map[string]interface{}, error) {

	//type = recurring
	//stripeKey := "sk_test_51OjqyFJFUQv2NTJsitgDUhNX3CPbns3eE3IyxSdTc8yEhI5p24SDyn9lyEI4AqaMSRghw6V25XoStkYa8Zl7zEOg006vuF1cTQ"
	successUrl := "https://lxroot.com/complete"
	var fmap = make(map[string]string)
	fmap["client_reference_id"] = docNumber

	fmap["line_items[0][price]"] = priceId
	fmap["line_items[0][quantity]"] = qty
	fmap["line_items[1][price]"] = "price_1PPddlJFUQv2NTJs4sxm013J"
	fmap["line_items[1][quantity]"] = "2"

	fmap["success_url"] = successUrl
	fmap["customer_email"] = customerEmail
	fmap["mode"] = "subscription" //payment,setup,subscription
	//fmap["customer_creation"]="always" //payment
	//set payment_intent_data.setup_future_usage to have Checkout automatically
	return apiPostRequest("https://api.stripe.com/v1/checkout/sessions", stripeSecretKey, fmap)
}
