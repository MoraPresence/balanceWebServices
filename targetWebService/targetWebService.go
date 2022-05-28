package main

import (
	"encoding/json"
	"log"
	"net/http"
	"sync"
	"time"
)

var loggerDelay time.Duration = 10
var targetWorkDelay time.Duration = 100

var mutex sync.Mutex

type Service struct {
	Host   string `json:host`
	Port   string `json:port`
	Health uint   `json:health`
}

func (service *Service) HomeHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)

	mutex.Lock()
	service.Health++
	mutex.Unlock()
	go func() {
		time.Sleep(targetWorkDelay * time.Second)
		mutex.Lock()
		if service.Health != 0 {
			service.Health--
		}
		mutex.Unlock()
	}()
}

func (service *Service) HealthHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	mutex.Lock()
	jdata, _ := json.Marshal(service)
	w.Write(jdata)
	mutex.Unlock()
	return
}

func (service *Service) Logger() {
	for true {
		start := time.Now()
		time.Sleep(loggerDelay * time.Second)

		if time.Since(start) >= (loggerDelay * time.Second) {
			mutex.Lock()
			log.Printf("Сurrent amount request: %v. From Service: %v\n", service.Health, service.Host)
			mutex.Unlock()
		}
	}
}
