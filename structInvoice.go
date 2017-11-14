package main

type Invoice struct {
	OrderNumber string      `json:"OrderNumber"`
	Items       []OrderItem `json:"items"`
}
