package main

import (
	"encoding/json"
	"fmt"

	"github.com/hyperledger/fabric/core/chaincode/shim"
	pb "github.com/hyperledger/fabric/protos/peer"
)
var (
	fileName = "chaincode"
)

//structure of chaincode
type CloudChaincode struct{
}

//structure of data hash
type Hash struct{
	Key string `json:"Key"`
	Value string `json:"Value"`
}

//initialization function
func (t *CloudChaincode) Init(stub shim.ChaincodeStubInterface)pb.Response{
	// Whatever variable initialisation you want can be done here //
	return shim.Success(nil)
}

// invoking functions
func  (t *CloudChaincode) Invoke(stub shim.ChaincodeStubInterface)pb.Response{
	fmt.Println("Entering Invoke")

	// IF-ELSE-IF all the functions 
	function, args := stub.GetFunctionAndParameters()
	if function == "PutHash" {
		return t.PutHash(stub, args)
	}else if function == "QueryHash" {
		return t.QueryHash(stub, args)
	}
	fmt.Println("invoke did not find func : " + function) //error
	return shim.Error("Received unknown function invocation")
	// end of all functions
}

//put state
func  (t *CloudChaincode) PutHash(stub shim.ChaincodeStubInterface, args []string)pb.Response{
	if len(args) != 2 {
		fmt.Println("Give arguments as: key---value ")
		return shim.Error("Incorrect number of arguments")
	}

	// ASSIGNING VARIABLES
	var Key = args[0]
	var Value = args[1]
	
	// CHECKING IF KEY ALREADY EXISTS
	HashAsBytes, err := stub.GetState(Key)
	if err != nil {
		return shim.Error("Failed to get key:" + err.Error())
	}else if HashAsBytes != nil {
		return shim.Error("Fatel error: network is compromised ")
	}

	// PUT DATA ON LEDGER
	var hash = &Hash{Key: Key, Value: Value}
	HashJSONasBytes, err := json.Marshal(hash)
	if err != nil {
		return shim.Error(err.Error())
	}
	err = stub.PutState(Key, HashJSONasBytes)
	if err != nil {
		return shim.Error(err.Error())
	}
	fmt.Println("Ledger Updated Successfully")
	return shim.Success(nil)
}

// QUERY FUNCTION
func  (t *CloudChaincode) QueryHash(stub shim.ChaincodeStubInterface, args []string)pb.Response{
	if len(args) != 1{
		fmt.Println("Give Valid Key")
		return shim.Error("Incorrect number of arguments")
	}
	var key = args[0]
	ValueAsBytes, err := stub.GetState(key)
	if err != nil {
		fmt.Println("Invalid key")
		return shim.Error(err.Error())
	}
	if ValueAsBytes == nil {
		fmt.Println("Data not found")
	}
	return shim.Success(ValueAsBytes)
}

// MAIN FUNCTION
func  main() {
	err := shim.Start(new(CloudChaincode))
	if err != nil {
		fmt.Printf("Error starting Chaincode: %s", err)
	}
}