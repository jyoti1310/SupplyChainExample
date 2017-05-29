/*
Copyright IBM Corp 2016 All Rights Reserved.

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

package main

import (
	"errors"
	"fmt"
	"strconv"
	"encoding/json"
	
	"github.com/hyperledger/fabric/core/chaincode/shim"
)

// SimpleChaincode example simple Chaincode implementation
type SimpleChaincode struct {
}

type Shipment struct{
		CurrentOwner string `json:"CurrentOwner"`
        MaximumTemperatureRecorded int `json:"MaximumTemperatureRecorded"`
        TemperatureThreshold int `json:"TemperatureThreshold"`
}

func main() {
	err := shim.Start(new(SimpleChaincode))
	if err != nil {
		fmt.Printf("Error starting Simple chaincode: %s", err)
	}
}

// Init resets all the things
func (t *SimpleChaincode) Init(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {
	if len(args) != 1 {
		return nil, errors.New("Incorrect number of arguments. Expecting 1")
	}

	err := stub.PutState("hello_Block", []byte(args[0]))
	if err != nil {
		return nil, err
	}

	return nil, nil
}

// Invoke isur entry point to invoke a chaincode function
func (t *SimpleChaincode) Invoke(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {
	fmt.Println("invoke is running " + function)

	// Handle different functions
	if function == "init" {
		return t.Init(stub, "init", args)
	} else if function == "write" {
		return t.write(stub, args)
	} else if  function == "startShipment" {
		return t.addNewShipment(stub, args)
	}else if  function == "transferOwner" {
		return t.transferOwner(stub, args)
	}
	fmt.Println("invoke did not find func: " + function)

	return nil, errors.New("Received unknown function invocation: " + function)
}

// Query is our entry point for queries
func (t *SimpleChaincode) Query(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {
	fmt.Println("query is running " + function)

	// Handle different functions
	if function == "read" { //read a variable
		return t.read(stub, args)
	} /*else if function == "searchLogBog" {
		return t.searchSKATEmployee(stub,args)
	}*/
	fmt.Println("query did not find func: " + function)

	return nil, errors.New("Received unknown function query: " + function)
}

// write - invoke function to write key/value pair
func (t *SimpleChaincode) write(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	var key, value string
	var err error
	fmt.Println("running write()")

	if len(args) != 2 {
		return nil, errors.New("Incorrect number of arguments. Expecting 2. name of the key and value to set")
	}

	key = args[0] //rename for funsies
	value = args[1]
	err = stub.PutState(key, []byte(value)) //write the variable into the chaincode state
	if err != nil {
		return nil, err
	}
	return nil, nil
}

// read - query function to read key/value pair
func (t *SimpleChaincode) read(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	var key, jsonResp string
	var err error

	if len(args) != 1 {
		return nil, errors.New("Incorrect number of arguments. Expecting name of the key to query")
	}

	key = args[0]
	valAsbytes, err := stub.GetState(key)
	if err != nil {
		jsonResp = "{\"Error\":\"Failed to get state for " + key + "\"}"
		return nil, errors.New(jsonResp)
	}

	return valAsbytes, nil
}

// ============================================================================================================================
// Init Employee - create a new Employee, store into chaincode state
// ============================================================================================================================
func (t *SimpleChaincode) addNewShipment(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	var err error
	var jsonResp string
	//,jsonResp string
	
	//   0       1       2     3
	// "asdf", "blue", "35", "bob"
	if len(args) != 3 {
		return nil, errors.New("Incorrect number of arguments. Expecting 9")
	}
	fmt.Println("- adding new Shipment")
	fmt.Println("CurrentOwner-"+args[0])
	
	
	/*if len(args[0]) <= 0 {
		return nil, errors.New("1st argument must be a non-empty string")
	}
	if len(args[1]) <= 0 {
		return nil, errors.New("2nd argument must be a non-empty string")
	}
	if len(args[2]) <= 0 {
		return nil, errors.New("3rd argument must be a non-empty string")
	}
	if len(args[3]) <= 0 {
		return nil, errors.New("4th argument must be a non-empty string")
	}*/
	NewShipment := Shipment{}

	NewShipment.CurrentOwner = args[0]
	
	NewShipment.MaximumTemperatureRecorded, err =strconv.Atoi(args[1])
	
	if err != nil {
		return nil, errors.New("Maximum temperature Recorded must be a numeric string")
	}
	NewShipment.TemperatureThreshold, err = strconv.Atoi(args[2])
	
	if err != nil {
		return nil, errors.New("TemperatureThreshold must be a numeric string")
	}

	fmt.Println("adding Shipment @ " + NewShipment.CurrentOwner + ", " + strconv.Itoa(NewShipment.MaximumTemperatureRecorded)+ ", " + strconv.Itoa(NewShipment.TemperatureThreshold));
	fmt.Println("- end adding new Shipment")
	jsonAsBytes, _ := json.Marshal(NewShipment)

	if err != nil {
		return jsonAsBytes, err
	}
	//Added for test purpose
	err = stub.PutState("BlueShipment", jsonAsBytes)	//store shipment
	
	if err != nil {	
		jsonResp = "{\"Error\":\"Failed to Add new Shipment" + "\"}"
		return jsonAsBytes, errors.New(jsonResp)
	}	
	
	fmt.Println("- end addShipment")
	return jsonAsBytes, nil
}

func (t *SimpleChaincode) transferOwner(stub shim.ChaincodeStubInterface, args []string) ([]byte, error){

var jsonResp string 
valAsbytes, err := stub.GetState("BlueShipment")
	if err != nil {
		jsonResp = "{\"Error\":\"Failed to get state for " + "BlueShipment" + "\"}"
		return nil, errors.New(jsonResp)
	}
	var currentShipment Shipment
	json.Unmarshal(valAsbytes, &currentShipment)	

	
	if(currentShipment.MaximumTemperatureRecorded < currentShipment.TemperatureThreshold){
		currentShipment.CurrentOwner =args[0]
		updatedShipmentJsonAsBytes, _  := json.Marshal(currentShipment)
		err = stub.PutState("BlueShipment", updatedShipmentJsonAsBytes)	
		fmt.Println("owner changed")
		return updatedShipmentJsonAsBytes, nil
		}else{
		fmt.Println("Contract Breached")
		}
	
	return nil, nil
}