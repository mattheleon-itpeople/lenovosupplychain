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

type InvokeFunc func(stub shim.ChaincodeStubInterface, args[] string) pb.Response

// SDMChaincode Chaincode implementation
type LenovoChainCode struct {
    tableMap map[string] int
    funcMap map[string] InvokeFunc
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
    OPEN = "open"
    CLOSED = "closed"
    REFUNDED = "refunded"
)

/////////////////////////////////////////////////////
// Constant for order type
/////////////////////////////////////////////////////
const (
    PURCHASE = "purchase"
    SALES = "sales"
)

/////////////////////////////////////////////////////
// Constant for All function name that will be called from invoke
/////////////////////////////////////////////////////
const (
<<<<<<< HEAD
    GV string = "getVersion"
    CPO string = "createPurchaseOrder"
    CSO string = "createShipment"
    SHPT string = "shipPart"
    CUR string = "createUser"
    UUR string = "updateUser"
    QUR string = "queryUser"
    DUR string = "deleteUser"
    QO string = "queryOrder"
    QOBN string = "queryOrderByOrderNumber"
)

func(t * LenovoChainCode) initMaps() {
    t.tableMap = make(map[string] int)
    t.tableMap[BIT] = 3
    t.funcMap = make(map[string] InvokeFunc)
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
=======
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
>>>>>>> d02955fd2757829de2f82dec5e91b55e1b638d45

}

////////////////////////////////////////////////////////////////////////////
// Query Version of BlockChain from Leadger
// This method for initial system test
////////////////////////////////////////////////////////////////////////////

func getVersion(stub shim.ChaincodeStubInterface, args[] string) pb.Response {
    logger.Info("-------- getVersion --------")
        // Get the cc version from the ledger
    version, err: = stub.GetState(VERSION)
    if err != nil {
        jsonResp: = "{\"Error\":\"Failed to get state for " + VERSION + "\"}"
        return shim.Error(jsonResp)
    }

    if version == nil {
        jsonResp: = "{\"Error\":\"" + VERSION + " is nil \"}"
        return shim.Error(jsonResp)
    }

    jsonResp: = "{\"Version\":\"" + string(version) + "\"}"
    logger.Infof("Query Response:%s\n", jsonResp)
    return shim.Success(version)
}

<<<<<<< HEAD
func createShipment(stub shim.ChaincodeStubInterface, args[] string) pb.Response {
    var err error
    var Avalbytes[] byte
    logger.Infof("CreateShipment : Arguments : %s", args[0])
    shipment: = Shipment {}
    err = json.Unmarshal([] byte(args[0]), & shipment)
    if err != nil {
        return shim.Error("CreateShipment : Failed to convert arg[0] to a Shipment object: " + err.Error())
    }

    // Query and Retrieve the Full BaicInfo
    keys: = [] string {
        shipment.ShipmentNumber
    }

    objectType: = "PO"
    Avalbytes, err = dbapi.QueryObject(stub, objectType, keys)
    if err != nil {
        return shim.Error("CreateShipment() : Failed to query shipment object")
    }

    if Avalbytes != nil {
        return shim.Error(fmt.Sprintf("CreateShipment() : " +
            "ID for Shipment Number: %s already exist ", shipment.ShipmentNumber))
    }

    err = dbapi.UpdateObject(stub, objectType, keys, [] byte(args[0]))
    if err != nil {
        logger.Errorf("CreateShipment : Error inserting Object into LedgerState %s", err)
        return shim.Error("CreateShipment : Shipment Update failed")
    }

    return shim.Success(nil)

}

func createOrder(stub shim.ChaincodeStubInterface, args[] string) pb.Response {
    var err error
    var Avalbytes[] byte
    logger.Infof("CreateOrder : Arguments : %s", args[0])
    Order: = Order {}
    err = json.Unmarshal([] byte(args[0]), & Order)
    if err != nil {
        return shim.Error("CreateOrder : Failed to convert arg[0] to a Order: " + err.Error())
    }

    // Query and Retrieve the Full BaicInfo
    keys: = [] string {
        Order.From, Order.To, Order.OrderNumber
    }

    orderType: = "PO"
    Avalbytes, err = dbapi.QueryObject(stub, orderType, keys)
    if err != nil {
        return shim.Error("CreateOrder() : Failed to query Order object")
    }

    if Avalbytes != nil {
        return shim.Error(fmt.Sprintf("CreateOrder() : " +
            "Order for Order Number: %s already exist ", Order.OrderNumber))
    }

    err = dbapi.UpdateObject(stub, orderType, keys, [] byte(args[0]))
    if err != nil {
        logger.Errorf("CreateOrder : Error inserting Object into LedgerState %s", err)
        return shim.Error("CreateOrder : POs Update failed")
    }

    return shim.Success(nil)

=======
func createShipment(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	var err error
	var Avalbytes []byte
	Shipment := Shipment{}
	Order := Order{}

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
	//TODO: QUERY ORDER LOOP THROUGH ORDER ITEM
	Avalbytes, err = dbapi.QueryObject(stub, objectType, keys)

	if err != nil {
		return shim.Error("Failed to retrieve order with provided shipping notice. " + err.Error())
	}
	if &Avalbytes == nil {
		return shim.Error("No order was retrieved. " + err.Error())
	}
	err = json.Unmarshal(Avalbytes, &Order)
	if err != nil {
		return shim.Error("Failed to marshal Sales Order. " + string(Avalbytes))
	}

	//items := Order.Items

	quantity := make(map[string][]int)

	for _, i := range Order.Items {
		quantity[i.PartNumber] = append(quantity[i.PartNumber], i.Quantity)
		fmt.Println(quantity[i.PartNumber])
	}

	/*
	   if len(Order.orderLine) != len(Shipment.ShippedItems) {
	       return shim.Error("***** Order quantity does not match shipping quantity. Changing order status to: pending review. *****")
	   }

	   for iterator < len(Order.orderLine) {

	       // TAKE ITERATOR CREATE MAP OF BOTH SIDES AND COMPARE

	       //orderQuantity := Order.quantity
	       //shipQuantity := Shipment.quantity
	           iterator += iterator
	       }
	       iterator += iterator


	   if orderQuantity != shipQuantity {

	   }
	*/

	err = dbapi.UpdateObject(stub, objectType, keys, []byte(args[0]))
	if err != nil {
		logger.Errorf("shipPart : Error inserting Shipment of parts into LedgerState %s", err)
		return shim.Error("shipPart : Shipping part failed")
	}

	return shim.Success(nil)
>>>>>>> d02955fd2757829de2f82dec5e91b55e1b638d45
}

/********************************************************************************************
 * 1 order 1 full shipment only quantity													*
 * TODO: Check price, quantity, delivery date, partial or full shipment, unit of measure	*																							*
 *	if len(orderline) len(ship) for(part number & orderline compare quant)																						*
 *																							*
 * Verify shipping notice against order - Send ASN and ship part							*
 * Give these arguments to the key array in this order:										*
 * - from -> Shipment from field															*	
 * - to -> Shipment to field																*
 * - orderNumber -> Shipment order number field												*
 * - shipQuantity -> Shipment quantity amount												*		
 *																							*
 * Query object using "SHP" object type and retrieve values based on:						*
 * @param = Shipment.OrderNumber															*
 * 																							*
 * Retrieve order json object and extrapolate quantity value								*
 * Validate quantity is equivalent to shipping notice and proceed with ASN and shipment.    *
 * Using UpdateObject from dbapi, write the ACK into the ledger. Otherwise return Errors.	*
 ********************************************************************************************/
func shipPart(stub shim.ChaincodeStubInterface, args[] string) pb.Response {
	var err error
    var Avalbytes[] byte
    shippingNotice: = ShippingNotice {}
    err = json.Unmarshal([] byte(args[0]), &ShippingNotice)
    if err != nil {
        return shim.Error("Failed to retrieve shipping notice with provided order number. " + err.Error())
	}

	from := Shipment.From
	to := Shipment.To
	orderNumber := Shipment.OrderNumber

	iterator := 0
	

    keys: = [] string {
		from,
		to,
		orderNumber
    }

	objectType: = "PO"
	//TODO: QUERY ORDER LOOP THROUGH ORDER ITEM
	Avalbytes, err = dbapi.QueryObject(stub, objectType, keys)
	err = json.Unmarshal(byte(args[0]), &Order)
	if err != nil {
        return shim.Error("Failed to retrieve order with provided shipping notice. " + err.Error())
	}
	
	Orders := dbapi.getList(stub, objectType, keys)

	quantity := make(map[string][]*Orders)
	for _, i := range Order {
			quantity[i.partNumber] = append(quantity[i.partNumber], i.quantity)
			fmt.println(quantity[i.partNumber])
	}

	/*
	if len(Order.orderLine) != len(Shipment.ShippedItems) {
		return shim.Error("***** Order quantity does not match shipping quantity. Changing order status to: pending review. *****")
	}
	
	for iterator < len(Order.orderLine) {
		
		// TAKE ITERATOR CREATE MAP OF BOTH SIDES AND COMPARE

		//orderQuantity := Order.quantity
		//shipQuantity := Shipment.quantity
			iterator += iterator
		}
		iterator += iterator


	if orderQuantity != shipQuantity {
       
	}
	*/
	
	err = dbapi.UpdateObject(stub, objectType, keys, [] byte(args[0]))
    if err != nil {
        logger.Errorf("shipPart : Error inserting Shipment of parts into LedgerState %s", err)
        return shim.Error("shipPart : Shipping part failed")
    }
	
    return shim.Success(nil)
}


/********************************************************************************************
 * Sends an acknowledgement upon the recieving of a particular Purchase Order.				*
 * Give these arguments to the args array in this order:									*
 * - args[0] -> order number																*	
 * - args[1] -> from whom																	*
 * - args[2] -> to whom																		*		
 *																							*
 * Using UpdateObject from dbapi, write the ACK into the ledger. Otherwise return Errors.	*
 ********************************************************************************************/
func sendAcknowledgement(stub shim.ChaincodeStubInterface, args[] string) pb.Response {
	// TODO: Pass order as object using query instead of args

	// check valid numnber of args
	if len(args) != 3 {
		return shim.Error("Incorrect number of arguments, expecting 3. ('Order Number', 'From', 'To')")
	}

	// extract args
	from := args[0]
	to := args[1]
	orderNum := args[2]

<<<<<<< HEAD
	// assign keys from args array
	keys = [] string {
		from, 
		to,
		orderNum
=======
	objectType := "return"
	Avalbytes, err = dbapi.QueryObject(stub, objectType, keys)
	if err != nil {
		return shim.Error("CreateReturnNotice() : Failed to query return order object")
>>>>>>> d02955fd2757829de2f82dec5e91b55e1b638d45
	}

	// pass in object type - invoke dbapi's UpdateObject query and validate
	objectType = "PO"
	err = dbapi.UpdateObject(stub, objectType, keys, [] byte(args[0]))
    if err != nil {
        logger.Errorf("sendAcknowledgement : Error inserting ACK into LedgerState %s", err)
        return shim.Error("sendAcknowledgement : Send ACK failed")
    }

    return shim.Success(nil)

    
}
func generateShippingNotice(stub shim.ChaincodeStubInterface, args[] string) pb.Response {
	var err error
<<<<<<< HEAD
    var Avalbytes[] byte
    logger.Infof("generateShippingNotice : Arguments : %s", args[0])
    shippingNotice: = ShippingNotice {}
    err = json.Unmarshal([] byte(args[0]), & ShippingNotice)
    if err != nil {
        return shim.Error("generateShippingNotice : Failed to convert arg[0] to a ACK object: " + err.Error())
    }

    keys: = [] string {
        Order.OrderNumber
    }

    objectType: = "PO"
    Avalbytes, err = dbapi.QueryObject(stub, objectType, keys)

    if err != nil {
        return shim.Error("Order does not exist or was not invoiced yet.")
    }

    if Avalbytes != nil {
        return shim.Error(fmt.Sprintf("generateShippingNotice() : " +
            "Acknolwedgement of Order was already sent ", Order.OrderNumber))
=======
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
>>>>>>> d02955fd2757829de2f82dec5e91b55e1b638d45
	}
	
	err = dbapi.UpdateObject(stub, objectType, keys, [] byte(args[0]))
    if err != nil {
        logger.Errorf("generateShippingNotice : Error inserting ACK into LedgerState %s", err)
        return shim.Error("generateShippingNotice : Send Shipping Notice failed")
    }
	
    return shim.Success(nil)
}
func generateGoodsReceived(stub shim.ChaincodeStubInterface, args[] string) pb.Response {
	var err error
    var Avalbytes[] byte
    logger.Infof("generateGoodsRecieved : Arguments : %s", args[0])
    goodsRecieved: = GoodsRecieved {}
    err = json.Unmarshal([] byte(args[0]), & goodsRecieved)
    if err != nil {
        return shim.Error("generateGoodsRecieved : Failed to convert arg[0] to a good recieved object: " + err.Error())
    }

    keys: = [] string {
        Shipment.OrderNumber
    }

    objectType: = "SPO"
    Avalbytes, err = dbapi.QueryObject(stub, objectType, keys)

    if err != nil {
        return shim.Error("Order does not exist or was not delivered.")
    }

    if Avalbytes != nil {
        return shim.Error(fmt.Sprintf("sendAcknowledgement() : " +
            "Goods recieved notice was already sent ", Shipment.OrderNumber))
	}
	
	err = dbapi.UpdateObject(stub, objectType, keys, [] byte(args[0]))
    if err != nil {
        logger.Errorf("sendAcknowledgement : Error inserting goods recieved into LedgerState %s", err)
        return shim.Error("sendAcknowledgement : Send goods recieved notice failed")
    }

    return shim.Success(nil)
}
func createReturnNotice(stub shim.ChaincodeStubInterface, args[] string) pb.Response {
    var err error
    var Avalbytes[] byte
    logger.Infof("CreateReturnNotice : Arguments : %s", args[0])
    returnNotice: = ReturnNotice {}
    err = json.Unmarshal([] byte(args[0]), & returnNotice)
    if err != nil {
        return shim.Error("CreateReturnNotice : Failed to convert arg[0] to a Return notice: " + err.Error())
    }

    // Query and Retrieve the Full BaicInfo
    keys: = [] string {
        returnNotice.OrderNumber
    }

    objectType: = "RET"
    Avalbytes, err = dbapi.QueryObject(stub, objectType, keys)
    if err != nil {
        return shim.Error("CreateReturnNotice() : Failed to query return order object")
    }

    if Avalbytes != nil {
        return shim.Error(fmt.Sprintf("CreateReturnNotice() : " +
            "ID for Return Invoice: %s already exist ", returnNotice.OrderNumber))
    }

    err = dbapi.UpdateObject(stub, objectType, keys, [] byte(args[0]))
    if err != nil {
        logger.Errorf("CreateReturnNotice : Error inserting Object into LedgerState %s", err)
        return shim.Error("CreateReturnNotice : Return Notice Create failed")
    }

    return shim.Success(nil)
}
func createInvoice(stub shim.ChaincodeStubInterface, args[] string) pb.Response {
    var err error
    var Avalbytes[] byte
    logger.Infof("CreateInvoice : Arguments : %s", args[0])
    invoice: = Invoice {}
    err = json.Unmarshal([] byte(args[0]), & invoice)
    if err != nil {
        return shim.Error("CreateInvoice : Failed to convert arg[0] to a Invoice object: " + err.Error())
    }

    // Query and Retrieve the Full BaicInfo
    keys: = [] string {
        invoice.OrderNumber
    }

    objectType: = "INV"
    Avalbytes, err = dbapi.QueryObject(stub, objectType, keys)
    if err != nil {
        return shim.Error("CreateInvoice() : Failed to query invoice")
    }

    if Avalbytes != nil {
        return shim.Error(fmt.Sprintf("CreateInvoice() : " +
            "ID for Invoice Number: %s already exist ", invoice.OrderNumber))
    }

    err = dbapi.UpdateObject(stub, objectType, keys, [] byte(args[0]))
    if err != nil {
        logger.Errorf("CreateInvoice : Error inserting Object into LedgerState %s", err)
        return shim.Error("CreateInvoice : Invoice Create failed")
    }

    return shim.Success(nil)
}
func sendPayment(stub shim.ChaincodeStubInterface, args[] string) pb.Response {
    return shim.Success(nil)
}

////////////////////////////////////////////////////////////////////////////
// Query Function
////////////////////////////////////////////////////////////////////////////
////////////////////////////////////////////////////////////////////////////
// Query Order given the Order Number and the 'From' organization
////////////////////////////////////////////////////////////////////////////
<<<<<<< HEAD
func queryOrderByOrderNumber(stub shim.ChaincodeStubInterface, args[] string) pb.Response {
    var query = QueryOrder {}
    var i = 0

    if len(args) < 1 {
        logger.Infof("queryOrderByOrderNumber requires request paramater")
        return shim.Error("queryOrderByOrderNumber requires request parameter")
    }

    err: = json.Unmarshal([] byte(args[0]), & query)

    if err != nil {
        logger.Infof("queryOrderByOrderNumber failed to unmarshal data :" + err.Error())
        return shim.Error("queryOrderByOrderNumber failed to unmarshal data : " + err.Error())
    }

    logger.Infof("queryOrderByOrderNumber : Arguments :" + query.OrderNumber + " : " + query.From)
    keys: = [] string {
        query.From
    }
    results, err: = dbapi.GetList(stub, "PO", keys)
    logger.Info("QueryByGetQuery - returned from dbapi")

    for i = 0;
    results.HasNext();
    i++{
        logger.Info("QueryByGetQuery - Iterating")
            // Retrieve the Key and Object
        myCompositeKey, err: = results.Next()
        if err != nil {
            return shim.Error(err.Error())
        }
        logger.Infof("QueryOrderByOrderNumber() : my Value : ", myCompositeKey)
    }
    return shim.Success(nil)
=======
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
>>>>>>> d02955fd2757829de2f82dec5e91b55e1b638d45
}

////////////////////////////////////////////////////////////////////////////
// Query a specific Order with a full key
////////////////////////////////////////////////////////////////////////////
func queryOrder(stub shim.ChaincodeStubInterface, args[] string) pb.Response {
    var err error
    var Orderbytes[] byte
    var query = QueryOrder {}

<<<<<<< HEAD
    logger.Infof("Received %s as arguments  ")

    if len(args) < 1 {
        logger.Infof("queryOrder requires request paramater")
        return shim.Error("queryOrder requires request parameter")
    }

    err = json.Unmarshal([] byte(args[0]), & query)
    logger.Infof("queryOrder : Arguments : %s", args[0])

    keys: = [] string {
        query.From, query.To, query.OrderNumber
    }
    Orderbytes, err = dbapi.QueryObject(stub, "PO", keys)

    if err != nil {
        logger.Infof("queryOrder fail to retrieve order (order number: %s, company %s )", query.OrderNumber, query.From)
        return shim.Error("queryOrder fail to retrieve order")
    }
=======
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
>>>>>>> d02955fd2757829de2f82dec5e91b55e1b638d45

    return shim.Success(Orderbytes)
}

<<<<<<< HEAD
func queryAllOrders(stub shim.ChaincodeStubInterface, args[] string) pb.Response {
    var err error
    var orders[] Order
    var order Order = Order {}
    var i = 0

    if len(args) < 1 {
        logger.Infof("queryOrder requires one parameter (originating company)")
        return shim.Error("queryOrder requires one parameter (originating company)")
    }
    logger.Infof("queryOrder : Arguments : %s", args[0])

    keys: = [] string {
        args[0]
    }

    rs, err: = dbapi.GetList(stub, "PO", keys)

    if err != nil {
        logger.Infof("queryOrder fail to retrieve orders: %s", args[0])
        return shim.Error("queryOrder fail to retrieve orders: )" + args[0])
    }

    for i = 0;
    rs.HasNext();
    i++{
        myKV, err: = rs.Next()
        if err != nil {
            return shim.Error(err.Error())
        }

        err = json.Unmarshal(myKV.Value, & order)

        if err != nil {
            logger.Infof("queryOrder fail to marshal order: %s", err.Error())
            return shim.Error("queryOrder fail to marshal order: " + err.Error())
        }

        orders = append(orders, order)
    }

    jsonRows, err: = json.Marshal(orders)
    return shim.Success(jsonRows)
}

func queryOrderStatus(stub shim.ChaincodeStubInterface, args[] string) pb.Response {
    return shim.Success(nil)
=======
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
>>>>>>> d02955fd2757829de2f82dec5e91b55e1b638d45
}
func queryShipmentByOrderNumber(stub shim.ChaincodeStubInterface, args[] string) pb.Response {
    return shim.Success(nil)
}

////////////////////////////////////////////////////////////////////////////
//Query by selector (rich query!)
////////////////////////////////////////////////////////////////////////////

////////////////////////////////////////////////////////////////////////////
// Helper functions
////////////////////////////////////////////////////////////////////////////

<<<<<<< HEAD
func getFormattedOrderQuery(orderNumber string, From string) string {
    return fmt.Sprintf("{\"selector\": { \"orderNumber\": \"%s\"}}", orderNumber)
=======
func getFormattedOrderQuery(orderNumber string) string {
	return fmt.Sprintf("{\\\"selector\\\": { \\\"orderNumber\\\": \\\"%s\\\"}}", orderNumber)
>>>>>>> d02955fd2757829de2f82dec5e91b55e1b638d45
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