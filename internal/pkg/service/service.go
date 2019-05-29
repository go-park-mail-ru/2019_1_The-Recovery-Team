package service

import (
	"fmt"
	"log"
	"strconv"

	consulapi "github.com/hashicorp/consul/api"
	"google.golang.org/grpc"
	"google.golang.org/grpc/balancer/roundrobin"
)

func Connect(addr string, name string) *grpc.ClientConn {
	conn, err := grpc.Dial(fmt.Sprintf("srv://%s/%s", addr, name),
		grpc.WithInsecure(), grpc.WithBalancerName(roundrobin.Name), grpc.WithBlock())
	if err != nil {
		log.Fatal("Can't connect to profile service:", err)
	}
	return conn
}

func createConsulClient(consulAddr string, consulPort int) *consulapi.Client {
	config := consulapi.DefaultConfig()
	config.Address = consulAddr + ":" + strconv.Itoa(consulPort)
	consul, err := consulapi.NewClient(config)
	if err != nil {
		log.Fatal("Can't connect to consul:", err)
	}
	return consul
}

func RegisterInConsul(consulAddr string, consulPort int, name, addr string, port int) {
	consul := createConsulClient(consulAddr, consulPort)
	err := consul.Agent().ServiceRegister(&consulapi.AgentServiceRegistration{
		ID:      name + strconv.Itoa(port),
		Name:    name,
		Port:    port,
		Address: addr,
	})
	if err != nil {
		log.Fatalf("Can't add %s to resolver: %s", name, err.Error())
	}
	log.Println("Registered in resolver", name, port)
}

func DeregisterInConsul(consulAddr string, consulPort int, name string, port int) {
	consul := createConsulClient(consulAddr, consulPort)
	err := consul.Agent().ServiceDeregister(name + strconv.Itoa(port))
	if err != nil {
		log.Fatal("Can't remove service from resolver:", err)
	}
	log.Println("Deregistered in resolver", name, port)
}
