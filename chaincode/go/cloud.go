package main

import (
	"encoding/json"
	"fmt"

	"github.com/hyperledger/fabric/core/chaincode/shim"
	pb "github.com/hyperledger/fabric/protos/peer"
)
var (
	fileName = "cloud"
)

//structure of chaincode
type CloudChaincode struct{
}

// Data regarding where the file is stored
type FileData struct{
	SecretKey string `json:"SecretKey"`
	ServiceProviderMap map[int]string
}

//structure of data 
type Data struct{
	CompositeKey string `json:"CompositeKey"`
	PublicKey string `json:"PublicKey"`
	FileData map[string]FileData `json:"FileData"`
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
	if function == "CreateUser" {
		return t.CreateUser(stub, args)
	}else if function == "UploadFile" {
		return t.UploadFile(stub, args)
	}else if function == "DownloadFile" {
		return t.DownloadFile(stub, args)
	}
	fmt.Println("invoke did not find func : " + function) //error
	return shim.Error("Received unknown function invocation")
	// end of all functions
}

// Creating data of user 
func  (t *CloudChaincode) CreateUser(stub shim.ChaincodeStubInterface, args []string)pb.Response{
	var CompositeKey = args[0]
	var PublicKey = args[1]
	// checking for an error or if the user already exists
	DataAsBytes, err := stub.GetState(CompositeKey)
	if err != nil {
		return shim.Error("Failed to get CompositeKey:" + err.Error())
	}else if DataAsBytes != nil{
		return shim.Error("User with current composite Key already exists")
	}
	var Data = &Data{CompositeKey:CompositeKey, PublicKey:PublicKey}
	DataJsonAsBytes, err :=json.Marshal(Data)
	if err != nil {
		shim.Error("Error encountered while Marshalling")
	}
	err = stub.PutState(CompositeKey, DataJsonAsBytes)
	if err != nil {
		shim.Error("Error encountered while putting data")
	}
	fmt.Println("Ledger Updated Successfully")
	return shim.Success(nil)
}

// Uploading a file
func  (t *CloudChaincode) UploadFile(stub shim.ChaincodeStubInterface, args []string)pb.Response{
	var CompositeKey = args[0]
	DataAsBytes, err := stub.GetState(CompositeKey)
	if err != nil {
		return shim.Error("Failed to get CompositeKey:" + err.Error())
	}else if DataAsBytes == nil{
		return shim.Error("Unkown composite key")
	}
	var SecretKey = args[1]
	var ServiceProviderMap = args[2]
	var FileDataNew = &FileData{SecretKey:SecretKey, ServiceProviderMap:ServiceProviderMap}
	var Data Data
	err = json.Unmarshal(DataAsBytes, &Data)
	if err != nil {
		return shim.Error("Error encountered during unmarshalling the data")
	}
	Data.FileData[SecretKey] = *FileDataNew
	DataJsonAsBytes, err :=json.Marshal(Data)
	if err != nil {
		return shim.Error("Error encountered while remarshalling")
	}
	err = stub.PutState(CompositeKey, DataJsonAsBytes)
	if err != nil {
		return shim.Error("error encountered while putting state")
	}
	fmt.Println("File uploaded successfully")
	return shim.Success(nil)
}

// Downloading a file
func  (t *CloudChaincode) DownloadFile(stub shim.ChaincodeStubInterface, args []string)pb.Response{
	if len(args) != 2 {
		fmt.Println("Incorrect number of arguments")
		return shim.Error("Incorrect number of arguments")
	}
	var CompositeKey = args[0]
	var SecretKey = args[1]
	DataAsBytes, err := stub.GetState(CompositeKey)
	if err != nil {
		return shim.Error("Error encountered")
	}else if DataAsBytes == nil {
		return shim.Error("No user with the given CompositeKey")
	}
	var Data Data
	err = json.Unmarshal(DataAsBytes, &Data)
	var FileData FileData
	FileData = Data.FileData[SecretKey]
	FileJsonAsBytes, err :=json.Marshal(FileData)
	if err != nil {
		return shim.Error("Error encountered")
	}
	return shim.Success(FileJsonAsBytes)
}

// MAIN FUNCTION
func  main() {
	err := shim.Start(new(CloudChaincode))
	if err != nil {
		fmt.Printf("Error starting Chaincode: %s", err)
	}
}