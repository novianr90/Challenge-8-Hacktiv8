package main

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"sync"
	"time"

	"github.com/go-resty/resty/v2"
)

type Data struct {
	Water int `json:"water"`
	Wind  int `json:"wind"`
}

func main() {
	var wg sync.WaitGroup
	dataChan := make(chan Data)

	go func() {
		for {
			data := Data{
				Water: rand.Intn(100),
				Wind:  rand.Intn(100),
			}

			dataChan <- data
			time.Sleep(15 * time.Second)
		}
	}()

	for {
		wg.Add(1)

		data := <-dataChan

		go PostRequest(&wg, data)
	}

}

func PostRequest(wg *sync.WaitGroup, jsonData Data) {
	defer wg.Done()

	client := resty.New()

	jsonPayload, err := json.Marshal(jsonData)

	if err != nil {
		fmt.Println("err di marshal")
		return
	}

	resp, err := client.R().
		SetHeader("Content-Type", "application/json").
		SetBody(jsonPayload).
		Post("https://jsonplaceholder.typicode.com/posts")

	if err != nil {
		fmt.Println("err di request")
		return
	}

	respBody := string(resp.Body())

	var data Data
	err = json.Unmarshal([]byte(respBody), &data)

	if err != nil {
		fmt.Println("err di unmarshal")
		return
	}

	waterStatus := ""
	if data.Water < 5 {
		waterStatus = "aman"
	} else if data.Water >= 5 && data.Water <= 8 {
		waterStatus = "siaga"
	} else {
		waterStatus = "bahaya"
	}

	windStatus := ""
	if data.Wind < 6 {
		windStatus = "aman"
	} else if data.Wind >= 6 && data.Wind <= 15 {
		windStatus = "siaga"
	} else {
		windStatus = "bahaya"
	}

	fmt.Printf("{\n\t\"water\": %d,\n\t\"wind\": %d\n}\n", data.Water, data.Wind)
	fmt.Println("status water :", waterStatus)
	fmt.Println("status wind :", windStatus)
}
