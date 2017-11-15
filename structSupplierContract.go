package main

type SupplierContract struct {
	OrderHeader
	EffectiveDate string            `json:"effectivedate"`
	CSAID         string            `json:"csaid"`
	CouncilName   string            `json:"councilname"`
	LaborType     string            `json:"labortype"`
	RequestorName string            `json:"requestorname"`
	RequestType   string            `json:"requesttype"`
	ClientName    string            `json:"clientname"`
	Items         []Item            `json:"items"`
	Contractors   map[string]string `json:"contractors"`
	Status        string            `json:"status"`
	Description   string            `json:"description"`
}

type OrderHeader struct {
	ObjectType        string `json:"doctype"`
	PONumber          string `json:"ponumber"`
	SupplierID        string `json:"supplierid"`
	VendorDescription string `json:"vendordescription"`
	POIssueDate       string `json:"poissuedate"`
	POStartDate       string `json:"postartdate"`
	POEndDate         string `json:"poenddate"`
	PaymentTerm       string `json:"paymentterm"`
	PaymentDays       string `json:"paymentdays"`
	SowNumber         string `json:"sownum"`
	Currency          string `json:"currency"`
}

type Item struct {
	UUID         string         `json:"uuid"`
	ShortText    string         `json:"shorttext"`
	JobRoleID    string         `json:"jobroleid"`
	JobRoleTitle string         `json:"jobroletitle"`
	SkillID      string         `json:"skillid"`
	SkillText    string         `json:"skilltitle"`
	StartDate    string         `json:"startdate"`
	EndDate      string         `json:"enddate"`
	StraightTime ItemDetails    `json:"straighttime"`
	OverTime     ItemDetails    `json:"overtime"`
	Expense      ItemDetails    `json:"expense"`
	Address      AddressDetails `json:"address"`
}

type ItemDetails struct {
	CommodityCode         string `json:"commoditycode"`
	EffectiveDate         string `json:"effectivedate"`
	UnitPrice             string `json:"unitprice"`
	OldPrice              string `json:"oldprice"`
	UOM                   string `json:"uom"`
	ShortText             string `json:"shorttext"`
	OrderedQuantity       string `json:"orderedquantity"`
	OrderedValue          string `json:"orderedvalue"`
	InvoicedQuantity      string `json:"invoicedquantity"`
	InvoicedValue         string `json:"invoicedvalue"`
	ToBeDeliveredValue    string `json:"tobedeliveredvalue"`
	ToBeDeliveredQuantity string `json:"tobedeliveredquantity"`
}

type AddressDetails struct {
	Name     string `json:"name"`
	Street   string `json:"street"`
	Postcode string `json:"postcode"`
	City     string `json:"city"`
	State    string `json:"state"`
	Country  string `json:"country"`
}
