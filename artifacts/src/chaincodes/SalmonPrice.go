package main


import (
	"encoding/json"
	"fmt"
	"strings"
	"github.com/hyperledger/fabric/core/chaincode/shim"
	pb "github.com/hyperledger/fabric/protos/peer"
)
type SalmonPrice struct {
	id 					string 			`json:"id"`
	seller 				string        	`json:"seller"` //field for couchdb
	buyer       		string          `json:"buyer"`      //the fieldtags are needed to keep case from bouncing around
	price 				string 			`json:"price"`
}
type SmartContract struct {
	
}

func main() {

	err := shim.Start(new(SmartContract))

	if err != nil {

		fmt.Println("Could not start SmartContract")

	} else {

		fmt.Println("SmartContract successfully started")

	}
}

func (s *SmartContract)Init(APIstub shim.ChaincodeStubInterface)pb.Response {
	return shim.Success(nil)
}

func (s *SmartContract)Invoke(APIstub shim.ChaincodeStubInterface)pb.Response {
	function, args := APIstub.GetFunctionAndParameters()
	fmt.Println("invoke is running " + function)
	if function == "price" {
		return s.priceAgreement(APIstub, args)
	}else if function == "delete" {
		return s.deletePriceAgreement(APIstub, args)
	}
	fmt.Println("invoke did not find func: " + function)//error
	return shim.Error("Received unknown function invocation")
}

func (t *SmartContract)priceAgreement(stub shim.ChaincodeStubInterface, args []string)pb.Response {
	var err error

	//   0          1          2           3
	// "id", "sellerName", "buyerName", "price"
	if len(args) != 4 {
		return shim.Error("Incorrect number of arguments. Expecting 4")
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
		return shim.Error("3rd argument must be a non-empty string")
	}
	id := args[0]
	sellerName := strings.ToLower(args[1])
	buyerName := strings.ToLower(args[2])
	price := strings.ToLower(args[3])


	// ==== Check if price has been recorded yet ====
	priceAsByte, err := stub.GetState(id)
	if err != nil {
		return shim.Error("Failed to get price: " + err.Error())
	}else if priceAsByte != nil {
		fmt.Println("This salmon price agreement already exists: " + id)
		return shim.Error("This salmon price agreement has been recorded: " + id)
	}

	// ==== Create SalmonPrice object and marshal to JSON ====

	salmonPrice :=  & SalmonPrice {id, sellerName, buyerName, price}
	salmonPriceJSONasBytes, err := json.Marshal(salmonPrice)
	if err != nil {
		return shim.Error(err.Error())
	}
	// === record SalmonPrice to state ===
	err = stub.PutState(id, salmonPriceJSONasBytes)
	if err != nil {
		return shim.Error(err.Error())
	}

	fmt.Println("- end record SalmonPrice")
	return shim.Success(nil)
}

func (s *SmartContract)deletePriceAgreement(stub shim.ChaincodeStubInterface, args []string)pb.Response {
	if len(args) != 1 {
		return shim.Error("Incorect number of arguments, expecting 1")
	}
	var jsonResp string
	var salmonPriceJSON SalmonPrice
	salmonPriceAsBytes, err := stub.GetState(args[0])
	if err != nil {
		jsonResp = "{\"Error\":\"Failed to get state at " + args[0] + "\"}"
		return shim.Error(jsonResp)
	}else if salmonPriceAsBytes == nil {
		jsonResp = "{\"Error\":\"Price Agreement does not exist at:" + args[0] + "\"}"
		return shim.Error(jsonResp)
	}

	err = json.Unmarshal([]byte(salmonPriceAsBytes),  & salmonPriceJSON)
	if err != nil {
		jsonResp = "{\"Error\":\"Failed to decode JSON of:" + args[0] + "\"}"
		return shim.Error(jsonResp)
	}

	err = stub.DelState(args[0])//remove the marble from chaincode state
	if err != nil {
		return shim.Error("Failed to delete state:" + err.Error())
	}
	return shim.Success(nil)
}

func (s *SmartContract)querySalmonPriceAgreement(stub shim.ChaincodeStubInterface, args []string)pb.Response {
	if len(args) != 1 {
		return shim.Error("Incorect number of arguments, expecting 1")
	}

	salmonPriceAsBytes, _ := stub.GetState(args[0])
	if salmonPriceAsBytes == nil {
		return shim.Error("Could not find Salmon Price Agreement")
	}

	salmonPrice := SalmonPrice {}
	err := json.Unmarshal(salmonPriceAsBytes,  & salmonPrice)
	if err != nil {
		return shim.Error(err.Error())
	}
	fmt.Printf("Found price agreement between %s and %s at %s.\n", salmonPrice.seller, salmonPrice.buyer, args[0])
	return shim.Success(salmonPriceAsBytes)
}

