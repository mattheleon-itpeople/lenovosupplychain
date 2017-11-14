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
		t.FailNow()
	}
}

func checkInvoke(t *testing.T, stub *shim.MockStub, args [][]byte) {
	res := stub.MockInvoke("1", args)
	if res.Status != shim.OK {
		fmt.Println("Invoke", args, "failed", string(res.Message))
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
	checkInvoke(t, stub, [][]byte{[]byte("createOrder"), []byte("{\"orderNumber\": \"1234\",\"supplierID\": \"Manu1\", \"items\": [{\"partNumber\":\"1111\",\"quantity\":100},{\"partNumber\": \"2222\", \"quantity\": 200}], \"From\": \"Manu1\", \"To\": \"Lenovo\"}")})

	// validate supplier details of org ITPC with querySupplierBasicInfo
	checkQuery(t, stub, "queryOrder", "{\"orderNumber\": \"1234\", \"Requestor\": \"Manu1\", \"Partner\": \"Lenovo\"}", "{\"orderNumber\": \"1234\",\"supplierID\": \"Manu1\", \"items\": [{\"partNumber\":\"1111\",\"quantity\":100},{\"partNumber\": \"2222\", \"quantity\": 200}], \"From\": \"Manu1\", \"To\": \"Lenovo\"}")

	// validate supplier details of org ITPC with querySupplierBasicInfo
	checkQuery(t, stub, "queryOrderByOrderNumber", "{\"orderNumber\":\"1234\", \"Requestor\":\"Manu1\"}", "[{\"orderNumber\":\"1234\",\"supplierId\":\"Manu1\",\"items\":[{\"partNumber\":\"1111\",\"itemCondition\":\"\",\"quantity\":100},{\"partNumber\":\"2222\",\"itemCondition\":\"\",\"quantity\":200}],\"auditInfo\":{\"createdBy\":\"\",\"updatedBy\":\"\",\"createdTS\":\"\",\"updatedTS\":\"\"},\"from\":\"Manu1\",\"to\":\"Lenovo\"}]")

}

func TestSCM_invoke_createShipment(t *testing.T) {
	lcc := new(LenovoChainCode)
	stub := shim.NewMockStub("ldm", lcc)

	// Init cc_version=v0
	checkInit(t, stub, [][]byte{[]byte("init"), []byte("v0")})

	// createSupplierBasicInfo for ITPC organization
<<<<<<< HEAD
	checkInvoke(t, stub, [][]byte{[]byte("createShipment"), []byte("{\"shipmentNumber\": \"1234\",\"trackingnumber\": \"4567\", \"supplierId\": \"supid1\", \"ordernumber\": \"1234\", \"from\": \"lenovo\", \"to\": \"somewhere\" }")})
=======
<<<<<<< HEAD
	checkInvoke(t, stub, [][]byte{[]byte("createShipment"), []byte("{\"shipmentNumber\": \"1236\",\"trackingnumber\": \"4567\", \"supplierId\": \"supid1\", \"ordernumber\": \"0001\", \"from\": \"lenovo\", \"to\": \"somewhere\" }")})
=======
	checkInvoke(t, stub, [][]byte{[]byte("createShipment"), []byte("{\"shipmentNumber\": \"4321\",\"partNumber\": \"4567\", \"supplierId\": \"supid1\", \"partSerialNumber\": \"pserial\", \"poNumber\": \"po1\"}")})
>>>>>>> d02955fd2757829de2f82dec5e91b55e1b638d45
>>>>>>> mattheleon-itpeople-leonDev

	// validate supplier details of org ITPC with querySupplierBasicInfo
	checkQuery(t, stub, "queryShipment", "{\"shipmentNumber\":\"4321\", \"Requestor\":\"supid1\"}", "[{\"orderNumber\":\"1234\",\"supplierId\":\"Manu1\",\"items\":[{\"partNumber\":\"1111\",\"itemCondition\":\"\",\"quantity\":100},{\"partNumber\":\"2222\",\"itemCondition\":\"\",\"quantity\":200}],\"auditInfo\":{\"createdBy\":\"\",\"updatedBy\":\"\",\"createdTS\":\"\",\"updatedTS\":\"\"},\"from\":\"Manu1\",\"to\":\"Lenovo\"}]")


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
