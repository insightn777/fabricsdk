package main
  
import(
        "encoding/json"
        "fmt"
        //"reflect"
        "strconv"
        "bytes"
        "github.com/hyperledger/fabric/core/chaincode/shim"
        pb "github.com/hyperledger/fabric/protos/peer"
)
type SimpleChaincode struct {}

type user struct {
        ObjectType string `json:"DocType"`
        ID string `json:"ID"`
        Name string `json:"Name"`
        Password string `json:"Password"`
        SocialNumber string `json:"SocialNumber"`
        Location string `json:"Location"`
        VoteResult []votevoteresult `json:"VoteResult"`
} // user information

type uservoteresult struct {
        ID string `json:"ID"`
        Result []int `json:"Result"`
} // user's vote result

type votevoteresult struct {
        Votename string `json:"Votename"`
        Result []int `json:"Result"`
}

type vote struct {
        ObjectType string `json:"DocType"`
        Votename string `json:"Votename"`
        StartDate string `json:"StartDate"`
        EndDate string `json:"StartDate"`
        Question []string `json:"Question"`
        UserResult []uservoteresult `json:"UserResult"`
} //have question and its result per user

func main () {
        err := shim.Start(new(SimpleChaincode))
        if err != nil {
                fmt.Printf("Error starting : %s", err)
        }
}
func (t *SimpleChaincode) Init(stub shim.ChaincodeStubInterface) pb.Response {
        fmt.Printf("Initialize\n")
        return shim.Success(nil)
}

func (t *SimpleChaincode) Invoke(stub shim.ChaincodeStubInterface) pb.Response {
        function, args := stub.GetFunctionAndParameters()
        if function == "insertUser" {
                return t.insertUser(stub, args)
        } else if function == "insertVote"{
                return t.insertVote(stub, args)
        } else if function == "insertVoteResult" {
                return t.insertVoteResult(stub, args)
        } else if function == "query" {
                return t.query(stub, args)
        } else if function == "delete" {
                return t.delete(stub, args)
        }
        fmt.Printf("not function");
        return shim.Error("not funtion");
}
func (t *SimpleChaincode) insertUser(stub shim.ChaincodeStubInterface, args []string) pb.Response{
        var err error
        fmt.Printf("insert...\n")
        //inser user information
        userdata := user{}
        userdata.ObjectType = "user"
        userdata.ID = args[0]
        userdata.Name = args[1]
        userdata.Password = args[2]
        userdata.SocialNumber = args[3]
        userdata.Location = args[4]
        //make structure and marshaling
        userJSONasBytes, err := json.Marshal(userdata)
        if err != nil {
                fmt.Printf("error")
                return shim.Error(err.Error())
        }
        fmt.Println(string(userJSONasBytes))
        if err != nil {
                fmt.Printf("insert error 1\n")
                return shim.Error(err.Error())
        }
        //insert DB
        err = stub.PutState(args[0], userJSONasBytes)
        if err != nil {
                fmt.Printf("insert error 2\n")
                return shim.Error(err.Error())
        }
  
        return shim.Success(nil)
}
func (t *SimpleChaincode) insertVote (stub shim.ChaincodeStubInterface, args []string) pb.Response{
        //Admin insert vote question
        var err error
        votename := args[0]
        questionNum := len(args) - 3
        voteAsByte, err := stub.GetState(votename)
        if voteAsByte != nil {
                fmt.Println("vote already exist")
                return shim.Error("vote already exist")
        }
        //inser vote qustion
        votedata := vote{}
        votedata.ObjectType = "vote"
        votedata.Votename = args[0]
        votedata.StartDate = args[1]
        votedata.StartDate = args[2]
        votedata.Question = make([]string, questionNum)
        for i := 0; i < questionNum; i++ {
                votedata.Question[i] = args[i + 3]
        }
        fmt.Println(votedata)
        // marshaling
        voteAsJSONBytes, err := json.Marshal(votedata)
        if err != nil {
                fmt.Println("Marshal error")
                return shim.Error("Marshal error")
        }
        //insert DB
        err = stub.PutState(votename, voteAsJSONBytes)
        if err != nil {
                fmt.Println("DB insert error")
                shim.Error("DB insert error")
        }
        fmt.Println("vote question insert success")
        return shim.Success(nil)
}
func (t *SimpleChaincode) insertVoteResult (stub shim.ChaincodeStubInterface, args []string) pb.Response {
        var err error
        //var temp tempvote
        //query key that check for vote question have existed
        votequeryString := "{\"selector\":{\"DocType\":\"vote\", \"Votename\":\"" + args[0] + "\"}}"
        userqueryString := "{\"selector\":{\"DocType\":\"user\", \"ID\":\"" + args[1] + "\"}}"
        existflag, err := stub.GetState(args[0])
        if existflag == nil {
                fmt.Println("this vote not exist")
                return shim.Error("this vote not exsit")
        }
        // prev query get key-vaule form
        voteresultIterator, err := stub.GetQueryResult(votequeryString)
        if err != nil {
                fmt.Println("get query error")
                return shim.Error("get query error")
        }
        userresultIterator, err := stub.GetQueryResult(userqueryString)
        voteresponse, err := voteresultIterator.Next()
        if err != nil {
                fmt.Println("not next")
                return shim.Error("not next")
        }
        userresponse, err := userresultIterator.Next()
  
        //unmarshialing prev and append new data
        votedata := vote{}
        userdata := user{}
        err = json.Unmarshal(voteresponse.Value, &votedata)
        if err != nil {
                fmt.Println("unmarshal error")
                return shim.Error("unmarshal error")
        }
        err = json.Unmarshal(userresponse.Value, &userdata)
        fmt.Println("prev struct : ", votedata)
        temp1 := uservoteresult{}
        temp2 := votevoteresult{}
        temp1.ID = args[1]
        temp2.Votename = args[0]
        temp1.Result = make([]int, len(args) - 2)
        temp2.Result = make([]int, len(args) - 2)
        for i := 0; i < len(args) - 2; i++ {
                temp1.Result[i], _ = strconv.Atoi(args[i + 2])
                temp2.Result[i], _ = strconv.Atoi(args[i + 2])
        }
        votedata.UserResult = append(votedata.UserResult, temp1)
        userdata.VoteResult = append(userdata.VoteResult, temp2) 
  
        //new vote structure marshalring
        fmt.Println("now struct : ", votedata)
        voteResultAsJSONBytes, err := json.Marshal(votedata)
        if err != nil {
                fmt.Println("marshal error")
                return shim.Error("marshal error")
        }
        userResultAsJSONBytes, err := json.Marshal(userdata)
  
        //insert DB
        err = stub.PutState(votedata.Votename, voteResultAsJSONBytes)
        if err != nil {
                fmt.Println("insert DB error")
                shim.Error("insert DB error")
        }
        err = stub.PutState(userdata.ID, userResultAsJSONBytes)
  
        return shim.Success(nil)
}

func(t *SimpleChaincode) query(stub shim.ChaincodeStubInterface, args []string) pb.Response {
        var err error
        queryString := args[0]
        queryResults, err := getQueryResultForQueryString(stub, queryString)
        if err != nil {
                return shim.Error(err.Error())
        } 
        /*for i := 0; i < len(queryResults); i++ {
                fmt.Println(string(queryResults[i]))
        }for check*/
        return shim.Success(nil)
}
func getQueryResultForQueryString(stub shim.ChaincodeStubInterface, queryString string) ([][]byte, error) {
        var buffer [][]uint8
  
        fmt.Printf("- getQueryResultForQueryString queryString:\n%s\n", queryString)
  
        resultsIterator, err := stub.GetQueryResult(queryString)
        if err != nil {
                fmt.Printf("quert error 1\n")
                return nil, err
        }
  
        defer resultsIterator.Close()
        
        for resultsIterator.HasNext() {
                response, _ := resultsIterator.Next()
                buffer = append(buffer, response.Value)
        }
        //make buffer -> string form
        /*buffer, err := constructQueryResponseFromIterator(resultsIterator)
        if err != nil {
                fmt.Printf("query error2\n")
                return nil, err
        }*/
        //fmt.Printf("- getQueryResultForQueryString queryResult:\n%s\n", buffer.String())
  
        return buffer, nil
}
func constructQueryResponseFromIterator(resultsIterator shim.StateQueryIteratorInterface) (*bytes.Buffer, error) {
        // buffer is a JSON array containing QueryResults
        var buffer bytes.Buffer
        buffer.WriteString("[")
        bArrayMemberAlreadyWritten := false
        for resultsIterator.HasNext() {
                queryResponse, err := resultsIterator.Next()
                if err != nil {
                        return nil, err
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
        return &buffer, nil
}

func (t *SimpleChaincode) delete(stub shim.ChaincodeStubInterface, args []string) pb.Response {
        var err error
        voteName := args[0]
        fmt.Printf("delete vote named : %s\n", voteName)
        // maybe i have to connect with vote database
        existflag, err := stub.GetState(voteName)
        if existflag == nil {
                fmt.Printf("this vote not exist")
                return shim.Error(err.Error())
        }
        if err != nil {
                fmt.Printf("GetState() error")
                return shim.Error(err.Error())
        }
        err = stub.DelState(voteName)
        if err != nil {
          fmt.Printf("DelState() error")
                return shim.Error(err.Error())
        }
        fmt.Printf("deletion successed.\n")
        return shim.Success(nil)
}
