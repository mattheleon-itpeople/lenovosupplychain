package main

type ReturnNotice struct {
	OrderNumber string      `json:"OrderNumber"`
	Items       []OrderItem `json:"items"`
}
