package blockchain

import (
	"fmt"
	"github.com/hyperledger/fabric-sdk-go/pkg/client/channel"
)

// QueryHello query the chaincode to get the state of hello
func (setup *FabricSetup) QueryAllUser() (string, error) {

	// Prepare arguments
	var args []string
	args = append(args, "query")
	args = append(args, "QueryAllUser")
	response, err := setup.client.Query(channel.Request{ChaincodeID: setup.ChainCodeID, Fcn: args[0], Args: [][]byte{[]byte(args[1])}})
	if err != nil {
		return "", fmt.Errorf("failed to query: %v", err)
	}

	fmt.Println("data : " + string(response.Payload))
	return string(response.Payload), nil
}

func (setup *FabricSetup) QueryUserByName(name string) (string, error) {

	// Prepare arguments
	var args []string
	args = append(args, "query")
	args = append(args, "QueryUserByName")
	args = append(args, name)

	fmt.Println("passed name : " + name)
	response, err := setup.client.Query(channel.Request{ChaincodeID: setup.ChainCodeID, Fcn: args[0], Args: [][]byte{[]byte(args[1]), []byte(args[2])}})
	if err != nil {
		fmt.Println("this error first?")
		return "", fmt.Errorf("failed to query: %v", err)
	}

	fmt.Println("dat : " + string(response.Payload))
	return string(response.Payload), nil
}

func (setup *FabricSetup) QueryAllVote() (string, error) {
	// Prepare arguments
	var args []string
	args = append(args, "query")
	args = append(args, "QueryAllVote")

	response, err := setup.client.Query(channel.Request{ChaincodeID: setup.ChainCodeID, Fcn: args[0], Args: [][]byte{[]byte(args[1])}})
	if err != nil {
		return "", fmt.Errorf("failed to query: %v", err)
	}
	fmt.Println("data : " + string(response.Payload))
	return string(response.Payload), nil
}

func (setup *FabricSetup) QueryVoteByName(name string) (string, error) {

	// Prepare arguments
	var args []string
	args = append(args, "query")
	args = append(args, "QueryVoteByName")
	args = append(args, name)

	response, err := setup.client.Query(channel.Request{ChaincodeID: setup.ChainCodeID, Fcn: args[0], Args: [][]byte{[]byte(args[1]), []byte(args[2])}})
	if err != nil {
		return "", fmt.Errorf("failed to query: %v", err)
	}

	fmt.Println("data : " + string(response.Payload))
	return string(response.Payload), nil
}

func (setup *FabricSetup) QueryVoteByLocation(location string) (string, error) {

        // Prepare arguments
        var args []string
        args = append(args, "query")
        args = append(args, "QueryVoteByLocation")
        args = append(args, location)

        response, err := setup.client.Query(channel.Request{ChaincodeID: setup.ChainCodeID, Fcn: args[0], Args: [][]byte{[]byte(args[1]), []byte(args[2])}})
        if err != nil {
                return "", fmt.Errorf("failed to query: %v", err)
        }

        fmt.Println("data : " + string(response.Payload))
        return string(response.Payload), nil
}
