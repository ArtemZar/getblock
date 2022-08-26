package jrpc

import (
	"context"
	"getblock/configs"
	log "github.com/sirupsen/logrus"
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
func (rc *RpcClient) GetLastBlockNumber(ctx context.Context) string {
	response, err := rc.RpcClient.Call(ctx, configs.MethodGetLastBlockNumber)
	if err != nil {
		log.Info("Something is wrong with the json-rpc response", err)
	}
	return response.Result.(string)
}

// GetBlockData gets the data of the next block via the json-rpc method by hash
func (rc *RpcClient) GetBlockData(ctx context.Context, blockID, method string, fullTxObj bool) *jsonrpc.RPCResponse {
	response, err := rc.RpcClient.Call(ctx, method, blockID, fullTxObj)
	if err != nil {
		log.Info("Something is wrong with the json-rpc response", err)
	}
	return response
}
