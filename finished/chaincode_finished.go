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

var logger = shim.NewLogger("CLDChaincode")

type Consignment struct {
	PackageID    string `json: "packageID"`
	PackageType  string `json: "packageType"`
	BookedOn     string `json: "bookedOn"`
	From         string `json: "from"`
	To           string `json: "to"`
	FlightNumber string `json: "flightNumber"`
	Date         string `json: "date"`
}
type packageID_holder struct {
	packageIDs []string `json: "packageIDs"`
}

func main() {
	err := shim.Start(new(SimpleChaincode))
	if err != nil {
		fmt.Printf("Error starting Simple chaincode: %s", err)
	}
}

// Init resets all the things
func (t *SimpleChaincode) Init(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {
	// if len(args) != 3 {
	// 	return nil, errors.New("Incorrect number of arguments. Expecting 1")
	// }
	var packageIDs packageID_holder
	bytes, err := json.Marshal(&packageIDs)
	if err != nil {
		return nil, err
	}

	err = stub.PutState("packageIDs", bytes)

	return nil, nil
}

// Invoke isur entry point to invoke a chaincode function
func (t *SimpleChaincode) Invoke(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {
	fmt.Println("invoke is running " + function)

	// Handle different functions
	if function == "init" {
		return t.Init(stub, "init", args)
	} else if function == "createpackage" {
		return t.CreatePackage(stub, args)
	}
	fmt.Println("invoke did not find func: " + function)

	return nil, errors.New("Received unknown function invocation: " + function)
}
func (t *SimpleChaincode) CreatePackage(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	var c Consignment
	packageId := "\"packageID\":\"" + args[0] + "\", "
	packageType := "\"packageType\":\"" + args[1] + "\", "
	bookedOn := "\"bookedOn\":\"" + args[2] + "\", "
	from := "\"from\":\"" + args[3] + "\", "
	to := "\"to\":\"" + args[4] + "\", "
	flightNumber := "\"flightNumber\":\"\", "
	date := "\"date\":\"\" "
	packageJson := "{" + packageId + packageType + bookedOn + from + to + flightNumber + date + "}"

	err := json.Unmarshal([]byte(packageJson), &c)

	if err != nil {
		return nil, err
	}

	fmt.Println("Inside create package")
	_, err = t.save_changes(stub, c)
	if err != nil {
		return nil, err
	}

	bytes, error := stub.GetState("packageIDs")
	if error != nil {
		return nil, error
	}

	var packageIDs packageID_holder

	err = json.Unmarshal(bytes, &packageIDs)
	if err != nil {
		return nil, err
	}

	packageIDs.packageIDs = append(packageIDs.packageIDs, args[0])

	bytes, err = json.Marshal(packageIDs)
	if err != nil {
		return nil, err
	}

	err = stub.PutState("packageIDs", bytes)

	return nil, nil

}
func (t *SimpleChaincode) save_changes(stub shim.ChaincodeStubInterface, c Consignment) (bool, error) {

	bytes, err := json.Marshal(c)

	if err != nil {
		fmt.Printf("SAVE_CHANGES: Error converting vehicle record: %s", err)
		return false, errors.New("Error converting vehicle record")
	}

	err = stub.PutState(c.PackageID, bytes)

	if err != nil {
		fmt.Printf("SAVE_CHANGES: Error storing vehicle record: %s", err)
		return false, errors.New("Error storing vehicle record")
	}

	return true, nil
}

// Query is our entry point for queries
func (t *SimpleChaincode) Query(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {
	fmt.Println("query is running " + function)
	logger.Debug("query function: ", function)
	// Handle different functions
	if function == "read" { //read a variable
		return t.get_packages(stub)
	}
	fmt.Println("query did not find func: " + function)

	return nil, errors.New("Received unknown function query: " + function)
}

//user_type1_0 18e23bffec
//user_type1_3 ba7469144a
//user_type1_4 427c6df636

//=================================================================================================================================
//	 get_vehicles
//=================================================================================================================================

func (t *SimpleChaincode) get_packages(stub shim.ChaincodeStubInterface) ([]byte, error) {
	logger.Debug("get_packages1: ")
	bytes, err := stub.GetState("packageIDs")
	logger.Debug("get_packages1: ", bytes)
	if err != nil {
		return nil, errors.New("Unable to get v5cIDs")
	}
	fmt.Println("Inside read package1")
	var packageIDs packageID_holder

	err = json.Unmarshal(bytes, &packageIDs)

	if err != nil {
		return nil, errors.New("Corrupt V5C_Holder")
	}

	result := "["

	var temp []byte
	var c Consignment

	for _, packageID := range packageIDs.packageIDs {

		c, err = t.retrieve_id(stub, packageID)
		logger.Debug("inside loop function: ", packageID)

		if err != nil {
			return nil, errors.New("Failed to retrieve V5C")
		}
		temp, err = t.get_package_details(stub, c)

		if err == nil {
			result += string(temp) + ","
		}
	}

	if len(result) == 1 {
		result = string(bytes)
	} else {
		logger.Debug("inside else function: ", result)
		result = result[:len(result)-1] + "]"
	}
	//return bytes, nil
	return []byte(result), nil
}
func (t *SimpleChaincode) get_package_details(stub shim.ChaincodeStubInterface, c Consignment) ([]byte, error) {
	logger.Debug("get_package_details function: ", c)
	bytes, err := json.Marshal(c)

	if err != nil {
		return nil, errors.New("GET_VEHICLE_DETAILS: Invalid vehicle object")
	}
	return bytes, nil

}

func (t *SimpleChaincode) retrieve_id(stub shim.ChaincodeStubInterface, packageID string) (Consignment, error) {
	logger.Debug("retrieve_id function: ", packageID)
	var c Consignment

	bytes, err := stub.GetState(packageID)

	if err != nil {
		fmt.Printf("RETRIEVE_V5C: Failed to invoke vehicle_code: %s", err)
		return c, errors.New("RETRIEVE_V5C: Error retrieving vehicle with v5cID = " + packageID)
	}

	err = json.Unmarshal(bytes, &c)

	if err != nil {
		fmt.Printf("RETRIEVE_V5C: Corrupt vehicle record "+string(bytes)+": %s", err)
		return c, errors.New("RETRIEVE_V5C: Corrupt vehicle record" + string(bytes))
	}

	return c, nil
}

// write - invoke function to write key/value pair
func (t *SimpleChaincode) updateDate(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	// fmt.Println("running write()")
	//
	// if len(args) != 2 {
	// 	return nil, errors.New("Incorrect number of arguments. Expecting 1. Date of travel")
	// }
	// userName := args[0]
	// date := args[1]
	//
	// userArray, _ := t.unpack(stub, userName)
	//
	// for i, value := range userArray {
	// 	if value.UserName == userName {
	// 		value.Details.Date = date
	// 	}
	// 	userArray[i] = value
	// }
	// val, err := t.repack(stub, userArray)
	// if err != nil {
	// 	return nil, err
	// }

	return nil, nil
}

//func (t *SimpleChaincode) repack(stub shim.ChaincodeStubInterface, userArray []User) ([]byte, error) {

// 	value, err := json.Marshal(userArray)
//
// 	if err != nil {
// 		return nil, err
// 	}
//
// 	err1 := stub.PutState("User", value) //write the variable into the chaincode state
// 	if err1 != nil {
// 		return nil, err1
// 	}
// 	return nil, nil
//
// }
//
// //Read from the stub and return the user array
//  func (t *SimpleChaincode) unpack(stub shim.ChaincodeStubInterface, userName string) (string, error) {
// // 	valuAsBytes, err := stub.GetState(userName)
// // 	if err != nil {
// // 		return nil, err
// // 	}
// // 	Users := make([]User, 0)
// // 	jsonError := json.Unmarshal(valuAsBytes, &Users)
// // 	if jsonError != nil {
// // 		return nil, jsonError
// // 	}
// 	return nil, nil
// }
//
// // read - query function to read key/value pair
func (t *SimpleChaincode) read(stub shim.ChaincodeStubInterface, userName string) ([]byte, error) {
	//
	// 	valAsbytes, err := stub.GetState(userName)
	// 	if err != nil {
	//
	// 		return nil, err
	// 	}

	return nil, nil
}
