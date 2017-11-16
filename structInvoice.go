package main

type Invoice struct {
	ObjectType          string        `json:"doctype"`
	InvoiceNumber       string        `json:"invoiceNumber"`
	From                string        `json:"from"`
	To                  string        `json:"to"`
	OriginalOrderNumber string        `json:"originalOrderNumber"`
	OriginalOrderType   string        `json:orginalOrderType`
	Items               []ItemDetails `json:"items"`
}
