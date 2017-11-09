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
)

func (t *LenovoChainCode) initMaps() {
	t.tableMap = make(map[string]int)
	t.tableMap[BIT] = 3
	t.funcMap = make(map[string]InvokeFunc)
	t.funcMap[GV] = getVersion
	t.funcMap[CPO] = CreatePurchaseOrder
	t.funcMap[CSO] = CreateShipment
	t.funcMap[SHPT] = ShipPart
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

func createPurchaseOrder(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	var err error
	var Avalbytes []byte
	logger.Infof("CreatePurchaseOrder : Arguments : %s", args[0])
	purchaseOrder := PurchaseOrder{}
	err = json.Unmarshal([]byte(args[0]), &purchaseOrder)
	if err != nil {
		return shim.Error("CreatePurchaseOrder : Failed to convert arg[0] to a PO object: " + err.Error())
	}

	// Query and Retrieve the Full BaicInfo
	keys := []string{purchaseOrder.PONumber}

	objectType := "PO"
	Avalbytes, err = dbapi.QueryObject(stub, objectType, keys)
	if err != nil {
		return shim.Error("CreatePurchaseOrder() : Failed to query POs object")
	}

	if Avalbytes != nil {
		return shim.Error(fmt.Sprintf("CreatePurchaseOrder() : "+
			"PO for PONumber: %s already exist ", purchaseOrder.PONumber))
	}

	err = dbapi.UpdateObject(stub, objectType, keys, []byte(args[0]))
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

	var err error
	var Avalbytes []byte
	logger.Infof("ShipPart : Arguments : %s", args[0])
	shipment := Shipment{}
	err = json.Unmarshal([]byte(args[0]), &shipment)
	if err != nil {
		return shim.Error("ShipPart : Failed to convert arg[0] to a Shipment object")
	}

	// Query and Retrieve the Full BaicInfo
	keys := []string{shipment.PartSerialNumber, shipment.SupplierID}
	logger.Infof("Keys for ShipPart : %v", keys)

	objectType := "RC"
	Avalbytes, err = dbapi.QueryObject(stub, objectType, keys)
	if err != nil {
		return shim.Error("ShipPart() : Failed to query Shipment object")
	}

	if Avalbytes != nil {
		return shim.Error(fmt.Sprintf("ShipPart() : "+
			"Shipment for PartNumber: %s and SupplierId: %s and PartSupplierNumber: %s already exist ",
			shipment.PartNumber, shipment.SupplierID, shipment.PartSerialNumber))
	}

	err = dbapi.UpdateObject(stub, objectType, keys, []byte(args[0]))
	if err != nil {
		logger.Errorf("ShipPart : Error inserting Object into LedgerState %s", err)
		return shim.Error("ShipPart : Shipment object Update failed")
	}

	return shim.Success(nil)

}

func sendAcknowledgement(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	return shim.Success(nil)
}
func createSalesOrder(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	return shim.Success(nil)
}
func generateSippingNotice(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	return shim.Success(nil)
}
func generateGoodsReceived(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	return shim.Success(nil)
}
func createReturnNotice(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	return shim.Success(nil)
}
func createInvoice(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	return shim.Success(nil)
}
func sendPayment(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	return shim.Success(nil)
}

////////////////////////////////////////////////////////////////////////////
// Query Function
////////////////////////////////////////////////////////////////////////////

func queryOrderByKey(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	return shim.Success(nil)
}
func queryPartsByOrderNumber(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	return shim.Success(nil)
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
//	}
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
