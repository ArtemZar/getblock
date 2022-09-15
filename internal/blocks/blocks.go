package blocks

import (
	"context"
	"fmt"
	"getblock/configs"
	"getblock/internal/jrpc"
	"getblock/internal/transactions"
	"math/big"
	"strconv"
)

type Block struct {
	BaseFeePerGas    string                      `json:"baseFeePerGas"`
	Difficulty       string                      `json:"difficulty"`
	ExtraData        string                      `json:"extraData"`
	GasLimit         string                      `json:"gasLimit"`
	GasUsed          string                      `json:"gasUsed"`
	Hash             string                      `json:"hash"`
	LogsBloom        string                      `json:"logsBloom"`
	Miner            string                      `json:"miner"`
	MixHash          string                      `json:"mixHash"`
	Nonce            string                      `json:"nonce"`
	Number           string                      `json:"number"`
	ParentHash       string                      `json:"parentHash"`
	ReceiptsRoot     string                      `json:"receiptsRoot"`
	Sha3Uncles       string                      `json:"sha3Uncles"`
	Size             string                      `json:"size"`
	StateRoot        string                      `json:"stateRoot"`
	Timestamp        string                      `json:"timestamp"`
	TotalDifficulty  string                      `json:"totalDifficulty"`
	Transactions     []*transactions.Transaction `json:"transactions"`
	TransactionsRoot string                      `json:"transactionsRoot"`
	Uncles           []interface{}               `json:"uncles"`
}

// GetTrxsFromBlock getting one block by number and analyzing transaction by it.
// Collection with addresses and amounts sending to chenal for merge after.
func GetTrxsFromBlock(ctx context.Context, rpcClient *jrpc.RpcClient, storCh chan map[string]*big.Int, blockNum int64) error {
	blockNumStr := "0x" + strconv.FormatInt(blockNum, 16)
	blockData, err := rpcClient.GetBlockData(ctx, blockNumStr, configs.MethodGetBlockByNumber, true)
	if err != nil {
		return fmt.Errorf("can't take data from block %s. Get error: %v", blockNumStr, err)
	}
	var block = Block{}
	if err := blockData.GetObject(&block); err != nil {
		return fmt.Errorf("some error on json unmarshal level or json result field was null. Get error: %v", err)
	}

	addressesStorage, err := transactions.AddressRebalancing(block.Transactions)
	if err != nil {
		return fmt.Errorf("can't fills the address storage  from block %s. Get error: %v", blockNumStr, err)
	}
	storCh <- addressesStorage

	return nil
}
