package chaincode_test

import (
	"car_transfer/chaincode"
	"testing"

	"github.com/hyperledger/fabric-chaincode-go/shim"
	"github.com/hyperledger/fabric-chaincode-go/shimtest"
	pb "github.com/hyperledger/fabric-protos-go/peer"
	"github.com/stretchr/testify/assert"
)

const (
	alice = `{"Id":"1", "Name":"Alice"}`
	bob   = `{"Id":"2", "Name":"Bob"}`

	emptyOwners = "[]"
	oneOwners   = "[" + alice + "]"
	twoOwners   = "[" + alice + "," + bob + "]"

	timestamp = `"2018-01-01T12:34:56Z"`

	car1  = `{"Id":"1", "Name":"E-AE86", "OwnerId":"1", "Timestamp":` + timestamp + `}`
	car1b = `{"Id":"1", "Name":"E-AE86", "OwnerId":"2", "Timestamp":` + timestamp + `}`
	car2  = `{"Id":"2", "Name":"GF-FD3S", "OwnerId":"1", "Timestamp":` + timestamp + `}`

	oneCars = "[" + car1 + "]"
	twoCars = "[" + car1 + "," + car2 + "]"

	one = `"1"`
	two = `"2"`
)

//
// test utilities
//

// custom assertion that checks if the response is OK.
func responseOK(res pb.Response) func() bool {
	return func() bool { return res.Status < shim.ERRORTHRESHOLD }
}

// custom assertion that checks if the response is FAIL.
func responseFail(res pb.Response) func() bool {
	return func() bool { return res.Status >= shim.ERRORTHRESHOLD }
}

// converts function name and arguments into the format that MockStub accepts.
// This function was copied and slightly modified from mockstub.go.
func getBytes(function string, args ...string) [][]byte {
	bytes := make([][]byte, 0, len(args)+1)
	bytes = append(bytes, []byte(function))
	for _, s := range args {
		bytes = append(bytes, []byte(s))
	}
	return bytes
}

//
// Testcases
//

func getStub() *shimtest.MockStub {
	cc := &chaincode.CarTransferCC{}
	stub := shimtest.NewMockStub("cartransfer", cc)
	return stub
}

// OK1: normal Init()
func TestInit_OK1(t *testing.T) {
	stub := getStub()
	if assert.NotNil(t, stub) {
		res := stub.MockInit("5", nil)
		assert.Condition(t, responseOK(res))
	}
}

// NG1: unknown method Invoke()
func TestInvoke_NG1(t *testing.T) {
	stub := getStub()
	if assert.NotNil(t, stub) &&
		assert.Condition(t, responseOK(stub.MockInit("5", nil))) {
		res := stub.MockInvoke("5", getBytes("BadMethod"))
		assert.Condition(t, responseFail(res))
	}
}

// OK1: success
func TestAddOwner_OK(t *testing.T) {
	stub := getStub()
	if assert.NotNil(t, stub) &&
		assert.Condition(t, responseOK(stub.MockInit("5", nil))) {
		res := stub.MockInvoke("5", getBytes("ListOwners", one))
		assert.Condition(t, responseOK(res))
		assert.JSONEq(t, emptyOwners, string(res.Payload))

		res = stub.MockInvoke("5", getBytes("AddOwner", alice))
		assert.Condition(t, responseOK(res))

		res = stub.MockInvoke("5", getBytes("ListOwners", one))
		assert.Condition(t, responseOK(res))
		assert.JSONEq(t, oneOwners, string(res.Payload))
	}
}

// NG1: less arguments
func TestAddOwner_NG1(t *testing.T) {
	stub := getStub()
	if assert.NotNil(t, stub) &&
		assert.Condition(t, responseOK(stub.MockInit("5", nil))) {
		res := stub.MockInvoke("5", getBytes("AddOwner"))
		assert.Condition(t, responseFail(res))
	}
}

// NG2: illegal JSON argument
func TestAddOwner_NG2(t *testing.T) {
	stub := getStub()
	if assert.NotNil(t, stub) &&
		assert.Condition(t, responseOK(stub.MockInit("5", nil))) {
		res := stub.MockInvoke("5", getBytes("AddOwner", "bad"))
		assert.Condition(t, responseFail(res))
	}
}

// OK1: 1 Owner
func TestListOwners_OK1(t *testing.T) {
	stub := getStub()
	if assert.NotNil(t, stub) &&
		assert.Condition(t, responseOK(stub.MockInit("5", nil))) {
		res := stub.MockInvoke("5", getBytes("AddOwner", alice))
		assert.Condition(t, responseOK(res))

		res = stub.MockInvoke("5", getBytes("ListOwners"))
		assert.Condition(t, responseOK(res))
		t.Logf("%s", res.Payload)
		assert.JSONEq(t, oneOwners, string(res.Payload))
	}
}

// OK2: 2 Owners
func TestListOwners_OK2(t *testing.T) {
	stub := getStub()
	if assert.NotNil(t, stub) &&
		assert.Condition(t, responseOK(stub.MockInit("5", nil))) {
		res := stub.MockInvoke("5", getBytes("AddOwner", alice))
		assert.Condition(t, responseOK(res))

		res = stub.MockInvoke("5", getBytes("AddOwner", bob))
		assert.Condition(t, responseOK(res))

		res = stub.MockInvoke("5", getBytes("ListOwners"))
		assert.Condition(t, responseOK(res))
		assert.JSONEq(t, twoOwners, string(res.Payload))
	}
}

// OK1: a single Car
func TestAddCar_OK1(t *testing.T) {
	stub := getStub()
	if assert.NotNil(t, stub) &&
		assert.Condition(t, responseOK(stub.MockInit("5", nil))) {
		res := stub.MockInvoke("5", getBytes("GetCar", one))
		assert.Condition(t, responseFail(res))

		res = stub.MockInvoke("5", getBytes("AddOwner", alice))
		assert.Condition(t, responseOK(res))
		res = stub.MockInvoke("5", getBytes("AddCar", car1))
		assert.Condition(t, responseOK(res))

		res = stub.MockInvoke("5", getBytes("ListCars", one))
		if assert.Condition(t, responseOK(res)) {
			assert.JSONEq(t, oneCars, string(res.Payload))
		}
	}
}

// OK2: two Cars
func TestListCars_OK2(t *testing.T) {
	stub := getStub()
	if assert.NotNil(t, stub) &&
		assert.Condition(t, responseOK(stub.MockInit("5", nil))) {
		res := stub.MockInvoke("5", getBytes("AddOwner", alice))
		assert.Condition(t, responseOK(res))

		res = stub.MockInvoke("5", getBytes("AddCar", car1))
		assert.Condition(t, responseOK(res))

		res = stub.MockInvoke("5", getBytes("AddCar", car2))
		assert.Condition(t, responseOK(res))

		res = stub.MockInvoke("5", getBytes("ListCars"))
		assert.Condition(t, responseOK(res))
		assert.JSONEq(t, twoCars, string(res.Payload))
	}
}

// OK1: change owner from Alice to Bob
func TestUpdateCar_OK1(t *testing.T) {
	stub := getStub()
	if assert.NotNil(t, stub) &&
		assert.Condition(t, responseOK(stub.MockInit("5", nil))) {
		res := stub.MockInvoke("5", getBytes("AddOwner", alice))
		assert.Condition(t, responseOK(res))
		res = stub.MockInvoke("5", getBytes("AddOwner", bob))
		assert.Condition(t, responseOK(res))
		res = stub.MockInvoke("5", getBytes("AddCar", car1))
		assert.Condition(t, responseOK(res))
		res = stub.MockInvoke("5", getBytes("UpdateCar", car1b))
		assert.Condition(t, responseOK(res))

		res = stub.MockInvoke("5", getBytes("GetCar", one))
		if assert.Condition(t, responseOK(res)) {
			assert.JSONEq(t, car1b, string(res.Payload))
		}
	}
}

// NG1: specified car does not exist
func TestUpdateCar_NG1(t *testing.T) {
	stub := getStub()
	if assert.NotNil(t, stub) &&
		assert.Condition(t, responseOK(stub.MockInit("5", nil))) {
		res := stub.MockInvoke("5", getBytes("UpdateCar", car1b))
		assert.Condition(t, responseFail(res))
	}
}

// OK2: transfer from Alice to Bob
func TestTransferCar_OK1(t *testing.T) {
	stub := getStub()
	if assert.NotNil(t, stub) &&
		assert.Condition(t, responseOK(stub.MockInit("5", nil))) {
		res := stub.MockInvoke("5", getBytes("AddOwner", alice))
		assert.Condition(t, responseOK(res))
		res = stub.MockInvoke("5", getBytes("AddOwner", bob))
		assert.Condition(t, responseOK(res))

		res = stub.MockInvoke("5", getBytes("AddCar", car1))
		assert.Condition(t, responseOK(res))
		res = stub.MockInvoke("5", getBytes("TransferCar", one, two))
		assert.Condition(t, responseOK(res))

		res = stub.MockInvoke("5", getBytes("GetCar", one))
		if assert.Condition(t, responseOK(res)) {
			assert.JSONEq(t, car1b, string(res.Payload))
		}
	}
}

// NG1: specified Car does not exist
func TestTransferCar_NG1(t *testing.T) {
	stub := getStub()
	if assert.NotNil(t, stub) &&
		assert.Condition(t, responseOK(stub.MockInit("5", nil))) {
		res := stub.MockInvoke("5", getBytes("AddOwner", alice))
		res = stub.MockInvoke("5", getBytes("TransferCar", one, two))
		assert.Condition(t, responseFail(res))
	}
}

// NG2: new Owner not found
func TestTransferCar_NG2(t *testing.T) {
	stub := getStub()
	if assert.NotNil(t, stub) &&
		assert.Condition(t, responseOK(stub.MockInit("5", nil))) {
		res := stub.MockInvoke("5", getBytes("AddOwner", alice))
		res = stub.MockInvoke("5", getBytes("AddCar", car1))

		res = stub.MockInvoke("5", getBytes("TransferCar", one, two))
		assert.Condition(t, responseFail(res))
	}
}

// NG3: less arguments
func TestTransferCar_NG3(t *testing.T) {
	stub := getStub()
	if assert.NotNil(t, stub) &&
		assert.Condition(t, responseOK(stub.MockInit("5", nil))) {
		res := stub.MockInvoke("5", getBytes("AddOwner", alice))
		res = stub.MockInvoke("5", getBytes("TransferCar", one))
		assert.Condition(t, responseFail(res))
	}
}
