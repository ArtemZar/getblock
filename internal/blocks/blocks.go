package blocks

import (
	"context"
	"getblock/configs"
	"getblock/internal/jrpc"
	"getblock/internal/transactions"
	log "github.com/sirupsen/logrus"
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

// GetLastBlockData gets the data of the last block via the json-rpc method by number
func GetLastBlockData(ctx context.Context, rpcClient *jrpc.RpcClient) *Block {
	lastBlockNumber := rpcClient.GetLastBlockNumber(ctx)
	blockData := rpcClient.GetBlockData(ctx, lastBlockNumber, configs.MethodGetBlockByNumber, true)
	block := &Block{}
	if err := blockData.GetObject(&block); err != nil || block == nil {
		log.Info("some error on json unmarshal level or json result field was null ", err)
	}
	return block
}

// CreateBlocsStorege сreates a storege of chain block data
func CreateBlocsStorege(ctx context.Context, rpcClient *jrpc.RpcClient) [configs.NumOfBlocs]Block {
	var blocksStorage [configs.NumOfBlocs]Block
	nextBlock := GetLastBlockData(ctx, rpcClient)
	blocksStorage[0] = *nextBlock
	for i := 1; i < configs.NumOfBlocs; i++ {
		blockData := rpcClient.GetBlockData(ctx, nextBlock.ParentHash, configs.MethodGetBlockByHash, true)
		nextBlock = &Block{}
		if err := blockData.GetObject(&nextBlock); err != nil || nextBlock == nil {
			log.Info("some error on json unmarshal level or json result field was null ", err)
		}
		blocksStorage[i] = *nextBlock
	}
	return blocksStorage
}
