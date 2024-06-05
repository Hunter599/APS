package main

import (
	"database/sql"
	"fmt"
	"log"
	"syscall"

	"golang.org/x/term"
)


type Storage interface {
    CreateUser(*User) error
    GetUserById(int) (*User, error)
    GetProfilePicPathById(int) (string, error)
    GetUserPostsByUserId(int) ([]*GetPostsRequest, error)
    /*
    CreatePost(*Post) error
    GetPostById(int) (*Post, error)
    CreateComment(*Comment) error
    GetPostComments(int) ([]*Post, error)
    */
}

type MySQLStorage struct {
    db *sql.DB
}

func NewMySQLStorage() (*MySQLStorage, error) {
    user, password := authenticate()
    connection := user + ":" + password + "@/" + DB_NAME

    db, err := sql.Open("mysql", connection)
    if err != nil {
        return nil, err
    }

    return &MySQLStorage{db: db}, nil
}

func (s *MySQLStorage) CreateUser(u *User) error {
    q := `INSERT INTO app_user(username, email, password_hash) VALUES (?, ?, ?);`

    result, err := s.db.Exec(q, u.Username, u.Email, u.Password)
    if err != nil {
        log.Println(err, "here")
        return err
    }
    
    id, _ := result.LastInsertId()
    u.Id = int(id)
    return nil
}

func (s *MySQLStorage) GetUserById(id int) (*User, error) {
    u := new(User)
    q := `
    SELECT
    username, email, first_name, last_name, sex, profile_picture_path
    FROM app_user
    WHERE id = ?;`

    err := s.db.QueryRow(q, id).Scan(
        &u.Username,
        &u.Email,
        &u.FirstName,
        &u.LastName,
        &u.Sex,
        &u.ProfilePicture)

    if err != nil {
        return nil, err
    }

    u.Id = id
    return u, nil
}

func (s *MySQLStorage) GetProfilePicPathById(id int) (string, error) {
    var path string
    q := "SELECT profile_picture_path FROM app_user WHERE id = ?"
    err := s.db.QueryRow(q, id).Scan(&path)
    if err != nil {
        return "", err
    }
    
    return path, nil
}

func (s * MySQLStorage) GetUserPostsByUserId(id int) ([]*GetPostsRequest, error) {
    posts := make([]*GetPostsRequest, 0)
    q := `
    SELECT title, body, app_user.username, community.community_name
    FROM post JOIN community ON post.community_id = community.id
    JOIN app_user ON post.user_id = app_user.id
    WHERE user_id = ?;
    `
    rows, err := s.db.Query(q, id)
    if err != nil {
        return nil, err
    }

    for rows.Next() {
        p := new(GetPostsRequest)
        err := rows.Scan(
            &p.Post.Title,
            &p.Post.Body,
            &p.User.Username,
            &p.Community.Name,
        )
        if err != nil {
            return nil, err
        }

        p.User.Id = id
        posts = append(posts, p)
    }

    return posts, nil
}

func authenticate() (user string, password string) {

    fmt.Print("database user: ")
    fmt.Scan(&user)

    fmt.Print("database password: ")
    bytePassword, err := term.ReadPassword(int(syscall.Stdin))
    if err != nil {
        log.Fatal("Could not read database password.", err)
        return
    }

    password = string(bytePassword)
    return 
}
