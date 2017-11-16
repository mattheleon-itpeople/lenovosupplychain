package main

type ShippedItem struct {
	PartNumber      string `json:"partNumber"`
	SerialNumber    string `json:"serialNumber"`
	UnitOfMeasure   string `json:"unitOfMeasure"`
	ShippedQuantity string `json:"shippedQunatity"`
}

type Shipment struct {
	ObjectType     string        `json:"doctype"`
	ShipmentNumber string        `json:"shipmentNumber"`
	TrackingNumber string        `json:"trackingNumber"`
	DistributorID  string        `json:"distributorId"`
	ShippedItems   []ShippedItem `json:"shippedItems"`
	OrderNumber    string        `json:"OrderNumber"`
	From           string        `json:"from"`
	To             string        `json:"to"`
}
