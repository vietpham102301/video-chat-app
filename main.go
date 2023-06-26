package main

import (
	"log"
	"net/http"
	"video-chat-app/server"
)

func main() {
	server.AllRooms.Init()
	server.AllUsers.Init()

	http.HandleFunc("/create", server.CreateRoomRequestHandler)
	http.HandleFunc("/join", server.JoinRoomRequestHandler)
	http.HandleFunc("/notify", server.NotifyUserHandler)
	http.HandleFunc("/establish", server.HandleWebSocketConnection)



	log.Println("Starting Server on Port 8000")
	err := http.ListenAndServeTLS(":8000", "./cert.pem", "./key.pem", nil)
	if err != nil {
		log.Fatal(err)
	}
	//http.ListenAndServe(":8000", nil)
}
