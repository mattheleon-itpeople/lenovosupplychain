package main

import (
	"dbapi"
	"encoding/json"
	"fmt"

	"github.com/hyperledger/fabric/core/chaincode/shim"
)

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
