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

var logger = shim.NewLogger("lenovo_chaincode")

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
	CPO  string = "createPurchaseOrder"
	CSO  string = "createShipment"
	SHPT string = "shipPart"
	CUR  string = "createUser"
	UUR  string = "updateUser"
	QUR  string = "queryUser"
	DUR  string = "deleteUser"
	QO   string = "queryOrder"
	QOBN string = "queryOrderByOrderNumber"
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

	objectType := "PO"
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

	objectType := "return"
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
// Query Order given the Order Number and the 'From' organization
////////////////////////////////////////////////////////////////////////////
func queryOrderByOrderNumber(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	var query = QueryOrder{}
	var i = 0

	if len(args) < 1 {
		logger.Infof("queryOrderByOrderNumber requires request paramater")
		return shim.Error("queryOrderByOrderNumber requires request parameter")
	}

	err := json.Unmarshal([]byte(args[0]), &query)

	if err != nil {
		logger.Infof("queryOrderByOrderNumber failed to unmarshal data :" + err.Error())
		return shim.Error("queryOrderByOrderNumber failed to unmarshal data : " + err.Error())
	}

	logger.Infof("queryOrderByOrderNumber : Arguments :" + query.OrderNumber + " : " + query.From)
	keys := []string{query.From}
	results, err := dbapi.GetList(stub, "PO", keys)
	logger.Info("QueryByGetQuery - returned from dbapi")

	for i = 0; results.HasNext(); i++ {
		logger.Info("QueryByGetQuery - Iterating")
		// Retrieve the Key and Object
		myCompositeKey, err := results.Next()
		if err != nil {
			return shim.Error(err.Error())
		}
		logger.Infof("QueryOrderByOrderNumber() : my Value : ", myCompositeKey)
	}
	return shim.Success(nil)
}

////////////////////////////////////////////////////////////////////////////
// Query a specific Order with a full key
////////////////////////////////////////////////////////////////////////////
func queryOrder(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	var err error
	var Orderbytes []byte
	var query = QueryOrder{}

	logger.Infof("Received %s as arguments  ")

	if len(args) < 1 {
		logger.Infof("queryOrder requires request paramater")
		return shim.Error("queryOrder requires request parameter")
	}

	err = json.Unmarshal([]byte(args[0]), &query)
	logger.Infof("queryOrder : Arguments : %s", args[0])

	keys := []string{query.From, query.To, query.OrderNumber}
	Orderbytes, err = dbapi.QueryObject(stub, "PO", keys)

	if err != nil {
		logger.Infof("queryOrder fail to retrieve order (order number: %s, company %s )", query.OrderNumber, query.From)
		return shim.Error("queryOrder fail to retrieve order")
	}

	return shim.Success(Orderbytes)
}

func queryAllOrders(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	var err error
	var orders []Order
	var order Order = Order{}
	var i = 0

	if len(args) < 1 {
		logger.Infof("queryOrder requires one parameter (originating company)")
		return shim.Error("queryOrder requires one parameter (originating company)")
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

		orders = append(orders, order)
	}

	jsonRows, err := json.Marshal(orders)
	return shim.Success(jsonRows)
}

func queryOrderStatus(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	return shim.Success(nil)
}
func queryShipmentByOrderNumber(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	return shim.Success(nil)
}

////////////////////////////////////////////////////////////////////////////
// Helper functions
////////////////////////////////////////////////////////////////////////////

func getFormattedOrderQuery(orderNumber string, From string) string {
	return fmt.Sprintf("{\"selector\": { \"orderNumber\": \"%s\"}}", orderNumber)
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
//		logger.Errorf("DeleteUser: Error Deletng  Object from  LedgerState %s", err)
//		return shim.Error("DeleteUser : User object Delete failed")
//	}

//	return shim.Success([]byte(Avalbytes))

//}
