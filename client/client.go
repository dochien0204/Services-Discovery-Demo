package main

import (
	"fmt"
	"io"
	"net/http"
	"time"

	consulapi "github.com/hashicorp/consul/api"
)

var url string

func main() {
	fmt.Println("Nhuw cc")
	serviceDiscoveryWithConsul()
	fmt.Println("Starting Client.")
	var client = &http.Client{
		Timeout: time.Second * 30,
	}
	callServerEvery(10*time.Second, client)
}

func serviceDiscoveryWithConsul() {
	config := consulapi.DefaultConfig()
	consul, error := consulapi.NewClient(config)
	if error != nil {
		fmt.Println(error)
	}
	services, error := consul.Agent().Services()
	if error != nil {
		fmt.Println(error)
	}

	service := services["hello"]
	address := service.Address
	port := service.Port
	url = fmt.Sprintf("http://%s:%v/hello", address, port)
}

func Hello(t time.Time, client *http.Client) {
	response, err := client.Get(url)
	if err != nil {
		fmt.Println(err)
		return
	}
	body, _ := io.ReadAll(response.Body)
	fmt.Printf("%s. Time is %v\n", body, t)
}

func callServerEvery(d time.Duration, client *http.Client) {
	for x := range time.Tick(d) {
		Hello(x, client)
	}
}
