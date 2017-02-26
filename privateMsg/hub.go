package main

import (
	"sync"
)

// hub maintains the set of active clients
type Hub struct {
	// Registered clients.
	clients      map[int]*Client
	tokens       map[int]string
	mutex_client sync.RWMutex
	mutex_token  sync.RWMutex
}

func InitHub() *Hub {
	return &Hub{
		clients: make(map[int]*Client),
		tokens:  make(map[int]string),
	}
}

func (self *Hub) Register(c *Client) {

	self.mutex_client.Lock()
	self.clients[c.userid] = c
	self.mutex_client.Unlock()
}

func (self *Hub) UnRegister(c *Client) {

	self.mutex_client.Lock()
	delete(self.clients, c.userid)
	close(c.send)
	self.mutex_client.Unlock()

}

func (self *Hub) GetClient(userid int) *Client {

	self.mutex_client.RLock()
	var client *Client = nil
	if c, ok := self.clients[userid]; ok {
		client = c
	}
	self.mutex_client.RUnlock()

	return client
}

func (self *Hub) BindToken(uid int, token string) {
	self.mutex_token.Lock()
	self.tokens[uid] = token
	self.mutex_token.Unlock()
}

func (self *Hub) isValidToken(uid int, token string) bool {
	self.mutex_token.RLock()
	t := self.tokens[uid]
	self.mutex_token.RUnlock()

	if t == token {
		return true
	}

	return false
}
