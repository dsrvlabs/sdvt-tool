package calldata

import (
	"encoding/hex"
	"fmt"
	"math/big"
	"strings"

	"golang.org/x/crypto/sha3"
)

type ABIFunction struct {
	Constant        bool          `json:"constant"`
	Inputs          []ABIArgument `json:"inputs"`
	Name            string        `json:"name"`
	Outputs         []ABIArgument `json:"outputs"`
	Payable         bool          `json:"payable"`
	StateMutability string        `json:"stateMutability"`
	Type            string        `json:"type"`
}

func (a ABIFunction) Signature() (string, error) {
	argStr := ""
	for _, arg := range a.Inputs {
		if arg.Type == "tuple" {
			argStr += "("
			for _, tupleArg := range arg.Components {
				argStr += tupleArg.Type + ","
			}
			argStr = argStr[:len(argStr)-1]
			argStr += "),"
		} else {
			argStr += arg.Type + ","
		}
	}

	if len(argStr) > 0 {
		argStr = argStr[:len(argStr)-1]
	}

	funcSpec := fmt.Sprintf("%s(%s)", a.Name, argStr)

	sha := sha3.NewLegacyKeccak256()
	_, err := sha.Write([]byte(funcSpec))
	if err != nil {
		return "", err
	}

	funcSignature := hex.EncodeToString(sha.Sum(nil))[:8]

	return funcSignature, nil
}

type ABIArgument struct {
	InternalType string        `json:"internalType"`
	Name         string        `json:"name"`
	Type         string        `json:"type"`
	Components   []ABIArgument `json:"components"`
}

type Calldata struct {
	rawData string
}

func (c *Calldata) LoadString(data string) {
	if strings.HasPrefix(data, "0x") {
		c.rawData = data[2:]
	} else {
		c.rawData = data
	}
}

func (c *Calldata) FunctionSignature() string {
	return c.rawData[:8]
}

func (c *Calldata) FunctionData(spec *ABIFunction) error {
	//header := c.rawData[8:64+8]
	//fmt.Println(header)

	offset := c.rawData[8:64+8]
	offsetData, err := hex.DecodeString(offset)
	if err != nil {
		return err
	}

	offsetValue := big.NewInt(0).SetBytes(offsetData)

	fmt.Println(offsetValue)

	temp := c.rawData[(288*2):(288*2)+256]
	fmt.Println(temp)

	for _, arg := range spec.Inputs {
		switch arg.Type {
		case "bytes[]":
			break
		case "uint32":
			break
		case "uint64":
			break
		case "uint256":
			break
		default:
			break
		}
	}

	return nil
}
