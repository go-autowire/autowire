package service

import (
	"autowire"
	"log"
)

func init() {
	autowire.Autowire(&AppClient{})
}

type Client interface {
	Connect()
}

type AppClient struct {
	Type string
}

func (AppClient) Connect() {
	log.Println("Connected")
}
