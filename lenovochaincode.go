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
	"strconv"

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
	CSO  string = "createSalesOrder"
	CSP  string = "createShipment"
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
	QSO  string = "querySalesOrder"
)

func (t *LenovoChainCode) initMaps() {
	t.tableMap = make(map[string]int)
	t.tableMap[BIT] = 3
	t.funcMap = make(map[string]InvokeFunc)
	t.funcMap[GV] = getVersion
	t.funcMap[CPO] = createPurchaseOrder
	t.funcMap[CSP] = createShipment
	t.funcMap[CSO] = createSalesOrder
	t.funcMap[CACK] = createAcknowledgement
	t.funcMap[SHPT] = shipPart
	t.funcMap[QO] = queryPurchaseOrder
	t.funcMap[QOBN] = queryOrderByOrderNumber
	t.funcMap[QRQ] = queryRichQuery
	t.funcMap[QS] = queryShipment
	t.funcMap[QSO] = querySalesOrder
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
	funcName := getFunctionName()
	logger.Info("-------- " + funcName + " --------")
	// Get the cc version from the ledger
	version, err := stub.GetState(VERSION)
	if err != nil {
		jsonResp := "{\"Error\":\"" + funcName + " - Failed to get state for " + VERSION + "\"}"
		return shim.Error(jsonResp)
	}

	if version == nil {
		jsonResp := "{\"Error\":\"" + funcName + " : " + VERSION + " is nil \"}"
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
	funcName := getFunctionName()

	if len(args) < 1 {
		return shim.Error("Not enough parameters")
	}

	fmt.Println("createShipment : Arguments : %s", args[0])

	err = json.Unmarshal([]byte(args[0]), &Shipment)
	if err != nil {
		return shim.Error(funcName + " : Failed to unmarshal shipment. " + err.Error())
	}

	distributerID := Shipment.DistributorID
	to := Shipment.To
	shippingNumber := Shipment.ShipmentNumber
	keys := []string{to, distributerID, shippingNumber}

	objectType := "PO"
	Avalbytes, err = dbapi.QueryObject(stub, objectType, keys)

	if err != nil {
		return shim.Error(funcName + " : Failed to retrieve order with provided order number. " + err.Error())
	}
	if &Avalbytes == nil {
		return shim.Error(funcName + " : No order was retrieved. " + err.Error())
	}

	err = json.Unmarshal(Avalbytes, &PurchaseOrder)
	if err != nil {
		return shim.Error(funcName + " : Failed to marshal Sales Order. " + string(Avalbytes))
	}

	if len(PurchaseOrder.Items) != len(Shipment.ShippedItems) {
		return shim.Error(funcName + " : Order quantity does not match shipping quantity. Changing order status to: pending review. *****")
	}

	orderedquantity := make(map[string][]int)

	for _, i := range PurchaseOrder.Items {
		fmt.Println("jkfnfkjnhns : " + i.OrderedQuantity)
		quant, err := strconv.Atoi(i.OrderedQuantity)
		orderedquantity[i.CommodityCode] = append(orderedquantity[i.CommodityCode], quant)
		fmt.Println(orderedquantity[i.CommodityCode])

		if err != nil {
			return shim.Error(funcName + " : ***** Error converting quantity to int *****")
		}
	}

	objectType = "SHP"
	err = dbapi.UpdateObject(stub, objectType, keys, []byte(args[0]))

	if err != nil {
		logger.Errorf(funcName+" :  Error inserting Shipment of parts into LedgerState %s", err)
		return shim.Error(funcName + " :  Shipping part failed")
	}

	return shim.Success(nil)
}
func createPurchaseOrder(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	var err error
	var Avalbytes []byte
	var PurchaseOrderBytes []byte
	PurchaseOrder := PurchaseOrder{}
	funcName := getFunctionName()

	/*Unmarshal the payload into a PurchaseOrder struct*/
	err = json.Unmarshal([]byte(args[0]), &PurchaseOrder)
	if err != nil {
		return shim.Error(funcName + " : Failed to convert arg[0] to a PurchaseOrder: " + err.Error())
	}

	/*=================================================================
		    Determine if the Purchase Order is already in the Ledger by:
				- Adding the From, To and OrderNumber to the keys, along with
				    type 'PO'
				- Doing a QueryObject to return the object and failing if we
				    receive anything in err, or if there are bytes returned
	    =================================================================
	*/
	keys := []string{PurchaseOrder.From, PurchaseOrder.To, PurchaseOrder.PONumber}

	PurchaseOrderType := "PO"
	Avalbytes, err = dbapi.QueryObject(stub, PurchaseOrderType, keys)
	if err != nil {
		return shim.Error(funcName + " :  Failed to query PurchaseOrder object")
	}

	if Avalbytes != nil {
		return shim.Error(fmt.Sprintf(funcName+" : "+
			"PurchaseOrder for PurchaseOrder Number: %s already exists ", PurchaseOrder.PONumber))
	}

	/*If the incoming Status  of the new Purchase Order is not OPEN, then reset it to OPEN */
	if PurchaseOrder.Status != OPEN {
		PurchaseOrder.Status = OPEN
	}

	/*Serialize the Purchase Order - as we have updated the STATUS - and store in the Ledger*/
	PurchaseOrderBytes, err = json.Marshal(PurchaseOrder)
	if err != nil {
		return shim.Error(funcName + " :  Failed to serialize PurchaseOrder object")
	}

	err = dbapi.UpdateObject(stub, PurchaseOrderType, keys, PurchaseOrderBytes)
	if err != nil {
		logger.Errorf(funcName+" :  Error inserting Object into LedgerState %s", err)
		return shim.Error(funcName + " :  POs Update failed")
	}

	return shim.Success(nil)

}

func createSalesOrder(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	var err error
	var Avalbytes []byte
	var SalesOrderBytes []byte
	var PurchaseOrderBytes []byte
	SalesOrder := SalesOrder{}
	funcName := getFunctionName()

	/*Unmarshal the payload into a PurchaseOrder struct*/
	err = json.Unmarshal([]byte(args[0]), &SalesOrder)
	if err != nil {
		return shim.Error(funcName + " : Failed to convert arg[0] to a SalesOrder: " + err.Error())
	}

	/*=================================================================
		    Determine if the Sales Order is already in the Ledger by:
				- Adding the From, To and OrderNumber to the keys, along with
				    type 'PO'
				- Doing a QueryObject to return the object and failing if we
				    receive anything in err, or if there are bytes returned
	    =================================================================*/
	keys := []string{SalesOrder.From, SalesOrder.To, SalesOrder.PONumber}

	SalesOrderType := "SO"
	Avalbytes, err = dbapi.QueryObject(stub, SalesOrderType, keys)
	if err != nil {
		return shim.Error(funcName + " : Failed to query SalesOrder object")
	}

	if Avalbytes != nil {
		return shim.Error(fmt.Sprintf(funcName+" : "+
			"SalesOrder for SalesOrder Number: %s already exists ", SalesOrder.PONumber))
	}

	/* Check the original purchase order based on the original information in the Sales Order*/
	PurchaseOrderType := "PO"
	poKeys := []string{SalesOrder.OriginalOrderer, SalesOrder.From, SalesOrder.OriginalPONumber}

	PurchaseOrderBytes, err = dbapi.QueryObject(stub, PurchaseOrderType, poKeys)

	if err != nil {
		return shim.Error(funcName + " : Failed to query original PurchaseOrder object")
	}

	if PurchaseOrderBytes == nil {
		return shim.Error(funcName + " : Failed to query original PurchaseOrder object")
	}

	/*If the incoming Status  of the new Purchase Order is not OPEN, then reset it to OPEN */
	if SalesOrder.Status != OPEN {
		SalesOrder.Status = OPEN
	}

	/*Serialize the Purchase Order - as we have updated the STATUS - and store in the Ledger*/
	SalesOrderBytes, err = json.Marshal(SalesOrder)
	if err != nil {
		return shim.Error(funcName + " : Failed to serialize SalesOrder object")
	}

	err = dbapi.UpdateObject(stub, SalesOrderType, keys, SalesOrderBytes)
	if err != nil {
		logger.Errorf(funcName+" : Error inserting Object into LedgerState %s", err)
		return shim.Error(funcName + " : POs Update failed")
	}

	return shim.Success(nil)

}

////////////////////////////////////////////////////////////////////////////
// ShipPart - not implemented for first version. Will be used to update
//            partial or invcremental shipping notices
////////////////////////////////////////////////////////////////////////////

func shipPart(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	return shim.Success(nil)
}

/*=================================================================
  createAcknowledgement - Used to store acknowledgements in the
	                        ledger.
================================================================= */

func createAcknowledgement(stub shim.ChaincodeStubInterface, args []string) pb.Response {

	var orderBytes []byte
	var ackBytes []byte
	acknowledgement := Acknowledgement{}
	funcName := getFunctionName()
	/*=================================================================
		  Unmarshal the Acknowledgement object and error if that fails
			or we do not receieve a payload
	    ================================================================= */
	err := json.Unmarshal([]byte(args[0]), &acknowledgement)

	if len(args) < 1 {
		return shim.Error(funcName + " : requires an Acknowledgement document")
	}
	logger.Infof(funcName+" : Arguments : %s", args[0])

	if err != nil {
		return shim.Error(funcName + " : Failed to convert arg[0] to an Acknolwdgement notice: " + err.Error())
	}

	/*=================================================================
		  we can receive different acks for different objects in upcoming
			releases, so we determine the object from the acks DocumentType
			and perform document specific ack processing
	  ================================================================= */
	switch acknowledgement.DocumentType {
	case "PO":
		/*=================================================================
			  Retrieve the Purchase Order using the keys from the Ack
				(swapping the From and To). We error if we cannot marshal or
				recieve no bytes
		   ================================================================= */
		keys := []string{acknowledgement.To, acknowledgement.From, acknowledgement.DocumentNumber}
		order, err := retrieveAndMarshalPOObject(stub, keys)
		if err != nil {
			return shim.Error(funcName + " : no existing Purchase Order Number " + acknowledgement.DocumentNumber)
		}

		if order.Status != OPEN {
			return shim.Error(funcName + " :  Sales Order Number " + acknowledgement.DocumentNumber + " not in OPEN state")
		}

		/*=================================================================
		  We have the Purchase Order, and we set the STATUS to 'ACKNOWLEDGED'
			The updated Purchase Order is updated back into the ledger
		================================================================= */
		order.Status = ACKNOWLEDGED
		orderBytes, err = json.Marshal(order)
		if err != nil {
			return shim.Error(funcName + " : failed to unmarshal existing Purchase Order Number " + acknowledgement.DocumentNumber)
		}
		err = dbapi.UpdateObject(stub, "PO", keys, orderBytes)
		if err != nil {
			return shim.Error(funcName + " : failed to update existing Purchase Order Number " + acknowledgement.DocumentNumber)
		}

		/*=================================================================
			  The Acknowledgement document is now stored in the ledger
				(errors on failure to store)
		   ================================================================= */
		aKeys := []string{acknowledgement.From, acknowledgement.To, acknowledgement.AckNumber}
		ackBytes, err = json.Marshal(acknowledgement)
		if err != nil {
			return shim.Error(funcName + " : failed to update existing Purchase Order Number " + acknowledgement.DocumentNumber)
		}
		dbapi.UpdateObject(stub, "ACK", aKeys, ackBytes)
	/*=================================================================
	  If the documentType is SO, then ack the Sales Order
		================================================================= */
	case "SO":
		/*=================================================================
			Retrieve the Purchase Order using the keys from the Ack
			(swapping the From and To). We error if we cannot marshal or
			recieve no bytes
		 ================================================================= */
		keys := []string{acknowledgement.To, acknowledgement.From, acknowledgement.DocumentNumber}
		order, err := retrieveAndMarshalPOObject(stub, keys)
		if err != nil {
			return shim.Error(funcName + " : no existing Sales Order Number " + acknowledgement.DocumentNumber)
		}

		if order.Status != OPEN {
			return shim.Error(funcName + " :  Sales Order Number " + acknowledgement.DocumentNumber + " not in OPEN state")
		}

		/*=================================================================
			We have the Purchase Order, and we set the STATUS to 'ACKNOWLEDGED'
			The updated Purchase Order is updated into the ledger
		================================================================= */
		order.Status = ACKNOWLEDGED
		orderBytes, err = json.Marshal(order)
		if err != nil {
			return shim.Error(funcName + " : failed to unmarshal existing Sales Order Number " + acknowledgement.DocumentNumber)
		}
		err = dbapi.UpdateObject(stub, "PO", keys, orderBytes)
		if err != nil {
			return shim.Error(funcName + " : failed to update existing Sales Order Number " + acknowledgement.DocumentNumber)
		}

		/*=================================================================
			The Acknowledgement document is now stored in the ledger
			(errors on failure to store)
		 ================================================================= */
		aKeys := []string{acknowledgement.From, acknowledgement.To, acknowledgement.AckNumber}
		ackBytes, err = json.Marshal(acknowledgement)
		if err != nil {
			return shim.Error(funcName + " : failed to update existing Purchase Order Number " + acknowledgement.DocumentNumber)
		}
		dbapi.UpdateObject(stub, "ACK", aKeys, ackBytes)
	default:
		return shim.Error(funcName + " : Ack for Doctype " + acknowledgement.DocumentType + " not yet implemented")
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
	funcName := getFunctionName()

	err = json.Unmarshal([]byte(args[0]), &returnNotice)
	if err != nil {
		return shim.Error(funcName + " :Failed to convert arg[0] to a Return notice: " + err.Error())
	}

	// Query and Retrieve the Full BaicInfo
	keys := []string{returnNotice.OrderNumber}

	objectType := "return"
	Avalbytes, err = dbapi.QueryObject(stub, objectType, keys)
	if err != nil {
		return shim.Error(funcName + " : Failed to query return PurchaseOrder object")
	}

	if Avalbytes != nil {
		return shim.Error(fmt.Sprintf(funcName+" : "+
			"ID for Return Invoice: %s already exist ", returnNotice.OrderNumber))
	}

	err = dbapi.UpdateObject(stub, objectType, keys, []byte(args[0]))
	if err != nil {
		logger.Errorf(funcName+" : Error inserting Object into LedgerState %s", err)
		return shim.Error(funcName + " : Return Notice Create failed")
	}

	return shim.Success(nil)
}
func createInvoice(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	var err error
	var Avalbytes []byte
	logger.Infof("CreateInvoice : Arguments : %s", args[0])
	invoice := Invoice{}
	funcName := getFunctionName()

	err = json.Unmarshal([]byte(args[0]), &invoice)
	if err != nil {
		return shim.Error(funcName + " :Failed to convert arg[0] to a Invoice object: " + err.Error())
	}

	// Query and Retrieve the Full BaicInfo
	keys := []string{invoice.OrderNumber}

	objectType := "invoice"
	Avalbytes, err = dbapi.QueryObject(stub, objectType, keys)
	if err != nil {
		return shim.Error(funcName + " : Failed to query shipment object")
	}

	if Avalbytes != nil {
		return shim.Error(fmt.Sprintf("CreateInvoice() : "+
			"ID for Invoice Number: %s already exist ", invoice.OrderNumber))
	}

	err = dbapi.UpdateObject(stub, objectType, keys, []byte(args[0]))
	if err != nil {
		logger.Errorf(funcName+" : Error inserting Object into LedgerState %s", err)
		return shim.Error(funcName + " : Invoice Create failed")
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
	funcName := getFunctionName()

	var i = 0

	if len(args) < 1 {
		logger.Infof(funcName + " : requires one parameter (originating company)")
		return shim.Error(funcName + " : requires one parameter (originating company)")
	}

	err = json.Unmarshal([]byte(args[0]), &query)

	if err != nil {
		logger.Infof(funcName + " :failed to marshal query object")
		return shim.Error(funcName + " : failed to marshal query object")
	}
	logger.Infof(funcName+" : Arguments : %s", query.Requestor)

	keys := []string{query.Requestor}

	rs, err := dbapi.GetList(stub, "PO", keys)

	if err != nil {
		logger.Infof(funcName+" : failed to retrieve PurchaseOrders: %s", keys[0])
		return shim.Error(funcName + " : failed to retrieve PurchaseOrders: %s)" + keys[0])
	}

	for i = 0; rs.HasNext(); i++ {
		myKV, err := rs.Next()
		if err != nil {
			return shim.Error(err.Error())
		}

		err = json.Unmarshal(myKV.Value, &PurchaseOrder)

		if err != nil {
			logger.Infof(funcName+" : failed to marshal PurchaseOrder: %s", err.Error())
			return shim.Error(funcName + " : failed to marshal PurchaseOrder: " + err.Error())
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
	funcName := getFunctionName()

	if len(args) < 1 {
		logger.Infof(funcName + " : requires request paramater")
		return shim.Error(funcName + " : requires request parameter")
	}
	logger.Infof(funcName+" : Arguments : %s", args[0])
	err = json.Unmarshal([]byte(args[0]), &query)

	if err != nil {
		logger.Infof(funcName+" : Arguments : %s", args[0])
	}

	keys := []string{query.Requestor, query.Partner, query.OrderNumber}
	PurchaseOrderbytes, err = dbapi.QueryObject(stub, "PO", keys)

	if err != nil {
		logger.Infof(funcName+" : failed to retrieve PurchaseOrder (PurchaseOrder number: %s, company %s )", query.OrderNumber, query.Requestor)
		return shim.Error(funcName + " : failed to retrieve PurchaseOrder")
	}

	return shim.Success(PurchaseOrderbytes)
}

func querySalesOrder(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	var err error
	var PurchaseOrderbytes []byte
	var query = QueryOrder{}
	funcName := getFunctionName()

	if len(args) < 1 {
		logger.Infof(funcName + " :requires request paramater")
		return shim.Error(funcName + " :requires request parameter")
	}
	logger.Infof(funcName+" : Arguments : %s", args[0])
	err = json.Unmarshal([]byte(args[0]), &query)

	if err != nil {
		logger.Infof(funcName+" : Arguments : %s", args[0])
	}

	keys := []string{query.Requestor, query.Partner, query.OrderNumber}
	PurchaseOrderbytes, err = dbapi.QueryObject(stub, "PO", keys)

	if err != nil {
		logger.Infof(funcName+" :fail to retrieve PurchaseOrder (PurchaseOrder number: %s, company %s )", query.OrderNumber, query.Requestor)
		return shim.Error(funcName + " : fail to retrieve PurchaseOrder")
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
	funcName := getFunctionName()

	logger.Infof("Received %s as arguments  ")

	if len(args) < 1 {
		logger.Infof(funcName + " :requires request paramater")
		return shim.Error(funcName + " :requires request parameter")
	}

	logger.Infof(funcName+" : Arguments : %s", args[0])

	keys := []string{args[0]}

	rs, err := dbapi.GetList(stub, "PO", keys)

	if err != nil {
		logger.Infof(funcName+" :fail to retrieve PurchaseOrders: %s", args[0])
		return shim.Error(funcName + " :fail to retrieve PurchaseOrders: )" + args[0])
	}

	for i = 0; rs.HasNext(); i++ {
		myKV, err := rs.Next()
		if err != nil {
			return shim.Error(err.Error())
		}

		err = json.Unmarshal(myKV.Value, &PurchaseOrder)

		if err != nil {
			logger.Infof(funcName+" :fail to marshal PurchaseOrder: %s", err.Error())
			return shim.Error(funcName + " :fail to marshal PurchaseOrder: " + err.Error())
		}

		if PurchaseOrder.From == query.Requestor || PurchaseOrder.To == query.Requestor {
			PurchaseOrders = append(PurchaseOrders, PurchaseOrder)
		}
	}

	jsonRows, err := json.Marshal(PurchaseOrders)
	return shim.Success(jsonRows)
}

func queryPurchaseOrderStatus(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	return shim.Success(nil)
}

////////////////////////////////////////////////////////////////////////////
//   Shipment Queries
////////////////////////////////////////////////////////////////////////////

////////////////////////////////////////////////////////////////////////////
// Query a specific Shipment with a full key
////////////////////////////////////////////////////////////////////////////
func queryShipment(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	var err error
	var Shipmentbytes []byte
	var query = QueryShipment{}
	funcName := getFunctionName()

	if len(args) < 1 {
		logger.Infof(funcName + " :requires request parameter")
		return shim.Error(funcName + " :requires request parameter")
	}
	logger.Infof(funcName+" :Arguments : %s", args[0])
	err = json.Unmarshal([]byte(args[0]), &query)

	if err != nil {
		logger.Infof(funcName+" : Arguments : %s", args[0])
	}

	keys := []string{query.Requestor, query.Partner, query.ShipmentNumber}
	Shipmentbytes, err = dbapi.QueryObject(stub, "SHP", keys)

	if err != nil {
		logger.Infof(funcName+" :fail to retrieve shipment (shipment number: %s, company %s )", query.ShipmentNumber, query.Requestor)
		return shim.Error(funcName + " :fail to retrieve shipment")
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
	funcName := getFunctionName()

	if len(args) < 1 {
		logger.Infof(funcName + " :requires request paramater")
		return shim.Error(funcName + " :requires request parameter")
	}

	logger.Infof(funcName+" :Arguments : %s", args[0])
	err = json.Unmarshal([]byte(args[0]), &query)

	if err != nil {
		logger.Infof(funcName+" :Arguments : %s, %s", args[0], err.Error())
		return shim.Error(funcName + " : error unmarshalling: " + err.Error())
	}

	logger.Infof("queryRichQuery : query name %s", query.QueryName)
	if len(query.QueryFields) < 1 {
		strAsBytes, _ := json.Marshal(query)
		logger.Infof(funcName+" :Arguments : %s", strAsBytes)
		return shim.Error(funcName + " : no fields requested")
	}

	queryfield = query.QueryFields[0].FieldValue

	formattedstring = getFormattedPurchaseOrderQuery(queryfield)
	logger.Infof(funcName+" : Query : %s", formattedstring)
	querybytes, err := dbapi.GetQueryResultForQueryString(stub, formattedstring)

	if err != nil {
		logger.Infof(funcName+" : fail to retrieve PurchaseOrders: %s", err.Error())
		return shim.Error(funcName + " : fail to retrieve PurchaseOrders: )" + err.Error())
	}

	return shim.Success(querybytes)
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
