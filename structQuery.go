//QueryOrder : definition of Query Order
//structure to retrieve order from ledger
//###########################################

package main

type QueryOrder struct {
	OrderNumber string `json:"orderNumber"`
	Requestor   string `json:"requestor"`
	Partner     string `json:"partner"`
}

type QueryShipment struct {
	ShipmentNumber string `json:"shipmentNumber"`
	Requestor      string `json:"requestor"`
	Partner        string `json:"partner"`
}

type QueryField struct {
	FieldName  string `json:"fieldName"`
	FieldValue string `json:"fieldValue"`
}

type RichQuery struct {
	QueryName   string       `json:"queryName"`
	QueryFields []QueryField `json:"queryFields"`
}
