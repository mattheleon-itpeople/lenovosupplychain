package main

type Part struct {
	PartNumber string    `json:"partNumber"`
	ManufID    string    `json:"manufId"`
	RCType     string    `json:"rcType"`
	AuditInfo  AuditInfo `json:"auditInfo"`
}

//TODO
/*
type OrderItem struct {
	PartNumber     string  `json:"partNumber"`
	ItemCondition  string  `json:"itemCondition"` //New or Refurbished
	Quantity       int     `json:"quantity"`      //add unit of measurement
	PricePerUnit   float64 `json:"pricePerUnit"`
	UnitOfMeasure  string  `json:"unitOfMeasure"`
	TotalLinePrice float64 `json:"totalLinePrice"`
}*/

type OrderLine struct {
	Items     []ItemDetails `json:"items"`
	AuditInfo AuditInfo     `json:"auditInfo"`
	From      string        `json:"from"`
	To        string        `json:"to"`
}

type PurchaseOrder struct {
	OrderHeader
	OrderLine
}
