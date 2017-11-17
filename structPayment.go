package main

type Payment struct {
	ObjectType    string `json:"doctype"`
	PaymentNumber string `json:"paymentNumber"`
	InvoiceNumber string `json:"invoiceNumber"`
	Remittance    string `json:"remittance"` //Amount to be pait
	From          string `json:"from"`
	To            string `json:"to"`
}
