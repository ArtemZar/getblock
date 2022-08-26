package transactions

import (
	"math/big"
	"sync"
)

type Transaction struct {
	BlockHash            string        `json:"blockHash"`
	BlockNumber          string        `json:"blockNumber"`
	From                 string        `json:"from"`
	Gas                  string        `json:"gas"`
	GasPrice             string        `json:"gasPrice"`
	Hash                 string        `json:"hash"`
	Input                string        `json:"input"`
	Nonce                string        `json:"nonce"`
	To                   string        `json:"to"`
	TransactionIndex     string        `json:"transactionIndex"`
	Value                string        `json:"value"`
	Type                 string        `json:"type"`
	V                    string        `json:"v"`
	R                    string        `json:"r"`
	S                    string        `json:"s"`
	MaxFeePerGas         string        `json:"maxFeePerGas"`
	MaxPriorityFeePerGas string        `json:"maxPriorityFeePerGas"`
	AccessList           []interface{} `json:"accessList"`
	ChainID              string        `json:"chainId"`
}

// AddressRebalancing fills the address storage and totals the volume of transactions on them
func AddressRebalancing(storage *map[string]*big.Int, txs []*Transaction, wg *sync.WaitGroup, mu *sync.Mutex) {
	defer wg.Done()
	var min = big.NewInt(0)

	for _, tx := range txs {
		val := new(big.Int)
		val.SetString(tx.Value, 0)
		if val.Cmp(min) == 0 {
			return
		}

		mu.Lock()
		if _, ok := (*storage)[tx.From]; !ok {
			(*storage)[tx.From] = val
		} else {
			(*storage)[tx.From].Add((*storage)[tx.From], val)
		}

		if _, ok := (*storage)[tx.To]; !ok {
			(*storage)[tx.To] = val
		} else {
			(*storage)[tx.To].Add((*storage)[tx.To], val)
		}
		mu.Unlock()

	}
}
