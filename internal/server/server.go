package server

import (
	"context"
	"fmt"
	"getblock/internal/blocks"
	"getblock/internal/jrpc"
	"getblock/internal/transactions"
	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"
	"math/big"
	"net/http"
	"sync"
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

func (s *Server) Start() error {
	if err := s.configureLogger(); err != nil {
		return err
	}

	s.configureRouter()

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

func (s *Server) configureRouter() {
	s.router.HandleFunc("/", s.handleIndex())
	s.router.HandleFunc("/find_address/{apikey}", s.handleFindAddress())

}

func (s *Server) handleIndex() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("test")) // nolint:errcheck // ok
	}
}

func (s *Server) handleFindAddress() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		apiKey := vars["apikey"] // or os.Args[1]
		client := jrpc.Init(apiKey)
		rpcClient := &jrpc.RpcClient{
			RpcClient: client.NewRpc(),
		}

		BlocksStorage := blocks.CreateBlocsStorege(context.Background(), rpcClient)

		addressesStorage := make(map[string]*big.Int)
		wg := &sync.WaitGroup{}
		mu := &sync.Mutex{}

		for _, vol := range BlocksStorage {
			wg.Add(1)
			go transactions.AddressRebalancing(&addressesStorage, vol.Transactions, wg, mu)
		}
		wg.Wait()

		address := ""
		change := big.NewInt(0)
		for key, val := range addressesStorage {
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
