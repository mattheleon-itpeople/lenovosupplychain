package main

/*#########################################
Query Struct definitiosn
###########################################
*/
type QueryOrder struct {
	OrderNumber string `json:"OrderNumber"`
	Requestor   string `json:"Requestor"`
	Partner     string `json:"Partner"`
}

type QueryShipment struct {
	ShipmentNumber string `json:"shipmentnumber"`
	Requestor      string `json:"Requestor"`
	Partner        string `json:"Partner"`
}

type QueryField struct {
	FieldName  string `json:"fieldname"`
	FieldValue string `json:"fieldvalue"`
}

type RichQuery struct {
	QueryName   string       `json:"queryname"`
	QueryFields []QueryField `json:"queryfields"`
}
