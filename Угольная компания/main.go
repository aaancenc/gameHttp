package main

import (
	"coal_company/company"
	"coal_company/http"
	"context"
	"fmt"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	company := company.NewCompany(ctx)
	handlers := http.NewHTTPHandlers(company)
	server := http.NewHTTPServer(handlers)

	fmt.Println("game started!")
	if err := server.Run(); err != nil {
		fmt.Println("error while running HTTP server:", err)
	}
}
