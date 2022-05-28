package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"math"
	"net/http"
	"net/http/httputil"
	"net/url"
	"sync"
	"time"
)

var respDelay int64 = 10

var mutex sync.Mutex

type Service struct {
	Host   string `json:host`
	Port   string `json:port`
	Health int    `json:health`
}

func selectionsort(items []Service) {
	var n = len(items)
	for i := 0; i < n; i++ {
		var minIdx = i
		for j := i; j < n; j++ {
			if items[j].Health < items[minIdx].Health {
				minIdx = j
			}
		}
		items[i], items[minIdx] = items[minIdx], items[i]
	}
}

func getHealthServices(arrayWS []Service) {
	for true {
		for i := 0; i < len(arrayWS); i++ {
			var tmpService Service
			var trueServer string

			mutex.Lock()
			trueServer = arrayWS[i].Host + ":" + arrayWS[i].Port
			mutex.Unlock()

			if len(trueServer) == 0 {
				return
			}
			if ((time.Now().Unix() % respDelay) != 0) && (arrayWS[i].Health == math.MaxInt) {
				continue
			}
			resp, _ := http.Get(trueServer + "/health")
			if resp == nil {
				mutex.Lock()
				arrayWS[i].Health = math.MaxInt
				mutex.Unlock()
				continue
			} else {
				defer resp.Body.Close()
			}

			if resp.StatusCode == http.StatusOK {
				body, err := ioutil.ReadAll(resp.Body)
				if err != nil {
					log.Fatal(err)
					return
				}
				err = json.Unmarshal(body, &tmpService)
				if err != nil {
					return
				}
			}
			mutex.Lock()
			arrayWS[i].Health = tmpService.Health
			mutex.Unlock()
		}
		mutex.Lock()
		selectionsort(arrayWS)
		mutex.Unlock()
	}
}

func balanceProxyWebService(w http.ResponseWriter, r *http.Request, arrayWS []Service) {
	var trueServer string

	mutex.Lock()
	for i := 0; i < len(arrayWS); i++ {
		if arrayWS[i].Health == math.MaxInt {
			continue
		}
		trueServer = arrayWS[i].Host + ":" + arrayWS[i].Port
		break
	}
	mutex.Unlock()

	url, err := url.Parse(trueServer)
	if err != nil {
		log.Println(err)
		return
	}

	proxy := httputil.NewSingleHostReverseProxy(url)
	proxy.ServeHTTP(w, r)

}
