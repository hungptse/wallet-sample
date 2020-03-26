package main

import (
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/hyperledger/fabric-contract-api-go/contractapi"
)

type SmartContract struct {
	contractapi.Contract
}

type Wallet struct {
	Address string `json:"address"`
	Amount int64 `json:"amount"`
	Name string `json:"name"`
}

type QueryResult struct {
	Key    string `json:"Key"`
	Record *Wallet
}


func (s *SmartContract) InitLedger(ctx contractapi.TransactionContextInterface) error  {
	wallets := []Wallet{
		{
			Address: "27461686465BA4FF87377D8F3914F",
			Amount:  0,
			Name:    "Main",
		},
	}

	for i, wallet := range wallets {
		var walletAsBytes, _ = json.Marshal(wallet)
		err := ctx.GetStub().PutState("WALLET" + strconv.Itoa(i), walletAsBytes)

		if err != nil {
			return fmt.Errorf("Failed to put to world state. %s", err.Error())
		}
	}
	return nil
}

func (s *SmartContract) CreateWallet(ctx contractapi.TransactionContextInterface, walletKey string,address string, amount int64,name string) error {
	 wallet := Wallet{
		Address: address,
		Amount:  amount,
		Name:   name,
	}

	var walletAsBytes, _  = json.Marshal(wallet)

	return ctx.GetStub().PutState(walletKey, walletAsBytes)
}

func (s *SmartContract) QueryWallet(ctx contractapi.TransactionContextInterface, walletKey string) (*Wallet, error)  {
	walletAsBytes, err := ctx.GetStub().GetState(walletKey)

	if err != nil {
		return nil,  fmt.Errorf("Failed to read from world state. %s", err.Error())
	}

	if walletAsBytes == nil{
		return nil, fmt.Errorf("%s does not exist", walletKey)
	}
	var wallet = new(Wallet)
	_ = json.Unmarshal(walletAsBytes,wallet)
	return wallet, nil
}

func (s *SmartContract) QueryAllWallet(ctx contractapi.TransactionContextInterface) ([]QueryResult, error)  {
	resultsIterator, err := ctx.GetStub().GetStateByRange("WALLET1","WALLET100")

	if err != nil {
		return nil, err
	}
	defer resultsIterator.Close()

	results := []QueryResult{}

	for resultsIterator.HasNext() {
		queryRespone, err := resultsIterator.Next()

		if err != nil {
			return nil, err
		}

		wallet := new(Wallet)
		_ = json.Unmarshal(queryRespone.Value, wallet)

		queryResult := QueryResult{
			Key:    queryRespone.Key,
			Record: wallet,
		}
		results = append(results, queryResult)
	}
	return results, nil
}

//func (s *SmartContract) TransferAmount(ctx contractapi.TransactionContextInterface, fromWallet string, toWallet string,transactionType string, amount string) error  {
//
//}


func main() {
	chaincode, err := contractapi.NewChaincode(new(SmartContract))

	if err != nil {
		fmt.Printf("Error create [WALLET] chaincode: %s", err.Error())
	}

	if err := chaincode.Start(); err != nil {
		fmt.Printf("Error create [WALLET] chaincode: %s", err.Error())
	}
}
