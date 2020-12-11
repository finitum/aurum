package main

import (
	"errors"
	aurum "github.com/finitum/aurum/clients/go"
	"github.com/finitum/aurum/pkg/jwt"
)

type Client struct {
	aurum.Client
	tp jwt.TokenPair
}

type ClientManager struct {
	clients []Client

	active int
}

func NewClientManager() *ClientManager {
	return &ClientManager{
		clients: nil,
		active: 0,
	}
}


func (c *ClientManager) NumClients() int {
	return len(c.clients)
}

func (c *ClientManager) AddClient(url string) error {
	for i, client := range c.clients {
		if client.GetUrl() == url {
			c.active = i
			return nil
		}
	}

	client, err := aurum.NewRemoteClient(url, true)
	if err != nil {
		return err
	}
	c.clients = append(c.clients, Client {
		Client: client,
	})

	c.active = len(c.clients) - 1

	return nil
}

func (c *ClientManager) GetActiveClient() (*Client, error) {
	if c.active >= len(c.clients) {
		return nil, errors.New("no clients added yet")
	}

	return &c.clients[c.active], nil
}

func (c *ClientManager) GetActiveClientIndex() int {
	return c.active
}

func (c *ClientManager) SetActiveClient(i int) error {
	if i >= len(c.clients) {
		return errors.New("no client found with this index")
	}

	c.active = i

	return nil
}

func (c *ClientManager) GetClient(i int) (Client, error){
	if i >= len(c.clients) {
		return Client{}, errors.New("no client found with this index")
	}

	return c.clients[i], nil
}