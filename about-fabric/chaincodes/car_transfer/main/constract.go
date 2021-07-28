package main

import (
	cartransfer "car_transfer"
	"car_transfer/chaincode"
	"fmt"

	"github.com/hyperledger/fabric-chaincode-go/shim"
)

func main() {
	var _ cartransfer.CarTransfer = (*chaincode.CarTransferCC)(nil)

	err := shim.Start(&chaincode.CarTransferCC{})
	if err != nil {
		fmt.Printf("Error in chaincode process: %s", err)
	}
	fmt.Println("start!")
}
