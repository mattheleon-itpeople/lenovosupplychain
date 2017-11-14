package main

type RecievedItem struct {
	PartNumber      string `json:"partNumber"`
	ShipmentNumber  string `json:"shipmentNumber"`
	ReceivingStatus string `json:"receivingStatus"`
}
