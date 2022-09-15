package server

import (
	"context"
	"fmt"
	"getblock/configs"
	"getblock/internal/blocks"
	"getblock/internal/jrpc"
	"math/big"
	"net/http"
	"sync"

	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"
)

type Server struct {
	config *Config
	logger *log.Logger
	router *mux.Router
}

func New(config *Config) *Server {
	return &Server{
		config: config,
		logger: log.New(),
		router: mux.NewRouter(),
	}
}

func (s *Server) Start(ctx context.Context) error {
	if err := s.configureLogger(); err != nil {
		return err
	}

	s.configureRouter(ctx)

	s.logger.Info("The Server is starting")

	return http.ListenAndServe(s.config.BindAddr, s.router)
}

func (s *Server) configureLogger() error {
	level, err := log.ParseLevel(s.config.Loglevel)
	if err != nil {
		return err
	}
	s.logger.SetLevel(level)
	return nil
}

func (s *Server) configureRouter(ctx context.Context) {
	s.router.HandleFunc("/", s.handleIndex())
	s.router.HandleFunc("/find_address/{apikey}", s.handleFindAddress(ctx))

}

func (s *Server) handleIndex() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("test")) // nolint:errcheck // ok
	}
}

func (s *Server) handleFindAddress(ctx context.Context) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// get Api key
		vars := mux.Vars(r)
		apiKey := vars["apikey"] // or os.Args[1]

		// init json rpc client
		client := jrpc.Init(apiKey)
		rpcClient := &jrpc.RpcClient{
			RpcClient: client.NewRpc(),
		}

		fullAddrStorage := make(map[string]*big.Int)
		storCh := make(chan map[string]*big.Int, 10)
		wg := sync.WaitGroup{}

		lastBlockNumber, err := rpcClient.GetLastBlockNumber(context.Background())
		if err != nil {
			s.logger.Error(err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		// Bypass of blocks
		for i := int64(0); i < configs.NumOfBlocs; i++ {
			wg.Add(1)
			go func(blockNum int64) {
				defer wg.Done()
				err := blocks.GetTrxsFromBlock(ctx, rpcClient, storCh, blockNum)
				if err != nil {
					s.logger.Errorf("error get trxs from block, %v", err)
					//TODO abort process or not
				}
			}(lastBlockNumber - i)
		}

		// merging all collections in one
		readCh := make(chan struct{})
		go func() {
		Loop:
			for {
				newData, ok := <-storCh
				if newData == nil && !ok {
					break Loop
				}
				for k, v := range newData {
					if _, ok := fullAddrStorage[k]; !ok {
						fullAddrStorage[k] = v
					} else {
						fullAddrStorage[k].Add(fullAddrStorage[k], v)
					}
				}
			}
			readCh <- struct{}{}
		}()

		wg.Wait()
		close(storCh)
		<-readCh
		close(readCh)

		// finding max value
		address := ""
		change := big.NewInt(0)
		for key, val := range fullAddrStorage {
			if val.Cmp(change) == 1 {
				change = val
				address = key
			}
		}
		//nolint:errcheck // ok
		w.Write(
			[]byte(
				fmt.Sprintf(
					"Aдрес, баланс которого изменился (в любую сторону) больше остальных за последние сто блоков\n%s",
					address)))
	}
}
