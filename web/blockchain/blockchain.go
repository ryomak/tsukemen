package blockchain

import (
	"github.com/hyperledger/fabric-sdk-go/pkg/client/channel"
	"github.com/hyperledger/fabric-sdk-go/pkg/client/event"
	"github.com/hyperledger/fabric-sdk-go/pkg/client/resmgmt"
	"github.com/hyperledger/fabric-sdk-go/pkg/fabsdk"
	"github.com/ryomak/tsukemen/web/model"
	"encoding/json"
	"strconv"
	"time"
	"os"
	"fmt"
)

type BlockchainSession struct {
  ConfigFile      string
	OrgID           string
	OrdererID       string
	ChannelID       string
	ChainCodeID     string
	initialized     bool
	ChannelConfig   string
	ChaincodeGoPath string
	ChaincodePath   string
	OrgAdmin        string
	OrgName         string
	UserName        string
	client          *channel.Client
	admin           *resmgmt.Client
	sdk             *fabsdk.FabricSDK
	event           *event.Client
}

func NewBlockchainSession() *BlockchainSession {
    session := BlockchainSession{
		// Network parameters
		OrdererID: "orderer.hf.chainhero.io",

		// Channel parameters
		ChannelID:     "chainhero",
		ChannelConfig: os.Getenv("GOPATH") + "/src/github.com/ryomak/tsukemen/fixtures/artifacts/chainhero.channel.tx",

		// Chaincode parameters
		ChainCodeID:     "tsukemen",
		ChaincodeGoPath: os.Getenv("GOPATH"),
		ChaincodePath:   "github.com/ryomak/tsukemen/web/blockchain/chaincode/",
		OrgAdmin:        "Admin",
		OrgName:         "org1",
		ConfigFile:      "config.yaml",

		// User parameters
		UserName: "User1",
	}
  err := session.Initialize()
	if err != nil {
		fmt.Printf("Unable to initialize the Fabric SDK: %v\n", err)
		return nil
	}
	// Close SDK
	defer session.CloseSDK()

	// Install and instantiate the chaincode
	err = session.InstallAndInstantiateCC()
	if err != nil {
		fmt.Printf("Unable to install and instantiate the chaincode: %v\n", err)
		return  nil
	}
	return &session
}
var number = 0

func (b *BlockchainSession) VoteForCandidate(v model.Vote) error {
	number++
	function := "createVote"
	var args []string
	args = append(args, "user"+strconv.Itoa(number))
	args = append(args, v.User)
	args = append(args, strconv.Itoa(v.CandidateID))

	eventID := "voteForInvoke"

	// Add data that will be visible in the proposal, like a description of the invoke request
	transientDataMap := make(map[string][]byte)
	transientDataMap["result"] = []byte("invoke vote")

	reg, notifier, err := b.event.RegisterChaincodeEvent(b.ChainCodeID, eventID)
	if err != nil {
		return err
	}
	defer b.event.Unregister(reg)
	response, err := b.client.Execute(channel.Request{ChaincodeID: b.ChainCodeID, Fcn: function, Args: [][]byte{[]byte(args[1]), []byte(args[2]), []byte(args[2])}, TransientMap: transientDataMap})
	if err != nil {
		return fmt.Errorf("failed to move funds: %v", err)
	}
		// Wait for the result of the submission
	select {
	case ccEvent := <-notifier:
		fmt.Printf("Received CC event: %v\n", ccEvent)
	case <-time.After(time.Second * 20):
		return fmt.Errorf("did NOT receive CC event for eventId(%s)", eventID)
	}
	fmt.Println(response)
	return nil
}
func (b *BlockchainSession) Result() ([]model.Vote, error) {
	// Prepare arguments
	function := "queryAllVotes"
	response, err := b.client.Query(channel.Request{ChaincodeID: b.ChainCodeID, Fcn: function, Args: [][]byte{}})
	if err != nil {
		return nil, fmt.Errorf("failed to query: %v", err)
	}
	var v []model.Vote
    if err := json.Unmarshal(response.Payload, &v); err != nil {
        return nil ,err
    }
    fmt.Print("result:")
  	fmt.Println(response.Payload)
	//return string(response.Payload), nil
	return v, nil
}
