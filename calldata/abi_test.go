package calldata

import (
	"bytes"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"math/big"
	"os"
	"strings"
	"testing"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/stretchr/testify/assert"
)

type Cluster struct {
	ValidatorCount  *uint32  // uint32
	NetworkFeeIndex *uint64  // uint64
	Index           *uint64  // uint64
	Active          *bool    // bool
	Balance         *big.Int // uint256
}

type BulkRegisterValidatorInput struct {
	PublicKeys  [][]byte // bytes[]
	OperatorIds []uint64 // uint64[]
	SharesData  [][]byte // bytes[]
	Amount      big.Int  // uint256
	//Cluster     *Cluster // tuple
}

func Test_Temp(t *testing.T) {
	d, err := os.ReadFile("../abis/ssv.json")
	assert.Nil(t, err)

	parsedABI, err := abi.JSON(bytes.NewReader(d))
	assert.Nil(t, err)

	txDataHex := txData
	txDataHex = strings.TrimPrefix(txDataHex, "0x")
	txData, err := hex.DecodeString(txDataHex)
	assert.Nil(t, err)

	// bulkRegisterValidator
	method, err := parsedABI.MethodById(txData[:4])
	assert.Nil(t, err)

	// Remove the function selector (first 4 bytes of txData)
	inputData := txData[4:]

	// Prepare variables to hold the decoded input parameters
	outputs := []interface{}{
		new([][]byte),
		new([]uint64),
		new([][]byte),
		new(big.Int),
		new(big.Int),
		//new(Cluster),
	}

	// Unpack the input data into the variables
	outputs, err = method.Inputs.Unpack(inputData)
	assert.Nil(t, err)

	for _, pubkey := range outputs[0].([][]byte) {
		fmt.Println(hex.EncodeToString(pubkey))
	}

	for _, operatorId := range outputs[1].([]uint64) {
		fmt.Println(operatorId)
	}

	fmt.Println(outputs[3].(*big.Int).String())

	// TODO: asserts
}

func Test_DecodeTX(t *testing.T) {
	// TODO: Find right place to store json file.
	// or. what about loading abi from etherscan directly?
	d, err := os.ReadFile("../abis/ssv.json")
	assert.Nil(t, err)

	abiFunctions := []ABIFunction{}

	err = json.Unmarshal(d, &abiFunctions)
	assert.Nil(t, err)

	funcBySignature := map[string]ABIFunction{}

	for _, function := range abiFunctions {
		funcSignature, err := function.Signature()
		assert.Nil(t, err)

		funcBySignature[funcSignature] = function
	}

	// Captured from
	// https://etherscan.io/address/0xDD9BC35aE942eF0cFa76930954a156B3fF30a4E1#writeProxyContract
	assert.Equal(t, funcBySignature["22f18bf5"].Name, "bulkRegisterValidator")
	assert.Equal(t, funcBySignature["79ba5097"].Name, "acceptOwnership")

	// TODO: Get TX data

	spec, ok := funcBySignature["22f18bf5"]
	_ = ok

	callData := Calldata{}
	callData.LoadString(txData)

	assert.Equal(t, callData.FunctionSignature(), "22f18bf5")

	err = callData.FunctionData(&spec)
}
