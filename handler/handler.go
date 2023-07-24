package handler

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"real-time-forum/database"
	"real-time-forum/utils"

	"github.com/gorilla/websocket"
	"github.com/mattn/go-sqlite3"
)

const (
	Reset = "\033[0m"
	Red   = "\033[31m"
)

func ProfileHandler(w http.ResponseWriter, r *http.Request) {
	userId := utils.AuthentifiedUser(r)
	if userId == -1 {
		http.Redirect(w, r, "/", http.StatusSeeOther)
	}
	p := utils.GetProfile(userId)
	json.NewEncoder(w).Encode(p)
}

func NickNameHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		newNick := r.FormValue("nick")
		userId := utils.AuthentifiedUser(r)
		if userId == -1 {
			http.Redirect(w, r, "/", http.StatusSeeOther)
		}
		utils.UpdateNick(userId, newNick)
		http.Redirect(w, r, "/", http.StatusSeeOther)
	}
}

func BiographyHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		newBio := r.FormValue("bio")
		userId := utils.AuthentifiedUser(r)
		if userId == -1 {
			http.Redirect(w, r, "/", http.StatusSeeOther)
		}
		utils.UpdateBio(userId, newBio)
		http.Redirect(w, r, "/", http.StatusSeeOther)
	}
}

var fileIndex = 0

func HandleUploadImage(w http.ResponseWriter, r *http.Request) {
	userId := utils.AuthentifiedUser(r)
	filepath := ""
	file, _, err := r.FormFile("file")
	if err == nil {
		filepath = fmt.Sprintf("public/profilUploads/%d.png", fileIndex)
		dst, err := os.Create(filepath)
		if err != nil {
			log.Fatalln(err)
		}
		io.Copy(dst, file)
		fileIndex++
		utils.InsertFilePath(filepath, userId)
		http.Redirect(w, r, "/", http.StatusSeeOther)
	}
}

func BirthdayHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		newBirthday := r.FormValue("birthday")
		userId := utils.AuthentifiedUser(r)
		if userId == -1 {
			http.Redirect(w, r, "/", http.StatusSeeOther)
		}
		utils.UpdateBirthday(userId, newBirthday)
		utils.UpdateAge(newBirthday, userId)
		http.Redirect(w, r, "/", http.StatusSeeOther)
	}
}

func HobbyHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		userId := utils.AuthentifiedUser(r)
		if userId == -1 {
			http.Redirect(w, r, "/", http.StatusSeeOther)
		}
		hobby := r.FormValue("hobby")
		utils.UpdateHobby(userId, hobby)
		http.Redirect(w, r, "/", http.StatusSeeOther)
	}
}

func RegisterHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		username := r.FormValue("username")
		mail := r.FormValue("email")
		password := r.FormValue("password")
		male := r.FormValue("radioMale")
		female := r.FormValue("radioFemale")
		otherGender := r.FormValue("radioOther")
		gender := ""
		if male != "" {
			gender = "male"
		}
		if female != "" {
			gender = "female"
		}
		if otherGender != "" {
			gender = "other"
		}
		err := utils.Insert2User(username, mail, password, gender)
		if sqliteErr, ok := err.(sqlite3.Error); ok {
			if sqliteErr.ExtendedCode == sqlite3.ErrConstraintUnique {
				http.SetCookie(w, &http.Cookie{
					Name:  "error",
					Value: "Email already exists. Please try again, or Sign in.",
					Path:  "/",
				})
			} else {
				panic(err)
			}
		}
		http.Redirect(w, r, "/", http.StatusSeeOther)

	}
}

func LoginHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		var mail string
		var username string
		password := r.FormValue("password")
		if strings.Contains(r.FormValue("email"), "@") {
			mail = r.FormValue("email")
		} else {
			username = r.FormValue("email")
		}
		fmt.Println(mail, username)
		if utils.CompareCoord(username, mail, password) == true {
			userID := utils.GetUserID(username, mail, password)
			uuid := utils.SessionUuid(userID)
			fmt.Println("userID:", userID, "UUID:", uuid)
			http.SetCookie(w, &http.Cookie{
				Name:    "session",
				Value:   uuid,
				Expires: time.Now().Add(time.Hour),
				Path:    "/",
			})
			http.SetCookie(w, &http.Cookie{
				Name:  "error",
				Value: " ",
				Path:  "/",
			})
			// set cookie username
			http.SetCookie(w, &http.Cookie{
				Name:  "username",
				Value: username,
				Path:  "/",
			})
			http.Redirect(w, r, "/", http.StatusSeeOther)
		} else {
			http.SetCookie(w, &http.Cookie{
				Name:  "error",
				Value: "Incorrect username, email or password. Please try again.",
				Path:  "/",
			})
			http.Redirect(w, r, "/", http.StatusSeeOther)
		}
	}
}

func WsHandler(w http.ResponseWriter, r *http.Request) {
	jesus, err := database.Upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}
	go func() {
		defer jesus.Close()
		database.Clients[jesus] = true

		fmt.Println("new connection")
		broadcastPosts()
		BroadcastMessage(true)
		BroadCastConnected()
		for {
			_, msg, err := jesus.ReadMessage()
			if err != nil {
				delete(database.Clients, jesus)
				jesus.Close()
				log.Printf("connection ReadMessage error: %s\n", err.Error())
				return
			}

			log.Println(string(msg), "here is the message")
			var message database.WsMessage
			err = json.Unmarshal(msg, &message)
			if err != nil {
				log.Println(err, "line 57")
				return
			}
			switch message.Type {
			case "post":
				var post database.Post
				data, ok := message.Data.(map[string]interface{})
				post.UserId = utils.AuthentifiedUser(r)
				post.Title = data["title"].(string)
				post.Body = data["body"].(string)
				post.Category = data["category"].(string)
				post.Username = "By " + utils.GetUsernameFromUserId(utils.AuthentifiedUser(r))
				post.PostTime = "Posted the: " + time.Now().Format("2006-01-02 at 15:04:05")
				if !ok {
					log.Println("error")
					return
				}
				if err != nil {
					fmt.Println("err is not nil")
					log.Println(err)
					return
				}
				err = utils.Insert2Post(post)
				if err != nil {
					log.Println(err)
					return
				}
				broadcastPosts()
			case "comment":
				var comment database.Comment
				data, ok := message.Data.(map[string]interface{})
				comment.UserId = utils.AuthentifiedUser(r)
				comment.PostId = int(data["postid"].(float64))
				comment.Body = data["body"].(string)
				comment.Username = "Comment from " + utils.GetUsernameFromUserId(utils.AuthentifiedUser(r))
				if err != nil {
					log.Println(err)
					return
				}
				if !ok {
					log.Println("error")
					return
				}
				if err != nil {
					fmt.Println("err is not nil")
					log.Println(err)
					return
				}
				err = utils.Insert2Comment(comment)
				if err != nil {
					log.Println(err)
					return
				}
				broadcastPosts()
			case "message":
				var chat database.Message
				data, ok := message.Data.(map[string]interface{})
				chat.UserId = utils.AuthentifiedUser(r)
				chat.Body = data["body"].(string)
				chat.Username = "Message from " + utils.GetUsernameFromUserId(utils.AuthentifiedUser(r))
				chat.CreatedAt = "Sent the: " + time.Now().Format("2006-01-02 at 15:04:05")
				if err != nil {
					log.Println(err)
					return
				}
				if !ok {
					log.Println("error")
					return
				}
				if err != nil {
					fmt.Println("err is not nil")
					log.Println(err)
					return
				}
				err = utils.Insert2Message(chat)
				if err != nil {
					log.Println(err)
					return
				}
				BroadcastMessage(false)
			case "typing":
				BroadcastTyping(jesus, "is typing...", r)
			case "login":
				fmt.Println("Login Successfull")
				user_id := utils.AuthentifiedUser(r)
				username := utils.GetUsernameFromUserId(user_id)
				utils.UpdateConnected(username, user_id)
				BroadCastConnected()
			case "logout":
				fmt.Println("Logout Successfull")
				user_id := utils.AuthentifiedUser(r)
				username := utils.GetUsernameFromUserId(user_id)
				utils.UpdateDisconnected(username, user_id)
				BroadCastConnected()
			}
		}
	}()
}

func broadcastPosts() {
	for client := range database.Clients {

		posts, err := utils.GetPosts(utils.OpenDatabase())
		if err != nil {
			log.Println(err)
			return
		}

		jsonPosts := database.WsMessage{
			Type: "posts",
			Data: posts,
		}
		jsonPostsBytes, err := json.Marshal(jsonPosts)
		if err != nil {
			log.Println(err)
			return
		}
		fmt.Println(Red, "jsonPostsBytes", string(jsonPostsBytes), Reset)
		err = client.WriteMessage(websocket.TextMessage, jsonPostsBytes)
		if err != nil {
			log.Println(err)
			return
		}

	}
}

func BroadcastTyping(sender *websocket.Conn, message string, r *http.Request) {
	for conn := range database.Clients {
		if conn != sender {
			jsonTyping := database.WsMessage{
				Type: "typing",
				Data: map[string]any{"message": message, "username": utils.GetUsernameFromUserId(utils.AuthentifiedUser(r))},
			}
			jsonTypingBytes, err := json.Marshal(jsonTyping)
			if err != nil {
				log.Println(err)
				return
			}
			err = conn.WriteMessage(websocket.TextMessage, jsonTypingBytes)
			if err != nil {
				log.Println(err)
				return
			}
		}
	}
}

func BroadcastMessage(initial bool) {
	for client := range database.Clients {

		messages, err := utils.GetMessages(utils.OpenDatabase())
		if err != nil {
			log.Println(err)
			return
		}
		fmt.Println(Red, "messages", messages, Reset)
		jsonMessages := database.WsMessage{
			Type: "messages",
			Data: map[string]any{"messages": messages, "initial": initial},
		}
		jsonMessagesBytes, err := json.Marshal(jsonMessages)
		if err != nil {
			log.Println(err)
			return
		}
		fmt.Println(Red, "jsonMessagesBytes", string(jsonMessagesBytes), Reset)
		err = client.WriteMessage(websocket.TextMessage, jsonMessagesBytes)
		if err != nil {
			log.Println(err)
			return
		}

	}
}

func BroadCastConnected() {
	for client := range database.Clients {

		connected, err := utils.GetConnected(utils.OpenDatabase())
		if err != nil {
			log.Println(err)
			return
		}
		fmt.Println("connected", connected)
		jsonConnected := database.WsMessage{
			Type: "status",
			Data: map[string]any{"connected": connected},
		}
		jsonConnectedBytes, err := json.Marshal(jsonConnected)
		if err != nil {
			log.Println(err)
			return
		}
		fmt.Println(Red, "jsonConnectedBytes", string(jsonConnectedBytes), Reset)
		err = client.WriteMessage(websocket.TextMessage, jsonConnectedBytes)
		if err != nil {
			log.Println(err)
			return
		}

	}
}
