package main

import (
	"fmt"
	"net/http"
	"os"
)

const path string = "/"

func main() {
	serv_mux := http.NewServeMux()
	serv_mux.Handle(path, http.FileServer(http.Dir(".")))
	server := http.Server{
		Addr:                         ":8080",
		Handler:                      serv_mux,
		DisableGeneralOptionsHandler: true,
	}

	err := server.ListenAndServe()

	if err != nil {
		fmt.Println("fail to start the server")
		os.Exit(1)
	}
}
