package main

//stripe balance object
type Balance struct {
	Object    string                   `json:"object"`
	Available []map[string]interface{} `json:"available"`
	Livemode  bool                     `json:"livemode"`
	Pending   []map[string]interface{} `json:"pending"`
}

//stripe event object https://stripe.com/docs/api/events/object
type Event struct {
	ID              string                 `json:"id"`
	ApiVersion      string                 `json:"api_version"`
	Data            map[string]interface{} `json:"data"`
	Request         map[string]interface{} `json:"request"`
	Type            string                 `json:"type"` //?? invoice.created | charge.refunded
	Livemode        bool                   `json:"livemode"`
	Created         int64                  `json:"created"`
	Object          string                 `json:"object"`
	PendingWebhooks int                    `json:"pending_webhooks"`
}

type BillingAddress struct {
	City       string `json:"city"`
	Country    string `json:"country"`
	Line1      string `json:"line1"`
	Line2      string `json:"line2"`
	PostalCode string `json:"postal_code"`
	State      string `json:"state"`
}

type Customer struct {
	ID            string         `json:"id"`
	Name          string         `json:"name"`
	Phone         string         `json:"phone"`
	Shipping      string         `json:"shipping"`
	Address       BillingAddress `json:"address"`
	Description   string         `json:"description"`
	Email         string         `json:"email"`
	Currency      string         `json:"currency"`
	Created       int64          `json:"created"`
	InvoicePrefix string         `json:"invoice_prefix"`
}

type List struct {
	Object  string                   `json:"object"`
	Data    []map[string]interface{} `json:"data"`
	HasMore bool                     `json:"has_more"`
	Url     string                   `json:"url"`
}

type sInvoice struct {
	Object  string    `json:"object"`
	Data    []Invoice `json:"data"`
	HasMore bool      `json:"has_more"`
	Url     string    `json:"url"`
}

type Invoice struct {
	ID                  string         `json:"id"`
	Number              string         `json:"number"`
	ReceiptNumber       string         `json:"receipt_number"`
	Paid                bool           `json:"paid"`
	PaymentIntent       string         `json:"payment_intent"`
	AccountCountry      string         `json:"account_country"`
	AccountName         string         `json:"account_name"`
	Customer            string         `json:"customer"`
	CustomerAddress     BillingAddress `json:"customer_address"`
	CustomerEmail       string         `json:"customer_email"`
	CustomerName        string         `json:"customer_name"`
	CustomerPhone       string         `json:"customer_phone"`
	Created             int64          `json:"created"`
	Subscription2       string         `json:"subscription"`
	PeriodStart         int64          `json:"period_start"`
	PeriodEnd           int64          `json:"period_end"`
	Charge              string         `json:"charge"`
	Currency            string         `json:"currency"`
	CollectionMethod    string         `json:"collection_method"`
	AmountPaid          int64          `json:"amount_paid"` //The amount, in cents, that was paid.
	AmountDue           int64          `json:"amount_due"`
	AmountRemaining     int64          `json:"amount_remaining"` //The amount remaining, in cents, that is due.
	BillingReason       string         `json:"billing_reason"`
	Total               int64          `json:"total"`
	HostedInvoiceUrl    string         `json:"hosted_invoice_url"`
	InvoicePdf          string         `json:"invoice_pdf"`
	WebhooksDeliveredAt int64          `json:"webhooks_delivered_at"`
	Status              string         `json:"status"`
}

//https://stripe.com/docs/api/payment_intents/object
//payment_intent
//charge
type PaymentIntent struct {
	ID                 string   `json:"id"`
	Amount             int64    `json:"amount"`
	AmountReceived     int64    `json:"amount_received"`
	ClientSecret       string   `json:"client_secret"`
	Currency           string   `json:"currency"`
	Customer           string   `json:"customer"`
	Description        string   `json:"description"`
	Charges            List     `json:"charges"`
	Invoice            string   `json:"invoice"`
	PaymentMethod      string   `json:"payment_method"`
	PaymentMethodTypes []string `json:"payment_method_types"`
	ReceiptEmail       string   `json:"receipt_email"`
	Created            int64    `json:"created"`
	Status             string   `json:"status"`
}

type Subscription2 struct {
	ID                     string `json:"id"`
	Cancel_at_period_end   bool   `json:"cancel_at_period_end"`
	Current_period_start   int64  `json:"current_period_start"`
	Current_period_end     int64  `json:"current_period_end"`
	Canceled_at            int64  `json:"canceled_at"`
	Customer               string `json:"customer"`
	Default_payment_method string `json:"default_payment_method"`
	Latest_invoice         string `json:"latest_invoice"`
	Start_date             int64  `json:"start_date"`
	Ended_at               int64  `json:"ended_at"`
	Collection_method      string `json:"collection_method"`
	Created                int64  `json:"created"`
	Livemode               bool   `json:"livemode"`
	Items                  List   `json:"items"`
	Status                 string `json:"status"`
}

type CheckoutCustomerDetails struct {
	Email      string `json:"email"`
	Phone      string `json:"phone"`
	Tax_exempt string `json:"tax_exempt"`
}

//checkout.session
type Session struct {
	ID               string                  `json:"id"`
	ClentReferenceId string                  `json:"client_reference_id"`
	CancelUrl        string                  `json:"cancel_url"`
	Currency         string                  `json:"currency"`
	Customer         string                  `json:"customer"`
	Invoice          string                  `json:"invoice"`
	CustomerDetails  CheckoutCustomerDetails `json:"customer_details"` //map[string]interface{}
	CustomerEmail    string                  `json:"customer_email"`
	PaymentIntent    string                  `json:"payment_intent"`
	PaymentStatus    string                  `json:"payment_status"`
	SuccessUrl       string                  `json:"success_url"`
	Mode             string                  `json:"mode"` //subscription
	AmountTotal      int64                   `json:"amount_total"`
	Subscription2    string                  `json:"subscription"`
	Status           string                  `json:"status"` //open,complete,expired
}

type BillingDetails struct {
	Address BillingAddress `json:"address"`
	Email   string         `json:"email"`
	Name    string         `json:"name"`
	Phone   string         `json:"phone"`
}

type Charge struct {
	ID                   string                 `json:"id"`
	Amount               int64                  `json:"amount"`
	AmountCaptured       int64                  `json:"amount_captured"`
	AmountRefunded       int64                  `json:"amount_refunded"`
	BalanceTransaction   string                 `json:"balance_transaction"`
	BillingDetails       BillingDetails         `json:"billing_details"`
	Invoice              string                 `json:"invoice"`
	Currency             string                 `json:"currency"`
	Customer             string                 `json:"customer"`
	Description          string                 `json:"description"`
	Paid                 bool                   `json:"paid"`
	Refunded             bool                   `json:"refunded"`
	PaymentIntent        string                 `json:"payment_intent"`
	PaymentMethod        string                 `json:"payment_method"`
	PaymentMethodDetials map[string]interface{} `json:"payment_method_details"`
	ReceiptUrl           string                 `json:"receipt_url"`
	Created              int64                  `json:"created"`
	Status               string                 `json:"status"`
}
