package main

import (
	"bufio"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"math/big"
	"net"
	"strconv"

	"github.com/hyperledger/fabric/core/chaincode/shim"
	pb "github.com/hyperledger/fabric/protos/peer"
)

var (
	fileName = "interoperability"
)

type Interoperability struct {
}

type Asset struct {
	ID    string `json:"ID"`
	Owner string `json:"Owner"`
}

func GetHeight(IP string) int {
	url := IP + ":8080"
	conn, _ := net.Dial("tcp", url)
	message, _ := bufio.NewReader(conn).ReadString('\n')
	height, _ := strconv.Atoi(message)
	return height
}

type Proposal struct {
	PrivateHashX        string  `json:"PrivateHashX"`
	PrivateHashY        string  `json:"PrivateHashY"`
	TransactionHash     string  `json:"TransactionHash"`
	TransactionLockHash string  `json:"TransactionLockHash"`
	From                string  `json:"From"`
	To                  string  `json:"To"`
	Value               string  `json:"Value"`
	X                   string  `json:"x"`
	Y                   string  `json:"y"`
	Secret              string  `json"s"`
	BlockHeight         int     `json"BlockHeight"`
	IsConfirmed         float64 `json:"IsConfirmed"`
	IsInvalidated       float64 `json:"IsInvalidated"`
}

//initialization function
func (t *Interoperability) Init(stub shim.ChaincodeStubInterface) pb.Response {
	// Whatever variable initialisation you want can be done here //
	return shim.Success(nil)
}

// invoking functions
func (t *Interoperability) Invoke(stub shim.ChaincodeStubInterface) pb.Response {
	// IF-ELSE-IF all the functions
	function, args := stub.GetFunctionAndParameters()
	if function == "CreateProposal" {
		return t.CreateProposal(stub, args)
	} else if function == "StartTx" {
		return t.StartTx(stub, args)
	} else if function == "CommitLockSecret" {
		return t.CommitLockSecret(stub, args)
	} else if function == "CommitPrivateSecret" {
		return t.CommitPrivateSecret(stub, args)
	} else if function == "ConfirmProposal" {
		return t.ConfirmProposal(stub, args)
	} else if function == "Query" {
		return t.Query(stub, args)
	} else if function == "CreateAsset" {
		return t.CreateAsset(stub, args)
	} else if function == "InvalidateProposal" {
		return t.InvalidateProposal(stub, args)
	}

	fmt.Println("invoke did not find func : " + function) //error
	return shim.Error("Received unknown function invocation")
	// end of all functions
}

// CreateAsset issues a new asset to the world state with given details.
func (t *Interoperability) CreateAsset(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	var ID = args[0]
	var Owner = args[1]

	AssetAsBytes, err := stub.GetState(ID)
	if err != nil {
		return shim.Error("Failed to get Proposal:" + err.Error())
	} else if AssetAsBytes != nil {
		return shim.Error("Proposal already exists")
	}
	var Asset = &Asset{ID: ID, Owner: Owner}
	AssetJsonAsBytes, err := json.Marshal(Asset)
	if err != nil {
		shim.Error("Error encountered while Marshalling")
	}
	err = stub.PutState(ID, AssetJsonAsBytes)
	if err != nil {
		shim.Error("Error encountered while Creating Proposal")
	}
	fmt.Println("Ledger Updated Successfully")
	return shim.Success(nil)
}

// Initiate transaction by commiting g^x or g^y against transaction-id
func (t *Interoperability) StartTx(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	var TxID = args[0]
	var PrivateHashX = args[1]
	var PrivateHashY = args[2]

	ProposalAsBytes, err := stub.GetState(TxID)
	if err != nil {
		return shim.Error("Failed to get Proposal:" + err.Error())
	} else if ProposalAsBytes != nil {
		return shim.Error("Proposal already exists")
	}
	var Proposal = &Proposal{PrivateHashX: PrivateHashX, PrivateHashY: PrivateHashY, TransactionHash: "", From: "", To: "", Value: "", X: "", Y: "", Secret: "", BlockHeight: 0, IsConfirmed: 0, IsInvalidated: 0}
	ProposalJsonAsBytes, err := json.Marshal(Proposal)
	if err != nil {
		shim.Error("Error encountered while Marshalling")
	}
	err = stub.PutState(TxID, ProposalJsonAsBytes)
	if err != nil {
		shim.Error("Error encountered while Creating Proposal")
	}
	fmt.Println("Ledger Updated Successfully")
	return shim.Success(nil)
}

// Create a proposal to transfer assets
func (t *Interoperability) CreateProposal(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	var TxID = args[0]
	var TransactionHash = args[1]
	var TransactionLockHash = args[2]
	var From = args[3]
	var To = args[4]
	var AssetValue = args[5]
	var IP = args[6]

	BlockHeight := GetHeight(IP)

	ProposalAsBytes, err := stub.GetState(TxID)
	if err != nil {
		return shim.Error("Failed to get Proposal:" + err.Error())
	} else if ProposalAsBytes == nil {
		return shim.Error("No proposal with the current hash")
	}

	var Proposal Proposal
	err = json.Unmarshal(ProposalAsBytes, &Proposal)
	if err != nil {
		return shim.Error("Error encountered during unmarshalling the data")
	}

	Proposal.TransactionHash = TransactionHash
	Proposal.TransactionLockHash = TransactionLockHash
	Proposal.From = From
	Proposal.To = To
	Proposal.Value = AssetValue
	Proposal.BlockHeight = BlockHeight
	Proposal.IsConfirmed = 0
	Proposal.IsInvalidated = 0

	ProposalJsonAsBytes, err := json.Marshal(Proposal)
	if err != nil {
		shim.Error("Error encountered while Re-marshalling")
	}
	err = stub.PutState(TxID, ProposalJsonAsBytes)
	if err != nil {
		shim.Error("Error encountered while Creating Proposal")
	}
	fmt.Println("Ledger Updated Successfully")
	return shim.Success(nil)
}

func (t *Interoperability) CommitLockSecret(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	var TxID = args[0]
	var Secret = args[1]
	var IP = args[2]
	SecretHashBytes := sha256.Sum256([]byte(Secret))
	SecretHash := fmt.Sprintf("%x", SecretHashBytes)

	BlockHeight := GetHeight(IP)
	ProposalAsBytes, err := stub.GetState(TxID)
	if err != nil {
		return shim.Error("Failed to get Proposal:" + err.Error())
	} else if ProposalAsBytes == nil {
		return shim.Error("No proposal with the current hash")
	}

	var Proposal Proposal
	err = json.Unmarshal(ProposalAsBytes, &Proposal)
	if err != nil {
		return shim.Error("Error encountered during unmarshalling the data")
	}
	if BlockHeight-Proposal.BlockHeight > 5 {
		return shim.Error("Too late to commit secret")
	}
	if Proposal.TransactionLockHash == SecretHash {

		// Implement the Time Lock concepts
		// Block Height should not be more than n/2

		Proposal.Secret = Secret
		ProposalJsonAsBytes, err := json.Marshal(Proposal)
		if err != nil {
			return shim.Error("Error encountered while remarshalling")
		}
		err = stub.PutState(TxID, ProposalJsonAsBytes)
		if err != nil {
			return shim.Error("error encountered while putting state")
		}
		fmt.Println("UNLOCKED")
		return shim.Success(nil)
	}
	fmt.Println("Wrong Secret. Couldn't unlock.")
	return shim.Error("Hash does not match\n Real Hash: " + Proposal.TransactionLockHash + "\n Secret Hash" + SecretHash)
}

func (t *Interoperability) CommitPrivateSecret(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	var TxID = args[0]
	var X = args[1]

	ProposalAsBytes, err := stub.GetState(TxID)
	if err != nil {
		return shim.Error("Failed to get Proposal:" + err.Error())
	} else if ProposalAsBytes == nil {
		return shim.Error("No proposal with the current hash")
	}

	var Proposal Proposal
	err = json.Unmarshal(ProposalAsBytes, &Proposal)
	if err != nil {
		return shim.Error("Error encountered during unmarshalling the data")
	}

	PrivateBigInt := new(big.Int)
	PrivateBigInt.SetString(X, 16)

	PubXBigInt := new(big.Int)
	PubXBigInt.SetString(Proposal.PrivateHashX, 16)

	PubYBigInt := new(big.Int)
	PubYBigInt.SetString(Proposal.PrivateHashY, 16)

	privc, _ := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	pubc := privc.PublicKey

	SecretHashBytes, _ := pubc.Curve.ScalarMult(PubXBigInt, PubYBigInt, PrivateBigInt.Bytes())
	shared := sha256.Sum256(SecretHashBytes.Bytes())
	SecretHash := fmt.Sprintf("%x", shared)

	if Proposal.TransactionHash == SecretHash && Proposal.Secret != "" {

		// Implement the Time Lock concepts
		// Block Height should not be more than n

		Proposal.X = X
		ProposalJsonAsBytes, err := json.Marshal(Proposal)
		if err != nil {
			return shim.Error("Error encountered while remarshalling")
		}
		err = stub.PutState(TxID, ProposalJsonAsBytes)
		if err != nil {
			return shim.Error("error encountered while putting state")
		}
		fmt.Println("UNLOCKED")
		return shim.Success(nil)
	}
	fmt.Println("Wrong Secret or prereq not finished")
	return shim.Error("Hash does not match")
}

// Confirm Proposal
func (t *Interoperability) ConfirmProposal(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	var TxID = args[0]

	ProposalAsBytes, err := stub.GetState(TxID)
	if err != nil {
		return shim.Error("Failed to get Proposal:" + err.Error())
	} else if ProposalAsBytes == nil {
		return shim.Error("Please give a valid TransactionHash")
	}
	var Proposal Proposal
	err = json.Unmarshal(ProposalAsBytes, &Proposal)
	if err != nil {
		return shim.Error("Error encountered during unmarshalling the data")
	}

	if Proposal.Secret != "" && Proposal.X != "" {

		// Can include transaction based on the use case
		// For example instead of just flagging the transaction
		// as complete you can transfer the tokens.
		AssetAsBytes, err := stub.GetState(Proposal.Value)
		if err != nil {
			return shim.Error("Failed to get Proposal:" + err.Error())
		} else if ProposalAsBytes == nil {
			return shim.Error("Please give a valid Asset does not exist")
		}
		var Asset Asset
		err = json.Unmarshal(AssetAsBytes, &Asset)
		if err != nil {
			return shim.Error("Error encountered during unmarshalling the data")
		}
		if Asset.Owner != Proposal.From {
			return shim.Error("Owner of the asset is not the same person to create the proposal")
		}
		Asset.Owner = Proposal.To
		AssetJsonAsBytes, err := json.Marshal(Asset)
		err = stub.PutState(Proposal.Value, AssetJsonAsBytes)
		if err != nil {
			return shim.Error("error encountered while putting state")
		}

		Proposal.IsConfirmed = 1
		ProposalJsonAsBytes, err := json.Marshal(Proposal)
		if err != nil {
			return shim.Error("Error encountered while remarshalling")
		}
		err = stub.PutState(TxID, ProposalJsonAsBytes)
		if err != nil {
			return shim.Error("error encountered while putting state")
		}
		fmt.Println("VERIFIED!!")
		return shim.Success(nil)
	}
	fmt.Println("Wrong secret or prereq not finished")
	return shim.Error("Tx not complete")
}

// Invalidate Proposal
func (t *Interoperability) InvalidateProposal(stub shim.ChaincodeStubInterface, args []string) pb.Response {

	// Check if Time lock conditions holds good.
	// You can implement this by either using certificates or Block Height.

	var TxID = args[0]
	var IP = args[1]

	BlockHeight := GetHeight(IP)
	ProposalAsBytes, err := stub.GetState(TxID)
	if err != nil {
		return shim.Error("Failed to get Proposal:" + err.Error())
	} else if ProposalAsBytes == nil {
		return shim.Error("Please give a valid TxID")
	}
	var Proposal Proposal
	err = json.Unmarshal(ProposalAsBytes, &Proposal)
	if err != nil {
		return shim.Error("Error encountered during unmarshalling the data")
	}

	if BlockHeight-Proposal.BlockHeight < 10 {
		return shim.Error("Have to wait before invalidating proposal")
	}
	Proposal.IsInvalidated = 1
	ProposalJsonAsBytes, err := json.Marshal(Proposal)
	if err != nil {
		return shim.Error("Error encountered while remarshalling")
	}
	err = stub.PutState(TxID, ProposalJsonAsBytes)
	if err != nil {
		return shim.Error("error encountered while putting state")
	}
	fmt.Println("INVALIDATED!!")
	return shim.Success(nil)
}
func (t *Interoperability) Query(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	DataAsBytes, err := stub.GetState(args[0])
	if err != nil {
		return shim.Error("Error encountered")
	} else if DataAsBytes == nil {
		return shim.Error("No Data")
	}
	return shim.Success(DataAsBytes)
}

func main() {
	err := shim.Start(new(Interoperability))
	if err != nil {
		fmt.Printf("Error starting Chaincode: %s", err)
	}
}
