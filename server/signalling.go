package server

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/websocket"
	"io/ioutil"
	"log"
	"net/http"
	"sync"
)

// AllRooms is the global hashmap for the server
var AllRooms RoomMap

type UserMap struct {
	Mutex sync.RWMutex
	Map   map[string]*websocket.Conn
}

func (u *UserMap) Init() {
	u.Map = make(map[string]*websocket.Conn)
}

var AllUsers UserMap

//var UserConnections = make(map[string]*websocket.Conn)

// CreateRoomRequestHandler Create a Room and return roomID
func CreateRoomRequestHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	roomID := AllRooms.CreateRoom()

	type resp struct {
		Room string `json:"room"`
	}

	websocketURL := fmt.Sprintf("wss://%s/join?roomID=%s", r.Host, roomID)

	response := resp{
		Room: websocketURL,
	}

	log.Println(AllRooms.Map)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func HandleWebSocketConnection(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	// Get the user ID from the request or any other identifier
	userID := r.URL.Query().Get("userID")

	// Upgrade the HTTP connection to WebSocket
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		// Handle the error
		return
	}

	// Store the WebSocket connection
	AllUsers.Mutex.Lock()
	AllUsers.Map[userID] = conn
	AllUsers.Mutex.Unlock()

	// Start listening for messages from the user
	//go ListenForUserMessages(conn, userID)
}



type NotificationRequest struct {
	RecipientId  string `json:"recipientId"`
	Notification string `json:"notification"`
}

func NotifyUserHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "POST")
	// Read the request body
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Failed to read request body"))
		return
	}

	// Parse the request body into a NotificationRequest struct
	var request NotificationRequest
	err = json.Unmarshal(body, &request)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Failed to parse request body"))
		return
	}

	AllUsers.Mutex.RLock()
	conn, ok := AllUsers.Map[request.RecipientId]
	AllUsers.Mutex.RUnlock()
	// Validate and process the user ID and message
	// ...

	// Notify the user

	log.Println("sms d")
	if !ok {
		log.Println("not ok")
		// User not found or WebSocket connection not established
		return
	}

	// Send the notification message to the user
	err = conn.WriteJSON(request.Notification)
	if err != nil {
		log.Println("error here!")
	}

	// Respond with a success message
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Notification sent successfully"))
}


var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

type broadcastMsg struct {
	Message map[string]interface{}
	RoomID  string
	Client  *websocket.Conn
}

var broadcast = make(chan broadcastMsg, 100)

func broadcaster() {
	for msg := range broadcast {
		clients := AllRooms.Get(msg.RoomID)

		for _, client := range clients {
			if client.Conn != msg.Client {
				err := client.Conn.WriteJSON(msg.Message)

				if err != nil {
					log.Println("Error writing to WebSocket:", err)

					// Close the WebSocket connection on error
					client.Conn.Close()

					// Remove the disconnected participant from the room
					AllRooms.RemoveParticipant(msg.RoomID, client.Conn)
				}
			}
		}
	}
}

// JoinRoomRequestHandler will join the client in a particular room
func JoinRoomRequestHandler(w http.ResponseWriter, r *http.Request) {
	roomID, ok := r.URL.Query()["roomID"]

	if !ok {
		log.Println("roomID missing in URL Parameters")
		return
	}

	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Fatal("Web Socket Upgrade Error", err)
	}

	AllRooms.InsertIntoRoom(roomID[0], false, ws)

	go broadcaster()

	defer func() {
		// Close the WebSocket connection

		// Remove the disconnected participant from the room
		AllRooms.RemoveParticipant(roomID[0], ws)
	}()

	for {
		var msg broadcastMsg

		err := ws.ReadJSON(&msg.Message)
		if err != nil {
			log.Println("Read Error:", err)
			break

		}

		msg.Client = ws
		msg.RoomID = roomID[0]

		log.Println(msg.Message)

		broadcast <- msg
	}

	err = ws.Close()
	if err != nil {
		log.Println("Error closing WebSocket:", err)
	}
}
