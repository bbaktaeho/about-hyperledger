package cartransfer

import (
	"time"

	"github.com/hyperledger/fabric-chaincode-go/shim"
)

type Owner struct {
	Id   string
	Name string
}

type Car struct {
	Id        string
	Name      string
	OwnerId   string
	Timestamp time.Time
}

type CarTransfer interface {
	AddOwner(shim.ChaincodeStubInterface, *Owner) error
	CheckOwner(shim.ChaincodeStubInterface, string) (bool, error)
	ListOwners(shim.ChaincodeStubInterface) ([]*Owner, error)

	AddCar(shim.ChaincodeStubInterface, *Car) error
	CheckCar(shim.ChaincodeStubInterface, string) (bool, error)
	ValidateCar(shim.ChaincodeStubInterface, *Car) (bool, error)
	GetCar(shim.ChaincodeStubInterface, string) (*Car, error)
	UpdateCar(shim.ChaincodeStubInterface, *Car) error
	ListCars(shim.ChaincodeStubInterface) ([]*Car, error)
	ListOwnerIdCars(shim.ChaincodeStubInterface, string) ([]*Car, error)

	TransferCar(stub shim.ChaincodeStubInterface, carId string, newOwnerId string) error
}
