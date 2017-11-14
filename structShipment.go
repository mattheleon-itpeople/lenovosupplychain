package main

type ShippedItem struct {
	PartNumber    string `json:"partNumber"`
	SerialNumber  string `json:"serialNumber"`
	UnitOfMeasure string `json:"unitofmeasure"`
}

type Shipment struct {
	ShipmentNumber string        `json:"shipmentNumber"`
	TrackingNumber string        `json:"trackingNumber"`
	SupplierID     string        `json:"supplierId"`
	ShippedItems   []ShippedItem `json:"shippedItems"`
	OrderNumber    string        `json:"OrderNumber"`
}
