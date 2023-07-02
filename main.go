package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"
)

type ViaCep struct {
	Cep         string `json:"cep"`
	Logradouro  string `json:"logradouro"`
	Complemento string `json:"complemento"`
	Bairro      string `json:"bairro"`
	Localidade  string `json:"localidade"`
	Uf          string `json:"uf"`
	Ibge        string `json:"ibge"`
	Gia         string `json:"gia"`
	Ddd         string `json:"ddd"`
	Siafi       string `json:"siafi"`
}

type ApiCEP struct {
	Code       string `json:"code"`
	State      string `json:"state"`
	City       string `json:"city"`
	District   string `json:"district"`
	Address    string `json:"address"`
	Status     int    `json:"status"`
	Ok         bool   `json:"ok"`
	StatusText string `json:"statusText"`
}

type Message struct {
	Name    string
	Payload interface{}
}

func requestApiCep(ch chan Message, cep string) {
	contextWithTimeout, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	url := "https://cdn.apicep.com/file/apicep/" + cep + ".json"
	request, errorRequest := http.NewRequestWithContext(contextWithTimeout, "GET", url, nil)
	if errorRequest != nil {
		panic(errorRequest)
	}

	response, errorResponse := http.DefaultClient.Do(request)
	if errorResponse != nil {
		panic(errorResponse)
	}
	defer response.Body.Close()

	body, errorReadAll := io.ReadAll(response.Body)
	if errorReadAll != nil {
		panic(errorReadAll)
	}
	var apicep ApiCEP
	errorUnmarshal := json.Unmarshal(body, &apicep)
	if errorUnmarshal != nil {
		panic(errorUnmarshal)
	}
	message := Message{Name: "API CEP", Payload: apicep}
	ch <- message
	close(ch)
}

func requestViaCep(ch chan Message, cep string) {
	contextWithTimeout, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	url := "http://viacep.com.br/ws/" + cep + "/json/"
	request, errorRequest := http.NewRequestWithContext(contextWithTimeout, "GET", url, nil)
	if errorRequest != nil {
		panic(errorRequest)
	}
	var viaCEP ViaCep
	response, errorResponse := http.DefaultClient.Do(request)
	if errorResponse != nil {
		panic(errorResponse)
	}
	defer response.Body.Close()

	body, errorReadAll := io.ReadAll(response.Body)
	if errorReadAll != nil {
		panic(errorReadAll)
	}

	errorUnmarshal := json.Unmarshal(body, &viaCEP)
	if errorUnmarshal != nil {
		panic(errorUnmarshal)
	}

	message := Message{Name: "Via CEP", Payload: viaCEP}
	ch <- message
	close(ch)
}

func main() {
	cep := os.Args[1:]
	channelMessage := make(chan Message)

	go requestApiCep(channelMessage, cep[0])
	go requestViaCep(channelMessage, cep[0])

	select {
	case ch := <-channelMessage:
		fmt.Printf("%s: \n", ch.Name)
		fmt.Println(ch.Payload)

	case <-time.After(time.Second):
		panic("timeout")
	}
}
