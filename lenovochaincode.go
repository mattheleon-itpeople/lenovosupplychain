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
// Constant for order status
/////////////////////////////////////////////////////
const (
	OPEN     = "open"
	CLOSED   = "closed"
	REFUNDED = "refunded"
)

/////////////////////////////////////////////////////
// Constant for order type
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
	CPO  string = "createOrder"
	CSO  string = "createShipment"
	SHPT string = "shipPart"
	CUR  string = "createUser"
	UUR  string = "updateUser"
	QUR  string = "queryUser"
	DUR  string = "deleteUser"
	QO   string = "queryOrder"
	QOBN string = "queryOrderByOrderNumber"
	QRQ  string = "queryRichQuery"
	QS   string = "queryShipment"
)

func (t *LenovoChainCode) initMaps() {
	t.tableMap = make(map[string]int)
	t.tableMap[BIT] = 3
	t.funcMap = make(map[string]InvokeFunc)
	t.funcMap[GV] = getVersion
	t.funcMap[CPO] = createOrder
	t.funcMap[CSO] = createShipment
	t.funcMap[SHPT] = shipPart
	t.funcMap[QO] = queryOrder
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
	logger.Infof("CreateShipment : Arguments : %s", args[0])
	shipment := Shipment{}
	err = json.Unmarshal([]byte(args[0]), &shipment)
	if err != nil {
		return shim.Error("CreateShipment : Failed to convert arg[0] to a Shipment object: " + err.Error())
	}

	// Query and Retrieve the Full BaicInfo
	keys := []string{shipment.ShipmentNumber}

	objectType := "SHP"
	Avalbytes, err = dbapi.QueryObject(stub, objectType, keys)
	if err != nil {
		return shim.Error("CreateShipment() : Failed to query shipment object")
	}

	if Avalbytes != nil {
		return shim.Error(fmt.Sprintf("CreateShipment() : "+
			"ID for Shipment Number: %s already exist ", shipment.ShipmentNumber))
	}

	err = dbapi.UpdateObject(stub, objectType, keys, []byte(args[0]))
	if err != nil {
		logger.Errorf("CreateShipment : Error inserting Object into LedgerState %s", err)
		return shim.Error("CreateShipment : Shipment Update failed")
	}

	return shim.Success(nil)

}

func createOrder(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	var err error
	var Avalbytes []byte
	logger.Infof("CreateOrder : Arguments : %s", args[0])
	Order := Order{}
	err = json.Unmarshal([]byte(args[0]), &Order)
	if err != nil {
		return shim.Error("CreateOrder : Failed to convert arg[0] to a Order: " + err.Error())
	}

	// Query and Retrieve the Full BaicInfo
	keys := []string{Order.From, Order.To, Order.OrderNumber}

	orderType := "PO"
	Avalbytes, err = dbapi.QueryObject(stub, orderType, keys)
	if err != nil {
		return shim.Error("CreateOrder() : Failed to query Order object")
	}

	if Avalbytes != nil {
		return shim.Error(fmt.Sprintf("CreateOrder() : "+
			"Order for Order Number: %s already exist ", Order.OrderNumber))
	}

	err = dbapi.UpdateObject(stub, orderType, keys, []byte(args[0]))
	if err != nil {
		logger.Errorf("CreateOrder : Error inserting Object into LedgerState %s", err)
		return shim.Error("CreateOrder : POs Update failed")
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

func sendAcknowledgement(stub shim.ChaincodeStubInterface, args []string) pb.Response {
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

	objectType := "RET"
	Avalbytes, err = dbapi.QueryObject(stub, objectType, keys)
	if err != nil {
		return shim.Error("CreateReturnNotice() : Failed to query return order object")
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

	objectType := "INV"
	Avalbytes, err = dbapi.QueryObject(stub, objectType, keys)
	if err != nil {
		return shim.Error("CreateInvoice() : Failed to query invoice")
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
// Query Order given the Order Number and the 'From' organization
////////////////////////////////////////////////////////////////////////////
func queryOrderByOrderNumber(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	var err error
	var orders []Order
	var order Order = Order{}
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
		logger.Infof("queryOrderByOrderNumber - failed to retrieve orders: %s", keys[0])
		return shim.Error("queryOrderByOrderNumber - failed to retrieve orders: %s)" + keys[0])
	}

	for i = 0; rs.HasNext(); i++ {
		myKV, err := rs.Next()
		if err != nil {
			return shim.Error(err.Error())
		}
		err = json.Unmarshal(myKV.Value, &order)

		if err != nil {
			logger.Infof("queryOrderByOrderNumber - failed to marshal order: %s", err.Error())
			return shim.Error("queryOrderByOrderNumber - failed to marshal order: " + err.Error())
		}
		if order.OrderNumber == query.OrderNumber {
			orders = append(orders, order)
		}
	}

	jsonRows, err := json.Marshal(orders)
	return shim.Success(jsonRows)
}

////////////////////////////////////////////////////////////////////////////
// Query a specific Order with a full key
////////////////////////////////////////////////////////////////////////////
func queryOrder(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	var err error
	var Orderbytes []byte
	var query = QueryOrder{}

	if len(args) < 1 {
		logger.Infof("queryOrder requires request paramater")
		return shim.Error("queryOrder requires request parameter")
	}
	logger.Infof("queryOrder : Arguments : %s", args[0])
	err = json.Unmarshal([]byte(args[0]), &query)

	if err != nil {
		logger.Infof("queryOrder : Arguments : %s", args[0])
	}

	keys := []string{query.Requestor, query.Partner, query.OrderNumber}
	Orderbytes, err = dbapi.QueryObject(stub, "PO", keys)

	if err != nil {
		logger.Infof("queryOrder fail to retrieve order (order number: %s, company %s )", query.OrderNumber, query.Requestor)
		return shim.Error("queryOrder fail to retrieve order")
	}

	return shim.Success(Orderbytes)
}

////////////////////////////////////////////////////////////////////////////
// Query All Orders for a specific company (in the 'From')
////////////////////////////////////////////////////////////////////////////
func queryAllOrders(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	var err error
	var orders []Order
	var order Order = Order{}
	var query = QueryOrder{}
	var i = 0

	logger.Infof("Received %s as arguments  ")

	if len(args) < 1 {
		logger.Infof("queryOrder requires request paramater")
		return shim.Error("queryOrder requires request parameter")
	}

	logger.Infof("queryOrder : Arguments : %s", args[0])

	keys := []string{args[0]}

	rs, err := dbapi.GetList(stub, "PO", keys)

	if err != nil {
		logger.Infof("queryOrder fail to retrieve orders: %s", args[0])
		return shim.Error("queryOrder fail to retrieve orders: )" + args[0])
	}

	for i = 0; rs.HasNext(); i++ {
		myKV, err := rs.Next()
		if err != nil {
			return shim.Error(err.Error())
		}

		err = json.Unmarshal(myKV.Value, &order)

		if err != nil {
			logger.Infof("queryOrder fail to marshal order: %s", err.Error())
			return shim.Error("queryOrder fail to marshal order: " + err.Error())
		}

		if order.From == query.Requestor || order.To == query.Requestor {
			orders = append(orders, order)
		}
	}

	jsonRows, err := json.Marshal(orders)
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
// Query a specific Shipment with a partial key
////////////////////////////////////////////////////////////////////////////
func queryAllShipments(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	var err error
	var shipments []Shipment
	var shipment Shipment = Shipment{}
	var query = QueryShipment{}
	var i = 0

	logger.Infof("Received %s as arguments  ")

	if len(args) < 1 {
		logger.Infof("queryAllShipments requires request paramater")
		return shim.Error("queryAllShipments requires request parameter")
	}

	logger.Infof("queryAllShipments : Arguments : %s", args[0])

	keys := []string{args[0]}

	rs, err := dbapi.GetList(stub, "SHP", keys)

	if err != nil {
		logger.Infof("queryAllShipments fail to retrieve orders: %s", args[0])
		return shim.Error("queryAllShipments fail to retrieve orders: )" + args[0])
	}

	for i = 0; rs.HasNext(); i++ {
		myKV, err := rs.Next()
		if err != nil {
			return shim.Error(err.Error())
		}

		err = json.Unmarshal(myKV.Value, &shipment)

		if err != nil {
			logger.Infof("queryAllShipments fail to marshal order: %s", err.Error())
			return shim.Error("queryAllShipments fail to marshal order: " + err.Error())
		}

		if shipment.From == query.Requestor || shipment.To == query.Requestor {
			shipments = append(shipments, shipment)
		}
	}

	jsonRows, err := json.Marshal(shipments)
	return shim.Success(jsonRows)
}

////////////////////////////////////////////////////////////////////////////
//  Rich query for all orders
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

	formattedstring = getFormattedOrderQuery(queryfield)
	logger.Infof("queryRichQuery : Query : %s", formattedstring)
	querybytes, err := dbapi.GetQueryResultForQueryString(stub, formattedstring)

	if err != nil {
		logger.Infof("queryRichQuery fail to retrieve orders: %s", err.Error())
		return shim.Error("queryRichQuery fail to retrieve orders: )" + err.Error())
	}

	return shim.Success(querybytes)
}

func queryOrderStatus(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	return shim.Success(nil)
}
func queryShipmentByOrderNumber(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	return shim.Success(nil)
}

////////////////////////////////////////////////////////////////////////////
//Query by selector (rich query!)
////////////////////////////////////////////////////////////////////////////

////////////////////////////////////////////////////////////////////////////
// Helper functions
////////////////////////////////////////////////////////////////////////////

func getFormattedOrderQuery(orderNumber string) string {
	return fmt.Sprintf("{\\\"selector\\\": { \\\"orderNumber\\\": \\\"%s\\\"}}", orderNumber)
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
