package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
)

func main() {
	jfile, err := ioutil.ReadFile("./config.json")
	if err != nil {
		log.Fatal(err)
	}
	data := make([]Service, 5)
	err2 := json.Unmarshal(jfile, &data)

	if err2 != nil {
		log.Fatal(err2)
	}

	go getHealthServices(data)
	for i := 0; i < 5; i++ {
		log.Printf("Balancer target: %v : %v\n", data[i].Host, data[i].Port)
	}

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		balanceProxyWebService(w, r, data)
	})
	if err := http.ListenAndServe(":8080", nil); err != nil && err != http.ErrServerClosed {
		log.Fatalf("listen: %s\n", err)
	}
}
