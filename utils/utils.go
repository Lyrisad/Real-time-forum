package utils

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	"real-time-forum/database"

	"github.com/gofrs/uuid"
)

var (
	posts database.Post
	db    *sql.DB
)

func Insert2Post(post database.Post) error {
	db := OpenDatabase()
	defer db.Close()

	stmt, err := db.Prepare("INSERT INTO posts (user_id, title, body, category, username, postime, likes, dislikes) VALUES (?, ?, ?, ?, ?, ?, ?, ?)")
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(post.UserId, post.Title, post.Body, post.Category, post.Username, post.PostTime, post.Likes, post.Dislikes)
	if err != nil {
		return err
	}

	return nil
}

func Insert2Comment(comment database.Comment) error {
	db := OpenDatabase()
	defer db.Close()

	stmt, err := db.Prepare("INSERT INTO comments (body, user_id, username, post_id, likes, dislikes) VALUES (?, ?, ?, ?, ?, ?)")
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(comment.Body, comment.UserId, comment.Username, comment.PostId, comment.Likes, comment.Dislikes)
	if err != nil {
		return err
	}

	return nil
}

func Insert2Message(message database.Message) error {
	db := OpenDatabase()
	defer db.Close()

	stmt, err := db.Prepare("INSERT INTO messages (body, user_id, username, postime) VALUES (?, ?, ?, ?)")
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(message.Body, message.UserId, message.Username, message.CreatedAt)
	if err != nil {
		return err
	}

	return nil
}

func OpenDatabase() *sql.DB {
	// open database
	var err error
	db, err = sql.Open("sqlite3", "./data.db")
	if err != nil {
		log.Fatal(err)
	}
	return db
}

func CreateTables() {
	_, err := db.Exec("CREATE TABLE IF NOT EXISTS users (id INTEGER PRIMARY KEY, username TEXT NOT NULL, mail TEXT NOT NULL UNIQUE, password TEXT NOT NULL, image TEXT, biography TEXT, created_at INT, nickname TEXT, gender TEXT, age INT, hobby TEXT, birthday TEXT)")
	if err != nil {
		fmt.Println("Error while creating table users")
		log.Fatal(err)
	}
	_, err = db.Exec("CREATE TABLE IF NOT EXISTS sessions (user_id INTEGER, session_uuid TEXT)")
	if err != nil {
		fmt.Println("Error while creating table sessions")
		log.Fatal(err)
	}
	_, err = db.Exec("CREATE TABLE IF NOT EXISTS userConnected (user_id INTEGER, username TEXT, status TEXT)")
	if err != nil {
		fmt.Println("Error while creating table sessions")
		log.Fatal(err)
	}
	// create post table
	_, err = db.Exec("CREATE TABLE IF NOT EXISTS posts (id INTEGER PRIMARY KEY AUTOINCREMENT, user_id INTEGER, username TEXT, title TEXT, body TEXT, likes INTEGER, dislikes INTEGER NOT NULL, postime TEXT, category TEXT, image TEXT)")
	if err != nil {
		fmt.Println("Error while creating table posts")
		log.Fatal(err)
	}
	_, err = db.Exec("CREATE TABLE IF NOT EXISTS comments (id INTEGER PRIMARY KEY AUTOINCREMENT, user_id INTEGER, post_id INTEGER, username TEXT, body TEXT, likes INTEGER, dislikes INTEGER NOT NULL, postime TEXT)")
	if err != nil {
		fmt.Println("Error while creating table posts")
		log.Fatal(err)
	}
	_, err = db.Exec("CREATE TABLE IF NOT EXISTS messages (id INTEGER PRIMARY KEY AUTOINCREMENT, user_id INTEGER,  username TEXT, body TEXT, postime TEXT)")
	if err != nil {
		fmt.Println("Error while creating table posts")
		log.Fatal(err)
	}
	// create private message table
	_, err = db.Exec("CREATE TABLE IF NOT EXISTS private_messages (sender_id INTEGER, receiver_id INTEGER, senderUsername TEXT, receiverName TEXT, body TEXT, postime TEXT)")
	if err != nil {
		fmt.Println("Error while creating table posts")
		log.Fatal(err)
	}
}

func Insert2User(username string, mail string, password string, gender string) error {
	CreatedAt := time.Now().Unix()
	_, err := db.Exec("INSERT INTO users (username, mail, password, image, biography, created_at, nickname, gender, age, hobby, birthday) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?,?)",
		username, mail, password, "/images/user/default.png", "Enter your biography", CreatedAt, "Enter your nickname", gender, 0, "Enter your hobby", "Enter your birthday")
	if err != nil {
		return err
	}
	return nil
}

func CompareCoord(username string, mail string, password string) bool {
	var dbId int
	var dbMail string
	var dbPassword string
	var dbUsername string

	if mail != "" {
		db.QueryRow("SELECT id, mail, password, username FROM users WHERE mail = ?", mail).Scan(&dbId, &dbMail, &dbPassword, &dbUsername)
	} else {
		db.QueryRow("SELECT id, mail, password, username FROM users WHERE username = ?", username).Scan(&dbId, &dbMail, &dbPassword, &dbUsername)
	}
	fmt.Println(dbMail, dbPassword, dbUsername)
	if dbMail == mail && dbPassword == password || dbUsername == username && dbPassword == password {
		uuid, _ := uuid.NewV4()
		if _, err := db.Exec("INSERT INTO sessions (user_id, session_uuid) VALUES (?,?)", dbId, uuid.String()); err != nil {
			panic(err)
		}
		return true
	}
	return false
}

func UpdateConnected(username string, user_id int) {
	// check if user_id already exists
	var count int
	err := db.QueryRow("SELECT COUNT(*) FROM userConnected WHERE user_id = ?", user_id).Scan(&count)
	if err != nil {
		fmt.Println("Error while checking if user_id already exists")
		log.Fatal(err)
	}

	if count > 0 {
		// user_id already exists, update the row
		_, err = db.Exec("UPDATE userConnected SET username = ?, status = ? WHERE user_id = ?", username, "Online", user_id)
		if err != nil {
			fmt.Println("Error while updating session")
			log.Fatal(err)
		}
	} else {
		// user_id does not exist, insert a new row
		_, err = db.Exec("INSERT INTO userConnected (user_id, username, status) VALUES (?, ?, ?)", user_id, username, "Online")
		if err != nil {
			fmt.Println("Error while inserting in session")
			log.Fatal(err)
		}
	}
}

func UpdateDisconnected(username string, user_id int) {
	// update status to "offline" in db where user_id and username
	_, err := db.Exec("UPDATE userConnected SET status = ? WHERE user_id = ? AND username = ?", "Offline", user_id, username)
	if err != nil {
		fmt.Println("Error while updating session status to offline")
		log.Fatal(err)
	}
}

func GetUserID(username string, mail string, password string) int {
	var dbID int
	if mail != "" {
		db.QueryRow("SELECT id FROM users WHERE mail = ? AND password = ?", mail, password).Scan(&dbID)
	} else {
		db.QueryRow("SELECT id FROM users WHERE username = ? AND password = ?", username, password).Scan(&dbID)
	}
	return dbID
}

func SessionUuid(userId int) string {
	var dbSessionUUID string
	db.QueryRow("SELECT session_uuid FROM sessions WHERE user_id = ?", userId).Scan(&dbSessionUUID)
	return dbSessionUUID
}

func AuthentifiedUser(r *http.Request) int {
	cookie, err := r.Cookie("session")
	if err != nil {
		return -1
	}
	var user_id int
	err = db.QueryRow("SELECT user_id FROM sessions WHERE session_uuid = ?", cookie.Value).Scan(&user_id)
	if err != nil {
		return -1
	}
	fmt.Println(cookie.Value, user_id, "cookie value and user id")
	return user_id
}

func GetProfile(user_id int) database.UserProfile {
	var up database.UserProfile
	err := db.QueryRow("SELECT username, image, biography, nickname, created_at, gender, age, hobby, birthday FROM users WHERE id = ?",
		user_id).Scan(&up.Username, &up.Image, &up.Biography, &up.Nickname, &up.CreatedAt, &up.Gender, &up.Age, &up.Hobby, &up.Birthday)
	if err != nil {
		log.Println(err)
		return database.UserProfile{}
	}
	up.CreatedAtStr = time.Unix(int64(up.CreatedAt), 0).Format("2006-01-02 at 15:04:05")
	return up
}

func UpdateNick(user_id int, newNickname string) {
	_, err := db.Exec("UPDATE users SET nickname = ? WHERE id = ?", newNickname, user_id)
	if err != nil {
		log.Println(err)
	}
}

func UpdateHobby(user_id int, newHobby string) {
	_, err := db.Exec("UPDATE users SET hobby = ? WHERE id = ?", newHobby, user_id)
	if err != nil {
		log.Println(err)
	}
}

func UpdateBio(user_id int, newBio string) {
	_, err := db.Exec("UPDATE users SET biography = ? WHERE id = ?", newBio, user_id)
	if err != nil {
		log.Println(err)
	}
}

func InsertFilePath(filePath string, user_id int) {
	_, err := db.Exec("UPDATE users SET image = ? WHERE id = ?", strings.TrimPrefix(filePath, "public"), user_id)
	if err != nil {
		log.Println(err)
	}
}

func UpdateBirthday(user_id int, newBirthday string) {
	_, err := db.Exec("UPDATE users SET birthday = ? WHERE id = ?", newBirthday, user_id)
	if err != nil {
		log.Println(err)
	}
}

func UpdateAge(newBirthday string, user_id int) {
	// birthday format : 2006-01-02
	year, _ := time.Parse("2006-01-02", newBirthday)
	age := time.Now().Year() - year.Year()
	_, err := db.Exec("UPDATE users SET age = ? WHERE id = ?", age, user_id)
	if err != nil {
		log.Println(err)
	}
}

func GetUsernameFromUserId(uid int) string {
	var username string
	row := db.QueryRow("SELECT username FROM users WHERE id = ?", uid)
	err := row.Scan(&username)
	if err != nil {
		panic(err)
	}
	return username
}

func GetPosts(db *sql.DB) ([]database.Post, error) {
	rows, err := db.Query("SELECT id, title, body, username, postime, category FROM posts")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var posts []database.Post
	for rows.Next() {
		var post database.Post
		err := rows.Scan(&post.Id, &post.Title, &post.Body, &post.Username, &post.PostTime, &post.Category)
		if err != nil {
			return nil, err
		}
		post.Comments, err = GetComments(db, post.Id)
		if err != nil {
			return nil, err
		}
		posts = append(posts, post)
		fmt.Println(post, "here are the posts")
	}

	err = rows.Err()
	if err != nil {
		return nil, err
	}

	return posts, nil
}

func GetComments(db *sql.DB, post_id int) ([]database.Comment, error) {
	rows, err := db.Query("SELECT id, post_id, username, body, likes, dislikes FROM comments WHERE post_id = ?", post_id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var comments []database.Comment
	for rows.Next() {
		var comment database.Comment
		err := rows.Scan(&comment.Id, &comment.PostId, &comment.Username, &comment.Body, &comment.Likes, &comment.Dislikes)
		if err != nil {
			return nil, err
		}
		comments = append(comments, comment)
	}

	err = rows.Err()
	if err != nil {
		return nil, err
	}

	return comments, nil
}

func APINewCommentHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		uid := AuthentifiedUser(r)
		if uid == -1 {
			w.WriteHeader(http.StatusForbidden)
			return
		}

		postId := r.FormValue("id")
		body := r.FormValue("commentBox")

		postIdInt, err := strconv.Atoi(postId)
		if err != nil {
			log.Println(err)
			http.Error(w, "Invalid post id.", http.StatusBadRequest)
			return
		}

		username := GetUsernameFromUserId(uid)

		_, err = db.Exec("INSERT INTO comments(post_id, user_id, username, body, likes, dislikes) VALUES(?,?,?,?,?,?)", postIdInt, uid, username, body, 0, 0)
		if err != nil {
			log.Println(err)
			http.Error(w, "Failed to insert comment into the database.", http.StatusInternalServerError)
			return
		}
		return
	}
}

func GetUserProfile(username string) database.UserProfile {
	var profile database.UserProfile

	query := "SELECT username, image, biography, nickname, created_at, created_at_str, gender, age, hobby, birthday FROM user_profiles WHERE username = $1"
	err := db.QueryRow(query, username).Scan(
		&profile.Username,
		&profile.Image,
		&profile.Biography,
		&profile.Nickname,
		&profile.CreatedAt,
		&profile.CreatedAtStr,
		&profile.Gender,
		&profile.Age,
		&profile.Hobby,
		&profile.Birthday,
	)
	if err != nil {
		// Handle the error (e.g., return a default or empty profile)
		log.Println("Error retrieving user profile:", err)
		return database.UserProfile{}
	}

	return profile
}

func GetMessages(db *sql.DB) ([]database.Message, error) {
	var messages []database.Message
	rows, err := db.Query("SELECT id, user_id, body, username, postime FROM messages")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var message database.Message
		err := rows.Scan(&message.Id, &message.UserId, &message.Body, &message.Username, &message.CreatedAt)
		if err != nil {
			return nil, err
		}
		messages = append(messages, message)
	}
	return messages, nil
}

func GetConnected(db *sql.DB) ([]database.ConnectedUser, error) {
	var connectedUsers []database.ConnectedUser
	rows, err := db.Query("SELECT user_id, username, status FROM UserConnected")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var connectedUser database.ConnectedUser
		err := rows.Scan(&connectedUser.Id, &connectedUser.Username, &connectedUser.Status)
		if err != nil {
			return nil, err
		}
		connectedUsers = append(connectedUsers, connectedUser)
	}
	return connectedUsers, nil
}
