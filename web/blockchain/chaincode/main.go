package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/hyperledger/fabric/core/chaincode/shim"
	pb "github.com/hyperledger/fabric/protos/peer"
)

//same as 	"github.com/ryomak/tsukemen/web/model"
type Vote struct {
	UserName      string `json:"user_name"`
	CandidateName string `json:"candidate_name"`
}

// HeroesServiceChaincode implementation of Chaincode
type VoteChainCode struct {
}

// Init of the chaincode
// This function is called only one when the chaincode is instantiated.
// So the goal is to prepare the ledger to handle future requests.
func (v *VoteChainCode) Init(stub shim.ChaincodeStubInterface) pb.Response {
	fmt.Println("########### HeroesServiceChaincode Init ###########")

	// Get the function and arguments from the request
	function, _ := stub.GetFunctionAndParameters()

	// Check if the request is the init function
	if function != "init" {
		return shim.Error("Unknown function call")
	}
	return v.initVotes(stub)
}

// Invoke
func (v *VoteChainCode) Invoke(stub shim.ChaincodeStubInterface) pb.Response {
	fmt.Println("########### HeroesServiceChaincode Invoke ###########")

	function, args := stub.GetFunctionAndParameters()

	if function == "initVotes" {
		return v.initVotes(stub)
	} else if function == "createVote" {
		return v.createVote(stub, args)
	} else if function == "queryAllVotes" {
		return v.queryAllVotes(stub)
	}
	return shim.Error("Invalid Smart Contract function name")
}

func (v *VoteChainCode) initVotes(stub shim.ChaincodeStubInterface) pb.Response {
	votes := []Vote{
		Vote{UserName: "user0", CandidateName: "candidate1"},
	}

	i := 0
	for i < len(votes) {
		fmt.Println("i is ", i)
		voteAsBytes, _ := json.Marshal(votes[i])
		stub.PutState("Vote"+strconv.Itoa(i), voteAsBytes)
		fmt.Println("Added", votes[i])
		i = i + 1
	}

	return shim.Success(nil)
}

func (v *VoteChainCode) createVote(stub shim.ChaincodeStubInterface, args []string) pb.Response {

	if len(args) != 3 {
		return shim.Error("Incorrect number of arguments. Expecting 2")
	}
	vote := Vote{UserName: args[1], CandidateName: args[2]}
	voteAsBytes, _ := json.Marshal(vote)
	err := stub.PutState(args[0], voteAsBytes)
	if err != nil {
		return shim.Error("Failed to update state of hello")
	}
	return shim.Success(nil)
}

func (v *VoteChainCode) queryAllVotes(stub shim.ChaincodeStubInterface) pb.Response {

	startKey := "Vote1"
	endKey := "Vote999"
	resultsIterator, err := stub.GetStateByRange(startKey, endKey)
	if err != nil {
		return shim.Error(err.Error())
	}
	defer resultsIterator.Close()

	// buffer is a JSON array containing QueryResults
	var buffer bytes.Buffer

	bArrayMemberAlreadyWritten := false
	buffer.WriteString("[")
	for resultsIterator.HasNext() {
		queryResponse, err := resultsIterator.Next()
		if err != nil {
			return shim.Error(err.Error())
		}
		if bArrayMemberAlreadyWritten == true {
			buffer.WriteString(",")
		}
		buffer.WriteString(string(queryResponse.Value))
		bArrayMemberAlreadyWritten = true
	}
	buffer.WriteString("]")

	return shim.Success(buffer.Bytes())
}
func main() {
	// Start the chaincode and make it ready for futures requests
	err := shim.Start(new(VoteChainCode))
	if err != nil {
		fmt.Printf("Error starting Heroes Service chaincode: %s", err)
	}
}
