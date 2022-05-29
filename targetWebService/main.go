package main

import (
	"flag"
	"fmt"
	"github.com/gorilla/mux"
	"net/http"
)

func main() {
	var serviceName = flag.String("s", "", "Name of this Service")
	var port = flag.String("p", "", "HTTP port to listen to")
	flag.Parse()
	router := mux.NewRouter()
	service := Service{*serviceName, *port, 0}
	router.HandleFunc("/", service.HomeHandler)
	router.HandleFunc("/health", service.HealthHandler)

	http.Handle("/", router)

	go http.ListenAndServe(fmt.Sprintf(":%v", service.Port), nil)
	go service.Logger()
	fmt.Printf("%v listening on port: %v, Press <Enter to exit>\n", service.Host, service.Port)
	fmt.Scanln()
}
