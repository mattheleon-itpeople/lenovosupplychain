/*
Copyright IBM Corp. 2016 All Rights Reserved.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

		 http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

//this file will contain all logic related to chain code execution

package main

import (
	"encoding/json"
	"fmt"

	"dbapi"

	"github.com/hyperledger/fabric/core/chaincode/shim"
	pb "github.com/hyperledger/fabric/protos/peer"
)

type InvokeFunc func(stub shim.ChaincodeStubInterface, args []string) pb.Response

// SDMChaincode Chaincode implementation
type LenovoChainCode struct {
	tableMap map[string]int
	funcMap  map[string]InvokeFunc
}

// version
const VERSION string = "version"

/////////////////////////////////////////////////////
// Constant for table names
/////////////////////////////////////////////////////
const (
	BIT string = "BasicEntityInfoTable"
	USR string = "User"
)

/////////////////////////////////////////////////////
// Constant for different persona
/////////////////////////////////////////////////////
const (
	ADMIN string = "Admin"
)

/////////////////////////////////////////////////////
// Constant for PurchaseOrder status
/////////////////////////////////////////////////////
const (
	OPEN         = "open"
	CLOSED       = "closed"
	REFUNDED     = "refunded"
	ACKNOWLEDGED = "acknowledged"
)

/////////////////////////////////////////////////////
// Constant for PurchaseOrder type
/////////////////////////////////////////////////////
const (
	PURCHASE = "purchase"
	SALES    = "sales"
)

/////////////////////////////////////////////////////
// Constant for All function name that will be called from invoke
/////////////////////////////////////////////////////
const (
	GV   string = "getVersion"
	CPO  string = "createPurchaseOrder"
	CSO  string = "createShipment"
	CACK string = "createAcknowledgement"
	SHPT string = "shipPart"
	CUR  string = "createUser"
	UUR  string = "updateUser"
	QUR  string = "queryUser"
	DUR  string = "deleteUser"
	QO   string = "queryPurchaseOrder"
	QOBN string = "queryOrderByOrderNumber"
	QRQ  string = "queryRichQuery"
	QS   string = "queryShipment"
)

func (t *LenovoChainCode) initMaps() {
	t.tableMap = make(map[string]int)
	t.tableMap[BIT] = 3
	t.funcMap = make(map[string]InvokeFunc)
	t.funcMap[GV] = getVersion
	t.funcMap[CPO] = createPurchaseOrder
	t.funcMap[CSO] = createShipment
	t.funcMap[CACK] = createAcknowledgement
	t.funcMap[SHPT] = shipPart
	t.funcMap[QO] = queryPurchaseOrder
	t.funcMap[QOBN] = queryOrderByOrderNumber
	t.funcMap[QRQ] = queryRichQuery
	t.funcMap[QS] = queryShipment
	//	t.funcMap[CUR] = CreateUser
	//	t.funcMap[UUR] = UpdateUser
	//	t.funcMap[QUR] = QueryUser
	//	t.funcMap[DUR] = DeleteUser

}

////////////////////////////////////////////////////////////////////////////
// Query Version of BlockChain from Leadger
// This method for initial system test
////////////////////////////////////////////////////////////////////////////

func getVersion(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	logger.Info("-------- getVersion --------")
	// Get the cc version from the ledger
	version, err := stub.GetState(VERSION)
	if err != nil {
		jsonResp := "{\"Error\":\"Failed to get state for " + VERSION + "\"}"
		return shim.Error(jsonResp)
	}

	if version == nil {
		jsonResp := "{\"Error\":\"" + VERSION + " is nil \"}"
		return shim.Error(jsonResp)
	}

	jsonResp := "{\"Version\":\"" + string(version) + "\"}"
	logger.Infof("Query Response:%s\n", jsonResp)
	return shim.Success(version)
}

func createShipment(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	var err error
	var Avalbytes []byte
	Shipment := Shipment{}
	PurchaseOrder := PurchaseOrder{}

	if len(args) < 1 {
		return shim.Error("Not enough parameters")
	}
	err = json.Unmarshal([]byte(args[0]), &Shipment)
	if err != nil {
		return shim.Error("Failed to unmarshal shipment. " + args[0])
	}

	from := Shipment.From
	to := Shipment.To
	shipNumber := Shipment.ShipmentNumber

	keys := []string{to, from, shipNumber}

	objectType := "PO"
	//TODO: QUERY PurchaseOrder LOOP THROUGH PurchaseOrder ITEM
	Avalbytes, err = dbapi.QueryObject(stub, objectType, keys)

	if err != nil {
		return shim.Error("Failed to retrieve PurchaseOrder with provided shipping notice. " + err.Error())
	}
	if &Avalbytes == nil {
		return shim.Error("No PurchaseOrder was retrieved. " + err.Error())
	}
	err = json.Unmarshal(Avalbytes, &PurchaseOrder)
	if err != nil {
		return shim.Error("Failed to marshal Sales PurchaseOrder. " + string(Avalbytes))
	}

	//items := PurchaseOrder.Items

	quantity := make(map[string][]string)

	for _, i := range PurchaseOrder.Items {
		quantity[i.CommodityCode] = append(quantity[i.CommodityCode], i.OrderedQuantity)
		fmt.Println(quantity[i.CommodityCode])
	}

	/*
	   if len(PurchaseOrder.PurchaseOrderLine) != len(Shipment.ShippedItems) {
	       return shim.Error("***** PurchaseOrder quantity does not match shipping quantity. Changing PurchaseOrder status to: pending review. *****")
	   }

	   for iterator < len(PurchaseOrder.PurchaseOrderLine) {

	       // TAKE ITERATOR CREATE MAP OF BOTH SIDES AND COMPARE

	       //PurchaseOrderQuantity := PurchaseOrder.quantity
	       //shipQuantity := Shipment.quantity
	           iterator += iterator
	       }
	       iterator += iterator


	   if PurchaseOrderQuantity != shipQuantity {

	   }
	*/

	err = dbapi.UpdateObject(stub, objectType, keys, []byte(args[0]))
	if err != nil {
		logger.Errorf("shipPart : Error inserting Shipment of parts into LedgerState %s", err)
		return shim.Error("shipPart : Shipping part failed")
	}

	return shim.Success(nil)
}

func createPurchaseOrder(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	var err error
	var Avalbytes []byte
	logger.Infof("CreatePurchaseOrder : Arguments : %s", args[0])
	PurchaseOrder := PurchaseOrder{}
	err = json.Unmarshal([]byte(args[0]), &PurchaseOrder)
	if err != nil {
		return shim.Error("CreatePurchaseOrder : Failed to convert arg[0] to a PurchaseOrder: " + err.Error())
	}
	fmt.Println("Purchase order " + PurchaseOrder.To)

	// Query and Retrieve the Full BaicInfo
	keys := []string{PurchaseOrder.From, PurchaseOrder.To, PurchaseOrder.PONumber}

	PurchaseOrderType := "PO"
	Avalbytes, err = dbapi.QueryObject(stub, PurchaseOrderType, keys)
	if err != nil {
		return shim.Error("CreatePurchaseOrder() : Failed to query PurchaseOrder object")
	}

	if Avalbytes != nil {
		return shim.Error(fmt.Sprintf("CreatePurchaseOrder() : "+
			"PurchaseOrder for PurchaseOrder Number: %s already exist ", PurchaseOrder.PONumber))
	}

	err = dbapi.UpdateObject(stub, PurchaseOrderType, keys, []byte(args[0]))
	if err != nil {
		logger.Errorf("CreatePurchaseOrder : Error inserting Object into LedgerState %s", err)
		return shim.Error("CreatePurchaseOrder : POs Update failed")
	}

	return shim.Success(nil)

}

////////////////////////////////////////////////////////////////////////////
// List BasicEntityInfo data  about supplyer owned/invited by  specific Buyer.
// Key will be buyer unique-id
////////////////////////////////////////////////////////////////////////////

func shipPart(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	return shim.Success(nil)
}

func createAcknowledgement(stub shim.ChaincodeStubInterface, args []string) pb.Response {

	var orderBytes []byte
	var ackBytes []byte
	acknowledgement := Acknowledgement{}
	err := json.Unmarshal([]byte(args[0]), &acknowledgement)

	if len(args) < 1 {
		return shim.Error("sendAcknowledgement : requires an Acknowledgement document")
	}
	logger.Infof("CreateAckowledgement : Arguments : %s", args[0])

	if err != nil {
		return shim.Error("sendAcknowledgement : Failed to convert arg[0] to an Acknolwdgement notice: " + err.Error())
	}
	switch acknowledgement.DocumentType {
	case "PO":
		keys := []string{acknowledgement.To, acknowledgement.From, acknowledgement.DocumentNumber}
		order, err := retrieveAndMarshalPOObject(stub, keys)
		if err != nil {
			return shim.Error("sendAcknowledgement() - no existing Purchase Order Number " + acknowledgement.DocumentNumber)
		}

		/*Update order to show Acknowledgement */
		order.Status = ACKNOWLEDGED
		orderBytes, err = json.Marshal(order)
		if err != nil {
			return shim.Error("sendAcknowledgement() - failed to unmarshal existing Purchase Order Number " + acknowledgement.DocumentNumber)
		}
		err = dbapi.UpdateObject(stub, "PO", keys, orderBytes)
		if err != nil {
			return shim.Error("sendAcknowledgement() - failed to update existing Purchase Order Number " + acknowledgement.DocumentNumber)
		}
		aKeys := []string{acknowledgement.From, acknowledgement.To, acknowledgement.DocumentNumber}
		ackBytes, err = json.Marshal(acknowledgement)
		if err != nil {
			return shim.Error("sendAcknowledgement() - failed to update existing Purchase Order Number " + acknowledgement.DocumentNumber)
		}
		dbapi.UpdateObject(stub, "ACK", aKeys, ackBytes)

	default:
		return shim.Error("sendAcknowledgement - Ack for Doctype " + acknowledgement.DocumentType + " not yet implemented")
	}
	return shim.Success(nil)
}
func generateShippingNotice(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	return shim.Success(nil)
}
func generateGoodsReceived(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	return shim.Success(nil)
}
func createReturnNotice(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	var err error
	var Avalbytes []byte
	logger.Infof("CreateReturnNotice : Arguments : %s", args[0])
	returnNotice := ReturnNotice{}
	err = json.Unmarshal([]byte(args[0]), &returnNotice)
	if err != nil {
		return shim.Error("CreateReturnNotice : Failed to convert arg[0] to a Return notice: " + err.Error())
	}

	// Query and Retrieve the Full BaicInfo
	keys := []string{returnNotice.OrderNumber}

	objectType := "return"
	Avalbytes, err = dbapi.QueryObject(stub, objectType, keys)
	if err != nil {
		return shim.Error("CreateReturnNotice() : Failed to query return PurchaseOrder object")
	}

	if Avalbytes != nil {
		return shim.Error(fmt.Sprintf("CreateReturnNotice() : "+
			"ID for Return Invoice: %s already exist ", returnNotice.OrderNumber))
	}

	err = dbapi.UpdateObject(stub, objectType, keys, []byte(args[0]))
	if err != nil {
		logger.Errorf("CreateReturnNotice : Error inserting Object into LedgerState %s", err)
		return shim.Error("CreateReturnNotice : Return Notice Create failed")
	}

	return shim.Success(nil)
}
func createInvoice(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	var err error
	var Avalbytes []byte
	logger.Infof("CreateInvoice : Arguments : %s", args[0])
	invoice := Invoice{}
	err = json.Unmarshal([]byte(args[0]), &invoice)
	if err != nil {
		return shim.Error("CreateInvoice : Failed to convert arg[0] to a Invoice object: " + err.Error())
	}

	// Query and Retrieve the Full BaicInfo
	keys := []string{invoice.OrderNumber}

	objectType := "invoice"
	Avalbytes, err = dbapi.QueryObject(stub, objectType, keys)
	if err != nil {
		return shim.Error("CreateInvoice() : Failed to query shipment object")
	}

	if Avalbytes != nil {
		return shim.Error(fmt.Sprintf("CreateInvoice() : "+
			"ID for Invoice Number: %s already exist ", invoice.OrderNumber))
	}

	err = dbapi.UpdateObject(stub, objectType, keys, []byte(args[0]))
	if err != nil {
		logger.Errorf("CreateInvoice : Error inserting Object into LedgerState %s", err)
		return shim.Error("CreateInvoice : Invoice Create failed")
	}

	return shim.Success(nil)
}
func sendPayment(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	return shim.Success(nil)
}

////////////////////////////////////////////////////////////////////////////
// Query Function
////////////////////////////////////////////////////////////////////////////
////////////////////////////////////////////////////////////////////////////
// Query PurchaseOrder given the PurchaseOrder Number and the 'From' organization
////////////////////////////////////////////////////////////////////////////
func queryOrderByOrderNumber(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	var err error
	var PurchaseOrders []PurchaseOrder
	var PurchaseOrder PurchaseOrder = PurchaseOrder{}
	var query QueryOrder = QueryOrder{}

	var i = 0

	if len(args) < 1 {
		logger.Infof("queryOrderByOrderNumber - requires one parameter (originating company)")
		return shim.Error("queryOrderByOrderNumber - requires one parameter (originating company)")
	}

	err = json.Unmarshal([]byte(args[0]), &query)

	if err != nil {
		logger.Infof("queryOrderByOrderNumber - failed to marshal query object")
		return shim.Error("queryOrderByOrderNumber - failed to marshal query object")
	}
	logger.Infof("queryOrderByOrderNumber : Arguments : %s", query.Requestor)

	keys := []string{query.Requestor}

	rs, err := dbapi.GetList(stub, "PO", keys)

	if err != nil {
		logger.Infof("queryOrderByOrderNumber - failed to retrieve PurchaseOrders: %s", keys[0])
		return shim.Error("queryOrderByOrderNumber - failed to retrieve PurchaseOrders: %s)" + keys[0])
	}

	for i = 0; rs.HasNext(); i++ {
		myKV, err := rs.Next()
		if err != nil {
			return shim.Error(err.Error())
		}

		err = json.Unmarshal(myKV.Value, &PurchaseOrder)

		if err != nil {
			logger.Infof("queryOrderByOrderNumber - failed to marshal PurchaseOrder: %s", err.Error())
			return shim.Error("queryOrderByOrderNumber - failed to marshal PurchaseOrder: " + err.Error())
		}

		if PurchaseOrder.PONumber == query.OrderNumber {
			PurchaseOrders = append(PurchaseOrders, PurchaseOrder)
		}
	}

	jsonRows, err := json.Marshal(PurchaseOrders)
	return shim.Success(jsonRows)
}

////////////////////////////////////////////////////////////////////////////
// Query a specific PurchaseOrder with a full key
////////////////////////////////////////////////////////////////////////////
func queryPurchaseOrder(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	var err error
	var PurchaseOrderbytes []byte
	var query = QueryOrder{}

	if len(args) < 1 {
		logger.Infof("queryPurchaseOrder requires request paramater")
		return shim.Error("queryPurchaseOrder requires request parameter")
	}
	logger.Infof("queryPurchaseOrder : Arguments : %s", args[0])
	err = json.Unmarshal([]byte(args[0]), &query)

	if err != nil {
		logger.Infof("queryPurchaseOrder : Arguments : %s", args[0])
	}

	keys := []string{query.Requestor, query.Partner, query.OrderNumber}
	PurchaseOrderbytes, err = dbapi.QueryObject(stub, "PO", keys)

	if err != nil {
		logger.Infof("queryPurchaseOrder fail to retrieve PurchaseOrder (PurchaseOrder number: %s, company %s )", query.OrderNumber, query.Requestor)
		return shim.Error("queryPurchaseOrder fail to retrieve PurchaseOrder")
	}

	return shim.Success(PurchaseOrderbytes)
}

////////////////////////////////////////////////////////////////////////////
// Query All PurchaseOrders for a specific company (in the 'From')
////////////////////////////////////////////////////////////////////////////
func queryAllPurchaseOrders(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	var err error
	var PurchaseOrders []PurchaseOrder
	var PurchaseOrder PurchaseOrder = PurchaseOrder{}
	var query = QueryOrder{}
	var i = 0

	logger.Infof("Received %s as arguments  ")

	if len(args) < 1 {
		logger.Infof("queryPurchaseOrder requires request paramater")
		return shim.Error("queryPurchaseOrder requires request parameter")
	}

	logger.Infof("queryPurchaseOrder : Arguments : %s", args[0])

	keys := []string{args[0]}

	rs, err := dbapi.GetList(stub, "PO", keys)

	if err != nil {
		logger.Infof("queryPurchaseOrder fail to retrieve PurchaseOrders: %s", args[0])
		return shim.Error("queryPurchaseOrder fail to retrieve PurchaseOrders: )" + args[0])
	}

	for i = 0; rs.HasNext(); i++ {
		myKV, err := rs.Next()
		if err != nil {
			return shim.Error(err.Error())
		}

		err = json.Unmarshal(myKV.Value, &PurchaseOrder)

		if err != nil {
			logger.Infof("queryPurchaseOrder fail to marshal PurchaseOrder: %s", err.Error())
			return shim.Error("queryPurchaseOrder fail to marshal PurchaseOrder: " + err.Error())
		}

		if PurchaseOrder.From == query.Requestor || PurchaseOrder.To == query.Requestor {
			PurchaseOrders = append(PurchaseOrders, PurchaseOrder)
		}
	}

	jsonRows, err := json.Marshal(PurchaseOrders)
	return shim.Success(jsonRows)
}

////////////////////////////////////////////////////////////////////////////
//  Get Shipment
////////////////////////////////////////////////////////////////////////////

////////////////////////////////////////////////////////////////////////////
// Query a specific Shipment with a full key
////////////////////////////////////////////////////////////////////////////
func queryShipment(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	var err error
	var Shipmentbytes []byte
	var query = QueryShipment{}

	if len(args) < 1 {
		logger.Infof("queryShipment requires request paramater")
		return shim.Error("queryShipment requires request parameter")
	}
	logger.Infof("queryShipment : Arguments : %s", args[0])
	err = json.Unmarshal([]byte(args[0]), &query)

	if err != nil {
		logger.Infof("queryShipment : Arguments : %s", args[0])
	}

	keys := []string{query.Requestor, query.Partner, query.ShipmentNumber}
	Shipmentbytes, err = dbapi.QueryObject(stub, "SHP", keys)

	if err != nil {
		logger.Infof("queryShipment fail to retrieve shipment (shipment number: %s, company %s )", query.ShipmentNumber, query.Requestor)
		return shim.Error("queryShipment fail to retrieve shipment")
	}

	return shim.Success(Shipmentbytes)
}

////////////////////////////////////////////////////////////////////////////
//  Rich query for all PurchaseOrders
////////////////////////////////////////////////////////////////////////////
func queryRichQuery(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	var err error
	var query RichQuery = RichQuery{}
	var queryfield string
	var formattedstring string

	if len(args) < 1 {
		logger.Infof("queryRichQuery requires request paramater")
		return shim.Error("queryRichQuery requires request parameter")
	}

	logger.Infof("queryRichQuery : Arguments : %s", args[0])
	err = json.Unmarshal([]byte(args[0]), &query)

	if err != nil {
		logger.Infof("queryRichQuery : Arguments : %s, %s", args[0], err.Error())
		return shim.Error("queryRichQuery : error unmarshalling: " + err.Error())
	}

	logger.Infof("queryRichQuery : query name %s", query.QueryName)
	if len(query.QueryFields) < 1 {
		strAsBytes, _ := json.Marshal(query)
		logger.Infof("queryRichQuery : Arguments : %s", strAsBytes)
		return shim.Error("queryRichQuery : no fields requested")
	}

	queryfield = query.QueryFields[0].FieldValue

	formattedstring = getFormattedPurchaseOrderQuery(queryfield)
	logger.Infof("queryRichQuery : Query : %s", formattedstring)
	querybytes, err := dbapi.GetQueryResultForQueryString(stub, formattedstring)

	if err != nil {
		logger.Infof("queryRichQuery fail to retrieve PurchaseOrders: %s", err.Error())
		return shim.Error("queryRichQuery fail to retrieve PurchaseOrders: )" + err.Error())
	}

	return shim.Success(querybytes)
}

func queryPurchaseOrderStatus(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	return shim.Success(nil)
}
func queryShipmentByPurchaseOrderNumber(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	return shim.Success(nil)
}

////////////////////////////////////////////////////////////////////////////
//Query by selector (rich query!)
////////////////////////////////////////////////////////////////////////////

////////////////////////////////////////////////////////////////////////////
// Helper functions
////////////////////////////////////////////////////////////////////////////

func getFormattedPurchaseOrderQuery(PurchaseOrderNumber string) string {
	return fmt.Sprintf("{\\\"selector\\\": { \\\"PurchaseOrderNumber\\\": \\\"%s\\\"}}", PurchaseOrderNumber)
}

////////////////////////////////////////////////////////////////////////////
// Inserting Data in to user table
////////////////////////////////////////////////////////////////////////////

//func CreateUser(stub shim.ChaincodeStubInterface, args []string) pb.Response {

//	var Avalbytes []byte
//	// Convert the arg to a userrequest Object
//	logger.Infof("CreateUser : Arguments for dbapi : %s", args[0])

//	userreq, err := common.JSONtoUserRequest([]byte(args[0]))
//	userreq.UseInfo.ObjectType = USR
//	if err != nil {
//		return shim.Error("CreateUser : Failed to convert arg[0] to a UserRequest object")
//	}
//	if !strings.EqualFold(userreq.UserPersona, ADMIN) {
//		return shim.Error("CreateUser : User is not authorized to create this user")
//	}

//	// Query and Retrieve the Full User Data
//	keys := []string{userreq.UseInfo.EnrollmentID}
//	logger.Infof("Keys for user : %v", keys)

//	Avalbytes, err = sdmdbapi.QueryObject(stub, USR, keys)
//	if err != nil {
//		return shim.Error("CreateUser() : Failed to query user object")
//	}

//	if Avalbytes != nil {
//		return shim.Error(fmt.Sprintf("CreateUser() : user for EnrollmentID: %s already exist ", userreq.UseInfo.EnrollmentID))
//	userjson, err := common.UsertoJSON(userreq.UseInfo)
//	err = sdmdbapi.UpdateObject(stub, USR, keys, []byte(userjson))
//	if err != nil {
//		logger.Errorf("CreateUser : Error inserting Object into LedgerState %s", err)
//		return shim.Error("CreateUser : user object Update failed")
//	}

//	return shim.Success([]byte(userjson))

//}

//////////////////////////////////////////////////////////////////////////////
//// Query User data  for buyer/supplyer.
//// Key will be buyer/supllier Enrollment Id
//////////////////////////////////////////////////////////////////////////////

//func QueryUser(stub shim.ChaincodeStubInterface, args []string) pb.Response {

//	var Avalbytes []byte

//	userreq, err := common.JSONtoUserRequest([]byte(args[0]))
//	userreq.UseInfo.ObjectType = USR
//	if err != nil {
//		return shim.Error("QueryUser : Failed to convert arg[0] to a UserRequest object")
//	}

//	if !strings.EqualFold(userreq.UserPersona, ADMIN) {
//		return shim.Error("QueryUser : User is not authorized to create this user")
//	}

//	userkey := []string{userreq.UseInfo.EnrollmentID}

//	Avalbytes, err = sdmdbapi.QueryObject(stub, USR, userkey)

//	me, _ := common.JSONtoUser(Avalbytes)
//	logger.Infof("QueryUser() : **** User ****: %v", me)

//	if err != nil {
//		return shim.Error("QueryUser() : Failed to query object successfully")
//	}

//	if Avalbytes == nil {
//		return shim.Error("QueryUser() : User Request id not found " + args[0])
//	}

//	return shim.Success(Avalbytes)

//}

//////////////////////////////////////////////////////////////////////////////
//// Updating User data belongs for buyer/supplier
////
//////////////////////////////////////////////////////////////////////////////

//func UpdateUser(stub shim.ChaincodeStubInterface, args []string) pb.Response {

//	var Avalbytes []byte
//	// Convert the arg to a userrequest Object
//	logger.Infof("UpdateUser : Arguments for dbapi : %s", args[0])

//	userreq, err := common.JSONtoUserRequest([]byte(args[0]))
//	userreq.UseInfo.ObjectType = USR

//	if err != nil {
//		return shim.Error("UpdateUser : Failed to convert arg[0] to a UserRequest object")
//	}
//	if !strings.EqualFold(userreq.UserPersona, ADMIN) {
//		return shim.Error("UpdateUser : User is not authorized to create this user")
//	}

//	// Query and Retrieve the Full user
//	keys := []string{userreq.UseInfo.EnrollmentID}
//	logger.Infof("Keys for user : %v", keys)
//	if err != nil {
//		return shim.Error("UpdateUser : Failed to convert user to a json")
//	}
//	Avalbytes, err = sdmdbapi.QueryObject(stub, USR, keys)
//	if err != nil {
//		return shim.Error("UpdateUser() : Failed to query user object")
//	}

//	if Avalbytes == nil {
//		return shim.Error(fmt.Sprintf("UpdateUser() : user for EnrollmentID: %s does not exist ", userreq.UseInfo.EnrollmentID))
//	}
//	userjson, err := common.UsertoJSON(userreq.UseInfo)
//	err = sdmdbapi.UpdateObject(stub, USR, keys, []byte(userjson))
//	if err != nil {
//		logger.Errorf("UpdateUser: Error Updating Object into LedgerState %s", err)
//		return shim.Error("UpdateUser : User object Update failed")
//	}

//	return shim.Success([]byte(userjson))

//}

//////////////////////////////////////////////////////////////////////////////
//// Delete User  belongs for buyer/supplier
////
//////////////////////////////////////////////////////////////////////////////

//func DeleteUser(stub shim.ChaincodeStubInterface, args []string) pb.Response {

//	var Avalbytes []byte
//	// Convert the arg to a userrequest Object
//	logger.Infof("DeleteUser : Arguments for dbapi : %s", args[0])

//	userreq, err := common.JSONtoUserRequest([]byte(args[0]))
//	userreq.UseInfo.ObjectType = USR

//	if err != nil {
//		return shim.Error("DeleteUser : Failed to convert arg[0] to a UserRequest object")
//	}
//	if !strings.EqualFold(userreq.UserPersona, ADMIN) {
//		return shim.Error("DeleteUser : User is not authorized to Delete this user")
//	}

//	// Query and Retrieve the Full user
//	keys := []string{userreq.UseInfo.EnrollmentID}
//	logger.Infof("Keys for user : %v", keys)
//	Avalbytes, err = sdmdbapi.QueryObject(stub, USR, keys)
//	if err != nil {
//		return shim.Error("DeleteUser() : Failed to query user object")
//	}

//	if Avalbytes == nil {
//		return shim.Error(fmt.Sprintf("DeleteUser() : user for EnrollmentID: %s does not exist ", userreq.UseInfo.EnrollmentID))
//	}

//	err = sdmdbapi.DeleteObject(stub, USR, keys)
//	if err != nil {
//		logger.Errorf("DeleteUser: Error Deletng  Object from  LedgerState %s", err)i
//		return shim.Error("DeleteUser : User object Delete failed")
//	}

//	return shim.Success([]byte(Avalbytes))

//}
