package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	consulapi "github.com/hashicorp/consul/api"
)

func main() {
	servicesRegistryConsul()
	log.Println("Starting hello server")
	http.HandleFunc("/hello", HelloHi)
	http.HandleFunc("/check", Check)
	http.ListenAndServe(":8080", nil)
}

func servicesRegistryConsul() {
	config := consulapi.DefaultConfig()

	consul, err := consulapi.NewClient(config)

	if err != nil {
		log.Println(err)
	}

	servicesId := "hello"
	port := 8080
	address, _ := os.Hostname()

	registration := &consulapi.AgentServiceRegistration{
		ID:      servicesId,
		Name:    "hello",
		Port:    port,
		Address: address,
		Check: &consulapi.AgentServiceCheck{
			HTTP:     fmt.Sprintf("http://%s:%v/check", address, port),
			Interval: "10s",
			Timeout:  "30s",
		},
		Tags: []string{"hello"},
	}

	fmt.Printf("http://%s:%v/", address, port)

	regisErr := consul.Agent().ServiceRegister(registration)

	if regisErr != nil {
		log.Printf("Failed to register services %s:%v", address, port)
		log.Println(regisErr)
	} else {
		log.Printf("Successfully register services: %s:%v", address, port)
	}
}

func HelloHi(w http.ResponseWriter, r *http.Request) {
	log.Println("helloworld service is called.")
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "Hello world.")
}

func Check(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "Consul check")
}
