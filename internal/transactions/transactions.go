package transactions

import (
	"fmt"
	"math/big"
)

type Transaction struct {
	BlockHash        string `json:"blockHash"`
	BlockNumber      string `json:"blockNumber"`
	From             string `json:"from"`
	Gas              string `json:"gas"`
	GasPrice         string `json:"gasPrice"`
	Hash             string `json:"hash"`
	Input            string `json:"input"`
	Nonce            string `json:"nonce"`
	To               string `json:"to"`
	TransactionIndex string `json:"transactionIndex"`
	Value            string `json:"value"`
	//Type                 string        `json:"type"`
	V string `json:"v"`
	R string `json:"r"`
	S string `json:"s"`
	//MaxFeePerGas         string        `json:"maxFeePerGas"`
	//MaxPriorityFeePerGas string        `json:"maxPriorityFeePerGas"`
	//AccessList           []interface{} `json:"accessList"`
	//ChainID              string        `json:"chainId"`
}

// AddressRebalancing fills the address storage and totals the volume of transactions on them
func AddressRebalancing(txs []*Transaction) (map[string]*big.Int, error) {
	storage := make(map[string]*big.Int)

	for _, tx := range txs {

		val := new(big.Int)
		gasLimit := new(big.Int)
		gasPrice := new(big.Int)
		fee := new(big.Int)

		var ok bool

		if _, ok = val.SetString(tx.Value, 0); !ok {
			return nil, fmt.Errorf("couldn't parse tx value of '%s'", gasLimit)
		}
		if _, ok = gasPrice.SetString(tx.GasPrice, 0); !ok {
			return nil, fmt.Errorf("couldn't parse gas price of '%s'", gasPrice)
		}
		if _, ok = gasLimit.SetString(tx.Gas, 0); !ok {
			return nil, fmt.Errorf("couldn't parse gas limit of '%s'", gasLimit)
		}

		fee.Mul(gasLimit, gasPrice)

		if _, ok := storage[tx.From]; !ok {
			storage[tx.From] = val.Sub(val, fee)
		} else {
			storage[tx.From].Add(storage[tx.From], val.Sub(val, fee))
		}

		if _, ok := storage[tx.To]; !ok {
			storage[tx.To] = val
		} else {
			storage[tx.To].Add(storage[tx.To], val)
		}
	}
	return storage, nil
}
