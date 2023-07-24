package database

import (
	"net/http"

	"github.com/gorilla/websocket"
)

type UserProfile struct {
	Username     string
	Image        string
	Biography    string
	Nickname     string
	CreatedAt    int
	CreatedAtStr string
	Gender       string
	Age          int
	Hobby        string
	Birthday     string
}

type Post struct {
	Id        int       `json:"id"`
	UserId    int       `json:"userid"`
	Username  string    `json:"username"`
	Title     string    `json:"title"`
	Body      string    `json:"body"`
	Likes     int       `json:"likes"`
	Dislikes  int       `json:"dislikes"`
	Comments  []Comment `json:"comments"`
	PostTime  string    `json:"postime"`
	Category  string    `json:"category"`
	ImagePath string    `json:"image"`
}

type Comment struct {
	Id       int    `json:"id"`
	UserId   int    `json:"userid"`
	Username string `json:"username"`
	PostId   int    `json:"postid"`
	Body     string `json:"body"`
	Likes    int    `json:"likes"`
	Dislikes int    `json:"dislikes"`
}

type Message struct {
	Id        int    `json:"id"`
	UserId    int    `json:"userid"`
	Username  string `json:"username"`
	CreatedAt string `json:"created_at"`
	Body      string `json:"body"`
}

type WsMessage struct {
	Type string      `json:"type"`
	Data interface{} `json:"data"`
}

type ConnectedUser struct {
	Username string `json:"username"`
	Id       int    `json:"id"`
	Status   string `json:"status"`
}

type Client struct {
	Connection *websocket.Conn
	Username   string `json:"username"`
	IsTyping   bool   `json:"is_typing"`
}

type PrivateMessage struct {
	SenderID       int    `json:"senderID"`
	RecipientID    int    `json:"recipientID"`
	SenderUsername string `json:"senderUsername"`
	RecipientName  string `json:"recipientName"`
	Body           string `json:"body"`
	SentTime       string `json:"sentTime"`
}

var (
	Clients  = make(map[*websocket.Conn]bool)
	Cli      = make(map[*websocket.Conn]Client)
	Upgrader = websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}
)
