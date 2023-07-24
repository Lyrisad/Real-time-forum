package main

import (
	"fmt"
	"net/http"

	"real-time-forum/handler"
	"real-time-forum/utils"
)

const (
	Reset = "\033[0m"
	Red   = "\033[31m"
)

func main() {
	utils.OpenDatabase()
	utils.CreateTables()
	defer utils.OpenDatabase().Close()

	fs := http.FileServer(http.Dir("./public"))
	http.Handle("/", fs)
	http.HandleFunc("/ws", handler.WsHandler)

	http.HandleFunc("/register", handler.RegisterHandler)
	http.HandleFunc("/login", handler.LoginHandler)
	http.HandleFunc("/profile", handler.ProfileHandler)
	http.HandleFunc("/nick", handler.NickNameHandler)
	http.HandleFunc("/bio", handler.BiographyHandler)
	http.HandleFunc("/upload", handler.HandleUploadImage)
	http.HandleFunc("/birthday", handler.BirthdayHandler)
	http.HandleFunc("/hobby", handler.HobbyHandler)

	fmt.Println("http://localhost:8080")
	http.ListenAndServe(":8080", nil)
}
