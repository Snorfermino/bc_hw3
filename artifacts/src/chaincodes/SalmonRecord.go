package main


import (
	"encoding/json"
	"fmt"
	"strings"
	"github.com/hyperledger/fabric/core/chaincode/shim"
	pb "github.com/hyperledger/fabric/protos/peer"
)
// ----- Salmon ----- //
type Salmon struct {
	id string `json:"id"`
	vessel 				string        	`json:"vessel"` 
	datetime       		string          `json:"datetime"`
	location 			string			`json:"location"`
	holder 				string 			`json:"holder"`
}


type SmartContract struct {
	
}

func main() {

    err := shim.Start(new(SmartContract))

    if err != nil {

        fmt.Println("Could not start SampleChaincode")

    } else {

        fmt.Println("SampleChaincode successfully started")

    }

}

func (s *SmartContract)Init(APIstub shim.ChaincodeStubInterface) pb.Response {

	
		return shim.Success(nil)
}

func (s *SmartContract)Invoke(APIstub shim.ChaincodeStubInterface) pb.Response {
	function, args := APIstub.GetFunctionAndParameters()
	fmt.Println("invoke is running " + function)
	if function == "recordSalmon" {
		return s.recordSalmon(APIstub, args)
	} else if function == "changeSalmonHolder" {
		return s.changeSalmonHolder(APIstub, args)
	} else if function == "querySalmon" {
		return s.querySalmon(APIstub, args)
	} else if function == "queryAllSalmon" {
		return s.queryAllSalmon(APIstub, args)
	}
	fmt.Println("invoke did not find func: " + function) //error
	return shim.Error("Received unknown function invocation")
}
func (t *SmartContract)recordSalmon(stub shim.ChaincodeStubInterface, args []string)pb.Response {
	var err error

	//   0       1       2           3            4
	// "Id", "vessel", "datetime", "location", "holder"
	if len(args) != 5 {
		return shim.Error("Incorrect number of arguments. Expecting 5")
	}

	// ==== Input sanitation ====
	fmt.Println("- start recording Salmon")
	if len(args[0]) <= 0 {
		return shim.Error("1st argument must be a non-empty string")
	}
	if len(args[1]) <= 0 {
		return shim.Error("2nd argument must be a non-empty string")
	}
	if len(args[2]) <= 0 {
		return shim.Error("3rd argument must be a non-empty string")
	}
	if len(args[3]) <= 0 {
		return shim.Error("4th argument must be a non-empty string")
	}
	if len(args[4]) <= 0 {
		return shim.Error("5th argument must be a non-empty string")
	}
	id := args[0]
	vessel := strings.ToLower(args[1])
	datetime := strings.ToLower(args[2])
	location := strings.ToLower(args[3])
	holder := strings.ToLower(args[4])


	// ==== Check if salmon has been recorded yet ====
	salmonAsBytes, err := stub.GetState(id)
	if err != nil {
		return shim.Error("Failed to get salmon: " + err.Error())
	}else if salmonAsBytes != nil {
		fmt.Println("This marble already exists: " + id)
		return shim.Error("This salmon has been recorded: " + id)
	}

	// ==== Create salmon object and marshal to JSON ====
	// objectType := "Salmon"
	salmon :=  & Salmon {id, vessel, datetime, location, holder}
	salmonJSONasBytes, err := json.Marshal(salmon)
	if err != nil {
		return shim.Error(err.Error())
	}
	// === record salmon to state ===
	err = stub.PutState(id, salmonJSONasBytes)
	if err != nil {
		return shim.Error(err.Error())
	}

	fmt.Println("- end record salmon")
	return shim.Success(nil)
}

func (s * SmartContract)querySalmon(stub shim.ChaincodeStubInterface, args []string)pb.Response {
	if len(args) != 1 {
		return shim.Error("Incorect number of arguments, expecting 1")
	}

	salmonAsBytes, _ := stub.GetState(args[0])
	if salmonAsBytes == nil {
		return shim.Error("Could not find Salmon")
	}

	salmon := Salmon {}
	err := json.Unmarshal(salmonAsBytes,  &salmon)
	if err != nil {
		return shim.Error(err.Error())
	}
	fmt.Printf("I saw %s at %s.\n", salmon.vessel, salmon.holder)
	return shim.Success(salmonAsBytes)
}

func (s * SmartContract)queryAllSalmon(stub shim.ChaincodeStubInterface, args []string)pb.Response {
	if len(args) != 1 {
		return shim.Error("Incorect number of arguments, expecting 1")
	}

	salmonAsBytes, _ := stub.GetState(args[0])
	if salmonAsBytes == nil {
		return shim.Error("Could not find Salmon")
	}
	return shim.Success(salmonAsBytes)
}

func (s * SmartContract)changeSalmonHolder(stub shim.ChaincodeStubInterface, args []string)pb.Response {
	if len(args) != 2 {
		return shim.Error("Incorrect number of arguments. Expecting 2")
	}
	if len(args[0]) <= 0 {
		return shim.Error("1st argument must be a non-empty string")
	}
	if len(args[1]) <= 0 {
		return shim.Error("2nd argument must be a non-empty string")
	}
	salmonId := args[0]
	ownerName := strings.ToLower(args[1])

	salmonAsBytes, err := stub.GetState(salmonId)
	if err != nil {
		return shim.Error("Failed to get salmon: " + err.Error())
	}
	salmon := Salmon {}
	err = json.Unmarshal(salmonAsBytes,  & salmon)
	if err != nil {
		return shim.Error(err.Error())
	}
	salmon.holder = ownerName 

	salmonAsBytes, err = json.Marshal(salmon)
	err = stub.PutState(salmonId, salmonAsBytes)//rewrite the marble
	if err != nil {
		return shim.Error(err.Error())
	}

	return shim.Success(nil)
}