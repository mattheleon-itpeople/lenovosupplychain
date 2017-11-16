package main

import (
	"dbapi"
	"encoding/json"
	"fmt"
	"runtime"
	"strings"

	"github.com/hyperledger/fabric/core/chaincode/shim"
)

func getFunctionName() string {
	pc, _, _, _ := runtime.Caller(1)
	splits := strings.Split(runtime.FuncForPC(pc).Name(), "/")
	name := strings.Split(splits[len(splits)-1], ".")
	return name[len(name)-1]
}

func retrieveAndMarshalPOObject(stub shim.ChaincodeStubInterface, args []string) (PurchaseOrder, error) {
	order := PurchaseOrder{}
	var err error
	var objectBytes []byte

	keys := args

	objectBytes, err = dbapi.QueryObject(stub, "PO", keys)
	if err != nil {
		err = fmt.Errorf("retrieveAndMarshal() : Failed to query PurchaseOrder object")
		return order, err
	}

	if objectBytes == nil {
		return order, fmt.Errorf("retrieveAndMarshal() : "+"PurchaseOrder Number: %s does not exist ", order.PONumber)
	}
	err = json.Unmarshal(objectBytes, &order)
	if err != nil {
		return order, fmt.Errorf("retrieveAndMarshal()  : marshalling PO failed")
	}

	err = json.Unmarshal(objectBytes, &order)
	if err != nil {
		return order, fmt.Errorf("retrieveAndMarshal()  : marshalling PO failed")
	}
	return order, nil

}

func retrieveAndMarshalSOObject(stub shim.ChaincodeStubInterface, args []string) (SalesOrder, error) {
	order := SalesOrder{}
	var err error
	var objectBytes []byte

	keys := args

	objectBytes, err = dbapi.QueryObject(stub, "SO", keys)
	if err != nil {
		err = fmt.Errorf("retrieveAndMarshal() : Failed to query SalesOrder object")
		return order, err
	}

	if objectBytes == nil {
		return order, fmt.Errorf("retrieveAndMarshal() : "+"SalesOrder Number: %s does not exist ", order.PONumber)
	}
	err = json.Unmarshal(objectBytes, &order)
	if err != nil {
		return order, fmt.Errorf("retrieveAndMarshal()  : marshalling SO failed")
	}

	err = json.Unmarshal(objectBytes, &order)
	if err != nil {
		return order, fmt.Errorf("retrieveAndMarshal()  : marshalling PO failed")
	}
	return order, nil

}
