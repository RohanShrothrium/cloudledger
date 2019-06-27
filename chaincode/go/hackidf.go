package main

import (
	"encoding/json"
	"fmt"
	"github.com/hyperledger/fabric/core/chaincode/shim"
	pb "github.com/hyperledger/fabric/protos/peer"
	"crypto/sha256"
)
var (
	fileName = "hackidf"
)

//structure of chaincode
type HackidfChaincode struct{
}

// User Structure
type User struct{
	Username string `json:"Username"`
	PasswordHash [32]byte `json:"PasswordHash"`
	Email string `json:"Email"`
	Ph string `json:"Ph"`
	IsVerified string `json:"IsVerified"`
}

// Organisation Structure
type Organisation struct{
	OrgName string `json:"OrgName"`
	PasswordHash [32]byte `json:"PasswordHash"`
	IsVerified string `json:"IsVerified"`
}

// Claim 
type Claim struct{
	UserID string `json:"UserID"`
	OrgID string `json:"OrgID"`
	Skill string `json:"Skill"`
	Comments string `json:"Comments"`
	Timestamp string `json:"Timestamp"`
	IsVerified string `json:"IsVerified"`
}

//initialization function
func (t *HackidfChaincode) Init(stub shim.ChaincodeStubInterface)pb.Response{
	// Whatever variable initialisation you want can be done here //
	return shim.Success(nil)
}

// invoking functions
func  (t *HackidfChaincode) Invoke(stub shim.ChaincodeStubInterface)pb.Response{
	// IF-ELSE-IF all the functions 
	function, args := stub.GetFunctionAndParameters()
	if function == "CreateUser" {
		return t.CreateUser(stub, args)
	}else if function == "VerifyUser" {
		return t.VerifyUser(stub, args)
	}else if function == "CreateOrg" {
		return t.CreateOrg(stub, args)
	}else if function == "VerifyOrg" {
		return t.VerifyOrg(stub, args)
	}else if function == "MakeClaim" {
		return t.MakeClaim(stub, args)
	}else if function == "VerifyClaim"{
		return t.VerifyClaim(stub, args)
	}else if function == "Query"{
		return t.Query(stub, args)
	}
	fmt.Println("invoke did not find func : " + function) //error
	return shim.Error("Received unknown function invocation")
	// end of all functions
}

// Check KYC
func CheckUser(stub shim.ChaincodeStubInterface, UserID string) int {
	UserAsBytes, err := stub.GetState(UserID)
	var User User
	err = json.Unmarshal(UserAsBytes, &User)
	if err != nil {
		return 0
	}
	if User.IsVerified == "True" {
        return 1
	}
	return 0
}

// Check If Org exists and is verified
func CheckOrg(stub shim.ChaincodeStubInterface, OrgID string) int {
	OrgAsBytes, err := stub.GetState(OrgID)
	if err != nil {
		return 0
	}else if OrgAsBytes == nil{
		return 0
	}
	var Organisation Organisation
	err = json.Unmarshal(OrgAsBytes, &Organisation)
	if err != nil {
		return 0
	}
	if Organisation.IsVerified == "True" {
        return 1
    }
	return 0
}

// Org Verify Password
func OrgVerifyPassword(stub shim.ChaincodeStubInterface, OrgID string, Password string) int {
	OrgAsBytes, err := stub.GetState(OrgID)
	if err != nil {
		return 0
	}else if OrgAsBytes == nil{
		return 0
	}
	var Organisation Organisation
	err = json.Unmarshal(OrgAsBytes, &Organisation)
	if err != nil {
		return 0
	}
	if Organisation.PasswordHash == sha256.Sum256([]byte(Password)) {
        return 1
    }
	return 0
}

// User Verify Password
func UserVerifyPassword(stub shim.ChaincodeStubInterface, UserID string, Password string) int {
	UserAsBytes, err := stub.GetState(UserID)
	if err != nil {
		return 0
	}else if UserAsBytes == nil{
		return 0
	}
	var User User
	err = json.Unmarshal(UserAsBytes, &User)
	if err != nil {
		return 0
	}
	if User.PasswordHash == sha256.Sum256([]byte(Password)) {
        return 1
    }
	return 0
}

// Adding info about a user
func  (t *HackidfChaincode) CreateUser(stub shim.ChaincodeStubInterface, args []string)pb.Response{
	var UserID = args[0]
	var Username = args[1]
	var Password = args[2]
	var Email = args[3]
	var Ph = args[4]
	var IsVerified = "False"
	PasswordHash := sha256.Sum256([]byte(Password))
	// checking for an error or if the user already exists
	UserAsBytes, err := stub.GetState(Username)
	if err != nil {
		return shim.Error("Failed to get Username:" + err.Error())
	}else if UserAsBytes != nil{
		return shim.Error("User with current username already exists")
	}

	var User = &User{Username:Username, PasswordHash:PasswordHash, Email:Email, Ph:Ph, IsVerified:IsVerified}
	UserJsonAsBytes, err :=json.Marshal(User)
	if err != nil {
		shim.Error("Error encountered while Marshalling")
	}
	err = stub.PutState(UserID, UserJsonAsBytes)
	if err != nil {
		shim.Error("Error encountered while Creating User")
	}
	fmt.Println("Ledger Updated Successfully")
	return shim.Success(nil)
}

// Do KYC for peer
func  (t *HackidfChaincode) VerifyUser(stub shim.ChaincodeStubInterface, args []string)pb.Response{
	var UserID = args[0]
	var Password = args[1]
	PasswordHash := sha256.Sum256([]byte(Password))
	if PasswordHash != sha256.Sum256([]byte("Password")) {
		return shim.Error("WRONG PASSWORD ALERT!")
	}
	UserAsBytes, err := stub.GetState(UserID)
	if err != nil {
		return shim.Error("Failed to get User:" + err.Error())
	}else if UserAsBytes == nil{
		return shim.Error("Please give a valid User-ID")
	}
	if CheckUser(stub, UserID) == 1 {
		return shim.Error("The user is already verified")
	}
	var User User
	err = json.Unmarshal(UserAsBytes, &User)
	if err != nil {
		return shim.Error("Error encountered during unmarshalling the data")
	}
	User.IsVerified = "True"
	UserJsonAsBytes, err :=json.Marshal(User)
	if err != nil {
		return shim.Error("Error encountered while remarshalling")
	}
	err = stub.PutState(UserID, UserJsonAsBytes)
	if err != nil {
		return shim.Error("error encountered while putting state")
	}
	fmt.Println("VERIFIED!!")
	return shim.Success(nil)
}

// Adding info about an Organisations
func  (t *HackidfChaincode) CreateOrg(stub shim.ChaincodeStubInterface, args []string)pb.Response{
	var OrgID = args[0]
	var OrgName = args[1]
	var Password = args[2]
	var IsVerified = "False"
	PasswordHash := sha256.Sum256([]byte(Password))
	// checking for an error or if the user already exists
	OrgAsBytes, err := stub.GetState(OrgID)
	if err != nil {
		return shim.Error("Failed to get Organisation:" + err.Error())
	}else if OrgAsBytes != nil{
		return shim.Error("Organisation is already registered")
	}
	var Organisation = &Organisation{OrgName:OrgName, PasswordHash:PasswordHash, IsVerified:IsVerified}
	OrgJsonAsBytes, err :=json.Marshal(Organisation)
	if err != nil {
		shim.Error("Error encountered while Marshalling")
	}
	err = stub.PutState(OrgID, OrgJsonAsBytes)
	if err != nil {
		shim.Error("Error encountered while Creating Organisation")
	}
	fmt.Println("Ledger Updated Successfully")
	return shim.Success(nil)
}

// Verify Organisation
func  (t *HackidfChaincode) VerifyOrg(stub shim.ChaincodeStubInterface, args []string)pb.Response{
	var OrgID = args[0]
	var Password = args[1]
	PasswordHash := sha256.Sum256([]byte(Password))
	if PasswordHash != sha256.Sum256([]byte("Password")) {
		return shim.Error("WRONG PASSWORD ALERT!")
	}
	OrgAsBytes, err := stub.GetState(OrgID)
	if err != nil {
		return shim.Error("Failed to get Organisation:" + err.Error())
	}else if OrgAsBytes == nil{
		return shim.Error("Organisation not registered")
	}
	var Organisation Organisation
	err = json.Unmarshal(OrgAsBytes, &Organisation)
	if err != nil {
		return shim.Error("Error encountered during unmarshalling the data")
	}
	Organisation.IsVerified = "True"
	OrgJsonAsBytes, err :=json.Marshal(Organisation)
	if err != nil {
		return shim.Error("Error encountered while remarshalling")
	}
	err = stub.PutState(OrgID, OrgJsonAsBytes)
	if err != nil {
		return shim.Error("error encountered while putting state")
	}
	fmt.Println("VERIFIED!!")
	return shim.Success(nil)
}

// Make Claim
func  (t *HackidfChaincode) MakeClaim(stub shim.ChaincodeStubInterface, args []string)pb.Response{
	var Hash = args[0]
	var UserID = args[1]
	var Password = args[2]
	var OrgID = args[3]
	var Skill = args[4]
	var Timestamp = args[5]
	var IsVerified = "False"
	if CheckUser(stub, UserID) == 0 {
		return shim.Error("Please finish your KYC procedure with InfoEaze.")
	}
	if CheckOrg(stub, OrgID) == 0 {
		return shim.Error("Org isn't verified by InfoEaze.")
	}
	if UserVerifyPassword(stub, UserID, Password) == 0 {
		return shim.Error("Password doesn't match user")
	}
	var Claim = &Claim{UserID:UserID, OrgID:OrgID, Skill:Skill, Comments:"NIL", Timestamp:Timestamp ,IsVerified:IsVerified}
	ClaimJsonAsBytes, err :=json.Marshal(Claim)
	if err != nil {
		shim.Error("Error encountered while Marshalling")
	}
	err = stub.PutState(Hash, ClaimJsonAsBytes)
	if err != nil {
		shim.Error("Error encountered while Making Claim")
	}
	fmt.Println("Ledger Updated Successfully")
	return shim.Success(nil)
}

// Verify Claim
func  (t *HackidfChaincode) VerifyClaim(stub shim.ChaincodeStubInterface, args []string)pb.Response{
	var Hash = args[0]
	var OrgID = args[1]
	var Password = args[2]
	if OrgVerifyPassword(stub, OrgID, Password) == 0 {
		return shim.Error("Password doesn't match Organisation")
	}
	ClaimAsBytes, err := stub.GetState(Hash)
	if err != nil {
		return shim.Error("Failed to get Claim:" + err.Error())
	}else if ClaimAsBytes == nil{
		return shim.Error("Claim not made")
	}
	var Claim Claim
	err = json.Unmarshal(ClaimAsBytes, &Claim)
	if err != nil {
		return shim.Error("Error encountered during unmarshalling the data")
	}
	Claim.IsVerified = "True"
	ClaimJsonAsBytes, err :=json.Marshal(Claim)
	if err != nil {
		return shim.Error("Error encountered while remarshalling")
	}
	err = stub.PutState(Hash, ClaimJsonAsBytes)
	if err != nil {
		return shim.Error("error encountered while putting state")
	}
	fmt.Println("VERIFIED!!")
	return shim.Success(nil)
}

// Query Function
func  (t *HackidfChaincode) Query(stub shim.ChaincodeStubInterface, args []string)pb.Response{
	DataAsBytes, err := stub.GetState(args[0])
	if err != nil {
		return shim.Error("Error encountered")
	}else if DataAsBytes == nil {
		return shim.Error("No Data")
	}
	return shim.Success(DataAsBytes)
}
// MAIN FUNCTION
func  main() {
	err := shim.Start(new(HackidfChaincode))
	if err != nil {
		fmt.Printf("Error starting Chaincode: %s", err)
	}
}