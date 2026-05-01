package http

import (
	"errors"
	"net/http"

	"github.com/gorilla/mux"
)

type HTTPServer struct {
	handlers *HTTPHandlers
	server   *http.Server
}

func NewHTTPServer(handlers *HTTPHandlers) *HTTPServer {
	return &HTTPServer{
		handlers: handlers,
	}
}

func (s *HTTPServer) Run() error {
	router := mux.NewRouter()

	router.Path("/miners").Methods(http.MethodPost).HandlerFunc(s.handlers.HandleCreateNewMiner)
	router.Path("/miners").Methods(http.MethodGet).HandlerFunc(s.handlers.HandlerGetMiners)
	router.Path("/miners/salaries").Methods(http.MethodGet).HandlerFunc(s.handlers.HandleGetMinersSalaries)

	router.Path("/equipment").Methods(http.MethodPost).HandlerFunc(s.handlers.HandleBuyEquipment)
	router.Path("/equipment").Methods(http.MethodGet).HandlerFunc(s.handlers.HandleCheckEquipment)
	router.Path("/equipment/prices").Methods(http.MethodGet).HandlerFunc(s.handlers.HandleGetEquipmentPrices)

	router.Path("/company").Methods(http.MethodGet).HandlerFunc(s.handlers.HandleGetCompanyStatistics)
	router.Path("/company/complete").Methods(http.MethodPost).HandlerFunc(s.handlers.HandleCompleteGame)

	server := http.Server{
		Addr:    ":9091",
		Handler: router,
	}

	s.handlers.SetCloseServerFunc(server.Close)

	err := server.ListenAndServe()
	if errors.Is(err, http.ErrServerClosed) {
		return nil
	}

	return err
}
