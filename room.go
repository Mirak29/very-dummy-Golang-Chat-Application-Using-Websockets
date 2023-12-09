package main

import(
	"log"
	"net/http"
	"github.com/gorilla/websocket"
)

type room struct{
	clients map[*client]bool
	join chan *client
	leave chan *client
	forward chan []byte
}

func newRoom() *room{
	return &room{
		forward: make(chan []byte),
		join: make(chan *client),
		leave: make(chan *client),
		clients: map[*client]bool{},
	}
}

func (r *room) run(){
	for{
		select{
			case client := <- r.join: 
				r.clients[client] = true
			case client := <- r.leave:
				delete(r.clients, client)
				close(client.receive)
			case msg := <- r.forward:
				for client := range r.clients{
					client.receive <- msg
				}
		}
	}
}