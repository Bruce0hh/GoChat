package websocket

type Hub struct {
	//上线
	Login chan *Client
	//下线
	Logout chan *Client
	//所有在线客户端的内存地址
	Clients map[*Client]bool
}

func CreateHubFactory() *Hub {
	return &Hub{
		Login:   make(chan *Client),
		Logout:  make(chan *Client),
		Clients: make(map[*Client]bool),
	}
}

func (h *Hub) Run() {
	for {
		select {
		case client := <-h.Login:
			h.Clients[client] = true
		case client := <-h.Logout:
			if _, ok := h.Clients[client]; ok {
				_ = client.Conn.Close()
				delete(h.Clients, client)
			}
		}
	}
}
