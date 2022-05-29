package main

/*
Распределитель нагрузки

Задача:
Написать веб-сервис, распределяющий нагрузку между N веб сервисами равномерно.
Если один из веб сервисов перестает отвечать, то исключить его из массива веб сервисов и
распределять нагрузку между N-1 сервисами.
Нагрузкой для каждого сервиса является количество запросов, находящихся у него в обработке.
Относительно этого числа распределение должно быть равномерным.

На выходе у вас должен получиться HTTP прокси сервер, который равномерно
распределит нагрузку между N сервисами.


Требования к ЯП: -
Требования к развертке: Docker
Требования к конфигурации: список веб сервисов с их адресами прописать в отдельном
конфигурационном файле.


Шаги тестирования:
Поднять минимум 5 сервисов таргетов
Поднять сервис распределитель
Запустить скрипт спама запросами
Каждые 10 секунд каждый из сервисов-таргетов должен логировать текущее
количество запросов в обработке
На 100-ой секунде отключить один из сервисов-таргетов
На 200-ой включить его обратно
По итогу тестирования лог файлы таргетов должны свидетельствовать о равномерном
распределении нагрузки, в момент отключения таргета n,
таргеты N-1 должны разделить между собой нагрузку.
В момент включения таргета n, таргеты N должны освободить часть нагрузки под поднявшийся таргет.

*/
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
	Name   string `json:name`
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
			trueServer = "http://" + arrayWS[i].Name + ":" + arrayWS[i].Port
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
		trueServer = "http://" + arrayWS[i].Name + ":" + arrayWS[i].Port
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
