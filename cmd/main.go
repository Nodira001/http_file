package main

import (
	"log"
	"net"
	"net/http"
	"os"

	"github.com/Nodira001/http/cmd/app"
	"github.com/Nodira001/http/pkg/banners"
)

func main() {
	host := "127.0.0.1"
	port := "9999"
	if err := execute(host, port); err != nil {
		log.Println(err)
		os.Exit(1)
	}
}
func execute(host string, port string) (err error) {
	mux := http.NewServeMux()
	service := banners.NewService()
	server := app.NewServer(mux, service)
	server.Init()
	var svr = &http.Server{
		Addr:    net.JoinHostPort(host, port),
		Handler: server,
	}

	return svr.ListenAndServe()
}
