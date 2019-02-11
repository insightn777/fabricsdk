package blockchain

import (
	"fmt"
	"github.com/hyperledger/fabric-sdk-go/pkg/client/channel"
	//"time"
)

// InvokeHello
func (setup *FabricSetup) InsertUser(value ...string) (string, error) {

	// Prepare arguments
	var args []string
	args = append(args, "insertUser")
	args = append(args, value[0]) //ID
	args = append(args, value[1]) //Name
	args = append(args, value[2]) //Password
	args = append(args, value[3]) //SocialNumber
	args = append(args, value[4]) //Location

	//eventID := "eventInvokeUser"

	// Add data that will be visible in the proposal, like a description of the invoke request
	transientDataMap := make(map[string][]byte)
	transientDataMap["result"] = []byte("Transient data in hello invoke")

	/*reg, notifier, err := setup.event.RegisterChaincodeEvent(setup.ChainCodeID, eventID)
	if err != nil {
		return "", err
	}
	defer setup.event.Unregister(reg)*/

	for i := 0; i < 6; i++ {
		fmt.Printf(args[i])
	}
	// Create a request (proposal) and send it
	response, err := setup.client.Execute(channel.Request{ChaincodeID: setup.ChainCodeID, Fcn: args[0], Args: [][]byte{[]byte(args[1]), []byte(args[2]), []byte(args[3]), []byte(args[4]), []byte(args[5])}, TransientMap: transientDataMap})
	if err != nil {
		return "", fmt.Errorf("failed to move funds: %v", err)
	}
	// Wait for the result of the submission
	/*select {
	case ccEvent := <-notifier:
		fmt.Printf("Received CC event: %v\n", ccEvent)
	case <-time.After(time.Second * 20):
		fmt.Println("this error!!!")
		return "", fmt.Errorf("did NOT receive CC event for eventId(%s)", eventID)
	}*/
return string(response.TransactionID), nil
}

func (setup *FabricSetup) InsertVote(value ...string) (string, error) {
	// Prepare arguments
	var args []string
	args = append(args, "insertVote")
	args = append(args, value[0]) //Votename
	args = append(args, value[1]) //StartDate
	args = append(args, value[2]) //EndDate
	args = append(args, value[3]) //Location
	for i := 0; i < len(value) - 4; i++ {
		args = append(args, value[i + 4]) //Questiins...
	}

	fmt.Printf("Argument : ")
	fmt.Println(args[1] + args[2] + args[3] + args[4] + args[5] + args[6] + args[7] + args[8])

	//eventID := "eventInsertVote"
	// Add data that will be visible in the proposal, like a description of the invoke request
	transientDataMap := make(map[string][]byte)
	transientDataMap["result"] = []byte("Transient data in hello invoke")
	/*reg, notifier, err := setup.event.RegisterChaincodeEvent(setup.ChainCodeID, eventID)
	if err != nil {
		return "", err
	}
	defer setup.event.Unregister(reg)*/

	// Create a request (proposal) and send it
	response, err := setup.client.Execute(channel.Request{ChaincodeID: setup.ChainCodeID, Fcn: args[0], Args: [][]byte{[]byte(args[1]), []byte(args[2]), []byte(args[3]), []byte(args[4]), []byte(args[5]), []byte(args[6]), []byte(args[7]), []byte(args[8])}, TransientMap: transientDataMap})
	if err != nil {
		return "", fmt.Errorf("failed to move funds: %v", err)
	}
	// Wait for the result of the submission
	/*select {
	case ccEvent := <-notifier:
		fmt.Printf("Received CC event: %v\n", ccEvent)
	case <-time.After(time.Second * 20):
		return "", fmt.Errorf("did NOT receive CC event for eventId(%s)", eventID)
	}*/
	return string(response.TransactionID), nil
}

func (setup *FabricSetup) InsertVoteResult(value ...string) (string, error) {
	// Prepare arguments
	var args []string
	args = append(args, "insertVoteResult")
	args = append(args, value[0]) //Votename
	args = append(args, value[1]) //ID
	args = append(args, value[2]) //Result

	//eventID := "eventInsertVoteResult"
	// Add data that will be visible in the proposal, like a description of the invoke request
	transientDataMap := make(map[string][]byte)
	transientDataMap["result"] = []byte("Transient data in hello invoke")
	/*reg, notifier, err := setup.event.RegisterChaincodeEvent(setup.ChainCodeID, eventID)
	if err != nil {
		return "", err
	}
	defer setup.event.Unregister(reg)*/
	// Create a request (proposal) and send it
	response, err := setup.client.Execute(channel.Request{ChaincodeID: setup.ChainCodeID, Fcn: args[0], Args: [][]byte{[]byte(args[1]), []byte(args[2]), []byte(args[3])}, TransientMap: transientDataMap})
	if err != nil {
		return "", fmt.Errorf("failed to move funds: %v", err)
	}
	// Wait for the result of the submission
	/*select {
	case ccEvent := <-notifier:
		fmt.Printf("Received CC event: %v\n", ccEvent)
	case <-time.After(time.Second * 20):
		return "", fmt.Errorf("did NOT receive CC event for eventId(%s)", eventID)
	}*/
	return string(response.TransactionID), nil
}

func (setup *FabricSetup) DeleteUser(value string) (string, error) {
        // Prepare arguments
        var args []string
        args = append(args, "delete")
        args = append(args, value) //ID

        //eventID := "eventInsertVoteResult"
        // Add data that will be visible in the proposal, like a description of the invoke request
        transientDataMap := make(map[string][]byte)
        transientDataMap["result"] = []byte("Transient data in hello invoke")
        /*reg, notifier, err := setup.event.RegisterChaincodeEvent(setup.ChainCodeID, eventID)
        if err != nil {
                return "", err
        }
        defer setup.event.Unregister(reg)*/
        // Create a request (proposal) and send it
        response, err := setup.client.Execute(channel.Request{ChaincodeID: setup.ChainCodeID, Fcn: args[0], Args: [][]byte{[]byte(args[1])}, TransientMap: transientDataMap})
        if err != nil {
                return "", fmt.Errorf("failed to move funds: %v", err)
        }
        // Wait for the result of the submission
        /*select {
        case ccEvent := <-notifier:
                fmt.Printf("Received CC event: %v\n", ccEvent)
        case <-time.After(time.Second * 20):
                return "", fmt.Errorf("did NOT receive CC event for eventId(%s)", eventID)
        }*/
        return string(response.TransactionID), nil
}

