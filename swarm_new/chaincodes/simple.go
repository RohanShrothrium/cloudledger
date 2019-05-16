package main

import (
	"encoding/json"
	"fmt"

	"github.com/hyperledger/fabric/core/chaincode/shim"
	pb "github.com/hyperledger/fabric/protos/peer"
)
import "strconv"
var (
	fileName = "simple"
)

//structure of chaincode
type CloudChaincode struct{
}

// Data regarding where the file is stored
type FileData struct{
	SecretKey string `json:"SecretKey"`
	ServiceProviderMap map[int][]string
}

//structure of data 
type Data struct{
	CompositeKey string `json:"CompositeKey"`
	PublicKey string `json:"PublicKey"`
	FileData map[string]FileData `json:"FileData"`
}

//structure for sharing data
type ShareData struct{
	EdgeKey string `json:"EdgeKey"`
	FileData map[string]FileData `json:"FileData"`
}

//initialization function
func (t *CloudChaincode) Init(stub shim.ChaincodeStubInterface)pb.Response{
	// Whatever variable initialisation you want can be done here //
	return shim.Success(nil)
}

// invoking functions
func  (t *CloudChaincode) Invoke(stub shim.ChaincodeStubInterface)pb.Response{
	// IF-ELSE-IF all the functions 
	function, args := stub.GetFunctionAndParameters()
	if function == "CreateUser" {
		return t.CreateUser(stub, args)
	}else if function == "UploadFile" {
		return t.UploadFile(stub, args)
	}else if function == "DownloadFile" {
		return t.DownloadFile(stub, args)
	}else if function == "DeleteFile" {
		return t.DeleteFile(stub, args)
	}else if function == "ShareFile" {
		return t.ShareFile(stub, args)
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
	ServiceProviderMap := make(map[int][]string)
	var SecretKey = args[1]
	var N = len(args)
	fmt.Println(N)
	var count = 1
	for i := 2; i < N; i++ {
		k, err := strconv.Atoi(args[i])
		if err != nil {
			ServiceProviderMap[count] = append(ServiceProviderMap[count],args[i])
		}else {
			count = k
		}	
	}
	fmt.Println(ServiceProviderMap)
	var FileDataNew = &FileData{SecretKey:SecretKey, ServiceProviderMap:ServiceProviderMap}
	fmt.Println(FileDataNew)
	var Data Data
	err = json.Unmarshal(DataAsBytes, &Data)
	if err != nil {
		return shim.Error("Error encountered during unmarshalling the data")
	}
	if Data.FileData == nil {
		FileDataUploaded := make(map[string]FileData)
		FileDataUploaded[SecretKey] = *FileDataNew
		Data.FileData = FileDataUploaded
	}else{
		Data.FileData[SecretKey] = *FileDataNew
	}
	DataJsonAsBytes, err :=json.Marshal(Data)
	fmt.Println(Data)
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

// Deleting a file 
func  (t *CloudChaincode) DeleteFile(stub shim.ChaincodeStubInterface, args []string)pb.Response{
	if len(args) != 2 {
		fmt.Println("Incorrect number of arguments")
		return shim.Error("Error encountered")
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
	// 
	// CHECK THIS OUT
	delete(Data.FileData, SecretKey)

	DataJsonAsBytes, err :=json.Marshal(Data)
	fmt.Println(Data)
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

// Sharing a file
func  (t *CloudChaincode) ShareFile(stub shim.ChaincodeStubInterface, args []string)pb.Response{
	if len(args) != 3 {
		fmt.Println("Incorrect number of arguments")
		return shim.Error("Error encountered")
	}
	var CompositeKey = args[0]
	var SecretKey = args[1]
	var EdgeKey = args[2]
	DataAsBytes, err := stub.GetState(CompositeKey)
	if err != nil {
		return shim.Error("Error encountered")
	}else if DataAsBytes == nil {
		return shim.Error("No user with the given CompositeKey")
	}
	var Data Data
	var ShareData ShareData
	err = json.Unmarshal(DataAsBytes, &Data)
	ShareDataAsBytes, err := stub.GetState(EdgeKey)
	if err != nil {
		return shim.Error("Error encountered")
	}else if ShareDataAsBytes == nil {
		ShareData.EdgeKey = EdgeKey
		FileDataShare := make(map[string]FileData)
		FileDataShare[SecretKey] = Data.FileData[SecretKey]
		ShareData.FileData = FileDataShare
	}else if ShareDataAsBytes != nil {
		err = json.Unmarshal(ShareDataAsBytes, &ShareData)
		ShareData.FileData[SecretKey] = Data.FileData[SecretKey]
	}
		
	ShareDataJsonAsBytes, err := json.Marshal(ShareData)
	fmt.Println(ShareData)
	if err != nil {
		return shim.Error("Error encountered while remarshalling")
	}
	err = stub.PutState(EdgeKey, ShareDataJsonAsBytes)
	if err != nil {
		return shim.Error("error encountered while putting state")
	}
	fmt.Println("File shared successfully")
	return shim.Success(nil)
}

// MAIN FUNCTION
func  main() {
	err := shim.Start(new(CloudChaincode))
	if err != nil {
		fmt.Printf("Error starting Chaincode: %s", err)
	}
}
