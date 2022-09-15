package jrpc

import (
	"context"
	"fmt"
	"getblock/configs"
	"strconv"

	"github.com/ybbus/jsonrpc/v3"
)

type Client struct {
	Endpoint string
	Options  *jsonrpc.RPCClientOpts
}

// Init
func Init(apikey string) *Client {
	return &Client{
		Endpoint: configs.Endpoint,
		Options: &jsonrpc.RPCClientOpts{
			CustomHeaders: map[string]string{
				"x-api-key":    apikey,
				"Content-Type": "application/json",
			},
		},
	}
}

func (c *Client) NewRpc() jsonrpc.RPCClient {
	return jsonrpc.NewClientWithOpts(c.Endpoint, c.Options)
}

type RpcClient struct {
	RpcClient jsonrpc.RPCClient
}

// GetLastBlockNumber gets the number of the last block via the json-rpc method
func (rc *RpcClient) GetLastBlockNumber(ctx context.Context) (int64, error) {
	response, err := rc.RpcClient.Call(ctx, configs.MethodGetLastBlockNumber)
	if err != nil {
		return 0, fmt.Errorf("something is wrong with the json-rpc response. Error: %v", err)
	}
	numberInt, err := strconv.ParseInt(response.Result.(string)[2:], 16, 64)
	if err != nil {
		return 0, fmt.Errorf("block number is not parse to hex int: %v", err)
	}

	return numberInt, nil
}

// GetBlockData gets the data of the next block via the json-rpc method by hash
func (rc *RpcClient) GetBlockData(ctx context.Context, blockID string, method string, fullTxObj bool) (*jsonrpc.RPCResponse, error) {
	response, err := rc.RpcClient.Call(ctx, method, blockID, fullTxObj)
	if err != nil {
		return nil, fmt.Errorf("something is wrong with the json-rpc response, with error: %v", err)
	}
	return response, nil
}
