package main

type OrderStatus struct {
	ObjectType  string `json:"doctype"`
	OrderNumber string `json:"orderNumber"`
	Status      string `json:"status"`
}
