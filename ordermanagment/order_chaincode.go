/*
 * Copyright IBM Corp All Rights Reserved
 *
 * SPDX-License-Identifier: Apache-2.0
 */

package main

import (
	"encoding/json"
	"fmt"

	"github.com/hyperledger/fabric/core/chaincode/shim"
	pd "github.com/hyperledger/fabric/protos/peer"
)

// SimpleAsset implements a simple chaincode to manage an asset
type OrderAsset struct {
}

//Order sturct holkding both the Purchase Order and Sales Orders - ditinquished by the object type
type Order struct {
	ObjectType      string      `json:"docType"`         //docType is used to distinguish the various types of objects in state database
	OrderNumber     string      `json:"ordernumber"`     //the fieldtags are needed to keep case from bouncing around
	ReferenceNumber string      `json:"referencenumber"` //the fieldtags are needed to keep case from bouncing around
	ReferenceType   string      `json:"referencetype"`
	From            string      `json:"from"`
	To              string      `json:"to"`
	Part            []OrderLine `json:"part"`
	PODate          string      `json:"poDate"`
}

/*Organization of the user*/
type Organization struct {
	ObjectType string `json:"docType"` //docType is used to distinguish the various types of objects in state database
	OrgID      string `json:"orgid"`
	Name       string `json:"name"`
	Details    string `json:"details"`
}

/*User - associated with an Organization*/
type User struct {
	ObjectType string `json:"docType"` //docType is used to distinguish the various types of objects in state database
	UserName   string `json:"username"`
	FirstName  string `json:"firstname"`
	LastName   string `json:"lastname"`
	OrgID      string `json:"orgid"`
}

/*Actual part detail */
type PartDetail struct {
	ObjectType    string `json:"docType"` //docType is used to distinguish the various types of objects in state database
	PartNumber    string `json:"partnumber"`
	Description   string `json:"description"`
	UnitOfMeasure string `json:"unitofmeasure"`
}

/*Order line  - used in the order*/
type OrderLine struct {
	ObjectType   string     `json:"docType"` //docType is used to distinguish the various types of objects in state database
	Part         PartDetail `json:"part"`
	Price        float64    `json:"price"`
	Quantity     int        `json:"quantity"`
	DeliveryDate string     `json:"deliverydate"`
}

/*Status line - updates the status of each line in the purchase order*/
type StatusLine struct {
	ObjectType           string     `json:"docType"` //docType is used to distinguish the various types of objects in state database
	Part                 PartDetail `json:"part"`
	ExpectedDeliveryDate string     `json:"expecteddeliverydate"`
}

/*Actual status transaction object */
type POStatus struct {
	ObjectType      string       `json:"docType"` //docType is used to distinguish the various types of objects in state database
	PONumber        string       `json:"ponumber"`
	ReferenceNumber string       `json:"referencenumber"`
	Status          []StatusLine `json:"status"`
}

// Init is called during chaincode instantiation to initialize any
// data. Note that chaincode upgrade also calls this function to reset
// or to migrate data.
func (t *OrderAsset) Init(stub shim.ChaincodeStubInterface) pd.Response {
	// Get the args from the transaction proposal
	return shim.Success(nil)
}

// Invoke is called per transaction on the chaincode. Each transaction is
// either a 'get' or a 'set' on the asset created by Init function. The Set
// method may create a new asset by specifying a new key-value pair.
func (t *OrderAsset) Invoke(stub shim.ChaincodeStubInterface) pd.Response {
	// Extract the function and args from the transaction proposal
	fn, args := stub.GetFunctionAndParameters()

	var result pd.Response
	var err error
	if fn == "createOrder" {
		result = t.createOrder(stub, args)
	} else if fn == "getPurchaseOrder" {
		result = t.getPurchaseOrder(stub, args)
	} else {
		return shim.Error("No such function: " + fn)
	}
	if err != nil {
		return result
	}

	// Return the result as success payload
	return result
}

//Creates a computer
func (t *OrderAsset) createOrder(stub shim.ChaincodeStubInterface, args []string) pd.Response {
	var err error
	var order Order = Order{}
	if len(args) < 1 {
		jsonError := "\"Error\": \"Expecting order to be provided in parameter list\""
		return shim.Error(jsonError)
	}
	err = json.Unmarshal([]byte(args[0]), &order)
	if err != nil {
		fmt.Printf(err.Error())
		jsonError := "\"Error\": \"" + err.Error() + "\""
		return shim.Error(jsonError)
	}

	err = stub.PutState(order.OrderNumber, []byte(args[0]))

	fmt.Printf(">>> I have Created order " + order.OrderNumber)
	return shim.Success(nil)

}

//Creates a computer
func (t *OrderAsset) getPurchaseOrder(stub shim.ChaincodeStubInterface, args []string) pd.Response {
	var order Order = Order{}
	if len(args) < 1 {
		jsonError := "\"Error\": \"Expecting order number to be provided in parameter list\""
		return shim.Error(jsonError)
	}

	ordernumber := args[0]
	orderAsBytes, err := stub.GetState(ordernumber)
	if err != nil {
		jsonError := "\"Error\": \"" + err.Error() + "\""
		return shim.Error(jsonError)
	}
	if orderAsBytes == nil {
		jsonError := "{\"Error\":\"Failed to get state for order number" + ordernumber + "\"}"
		return shim.Error(jsonError)
	}
	err = json.Unmarshal(orderAsBytes, &order)
	if err != nil {
		jsonError := "\"Error\": \"" + err.Error() + "\""
		return shim.Error(jsonError)
	}

	fmt.Printf(">>> I  have retrieved order " + ordernumber)

	return shim.Success(orderAsBytes)

}

// ===================================================================================
// Main
// ===================================================================================
func main() {
	err := shim.Start(new(OrderAsset))
	if err != nil {
		fmt.Printf("Error starting Computer chaincode: %s", err)
	}

	fmt.Printf("**********Order Management Chaincode started**********")
}
