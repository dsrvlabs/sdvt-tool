package ssv

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type ValidatorsResponse struct {
	Validators []Validator `json:"validators"`
	Pagination Pagination  `json:"pagination"`
}

type Validator struct {
	ID               int           `json:"id"`
	PublicKey        string        `json:"public_key"`
	Cluster          string        `json:"cluster"`
	OwnerAddress     string        `json:"owner_address"`
	Status           string        `json:"status"`
	IsValid          bool          `json:"is_valid"`
	IsDeleted        bool          `json:"is_deleted"`
	IsPublicKeyValid bool          `json:"is_public_key_valid"`
	IsSharesValid    bool          `json:"is_shares_valid"`
	IsOperatorsValid bool          `json:"is_operators_valid"`
	Operators        []int         `json:"operators"`
	ValidatorInfo    ValidatorInfo `json:"validator_info"`
	Version          string        `json:"version"`
	Network          string        `json:"network"`
}

type ValidatorInfo struct {
	Index           int    `json:"index,omitempty"`
	Status          string `json:"status,omitempty"`
	ActivationEpoch int    `json:"activation_epoch,omitempty"`
}

type Pagination struct {
	Total        int `json:"total"`
	Pages        int `json:"pages"`
	PerPage      int `json:"per_page"`
	Page         int `json:"page"`
	CurrentFirst int `json:"current_first"`
	CurrentLast  int `json:"current_last"`
}

type Client interface {
	ListValidators(operatorID int) ([]Validator, error)
}

type ssvClient struct {
	rpc string
}

func (c ssvClient) ListValidators(operatorID int) ([]Validator, error) {
	validators := []Validator{}
	currentPage := 1
	for {
		url := fmt.Sprintf("%s/api/v4/mainnet/validators/in_operator/%d?page=%d&perPage=100", c.rpc, operatorID, currentPage)
		req, err := http.NewRequest(http.MethodGet, url, nil)
		if err != nil {
			return nil, err
		}

		req.Header.Add("accept", "*/*")

		cli := http.Client{}
		resp, err := cli.Do(req)
		if err != nil {
			return nil, err
		}

		rawData, err := io.ReadAll(resp.Body)
		if err != nil {
			return nil, err
		}

		valResp := ValidatorsResponse{}
		err = json.Unmarshal(rawData, &valResp)
		if err != nil {
			return nil, err
		}

		if valResp.Pagination.Total != 0 {
			validators = append(validators, valResp.Validators...)
		}

		// Exit condition
		if valResp.Pagination.Total < valResp.Pagination.PerPage {
			break
		}

		currentPage++
	}

	return validators, nil
}

func NewClient(rpc string) Client {
	return &ssvClient{
		rpc: rpc,
	}
}
