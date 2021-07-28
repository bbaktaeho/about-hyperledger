package chaincode

import (
	cartransfer "car_transfer"
	"car_transfer/utils"
	"encoding/json"
	"errors"
	"fmt"
	"strings"

	"github.com/hyperledger/fabric-chaincode-go/shim"
	pb "github.com/hyperledger/fabric-protos-go/peer"
	"github.com/jinzhu/inflection"
)

//
// Chaincode interface implementation
//
type CarTransferCC struct{}

func (c *CarTransferCC) Init(stub shim.ChaincodeStubInterface) pb.Response {
	utils.Log("Init", 1, "init")
	return shim.Success([]byte{})
}

func (c *CarTransferCC) Invoke(stub shim.ChaincodeStubInterface) pb.Response {
	//sample of API use: show tX timestamp
	timestamp, err := stub.GetTxTimestamp()
	if err != nil {
		return shim.Error(fmt.Sprintf("failed to get TX timestamp: %s", err))
	}

	var (
		fcn  string
		args []string
	)
	fcn, args = stub.GetFunctionAndParameters()
	utils.Log("Invoke", 1, fmt.Sprintf("timestamp = %s", timestamp))
	utils.Log("Invoke", 1, fmt.Sprintf("function name = %s, args = %s", fcn, args))

	switch fcn {
	// adds a new Owner
	case "AddOwner":
		// checks arguments length
		if err := utils.CheckLen(1, args); err != nil {
			return shim.Error(err.Error())
		}

		// unmarshal
		owner := &cartransfer.Owner{}
		err := json.Unmarshal([]byte(args[0]), owner)
		if err != nil {
			mes := fmt.Sprintf("failed to unmarshal Owner JSON: %s", err.Error())
			return shim.Error(mes)
		}

		err = c.AddOwner(stub, owner)
		if err != nil {
			return shim.Error(err.Error())
		}

		// returns a success value
		return shim.Success([]byte{})

		// lists Owners
	case "ListOwners":
		owners, err := c.ListOwners(stub)
		if err != nil {
			return shim.Error(err.Error())
		}
		// marshal
		b, err := json.Marshal(owners)
		if err != nil {
			mes := fmt.Sprintf("failed to marshal Owners: %s", err.Error())
			return shim.Error(mes)
		}
		// returns a success value
		return shim.Success(b)

	// adds a new Car
	case "AddCar":
		// checks arguments length
		if err := utils.CheckLen(1, args); err != nil {
			return shim.Error(err.Error())
		}

		// unmarshal
		car := &cartransfer.Car{}
		err := json.Unmarshal([]byte(args[0]), car)
		if err != nil {
			mes := fmt.Sprintf("failed to unmarshal Car JSON: %s", err.Error())
			return shim.Error(mes)
		}

		err = c.AddCar(stub, car)
		if err != nil {
			return shim.Error(err.Error())
		}

		// returns a success value
		return shim.Success([]byte{})

		// lists Cars
	case "ListCars":
		cars, err := c.ListCars(stub)
		if err != nil {
			return shim.Error(err.Error())
		}

		// marshal
		b, err := json.Marshal(cars)
		if err != nil {
			mes := fmt.Sprintf("failed to marshal Cars: %s", err.Error())
			return shim.Error(mes)
		}

		// returns a success value
		return shim.Success(b)

	case "ListOwnerIdCars":
		// unmarshal
		var owner string
		err := json.Unmarshal([]byte(args[0]), &owner)
		if err != nil {
			mes := fmt.Sprintf("failed to unmarshal the 1st argument: %s", err.Error())
			return shim.Error(mes)
		}

		cars, err := c.ListOwnerIdCars(stub, owner)
		if err != nil {
			return shim.Error(err.Error())
		}

		// marshal
		b, err := json.Marshal(cars)
		if err != nil {
			mes := fmt.Sprintf("failed to marshal Cars: %s", err.Error())
			return shim.Error(mes)
		}

		// returns a success value
		return shim.Success(b)

		// gets an existing Car
	case "GetCar":
		// checks arguments length
		if err := utils.CheckLen(1, args); err != nil {
			return shim.Error(err.Error())
		}

		// unmarshal
		var id string
		err := json.Unmarshal([]byte(args[0]), &id)
		if err != nil {
			mes := fmt.Sprintf("failed to unmarshal the 1st argument: %s", err.Error())
			return shim.Error(mes)
		}

		car, err := c.GetCar(stub, id)
		if err != nil {
			return shim.Error(err.Error())
		}

		// marshal
		b, err := json.Marshal(car)
		if err != nil {
			mes := fmt.Sprintf("failed to marshal Car: %s", err.Error())
			return shim.Error(mes)
		}

		// returns a success value
		return shim.Success(b)

		// updates an existing Car
	case "UpdateCar":
		// checks arguments length
		if err := utils.CheckLen(1, args); err != nil {
			return shim.Error(err.Error())
		}

		// unmarshal
		car := new(cartransfer.Car)
		err := json.Unmarshal([]byte(args[0]), car)
		if err != nil {
			mes := fmt.Sprintf("failed to unmarshal Car JSON: %s", err.Error())
			return shim.Error(mes)
		}

		err = c.UpdateCar(stub, car)
		if err != nil {
			return shim.Error(err.Error())
		}

		// returns a success value
		return shim.Success([]byte{})

		// transfers an existing Car to an existing Owner
	case "TransferCar":
		// checks arguments length
		if err := utils.CheckLen(2, args); err != nil {
			return shim.Error(err.Error())
		}

		// unmarshal
		var carId, newOwnerId string
		err := json.Unmarshal([]byte(args[0]), &carId)
		if err != nil {
			mes := fmt.Sprintf(
				"failed to unmarshal the 1st argument: %s",
				err.Error(),
			)
			return shim.Error(mes)
		}

		err = json.Unmarshal([]byte(args[1]), &newOwnerId)
		if err != nil {
			mes := fmt.Sprintf(
				"failed to unmarshal the 2nd argument: %s",
				err.Error(),
			)
			return shim.Error(mes)
		}

		err = c.TransferCar(stub, carId, newOwnerId)
		if err != nil {
			return shim.Error(err.Error())
		}

		// returns a success valuee
		return shim.Success([]byte{})
	}

	// if the function name is unknown
	mes := fmt.Sprintf("Unknown method: %s", fcn)
	return shim.Error(mes)
}

//
// methos implementing CarTransfer interface
//

// Adds a new Owner
func (c *CarTransferCC) AddOwner(stub shim.ChaincodeStubInterface, owner *cartransfer.Owner) error {

	// checks if the specified Owner exists
	found, err := c.CheckOwner(stub, owner.Id)
	if err != nil {
		return err
	}
	if found {
		mes := fmt.Sprintf("an Owner with Id = %s alerady exists", owner.Id)
		return errors.New(mes)
	}

	// converts to JSON
	b, err := json.Marshal(owner)
	if err != nil {
		return err
	}

	// creates a composite key
	key, err := stub.CreateCompositeKey("Owner", []string{owner.Id})
	if err != nil {
		return err
	}

	// stores to the State DB
	err = stub.PutState(key, b)
	if err != nil {
		return err
	}

	// returns successfully
	return nil
}

// Checks existence of the specified Owner
func (c *CarTransferCC) CheckOwner(stub shim.ChaincodeStubInterface, id string) (bool, error) {
	// creates a composite key
	key, err := stub.CreateCompositeKey("Owner", []string{id})
	if err != nil {
		return false, err
	}

	// loads from the State DB
	jsonBytes, err := stub.GetState(key)
	if err != nil {
		return false, err
	}

	// returns successfully
	return jsonBytes != nil, nil
}

// Lists Owners
func (c *CarTransferCC) ListOwners(stub shim.ChaincodeStubInterface) ([]*cartransfer.Owner, error) {
	// executes a range query, which returns an iterator
	iter, err := stub.GetStateByPartialCompositeKey("Owner", []string{})
	if err != nil {
		return nil, err
	}

	// will close the iterator when returned from c method
	defer iter.Close()
	owners := []*cartransfer.Owner{}

	// loops over the iterator
	for iter.HasNext() {
		kv, err := iter.Next()
		if err != nil {
			return nil, err
		}
		owner := new(cartransfer.Owner)
		err = json.Unmarshal(kv.Value, owner)
		if err != nil {
			return nil, err
		}
		owners = append(owners, owner)
	}

	var owners_string string
	for _, owner := range owners {
		b, _ := json.Marshal(owner)
		owners_string += string(b)
	}

	utils.Log("ListOwners", 1, owners_string)
	// returns successfully
	if len(owners) > 1 {
		utils.Log("ListOwners", 1, fmt.Sprintf("%d %s found", len(owners), inflection.Plural("Owner")))
	} else {
		utils.Log("ListOwners", 1, fmt.Sprintf("%d %s found", len(owners), "Owner"))
	}
	return owners, nil
}

// Adds a new Car
func (c *CarTransferCC) AddCar(stub shim.ChaincodeStubInterface, car *cartransfer.Car) error {
	// creates a composite key
	key, err := stub.CreateCompositeKey("Car", []string{car.Id})
	if err != nil {
		return err
	}

	// checks if the specified Car exists
	found, err := c.CheckCar(stub, car.Id)
	if err != nil {
		return err
	}
	if found {
		mes := fmt.Sprintf("Car with Id = %s already exists", car.Id)
		return errors.New(mes)
	}

	// validates the Car
	ok, err := c.ValidateCar(stub, car)
	if err != nil {
		return err
	}
	if !ok {
		mes := "Validation of the Car failed"
		return errors.New(mes)
	}

	// converts to JSON
	b, err := json.Marshal(car)
	if err != nil {
		return err
	}

	// stores to the State DB
	err = stub.PutState(key, b)
	if err != nil {
		return err
	}

	// returns successfully
	return nil
}

// Checks existence of the specified Car
func (c *CarTransferCC) CheckCar(stub shim.ChaincodeStubInterface, id string) (bool, error) {

	// creates a composite key
	key, err := stub.CreateCompositeKey("Car", []string{id})
	if err != nil {
		return false, err
	}

	// loads from the State DB
	jsonBytes, err := stub.GetState(key)
	if err != nil {
		return false, err
	}

	// returns successfully
	return jsonBytes != nil, nil
}

// Validates the content of the specified Car
func (c *CarTransferCC) ValidateCar(stub shim.ChaincodeStubInterface, car *cartransfer.Car) (bool, error) {
	// checks existence of the Owner with the OwnerId
	found, err := c.CheckOwner(stub, car.OwnerId)
	if err != nil {
		return false, err
	}

	// returns successfully
	return found, nil
}

// Gets the specified Car
func (c *CarTransferCC) GetCar(stub shim.ChaincodeStubInterface, id string) (*cartransfer.Car, error) {
	// creates a composite key
	key, err := stub.CreateCompositeKey("Car", []string{id})
	if err != nil {
		return nil, err
	}

	// loads from the state DB
	jsonBytes, err := stub.GetState(key)
	if err != nil {
		return nil, err
	}
	if jsonBytes == nil {
		mes := fmt.Sprintf("Car with Id = %s was not found", id)
		return nil, errors.New(mes)
	}

	// unmarshal
	car := new(cartransfer.Car)
	err = json.Unmarshal(jsonBytes, car)
	if err != nil {
		return nil, err
	}

	// returns successfully
	return car, nil
}

// Updates the content of the specified Car
func (c *CarTransferCC) UpdateCar(stub shim.ChaincodeStubInterface, car *cartransfer.Car) error {
	// checks existence of the specified Car
	found, err := c.CheckCar(stub, car.Id)
	if err != nil {
		return err
	}
	if !found {
		mes := fmt.Sprintf("Car with Id = %s does not exist", car.Id)
		return errors.New(mes)
	}

	// validates the Car
	ok, err := c.ValidateCar(stub, car)
	if err != nil {
		return err
	}
	if !ok {
		mes := "Validation of the Car failed"
		return errors.New(mes)
	}

	// creates a composite key
	key, err := stub.CreateCompositeKey("Car", []string{car.Id})
	if err != nil {
		return err
	}

	// converts to JSON
	b, err := json.Marshal(car)
	if err != nil {
		return err
	}

	// stores to the State DB
	err = stub.PutState(key, b)
	if err != nil {
		return err
	}

	// returns successfully
	return nil
}

// Lists Cars
func (c *CarTransferCC) ListCars(stub shim.ChaincodeStubInterface) ([]*cartransfer.Car, error) {
	// executes a range query, which returns an iterator
	iter, err := stub.GetStateByPartialCompositeKey("Car", []string{})
	if err != nil {
		return nil, err
	}

	// will close the iterator when returned from c method
	defer iter.Close()

	// loops over the iterator
	cars := []*cartransfer.Car{}
	for iter.HasNext() {
		kv, err := iter.Next()
		if err != nil {
			return nil, err
		}
		car := new(cartransfer.Car)
		err = json.Unmarshal(kv.Value, car)
		if err != nil {
			return nil, err
		}
		cars = append(cars, car)
	}

	// returns successfully
	// returns successfully
	if len(cars) > 1 {
		utils.Log("ListCars", 1, fmt.Sprintf("%d %s found", len(cars), inflection.Plural("Car")))
	} else {
		utils.Log("ListCars", 1, fmt.Sprintf("%d %s found", len(cars), "Car"))
	}
	return cars, nil
}

// Lists OwnerId Cars
func (c *CarTransferCC) ListOwnerIdCars(stub shim.ChaincodeStubInterface, ownerId string) ([]*cartransfer.Car, error) {
	// executes a range query, which returns an iterator
	iter, err := stub.GetStateByPartialCompositeKey("Car", []string{})
	if err != nil {
		return nil, err
	}

	// will close the iterator when returned from c method
	defer iter.Close()

	// loops over the iterator
	cars := []*cartransfer.Car{}
	for iter.HasNext() {
		kv, err := iter.Next()
		if err != nil {
			return nil, err
		}
		car := new(cartransfer.Car)
		err = json.Unmarshal(kv.Value, car)
		if err != nil {
			return nil, err
		}
		if strings.Index(ownerId, "admin") != -1 {
			cars = append(cars, car)
		} else {
			if car.OwnerId == ownerId {
				cars = append(cars, car)
			}

		}

	}

	// returns successfully
	if len(cars) > 1 {
		utils.Log("ListCars", 1, fmt.Sprintf("%d %s found", len(cars), inflection.Plural("Car")))
	} else {
		utils.Log("ListCars", 1, fmt.Sprintf("%d %s found", len(cars), "Car"))
	}
	return cars, nil
}

// Transfers the specified Car to the specified Owner
func (c *CarTransferCC) TransferCar(stub shim.ChaincodeStubInterface, carId string, newOwnerId string) error {

	// gets the specified Car (err returned if it does not exist)
	car, err := c.GetCar(stub, carId)
	if err != nil {
		return err
	}

	// updates OwnerId field
	car.OwnerId = newOwnerId

	// stores the updated Car back to the State DB
	err = c.UpdateCar(stub, car)
	if err != nil {
		return err
	}

	// returns successfully
	return nil
}
