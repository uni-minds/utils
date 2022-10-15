package tools

import (
	"crypto/tls"
	"fmt"
	"net/http"
	"time"
)

const timeout = 10 * time.Second

type ClientPool struct {
	clients    chan *http.Client
	clientUsed int
	clientMax  int
}

func NewClientPool(maxClients int) *ClientPool {
	pool := new(ClientPool)
	pool.Init(maxClients)
	return pool
}

func (p *ClientPool) Init(maxClients int) {
	p.clientMax = maxClients
	p.clients = make(chan *http.Client, maxClients)
}

func (p *ClientPool) GetClient() (client *http.Client) {
	t := time.NewTicker(timeout)
	for {
		select {
		case client = <-p.clients:
			p.clientUsed += 1
			return client

		case <-t.C:
			fmt.Println("number of clients reach max, no more available")
			return nil

		default:
			if p.clientUsed < p.clientMax {
				tr := &http.Transport{
					TLSClientConfig:       &tls.Config{InsecureSkipVerify: true},
					TLSHandshakeTimeout:   timeout,
					ResponseHeaderTimeout: timeout,
					ExpectContinueTimeout: timeout,
				}

				client = &http.Client{
					Transport: tr,
					Timeout:   timeout,
				}

				p.clientUsed += 1
				fmt.Printf("generate new client, %d use / %d max\n", p.clientUsed, p.clientMax)
				return client
			}
		}
	}
}

func (p *ClientPool) Recycle(client *http.Client) {
	if client != nil {
		p.clientUsed -= 1
		p.clients <- client
	}
}
