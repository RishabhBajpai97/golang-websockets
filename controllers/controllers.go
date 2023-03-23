package controllers

import (
	"fmt"
	"log"
	"net/http"
	"sort"

	"github.com/CloudyKit/jet/v6"
	"github.com/gorilla/websocket"
)

var wsChan = make(chan WsPayload)
var clients = make(map[WebSocketConnection]string)

var views = jet.NewSet(
	jet.NewOSFileSystemLoader("./html"),
	jet.InDevelopmentMode(),
)

var upgradeConnection = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func Home(w http.ResponseWriter, r *http.Request) {
	err := renderPage(w, "home.jet", nil)
	if err != nil {
		log.Println(err)
	}
}

type WebSocketConnection struct {
	*websocket.Conn
}
type WsResponse struct {
	Action      string `json:"action"`
	Message     string `json:"message"`
	MessageType string `json:"message_type"`
	ConnectedUsers []string `json:"connected_users"`
}

type WsPayload struct {
	Action      string              `json:"action"`
	Username    string              `json:"username"`
	Message     string              `json:"message"`
	MessageType string              `json:"message_type"`
	Conn        WebSocketConnection `json:"-"`
}

func WsEndpoint(w http.ResponseWriter, r *http.Request) {
	ws, err := upgradeConnection.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
	}
	log.Println("Connected to the client")
	var response WsResponse
	response.Message = "Hello"

	conn := WebSocketConnection{Conn: ws}
	clients[conn] = ""
	err = ws.WriteJSON(response)
	if err != nil {
		log.Println(err)
	}

	go ListenforWs(&conn)

}

func ListentoWs(){
	var response WsResponse
	for{
		e:= <- wsChan
	switch e.Action{
	case "username":
		clients[e.Conn] = e.Username
		userList := getUserList();
		response.Action="users_list"
		response.ConnectedUsers=userList
		broadcastAll(response)
	case "left":
		response.Action = "users_list"
		delete(clients, e.Conn)
		users := getUserList()
		response.ConnectedUsers = users
		broadcastAll(response)
	case "broadcast":
		response.Action = "broadcast"
		response.Message = fmt.Sprintf("<strong>%s</strong> : %s", e.Username, e.Message)
		broadcastAll(response)

	}
	}

}

func getUserList()[]string{
	var userList []string

	for _,x := range clients{
		if x!=""{
			userList = append(userList, x)
		}
	}
	sort.Strings(userList)
	return userList
}


func broadcastAll(response WsResponse){
		for client:= range clients{
			err:=client.WriteJSON(response);
			if err!=nil{
				log.Println("Websocket error")
				_=client.Close()
				delete(clients,client)
			}
		}
}

func ListenforWs(conn *WebSocketConnection) {
	defer func() {
		if r := recover(); r != nil {
			log.Println("error")
		}
	}()
	var payload WsPayload
	for {
		err := conn.ReadJSON(&payload);

		if err!=nil{
			log.Println("Error in payload conversion")
		}else{
			fmt.Println(payload)
			payload.Conn = *conn

			wsChan <- payload
		}
	}

}

func renderPage(w http.ResponseWriter, templ string, data jet.VarMap) error {
	view, err := views.GetTemplate(templ)
	if err != nil {
		log.Println(err)
		return err
	}
	err = view.Execute(w, data, nil)
	if err != nil {
		log.Println(err)
		return err
	}
	return nil
}
