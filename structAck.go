package main

type Acknowledgement struct {
	ObjectType     string `json:"doctype"`
	AckNumber      string `json:"ackNumber"`
	DocumentType   string `json:"documentType"`
	DocumentNumber string `json:"documentNumber"`
	From           string `json:"from"`
	To             string `json:"to"`
}
