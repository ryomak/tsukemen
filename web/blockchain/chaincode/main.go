package main

import (
	"fmt"
	"github.com/hyperledger/fabric/core/chaincode/shim"
	pb "github.com/hyperledger/fabric/protos/peer"
	"github.com/ryomak/tsukemen/web/model"
	"strconv"
)

// HeroesServiceChaincode implementation of Chaincode
type VoteChainCode struct {
}

// Init of the chaincode
// This function is called only one when the chaincode is instantiated.
// So the goal is to prepare the ledger to handle future requests.
func (t *VoteChainCode) Init(stub shim.ChaincodeStubInterface) pb.Response {
	fmt.Println("########### HeroesServiceChaincode Init ###########")

	// Get the function and arguments from the request
	function, _ := stub.GetFunctionAndParameters()

	// Check if the request is the init function
	if function != "init" {
		return shim.Error("Unknown function call")
	}
	// Return a successful message
	return shim.Success(nil)
}

// Invoke
// All future requests named invoke will arrive here.
func (t *VoteChainCode) Invoke(stub shim.ChaincodeStubInterface) pb.Response {
	fmt.Println("########### HeroesServiceChaincode Invoke ###########")

	// Get the function and arguments from the request
	function, args := stub.GetFunctionAndParameters()

	if function == "queryVote" {
		return s.queryVote(APIstub, args)
	} else if function == "initVotes" {
		return s.initVotes(APIstub)
	} else if function == "createVote" {
		return s.createVote(APIstub, args)
	} else if function == "queryAllVotes" {
		return s.queryAllVotes(APIstub)
	} else if function == "changeVote" {
		return s.changeVote(APIstub, args)
	}
	// If the arguments given don’t match any function, we return an error
	return shim.Error("Invalid Smart Contract function name")
}

func (t *VoteChainCode) queryVote(APIstub shim.ChaincodeStubInterface, args []string) sc.Response {

	if len(args) != 1 {
		return shim.Error("Incorrect number of arguments. Expecting 1")
	}

	voteAsBytes, _ := APIstub.GetState(args[0])
	return shim.Success(voteAsBytes)
}

func (s *SmartContract) initVotes(APIstub shim.ChaincodeStubInterface) sc.Response {
	votes := []model.Vote{
		model.Vote{User:"user0",CandidateID:0},
	}

	i := 0
	for i < len(votes) {
		fmt.Println("i is ", i)
		voteAsBytes, _ := json.Marshal(votes[i])
		APIstub.PutState("Vote"+strconv.Itoa(i), voteAsBytes)
		fmt.Println("Added", votes[i])
		i = i + 1
	}

	return shim.Success(nil)
}

func (s *SmartContract) createVote(APIstub shim.ChaincodeStubInterface, args []string) sc.Response {

	if len(args) != 3 {
		return shim.Error("Incorrect number of arguments. Expecting 2")
	}
	vote := model.Vote{User:args[1],CandidateID:strconv.Itoa(args[2])}
	voteAsBytes, _ := json.Marshal(vote)
	err := stub.PutState(args[0], [voteAsBytes)

	if err != nil {
		return shim.Error("Failed to update state of hello")
	}

	
	err = stub.SetEvent("voteForInvoke", []byte{})
	if err != nil {
		return shim.Error(err.Error())
	}

	return shim.Success(nil)
}

func (s *SmartContract) queryAllVotes(APIstub shim.ChaincodeStubInterface) sc.Response {

	startKey := "Vote0"
	endKey := "Vote99999"

	resultsIterator, err := APIstub.GetStateByRange(startKey, endKey)
	if err != nil {
		return shim.Error(err.Error())
	}
	defer resultsIterator.Close()

	// buffer is a JSON array containing QueryResults
	var buffer bytes.Buffer
	buffer.WriteString("[")

	bArrayMemberAlreadyWritten := false
	for resultsIterator.HasNext() {
		queryResponse, err := resultsIterator.Next()
		if err != nil {
			return shim.Error(err.Error())
		}
		// Add a comma before array members, suppress it for the first array member
		if bArrayMemberAlreadyWritten == true {
			buffer.WriteString(",")
		}
		buffer.WriteString("{\"Key\":")
		buffer.WriteString("\"")
		buffer.WriteString(queryResponse.Key)
		buffer.WriteString("\"")

		buffer.WriteString(", \"Record\":")
		// Record is a JSON object, so we write as-is
		buffer.WriteString(string(queryResponse.Value))
		buffer.WriteString("}")
		bArrayMemberAlreadyWritten = true
	}
	buffer.WriteString("]")

	fmt.Printf("- queryAllVotes:\n%s\n", buffer.String())

	return shim.Success(buffer.Bytes())
}

func (s *SmartContract) changeVote(APIstub shim.ChaincodeStubInterface, args []string) sc.Response {

	if len(args) != 2 {
		return shim.Error("Incorrect number of arguments. Expecting 2")
	}

	voteAsBytes, _ := APIstub.GetState(args[0])
	vote := model.Vote{}

	json.Unmarshal(carAsBytes, &vote)
	vote.CandidateID = args[1]

	voteAsBytes, _ = json.Marshal(vote)
	APIstub.PutState(args[0], voteAsBytes)

	return shim.Success(nil)
}

func main() {
	// Start the chaincode and make it ready for futures requests
	err := shim.Start(new(VoteChainCode))
	if err != nil {
		fmt.Printf("Error starting Heroes Service chaincode: %s", err)
	}
}
