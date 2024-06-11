package main

//Company #1 wise database <> META_INFO
type Company struct {
	ID          string `json:"id"`
	Type        string `json:"type"`
	CompanyName string `json:"company_name,omitempty"`
	Website     string `json:"website,omitempty"`
	Status      int    `json:"status"`
}

//Access #2 login type admin,super,user,etc <> META_INFO
type Access struct {
	ID         string `json:"id"`
	Type       string `json:"type"`
	CompanyID  string `json:"cid"` //foreign key
	AccessName string `json:"access_name"`
	Status     int    `json:"status"`
}

//<> META_INFO
type Country struct {
	ID             string `json:"id"`
	Type           string `json:"type"`
	CompanyID      string `json:"cid"`                     //foreign key
	Name           string `json:"name"`                    //country name
	IsoCode        string `json:"iso_code"`                //2 char = bd
	CountryCode    string `json:"country_code"`            // +880
	CurrencyName   string `json:"currency_name,omitempty"` //BDT
	CurrencySymbol string `json:"currency_symbol"`         //TK
	Status         int    `json:"status"`
}

//Account #3 ::cid::account_id
type Account struct {
	//ContactInfo []Contact `json:"contact_info,omitempty"`
	ID          string `json:"id"`                     //system auto generated
	Type        string `json:"type"`                   //account
	CompanyID   string `json:"cid"`                    //foreign key
	ParentID    string `json:"parent_id"`              //if any parent
	Photo       string `json:"photo,omitempty"`        //account owner photo
	RateplanID  string `json:"rateplan_id,omitempty"`  //
	AccountType string `json:"account_type"`           //vendor,customer,goods_supplier,consumer,payment_provider,shipping_provider
	AccountName string `json:"account_name"`           //supplier business name or username
	CustomID    string `json:"customid,omitempty"`     //unique customer IDENTIFICATION
	Code        string `json:"code"`                   //Ledger | supplier | customer code
	FirstName   string `json:"first_name"`             //
	LastName    string `json:"last_name"`              //
	DateOfBirth string `json:"dob,omitempty"`          //
	Gender      string `json:"gender,omitempty"`       //female,male,other
	Phone       string `json:"phone"`                  //phone
	Email       string `json:"email"`                  //
	ReferralUrl string `json:"referral_url,omitempty"` //referral_url
	Industry    string `json:"industry"`               //industry
	Remarks     string `json:"remarks"`                //remarks
	CreateDate  string `json:"create_date"`
	UpdateDate  string `json:"update_date"`
	Status      int    `json:"status"`
}

//Address ...
type Address struct {
	ID          string `json:"id"`
	Type        string `json:"type"`
	CompanyID   string `json:"cid"`          //foreign key
	AccountID   string `json:"account_id"`   //foreign_key
	AddressType string `json:"address_type"` //billing,shipping
	Country     string `json:"country"`      //iso_code
	State       string `json:"state"`
	City        string `json:"city"`
	Address1    string `json:"address1"`
	Address2    string `json:"address2"`
	Zip         string `json:"zip"`
	Status      int    `json:"status"`
}

//CustomField ...
// type CustomField struct {
// 	ID         string `json:"id"`
// 	Type       string `json:"type"` //account
// 	CompanyID  string `json:"cid"`  //foreign key
// 	AccountID  string `json:"account_id"`
// 	Owner      string `json:"owner"`    //account
// 	FieldID    string `json:"field_id"` //extra
// 	FieldLabel string `json:"field_label"`
// 	FieldName  string `json:"field_name"`
// 	Status     int    `json:"status"`
// }

// //CustomFieldValue ...
// type CustomFieldValue struct {
// 	ID         string `json:"id"`
// 	Type       string `json:"type"`        //account
// 	CompanyID  string `json:"cid"`         //foreign key
// 	FieldID    string `json:"field_id"`    //CustomFieldid
// 	FieldValue string `json:"field_value"` //account
// 	Status     int    `json:"status"`
// }

//Contact info for account and user
// type Contact struct {
// 	ID          string `json:"id"`
// 	Type        string `json:"type"`       //table_name
// 	CompanyID   string `json:"cid"`        //foreign key
// 	ContactID   int64  `json:"contact_id"` //mobile,email,pager,phone,website,socialmedia
// 	ContactType string `json:"contact_type"`
// 	ContactData string `json:"contact_date"`
// 	Owner       string `json:"owner"`     ////owner table info, who owns this data
// 	Parameter   string `json:"parameter"` //docID
// 	Status      int    `json:"status"`
// }

//Login #4  all user account login table
type Login struct {
	ID          string `json:"id"`
	Type        string `json:"type"`
	CompanyID   string `json:"cid"`          //foreign key
	AccountID   string `json:"account_id"`   //foreign key
	AccessID    string `json:"access_id"`    //foreign key
	AccessName  string `json:"access_name"`  //customer type
	UserName    string `json:"username"`     //email or mobile as username
	Password    string `json:"passw"`        //password
	TfaStatus   int    `json:"tfa_status"`   //TFA = 0,1
	TfaMedium   string `json:"tfa_medium"`   //TFA
	TfaSetupkey string `json:"tfa_setupkey"` //TFA
	CreateDate  string `json:"create_date"`
	LastLogin   string `json:"last_login,omitempty"` //update date
	Status      int    `json:"status"`
}

//ActivityLog ***
type ActivityLog struct {
	ID           string `json:"id"`
	Type         string `json:"type"`
	CompanyID    string `json:"cid"` //foreign key
	ActivityType string `json:"activity_type"`
	OwnerTable   string `json:"table_name"` //owner_table
	Parameter    string `json:"parameter"`  //key=val or id
	LogDetails   string `json:"log_details"`
	IPAddress    string `json:"ip_address"`
	LoginID      string `json:"login_id"` //optional foreign key
	CreateDate   string `json:"create_date"`
	Status       int    `json:"status"`
}

//VisitorSession info
type VisitorSession struct {
	//Platform       string `json:"platform"`
	ID             string `json:"id"`
	Type           string `json:"type"`
	CompanyID      string `json:"cid"` //foreign key
	SessionCode    string `json:"session_code"`
	Device         string `json:"device"`
	ScreenSize     string `json:"screen_size"`
	BrowserVersion string `json:"browser_version"`
	OsVersion      string `json:"os_version"`
	IPAddress      string `json:"ip_address"`
	GeoLocation    string `json:"geo_location"` //
	City           string `json:"city"`
	Country        string `json:"country"`
	VisitorCount   int    `json:"vcount"`
	CreateDate     string `json:"create_date"`
	UpdateDate     string `json:"update_date"`
	Status         int    `json:"status"`
}

//DeviceLog tracks user device info (where they login to the system)
// type DeviceLog struct {
// 	ID          string `json:"id"`
// 	Type        string `json:"type"`
// 	CompanyID   string `json:"cid"`      //foreign key
// 	LoginID     string `json:"login_id"` //foreign key
// 	Browser     string `json:"browser"`
// 	DeviceType  string `json:"device_type"`
// 	Os          string `json:"os"`
// 	Platform    string `json:"platform"`
// 	ScreenSize  string `json:"screen_size"`
// 	GeoLocation string `json:"geolocation,omitempty"`
// 	CreateDate  string `json:"create_date"`
// 	Status      int    `json:"status"`
// }

//LoginSession keeps user login session for 24 hours or more
type LoginSession struct {
	//DeviceInfo  string `json:"device_log"` //foreign key device_log::1
	ID          string `json:"id"`
	Type        string `json:"type"`
	CompanyID   string `json:"cid"` //foreign key
	SessionCode string `json:"session_code"`
	LoginID     string `json:"login_id"` //foreign key
	IPAddress   string `json:"ip_address"`
	IPCity      string `json:"city"`
	IPCountry   string `json:"country"`
	UserAgent   string `json:"user_agent"`
	LoginTime   string `json:"login_time"`  //
	LogoutTime  string `json:"logout_time"` //
	CreateDate  string `json:"create_date"`
	Status      int    `json:"status"`
}

type Authc struct {
	ID         string `json:"id"`
	Type       string `json:"type"`
	CompanyID  string `json:"cid"`      //foreign key
	LoginID    string `json:"login_id"` //foreign key
	Token      string `json:"token"`
	IpAddress  string `json:"ip_address"`
	CreateDate string `json:"create_date"`
	ExpireDate string `json:"expire_date"`
	UpdateDate string `json:"update_date"`
	Status     int    `json:"status"`
}

//Setting table keep only one file <> META_INFO
type Settings struct {
	ID         string `json:"id"`
	Type       string `json:"type"`
	CompanyID  string `json:"cid"` //foreign key
	FieldName  string `json:"field_name"`
	FieldValue string `json:"field_value"`
	Purpose    string `json:"purpose"`
	Status     int    `json:"status"`
}

//Message table
// type Message struct {
// 	ID          string `json:"id"`
// 	Type        string `json:"type"`
// 	CompanyID   string `json:"cid"`          //foreign key
// 	MessageType string `json:"message_type"` //contactus,inapp,email,sms,verify
// 	Sender      string `json:"sender"`       //sender = login.username
// 	Receiver    string `json:"receiver"`     //receiver = login.username
// 	Subject     string `json:"subject"`
// 	MessageBody string `json:"message_body"`
// 	CreateDate  string `json:"create_date"`
// 	Status      int    `json:"status"`
// }

//Verification message
type Verification struct {
	ID                  string `json:"id"`
	Type                string `json:"type"`
	CompanyID           string `json:"cid"`                  //foreign key
	MessageID           string `json:"message_id,omitempty"` //foreign key
	Username            string `json:"username"`             //login.username
	VerificationPurpose string `json:"verification_purpose"` //TFA|SIGNUP
	VerificationCOde    string `json:"verification_code"`
	CreateDate          string `json:"create_date"`
	UpdateDate          string `json:"update_date"`
	Status              int    `json:"status"`
}

//AccountHead+AccountGroup merged into achead Achead <> META_INFO
type Achead struct {
	ID                   string  `json:"id"`            //unique id*
	Type                 string  `json:"type"`          //table
	CompanyID            string  `json:"cid"`           //foreign key
	AccountGroup         string  `json:"account_group"` //*AccountGroup= Asset|Liability|Equity|Revenue|Expense
	AccountType          string  `json:"account_type"`  //* group|head
	Name                 string  `json:"name"`          //ledger name *
	Description          string  `json:"description"`   //*ledger_details
	Identifier           string  `json:"identifier"`    //* for ensuring no ledgers are duplicate
	LedgerCode           string  `json:"code"`          //ledger number *
	ParentID             string  `json:"parent_id"`     //parent account *
	CurrentBalance       float64 `json:"balance"`       //ledger balance *
	CurrentBalanceType   string  `json:"baltype"`       //ledger balance type Dr or Cr *
	Restricted           int     `json:"restricted"`    //1=Yes, No=0 *
	CostCenterApplicable int     `json:"cost_center"`   //1=Yes, No=0 *
	Remarks              string  `json:"remarks"`       //*
	CreateDate           string  `json:"create_date"`   //insert date *
	Status               int     `json:"status"`        //0=Inactive, 1=Active, 9=Deleted *
}

//Warehouse table xlsx file 07
// type Warehouse struct {
// 	ID               string `json:"id"`
// 	Type             string `json:"type"`
// 	CompanyID        string `json:"cid"`            //foreign key
// 	WarehouseName    string `json:"warehouse_name"` //warehouse_name
// 	WarehouseDetails string `json:"details"`        //warehouse_details
// 	IsDefault        bool   `json:"isdefault"`      //true|false == 0|1
// 	Status           int    `json:"status"`
// }

//UOM = Unit of Measurement
// type UOM struct {
// 	ID         string `json:"id"`
// 	Type       string `json:"type"`
// 	CompanyID  string `json:"cid"` //foreign key
// 	UnitName   string `json:"unit_name"`
// 	UnitSymbol string `json:"symbol"`
// 	Status     int    `json:"status"`
// }

//Tax .. <> META_INFO
type Tax struct {
	ID                string  `json:"id"`
	Type              string  `json:"type"`
	CompanyID         string  `json:"cid"`            //foreign key
	Name              string  `json:"name"`           //tax name
	DisplayName       string  `json:"display_name"`   //display title
	TaxRegNumber      string  `json:"tax_reg_number"` //taxid provided by govt
	TaxType           string  `json:"tax_type"`       //vat | tax | income_tax
	Rate              float64 `json:"rate"`           //tax rate applicable in percentage
	AccountHeaderCode string  `json:"account_code"`   //ledger number attached to it
	Remarks           string  `json:"remarks"`        //remarks=default auto selected from item creation
	Status            int     `json:"status"`
}

//Department itemDepartment -> ItemLine or department are same just one under another
//TableMap::-> Deprtment & ItemLine <> META_INFO
type Department struct {
	ID          string `json:"id"`
	Type        string `json:"type"`
	CompanyID   string `json:"cid"`         //foreign key
	ParentID    string `json:"parent_id"`   //same table = parent_id="" parent department id
	Name        string `json:"name"`        //child department
	Code        string `json:"code"`        //unique code
	Description string `json:"description"` //default = when add from item page
	Owner       string `json:"owner"`       //item | item_line | office | ticket
	Status      int    `json:"status"`
}

//Category or ItemCategory/ItemGroup keeps product group/category info
//TableMap::-> ItemGroup > website /shop <> META_INFO
type Category struct {
	ID            string `json:"id"`
	Type          string `json:"type"`
	CompanyID     string `json:"cid"`           //foreign key
	DepartmentID  string `json:"department_id"` //foreign key itemLine
	Name          string `json:"name"`
	Description   string `json:"description"`
	Code          string `json:"code"`
	Position      int64  `json:"position"`
	CategoryImage string `json:"image"`
	Owner         string `json:"owner"` //item / product | service
	Status        int    `json:"status"`
}

//ItemReward ... invoice qty 1 == 50 points
// type ItemReward struct {
// 	ID          string `json:"id"`
// 	Type        string `json:"type"`         //item_reward
// 	CompanyID   string `json:"cid"`          //foreign key
// 	ItemID      string `json:"item_id"`      //foreign key
// 	RewardField string `json:"reward_field"` //qty,amount,amountPercentAsPoint, ProfitSharePercent
// 	RewardValue int64  `json:"reward_value"` //10 pcs
// 	RewardPoint int64  `json:"reward_point"` //50 points
// 	Status      int    `json:"status"`
// }

//Item keeps product/Item info
type Item struct {
	ID              string `json:"id"`
	Type            string `json:"type"`
	CompanyID       string `json:"cid"`         //foreign key
	ItemCode        string `json:"item_code"`   //stripe.priceId = price_1PPlulJFUQv2NTJsqGsPFpLa
	CategoryID      string `json:"category_id"` //foreign key as item_group
	ItemType        string `json:"item_type"`   //raw_material,stockable,consumable,service|license_key
	ItemName        string `json:"item_name"`
	ItemDescription string `json:"item_description"`
	//ItemURL                     string  `json:"item_url"`
	//ItemLine                    string  `json:"item_line,omitempty"` //*** put extra to avoid complexity in join
	ItemImage string `json:"item_image"`
	//Barcode                     string  `json:"barcode"`
	InventoryAccount            string  `json:"inventory_account"`              //Asset / inventory ledger
	COGSAccount                 string  `json:"cogs_account"`                   //Expense = ledger number
	SalesAccount                string  `json:"sales_account"`                  //Income = ledger number
	OpeningBalanceEquityAccount string  `json:"opening_balance_equity_account"` //AccountHead = ledger number
	VatID                       string  `json:"tax_id"`                         //foreign key
	VatPercent                  float64 `json:"vat"`                            //vat percent
	BuyPrice                    float64 `json:"buy_price"`                      //cost price, Trade Price
	SalePrice                   float64 `json:"sale_price"`                     //MRP
	Tags                        string  `json:"tags,omitempty"`                 //?? department
	SupplierID                  string  `json:"supplier,omitempty"`             //AccountTableID*account_type = supplier
	UnitOfMeasure               string  `json:"uom,omitempty"`                  //1 unit = license
	TrackingBy                  string  `json:"tracking"`                       //tracking by license_key, unique_serial, lot_number, no_tracking
	//ReorderLevel                int64   `json:"reorder_level"`
	//ReorderQty                  int64   `json:"reorder_qty"`
	//BatchNo                     string  `json:"batch"`
	//ExpireDate                  string  `json:"expire_date"`
	//StockQty                    int64   `json:"stock_qty"`
	//PublishOnWebsite            int     `json:"publish_on_web"` //published | unpublished
	//DisplayOnSales              int     `json:"display_sales"`
	//DisplayOnPurchase           int     `json:"display_purchase"`
	Availability string `json:"availability,omitempty"` //coming soon | available | for website
	Status       int    `json:"status"`
}

//ItemAttribute ..
// type ItemAttribute struct {
// 	ID           string `json:"id"`
// 	Type         string `json:"type"`
// 	CompanyID    string `json:"cid"`           //foreign key
// 	ItemID       string `json:"item_id"`       //item_id foreign key
// 	AttributeKey string `json:"attribute_key"` //attr key
// 	Position     int64  `json:"position"`      //serial / position
// 	KeyType      string `json:"key_type"`      //select, text, radio, color
// 	DefaultValue string `json:"default_value"`
// 	Status       int    `json:"status"`
// }

// //ItemAttributeValue ..
// type ItemAttributeValue struct {
// 	ID              string `json:"id"`
// 	Type            string `json:"type"`
// 	CompanyID       string `json:"cid"`             //foreign key
// 	ItemID          string `json:"item_id"`         //foreign key item_id
// 	ItemAttributeID string `json:"attribute_id"`    //foreign key ItemAttribute
// 	AttributeValue  string `json:"attribute_value"` //select, text, radio, color
// 	Position        int64  `json:"position"`        //serial / position
// 	Status          int    `json:"status"`
// }

//Rateplan ...
// type Rateplan struct {
// 	ID         string `json:"id"`
// 	Type       string `json:"type"` //type=rateplan
// 	CompanyID  string `json:"cid"`
// 	CustomerID string `json:"account_id"` //AccountID
// 	Name       string `json:"name"`
// 	Owner      string `json:"owner"`
// 	Status     string `json:"status"`
// }

// //Rate ...
// type Rate struct {
// 	ID               string `json:"id"`
// 	Type             string `json:"type"` //type=rate
// 	CompanyID        string `json:"cid"`
// 	RateplanID       string `json:"rateplan_id"`
// 	ItemID           string `json:"item_id"`
// 	ItemLoyaltyPoint int64  `json:"loyalty_point"` //bonus point
// 	Rate             string `json:"rate"`
// 	Rebate           string `json:"rebate"`
// 	Charge           string `json:"charge"`
// 	Remarks          string `json:"remarks"`
// 	Status           string `json:"status"`
// }

//FileStore general table ***
type FileStore struct {
	ID         string `json:"id"`
	Type       string `json:"type"`        //table
	CompanyID  string `json:"cid"`         //company
	Reference  string `json:"reference"`   //foreign key tableOwnerId
	OwnerTable string `json:"owner_table"` //table_name
	FileType   string `json:"file_type"`   //jpeg,png,pdf
	Filepath   string `json:"filepath"`    //file_location -> filepath
	Remarks    string `json:"remarks"`
	Status     int    `json:"status"`
}

//DocKeeper keeps all document info
type DocKeeper struct {
	ID        string `json:"id"`
	Type      string `json:"type"`
	CompanyID string `json:"cid"` //foreign key
	//WarehouseID    string  `json:"warehouse_id,omitempty"` //foreign key
	DocName        string  `json:"doc_name"`        //Shopping Cart
	DocType        string  `json:"doc_type"`        //cart,purchase,sales,voucher
	DocRef         string  `json:"doc_ref"`         //visitorSession or stripe reference
	DocNumber      string  `json:"doc_number"`      //auto generated
	DocDescription string  `json:"doc_description"` //stripe > session.id
	PostingDate    string  `json:"posting_date"`
	LoginID        string  `json:"login_id"`       //foreign key
	AccountID      string  `json:"account_id"`     //foreign key
	TaxRule        string  `json:"tax_rule"`       //*** foreign key tax.id
	TotalTax       float64 `json:"total_tax"`      //***
	TotalDiscount  float64 `json:"total_discount"` //*** may be applied discount_code
	TotalPayable   float64 `json:"total_payable"`  //***
	PaymentStatus  string  `json:"payment_status"` //paid,unpaid,no_payment_required
	ReceiptUrl     string  `json:"receipt_url"`    //
	DocStatus      string  `json:"doc_status"`     //pending|open,checkout_session -> open,complete,expired
	CreateDate     string  `json:"create_date"`
	Status         int     `json:"status"`
}

//TransactionRecord keeps all transaction info like sales, purchase, online oreder
type TransactionRecord struct {
	ID             string  `json:"id"`
	Type           string  `json:"type"`
	CompanyID      string  `json:"cid"`        //foreign key
	TrxType        string  `json:"trx_type"`   // cart/purchase/sales/purchase_return/sales_return/rejected
	DocNumber      string  `json:"doc_number"` //foriegn key
	ItemID         string  `json:"item_id"`
	ItemInfo       string  `json:"item_info"`
	StockInfo      string  `json:"stock_info"`  //stripePriceId => price_1PPlulJFUQv2NTJsqGsPFpLa
	ProductSerial  string  `json:"item_serial"` //SKU or barcode | license_key
	Quantity       int     `json:"quantity"`
	Rate           float64 `json:"rate"`  //per unit price, GAAP Compliance
	Price          float64 `json:"price"` //price = qty x rate'
	DiscountRate   float64 `json:"discount_rate"`
	TaxRate        float64 `json:"tax_rate"`
	DiscountAmount float64 `json:"discount_amount"` //(price x discount_rate)/100
	TaxableAmount  float64 `json:"taxable_amount"`  //price-total_discount
	TaxAmount      float64 `json:"tax_amount"`
	PayableAmount  float64 `json:"payable_amount"` //taxable_amount + total_tax | price-total_discount+total_tax
	CreateDate     string  `json:"create_date"`
	Status         int     `json:"status"` //0=Inactive, 1=Active, 9=Deleted
}

//StockMovement keeps stock transfer from warehouse to warehouse or other
// type StockMovement struct {
// 	ID            string `json:"id"`
// 	Type          string `json:"type"`
// 	CompanyID     string `json:"cid"`        //foreign key
// 	DocNumber     string `json:"doc_number"` //foriegn key
// 	StockNote     string `json:"stock_note"`
// 	ItemID        string `json:"item_id"`      //foreign key
// 	ProductSerial string `json:"item_serial"`  //SKU or barcode
// 	WarehouseID   string `json:"warehouse_id"` //foreign key
// 	MovementType  string `json:"mtype"`        //movement type = IN/OUT
// 	Quantity      int    `json:"quantity"`
// 	StockBalance  int    `json:"stock_balance"`
// 	Status        int    `json:"status"` //0=Inactive, 1=Active, 9=Deleted
// }

//LedgerTransaction table stores accounting transaction between ledgers/accounting_header
//Debit | Credit note
type LedgerTransaction struct {
	ID           string  `json:"id"`
	Type         string  `json:"type"`
	CompanyID    string  `json:"cid"`          //foreign key
	DocNumber    string  `json:"doc_number"`   //foriegn key voucherNo
	Particulars  string  `json:"particulars"`  //opposite ledger
	AcheadID     string  `json:"acid"`         //achead.ID > account_id
	VoucherName  string  `json:"voucher_name"` //
	LedgerNumber string  `json:"ledger"`       //ledger code on AccountHeader
	LedgerName   string  `json:"ledger_name"`
	Description  string  `json:"description"`
	Debit        float64 `json:"debit"`
	Credit       float64 `json:"credit"`
	Balance      float64 `json:"balance"`
	BalanceType  string  `json:"baltype"` //balance type Dr/Cr/Eq
	CreateDate   string  `json:"create_date"`
	Source       string  `json:"source"` //monthly_checklist::1
	Status       int     `json:"status"` //0=Inactive, 1=Active, 9=Deleted
}

//ShippingAddress ...
// type ShippingAddress struct {
// 	ID               string `json:"id"`
// 	Type             string `json:"type"`
// 	CompanyID        string `json:"cid"`        //foreign key
// 	DocNumber        string `json:"doc_number"` //foriegn key
// 	ReciverFirstName string `json:"first_name"`
// 	ReciverLastName  string `json:"last_name"`
// 	CompanyName      string `json:"company_name"`
// 	Phone            string `json:"phone"`
// 	Email            string `json:"email"`
// 	Address1         string `json:"address1"`
// 	Address2         string `json:"address2"`
// 	City             string `json:"city"`
// 	ZipCode          string `json:"zip"`
// 	Country          string `json:"country"` //foreign key
// 	Remarks          string `json:"remarks"`
// 	Status           int    `json:"status"` //0=Inactive, 1=Active, 9=Deleted
// }

//PaymentOption ...META_INFO
type PaymentOption struct {
	ID               string `json:"id"`
	Type             string `json:"type"`
	CompanyID        string `json:"cid"` //foreign key
	Position         int    `json:"position"`
	IsDefault        int    `json:"isdefault"`
	ProviderID       string `json:"account_id"` //foriegn key, payment provider Accountid
	OptionName       string `json:"name"`       //unique value > Stripe, Paypal
	OptionDetails    string `json:"details"`
	ChargeApplicable int    `json:"charge_applicable"`
	ChargeType       string `json:"charge_type"`
	ChargeAmount     string `json:"charge_amount"`
	LedgerNumber     string `json:"ledger"` //ledger AccountHead
	Status           int    `json:"status"` //0=Inactive, 1=Active, 9=Deleted
}

//ShippingOption ...
// type ShippingOption struct {
// 	ID              string  `json:"id"`
// 	Type            string  `json:"type"`
// 	CompanyID       string  `json:"cid"` //foreign key
// 	Position        int     `json:"position"`
// 	IsDefault       int     `json:"isdefault"`
// 	SupplierID      string  `json:"account_id"` //*** Account Table id
// 	ProviderName    string  `json:"name"`       //account_name
// 	ProviderDetails string  `json:"details"`
// 	DeliveryTime    int     `json:"delivery_time"` //'how many days it takes to deliver',
// 	ChargeAmount    float64 `json:"charge_amount"` //Shhipping charge == 'handling fees',
// 	LedgerNumber    string  `json:"ledger"`        //*** ledger AccountHead
// 	Remarks         string  `json:"remarks"`
// 	Status          int     `json:"status"` //0=Inactive, 1=Active, 9=Deleted
// }

//DocPayShipInfo ...
type DocPayShipInfo struct {
	ID            string  `json:"id"`
	Type          string  `json:"type"`
	CompanyID     string  `json:"cid"`            //foreign key
	DocNumber     string  `json:"doc_number"`     //foriegn key
	PaymentMethod string  `json:"payment_method"` //foreign key -> PaymentOption.ID
	PaymentCharge float64 `json:"payment_charge"` //extra payment fees
	PaymentFrom   string  `json:"payment_from"`   //client email
	PaymentTrxID  string  `json:"payment_trxid"`  //stripe.payment_intent
	ReceiptURL    string  `json:"receipt_url"`    //stripe invoice url
	// ShippingMethod    string  `json:"shipping_method"` //foriegn key == ShippingOption id
	// ShippingCharge    float64 `json:"shipping_charge"`
	// ShippingAddressid string  `json:"shipping_address"` //foreign key
	// ShippingStatus    string  `json:"shipping_status"`  //shipped,processing,stuck,delivered
	// ShippingDate      string  `json:"shipping_date"`
	// DeliveryDate      string  `json:"delivery_date"`
	Remarks    string `json:"remarks"`
	CreateDate string `json:"create_date"`
	UpdateDate string `json:"update_date"`
	Status     int    `json:"status"` //0=Inactive, 1=Active, 9=Deleted

}

type WaitList struct {
	ID             string `json:"id"`
	Type           string `json:"type"`
	CompanyID      string `json:"cid"` //foreign key
	FirstName      string `json:"first_name"`
	LastName       string `json:"last_name"`
	CompanyWebsite string `json:"company_url"`
	Email          string `json:"email"`
	IpAddress      string `json:"ip_address"`
	CreateDate     string `json:"create_date"`
	Status         int    `json:"status"` //0=Inactive, 1=Active, 9=Deleted
}

type Ticket struct {
	ID             string `json:"id"`
	Type           string `json:"type"`
	CompanyID      string `json:"cid"` //foreign key
	Department     string `json:"department"`
	Subject        string `json:"subject"`
	Message        string `json:"message"`
	Referece       string `json:"reference"` //ticket_number
	LoginId        string `json:"login_id"`  //foreign key
	TicketPriority string `json:"priority"`  //High, Low, Medium
	TicketStatus   string `json:"ticket_status"`
	IpAddress      string `json:"ip_address"`
	CreateDate     string `json:"create_date"`
	UpdateDate     string `json:"update_date"` //last update
	Status         int    `json:"status"`      //0=Inactive, 1=Active, 9=Deleted
}

type TicketResponse struct {
	ID         string `json:"id"`
	Type       string `json:"type"`
	CompanyID  string `json:"cid"`        //foreign key
	TicketId   string `json:"ticket_id"`  //foreign key
	RespondBy  string `json:"respond_by"` //login_id
	Message    string `json:"message"`
	Ratings    int    `json:"ratings"` //1,2,3,4,5 = star
	IpAddress  string `json:"ip_address"`
	CreateDate string `json:"create_date"`
	Status     int    `json:"status"` //0=Inactive, 1=Active, 9=Deleted
}

type Subscription struct {
	ID                string `json:"id"`
	Type              string `json:"type"`
	CompanyID         string `json:"cid"`                //foreign key
	ItemId            string `json:"item_id"`            //foreign key
	Name              string `json:"name"`               //Lxroot License
	AccountId         string `json:"account_id"`         //foreign key
	Subscriber        string `json:"subscriber"`         //login_id
	Domain            string `json:"domain"`             //where lince used > wget -qO- lxr.sh | bash -s yourdomain.com
	Renews            string `json:"billing"`            //monthly|yearly
	Price             int64  `json:"price"`              //price
	PaymentStatus     string `json:"payment_status"`     //paid|refunded
	SubscriptionStart string `json:"subscription_start"` //purchase date by customer
	SubscriptionEnd   string `json:"subscription_end"`   //next billing date
	Createdate        string `json:"create_date"`        //this row create date
	Updatedate        string `json:"update_date"`
	Remarks           string `json:"remarks"`
	Status            int    `json:"status"` //1=active,0=pending,2=cancelled
}
