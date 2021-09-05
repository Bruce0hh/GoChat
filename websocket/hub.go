package websocket

import "sync"

type Hub struct {
	//上线
	Login chan *Client
	//下线
	Logout chan *Client
	//所有在线客户端的内存地址
	Clients map[string]*Client

	RWLock sync.RWMutex
}

func CreateHubFactory() *Hub {
	return &Hub{
		Login:   make(chan *Client),
		Logout:  make(chan *Client),
		Clients: make(map[string]*Client),
	}
}

func (h *Hub) Run() {
	for {
		select {
		case client := <-h.Login:
			h.RWLock.Lock()
			h.Clients[client.flag] = client
			h.RWLock.Unlock()
		case client := <-h.Logout:
			if _, ok := h.Clients[client.flag]; ok {
				_ = client.Conn.Close()
				delete(h.Clients, client.flag)
			}
		}
	}
}
