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
	"encoding/json"
	"errors"
	"fmt"

	"github.com/hyperledger/fabric/core/chaincode/shim"
)

// SimpleChaincode example simple Chaincode implementation
type SimpleChaincode struct {
}

type UserDetails struct {
	FirstName string
	LastName  string
	Date      string
}
type User struct {
	UserName string
	Details  UserDetails
}

func main() {
	err := shim.Start(new(SimpleChaincode))
	if err != nil {
		fmt.Printf("Error starting Simple chaincode: %s", err)
	}
}

// Init resets all the things
func (t *SimpleChaincode) Init(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {
	if len(args) != 3 {
		return nil, errors.New("Incorrect number of arguments. Expecting 1")
	}
	Users := make([]User, 0)
	details := UserDetails{FirstName: args[1], LastName: args[2], Date: ""}
	user := User{UserName: args[0], Details: details}
	Users = append(Users, user)
	userJson, err := json.Marshal(Users)
	if err != nil {
		return nil, err
	}
	err1 := stub.PutState("User", userJson)
	if err1 != nil {
		return nil, err1
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
		return t.updateDate(stub, args)
	}
	fmt.Println("invoke did not find func: " + function)

	return nil, errors.New("Received unknown function invocation: " + function)
}

// Query is our entry point for queries
func (t *SimpleChaincode) Query(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {
	fmt.Println("query is running " + function)

	// Handle different functions
	if function == "read" { //read a variable
		return t.read(stub)
	}
	fmt.Println("query did not find func: " + function)

	return nil, errors.New("Received unknown function query: " + function)
}

// write - invoke function to write key/value pair
func (t *SimpleChaincode) updateDate(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	fmt.Println("running write()")

	if len(args) != 2 {
		return nil, errors.New("Incorrect number of arguments. Expecting 1. Date of travel")
	}
	userName := args[0]
	date := args[1]

	userArray, _ := t.unpack(stub)

	for i, value := range userArray {
		if value.UserName == userName {
			value.Details.Date = date
		}
		userArray[i] = value
	}

	return nil, nil
}
func (t *SimpleChaincode) repack(stub shim.ChaincodeStubInterface, userArray []User) ([]byte, error) {

	value, err := json.Marshal(userArray)

	if err != nil {
		return nil, err
	}

	err1 := stub.PutState("User", value) //write the variable into the chaincode state
	if err1 != nil {
		return nil, err1
	}
	return nil, nil

}

//Read from the stub and return the user array
func (t *SimpleChaincode) unpack(stub shim.ChaincodeStubInterface) ([]User, error) {
	valuAsBytes, err := stub.GetState("User")
	if err != nil {
		return nil, err
	}
	Users := make([]User, 0)
	jsonError := json.Unmarshal(valuAsBytes, &Users)
	if jsonError != nil {
		return nil, jsonError
	}
	return Users, nil
}

// read - query function to read key/value pair
func (t *SimpleChaincode) read(stub shim.ChaincodeStubInterface) ([]byte, error) {

	valAsbytes, err := stub.GetState("User")
	if err != nil {

		return nil, err
	}

	return valAsbytes, nil
}
