package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"tiny-ledger/internal/app"
	handler "tiny-ledger/internal/handler/http"
	"tiny-ledger/internal/storage"
)

func main() {
	port := flag.Int("port", 8080, "port to listen on")
	flag.Parse()

	memStore := storage.NewMemStore()
	ledgerSvc := app.NewLedger(memStore)
	httpHandler := handler.NewHTTPHandler(ledgerSvc)

	addr := fmt.Sprintf(":%d", *port)
	log.Printf("tiny-ledger listening on %s", addr)
	log.Fatal(http.ListenAndServe(addr, httpHandler))
}
