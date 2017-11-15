package main

type ReturnNotice struct {
	ObjectType  string        `json:"doctype"`
	OrderNumber string        `json:"orderNumber"`
	Items       []ItemDetails `json:"items"`
}
