package game

import "com.pawmot/goball/ws"

type Room struct {
	clients map[*ws.Client]bool

	Broadcast chan []byte

	Register chan *ws.Client

	Unregister chan *ws.Client
}

func NewRoom() *Room {
	return &Room{
		Broadcast:  make(chan []byte),
		Register:   make(chan *ws.Client),
		Unregister: make(chan *ws.Client),
		clients:    make(map[*ws.Client]bool),
	}
}

func (r *Room) Run() {
	for {
		select {
		case client := <-r.Register:
			r.clients[client] = true
		case client := <-r.Unregister:
			if _, ok := r.clients[client]; ok {
				delete(r.clients, client)
				close(client.Send)
			}
		case message := <-r.Broadcast:
			for client := range r.clients {
				select {
				case client.Send <- message:
				default:
					close(client.Send)
					delete(r.clients, client)
				}
			}
		}
	}
}
