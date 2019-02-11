package blockchain

import (
        "fmt"
        "github.com/hyperledger/fabric-sdk-go/test/integration"
        "github.com/hyperledger/fabric-sdk-go/pkg/client/channel"
        "github.com/hyperledger/fabric-sdk-go/pkg/client/event"
        mspclient "github.com/hyperledger/fabric-sdk-go/pkg/client/msp"
        "github.com/hyperledger/fabric-sdk-go/pkg/client/resmgmt"
        "github.com/hyperledger/fabric-sdk-go/pkg/common/errors/retry"
        "github.com/hyperledger/fabric-sdk-go/pkg/common/providers/msp"
        "github.com/hyperledger/fabric-sdk-go/pkg/core/config"
        packager "github.com/hyperledger/fabric-sdk-go/pkg/fab/ccpackager/gopackager"
        "github.com/hyperledger/fabric-sdk-go/pkg/fabsdk"
        "github.com/hyperledger/fabric-sdk-go/third_party/github.com/hyperledger/fabric/common/cauthdsl"
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
        UserName        string
        client          *channel.Client
        admin           *resmgmt.Client
        sdk             *fabsdk.FabricSDK
        event           *event.Client
}

// Initialize reads the configuration file and sets up the client, chain and event hub
func (setup *FabricSetup) Initialize() error {
        var setup2, setup3 FabricSetup
        
        //조직1
        // Add parameters for the initialization
        if setup.initialized {
                return errors.New("1. sdk already initialized")
        }

        // Initialize the SDK with the configuration file
        sdk, err := fabsdk.New(config.FromFile(setup.ConfigFile))
        if err != nil {
                return errors.WithMessage(err, "1. failed to create SDK")
        }
        setup.sdk = sdk
        fmt.Println("1. SDK created")

        // The resource management client is responsible for managing channels (create/update channel)
        resourceManagerClientContext := setup.sdk.Context(fabsdk.WithUser(setup.OrgAdmin), fabsdk.WithOrg(setup.OrgName))
        if err != nil {
                return errors.WithMessage(err, "1. failed to load Admin identity")
        }
        resMgmtClient, err := resmgmt.New(resourceManagerClientContext)
        if err != nil {
                return errors.WithMessage(err, "1. failed to create channel management client from Admin identity")
        }
        setup.admin = resMgmtClient
        fmt.Println("1. Ressource management client created")

        // The MSP client allow us to retrieve user information from their identity, like its signing identity which we will need to save the channel
        mspClient, err := mspclient.New(sdk.Context(), mspclient.WithOrg(setup.OrgName))
        if err != nil {
                return errors.WithMessage(err, "1. failed to create MSP client")
        }
        adminIdentity, err := mspClient.GetSigningIdentity(setup.OrgAdmin)
        if err != nil {
                return errors.WithMessage(err, "1. failed to get admin signing identity")
        }
        
        //조직2
        setup2 = FabricSetup{
                "config.yaml",
                "",
                "orderer.hf.chainhero.io",
                "chainhero",
                "test1",
                false,
                "/root/go/src/github.com/chainHero/heroes-service/fixtures/artifacts/chainhero.channel.tx",
                "/root/go", 
                "github.com/chainHero/heroes-service/chaincode/",
                "Admin",
                "org2",
                "User1",
                nil,
                nil,
                nil,
                nil,
        }
        sdk2, err := fabsdk.New(config.FromFile(setup2.ConfigFile))
        if err != nil {
                return errors.WithMessage(err, "2. failed to create SDK")
        }
        setup2.sdk = sdk2
        fmt.Println("2. SDK created")

        // The resource management client is responsible for managing channels (create/update channel)
        resourceManagerClientContext2 := setup2.sdk.Context(fabsdk.WithUser(setup2.OrgAdmin), fabsdk.WithOrg(setup2.OrgName))
        if err != nil {
                return errors.WithMessage(err, "2. failed to load Admin identity")
        }
        resMgmtClient2, err := resmgmt.New(resourceManagerClientContext2)
        if err != nil {
                return errors.WithMessage(err, "2. failed to create channel management client from Admin identity")
        }
        setup2.admin = resMgmtClient2
        fmt.Println("2. Ressource management client created")

        // The MSP client allow us to retrieve user information from their identity, like its signing identity which we will need to save the channel
        mspClient2, err := mspclient.New(sdk2.Context(), mspclient.WithOrg(setup2.OrgName))
        if err != nil {
                return errors.WithMessage(err, "2. failed to create MSP client")
        }
        adminIdentity2, err := mspClient2.GetSigningIdentity(setup2.OrgAdmin)
        if err != nil {
                return errors.WithMessage(err, "2. failed to get admin signing identity")
        }
        
        //조직3
        setup3 = FabricSetup{
                "config.yaml",
                "",
                "orderer.hf.chainhero.io",
                "chainhero",
                "test1",
                false,
                "/root/go/src/github.com/chainHero/heroes-service/fixtures/artifacts/chainhero.channel.tx",
                "/root/go", 
                "github.com/chainHero/heroes-service/chaincode/",
                "Admin",
                "org3",
                "User1",
                nil,
                nil,
                nil,
                nil,
        }
        sdk3, err := fabsdk.New(config.FromFile(setup3.ConfigFile))
        if err != nil {
                return errors.WithMessage(err, "3. failed to create SDK")
        }
        setup3.sdk = sdk3
        fmt.Println("3. SDK created")

        // The resource management client is responsible for managing channels (create/update channel)
        resourceManagerClientContext3 := setup3.sdk.Context(fabsdk.WithUser(setup3.OrgAdmin), fabsdk.WithOrg(setup3.OrgName))
        if err != nil {
                return errors.WithMessage(err, "3. failed to load Admin identity")
        }
        resMgmtClient3, err := resmgmt.New(resourceManagerClientContext3)
        if err != nil {
                return errors.WithMessage(err, "3. failed to create channel management client from Admin identity")
        }
        setup3.admin = resMgmtClient3
        fmt.Println("3. Ressource management client created")

        // The MSP client allow us to retrieve user information from their identity, like its signing identity which we will need to save the channel
        mspClient3, err := mspclient.New(sdk3.Context(), mspclient.WithOrg(setup3.OrgName))
        if err != nil {
                return errors.WithMessage(err, "3. failed to create MSP client")
        }
        adminIdentity3, err := mspClient3.GetSigningIdentity(setup3.OrgAdmin)
        if err != nil {
                return errors.WithMessage(err, "3. failed to get admin signing identity")
        }

        req := resmgmt.SaveChannelRequest{ChannelID: "chainhero", ChannelConfigPath: setup.ChannelConfig, SigningIdentities: []msp.SigningIdentity{adminIdentity, adminIdentity2, adminIdentity3}}
        txID, err := setup.admin.SaveChannel(req, resmgmt.WithOrdererEndpoint(setup.OrdererID))
        if err != nil || txID.TransactionID == "" {
                return errors.WithMessage(err, "failed to save channel")
        }
        fmt.Println("Channel created")

        // Make admin user join the previously created channel
        if err = setup.admin.JoinChannel(setup.ChannelID, resmgmt.WithRetry(retry.DefaultResMgmtOpts), resmgmt.WithOrdererEndpoint(setup.OrdererID)); err != nil {
                return errors.WithMessage(err, "1. failed to make admin join channel")
        }
        fmt.Println("1. Channel joined")

        if err = setup2.admin.JoinChannel(setup.ChannelID, resmgmt.WithRetry(retry.DefaultResMgmtOpts), resmgmt.WithOrdererEndpoint(setup.OrdererID)); err != nil {
                return errors.WithMessage(err, "failed to make admin join channel")
        }
        fmt.Println("2. Channel joined")

        if err = setup3.admin.JoinChannel(setup.ChannelID, resmgmt.WithRetry(retry.DefaultResMgmtOpts), resmgmt.WithOrdererEndpoint(setup.OrdererID)); err != nil {
                return errors.WithMessage(err, "failed to make admin join channel")
        }
        fmt.Println("3. Channel joined")
        fmt.Println("Initialization Successful")
        setup.initialized = true

        //install
        fmt.Println("Try install")
        ccPkg, err := packager.NewCCPackage(setup.ChaincodePath, setup.ChaincodeGoPath)
        if err != nil {
                return errors.WithMessage(err, "failed to create chaincode package")
        }
        installCCReq := resmgmt.InstallCCRequest{Name: setup.ChainCodeID, Path: setup.ChaincodePath, Version: "0", Package: ccPkg}
        fmt.Println("ccPkg created")
        _, err = integration.DiscoverLocalPeers(resourceManagerClientContext, 2)
        _, err = integration.DiscoverLocalPeers(resourceManagerClientContext2, 2)
        _, err = integration.DiscoverLocalPeers(resourceManagerClientContext3, 2)

        _, err = resMgmtClient.InstallCC(installCCReq, resmgmt.WithRetry(retry.DefaultResMgmtOpts))
        _, err = resMgmtClient2.InstallCC(installCCReq, resmgmt.WithRetry(retry.DefaultResMgmtOpts))
        _, err = resMgmtClient3.InstallCC(installCCReq, resmgmt.WithRetry(retry.DefaultResMgmtOpts))
        if err != nil {
		fmt.Println("install error : ", err)
		return errors.WithMessage(err, "install error")
	}
	fmt.Println("install cc")

        //instantiate
        fmt.Println("Try instantiate")
        ccPolicy := cauthdsl.SignedByAnyMember([]string{"org1.hf.chainhero.io", "org2.hf.chainhero.io", "org3.hf.chainhero.io"})

        resp, err := setup.admin.InstantiateCC(setup.ChannelID, resmgmt.InstantiateCCRequest{Name: setup.ChainCodeID, Path: setup.ChaincodeGoPath, Version: "0", Args: [][]byte{[]byte("init")}, Policy: ccPolicy})
        if err != nil {
                fmt.Println("instantiat error : ", err)
		return errors.WithMessage(err, "instantiate error")
        }
	fmt.Println("instantiate Response : ", resp)
        fmt.Println("instantiate end")

        // Channel client is used to query and execute transactions
        clientContext := setup.sdk.ChannelContext(setup.ChannelID, fabsdk.WithUser(setup.UserName), fabsdk.WithOrg(setup.OrgName))
        setup.client, err = channel.New(clientContext)
        if err != nil {
                return errors.WithMessage(err, "failed to create new channel client")
        }
        fmt.Println("Channel client created")

        // Creation of the client which will enables access to our channel events
        setup.event, err = event.New(clientContext)
        if err != nil {
                return errors.WithMessage(err, "failed to create new event client")
        }
        fmt.Println("Event client created")

        fmt.Println("Chaincode Installation & Instantiation Successful")

        return nil
}

func (setup *FabricSetup) CloseSDK() {
        setup.sdk.Close()
}
