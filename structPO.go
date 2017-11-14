package main

type Part struct {
	PartNumber string    `json:"partnumber"`
	ManufID    string    `json:"manufid"`
	RCType     string    `json:"rctype"`
	AuditInfo  AuditInfo `json:"auditInfo"`
}

//TODO
type OrderItem struct {
	PartNumber     string  `json:"partNumber"`
	ItemCondition  string  `json:"itemCondition"` //New or Refurbished
	Quantity       int     `json:"quantity"`      //add unit of measurement
	PricePerUnit   float64 `json:"priceperunit"`
	UnitOfMeasure  string  `json:"unitofmeasure"`
	TotalLinePrice float64 `json:"totallineprice"`
}

type UnitOfMeasure struct {
}

type Order struct {
	OrderNumber string      `json:"orderNumber"`
	SupplierID  string      `json:"supplierId"`
	Items       []OrderItem `json:"items"`
	AuditInfo   AuditInfo   `json:"auditInfo"`
	From        string      `json:"from"`
	To          string      `json:"to"`
}
