package client

import (
	"encoding/json"

	"github.com/AthanorLabs/go-relayer/common"
)

type Client struct {
	endpoint string
}

func NewClient(endpoint string) *Client {
	return &Client{
		endpoint: endpoint,
	}
}

// Cancel calls relayer_submitTransaction.
func (c *Client) SubmitTransaction(req *common.SubmitTransactionRequest) (*common.SubmitTransactionResponse, error) {
	const method = "relayer_submitTransaction"

	params, err := json.Marshal(req)
	if err != nil {
		return nil, err
	}

	resp, err := PostRPC(c.endpoint, method, string(params))
	if err != nil {
		return nil, err
	}

	if resp.Error != nil {
		return nil, resp.Error
	}

	var res *common.SubmitTransactionResponse
	if err = json.Unmarshal(resp.Result, &res); err != nil {
		return nil, err
	}

	return res, nil
}
