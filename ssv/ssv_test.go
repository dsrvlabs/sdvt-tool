package ssv

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"testing"

	"github.com/jarcoal/httpmock"
	"github.com/stretchr/testify/assert"
)

func TestSSV_Operator(t *testing.T) {
	// Activate httpmock
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	// Load fixture
	fixtureData, err := os.ReadFile("fixtures/validators_1044.json")
	assert.Nil(t, err)

	rpcURL := "https://api.ssv.network"
	operatorID := 1044

	httpmock.RegisterResponder(
		http.MethodGet,
		fmt.Sprintf("%s/api/v4/mainnet/validators/in_operator/%d?page=1&perPage=100", rpcURL, operatorID),
		httpmock.NewStringResponder(200, string(fixtureData)))

	// Test
	client := NewClient(rpcURL)
	validators, err := client.ListValidators(1044)
	assert.Nil(t, nil)

	// Asserts
	valResp := ValidatorsResponse{}
	err = json.Unmarshal(fixtureData, &valResp)
	assert.Nil(t, err)

	// Asserts
	for i, validator := range validators {
		assert.Equal(t, valResp.Validators[i].ID, validator.ID)
		assert.Equal(t, valResp.Validators[i].PublicKey, validator.PublicKey)
	}
}
