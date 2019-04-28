package blockchain

import (
        "fmt"
        "github.com/hyperledger/fabric-sdk-go/pkg/client/channel"
        "github.com/hyperledger/fabric-sdk-go/pkg/client/event"
        "github.com/hyperledger/fabric-sdk-go/pkg/core/config"
        "github.com/hyperledger/fabric-sdk-go/pkg/fabsdk"
        "github.com/pkg/errors"
)

// FabricSetup implementation
type FabricSetup struct {
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
        client          *channel.Client
        sdk             *fabsdk.FabricSDK
        event           *event.Client
}

// Initialize reads the configuration file and sets up the client, chain and event hub
func (setup *FabricSetup) Initialize() error {

	// Add parameters for the initialization
	if setup.initialized {
		return errors.New("sdk already initialized")
	}

//  1) Instantiate a fabsdk instance using a configuration.

	// Initialize the SDK with the configuration file
	sdk, err := fabsdk.New(config.FromFile(setup.ConfigFile))
	if err != nil {
		return errors.WithMessage(err, "failed to create SDK")
	}
	setup.sdk = sdk
	fmt.Println("SDK created")

	fmt.Println("Initialization Successful")
	setup.initialized = true
	return nil
}

func (setup *FabricSetup) InstallAndInstantiateCC() error {

	/////////////////           클라이언트 만드는 부분           /////////////////

	//  2) Create a context based on a user and organization, using your fabsdk instance.

	// Channel client is used to query and execute transactions
	clientContext := setup.sdk.ChannelContext(setup.ChannelID, fabsdk.WithUser(setup.OrgAdmin), fabsdk.WithOrg(setup.OrgName))
	// channel client

	client, err := channel.New(clientContext)
	if err != nil {
		return errors.WithMessage(err, "failed to create new channel client")
	}
	setup.client = client
	fmt.Println("Channel client created")

	// Creation of the client which will enables access to our channel events
	// event client
	setup.event, err = event.New(clientContext)
	if err != nil {
		return errors.WithMessage(err, "failed to create new event client")
	}
	fmt.Println("Event client created")

	return nil
}

func (setup *FabricSetup) CloseSDK() {
        setup.sdk.Close()
}
