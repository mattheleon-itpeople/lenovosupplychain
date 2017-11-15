package main

type Acknowledgement struct {
	ObjectType     string `json:"doctype"`
	DocumentType   string `json:"documentType"`
	documentNumber string `json:"documentNumber"`
	From           string `json:"from"`
	To             string `json:"to"`
}
