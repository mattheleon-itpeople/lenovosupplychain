/******************************************************************
Copyright IT People Corp. 2017 All Rights Reserved.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

                 http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.

******************************************************************/

///////////////////////////////////////////////////////////////////////
// chaincode unit tests
///////////////////////////////////////////////////////////////////////
package main

import (
	"fmt"
	"testing"

	"github.com/hyperledger/fabric/core/chaincode/shim"
)

func checkInit(t *testing.T, stub *shim.MockStub, args [][]byte) {
	res := stub.MockInit("1", args)
	if res.Status != shim.OK {
		fmt.Println("Init failed", string(res.Message))
		t.FailNow()
	}
}

func checkState(t *testing.T, stub *shim.MockStub, name string, value string) {
	bytes := stub.State[name]
	if bytes == nil {
		fmt.Println("State", name, "failed to get value")
		t.FailNow()
	}
	if string(bytes) != value {
		fmt.Println("State value", name, "was not", value, "as expected")
		t.FailNow()
	}
}

func checkQuery(t *testing.T, stub *shim.MockStub, fcn string, name string, value string) {
	res := stub.MockInvoke("1", [][]byte{[]byte(fcn), []byte(name)})
	if res.Status != shim.OK {
		fmt.Println("Query ", name, "failed", string(res.Message))
		t.FailNow()
	}
	if res.Payload == nil {
		fmt.Println("Query ", name, "failed to get value")
		t.FailNow()
	}

	if string(res.Payload) != value {
		fmt.Println("Query value", name, "was not", value, "as expected")
		fmt.Println("Payload : " + string(res.Payload))
		t.FailNow()
	}
}

func checkInvoke(t *testing.T, stub *shim.MockStub, args [][]byte) {
	res := stub.MockInvoke("1", args)
	if res.Status != shim.OK {
		fmt.Println("Invoke", string(args[0]), "failed", string(res.Message))
		t.FailNow()
	}
}

func TestSCM_Init(t *testing.T) {
	lcc := new(LenovoChainCode)
	stub := shim.NewMockStub("ldm", lcc)

	// Init cc_version=v0
	checkInit(t, stub, [][]byte{[]byte("init"), []byte("v0")})

	checkState(t, stub, "version", "v0")
}

func TestSCM_query_ccversion(t *testing.T) {
	lcc := new(LenovoChainCode)
	stub := shim.NewMockStub("ldm", lcc)

	// Init cc_version=v0
	checkInit(t, stub, [][]byte{[]byte("init"), []byte("v0")})

	// query cc version
	checkQuery(t, stub, "getVersion", "version", "v0")

}

func TestSCM_invoke_createPurchaseOrder(t *testing.T) {
	lcc := new(LenovoChainCode)
	stub := shim.NewMockStub("ldm", lcc)

	// Init cc_version=v0
	checkInit(t, stub, [][]byte{[]byte("init"), []byte("v0")})

	// createSupplierBasicInfo for ITPC organization
	checkInvoke(t, stub, [][]byte{[]byte("createPurchaseOrder"), []byte("{\"doctype\": \"PO\",\"ponumber\": \"1234\",\"supplierid\": \"supplier\",\"vendordescription\": \"Lenovo Laptop Builder\",\"from\": \"Manu1\", \"to\": \"Lenovo\", \"items\": [{\"commoditycode\": \"1234\",\"unitprice\": \"12.50\",\"uom\": \"EA\",\"shorttext\": \"widget1\",\"orderedquantity\": \"100\",\"orderedvalue\": \"1250.00\"},{\"commoditycode\": \"1235\",\"unitprice\": \"15.50\",\"uom\": \"EA\",\"shorttext\": \"widget2\",\"orderedquantity\": \"100\",\"orderedvalue\": \"1500.00\"}]}")})

	// validate supplier details of org ITPC with querySupplierBasicInfo
	checkQuery(t, stub, "queryPurchaseOrder", "{\"orderNumber\": \"1234\", \"requestor\": \"Manu1\", \"partner\": \"Lenovo\"}", "{\"doctype\": \"PO\",\"ponumber\": \"1234\",\"supplierid\": \"supplier\",\"vendordescription\": \"Lenovo Laptop Builder\",\"from\": \"Manu1\", \"to\": \"Lenovo\", \"items\": [{\"commoditycode\": \"1234\",\"unitprice\": \"12.50\",\"uom\": \"EA\",\"shorttext\": \"widget1\",\"orderedquantity\": \"100\",\"orderedvalue\": \"1250.00\"},{\"commoditycode\": \"1235\",\"unitprice\": \"15.50\",\"uom\": \"EA\",\"shorttext\": \"widget2\",\"orderedquantity\": \"100\",\"orderedvalue\": \"1500.00\"}]}")

	// validate supplier details of org ITPC with querySupplierBasicInfo
	checkQuery(t, stub, "queryOrderByOrderNumber", "{\"orderNumber\":\"1234\", \"requestor\":\"Manu1\"}", "[{\"doctype\":\"PO\",\"ponumber\":\"1234\",\"supplierid\":\"supplier\",\"vendordescription\":\"Lenovo Laptop Builder\",\"poissuedate\":\"\",\"postartdate\":\"\",\"poenddate\":\"\",\"paymentterm\":\"\",\"paymentdays\":\"\",\"sownum\":\"\",\"currency\":\"\",\"from\":\"Manu1\",\"to\":\"Lenovo\",\"items\":[{\"commoditycode\":\"1234\",\"effectivedate\":\"\",\"unitprice\":\"12.50\",\"oldprice\":\"\",\"uom\":\"EA\",\"shorttext\":\"widget1\",\"orderedquantity\":\"100\",\"orderedvalue\":\"1250.00\",\"invoicedquantity\":\"\",\"invoicedvalue\":\"\",\"tobedeliveredvalue\":\"\",\"tobedeliveredquantity\":\"\"},{\"commoditycode\":\"1235\",\"effectivedate\":\"\",\"unitprice\":\"15.50\",\"oldprice\":\"\",\"uom\":\"EA\",\"shorttext\":\"widget2\",\"orderedquantity\":\"100\",\"orderedvalue\":\"1500.00\",\"invoicedquantity\":\"\",\"invoicedvalue\":\"\",\"tobedeliveredvalue\":\"\",\"tobedeliveredquantity\":\"\"}],\"auditInfo\":{\"createdBy\":\"\",\"updatedBy\":\"\",\"createdTS\":\"\",\"updatedTS\":\"\"},\"status\":\"\"}]")

}

func TestSCM_invoke_createAcknowledgement(t *testing.T) {
	lcc := new(LenovoChainCode)
	stub := shim.NewMockStub("ldm", lcc)

	// Init cc_version=v0
	checkInit(t, stub, [][]byte{[]byte("init"), []byte("v0")})

	// createSupplierBasicInfo for ITPC organization
	checkInvoke(t, stub, [][]byte{[]byte("createPurchaseOrder"), []byte("{\"doctype\": \"PO\",\"ponumber\": \"1234\",\"supplierid\": \"supplier\",\"vendordescription\": \"Lenovo Laptop Builder\",\"from\": \"Manu1\", \"to\": \"Lenovo\", \"items\": [{\"commoditycode\": \"1234\",\"unitprice\": \"12.50\",\"uom\": \"EA\",\"shorttext\": \"widget1\",\"orderedquantity\": \"100\",\"orderedvalue\": \"1250.00\"},{\"commoditycode\": \"1235\",\"unitprice\": \"15.50\",\"uom\": \"EA\",\"shorttext\": \"widget2\",\"orderedquantity\": \"100\",\"orderedvalue\": \"1500.00\"}]}")})

	// validate supplier details of org ITPC with querySupplierBasicInfo
	checkQuery(t, stub, "queryPurchaseOrder", "{\"orderNumber\": \"1234\", \"requestor\": \"Manu1\", \"partner\": \"Lenovo\"}", "{\"doctype\": \"PO\",\"ponumber\": \"1234\",\"supplierid\": \"supplier\",\"vendordescription\": \"Lenovo Laptop Builder\",\"from\": \"Manu1\", \"to\": \"Lenovo\", \"items\": [{\"commoditycode\": \"1234\",\"unitprice\": \"12.50\",\"uom\": \"EA\",\"shorttext\": \"widget1\",\"orderedquantity\": \"100\",\"orderedvalue\": \"1250.00\"},{\"commoditycode\": \"1235\",\"unitprice\": \"15.50\",\"uom\": \"EA\",\"shorttext\": \"widget2\",\"orderedquantity\": \"100\",\"orderedvalue\": \"1500.00\"}]}")

	checkInvoke(t, stub, [][]byte{[]byte("createAcknowledgement"), []byte("{\"doctype\":\"Acknowledgement\",\"documentType\":\"PO\",\"documentNumber\":\"1234\",\"from\":\"Lenovo\",\"to\":\"Manu1\"}")})

	checkQuery(t, stub, "queryPurchaseOrder", "{\"orderNumber\": \"1234\", \"requestor\": \"Manu1\", \"partner\": \"Lenovo\"}", "{\"doctype\":\"PO\",\"ponumber\":\"1234\",\"supplierid\":\"supplier\",\"vendordescription\":\"Lenovo Laptop Builder\",\"poissuedate\":\"\",\"postartdate\":\"\",\"poenddate\":\"\",\"paymentterm\":\"\",\"paymentdays\":\"\",\"sownum\":\"\",\"currency\":\"\",\"from\":\"Manu1\",\"to\":\"Lenovo\",\"items\":[{\"commoditycode\":\"1234\",\"effectivedate\":\"\",\"unitprice\":\"12.50\",\"oldprice\":\"\",\"uom\":\"EA\",\"shorttext\":\"widget1\",\"orderedquantity\":\"100\",\"orderedvalue\":\"1250.00\",\"invoicedquantity\":\"\",\"invoicedvalue\":\"\",\"tobedeliveredvalue\":\"\",\"tobedeliveredquantity\":\"\"},{\"commoditycode\":\"1235\",\"effectivedate\":\"\",\"unitprice\":\"15.50\",\"oldprice\":\"\",\"uom\":\"EA\",\"shorttext\":\"widget2\",\"orderedquantity\":\"100\",\"orderedvalue\":\"1500.00\",\"invoicedquantity\":\"\",\"invoicedvalue\":\"\",\"tobedeliveredvalue\":\"\",\"tobedeliveredquantity\":\"\"}],\"auditInfo\":{\"createdBy\":\"\",\"updatedBy\":\"\",\"createdTS\":\"\",\"updatedTS\":\"\"},\"status\":\"acknowledged\"}")

}

func TestSCM_invoke_createShipment(t *testing.T) {
	lcc := new(LenovoChainCode)
	stub := shim.NewMockStub("ldm", lcc)

	// Init cc_version=v0
	checkInit(t, stub, [][]byte{[]byte("init"), []byte("v0")})

	// createSupplierBasicInfo for ITPC organization
	checkInvoke(t, stub, [][]byte{[]byte("createPurchaseOrder"), []byte("{\"doctype\": \"PO\",\"ponumber\": \"1234\",\"supplierid\": \"supplier\",\"vendordescription\": \"Lenovo Laptop Builder\",\"from\": \"Manu1\", \"to\": \"Lenovo\", \"items\": [{\"commoditycode\": \"1234\",\"unitprice\": \"12.50\",\"uom\": \"EA\",\"shorttext\": \"widget1\",\"orderedquantity\": \"100\",\"orderedvalue\": \"1250.00\"},{\"commoditycode\": \"1235\",\"unitprice\": \"15.50\",\"uom\": \"EA\",\"shorttext\": \"widget2\",\"orderedquantity\": \"100\",\"orderedvalue\": \"1500.00\"}]}")})

	checkQuery(t, stub, "queryPurchaseOrder", "{\"orderNumber\": \"1234\", \"requestor\": \"Manu1\", \"partner\": \"Lenovo\"}", "{\"doctype\": \"PO\",\"ponumber\": \"1234\",\"supplierid\": \"supplier\",\"vendordescription\": \"Lenovo Laptop Builder\",\"from\": \"Manu1\", \"to\": \"Lenovo\", \"items\": [{\"commoditycode\": \"1234\",\"unitprice\": \"12.50\",\"uom\": \"EA\",\"shorttext\": \"widget1\",\"orderedquantity\": \"100\",\"orderedvalue\": \"1250.00\"},{\"commoditycode\": \"1235\",\"unitprice\": \"15.50\",\"uom\": \"EA\",\"shorttext\": \"widget2\",\"orderedquantity\": \"100\",\"orderedvalue\": \"1500.00\"}]}")

	// createSupplierBasicInfo for ITPC organization
	checkInvoke(t, stub, [][]byte{[]byte("createShipment"), []byte("{\"shipmentNumber\": \"1234\",\"trackingnumber\": \"4567\", \"supplierId\": \"supid1\", \"ordernumber\": \"0001\", \"from\": \"Lenovo\", \"to\": \"Manu1\"}")})
	// validate supplier details of org ITPC with querySupplierBasicInfo
	//checkQuery(t, stub, "querySupplierBasicInfo", "{\"Orgname\": \"ITPC\"}", "{\"Orgname\": \"ITPC\",\"Requestedby\": \"Lenovo\",\"Providedby\": \"IBM\",\"address\": {\"street\": \"11,abcd dr\",\"zip\": \"33647\",\"city\": \"Tampa\",\"country\": \"USA\",\"state\": \"Florida\",\"timezone\": \"EST\"},\"contacts\": [{\"type\": \"mobile\",\"cvalue\": \"+1-813-499-3389\"}, {\"type\": \"Email\",\"cvalue\": \"abc@gmail.com\"}],\"orgtype\": \"0\",\"hashedbuyerinfo\": \"\",\"hashedsupinfo\": \"\"}")
}

// func TestSDM_invoke_updateSupplier(t *testing.T) {
// 	scc := new(SDMChaincode)
// 	stub := shim.NewMockStub("sdm", scc)
//
// 	// Init cc_version=v0
// 	checkInit(t, stub, [][]byte{[]byte("init"), []byte("v0")})
//
// 	// createSupplierBasicInfo for ITPC organization
// 	checkInvoke(t, stub, [][]byte{[]byte("createSupplierBasicInfo"), []byte("{\"Orgname\": \"ITPeople\",\"Requestedby\": \"Lenovo\",\"Providedby\": \"IBM\",\"address\": {\"street\": \"11,abcd dr\",\"zip\": \"33647\",\"city\": \"Tampa\",\"country\": \"USA\",\"state\": \"Florida\",\"timezone\": \"EST\"},\"contacts\": [{\"type\": \"mobile\",\"cvalue\": \"+1-813-499-3389\"}, {\"type\": \"Email\",\"cvalue\": \"abc@gmail.com\"}],\"orgtype\": \"0\",\"hashedbuyerinfo\": \"\",\"hashedsupinfo\": \"\"}")})
//
// 	// createSupplierBasicInfo for ITPC organization
// 	checkInvoke(t, stub, [][]byte{[]byte("updateSupplierBasicInfo"), []byte("{\"Orgname\": \"ITPeople\",\"Requestedby\": \"Fedility\",\"Providedby\": \"XYZcorp\",\"address\": {\"street\": \"11,abcd dr\",\"zip\": \"33647\",\"city\": \"Tampa\",\"country\": \"USA\",\"state\": \"Florida\",\"timezone\": \"EST\"},\"contacts\": [{\"type\": \"mobile\",\"cvalue\": \"+1-813-499-3389\"}, {\"type\": \"Email\",\"cvalue\": \"abc@gmail.com\"}],\"orgtype\": \"0\",\"hashedbuyerinfo\": \"\",\"hashedsupinfo\": \"\"}")})
//
// 	// validate supplier details of org ITPC with querySupplierBasicInfo
// 	checkQuery(t, stub, "querySupplierBasicInfo", "{\"Orgname\": \"ITPeople\"}", "{\"Orgname\": \"ITPeople\",\"Requestedby\": \"Fedility\",\"Providedby\": \"XYZcorp\",\"address\": {\"street\": \"11,abcd dr\",\"zip\": \"33647\",\"city\": \"Tampa\",\"country\": \"USA\",\"state\": \"Florida\",\"timezone\": \"EST\"},\"contacts\": [{\"type\": \"mobile\",\"cvalue\": \"+1-813-499-3389\"}, {\"type\": \"Email\",\"cvalue\": \"abc@gmail.com\"}],\"orgtype\": \"0\",\"hashedbuyerinfo\": \"\",\"hashedsupinfo\": \"\"}")
// }
